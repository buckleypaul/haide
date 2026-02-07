package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/buckleypaul/haide/internal/config"
	"github.com/buckleypaul/haide/internal/exclude"
	"github.com/buckleypaul/haide/internal/git"
	"github.com/spf13/cobra"
)

var initDryRun bool

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize haide in the current repository",
	Long: `Initialize haide in the current repository by creating/updating .git/info/exclude
with haide-managed entries. Applies global exclusions and project-specific overrides.`,
	RunE: runInit,
}

func init() {
	initCmd.Flags().BoolVar(&initDryRun, "dry-run", false, "Preview changes without applying them")
}

func runInit(cmd *cobra.Command, args []string) error {
	// Check if in git repository
	repoRoot, err := git.GetRepoRoot()
	if err != nil {
		return fmt.Errorf("not in a git repository")
	}

	repoName, err := git.GetRepoBasename()
	if err != nil {
		return err
	}

	// Load config (creates default if doesn't exist)
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Load local config
	localCfg, err := config.LoadLocal(repoRoot)
	if err != nil {
		return fmt.Errorf("failed to load local config: %w", err)
	}

	// Migrate legacy [project:name] patterns to local config
	globalConfigModified := false
	if legacyPatterns := cfg.MigrateProjectToLocal(repoName); len(legacyPatterns) > 0 {
		for _, p := range legacyPatterns {
			localCfg.AddPattern(p)
		}
		globalConfigModified = true
		fmt.Printf("Migrated %d patterns from global config to local .haide file\n", len(legacyPatterns))
	}

	// Get effective patterns
	effectivePatterns := cfg.GetEffectivePatternsWithLocal(repoName, localCfg.Patterns)

	// Check which excluded files are already tracked
	trackedFiles, err := git.GetTrackedFiles(effectivePatterns)
	if err != nil {
		return fmt.Errorf("failed to check tracked files: %w", err)
	}

	// Add overrides for tracked files
	localConfigModified := false
	for _, file := range trackedFiles {
		for _, pattern := range effectivePatterns {
			matches, _ := matchFile(file, pattern)
			if matches {
				fmt.Printf("File '%s' is tracked but matches exclusion pattern '%s'\n", file, pattern)
				localCfg.AddPattern("+" + pattern)
				localConfigModified = true

				// Ask if user wants to untrack
				if !initDryRun {
					fmt.Printf("Do you want to untrack this file? (y/N): ")
					reader := bufio.NewReader(os.Stdin)
					response, _ := reader.ReadString('\n')
					response = strings.TrimSpace(strings.ToLower(response))
					if response == "y" || response == "yes" {
						if err := git.UntrackFile(file); err != nil {
							fmt.Printf("Warning: failed to untrack file: %v\n", err)
						} else {
							fmt.Printf("Untracked '%s'\n", file)
						}
					}
				}
				break
			}
		}
	}

	// Save configs if modified
	if !initDryRun {
		if globalConfigModified {
			if err := cfg.Save(); err != nil {
				return fmt.Errorf("failed to save config: %w", err)
			}
		}
		if localConfigModified || globalConfigModified {
			if err := localCfg.Save(); err != nil {
				return fmt.Errorf("failed to save local config: %w", err)
			}
			fmt.Println("Saved overrides to local .haide file")
		}
	}

	// Recalculate effective patterns after adding overrides
	effectivePatterns = cfg.GetEffectivePatternsWithLocal(repoName, localCfg.Patterns)
	excludePatterns := buildExcludePatterns(effectivePatterns)

	// Load exclude file
	excludeFile, err := exclude.Load(repoRoot)
	if err != nil {
		return fmt.Errorf("failed to load exclude file: %w", err)
	}

	// Show diff
	diff := excludeFile.GetDiff(excludePatterns)
	fmt.Println("Changes to .git/info/exclude:")
	fmt.Println(diff)

	if initDryRun {
		fmt.Println("\n(Dry run - no changes applied)")
		return nil
	}

	// Update exclude file
	excludeFile.SetManagedContent(excludePatterns)
	if err := excludeFile.Save(); err != nil {
		return fmt.Errorf("failed to save exclude file: %w", err)
	}

	fmt.Println("\nSuccessfully initialized haide")
	fmt.Printf("Config: %s\n", cfg.Path())
	fmt.Printf("Local config: %s\n", localCfg.Path())
	fmt.Printf("Exclude file: %s\n", excludeFile.Path())

	return nil
}

func matchFile(file, pattern string) (bool, error) {
	// Simple pattern matching - can be enhanced
	if strings.HasSuffix(pattern, "/") {
		dirPattern := strings.TrimSuffix(pattern, "/")
		return strings.HasPrefix(file, dirPattern+"/"), nil
	}
	return strings.Contains(file, pattern) || strings.HasSuffix(file, pattern), nil
}

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

var (
	updateDryRun bool
	updateYes    bool
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update repository exclusions from config",
	Long: `Update the current repository's .git/info/exclude file to match the config.
Shows a diff and prompts for confirmation before applying changes.`,
	RunE: runUpdate,
}

func init() {
	updateCmd.Flags().BoolVar(&updateDryRun, "dry-run", false, "Preview changes without applying them")
	updateCmd.Flags().BoolVarP(&updateYes, "yes", "y", false, "Auto-confirm without prompting")
}

func runUpdate(cmd *cobra.Command, args []string) error {
	// Check if in git repository
	repoRoot, err := git.GetRepoRoot()
	if err != nil {
		return fmt.Errorf("not in a git repository")
	}

	repoName, err := git.GetRepoBasename()
	if err != nil {
		return err
	}

	// Load config
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

	if diff == "No changes" && !globalConfigModified {
		fmt.Println("Already up to date")
		return nil
	}

	if updateDryRun {
		fmt.Println("\n(Dry run - no changes applied)")
		return nil
	}

	// Prompt for confirmation unless --yes
	if !updateYes && diff != "No changes" {
		fmt.Print("\nApply these changes? (y/N): ")
		reader := bufio.NewReader(os.Stdin)
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(strings.ToLower(response))
		if response != "y" && response != "yes" {
			fmt.Println("Aborted")
			return nil
		}
	}

	// Save configs if migration happened
	if globalConfigModified {
		if err := cfg.Save(); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}
		if err := localCfg.Save(); err != nil {
			return fmt.Errorf("failed to save local config: %w", err)
		}
	}

	// Update exclude file
	excludeFile.SetManagedContent(excludePatterns)
	if err := excludeFile.Save(); err != nil {
		return fmt.Errorf("failed to save exclude file: %w", err)
	}

	fmt.Println("\nSuccessfully updated exclusions")
	return nil
}

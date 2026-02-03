package commands

import (
	"fmt"

	"github.com/buckleypaul/haide/internal/config"
	"github.com/buckleypaul/haide/internal/exclude"
	"github.com/buckleypaul/haide/internal/git"
	"github.com/spf13/cobra"
)

var addDryRun bool

var addCmd = &cobra.Command{
	Use:   "add PATTERN",
	Short: "Add a pattern to project-specific exclusions",
	Long:  `Add a file or directory pattern to the current project's exclusion list.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runAdd,
}

func init() {
	addCmd.Flags().BoolVar(&addDryRun, "dry-run", false, "Preview changes without applying them")
}

func runAdd(cmd *cobra.Command, args []string) error {
	pattern := args[0]

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

	// Add pattern to project
	cfg.AddProjectPattern(repoName, pattern)

	// Get effective patterns
	effectivePatterns := cfg.GetEffectivePatterns(repoName)

	// Load exclude file
	excludeFile, err := exclude.Load(repoRoot)
	if err != nil {
		return fmt.Errorf("failed to load exclude file: %w", err)
	}

	// Show diff
	diff := excludeFile.GetDiff(effectivePatterns)
	fmt.Println("Changes to .git/info/exclude:")
	fmt.Println(diff)

	if addDryRun {
		fmt.Println("\n(Dry run - no changes applied)")
		return nil
	}

	// Save config
	if err := cfg.Save(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	// Update exclude file
	excludeFile.SetManagedContent(effectivePatterns)
	if err := excludeFile.Save(); err != nil {
		return fmt.Errorf("failed to save exclude file: %w", err)
	}

	fmt.Printf("\nAdded pattern '%s' to project '%s'\n", pattern, repoName)
	return nil
}

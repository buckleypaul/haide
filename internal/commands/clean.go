package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/buckleypaul/haide/internal/exclude"
	"github.com/buckleypaul/haide/internal/git"
	"github.com/spf13/cobra"
)

var cleanDryRun bool

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Remove haide management from repository",
	Long: `Remove the haide-managed section from .git/info/exclude.
Preserves non-haide entries in the file.`,
	RunE: runClean,
}

func init() {
	cleanCmd.Flags().BoolVar(&cleanDryRun, "dry-run", false, "Preview changes without applying them")
}

func runClean(cmd *cobra.Command, args []string) error {
	// Check if in git repository
	repoRoot, err := git.GetRepoRoot()
	if err != nil {
		return fmt.Errorf("not in a git repository")
	}

	// Load exclude file
	excludeFile, err := exclude.Load(repoRoot)
	if err != nil {
		return fmt.Errorf("failed to load exclude file: %w", err)
	}

	// Show what will be removed
	managedContent := excludeFile.GetManagedContent()
	if len(managedContent) == 0 {
		fmt.Println("No haide-managed content found")
		return nil
	}

	fmt.Println("Will remove the following patterns:")
	for _, pattern := range managedContent {
		fmt.Printf("  %s\n", pattern)
	}

	if cleanDryRun {
		fmt.Println("\n(Dry run - no changes applied)")
		return nil
	}

	// Prompt for confirmation
	fmt.Print("\nRemove haide management from this repository? (y/N): ")
	reader := bufio.NewReader(os.Stdin)
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(strings.ToLower(response))
	if response != "y" && response != "yes" {
		fmt.Println("Aborted")
		return nil
	}

	// Clear managed content
	excludeFile.ClearManagedContent()
	if err := excludeFile.Save(); err != nil {
		return fmt.Errorf("failed to save exclude file: %w", err)
	}

	fmt.Println("\nSuccessfully removed haide management")
	fmt.Println("Note: Config file is preserved for other repositories")
	return nil
}

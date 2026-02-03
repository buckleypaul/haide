package commands

import (
	"fmt"

	"github.com/buckleypaul/haide/internal/config"
	"github.com/spf13/cobra"
)

var removeGlobalDryRun bool

var removeGlobalCmd = &cobra.Command{
	Use:   "remove-global PATTERN",
	Short: "Remove a pattern from global exclusions",
	Long:  `Remove a file or directory pattern from the global exclusion list.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runRemoveGlobal,
}

func init() {
	removeGlobalCmd.Flags().BoolVar(&removeGlobalDryRun, "dry-run", false, "Preview changes without applying them")
}

func runRemoveGlobal(cmd *cobra.Command, args []string) error {
	pattern := args[0]

	// Load config
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Check if pattern exists
	found := false
	for _, p := range cfg.Global {
		if p == pattern {
			found = true
			break
		}
	}

	if !found {
		fmt.Printf("Pattern '%s' not found in global exclusions\n", pattern)
		return nil
	}

	// Remove pattern from global
	cfg.RemoveGlobalPattern(pattern)

	if removeGlobalDryRun {
		fmt.Printf("Would remove pattern '%s' from global exclusions\n", pattern)
		fmt.Println("\n(Dry run - no changes applied)")
		return nil
	}

	// Save config
	if err := cfg.Save(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Printf("Removed pattern '%s' from global exclusions\n", pattern)
	fmt.Println("Run 'haide update' in each repository to apply the change")
	return nil
}

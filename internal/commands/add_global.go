package commands

import (
	"fmt"

	"github.com/buckleypaul/haide/internal/config"
	"github.com/spf13/cobra"
)

var addGlobalDryRun bool

var addGlobalCmd = &cobra.Command{
	Use:   "add-global PATTERN",
	Short: "Add a pattern to global exclusions",
	Long:  `Add a file or directory pattern to the global exclusion list.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runAddGlobal,
}

func init() {
	addGlobalCmd.Flags().BoolVar(&addGlobalDryRun, "dry-run", false, "Preview changes without applying them")
}

func runAddGlobal(cmd *cobra.Command, args []string) error {
	pattern := args[0]

	// Load config
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Check if pattern already exists
	for _, p := range cfg.Global {
		if p == pattern {
			fmt.Printf("Pattern '%s' already exists in global exclusions\n", pattern)
			return nil
		}
	}

	// Add pattern to global
	cfg.AddGlobalPattern(pattern)

	if addGlobalDryRun {
		fmt.Printf("Would add pattern '%s' to global exclusions\n", pattern)
		fmt.Println("\n(Dry run - no changes applied)")
		return nil
	}

	// Save config
	if err := cfg.Save(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Printf("Added pattern '%s' to global exclusions\n", pattern)
	fmt.Println("Run 'haide update' in each repository to apply the change")
	return nil
}

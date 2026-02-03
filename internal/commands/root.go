package commands

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "haide",
	Short: "AI File Exclusion Manager for Git",
	Long: `haide manages AI-specific file exclusions across git repositories using
a centralized configuration with global and per-project overrides.`,
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(addGlobalCmd)
	rootCmd.AddCommand(removeCmd)
	rootCmd.AddCommand(removeGlobalCmd)
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(infoCmd)
	rootCmd.AddCommand(cleanCmd)
}

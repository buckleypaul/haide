package commands

import (
	"fmt"
	"strings"

	"github.com/buckleypaul/haide/internal/config"
	"github.com/buckleypaul/haide/internal/git"
	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Display current configuration",
	Long:  `Display the current exclusion configuration including global and project-specific patterns.`,
	RunE:  runInfo,
}

func runInfo(cmd *cobra.Command, args []string) error {
	// Load config
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	fmt.Printf("Config file: %s\n\n", cfg.Path())

	// Display global exclusions
	fmt.Println("Global exclusions:")
	if len(cfg.Global) == 0 {
		fmt.Println("  (none)")
	} else {
		for _, pattern := range cfg.Global {
			fmt.Printf("  %s\n", pattern)
		}
	}

	// Check if in a git repository
	repoRoot, err := git.GetRepoRoot()
	if err != nil {
		// Not in a repo, just show global
		return nil
	}

	repoName, err := git.GetRepoBasename()
	if err != nil {
		return err
	}

	fmt.Printf("\nProject: %s\n", repoName)

	// Display legacy global project patterns (if any remain)
	projectPatterns := cfg.GetProjectPatterns(repoName)
	if len(projectPatterns) > 0 {
		fmt.Println("\nProject-specific patterns (legacy - in global config):")
		for _, pattern := range projectPatterns {
			if strings.HasPrefix(pattern, "+") {
				fmt.Printf("  %s (override - include despite global exclusion)\n", pattern)
			} else {
				fmt.Printf("  %s\n", pattern)
			}
		}
	}

	// Load and display local config
	localCfg, err := config.LoadLocal(repoRoot)
	if err != nil {
		return fmt.Errorf("failed to load local config: %w", err)
	}

	if len(localCfg.Patterns) > 0 {
		fmt.Printf("\nLocal patterns (%s):\n", localCfg.Path())
		for _, pattern := range localCfg.Patterns {
			if strings.HasPrefix(pattern, "+") {
				fmt.Printf("  %s (override - include despite global exclusion)\n", pattern)
			} else {
				fmt.Printf("  %s\n", pattern)
			}
		}
	} else if localCfg.Exists() {
		fmt.Printf("\nLocal patterns (%s):\n", localCfg.Path())
		fmt.Println("  (none)")
	}

	// Display effective exclusions
	effectivePatterns := cfg.GetEffectivePatternsWithLocal(repoName, localCfg.Patterns)
	fmt.Println("\nEffective exclusions:")
	if len(effectivePatterns) == 0 {
		fmt.Println("  (none)")
	} else {
		for _, pattern := range effectivePatterns {
			fmt.Printf("  %s\n", pattern)
		}
	}

	return nil
}

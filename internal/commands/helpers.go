package commands

import "github.com/buckleypaul/haide/internal/config"

// buildExcludePatterns prepends the local config filename to effective patterns
// so that the .haide file itself is excluded from git.
func buildExcludePatterns(effective []string) []string {
	return append([]string{config.LocalConfigFilename}, effective...)
}

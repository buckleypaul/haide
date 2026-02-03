package patterns

import (
	"path/filepath"
	"strings"

	"github.com/gobwas/glob"
)

// Match checks if a file path matches a pattern
func Match(pattern, path string) (bool, error) {
	// Handle directory patterns (ending with /)
	if strings.HasSuffix(pattern, "/") {
		dirPattern := strings.TrimSuffix(pattern, "/")
		// Check if path is in this directory
		return strings.HasPrefix(path, dirPattern+"/") || path == dirPattern, nil
	}

	// Handle ** recursive patterns
	if strings.Contains(pattern, "**") {
		g, err := glob.Compile(pattern, '/')
		if err != nil {
			return false, err
		}
		return g.Match(path), nil
	}

	// Handle patterns with directory components
	if strings.Contains(pattern, "/") {
		match, err := filepath.Match(pattern, path)
		if err != nil {
			return false, err
		}
		if match {
			return true, nil
		}

		// Also try matching against just the filename part
		match, err = filepath.Match(pattern, filepath.Base(path))
		return match, err
	}

	// Simple filename pattern
	return filepath.Match(pattern, filepath.Base(path))
}

// MatchAny checks if a path matches any of the patterns
func MatchAny(patterns []string, path string) (bool, error) {
	for _, pattern := range patterns {
		matched, err := Match(pattern, path)
		if err != nil {
			return false, err
		}
		if matched {
			return true, nil
		}
	}
	return false, nil
}

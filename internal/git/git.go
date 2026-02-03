package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

// GetRepoRoot returns the root directory of the git repository
func GetRepoRoot() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("not in a git repository")
	}
	return strings.TrimSpace(out.String()), nil
}

// GetRepoBasename returns the basename of the repository
func GetRepoBasename() (string, error) {
	root, err := GetRepoRoot()
	if err != nil {
		return "", err
	}
	return filepath.Base(root), nil
}

// IsFileTracked checks if a file is tracked by git
func IsFileTracked(path string) (bool, error) {
	cmd := exec.Command("git", "ls-files", "--error-unmatch", path)
	err := cmd.Run()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// UntrackFile removes a file from git tracking
func UntrackFile(path string) error {
	cmd := exec.Command("git", "rm", "--cached", path)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to untrack file: %s", stderr.String())
	}
	return nil
}

// GetTrackedFiles returns a list of files tracked by git that match any of the patterns
func GetTrackedFiles(patterns []string) ([]string, error) {
	// Get all tracked files
	cmd := exec.Command("git", "ls-files")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to list tracked files: %w", err)
	}

	allFiles := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(allFiles) == 1 && allFiles[0] == "" {
		return []string{}, nil
	}

	// Filter files that match any pattern
	var matched []string
	for _, file := range allFiles {
		for _, pattern := range patterns {
			matches, err := filepath.Match(pattern, filepath.Base(file))
			if err != nil {
				continue
			}
			if matches || matchesPath(file, pattern) {
				matched = append(matched, file)
				break
			}
		}
	}

	return matched, nil
}

// matchesPath checks if a file path matches a pattern that may include directories
func matchesPath(path, pattern string) bool {
	// Handle directory patterns
	if strings.HasSuffix(pattern, "/") {
		return strings.HasPrefix(path, strings.TrimSuffix(pattern, "/"))
	}

	// Handle patterns with wildcards in directory components
	if strings.Contains(pattern, "/") {
		match, err := filepath.Match(pattern, path)
		if err != nil {
			return false
		}
		return match
	}

	return false
}

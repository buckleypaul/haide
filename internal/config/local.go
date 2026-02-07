package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// LocalConfigFilename is the name of the local config file in a project root
const LocalConfigFilename = ".haide"

// LocalConfig represents the local per-project .haide configuration file
type LocalConfig struct {
	Patterns []string
	path     string
}

// LoadLocal loads the local .haide config file from a repository root.
// Returns empty patterns (not an error) if the file doesn't exist.
func LoadLocal(repoRoot string) (*LocalConfig, error) {
	path := filepath.Join(repoRoot, LocalConfigFilename)
	lc := &LocalConfig{
		path: path,
	}

	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return lc, nil
		}
		return nil, fmt.Errorf("failed to open local config: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		lc.Patterns = append(lc.Patterns, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read local config: %w", err)
	}

	return lc, nil
}

// Save writes the local config to disk using atomic temp+rename.
func (lc *LocalConfig) Save() error {
	tmpPath := lc.path + ".tmp"
	f, err := os.Create(tmpPath)
	if err != nil {
		return fmt.Errorf("failed to create temp local config: %w", err)
	}
	defer func() {
		f.Close()
		os.Remove(tmpPath)
	}()

	if _, err := f.WriteString("# haide local project exclusions\n"); err != nil {
		return fmt.Errorf("failed to write local config header: %w", err)
	}
	if _, err := f.WriteString("# Use +PATTERN to override a global exclusion\n\n"); err != nil {
		return fmt.Errorf("failed to write local config header: %w", err)
	}

	for _, pattern := range lc.Patterns {
		if _, err := f.WriteString(pattern + "\n"); err != nil {
			return fmt.Errorf("failed to write local config pattern: %w", err)
		}
	}

	if err := f.Sync(); err != nil {
		return fmt.Errorf("failed to sync local config: %w", err)
	}
	if err := f.Close(); err != nil {
		return fmt.Errorf("failed to close local config: %w", err)
	}

	if err := os.Rename(tmpPath, lc.path); err != nil {
		return fmt.Errorf("failed to save local config: %w", err)
	}

	return nil
}

// AddPattern adds a pattern with deduplication.
func (lc *LocalConfig) AddPattern(pattern string) {
	for _, p := range lc.Patterns {
		if p == pattern {
			return
		}
	}
	lc.Patterns = append(lc.Patterns, pattern)
}

// RemovePattern removes a pattern and returns whether it was found.
func (lc *LocalConfig) RemovePattern(pattern string) bool {
	for i, p := range lc.Patterns {
		if p == pattern {
			lc.Patterns = append(lc.Patterns[:i], lc.Patterns[i+1:]...)
			return true
		}
	}
	return false
}

// HasPattern returns whether the local config contains the given pattern.
func (lc *LocalConfig) HasPattern(pattern string) bool {
	for _, p := range lc.Patterns {
		if p == pattern {
			return true
		}
	}
	return false
}

// Path returns the path to the local config file.
func (lc *LocalConfig) Path() string {
	return lc.path
}

// Exists returns whether the local config file exists on disk.
func (lc *LocalConfig) Exists() bool {
	_, err := os.Stat(lc.path)
	return err == nil
}

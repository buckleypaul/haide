package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Config represents the haide configuration structure
type Config struct {
	Global   []string
	Projects map[string][]string
	path     string
}

// DefaultGlobalExclusions returns the default list of AI files to exclude
func DefaultGlobalExclusions() []string {
	return []string{
		// Claude AI
		"CLAUDE.md",
		"AGENTS.md",
		".claude/",
		".clauderc",

		// Cursor
		".cursorrules",
		".cursorignore",

		// GitHub Copilot
		".github/copilot-instructions.md",

		// Cody
		".cody/",

		// AI-generated documentation
		"AI_NOTES.md",
		"AI_TODO.md",
		".ai/",

		// Common AI skill/plugin directories
		"skills/",
		".skills/",
		"ai-skills/",

		// AI conversation logs
		"ai-conversations/",
		".ai-logs/",

		// AI prompts and templates
		"prompts/",
		".prompts/",
		"ai-prompts/",

		// AI config files
		".aiconfig",
		".ai-config.json",
		".ai-config.yaml",
	}
}

// GetConfigPath returns the path to the config file
func GetConfigPath() string {
	if haideHome := os.Getenv("HAIDE_HOME"); haideHome != "" {
		return filepath.Join(haideHome, "config.ini")
	}
	home, err := os.UserHomeDir()
	if err != nil {
		panic(fmt.Sprintf("failed to get user home directory: %v", err))
	}
	return filepath.Join(home, ".haide", "config.ini")
}

// Load loads the configuration from the config file
func Load() (*Config, error) {
	configPath := GetConfigPath()

	// If config doesn't exist, create default
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return createDefault()
	}

	file, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open config: %w", err)
	}
	defer file.Close()

	config := &Config{
		Projects: make(map[string][]string),
		path:     configPath,
	}

	scanner := bufio.NewScanner(file)
	var currentSection string
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Check for section header
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			sectionName := strings.TrimSpace(line[1 : len(line)-1])

			// Handle [global] section
			if sectionName == "global" {
				currentSection = "global"
				continue
			}

			// Handle [project:name] section
			if strings.HasPrefix(sectionName, "project:") {
				currentSection = strings.TrimPrefix(sectionName, "project:")
				continue
			}

			return nil, fmt.Errorf("line %d: invalid section name %q (must be [global] or [project:name])", lineNum, sectionName)
		}

		// Pattern line
		if currentSection == "" {
			return nil, fmt.Errorf("line %d: pattern %q found outside of section", lineNum, line)
		}

		pattern := strings.TrimSpace(line)
		if currentSection == "global" {
			config.Global = append(config.Global, pattern)
		} else {
			config.Projects[currentSection] = append(config.Projects[currentSection], pattern)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	return config, nil
}

// createDefault creates a default configuration file
func createDefault() (*Config, error) {
	configPath := GetConfigPath()
	configDir := filepath.Dir(configPath)

	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}

	config := &Config{
		Global:   DefaultGlobalExclusions(),
		Projects: make(map[string][]string),
		path:     configPath,
	}

	if err := config.Save(); err != nil {
		return nil, err
	}

	return config, nil
}

// Save saves the configuration to disk
func (c *Config) Save() error {
	// Create temp file for atomic write
	tmpPath := c.path + ".tmp"
	f, err := os.Create(tmpPath)
	if err != nil {
		return fmt.Errorf("failed to create temp config file: %w", err)
	}
	defer func() {
		f.Close()
		os.Remove(tmpPath) // Clean up on error
	}()

	// Write header comments
	if _, err := f.WriteString("# haide configuration file\n"); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}
	if _, err := f.WriteString("# Global exclusions apply to all repositories\n"); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}
	if _, err := f.WriteString("# Project-specific sections use [project:name] format\n"); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}
	if _, err := f.WriteString("# Use +PATTERN to override a global exclusion for specific project\n\n"); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	// Write global section
	if _, err := f.WriteString("[global]\n"); err != nil {
		return fmt.Errorf("failed to write global section: %w", err)
	}
	for _, pattern := range c.Global {
		if _, err := f.WriteString(pattern + "\n"); err != nil {
			return fmt.Errorf("failed to write global pattern: %w", err)
		}
	}

	// Write project sections
	for projectName, patterns := range c.Projects {
		if _, err := f.WriteString("\n[project:" + projectName + "]\n"); err != nil {
			return fmt.Errorf("failed to write project section: %w", err)
		}
		for _, pattern := range patterns {
			if _, err := f.WriteString(pattern + "\n"); err != nil {
				return fmt.Errorf("failed to write project pattern: %w", err)
			}
		}
	}

	// Sync and close
	if err := f.Sync(); err != nil {
		return fmt.Errorf("failed to sync config: %w", err)
	}
	if err := f.Close(); err != nil {
		return fmt.Errorf("failed to close config: %w", err)
	}

	// Atomic rename
	if err := os.Rename(tmpPath, c.path); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	return nil
}

// GetProjectPatterns returns patterns for a specific project
func (c *Config) GetProjectPatterns(projectName string) []string {
	return c.Projects[projectName]
}

// AddProjectPattern adds a pattern to a project's exclusion list
func (c *Config) AddProjectPattern(projectName, pattern string) {
	if c.Projects == nil {
		c.Projects = make(map[string][]string)
	}

	// Check if pattern already exists
	for _, p := range c.Projects[projectName] {
		if p == pattern {
			return
		}
	}

	c.Projects[projectName] = append(c.Projects[projectName], pattern)
}

// RemoveProjectPattern removes a pattern from a project's exclusion list
func (c *Config) RemoveProjectPattern(projectName, pattern string) {
	patterns := c.Projects[projectName]
	newPatterns := make([]string, 0, len(patterns))
	for _, p := range patterns {
		if p != pattern {
			newPatterns = append(newPatterns, p)
		}
	}
	if len(newPatterns) > 0 {
		c.Projects[projectName] = newPatterns
	} else {
		delete(c.Projects, projectName)
	}
}

// AddGlobalPattern adds a pattern to the global exclusion list
func (c *Config) AddGlobalPattern(pattern string) {
	// Check if pattern already exists
	for _, p := range c.Global {
		if p == pattern {
			return
		}
	}
	c.Global = append(c.Global, pattern)
}

// RemoveGlobalPattern removes a pattern from the global exclusion list
func (c *Config) RemoveGlobalPattern(pattern string) {
	newGlobal := make([]string, 0, len(c.Global))
	for _, p := range c.Global {
		if p != pattern {
			newGlobal = append(newGlobal, p)
		}
	}
	c.Global = newGlobal
}

// GetEffectivePatterns returns the effective patterns for a project
// (global patterns + project patterns - overrides)
func (c *Config) GetEffectivePatterns(projectName string) []string {
	projectPatterns := c.GetProjectPatterns(projectName)

	// Build override set
	overrides := make(map[string]bool)
	for _, pattern := range projectPatterns {
		if strings.HasPrefix(pattern, "+") {
			overrides[strings.TrimPrefix(pattern, "+")] = true
		}
	}

	// Start with global patterns, excluding overridden ones
	var effective []string
	for _, pattern := range c.Global {
		if !overrides[pattern] {
			effective = append(effective, pattern)
		}
	}

	// Add project-specific patterns (non-override)
	for _, pattern := range projectPatterns {
		if !strings.HasPrefix(pattern, "+") {
			effective = append(effective, pattern)
		}
	}

	return effective
}

// Path returns the path to the config file
func (c *Config) Path() string {
	return c.path
}

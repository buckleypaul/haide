package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

// Config represents the haide configuration structure
type Config struct {
	Global   []string            `toml:"global"`
	Projects map[string][]string `toml:"-"`
	path     string
	raw      map[string]interface{}
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
		return filepath.Join(haideHome, "config.toml")
	}
	home, err := os.UserHomeDir()
	if err != nil {
		panic(fmt.Sprintf("failed to get user home directory: %v", err))
	}
	return filepath.Join(home, ".haide", "config.toml")
}

// Load loads the configuration from the config file
func Load() (*Config, error) {
	configPath := GetConfigPath()

	// If config doesn't exist, create default
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return createDefault()
	}

	var raw map[string]interface{}
	if _, err := toml.DecodeFile(configPath, &raw); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	config := &Config{
		Projects: make(map[string][]string),
		path:     configPath,
		raw:      raw,
	}

	// Parse global section
	if globalRaw, ok := raw["global"].([]interface{}); ok {
		for _, item := range globalRaw {
			if str, ok := item.(string); ok {
				config.Global = append(config.Global, str)
			}
		}
	}

	// Parse project sections
	for key, value := range raw {
		if key == "global" {
			continue
		}
		if patterns, ok := value.([]interface{}); ok {
			var projectPatterns []string
			for _, item := range patterns {
				if str, ok := item.(string); ok {
					projectPatterns = append(projectPatterns, str)
				}
			}
			config.Projects[key] = projectPatterns
		}
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
		raw:      make(map[string]interface{}),
	}

	if err := config.Save(); err != nil {
		return nil, err
	}

	return config, nil
}

// Save saves the configuration to disk
func (c *Config) Save() error {
	// Build the TOML structure
	data := make(map[string]interface{})
	data["global"] = c.Global

	// Add project sections
	for project, patterns := range c.Projects {
		data[project] = patterns
	}

	// Open file for writing
	f, err := os.Create(c.path)
	if err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}
	defer f.Close()

	// Write header comment
	f.WriteString("# haide configuration file\n")
	f.WriteString("# Global exclusions apply to all repositories\n")
	f.WriteString("# Project-specific sections are keyed by repository basename\n")
	f.WriteString("# Use +PATTERN to override a global exclusion for a specific project\n\n")

	// Encode TOML
	encoder := toml.NewEncoder(f)
	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("failed to encode config: %w", err)
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

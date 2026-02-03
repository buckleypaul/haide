package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetEffectivePatterns(t *testing.T) {
	cfg := &Config{
		Global: []string{"CLAUDE.md", "AGENTS.md", ".ai/"},
		Projects: map[string][]string{
			"test-repo": {"+CLAUDE.md", "custom.md"},
		},
	}

	effective := cfg.GetEffectivePatterns("test-repo")

	// Should include: AGENTS.md, .ai/ (from global, minus override), custom.md (project-specific)
	expected := map[string]bool{
		"AGENTS.md":  true,
		".ai/":       true,
		"custom.md":  true,
	}

	// Should NOT include CLAUDE.md (overridden)
	for _, pattern := range effective {
		if pattern == "CLAUDE.md" {
			t.Errorf("CLAUDE.md should be overridden but was included")
		}
		if pattern == "+CLAUDE.md" {
			t.Errorf("Override pattern should not appear in effective list")
		}
	}

	// Check expected patterns are present
	for _, pattern := range effective {
		if expected[pattern] {
			delete(expected, pattern)
		}
	}

	if len(expected) > 0 {
		t.Errorf("Missing expected patterns: %v", expected)
	}
}

func TestAddRemoveProjectPattern(t *testing.T) {
	cfg := &Config{
		Global:   []string{"CLAUDE.md"},
		Projects: make(map[string][]string),
	}

	// Add pattern
	cfg.AddProjectPattern("test-repo", "custom.md")
	patterns := cfg.GetProjectPatterns("test-repo")
	if len(patterns) != 1 || patterns[0] != "custom.md" {
		t.Errorf("Expected custom.md, got %v", patterns)
	}

	// Add duplicate (should be ignored)
	cfg.AddProjectPattern("test-repo", "custom.md")
	patterns = cfg.GetProjectPatterns("test-repo")
	if len(patterns) != 1 {
		t.Errorf("Duplicate pattern was added")
	}

	// Remove pattern
	cfg.RemoveProjectPattern("test-repo", "custom.md")
	patterns = cfg.GetProjectPatterns("test-repo")
	if len(patterns) != 0 {
		t.Errorf("Pattern was not removed")
	}
}

func TestAddRemoveGlobalPattern(t *testing.T) {
	cfg := &Config{
		Global:   []string{"CLAUDE.md"},
		Projects: make(map[string][]string),
	}

	// Add pattern
	cfg.AddGlobalPattern("AGENTS.md")
	if len(cfg.Global) != 2 {
		t.Errorf("Expected 2 global patterns, got %d", len(cfg.Global))
	}

	// Add duplicate (should be ignored)
	cfg.AddGlobalPattern("AGENTS.md")
	if len(cfg.Global) != 2 {
		t.Errorf("Duplicate global pattern was added")
	}

	// Remove pattern
	cfg.RemoveGlobalPattern("AGENTS.md")
	if len(cfg.Global) != 1 || cfg.Global[0] != "CLAUDE.md" {
		t.Errorf("Pattern was not removed correctly")
	}
}

func TestSaveLoad(t *testing.T) {
	// Create temp directory for config
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.toml")

	cfg := &Config{
		Global: []string{"CLAUDE.md", "AGENTS.md"},
		Projects: map[string][]string{
			"test-repo": {"+CLAUDE.md", "custom.md"},
		},
		path: configPath,
	}

	// Save
	if err := cfg.Save(); err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Set environment variable to use temp directory
	oldHaideHome := os.Getenv("HAIDE_HOME")
	os.Setenv("HAIDE_HOME", tmpDir)
	defer os.Setenv("HAIDE_HOME", oldHaideHome)

	// Load
	loaded, err := Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Verify global patterns
	if len(loaded.Global) != len(cfg.Global) {
		t.Errorf("Expected %d global patterns, got %d", len(cfg.Global), len(loaded.Global))
	}

	// Verify project patterns
	testRepoPatterns := loaded.GetProjectPatterns("test-repo")
	if len(testRepoPatterns) != 2 {
		t.Errorf("Expected 2 project patterns, got %d", len(testRepoPatterns))
	}
}

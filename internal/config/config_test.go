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
	configPath := filepath.Join(tmpDir, "config.ini")

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

func TestGetEffectivePatternsWithLocal(t *testing.T) {
	cfg := &Config{
		Global: []string{"CLAUDE.md", "AGENTS.md", ".ai/"},
		Projects: map[string][]string{
			"test-repo": {"legacy-extra.md"},
		},
	}

	localPatterns := []string{"+CLAUDE.md", "custom.md"}

	effective := cfg.GetEffectivePatternsWithLocal("test-repo", localPatterns)

	expected := map[string]bool{
		"AGENTS.md":       true,
		".ai/":            true,
		"legacy-extra.md": true,
		"custom.md":       true,
	}

	for _, pattern := range effective {
		if pattern == "CLAUDE.md" {
			t.Errorf("CLAUDE.md should be overridden but was included")
		}
		if pattern == "+CLAUDE.md" {
			t.Errorf("Override pattern should not appear in effective list")
		}
		delete(expected, pattern)
	}

	if len(expected) > 0 {
		t.Errorf("Missing expected patterns: %v", expected)
	}
}

func TestGetEffectivePatternsWithLocalDedup(t *testing.T) {
	cfg := &Config{
		Global: []string{"CLAUDE.md", "AGENTS.md"},
		Projects: map[string][]string{
			"test-repo": {"custom.md"},
		},
	}

	// local also has custom.md — should not duplicate
	localPatterns := []string{"custom.md", "extra.md"}

	effective := cfg.GetEffectivePatternsWithLocal("test-repo", localPatterns)

	count := 0
	for _, p := range effective {
		if p == "custom.md" {
			count++
		}
	}
	if count != 1 {
		t.Errorf("Expected custom.md once, got %d times", count)
	}
}

func TestMigrateProjectToLocal(t *testing.T) {
	cfg := &Config{
		Global: []string{"CLAUDE.md"},
		Projects: map[string][]string{
			"test-repo": {"+CLAUDE.md", "custom.md"},
			"other":     {"other.md"},
		},
	}

	extracted := cfg.MigrateProjectToLocal("test-repo")

	if len(extracted) != 2 {
		t.Fatalf("Expected 2 extracted patterns, got %d", len(extracted))
	}
	if extracted[0] != "+CLAUDE.md" || extracted[1] != "custom.md" {
		t.Errorf("Unexpected extracted patterns: %v", extracted)
	}

	// Section should be removed
	if _, ok := cfg.Projects["test-repo"]; ok {
		t.Error("Project section should be removed after migration")
	}

	// Other project should be unaffected
	if _, ok := cfg.Projects["other"]; !ok {
		t.Error("Other project section should still exist")
	}
}

func TestMigrateProjectToLocalEmpty(t *testing.T) {
	cfg := &Config{
		Global:   []string{"CLAUDE.md"},
		Projects: map[string][]string{},
	}

	extracted := cfg.MigrateProjectToLocal("nonexistent")
	if extracted != nil {
		t.Errorf("Expected nil for nonexistent project, got %v", extracted)
	}
}

func TestParseINIFormat(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		wantGlobal  []string
		wantProject map[string][]string
		wantErr     bool
	}{
		{
			name: "valid config with comments",
			content: `# Comment at top
[global]
CLAUDE.md
# Comment in section
AGENTS.md

[project:myrepo]
custom.md
+CLAUDE.md
`,
			wantGlobal: []string{"CLAUDE.md", "AGENTS.md"},
			wantProject: map[string][]string{
				"myrepo": {"+CLAUDE.md", "custom.md"},
			},
		},
		{
			name: "empty lines and whitespace",
			content: `[global]

CLAUDE.md

AGENTS.md

[project:test]

  custom.md
`,
			wantGlobal: []string{"CLAUDE.md", "AGENTS.md"},
			wantProject: map[string][]string{
				"test": {"custom.md"},
			},
		},
		{
			name: "pattern outside section",
			content: `CLAUDE.md
[global]
AGENTS.md
`,
			wantErr: true,
		},
		{
			name: "invalid section name",
			content: `[invalid]
pattern.md
`,
			wantErr: true,
		},
		{
			name: "only global section",
			content: `[global]
CLAUDE.md
AGENTS.md
`,
			wantGlobal:  []string{"CLAUDE.md", "AGENTS.md"},
			wantProject: map[string][]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			configPath := filepath.Join(tmpDir, "config.ini")

			if err := os.WriteFile(configPath, []byte(tt.content), 0644); err != nil {
				t.Fatalf("Failed to write test config: %v", err)
			}

			oldHaideHome := os.Getenv("HAIDE_HOME")
			os.Setenv("HAIDE_HOME", tmpDir)
			defer os.Setenv("HAIDE_HOME", oldHaideHome)

			cfg, err := Load()
			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			// Check global patterns
			if len(cfg.Global) != len(tt.wantGlobal) {
				t.Errorf("Global patterns count: got %d, want %d", len(cfg.Global), len(tt.wantGlobal))
			}
			for i, want := range tt.wantGlobal {
				if i >= len(cfg.Global) || cfg.Global[i] != want {
					t.Errorf("Global pattern %d: got %q, want %q", i, cfg.Global[i], want)
				}
			}

			// Check project patterns
			if len(cfg.Projects) != len(tt.wantProject) {
				t.Errorf("Project count: got %d, want %d", len(cfg.Projects), len(tt.wantProject))
			}
			for projName, wantPatterns := range tt.wantProject {
				gotPatterns := cfg.Projects[projName]
				if len(gotPatterns) != len(wantPatterns) {
					t.Errorf("Project %q patterns count: got %d, want %d", projName, len(gotPatterns), len(wantPatterns))
				}
			}
		})
	}
}

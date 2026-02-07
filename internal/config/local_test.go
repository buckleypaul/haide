package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadLocalNonexistent(t *testing.T) {
	tmpDir := t.TempDir()
	lc, err := LoadLocal(tmpDir)
	if err != nil {
		t.Fatalf("Expected no error for nonexistent file, got: %v", err)
	}
	if len(lc.Patterns) != 0 {
		t.Errorf("Expected 0 patterns, got %d", len(lc.Patterns))
	}
	if lc.Exists() {
		t.Error("Exists() should return false for nonexistent file")
	}
}

func TestLoadLocalWithPatterns(t *testing.T) {
	tmpDir := t.TempDir()
	content := `# haide local project exclusions
# Use +PATTERN to override a global exclusion

custom-ai-tool.md
my-prompts/

+CLAUDE.md
# another comment
some-file.txt
`
	if err := os.WriteFile(filepath.Join(tmpDir, LocalConfigFilename), []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	lc, err := LoadLocal(tmpDir)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expected := []string{"custom-ai-tool.md", "my-prompts/", "+CLAUDE.md", "some-file.txt"}
	if len(lc.Patterns) != len(expected) {
		t.Fatalf("Expected %d patterns, got %d: %v", len(expected), len(lc.Patterns), lc.Patterns)
	}
	for i, want := range expected {
		if lc.Patterns[i] != want {
			t.Errorf("Pattern %d: got %q, want %q", i, lc.Patterns[i], want)
		}
	}
}

func TestLocalSaveLoadRoundTrip(t *testing.T) {
	tmpDir := t.TempDir()
	lc := &LocalConfig{
		Patterns: []string{"foo.md", "bar/", "+CLAUDE.md"},
		path:     filepath.Join(tmpDir, LocalConfigFilename),
	}

	if err := lc.Save(); err != nil {
		t.Fatalf("Failed to save: %v", err)
	}

	loaded, err := LoadLocal(tmpDir)
	if err != nil {
		t.Fatalf("Failed to load: %v", err)
	}

	if len(loaded.Patterns) != len(lc.Patterns) {
		t.Fatalf("Expected %d patterns, got %d", len(lc.Patterns), len(loaded.Patterns))
	}
	for i, want := range lc.Patterns {
		if loaded.Patterns[i] != want {
			t.Errorf("Pattern %d: got %q, want %q", i, loaded.Patterns[i], want)
		}
	}
}

func TestLocalAddPatternDedup(t *testing.T) {
	lc := &LocalConfig{}

	lc.AddPattern("foo.md")
	if len(lc.Patterns) != 1 {
		t.Fatalf("Expected 1 pattern, got %d", len(lc.Patterns))
	}

	lc.AddPattern("foo.md")
	if len(lc.Patterns) != 1 {
		t.Errorf("Duplicate pattern was added, expected 1 got %d", len(lc.Patterns))
	}

	lc.AddPattern("bar.md")
	if len(lc.Patterns) != 2 {
		t.Errorf("Expected 2 patterns, got %d", len(lc.Patterns))
	}
}

func TestLocalRemovePattern(t *testing.T) {
	lc := &LocalConfig{
		Patterns: []string{"foo.md", "bar.md", "baz.md"},
	}

	if !lc.RemovePattern("bar.md") {
		t.Error("RemovePattern should return true for existing pattern")
	}
	if len(lc.Patterns) != 2 {
		t.Errorf("Expected 2 patterns after removal, got %d", len(lc.Patterns))
	}

	if lc.RemovePattern("nonexistent.md") {
		t.Error("RemovePattern should return false for nonexistent pattern")
	}
}

func TestLocalHasPattern(t *testing.T) {
	lc := &LocalConfig{
		Patterns: []string{"foo.md", "+CLAUDE.md"},
	}

	if !lc.HasPattern("foo.md") {
		t.Error("HasPattern should return true for existing pattern")
	}
	if !lc.HasPattern("+CLAUDE.md") {
		t.Error("HasPattern should return true for override pattern")
	}
	if lc.HasPattern("nonexistent.md") {
		t.Error("HasPattern should return false for nonexistent pattern")
	}
}

func TestLocalExists(t *testing.T) {
	tmpDir := t.TempDir()
	lc := &LocalConfig{
		Patterns: []string{"test.md"},
		path:     filepath.Join(tmpDir, LocalConfigFilename),
	}

	if lc.Exists() {
		t.Error("Exists should return false before save")
	}

	if err := lc.Save(); err != nil {
		t.Fatalf("Failed to save: %v", err)
	}

	if !lc.Exists() {
		t.Error("Exists should return true after save")
	}
}

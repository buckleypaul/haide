package exclude

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadNonExistent(t *testing.T) {
	tmpDir := t.TempDir()
	repoRoot := filepath.Join(tmpDir, "test-repo")
	os.MkdirAll(filepath.Join(repoRoot, ".git", "info"), 0755)

	f, err := Load(repoRoot)
	if err != nil {
		t.Fatalf("Failed to load non-existent exclude file: %v", err)
	}

	if f.hasManagedBlock {
		t.Error("Should not have managed block for non-existent file")
	}
}

func TestSaveLoad(t *testing.T) {
	tmpDir := t.TempDir()
	repoRoot := filepath.Join(tmpDir, "test-repo")
	os.MkdirAll(filepath.Join(repoRoot, ".git", "info"), 0755)

	// Create and save
	f := &File{
		path:            filepath.Join(repoRoot, ".git", "info", "exclude"),
		beforeHaide:     []string{"# Existing comment", "*.log"},
		haideContent:    []string{"CLAUDE.md", "AGENTS.md"},
		afterHaide:      []string{"# After comment"},
		hasManagedBlock: true,
	}

	if err := f.Save(); err != nil {
		t.Fatalf("Failed to save: %v", err)
	}

	// Load
	loaded, err := Load(repoRoot)
	if err != nil {
		t.Fatalf("Failed to load: %v", err)
	}

	if !loaded.hasManagedBlock {
		t.Error("Should have managed block")
	}

	if len(loaded.haideContent) != 2 {
		t.Errorf("Expected 2 haide patterns, got %d", len(loaded.haideContent))
	}

	// Note: Save() adds blank lines between sections, so we check for content
	hasExistingComment := false
	hasLogPattern := false
	for _, line := range loaded.beforeHaide {
		if line == "# Existing comment" {
			hasExistingComment = true
		}
		if line == "*.log" {
			hasLogPattern = true
		}
	}
	if !hasExistingComment || !hasLogPattern {
		t.Error("Before-haide content not preserved correctly")
	}

	hasAfterComment := false
	for _, line := range loaded.afterHaide {
		if line == "# After comment" {
			hasAfterComment = true
		}
	}
	if !hasAfterComment {
		t.Error("After-haide content not preserved correctly")
	}
}

func TestSetManagedContent(t *testing.T) {
	tmpDir := t.TempDir()
	repoRoot := filepath.Join(tmpDir, "test-repo")
	os.MkdirAll(filepath.Join(repoRoot, ".git", "info"), 0755)

	f := &File{
		path: filepath.Join(repoRoot, ".git", "info", "exclude"),
	}

	patterns := []string{"CLAUDE.md", "AGENTS.md"}
	f.SetManagedContent(patterns)

	if !f.hasManagedBlock {
		t.Error("Should have managed block after setting content")
	}

	if len(f.haideContent) != 2 {
		t.Errorf("Expected 2 patterns, got %d", len(f.haideContent))
	}
}

func TestClearManagedContent(t *testing.T) {
	f := &File{
		haideContent:    []string{"CLAUDE.md"},
		hasManagedBlock: true,
	}

	f.ClearManagedContent()

	if f.hasManagedBlock {
		t.Error("Should not have managed block after clearing")
	}

	if len(f.haideContent) != 0 {
		t.Error("Content should be empty after clearing")
	}
}

func TestGetDiff(t *testing.T) {
	f := &File{
		haideContent: []string{"CLAUDE.md", "AGENTS.md"},
	}

	newPatterns := []string{"CLAUDE.md", ".ai/"}

	diff := f.GetDiff(newPatterns)

	if !strings.Contains(diff, "+ .ai/") {
		t.Error("Diff should show .ai/ as addition")
	}

	if !strings.Contains(diff, "- AGENTS.md") {
		t.Error("Diff should show AGENTS.md as removal")
	}

	if strings.Contains(diff, "CLAUDE.md") {
		t.Error("Diff should not include unchanged pattern")
	}
}

func TestGetDiffNoChanges(t *testing.T) {
	f := &File{
		haideContent: []string{"CLAUDE.md", "AGENTS.md"},
	}

	newPatterns := []string{"CLAUDE.md", "AGENTS.md"}

	diff := f.GetDiff(newPatterns)

	if diff != "No changes" {
		t.Errorf("Expected 'No changes', got: %s", diff)
	}
}

package patterns

import "testing"

func TestMatch(t *testing.T) {
	tests := []struct {
		pattern string
		path    string
		want    bool
	}{
		// Simple filename patterns
		{"CLAUDE.md", "CLAUDE.md", true},
		{"CLAUDE.md", "AGENTS.md", false},
		{"*.md", "README.md", true},
		{"*.md", "README.txt", false},

		// Directory patterns
		{".ai/", ".ai/notes.txt", true},
		{".ai/", ".ai/", true},
		{".ai/", "other/file.txt", false},

		// Path patterns
		{".github/copilot-instructions.md", ".github/copilot-instructions.md", true},
		{".github/*.md", ".github/README.md", true},

		// Recursive patterns
		{"**/*.json", "deep/nested/config.json", true},
		{"**/*.json", "file.txt", false},

		// Prefix patterns
		{"ai-*", "ai-notes.md", true},
		{"ai-*", "notes.md", false},
	}

	for _, tt := range tests {
		got, err := Match(tt.pattern, tt.path)
		if err != nil {
			t.Errorf("Match(%q, %q) error: %v", tt.pattern, tt.path, err)
			continue
		}
		if got != tt.want {
			t.Errorf("Match(%q, %q) = %v, want %v", tt.pattern, tt.path, got, tt.want)
		}
	}
}

func TestMatchAny(t *testing.T) {
	patterns := []string{"*.md", ".ai/", "CLAUDE.md"}

	tests := []struct {
		path string
		want bool
	}{
		{"README.md", true},
		{".ai/notes.txt", true},
		{"CLAUDE.md", true},
		{"file.txt", false},
		{"other/README.md", true},
	}

	for _, tt := range tests {
		got, err := MatchAny(patterns, tt.path)
		if err != nil {
			t.Errorf("MatchAny(%v, %q) error: %v", patterns, tt.path, err)
			continue
		}
		if got != tt.want {
			t.Errorf("MatchAny(%v, %q) = %v, want %v", patterns, tt.path, got, tt.want)
		}
	}
}

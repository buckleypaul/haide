package exclude

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	BeginMarker = "# BEGIN HAIDE"
	EndMarker   = "# END HAIDE"
)

// File represents a git info/exclude file
type File struct {
	path            string
	beforeHaide     []string
	haideContent    []string
	afterHaide      []string
	hasManagedBlock bool
}

// Load loads the exclude file from a repository
func Load(repoRoot string) (*File, error) {
	path := filepath.Join(repoRoot, ".git", "info", "exclude")

	f := &File{
		path: path,
	}

	// If file doesn't exist, that's okay
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return f, nil
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open exclude file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	inManagedBlock := false
	var currentSection *[]string

	for scanner.Scan() {
		line := scanner.Text()

		if line == BeginMarker {
			inManagedBlock = true
			f.hasManagedBlock = true
			continue
		}

		if line == EndMarker {
			inManagedBlock = false
			currentSection = &f.afterHaide
			continue
		}

		if inManagedBlock {
			f.haideContent = append(f.haideContent, line)
		} else if currentSection == nil {
			f.beforeHaide = append(f.beforeHaide, line)
		} else {
			*currentSection = append(*currentSection, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read exclude file: %w", err)
	}

	return f, nil
}

// SetManagedContent sets the content managed by haide
func (f *File) SetManagedContent(patterns []string) {
	f.haideContent = patterns
	f.hasManagedBlock = true
}

// GetManagedContent returns the content currently managed by haide
func (f *File) GetManagedContent() []string {
	return f.haideContent
}

// ClearManagedContent removes the haide-managed block
func (f *File) ClearManagedContent() {
	f.haideContent = nil
	f.hasManagedBlock = false
}

// Save writes the exclude file back to disk
func (f *File) Save() error {
	// Ensure the directory exists
	dir := filepath.Dir(f.path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create .git/info directory: %w", err)
	}

	// Create temporary file
	tmpPath := f.path + ".tmp"
	file, err := os.Create(tmpPath)
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %w", err)
	}

	// Write content
	writer := bufio.NewWriter(file)

	// Write before-haide content
	for _, line := range f.beforeHaide {
		fmt.Fprintln(writer, line)
	}

	// Write haide-managed content if present
	if f.hasManagedBlock && len(f.haideContent) > 0 {
		// Add blank line before haide block if there's content before
		if len(f.beforeHaide) > 0 && f.beforeHaide[len(f.beforeHaide)-1] != "" {
			fmt.Fprintln(writer)
		}

		fmt.Fprintln(writer, BeginMarker)
		for _, pattern := range f.haideContent {
			if pattern != "" { // Skip empty lines in managed content
				fmt.Fprintln(writer, pattern)
			}
		}
		fmt.Fprintln(writer, EndMarker)
	}

	// Write after-haide content
	if len(f.afterHaide) > 0 {
		// Add blank line before after content if there's haide content
		if f.hasManagedBlock && len(f.haideContent) > 0 {
			fmt.Fprintln(writer)
		}
		for _, line := range f.afterHaide {
			fmt.Fprintln(writer, line)
		}
	}

	if err := writer.Flush(); err != nil {
		file.Close()
		os.Remove(tmpPath)
		return fmt.Errorf("failed to write exclude file: %w", err)
	}

	file.Close()

	// Atomic rename
	if err := os.Rename(tmpPath, f.path); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("failed to save exclude file: %w", err)
	}

	return nil
}

// GetDiff returns a human-readable diff of what will change
func (f *File) GetDiff(newPatterns []string) string {
	var diff strings.Builder

	// Get current patterns as a set
	currentSet := make(map[string]bool)
	for _, p := range f.haideContent {
		if p != "" {
			currentSet[p] = true
		}
	}

	// Get new patterns as a set
	newSet := make(map[string]bool)
	for _, p := range newPatterns {
		if p != "" {
			newSet[p] = true
		}
	}

	// Find additions
	var additions []string
	for pattern := range newSet {
		if !currentSet[pattern] {
			additions = append(additions, pattern)
		}
	}

	// Find removals
	var removals []string
	for pattern := range currentSet {
		if !newSet[pattern] {
			removals = append(removals, pattern)
		}
	}

	if len(additions) > 0 {
		diff.WriteString("Additions:\n")
		for _, pattern := range additions {
			diff.WriteString(fmt.Sprintf("  + %s\n", pattern))
		}
	}

	if len(removals) > 0 {
		if len(additions) > 0 {
			diff.WriteString("\n")
		}
		diff.WriteString("Removals:\n")
		for _, pattern := range removals {
			diff.WriteString(fmt.Sprintf("  - %s\n", pattern))
		}
	}

	if len(additions) == 0 && len(removals) == 0 {
		return "No changes"
	}

	return diff.String()
}

// Path returns the path to the exclude file
func (f *File) Path() string {
	return f.path
}

# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

haide is a CLI tool written in Go that manages AI-specific file exclusions across git repositories. It uses a centralized INI-style configuration file (~/.haide/config.ini) that maintains global exclusion patterns and per-project overrides, writing patterns to each repository's `.git/info/exclude` file.

## Development Commands

### Building
```bash
make build           # Build binary to ./haide
go build -o haide ./cmd/haide
```

### Testing
```bash
make test            # Run all unit tests
go test ./...        # Run unit tests directly
go test -v ./...     # Run unit tests with verbose output
./integration_test.sh # Run integration tests
```

### Code Quality
```bash
make fmt             # Format all Go files
make lint            # Run golangci-lint
go fmt ./...         # Format code directly
golangci-lint run    # Run linter directly
```

### Installation
```bash
make install         # Build and copy to /usr/local/bin/
```

## Architecture

### Core Components

The application follows a clean layered architecture:

**Commands Layer** (`internal/commands/`): Cobra-based CLI commands that orchestrate operations. Each command (init, add, add-global, remove, remove-global, update, info, clean) handles user interaction, dry-run logic, and coordinates between config, exclude, and git packages.

**Config Layer** (`internal/config/`): Manages the centralized INI-style configuration at `~/.haide/config.ini`. Uses a custom parser (no external dependencies) to read simple section-based format. The Config struct maintains global exclusions and per-project patterns. Key concept: projects are keyed by repository basename (from `filepath.Base(gitRoot)`). Override patterns are prefixed with `+` (e.g., `+CLAUDE.md`) to include a globally-excluded file in a specific project.

**Exclude Layer** (`internal/exclude/`): Manages `.git/info/exclude` files within git repositories. Uses delimited blocks (`# BEGIN HAIDE` / `# END HAIDE`) to mark managed content while preserving external entries. The File struct maintains three sections: beforeHaide, haideContent, and afterHaide.

**Git Layer** (`internal/git/`): Wraps git operations via exec.Command. Provides utilities for finding repo root, checking if files are tracked, and listing tracked files matching patterns.

**Patterns Layer** (`internal/patterns/`): Pattern matching logic using filepath.Match and gobwas/glob for recursive patterns. Handles directory patterns (trailing `/`), recursive patterns (`**`), and simple glob patterns.

### Key Design Decisions

**One-Way Sync**: Configuration is the source of truth. Changes flow from `~/.haide/config.ini` → `.git/info/exclude` only. Manual edits to `.git/info/exclude` within haide markers are overwritten on next update.

**Repository Identification**: Projects are keyed by repository basename. This means two repos with the same basename share configuration (documented limitation).

**Override Mechanism**: When a file is already tracked but matches a global exclusion, haide automatically adds a `+PATTERN` override to the project config during `init`. This allows per-project exceptions to global rules.

**Atomic Operations**: Exclude file updates use temp files with atomic rename to prevent corruption.

### Testing Strategy

Tests use Go's standard testing package with table-driven tests. Each package has corresponding `_test.go` files:
- `config_test.go`: Configuration loading, saving, pattern management
- `exclude_test.go`: Exclude file parsing, content management, diff generation
- `patterns_test.go`: Pattern matching for various glob scenarios

Integration tests in `integration_test.sh` test end-to-end workflows in temporary git repositories.

## Important Patterns

### Pattern Matching Behavior
- `CLAUDE.md` - matches file in any directory (checked against basename)
- `.claude/` - matches directory (trailing slash required)
- `**/*.json` - recursive matching using gobwas/glob
- `ai-*` - glob patterns supported via filepath.Match

### Dry-Run Implementation
All commands support `--dry-run` flag. Commands should check the flag and show what would happen without making changes. Use pattern:
```go
if dryRun {
    fmt.Println("Would perform operation X")
    return nil
}
// perform operation
```

### Error Handling
Commands return errors to be handled by main(). User-facing messages go to stderr via fmt.Fprintf(os.Stderr, ...).

## Release Process

Releases use GitHub Actions + GoReleaser:
1. Update CHANGELOG.md
2. Create and push tag: `git tag vX.Y.Z && git push origin vX.Y.Z`
3. GitHub Actions builds binaries for multiple platforms
4. Homebrew formula auto-updates via buckleypaul/homebrew-tap

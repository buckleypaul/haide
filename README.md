# haide - AI File Exclusion Manager for Git

haide manages AI-specific file exclusions across git repositories using a centralized configuration with global and per-project overrides.

## Features

- Centralized management of AI-specific file exclusions
- Global exclusions applied to all repositories
- Per-project overrides to include files that are globally excluded
- Pattern matching support (wildcards)
- Automatic detection of already-tracked files
- Dry-run mode for all commands
- Simple CLI interface

## Installation

### Using Homebrew (macOS/Linux)

```bash
brew install buckleypaul/tap/haide
```

### From Source

```bash
go install github.com/buckleypaul/haide/cmd/haide@latest
```

Or build locally:

```bash
git clone https://github.com/buckleypaul/haide.git
cd haide
go build -o haide ./cmd/haide
```

## Quick Start

1. Initialize haide in your repository:
```bash
cd your-repo
haide init
```

2. View current configuration:
```bash
haide info
```

3. Add a project-specific exclusion:
```bash
haide add my-ai-notes.md
```

4. Add a global exclusion (applies to all repos):
```bash
haide add-global '*.ai'
```

5. Update other repositories:
```bash
cd other-repo
haide update
```

## Commands

### `haide init`
Initialize haide in the current repository. Creates/updates `.git/info/exclude` with haide-managed entries.

Options:
- `--dry-run`: Preview changes without applying them

### `haide add PATTERN`
Add a pattern to project-specific exclusions.

Options:
- `--dry-run`: Preview changes without applying them

### `haide add-global PATTERN`
Add a pattern to global exclusions (applies to all repositories).

Options:
- `--dry-run`: Preview changes without applying them

### `haide remove PATTERN`
Remove a pattern from project-specific exclusions.

Options:
- `--dry-run`: Preview changes without applying them

### `haide remove-global PATTERN`
Remove a pattern from global exclusions.

Options:
- `--dry-run`: Preview changes without applying them

### `haide update`
Update repository exclusions from config. Shows diff and prompts for confirmation.

Options:
- `--dry-run`: Preview changes without applying them
- `-y, --yes`: Auto-confirm without prompting

### `haide info`
Display current configuration including global and project-specific patterns.

### `haide clean`
Remove haide management from repository. Preserves non-haide entries.

Options:
- `--dry-run`: Preview changes without applying them

## Configuration

haide stores configuration at `~/.haide/config.ini` (or `$HAIDE_HOME/config.ini` if set).

### Default Global Exclusions

On first use, haide creates a default configuration with common AI coding artifacts:

```ini
[global]
CLAUDE.md
AGENTS.md
.claude/
.clauderc
.cursorrules
.cursorignore
.github/copilot-instructions.md
.cody/
AI_NOTES.md
AI_TODO.md
.ai/
skills/
.skills/
ai-skills/
ai-conversations/
.ai-logs/
prompts/
.prompts/
ai-prompts/
.aiconfig
.ai-config.json
.ai-config.yaml
```

### Project-Specific Overrides

To include a globally-excluded file in a specific project, haide automatically adds override entries (prefixed with `+`):

```ini
[project:myproject]
+CLAUDE.md
my-notes.md
```

### Pattern Matching

haide supports standard glob patterns:

- `*.md` - Match all markdown files
- `.ai/` - Match directory (trailing slash)
- `ai-*` - Match files starting with "ai-"
- `**/*.json` - Recursive directory matching

## How It Works

1. haide manages `.git/info/exclude` (not `.gitignore`) for per-repository exclusions that shouldn't be shared
2. Managed content is delimited by `# BEGIN HAIDE` and `# END HAIDE` markers
3. Content outside markers is preserved
4. Config is the source of truth - changes sync one-way to `.git/info/exclude`
5. Repository sections are keyed by repository basename

## Examples

### Initialize a new repository
```bash
cd my-project
haide init
```

### Add a custom AI file pattern
```bash
haide add "ai-logs/*.txt"
```

### Add a global pattern for all repos
```bash
haide add-global ".anthropic/"
cd other-repo
haide update -y
```

### Include a globally-excluded file in one repo
When you run `haide init` and a file is already tracked but matches a global exclusion, haide automatically adds an override:

```bash
# CLAUDE.md is tracked but globally excluded
haide init
# Automatically adds +CLAUDE.md override to project config
```

### Remove haide from a repository
```bash
haide clean
```

## Limitations

- Unix-like systems only (macOS, Linux)
- One-way sync: config → `.git/info/exclude` only
- Repository basename collisions not supported (same basename = shared config)
- Pattern matching uses standard glob syntax (not gitignore-style patterns)

## License

MIT

## Contributing

Contributions welcome! Please open an issue or pull request.

## Support

For issues and feature requests, please use the GitHub issue tracker:
https://github.com/buckleypaul/haide/issues

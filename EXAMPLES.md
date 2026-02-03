# haide Examples

## Basic Usage

### Initialize in a new repository
```bash
cd my-project
haide init
```

### Add project-specific exclusions
```bash
# Exclude a specific file
haide add AI_NOTES.md

# Exclude a directory
haide add my-ai-logs/

# Exclude with wildcard
haide add "ai-*.txt"
```

### Add global exclusions
```bash
# Add to all repositories
haide add-global ".anthropic/"
haide add-global "*.ai-backup"

# Update other repositories
cd other-project
haide update -y
```

### View configuration
```bash
# Show global and project-specific config
haide info
```

### Remove exclusions
```bash
# Remove project-specific pattern
haide remove "ai-*.txt"

# Remove global pattern
haide remove-global ".anthropic/"
```

### Clean up
```bash
# Remove haide from repository
haide clean
```

## Advanced Scenarios

### Override global exclusions for specific projects

Sometimes you want to track a file in one project that's globally excluded:

```bash
# In a documentation repository where CLAUDE.md should be tracked
cd my-docs-repo

# Create and track CLAUDE.md
echo "# Claude Documentation" > CLAUDE.md
git add CLAUDE.md
git commit -m "Add Claude docs"

# Initialize haide - it will detect the tracked file
haide init
# Prompt: File 'CLAUDE.md' matches exclusion pattern 'CLAUDE.md'
# Automatically adds +CLAUDE.md override to config
```

Or manually add an override:

```bash
# Edit ~/.haide/config.ini
[project:my-docs-repo]
+CLAUDE.md
```

### Pattern matching examples

```bash
# Match all markdown files
haide add "*.md"

# Match specific directory
haide add ".ai-workspace/"

# Match files with prefix
haide add "ai-*"

# Match nested files (recursive)
haide add-global "**/*.ai-backup"

# Match files in subdirectories
haide add ".claude/**"
```

### Dry-run everything first

```bash
# Preview changes before applying
haide init --dry-run
haide add "pattern" --dry-run
haide update --dry-run
haide clean --dry-run
```

### Using custom config location

```bash
# Set HAIDE_HOME environment variable
export HAIDE_HOME=~/my-custom-config
haide init

# Config will be at ~/my-custom-config/config.ini
```

## Workflow Examples

### Team workflow

1. One team member sets up global exclusions:
```bash
haide add-global ".cursor/"
haide add-global ".github/copilot-*"
haide add-global "AI_SCRATCH.md"
```

2. Share the config file:
```bash
# Copy ~/.haide/config.ini to team wiki or shared location
cat ~/.haide/config.ini
```

3. Team members copy the config and initialize:
```bash
# Copy shared config to ~/.haide/config.ini
haide init
```

### Managing multiple repositories

```bash
# Add global pattern
haide add-global "new-pattern"

# Update all repositories
cd ~/projects
for dir in */; do
  (cd "$dir" && [ -d .git ] && haide update -y)
done
```

### Migrating existing repositories

```bash
# For repositories with AI files already tracked
cd existing-repo
haide init
# Responds to prompts about tracked files
# Either untrack or add override

# Batch untrack all AI files
git rm --cached CLAUDE.md AGENTS.md .claude/
git commit -m "Untrack AI files (now managed by haide)"
```

### Project-specific exclusions only

```bash
# Remove global exclusions (edit config)
# Keep only project-specific patterns
[project:myproject]
custom-ai-notes.md
workspace/ai/
```

## Troubleshooting

### Config not found
```bash
# Check config location
haide info

# Verify HAIDE_HOME
echo $HAIDE_HOME

# Recreate config
rm ~/.haide/config.ini
haide init  # Creates default config
```

### Patterns not working
```bash
# Test with dry-run
haide init --dry-run

# Check effective patterns
haide info

# Verify git recognizes exclusions
git status
git check-ignore -v CLAUDE.md
```

### Repository not detected
```bash
# Ensure you're in a git repository
git rev-parse --show-toplevel

# Check if .git exists
ls -la .git
```

### Tracked files still showing
```bash
# Files must be untracked from git
git rm --cached filename

# Or use git status to verify
git status --ignored
```

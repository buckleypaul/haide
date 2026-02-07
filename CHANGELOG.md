# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.3.0] - 2026-02-07

### Added
- Local `.haide` config file for per-project patterns at repository root
- `--remove-local` flag to `clean` command to delete `.haide` file
- Automatic migration from global `[project:name]` sections to local `.haide`
- Local config now appears in `.git/info/exclude` automatically

### Changed
- **Project-specific patterns now stored in `.haide` file instead of global config**
- `add` command writes to local `.haide` file instead of global config
- `remove` command removes from local config (with fallback to legacy global config)
- `init` and `update` commands auto-migrate legacy `[project:name]` sections
- `info` command displays both legacy and local patterns separately
- Per-project configuration is now truly local to each repository

### Migration Note
Existing `[project:name]` sections in `~/.haide/config.ini` will be automatically
migrated to `.haide` files in each repository when running `init` or `update`.
The global config will be cleaned up automatically during migration.

## [0.2.0] - 2026-02-03

### Changed
- **BREAKING**: Config file format changed from TOML to simple INI-style format
- Config file extension: `config.toml` → `config.ini`
- Removed TOML dependency in favor of custom INI parser
- Simplified config syntax: patterns no longer require quotes
- Section format: `[global]` and `[project:name]` instead of TOML tables

### Improved
- More readable config format with one pattern per line
- Cleaner syntax without quotes or commas
- Reduced external dependencies

### Migration Note
Users upgrading from v0.1.0 will need to:
1. Back up `~/.haide/config.toml` if customized
2. Delete old config: `rm ~/.haide/config.toml`
3. Run `haide init` to create new `config.ini` with same patterns

## [0.1.0] - 2026-02-03

### Added
- Initial release of haide
- Core CLI commands: init, add, add-global, remove, remove-global, update, info, clean
- Configuration management with TOML format
- Global and per-project exclusion patterns
- Override support with + prefix
- Pattern matching with wildcards
- Automatic detection of tracked files
- Dry-run mode for all commands
- .git/info/exclude management with markers
- Comprehensive default AI file exclusions
- Integration tests
- GitHub Actions CI/CD
- Homebrew formula support via goreleaser

[Unreleased]: https://github.com/buckleypaul/haide/compare/v0.3.0...HEAD
[0.3.0]: https://github.com/buckleypaul/haide/compare/v0.2.0...v0.3.0
[0.2.0]: https://github.com/buckleypaul/haide/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/buckleypaul/haide/releases/tag/v0.1.0

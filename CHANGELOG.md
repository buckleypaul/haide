# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

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

[Unreleased]: https://github.com/buckleypaul/haide/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/buckleypaul/haide/releases/tag/v0.1.0

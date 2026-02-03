# Contributing to haide

Thank you for your interest in contributing to haide!

## Development Setup

1. Clone the repository:
```bash
git clone https://github.com/buckleypaul/haide.git
cd haide
```

2. Install dependencies:
```bash
go mod download
```

3. Build the project:
```bash
make build
# or
go build -o haide ./cmd/haide
```

4. Run tests:
```bash
make test
# or
go test ./...
```

5. Run integration tests:
```bash
./integration_test.sh
```

## Code Style

- Follow standard Go conventions
- Run `go fmt` before committing
- Run `golangci-lint run` to check for issues
- Write tests for new features
- Update documentation for user-facing changes

## Pull Request Process

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass (`go test ./...`)
6. Run linter (`golangci-lint run`)
7. Commit your changes (`git commit -m 'Add amazing feature'`)
8. Push to your branch (`git push origin feature/amazing-feature`)
9. Open a Pull Request

## Commit Message Guidelines

- Use present tense ("Add feature" not "Added feature")
- Use imperative mood ("Move cursor to..." not "Moves cursor to...")
- Limit first line to 72 characters
- Reference issues and pull requests after the first line

## Testing

- Unit tests: `go test ./...`
- Integration tests: `./integration_test.sh`
- All tests must pass before merging

## Adding New Commands

1. Create command file in `internal/commands/`
2. Register command in `internal/commands/root.go`
3. Add tests
4. Update README.md with command documentation
5. Update EXAMPLES.md with usage examples

## Release Process

Releases are automated via GitHub Actions:

1. Update CHANGELOG.md
2. Create and push a new tag: `git tag v0.x.0 && git push origin v0.x.0`
3. GitHub Actions will build, test, and create a release
4. Homebrew formula will be automatically updated

## Questions?

Feel free to open an issue for:
- Bug reports
- Feature requests
- Documentation improvements
- Questions about contributing

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

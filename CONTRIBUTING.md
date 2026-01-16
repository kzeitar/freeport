# Contributing to freeport

Thank you for your interest in contributing to `freeport`! Contributions are welcome, and we appreciate your help in making this tool better.

## How to Contribute

### Reporting Bugs

Before creating bug reports, please check the existing issues to avoid duplicates. When creating a bug report, include:

- A clear and descriptive title
- Steps to reproduce the issue
- Expected behavior vs. actual behavior
- Your operating system and `freeport` version
- Any relevant error messages or logs

### Suggesting Enhancements

Enhancement suggestions are welcome! Please provide:

- A clear description of the proposed feature
- Use cases and benefits
- Examples of how the feature would work

### Pull Requests

1. **Fork the repository** and create your branch from `main`
2. **Make your changes** with clear, focused commits
3. **Add tests** for new functionality or bug fixes
4. **Ensure all tests pass**: `go test ./...`
5. **Format your code**: `go fmt ./...` or `make fmt`
6. **Update documentation** if needed (README, docs, code comments)
7. **Submit a pull request** with a descriptive title and description

#### Development Setup

```bash
# Clone your fork
git clone https://github.com/kzeitar/freeport.git
cd freeport

# Run tests
make test

# Build locally
make build
```

### Code Style

- Follow standard Go conventions and use `gofmt`
- Write clear, readable code with meaningful variable names
- Add godoc comments to exported functions and types
- Keep changes small and focused on a single issue
- Update tests and documentation as needed

### Testing

- Write tests for new features and bug fixes
- Ensure all existing tests pass before submitting
- Test on multiple platforms if possible (Linux, macOS, Windows)
- Consider edge cases and error conditions

## Getting Help

If you need help or have questions:

- Open an issue for bugs or feature requests
- Check existing issues and discussions for answers
- Read the [documentation](docs/usage.md) for usage guidance

## License

By contributing, you agree that your contributions will be licensed under the [MIT License](LICENSE).

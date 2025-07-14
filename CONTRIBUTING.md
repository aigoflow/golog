# Contributing to GoLog

Thank you for your interest in contributing to GoLog! We welcome contributions from the community to help make this project better.

## Code of Conduct

By participating in this project, you agree to abide by our Code of Conduct:
- Be respectful and inclusive
- Welcome newcomers and help them get started
- Focus on constructive criticism
- Respect differing viewpoints and experiences

## How to Contribute

### Reporting Bugs

1. Check if the bug has already been reported in [Issues](https://github.com/yourusername/golog/issues)
2. If not, create a new issue with:
   - Clear title and description
   - Steps to reproduce
   - Expected vs actual behavior
   - System information (OS, Go version)
   - Relevant logs or error messages

### Suggesting Features

1. Check existing issues and discussions
2. Open a new issue with the "enhancement" label
3. Describe the feature and its use case
4. Explain how it benefits LLM integration or Prolog learning

### Pull Requests

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass (`make test`)
6. Commit with clear messages (`git commit -m 'Add amazing feature'`)
7. Push to your branch (`git push origin feature/amazing-feature`)
8. Open a Pull Request

### Development Setup

```bash
# Clone your fork
git clone https://github.com/yourusername/golog.git
cd golog

# Install dependencies
make deps

# Run tests
make test

# Build
make build

# Run locally
make dev
```

### Coding Standards

- Follow Go standard formatting (`go fmt`)
- Write clear, self-documenting code
- Add comments for complex logic
- Keep functions small and focused
- Write unit tests for new features

### Testing

- Write tests for all new functionality
- Ensure existing tests pass
- Add integration tests for API changes
- Test both UI and API components

### Documentation

- Update README.md for new features
- Add API documentation for new endpoints
- Include examples for LLM integration
- Update tutorial for new Prolog features

## Release Process

1. Update version numbers
2. Update CHANGELOG.md
3. Create a pull request
4. After merge, tag the release
5. GitHub Actions will build and publish

## Need Help?

- Join our [Discussions](https://github.com/yourusername/golog/discussions)
- Check the [Wiki](https://github.com/yourusername/golog/wiki)
- Ask questions in issues with the "question" label

Thank you for contributing! 🎉
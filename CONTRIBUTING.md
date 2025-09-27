# Contributing to Green Orb

Thank you for your interest in contributing to Green Orb! This document provides guidelines and information for developers who want to contribute to the project.

## Table of Contents

- [Development Setup](#development-setup)
- [Project Structure](#project-structure)
- [Building the Project](#building-the-project)
- [Running Tests](#running-tests)
- [Code Style](#code-style)
- [Making Changes](#making-changes)
- [Submitting Pull Requests](#submitting-pull-requests)

## Development Setup

### Prerequisites

- Go 1.25 or later
- Git
- Make (optional, for convenience)

### Getting Started

1. Fork the repository on GitHub
2. Clone your fork:
   ```bash
   git clone https://github.com/YOUR-USERNAME/green-orb.git
   cd green-orb
   ```
3. Add the upstream remote:
   ```bash
   git remote add upstream https://github.com/atgreen/green-orb.git
   ```
4. Install dependencies:
   ```bash
   go mod download
   ```

## Project Structure

The Green Orb codebase is organized into several modules for better maintainability:

```
green-orb/
├── green-orb.go        # Main entry point and CLI
├── config.go           # Configuration loading and validation
├── worker.go           # Action processing workers
├── monitor.go          # Output monitoring from observed process
├── kafka.go            # Kafka client management
├── checks.go           # Health check scheduler
├── metrics.go          # Prometheus metrics
├── signals_unix.go     # Unix signal handling
├── signals_windows.go  # Windows signal handling
├── *_test.go          # Test files
└── .github/           # CI/CD workflows
```

### Key Components

- **Config Module**: Handles YAML configuration parsing and validation
- **Worker Pool**: Manages concurrent action processing with rate limiting
- **Monitor**: Processes stdout/stderr from the observed process
- **Kafka Manager**: Manages Kafka client connections with TLS/SASL support
- **Check Scheduler**: Runs periodic health checks (HTTP, TCP, flapping detection)
- **Metrics**: Exposes Prometheus metrics for monitoring

## Building the Project

### Basic Build

```bash
go build -o orb
```

### Cross-Platform Builds

```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o orb-linux-amd64

# Windows
GOOS=windows GOARCH=amd64 go build -o orb-windows-amd64.exe

# macOS
GOOS=darwin GOARCH=amd64 go build -o orb-darwin-amd64
```

### Build with Version Information

```bash
go build -ldflags "-X main.version=$(git describe --tags --always)" -o orb
```

## Running Tests

### Unit Tests

Run all unit tests:
```bash
go test ./...
```

Run tests with coverage:
```bash
go test -cover ./...
```

Generate coverage report:
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

### Benchmarks

Run benchmarks:
```bash
go test -bench=. -benchmem ./...
```

### End-to-End Tests

The E2E test is run in CI but can be executed locally:
```bash
go build -o orb
./scripts/e2e-test.sh
```

### Linting

We use golangci-lint for code quality checks:
```bash
golangci-lint run --timeout=5m
```

## Code Style

### Go Standards

- Follow the official [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Use `gofmt` for formatting
- Use meaningful variable and function names
- Keep functions small and focused
- Handle errors explicitly

### Project Conventions

- **Error Messages**: Start with "green-orb error:" or "green-orb warning:"
- **Metrics**: Use descriptive names with appropriate labels
- **Configuration**: Validate all configuration at startup
- **Concurrency**: Use channels for communication, mutexes for shared state
- **Testing**: Write tests for new functionality, aim for 80%+ coverage

### Commit Messages

- Use clear, descriptive commit messages
- Start with a capital letter
- Use imperative mood ("Add feature" not "Added feature")
- Reference issues when applicable (#123)

Example:
```
Add Kafka SASL/PLAIN authentication support

- Implement SASL mechanism configuration
- Add username/password fields to channel config
- Update documentation with authentication examples

Fixes #42
```

## Making Changes

### Development Workflow

1. Create a feature branch:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. Make your changes, following the code style guidelines

3. Add or update tests as needed

4. Run tests and linting:
   ```bash
   go test ./...
   golangci-lint run
   ```

5. Commit your changes with a descriptive message

6. Push to your fork:
   ```bash
   git push origin feature/your-feature-name
   ```

### Adding New Features

When adding new features:

1. **Update Configuration**: Add new fields to `Config` struct if needed
2. **Add Validation**: Ensure configuration validation in `validateConfig()`
3. **Write Tests**: Add unit tests for new functionality
4. **Update Documentation**: Update README.md and add examples
5. **Add Metrics**: Include appropriate Prometheus metrics
6. **Consider Backwards Compatibility**: Ensure existing configs still work

### Adding New Channel Types

To add a new channel type:

1. Update `Channel` struct in `config.go` if new fields needed
2. Add validation in `validateConfig()`
3. Implement execution in `worker.go`'s `processAction()`
4. Add tests in `green-orb_test.go`
5. Update README.md with documentation

## Submitting Pull Requests

### Before Submitting

- [ ] Code follows project style guidelines
- [ ] All tests pass
- [ ] Linting passes without errors
- [ ] Documentation is updated
- [ ] Commit messages are clear
- [ ] Branch is up to date with main

### Pull Request Process

1. Create a pull request from your feature branch to `main`
2. Fill out the PR template with:
   - Description of changes
   - Testing performed
   - Documentation updates
3. Wait for CI checks to pass
4. Address review feedback
5. Maintainers will merge once approved

### PR Guidelines

- Keep PRs focused on a single feature or fix
- Include tests for new functionality
- Update documentation as part of the PR
- Respond to feedback promptly
- Be patient and respectful

## Testing Guidelines

### Unit Test Structure

```go
func TestFunctionName(t *testing.T) {
    tests := []struct {
        name    string
        input   InputType
        want    OutputType
        wantErr bool
    }{
        {
            name:    "valid input",
            input:   validInput,
            want:    expectedOutput,
            wantErr: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := FunctionName(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("FunctionName() error = %v, wantErr %v", err, tt.wantErr)
            }
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("FunctionName() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

### Integration Tests

Integration tests should:
- Test complete workflows
- Use real configurations
- Verify side effects
- Clean up resources

## Getting Help

- Open an issue for bugs or feature requests
- Start a discussion for questions
- Check existing issues before creating new ones
- Join community discussions

## License

By contributing to Green Orb, you agree that your contributions will be licensed under the MIT License.

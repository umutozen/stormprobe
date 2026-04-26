# Contributing to StormProbe

Thank you for considering contributing to StormProbe.

## Getting Started

```bash
git clone https://github.com/umutozen/stormprobe.git
cd stormprobe
go build ./...
go test ./...
```

## Development Requirements

- Go 1.21+
- katana and httpx in PATH (optional, for discovery testing)

## Making Changes

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/your-feature`
3. Make your changes
4. Run tests: `go test ./...`
5. Run linter: `golangci-lint run ./...`
6. Commit with a clear message following [Conventional Commits](https://www.conventionalcommits.org/)
7. Push and open a Pull Request

## Commit Message Format

```
type: short description

types: feat, fix, docs, test, ci, refactor, perf
```

Examples:
```
feat: add --timeout flag for per-request timeout
fix: detect Edge from Program Files (x86) on Windows
docs: add Portuguese translation
```

## Code Standards

- All identifiers must be in English
- Functions must not exceed 40 lines
- No magic numbers -- use named constants
- All exported errors must be handled
- New features require unit tests

## Reporting Issues

Use the GitHub issue templates for bug reports and feature requests.

## Legal

By contributing, you agree your contributions will be licensed under the MIT License.

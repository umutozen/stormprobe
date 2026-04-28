# Contributing to StormProbe

StormProbe is a verdict-first HTTP load testing engine. Every contribution should move the tool closer to better operational decisions — not just more numbers.

> **Golden Rule:** If it doesn't improve the final verdict, it's not a priority.

---

## Project Architecture

```
stormprobe/
├── cmd/stormprobe/       # CLI entry point & flag parsing
├── internal/
│   ├── alert/            # CI/CD threshold evaluation (P99, error rate, RPS)
│   ├── compare/          # Before/after JSON report comparison
│   ├── config/           # Flag parsing, validation, defaults
│   ├── discovery/        # Katana crawl + Httpx probe orchestration
│   ├── metrics/          # Latency percentiles, throughput, error classification
│   ├── report/           # JSON + HTML report generation
│   ├── runner/           # 4-phase load executor (ramp/sustained/spike/recovery)
│   └── verdict/          # Safe concurrency, degradation point, bottleneck classification
├── scripts/              # i18n build tooling (maintainers only)
├── docs/i18n/            # Auto-generated localized README files (16 languages)
├── assets/               # Banner and static assets
└── outputs/              # Test output directory (gitignored)
```

Each `internal/` package has a single responsibility. Keep it that way.

---

## Development Setup

```bash
git clone https://github.com/umutozen/stormprobe.git
cd stormprobe
go build -o stormprobe.exe ./cmd/stormprobe
go test ./...
```

### Requirements

| Tool | Purpose |
|---|---|
| Go 1.21+ | Build and test |
| golangci-lint | Linting (`golangci-lint run ./...`) |
| Katana | Optional — endpoint crawling |
| Httpx | Optional — endpoint probing |

### Makefile Commands

```bash
make build    # Build binary
make test     # Run go vet + go test
make lint     # Run golangci-lint
make clean    # Remove binary + outputs/
```

---

## Making Changes

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/your-feature`
3. Make your changes
4. Verify:
   ```bash
   go test ./...
   go vet ./...
   golangci-lint run ./...
   ```
5. Commit with a clear message (see below)
6. Push and open a Pull Request

---

## Code Standards

### Hard Rules

- **No comments in source code** — neither Turkish nor English
- **Functions ≤ 40 lines** — split if longer
- **No magic numbers** — use named constants
- **All errors must be handled** — no silent failures
- **New features require unit tests**

### Architecture Rules

- **Single Responsibility** — one package, one job
- UI logic and business logic must be separated
- API calls go through service layers, never inside components
- No unnecessary global state — prefer local state
- Derived state is computed, not stored

### Linter Configuration

The project uses `golangci-lint` with these enabled linters:

`errcheck` · `gosimple` · `govet` · `ineffassign` · `staticcheck` · `unused` · `gofmt` · `goimports` · `misspell` · `unconvert` · `unparam`

All PRs must pass `golangci-lint run ./...` with zero warnings.

---

## Commit Message Format

Follow [Conventional Commits](https://www.conventionalcommits.org/):

```
type: short description
```

| Type | Use |
|---|---|
| `feat` | New feature or flag |
| `fix` | Bug fix |
| `docs` | Documentation changes |
| `test` | Adding or modifying tests |
| `ci` | CI/CD configuration |
| `refactor` | Code restructuring without behavior change |
| `perf` | Performance improvement |

**Examples:**

```
feat: add --timeout flag for per-request timeout
fix: correct throughput plateau calculation in recovery phase
docs: update TR and ES i18n README files
refactor: extract saturation logic from runner to metrics
```

---

## When Adding a New Feature or Flag

1. Implement the feature in the appropriate `internal/` package
2. Wire the flag in `cmd/stormprobe/main.go` and `internal/config/`
3. Update `README.md` — add a row to the Flags table
4. Update `RELEASE_NOTES.md` if it's a user-facing change
5. Update `config.example.json` if a new config field is introduced
6. Write unit tests covering the new behavior

> **Note:** Documentation translations (i18n) are maintained internally. Contributors do not need to run translation scripts.

---

## Pull Request Checklist

Before submitting, verify:

- [ ] `go test ./...` passes
- [ ] `go vet ./...` passes
- [ ] `golangci-lint run ./...` passes
- [ ] Functions are ≤ 40 lines
- [ ] No comments in source code
- [ ] No hardcoded values — constants are defined
- [ ] New features include unit tests
- [ ] `CHANGELOG.md` updated under `[Unreleased]`
- [ ] No sensitive data or credentials included

---

## Docker Testing

```bash
docker build -t stormprobe:dev .
docker run --rm stormprobe:dev --no-discovery --insecure https://example.com
```

The Docker image bundles Katana and Httpx automatically.

---

## Reporting Issues

Use the GitHub issue templates:

- **Bug Report** — for unexpected behavior, crashes, or incorrect verdicts
- **Feature Request** — for new capabilities that improve decision quality

Include reproduction steps, expected vs actual output, and OS/Go version.

---

## Security

Found a vulnerability? Do **not** open a public issue.

Follow the responsible disclosure process in [SECURITY.md](SECURITY.md).

---

## License

By contributing, you agree that your contributions will be licensed under the [MIT License](LICENSE).

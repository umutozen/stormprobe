# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [Unreleased]

---

## [1.0.0] - 2026-04-26

### Initial Release — StormProbe v1.0.0

StormProbe is a zero-dependency, production-grade HTTP load testing CLI tool with
autonomous endpoint discovery powered by Katana and Httpx.

### Added

#### Core Load Testing Engine
- 4-phase load test pipeline: **Ramp-Up → Sustained → Spike → Recovery**
- Configurable concurrency per phase (`--concurrency-ramp`, `--concurrency-sustained`, `--concurrency-spike`)
- Per-request timeout control (`--timeout`)
- Requests-per-worker tuning (`--req-per-worker`)
- Warm-up phase (3 requests) before test execution begins
- Real-time ASCII table output during each phase with: Concurrency, Success%, Avg/P50/P95/P99 latencies, req/s
- Detailed error analysis summary: Timeout, Connection Reset/Refused, HTTP 4xx/5xx, Other
- HTTP status code distribution report at test completion

#### Autonomous Endpoint Discovery
- Katana static + JS/headless crawl integration for deep endpoint discovery
- Httpx probing to filter and validate live endpoints
- Windows Chrome/Edge headless browser auto-detection from common install paths
- `--no-discovery` flag to skip crawl and test root path (`/`) only
- `--endpoints` flag to supply a custom endpoint list from file (one path per line, `#` comment support)
- Custom binary paths via `--katana-path` and `--httpx-path`

#### HTTP Client & Headers
- `--header` / `-H` flag for custom HTTP headers (repeatable, e.g. `-H "Authorization: Bearer TOKEN"`)
- Custom headers propagated to: warmup requests, all load test phases, Katana discovery, Httpx probing
- `--insecure` flag for TLS certificate verification bypass
- High-performance HTTP transport: 600 max idle connections, 30s idle timeout

#### Reporting
- Self-contained HTML report with embedded SVG latency chart (Ramp-Up/Sustained/Spike/Recovery per scenario)
- Error distribution bar chart in HTML report
- JSON report output with full per-phase result data
- `--format` flag: `json`, `html`, `both` (default: `both`)
- Configurable output directory (`--output`, default: `./outputs`)

#### Developer & Operations
- Zero external Go dependencies (`go.mod` stdlib-only)
- Cross-platform binary distribution via GoReleaser: Linux, macOS, Windows — amd64/arm64
- Docker image published to GHCR (`ghcr.io/umutozen/stormprobe`) on `main` push and release tags
- Multi-stage Dockerfile with minimal final image
- CI/CD via GitHub Actions: test, lint, Docker build, GoReleaser release
- Unit tests for `metrics` and `report` packages
- `.golangci.yml` linter configuration (errcheck, gofmt)
- `golangci-lint` errcheck fixes: `io.Copy` and `cmd.Process.Kill()` return values handled

#### Documentation & Community
- 16-language README: EN, TR, DE, FR, ES, IT, PT-BR, RU, ZH-CN, ZH-TW, JA, KO, AR, HI, NL, PL, UK
- `CONTRIBUTING.md` — contribution guidelines
- `SECURITY.md` — vulnerability reporting policy
- `LICENSE` — MIT license
- `config.example.yml` — annotated configuration reference
- GitHub issue templates: bug report, feature request
- GitHub PR template
- `FUNDING.yml` — GitHub Sponsors / Open Collective configuration
- `Makefile` — build, test, lint shortcuts

### Fixed
- SVG chart label truncation and bottom padding in HTML report
- Success badge thresholds in table output
- Turkish characters in language switcher and disclaimer sections of README
- Chrome/Edge headless browser detection on Windows when binaries are not on `PATH`

---

## Links

[Unreleased]: https://github.com/umutozen/stormprobe/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/umutozen/stormprobe/releases/tag/v1.0.0

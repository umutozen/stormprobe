# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [Unreleased]

### Added

- `--duration` flag: per-phase duration mode (`30s`, `1m`) using `context.WithTimeout`. Overrides `--req-per-worker` when set
- `--alert-p99` flag: exit code 1 if P99 latency exceeds threshold in any phase
- `--alert-error-rate` flag: exit code 1 if error rate exceeds threshold (%) in any phase
- `--alert-rps` flag: exit code 1 if req/s falls below threshold in any phase
- `internal/alert` package: threshold check logic, CI/CD integration

---

## [1.0.0] - 2026-04-26

### Added

**Core Load Testing Engine**
- 4-phase load test pipeline: Ramp-Up, Sustained, Spike, Recovery
- Configurable concurrency per phase (`--concurrency-ramp`, `--concurrency-sustained`, `--concurrency-spike`)
- Per-request timeout control (`--timeout`)
- Requests-per-worker tuning (`--req-per-worker`)
- Warm-up phase (3 requests) executed before test begins
- Real-time ASCII table output per phase: Concurrency, Success%, Avg/P50/P95/P99 latencies, req/s
- Detailed error analysis summary: Timeout, Connection Reset/Refused, HTTP 4xx/5xx, Other
- HTTP status code distribution report at test completion

**Autonomous Endpoint Discovery**
- Katana static and JS/headless crawl integration
- Httpx probing to filter and validate live endpoints
- `--no-discovery` flag to skip crawl and test root path only
- `--endpoints` flag to load a custom endpoint list from file (one path per line, `#` comment support)
- Custom binary paths via `--katana-path` and `--httpx-path`
- Windows Chrome/Edge headless browser auto-detection from common install paths

**HTTP Client and Headers**
- `--header` / `-H` flag for custom HTTP headers (repeatable)
- Custom headers propagated to warmup requests, all load test phases, Katana, and Httpx
- `--insecure` flag for TLS certificate verification bypass
- High-performance HTTP transport: 600 max idle connections, 30s idle timeout

**Reporting**
- Self-contained HTML report with embedded SVG latency chart
- Error distribution visualization in HTML report
- JSON report output with full per-phase result data
- `--format` flag: `json`, `html`, `both` (default: `both`)
- Configurable output directory via `--output`

**Distribution and CI/CD**
- Zero external Go dependencies (stdlib only)
- Cross-platform binaries via GoReleaser: Linux, macOS, Windows — amd64 and arm64
- Docker image published to GHCR (`ghcr.io/umutozen/stormprobe`)
- Multi-stage Dockerfile with minimal final image
- GitHub Actions CI/CD: test, lint, Docker build, GoReleaser release
- Unit tests for `metrics` and `report` packages
- `.golangci.yml` linter configuration

**Documentation and Community**
- 16-language README: EN, TR, DE, FR, ES, IT, PT-BR, RU, ZH-CN, ZH-TW, JA, KO, AR, HI, NL, PL, UK
- `CONTRIBUTING.md`, `SECURITY.md`, `LICENSE` (MIT)
- `config.example.yml` annotated configuration reference
- GitHub issue templates (bug report, feature request) and PR template
- `FUNDING.yml` for GitHub Sponsors / Open Collective

---

[Unreleased]: https://github.com/umutozen/stormprobe/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/umutozen/stormprobe/releases/tag/v1.0.0

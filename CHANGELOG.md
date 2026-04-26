# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.0.1] - 2026-04-26

### Added
- `--header` / `-H` flag for custom HTTP headers (repeatable, e.g. `-H "Authorization: Bearer TOKEN"`)
- Windows Chrome/Edge headless browser auto-detection from common install paths
- Docker image published to GHCR on every `main` branch push
- Standard project files: LICENSE (MIT), CONTRIBUTING.md, SECURITY.md, CHANGELOG.md
- GitHub issue templates (bug report, feature request) and PR template
- `.golangci.yml` linter configuration

### Fixed
- Custom headers now applied to warmup requests
- Katana and httpx discovery processes receive custom headers via `-H`
- gofmt spacing fix in `headerFlag.String()` method
- `golangci-lint` errcheck warnings for `io.Copy` and `cmd.Process.Kill()`
- SVG chart label truncation and bottom padding in HTML report

## [1.0.0] - 2026-04-26

### Added
- Autonomous endpoint discovery via Katana (static + JS/headless crawl) and Httpx
- 4-phase load test: Ramp-Up, Sustained, Spike, Recovery
- Self-contained HTML report with SVG latency chart and error distribution bars
- JSON report output
- `--insecure` flag for TLS certificate verification bypass
- `--no-discovery` flag to skip crawl and test root path only
- `--endpoints` flag to supply a custom endpoint list
- `--format` flag: json, html, both
- Docker image published to GHCR (`ghcr.io/umutozen/stormprobe`)
- Cross-platform binaries via GoReleaser (Linux, macOS, Windows -- amd64/arm64)
- 16-language README (EN, TR, DE, FR, ES, IT, PT-BR, RU, ZH-CN, ZH-TW, JA, KO, AR, HI, NL, PL, UK)
- CI/CD via GitHub Actions (test, lint, docker, release)
- Windows Chrome/Edge headless browser auto-detection from common install paths

[Unreleased]: https://github.com/umutozen/stormprobe/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/umutozen/stormprobe/releases/tag/v1.0.0

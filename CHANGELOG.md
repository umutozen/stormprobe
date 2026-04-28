# Changelog

All notable changes to StormProbe are documented in this file.

Format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).

---

## [1.0.0] - 2026-04-26

**Initial Release**

StormProbe is a zero-dependency, production-grade HTTP load testing engine written in pure Go. It autonomously discovers endpoints, runs phased stress tests, detects saturation patterns, and produces actionable production-readiness verdicts.

### Verdict Engine

4-branch classification (`Healthy` / `Degraded` / `Unstable` / `Critical`) with confidence qualifiers. Includes degradation point detection, throughput saturation analysis, and multi-metric recovery validation.

### Compare Mode

`stormprobe compare before.json after.json` — a verdict delta engine that classifies changes as `IMPROVED` / `STABLE` / `STABLE WITH RISK` / `DEGRADED` / `CRITICAL REGRESSION` with a "Why" explanation.

### 4-Phase Load Pipeline

Discovery → Ramp-Up → Sustained → Spike → Recovery — each phase targets a specific operational question about the system under test.

### Autonomous Discovery

Katana (static + headless) crawl combined with Httpx probe validation. Smart endpoint scoring filters static assets and ranks paths by test value.

### CI/CD Integration

`--alert-p99`, `--alert-error-rate`, `--alert-rps` flags produce exit code 1 when thresholds are breached. Duration-based phases via `--duration` for consistent benchmarks.

### Reports

Self-contained HTML dashboard with executive summary, SVG latency chart, and priority action list. JSON output with `report_id`, `confidence`, and null semantics for clean downstream parsing.

### CLI

Custom headers (`-H`), TLS bypass (`--insecure`), format selection (`--format`), custom endpoint lists (`--endpoints`), JSON config file (`--config`), response body validation (`--success-match`, `--fail-match`), and full concurrency tuning.

### Multi-Target Mode

Define multiple targets in a single JSON config, run tests sequentially with per-target summary table and automatic verdict comparison.

### Live Progress Bar

Real-time progress indicator showing completion percentage, request count, req/s, and elapsed time during test execution.

### Infrastructure

Cross-platform binaries (Linux, macOS, Windows — amd64 & arm64), Docker image with bundled Katana/Httpx, GoReleaser automation, GitHub Actions CI/CD, and 16-language documentation.

---

[1.0.0]: https://github.com/umutozen/stormprobe/releases/tag/v1.0.0

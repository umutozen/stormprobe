[English](../../README.md) · [Türkçe](docs/i18n/README.tr.md) · [Deutsch](docs/i18n/README.de.md) · [Francais](docs/i18n/README.fr.md) · [Espanol](docs/i18n/README.es.md) · [Italiano](docs/i18n/README.it.md) · [Portugues](docs/i18n/README.pt-BR.md) · [Russkij](docs/i18n/README.ru.md) · [Zhongwen Jian](docs/i18n/README.zh-CN.md) · [Zhongwen Fan](docs/i18n/README.zh-TW.md) · [Nihongo](docs/i18n/README.ja.md) · [Hangugeo](docs/i18n/README.ko.md) · [Al-Arabiyya](docs/i18n/README.ar.md) · [Hindi](docs/i18n/README.hi.md) · [Nederlands](docs/i18n/README.nl.md) · [Polski](docs/i18n/README.pl.md) · [Ukrainska](docs/i18n/README.uk.md)

<p align="center">
  <img src="assets/banner.png" alt="StormProbe Banner" width="100%">
</p>

<h1 align="center">StormProbe</h1>

<p align="center">
  <strong>discover · load · analyze</strong><br>
  Production-grade HTTP load tester with autonomous endpoint discovery. Zero external dependencies.
</p>

<p align="center">
  <a href="https://github.com/umutozen/stormprobe/actions/workflows/ci.yml">
    <img src="https://github.com/umutozen/stormprobe/actions/workflows/ci.yml/badge.svg" alt="CI">
  </a>
  <a href="https://goreportcard.com/report/github.com/umutozen/stormprobe">
    <img src="https://goreportcard.com/badge/github.com/umutozen/stormprobe" alt="Go Report Card">
  </a>
  <a href="https://github.com/umutozen/stormprobe/releases">
    <img src="https://img.shields.io/github/v/release/umutozen/stormprobe" alt="Latest Release">
  </a>
  <a href="LICENSE">
    <img src="https://img.shields.io/badge/license-MIT-blue.svg" alt="MIT License">
  </a>
  <a href="https://pkg.go.dev/github.com/umutozen/stormprobe">
    <img src="https://pkg.go.dev/badge/github.com/umutozen/stormprobe.svg" alt="Go Reference">
  </a>
</p>

---

## What is StormProbe?

StormProbe is a **zero-dependency** HTTP load testing tool written in pure Go. It automatically discovers your application's endpoints using [Katana](https://github.com/projectdiscovery/katana) and [Httpx](https://github.com/projectdiscovery/httpx), then runs a structured 4-phase load test — **Ramp-Up → Sustained → Spike → Recovery** — and produces JSON and HTML reports with P50/P95/P99 latency metrics.

Designed for:
- **CI/CD pipelines** — alert thresholds with exit codes
- **Penetration testers** — find live endpoints automatically, stress them
- **DevOps engineers** — benchmark before & after deployments
- **QA teams** — validate performance SLAs against real traffic patterns

---

## Features

| Feature | Details |
|---|---|
| **Auto-Discovery** | Katana crawl + Httpx probe, deduplicated live endpoint list |
| **4-Phase Load Test** | Ramp-Up → Sustained → Spike → Recovery |
| **Rich Metrics** | P50 / P95 / P99 latency, req/s, error classification per phase |
| **Duration Mode** | Time-based phases (`--duration 30s`) instead of request count |
| **Alert Thresholds** | CI/CD-ready `exit 1` on P99, error rate, or RPS breach |
| **Dual Reports** | JSON (machine-readable) + self-contained HTML dark dashboard |
| **Custom Headers** | Bearer tokens, tenant headers, cookies — all propagated everywhere |
| **Docker Ready** | Single image with bundled Katana + Httpx |
| **Cross-Platform** | Linux, macOS, Windows — amd64 & arm64 |
| **Zero Dependencies** | Pure Go 1.21+ stdlib, `go install` and done |

---

## Quick Start

```bash
# Install (Go 1.21+ required)
go install github.com/umutozen/stormprobe/cmd/stormprobe@latest

# Basic test — auto-discovers endpoints, all 4 phases
stormprobe https://example.com

# Skip discovery, just hammer the root path
stormprobe --no-discovery https://example.com

# TLS + custom Auth header + JSON-only report
stormprobe --insecure -H "Authorization: Bearer TOKEN" --format json https://api.example.com

# Duration-based test: 30 seconds per phase
stormprobe --insecure --duration 30s --no-discovery https://example.com

# CI/CD: fail if P99 > 500ms or error rate > 5%
stormprobe --alert-p99 500 --alert-error-rate 5 https://example.com
```

---

## Installation

### Option 1 — `go install` *(recommended)*

Requires Go 1.21+. Installs the binary directly to `$GOPATH/bin`.

```bash
go install github.com/umutozen/stormprobe/cmd/stormprobe@latest
stormprobe --help
```

### Option 2 — Pre-built Binary *(no Go required)*

Download from the [Releases page](https://github.com/umutozen/stormprobe/releases).

```bash
# Linux / macOS
chmod +x stormprobe
sudo mv stormprobe /usr/local/bin/

# Windows PowerShell
Move-Item stormprobe.exe C:\Windows\System32\stormprobe.exe
```

### Option 3 — Build from Source

```bash
git clone https://github.com/umutozen/stormprobe.git
cd stormprobe
go build -o stormprobe ./cmd/stormprobe      # Linux / macOS
go build -o stormprobe.exe ./cmd/stormprobe  # Windows
```

### Option 4 — Docker

```bash
# Minimal test, no saved report
docker run --rm ghcr.io/umutozen/stormprobe:latest --insecure https://example.com

# Save HTML + JSON reports to ./outputs
# Linux / macOS
docker run --rm -v $(pwd)/outputs:/app/outputs \
  ghcr.io/umutozen/stormprobe:latest --insecure --format both https://example.com

# Windows PowerShell
docker run --rm -v "${PWD}\outputs:/app/outputs" `
  ghcr.io/umutozen/stormprobe:latest --insecure --format both https://example.com
```

### Optional — Endpoint Discovery Tools

Auto-discovery requires Katana and Httpx. If not installed, use `--no-discovery` or provide `--endpoints`.

```bash
# Install via script
bash scripts/install-tools.sh

# Or manually
go install github.com/projectdiscovery/katana/cmd/katana@latest
go install github.com/projectdiscovery/httpx/cmd/httpx@latest
```

---

## CLI Reference

```
stormprobe [flags] <target-url>
```

### Concurrency & Load

| Flag | Default | Description |
|---|---|---|
| `--concurrency-ramp` | `50` | Peak concurrency for ramp-up phase |
| `--concurrency-sustained` | `75` | Concurrency for sustained phase |
| `--concurrency-spike` | `250` | Peak concurrency for spike phase |
| `--req-per-worker` | `15` | Requests per worker per step (ignored if `--duration` is set) |
| `--duration` | `0` (disabled) | Time-based phase duration, e.g. `30s`, `1m` |
| `--timeout` | `10s` | Per-request timeout |

### Discovery

| Flag | Default | Description |
|---|---|---|
| `--no-discovery` | `false` | Skip Katana+Httpx, test root path only |
| `--endpoints` | `""` | Load endpoints from a file (one path per line, `#` comments ok) |
| `--katana-path` | `""` | Custom Katana binary path |
| `--httpx-path` | `""` | Custom Httpx binary path |

### HTTP & Security

| Flag | Default | Description |
|---|---|---|
| `-H`, `--header` | — | Custom header (repeatable): `-H "Authorization: Bearer TOKEN"` |
| `--insecure` | `false` | Skip TLS certificate verification |

### Output

| Flag | Default | Description |
|---|---|---|
| `--format` | `both` | Report format: `json`, `html`, `both` |
| `--output` | `./outputs` | Output directory for reports |

### CI/CD Alerts

| Flag | Default | Description |
|---|---|---|
| `--alert-p99` | `0` (disabled) | Exit 1 if P99 latency (ms) exceeds threshold in any phase |
| `--alert-error-rate` | `0` (disabled) | Exit 1 if error rate (%) exceeds threshold in any phase |
| `--alert-rps` | `0` (disabled) | Exit 1 if req/s falls below threshold in any phase |

---

## Usage Examples

```bash
# Auto-discover + full test
stormprobe https://example.com

# Skip discovery (faster, useful for known targets)
stormprobe --no-discovery https://example.com

# Custom endpoint list (one URL path per line)
stormprobe --endpoints endpoints.txt https://example.com

# Authenticated API stress test
stormprobe --insecure -H "Authorization: Bearer TOKEN" -H "X-Tenant: acme" https://api.example.com

# High-concurrency spike (500 virtual users)
stormprobe --insecure --concurrency-spike 500 --req-per-worker 20 https://example.com

# Duration-based: each phase runs for 1 minute
stormprobe --insecure --duration 1m --concurrency-spike 300 https://example.com

# JSON-only report saved to custom directory
stormprobe --insecure --format json --output ./results https://example.com

# CI/CD: exit 1 if P99 > 500ms or error% > 5% or RPS < 100
stormprobe --alert-p99 500 --alert-error-rate 5 --alert-rps 100 https://example.com
```

### CI/CD Integration (GitHub Actions)

```yaml
- name: Performance gate
  run: |
    go install github.com/umutozen/stormprobe/cmd/stormprobe@latest
    stormprobe \
      --no-discovery \
      --duration 30s \
      --alert-p99 500 \
      --alert-error-rate 2 \
      --alert-rps 50 \
      https://staging.example.com
```

---

## Test Phases

```
Discovery ──► Ramp-Up ──► Sustained ──► Spike ──► Recovery
```

| Phase | Description |
|---|---|
| **0 — Discovery** | Katana crawl + Httpx probe, deduplicated endpoint list |
| **1 — Ramp-Up** | Gradual concurrency increase: 5 → 15 → 30 → peak users |
| **2 — Sustained** | 3 waves at target concurrency, measures steady-state degradation |
| **3 — Spike** | Sudden burst to peak concurrency, then cooldown |
| **4 — Recovery** | Post-spike health check after 10s cooldown |

---

## Output & Reports

Reports are saved to `./outputs/` (configurable via `--output`).

| File | Description |
|---|---|
| `stormprobe_report_<timestamp>.json` | All metrics per phase, machine-readable |
| `stormprobe_report_<timestamp>.html` | Self-contained dark-theme dashboard with SVG latency chart |

**Metrics captured per phase:**
- `success%` — successful request ratio
- `avg / P50 / P95 / P99` — latency percentiles (ms)
- `req/s` — throughput
- Error breakdown: Timeout, Connection Reset, Connection Refused, HTTP 4xx, HTTP 5xx, Other
- HTTP status code distribution

---

## Project Structure

```
stormprobe/
├── cmd/stormprobe/
│   └── main.go                  CLI entry point, flag parsing, phase orchestration
├── internal/
│   ├── config/
│   │   └── config.go            Shared types, phase generators, AlertConfig
│   ├── alert/
│   │   └── alert.go             Threshold checks, CI/CD exit code logic
│   ├── metrics/
│   │   ├── latency.go           Percentile calculation (P50/P95/P99)
│   │   └── errors.go            Error classification
│   ├── discovery/
│   │   ├── discovery.go         Discover() public interface
│   │   ├── katana.go            Katana crawler integration
│   │   └── httpx.go             Httpx probe integration
│   ├── runner/
│   │   ├── phase.go             Phase orchestration, duration/count modes
│   │   └── worker.go            Goroutine pool, per-worker HTTP execution
│   └── report/
│       ├── json.go              JSON report writer
│       └── html.go              Self-contained HTML dark dashboard
├── assets/
│   └── banner.png
├── scripts/
│   └── install-tools.sh         Katana + Httpx installer
├── Dockerfile                   Multi-stage: Go build + katana + httpx
├── .goreleaser.yml              Cross-platform release config
├── Makefile                     build, test, lint targets
└── config.example.yml           Annotated configuration reference
```

---

## Requirements

| Requirement | Notes |
|---|---|
| **Go 1.21+** | Required for `go install` or build from source |
| **Katana** | Optional — needed for auto-discovery |
| **Httpx** | Optional — needed for auto-discovery |
| **Docker** | Optional — alternative to Go install |

> Without Katana/Httpx, use `--no-discovery` or `--endpoints`.

---

## Disclaimer

This tool is intended for **authorized security testing and performance evaluation only**. You must obtain explicit written permission from the target system owner before running any load tests. Unauthorized use may violate local, national, or international laws. The authors assume **no liability** for misuse or damage caused by this tool.

---

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for contribution guidelines and development setup.

## Security

See [SECURITY.md](SECURITY.md) for vulnerability reporting policy.

## License

[MIT](LICENSE) © [Umut ÖZEN](https://github.com/umutozen)

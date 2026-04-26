[English](../../README.md) · [Türkçe](docs/i18n/README.tr.md) · [Deutsch](docs/i18n/README.de.md) · [Francais](docs/i18n/README.fr.md) · [Espanol](docs/i18n/README.es.md) · [Italiano](docs/i18n/README.it.md) · [Portugues](docs/i18n/README.pt-BR.md) · [Russkij](docs/i18n/README.ru.md) · [Zhongwen Jian](docs/i18n/README.zh-CN.md) · [Zhongwen Fan](docs/i18n/README.zh-TW.md) · [Nihongo](docs/i18n/README.ja.md) · [Hangugeo](docs/i18n/README.ko.md) · [Al-Arabiyya](docs/i18n/README.ar.md) · [Hindi](docs/i18n/README.hi.md) · [Nederlands](docs/i18n/README.nl.md) · [Polski](docs/i18n/README.pl.md) · [Ukrainska](docs/i18n/README.uk.md)

# StormProbe

<p align="center">
  <img src="assets/banner.png" alt="StormProbe" width="100%">
</p>

**discover. load. analyze.**

Production-grade HTTP load tester with autonomous endpoint discovery. Zero dependencies.

[![CI](https://github.com/umutozen/stormprobe/actions/workflows/ci.yml/badge.svg)](https://github.com/umutozen/stormprobe/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/umutozen/stormprobe)](https://goreportcard.com/report/github.com/umutozen/stormprobe)

---

## Features

- **Auto-Discovery** -- Katana crawl + Httpx probe finds live endpoints automatically
- **Multi-Phase Testing** -- Ramp-Up, Sustained, Spike, Recovery
- **Detailed Metrics** -- P50 / P95 / P99 latency, req/s, error classification
- **Dual Reports** -- JSON (machine-readable) + HTML (visual dark-theme dashboard)
- **Zero Dependencies** -- Pure Go standard library, no third-party modules
- **Docker Ready** -- Single command with bundled Katana + Httpx
- **Cross-Platform** -- Linux, macOS, Windows binaries via GoReleaser

## Quick Start

```bash
go run ./cmd https://example.com

go build -o stormprobe ./cmd
./stormprobe --insecure https://example.com
```

### Docker

```bash
# Quick test, no report saved
docker run --rm ghcr.io/umutozen/stormprobe:latest --insecure https://example.com

# Full test with reports saved to ./outputs (Linux / macOS)
docker run --rm -v $(pwd)/outputs:/app/outputs ghcr.io/umutozen/stormprobe:latest --insecure --format both https://example.com

# Full test with reports saved to ./outputs (Windows PowerShell)
docker run --rm -v "${PWD}\outputs:/app/outputs" ghcr.io/umutozen/stormprobe:latest --insecure --format both https://example.com

# High-load spike test (500 concurrent users)
docker run --rm ghcr.io/umutozen/stormprobe:latest --insecure --concurrency-spike 500 https://example.com

# Skip discovery, test root path only
docker run --rm ghcr.io/umutozen/stormprobe:latest --insecure --no-discovery https://example.com
```

## Usage

```
stormprobe [flags] <target-url>

Flags:
  --concurrency-ramp      int       Peak concurrency for ramp-up (default 50)
  --concurrency-sustained int       Concurrency for sustained phase (default 75)
  --concurrency-spike     int       Peak concurrency for spike (default 250)
  --req-per-worker        int       Requests per worker per step (default 15)
  --timeout               duration  Per-request timeout (default 10s)
  --endpoints             string    Endpoints file (one path per line, skips discovery)
  --no-discovery                    Skip katana+httpx, test root path only
  --output                string    Output directory for reports (default ./outputs)
  --format                string    Report format: json, html, both (default both)
  --katana-path           string    Custom katana binary path
  --httpx-path            string    Custom httpx binary path
  --insecure                        Skip TLS certificate verification
  --header, -H            string    Custom HTTP header (repeatable): -H 'Authorization: Bearer TOKEN'
```

### Examples

```bash
stormprobe --insecure --concurrency-spike 500 --req-per-worker 20 https://example.com

stormprobe --insecure --endpoints endpoints.txt https://example.com

stormprobe --insecure --format json --output ./results https://example.com

stormprobe --insecure -H "Authorization: Bearer TOKEN" -H "X-Tenant: acme" https://api.example.com
```

## Test Phases

| Phase | Description |
|---|---|
| **0 -- Discovery** | Katana crawl + Httpx probe, deduplicated endpoint list |
| **1 -- Ramp-Up** | Gradual increase: 5, 15, 30, 50 virtual users |
| **2 -- Sustained** | 3 waves at target concurrency, measures degradation |
| **3 -- Spike** | Sudden burst to peak, then cooldown |
| **4 -- Recovery** | Post-spike health check after 10s cooldown |

## Output

| File | Description |
|---|---|
| `stormprobe_report_<timestamp>.json` | Machine-readable results with all metrics |
| `stormprobe_report_<timestamp>.html` | Visual dashboard, open in any browser |

## Project Structure

```
stormprobe/
    cmd/main.go                    CLI entry point
    internal/
        config/config.go           Shared types, phase generators
        metrics/
            latency.go             Percentile calculation
            errors.go              Error classification
        discovery/
            discovery.go           Discover() public interface
            katana.go              Katana crawler integration
            httpx.go               Httpx probe integration
        runner/
            phase.go               Phase orchestration
            worker.go              Goroutine pool, per-worker RNG
        report/
            json.go                JSON report writer
            html.go                Self-contained HTML report
    Dockerfile                     Multi-stage with katana + httpx
    .goreleaser.yml                Cross-compilation config
    Makefile                       Build, test, lint targets
    config.example.yml             Phase defaults reference
```

## Requirements

- Go 1.21+
- [Katana](https://github.com/projectdiscovery/katana) (ProjectDiscovery) -- optional, for auto-discovery
- [Httpx](https://github.com/projectdiscovery/httpx) (ProjectDiscovery) -- optional, for endpoint validation

```bash
bash scripts/install-tools.sh
```

## Disclaimer

This tool is intended for **authorized security testing and performance evaluation only**. You must obtain explicit written permission from the target system owner before running any load tests. Unauthorized use of this tool against systems you do not own or have permission to test may violate local, national, or international laws. The authors assume no liability for misuse or damage caused by this tool.

## License

MIT

## Author

**Umut ÖZEN** -- [@umutozen](https://github.com/umutozen)

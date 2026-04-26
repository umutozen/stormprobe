# StormProbe v1.0.0 — Initial Release

**StormProbe** is a zero-dependency, production-grade HTTP load testing CLI tool with autonomous endpoint discovery.

---

## Highlights

- **4-Phase Load Test** — Ramp-Up, Sustained, Spike, Recovery pipeline
- **Autonomous Discovery** — Katana (static + JS/headless) crawl + Httpx probing to find all live endpoints automatically
- **Duration Mode** — `--duration 30s` runs each phase for a fixed time instead of fixed request count
- **Alert System** — `--alert-p99`, `--alert-error-rate`, `--alert-rps` produce exit code 1 for CI/CD integration
- **Custom Headers** — `-H "Authorization: Bearer TOKEN"` propagated across all phases and discovery
- **Rich Reports** — Self-contained HTML report with SVG latency chart and error distribution bars, plus JSON output
- **TLS Bypass** — `--insecure` flag for internal/self-signed certificate targets
- **Zero Dependencies** — Pure Go stdlib, no third-party Go modules
- **Cross-Platform** — Pre-built binaries for Linux, macOS, Windows (amd64 & arm64)
- **Docker Ready** — `docker pull ghcr.io/umutozen/stormprobe:v1.0.0`

---

## Installation

### Binary (recommended)
Download the archive for your platform from the Assets section below, extract, and run:

```bash
# Linux / macOS
./stormprobe https://your-target.com

# Windows (PowerShell)
.\stormprobe.exe https://your-target.com
```

### Docker
```bash
docker run --rm ghcr.io/umutozen/stormprobe:v1.0.0 https://your-target.com
```

---

## Quick Usage

```bash
# Full autonomous discovery + all phases
stormprobe https://your-target.com

# Skip discovery, test root only
stormprobe --no-discovery https://your-target.com

# Custom headers, TLS bypass, JSON-only report
stormprobe -H "Authorization: Bearer TOKEN" --insecure --format json https://your-target.com

# Supply your own endpoint list
stormprobe --endpoints ./paths.txt https://your-target.com

# Duration-based test (30s per phase)
stormprobe --duration 30s --no-discovery https://your-target.com

# CI/CD alert thresholds (exit 1 if breached)
stormprobe --alert-p99 500 --alert-error-rate 5 https://your-target.com
```

---

## CLI Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--concurrency-ramp` | 50 | Peak concurrency — Ramp-Up phase |
| `--concurrency-sustained` | 75 | Concurrency — Sustained phase |
| `--concurrency-spike` | 250 | Peak concurrency — Spike phase |
| `--req-per-worker` | 15 | Requests per worker per step |
| `--timeout` | 10s | Per-request timeout |
| `--duration` | — | Per-phase duration (e.g. `30s`, `1m`). Overrides req-per-worker |
| `--alert-p99` | — | Fail (exit 1) if P99 latency exceeds Xms in any phase |
| `--alert-error-rate` | — | Fail (exit 1) if error rate exceeds X% in any phase |
| `--alert-rps` | — | Fail (exit 1) if req/s falls below X in any phase |
| `--no-discovery` | false | Skip Katana+Httpx, test `/` only |
| `--endpoints` | — | Custom endpoints file (one path/line) |
| `-H` / `--header` | — | Custom HTTP header (repeatable) |
| `--insecure` | false | Skip TLS certificate verification |
| `--format` | both | Report format: `json`, `html`, `both` |
| `--output` | `./outputs` | Output directory for reports |

---

## See Also

- [README](https://github.com/umutozen/stormprobe#readme) — Full documentation (16 languages)
- [CHANGELOG](https://github.com/umutozen/stormprobe/blob/main/CHANGELOG.md) — What's new
- [CONTRIBUTING](https://github.com/umutozen/stormprobe/blob/main/CONTRIBUTING.md) — How to contribute

# StormProbe

<p align="center">
  <img src="assets/banner.png" alt="StormProbe Banner" width="100%">
</p>

<h1 align="center">StormProbe</h1>

<p align="center">
  <strong>crawl · discover · stress · decide</strong><br>
  Modern HTTP load testing with autonomous endpoint discovery and production readiness verdicts.
</p>

<p align="center">
  <strong>Not just benchmarking.</strong><br>
  StormProbe doesn't just test systems. It explains whether they actually became safer.
</p>

<p align="center">
  <a href="docs/i18n/README.tr.md">Türkçe</a> · 
  <a href="docs/i18n/README.de.md">Deutsch</a> · 
  <a href="docs/i18n/README.fr.md">Français</a> · 
  <a href="docs/i18n/README.es.md">Español</a> · 
  <a href="docs/i18n/README.it.md">Italiano</a> · 
  <a href="docs/i18n/README.pt-BR.md">Português</a> · 
  <a href="docs/i18n/README.ru.md">Русский</a> · 
  <a href="docs/i18n/README.zh-CN.md">简体中文</a> · 
  <a href="docs/i18n/README.zh-TW.md">繁體中文</a> · 
  <a href="docs/i18n/README.ja.md">日本語</a> · 
  <a href="docs/i18n/README.ko.md">한국어</a> · 
  <a href="docs/i18n/README.ar.md">العربية</a> · 
  <a href="docs/i18n/README.hi.md">हिन्दी</a> · 
  <a href="docs/i18n/README.nl.md">Nederlands</a> · 
  <a href="docs/i18n/README.pl.md">Polski</a> · 
  <a href="docs/i18n/README.uk.md">Українська</a>
</p>

---

## Why StormProbe?

Traditional load testers show numbers.

StormProbe gives decisions.

Instead of manually choosing endpoints and interpreting raw latency tables, StormProbe automatically discovers live endpoints, runs phased stress testing, detects saturation patterns, and produces a final verdict with actionable recommendations.

### Traditional Load Testers vs StormProbe

| Traditional Tools              | StormProbe                              |
|--------------------------------|-----------------------------------------|
| Manual endpoint selection      | Auto-discovers live endpoints           |
| Raw benchmark output           | Final verdict + bottleneck diagnosis    |
| Requires manual interpretation | Actionable recommendations              |
| Single stress pattern          | Ramp-Up + Sustained + Spike + Recovery  |
| "Something is slow"            | "Here is where and why it breaks"       |

---

## What It Does

StormProbe is a production-grade HTTP load testing CLI written in Go.

It combines:

- Endpoint discovery
- Phased concurrency testing
- Latency percentile analysis
- Throughput saturation detection
- Recovery validation
- CI/CD threshold alerts
- JSON + HTML reporting
- Production readiness verdicts

Designed for:

- DevOps engineers
- QA teams
- Penetration testers
- Infrastructure teams
- CI/CD performance gates
- Internal audit and platform teams

---

## Real Output Example

```text
FINAL VERDICT
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  Rating                          : Healthy (Observed Minor Transient Errors)
  Safe concurrency                : up to 250 virtual users
  Recovery                        : OK

  Bottleneck    : Minor transient failures under peak load
  Recommendation : Monitor P99 under real traffic.

  Production Recommendation : Safe for production at tested concurrency levels
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

**Compare Mode — The Decision Engine:**

```text
  stormprobe compare before.json after.json

  Rating           : High Risk → Not Production Ready
  Confidence       : Unstable — Recovery Failed → Low Confidence — Immediate Action Required
  Safe concurrency : 15 VU → 50 VU
  Recovery         : FAILED → OK

  Bottleneck:
    Before: Application/DB processing bottleneck
    After:  No sustained bottleneck detected

  ┌─────────────────────────────────────────────────────────────────┐
  │  Verdict: PARTIAL IMPROVEMENT — Still Not Production Ready     │
  └─────────────────────────────────────────────────────────────────┘

  Why:
    + Safe concurrency increased 233% (15 → 50 VU)
    + Recovery restored (was FAILED, now OK)
    + Confidence improved (Unstable → Low Confidence)
    + Lower transient error volume observed (6169 → 4 failed)
```

Not just metrics. A decision your team can act on.

---

## Features

| Feature              | Details                                                                              |
|----------------------|--------------------------------------------------------------------------------------|
| Auto Discovery       | Katana crawl + Httpx probe with deduplicated live endpoints                          |
| 4-Phase Load Test    | Ramp-Up → Sustained → Spike → Recovery                                               |
| Verdict Engine       | Healthy / Degraded Under Load / Not Production Ready / High Risk                    |
| Confidence Qualifier | `High Confidence` / `Observed Minor Transient Errors` / `Recovery Impaired`         |
| Compare Mode         | Verdict delta engine: IMPROVED / PARTIAL IMPROVEMENT / STABLE / DEGRADED            |
| Rich Metrics         | Avg / P50 / P95 / P99 latency + req/s + error classification                        |
| Throughput Analysis  | Detects saturation before hard failures happen                                       |
| Recovery Validation  | Multi-metric check: success rate + avg + P95 vs baseline                             |
| Alert Thresholds     | CI/CD-ready exit codes for latency, error rate, and RPS                              |
| JSON Semantics       | `report_id`, `confidence`, `null` degradation fields — dashboard-ready               |
| Dual Reports         | JSON + self-contained HTML dashboard                                                 |
| Custom Headers       | Auth tokens, tenant headers, cookies                                                 |
| Docker Ready         | Single image with bundled tools                                                      |
| Cross Platform       | Linux, macOS, Windows — amd64 & arm64                                                |
| Pure Go Core         | Fast install, minimal operational overhead                                           |
| **Adaptive Load Intelligence** | Auto-detects server stack (IIS/Tomcat/PHP) and applies stack-specific diagnostic profiles |
| **Session Lock Detector** | Identifies ASP.NET exclusive session lock via P50/P99 ratio serialization pattern |
| **Stack-Aware Verdict** | Primary Cause + Secondary Contributors + "Why this diagnosis?" per server type    |

---

## Quick Start

```bash
# Install
go install github.com/umutozen/stormprobe/cmd/stormprobe@latest

# Full test with auto-discovery
stormprobe https://example.com

# Skip discovery and test root path only
stormprobe --no-discovery https://example.com

# API test with auth header
stormprobe --insecure \
  -H "Authorization: Bearer TOKEN" \
  https://api.example.com

# CI/CD gate: fail build if latency exceeds threshold
stormprobe \
  --alert-p99 500 \
  --alert-error-rate 5 \
  https://example.com

# Use a JSON config file
stormprobe --config stormprobe.json

# Config file + CLI override (flag always wins)
stormprobe --config stormprobe.json --duration 10s

# Compare two runs — verdict delta engine
stormprobe compare outputs/before.json outputs/after.json
```

---

## Test Flow

```
Discovery → Ramp-Up → Sustained → Spike → Recovery → Verdict
```

| Phase     | Purpose                                |
|-----------|----------------------------------------|
| Discovery | Crawl and validate live endpoints      |
| Ramp-Up   | Gradual concurrency increase           |
| Sustained | Validate steady-state performance      |
| Spike     | Sudden burst traffic simulation        |
| Recovery  | Verify post-spike system health        |
| Verdict   | Safe concurrency + bottleneck analysis |

---

## Installation

### Go Install (Recommended)

```bash
go install github.com/umutozen/stormprobe/cmd/stormprobe@latest
stormprobe --help
```

### Docker

```bash
docker run --rm \
  ghcr.io/umutozen/stormprobe:latest \
  --insecure https://example.com
```

### Prebuilt Releases

Download binaries from the [GitHub Releases](https://github.com/umutozen/stormprobe/releases) page.

---

## Flags

| Flag                        | Description                                                  |
|-----------------------------|--------------------------------------------------------------|
| `--duration`                | Run each phase for X duration (e.g. `30s`, `1m`)            |
| `--alert-p99`               | Exit 1 if P99 latency exceeds threshold (ms)                 |
| `--alert-error-rate`        | Exit 1 if error rate exceeds threshold (%)                   |
| `--alert-rps`               | Exit 1 if req/s drops below threshold                        |
| `-H` / `--header`           | Custom HTTP header (repeatable)                              |
| `--insecure`                | Skip TLS verification                                        |
| `--no-discovery`            | Skip Katana+Httpx, test root path only                       |
| `--endpoints`               | Load custom endpoint list from file                          |
| `--format`                  | Report format: `json`, `html`, `both`                        |
| `--concurrency-ramp`        | Ramp-Up peak concurrency                                     |
| `--concurrency-sustained`   | Sustained concurrency                                        |
| `--concurrency-spike`       | Spike peak concurrency                                       |
| `--req-per-worker`          | Requests per worker                                          |
| `--timeout`                 | Per-request timeout                                          |
| `--config`                  | Load settings from a JSON config file                        |
| `--success-match`           | Regex: response body must match for success                  |
| `--fail-match`              | Regex: response body match means failure                     |
| `--auto-profile`            | Auto-detect server stack and apply diagnostic profile (default: true) |
| `--profile`                 | Manual profile override: `iis`, `tomcat`, `php`              |
| `--no-fingerprint`          | Skip server stack fingerprinting                             |
| `--show-profile`            | Show detected stack and diagnostic profile, then exit        |

### Compare Subcommand

```bash
stormprobe compare <before.json> <after.json>
```

Compares two JSON reports and outputs a **verdict delta** — not a raw diff.
Produces one of: `IMPROVED` / `STABLE` / `STABLE WITH RISK` / `DEGRADED` / `CRITICAL REGRESSION`.

---

## CI/CD Example

### GitHub Actions

```yaml
- name: Performance Gate
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

This allows performance regressions to fail deployments automatically.

---

## Example Use Cases

**Before Production Release** — Validate whether concurrency changes break latency targets.

**Security + Performance Testing** — Discover hidden live endpoints and stress them safely.

**Infrastructure Audits** — Measure saturation points before traffic spikes happen.

**Regression Detection** — Compare before/after deployment performance in CI.

---

## Requirements

| Requirement | Notes                            |
|-------------|----------------------------------|
| Go 1.21+    | Required for `go install`        |
| Katana      | Optional, for endpoint discovery |
| Httpx       | Optional, for endpoint discovery |
| Docker      | Optional alternative runtime     |

Without Katana/Httpx, use `--no-discovery` or `--endpoints endpoints.txt`.

---

## Disclaimer

This tool is intended for authorized testing only.

Always obtain explicit permission before running load tests against any target system.

Unauthorized use may violate local laws and policies.

---

## Contributing

Contributions are welcome. See `CONTRIBUTING.md` for development setup and contribution guidelines.

## Security

See `SECURITY.md` for responsible vulnerability disclosure.

## License

MIT License © Umut ÖZEN

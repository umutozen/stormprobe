# StormProbe

<p align="center">
  <img src="../../assets/banner.png" alt="StormProbe Banner" width="100%">
</p>

<h1 align="center">StormProbe</h1>

<p align="center">
  <strong>crawl · discover · stress · decide</strong><br>
  Modern HTTP load testing with autonomous endpoint discovery & production-readiness verdicts.
</p>

<p align="center">
  <strong>Not just benchmarking.</strong><br>
  StormProbe doesn't just test. It explains whether your system actually became safer.
</p>

---

## Why StormProbe?

Traditional load testers show numbers.<br>StormProbe gives decisions.

Instead of manually selecting endpoints and interpreting raw latency tables, StormProbe automatically discovers live endpoints, runs phased stress tests, detects saturation patterns, and produces a final verdict with actionable recommendations.

### Traditional Load Testers vs StormProbe

| Traditional Tools | StormProbe |
|---|---|
| Manual endpoint selection | Discovers live endpoints automatically |
| Raw benchmark output | Final verdict + bottleneck diagnosis |
| Requires manual interpretation | Actionable recommendations |
| Single stress pattern | Ramp-Up + Sustained + Spike + Recovery |
| "Something is slow" | "Here is where and why it breaks" |

---

## What does it do?

StormProbe is a production-grade HTTP load testing CLI written in Go.

- Endpoint discovery
- Phased concurrency testing
- Latency percentile analysis
- Throughput saturation detection
- Recovery validation
- CI/CD threshold alerts
- JSON + HTML reporting
- Production readiness verdicts

---

## Real Output Example

```text
FINAL VERDICT

Rating                          : Degraded Under Load
Safe concurrency                : up to 30 virtual users
Perf. degradation starts at     : 50 vu
Severe degradation starts at    : 150 vu
Throughput plateau              : ~36 req/s
Recovery                        : OK

Bottleneck    : Throughput saturation
Detail        : Latency grew 5.4x while throughput grew only 2.1x
                for 10x concurrency increase.

Recommendation           : Cap production traffic at ~30 concurrent users.
Production Recommendation : Safe for low-to-medium traffic only
```

**Compare Mode — The Decision Engine:**

```text
  stormprobe compare before.json after.json

  Rating           : High Risk → Not Production Ready
  Confidence       : Unstable — Recovery Failed → Low Confidence — Immediate Action Required
  Safe concurrency : 15 VU → 50 VU
  Recovery         : FAILED → OK

  ┌─────────────────────────────────────────────────────────────────┐
  │  Verdict: PARTIAL IMPROVEMENT — Still Not Production Ready     │
  └─────────────────────────────────────────────────────────────────┘

  Why:
    + Safe concurrency increased 233% (15 → 50 VU)
    + Recovery restored (was FAILED, now OK)
    + Confidence improved (Unstable → Low Confidence)
    + Lower transient error volume observed (6169 → 4 failed)
```

---

## Features

| Feature | Details |
|---|---|
| Auto Discovery | Katana crawl + Httpx probe with deduplicated live endpoints |
| 4-Phase Load Test | Ramp-Up → Sustained → Spike → Recovery |
| Verdict Engine | Safe concurrency, degradation point, bottleneck classification |
| Rich Metrics | Avg / P50 / P95 / P99 latency + req/s + error classification |
| Throughput Analysis | Detects saturation before hard failures happen |
| Recovery Validation | Confirms whether the system recovers after traffic spikes |
| Alert Thresholds | CI/CD-ready exit codes for latency, error rate, and RPS |
| Dual Reports | JSON + self-contained HTML dashboard |
| Custom Headers | Auth tokens, tenant headers, cookies |
| Docker Ready | Single image with bundled tools |
| Cross Platform | Linux, macOS, Windows — amd64 & arm64 |
| Pure Go Core | Fast install, minimal operational overhead |

---

## Quick Start

```bash
# Install
go install github.com/umutozen/stormprobe/cmd/stormprobe@latest

# Full test with auto-discovery
stormprobe https://example.com

# Skip discovery, test root path only
stormprobe --no-discovery https://example.com

# API test with auth header
stormprobe --insecure \
  -H "Authorization: Bearer TOKEN" \
  https://api.example.com

# CI/CD gate
stormprobe \\
  --alert-p99 500 \\
  --alert-error-rate 5 \\
  https://example.com

# Use JSON config file
stormprobe --config stormprobe.json
```

---

## Test Flow

```text
Discovery → Ramp-Up → Sustained → Spike → Recovery → Verdict
```

| Phase | Purpose |
|---|---|
| Discovery | Crawl and validate live endpoints |
| Ramp-Up | Gradual concurrency increase |
| Sustained | Validate steady-state performance |
| Spike | Sudden burst traffic simulation |
| Recovery | Verify post-spike system health |
| Verdict | Safe concurrency + bottleneck analysis |

---

## Installation

### Go Install

```bash
go install github.com/umutozen/stormprobe/cmd/stormprobe@latest
```

Recommended way.

### Docker

```bash
docker run --rm ghcr.io/umutozen/stormprobe:latest --insecure https://example.com
```

---

## CI/CD

```yaml
- name: Performance Gate
  run: |
    go install github.com/umutozen/stormprobe/cmd/stormprobe@latest
    stormprobe \
      --no-discovery \
      --duration 30s \
      --alert-p99 500 \
      --alert-error-rate 2 \
      https://staging.example.com
```

This enables performance regressions to fail deployments automatically.

---

## Use Cases

**Pre-Production Release** — Validate if concurrency changes break latency targets.

**Security + Performance** — Discover hidden live endpoints and safely stress test them.

**Infrastructure Audits** — Measure saturation points before traffic spikes.

**Regression Detection** — Compare before/after deployment performance in CI.

---

## Flags

| Flag | Description |
|---|---|
| `--duration` | Run each phase for X duration |
| `--alert-p99` | Exit 1 if P99 latency exceeds threshold |
| `--alert-error-rate` | Exit 1 if error rate exceeds threshold |
| `--alert-rps` | Exit 1 if req/s drops below threshold |
| `-H` / `--header` | Custom HTTP header |
| `--insecure` | Skip TLS verification |
| `--no-discovery` | Skip Katana+Httpx |
| `--endpoints` | Load custom endpoint list |
| `--format` | Report format |
| `--concurrency-ramp` | Ramp-Up peak concurrency |
| `--concurrency-sustained` | Sustained concurrency |
| `--concurrency-spike` | Spike peak concurrency |
| `--req-per-worker` | Requests per worker |
| `--timeout` | Per-request timeout |
| `--config` | Load settings from JSON config file |
| `--auto-profile` | Auto-detect server stack and apply diagnostic profile (default: true) |
| `--profile` | Manual profile override: iis, tomcat, php |
| `--no-fingerprint` | Skip server stack fingerprinting |
| `--show-profile` | Show detected stack and diagnostic profile, then exit |

---

## Requirements

| Requirement | Notes |
|---|---|
| Go 1.21+ | Required for go install |
| Katana | Optional, for discovery |
| Httpx | Optional, for discovery |
| Docker | Optional alternative runtime |

---

## Disclaimer

This tool is intended for authorized testing only.

Always obtain explicit permission.

Unauthorized use may violate local laws and policies.

---

## Contributing

Contributions are welcome.

## Security

See SECURITY.md.

## License

MIT License © Umut ÖZEN

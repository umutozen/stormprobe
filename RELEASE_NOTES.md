# StormProbe v1.0.0 — Release Notes

> **crawl · discover · stress · decide**

StormProbe is a zero-dependency, production-grade HTTP load testing engine written in pure Go. It discovers your endpoints, runs phased stress tests, detects saturation patterns, and delivers a final production-readiness verdict.

---

## What's New in v1.0.0

### Verdict Engine

The core differentiator. StormProbe doesn't just output numbers — it decides.

| Classification | Meaning |
|---|---|
| `Healthy` | System handles target concurrency within acceptable thresholds |
| `Degraded` | Performance degrades under load but remains functional |
| `Unstable` | Intermittent failures or severe latency growth detected |
| `Critical` | System cannot sustain the requested workload |

Each verdict includes a **confidence qualifier** — `High Confidence`, `Observed Minor Transient Errors`, or `Recovery Impaired` — so you know how much to trust the result.

### Compare Mode

```bash
stormprobe compare outputs/before.json outputs/after.json
```

Not a diff tool — a **verdict delta engine**. Outputs one of:

`IMPROVED` · `STABLE` · `STABLE WITH RISK` · `DEGRADED` · `CRITICAL REGRESSION`

Includes a **Why** section explaining the exact metrics that changed and their operational impact.

### Autonomous Discovery

StormProbe uses **Katana** (static + JS/headless crawl) and **Httpx** (probe validation) to automatically find all live endpoints. No manual URL lists required.

### 4-Phase Load Pipeline

```
Discovery → Ramp-Up → Sustained → Spike → Recovery → Verdict
```

| Phase | Purpose |
|---|---|
| Ramp-Up | Gradually increase concurrency to find the breaking point |
| Sustained | Validate steady-state performance under constant load |
| Spike | Simulate sudden traffic bursts to test resilience |
| Recovery | Confirm the system returns to healthy state after spike |

### Duration Mode

```bash
stormprobe --duration 30s --no-discovery https://example.com
```

Run each phase for a fixed time window instead of a fixed request count. Ideal for consistent CI/CD benchmarks.

### CI/CD Alert System

```bash
stormprobe --alert-p99 500 --alert-error-rate 5 --alert-rps 100 https://staging.example.com
```

Any threshold breach produces **exit code 1** — plug directly into your deployment pipeline as a performance gate.

### Rich Reports

- **HTML** — Self-contained dashboard with executive summary card, SVG latency chart, and priority action list
- **JSON** — Machine-readable output with `report_id`, `confidence`, `degradation_at_vu`, `failure_point_vu` fields
- Fields are `null` when not detected (not `0`) for clean downstream parsing

### JSON Config File

```bash
stormprobe --config stormprobe.json
```

Load all settings from a JSON config file. CLI flags override config values. Supports single-target and multi-target configurations.

### Multi-Target Mode

```bash
stormprobe --config belediyeler.json
```

Define multiple targets in a single config, run tests sequentially, get a per-target summary table, and automatic verdict comparison between targets.

### Response Body Validation

```bash
stormprobe --fail-match "error|timeout|fault" --success-match "OK" https://example.com
```

Regex-based response body inspection. Catches HTTP 200 responses that contain error content — directly improves verdict accuracy.

### Adaptive Load Intelligence

StormProbe now detects the target server stack before testing begins.

```
stormprobe --show-profile https://example.com

Edge Layer       : Nginx (High)
Application Layer: ASP.NET / IIS (Medium)
Diagnostic Profile: iis
```

Based on detection, it applies a stack-specific diagnostic profile and produces a stack-aware verdict:

```
Detected Stack : IIS + ASP.NET (High Confidence)

Primary Cause:
  IIS App Pool worker limit reached

Secondary Contributors:
  - DB connection pool pressure
  - Session state locking

Why this diagnosis?
  - Latency spike begins after 75 VU
  - Throughput plateaus at ~125 req/s despite 10x concurrency increase
  - P99 latency reaches 2414ms — queue buildup pattern, not hard limit
```

Supported profiles: `iis`, `tomcat`, `php`, `default`.
Confidence fallback: Low/Unknown → default profile (never forces wrong diagnosis).

### Session Lock Detector

Detects ASP.NET exclusive session lock contention automatically via P50/P99 ratio analysis.

When P99/P50 ratio is tight across all concurrency levels (requests serializing), and throughput fails to scale with concurrency, StormProbe emits:

```
Primary Cause: ASP.NET exclusive session lock detected
Secondary Contributors:
  - Session.Abandon() or async session provider missing
  - Authenticated endpoints serialize under IIS session mutex
  - Recommendation: enable read-only session or use distributed cache (Redis)
```

No configuration required — fires automatically during bottleneck analysis.

### Live Progress Bar

Real-time progress indicator during test execution showing completion percentage, request count, req/s, and elapsed time.

---

## Installation

### Go Install (recommended)

```bash
go install github.com/umutozen/stormprobe/cmd/stormprobe@latest
```

### Pre-built Binaries

Download the archive for your platform from the [Releases](https://github.com/umutozen/stormprobe/releases) page:

| Platform | Architecture |
|---|---|
| Linux | amd64, arm64 |
| macOS | amd64, arm64 |
| Windows | amd64, arm64 |

```bash
# Linux / macOS
./stormprobe https://your-target.com

# Windows
.\stormprobe.exe https://your-target.com
```

### Docker

```bash
docker run --rm ghcr.io/umutozen/stormprobe:v1.0.0 --insecure https://your-target.com
```

Image includes Katana and Httpx pre-installed.

---

## Usage Examples

```bash
# Full autonomous discovery + 4-phase load test
stormprobe https://your-target.com

# Skip discovery, test root path only
stormprobe --no-discovery https://your-target.com

# Custom auth header + TLS bypass + JSON-only output
stormprobe \
  -H "Authorization: Bearer TOKEN" \
  --insecure --format json \
  https://your-target.com

# Custom endpoint list
stormprobe --endpoints ./paths.txt https://your-target.com

# Duration-based phases (30s each)
stormprobe --duration 30s --no-discovery https://your-target.com

# CI/CD performance gate (exit 1 if breached)
stormprobe \
  --alert-p99 500 \
  --alert-error-rate 5 \
  --alert-rps 100 \
  https://staging.example.com

# Compare two test runs
stormprobe compare outputs/before.json outputs/after.json
```

---

## CLI Flags

| Flag | Default | Description |
|---|---|---|
| `--concurrency-ramp` | `50` | Peak concurrency during Ramp-Up phase |
| `--concurrency-sustained` | `75` | Concurrency during Sustained phase |
| `--concurrency-spike` | `250` | Peak concurrency during Spike phase |
| `--req-per-worker` | `15` | Requests per worker per step |
| `--timeout` | `10s` | Per-request timeout |
| `--duration` | — | Per-phase duration (e.g. `30s`, `1m`). Overrides `--req-per-worker` |
| `--alert-p99` | — | Exit 1 if P99 latency exceeds threshold (ms) |
| `--alert-error-rate` | — | Exit 1 if error rate exceeds threshold (%) |
| `--alert-rps` | — | Exit 1 if req/s drops below threshold |
| `--no-discovery` | `false` | Skip Katana+Httpx, test `/` only |
| `--endpoints` | — | Custom endpoint list file (one path per line) |
| `-H` / `--header` | — | Custom HTTP header (repeatable) |
| `--insecure` | `false` | Skip TLS certificate verification |
| `--format` | `both` | Report format: `json`, `html`, `both` |
| `--output` | `./outputs` | Output directory for reports |
| `--config` | — | Load settings from a JSON config file |
| `--success-match` | — | Regex: response body must match for success |
| `--fail-match` | — | Regex: response body match means failure |
| `--auto-profile` | `true` | Auto-detect server stack and apply diagnostic profile |
| `--profile` | — | Manual profile override: `iis`, `tomcat`, `php` |
| `--no-fingerprint` | `false` | Skip server stack fingerprinting |
| `--show-profile` | `false` | Show detected stack and profile, then exit |
| `compare <before> <after>` | — | Verdict delta engine: compare two JSON reports |

---

## Architecture

```
internal/
├── alert/       # CI/CD threshold evaluation
├── compare/     # Verdict delta engine
├── config/      # Flag parsing & validation
├── discovery/   # Katana + Httpx orchestration
├── metrics/     # Latency percentiles, throughput, error classification
├── report/      # JSON + HTML generation
├── runner/      # 4-phase load executor
└── verdict/     # Classification + confidence + bottleneck logic
```

**Zero Go dependencies.** Pure standard library. No vendor lock-in.

---

## Technical Details

- **Recovery validation** checks success rate + avg latency + P95 against baseline (not just P99)
- **JSON reports** include `report_id` (auto-derived from timestamp + hostname) for tracking
- **Throughput analysis** detects saturation patterns: latency growth vs throughput plateau ratio
- **Verdict confidence** is computed from error consistency, recovery success, and metric variance

---

## Links

| Resource | URL |
|---|---|
| Documentation | [README](https://github.com/umutozen/stormprobe#readme) (16 languages) |
| Changelog | [CHANGELOG.md](https://github.com/umutozen/stormprobe/blob/main/CHANGELOG.md) |
| Contributing | [CONTRIBUTING.md](https://github.com/umutozen/stormprobe/blob/main/CONTRIBUTING.md) |
| Security | [SECURITY.md](https://github.com/umutozen/stormprobe/blob/main/SECURITY.md) |
| License | [MIT](https://github.com/umutozen/stormprobe/blob/main/LICENSE) |

---

MIT License © Umut ÖZEN

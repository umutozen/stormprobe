[English](../../README.md) · [Türkçe](README.tr.md) · [Deutsch](README.de.md) · [Francais](README.fr.md) · [Espanol](README.es.md) · [Italiano](README.it.md) · [Portugues](README.pt-BR.md) · [Russkij](README.ru.md) · [Zhongwen Jian](README.zh-CN.md) · [Zhongwen Fan](README.zh-TW.md) · [Nihongo](README.ja.md) · [Hangugeo](README.ko.md) · [Al-Arabiyya](README.ar.md) · [Hindi](README.hi.md) · [Nederlands](README.nl.md) · [Polski](README.pl.md) · [Ukrainska](README.uk.md)

<p align="center">
  <img src="../../assets/banner.png" alt="StormProbe Banner" width="100%">
</p>

<h1 align="center">StormProbe</h1>

<p align="center">
  <strong>khoje · load · vishleshan</strong><br>
  Svayat endpoint khoj ke saath utpadan-grade HTTP load testing tool. Zero bahari dependencies.
</p>

<p align="center">
  <a href="https://github.com/umutozen/stormprobe/actions/workflows/ci.yml">
    <img src="https://github.com/umutozen/stormprobe/actions/workflows/ci.yml/badge.svg" alt="CI">
  </a>
  <a href="https://goreportcard.com/report/github.com/umutozen/stormprobe">
    <img src="https://goreportcard.com/badge/github.com/umutozen/stormprobe" alt="Go Report Card">
  </a>
  <a href="https://github.com/umutozen/stormprobe/releases">
    <img src="https://img.shields.io/github/v/release/umutozen/stormprobe" alt="Navanatam sanskarana">
  </a>
  <a href="../../LICENSE">
    <img src="https://img.shields.io/badge/license-MIT-blue.svg" alt="MIT laisens">
  </a>
</p>

---

## StormProbe kya hai?

StormProbe ek **shunya-dependency** HTTP load testing tool hai jo shuddh Go mein likha gaya hai. Yah [Katana](https://github.com/projectdiscovery/katana) aur [Httpx](https://github.com/projectdiscovery/httpx) ke saath aapke application ke endpoints ko svachalita roop se khojta hai, phir ek sanrachit 4-charan load test chalata hai — **Ramp-Up → Sustained → Spike → Recovery** — aur P50/P95/P99 latency metrics ke saath JSON aur HTML reports utpann karta hai.

Inke liye design kiya gaya:
- **CI/CD pipelines** — exit codes ke saath alert thresholds
- **Penetration testers** — live endpoints ki svachalita khoj
- **DevOps engineers** — deployment se pahle aur baad benchmark
- **QA teams** — performance SLA validation

---

## Visheshtaen

| Visheshata | Vivaran |
|---|---|
| **Svachalita Khoj** | Katana crawl + Httpx probe, deduplicated live endpoint list |
| **4-Charan Load Test** | Ramp-Up → Sustained → Spike → Recovery |
| **Samriddh Metrics** | P50/P95/P99 latency, req/s, prati charan error classification |
| **Avadhi Mode** | Samay-adharit charan (`--duration 30s`) anurodh count ki jagah |
| **Alert Thresholds** | CI/CD-ready `exit 1` P99, error rate, ya RPS ulanghan par |
| **Dviguna Reports** | JSON (machine-readable) + standalone HTML dark dashboard |
| **Custom Headers** | Bearer tokens, tenant headers, cookies — har jagah propagated |
| **Docker Ready** | Katana + Httpx sammilit single image |
| **Cross-Platform** | Linux, macOS, Windows — amd64 & arm64 |
| **Zero Dependencies** | Shuddh Go 1.21+ stdlib, `go install` kaafi hai |

---

## Jald Shuru

```bash
# Sthapana (Go 1.21+ avashyak)
go install github.com/umutozen/stormprobe/cmd/stormprobe@latest

# Buniyadi test — svachalita khoj, 4 charan
stormprobe https://example.com

# Khoj chhoden, sirf root path test karen
stormprobe --no-discovery https://example.com

# TLS + auth header + sirf JSON report
stormprobe --insecure -H "Authorization: Bearer TOKEN" --format json https://api.example.com

# Avadhi-adharit test: prati charan 30 second
stormprobe --insecure --duration 30s --no-discovery https://example.com

# CI/CD: exit 1 agar P99 > 500ms ya error > 5%
stormprobe --alert-p99 500 --alert-error-rate 5 https://example.com
```

---

## Sthapana

### Vikalp 1 — `go install` *(anushansit)*

```bash
go install github.com/umutozen/stormprobe/cmd/stormprobe@latest
stormprobe --help
```

### Vikalp 2 — Pre-built Binary *(Go avashyak nahin)*

[Releases page](https://github.com/umutozen/stormprobe/releases) se download karen.

```bash
# Linux / macOS
chmod +x stormprobe
sudo mv stormprobe /usr/local/bin/

# Windows PowerShell
Move-Item stormprobe.exe C:\Windows\System32\stormprobe.exe
```

### Vikalp 3 — Source se build karen

```bash
git clone https://github.com/umutozen/stormprobe.git
cd stormprobe
go build -o stormprobe ./cmd/stormprobe      # Linux / macOS
go build -o stormprobe.exe ./cmd/stormprobe  # Windows
```

### Vikalp 4 — Docker

```bash
docker run --rm ghcr.io/umutozen/stormprobe:latest --insecure https://example.com

# Linux / macOS
docker run --rm -v $(pwd)/outputs:/app/outputs \
  ghcr.io/umutozen/stormprobe:latest --insecure --format both https://example.com

# Windows PowerShell
docker run --rm -v "${PWD}\outputs:/app/outputs" `
  ghcr.io/umutozen/stormprobe:latest --insecure --format both https://example.com
```

### Vaikalpik — Khoj Upkaran

```bash
bash scripts/install-tools.sh
```

> Katana/Httpx ke bina `--no-discovery` ya `--endpoints` upayog karen.

---

## CLI Sandarbh

### Samantarta aur Load

| Flag | Default | Vivaran |
|---|---|---|
| `--concurrency-ramp` | `50` | Ramp-Up charan ke liye peak concurrency |
| `--concurrency-sustained` | `75` | Sustained charan ke liye concurrency |
| `--concurrency-spike` | `250` | Spike charan ke liye peak concurrency |
| `--req-per-worker` | `15` | Prati worker anurodh (`--duration` set hone par avaganya) |
| `--duration` | `0` (asamarth) | Samay-adharit charan avadhi, jaise `30s`, `1m` |
| `--timeout` | `10s` | Prati anurodh timeout |

### Khoj

| Flag | Default | Vivaran |
|---|---|---|
| `--no-discovery` | `false` | Katana+Httpx chhoden, sirf root test karen |
| `--endpoints` | `""` | File se endpoints load karen |
| `--katana-path` | `""` | Custom Katana binary path |
| `--httpx-path` | `""` | Custom Httpx binary path |

### HTTP aur Suraksha

| Flag | Default | Vivaran |
|---|---|---|
| `-H`, `--header` | — | Custom header (doharaya ja sakta hai) |
| `--insecure` | `false` | TLS certificate verification chhodna |

### Nikay

| Flag | Default | Vivaran |
|---|---|---|
| `--format` | `both` | Report format: `json`, `html`, `both` |
| `--output` | `./outputs` | Reports ke liye output directory |

### CI/CD Alerts

| Flag | Default | Vivaran |
|---|---|---|
| `--alert-p99` | `0` (asamarth) | exit 1 agar P99 (ms) threshold se adhik |
| `--alert-error-rate` | `0` (asamarth) | exit 1 agar error rate (%) threshold se adhik |
| `--alert-rps` | `0` (asamarth) | exit 1 agar req/s threshold se kam |

---

## Test Charan

| Charan | Vivaran |
|---|---|
| **0 — Discovery** | Katana crawl + Httpx probe, deduplicated list |
| **1 — Ramp-Up** | Dheere-dheere concurrency badhata hai: 5 → 15 → 30 → peak |
| **2 — Sustained** | Target concurrency par 3 lahrein |
| **3 — Spike** | Peak tak achanak vriddhi, phir thanda |
| **4 — Recovery** | 10s thande ke baad post-spike health check |

---

## Nikay aur Reports

| File | Vivaran |
|---|---|
| `stormprobe_report_<timestamp>.json` | Prati charan sabhi metrics, machine-readable |
| `stormprobe_report_<timestamp>.html` | SVG latency chart ke saath standalone dark dashboard |

---

## Avasyambhavi Suchana

Yah upkaran **keval adhikarit security testing aur performance evaluation** ke liye hai. Kisi bhi load test chalane se pahle target system ke malik se spast likhit anumati lena avashyak hai. Lekhak isa upkaran ke durupyog ya nuksan ke liye **koi jimmedari nahin leta**.

---

## Laisens

[MIT](../../LICENSE) © [Umut ÖZEN](https://github.com/umutozen)

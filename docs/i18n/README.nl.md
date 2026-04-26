[English](../../README.md) · [Türkçe](README.tr.md) · [Deutsch](README.de.md) · [Francais](README.fr.md) · [Espanol](README.es.md) · [Italiano](README.it.md) · [Portugues](README.pt-BR.md) · [Russkij](README.ru.md) · [Zhongwen Jian](README.zh-CN.md) · [Zhongwen Fan](README.zh-TW.md) · [Nihongo](README.ja.md) · [Hangugeo](README.ko.md) · [Al-Arabiyya](README.ar.md) · [Hindi](README.hi.md) · [Nederlands](README.nl.md) · [Polski](README.pl.md) · [Ukrainska](README.uk.md)

<p align="center">
  <img src="../../assets/banner.png" alt="StormProbe Banner" width="100%">
</p>

<h1 align="center">StormProbe</h1>

<p align="center">
  <strong>ontdekken · laden · analyseren</strong><br>
  HTTP-belastingstesttool op productieniveau met autonome endpoint-detectie. Nul externe afhankelijkheden.
</p>

<p align="center">
  <a href="https://github.com/umutozen/stormprobe/actions/workflows/ci.yml">
    <img src="https://github.com/umutozen/stormprobe/actions/workflows/ci.yml/badge.svg" alt="CI">
  </a>
  <a href="https://goreportcard.com/report/github.com/umutozen/stormprobe">
    <img src="https://goreportcard.com/badge/github.com/umutozen/stormprobe" alt="Go Report Card">
  </a>
  <a href="https://github.com/umutozen/stormprobe/releases">
    <img src="https://img.shields.io/github/v/release/umutozen/stormprobe" alt="Nieuwste versie">
  </a>
  <a href="../../LICENSE">
    <img src="https://img.shields.io/badge/license-MIT-blue.svg" alt="MIT-licentie">
  </a>
</p>

---

## Wat is StormProbe?

StormProbe is een **afhankelijkheidsvrije** HTTP-belastingstesttool geschreven in puur Go. Het ontdekt automatisch de endpoints van uw applicatie met [Katana](https://github.com/projectdiscovery/katana) en [Httpx](https://github.com/projectdiscovery/httpx), voert vervolgens een gestructureerde 4-fasen belastingstest uit — **Ramp-Up → Sustained → Spike → Recovery** — en produceert JSON- en HTML-rapporten met P50/P95/P99-latentiemetrieken.

Ontworpen voor:
- **CI/CD-pijplijnen** — waarschuwingsdrempels met exit-codes
- **Penetratietesters** — automatische detectie van actieve endpoints
- **DevOps-ingenieurs** — benchmarks voor en na deployment
- **QA-teams** — validatie van performance-SLA's

---

## Functies

| Functie | Details |
|---|---|
| **Automatische detectie** | Katana-crawl + Httpx-probe, gededupliceerde lijst |
| **4-fasen belastingstest** | Ramp-Up → Sustained → Spike → Recovery |
| **Rijke metrieken** | P50/P95/P99-latentie, req/s, foutclassificatie per fase |
| **Duurmodus** | Tijdsgebaseerde fasen (`--duration 30s`) in plaats van aantalverzoeken |
| **Waarschuwingsdrempels** | CI/CD-compatibele `exit 1` bij P99-, foutpercentage- of RPS-overtreding |
| **Dubbele rapporten** | JSON (machineleesbaar) + zelfstandig HTML-donker dashboard |
| **Aangepaste headers** | Bearer-tokens, tenant-headers, cookies — overal doorgegeven |
| **Docker-klaar** | Enkelvoudige image met ingebouwde Katana + Httpx |
| **Platformonafhankelijk** | Linux, macOS, Windows — amd64 & arm64 |
| **Nul afhankelijkheden** | Puur Go 1.21+ stdlib, `go install` volstaat |

---

## Snelstart

```bash
# Installatie (Go 1.21+ vereist)
go install github.com/umutozen/stormprobe/cmd/stormprobe@latest

# Basistest — automatische detectie, alle 4 fasen
stormprobe https://example.com

# Detectie overslaan, alleen rootpad testen
stormprobe --no-discovery https://example.com

# TLS + auth-header + alleen JSON-rapport
stormprobe --insecure -H "Authorization: Bearer TOKEN" --format json https://api.example.com

# Tijdsgebaseerde test: 30 seconden per fase
stormprobe --insecure --duration 30s --no-discovery https://example.com

# CI/CD: exit 1 als P99 > 500ms of fouten > 5%
stormprobe --alert-p99 500 --alert-error-rate 5 https://example.com
```

---

## Installatie

### Optie 1 — `go install` *(aanbevolen)*

```bash
go install github.com/umutozen/stormprobe/cmd/stormprobe@latest
stormprobe --help
```

### Optie 2 — Voorgebouwd binair bestand *(zonder Go)*

Downloaden van de [Releases-pagina](https://github.com/umutozen/stormprobe/releases).

```bash
# Linux / macOS
chmod +x stormprobe
sudo mv stormprobe /usr/local/bin/

# Windows PowerShell
Move-Item stormprobe.exe C:\Windows\System32\stormprobe.exe
```

### Optie 3 — Bouwen vanuit broncode

```bash
git clone https://github.com/umutozen/stormprobe.git
cd stormprobe
go build -o stormprobe ./cmd/stormprobe      # Linux / macOS
go build -o stormprobe.exe ./cmd/stormprobe  # Windows
```

### Optie 4 — Docker

```bash
docker run --rm ghcr.io/umutozen/stormprobe:latest --insecure https://example.com

# Linux / macOS
docker run --rm -v $(pwd)/outputs:/app/outputs \
  ghcr.io/umutozen/stormprobe:latest --insecure --format both https://example.com

# Windows PowerShell
docker run --rm -v "${PWD}\outputs:/app/outputs" `
  ghcr.io/umutozen/stormprobe:latest --insecure --format both https://example.com
```

### Optioneel — Detectiehulpmiddelen

```bash
bash scripts/install-tools.sh
```

> Zonder Katana/Httpx `--no-discovery` of `--endpoints` gebruiken.

---

## CLI-referentie

### Gelijktijdigheid en belasting

| Vlag | Standaard | Beschrijving |
|---|---|---|
| `--concurrency-ramp` | `50` | Maximale gelijktijdigheid voor Ramp-Up-fase |
| `--concurrency-sustained` | `75` | Gelijktijdigheid voor Sustained-fase |
| `--concurrency-spike` | `250` | Maximale gelijktijdigheid voor Spike-fase |
| `--req-per-worker` | `15` | Verzoeken per worker (genegeerd als `--duration` is ingesteld) |
| `--duration` | `0` (uitgeschakeld) | Tijdsgebaseerde faseduur, bijv. `30s`, `1m` |
| `--timeout` | `10s` | Time-out per verzoek |

### Detectie

| Vlag | Standaard | Beschrijving |
|---|---|---|
| `--no-discovery` | `false` | Katana+Httpx overslaan, alleen root testen |
| `--endpoints` | `""` | Endpoints laden vanuit bestand |
| `--katana-path` | `""` | Aangepast Katana-binair pad |
| `--httpx-path` | `""` | Aangepast Httpx-binair pad |

### HTTP en beveiliging

| Vlag | Standaard | Beschrijving |
|---|---|---|
| `-H`, `--header` | — | Aangepaste header (herhaalbaar) |
| `--insecure` | `false` | TLS-certificaatverificatie overslaan |

### Uitvoer

| Vlag | Standaard | Beschrijving |
|---|---|---|
| `--format` | `both` | Rapportformaat: `json`, `html`, `both` |
| `--output` | `./outputs` | Uitvoermap voor rapporten |

### CI/CD-waarschuwingen

| Vlag | Standaard | Beschrijving |
|---|---|---|
| `--alert-p99` | `0` (uitgeschakeld) | exit 1 als P99 (ms) drempel overschrijdt |
| `--alert-error-rate` | `0` (uitgeschakeld) | exit 1 als foutpercentage (%) drempel overschrijdt |
| `--alert-rps` | `0` (uitgeschakeld) | exit 1 als req/s onder drempel valt |

---

## Testfasen

| Fase | Beschrijving |
|---|---|
| **0 — Discovery** | Katana-crawl + Httpx-probe, gededupliceerde lijst |
| **1 — Ramp-Up** | Geleidelijke toename van gelijktijdigheid: 5 → 15 → 30 → piek |
| **2 — Sustained** | 3 golven op doelgelijktijdigheid |
| **3 — Spike** | Plotselinge toename tot piek, daarna afkoeling |
| **4 — Recovery** | Gezondheidscontrole na 10s afkoeling |

---

## Uitvoer en rapporten

| Bestand | Beschrijving |
|---|---|
| `stormprobe_report_<timestamp>.json` | Alle metrieken per fase, machineleesbaar |
| `stormprobe_report_<timestamp>.html` | Zelfstandig donker dashboard met SVG-latentiediagram |

---

## Disclaimer

Dit hulpmiddel is uitsluitend bedoeld voor **geautoriseerd beveiligingstesten en prestatie-evaluatie**. U moet expliciete schriftelijke toestemming verkrijgen van de eigenaar van het doelsysteem voordat u belastingstests uitvoert. De auteurs aanvaarden **geen aansprakelijkheid** voor misbruik of schade veroorzaakt door dit hulpmiddel.

---

## Licentie

[MIT](../../LICENSE) © [Umut ÖZEN](https://github.com/umutozen)

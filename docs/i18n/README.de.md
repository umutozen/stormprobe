[English](../../README.md) · [Türkçe](README.tr.md) · [Deutsch](README.de.md) · [Francais](README.fr.md) · [Espanol](README.es.md) · [Italiano](README.it.md) · [Portugues](README.pt-BR.md) · [Russkij](README.ru.md) · [Zhongwen Jian](README.zh-CN.md) · [Zhongwen Fan](README.zh-TW.md) · [Nihongo](README.ja.md) · [Hangugeo](README.ko.md) · [Al-Arabiyya](README.ar.md) · [Hindi](README.hi.md) · [Nederlands](README.nl.md) · [Polski](README.pl.md) · [Ukrainska](README.uk.md)

<p align="center">
  <img src="../../assets/banner.png" alt="StormProbe Banner" width="100%">
</p>

<h1 align="center">StormProbe</h1>

<p align="center">
  <strong>entdecken · laden · analysieren</strong><br>
  HTTP-Lasttest-Tool auf Produktionsniveau mit autonomer Endpunkterkennung. Keine externen Abhängigkeiten.
</p>

<p align="center">
  <a href="https://github.com/umutozen/stormprobe/actions/workflows/ci.yml">
    <img src="https://github.com/umutozen/stormprobe/actions/workflows/ci.yml/badge.svg" alt="CI">
  </a>
  <a href="https://goreportcard.com/report/github.com/umutozen/stormprobe">
    <img src="https://goreportcard.com/badge/github.com/umutozen/stormprobe" alt="Go Report Card">
  </a>
  <a href="https://github.com/umutozen/stormprobe/releases">
    <img src="https://img.shields.io/github/v/release/umutozen/stormprobe" alt="Neueste Version">
  </a>
  <a href="../../LICENSE">
    <img src="https://img.shields.io/badge/license-MIT-blue.svg" alt="MIT-Lizenz">
  </a>
</p>

---

## Was ist StormProbe?

StormProbe ist ein **abhängigkeitsfreies** HTTP-Lasttест-Tool in reinem Go. Es erkennt automatisch die Endpunkte Ihrer Anwendung mit [Katana](https://github.com/projectdiscovery/katana) und [Httpx](https://github.com/projectdiscovery/httpx), führt dann einen strukturierten 4-Phasen-Lasttest durch — **Ramp-Up → Sustained → Spike → Recovery** — und erzeugt JSON- und HTML-Berichte mit P50/P95/P99-Latenzkennzahlen.

Entwickelt für:
- **CI/CD-Pipelines** — Alert-Schwellenwerte mit Exit-Codes
- **Penetrationstester** — Live-Endpunkte automatisch finden und testen
- **DevOps-Ingenieure** — Vor- und Nach-Deployment-Benchmarks
- **QA-Teams** — Performance-SLA-Validierung gegen reale Traffic-Muster

---

## Funktionen

| Funktion | Details |
|---|---|
| **Auto-Erkennung** | Katana-Crawl + Httpx-Probe, deduplizierte Live-Endpunktliste |
| **4-Phasen-Lasttest** | Ramp-Up → Sustained → Spike → Recovery |
| **Umfangreiche Metriken** | P50 / P95 / P99 Latenz, req/s, Fehlerklassifizierung pro Phase |
| **Dauermodus** | Zeitbasierte Phasen (`--duration 30s`) statt Request-Anzahl |
| **Alert-Schwellenwerte** | CI/CD-fähiges `exit 1` bei P99-, Fehlerrate- oder RPS-Verletzung |
| **Doppelte Berichte** | JSON (maschinenlesbar) + eigenständiges HTML-Dunkeldashboard |
| **Benutzerdefinierte Header** | Bearer-Tokens, Tenant-Header, Cookies — überall weitergeleitet |
| **Docker-fertig** | Einzelnes Image mit gebündeltem Katana + Httpx |
| **Plattformübergreifend** | Linux, macOS, Windows — amd64 & arm64 |
| **Keine Abhängigkeiten** | Reines Go 1.21+ stdlib, `go install` reicht |

---

## Schnellstart

```bash
# Installation (Go 1.21+ erforderlich)
go install github.com/umutozen/stormprobe/cmd/stormprobe@latest

# Grundtest — erkennt Endpunkte automatisch, führt alle 4 Phasen aus
stormprobe https://example.com

# Erkennung überspringen, nur Root-Pfad testen
stormprobe --no-discovery https://example.com

# TLS + Auth-Header + nur JSON-Bericht
stormprobe --insecure -H "Authorization: Bearer TOKEN" --format json https://api.example.com

# Zeitbasierter Test: 30 Sekunden pro Phase
stormprobe --insecure --duration 30s --no-discovery https://example.com

# CI/CD: exit 1 wenn P99 > 500ms oder Fehler > 5%
stormprobe --alert-p99 500 --alert-error-rate 5 https://example.com
```

---

## Installation

### Option 1 — `go install` *(empfohlen)*

```bash
go install github.com/umutozen/stormprobe/cmd/stormprobe@latest
stormprobe --help
```

### Option 2 — Vorgefertigte Binary *(kein Go erforderlich)*

Von der [Releases-Seite](https://github.com/umutozen/stormprobe/releases) herunterladen.

```bash
# Linux / macOS
chmod +x stormprobe
sudo mv stormprobe /usr/local/bin/

# Windows PowerShell
Move-Item stormprobe.exe C:\Windows\System32\stormprobe.exe
```

### Option 3 — Aus Quellcode bauen

```bash
git clone https://github.com/umutozen/stormprobe.git
cd stormprobe
go build -o stormprobe ./cmd/stormprobe      # Linux / macOS
go build -o stormprobe.exe ./cmd/stormprobe  # Windows
```

### Option 4 — Docker

```bash
docker run --rm ghcr.io/umutozen/stormprobe:latest --insecure https://example.com

# Linux / macOS
docker run --rm -v $(pwd)/outputs:/app/outputs \
  ghcr.io/umutozen/stormprobe:latest --insecure --format both https://example.com

# Windows PowerShell
docker run --rm -v "${PWD}\outputs:/app/outputs" `
  ghcr.io/umutozen/stormprobe:latest --insecure --format both https://example.com
```

### Optional — Erkennungstools

```bash
bash scripts/install-tools.sh
```

> Ohne Katana/Httpx `--no-discovery` oder `--endpoints` verwenden.

---

## CLI-Referenz

### Parallelität & Last

| Flag | Standard | Beschreibung |
|---|---|---|
| `--concurrency-ramp` | `50` | Peak-Parallelität für Ramp-Up-Phase |
| `--concurrency-sustained` | `75` | Parallelität für Sustained-Phase |
| `--concurrency-spike` | `250` | Peak-Parallelität für Spike-Phase |
| `--req-per-worker` | `15` | Anfragen pro Worker (ignoriert wenn `--duration` gesetzt) |
| `--duration` | `0` (deaktiviert) | Zeitbasierte Phasendauer, z.B. `30s`, `1m` |
| `--timeout` | `10s` | Timeout pro Anfrage |

### Erkennung

| Flag | Standard | Beschreibung |
|---|---|---|
| `--no-discovery` | `false` | Katana+Httpx überspringen, nur Root testen |
| `--endpoints` | `""` | Endpunkte aus Datei laden (ein Pfad pro Zeile) |
| `--katana-path` | `""` | Benutzerdefinierter Katana-Binary-Pfad |
| `--httpx-path` | `""` | Benutzerdefinierter Httpx-Binary-Pfad |

### HTTP & Sicherheit

| Flag | Standard | Beschreibung |
|---|---|---|
| `-H`, `--header` | — | Benutzerdefinierter Header (wiederholbar) |
| `--insecure` | `false` | TLS-Zertifikatsprüfung überspringen |

### Ausgabe

| Flag | Standard | Beschreibung |
|---|---|---|
| `--format` | `both` | Berichtsformat: `json`, `html`, `both` |
| `--output` | `./outputs` | Ausgabeverzeichnis für Berichte |

### CI/CD-Alerts

| Flag | Standard | Beschreibung |
|---|---|---|
| `--alert-p99` | `0` (deaktiviert) | exit 1 wenn P99 (ms) Schwellenwert überschritten |
| `--alert-error-rate` | `0` (deaktiviert) | exit 1 wenn Fehlerrate (%) Schwellenwert überschritten |
| `--alert-rps` | `0` (deaktiviert) | exit 1 wenn req/s unter Schwellenwert fällt |

---

## Testphasen

| Phase | Beschreibung |
|---|---|
| **0 — Discovery** | Katana-Crawl + Httpx-Probe, deduplizierte Endpunktliste |
| **1 — Ramp-Up** | Schrittweise Parallelitätserhöhung: 5 → 15 → 30 → Peak |
| **2 — Sustained** | 3 Wellen bei Zielparallelität, misst Steady-State-Degradation |
| **3 — Spike** | Plötzlicher Anstieg auf Peak-Parallelität, dann Abkühlung |
| **4 — Recovery** | Post-Spike-Gesundheitscheck nach 10s Abkühlung |

---

## Ausgabe & Berichte

| Datei | Beschreibung |
|---|---|
| `stormprobe_report_<timestamp>.json` | Alle Metriken pro Phase, maschinenlesbar |
| `stormprobe_report_<timestamp>.html` | Eigenständiges Dunkeldashboard mit SVG-Latenzdiagramm |

---

## Haftungsausschluss

Dieses Tool ist ausschließlich für **autorisierte Sicherheitstests und Leistungsbewertungen** bestimmt. Vor der Durchführung von Lasttests ist eine ausdrückliche schriftliche Genehmigung des Zielsystembesitzers einzuholen. Die Autoren übernehmen **keine Haftung** für Missbrauch oder Schäden.

---

## Lizenz

[MIT](../../LICENSE) © [Umut ÖZEN](https://github.com/umutozen)

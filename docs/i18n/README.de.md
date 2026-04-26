[English](../../README.md) · [Türkçe](README.tr.md) · [Deutsch](README.de.md) · [Francais](README.fr.md) · [Espanol](README.es.md) · [Italiano](README.it.md) · [Portugues](README.pt-BR.md) · [Russkij](README.ru.md) · [Zhongwen Jian](README.zh-CN.md) · [Zhongwen Fan](README.zh-TW.md) · [Nihongo](README.ja.md) · [Hangugeo](README.ko.md) · [Al-Arabiyya](README.ar.md) · [Hindi](README.hi.md) · [Nederlands](README.nl.md) · [Polski](README.pl.md) · [Ukrainska](README.uk.md)

# StormProbe

<p align="center">
  <img src="../../assets/banner.png" alt="StormProbe" width="100%">
</p>

**entdecken. belasten. analysieren.**

Produktionsreifer HTTP-Lasttester mit autonomer Endpoint-Erkennung (Endpunkterkennung). Keine Abhaengigkeiten.

[![CI](https://github.com/umutozen/stormprobe/actions/workflows/ci.yml/badge.svg)](https://github.com/umutozen/stormprobe/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/umutozen/stormprobe)](https://goreportcard.com/report/github.com/umutozen/stormprobe)

---

## Funktionen

- **Automatische Erkennung** -- Katana-Crawl + Httpx-Probe findet aktive Endpoints automatisch
- **Mehrphasen-Tests** -- Ramp-Up (stufenweiser Anstieg), Sustained (Dauerlast), Spike (Lastspitze), Recovery (Erholung)
- **Detaillierte Metriken** -- P50 / P95 / P99 Latency (Latenz), req/s, Fehlerklassifizierung
- **Duale Berichte** -- JSON (maschinenlesbar) + HTML (visuelles Dark-Theme-Dashboard)
- **Keine Abhaengigkeiten** -- Reine Go-Standardbibliothek, keine Drittanbieter-Module
- **Docker-fertig** -- Einzelner Befehl mit gebundeltem Katana + Httpx
- **Plattformuebergreifend** -- Linux-, macOS-, Windows-Binaerdateien ueber GoReleaser

## Installation

### Binary (empfohlen)
Laden Sie die neueste Version für Ihre Plattform von der [Releases-Seite](https://github.com/umutozen/stormprobe/releases) herunter, entpacken Sie das Archiv und verschieben Sie die Binärdatei in ein Verzeichnis in Ihrem PATH:

```bash
# Linux / macOS
chmod +x stormprobe
sudo mv stormprobe /usr/local/bin/
stormprobe --no-discovery https://example.com

# Windows (PowerShell)
Move-Item stormprobe.exe C:\Windows\System32\stormprobe.exe
stormprobe --no-discovery https://example.com
```

### go install (requires Go 1.21+)
```bash
go install github.com/umutozen/stormprobe/cmd/stormprobe@latest
stormprobe --no-discovery https://example.com
```

### Aus dem Quellcode bauen (erfordert Go 1.21+)
```bash
git clone https://github.com/umutozen/stormprobe.git
cd stormprobe
go build -o stormprobe ./cmd/stormprobe   # Linux / macOS
go build -o stormprobe.exe ./cmd/stormprobe   # Windows
```

### Docker

```bash
# Schnelltest, kein Bericht gespeichert
docker run --rm ghcr.io/umutozen/stormprobe:latest --insecure https://example.com

# Vollständiger Test mit Berichten in ./outputs (Linux / macOS)
docker run --rm -v $(pwd)/outputs:/app/outputs ghcr.io/umutozen/stormprobe:latest --insecure --format both https://example.com

# Vollständiger Test mit Berichten in ./outputs (Windows PowerShell)
docker run --rm -v "${PWD}\outputs:/app/outputs" ghcr.io/umutozen/stormprobe:latest --insecure --format both https://example.com

# Hochlast-Spike-Test (500 gleichzeitige Benutzer)
docker run --rm ghcr.io/umutozen/stormprobe:latest --insecure --concurrency-spike 500 https://example.com

# Discovery überspringen, nur Root-Pfad testen
docker run --rm ghcr.io/umutozen/stormprobe:latest --insecure --no-discovery https://example.com

# Benutzerdefinierte HTTP-Header (Autorisierung, benutzerdefinierter Mandant usw.)
docker run --rm ghcr.io/umutozen/stormprobe:latest --insecure -H "Authorization: Bearer TOKEN" -H "X-Tenant: acme" https://api.example.com
```

## Verwendung

```
stormprobe [flags] <ziel-url>

Flags:
  --concurrency-ramp      int       Maximale Parallelität für Ramp-Up (Standard 50)
  --concurrency-sustained int       Parallelität für Sustained-Phase (Standard 75)
  --concurrency-spike     int       Maximale Parallelität für Spike (Standard 250)
  --req-per-worker        int       Anfragen pro Worker pro Schritt (Standard 15)
  --timeout               duration  Timeout pro Anfrage (Standard 10s)
  --endpoints             string    Endpunktlistendatei (ein Pfad pro Zeile, überspringt Discovery)
  --no-discovery                    Katana+httpx überspringen, nur Root-Pfad testen
  --output                string    Ausgabeverzeichnis für Berichte (Standard ./outputs)
  --format                string    Berichtsformat: json, html, both (Standard both)
  --katana-path           string    Benutzerdefinierter Katana-Binärpfad
  --httpx-path            string    Benutzerdefinierter httpx-Binärpfad
  --insecure                        TLS-Zertifikatsprüfung überspringen
  --header, -H            string    Benutzerdefinierter HTTP-Header (wiederholbar): -H 'Authorization: Bearer TOKEN'
  --duration              duration  Dauer pro Phase (z.B. 30s, 1m). Überschreibt req-per-worker
  --alert-p99             float     Exit 1 wenn P99-Latenz Xms überschreitet
  --alert-error-rate      float     Exit 1 wenn Fehlerrate X%% überschreitet
  --alert-rps             float     Exit 1 wenn req/s unter X fällt
```

### Beispiele

```bash
stormprobe --insecure --concurrency-spike 500 --req-per-worker 20 https://example.com

stormprobe --insecure --endpoints endpoints.txt https://example.com

stormprobe --insecure --format json --output ./results https://example.com

stormprobe --insecure -H "Authorization: Bearer TOKEN" -H "X-Tenant: acme" https://api.example.com

stormprobe --insecure --duration 30s --no-discovery https://example.com

stormprobe --insecure --duration 1m --concurrency-spike 300 https://example.com

stormprobe --insecure --alert-p99 500 --alert-error-rate 5 https://example.com
```

## Testphasen

| Phase | Beschreibung |
|---|---|
| **0 -- Discovery** | Katana-Crawl + Httpx-Probe, deduplizierte Endpoint-Liste |
| **1 -- Ramp-Up** | Stufenweiser Anstieg: 5, 15, 30, 50 Virtual Users (virtuelle Benutzer) |
| **2 -- Sustained** | 3 Wellen bei Ziel-Concurrency (Gleichzeitigkeit), misst Verschlechterung |
| **3 -- Spike** | Ploetzlicher Burst zum Spitzenwert, dann Abkuehlung |
| **4 -- Recovery** | Gesundheitspruefung nach 10s Abkuehlung nach Spike |

## Ausgabe

| Datei | Beschreibung |
|---|---|
| `stormprobe_report_<timestamp>.json` | Maschinenlesbare Ergebnisse mit allen Metriken |
| `stormprobe_report_<timestamp>.html` | Visuelles Dashboard, in jedem Browser oeffnen |

## Projektstruktur

```
stormprobe/
    cmd/main.go                    CLI-Einstiegspunkt
    internal/
        config/config.go           Gemeinsame Typen, Phasengeneratoren, alert config
        alert/
            alert.go               Threshold checks, CI/CD exit code logic
        metrics/
            latency.go             Perzentilberechnung
            errors.go              Fehlerklassifizierung
        discovery/
            discovery.go           Discover() oeffentliche Schnittstelle
            katana.go              Katana-Crawler-Integration
            httpx.go               Httpx-Probe-Integration
        runner/
            phase.go               Phasenorchestrierung
            worker.go              Goroutine-Pool, Pro-Worker-RNG
        report/
            json.go                JSON-Berichtschreiber
            html.go                Eigenstaendiger HTML-Bericht
    Dockerfile                     Mehrstufig mit Katana + Httpx
    .goreleaser.yml                Cross-Compilation-Konfiguration
    Makefile                       Build-, Test-, Lint-Ziele
    config.example.yml             Phasen-Standardwerte-Referenz
```

## Anforderungen

- Go 1.21+
- [Katana](https://github.com/projectdiscovery/katana) (ProjectDiscovery) -- optional, fuer automatische Erkennung
- [Httpx](https://github.com/projectdiscovery/httpx) (ProjectDiscovery) -- optional, fuer Endpoint-Validierung

```bash
bash scripts/install-tools.sh
```

## Haftungsausschluss

Dieses Tool ist ausschliesslich fuer **autorisierte Sicherheitstests und Leistungsbewertungen** bestimmt. Sie muessen vor der Durchfuehrung von Lasttests eine ausdrueckliche schriftliche Genehmigung des Zielsystemeigentuemers einholen. Die unbefugte Verwendung dieses Tools gegen Systeme, die Sie nicht besitzen oder fuer die Sie keine Testberechtigung haben, kann gegen lokale, nationale oder internationale Gesetze verstossen. Die Autoren uebernehmen keine Haftung fuer Missbrauch oder durch dieses Tool verursachte Schaeden.

## Lizenz

MIT

## Autor

**Umut ÖZEN** — [@umutozen](https://github.com/umutozen)

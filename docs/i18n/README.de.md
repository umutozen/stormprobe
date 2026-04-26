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

## Schnellstart

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

## Verwendung

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
```

### Beispiele

```bash
stormprobe --insecure --concurrency-spike 500 --req-per-worker 20 https://example.com

stormprobe --insecure --endpoints endpoints.txt https://example.com

stormprobe --insecure --format json --output ./results https://example.com
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
        config/config.go           Gemeinsame Typen, Phasengeneratoren
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

**Umut ÖZEN** -- [@umutozen](https://github.com/umutozen)

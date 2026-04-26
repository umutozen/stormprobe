[English](../../README.md) · [Türkçe](README.tr.md) · [Deutsch](README.de.md) · [Francais](README.fr.md) · [Espanol](README.es.md) · [Italiano](README.it.md) · [Portugues](README.pt-BR.md) · [Russkij](README.ru.md) · [Zhongwen Jian](README.zh-CN.md) · [Zhongwen Fan](README.zh-TW.md) · [Nihongo](README.ja.md) · [Hangugeo](README.ko.md) · [Al-Arabiyya](README.ar.md) · [Hindi](README.hi.md) · [Nederlands](README.nl.md) · [Polski](README.pl.md) · [Ukrainska](README.uk.md)

# StormProbe

<p align="center">
  <img src="../../assets/banner.png" alt="StormProbe" width="100%">
</p>

**ontdekken. belasten. analyseren.**

Productieklare HTTP-belastingtester met autonome endpoint-ontdekking (eindpuntdetectie). Nul afhankelijkheden.

[![CI](https://github.com/umutozen/stormprobe/actions/workflows/ci.yml/badge.svg)](https://github.com/umutozen/stormprobe/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/umutozen/stormprobe)](https://goreportcard.com/report/github.com/umutozen/stormprobe)

---

## Functies

- **Automatische ontdekking** -- Katana-crawl + Httpx-probe vindt automatisch actieve endpoints
- **Meerfasentests** -- Ramp-Up (geleidelijke toename), Sustained (aanhoudend), Spike (belastingspiek), Recovery (herstel)
- **Gedetailleerde metrieken** -- P50 / P95 / P99 latency (latentie), req/s, foutclassificatie
- **Dubbel rapport** -- JSON (machineleesbaar) + HTML (visueel donker thema dashboard)
- **Nul afhankelijkheden** -- Pure Go-standaardbibliotheek, geen externe modules
- **Docker-gereed** -- Enkel commando met ingebouwde Katana + Httpx
- **Platformonafhankelijk** -- Linux-, macOS-, Windows-binaries via GoReleaser

## Snelstart

```bash
go run ./cmd https://example.com

go build -o stormprobe ./cmd
./stormprobe --insecure https://example.com
```

### Docker

```bash
# Snelle test, geen rapport opgeslagen
docker run --rm ghcr.io/umutozen/stormprobe:latest --insecure https://example.com

# Volledige test met rapporten in ./outputs (Linux / macOS)
docker run --rm -v $(pwd)/outputs:/app/outputs ghcr.io/umutozen/stormprobe:latest --insecure --format both https://example.com

# Volledige test met rapporten in ./outputs (Windows PowerShell)
docker run --rm -v "${PWD}\outputs:/app/outputs" ghcr.io/umutozen/stormprobe:latest --insecure --format both https://example.com

# Hoge belasting piektest (500 gelijktijdige gebruikers)
docker run --rm ghcr.io/umutozen/stormprobe:latest --insecure --concurrency-spike 500 https://example.com

# Detectie overslaan, alleen rootpad testen
docker run --rm ghcr.io/umutozen/stormprobe:latest --insecure --no-discovery https://example.com

# Aangepaste HTTP-headers (autorisatie, aangepaste tenant, enz.)
docker run --rm ghcr.io/umutozen/stormprobe:latest --insecure -H "Authorization: Bearer TOKEN" -H "X-Tenant: acme" https://api.example.com
```

## Gebruik

```
stormprobe [flags] <doel-url>

Flags:
  --concurrency-ramp      int       Maximale gelijktijdigheid voor Ramp-Up (standaard 50)
  --concurrency-sustained int       Gelijktijdigheid voor Sustained-fase (standaard 75)
  --concurrency-spike     int       Maximale gelijktijdigheid voor Spike (standaard 250)
  --req-per-worker        int       Verzoeken per worker per stap (standaard 15)
  --timeout               duration  Time-out per verzoek (standaard 10s)
  --endpoints             string    Endpointlijstbestand (één pad per regel)
  --no-discovery                    Katana+httpx overslaan, alleen rootpad testen
  --output                string    Uitvoermap voor rapporten (standaard ./outputs)
  --format                string    Rapportindeling: json, html, both (standaard both)
  --katana-path           string    Aangepast katana-binair pad
  --httpx-path            string    Aangepast httpx-binair pad
  --insecure                        TLS-certificaatverificatie overslaan
  --header, -H            string    Aangepaste HTTP-header (herhaalbaar): -H 'Authorization: Bearer TOKEN'
  --duration              duration  Duur per fase (bijv. 30s, 1m). Vervangt req-per-worker
  --alert-p99             float     Exit 1 als P99-latentie Xms overschrijdt
  --alert-error-rate      float     Exit 1 als foutpercentage X%% overschrijdt
  --alert-rps             float     Exit 1 als req/s onder X valt
```

### Voorbeelden

```bash
stormprobe --insecure --concurrency-spike 500 --req-per-worker 20 https://example.com

stormprobe --insecure --endpoints endpoints.txt https://example.com

stormprobe --insecure --format json --output ./results https://example.com

stormprobe --insecure -H "Authorization: Bearer TOKEN" -H "X-Tenant: acme" https://api.example.com

stormprobe --insecure --duration 30s --no-discovery https://example.com

stormprobe --insecure --duration 1m --concurrency-spike 300 https://example.com

stormprobe --insecure --alert-p99 500 --alert-error-rate 5 https://example.com
```

## Testfasen

| Fase | Beschrijving |
|---|---|
| **0 -- Discovery** | Katana-crawl + Httpx-probe, ontdubbelde endpoint-lijst |
| **1 -- Ramp-Up** | Geleidelijke toename: 5, 15, 30, 50 virtual users (virtuele gebruikers) |
| **2 -- Sustained** | 3 golven bij doel-concurrency (gelijktijdigheid), meet degradatie |
| **3 -- Spike** | Plotselinge burst naar piek, daarna afkoeling |
| **4 -- Recovery** | Gezondheidscontrole na spike na 10s afkoeling |

## Uitvoer

| Bestand | Beschrijving |
|---|---|
| `stormprobe_report_<timestamp>.json` | Machineleesbare resultaten met alle metrieken |
| `stormprobe_report_<timestamp>.html` | Visueel dashboard, open in elke browser |

## Projectstructuur

```
stormprobe/
    cmd/main.go                    CLI-ingangspunt
    internal/
        config/config.go           Gedeelde typen, fasegeneratoren, alert config
        alert/
            alert.go               Threshold checks, CI/CD exit code logic
        metrics/
            latency.go             Percentielberekening
            errors.go              Foutclassificatie
        discovery/
            discovery.go           Discover() publieke interface
            katana.go              Katana-crawler-integratie
            httpx.go               Httpx-probe-integratie
        runner/
            phase.go               Fase-orchestratie
            worker.go              Goroutine-pool, RNG per worker
        report/
            json.go                JSON-rapportschrijver
            html.go                Zelfstandig HTML-rapport
    Dockerfile                     Meerfasig met Katana + Httpx
    .goreleaser.yml                Cross-compilatieconfiguratie
    Makefile                       Build-, test-, lint-doelen
    config.example.yml             Fase-standaardwaarden-referentie
```

## Vereisten

- Go 1.21+
- [Katana](https://github.com/projectdiscovery/katana) (ProjectDiscovery) -- optioneel, voor automatische ontdekking
- [Httpx](https://github.com/projectdiscovery/httpx) (ProjectDiscovery) -- optioneel, voor endpoint-validatie

```bash
bash scripts/install-tools.sh
```

## Disclaimer

Dit hulpmiddel is uitsluitend bedoeld voor **geautoriseerde beveiligingstests en prestatie-evaluatie**. U moet expliciete schriftelijke toestemming verkrijgen van de eigenaar van het doelsysteem voordat u belastingtests uitvoert. Ongeautoriseerd gebruik van dit hulpmiddel tegen systemen die u niet bezit of waarvoor u geen testtoestemming heeft, kan in strijd zijn met lokale, nationale of internationale wetgeving. De auteurs aanvaarden geen aansprakelijkheid voor misbruik of schade veroorzaakt door dit hulpmiddel.

## Licentie

MIT

## Auteur

**Umut ÖZEN** — [@umutozen](https://github.com/umutozen)

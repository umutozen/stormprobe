[English](../../README.md) · [Türkçe](README.tr.md) · [Deutsch](README.de.md) · [Francais](README.fr.md) · [Espanol](README.es.md) · [Italiano](README.it.md) · [Portugues](README.pt-BR.md) · [Russkij](README.ru.md) · [Zhongwen Jian](README.zh-CN.md) · [Zhongwen Fan](README.zh-TW.md) · [Nihongo](README.ja.md) · [Hangugeo](README.ko.md) · [Al-Arabiyya](README.ar.md) · [Hindi](README.hi.md) · [Nederlands](README.nl.md) · [Polski](README.pl.md) · [Ukrainska](README.uk.md)

# StormProbe

<p align="center">
  <img src="../../assets/banner.png" alt="StormProbe" width="100%">
</p>

**odkryj. obciaz. analizuj.**

Produkcyjne narzedzie do testow obciazeniowych HTTP z autonomicznym wykrywaniem endpoint (punktow koncowych). Zero zaleznosci.

[![CI](https://github.com/umutozen/stormprobe/actions/workflows/ci.yml/badge.svg)](https://github.com/umutozen/stormprobe/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/umutozen/stormprobe)](https://goreportcard.com/report/github.com/umutozen/stormprobe)

---

## Funkcje

- **Automatyczne wykrywanie** -- Skanowanie Katana + sonda Httpx automatycznie znajduje aktywne endpoint'y
- **Testy wielofazowe** -- Ramp-Up (stopniowy wzrost), Sustained (ciagle), Spike (szczyt obciazenia), Recovery (odzyskiwanie)
- **Szczegolowe metryki** -- P50 / P95 / P99 latency (opoznienie), req/s, klasyfikacja bledow
- **Podwojny raport** -- JSON (odczytywalny maszynowo) + HTML (wizualny dashboard w ciemnym motywie)
- **Zero zaleznosci** -- Czysta standardowa biblioteka Go, bez modulow zewnetrznych
- **Gotowy na Docker** -- Pojedyncze polecenie z wbudowanym Katana + Httpx
- **Wieloplatformowy** -- Pliki binarne Linux, macOS, Windows przez GoReleaser

## Szybki start

```bash
go run ./cmd https://example.com

go build -o stormprobe ./cmd
./stormprobe --insecure https://example.com
```

### Docker

```bash
# Szybki test, brak zapisanego raportu
docker run --rm ghcr.io/umutozen/stormprobe:latest --insecure https://example.com

# Pełny test z zapisaniem raportów do ./outputs (Linux / macOS)
docker run --rm -v $(pwd)/outputs:/app/outputs ghcr.io/umutozen/stormprobe:latest --insecure --format both https://example.com

# Pełny test z zapisaniem raportów do ./outputs (Windows PowerShell)
docker run --rm -v "${PWD}\outputs:/app/outputs" ghcr.io/umutozen/stormprobe:latest --insecure --format both https://example.com

# Test przeciążeniowy (500 równoczesnych użytkowników)
docker run --rm ghcr.io/umutozen/stormprobe:latest --insecure --concurrency-spike 500 https://example.com

# Pomiń wykrywanie, testuj tylko ścieżkę główną
docker run --rm ghcr.io/umutozen/stormprobe:latest --insecure --no-discovery https://example.com

# Niestandardowe nagłówki HTTP (autoryzacja, niestandardowy tenant itp.)
docker run --rm ghcr.io/umutozen/stormprobe:latest --insecure -H "Authorization: Bearer TOKEN" -H "X-Tenant: acme" https://api.example.com
```

## Uzycie

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
  --header, -H            string    Custom HTTP header (repeatable): -H 'Authorization: Bearer TOKEN'
  --duration              duration  Per-phase duration (e.g. 30s, 1m). Overrides req-per-worker when set
  --alert-p99             float     Fail (exit 1) if P99 latency exceeds Xms in any phase
  --alert-error-rate      float     Fail (exit 1) if error rate exceeds X%% in any phase
  --alert-rps             float     Fail (exit 1) if req/s falls below X in any phase
```

### Przyklady

```bash
stormprobe --insecure --concurrency-spike 500 --req-per-worker 20 https://example.com

stormprobe --insecure --endpoints endpoints.txt https://example.com

stormprobe --insecure --format json --output ./results https://example.com

stormprobe --insecure -H "Authorization: Bearer TOKEN" -H "X-Tenant: acme" https://api.example.com

stormprobe --insecure --duration 30s --no-discovery https://example.com

stormprobe --insecure --duration 1m --concurrency-spike 300 https://example.com

stormprobe --insecure --alert-p99 500 --alert-error-rate 5 https://example.com
```

## Fazy testow

| Faza | Opis |
|---|---|
| **0 -- Discovery** | Skanowanie Katana + sonda Httpx, lista endpoint'ow po deduplikacji |
| **1 -- Ramp-Up** | Stopniowy wzrost: 5, 15, 30, 50 virtual users (uzytkownicy wirtualni) |
| **2 -- Sustained** | 3 fale przy docelowej concurrency (wspolbieznosci), mierzy degradacje |
| **3 -- Spike** | Nagle wyrzut do szczytu, nastepnie chlodzenie |
| **4 -- Recovery** | Kontrola zdrowia po spike po 10s chlodzenia |

## Wynik

| Plik | Opis |
|---|---|
| `stormprobe_report_<timestamp>.json` | Wyniki odczytywalne maszynowo ze wszystkimi metrykami |
| `stormprobe_report_<timestamp>.html` | Wizualny dashboard, otworz w dowolnej przegladarce |

## Struktura projektu

```
stormprobe/
    cmd/main.go                    Punkt wejscia CLI
    internal/
        config/config.go           Typy wspoldzielone, generatory faz
        metrics/
            latency.go             Obliczanie percentyli
            errors.go              Klasyfikacja bledow
        discovery/
            discovery.go           Interfejs publiczny Discover()
            katana.go              Integracja crawlera Katana
            httpx.go               Integracja sondy Httpx
        runner/
            phase.go               Orkiestracja faz
            worker.go              Pula goroutine, RNG na worker
        report/
            json.go                Generator raportow JSON
            html.go                Samodzielny raport HTML
    Dockerfile                     Wieloetapowy z Katana + Httpx
    .goreleaser.yml                Konfiguracja kompilacji krzyzowej
    Makefile                       Cele build, test, lint
    config.example.yml             Referencja domyslnych wartosci faz
```

## Wymagania

- Go 1.21+
- [Katana](https://github.com/projectdiscovery/katana) (ProjectDiscovery) -- opcjonalnie, do automatycznego wykrywania
- [Httpx](https://github.com/projectdiscovery/httpx) (ProjectDiscovery) -- opcjonalnie, do walidacji endpoint'ow

```bash
bash scripts/install-tools.sh
```

## Zastrzezenie

To narzedzie jest przeznaczone wylacznie do **autoryzowanych testow bezpieczenstwa i oceny wydajnosci**. Przed uruchomieniem jakichkolwiek testow obciazeniowych nalezy uzyskac wyrazna pisemna zgode od wlasciciela systemu docelowego. Nieautoryzowane uzycie tego narzedzia przeciwko systemom, ktorych nie jestes wlascicielem lub na ktore nie masz pozwolenia na testowanie, moze naruszac przepisy lokalne, krajowe lub miedzynarodowe. Autorzy nie ponosz zadnej odpowiedzialnosci za niewlasciwe uzycie lub szkody spowodowane przez to narzedzie.

## Licencja

MIT

## Autor

**Umut ÖZEN** -- [@umutozen](https://github.com/umutozen)

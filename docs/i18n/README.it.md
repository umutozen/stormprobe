[English](../../README.md) · [Türkçe](README.tr.md) · [Deutsch](README.de.md) · [Francais](README.fr.md) · [Espanol](README.es.md) · [Italiano](README.it.md) · [Portugues](README.pt-BR.md) · [Russkij](README.ru.md) · [Zhongwen Jian](README.zh-CN.md) · [Zhongwen Fan](README.zh-TW.md) · [Nihongo](README.ja.md) · [Hangugeo](README.ko.md) · [Al-Arabiyya](README.ar.md) · [Hindi](README.hi.md) · [Nederlands](README.nl.md) · [Polski](README.pl.md) · [Ukrainska](README.uk.md)

# StormProbe

<p align="center">
  <img src="../../assets/banner.png" alt="StormProbe" width="100%">
</p>

**scopri. carica. analizza.**

Strumento di test di carico HTTP di livello produzione con scoperta autonoma degli endpoint (punti di accesso). Zero dipendenze.

[![CI](https://github.com/umutozen/stormprobe/actions/workflows/ci.yml/badge.svg)](https://github.com/umutozen/stormprobe/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/umutozen/stormprobe)](https://goreportcard.com/report/github.com/umutozen/stormprobe)

---

## Funzionalita

- **Scoperta automatica** -- Scansione Katana + sonda Httpx trova automaticamente gli endpoint attivi
- **Test multi-fase** -- Ramp-Up (aumento graduale), Sustained (sostenuto), Spike (picco di carico), Recovery (recupero)
- **Metriche dettagliate** -- P50 / P95 / P99 latency (latenza), req/s, classificazione degli errori
- **Doppio report** -- JSON (leggibile da macchina) + HTML (dashboard visuale tema scuro)
- **Zero dipendenze** -- Libreria standard Go pura, nessun modulo di terze parti
- **Pronto per Docker** -- Singolo comando con Katana + Httpx integrati
- **Multi-piattaforma** -- Binari Linux, macOS, Windows tramite GoReleaser

## Avvio rapido

```bash
go run ./cmd https://example.com

go build -o stormprobe ./cmd
./stormprobe --insecure https://example.com
```

### Docker

```bash
# Test rapido, nessun rapporto salvato
docker run --rm ghcr.io/umutozen/stormprobe:latest --insecure https://example.com

# Test completo con report in ./outputs (Linux / macOS)
docker run --rm -v $(pwd)/outputs:/app/outputs ghcr.io/umutozen/stormprobe:latest --insecure --format both https://example.com

# Test completo con report in ./outputs (Windows PowerShell)
docker run --rm -v "${PWD}\outputs:/app/outputs" ghcr.io/umutozen/stormprobe:latest --insecure --format both https://example.com

# Test di picco ad alto carico (500 utenti simultanei)
docker run --rm ghcr.io/umutozen/stormprobe:latest --insecure --concurrency-spike 500 https://example.com

# Salta il discovery, testa solo il percorso radice
docker run --rm ghcr.io/umutozen/stormprobe:latest --insecure --no-discovery https://example.com

# Intestazioni HTTP personalizzate (autorizzazione, tenant personalizzato, ecc.)
docker run --rm ghcr.io/umutozen/stormprobe:latest --insecure -H "Authorization: Bearer TOKEN" -H "X-Tenant: acme" https://api.example.com
```

## Utilizzo

```
stormprobe [flags] <url-destinazione>

Flags:
  --concurrency-ramp      int       Concorrenza massima per Ramp-Up (predefinito 50)
  --concurrency-sustained int       Concorrenza per la fase Sustained (predefinito 75)
  --concurrency-spike     int       Concorrenza massima per Spike (predefinito 250)
  --req-per-worker        int       Richieste per worker per passo (predefinito 15)
  --timeout               duration  Timeout per richiesta (predefinito 10s)
  --endpoints             string    File di elenco endpoint (un percorso per riga)
  --no-discovery                    Salta katana+httpx, testa solo il percorso radice
  --output                string    Directory di output per i report (predefinito ./outputs)
  --format                string    Formato del report: json, html, both (predefinito both)
  --katana-path           string    Percorso binario katana personalizzato
  --httpx-path            string    Percorso binario httpx personalizzato
  --insecure                        Salta la verifica del certificato TLS
  --header, -H            string    Header HTTP personalizzato (ripetibile): -H 'Authorization: Bearer TOKEN'
  --duration              duration  Durata per fase (es. 30s, 1m). Sostituisce req-per-worker
  --alert-p99             float     Exit 1 se la latenza P99 supera Xms
  --alert-error-rate      float     Exit 1 se il tasso di errore supera X%%
  --alert-rps             float     Exit 1 se req/s scende sotto X
```

### Esempi

```bash
stormprobe --insecure --concurrency-spike 500 --req-per-worker 20 https://example.com

stormprobe --insecure --endpoints endpoints.txt https://example.com

stormprobe --insecure --format json --output ./results https://example.com

stormprobe --insecure -H "Authorization: Bearer TOKEN" -H "X-Tenant: acme" https://api.example.com

stormprobe --insecure --duration 30s --no-discovery https://example.com

stormprobe --insecure --duration 1m --concurrency-spike 300 https://example.com

stormprobe --insecure --alert-p99 500 --alert-error-rate 5 https://example.com
```

## Fasi di test

| Fase | Descrizione |
|---|---|
| **0 -- Discovery** | Scansione Katana + sonda Httpx, lista endpoint deduplicata |
| **1 -- Ramp-Up** | Aumento graduale: 5, 15, 30, 50 virtual users (utenti virtuali) |
| **2 -- Sustained** | 3 ondate alla concurrency (concorrenza) obiettivo, misura il degrado |
| **3 -- Spike** | Burst improvviso al picco, poi raffreddamento |
| **4 -- Recovery** | Controllo di salute post-spike dopo 10s di raffreddamento |

## Output

| File | Descrizione |
|---|---|
| `stormprobe_report_<timestamp>.json` | Risultati leggibili da macchina con tutte le metriche |
| `stormprobe_report_<timestamp>.html` | Dashboard visuale, aprire in qualsiasi browser |

## Struttura del progetto

```
stormprobe/
    cmd/main.go                    Punto di ingresso CLI
    internal/
        config/config.go           Tipi condivisi, generatori di fasi, alert config
        alert/
            alert.go               Threshold checks, CI/CD exit code logic
        metrics/
            latency.go             Calcolo dei percentili
            errors.go              Classificazione degli errori
        discovery/
            discovery.go           Interfaccia pubblica Discover()
            katana.go              Integrazione crawler Katana
            httpx.go               Integrazione sonda Httpx
        runner/
            phase.go               Orchestrazione delle fasi
            worker.go              Pool di goroutine, RNG per worker
        report/
            json.go                Scrittore di report JSON
            html.go                Report HTML autonomo
    Dockerfile                     Multi-stage con Katana + Httpx
    .goreleaser.yml                Configurazione compilazione incrociata
    Makefile                       Obiettivi build, test, lint
    config.example.yml             Riferimento valori predefiniti delle fasi
```

## Requisiti

- Go 1.21+
- [Katana](https://github.com/projectdiscovery/katana) (ProjectDiscovery) -- opzionale, per la scoperta automatica
- [Httpx](https://github.com/projectdiscovery/httpx) (ProjectDiscovery) -- opzionale, per la validazione degli endpoint

```bash
bash scripts/install-tools.sh
```

## Dichiarazione di non responsabilita

Questo strumento e destinato esclusivamente a **test di sicurezza autorizzati e valutazione delle prestazioni**. E necessario ottenere un'autorizzazione scritta esplicita dal proprietario del sistema di destinazione prima di eseguire qualsiasi test di carico. L'uso non autorizzato di questo strumento contro sistemi che non si possiedono o per i quali non si dispone dell'autorizzazione al test puo violare le leggi locali, nazionali o internazionali. Gli autori non si assumono alcuna responsabilita per uso improprio o danni causati da questo strumento.

## Licenza

MIT

## Autore

**Umut ÖZEN** -- [@umutozen](https://github.com/umutozen)

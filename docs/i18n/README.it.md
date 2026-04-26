[English](../../README.md) · [Türkçe](README.tr.md) · [Deutsch](README.de.md) · [Francais](README.fr.md) · [Espanol](README.es.md) · [Italiano](README.it.md) · [Portugues](README.pt-BR.md) · [Russkij](README.ru.md) · [Zhongwen Jian](README.zh-CN.md) · [Zhongwen Fan](README.zh-TW.md) · [Nihongo](README.ja.md) · [Hangugeo](README.ko.md) · [Al-Arabiyya](README.ar.md) · [Hindi](README.hi.md) · [Nederlands](README.nl.md) · [Polski](README.pl.md) · [Ukrainska](README.uk.md)

<p align="center">
  <img src="../../assets/banner.png" alt="StormProbe Banner" width="100%">
</p>

<h1 align="center">StormProbe</h1>

<p align="center">
  <strong>scoprire · caricare · analizzare</strong><br>
  Strumento di test del carico HTTP di livello produttivo con scoperta autonoma degli endpoint. Zero dipendenze esterne.
</p>

<p align="center">
  <a href="https://github.com/umutozen/stormprobe/actions/workflows/ci.yml">
    <img src="https://github.com/umutozen/stormprobe/actions/workflows/ci.yml/badge.svg" alt="CI">
  </a>
  <a href="https://goreportcard.com/report/github.com/umutozen/stormprobe">
    <img src="https://goreportcard.com/badge/github.com/umutozen/stormprobe" alt="Go Report Card">
  </a>
  <a href="https://github.com/umutozen/stormprobe/releases">
    <img src="https://img.shields.io/github/v/release/umutozen/stormprobe" alt="Ultima versione">
  </a>
  <a href="../../LICENSE">
    <img src="https://img.shields.io/badge/license-MIT-blue.svg" alt="Licenza MIT">
  </a>
</p>

---

## Cos'è StormProbe?

StormProbe è uno strumento di test del carico HTTP **senza dipendenze** scritto in Go puro. Scopre automaticamente gli endpoint dell'applicazione con [Katana](https://github.com/projectdiscovery/katana) e [Httpx](https://github.com/projectdiscovery/httpx), poi esegue un test di carico strutturato in 4 fasi — **Ramp-Up → Sustained → Spike → Recovery** — e produce report JSON e HTML con metriche di latenza P50/P95/P99.

Progettato per:
- **Pipeline CI/CD** — soglie di allerta con codici di uscita
- **Penetration tester** — scoperta automatica degli endpoint attivi
- **Ingegneri DevOps** — benchmark prima e dopo il deployment
- **Team QA** — validazione degli SLA di performance

---

## Funzionalita

| Funzionalita | Dettagli |
|---|---|
| **Scoperta automatica** | Crawl Katana + sonda Httpx, lista endpoint deduplicata |
| **Test in 4 fasi** | Ramp-Up → Sustained → Spike → Recovery |
| **Metriche ricche** | Latenza P50/P95/P99, req/s, classificazione errori per fase |
| **Modalita durata** | Fasi basate sul tempo (`--duration 30s`) invece del conteggio richieste |
| **Soglie di allerta** | `exit 1` compatibile CI/CD su violazione P99, tasso di errore o RPS |
| **Report doppi** | JSON (leggibile da macchina) + dashboard HTML scuro autonomo |
| **Header personalizzati** | Token Bearer, header tenant, cookie — propagati ovunque |
| **Pronto per Docker** | Immagine singola con Katana + Httpx inclusi |
| **Multipiattaforma** | Linux, macOS, Windows — amd64 & arm64 |
| **Zero dipendenze** | Go 1.21+ stdlib puro, `go install` basta |

---

## Avvio rapido

```bash
# Installazione (richiede Go 1.21+)
go install github.com/umutozen/stormprobe/cmd/stormprobe@latest

# Test di base — scoperta automatica, 4 fasi
stormprobe https://example.com

# Salta la scoperta, testa solo il percorso radice
stormprobe --no-discovery https://example.com

# TLS + header auth + solo report JSON
stormprobe --insecure -H "Authorization: Bearer TOKEN" --format json https://api.example.com

# Test basato sulla durata: 30 secondi per fase
stormprobe --insecure --duration 30s --no-discovery https://example.com

# CI/CD: exit 1 se P99 > 500ms o errori > 5%
stormprobe --alert-p99 500 --alert-error-rate 5 https://example.com
```

---

## Installazione

### Opzione 1 — `go install` *(consigliato)*

```bash
go install github.com/umutozen/stormprobe/cmd/stormprobe@latest
stormprobe --help
```

### Opzione 2 — Binario precompilato *(senza Go)*

Scaricare dalla [pagina Releases](https://github.com/umutozen/stormprobe/releases).

```bash
# Linux / macOS
chmod +x stormprobe
sudo mv stormprobe /usr/local/bin/

# Windows PowerShell
Move-Item stormprobe.exe C:\Windows\System32\stormprobe.exe
```

### Opzione 3 — Compilare dai sorgenti

```bash
git clone https://github.com/umutozen/stormprobe.git
cd stormprobe
go build -o stormprobe ./cmd/stormprobe      # Linux / macOS
go build -o stormprobe.exe ./cmd/stormprobe  # Windows
```

### Opzione 4 — Docker

```bash
docker run --rm ghcr.io/umutozen/stormprobe:latest --insecure https://example.com

# Linux / macOS
docker run --rm -v $(pwd)/outputs:/app/outputs \
  ghcr.io/umutozen/stormprobe:latest --insecure --format both https://example.com

# Windows PowerShell
docker run --rm -v "${PWD}\outputs:/app/outputs" `
  ghcr.io/umutozen/stormprobe:latest --insecure --format both https://example.com
```

### Opzionale — Strumenti di scoperta

```bash
bash scripts/install-tools.sh
```

> Senza Katana/Httpx usare `--no-discovery` o `--endpoints`.

---

## Riferimento CLI

### Concorrenza e Carico

| Flag | Predefinito | Descrizione |
|---|---|---|
| `--concurrency-ramp` | `50` | Concorrenza massima per la fase Ramp-Up |
| `--concurrency-sustained` | `75` | Concorrenza per la fase Sustained |
| `--concurrency-spike` | `250` | Concorrenza massima per la fase Spike |
| `--req-per-worker` | `15` | Richieste per worker (ignorato se `--duration` e impostato) |
| `--duration` | `0` (disabilitato) | Durata fase basata sul tempo, es. `30s`, `1m` |
| `--timeout` | `10s` | Timeout per richiesta |

### Scoperta

| Flag | Predefinito | Descrizione |
|---|---|---|
| `--no-discovery` | `false` | Salta Katana+Httpx, testa solo la radice |
| `--endpoints` | `""` | Carica endpoint da file |
| `--katana-path` | `""` | Percorso binario Katana personalizzato |
| `--httpx-path` | `""` | Percorso binario Httpx personalizzato |

### HTTP e Sicurezza

| Flag | Predefinito | Descrizione |
|---|---|---|
| `-H`, `--header` | — | Header personalizzato (ripetibile) |
| `--insecure` | `false` | Salta la verifica del certificato TLS |

### Output

| Flag | Predefinito | Descrizione |
|---|---|---|
| `--format` | `both` | Formato report: `json`, `html`, `both` |
| `--output` | `./outputs` | Directory di output per i report |

### Alert CI/CD

| Flag | Predefinito | Descrizione |
|---|---|---|
| `--alert-p99` | `0` (disabilitato) | exit 1 se P99 (ms) supera la soglia |
| `--alert-error-rate` | `0` (disabilitato) | exit 1 se il tasso di errore (%) supera la soglia |
| `--alert-rps` | `0` (disabilitato) | exit 1 se req/s scende sotto la soglia |

---

## Fasi di test

| Fase | Descrizione |
|---|---|
| **0 — Discovery** | Crawl Katana + sonda Httpx, lista deduplicata |
| **1 — Ramp-Up** | Aumento graduale della concorrenza: 5 → 15 → 30 → picco |
| **2 — Sustained** | 3 ondate alla concorrenza target |
| **3 — Spike** | Aumento improvviso al picco, poi raffreddamento |
| **4 — Recovery** | Controllo salute post-spike dopo 10s di raffreddamento |

---

## Output e Report

| File | Descrizione |
|---|---|
| `stormprobe_report_<timestamp>.json` | Tutte le metriche per fase, leggibile da macchina |
| `stormprobe_report_<timestamp>.html` | Dashboard scuro autonomo con grafico SVG latenza |

---

## Disclaimer

Questo strumento e destinato **esclusivamente a test di sicurezza autorizzati e valutazione delle prestazioni**. E necessario ottenere un esplicito permesso scritto dal proprietario del sistema di destinazione prima di eseguire test di carico. Gli autori non assumono **alcuna responsabilita** per uso improprio o danni causati da questo strumento.

---

## Licenza

[MIT](../../LICENSE) © [Umut ÖZEN](https://github.com/umutozen)

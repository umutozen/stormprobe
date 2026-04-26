[English](../../README.md) · [Türkçe](README.tr.md) · [Deutsch](README.de.md) · [Francais](README.fr.md) · [Espanol](README.es.md) · [Italiano](README.it.md) · [Portugues](README.pt-BR.md) · [Russkij](README.ru.md) · [Zhongwen Jian](README.zh-CN.md) · [Zhongwen Fan](README.zh-TW.md) · [Nihongo](README.ja.md) · [Hangugeo](README.ko.md) · [Al-Arabiyya](README.ar.md) · [Hindi](README.hi.md) · [Nederlands](README.nl.md) · [Polski](README.pl.md) · [Ukrainska](README.uk.md)

# StormProbe

<p align="center">
  <img src="../../assets/banner.png" alt="StormProbe" width="100%">
</p>

**vyyavyty. navantazhyty. proanalizuvaty.**

Instrument navantazhuval'noho testuvannya HTTP vyrobnychoho rivnya z avtonomnym vyyavlennyam endpoint (kintsevykh tochok). Nul' zalezhnostey.

[![CI](https://github.com/umutozen/stormprobe/actions/workflows/ci.yml/badge.svg)](https://github.com/umutozen/stormprobe/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/umutozen/stormprobe)](https://goreportcard.com/report/github.com/umutozen/stormprobe)

---

## Mozhlyvosti

- **Avtomatychne vyyavlennya** -- Skanuvannnya Katana + zond Httpx avtomatychno znakhodyat aktyvni endpoint'y
- **Bahatofazne testuvannya** -- Ramp-Up (postupove zbil'shennya), Sustained (stale), Spike (shpychka navantazhennya), Recovery (vidnovlennya)
- **Detalizovani metryky** -- P50 / P95 / P99 latency (zatrymka), req/s, klasyfikatsiya pomylok
- **Podviynyy zvit** -- JSON (mashynochytnyy) + HTML (vizual'na panel' z temnoyu temoyu)
- **Nul' zalezhnostey** -- Chysta standartna biblioteka Go, bez moduliv tretikh storin
- **Hotovyy dlya Docker** -- Odna komanda z vbudovanymy Katana + Httpx
- **Krosplatformennyy** -- Binarni fayly Linux, macOS, Windows cherez GoReleaser

## Shvydkyy start

```bash
go run ./cmd https://example.com

go build -o stormprobe ./cmd
./stormprobe --insecure https://example.com
```

### Docker

```bash
# Швидкий тест, звіт не зберігається
docker run --rm ghcr.io/umutozen/stormprobe:latest --insecure https://example.com

# Повний тест зі збереженням звітів у ./outputs (Linux / macOS)
docker run --rm -v $(pwd)/outputs:/app/outputs ghcr.io/umutozen/stormprobe:latest --insecure --format both https://example.com

# Повний тест зі збереженням звітів у ./outputs (Windows PowerShell)
docker run --rm -v "${PWD}\outputs:/app/outputs" ghcr.io/umutozen/stormprobe:latest --insecure --format both https://example.com

# Навантажувальний тест (500 одночасних користувачів)
docker run --rm ghcr.io/umutozen/stormprobe:latest --insecure --concurrency-spike 500 https://example.com

# Пропустити виявлення, тестувати лише кореневий шлях
docker run --rm ghcr.io/umutozen/stormprobe:latest --insecure --no-discovery https://example.com

# Налаштовані HTTP-заголовки (авторизація, власний тенант тощо)
docker run --rm ghcr.io/umutozen/stormprobe:latest --insecure -H "Authorization: Bearer TOKEN" -H "X-Tenant: acme" https://api.example.com
```

## Vykorystannya

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
```

### Pryklady

```bash
stormprobe --insecure --concurrency-spike 500 --req-per-worker 20 https://example.com

stormprobe --insecure --endpoints endpoints.txt https://example.com

stormprobe --insecure --format json --output ./results https://example.com

stormprobe --insecure -H "Authorization: Bearer TOKEN" -H "X-Tenant: acme" https://api.example.com
```

## Fazy testuvannya

| Faza | Opys |
|---|---|
| **0 -- Discovery** | Skanuvannnya Katana + zond Httpx, deduplikovanyy spysok endpoint'iv |
| **1 -- Ramp-Up** | Postupove zbil'shennya: 5, 15, 30, 50 virtual users (virtual'ni korystuvachi) |
| **2 -- Sustained** | 3 khvyli pry tsil'oviy concurrency (odnochasnosti), vymyryuye dehradatsiyu |
| **3 -- Spike** | Raptovyy splezk do piku, potim okholodzhennya |
| **4 -- Recovery** | Perevirka zdorov'ya pislya spike cherez 10s okholodzhennya |

## Vykhid

| Fayl | Opys |
|---|---|
| `stormprobe_report_<timestamp>.json` | Mashynochytni rezul'taty z usima metrykamy |
| `stormprobe_report_<timestamp>.html` | Vizual'na panel', vidkryyte v bud'-yakomu brauzeri |

## Struktura proyektu

```
stormprobe/
    cmd/main.go                    Tochka vkhodu CLI
    internal/
        config/config.go           Spil'ni typy, heneratory faz
        metrics/
            latency.go             Obchyslennya protsentyliv
            errors.go              Klasyfikatsiya pomylok
        discovery/
            discovery.go           Publichnyy interfeys Discover()
            katana.go              Intehratsiya kraulera Katana
            httpx.go               Intehratsiya zonda Httpx
        runner/
            phase.go               Orkestratsiya faz
            worker.go              Pul horutyn, RNG na worker
        report/
            json.go                Generator zvitiv JSON
            html.go                Avtonomnyy HTML-zvit
    Dockerfile                     Bahatostupnevyy z Katana + Httpx
    .goreleaser.yml                Konfihuratsiya kros-kompilyatsiyi
    Makefile                       Tsili build, test, lint
    config.example.yml             Dovidka za zamovchuvannyyam faz
```

## Vymohy

- Go 1.21+
- [Katana](https://github.com/projectdiscovery/katana) (ProjectDiscovery) -- neobov'yazkovo, dlya avtomatychnoho vyyavlennya
- [Httpx](https://github.com/projectdiscovery/httpx) (ProjectDiscovery) -- neobov'yazkovo, dlya validatsiyi endpoint'iv

```bash
bash scripts/install-tools.sh
```

## Vidmova vid vidpovidal'nosti

Tsey instrument pryznachenyy vyklyuchno dlya **avtoryzovanoho testuvannya bezpeky ta otsinky produktyvnosti**. Vy povynni otrymaty yavnyy pys'movyy dozvil vid vlasnyky tsil'ovoyi systemy pered zapuskom bud'-yakykh testiv navantazhennya. Nesanktsionovane vykorystannya ts'oho instrumentu proty system, yakymy vy ne volodiyete abo ne mayete dozvolu na testuvannya, mozhe porushuvsty mistsevi, natsional'ni abo mizhnarodni zakony. Avtory ne nesut' vidpovidal'nosti za nepravyl'ne vykorystannya abo zbytky, zavdani tsym instrumentom.

## Litsenziya

MIT

## Avtor

**Umut ÖZEN** -- [@umutozen](https://github.com/umutozen)

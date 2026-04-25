[English](../../README.md) · [Turkce](README.tr.md) · [Deutsch](README.de.md) · [Francais](README.fr.md) · [Espanol](README.es.md) · [Italiano](README.it.md) · [Portugues](README.pt-BR.md) · [Russkij](README.ru.md) · [Zhongwen Jian](README.zh-CN.md) · [Zhongwen Fan](README.zh-TW.md) · [Nihongo](README.ja.md) · [Hangugeo](README.ko.md) · [Al-Arabiyya](README.ar.md) · [Hindi](README.hi.md) · [Nederlands](README.nl.md) · [Polski](README.pl.md) · [Ukrainska](README.uk.md)

# StormProbe

<p align="center">
  <img src="../../assets/banner.png" alt="StormProbe" width="100%">
</p>

**hakken suru. fuka wo kakeru. bunseki suru.**

Jiritsuteki na endpoint (endopointo) hakken kinoo wo sonaeta, honban-kyuu no HTTP fuka tesuto tsuuru. Zero izon.

[![CI](https://github.com/umutozen/stormprobe/actions/workflows/ci.yml/badge.svg)](https://github.com/umutozen/stormprobe/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/umutozen/stormprobe)](https://goreportcard.com/report/github.com/umutozen/stormprobe)

---

## Kinoo

- **Jidoo hakken** -- Katana kurooringu + Httpx puroobu de akutibu na endpoint wo jidoo de ni hakken
- **Tafeezu tesuto** -- Ramp-Up (dankai-teki zooka), Sustained (jizoku), Spike (fuka piiku), Recovery (kaifuku)
- **Shousai na metorikusu** -- P50 / P95 / P99 latency (reitenshi), req/s, eraa bunrui
- **Dyuaru repooto** -- JSON (kikai-kanou) + HTML (bijuaru daaku teema dasshubodo)
- **Zero izon** -- Junsui na Go hyoojun raiburarii, saadopaatii mojuuru nashi
- **Docker taioo** -- Katana + Httpx dookonde tan'itsu komando
- **Kurosu purattofoomu** -- GoReleaser ni yoru Linux, macOS, Windows bainarii

## Kuikku sutaato

```bash
go run ./cmd https://example.com

go build -o stormprobe ./cmd
./stormprobe --insecure https://example.com
```

### Docker

```bash
docker run --rm -v $(pwd)/outputs:/app/outputs \
  ghcr.io/umutozen/stormprobe --insecure https://example.com
```

## Shiyoo-hoo

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

### Rei

```bash
stormprobe --insecure --concurrency-spike 500 --req-per-worker 20 https://example.com

stormprobe --insecure --endpoints endpoints.txt https://example.com

stormprobe --insecure --format json --output ./results https://example.com
```

## Tesuto feezu

| Feezu | Setsumei |
|---|---|
| **0 -- Discovery** | Katana kurooringu + Httpx puroobu, jyuufuku haijo sareta endpoint risuto |
| **1 -- Ramp-Up** | Dankai-teki zooka: 5, 15, 30, 50 virtual users (kasoou yuuzaa) |
| **2 -- Sustained** | Taagetoo concurrency (doojisei) de 3 ha, retoka wo sokutei |
| **3 -- Spike** | Piiku he no totsuzen no baasuto, sono go reikyaku |
| **4 -- Recovery** | Spike go 10-byoo reikyaku go no kenzen-sei chekku |

## Shutsuryoku

| Fairu | Setsumei |
|---|---|
| `stormprobe_report_<timestamp>.json` | Subete no metorikusu wo fukumu kikai-kadoku na kekka |
| `stormprobe_report_<timestamp>.html` | Bijuaru dasshubodo, nin'i no burauza de hiraku |

## Purojekuto koozoo

```
stormprobe/
    cmd/main.go                    CLI entoriipointo
    internal/
        config/config.go           Kyooyuu-gata, feezu jenereetaa
        metrics/
            latency.go             Paasenta-iru keisan
            errors.go              Eraa bunrui
        discovery/
            discovery.go           Discover() paburikku intaafeisu
            katana.go              Katana kurooraa toogoo
            httpx.go               Httpx puroobu toogoo
        runner/
            phase.go               Feezu ookesutoreeshon
            worker.go              Goroutine puuru, waakaa goto no RNG
        report/
            json.go                JSON repooto raitaa
            html.go                Jiritsu-gata HTML repooto
    Dockerfile                     Katana + Httpx tsuki maruchi suteeji
    .goreleaser.yml                Kurosu konpairu settei
    Makefile                       Build, test, lint taagetoo
    config.example.yml             Feezu deforuto chi sanshoo
```

## Yooken

- Go 1.21+
- [Katana](https://github.com/projectdiscovery/katana) (ProjectDiscovery) -- opushon, jidoo hakken-yoo
- [Httpx](https://github.com/projectdiscovery/httpx) (ProjectDiscovery) -- opushon, endpoint barideshon-yoo

```bash
bash scripts/install-tools.sh
```

## Men-seki jikoo

Kono tsuuru wa **ninka sareta sekyuriti tesuto oyobi pafoomansu hyooka nomi** wo mokuteki to shite imasu. Fuka tesuto wo jikkou suru mae ni, taagetto shisutemu no shoyuusha kara meikaku na shomen ni yoru kyoka wo shutoku suru hitsuyoo ga arimasu. Shoyuu shite inai, matawa tesuto no kyoka wo ete inai shisutemu ni taisuru kono tsuuru no mukan'i na shiyoo wa, chiiki, kokuritsu, matawa kokusai-hoo ni ihan suru kanousei ga arimasu. Sakusha wa kono tsuuru no goyo matawa songai ni tsuite issai no sekinin wo oimasen.

## Raisensu

MIT

## Sakusha

**Umut ÖZEN** -- [@umutozen](https://github.com/umutozen)

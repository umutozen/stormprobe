[English](../../README.md) · [Türkçe](README.tr.md) · [Deutsch](README.de.md) · [Francais](README.fr.md) · [Espanol](README.es.md) · [Italiano](README.it.md) · [Portugues](README.pt-BR.md) · [Russkij](README.ru.md) · [Zhongwen Jian](README.zh-CN.md) · [Zhongwen Fan](README.zh-TW.md) · [Nihongo](README.ja.md) · [Hangugeo](README.ko.md) · [Al-Arabiyya](README.ar.md) · [Hindi](README.hi.md) · [Nederlands](README.nl.md) · [Polski](README.pl.md) · [Ukrainska](README.uk.md)

# StormProbe

<p align="center">
  <img src="../../assets/banner.png" alt="StormProbe" width="100%">
</p>

**balgyeon. buha. bunseog.**

Jayu endpoint (endeupoindeu) balgyeon gineungeul gajchun peurodeogsyeon deunggeub HTTP buha teseuteu dogu. Yero uijonseong.

[![CI](https://github.com/umutozen/stormprobe/actions/workflows/ci.yml/badge.svg)](https://github.com/umutozen/stormprobe/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/umutozen/stormprobe)](https://goreportcard.com/report/github.com/umutozen/stormprobe)

---

## Gineung

- **Jadong balgyeon** -- Katana keuroling + Httpx peurobeu ro hwalseonghwa doen endpoint jadong balgyeon
- **Da dangye teseuteu** -- Ramp-Up (dangyejeog jeungga), Sustained (jisog), Spike (buha choegojeom), Recovery (boggu)
- **Sangsaehan meteurig** -- P50 / P95 / P99 latency (jiyeon), req/s, oyu bunyu
- **Idung bogoseo** -- JSON (gigyepandogga-neung) + HTML (sigag dakeu tema daesibodeu)
- **Yero uijonseong** -- Sunsuhan Go pyojun raibeureori, sseodeuparti modyul eobseum
- **Docker junbi** -- Katana + Httpx pohamdoen daniryang myeongryeong
- **Keuroseu peullaeseupom** -- GoReleaser tonghan Linux, macOS, Windows baineorivon

## Bbalri sijag

```bash
go run ./cmd https://example.com

go build -o stormprobe ./cmd
./stormprobe --insecure https://example.com
```

### Docker

```bash
# 빠른 테스트, 보고서 저장 없음
docker run --rm ghcr.io/umutozen/stormprobe:latest --insecure https://example.com

# 보고서를 ./outputs에 저장하는 전체 테스트 (Linux / macOS)
docker run --rm -v $(pwd)/outputs:/app/outputs ghcr.io/umutozen/stormprobe:latest --insecure --format both https://example.com

# 보고서를 ./outputs에 저장하는 전체 테스트 (Windows PowerShell)
docker run --rm -v "${PWD}\outputs:/app/outputs" ghcr.io/umutozen/stormprobe:latest --insecure --format both https://example.com

# 높은 부하 스파이크 테스트 (500명 동시 사용자)
docker run --rm ghcr.io/umutozen/stormprobe:latest --insecure --concurrency-spike 500 https://example.com

# 디스커버리 건너뛰기, 루트 경로만 테스트
docker run --rm ghcr.io/umutozen/stormprobe:latest --insecure --no-discovery https://example.com

# 사용자 정의 HTTP 헤더 (인증, 커스텀 테넌트 등)
docker run --rm ghcr.io/umutozen/stormprobe:latest --insecure -H "Authorization: Bearer TOKEN" -H "X-Tenant: acme" https://api.example.com
```

## Sayong beob

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

### Yesi

```bash
stormprobe --insecure --concurrency-spike 500 --req-per-worker 20 https://example.com

stormprobe --insecure --endpoints endpoints.txt https://example.com

stormprobe --insecure --format json --output ./results https://example.com

stormprobe --insecure -H "Authorization: Bearer TOKEN" -H "X-Tenant: acme" https://api.example.com

stormprobe --insecure --duration 30s --no-discovery https://example.com

stormprobe --insecure --duration 1m --concurrency-spike 300 https://example.com

stormprobe --insecure --alert-p99 500 --alert-error-rate 5 https://example.com
```

## Teseuteu dagyeong

| Dangye | Seolmyeong |
|---|---|
| **0 -- Discovery** | Katana keuroling + Httpx peurobeu, junbog jegeo doen endpoint mogrog |
| **1 -- Ramp-Up** | Dangyejeog jeungga: 5, 15, 30, 50 virtual users (gasangui sayongjadeul) |
| **2 -- Sustained** | Mogpyo concurrency (dongsiseong) eseo 3hoe pado, seongneung jeohareul cheukjeong |
| **3 -- Spike** | Choegojeom euro gabjagseureoun beoseuteu, geu hu naenggag |
| **4 -- Recovery** | Spike hu 10cho naenggag hu geongangseonggeum geomsa |

## Chulryeog

| Pail | Seolmyeong |
|---|---|
| `stormprobe_report_<timestamp>.json` | Modeun meteurigeul pohamhan gigyepandogga-neung gyeolgwa |
| `stormprobe_report_<timestamp>.html` | Sigag daesibodeu, modeun beurauzeo eseo yeolgi |

## Peurojekteu gujo

```
stormprobe/
    cmd/main.go                    CLI jibryeom
    internal/
        config/config.go           Gongyudoen taipeu, dangye saengseonggi
        metrics/
            latency.go             Baegbunwi gyesan
            errors.go              Oyu bunyu
        discovery/
            discovery.go           Discover() gonggyeong inteopeiseu
            katana.go              Katana keuroleo tongrhab
            httpx.go               Httpx peurobeu tonghab
        runner/
            phase.go               Dangye okeseuteureyisyeon
            worker.go              Goroutine pul, worker dangui RNG
        report/
            json.go                JSON bogoseo gisulgi
            html.go                Jachejeogin HTML bogoseo
    Dockerfile                     Katana + Httpx wa hamge meolti seutaeji
    .goreleaser.yml                Keuroseu keompail seoljeong
    Makefile                       Build, test, lint daesang
    config.example.yml             Dangye gimbon gab chamjo
```

## Yogeon

- Go 1.21+
- [Katana](https://github.com/projectdiscovery/katana) (ProjectDiscovery) -- seontaeg, jadong balgyeon yong
- [Httpx](https://github.com/projectdiscovery/httpx) (ProjectDiscovery) -- seontaeg, endpoint geomjeung yong

```bash
bash scripts/install-tools.sh
```

## Myeon-chaeg johang

I doguneun **seungin-doen bo-an teseuteu mit seongneung pyeong-ga-e-man** sayong-doeeya hamnida. Buha teseuteureul silhaeng-hagi jeon-e daesang siseutem so-yujaro-buteo myeong-hwaghan seomyeon dong-uireul bad-a-ya hamnida. So-yuha-ji anh-geo-na teseuteu gwon-han-i eobs-neun siseutem-e daehan i dogue-ui musieon sayong-eun hyeon-ji, guggae ttoneun gugjae beob-eul wiban-hal su issseubnida. Jeo-janeun i doguui o-yong ttoneun sonhae-e daehan eo-tteo-han chaeg-im-do ji-ji anhseubnida.

## Raisenseu

MIT

## Jeogia

**Umut ÖZEN** -- [@umutozen](https://github.com/umutozen)

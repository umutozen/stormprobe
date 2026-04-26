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

go build -o stormprobe ./cmd/stormprobe
./stormprobe --insecure https://example.com
```

## 설치

### 바이너리 다운로드 *(Go 불필요)*
[Releases 페이지](https://github.com/umutozen/stormprobe/releases)에서 플랫폼에 맞는 최신 버전을 다운로드하고, 압축을 풀어 PATH에 있는 디렉토리로 이동합니다:

```bash
# Linux / macOS
chmod +x stormprobe
sudo mv stormprobe /usr/local/bin/
stormprobe --no-discovery https://example.com

# Windows (PowerShell)
Move-Item stormprobe.exe C:\Windows\System32\stormprobe.exe
stormprobe --no-discovery https://example.com
```

### go install *(Go 1.21+ 필요, 가장 간단)*
```bash
go install github.com/umutozen/stormprobe/cmd/stormprobe@latest
stormprobe --no-discovery https://example.com
```

### 소스에서 빌드 *(Go 1.21+ 필요)*
```bash
git clone https://github.com/umutozen/stormprobe.git
cd stormprobe
go build -o stormprobe ./cmd/stormprobe   # Linux / macOS
go build -o stormprobe.exe ./cmd/stormprobe   # Windows
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
stormprobe [flags] <대상-URL>

Flags:
  --concurrency-ramp      int       Ramp-Up 최대 동시 연결 수 (기본값 50)
  --concurrency-sustained int       Sustained 단계 동시 연결 수 (기본값 75)
  --concurrency-spike     int       Spike 최대 동시 연결 수 (기본값 250)
  --req-per-worker        int       단계당 워커당 요청 수 (기본값 15)
  --timeout               duration  요청당 타임아웃 (기본값 10s)
  --endpoints             string    엔드포인트 목록 파일 (줄당 경로 하나)
  --no-discovery                    katana+httpx 건너뛰기, 루트 경로만 테스트
  --output                string    보고서 출력 디렉토리 (기본값 ./outputs)
  --format                string    보고서 형식: json, html, both (기본값 both)
  --katana-path           string    사용자 정의 katana 바이너리 경로
  --httpx-path            string    사용자 정의 httpx 바이너리 경로
  --insecure                        TLS 인증서 검증 건너뛰기
  --header, -H            string    사용자 정의 HTTP 헤더 (반복 가능): -H 'Authorization: Bearer TOKEN'
  --duration              duration  단계별 지속 시간 (예: 30s, 1m). req-per-worker 대신 사용
  --alert-p99             float     P99 지연이 Xms 초과 시 exit 1
  --alert-error-rate      float     오류율이 X%% 초과 시 exit 1
  --alert-rps             float     req/s가 X 미만 시 exit 1
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
        config/config.go           Gongyudoen taipeu, dangye saengseonggi, alert config
        alert/
            alert.go               Threshold checks, CI/CD exit code logic
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

**Umut ÖZEN** — [@umutozen](https://github.com/umutozen)

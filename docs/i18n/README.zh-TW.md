[English](../../README.md) · [Turkce](README.tr.md) · [Deutsch](README.de.md) · [Francais](README.fr.md) · [Espanol](README.es.md) · [Italiano](README.it.md) · [Portugues](README.pt-BR.md) · [Russkij](README.ru.md) · [Zhongwen Jian](README.zh-CN.md) · [Zhongwen Fan](README.zh-TW.md) · [Nihongo](README.ja.md) · [Hangugeo](README.ko.md) · [Al-Arabiyya](README.ar.md) · [Hindi](README.hi.md) · [Nederlands](README.nl.md) · [Polski](README.pl.md) · [Ukrainska](README.uk.md)

# StormProbe

<p align="center">
  <img src="../../assets/banner.png" alt="StormProbe" width="100%">
</p>

**fa xian. jia zai. fen xi.**

sheng chan ji HTTP fu zai ce shi gong ju, ju bei zi zhu endpoint (duan dian) fa xian gong neng. ling yi lai.

[![CI](https://github.com/umutozen/stormprobe/actions/workflows/ci.yml/badge.svg)](https://github.com/umutozen/stormprobe/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/umutozen/stormprobe)](https://goreportcard.com/report/github.com/umutozen/stormprobe)

---

## gong neng te xing

- **zi dong fa xian** -- Katana pa qu + Httpx tan ce zi dong cha zhao huo yue endpoint
- **duo jie duan ce shi** -- Ramp-Up (zhu bu zeng jia), Sustained (chi xu), Spike (fu zai feng zhi), Recovery (hui fu)
- **xiang xi zhi biao** -- P50 / P95 / P99 latency (yan chi), req/s, cuo wu fen lei
- **shuang bao gao** -- JSON (ji qi ke du) + HTML (ke shi hua an se zhu ti yi biao ban)
- **ling yi lai** -- chun Go biao zhun ku, wu di san fang mo kuai
- **Docker jiu xu** -- dan yi ming ling, nei zhi Katana + Httpx
- **kua ping tai** -- tong guo GoReleaser sheng cheng Linux, macOS, Windows er jin zhi wen jian

## kuai su kai shi

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

## shi yong fang fa

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

### shi li

```bash
stormprobe --insecure --concurrency-spike 500 --req-per-worker 20 https://example.com

stormprobe --insecure --endpoints endpoints.txt https://example.com

stormprobe --insecure --format json --output ./results https://example.com
```

## ce shi jie duan

| jie duan | miao shu |
|---|---|
| **0 -- Discovery** | Katana pa qu + Httpx tan ce, qu chong endpoint lie biao |
| **1 -- Ramp-Up** | zhu bu zeng jia: 5, 15, 30, 50 virtual users (xu ni yong hu) |
| **2 -- Sustained** | zai mu biao concurrency (bing fa) xia 3 bo, ce liang tui hua |
| **3 -- Spike** | tu ran bao fa dao feng zhi, ran hou leng que |
| **4 -- Recovery** | spike hou 10 miao leng que qi hou de jian kang jian cha |

## shu chu

| wen jian | miao shu |
|---|---|
| `stormprobe_report_<timestamp>.json` | bao han suo you zhi biao de ji qi ke du jie guo |
| `stormprobe_report_<timestamp>.html` | ke shi hua yi biao ban, zai ren yi liu lan qi zhong da kai |

## xiang mu jie gou

```
stormprobe/
    cmd/main.go                    CLI ru kou dian
    internal/
        config/config.go           gong xiang lei xing, jie duan sheng cheng qi
        metrics/
            latency.go             bai fen wei ji suan
            errors.go              cuo wu fen lei
        discovery/
            discovery.go           Discover() gong gong jie kou
            katana.go              Katana pa chong qi ji cheng
            httpx.go               Httpx tan ce ji cheng
        runner/
            phase.go               jie duan bian pai
            worker.go              Goroutine chi, mei ge worker de RNG
        report/
            json.go                JSON bao gao xie ru qi
            html.go                du li HTML bao gao
    Dockerfile                     duo jie duan, nei zhi Katana + Httpx
    .goreleaser.yml                jiao cha bian yi pei zhi
    Makefile                       Build, test, lint mu biao
    config.example.yml             jie duan mo ren zhi can kao
```

## xu qiu

- Go 1.21+
- [Katana](https://github.com/projectdiscovery/katana) (ProjectDiscovery) -- ke xuan, yong yu zi dong fa xian
- [Httpx](https://github.com/projectdiscovery/httpx) (ProjectDiscovery) -- ke xuan, yong yu endpoint yan zheng

```bash
bash scripts/install-tools.sh
```

## Mian ze sheng ming

Ben gong ju jin yong yu **jing shou quan de an quan ce shi he xing neng ping gu**. Zai yun xing ren he fu zai ce shi zhi qian, nin bi xu huo de mu biao xi tong suo you zhe de ming que shu mian xu ke. Wei jing shou quan dui nin bu yong you huo wu quan ce shi de xi tong shi yong ben gong ju ke neng wei fan dang di, guo jia huo guo ji fa lv. Zuo zhe dui ben gong ju de wu yong huo zao cheng de sun hai bu cheng dan ren he ze ren.

## xu ke zheng

MIT

## zuo zhe

**Umut ÖZEN** -- [@umutozen](https://github.com/umutozen)

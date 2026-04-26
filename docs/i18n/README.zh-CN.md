[English](../../README.md) · [Türkçe](README.tr.md) · [Deutsch](README.de.md) · [Francais](README.fr.md) · [Espanol](README.es.md) · [Italiano](README.it.md) · [Portugues](README.pt-BR.md) · [Russkij](README.ru.md) · [Zhongwen Jian](README.zh-CN.md) · [Zhongwen Fan](README.zh-TW.md) · [Nihongo](README.ja.md) · [Hangugeo](README.ko.md) · [Al-Arabiyya](README.ar.md) · [Hindi](README.hi.md) · [Nederlands](README.nl.md) · [Polski](README.pl.md) · [Ukrainska](README.uk.md)

# StormProbe

<p align="center">
  <img src="../../assets/banner.png" alt="StormProbe" width="100%">
</p>

**fa xian. jia zai. fen xi.**

Sheng chan ji HTTP fu zai ce shi gong ju, ju bei zi zhu endpoint (duan dian) fa xian gong neng. Ling yi lai.

[![CI](https://github.com/umutozen/stormprobe/actions/workflows/ci.yml/badge.svg)](https://github.com/umutozen/stormprobe/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/umutozen/stormprobe)](https://goreportcard.com/report/github.com/umutozen/stormprobe)

---

## Gong neng te xing

- **Zi dong fa xian** -- Katana pa qu + Httpx tan ce zi dong cha zhao huo yue endpoint
- **Duo jie duan ce shi** -- Ramp-Up (zhu bu zeng jia), Sustained (chi xu), Spike (fu zai feng zhi), Recovery (hui fu)
- **Xiang xi zhi biao** -- P50 / P95 / P99 latency (yan chi), req/s, cuo wu fen lei
- **Shuang bao gao** -- JSON (ji qi ke du) + HTML (ke shi hua an se zhu ti yi biao ban)
- **Ling yi lai** -- Chun Go biao zhun ku, wu di san fang mo kuai
- **Docker jiu xu** -- Dan yi ming ling, nei zhi Katana + Httpx
- **Kua ping tai** -- Tong guo GoReleaser sheng cheng Linux, macOS, Windows er jin zhi wen jian

## Kuai su kai shi

```bash
go run ./cmd https://example.com

go build -o stormprobe ./cmd/stormprobe
./stormprobe --insecure https://example.com
```

## 安装

### 下载二进制文件 *（无需 Go）*
从 [Releases 页面](https://github.com/umutozen/stormprobe/releases) 下载适合您平台的最新版本，解压后将二进制文件移至 PATH 目录：

```bash
# Linux / macOS
chmod +x stormprobe
sudo mv stormprobe /usr/local/bin/
stormprobe --no-discovery https://example.com

# Windows (PowerShell)
Move-Item stormprobe.exe C:\Windows\System32\stormprobe.exe
stormprobe --no-discovery https://example.com
```

### go install *（需要 Go 1.21+，最简单）*
```bash
go install github.com/umutozen/stormprobe/cmd/stormprobe@latest
stormprobe --no-discovery https://example.com
```

### 从源码构建 *（需要 Go 1.21+）*
```bash
git clone https://github.com/umutozen/stormprobe.git
cd stormprobe
go build -o stormprobe ./cmd/stormprobe   # Linux / macOS
go build -o stormprobe.exe ./cmd/stormprobe   # Windows
```

### Docker

```bash
# 快速测试，不保存报告
docker run --rm ghcr.io/umutozen/stormprobe:latest --insecure https://example.com

# 完整测试并保存报告到 ./outputs（Linux / macOS）
docker run --rm -v $(pwd)/outputs:/app/outputs ghcr.io/umutozen/stormprobe:latest --insecure --format both https://example.com

# 完整测试并保存报告到 ./outputs（Windows PowerShell）
docker run --rm -v "${PWD}\outputs:/app/outputs" ghcr.io/umutozen/stormprobe:latest --insecure --format both https://example.com

# 高负载峰值测试（500个并发用户）
docker run --rm ghcr.io/umutozen/stormprobe:latest --insecure --concurrency-spike 500 https://example.com

# 跳过发现，仅测试根路径
docker run --rm ghcr.io/umutozen/stormprobe:latest --insecure --no-discovery https://example.com

# 自定义 HTTP 请求头（授权、自定义租户等）
docker run --rm ghcr.io/umutozen/stormprobe:latest --insecure -H "Authorization: Bearer TOKEN" -H "X-Tenant: acme" https://api.example.com
```

## Shi yong fang fa

```
stormprobe [flags] <目标URL>

Flags:
  --concurrency-ramp      int       爬升阶段最大并发数（默认 50）
  --concurrency-sustained int       持续阶段并发数（默认 75）
  --concurrency-spike     int       峰值阶段最大并发数（默认 250）
  --req-per-worker        int       每个 Worker 每步请求数（默认 15）
  --timeout               duration  每请求超时时间（默认 10s）
  --endpoints             string    端点列表文件（每行一个路径，跳过探测）
  --no-discovery                    跳过 katana+httpx，仅测试根路径
  --output                string    报告输出目录（默认 ./outputs）
  --format                string    报告格式：json、html、both（默认 both）
  --katana-path           string    自定义 katana 可执行文件路径
  --httpx-path            string    自定义 httpx 可执行文件路径
  --insecure                        跳过 TLS 证书验证
  --header, -H            string    自定义 HTTP 头部（可重复）：-H 'Authorization: Bearer TOKEN'
  --duration              duration  每阶段持续时间（如 30s、1m），优先于 req-per-worker
  --alert-p99             float     任意阶段 P99 延迟超过 Xms 时退出码为 1
  --alert-error-rate      float     任意阶段错误率超过 X%% 时退出码为 1
  --alert-rps             float     任意阶段 req/s 低于 X 时退出码为 1
```

### Shi li

```bash
stormprobe --insecure --concurrency-spike 500 --req-per-worker 20 https://example.com

stormprobe --insecure --endpoints endpoints.txt https://example.com

stormprobe --insecure --format json --output ./results https://example.com

stormprobe --insecure -H "Authorization: Bearer TOKEN" -H "X-Tenant: acme" https://api.example.com

stormprobe --insecure --duration 30s --no-discovery https://example.com

stormprobe --insecure --duration 1m --concurrency-spike 300 https://example.com

stormprobe --insecure --alert-p99 500 --alert-error-rate 5 https://example.com
```

## Ce shi jie duan

| Jie duan | Miao shu |
|---|---|
| **0 -- Discovery** | Katana pa qu + Httpx tan ce, qu chong endpoint lie biao |
| **1 -- Ramp-Up** | Zhu bu zeng jia: 5, 15, 30, 50 virtual users (xu ni yong hu) |
| **2 -- Sustained** | Zai mu biao concurrency (bing fa) xia 3 bo, ce liang tui hua |
| **3 -- Spike** | Tu ran bao fa dao feng zhi, ran hou leng que |
| **4 -- Recovery** | Spike hou 10 miao leng que qi hou de jian kang jian cha |

## Shu chu

| Wen jian | Miao shu |
|---|---|
| `stormprobe_report_<timestamp>.json` | Bao han suo you zhi biao de ji qi ke du jie guo |
| `stormprobe_report_<timestamp>.html` | Ke shi hua yi biao ban, zai ren yi liu lan qi zhong da kai |

## Xiang mu jie gou

```
stormprobe/
    cmd/main.go                    CLI ru kou dian
    internal/
        config/config.go           Gong xiang lei xing, jie duan sheng cheng qi, alert config
        alert/
            alert.go               Threshold checks, CI/CD exit code logic
        metrics/
            latency.go             Bai fen wei ji suan
            errors.go              Cuo wu fen lei
        discovery/
            discovery.go           Discover() gong gong jie kou
            katana.go              Katana pa chong qi ji cheng
            httpx.go               Httpx tan ce ji cheng
        runner/
            phase.go               Jie duan bian pai
            worker.go              Goroutine chi, mei ge worker de RNG
        report/
            json.go                JSON bao gao xie ru qi
            html.go                Du li HTML bao gao
    Dockerfile                     Duo jie duan, nei zhi Katana + Httpx
    .goreleaser.yml                Jiao cha bian yi pei zhi
    Makefile                       Build, test, lint mu biao
    config.example.yml             Jie duan mo ren zhi can kao
```

## Xu qiu

- Go 1.21+
- [Katana](https://github.com/projectdiscovery/katana) (ProjectDiscovery) -- ke xuan, yong yu zi dong fa xian
- [Httpx](https://github.com/projectdiscovery/httpx) (ProjectDiscovery) -- ke xuan, yong yu endpoint yan zheng

```bash
bash scripts/install-tools.sh
```

## Mian ze sheng ming

Ben gong ju jin yong yu **jing shou quan de an quan ce shi he xing neng ping gu**. Zai yun xing ren he fu zai ce shi zhi qian, nin bi xu huo de mu biao xi tong suo you zhe de ming que shu mian xu ke. Wei jing shou quan dui nin bu yong you huo wu quan ce shi de xi tong shi yong ben gong ju ke neng wei fan dang di, guo jia huo guo ji fa lv. Zuo zhe dui ben gong ju de wu yong huo zao cheng de sun hai bu cheng dan ren he ze ren.

## Xu ke zheng

MIT

## Zuo zhe

**Umut ÖZEN** — [@umutozen](https://github.com/umutozen)

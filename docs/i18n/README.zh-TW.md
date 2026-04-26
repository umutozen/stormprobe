[English](../../README.md) · [Türkçe](README.tr.md) · [Deutsch](README.de.md) · [Francais](README.fr.md) · [Espanol](README.es.md) · [Italiano](README.it.md) · [Portugues](README.pt-BR.md) · [Russkij](README.ru.md) · [Zhongwen Jian](README.zh-CN.md) · [Zhongwen Fan](README.zh-TW.md) · [Nihongo](README.ja.md) · [Hangugeo](README.ko.md) · [Al-Arabiyya](README.ar.md) · [Hindi](README.hi.md) · [Nederlands](README.nl.md) · [Polski](README.pl.md) · [Ukrainska](README.uk.md)

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

go build -o stormprobe ./cmd/stormprobe
./stormprobe --insecure https://example.com
```

## 安裝

### 二進位檔案（建議）
從 [Releases 頁面](https://github.com/umutozen/stormprobe/releases) 下載適合您平台的最新版本，解壓後將二進位檔案移至 PATH 目錄：

```bash
# Linux / macOS
chmod +x stormprobe
sudo mv stormprobe /usr/local/bin/
stormprobe --no-discovery https://example.com

# Windows (PowerShell)
Move-Item stormprobe.exe C:\Windows\System32\stormprobe.exe
stormprobe --no-discovery https://example.com
```

### go install (requires Go 1.21+)
```bash
go install github.com/umutozen/stormprobe/cmd/stormprobe@latest
stormprobe --no-discovery https://example.com
```

### 從原始碼構建（需要 Go 1.21+）
```bash
git clone https://github.com/umutozen/stormprobe.git
cd stormprobe
go build -o stormprobe ./cmd/stormprobe   # Linux / macOS
go build -o stormprobe.exe ./cmd/stormprobe   # Windows
```

### Docker

```bash
# 快速測試，不儲存報告
docker run --rm ghcr.io/umutozen/stormprobe:latest --insecure https://example.com

# 完整測試並儲存報告至 ./outputs（Linux / macOS）
docker run --rm -v $(pwd)/outputs:/app/outputs ghcr.io/umutozen/stormprobe:latest --insecure --format both https://example.com

# 完整測試並儲存報告至 ./outputs（Windows PowerShell）
docker run --rm -v "${PWD}\outputs:/app/outputs" ghcr.io/umutozen/stormprobe:latest --insecure --format both https://example.com

# 高負載峰值測試（500個並發用戶）
docker run --rm ghcr.io/umutozen/stormprobe:latest --insecure --concurrency-spike 500 https://example.com

# 跳過探索，僅測試根路徑
docker run --rm ghcr.io/umutozen/stormprobe:latest --insecure --no-discovery https://example.com

# 自訂 HTTP 標頭（授權、自訂租戶等）
docker run --rm ghcr.io/umutozen/stormprobe:latest --insecure -H "Authorization: Bearer TOKEN" -H "X-Tenant: acme" https://api.example.com
```

## shi yong fang fa

```
stormprobe [flags] <目標URL>

Flags:
  --concurrency-ramp      int       爬升階段最大並發數（預設 50）
  --concurrency-sustained int       持續階段並發數（預設 75）
  --concurrency-spike     int       峰值階段最大並發數（預設 250）
  --req-per-worker        int       每個 Worker 每步請求數（預設 15）
  --timeout               duration  每請求超時時間（預設 10s）
  --endpoints             string    端點列表檔案（每行一個路徑，跳過探測）
  --no-discovery                    跳過 katana+httpx，僅測試根路徑
  --output                string    報告輸出目錄（預設 ./outputs）
  --format                string    報告格式：json、html、both（預設 both）
  --katana-path           string    自訂 katana 可執行檔路徑
  --httpx-path            string    自訂 httpx 可執行檔路徑
  --insecure                        跳過 TLS 憑證驗證
  --header, -H            string    自訂 HTTP 標頭（可重複）：-H 'Authorization: Bearer TOKEN'
  --duration              duration  每階段持續時間（如 30s、1m），優先於 req-per-worker
  --alert-p99             float     任意階段 P99 延遲超過 Xms 時退出碼為 1
  --alert-error-rate      float     任意階段錯誤率超過 X%% 時退出碼為 1
  --alert-rps             float     任意階段 req/s 低於 X 時退出碼為 1
```

### shi li

```bash
stormprobe --insecure --concurrency-spike 500 --req-per-worker 20 https://example.com

stormprobe --insecure --endpoints endpoints.txt https://example.com

stormprobe --insecure --format json --output ./results https://example.com

stormprobe --insecure -H "Authorization: Bearer TOKEN" -H "X-Tenant: acme" https://api.example.com

stormprobe --insecure --duration 30s --no-discovery https://example.com

stormprobe --insecure --duration 1m --concurrency-spike 300 https://example.com

stormprobe --insecure --alert-p99 500 --alert-error-rate 5 https://example.com
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
        config/config.go           gong xiang lei xing, jie duan sheng cheng qi, alert config
        alert/
            alert.go               Threshold checks, CI/CD exit code logic
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

**Umut ÖZEN** — [@umutozen](https://github.com/umutozen)

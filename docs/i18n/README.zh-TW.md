[English](../../README.md) · [Türkçe](README.tr.md) · [Deutsch](README.de.md) · [Francais](README.fr.md) · [Espanol](README.es.md) · [Italiano](README.it.md) · [Portugues](README.pt-BR.md) · [Russkij](README.ru.md) · [Zhongwen Jian](README.zh-CN.md) · [Zhongwen Fan](README.zh-TW.md) · [Nihongo](README.ja.md) · [Hangugeo](README.ko.md) · [Al-Arabiyya](README.ar.md) · [Hindi](README.hi.md) · [Nederlands](README.nl.md) · [Polski](README.pl.md) · [Ukrainska](README.uk.md)

<p align="center">
  <img src="../../assets/banner.png" alt="StormProbe Banner" width="100%">
</p>

<h1 align="center">StormProbe</h1>

<p align="center">
  <strong>faxian · jiazai · fenxi</strong><br>
  具有自主端點發現功能的生產級HTTP負載測試工具。零外部依賴。
</p>

<p align="center">
  <a href="https://github.com/umutozen/stormprobe/actions/workflows/ci.yml">
    <img src="https://github.com/umutozen/stormprobe/actions/workflows/ci.yml/badge.svg" alt="CI">
  </a>
  <a href="https://goreportcard.com/report/github.com/umutozen/stormprobe">
    <img src="https://goreportcard.com/badge/github.com/umutozen/stormprobe" alt="Go Report Card">
  </a>
  <a href="https://github.com/umutozen/stormprobe/releases">
    <img src="https://img.shields.io/github/v/release/umutozen/stormprobe" alt="最新版本">
  </a>
  <a href="../../LICENSE">
    <img src="https://img.shields.io/badge/license-MIT-blue.svg" alt="MIT許可證">
  </a>
</p>

---

## 什麼是StormProbe？

StormProbe是一個用純Go編寫的**零依賴**HTTP負載測試工具。它使用[Katana](https://github.com/projectdiscovery/katana)和[Httpx](https://github.com/projectdiscovery/httpx)自動發現應用程序的端點，然後執行結構化的4階段負載測試——**Ramp-Up → Sustained → Spike → Recovery**——並生成包含P50/P95/P99延遲指標的JSON和HTML報告。

適用於：
- **CI/CD管道** — 帶有退出代碼的警報閾值
- **滲透測試人員** — 自動發現活躍端點
- **DevOps工程師** — 部署前後的基準測試
- **QA團隊** — 性能SLA驗證

---

## 功能

| 功能 | 詳情 |
|---|---|
| **自動發現** | Katana爬蟲 + Httpx探測，刪重後的活躍端點列表 |
| **4階段負載測試** | Ramp-Up → Sustained → Spike → Recovery |
| **豐富指標** | P50/P95/P99延遲、req/s、按階段分類錯誤 |
| **持續時間模式** | 基於時間的階段 (`--duration 30s`) 替代請求計數 |
| **警報閾值** | P99、錯誤率或RPS違規時觸發CI/CD兼容的`exit 1` |
| **雙重報告** | JSON（機器可讀）+ 獨立HTML暗色儀錶盤 |
| **自定義標頭** | Bearer令牌、租戶標頭、Cookie — 全部傳播 |
| **Docker就緒** | 內置Katana + Httpx的單一鏡像 |
| **跨平台** | Linux、macOS、Windows — amd64 & arm64 |
| **零依賴** | 純Go 1.21+ stdlib，`go install`即可完成 |

---

## 快速開始

```bash
# 安裝（需要Go 1.21+）
go install github.com/umutozen/stormprobe/cmd/stormprobe@latest

# 基本測試 — 自動發現，4個階段
stormprobe https://example.com

# 跳過發現，僅測試根路徑
stormprobe --no-discovery https://example.com

# TLS + 認證標頭 + 僅JSON報告
stormprobe --insecure -H "Authorization: Bearer TOKEN" --format json https://api.example.com

# 基於時間的測試：每個階段30秒
stormprobe --insecure --duration 30s --no-discovery https://example.com

# CI/CD：如果P99 > 500ms或錯誤 > 5%則exit 1
stormprobe --alert-p99 500 --alert-error-rate 5 https://example.com
```

---

## 安裝

### 選項1 — `go install` *(推薦)*

```bash
go install github.com/umutozen/stormprobe/cmd/stormprobe@latest
stormprobe --help
```

### 選項2 — 預構建二進制文件 *(無需Go)*

從[Releases頁面](https://github.com/umutozen/stormprobe/releases)下載。

```bash
# Linux / macOS
chmod +x stormprobe
sudo mv stormprobe /usr/local/bin/

# Windows PowerShell
Move-Item stormprobe.exe C:\Windows\System32\stormprobe.exe
```

### 選項3 — 從原始碼構建

```bash
git clone https://github.com/umutozen/stormprobe.git
cd stormprobe
go build -o stormprobe ./cmd/stormprobe      # Linux / macOS
go build -o stormprobe.exe ./cmd/stormprobe  # Windows
```

### 選項4 — Docker

```bash
docker run --rm ghcr.io/umutozen/stormprobe:latest --insecure https://example.com

# Linux / macOS
docker run --rm -v $(pwd)/outputs:/app/outputs \
  ghcr.io/umutozen/stormprobe:latest --insecure --format both https://example.com

# Windows PowerShell
docker run --rm -v "${PWD}\outputs:/app/outputs" `
  ghcr.io/umutozen/stormprobe:latest --insecure --format both https://example.com
```

### 可選 — 發現工具

```bash
bash scripts/install-tools.sh
```

> 沒有Katana/Httpx時使用`--no-discovery`或`--endpoints`。

---

## CLI參考

### 並發與負載

| 標誌 | 預設值 | 說明 |
|---|---|---|
| `--concurrency-ramp` | `50` | Ramp-Up階段的峰值並發 |
| `--concurrency-sustained` | `75` | Sustained階段的並發 |
| `--concurrency-spike` | `250` | Spike階段的峰值並發 |
| `--req-per-worker` | `15` | 每個worker的請求數（設定`--duration`時忽略） |
| `--duration` | `0` (停用) | 基於時間的階段持續時間，如`30s`、`1m` |
| `--timeout` | `10s` | 每個請求的逾時 |

### 發現

| 標誌 | 預設值 | 說明 |
|---|---|---|
| `--no-discovery` | `false` | 跳過Katana+Httpx，僅測試根路徑 |
| `--endpoints` | `""` | 從文件載入端點 |
| `--katana-path` | `""` | 自定義Katana二進制路徑 |
| `--httpx-path` | `""` | 自定義Httpx二進制路徑 |

### HTTP與安全

| 標誌 | 預設值 | 說明 |
|---|---|---|
| `-H`, `--header` | — | 自定義標頭（可重複） |
| `--insecure` | `false` | 跳過TLS憑證驗證 |

### 輸出

| 標誌 | 預設值 | 說明 |
|---|---|---|
| `--format` | `both` | 報告格式：`json`、`html`、`both` |
| `--output` | `./outputs` | 報告的輸出目錄 |

### CI/CD警報

| 標誌 | 預設值 | 說明 |
|---|---|---|
| `--alert-p99` | `0` (停用) | P99(ms)超出閾值時exit 1 |
| `--alert-error-rate` | `0` (停用) | 錯誤率(%)超出閾值時exit 1 |
| `--alert-rps` | `0` (停用) | req/s低於閾值時exit 1 |

---

## 測試階段

| 階段 | 說明 |
|---|---|
| **0 — Discovery** | Katana爬蟲 + Httpx探測，去重列表 |
| **1 — Ramp-Up** | 逐步增加並發：5 → 15 → 30 → 峰值 |
| **2 — Sustained** | 在目標並發下進行3次衝擊 |
| **3 — Spike** | 突然增加到峰值並發，然後冷卻 |
| **4 — Recovery** | 10秒冷卻後的峰值後健康檢查 |

---

## 輸出與報告

| 文件 | 說明 |
|---|---|
| `stormprobe_report_<timestamp>.json` | 每個階段的所有指標，機器可讀 |
| `stormprobe_report_<timestamp>.html` | 帶SVG延遲圖表的獨立暗色儀錶盤 |

---

## 免責聲明

本工具僅用於**授權的安全測試和性能評估**。在對目標系統執行任何負載測試之前，必須獲得系統所有者的明確書面許可。作者對本工具的誤用或造成的任何損害**不承擔任何責任**。

---

## 許可證

[MIT](../../LICENSE) © [Umut ÖZEN](https://github.com/umutozen)

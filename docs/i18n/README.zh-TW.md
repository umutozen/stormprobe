# StormProbe

<p align="center">
  <img src="../../assets/banner.png" alt="StormProbe Banner" width="100%">
</p>

<h1 align="center">StormProbe</h1>

<p align="center">
  <strong>crawl · discover · stress · decide</strong><br>
  現代 HTTP 壓力測試工具，支援自主 endpoint 發現與生產就緒性裁決。
</p>

<p align="center">
  <strong>不只是基準測試。</strong><br>
  StormProbe 不只是測試。它解釋您的系統是否真正變得更安全。
</p>

---

## 為何選擇 StormProbe？

傳統壓測工具只顯示數字。<br>StormProbe 給出明確決策。

無需手動選擇 endpoint，StormProbe 自動發現活躍端點，執行分階段壓力測試，檢測飽和點，並生成可操作的最終裁決與建議。

### 傳統工具 vs StormProbe

| 傳統工具 | StormProbe |
|---|---|
| 手動選擇 endpoint | 自動發現活躍 endpoint |
| 原始基準測試數據 | 最終裁決 + 瓶頸診斷 |
| 需要手動解讀 | 具體的操作建議 |
| 單一壓力模式 | 遞增 + 維持 + 峰值 + 恢復 |
| "系統變慢了" | "這裡出錯了，且原因在此" |

---

## 它的功能？

StormProbe 是一個以 Go 編寫的生產級 HTTP 壓力測試命令列工具。

- Endpoint 發現
- 分階段並發測試
- 延遲百分位分析
- 吞吐量飽和檢測
- 恢復驗證
- CI/CD 閾值警報
- JSON + HTML 報告
- 生產就緒性裁決

---

## 真實輸出範例

```text
FINAL VERDICT

Rating                          : Degraded Under Load
Safe concurrency                : up to 30 virtual users
Perf. degradation starts at     : 50 vu
Severe degradation starts at    : 150 vu
Throughput plateau              : ~36 req/s
Recovery                        : OK

Bottleneck    : Throughput saturation
Detail        : Latency grew 5.4x while throughput grew only 2.1x
                for 10x concurrency increase.

Recommendation           : Cap production traffic at ~30 concurrent users.
Production Recommendation : Safe for low-to-medium traffic only
```

**Compare Mode — The Decision Engine:**

```text
  stormprobe compare before.json after.json

  Rating           : High Risk → Not Production Ready
  Confidence       : Unstable — Recovery Failed → Low Confidence — Immediate Action Required
  Safe concurrency : 15 VU → 50 VU
  Recovery         : FAILED → OK

  ┌─────────────────────────────────────────────────────────────────┐
  │  Verdict: PARTIAL IMPROVEMENT — Still Not Production Ready     │
  └─────────────────────────────────────────────────────────────────┘

  Why:
    + Safe concurrency increased 233% (15 → 50 VU)
    + Recovery restored (was FAILED, now OK)
    + Confidence improved (Unstable → Low Confidence)
    + Lower transient error volume observed (6169 → 4 failed)
```

---

## 特色

| 特色 | 詳情 |
|---|---|
| 自動發現 | Katana crawl + Httpx probe 並對活躍 endpoint 去重 |
| 四階段負載測試 | Ramp-Up → Sustained → Spike → Recovery |
| 裁決引擎 | 安全併發數、效能退化點、瓶頸分類 |
| 豐富指標 | Avg / P50 / P95 / P99 延遲 + req/s + 錯誤分類 |
| 吞吐量分析 | 在嚴重故障發生前檢測系統飽和點 |
| 恢復驗證 | 確認系統在流量飆升後是否能夠恢復穩態 |
| 警報閾值 | 針對延遲、錯誤率與 RPS 的 CI/CD 就緒退出代碼 |
| 雙重報告 | JSON + 獨立的 HTML 儀表板 |
| 自訂 Headers | 授權 Token、租戶請求頭、Cookies |
| 支援 Docker | 包含綑綁工具的單一映像檔 |
| 跨平台 | Linux, macOS, Windows — amd64 & arm64 |
| 純 Go 核心 | 快速安裝，極低的維運成本 |

---

## 快速開始

```bash
# 安裝
go install github.com/umutozen/stormprobe/cmd/stormprobe@latest

# 全自動發現與測試
stormprobe https://example.com

# 跳過發現，僅測試根路徑
stormprobe --no-discovery https://example.com

# 帶 Auth header 的 API 測試
stormprobe --insecure \
  -H "Authorization: Bearer TOKEN" \
  https://api.example.com

# CI/CD 門檻設定
stormprobe \\
  --alert-p99 500 \\
  --alert-error-rate 5 \\
  https://example.com

# 使用 JSON 設定檔
stormprobe --config stormprobe.json
```

---

## 測試流程

```text
發現 → 預熱 → 維持 → 峰值 → 恢復 → 裁決
```

| 階段 | 目的 |
|---|---|
| Discovery | 爬取並驗證存活 endpoint |
| Ramp-Up | 並發數量漸增 |
| Sustained | 驗證穩態效能 |
| Spike | 突發流量模擬 |
| Recovery | 驗證峰值過後的系統健康 |
| Verdict | 安全並發數 + 瓶頸分析 |

---

## 安裝

### Go Install

```bash
go install github.com/umutozen/stormprobe/cmd/stormprobe@latest
```

推薦方式。

### Docker

```bash
docker run --rm ghcr.io/umutozen/stormprobe:latest --insecure https://example.com
```

---

## CI/CD

```yaml
- name: Performance Gate
  run: |
    go install github.com/umutozen/stormprobe/cmd/stormprobe@latest
    stormprobe \
      --no-discovery \
      --duration 30s \
      --alert-p99 500 \
      --alert-error-rate 2 \
      https://staging.example.com
```

這讓效能衰退時的部署會自動失敗(Fail)。

---

## 使用情境

**上線前** — 驗證並發變更是否破壞了延遲目標。

**安全性 + 效能** — 發現隱藏的活躍 endpoint 並安全地壓測它們。

**基礎設施稽核** — 在流量飆升前測量系統飽和點。

**退化檢測 (Regression)** — 比較部署前後在 CI 上的效能。

---

## Flags

| Flag | 描述 |
|---|---|
| `--duration` | 每階段持續時間 (例如 30s, 1m) |
| `--alert-p99` | P99 延遲若超過閾值(ms)則報錯退出 |
| `--alert-error-rate` | 錯誤率若超過閾值則報錯退出 |
| `--alert-rps` | Req/s 若低於閾值則報錯退出 |
| `-H` / `--header` | 自訂 HTTP 請求頭 |
| `--insecure` | 略過 TLS 驗證 |
| `--no-discovery` | 略過發現，僅測試根路徑 |
| `--endpoints` | 從檔案載入端點列表 |
| `--format` | 報告格式 |
| `--concurrency-ramp` | 預熱階段峰值 |
| `--concurrency-sustained` | 穩定維持的並發數 |
| `--concurrency-spike` | 突發峰值 |
| `--req-per-worker` | 每個執行緒的請求數 |
| `--timeout` | 單次請求超時 |
| `--config` | 從 JSON 設定檔載入設定 |
| `--auto-profile` | 自動偵測伺服器堆疊並應用診斷配置檔（預設：啟用） |
| `--profile` | 手動指定配置檔：iis、tomcat、php |
| `--no-fingerprint` | 跳過伺服器堆疊偵測 |
| `--show-profile` | 顯示偵測到的堆疊與配置檔，然後退出 |

---

## 需求

| 工具 | 備註 |
|---|---|
| Go 1.21+ | Go 環境 |
| Katana | 選用 |
| Httpx | 選用 |
| Docker | Docker 選用 |

---

## 免責聲明

本工具僅供經過授權的測試使用。

請確保你已獲得擁有人授權。

未經授權使用可能會觸法。

---

## 貢獻

請見 CONTRIBUTING 檔案。

## 安全性

回報安全漏洞請查看 SECURITY 檔案。

## 授權

MIT License © Umut ÖZEN

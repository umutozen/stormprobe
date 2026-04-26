[English](../../README.md) · [Türkçe](README.tr.md) · [Deutsch](README.de.md) · [Francais](README.fr.md) · [Espanol](README.es.md) · [Italiano](README.it.md) · [Portugues](README.pt-BR.md) · [Russkij](README.ru.md) · [Zhongwen Jian](README.zh-CN.md) · [Zhongwen Fan](README.zh-TW.md) · [Nihongo](README.ja.md) · [Hangugeo](README.ko.md) · [Al-Arabiyya](README.ar.md) · [Hindi](README.hi.md) · [Nederlands](README.nl.md) · [Polski](README.pl.md) · [Ukrainska](README.uk.md)

<p align="center">
  <img src="../../assets/banner.png" alt="StormProbe Banner" width="100%">
</p>

<h1 align="center">StormProbe</h1>

<p align="center">
  <strong>hakken · fuka · bunseki</strong><br>
  Jishu-tekina endopointohakken wo sonaeta seisan-kyuu no HTTP fuka tesutツール。gaibu izon zero。
</p>

<p align="center">
  <a href="https://github.com/umutozen/stormprobe/actions/workflows/ci.yml">
    <img src="https://github.com/umutozen/stormprobe/actions/workflows/ci.yml/badge.svg" alt="CI">
  </a>
  <a href="https://goreportcard.com/report/github.com/umutozen/stormprobe">
    <img src="https://goreportcard.com/badge/github.com/umutozen/stormprobe" alt="Go Report Card">
  </a>
  <a href="https://github.com/umutozen/stormprobe/releases">
    <img src="https://img.shields.io/github/v/release/umutozen/stormprobe" alt="Saishinsaban">
  </a>
  <a href="../../LICENSE">
    <img src="https://img.shields.io/badge/license-MIT-blue.svg" alt="MIT License">
  </a>
</p>

---

## StormProbeとは？

StormProbeはJure Go で書かれた**依存関係ゼロ**のHTTP負荷テストツールです。[Katana](https://github.com/projectdiscovery/katana)と[Httpx](https://github.com/projectdiscovery/httpx)を使用してアプリケーションのエンドポイントを自動的に検出し、**Ramp-Up → Sustained → Spike → Recovery**という4フェーズの負荷テストを実行し、P50/P95/P99レイテンシメトリクスを含むJSONおよびHTMLレポートを生成します。

対象ユーザー:
- **CI/CDパイプライン** — 終了コード付きアラート閾値
- **ペネトレーションテスター** — ライブエンドポイントの自動検出
- **DevOpsエンジニア** — デプロイ前後のベンチマーク
- **QAチーム** — パフォーマンスSLAの検証

---

## 機能

| 機能 | 詳細 |
|---|---|
| **自動検出** | Katanaクロール + Httpxプローブ、重複排除済みリスト |
| **4フェーズ負荷テスト** | Ramp-Up → Sustained → Spike → Recovery |
| **豊富なメトリクス** | P50/P95/P99レイテンシ、req/s、フェーズ別エラー分類 |
| **時間モード** | リクエスト数の代わりに時間ベースのフェーズ (`--duration 30s`) |
| **アラート閾値** | P99、エラー率、RPSの違反時にCI/CD対応の`exit 1` |
| **デュアルレポート** | JSON（機械可読）+ スタンドアロンHTMLダークダッシュボード |
| **カスタムヘッダー** | Bearerトークン、テナントヘッダー、Cookie — 全体に伝播 |
| **Docker対応** | KatanaとHttpxが組み込まれた単一イメージ |
| **クロスプラットフォーム** | Linux、macOS、Windows — amd64 & arm64 |
| **依存関係ゼロ** | 純粋なGo 1.21+ stdlib、`go install`だけで完了 |

---

## クイックスタート

```bash
# インストール（Go 1.21+が必要）
go install github.com/umutozen/stormprobe/cmd/stormprobe@latest

# 基本テスト — エンドポイント自動検出、4フェーズ実行
stormprobe https://example.com

# 検出をスキップ、ルートパスのみテスト
stormprobe --no-discovery https://example.com

# TLS + 認証ヘッダー + JSONレポートのみ
stormprobe --insecure -H "Authorization: Bearer TOKEN" --format json https://api.example.com

# 時間ベーステスト: フェーズごとに30秒
stormprobe --insecure --duration 30s --no-discovery https://example.com

# CI/CD: P99 > 500ms または エラー > 5% なら exit 1
stormprobe --alert-p99 500 --alert-error-rate 5 https://example.com
```

---

## インストール

### オプション1 — `go install` *(推奨)*

```bash
go install github.com/umutozen/stormprobe/cmd/stormprobe@latest
stormprobe --help
```

### オプション2 — ビルド済みバイナリ *(Go不要)*

[Releasesページ](https://github.com/umutozen/stormprobe/releases)からダウンロード。

```bash
# Linux / macOS
chmod +x stormprobe
sudo mv stormprobe /usr/local/bin/

# Windows PowerShell
Move-Item stormprobe.exe C:\Windows\System32\stormprobe.exe
```

### オプション3 — ソースからビルド

```bash
git clone https://github.com/umutozen/stormprobe.git
cd stormprobe
go build -o stormprobe ./cmd/stormprobe      # Linux / macOS
go build -o stormprobe.exe ./cmd/stormprobe  # Windows
```

### オプション4 — Docker

```bash
docker run --rm ghcr.io/umutozen/stormprobe:latest --insecure https://example.com

# Linux / macOS
docker run --rm -v $(pwd)/outputs:/app/outputs \
  ghcr.io/umutozen/stormprobe:latest --insecure --format both https://example.com

# Windows PowerShell
docker run --rm -v "${PWD}\outputs:/app/outputs" `
  ghcr.io/umutozen/stormprobe:latest --insecure --format both https://example.com
```

### オプション — 検出ツール

```bash
bash scripts/install-tools.sh
```

> Katana/Httpxなしの場合は`--no-discovery`または`--endpoints`を使用。

---

## CLIリファレンス

### 並行性と負荷

| フラグ | デフォルト | 説明 |
|---|---|---|
| `--concurrency-ramp` | `50` | Ramp-Upフェーズのピーク並行性 |
| `--concurrency-sustained` | `75` | Sustainedフェーズの並行性 |
| `--concurrency-spike` | `250` | Spikeフェーズのピーク並行性 |
| `--req-per-worker` | `15` | ワーカーあたりのリクエスト数 (`--duration`設定時は無視) |
| `--duration` | `0` (無効) | 時間ベースのフェーズ期間（例: `30s`, `1m`） |
| `--timeout` | `10s` | リクエストごとのタイムアウト |

### 検出

| フラグ | デフォルト | 説明 |
|---|---|---|
| `--no-discovery` | `false` | Katana+Httpxをスキップ、ルートのみテスト |
| `--endpoints` | `""` | ファイルからエンドポイントを読み込む |
| `--katana-path` | `""` | カスタムKatanaバイナリパス |
| `--httpx-path` | `""` | カスタムHttpxバイナリパス |

### HTTPとセキュリティ

| フラグ | デフォルト | 説明 |
|---|---|---|
| `-H`, `--header` | — | カスタムヘッダー（繰り返し可） |
| `--insecure` | `false` | TLS証明書の検証をスキップ |

### 出力

| フラグ | デフォルト | 説明 |
|---|---|---|
| `--format` | `both` | レポート形式: `json`, `html`, `both` |
| `--output` | `./outputs` | レポートの出力ディレクトリ |

### CI/CDアラート

| フラグ | デフォルト | 説明 |
|---|---|---|
| `--alert-p99` | `0` (無効) | P99 (ms)が閾値を超えたらexit 1 |
| `--alert-error-rate` | `0` (無効) | エラー率(%)が閾値を超えたらexit 1 |
| `--alert-rps` | `0` (無効) | req/sが閾値を下回ったらexit 1 |

---

## テストフェーズ

| フェーズ | 説明 |
|---|---|
| **0 — Discovery** | Katanaクロール + Httpxプローブ、重複排除リスト |
| **1 — Ramp-Up** | 段階的な並行性増加: 5 → 15 → 30 → ピーク |
| **2 — Sustained** | 目標並行性で3ウェーブ |
| **3 — Spike** | ピークへの急激な増加、その後クールダウン |
| **4 — Recovery** | 10秒クールダウン後のヘルスチェック |

---

## 出力とレポート

| ファイル | 説明 |
|---|---|
| `stormprobe_report_<timestamp>.json` | フェーズ別全メトリクス、機械可読 |
| `stormprobe_report_<timestamp>.html` | SVGレイテンシチャート付きスタンドアロンダークダッシュボード |

---

## 免責事項

このツールは**許可されたセキュリティテストおよびパフォーマンス評価のみ**を目的としています。負荷テストを実行する前に、対象システムの所有者から明示的な書面による許可を取得する必要があります。著者はこのツールの誤用または損害に対して**一切の責任を負いません**。

---

## ライセンス

[MIT](../../LICENSE) © [Umut ÖZEN](https://github.com/umutozen)

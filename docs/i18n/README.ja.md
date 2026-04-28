# StormProbe

<p align="center">
  <img src="../../assets/banner.png" alt="StormProbe Banner" width="100%">
</p>

<h1 align="center">StormProbe</h1>

<p align="center">
  <strong>crawl · discover · stress · decide</strong><br>
  エンドポイントの自動検出と本番対応評価を備えた最新のHTTP負荷テスト。
</p>

<p align="center">
  <strong>単なるベンチマークではありません。</strong><br>
  StormProbeはテストだけではありません。システムが本当に安全になったかどうかを説明します。
</p>

---

## StormProbeの理由

従来のテストツールは数字を表示します。<br>StormProbeは決断を下します。

エンドポイントの手動選択や生データの解釈を行う代わりに、StormProbeは自動的にエンドポイントを検出し、段階的な負荷テストを実行し、行動可能な提案を生成します。

### 従来ツールとの比較

| 従来のツール | StormProbe |
|---|---|
| 手動エンドポイント設定 | 稼働中のエンドポイントを自動検出 |
| 生ベンチマーク出力 | 最終評価とボトルネック診断 |
| 手動での解釈が必要 | 実用的な提案 |
| 単一のストレステスト | Ramp-Up + 維持 + スパイク + 復旧 |
| 「何かが遅い」 | 「ここで、こういう理由で壊れました」 |

---

## 何をするのか？

StormProbeはGoで書かれた本番環境向けの負荷テストCLIです。

- エンドポイントの自動検出
- 段階的負荷テスト
- パーセンタイル分析
- 飽和検出
- 復旧確認
- CI/CDアラート
- フルレポート

---

## 実際の出力例

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

## 機能

| 機能 | 詳細 |
|---|---|
| 自動検出 | Katanaクローラー + Httpxプローブ |
| 4段階負荷テスト | ランプアップ → 維持 → スパイク → 回復 |
| 評価エンジン | 安全な同時実行数、劣化ポイント、ボトルネックの特定 |
| 豊富な指標 | Avg / P50 / P95 / P99 待機時間 + req/s + エラー分類 |
| 処理能力分析 | 深刻な障害が起きる前に飽和を検出 |
| 回復検証 | スパイク後にシステムが回復するかを確認 |
| アラート閾値 | レイテンシ、エラー率、RPSに基づくCI/CD対応の終了コード |
| デュアルレポート | JSON + 独立したHTMLダッシュボード |
| カスタムヘッダー | 認証トークン、Cookieの設定 |
| Docker対応 | 全ツール内蔵の単一イメージ |
| クロスプラットフォーム | Linux, macOS, Windows — amd64 & arm64 |
| 純粋なGoコア | インストールが速く、運用が容易 |

---

## クイックスタート

```bash
# インストール
go install github.com/umutozen/stormprobe/cmd/stormprobe@latest

# フルテスト
stormprobe https://example.com

# 検出をスキップ
stormprobe --no-discovery https://example.com

# Authありテスト
stormprobe --insecure \
  -H "Authorization: Bearer TOKEN" \
  https://api.example.com

# CI/CDゲート
stormprobe \\
  --alert-p99 500 \\
  --alert-error-rate 5 \\
  https://example.com

# JSON設定ファイルを使用
stormprobe --config stormprobe.json
```

---

## テストの流れ

```text
検出 → ランプアップ → 維持 → スパイク → 回復 → 評価
```

| 段階 | 目的 |
|---|---|
| Discovery | URLのクロール |
| Ramp-Up | 負荷を緩やかに上げる |
| Sustained | 安定性を確認 |
| Spike | 突発的なアクセス |
| Recovery | 負荷後の状態 |
| Verdict | 最終分析 |

---

## インストール

### Go Install

```bash
go install github.com/umutozen/stormprobe/cmd/stormprobe@latest
```

推奨される方法です。

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

パフォーマンスが低下した場合、デプロイが失敗するようにします。

---

## 使用例

**リリース前** — スループット目標を確認

**セキュリティとパフォーマンス** — 隠しAPIの負荷テスト

**監査** — 飽和ポイントの測定

**回帰検出** — デプロイ前後の比較

---

## Flags

| Flag | 説明 |
|---|---|
| `--duration` | テスト時間(例 30s) |
| `--alert-p99` | P99レイテンシ許容値(ms) |
| `--alert-error-rate` | エラー率閾値(%) |
| `--alert-rps` | RPS閾値 |
| `-H` / `--header` | カスタムHTTPヘッダ |
| `--insecure` | TLSチェックをスキップ |
| `--no-discovery` | Rootのみテスト |
| `--endpoints` | エンドポイント一覧 |
| `--format` | 出力フォーマット |
| `--concurrency-ramp` | ランプアップ時の最大VU |
| `--concurrency-sustained` | 維持時のVU |
| `--concurrency-spike` | スパイク時のVU |
| `--req-per-worker` | 1Workerあたりのリクエスト |
| `--timeout` | タイムアウト秒数 |
| `--config` | JSON設定ファイルから設定を読み込む |
| `--auto-profile` | サーバースタックを自動検出し診断プロファイルを適用（デフォルト: 有効） |
| `--profile` | プロファイル指定: iis, tomcat, php |
| `--no-fingerprint` | サーバースタック検出をスキップ |
| `--show-profile` | 検出スタックとプロファイルを表示し終了 |

---

## 要件

| ツール | 備考 |
|---|---|
| Go 1.21+ | Go環境 |
| Katana | 推奨 |
| Httpx | 推奨 |
| Docker | Docker向け |

---

## 免責事項

当ツールは認可されたテスト専用です。

必ず許可を得てから実行してください。

無断使用は法律で禁止されています。

---

## 貢献

CONTRIBUTING.md を参照。

## セキュリティ

SECURITY.md を参照。

## ライセンス

MIT License © Umut ÖZEN

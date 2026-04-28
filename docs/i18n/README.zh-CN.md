# StormProbe

<p align="center">
  <img src="../../assets/banner.png" alt="StormProbe Banner" width="100%">
</p>

<h1 align="center">StormProbe</h1>

<p align="center">
  <strong>crawl · discover · stress · decide</strong><br>
  现代 HTTP 压力测试工具，支持自主 endpoint 发现与生产就绪性裁决。
</p>

<p align="center">
  <strong>不只是基准测试。</strong><br>
  StormProbe 不只是测试。它解释您的系统是否真正变得更安全。
</p>

---

## 为什么选择 StormProbe？

传统压测工具只显示数字。<br>StormProbe 给出决策。

无需手动选择 endpoint，StormProbe 自动发现活跃端点，执行分阶段压力测试，并生成可操作的最终裁决。

### 传统工具 vs StormProbe

| 传统工具 | StormProbe |
|---|---|
| 手动选择 endpoint | 自动发现活跃 endpoint |
| 原始输出 | 最终裁决 + 瓶颈诊断 |
| 需要解读 | 操作建议 |
| 单一模式 | 递增 + 保持 + 峰值 + 恢复 |
| "系统慢" | "这是根本原因" |

---

## 功能说明

StormProbe 是用 Go 编写的生产级 HTTP 压测工具。

- 自动发现
- 分阶段测试
- 延迟百分位
- 吞吐量分析
- 恢复验证
- CI/CD告警
- 报告格式

---

## 真实输出

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

## 特性

| 特性 | 详情 |
|---|---|
| 自动发现 | Katana crawl + Httpx probe，并对活跃端点去重 |
| 4阶段负载测试 | Ramp-Up → Sustained → Spike → Recovery |
| 裁决引擎 | 安全并发，性能退化点，瓶颈分类 |
| 丰富指标 | 平均 / P50 / P95 / P99 延迟 + 每秒请求 + 错误分类 |
| 吞吐量分析 | 在严重故障发生前检测系统饱和点 |
| 恢复验证 | 确认系统在流量飙升后是否能够恢复 |
| 告警阈值 | 针对延迟、错误率和 RPS 的 CI/CD 就绪退出代码 |
| 双重报告 | JSON + 独立的 HTML 仪表盘 |
| 自定义 Headers | 授权令牌、租户请求头、Cookies |
| Docker 就绪 | 包含捆绑工具的单一镜像 |
| 跨平台支持 | Linux, macOS, Windows — amd64 & arm64 |
| 纯 Go 核心 | 快速安装，最小的运维开销 |

---

## 快速开始

```bash
# 安装
go install github.com/umutozen/stormprobe/cmd/stormprobe@latest

# 全自动测试
stormprobe https://example.com

# 跳过发现
stormprobe --no-discovery https://example.com

# Token 测试
stormprobe --insecure \
  -H "Authorization: Bearer TOKEN" \
  https://api.example.com

# CI 自动化
stormprobe \\
  --alert-p99 500 \\
  --alert-error-rate 5 \\
  https://example.com

# 使用 JSON 配置文件
stormprobe --config stormprobe.json
```

---

## 测试流程

```text
发现阶段 → 负载渐增 → 维持负载 → 流量高峰 → 恢复验证 → 最终裁决
```

| 阶段 | 目的 |
|---|---|
| Discovery | 验证 API |
| Ramp-Up | 压力预热 |
| Sustained | 稳定分析 |
| Spike | 高并发模拟 |
| Recovery | 宕机验证 |
| Verdict | 产出结论 |

---

## 安装

### Go Install

```bash
go install github.com/umutozen/stormprobe/cmd/stormprobe@latest
```

推荐。

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

可以在性能下降时自动阻断发布流程。

---

## 场景

**上线前** — 验证延迟目标。

**安全性** — 隐蔽接口发现。

**架构审查** — 检测饱和点。

**回归测试** — CI/CD性能。

---

## Flags

| Flag | 描述 |
|---|---|
| `--duration` | 单阶段持续时间 |
| `--alert-p99` | P99 告警阈值 |
| `--alert-error-rate` | 错误率告警 |
| `--alert-rps` | 最低 RPS |
| `-H` / `--header` | 请求头 |
| `--insecure` | 忽略证书 |
| `--no-discovery` | 仅测试根路径 |
| `--endpoints` | 接口列表文件 |
| `--format` | 输出格式 |
| `--concurrency-ramp` | 峰值 VU |
| `--concurrency-sustained` | 稳定 VU |
| `--concurrency-spike` | 最大并发 VU |
| `--req-per-worker` | 单个线程请求数 |
| `--timeout` | 超时时间 |
| `--config` | 从 JSON 配置文件加载设置 |
| `--auto-profile` | 自动检测服务器栈并应用诊断配置文件（默认：启用） |
| `--profile` | 手动配置文件：iis、tomcat、php |
| `--no-fingerprint` | 跳过服务器栈检测 |
| `--show-profile` | 显示检测到的栈和配置文件，然后退出 |

---

## 前置要求

| 组件 | 声明 |
|---|---|
| Go 1.21+ | Go环境 |
| Katana | 可选 |
| Httpx | 可选 |
| Docker | Docker |

---

## 声明

仅允许授权攻击。

请求授权。

违例者法办。

---

## 贡献

查看代码。

## 安全

请报告漏洞。

## 许可

MIT License © Umut ÖZEN

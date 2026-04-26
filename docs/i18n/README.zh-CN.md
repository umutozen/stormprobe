[English](../../README.md) · [Türkçe](README.tr.md) · [Deutsch](README.de.md) · [Francais](README.fr.md) · [Espanol](README.es.md) · [Italiano](README.it.md) · [Portugues](README.pt-BR.md) · [Russkij](README.ru.md) · [Zhongwen Jian](README.zh-CN.md) · [Zhongwen Fan](README.zh-TW.md) · [Nihongo](README.ja.md) · [Hangugeo](README.ko.md) · [Al-Arabiyya](README.ar.md) · [Hindi](README.hi.md) · [Nederlands](README.nl.md) · [Polski](README.pl.md) · [Ukrainska](README.uk.md)

<p align="center">
  <img src="../../assets/banner.png" alt="StormProbe Banner" width="100%">
</p>

<h1 align="center">StormProbe</h1>

<p align="center">
  <strong>faxian · jiazai · fenxi</strong><br>
  具有自主端点发现功能的生产级HTTP负载测试工具。零外部依赖。
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
    <img src="https://img.shields.io/badge/license-MIT-blue.svg" alt="MIT许可证">
  </a>
</p>

---

## 什么是StormProbe？

StormProbe是一个用纯Go编写的**零依赖**HTTP负载测试工具。它使用[Katana](https://github.com/projectdiscovery/katana)和[Httpx](https://github.com/projectdiscovery/httpx)自动发现应用程序的端点，然后执行结构化的4阶段负载测试——**Ramp-Up → Sustained → Spike → Recovery**——并生成包含P50/P95/P99延迟指标的JSON和HTML报告。

适用于：
- **CI/CD管道** — 带有退出代码的警报阈值
- **渗透测试人员** — 自动发现活跃端点
- **DevOps工程师** — 部署前后的基准测试
- **QA团队** — 性能SLA验证

---

## 功能

| 功能 | 详情 |
|---|---|
| **自动发现** | Katana爬虫 + Httpx探测，删重后的活跃端点列表 |
| **4阶段负载测试** | Ramp-Up → Sustained → Spike → Recovery |
| **丰富指标** | P50/P95/P99延迟、req/s、按阶段分类错误 |
| **持续时间模式** | 基于时间的阶段 (`--duration 30s`) 替代请求计数 |
| **警报阈值** | P99、错误率或RPS违规时触发CI/CD兼容的`exit 1` |
| **双重报告** | JSON（机器可读）+ 独立HTML暗色仪表盘 |
| **自定义标头** | Bearer令牌、租户标头、Cookie — 全部传播 |
| **Docker就绪** | 内置Katana + Httpx的单一镜像 |
| **跨平台** | Linux、macOS、Windows — amd64 & arm64 |
| **零依赖** | 纯Go 1.21+ stdlib，`go install`即可完成 |

---

## 快速开始

```bash
# 安装（需要Go 1.21+）
go install github.com/umutozen/stormprobe/cmd/stormprobe@latest

# 基本测试 — 自动发现，4个阶段
stormprobe https://example.com

# 跳过发现，仅测试根路径
stormprobe --no-discovery https://example.com

# TLS + 认证标头 + 仅JSON报告
stormprobe --insecure -H "Authorization: Bearer TOKEN" --format json https://api.example.com

# 基于时间的测试：每个阶段30秒
stormprobe --insecure --duration 30s --no-discovery https://example.com

# CI/CD：如果P99 > 500ms或错误 > 5%则exit 1
stormprobe --alert-p99 500 --alert-error-rate 5 https://example.com
```

---

## 安装

### 选项1 — `go install` *(推荐)*

```bash
go install github.com/umutozen/stormprobe/cmd/stormprobe@latest
stormprobe --help
```

### 选项2 — 预构建二进制文件 *(无需Go)*

从[Releases页面](https://github.com/umutozen/stormprobe/releases)下载。

```bash
# Linux / macOS
chmod +x stormprobe
sudo mv stormprobe /usr/local/bin/

# Windows PowerShell
Move-Item stormprobe.exe C:\Windows\System32\stormprobe.exe
```

### 选项3 — 从源码构建

```bash
git clone https://github.com/umutozen/stormprobe.git
cd stormprobe
go build -o stormprobe ./cmd/stormprobe      # Linux / macOS
go build -o stormprobe.exe ./cmd/stormprobe  # Windows
```

### 选项4 — Docker

```bash
docker run --rm ghcr.io/umutozen/stormprobe:latest --insecure https://example.com

# Linux / macOS
docker run --rm -v $(pwd)/outputs:/app/outputs \
  ghcr.io/umutozen/stormprobe:latest --insecure --format both https://example.com

# Windows PowerShell
docker run --rm -v "${PWD}\outputs:/app/outputs" `
  ghcr.io/umutozen/stormprobe:latest --insecure --format both https://example.com
```

### 可选 — 发现工具

```bash
bash scripts/install-tools.sh
```

> 没有Katana/Httpx时使用`--no-discovery`或`--endpoints`。

---

## CLI参考

### 并发与负载

| 标志 | 默认值 | 说明 |
|---|---|---|
| `--concurrency-ramp` | `50` | Ramp-Up阶段的峰值并发 |
| `--concurrency-sustained` | `75` | Sustained阶段的并发 |
| `--concurrency-spike` | `250` | Spike阶段的峰值并发 |
| `--req-per-worker` | `15` | 每个worker的请求数（设置`--duration`时忽略） |
| `--duration` | `0` (禁用) | 基于时间的阶段持续时间，如`30s`、`1m` |
| `--timeout` | `10s` | 每个请求的超时 |

### 发现

| 标志 | 默认值 | 说明 |
|---|---|---|
| `--no-discovery` | `false` | 跳过Katana+Httpx，仅测试根路径 |
| `--endpoints` | `""` | 从文件加载端点 |
| `--katana-path` | `""` | 自定义Katana二进制路径 |
| `--httpx-path` | `""` | 自定义Httpx二进制路径 |

### HTTP与安全

| 标志 | 默认值 | 说明 |
|---|---|---|
| `-H`, `--header` | — | 自定义标头（可重复） |
| `--insecure` | `false` | 跳过TLS证书验证 |

### 输出

| 标志 | 默认值 | 说明 |
|---|---|---|
| `--format` | `both` | 报告格式：`json`、`html`、`both` |
| `--output` | `./outputs` | 报告的输出目录 |

### CI/CD警报

| 标志 | 默认值 | 说明 |
|---|---|---|
| `--alert-p99` | `0` (禁用) | P99(ms)超出阈值时exit 1 |
| `--alert-error-rate` | `0` (禁用) | 错误率(%)超出阈值时exit 1 |
| `--alert-rps` | `0` (禁用) | req/s低于阈值时exit 1 |

---

## 测试阶段

| 阶段 | 说明 |
|---|---|
| **0 — Discovery** | Katana爬虫 + Httpx探测，去重列表 |
| **1 — Ramp-Up** | 逐步增加并发：5 → 15 → 30 → 峰值 |
| **2 — Sustained** | 在目标并发下进行3次冲击 |
| **3 — Spike** | 突然增加到峰值并发，然后冷却 |
| **4 — Recovery** | 10秒冷却后的峰值后健康检查 |

---

## 输出与报告

| 文件 | 说明 |
|---|---|
| `stormprobe_report_<timestamp>.json` | 每个阶段的所有指标，机器可读 |
| `stormprobe_report_<timestamp>.html` | 带SVG延迟图表的独立暗色仪表盘 |

---

## 免责声明

本工具仅用于**授权的安全测试和性能评估**。在对目标系统执行任何负载测试之前，必须获得系统所有者的明确书面许可。作者对本工具的误用或造成的任何损害**不承担任何责任**。

---

## 许可证

[MIT](../../LICENSE) © [Umut ÖZEN](https://github.com/umutozen)

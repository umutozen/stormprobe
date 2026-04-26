[English](../../README.md) · [Türkçe](README.tr.md) · [Deutsch](README.de.md) · [Francais](README.fr.md) · [Espanol](README.es.md) · [Italiano](README.it.md) · [Portugues](README.pt-BR.md) · [Russkij](README.ru.md) · [Zhongwen Jian](README.zh-CN.md) · [Zhongwen Fan](README.zh-TW.md) · [Nihongo](README.ja.md) · [Hangugeo](README.ko.md) · [Al-Arabiyya](README.ar.md) · [Hindi](README.hi.md) · [Nederlands](README.nl.md) · [Polski](README.pl.md) · [Ukrainska](README.uk.md)

<p align="center">
  <img src="../../assets/banner.png" alt="StormProbe Banner" width="100%">
</p>

<h1 align="center">StormProbe</h1>

<p align="center">
  <strong>balgyon · buha · bunseok</strong><br>
  Jayul endpoint balgyon gineung-i namdeon saengsan-geup HTTP buhateseuteu doju. Oebujong-sonseong jero.
</p>

<p align="center">
  <a href="https://github.com/umutozen/stormprobe/actions/workflows/ci.yml">
    <img src="https://github.com/umutozen/stormprobe/actions/workflows/ci.yml/badge.svg" alt="CI">
  </a>
  <a href="https://goreportcard.com/report/github.com/umutozen/stormprobe">
    <img src="https://goreportcard.com/badge/github.com/umutozen/stormprobe" alt="Go Report Card">
  </a>
  <a href="https://github.com/umutozen/stormprobe/releases">
    <img src="https://img.shields.io/github/v/release/umutozen/stormprobe" alt="Choesinsillion">
  </a>
  <a href="../../LICENSE">
    <img src="https://img.shields.io/badge/license-MIT-blue.svg" alt="MIT License">
  </a>
</p>

---

## StormProbe란?

StormProbe는 순수 Go로 작성된 **의존성 제로** HTTP 부하 테스트 도구입니다. [Katana](https://github.com/projectdiscovery/katana)와 [Httpx](https://github.com/projectdiscovery/httpx)를 사용하여 애플리케이션의 엔드포인트를 자동으로 발견하고, **Ramp-Up → Sustained → Spike → Recovery**라는 체계적인 4단계 부하 테스트를 실행하며, P50/P95/P99 레이턴시 메트릭이 포함된 JSON 및 HTML 보고서를 생성합니다.

대상 사용자:
- **CI/CD 파이프라인** — 종료 코드를 포함한 알림 임계값
- **침투 테스터** — 라이브 엔드포인트 자동 발견
- **DevOps 엔지니어** — 배포 전후 벤치마크
- **QA 팀** — 성능 SLA 검증

---

## 기능

| 기능 | 세부 정보 |
|---|---|
| **자동 발견** | Katana 크롤 + Httpx 프로브, 중복 제거된 라이브 엔드포인트 목록 |
| **4단계 부하 테스트** | Ramp-Up → Sustained → Spike → Recovery |
| **풍부한 메트릭** | P50/P95/P99 레이턴시, req/s, 단계별 오류 분류 |
| **시간 모드** | 요청 수 대신 시간 기반 단계 (`--duration 30s`) |
| **알림 임계값** | P99, 오류율, RPS 위반 시 CI/CD 호환 `exit 1` |
| **이중 보고서** | JSON (기계 판독 가능) + 독립형 HTML 다크 대시보드 |
| **사용자 정의 헤더** | Bearer 토큰, 테넌트 헤더, 쿠키 — 모든 곳에 전파 |
| **Docker 준비** | Katana + Httpx가 내장된 단일 이미지 |
| **크로스 플랫폼** | Linux, macOS, Windows — amd64 & arm64 |
| **의존성 제로** | 순수 Go 1.21+ stdlib, `go install`만으로 충분 |

---

## 빠른 시작

```bash
# 설치 (Go 1.21+ 필요)
go install github.com/umutozen/stormprobe/cmd/stormprobe@latest

# 기본 테스트 — 자동 발견, 4단계 실행
stormprobe https://example.com

# 발견 건너뛰기, 루트 경로만 테스트
stormprobe --no-discovery https://example.com

# TLS + 인증 헤더 + JSON 보고서만
stormprobe --insecure -H "Authorization: Bearer TOKEN" --format json https://api.example.com

# 시간 기반 테스트: 단계당 30초
stormprobe --insecure --duration 30s --no-discovery https://example.com

# CI/CD: P99 > 500ms 또는 오류 > 5%이면 exit 1
stormprobe --alert-p99 500 --alert-error-rate 5 https://example.com
```

---

## 설치

### 옵션 1 — `go install` *(권장)*

```bash
go install github.com/umutozen/stormprobe/cmd/stormprobe@latest
stormprobe --help
```

### 옵션 2 — 빌드된 바이너리 *(Go 불필요)*

[Releases 페이지](https://github.com/umutozen/stormprobe/releases)에서 다운로드.

```bash
# Linux / macOS
chmod +x stormprobe
sudo mv stormprobe /usr/local/bin/

# Windows PowerShell
Move-Item stormprobe.exe C:\Windows\System32\stormprobe.exe
```

### 옵션 3 — 소스에서 빌드

```bash
git clone https://github.com/umutozen/stormprobe.git
cd stormprobe
go build -o stormprobe ./cmd/stormprobe      # Linux / macOS
go build -o stormprobe.exe ./cmd/stormprobe  # Windows
```

### 옵션 4 — Docker

```bash
docker run --rm ghcr.io/umutozen/stormprobe:latest --insecure https://example.com

# Linux / macOS
docker run --rm -v $(pwd)/outputs:/app/outputs \
  ghcr.io/umutozen/stormprobe:latest --insecure --format both https://example.com

# Windows PowerShell
docker run --rm -v "${PWD}\outputs:/app/outputs" `
  ghcr.io/umutozen/stormprobe:latest --insecure --format both https://example.com
```

### 선택 사항 — 발견 도구

```bash
bash scripts/install-tools.sh
```

> Katana/Httpx 없이는 `--no-discovery` 또는 `--endpoints` 사용.

---

## CLI 참조

### 동시성 및 부하

| 플래그 | 기본값 | 설명 |
|---|---|---|
| `--concurrency-ramp` | `50` | Ramp-Up 단계의 최대 동시성 |
| `--concurrency-sustained` | `75` | Sustained 단계의 동시성 |
| `--concurrency-spike` | `250` | Spike 단계의 최대 동시성 |
| `--req-per-worker` | `15` | 워커당 요청 수 (`--duration` 설정 시 무시) |
| `--duration` | `0` (비활성) | 시간 기반 단계 지속 시간, 예: `30s`, `1m` |
| `--timeout` | `10s` | 요청당 타임아웃 |

### 발견

| 플래그 | 기본값 | 설명 |
|---|---|---|
| `--no-discovery` | `false` | Katana+Httpx 건너뛰기, 루트만 테스트 |
| `--endpoints` | `""` | 파일에서 엔드포인트 로드 |
| `--katana-path` | `""` | 사용자 정의 Katana 바이너리 경로 |
| `--httpx-path` | `""` | 사용자 정의 Httpx 바이너리 경로 |

### HTTP 및 보안

| 플래그 | 기본값 | 설명 |
|---|---|---|
| `-H`, `--header` | — | 사용자 정의 헤더 (반복 가능) |
| `--insecure` | `false` | TLS 인증서 검증 건너뛰기 |

### 출력

| 플래그 | 기본값 | 설명 |
|---|---|---|
| `--format` | `both` | 보고서 형식: `json`, `html`, `both` |
| `--output` | `./outputs` | 보고서 출력 디렉토리 |

### CI/CD 알림

| 플래그 | 기본값 | 설명 |
|---|---|---|
| `--alert-p99` | `0` (비활성) | P99 (ms)가 임계값 초과 시 exit 1 |
| `--alert-error-rate` | `0` (비활성) | 오류율(%)이 임계값 초과 시 exit 1 |
| `--alert-rps` | `0` (비활성) | req/s가 임계값 미만 시 exit 1 |

---

## 테스트 단계

| 단계 | 설명 |
|---|---|
| **0 — Discovery** | Katana 크롤 + Httpx 프로브, 중복 제거 목록 |
| **1 — Ramp-Up** | 점진적 동시성 증가: 5 → 15 → 30 → 피크 |
| **2 — Sustained** | 목표 동시성에서 3개 웨이브 |
| **3 — Spike** | 피크로의 급격한 증가, 이후 쿨다운 |
| **4 — Recovery** | 10초 쿨다운 후 상태 확인 |

---

## 출력 및 보고서

| 파일 | 설명 |
|---|---|
| `stormprobe_report_<timestamp>.json` | 단계별 모든 메트릭, 기계 판독 가능 |
| `stormprobe_report_<timestamp>.html` | SVG 레이턴시 차트가 있는 독립형 다크 대시보드 |

---

## 면책 조항

이 도구는 **승인된 보안 테스트 및 성능 평가 전용**입니다. 부하 테스트를 실행하기 전에 대상 시스템 소유자로부터 명시적인 서면 허가를 받아야 합니다. 저자는 이 도구의 오용 또는 피해에 대해 **어떠한 책임도 지지 않습니다**.

---

## 라이선스

[MIT](../../LICENSE) © [Umut ÖZEN](https://github.com/umutozen)

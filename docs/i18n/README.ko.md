# StormProbe

<p align="center">
  <img src="../../assets/banner.png" alt="StormProbe Banner" width="100%">
</p>

<h1 align="center">StormProbe</h1>

<p align="center">
  <strong>crawl · discover · stress · decide</strong><br>
  엔드포인트 자동화 및 프로덕션 평가가 포함된 모던 HTTP 부하 테스트.
</p>

<p align="center">
  <strong>단순한 벤치마킹이 아닙니다.</strong><br>
  StormProbe는 단순히 테스트하지 않습니다. 시스템이 실제로 더 안전해졌는지 설명합니다.
</p>

---

## 왜 StormProbe 인가요?

기존 도구는 숫자를 보여줍니다.<br>StormProbe는 결정을 내립니다.

StormProbe는 자동으로 엔드포인트를 찾고 부하 테스트를 실행합니다.

### 기존 도구 vs StormProbe

| 기존 도구 | StormProbe |
|---|---|
| 수동 엔드포인트 선택 | 엔드포인트 자동 검색 |
| 원시 벤치마크 데이터 | 최종 판정 및 병목 원인 진단 |
| 수동 해석 필요 | 실행 가능한 조치 권장 |
| 단일 스트레스 패턴 | 증가 + 유지 + 급증 + 회복 |
| "뭔가 느림" | "어디서 왜 망가졌는지 정확히 알려줌" |

---

## 주요 기능

프로덕션 레벨의 HTTP CLI 도구입니다.

- 자동 발견
- 퍼센타일 분석
- 포화도 확인
- HTML 리포트

---

## 출력 예시

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

## 특징

| 특징 | 설명 |
|---|---|
| 자동 검색 | 중복 제거가 포함된 Katana 및 Httpx 기반 검색 |
| 4단계 부하 테스트 | 증가 → 유지 → 급증 → 회복 |
| 판정 엔진 | 안전 동시성, 성능 저하 지점, 병목 분류 |
| 세부 메트릭 | Avg / P50 / P95 / P99 레이턴시 및 초당 요청 수 |
| 처리량 분석 | 하드 다운 전에 포화 지점을 사전에 감지 |
| 경고 임계값 | CI/CD 파이프라인과 통합 가능한 종료 코드 |
| 이중 리포트 | JSON과 HTML 대시보드 |
| Docker | 모든 도구가 포함된 단일 이미지 |

---

## 빠른 시작

```bash
# 설치
go install github.com/umutozen/stormprobe/cmd/stormprobe@latest

# 완전 탐색 테스트
stormprobe https://example.com

# 탐색 건너뛰기
stormprobe --no-discovery https://example.com

# Auth
stormprobe --insecure \
  -H "Authorization: Bearer TOKEN" \
  https://api.example.com

# CI/CD 게이트
stormprobe \\
  --alert-p99 500 \\
  --alert-error-rate 5 \\
  https://example.com

# JSON 설정 파일 사용
stormprobe --config stormprobe.json
```

---

## 테스트 흐름

```text
검색 → 증가 → 유지 → 급증 → 회복 → 판정
```

| 단계 | 목표 |
|---|---|
| Discovery | URL 수집 |
| Ramp-Up | 부하 서서히 증가 |
| Sustained | 성능 검증 |
| Spike | 트래픽 피크 |
| Recovery | 회복 검증 |
| Verdict | 최종 분석 |

---

## Installation

### Go Install

```bash
go install github.com/umutozen/stormprobe/cmd/stormprobe@latest
```

권장되는 방식입니다.

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

성능 저하 발생 시 데플로이를 자동 취소합니다.

---

## 활용 시나리오

**배포 전** — 지연 시간 목표를 확인합니다.

**보안+부하** — 숨겨진 API 부하 테스트.

**인프라 감사** — 포화 구간 측정.

**회귀 탐지** — 배포 전후 차이 확인.

---

## Flags

| Flag | 설명 |
|---|---|
| `--duration` | 수행 시간 |
| `--alert-p99` | P99 레이턴시 경고 (ms) |
| `--alert-error-rate` | 에러 비율 |
| `--alert-rps` | 요청 감소 |
| `-H` / `--header` | 헤더 |
| `--insecure` | TLS 무시 |
| `--no-discovery` | Root 라우터 테스트 |
| `--endpoints` | 엔드포인트 목록 |
| `--format` | 포맷 |
| `--concurrency-ramp` | Ramp 가상유저 |
| `--concurrency-sustained` | Sustained 가상유저 |
| `--concurrency-spike` | Spike 가상유저 |
| `--req-per-worker` | 쓰레드당 요청 수 |
| `--timeout` | 타임아웃 |
| `--config` | JSON 설정 파일에서 설정 로드 |
| `--auto-profile` | 서버 스택 자동 감지 및 진단 프로파일 적용 (기본: 활성화) |
| `--profile` | 수동 프로파일: iis, tomcat, php |
| `--no-fingerprint` | 서버 스택 감지 건너뛰기 |
| `--show-profile` | 감지된 스택과 프로파일 표시 후 종료 |

---

## 요구사항

| 종류 | 비고 |
|---|---|
| Go 1.21+ | 필수 |
| Katana | 옵션 |
| Httpx | 옵션 |
| Docker | Docker용 |

---

## 경고

권한이 있는 대상에게만 사용하세요.

사전에 반드시 허가를 받아야 합니다.

불법적인 사용을 금합니다.

---

## 기여

문서를 참조하세요.

## 보안

안내문 확인.

## 라이선스

MIT License © Umut ÖZEN

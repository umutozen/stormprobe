# StormProbe

<p align="center">
  <img src="../../assets/banner.png" alt="StormProbe Banner" width="100%">
</p>

<h1 align="center">StormProbe</h1>

<p align="center">
  <strong>crawl · discover · stress · decide</strong><br>
  اختبار تحميل HTTP مع اكتشاف روابط وقرارات
</p>

<p align="center">
  <strong></strong><br>
  لا يقتصر StormProbe على الاختبار فقط. بل يوضح ما إذا كان نظامك أصبح أكثر أمانًا بالفعل.
</p>

---

## لماذا؟

يقدم حلول.

يكتشف ويحلل.

### المقارنة

| تقليدي | بطل |
|---|---|
| - | آلي |
| - | - |
| - | - |
| - | - |
| - | - |

---

## المهام

-

-

---

## -

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

## المميزات

|Feature|D|
|-|-|

---

## -

```bash
# 1
go install github.com/umutozen/stormprobe/cmd/stormprobe@latest

# 2
stormprobe https://example.com

# 3
stormprobe --no-discovery https://example.com

# 4
stormprobe --insecure \
  -H "Authorization: Bearer TOKEN" \
  https://api.example.com

# 5
stormprobe \\
  --alert-p99 500 \\
  --alert-error-rate 5 \\
  https://example.com

# استخدام ملف تكوين JSON
stormprobe --config stormprobe.json
```

---

## -

```text
-
```

| - | - |
|---|---|
| Discovery | 1 |
| Ramp-Up | 2 |
| Sustained | 3 |
| Spike | 4 |
| Recovery | 5 |
| Verdict | 6 |

---

## -

### Go Install

```bash
go install github.com/umutozen/stormprobe/cmd/stormprobe@latest
```

-

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

-

---

## -

**-** — -

**-** — -

**-** — -

**-** — -

---

## Flags

| Flag | - |
|---|---|
| `--duration` | - |
| `--alert-p99` | - |
| `--alert-error-rate` | - |
| `--alert-rps` | - |
| `-H` / `--header` | - |
| `--insecure` | - |
| `--no-discovery` | - |
| `--endpoints` | - |
| `--format` | - |
| `--concurrency-ramp` | - |
| `--concurrency-sustained` | - |
| `--concurrency-spike` | - |
| `--req-per-worker` | - |
| `--timeout` | - |
| `--config` | تحميل الإعدادات من ملف تكوين JSON |
| `--auto-profile` | اكتشاف الحزمة تلقائياً وتطبيق بروفيل التشخيص |
| `--profile` | بروفيل يدوي: iis, tomcat, php |
| `--no-fingerprint` | تخطي اكتشاف الحزمة |
| `--show-profile` | عرض الحزمة والبروفيل ثم الخروج |

---

## -

| - | - |
|---|---|
| Go 1.21+ | - |
| Katana | - |
| Httpx | - |
| Docker | - |

---

## -

-

-

-

---

## -

-

## -

-

## -

MIT License © Umut ÖZEN

# StormProbe

<p align="center">
  <img src="../../assets/banner.png" alt="StormProbe Banner" width="100%">
</p>

<h1 align="center">StormProbe</h1>

<p align="center">
  <strong>crawl · discover · stress · decide</strong><br>
  Nowoczesne testy obciążeniowe HTTP z autonomicznym wykrywaniem punktów końcowych i werdyktami gotowości produkcyjnej.
</p>

<p align="center">
  <strong>Nie tylko benchmarking.</strong><br>
  StormProbe nie tylko testuje. Wyjaśnia, czy Twój system naprawdę stał się bezpieczniejszy.
</p>

---

## Dlaczego StormProbe?

Tradycyjne testery pokazują liczby.<br>StormProbe podejmuje decyzje.

Pomija zgadywanie: sam szuka URL-i, testuje, wykrywa załamania i daje rekomendacje.

### Klasyczne aplikacje vs StormProbe

| Klasyczne aplikacje | StormProbe |
|---|---|
| Ręczny dobór tras | Autonomiczne skanowanie |
| Surowe liczby | Werdykt + diagnoza problemu |
| Trudna interpretacja | Użyteczne zalecenia |
| Stres test na płasko | Ramp-Up + Sustained + Spike + Recovery |
| "Coś zwolniło" | "Wąskie gardło objawiło się tutaj" |

---

## Co on robi?

Narzędzie HTTP w Go stworzone do produkcji.

- Testy w 4 fazach
- Wykrywanie nasycenia
- Weryfikacja powrotu
- Raporty JSON / HTML

---

## Przykład z konsoli

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

## Funkcje

| Funkcja | Szczegóły |
|---|---|
| Auto Discovery | Katana + Httpx |
| Test 4 Fazowy | Ramp-Up → Sustained → Spike → Recovery |
| Silnik Oceniający | Diagnoza i rekomendacje |
| Bogate Metryki | P50 / P95 / P99 |
| Analiza Przepustowości | Sprawdzanie punktu załamania |
| Limity i Alarmy | Integracja CI/CD z exit-codes |
| Raporty | JSON / HTML |
| Header | Pamięć cookies i auth |
| Multiplatform | Win, Linux, Mac |

---

## Szybki start

```bash
# Instalacja
go install github.com/umutozen/stormprobe/cmd/stormprobe@latest

# Pełny test
stormprobe https://example.com

# Test wyrywkowy root
stormprobe --no-discovery https://example.com

# Auth
stormprobe --insecure \
  -H "Authorization: Bearer TOKEN" \
  https://api.example.com

# Do CI/CD
stormprobe \\
  --alert-p99 500 \\
  --alert-error-rate 5 \\
  https://example.com

# Użyj pliku konfiguracyjnego JSON
stormprobe --config stormprobe.json
```

---

## Cykl operacyjny

```text
Szukaj → Wznoś → Utrzymaj → Strzelaj → Regeneruj → Oceń
```

| Faza | Cel |
|---|---|
| Discovery | Zbieraj Endpointy |
| Ramp-Up | Łagodne wznoszenie limitu |
| Sustained | Próba wydajnościowa |
| Spike | Szok |
| Recovery | Stabilizacja potoków |
| Verdict | Konkluzja |

---

## Instalacja

### Go Install

```bash
go install github.com/umutozen/stormprobe/cmd/stormprobe@latest
```

Gotowy i polecany.

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

Do bramkowania regresji wydajności.

---

## Przykłady

**Przed wydaniem** — Czy testy SLI nie popsuły?

**Audyt Security** — Ukryte URL.

**Testowanie Infrastruktury** — Lokalny benchmark.

**Regresja wydajnościowa** — Porównaj logiki biznesowe.

---

## Flags

| Flag | Instrukcja |
|---|---|
| `--duration` | Czas trwania próbki |
| `--alert-p99` | Odetnij pow. P99 (ms) |
| `--alert-error-rate` | Tolerancja błędu (%) |
| `--alert-rps` | Tolerancja minimalnego RPS |
| `-H` / `--header` | Niestandardowy nagłówek |
| `--insecure` | Pomiń cert. |
| `--no-discovery` | Zbadaj tylko początek sieci |
| `--endpoints` | Z listą endpointów |
| `--format` | Pokaż w formacie |
| `--concurrency-ramp` | Concurrency startowa |
| `--concurrency-sustained` | Docelowa |
| `--concurrency-spike` | Atak pik |
| `--req-per-worker` | Zapytań dla wątku |
| `--timeout` | Timeout na request |
| `--config` | Wczytaj ustawienia z pliku konfiguracyjnego JSON |
| `--auto-profile` | Automatyczne wykrywanie stosu serwera i zastosowanie profilu diagnostycznego (domyślnie: aktywne) |
| `--profile` | Ręczny profil: iis, tomcat, php |
| `--no-fingerprint` | Pominąć wykrywanie stosu serwera |
| `--show-profile` | Pokaż wykryty stos i profil, następnie zakończ |

---

## Pakiety C

| Paczka | Stan |
|---|---|
| Go 1.21+ | Go Core |
| Katana | Opcja |
| Httpx | Opcja |
| Docker | Docker |

---

## Polisa

Narzędzie autorskie.

Zapytaj o zgodę atakowanych przed testem.

Nie używaj na cudzych maszynach bo pójdziesz siedzieć.

---

## Współtwórz

Przeczytaj CONTRIBUTING.

## Zgłoś dziurę

SECURITY.md otwarte.

## Prawo MIT

MIT License © Umut ÖZEN

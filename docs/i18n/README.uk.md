# StormProbe

<p align="center">
  <img src="../../assets/banner.png" alt="StormProbe Banner" width="100%">
</p>

<h1 align="center">StormProbe</h1>

<p align="center">
  <strong>crawl · discover · stress · decide</strong><br>
  Сучасне навантажувальне тестування HTTP з автономним виявленням ендпоінтів та оцінками.
</p>

<p align="center">
  <strong>Не просто бенчмаркінг.</strong><br>
  StormProbe не просто тестує. Він пояснює, чи дійсно ваша система стала безпечнішою.
</p>

---

## Чому StormProbe?

Традиційні інструменти показують цифри.<br>StormProbe надає рішення.

Замість ручного вибору ендпоінтів StormProbe автоматично їх виявляє, проводить поетапне стрес-тестування та видає рекомендації.

### Традиційні інструменти

| Традиційні інструменти | StormProbe |
|---|---|
| Ручний вибір ендпоінтів | Автоматичне виявлення активних ендпоінтів |
| Сирий вивід | Фінальний вердикт + діагностика |
| Потребує ручної інтерпретації | Готові рекомендації до дії |
| Один патерн стресу | Нагрівання + Утримання + Пік + Відновлення |
| "Щось працює повільно" | "Ось де і чому все ламається" |

---

## Можливості

StormProbe — це production-grade HTTP CLI для навантажувального тестування.

- Виявлення ендпоінтів
- Тестування по фазам
- Аналіз перцентилів
- Класифікація помилок

---

## Приклад виводу

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

## Особливості

| Функція | Деталі |
|---|---|
| Автовиявлення | Katana crawl + Httpx probe з дедуплікацією активних ендпоінтів |
| 4-фазний тест навантаження | Ramp-Up → Sustained → Spike → Recovery |
| Механізм вердиктів | Безпечне навантаження, точка деградації, класифікація вузьких місць |
| Багаті метрики | Затримка Avg / P50 / P95 / P99 + req/s + класифікація помилок |
| Аналіз пропускної здатності | Виявляє насичення до серйозних збоїв |
| Перевірка відновлення | Підтверджує відновлення системи після сплесків трафіку |
| Пороги сповіщень | Коди помилок для CI/CD щодо затримки, рівня помилок та RPS |
| Подвійні звіти | JSON + незалежний HTML-дашборд |
| Власні заголовки | Токени автентифікації, tenant-заголовки, cookies |
| Готовність до Docker | Єдиний образ із вбудованими інструментами |
| Кросплатформність | Linux, macOS, Windows — amd64 & arm64 |
| Чисте ядро на Go | Швидке встановлення, мінімальні накладні витрати |

---

## Швидкий старт

```bash
# Встановлення
go install github.com/umutozen/stormprobe/cmd/stormprobe@latest

# Повний тест
stormprobe https://example.com

# Без виявлення
stormprobe --no-discovery https://example.com

# Для API
stormprobe --insecure \
  -H "Authorization: Bearer TOKEN" \
  https://api.example.com

# Для CI/CD
stormprobe \\
  --alert-p99 500 \\
  --alert-error-rate 5 \\
  https://example.com

# Використати JSON-файл конфігурації
stormprobe --config stormprobe.json
```

---

## Потік тестування

```text
Виявлення → Розгін → Стійкість → Сплеск → Відновлення → Вердикт
```

| Етап | Мета |
|---|---|
| Discovery | Збір URL |
| Ramp-Up | Навантаження |
| Sustained | Утримання |
| Spike | Пік |
| Recovery | Відновлення |
| Verdict | Вердикт |

---

## Встановлення

### Go Install

```bash
go install github.com/umutozen/stormprobe/cmd/stormprobe@latest
```

Рекомендований спосіб.

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

Дозволяє блокувати розгортання.

---

## Випадки використання

**Перед релізом** — Перевірка SLI

**Безпека** — Аналіз кінцевих точок

**Аудит** — Вузькі місця

**Регресія** — CI/CD

---

## Flags

| Flag | Опис |
|---|---|
| `--duration` | Час (30s) |
| `--alert-p99` | Помилка P99 |
| `--alert-error-rate` | Помилка % |
| `--alert-rps` | Min RPS |
| `-H` / `--header` | Заголовок |
| `--insecure` | Без TLS |
| `--no-discovery` | Лише корінь |
| `--endpoints` | Файл URL |
| `--format` | Формат |
| `--concurrency-ramp` | Пік |
| `--concurrency-sustained` | Стейбл |
| `--concurrency-spike` | Сплеск |
| `--req-per-worker` | Запити з потоку |
| `--timeout` | Таймаут |
| `--config` | Завантажити налаштування з JSON-файлу конфігурації |
| `--auto-profile` | Автоматично визначити стек сервера і застосувати діагностичний профіль (за замовчуванням: активно) |
| `--profile` | Профіль вручну: iis, tomcat, php |
| `--no-fingerprint` | Пропустити визначення стека сервера |
| `--show-profile` | Показати виявлений стек і профіль, потім вийти |

---

## Вимоги

| Вимога | Примітка |
|---|---|
| Go 1.21+ | Go |
| Katana | Скан |
| Httpx | Веб |
| Docker | Докер |

---

## Відмова

Тільки для етичного тестування.

Завжди отримуйте дозвіл.

Інакше незаконно.

---

## Участь

Див CONTRIBUTING.

## Безпека

Повідомте.

## Ліцензія

MIT License © Umut ÖZEN

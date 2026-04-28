# StormProbe

<p align="center">
  <img src="../../assets/banner.png" alt="StormProbe Banner" width="100%">
</p>

<h1 align="center">StormProbe</h1>

<p align="center">
  <strong>crawl · discover · stress · decide</strong><br>
  Тестирование производительности HTTP с автономным обнаружением эндпоинтов и оценкой готовности к release.
</p>

<p align="center">
  <strong>Не просто бенчмаркинг.</strong><br>
  StormProbe не просто тестирует. Он объясняет, стала ли ваша система действительно безопаснее.
</p>

---

## Почему StormProbe?

Обычные инструменты показывают числа.<br>StormProbe выдает решения.

Вместо ручного выбора, StormProbe обнаруживает маршруты, тестирует их и предлагает план действий.

### Традиционные инструменты

| Обычные утилиты | StormProbe |
|---|---|
| Ручной выбор эндпоинтов | Автоматическое обнаружение |
| Сырые цифры | Финальный вердикт + Диагноз узких мест |
| Требует интерпретации | Готовые рекомендации |
| Одиночный шаблон | Ramp-Up + Sustained + Spike + Recovery |
| "Что-то работает медленно" | "Вот где произошел провал" |

---

## Возможности

StormProbe — это production-grade HTTP CLI написанный на Go.

- Обнаружение маршрутов
- Тестирование по фазам
- Анализ перцентилей
- Анализ насыщения
- Проверка восстановления
- Пороги для CI/CD
- Отчеты JSON+HTML

---

## Пример вывода

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

## Функции

| Функция | Детали |
|---|---|
| Автообнаружение | Katana crawl + Httpx probe с дедупликацией живых эндпоинтов |
| 4-фазный нагрузочный тест | Ramp-Up → Sustained → Spike → Recovery |
| Механизм вердиктов | Безопасная нагрузка, точка деградации, классификация узких мест |
| Подробные метрики | Задержка Avg / P50 / P95 / P99 + req/s + классификация ошибок |
| Анализ пропускной способности | Обнаруживает насыщение до серьезных сбоев |
| Проверка восстановления | Подтверждает восстановление системы после скачков трафика |
| Пороги оповещений | Exit-коды для CI/CD по задержкам, уровню ошибок и RPS |
| Двойные отчеты | JSON + независимый HTML-дашборд |
| Свои заголовки | Токены авторизации, tenant-заголовки, cookies |
| Поддержка Docker | Единый образ со встроенными утилитами |
| Кроссплатформенность | Linux, macOS, Windows — amd64 & arm64 |
| Ядро на Go | Быстрая установка, минимальные накладные расходы |

---

## Быстрый старт

```bash
# Установка
go install github.com/umutozen/stormprobe/cmd/stormprobe@latest

# Полный тест
stormprobe https://example.com

# Без обнаружения
stormprobe --no-discovery https://example.com

# Auth-тест
stormprobe --insecure \
  -H "Authorization: Bearer TOKEN" \
  https://api.example.com

# CI/CD-шлюз
stormprobe \\
  --alert-p99 500 \\
  --alert-error-rate 5 \\
  https://example.com

# Использовать JSON-файл конфигурации
stormprobe --config stormprobe.json
```

---

## Поток Тестирования

```text
Обнаружение → Нарастание → Плато → Пик → Восстановление → Вердикт
```

| Этап | Цель |
|---|---|
| Discovery | Сбор URL |
| Ramp-Up | Нагрузка |
| Sustained | Удержание |
| Spike | Скачок |
| Recovery | После скачка |
| Verdict | Финальный анализ |

---

## Установка

### Go Install

```bash
go install github.com/umutozen/stormprobe/cmd/stormprobe@latest
```

Рекомендуемый способ.

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

Позволяет заблокировать деплой при провале.

---

## Сценарии применения

**Перед релизом** — Проверка SLA.

**Защита + Нагрузка** — Скрытые URL.

**Аудит** — Точки отказа.

**Регресс** — Сравнения сборок.

---

## Flags

| Flag | Описание |
|---|---|
| `--duration` | Длительность (напр. 30s) |
| `--alert-p99` | Порог P99 (ms) |
| `--alert-error-rate` | Порог ошибок (%) |
| `--alert-rps` | Отсечка RPS |
| `-H` / `--header` | Заголовок HTTP |
| `--insecure` | Откл TLS проверки |
| `--no-discovery` | Только Root |
| `--endpoints` | Загрузить файл URL |
| `--format` | Формат отчета |
| `--concurrency-ramp` | Рамп Пользователей |
| `--concurrency-sustained` | Стабильных Пользователей |
| `--concurrency-spike` | Пиковых Пользователей |
| `--req-per-worker` | Запросы от worker-а |
| `--timeout` | Тайм-аут API |
| `--config` | Загрузить настройки из JSON-файла конфигурации |
| `--auto-profile` | Автоматически определить стек сервера и применить диагностический профиль (по умолчанию: активно) |
| `--profile` | Профиль вручную: iis, tomcat, php |
| `--no-fingerprint` | Пропустить определение стека сервера |
| `--show-profile` | Показать обнаруженный стек и профиль, затем выйти |

---

## Требования

| Инструмент | Примечание |
|---|---|
| Go 1.21+ | Go |
| Katana | Опц |
| Httpx | Опц |
| Docker | Docker |

---

## Предупреждение

Для этичных тестов.

Согласуйте тестирование.

Незаконно без разрешения.

---

## Создание

См. CONTRIBUTING.md.

## Безопасность

См. SECURITY.md.

## Лицензия

MIT License © Umut ÖZEN

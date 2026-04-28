# StormProbe

<p align="center">
  <img src="../../assets/banner.png" alt="StormProbe Banner" width="100%">
</p>

<h1 align="center">StormProbe</h1>

<p align="center">
  <strong>crawl · discover · stress · decide</strong><br>
  Moderne HTTP load testing met automatische eindpuntherkenning en oordelen over productiegereedheid.
</p>

<p align="center">
  <strong>Niet alleen bankmetingen.</strong><br>
  StormProbe test niet alleen. Het legt uit of uw systeem daadwerkelijk veiliger is geworden.
</p>

---

## Waarom StormProbe?

Traditionele load testers tonen cijfers.<br>StormProbe maakt beslissingen.

In plaats van handmatig eindpunten te selecteren...

### Klassiek tgo StormProbe

| Klassiek | StormProbe |
|---|---|
| Handmatig | Automatisch actieve url ontdekken |
| Ruwe data | Conclusies |
| Zelf besturen | Aanbevelingen |
| 1 actiepad | Groei + Continu + Piek + Herstel |
| Iets draait erg langzaam | Trage prestatie wegens... |

---

## Wat doet het?

HTTP belasting versterken en analyseren.

- Caching verificatie
- Auto scaling
- ...

---

## Resultaat Voorbeeld

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

## Kenmerken

| Eigenschap | Detail |
|---|---|
| Autodiscovery | Via scanner |
| Fased loads | 4 manieren |
| Conclusiemotor | Identificeert bottelnek |
| Meldingen voor CI/CD | Via error code |

---

## Snel van start

```bash
# Installeer
go install github.com/umutozen/stormprobe/cmd/stormprobe@latest

# Volledig
stormprobe https://example.com

# Root tests
stormprobe --no-discovery https://example.com

# Authenticeer
stormprobe --insecure \
  -H "Authorization: Bearer TOKEN" \
  https://api.example.com

# Samen in pijplijn
stormprobe \\
  --alert-p99 500 \\
  --alert-error-rate 5 \\
  https://example.com

# JSON-configuratiebestand gebruiken
stormprobe --config stormprobe.json
```

---

## Werkvloei

```text
Vooraf
```

| Fases | Doelen |
|---|---|
| Discovery | Zoek netwerken |
| Ramp-Up | Boosters aan |
| Sustained | Uithoudingsvermogen |
| Spike | Verkeersspijkers |
| Recovery | Hersteltijd |
| Verdict | Sluiter |

---

## Inrichten

### Go Install

```bash
go install github.com/umutozen/stormprobe/cmd/stormprobe@latest
```

Beste manier.

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

Voorkomt implementaties die de boel breken.

---

## Scenario's

**Testomgeving** — Verg.

**Hacks** — Achterdeuren

**Architectuur** — Stresstest

**CI Automatisatie** — Vergelijk via builds

---

## Flags

| Flag | Uitleg |
|---|---|
| `--duration` | Doorlooptijd per test |
| `--alert-p99` | Verlaat bij p99 in de knel |
| `--alert-error-rate` | Foutenpercentage |
| `--alert-rps` | Min aantal RPS |
| `-H` / `--header` | HTTP kopjes |
| `--insecure` | Sla TLS beveil. over |
| `--no-discovery` | Root adres test |
| `--endpoints` | Doelenbestand |
| `--format` | Visie via: |
| `--concurrency-ramp` | Maximum concur. |
| `--concurrency-sustained` | Blijvende concur. |
| `--concurrency-spike` | Aanval concur. |
| `--req-per-worker` | Threads taken |
| `--timeout` | Timeout reactie |
| `--config` | Instellingen laden uit JSON-configuratiebestand |
| `--auto-profile` | Automatisch serverstack detecteren en diagnostisch profiel toepassen (standaard: actief) |
| `--profile` | Handmatig profiel: iis, tomcat, php |
| `--no-fingerprint` | Serverstack-detectie overslaan |
| `--show-profile` | Gedetecteerde stack en profiel weergeven, dan afsluiten |

---

## Pakketvereisten

| Bibliotheek | Info |
|---|---|
| Go 1.21+ | Go nodig |
| Katana | Tools |
| Httpx | Web detectie |
| Docker | Container |

---

## Waarschuwing

Enkel legale testen.

Pols verantwoordelijken.

Illegaal anders.

---

## Samenwerking

Bekijk tekst file.

## Meld Kwetsbaarheden

Security.

## Licenties

MIT License © Umut ÖZEN

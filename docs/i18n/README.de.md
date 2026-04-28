# StormProbe

<p align="center">
  <img src="../../assets/banner.png" alt="StormProbe Banner" width="100%">
</p>

<h1 align="center">StormProbe</h1>

<p align="center">
  <strong>crawl · discover · stress · decide</strong><br>
  Modernes HTTP-Lasttesten mit autonomer Endpunkt-Erkennung und Produktionsbereitschaftsbewertung.
</p>

<p align="center">
  <strong>Nicht nur Benchmarking.</strong><br>
  StormProbe testet nicht nur. Es erklärt, ob Ihr System tatsächlich sicherer geworden ist.
</p>

---

## Warum StormProbe?

Herkömmliche Lasttest-Tools zeigen Zahlen.<br>StormProbe fällt Entscheidungen.

Anstatt Endpunkte manuell auszuwählen und rohe Tabellen zu interpretieren, erkennt StormProbe automatisch Endpunkte, führt phasenweise Stresstests durch, identifiziert Engpässe und liefert handlungsorientierte Empfehlungen.

### Herkömmliche Lasttester vs StormProbe

| Herkömmliche Werkzeuge | StormProbe |
|---|---|
| Manuelle Endpunktauswahl | Erkennt automatisch Live-Endpunkte |
| Rohes Benchmark-Ergebnis | Endgültige Bewertung + Engpassdiagnose |
| Manuelle Interpretation nötig | Konkrete Handlungsempfehlungen |
| Einfaches Stresstesten | Ramp-Up + Sustained + Spike + Recovery |
| "Etwas ist langsam" | "Hier ist, wo und warum es zusammenbricht" |

---

## Was macht es?

StormProbe ist ein in Go geschriebenes HTTP-Lasttest-CLI für die Produktion.

- Endpunkt-Erkennung
- Phasenweises Auslastungstesten
- Latenz-Perzentil-Analyse
- Durchsatz-Sättigungs-Erkennung
- Wiederherstellungsvalidierung
- CI/CD-Schwellenwertalarme
- JSON + HTML Berichte
- Produktionsbereitschaftsbewertungen

---

## Beispiel einer realen Ausgabe

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

## Funktionen

| Funktion | Details |
|---|---|
| Auto Discovery | Katana Crawl + Httpx Probe mit deduplizierten Live-Endpunkten |
| 4-Phasen-Lasttest | Ramp-Up → Sustained → Spike → Recovery |
| Verdict Engine | Sichere Nebenläufigkeit, Verschlechterungspunkt, Engpassklassifizierung |
| Umfangreiche Metriken | Avg / P50 / P95 / P99 Latenz + Req/s + Fehlerklassifizierung |
| Durchsatzanalyse | Erkennt Sättigung, bevor harte Ausfälle auftreten |
| Recovery-Validierung | Bestätigt, ob sich das System nach Verkehrsspitzen erholt |
| Alarmgrenzwerte | CI/CD-bereite Exit-Codes für Latenz, Fehlerrate und RPS |
| Duale Berichte | JSON + eigenständiges HTML-Dashboard |
| Benutzerdefinierte Header | Auth-Token, Tenant-Header, Cookies |
| Docker Ready | Einzelnes Image mit gebündelten Tools |
| Plattformübergreifend | Linux, macOS, Windows — amd64 & arm64 |
| Reiner Go-Kern | Schnelle Installation, minimaler Betriebsaufwand |

---

## Schnellstart

```bash
# Installieren
go install github.com/umutozen/stormprobe/cmd/stormprobe@latest

# Vollständiger Test mit automatischer Erkennung
stormprobe https://example.com

# Erkennung überspringen, nur Root-Pfad testen
stormprobe --no-discovery https://example.com

# API-Test mit Auth-Header
stormprobe --insecure \
  -H "Authorization: Bearer TOKEN" \
  https://api.example.com

# CI/CD-Integration
stormprobe \\
  --alert-p99 500 \\
  --alert-error-rate 5 \\
  https://example.com

# Mit JSON-Konfigurationsdatei starten
stormprobe --config stormprobe.json
```

---

## Testablauf

```text
Erkennung → Hochfahren → Dauerhaft → Spitzenbelastung → Wiederherstellung → Urteil
```

| Phase | Zweck |
|---|---|
| Discovery | Crawl & Validierung aktiver Endpunkte |
| Ramp-Up | Allmähliche Erhöhung der Nebenläufigkeit |
| Sustained | Stetige Leistung verifizieren |
| Spike | Plötzliche Traffic-Berechnung |
| Recovery | Gesundheitscheck nach der Spitze |
| Verdict | Sichere Concurrency- & Engpass-Analyse |

---

## Installation

### Go Install

```bash
go install github.com/umutozen/stormprobe/cmd/stormprobe@latest
```

Empfohlener Weg.

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

Dadurch fallen Deployments bei Leistungsregressionen automatisch durch.

---

## Anwendungsfälle

**Vor dem Produktivbetrieb** — Prüfen Sie, ob Nebenläufigkeitsänderungen Latenzziele brechen.

**Sicherheit + Leistung** — Verborgene Endpunkte entdecken und testen.

**Infrastrukturprüfungen** — Messen Sie Sättigungspunkte vor Spitzenverkehr.

**Regressionserkennung** — Leistung vor und nach einem Deployment vergleichen.

---

## Flags

| Flag | Beschreibung |
|---|---|
| `--duration` | Führe jede Phase für X aus (z.B. `30s`, `1m`) |
| `--alert-p99` | Verlassen mit Fehler (1), falls P99 Latenz Schwellwert übersteigt |
| `--alert-error-rate` | Fehler mit 1, wenn Fehlerquote Schwellwert überschreitet |
| `--alert-rps` | Fehler wenn Req/s sinkt |
| `-H` / `--header` | Manueller Header |
| `--insecure` | Ohne TLS-Check |
| `--no-discovery` | Root pfad testen ohne Erkennung |
| `--endpoints` | Endpunkt-Liste hochladen |
| `--format` | Ausgabe Format |
| `--concurrency-ramp` | Maximum Ramp-Up concurrency |
| `--concurrency-sustained` | Sustained concurrency |
| `--concurrency-spike` | Spike peak concurrency |
| `--req-per-worker` | Anfragen pro Worker |
| `--timeout` | Timeout |
| `--config` | Einstellungen aus JSON-Konfigurationsdatei laden |
| `--auto-profile` | Server-Stack automatisch erkennen und Diagnoseprofil anwenden (Standard: aktiv) |
| `--profile` | Manuelles Profil: iis, tomcat, php |
| `--no-fingerprint` | Server-Fingerprinting überspringen |
| `--show-profile` | Erkannten Stack und Profil anzeigen, dann beenden |

---

## Anforderungen

| Paket | Bemerkung |
|---|---|
| Go 1.21+ | Für Installation benötigt |
| Katana | Optional für Erkennung |
| Httpx | Optional für Erkennung |
| Docker | Zum Docker-Starten erforderlich |

---

## Haftungsausschluss

Lediglich zu Testzwecken freigegeben.

Holen Sie immer die Erlaubnis ein.

Eine ungenehmigte Nutzung ist rechtswidrig.

---

## Beitragen

Lesen Sie `CONTRIBUTING.md`.

## Sicherheit

Informationen zu Fehlern in `SECURITY.md`.

## Lizenz

MIT License © Umut ÖZEN

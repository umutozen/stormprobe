# StormProbe

<p align="center">
  <img src="../../assets/banner.png" alt="StormProbe Banner" width="100%">
</p>

<h1 align="center">StormProbe</h1>

<p align="center">
  <strong>crawl · discover · stress · decide</strong><br>
  Test di carico HTTP moderni con scoperta autonoma degli endpoint e verdetti di prontezza alla produzione.
</p>

<p align="center">
  <strong>Non solo benchmarking.</strong><br>
  StormProbe non si limita a testare. Spiega se il tuo sistema è diventato effettivamente più sicuro.
</p>

---

## Perché StormProbe?

Gli strumenti di test del carico tradizionali mostrano numeri.<br>StormProbe fornisce decisioni.

Invece di selezionare manualmente gli endpoint, StormProbe li scopre automaticamente, esegue test di stress a fasi e produce raccomandazioni attuabili.

### Strumenti Tradizionali vs StormProbe

| Strumenti Tradizionali | StormProbe |
|---|---|
| Selezione manuale degli endpoint | Scopre automaticamente gli endpoint attivi |
| Output benchmark grezzo | Verdetto finale + diagnosi dei colli di bottiglia |
| Richiede interpretazione manuale | Raccomandazioni attuabili |
| Schema di stress singolo | Ramp-Up + Sustained + Spike + Recovery |
| "Qualcosa è lento" | "Ecco dove e perché si rompe" |

---

## Cosa fa?

StormProbe è un CLI di test del carico HTTP di qualità produzione scritto in Go.

- Scoperta degli endpoint
- Test di concorrenza a fasi
- Analisi dei percentili di latenza
- Rilevamento della saturazione
- Validazione del recupero
- Avvisi soglia CI/CD
- Report JSON + HTML
- Verdetti di prontezza alla produzione

---

## Esempio di Output Reale

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

## Funzionalità

| Funzionalità | Dettagli |
|---|---|
| Scoperta Automatica | Katana crawl + Httpx probe con endpoint live deduplicati |
| Test di Carico a 4 Fasi | Ramp-Up → Sustained → Spike → Recovery |
| Motore di Verdetto | Concorrenza sicura, punto di degrado, classificazione colli di bottiglia |
| Metriche Ricche | Latenza Avg / P50 / P95 / P99 + req/s + classificazione errori |
| Analisi Throughput | Rileva la saturazione prima di guasti critici |
| Validazione di Recupero | Conferma se il sistema si riprende dopo i picchi di traffico |
| Soglie di Avviso | Codici di uscita pronti per CI/CD su latenza, tasso di errore e RPS |
| Report Doppi | JSON + dashboard HTML autonoma |
| Header Personalizzati | Token di auth, header tenant, cookie |
| Pronto per Docker | Immagine singola con strumenti raggruppati |
| Multipiattaforma | Linux, macOS, Windows — amd64 & arm64 |
| Core puro in Go | Installazione veloce, carico operativo minimo |

---

## Avvio Rapido

```bash
# Installazione
go install github.com/umutozen/stormprobe/cmd/stormprobe@latest

# Test completo con scoperta
stormprobe https://example.com

# Salta scoperta, test root
stormprobe --no-discovery https://example.com

# Test API con auth
stormprobe --insecure \
  -H "Authorization: Bearer TOKEN" \
  https://api.example.com

# Controllo CI/CD
stormprobe \\
  --alert-p99 500 \\
  --alert-error-rate 5 \\
  https://example.com

# Usa file di configurazione JSON
stormprobe --config stormprobe.json
```

---

## Flusso di Test

```text
Scoperta → Ramp-Up → Sustained → Spike → Recupero → Verdetto
```

| Fase | Scopo |
|---|---|
| Discovery | Identifica endpoint |
| Ramp-Up | Aumento graduale |
| Sustained | Verifica prestazioni reali |
| Spike | Picco di traffico |
| Recovery | Verifica stato dopo picco |
| Verdict | Analisi collo di bottiglia |

---

## Installazione

### Go Install

```bash
go install github.com/umutozen/stormprobe/cmd/stormprobe@latest
```

Metodo raccomandato.

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

Ciò consente alle regressioni di prestazioni di fallire automaticamente i deployment.

---

## Casi d'Uso

**Prima della Produzione** — Verifica limiti di latenza.

**Sicurezza + Prestazioni** — Scopri endpoint nascosti.

**Audit Infrastruttura** — Misura i punti di saturazione.

**Rilevamento Regressioni** — Confronta prestazioni CI/CD.

---

## Flags

| Flag | Descrizione |
|---|---|
| `--duration` | Durata per fase (es. `30s`, `1m`) |
| `--alert-p99` | Esci con 1 se P99 latenza supera limite (ms) |
| `--alert-error-rate` | Esci con 1 se errori superano limite (%) |
| `--alert-rps` | Esci con 1 se req/s scende |
| `-H` / `--header` | Header personalizzato |
| `--insecure` | Salta verifica TLS |
| `--no-discovery` | Salta scoperta |
| `--endpoints` | Carica endpoint file |
| `--format` | Formato: `json`, `html`, `both` |
| `--concurrency-ramp` | Picco Ramp-Up |
| `--concurrency-sustained` | Picco Sostenuto |
| `--concurrency-spike` | Picco Spike |
| `--req-per-worker` | Richieste/worker |
| `--timeout` | Timeout/richiesta |
| `--config` | Carica impostazioni da file di configurazione JSON |
| `--auto-profile` | Rileva automaticamente lo stack server e applica il profilo diagnostico (predefinito: attivo) |
| `--profile` | Profilo manuale: iis, tomcat, php |
| `--no-fingerprint` | Salta il rilevamento dello stack server |
| `--show-profile` | Mostra stack rilevato e profilo, poi esci |

---

## Requisiti

| Requisito | Note |
|---|---|
| Go 1.21+ | Per installazione |
| Katana | Opzionale per scoperta |
| Httpx | Opzionale |
| Docker | Opzionale per runtime |

---

## Limitazione

Solo test autorizzati.

In caso contrario illegale.

Declina responsabilità.

---

## Contribuire

Vedi `CONTRIBUTING.md`.

## Sicurezza

Vedi `SECURITY.md`.

## Licenza

MIT License © Umut ÖZEN

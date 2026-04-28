# StormProbe

<p align="center">
  <img src="../../assets/banner.png" alt="StormProbe Banner" width="100%">
</p>

<h1 align="center">StormProbe</h1>

<p align="center">
  <strong>crawl · discover · stress · decide</strong><br>
  Pruebas de carga HTTP modernas con descubrimiento de endpoints y diagnóstico de producción.
</p>

<p align="center">
  <strong>No se trata de un simple benchmark.</strong><br>
  StormProbe no solo testea. Explica si tu sistema realmente se volvió más seguro.
</p>

---

## ¿Por qué StormProbe?

Las herramientas clásicas muestran números.<br>StormProbe toma decisiones.

En lugar de elegir a ciegas y leer tablas, descubre URLs, ataca por etapas y produce recomendaciones.

### Básico frente a StormProbe

| Herramientas Básicas | StormProbe |
|---|---|
| Selección manual | Descubrimiento en vivo |
| Datos brutos | Diagnóstico y cuellos de botella |
| Interpretación humana | Recomendaciones claras |
| Ataque continuo | Ramp / Sus / Spike / Rec |
| "Va lento" | "Aquí se rompe y por esto" |

---

## ¿Qué es?

Es un CLI de Go pensado para producción.

- Crawler de Endpoints
- Escalada progresiva
- Percentiles
- Saturación
- Recuperación
- CI/CD
- JSON+HTML
- Predicción

---

## Ejemplo en TTY

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

## Características

| Característica | Detalles |
|---|---|
| Descubrimiento Automático | Katana crawl + Httpx probe con endpoints deduplicados |
| Prueba de 4 Fases | Ramp-Up → Sustained → Spike → Recovery |
| Motor de Veredicto | Concurrencia segura, punto de degradación, clasificación de cuello de botella |
| Métricas Enriquecidas | Latencia Avg / P50 / P95 / P99 + req/s + clasificación de errores |
| Análisis de Rendimiento | Detecta la saturación antes de fallas críticas |
| Validación de Recuperación | Confirma si el sistema se recupera después de picos de tráfico |
| Umbrales de Alerta | Códigos de salida para CI/CD basados en latencia, tasa de errores y RPS |
| Reportes Duales | JSON + dashboard HTML independiente |
| Cabeceras Personalizadas | Tokens de autenticación, headers de inquilino, cookies |
| Listo para Docker | Imagen única con herramientas incluidas |
| Multiplataforma | Linux, macOS, Windows — amd64 & arm64 |
| Núcleo puro en Go | Instalación rápida, mínima sobrecarga operativa |

---

## Inicio Rápido

```bash
# Instalar
go install github.com/umutozen/stormprobe/cmd/stormprobe@latest

# Prueba total con descubrimiento
stormprobe https://example.com

# Sólo raís
stormprobe --no-discovery https://example.com

# Auth
stormprobe --insecure \
  -H "Authorization: Bearer TOKEN" \
  https://api.example.com

# Test CI/CD
stormprobe \\
  --alert-p99 500 \\
  --alert-error-rate 5 \\
  https://example.com

# Usar archivo de configuración JSON
stormprobe --config stormprobe.json
```

---

## Fases de Carga

```text
Descubrir → Subir → Mantener → Pico → Recuperar → Fin
```

| Fase | Fin |
|---|---|
| Discovery | Encontrar rutas operativas |
| Ramp-Up | Subir uso |
| Sustained | Estable |
| Spike | Pico masivo |
| Recovery | Recuperación |
| Verdict | Análisis y Cuello |

---

## Instalación

### Go Install

```bash
go install github.com/umutozen/stormprobe/cmd/stormprobe@latest
```

Forma recomendada.

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

Configurado para frenar despliegues rotos.

---

## Casos de Uso

**Deploy** — Verificación final.

**Seguridad** — APIs ocultas.

**Auditorías** — Medir puntos ciegos.

**Regresión** — CI y Pull Requests.

---

## Flags

| Flag | Dec |
|---|---|
| `--duration` | Tiempo base |
| `--alert-p99` | Error latency |
| `--alert-error-rate` | Error rate |
| `--alert-rps` | Req/s min |
| `-H` / `--header` | Header auth |
| `--insecure` | Sin certs |
| `--no-discovery` | Root test |
| `--endpoints` | Archivo lista |
| `--format` | Salida |
| `--concurrency-ramp` | Ramp vu |
| `--concurrency-sustained` | Sus vu |
| `--concurrency-spike` | Spike vu |
| `--req-per-worker` | Req worker |
| `--timeout` | Net timeout |
| `--config` | Cargar configuración desde archivo JSON |
| `--auto-profile` | Detectar stack del servidor automáticamente y aplicar perfil (por defecto: activo) |
| `--profile` | Perfil manual: iis, tomcat, php |
| `--no-fingerprint` | Omitir detección de stack del servidor |
| `--show-profile` | Mostrar stack detectado y perfil, luego salir |

---

## Requiere

| Tool | Uso |
|---|---|
| Go 1.21+ | Go |
| Katana | Crawl |
| Httpx | Web |
| Docker | Docker |

---

## Aviso

Solo ético.

Con permiso.

Bajo la ley.

---

## Comunidad

Ver contribuciones.

## Bugs

En archivo security.

## Licencia

MIT License © Umut ÖZEN

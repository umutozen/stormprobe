[English](../../README.md) · [Türkçe](README.tr.md) · [Deutsch](README.de.md) · [Francais](README.fr.md) · [Espanol](README.es.md) · [Italiano](README.it.md) · [Portugues](README.pt-BR.md) · [Russkij](README.ru.md) · [Zhongwen Jian](README.zh-CN.md) · [Zhongwen Fan](README.zh-TW.md) · [Nihongo](README.ja.md) · [Hangugeo](README.ko.md) · [Al-Arabiyya](README.ar.md) · [Hindi](README.hi.md) · [Nederlands](README.nl.md) · [Polski](README.pl.md) · [Ukrainska](README.uk.md)

<p align="center">
  <img src="../../assets/banner.png" alt="StormProbe Banner" width="100%">
</p>

<h1 align="center">StormProbe</h1>

<p align="center">
  <strong>descubrir · cargar · analizar</strong><br>
  Herramienta de prueba de carga HTTP de nivel productivo con descubrimiento autónomo de endpoints. Sin dependencias externas.
</p>

<p align="center">
  <a href="https://github.com/umutozen/stormprobe/actions/workflows/ci.yml">
    <img src="https://github.com/umutozen/stormprobe/actions/workflows/ci.yml/badge.svg" alt="CI">
  </a>
  <a href="https://goreportcard.com/report/github.com/umutozen/stormprobe">
    <img src="https://goreportcard.com/badge/github.com/umutozen/stormprobe" alt="Go Report Card">
  </a>
  <a href="https://github.com/umutozen/stormprobe/releases">
    <img src="https://img.shields.io/github/v/release/umutozen/stormprobe" alt="Ultima version">
  </a>
  <a href="../../LICENSE">
    <img src="https://img.shields.io/badge/license-MIT-blue.svg" alt="Licencia MIT">
  </a>
</p>

---

## ¿Qué es StormProbe?

StormProbe es una herramienta de prueba de carga HTTP **sin dependencias** escrita en Go puro. Descubre automáticamente los endpoints de su aplicación con [Katana](https://github.com/projectdiscovery/katana) y [Httpx](https://github.com/projectdiscovery/httpx), luego ejecuta una prueba de carga estructurada en 4 fases — **Ramp-Up → Sustained → Spike → Recovery** — y produce informes JSON y HTML con métricas de latencia P50/P95/P99.

Diseñado para:
- **Pipelines CI/CD** — umbrales de alerta con códigos de salida
- **Pentesters** — descubrimiento automático de endpoints activos
- **Ingenieros DevOps** — benchmarks antes y después del despliegue
- **Equipos QA** — validación de SLA de rendimiento

---

## Características

| Característica | Detalles |
|---|---|
| **Descubrimiento automático** | Rastreo Katana + sonda Httpx, lista de endpoints activos deduplicada |
| **Prueba en 4 fases** | Ramp-Up → Sustained → Spike → Recovery |
| **Métricas ricas** | Latencia P50/P95/P99, req/s, clasificación de errores por fase |
| **Modo duración** | Fases basadas en tiempo (`--duration 30s`) en lugar de conteo de solicitudes |
| **Umbrales de alerta** | `exit 1` compatible con CI/CD en violación de P99, tasa de error o RPS |
| **Informes dobles** | JSON (legible por máquina) + panel HTML oscuro autónomo |
| **Encabezados personalizados** | Tokens Bearer, encabezados de tenant, cookies — propagados en todas partes |
| **Listo para Docker** | Imagen única con Katana + Httpx incluidos |
| **Multiplataforma** | Linux, macOS, Windows — amd64 & arm64 |
| **Sin dependencias** | Go 1.21+ stdlib puro, `go install` es suficiente |

---

## Inicio rápido

```bash
# Instalación (requiere Go 1.21+)
go install github.com/umutozen/stormprobe/cmd/stormprobe@latest

# Prueba básica — descubrimiento automático, 4 fases
stormprobe https://example.com

# Omitir descubrimiento, probar solo la ruta raíz
stormprobe --no-discovery https://example.com

# TLS + encabezado auth + solo informe JSON
stormprobe --insecure -H "Authorization: Bearer TOKEN" --format json https://api.example.com

# Prueba basada en duración: 30 segundos por fase
stormprobe --insecure --duration 30s --no-discovery https://example.com

# CI/CD: exit 1 si P99 > 500ms o errores > 5%
stormprobe --alert-p99 500 --alert-error-rate 5 https://example.com
```

---

## Instalación

### Opción 1 — `go install` *(recomendado)*

```bash
go install github.com/umutozen/stormprobe/cmd/stormprobe@latest
stormprobe --help
```

### Opción 2 — Binario precompilado *(sin Go)*

Descargar desde la [página de Releases](https://github.com/umutozen/stormprobe/releases).

```bash
# Linux / macOS
chmod +x stormprobe
sudo mv stormprobe /usr/local/bin/

# Windows PowerShell
Move-Item stormprobe.exe C:\Windows\System32\stormprobe.exe
```

### Opción 3 — Compilar desde fuente

```bash
git clone https://github.com/umutozen/stormprobe.git
cd stormprobe
go build -o stormprobe ./cmd/stormprobe      # Linux / macOS
go build -o stormprobe.exe ./cmd/stormprobe  # Windows
```

### Opción 4 — Docker

```bash
docker run --rm ghcr.io/umutozen/stormprobe:latest --insecure https://example.com

# Linux / macOS
docker run --rm -v $(pwd)/outputs:/app/outputs \
  ghcr.io/umutozen/stormprobe:latest --insecure --format both https://example.com

# Windows PowerShell
docker run --rm -v "${PWD}\outputs:/app/outputs" `
  ghcr.io/umutozen/stormprobe:latest --insecure --format both https://example.com
```

### Opcional — Herramientas de descubrimiento

```bash
bash scripts/install-tools.sh
```

> Sin Katana/Httpx usar `--no-discovery` o `--endpoints`.

---

## Referencia CLI

### Concurrencia y Carga

| Flag | Predeterminado | Descripción |
|---|---|---|
| `--concurrency-ramp` | `50` | Concurrencia máxima para la fase Ramp-Up |
| `--concurrency-sustained` | `75` | Concurrencia para la fase Sustained |
| `--concurrency-spike` | `250` | Concurrencia máxima para la fase Spike |
| `--req-per-worker` | `15` | Solicitudes por worker (ignorado si `--duration` está configurado) |
| `--duration` | `0` (desactivado) | Duración de fase basada en tiempo, ej. `30s`, `1m` |
| `--timeout` | `10s` | Tiempo de espera por solicitud |

### Descubrimiento

| Flag | Predeterminado | Descripción |
|---|---|---|
| `--no-discovery` | `false` | Omitir Katana+Httpx, probar solo la raíz |
| `--endpoints` | `""` | Cargar endpoints desde un archivo |
| `--katana-path` | `""` | Ruta personalizada del binario Katana |
| `--httpx-path` | `""` | Ruta personalizada del binario Httpx |

### HTTP y Seguridad

| Flag | Predeterminado | Descripción |
|---|---|---|
| `-H`, `--header` | — | Encabezado personalizado (repetible) |
| `--insecure` | `false` | Omitir verificación del certificado TLS |

### Salida

| Flag | Predeterminado | Descripción |
|---|---|---|
| `--format` | `both` | Formato de informe: `json`, `html`, `both` |
| `--output` | `./outputs` | Directorio de salida para informes |

### Alertas CI/CD

| Flag | Predeterminado | Descripción |
|---|---|---|
| `--alert-p99` | `0` (desactivado) | exit 1 si P99 (ms) supera el umbral |
| `--alert-error-rate` | `0` (desactivado) | exit 1 si la tasa de error (%) supera el umbral |
| `--alert-rps` | `0` (desactivado) | exit 1 si req/s cae por debajo del umbral |

---

## Fases de prueba

| Fase | Descripción |
|---|---|
| **0 — Discovery** | Rastreo Katana + sonda Httpx, lista deduplicada |
| **1 — Ramp-Up** | Aumento gradual de concurrencia: 5 → 15 → 30 → pico |
| **2 — Sustained** | 3 oleadas a la concurrencia objetivo |
| **3 — Spike** | Aumento repentino al pico de concurrencia, luego enfriamiento |
| **4 — Recovery** | Verificación de salud post-spike tras 10s de enfriamiento |

---

## Salida e Informes

| Archivo | Descripción |
|---|---|
| `stormprobe_report_<timestamp>.json` | Todas las métricas por fase, legible por máquina |
| `stormprobe_report_<timestamp>.html` | Panel oscuro autónomo con gráfico SVG de latencia |

---

## Aviso legal

Esta herramienta está destinada **únicamente a pruebas de seguridad autorizadas y evaluación de rendimiento**. Debe obtener permiso escrito explícito del propietario del sistema objetivo antes de realizar pruebas de carga. Los autores no asumen **ninguna responsabilidad** por mal uso o daños causados por esta herramienta.

---

## Licencia

[MIT](../../LICENSE) © [Umut ÖZEN](https://github.com/umutozen)

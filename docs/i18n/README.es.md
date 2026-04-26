[English](../../README.md) · [Türkçe](README.tr.md) · [Deutsch](README.de.md) · [Francais](README.fr.md) · [Espanol](README.es.md) · [Italiano](README.it.md) · [Portugues](README.pt-BR.md) · [Russkij](README.ru.md) · [Zhongwen Jian](README.zh-CN.md) · [Zhongwen Fan](README.zh-TW.md) · [Nihongo](README.ja.md) · [Hangugeo](README.ko.md) · [Al-Arabiyya](README.ar.md) · [Hindi](README.hi.md) · [Nederlands](README.nl.md) · [Polski](README.pl.md) · [Ukrainska](README.uk.md)

# StormProbe

<p align="center">
  <img src="../../assets/banner.png" alt="StormProbe" width="100%">
</p>

**descubrir. cargar. analizar.**

Probador de carga HTTP de grado de produccion con descubrimiento autonomo de endpoints (puntos finales). Cero dependencias.

[![CI](https://github.com/umutozen/stormprobe/actions/workflows/ci.yml/badge.svg)](https://github.com/umutozen/stormprobe/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/umutozen/stormprobe)](https://goreportcard.com/report/github.com/umutozen/stormprobe)

---

## Caracteristicas

- **Descubrimiento automatico** -- Rastreo Katana + sonda Httpx encuentra endpoints activos automaticamente
- **Pruebas multifase** -- Ramp-Up (aumento gradual), Sustained (sostenido), Spike (pico de carga), Recovery (recuperacion)
- **Metricas detalladas** -- P50 / P95 / P99 latency (latencia), req/s, clasificacion de errores
- **Informe doble** -- JSON (legible por maquina) + HTML (panel visual con tema oscuro)
- **Cero dependencias** -- Biblioteca estandar de Go pura, sin modulos de terceros
- **Listo para Docker** -- Comando unico con Katana + Httpx incluidos
- **Multiplataforma** -- Binarios para Linux, macOS, Windows via GoReleaser

## Inicio rapido

```bash
go run ./cmd https://example.com

go build -o stormprobe ./cmd
./stormprobe --insecure https://example.com
```

### Docker

```bash
# Prueba rápida, sin informe guardado
docker run --rm ghcr.io/umutozen/stormprobe:latest --insecure https://example.com

# Prueba completa con informes en ./outputs (Linux / macOS)
docker run --rm -v $(pwd)/outputs:/app/outputs ghcr.io/umutozen/stormprobe:latest --insecure --format both https://example.com

# Prueba completa con informes en ./outputs (Windows PowerShell)
docker run --rm -v "${PWD}\outputs:/app/outputs" ghcr.io/umutozen/stormprobe:latest --insecure --format both https://example.com

# Prueba de pico de alta carga (500 usuarios simultáneos)
docker run --rm ghcr.io/umutozen/stormprobe:latest --insecure --concurrency-spike 500 https://example.com

# Omitir descubrimiento, probar solo la ruta raíz
docker run --rm ghcr.io/umutozen/stormprobe:latest --insecure --no-discovery https://example.com

# Encabezados HTTP personalizados (autorización, tenant personalizado, etc.)
docker run --rm ghcr.io/umutozen/stormprobe:latest --insecure -H "Authorization: Bearer TOKEN" -H "X-Tenant: acme" https://api.example.com
```

## Uso

```
stormprobe [flags] <url-destino>

Flags:
  --concurrency-ramp      int       Concurrencia máxima para Ramp-Up (predeterminado 50)
  --concurrency-sustained int       Concurrencia para fase Sustained (predeterminado 75)
  --concurrency-spike     int       Concurrencia máxima para Spike (predeterminado 250)
  --req-per-worker        int       Solicitudes por worker por paso (predeterminado 15)
  --timeout               duration  Tiempo de espera por solicitud (predeterminado 10s)
  --endpoints             string    Archivo de lista de endpoints (una ruta por línea)
  --no-discovery                    Omitir katana+httpx, probar solo ruta raíz
  --output                string    Directorio de salida para informes (predeterminado ./outputs)
  --format                string    Formato del informe: json, html, both (predeterminado both)
  --katana-path           string    Ruta binaria de katana personalizada
  --httpx-path            string    Ruta binaria de httpx personalizada
  --insecure                        Omitir verificación de certificado TLS
  --header, -H            string    Encabezado HTTP personalizado (repetible): -H 'Authorization: Bearer TOKEN'
  --duration              duration  Duración por fase (ej. 30s, 1m). Anula req-per-worker
  --alert-p99             float     Exit 1 si la latencia P99 supera Xms
  --alert-error-rate      float     Exit 1 si la tasa de error supera X%%
  --alert-rps             float     Exit 1 si req/s cae por debajo de X
```

### Ejemplos

```bash
stormprobe --insecure --concurrency-spike 500 --req-per-worker 20 https://example.com

stormprobe --insecure --endpoints endpoints.txt https://example.com

stormprobe --insecure --format json --output ./results https://example.com

stormprobe --insecure -H "Authorization: Bearer TOKEN" -H "X-Tenant: acme" https://api.example.com

stormprobe --insecure --duration 30s --no-discovery https://example.com

stormprobe --insecure --duration 1m --concurrency-spike 300 https://example.com

stormprobe --insecure --alert-p99 500 --alert-error-rate 5 https://example.com
```

## Fases de prueba

| Fase | Descripcion |
|---|---|
| **0 -- Discovery** | Rastreo Katana + sonda Httpx, lista de endpoints deduplicada |
| **1 -- Ramp-Up** | Aumento gradual: 5, 15, 30, 50 virtual users (usuarios virtuales) |
| **2 -- Sustained** | 3 oleadas en la concurrency (concurrencia) objetivo, mide la degradacion |
| **3 -- Spike** | Rafaga repentina al pico, luego enfriamiento |
| **4 -- Recovery** | Verificacion de salud post-spike tras 10s de enfriamiento |

## Salida

| Archivo | Descripcion |
|---|---|
| `stormprobe_report_<timestamp>.json` | Resultados legibles por maquina con todas las metricas |
| `stormprobe_report_<timestamp>.html` | Panel visual, abrir en cualquier navegador |

## Estructura del proyecto

```
stormprobe/
    cmd/main.go                    Punto de entrada CLI
    internal/
        config/config.go           Tipos compartidos, generadores de fases, alert config
        alert/
            alert.go               Threshold checks, CI/CD exit code logic
        metrics/
            latency.go             Calculo de percentiles
            errors.go              Clasificacion de errores
        discovery/
            discovery.go           Interfaz publica Discover()
            katana.go              Integracion del rastreador Katana
            httpx.go               Integracion de la sonda Httpx
        runner/
            phase.go               Orquestacion de fases
            worker.go              Pool de goroutines, RNG por worker
        report/
            json.go                Escritor de informes JSON
            html.go                Informe HTML autonomo
    Dockerfile                     Multi-etapa con Katana + Httpx
    .goreleaser.yml                Configuracion de compilacion cruzada
    Makefile                       Objetivos de build, test, lint
    config.example.yml             Referencia de valores predeterminados de fases
```

## Requisitos

- Go 1.21+
- [Katana](https://github.com/projectdiscovery/katana) (ProjectDiscovery) -- opcional, para descubrimiento automatico
- [Httpx](https://github.com/projectdiscovery/httpx) (ProjectDiscovery) -- opcional, para validacion de endpoints

```bash
bash scripts/install-tools.sh
```

## Descargo de responsabilidad

Esta herramienta esta destinada unicamente a **pruebas de seguridad autorizadas y evaluacion de rendimiento**. Debe obtener un permiso escrito explicito del propietario del sistema objetivo antes de ejecutar cualquier prueba de carga. El uso no autorizado de esta herramienta contra sistemas que no posee o para los que no tiene permiso de prueba puede violar leyes locales, nacionales o internacionales. Los autores no asumen ninguna responsabilidad por el mal uso o los danos causados por esta herramienta.

## Licencia

MIT

## Autor

**Umut ÖZEN** — [@umutozen](https://github.com/umutozen)

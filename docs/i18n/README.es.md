[English](../../README.md) · [Turkce](README.tr.md) · [Deutsch](README.de.md) · [Francais](README.fr.md) · [Espanol](README.es.md) · [Italiano](README.it.md) · [Portugues](README.pt-BR.md) · [Russkij](README.ru.md) · [Zhongwen Jian](README.zh-CN.md) · [Zhongwen Fan](README.zh-TW.md) · [Nihongo](README.ja.md) · [Hangugeo](README.ko.md) · [Al-Arabiyya](README.ar.md) · [Hindi](README.hi.md) · [Nederlands](README.nl.md) · [Polski](README.pl.md) · [Ukrainska](README.uk.md)

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
docker run --rm -v $(pwd)/outputs:/app/outputs \
  ghcr.io/umutozen/stormprobe --insecure https://example.com
```

## Uso

```
stormprobe [flags] <target-url>

Flags:
  --concurrency-ramp      int       Peak concurrency for ramp-up (default 50)
  --concurrency-sustained int       Concurrency for sustained phase (default 75)
  --concurrency-spike     int       Peak concurrency for spike (default 250)
  --req-per-worker        int       Requests per worker per step (default 15)
  --timeout               duration  Per-request timeout (default 10s)
  --endpoints             string    Endpoints file (one path per line, skips discovery)
  --no-discovery                    Skip katana+httpx, test root path only
  --output                string    Output directory for reports (default ./outputs)
  --format                string    Report format: json, html, both (default both)
  --katana-path           string    Custom katana binary path
  --httpx-path            string    Custom httpx binary path
  --insecure                        Skip TLS certificate verification
```

### Ejemplos

```bash
stormprobe --insecure --concurrency-spike 500 --req-per-worker 20 https://example.com

stormprobe --insecure --endpoints endpoints.txt https://example.com

stormprobe --insecure --format json --output ./results https://example.com
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
        config/config.go           Tipos compartidos, generadores de fases
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

**Umut ÖZEN** -- [@umutozen](https://github.com/umutozen)

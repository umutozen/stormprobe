[English](../../README.md) · [Türkçe](README.tr.md) · [Deutsch](README.de.md) · [Francais](README.fr.md) · [Espanol](README.es.md) · [Italiano](README.it.md) · [Portugues](README.pt-BR.md) · [Russkij](README.ru.md) · [Zhongwen Jian](README.zh-CN.md) · [Zhongwen Fan](README.zh-TW.md) · [Nihongo](README.ja.md) · [Hangugeo](README.ko.md) · [Al-Arabiyya](README.ar.md) · [Hindi](README.hi.md) · [Nederlands](README.nl.md) · [Polski](README.pl.md) · [Ukrainska](README.uk.md)

<p align="center">
  <img src="../../assets/banner.png" alt="StormProbe Banner" width="100%">
</p>

<h1 align="center">StormProbe</h1>

<p align="center">
  <strong>descobrir · carregar · analisar</strong><br>
  Ferramenta de teste de carga HTTP de nivel de producao com descoberta autonoma de endpoints. Zero dependencias externas.
</p>

<p align="center">
  <a href="https://github.com/umutozen/stormprobe/actions/workflows/ci.yml">
    <img src="https://github.com/umutozen/stormprobe/actions/workflows/ci.yml/badge.svg" alt="CI">
  </a>
  <a href="https://goreportcard.com/report/github.com/umutozen/stormprobe">
    <img src="https://goreportcard.com/badge/github.com/umutozen/stormprobe" alt="Go Report Card">
  </a>
  <a href="https://github.com/umutozen/stormprobe/releases">
    <img src="https://img.shields.io/github/v/release/umutozen/stormprobe" alt="Ultima versao">
  </a>
  <a href="../../LICENSE">
    <img src="https://img.shields.io/badge/license-MIT-blue.svg" alt="Licenca MIT">
  </a>
</p>

---

## O que e o StormProbe?

StormProbe e uma ferramenta de teste de carga HTTP **sem dependencias** escrita em Go puro. Descobre automaticamente os endpoints da sua aplicacao com [Katana](https://github.com/projectdiscovery/katana) e [Httpx](https://github.com/projectdiscovery/httpx), depois executa um teste de carga estruturado em 4 fases — **Ramp-Up → Sustained → Spike → Recovery** — e produz relatorios JSON e HTML com metricas de latencia P50/P95/P99.

Desenvolvido para:
- **Pipelines CI/CD** — limiares de alerta com codigos de saida
- **Pentesters** — descoberta automatica de endpoints ativos
- **Engenheiros DevOps** — benchmarks antes e apos implantacao
- **Equipes QA** — validacao de SLAs de desempenho

---

## Recursos

| Recurso | Detalhes |
|---|---|
| **Descoberta automatica** | Rastreamento Katana + sonda Httpx, lista deduplicada |
| **Teste em 4 fases** | Ramp-Up → Sustained → Spike → Recovery |
| **Metricas ricas** | Latencia P50/P95/P99, req/s, classificacao de erros por fase |
| **Modo duracao** | Fases baseadas em tempo (`--duration 30s`) em vez de contagem |
| **Limiares de alerta** | `exit 1` compativel CI/CD em violacao de P99, taxa de erro ou RPS |
| **Relatorios duplos** | JSON (legivel por maquina) + painel HTML escuro autonomo |
| **Cabecalhos personalizados** | Tokens Bearer, cabecalhos tenant, cookies — propagados em todos os lugares |
| **Pronto para Docker** | Imagem unica com Katana + Httpx incluidos |
| **Multiplataforma** | Linux, macOS, Windows — amd64 & arm64 |
| **Zero dependencias** | Go 1.21+ stdlib puro, `go install` e suficiente |

---

## Inicio rapido

```bash
# Instalacao (requer Go 1.21+)
go install github.com/umutozen/stormprobe/cmd/stormprobe@latest

# Teste basico — descoberta automatica, 4 fases
stormprobe https://example.com

# Pular descoberta, testar apenas o caminho raiz
stormprobe --no-discovery https://example.com

# TLS + cabecalho auth + apenas relatorio JSON
stormprobe --insecure -H "Authorization: Bearer TOKEN" --format json https://api.example.com

# Teste baseado em duracao: 30 segundos por fase
stormprobe --insecure --duration 30s --no-discovery https://example.com

# CI/CD: exit 1 se P99 > 500ms ou erros > 5%
stormprobe --alert-p99 500 --alert-error-rate 5 https://example.com
```

---

## Instalacao

### Opcao 1 — `go install` *(recomendado)*

```bash
go install github.com/umutozen/stormprobe/cmd/stormprobe@latest
stormprobe --help
```

### Opcao 2 — Binario pre-compilado *(sem Go)*

Baixar da [pagina de Releases](https://github.com/umutozen/stormprobe/releases).

```bash
# Linux / macOS
chmod +x stormprobe
sudo mv stormprobe /usr/local/bin/

# Windows PowerShell
Move-Item stormprobe.exe C:\Windows\System32\stormprobe.exe
```

### Opcao 3 — Compilar a partir do codigo-fonte

```bash
git clone https://github.com/umutozen/stormprobe.git
cd stormprobe
go build -o stormprobe ./cmd/stormprobe      # Linux / macOS
go build -o stormprobe.exe ./cmd/stormprobe  # Windows
```

### Opcao 4 — Docker

```bash
docker run --rm ghcr.io/umutozen/stormprobe:latest --insecure https://example.com

# Linux / macOS
docker run --rm -v $(pwd)/outputs:/app/outputs \
  ghcr.io/umutozen/stormprobe:latest --insecure --format both https://example.com

# Windows PowerShell
docker run --rm -v "${PWD}\outputs:/app/outputs" `
  ghcr.io/umutozen/stormprobe:latest --insecure --format both https://example.com
```

### Opcional — Ferramentas de descoberta

```bash
bash scripts/install-tools.sh
```

> Sem Katana/Httpx usar `--no-discovery` ou `--endpoints`.

---

## Referencia CLI

### Concorrencia e Carga

| Flag | Padrao | Descricao |
|---|---|---|
| `--concurrency-ramp` | `50` | Concorrencia maxima para a fase Ramp-Up |
| `--concurrency-sustained` | `75` | Concorrencia para a fase Sustained |
| `--concurrency-spike` | `250` | Concorrencia maxima para a fase Spike |
| `--req-per-worker` | `15` | Requisicoes por worker (ignorado se `--duration` estiver definido) |
| `--duration` | `0` (desativado) | Duracao de fase baseada em tempo, ex. `30s`, `1m` |
| `--timeout` | `10s` | Tempo limite por requisicao |

### Descoberta

| Flag | Padrao | Descricao |
|---|---|---|
| `--no-discovery` | `false` | Pular Katana+Httpx, testar apenas a raiz |
| `--endpoints` | `""` | Carregar endpoints de arquivo |
| `--katana-path` | `""` | Caminho personalizado do binario Katana |
| `--httpx-path` | `""` | Caminho personalizado do binario Httpx |

### HTTP e Seguranca

| Flag | Padrao | Descricao |
|---|---|---|
| `-H`, `--header` | — | Cabecalho personalizado (repetivel) |
| `--insecure` | `false` | Pular verificacao do certificado TLS |

### Saida

| Flag | Padrao | Descricao |
|---|---|---|
| `--format` | `both` | Formato do relatorio: `json`, `html`, `both` |
| `--output` | `./outputs` | Diretorio de saida para relatorios |

### Alertas CI/CD

| Flag | Padrao | Descricao |
|---|---|---|
| `--alert-p99` | `0` (desativado) | exit 1 se P99 (ms) exceder o limiar |
| `--alert-error-rate` | `0` (desativado) | exit 1 se a taxa de erro (%) exceder o limiar |
| `--alert-rps` | `0` (desativado) | exit 1 se req/s cair abaixo do limiar |

---

## Fases de teste

| Fase | Descricao |
|---|---|
| **0 — Discovery** | Rastreamento Katana + sonda Httpx, lista deduplicada |
| **1 — Ramp-Up** | Aumento gradual de concorrencia: 5 → 15 → 30 → pico |
| **2 — Sustained** | 3 ondas na concorrencia-alvo |
| **3 — Spike** | Aumento repentino ao pico, depois resfriamento |
| **4 — Recovery** | Verificacao de saude pos-spike apos 10s de resfriamento |

---

## Saida e Relatorios

| Arquivo | Descricao |
|---|---|
| `stormprobe_report_<timestamp>.json` | Todas as metricas por fase, legivel por maquina |
| `stormprobe_report_<timestamp>.html` | Painel escuro autonomo com grafico SVG de latencia |

---

## Aviso legal

Esta ferramenta destina-se **exclusivamente a testes de seguranca autorizados e avaliacao de desempenho**. Voce deve obter permissao por escrito explicita do proprietario do sistema de destino antes de executar testes de carga. Os autores nao assumem **nenhuma responsabilidade** por uso indevido ou danos causados por esta ferramenta.

---

## Licenca

[MIT](../../LICENSE) © [Umut ÖZEN](https://github.com/umutozen)

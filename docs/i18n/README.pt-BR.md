[English](../../README.md) · [Türkçe](README.tr.md) · [Deutsch](README.de.md) · [Francais](README.fr.md) · [Espanol](README.es.md) · [Italiano](README.it.md) · [Portugues](README.pt-BR.md) · [Russkij](README.ru.md) · [Zhongwen Jian](README.zh-CN.md) · [Zhongwen Fan](README.zh-TW.md) · [Nihongo](README.ja.md) · [Hangugeo](README.ko.md) · [Al-Arabiyya](README.ar.md) · [Hindi](README.hi.md) · [Nederlands](README.nl.md) · [Polski](README.pl.md) · [Ukrainska](README.uk.md)

# StormProbe

<p align="center">
  <img src="../../assets/banner.png" alt="StormProbe" width="100%">
</p>

**descobrir. carregar. analisar.**

Testador de carga HTTP de nivel de producao com descoberta autonoma de endpoints (pontos de extremidade). Zero dependencias.

[![CI](https://github.com/umutozen/stormprobe/actions/workflows/ci.yml/badge.svg)](https://github.com/umutozen/stormprobe/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/umutozen/stormprobe)](https://goreportcard.com/report/github.com/umutozen/stormprobe)

---

## Recursos

- **Descoberta automatica** -- Varredura Katana + sonda Httpx encontra endpoints ativos automaticamente
- **Testes multifase** -- Ramp-Up (aumento gradual), Sustained (sustentado), Spike (pico de carga), Recovery (recuperacao)
- **Metricas detalhadas** -- P50 / P95 / P99 latency (latencia), req/s, classificacao de erros
- **Relatorio duplo** -- JSON (legivel por maquina) + HTML (painel visual tema escuro)
- **Zero dependencias** -- Biblioteca padrao Go pura, sem modulos de terceiros
- **Pronto para Docker** -- Comando unico com Katana + Httpx inclusos
- **Multiplataforma** -- Binarios Linux, macOS, Windows via GoReleaser

## Inicio rapido

```bash
go run ./cmd https://example.com

go build -o stormprobe ./cmd/stormprobe
./stormprobe --insecure https://example.com
```

## Instalação

### Baixar binário *(sem Go)*
Baixe a versão mais recente para sua plataforma na [página de Releases](https://github.com/umutozen/stormprobe/releases), extraia o arquivo e mova o binário para um diretório no seu PATH:

```bash
# Linux / macOS
chmod +x stormprobe
sudo mv stormprobe /usr/local/bin/
stormprobe --no-discovery https://example.com

# Windows (PowerShell)
Move-Item stormprobe.exe C:\Windows\System32\stormprobe.exe
stormprobe --no-discovery https://example.com
```

### go install *(requer Go 1.21+, método mais fácil)*
```bash
go install github.com/umutozen/stormprobe/cmd/stormprobe@latest
stormprobe --no-discovery https://example.com
```

### Compilar a partir do código-fonte *(requer Go 1.21+)*
```bash
git clone https://github.com/umutozen/stormprobe.git
cd stormprobe
go build -o stormprobe ./cmd/stormprobe   # Linux / macOS
go build -o stormprobe.exe ./cmd/stormprobe   # Windows
```

### Docker

```bash
# Teste rápido, sem relatório salvo
docker run --rm ghcr.io/umutozen/stormprobe:latest --insecure https://example.com

# Teste completo com relatórios em ./outputs (Linux / macOS)
docker run --rm -v $(pwd)/outputs:/app/outputs ghcr.io/umutozen/stormprobe:latest --insecure --format both https://example.com

# Teste completo com relatórios em ./outputs (Windows PowerShell)
docker run --rm -v "${PWD}\outputs:/app/outputs" ghcr.io/umutozen/stormprobe:latest --insecure --format both https://example.com

# Teste de pico de alta carga (500 usuários simultâneos)
docker run --rm ghcr.io/umutozen/stormprobe:latest --insecure --concurrency-spike 500 https://example.com

# Pular descoberta, testar apenas o caminho raiz
docker run --rm ghcr.io/umutozen/stormprobe:latest --insecure --no-discovery https://example.com

# Cabeçalhos HTTP personalizados (autorização, tenant personalizado, etc.)
docker run --rm ghcr.io/umutozen/stormprobe:latest --insecure -H "Authorization: Bearer TOKEN" -H "X-Tenant: acme" https://api.example.com
```

## Uso

```
stormprobe [flags] <url-alvo>

Flags:
  --concurrency-ramp      int       Concorrência máxima para Ramp-Up (padrão 50)
  --concurrency-sustained int       Concorrência para fase Sustained (padrão 75)
  --concurrency-spike     int       Concorrência máxima para Spike (padrão 250)
  --req-per-worker        int       Requisições por worker por etapa (padrão 15)
  --timeout               duration  Timeout por requisição (padrão 10s)
  --endpoints             string    Arquivo de lista de endpoints (um caminho por linha)
  --no-discovery                    Ignorar katana+httpx, testar apenas caminho raiz
  --output                string    Diretório de saída para relatórios (padrão ./outputs)
  --format                string    Formato do relatório: json, html, both (padrão both)
  --katana-path           string    Caminho binário katana personalizado
  --httpx-path            string    Caminho binário httpx personalizado
  --insecure                        Ignorar verificação de certificado TLS
  --header, -H            string    Cabeçalho HTTP personalizado (repetível): -H 'Authorization: Bearer TOKEN'
  --duration              duration  Duração por fase (ex. 30s, 1m). Substitui req-per-worker
  --alert-p99             float     Exit 1 se latência P99 exceder Xms
  --alert-error-rate      float     Exit 1 se taxa de erro exceder X%%
  --alert-rps             float     Exit 1 se req/s cair abaixo de X
```

### Exemplos

```bash
stormprobe --insecure --concurrency-spike 500 --req-per-worker 20 https://example.com

stormprobe --insecure --endpoints endpoints.txt https://example.com

stormprobe --insecure --format json --output ./results https://example.com

stormprobe --insecure -H "Authorization: Bearer TOKEN" -H "X-Tenant: acme" https://api.example.com

stormprobe --insecure --duration 30s --no-discovery https://example.com

stormprobe --insecure --duration 1m --concurrency-spike 300 https://example.com

stormprobe --insecure --alert-p99 500 --alert-error-rate 5 https://example.com
```

## Fases de teste

| Fase | Descricao |
|---|---|
| **0 -- Discovery** | Varredura Katana + sonda Httpx, lista de endpoints deduplicada |
| **1 -- Ramp-Up** | Aumento gradual: 5, 15, 30, 50 virtual users (usuarios virtuais) |
| **2 -- Sustained** | 3 ondas na concurrency (concorrencia) alvo, mede a degradacao |
| **3 -- Spike** | Rajada repentina ao pico, seguida de resfriamento |
| **4 -- Recovery** | Verificacao de saude pos-spike apos 10s de resfriamento |

## Saida

| Arquivo | Descricao |
|---|---|
| `stormprobe_report_<timestamp>.json` | Resultados legiveis por maquina com todas as metricas |
| `stormprobe_report_<timestamp>.html` | Painel visual, abrir em qualquer navegador |

## Estrutura do projeto

```
stormprobe/
    cmd/main.go                    Ponto de entrada CLI
    internal/
        config/config.go           Tipos compartilhados, geradores de fases, alert config
        alert/
            alert.go               Threshold checks, CI/CD exit code logic
        metrics/
            latency.go             Calculo de percentis
            errors.go              Classificacao de erros
        discovery/
            discovery.go           Interface publica Discover()
            katana.go              Integracao do crawler Katana
            httpx.go               Integracao da sonda Httpx
        runner/
            phase.go               Orquestracao de fases
            worker.go              Pool de goroutines, RNG por worker
        report/
            json.go                Escritor de relatorio JSON
            html.go                Relatorio HTML autonomo
    Dockerfile                     Multi-estagio com Katana + Httpx
    .goreleaser.yml                Configuracao de compilacao cruzada
    Makefile                       Alvos de build, test, lint
    config.example.yml             Referencia de valores padrao das fases
```

## Requisitos

- Go 1.21+
- [Katana](https://github.com/projectdiscovery/katana) (ProjectDiscovery) -- opcional, para descoberta automatica
- [Httpx](https://github.com/projectdiscovery/httpx) (ProjectDiscovery) -- opcional, para validacao de endpoints

```bash
bash scripts/install-tools.sh
```

## Aviso legal

Esta ferramenta e destinada exclusivamente a **testes de seguranca autorizados e avaliacao de desempenho**. Voce deve obter permissao por escrito explicita do proprietario do sistema alvo antes de executar qualquer teste de carga. O uso nao autorizado desta ferramenta contra sistemas que voce nao possui ou para os quais nao tem permissao de teste pode violar leis locais, nacionais ou internacionais. Os autores nao assumem nenhuma responsabilidade por uso indevido ou danos causados por esta ferramenta.

## Licenca

MIT

## Autor

**Umut ÖZEN** — [@umutozen](https://github.com/umutozen)

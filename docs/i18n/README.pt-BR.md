# StormProbe

<p align="center">
  <img src="../../assets/banner.png" alt="StormProbe Banner" width="100%">
</p>

<h1 align="center">StormProbe</h1>

<p align="center">
  <strong>crawl · discover · stress · decide</strong><br>
  Testes de carga HTTP modernos com descoberta autônoma de endpoints e veredictos.
</p>

<p align="center">
  <strong>Não apenas benchmarking.</strong><br>
  StormProbe não apenas testa. Ele explica se seu sistema realmente se tornou mais seguro.
</p>

---

## Por que StormProbe?

Ferramentas tradicionais mostram números.<br>StormProbe fornece decisões.

Em vez de selecionar endpoints manualmente, StormProbe os descobre automaticamente, executa testes de stress em fases e produz recomendações acionáveis.

### Ferramentas Tradicionais vs StormProbe

| Ferramentas Tradicionais | StormProbe |
|---|---|
| Seleção manual de endpoints | Descobre endpoints ativos automaticamente |
| Saída de benchmark bruta | Veredicto final + diagnóstico de gargalos |
| Requer interpretação manual | Recomendações acionáveis |
| Padrão de stress único | Ramp-Up + Sustained + Spike + Recovery |
| "Algo está lento" | "Aqui está onde e por que falha" |

---

## O que faz?

StormProbe é um CLI de teste de carga HTTP de qualidade de produção escrito em Go.

- Descoberta de endpoints
- Teste de concorrência em fases
- Análise de percentis
- Detecção de saturação
- Validação de recuperação
- Alertas CI/CD
- Relatórios JSON + HTML
- Veredictos de prontidão

---

## Exemplo de Saída Real

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

## Recursos

| Recursos | Detalhes |
|---|---|
| Descoberta Automática | Katana crawl + Httpx probe com endpoints deduplicados |
| Teste em 4 Fases | Ramp-Up → Sustained → Spike → Recovery |
| Motor de Veredicto | Concorrência segura, ponto de degradação, classificação de gargalos |
| Métricas Ricas | Latência Avg / P50 / P95 / P99 + req/s + classificação de erros |
| Análise de Throughput | Detecta saturação antes de falhas graves |
| Validação de Recuperação | Confirma se o sistema se recupera após picos de tráfego |
| Limites de Alerta | Códigos de saída para CI/CD baseados em latência, taxa de erro e RPS |
| Relatórios Duplos | JSON + painel HTML autônomo |
| Cabeçalhos Customizados | Tokens de auth, headers de tenant, cookies |
| Pronto para Docker | Imagem única com ferramentas integradas |
| Multiplataforma | Linux, macOS, Windows — amd64 & arm64 |
| Núcleo em Go | Instalação rápida, sobrecarga operacional mínima |

---

## Início Rápido

```bash
# Instalar
go install github.com/umutozen/stormprobe/cmd/stormprobe@latest

# Teste autom.
stormprobe https://example.com

# Teste raíz
stormprobe --no-discovery https://example.com

# Auth
stormprobe --insecure \
  -H "Authorization: Bearer TOKEN" \
  https://api.example.com

# Esteira CI/CD
stormprobe \\
  --alert-p99 500 \\
  --alert-error-rate 5 \\
  https://example.com

# Usar arquivo de configuração JSON
stormprobe --config stormprobe.json
```

---

## Fluxo de Teste

```text
Descoberta → Ramp-Up → Sustained → Spike → Recovery → Veredicto
```

| Fase | Objetivo |
|---|---|
| Discovery | Rastreio e validação |
| Ramp-Up | Aumento gradual |
| Sustained | Validar estabilidade |
| Spike | Pico de tráfego |
| Recovery | Verificar pós-pico |
| Verdict | Análise segura de concorrência |

---

## Instalação

### Go Install

```bash
go install github.com/umutozen/stormprobe/cmd/stormprobe@latest
```

Método recomendado.

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

Isso permite que regressões de desempenho falhem implantações automaticamente.

---

## Casos de Uso

**Antes do Lançamento** — Valide se as mudanças quebram metas.

**Segurança + Desempenho** — Descubra URLs ocultas.

**Auditorias de Infraestrutura** — Meça saturação.

**Detecção de Regressão** — Compare na CI.

---

## Flags

| Flag | Descrição |
|---|---|
| `--duration` | Duração por fase |
| `--alert-p99` | Erro 1 para P99 alto |
| `--alert-error-rate` | Erro 1 para Taxa de Erro |
| `--alert-rps` | Erro 1 para RPS baixo |
| `-H` / `--header` | Cabeçalho HTTP |
| `--insecure` | Pular TLS |
| `--no-discovery` | Pular descoberta |
| `--endpoints` | Arquivo endpoint |
| `--format` | Formato do Relatório |
| `--concurrency-ramp` | Ramp-Up VU |
| `--concurrency-sustained` | Sustained VU |
| `--concurrency-spike` | Spike VU |
| `--req-per-worker` | Reqs por worker |
| `--timeout` | Timeout de req |
| `--config` | Carregar configurações de um arquivo JSON |
| `--auto-profile` | Detectar automaticamente o stack do servidor e aplicar perfil (padrão: ativo) |
| `--profile` | Perfil manual: iis, tomcat, php |
| `--no-fingerprint` | Ignorar detecção do stack do servidor |
| `--show-profile` | Mostrar stack detectado e perfil, depois sair |

---

## Requisitos

| Requ. | Notas |
|---|---|
| Go 1.21+ | Go |
| Katana | Opcional |
| Httpx | Opcional |
| Docker | Docker opcional |

---

## Aviso

Apenas testes autorizados.

Peça permissão.

Uso não autorizado é ilegal.

---

## Contribuindo

Veja o guia de contribuição.

## Segurança

Veja arquivo de segurança.

## Licença

MIT License © Umut ÖZEN

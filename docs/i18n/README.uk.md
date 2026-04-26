[English](../../README.md) · [Türkçe](README.tr.md) · [Deutsch](README.de.md) · [Francais](README.fr.md) · [Espanol](README.es.md) · [Italiano](README.it.md) · [Portugues](README.pt-BR.md) · [Russkij](README.ru.md) · [Zhongwen Jian](README.zh-CN.md) · [Zhongwen Fan](README.zh-TW.md) · [Nihongo](README.ja.md) · [Hangugeo](README.ko.md) · [Al-Arabiyya](README.ar.md) · [Hindi](README.hi.md) · [Nederlands](README.nl.md) · [Polski](README.pl.md) · [Ukrainska](README.uk.md)

<p align="center">
  <img src="../../assets/banner.png" alt="StormProbe Banner" width="100%">
</p>

<h1 align="center">StormProbe</h1>

<p align="center">
  <strong>vyyavyty · navantazhyty · proanalizuvaty</strong><br>
  HTTP-instrument navantazhuvalnoho testuvannya vyrobnychoho rivnya z avtonomnym vyyavlennyam kintsevykh tochok. Nul'ovi zovnishni zalezhnosti.
</p>

<p align="center">
  <a href="https://github.com/umutozen/stormprobe/actions/workflows/ci.yml">
    <img src="https://github.com/umutozen/stormprobe/actions/workflows/ci.yml/badge.svg" alt="CI">
  </a>
  <a href="https://goreportcard.com/report/github.com/umutozen/stormprobe">
    <img src="https://goreportcard.com/badge/github.com/umutozen/stormprobe" alt="Go Report Card">
  </a>
  <a href="https://github.com/umutozen/stormprobe/releases">
    <img src="https://img.shields.io/github/v/release/umutozen/stormprobe" alt="Ostannya versiya">
  </a>
  <a href="../../LICENSE">
    <img src="https://img.shields.io/badge/license-MIT-blue.svg" alt="Litsenziya MIT">
  </a>
</p>

---

## Shcho take StormProbe?

StormProbe — tse **instrument navantazhuvalnoho HTTP-testuvannya bez zalezhnostej**, napysanyy na chystomu Go. Vin avtomatychno vyyavlyaye kintsevi tochky vashoho zastosunку za dopomohoyu [Katana](https://github.com/projectdiscovery/katana) i [Httpx](https://github.com/projectdiscovery/httpx), potim vykonuye strukturovanyy 4-fazovyy navantazhuvalnyy test — **Ramp-Up → Sustained → Spike → Recovery** — i stvoryuye zvity JSON i HTML z metrykamy zatrymky P50/P95/P99.

Rozroblено dlya:
- **Konveyeriv CI/CD** — porogy spovishchennya z kodamy vyhodu
- **Pentesteriv** — avtomatychne vyyavlennya aktyvnykh kintsevykh tochok
- **DevOps-inzheneriv** — benchmarky do i pislya rozghortannya
- **Komand QA** — validatsiya SLA produktyvnosti

---

## Mozhlyvosti

| Mozhlyvistʹ | Detali |
|---|---|
| **Avtomatychne vyyavlennya** | Kraulinh Katana + zond Httpx, deduplikovanyi spysok |
| **4-fazove navantazhuvalʹne testuvannya** | Ramp-Up → Sustained → Spike → Recovery |
| **Bahati metryky** | Zatrymka P50/P95/P99, req/s, klasyfikatsiya pomylok za fazamy |
| **Rezhym tryvalosti** | Fazy na osnovi chasu (`--duration 30s`) zamist' kilkosti zapytiv |
| **Porohy spovishchennya** | `exit 1` sumisnyy z CI/CD pry porushanni P99, chastoty pomylok abo RPS |
| **Podviyni zvity** | JSON (mashinochytnyy) + avtonomna temna panel' HTML |
| **Korystuvatski zaholovky** | Bearer-tokeny, zaholovky orendarya, kuki — poshyreni skriz' |
| **Hotovyy do Docker** | Odne zobrazhennya z vbudovanym Katana + Httpx |
| **Krossplatformenyy** | Linux, macOS, Windows — amd64 & arm64 |
| **Nulʹovi zalezhnosti** | Chystyy Go 1.21+ stdlib, dostatno `go install` |

---

## Shvydkyy start

```bash
# Vstanovlennya (potribno Go 1.21+)
go install github.com/umutozen/stormprobe/cmd/stormprobe@latest

# Bazovyy test — avtomatychne vyyavlennya, 4 fazy
stormprobe https://example.com

# Propustyty vyyavlennya, testuvaty lyshe korenevyy shlyakh
stormprobe --no-discovery https://example.com

# TLS + zaholovok auth + lyshe zvit JSON
stormprobe --insecure -H "Authorization: Bearer TOKEN" --format json https://api.example.com

# Test na osnovi tryvalosti: 30 sekund na fazu
stormprobe --insecure --duration 30s --no-discovery https://example.com

# CI/CD: exit 1 yakshcho P99 > 500ms abo pomylky > 5%
stormprobe --alert-p99 500 --alert-error-rate 5 https://example.com
```

---

## Vstanovlennya

### Variant 1 — `go install` *(rekomenduyetsya)*

```bash
go install github.com/umutozen/stormprobe/cmd/stormprobe@latest
stormprobe --help
```

### Variant 2 — Gotovyy binarnyy fayl *(bez Go)*

Zavantazhyty zi storinky [Releases](https://github.com/umutozen/stormprobe/releases).

```bash
# Linux / macOS
chmod +x stormprobe
sudo mv stormprobe /usr/local/bin/

# Windows PowerShell
Move-Item stormprobe.exe C:\Windows\System32\stormprobe.exe
```

### Variant 3 — Zbirka zi dzherel

```bash
git clone https://github.com/umutozen/stormprobe.git
cd stormprobe
go build -o stormprobe ./cmd/stormprobe      # Linux / macOS
go build -o stormprobe.exe ./cmd/stormprobe  # Windows
```

### Variant 4 — Docker

```bash
docker run --rm ghcr.io/umutozen/stormprobe:latest --insecure https://example.com

# Linux / macOS
docker run --rm -v $(pwd)/outputs:/app/outputs \
  ghcr.io/umutozen/stormprobe:latest --insecure --format both https://example.com

# Windows PowerShell
docker run --rm -v "${PWD}\outputs:/app/outputs" `
  ghcr.io/umutozen/stormprobe:latest --insecure --format both https://example.com
```

### Optsionalno — Instrumenty vyyavlennya

```bash
bash scripts/install-tools.sh
```

> Bez Katana/Httpx vykorystovuvaty `--no-discovery` abo `--endpoints`.

---

## Dovidka CLI

### Paralelizm i navantazhennya

| Prapor | Za zamovchuvannyam | Opys |
|---|---|---|
| `--concurrency-ramp` | `50` | Maksymalnyy paralelizm dlya fazy Ramp-Up |
| `--concurrency-sustained` | `75` | Paralelizm dlya fazy Sustained |
| `--concurrency-spike` | `250` | Maksymalnyy paralelizm dlya fazy Spike |
| `--req-per-worker` | `15` | Zapyty na voronku (ihnoruyet'sya pry `--duration`) |
| `--duration` | `0` (vymkneno) | Tryvalist' fazy, napr. `30s`, `1m` |
| `--timeout` | `10s` | Taymaut na zapyt |

### Vyyavlennya

| Prapor | Za zamovchuvannyam | Opys |
|---|---|---|
| `--no-discovery` | `false` | Propustyty Katana+Httpx, testuvaty lyshe koren' |
| `--endpoints` | `""` | Zavantazhyty kintsevi tochky z faylu |
| `--katana-path` | `""` | Korystuvatsky shlyakh do binarnoho faylu Katana |
| `--httpx-path` | `""` | Korystuvatsky shlyakh do binarnoho faylu Httpx |

### HTTP i bezpeka

| Prapor | Za zamovchuvannyam | Opys |
|---|---|---|
| `-H`, `--header` | — | Korystuvatsky zaholovok (povtoryuvanyy) |
| `--insecure` | `false` | Propustyty perevirku sertyfikata TLS |

### Vykhid

| Prapor | Za zamovchuvannyam | Opys |
|---|---|---|
| `--format` | `both` | Format zvitu: `json`, `html`, `both` |
| `--output` | `./outputs` | Kataloh vykhodu dlya zvitiv |

### Spovishchennya CI/CD

| Prapor | Za zamovchuvannyam | Opys |
|---|---|---|
| `--alert-p99` | `0` (vymkneno) | exit 1 yakshcho P99 (ms) perevishuye porih |
| `--alert-error-rate` | `0` (vymkneno) | exit 1 yakshcho chastota pomylok (%) perevishuye porih |
| `--alert-rps` | `0` (vymkneno) | exit 1 yakshcho req/s padaye nyzhche poroh |

---

## Fazy testuvannya

| Faza | Opys |
|---|---|
| **0 — Discovery** | Kraulinh Katana + zond Httpx, deduplikovanyi spysok |
| **1 — Ramp-Up** | Postupove zbilshennya paralelizmu: 5 → 15 → 30 → pik |
| **2 — Sustained** | 3 khvyli pry tsilovomu navantazhenni |
| **3 — Spike** | Raptove zbilshennya do piku, potim okholodzhennya |
| **4 — Recovery** | Perevirka stanu pislya 10s okholodzhennya |

---

## Vykhid i zvity

| Fayl | Opys |
|---|---|
| `stormprobe_report_<timestamp>.json` | Vsi metryky za fazamy, mashinochytnyy |
| `stormprobe_report_<timestamp>.html` | Avtonomna temna panel' z hrafikom SVG zatrymky |

---

## Vidmova vid vidpovidalnosti

Tsej instrument pryznachenyy **vyklyuchno dlya avtoryzovanoho testuvannya bezpeky ta otsinky produktyvnosti**. Pered provedennyam navantazhuvalʹnykh testiv neobkhidno otrymaty yavnyy pysʹmovyy dozvil vid vlasny ka tsilovoyi systemy. Avtory ne nesutʹ **zhod noyi vidpovidalnosti** za nenal ezh ne vykorystannya abo shkodu, zavdanu tsym instrumentom.

---

## Litsenziya

[MIT](../../LICENSE) © [Umut ÖZEN](https://github.com/umutozen)

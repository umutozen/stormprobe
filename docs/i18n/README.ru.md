[English](../../README.md) · [Türkçe](README.tr.md) · [Deutsch](README.de.md) · [Francais](README.fr.md) · [Espanol](README.es.md) · [Italiano](README.it.md) · [Portugues](README.pt-BR.md) · [Russkij](README.ru.md) · [Zhongwen Jian](README.zh-CN.md) · [Zhongwen Fan](README.zh-TW.md) · [Nihongo](README.ja.md) · [Hangugeo](README.ko.md) · [Al-Arabiyya](README.ar.md) · [Hindi](README.hi.md) · [Nederlands](README.nl.md) · [Polski](README.pl.md) · [Ukrainska](README.uk.md)

<p align="center">
  <img src="../../assets/banner.png" alt="StormProbe Banner" width="100%">
</p>

<h1 align="center">StormProbe</h1>

<p align="center">
  <strong>obnaruzhit · nagruzit · proanalizirovat</strong><br>
  HTTP-instrument nagruzochnogo testirovaniya proizvodstvennogo urovnya s avtonomnym obnaruzheniem konechnyh tochek. Nulевые vneshnie zavisimosti.
</p>

<p align="center">
  <a href="https://github.com/umutozen/stormprobe/actions/workflows/ci.yml">
    <img src="https://github.com/umutozen/stormprobe/actions/workflows/ci.yml/badge.svg" alt="CI">
  </a>
  <a href="https://goreportcard.com/report/github.com/umutozen/stormprobe">
    <img src="https://goreportcard.com/badge/github.com/umutozen/stormprobe" alt="Go Report Card">
  </a>
  <a href="https://github.com/umutozen/stormprobe/releases">
    <img src="https://img.shields.io/github/v/release/umutozen/stormprobe" alt="Poslednyaya versiya">
  </a>
  <a href="../../LICENSE">
    <img src="https://img.shields.io/badge/license-MIT-blue.svg" alt="Litsenziya MIT">
  </a>
</p>

---

## Chto takoe StormProbe?

StormProbe — eto **instrument nagruzochnogo HTTP-testirovaniya bez zavisimostej**, napisannyj na chist om Go. On avtomaticheski obnaruzhivaet konechnye tochki vashego prilozheniya s pomoshh'yu [Katana](https://github.com/projectdiscovery/katana) i [Httpx](https://github.com/projectdiscovery/httpx), zatem vypolnyaet strukturirovannyj 4-fazovyj nagruzochnyj test — **Ramp-Up → Sustained → Spike → Recovery** — i sozdaet otchety JSON i HTML s metrikami zaderzhki P50/P95/P99.

Razrabotan dlya:
- **Konvejerov CI/CD** — porogi opoveshcheniya s kodami vyhoda
- **Pentesterov** — avtomaticheskoe obnaruzhenie aktivnyh konechnyh tochek
- **DevOps-inzhenerov** — benchmarki do i posle razvyortyvaniya
- **Komand QA** — validaciya SLA proizvoditel'nosti

---

## Vozmozhnosti

| Vozmozhnost' | Podrobnosti |
|---|---|
| **Avtomaticheskoe obnaruzhenie** | Krauling Katana + zond Httpx, deduplirovanny spisok |
| **4-fazovyj nagruzochnyj test** | Ramp-Up → Sustained → Spike → Recovery |
| **Bogatye metriki** | Zaderzhka P50/P95/P99, req/s, klassifikaciya oshibok po fazam |
| **Rezhim dlitel'nosti** | Fazy na osnove vremeni (`--duration 30s`) vmesto kolichestva zaprosov |
| **Porogi opoveshcheniya** | `exit 1` sovmestimy s CI/CD pri narushenii P99, chastoty oshibok ili RPS |
| **Dvojnye otchety** | JSON (mashinochitaemy) + avtonomny temnaya panel' HTML |
| **Pol'zovatel'skie zagolovki** | Bearer-tokeny, zagolovki arendatora, kuki — rasprostraneny povsyudu |
| **Gotov k Docker** | Odino izobrazheniye s vstroennym Katana + Httpx |
| **Krossplatformenny** | Linux, macOS, Windows — amd64 & arm64 |
| **Nul evye zavisimosti** | Chisty Go 1.21+ stdlib, dostatochno `go install` |

---

## Bystry start

```bash
# Ustanovka (trebuetsya Go 1.21+)
go install github.com/umutozen/stormprobe/cmd/stormprobe@latest

# Bazovy test — avtomaticheskoe obnaruzhenie, 4 fazy
stormprobe https://example.com

# Propustit' obnaruzhenie, testirovat' tol'ko korenevoj put'
stormprobe --no-discovery https://example.com

# TLS + zagolovok auth + tol'ko otchet JSON
stormprobe --insecure -H "Authorization: Bearer TOKEN" --format json https://api.example.com

# Test na osnove dlitel'nosti: 30 sekund na fazu
stormprobe --insecure --duration 30s --no-discovery https://example.com

# CI/CD: exit 1 esli P99 > 500ms ili oshibki > 5%
stormprobe --alert-p99 500 --alert-error-rate 5 https://example.com
```

---

## Ustanovka

### Variant 1 — `go install` *(rekomenduetsya)*

```bash
go install github.com/umutozen/stormprobe/cmd/stormprobe@latest
stormprobe --help
```

### Variant 2 — Gotovy binarny fajl *(bez Go)*

Zagruzit' so stranitsy [Releases](https://github.com/umutozen/stormprobe/releases).

```bash
# Linux / macOS
chmod +x stormprobe
sudo mv stormprobe /usr/local/bin/

# Windows PowerShell
Move-Item stormprobe.exe C:\Windows\System32\stormprobe.exe
```

### Variant 3 — Sborka iz ishodnikov

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

### Opcionalno — Instrumenty obnaruzheniya

```bash
bash scripts/install-tools.sh
```

> Bez Katana/Httpx ispol'zovat' `--no-discovery` ili `--endpoints`.

---

## Spravka po CLI

### Parallelizm i nagruzka

| Flag | Po umolchaniyu | Opisanie |
|---|---|---|
| `--concurrency-ramp` | `50` | Maksimal'ny parallelizm dlya fazy Ramp-Up |
| `--concurrency-sustained` | `75` | Parallelizm dlya fazy Sustained |
| `--concurrency-spike` | `250` | Maksimal'ny parallelizm dlya fazy Spike |
| `--req-per-worker` | `15` | Zaprosy na voronku (ignoriruetsya pri `--duration`) |
| `--duration` | `0` (otklyucheno) | Dlitel'nost' fazy, napr. `30s`, `1m` |
| `--timeout` | `10s` | Tajm-aut na zapros |

### Obnaruzhenie

| Flag | Po umolchaniyu | Opisanie |
|---|---|---|
| `--no-discovery` | `false` | Propustit' Katana+Httpx, testirovat' tol'ko koren' |
| `--endpoints` | `""` | Zagruzit' konechnye tochki iz fajla |
| `--katana-path` | `""` | Pol'zovatel'sky put' k binarnому fajlu Katana |
| `--httpx-path` | `""` | Pol'zovatel'sky put' k binarnому fajlu Httpx |

### HTTP i bezopasnost'

| Flag | Po umolchaniyu | Opisanie |
|---|---|---|
| `-H`, `--header` | — | Pol'zovatel'sky zagolovok (povtoryaemy) |
| `--insecure` | `false` | Propustit' proverku sertifikata TLS |

### Vyvod

| Flag | Po umolchaniyu | Opisanie |
|---|---|---|
| `--format` | `both` | Format otcheta: `json`, `html`, `both` |
| `--output` | `./outputs` | Katalog vyvoda dlya otchetov |

### Opoveshcheniya CI/CD

| Flag | Po umolchaniyu | Opisanie |
|---|---|---|
| `--alert-p99` | `0` (otklyucheno) | exit 1 esli P99 (ms) prevyshaet porog |
| `--alert-error-rate` | `0` (otklyucheno) | exit 1 esli chastota oshibok (%) prevyshaet porog |
| `--alert-rps` | `0` (otklyucheno) | exit 1 esli req/s padaet nizhe poroga |

---

## Fazy testirovaniya

| Faza | Opisanie |
|---|---|
| **0 — Discovery** | Krauling Katana + zond Httpx, deduplirovanny spisok |
| **1 — Ramp-Up** | Postupennoe uvelichenie parallelizma: 5 → 15 → 30 → pik |
| **2 — Sustained** | 3 volny pri celevoj nagruzke |
| **3 — Spike** | Vneznoe uvelichenie do pika, zatem ostyvanie |
| **4 — Recovery** | Proverka sostoyaniya posle 10s ostyvanya |

---

## Vyvod i otchety

| Fajl | Opisanie |
|---|---|
| `stormprobe_report_<timestamp>.json` | Vse metriki po fazam, mashinochitaemy |
| `stormprobe_report_<timestamp>.html` | Avtonomna temnaya panel' s grafikom SVG zaderzh ki |

---

## Otkaznichestvo

Etot instrument prednaznachen **isklyuchitel'no dlya avtorizovannogo testirovaniya bezopasnosti i ocenki proizvoditel'nosti**. Pered provedeniem nagruzochnyh testov neobhodimo poluchit' yavnoe pis'mennoe razreshenie ot vladel'ca celevoj sistemy. Avtory ne nesut **nikakoj otvetstvennosti** za nenadzhezhдное ispol'zovanie ili ushcherb, prichinennyj etim instrumentom.

---

## Litsenziya

[MIT](../../LICENSE) © [Umut ÖZEN](https://github.com/umutozen)

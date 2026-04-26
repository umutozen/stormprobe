[English](../../README.md) · [Türkçe](README.tr.md) · [Deutsch](README.de.md) · [Francais](README.fr.md) · [Espanol](README.es.md) · [Italiano](README.it.md) · [Portugues](README.pt-BR.md) · [Russkij](README.ru.md) · [Zhongwen Jian](README.zh-CN.md) · [Zhongwen Fan](README.zh-TW.md) · [Nihongo](README.ja.md) · [Hangugeo](README.ko.md) · [Al-Arabiyya](README.ar.md) · [Hindi](README.hi.md) · [Nederlands](README.nl.md) · [Polski](README.pl.md) · [Ukrainska](README.uk.md)

<p align="center">
  <img src="../../assets/banner.png" alt="StormProbe Banner" width="100%">
</p>

<h1 align="center">StormProbe</h1>

<p align="center">
  <strong>iktishaf · tahmil · tahlil</strong><br>
  Adat ikhtibarat al-himl HTTP bitawir intaji ma' iktishaf mustaqil li-nuqat al-nihaya. Sifr taba'iyyat kharijiyya.
</p>

<p align="center">
  <a href="https://github.com/umutozen/stormprobe/actions/workflows/ci.yml">
    <img src="https://github.com/umutozen/stormprobe/actions/workflows/ci.yml/badge.svg" alt="CI">
  </a>
  <a href="https://goreportcard.com/report/github.com/umutozen/stormprobe">
    <img src="https://goreportcard.com/badge/github.com/umutozen/stormprobe" alt="Go Report Card">
  </a>
  <a href="https://github.com/umutozen/stormprobe/releases">
    <img src="https://img.shields.io/github/v/release/umutozen/stormprobe" alt="Akhir isdar">
  </a>
  <a href="../../LICENSE">
    <img src="https://img.shields.io/badge/license-MIT-blue.svg" alt="Rukhsa MIT">
  </a>
</p>

---

## Ma huwa StormProbe?

StormProbe huwa adat ikhtibarat himl HTTP **bila taba'iyyat** maktuba bi-Go al-naqiyya. Yaktashif taqri'an nuqat nihayat tatbiqak bi-[Katana](https://github.com/projectdiscovery/katana) wa [Httpx](https://github.com/projectdiscovery/httpx), thumma yunaffidh ikhtibaran muntazaman min 4 marahil — **Ramp-Up → Sustained → Spike → Recovery** — wa yuntij taqariran JSON wa HTML ma' maqayis kahna P50/P95/P99.

Summy li:
- **Khattat CI/CD** — 'atabat tanbih ma' akwad khurul
- **Mukhtabiri al-ikhtiraq** — iktishaf taqri'i li-nuqat nihaya fa'ila
- **Muhandisi DevOps** — mu'ayaraat qabl wa ba'd al-nushr
- **Fariq QA** — tahdid SLA al-ada'

---

## Al-Khasa'is

| Al-Khassiya | Al-Tafasil |
|---|---|
| **Al-Iktishaf al-Taqri'i** | Zaht Katana + misfar Httpx, qa'ima muza'a |
| **Ikhtibarat min 4 marahil** | Ramp-Up → Sustained → Spike → Recovery |
| **Maqayis wafira** | Kahna P50/P95/P99, req/s, tasnif al-akhta' lil-marhala |
| **Wad' al-mudda** | Marahil qa'ima 'ala al-waqt (`--duration 30s`) badalan min 'adad al-talabat |
| **'Atabat al-tanbih** | `exit 1` mutawafiq CI/CD 'ind khark P99 aw ns. al-akhta' aw RPS |
| **Taqarir muzdawija** | JSON (li-al-aja'z) + lawhat HTML dhalamiyya mustaqilla |
| **Ru'us mukhassasa** | Rumuz Bearer, ru'us al-mustaj'ir, al-ka'inat — muntashira fi kull makan |
| **Jahiz li-Docker** | Sura wahida ma' Katana + Httpx mudarrajin |
| **Muta'adid al-manassat** | Linux, macOS, Windows — amd64 & arm64 |
| **Sifr taba'iyyat** | Go 1.21+ stdlib naqiyy, yakfi `go install` |

---

## Al-Bida' al-Sari'

```bash
# Al-Tansib (yahtaj Go 1.21+)
go install github.com/umutozen/stormprobe/cmd/stormprobe@latest

# Ikhtibarat asasi — iktishaf taqri'i, 4 marahil
stormprobe https://example.com

# Tajawwuz al-iktishaf, ikhtibarat masarat al-jidhr faqat
stormprobe --no-discovery https://example.com

# TLS + ra's al-musadaqa + taqrir JSON faqat
stormprobe --insecure -H "Authorization: Bearer TOKEN" --format json https://api.example.com

# Ikhtibarat qa'im 'ala al-mudda: 30 thaniya lil-marhala
stormprobe --insecure --duration 30s --no-discovery https://example.com

# CI/CD: exit 1 idha kana P99 > 500ms aw al-akhta' > 5%
stormprobe --alert-p99 500 --alert-error-rate 5 https://example.com
```

---

## Al-Tansib

### Al-Khiyar 1 — `go install` *(al-muwassa)*

```bash
go install github.com/umutozen/stormprobe/cmd/stormprobe@latest
stormprobe --help
```

### Al-Khiyar 2 — Malaf thanawi mu'add (bila Go)

Tahmil min [safhat al-isdaraat](https://github.com/umutozen/stormprobe/releases).

```bash
# Linux / macOS
chmod +x stormprobe
sudo mv stormprobe /usr/local/bin/

# Windows PowerShell
Move-Item stormprobe.exe C:\Windows\System32\stormprobe.exe
```

### Al-Khiyar 3 — Al-Bina' min al-masadir

```bash
git clone https://github.com/umutozen/stormprobe.git
cd stormprobe
go build -o stormprobe ./cmd/stormprobe      # Linux / macOS
go build -o stormprobe.exe ./cmd/stormprobe  # Windows
```

### Al-Khiyar 4 — Docker

```bash
docker run --rm ghcr.io/umutozen/stormprobe:latest --insecure https://example.com

# Linux / macOS
docker run --rm -v $(pwd)/outputs:/app/outputs \
  ghcr.io/umutozen/stormprobe:latest --insecure --format both https://example.com

# Windows PowerShell
docker run --rm -v "${PWD}\outputs:/app/outputs" `
  ghcr.io/umutozen/stormprobe:latest --insecure --format both https://example.com
```

### Ikhtiyari — Adawat al-iktishaf

```bash
bash scripts/install-tools.sh
```

> Bila Katana/Httpx, istkhadim `--no-discovery` aw `--endpoints`.

---

## Marji' CLI

### Al-Tasamut wa al-Himl

| Al-'alam | Al-Iftiradi | Al-Wasf |
|---|---|---|
| `--concurrency-ramp` | `50` | Aqsa tasamut li-marhalt Ramp-Up |
| `--concurrency-sustained` | `75` | Al-Tasamut li-marhalt Sustained |
| `--concurrency-spike` | `250` | Aqsa tasamut li-marhalt Spike |
| `--req-per-worker` | `15` | Talabat lil-'amil (yuhmal idha tudin `--duration`) |
| `--duration` | `0` (mu'attal) | Muddat marhala, mathal `30s`, `1m` |
| `--timeout` | `10s` | Muhla lil-talab al-wahid |

### Al-Iktishaf

| Al-'alam | Al-Iftiradi | Al-Wasf |
|---|---|---|
| `--no-discovery` | `false` | Tajawwuz Katana+Httpx, ikhtibarat al-jidhr faqat |
| `--endpoints` | `""` | Tahmil nuqat al-nihaya min malaf |
| `--katana-path` | `""` | Masarat malaf Katana thanawi mukhassis |
| `--httpx-path` | `""` | Masarat malaf Httpx thanawi mukhassis |

### HTTP wa al-Aman

| Al-'alam | Al-Iftiradi | Al-Wasf |
|---|---|---|
| `-H`, `--header` | — | Ra's mukhassis (mutakarrir) |
| `--insecure` | `false` | Tajawwuz التحقق min shahadat TLS |

### Al-Ikhraja

| Al-'alam | Al-Iftiradi | Al-Wasf |
|---|---|---|
| `--format` | `both` | Siyaghat al-taqrir: `json`, `html`, `both` |
| `--output` | `./outputs` | Dalil al-ikhraja lil-taqarir |

### Tanbihat CI/CD

| Al-'alam | Al-Iftiradi | Al-Wasf |
|---|---|---|
| `--alert-p99` | `0` (mu'attal) | exit 1 idha tajawwaz P99 (ms) al-'ataba |
| `--alert-error-rate` | `0` (mu'attal) | exit 1 idha tajawwaz ns. al-akhta' (%) al-'ataba |
| `--alert-rps` | `0` (mu'attal) | exit 1 idha habbata req/s dun al-'ataba |

---

## Marahil al-Ikhtibarat

| Al-Marhala | Al-Wasf |
|---|---|
| **0 — Discovery** | Zaht Katana + misfar Httpx, qa'ima muza'a |
| **1 — Ramp-Up** | Ziyada tadrijiyya: 5 → 15 → 30 → dhurwa |
| **2 — Sustained** | 3 mawjat 'ind al-tasamut al-maqsud |
| **3 — Spike** | Irtifa' mufaji' lil-dhurwa, thumma tatrid |
| **4 — Recovery** | Fihras al-sihha ba'd 10 thawani tatrid |

---

## Al-Ikhraja wa Al-Taqarir

| Al-Malaf | Al-Wasf |
|---|---|
| `stormprobe_report_<timestamp>.json` | Jami' al-maqayis bil-marhala, li-al-aja'z |
| `stormprobe_report_<timestamp>.html` | Lawha dhalamiyya mustaqilla ma' rasm SVG lil-kahna |

---

## Tanbih

Hadhihi al-adat mukhassasa **hassa li-ikhtibarat al-aman al-musallahin wa taqyim al-ada' faqat**. Yajibu al-husul 'ala ijaza katibiyya sarriha min sahib al-nizam al-mustahdaq qabl ijra' ay ikhtibarat himl. La yahmil al-mu'allifun **ayy mas'uliyya** 'an isa'at al-isti'mal aw al-darar al-nasij 'an hadhihi al-adat.

---

## Al-Rukhsa

[MIT](../../LICENSE) © [Umut ÖZEN](https://github.com/umutozen)

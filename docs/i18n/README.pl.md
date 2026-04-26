[English](../../README.md) · [Türkçe](README.tr.md) · [Deutsch](README.de.md) · [Francais](README.fr.md) · [Espanol](README.es.md) · [Italiano](README.it.md) · [Portugues](README.pt-BR.md) · [Russkij](README.ru.md) · [Zhongwen Jian](README.zh-CN.md) · [Zhongwen Fan](README.zh-TW.md) · [Nihongo](README.ja.md) · [Hangugeo](README.ko.md) · [Al-Arabiyya](README.ar.md) · [Hindi](README.hi.md) · [Nederlands](README.nl.md) · [Polski](README.pl.md) · [Ukrainska](README.uk.md)

<p align="center">
  <img src="../../assets/banner.png" alt="StormProbe Banner" width="100%">
</p>

<h1 align="center">StormProbe</h1>

<p align="center">
  <strong>odkryc · zaladowac · przeanalizowac</strong><br>
  Narzedzie do testowania obciazenia HTTP na poziomie produkcyjnym z autonomicznym wykrywaniem punktow koncowych. Zero zewnetrznych zaleznosci.
</p>

<p align="center">
  <a href="https://github.com/umutozen/stormprobe/actions/workflows/ci.yml">
    <img src="https://github.com/umutozen/stormprobe/actions/workflows/ci.yml/badge.svg" alt="CI">
  </a>
  <a href="https://goreportcard.com/report/github.com/umutozen/stormprobe">
    <img src="https://goreportcard.com/badge/github.com/umutozen/stormprobe" alt="Go Report Card">
  </a>
  <a href="https://github.com/umutozen/stormprobe/releases">
    <img src="https://img.shields.io/github/v/release/umutozen/stormprobe" alt="Najnowsza wersja">
  </a>
  <a href="../../LICENSE">
    <img src="https://img.shields.io/badge/license-MIT-blue.svg" alt="Licencja MIT">
  </a>
</p>

---

## Czym jest StormProbe?

StormProbe to **narzedzie do testowania obciazenia HTTP bez zaleznosci** napisane w czystym Go. Automatycznie wykrywa punkty koncowe aplikacji za pomoca [Katana](https://github.com/projectdiscovery/katana) i [Httpx](https://github.com/projectdiscovery/httpx), nastepnie przeprowadza strukturyzowany 4-fazowy test obciazenia — **Ramp-Up → Sustained → Spike → Recovery** — i produkuje raporty JSON i HTML z metrykami opoznien P50/P95/P99.

Zaprojektowane dla:
- **Potoków CI/CD** — progi alertów z kodami wyjscia
- **Pentesterów** — automatyczne wykrywanie aktywnych punktów koncowych
- **Inzynierow DevOps** — benchmarki przed i po wdrozeniu
- **Zespolow QA** — walidacja SLA wydajnosci

---

## Funkcje

| Funkcja | Szczegoly |
|---|---|
| **Automatyczne wykrywanie** | Indeksowanie Katana + sonda Httpx, deduplikowana lista |
| **4-fazowy test obciazenia** | Ramp-Up → Sustained → Spike → Recovery |
| **Bogate metryki** | Opoznienie P50/P95/P99, req/s, klasyfikacja bledów na faze |
| **Tryb czasu trwania** | Fazy oparte na czasie (`--duration 30s`) zamiast liczby zapytan |
| **Progi alertów** | `exit 1` kompatybilny z CI/CD przy naruszeniu P99, wskaznika bledów lub RPS |
| **Podwójne raporty** | JSON (czytalny maszynowo) + autonomiczny ciemny pulpit HTML |
| **Niestandardowe naglówki** | Tokeny Bearer, naglówki dzierzawcy, cookies — propagowane wszedie |
| **Gotowy na Docker** | Pojedynczy obraz z wbudowanym Katana + Httpx |
| **Wieloplatformowy** | Linux, macOS, Windows — amd64 & arm64 |
| **Zero zaleznosci** | Czysty Go 1.21+ stdlib, `go install` wystarczy |

---

## Szybki start

```bash
# Instalacja (wymagane Go 1.21+)
go install github.com/umutozen/stormprobe/cmd/stormprobe@latest

# Podstawowy test — automatyczne wykrywanie, 4 fazy
stormprobe https://example.com

# Pominiecie wykrywania, testowanie tylko sciezki glównej
stormprobe --no-discovery https://example.com

# TLS + naglówek auth + tylko raport JSON
stormprobe --insecure -H "Authorization: Bearer TOKEN" --format json https://api.example.com

# Test oparty na czasie trwania: 30 sekund na faze
stormprobe --insecure --duration 30s --no-discovery https://example.com

# CI/CD: exit 1 jesli P99 > 500ms lub bledy > 5%
stormprobe --alert-p99 500 --alert-error-rate 5 https://example.com
```

---

## Instalacja

### Opcja 1 — `go install` *(zalecane)*

```bash
go install github.com/umutozen/stormprobe/cmd/stormprobe@latest
stormprobe --help
```

### Opcja 2 — Gotowy plik binarny *(bez Go)*

Pobrac ze strony [Releases](https://github.com/umutozen/stormprobe/releases).

```bash
# Linux / macOS
chmod +x stormprobe
sudo mv stormprobe /usr/local/bin/

# Windows PowerShell
Move-Item stormprobe.exe C:\Windows\System32\stormprobe.exe
```

### Opcja 3 — Kompilacja ze zrodel

```bash
git clone https://github.com/umutozen/stormprobe.git
cd stormprobe
go build -o stormprobe ./cmd/stormprobe      # Linux / macOS
go build -o stormprobe.exe ./cmd/stormprobe  # Windows
```

### Opcja 4 — Docker

```bash
docker run --rm ghcr.io/umutozen/stormprobe:latest --insecure https://example.com

# Linux / macOS
docker run --rm -v $(pwd)/outputs:/app/outputs \
  ghcr.io/umutozen/stormprobe:latest --insecure --format both https://example.com

# Windows PowerShell
docker run --rm -v "${PWD}\outputs:/app/outputs" `
  ghcr.io/umutozen/stormprobe:latest --insecure --format both https://example.com
```

### Opcjonalnie — Narzedzia wykrywania

```bash
bash scripts/install-tools.sh
```

> Bez Katana/Httpx uzywac `--no-discovery` lub `--endpoints`.

---

## Dokumentacja CLI

### Wspólbieznosc i obciazenie

| Flaga | Domyslnie | Opis |
|---|---|---|
| `--concurrency-ramp` | `50` | Maksymalna wspólbieznosc dla fazy Ramp-Up |
| `--concurrency-sustained` | `75` | Wspólbieznosc dla fazy Sustained |
| `--concurrency-spike` | `250` | Maksymalna wspólbieznosc dla fazy Spike |
| `--req-per-worker` | `15` | Zapytania na pracownika (ignorowane jesli `--duration` jest ustawione) |
| `--duration` | `0` (wylaczone) | Czas trwania fazy, np. `30s`, `1m` |
| `--timeout` | `10s` | Limit czasu na zapytanie |

### Wykrywanie

| Flaga | Domyslnie | Opis |
|---|---|---|
| `--no-discovery` | `false` | Pominiecie Katana+Httpx, testowanie tylko korzenia |
| `--endpoints` | `""` | Ladowanie punktów koncowych z pliku |
| `--katana-path` | `""` | Niestandardowa sciezka do pliku binarnego Katana |
| `--httpx-path` | `""` | Niestandardowa sciezka do pliku binarnego Httpx |

### HTTP i bezpieczenstwo

| Flaga | Domyslnie | Opis |
|---|---|---|
| `-H`, `--header` | — | Niestandardowy naglówek (powtarzalny) |
| `--insecure` | `false` | Pominiecie weryfikacji certyfikatu TLS |

### Wyjscie

| Flaga | Domyslnie | Opis |
|---|---|---|
| `--format` | `both` | Format raportu: `json`, `html`, `both` |
| `--output` | `./outputs` | Katalog wyjsciowy dla raportów |

### Alerty CI/CD

| Flaga | Domyslnie | Opis |
|---|---|---|
| `--alert-p99` | `0` (wylaczone) | exit 1 jesli P99 (ms) przekroczy próg |
| `--alert-error-rate` | `0` (wylaczone) | exit 1 jesli wskaznik bledów (%) przekroczy próg |
| `--alert-rps` | `0` (wylaczone) | exit 1 jesli req/s spadnie ponizej progu |

---

## Fazy testów

| Faza | Opis |
|---|---|
| **0 — Discovery** | Indeksowanie Katana + sonda Httpx, deduplikowana lista |
| **1 — Ramp-Up** | Stopniowy wzrost wspólbieznosci: 5 → 15 → 30 → szczyt |
| **2 — Sustained** | 3 fale przy docelowej wspólbieznosci |
| **3 — Spike** | Nagly wzrost do szczytu, nastepnie schladzanie |
| **4 — Recovery** | Kontrola stanu po 10s schladzania |

---

## Wyjscie i raporty

| Plik | Opis |
|---|---|
| `stormprobe_report_<timestamp>.json` | Wszystkie metryki na faze, czytalny maszynowo |
| `stormprobe_report_<timestamp>.html` | Autonomiczny ciemny pulpit z wykresem SVG opoznien |

---

## Zastrzezenie

Narzedzie to jest przeznaczone **wylacznie do autoryzowanych testów bezpieczenstwa i oceny wydajnosci**. Przed przeprowadzeniem testów obciazeniowych nalezy uzyskac wyrazna pisemna zgode wlasciciela systemu docelowego. Autorzy nie ponosza **zadnej odpowiedzialnosci** za niewlasciwe uzycie lub szkody spowodowane tym narzedziem.

---

## Licencja

[MIT](../../LICENSE) © [Umut ÖZEN](https://github.com/umutozen)

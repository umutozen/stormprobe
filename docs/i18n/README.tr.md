[English](../../README.md) · [Türkçe](README.tr.md) · [Deutsch](README.de.md) · [Francais](README.fr.md) · [Espanol](README.es.md) · [Italiano](README.it.md) · [Portugues](README.pt-BR.md) · [Russkij](README.ru.md) · [Zhongwen Jian](README.zh-CN.md) · [Zhongwen Fan](README.zh-TW.md) · [Nihongo](README.ja.md) · [Hangugeo](README.ko.md) · [Al-Arabiyya](README.ar.md) · [Hindi](README.hi.md) · [Nederlands](README.nl.md) · [Polski](README.pl.md) · [Ukrainska](README.uk.md)

<p align="center">
  <img src="../../assets/banner.png" alt="StormProbe Banner" width="100%">
</p>

<h1 align="center">StormProbe</h1>

<p align="center">
  <strong>kesfet · yukle · analiz et</strong><br>
  Otonom endpoint keşfi ile üretim düzeyinde HTTP yük test aracı. Sıfır harici bağımlılık.
</p>

<p align="center">
  <a href="https://github.com/umutozen/stormprobe/actions/workflows/ci.yml">
    <img src="https://github.com/umutozen/stormprobe/actions/workflows/ci.yml/badge.svg" alt="CI">
  </a>
  <a href="https://goreportcard.com/report/github.com/umutozen/stormprobe">
    <img src="https://goreportcard.com/badge/github.com/umutozen/stormprobe" alt="Go Report Card">
  </a>
  <a href="https://github.com/umutozen/stormprobe/releases">
    <img src="https://img.shields.io/github/v/release/umutozen/stormprobe" alt="Son Sürüm">
  </a>
  <a href="../../LICENSE">
    <img src="https://img.shields.io/badge/license-MIT-blue.svg" alt="MIT Lisansı">
  </a>
</p>

---

## StormProbe Nedir?

StormProbe, saf Go ile yazılmış **sıfır bağımlılıklı** bir HTTP yük test aracıdır. [Katana](https://github.com/projectdiscovery/katana) ve [Httpx](https://github.com/projectdiscovery/httpx) kullanarak uygulamanızın endpoint'lerini otomatik olarak keşfeder, ardından yapılandırılmış 4 aşamalı yük testi çalıştırır — **Ramp-Up → Sustained → Spike → Recovery** — ve P50/P95/P99 gecikme metrikleri içeren JSON ve HTML raporları üretir.

Şunlar için tasarlanmıştır:
- **CI/CD pipeline'ları** — çıkış kodları ile alert eşikleri
- **Sızma testçileri** — canlı endpoint'leri otomatik bulma, stres testi
- **DevOps mühendisleri** — dağıtım öncesi ve sonrası karşılaştırmalı test
- **QA ekipleri** — gerçek trafik örüntülerine karşı performans SLA doğrulaması

---

## Özellikler

| Özellik | Detay |
|---|---|
| **Otomatik Keşif** | Katana tarama + Httpx prob, tekilleştirilmiş canlı endpoint listesi |
| **4 Aşamalı Yük Testi** | Ramp-Up → Sustained → Spike → Recovery |
| **Zengin Metrikler** | P50 / P95 / P99 gecikme, req/s, faz bazında hata sınıflandırması |
| **Süre Modu** | İstek sayısı yerine zaman bazlı fazlar (`--duration 30s`) |
| **Alert Eşikleri** | CI/CD uyumlu P99, hata oranı veya RPS ihlallerinde `exit 1` |
| **Çift Rapor** | JSON (makine okunabilir) + bağımsız HTML koyu tema panosu |
| **Özel Header** | Bearer token, tenant header, cookie — her yere iletilir |
| **Docker Hazır** | Katana + Httpx dahil tek imaj |
| **Çok Platformlu** | Linux, macOS, Windows — amd64 & arm64 |
| **Sıfır Bağımlılık** | Saf Go 1.21+ stdlib, `go install` yeterli |

---

## Hızlı Başlangıç

```bash
# Kurulum (Go 1.21+ gerekli)
go install github.com/umutozen/stormprobe/cmd/stormprobe@latest

# Temel test — endpoint'leri otomatik keşfeder, 4 fazı çalıştırır
stormprobe https://example.com

# Keşfi atla, sadece kök yolu test et
stormprobe --no-discovery https://example.com

# TLS + auth header + sadece JSON rapor
stormprobe --insecure -H "Authorization: Bearer TOKEN" --format json https://api.example.com

# Süre bazlı test: her faz 30 saniye çalışır
stormprobe --insecure --duration 30s --no-discovery https://example.com

# CI/CD: P99 > 500ms veya hata > %5 ise exit 1
stormprobe --alert-p99 500 --alert-error-rate 5 https://example.com
```

---

## Kurulum

### Seçenek 1 — `go install` *(önerilen)*

Go 1.21+ gerektirir. Binary doğrudan `$GOPATH/bin`'e kurulur.

```bash
go install github.com/umutozen/stormprobe/cmd/stormprobe@latest
stormprobe --help
```

### Seçenek 2 — Hazır Binary *(Go gerektirmez)*

[Releases sayfasından](https://github.com/umutozen/stormprobe/releases) indirin.

```bash
# Linux / macOS
chmod +x stormprobe
sudo mv stormprobe /usr/local/bin/

# Windows PowerShell
Move-Item stormprobe.exe C:\Windows\System32\stormprobe.exe
```

### Seçenek 3 — Kaynaktan Derleme

```bash
git clone https://github.com/umutozen/stormprobe.git
cd stormprobe
go build -o stormprobe ./cmd/stormprobe      # Linux / macOS
go build -o stormprobe.exe ./cmd/stormprobe  # Windows
```

### Seçenek 4 — Docker

```bash
# Rapor kaydedilmeden hızlı test
docker run --rm ghcr.io/umutozen/stormprobe:latest --insecure https://example.com

# HTML + JSON rapor kaydet
# Linux / macOS
docker run --rm -v $(pwd)/outputs:/app/outputs \
  ghcr.io/umutozen/stormprobe:latest --insecure --format both https://example.com

# Windows PowerShell
docker run --rm -v "${PWD}\outputs:/app/outputs" `
  ghcr.io/umutozen/stormprobe:latest --insecure --format both https://example.com
```

### Isteğe Bağlı — Keşif Araçları

```bash
bash scripts/install-tools.sh
```

> Katana/Httpx yoksa `--no-discovery` veya `--endpoints` kullanın.

---

## CLI Referansı

```
stormprobe [flags] <hedef-url>
```

### Eş Zamanlılık ve Yük

| Flag | Varsayılan | Açıklama |
|---|---|---|
| `--concurrency-ramp` | `50` | Ramp-Up fazı için maksimum eş zamanlılık |
| `--concurrency-sustained` | `75` | Sustained fazı için eş zamanlılık |
| `--concurrency-spike` | `250` | Spike fazı için maksimum eş zamanlılık |
| `--req-per-worker` | `15` | Worker başına istek sayısı (`--duration` ayarlıysa yok sayılır) |
| `--duration` | `0` (kapalı) | Zaman bazlı faz süresi, ör. `30s`, `1m` |
| `--timeout` | `10s` | İstek başına zaman aşımı |

### Keşif

| Flag | Varsayılan | Açıklama |
|---|---|---|
| `--no-discovery` | `false` | Katana+Httpx'i atla, sadece kök yolu test et |
| `--endpoints` | `""` | Dosyadan endpoint listesi yükle (satır başına bir yol) |
| `--katana-path` | `""` | Özel Katana binary yolu |
| `--httpx-path` | `""` | Özel Httpx binary yolu |

### HTTP ve Güvenlik

| Flag | Varsayılan | Açıklama |
|---|---|---|
| `-H`, `--header` | — | Özel header (tekrarlanabilir): `-H "Authorization: Bearer TOKEN"` |
| `--insecure` | `false` | TLS sertifika doğrulamasını atla |

### Çıktı

| Flag | Varsayılan | Açıklama |
|---|---|---|
| `--format` | `both` | Rapor formatı: `json`, `html`, `both` |
| `--output` | `./outputs` | Raporlar için çıktı dizini |

### CI/CD Alert'leri

| Flag | Varsayılan | Açıklama |
|---|---|---|
| `--alert-p99` | `0` (kapalı) | Herhangi bir fazda P99 (ms) eşiği aşılırsa exit 1 |
| `--alert-error-rate` | `0` (kapalı) | Herhangi bir fazda hata oranı (%) eşiği aşılırsa exit 1 |
| `--alert-rps` | `0` (kapalı) | Herhangi bir fazda req/s eşiğin altına düşerse exit 1 |

---

## Kullanım Örnekleri

```bash
# Otomatik keşif + tam test
stormprobe https://example.com

# Keşfi atla (bilinen hedefler için daha hızlı)
stormprobe --no-discovery https://example.com

# Özel endpoint listesi
stormprobe --endpoints endpoints.txt https://example.com

# Kimlik doğrulamalı API stres testi
stormprobe --insecure -H "Authorization: Bearer TOKEN" -H "X-Tenant: acme" https://api.example.com

# Yüksek eş zamanlılıklı spike (500 sanal kullanıcı)
stormprobe --insecure --concurrency-spike 500 --req-per-worker 20 https://example.com

# Süre bazlı test: her faz 1 dakika çalışır
stormprobe --insecure --duration 1m --concurrency-spike 300 https://example.com

# Sadece JSON rapor, özel dizine kaydet
stormprobe --insecure --format json --output ./results https://example.com

# CI/CD: P99 > 500ms veya hata > %5 veya RPS < 100 ise exit 1
stormprobe --alert-p99 500 --alert-error-rate 5 --alert-rps 100 https://example.com
```

### CI/CD Entegrasyonu (GitHub Actions)

```yaml
- name: Performance gate
  run: |
    go install github.com/umutozen/stormprobe/cmd/stormprobe@latest
    stormprobe \
      --no-discovery \
      --duration 30s \
      --alert-p99 500 \
      --alert-error-rate 2 \
      --alert-rps 50 \
      https://staging.example.com
```

---

## Test Fazları

```
Keşif ──► Ramp-Up ──► Sustained ──► Spike ──► Recovery
```

| Faz | Açıklama |
|---|---|
| **0 — Discovery** | Katana tarama + Httpx prob, tekilleştirilmiş endpoint listesi |
| **1 — Ramp-Up** | Kademeli eş zamanlılık artışı: 5 → 15 → 30 → zirve kullanıcı |
| **2 — Sustained** | Hedef eş zamanlılıkta 3 dalga, kararlı hal bozulmasını ölçer |
| **3 — Spike** | Zirve eş zamanlılığa ani artış, ardından soğuma |
| **4 — Recovery** | 10 saniyelik soğuma sonrası spike sonrası sağlık kontrolü |

---

## Çıktı ve Raporlar

Raporlar `./outputs/` klasörüne kaydedilir (`--output` ile değiştirilebilir).

| Dosya | Açıklama |
|---|---|
| `stormprobe_report_<timestamp>.json` | Faz başına tüm metrikler, makine okunabilir |
| `stormprobe_report_<timestamp>.html` | SVG gecikme grafiği içeren bağımsız koyu tema panosu |

**Faz başına ölçülen metrikler:** basari%, ort / P50 / P95 / P99 gecikme (ms), req/s, hata dağılımı (Timeout, Connection Reset, Connection Refused, HTTP 4xx, HTTP 5xx, Diger), HTTP durum kodu dağılımı.

---

## Proje Yapısı

```
stormprobe/
├── cmd/stormprobe/
│   └── main.go                  CLI giriş noktası, flag ayrıştırma, faz orkestrasyonu
├── internal/
│   ├── config/config.go         Paylaşılan tipler, faz üreticileri, AlertConfig
│   ├── alert/alert.go           Eşik kontrolleri, CI/CD exit kodu mantığı
│   ├── metrics/
│   │   ├── latency.go           Yüzdelik hesaplama (P50/P95/P99)
│   │   └── errors.go            Hata sınıflandırması
│   ├── discovery/
│   │   ├── discovery.go         Discover() genel arayüz
│   │   ├── katana.go            Katana tarayıcı entegrasyonu
│   │   └── httpx.go             Httpx prob entegrasyonu
│   ├── runner/
│   │   ├── phase.go             Faz orkestrasyonu, süre/sayı modları
│   │   └── worker.go            Goroutine havuzu, worker bazında HTTP çalıştırma
│   └── report/
│       ├── json.go              JSON rapor yazıcı
│       └── html.go              Bağımsız HTML koyu tema panosu
├── Dockerfile                   Çok aşamalı: Go derleme + katana + httpx
├── .goreleaser.yml              Çapraz platform sürüm yapılandırması
├── Makefile                     build, test, lint hedefleri
└── config.example.yml           Açıklamalı yapılandırma referansı
```

---

## Gereksinimler

| Gereksinim | Notlar |
|---|---|
| **Go 1.21+** | `go install` veya kaynaktan derleme için gerekli |
| **Katana** | Isteğe bağlı — otomatik keşif için gerekli |
| **Httpx** | Isteğe bağlı — otomatik keşif için gerekli |
| **Docker** | Isteğe bağlı — Go kurulumuna alternatif |

> Katana/Httpx olmadan `--no-discovery` veya `--endpoints` kullanın.

---

## Sorumluluk Reddi

Bu araç **yalnızca yetkili güvenlik testleri ve performans değerlendirmesi** içindir. Herhangi bir yük testi çalıştırmadan önce hedef sistem sahibinden açık yazılı izin almanız gerekmektedir. Yetkisiz kullanım yerel, ulusal veya uluslararası yasaları ihlal edebilir. Yazarlar kötüye kullanımdan doğan hiçbir zarardan **sorumlu değildir**.

---

## Katkı

Katkı kılavuzu ve geliştirme kurulumu için [CONTRIBUTING.md](../../CONTRIBUTING.md) dosyasına bakın.

## Lisans

[MIT](../../LICENSE) © [Umut ÖZEN](https://github.com/umutozen)

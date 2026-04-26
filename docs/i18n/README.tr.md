[English](../../README.md) · [Türkçe](README.tr.md) · [Deutsch](README.de.md) · [Francais](README.fr.md) · [Espanol](README.es.md) · [Italiano](README.it.md) · [Portugues](README.pt-BR.md) · [Russkij](README.ru.md) · [Zhongwen Jian](README.zh-CN.md) · [Zhongwen Fan](README.zh-TW.md) · [Nihongo](README.ja.md) · [Hangugeo](README.ko.md) · [Al-Arabiyya](README.ar.md) · [Hindi](README.hi.md) · [Nederlands](README.nl.md) · [Polski](README.pl.md) · [Ukrainska](README.uk.md)

# StormProbe

<p align="center">
  <img src="../../assets/banner.png" alt="StormProbe" width="100%">
</p>

**keşfet. yükle. analiz et.**

Otonom endpoint (uç nokta) keşfi ile üretim düzeyinde HTTP yük test aracı. Sıfır bağımlılık.

[![CI](https://github.com/umutozen/stormprobe/actions/workflows/ci.yml/badge.svg)](https://github.com/umutozen/stormprobe/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/umutozen/stormprobe)](https://goreportcard.com/report/github.com/umutozen/stormprobe)

---

## Özellikler

- **Otomatik Keşif** -- Katana tarama + Httpx prob ile canlı endpoint'ler otomatik bulunur
- **Çok Fazlı Test** -- Ramp-Up (kademeli artış), Sustained (sürekli), Spike (ani artış), Recovery (toparlanma)
- **Detaylı Metrikler** -- P50 / P95 / P99 latency (gecikme), req/s, hata sınıflandırması
- **Çift Rapor** -- JSON (makine tarafından okunabilir) + HTML (görsel koyu tema panosu)
- **Sıfır Bağımlılık** -- Saf Go standart kütüphanesi, üçüncü parti modül yok
- **Docker Hazır** -- Katana + Httpx dahil tek komut
- **Çok Platformlu** -- GoReleaser ile Linux, macOS, Windows ikili dosyaları

## Hızlı Başlangıç

```bash
go run ./cmd https://example.com

go build -o stormprobe ./cmd
./stormprobe --insecure https://example.com
```

### Docker

```bash
# Hızlı test, rapor kaydedilmez
docker run --rm ghcr.io/umutozen/stormprobe:latest --insecure https://example.com

# Rapor kaydederek tam test (Linux / macOS)
docker run --rm -v $(pwd)/outputs:/app/outputs ghcr.io/umutozen/stormprobe:latest --insecure --format both https://example.com

# Rapor kaydederek tam test (Windows PowerShell)
docker run --rm -v "${PWD}\outputs:/app/outputs" ghcr.io/umutozen/stormprobe:latest --insecure --format both https://example.com

# Yüksek yük ani artış testi (500 eş zamanlı kullanıcı)
docker run --rm ghcr.io/umutozen/stormprobe:latest --insecure --concurrency-spike 500 https://example.com

# Keşfi atla, yalnızca kök yolu test et
docker run --rm ghcr.io/umutozen/stormprobe:latest --insecure --no-discovery https://example.com

# Özel HTTP başlıkları (Yetkilendirme, özel kiracı, vb.)
docker run --rm ghcr.io/umutozen/stormprobe:latest --insecure -H "Authorization: Bearer TOKEN" -H "X-Tenant: acme" https://api.example.com
```

## Kullanım

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
  --header, -H            string    Custom HTTP header (repeatable): -H 'Authorization: Bearer TOKEN'
  --duration              duration  Per-phase duration (e.g. 30s, 1m). Overrides req-per-worker when set
  --alert-p99             float     Fail (exit 1) if P99 latency exceeds Xms in any phase
  --alert-error-rate      float     Fail (exit 1) if error rate exceeds X%% in any phase
  --alert-rps             float     Fail (exit 1) if req/s falls below X in any phase
```

### Örnekler

```bash
stormprobe --insecure --concurrency-spike 500 --req-per-worker 20 https://example.com

stormprobe --insecure --endpoints endpoints.txt https://example.com

stormprobe --insecure --format json --output ./results https://example.com

stormprobe --insecure -H "Authorization: Bearer TOKEN" -H "X-Tenant: acme" https://api.example.com

stormprobe --insecure --duration 30s --no-discovery https://example.com

stormprobe --insecure --duration 1m --concurrency-spike 300 https://example.com

stormprobe --insecure --alert-p99 500 --alert-error-rate 5 https://example.com
```

## Test Fazları

| Faz | Açıklama |
|---|---|
| **0 -- Discovery** | Katana tarama + Httpx prob, tekilleştirilmiş endpoint listesi |
| **1 -- Ramp-Up** | Kademeli artış: 5, 15, 30, 50 virtual users (sanal kullanıcılar) |
| **2 -- Sustained** | Hedef concurrency (eş zamanlılık) seviyesinde 3 dalga, bozulma ölçümü |
| **3 -- Spike** | Zirveye ani patlama, ardından soğuma |
| **4 -- Recovery** | 10 saniye soğuma sonrası spike sonrası sağlık kontrolü |

## Çıktı

| Dosya | Açıklama |
|---|---|
| `stormprobe_report_<timestamp>.json` | Tüm metrikleri içeren makine tarafından okunabilir sonuçlar |
| `stormprobe_report_<timestamp>.html` | Görsel pano, herhangi bir tarayıcıda açılır |

## Proje Yapısı

```
stormprobe/
    cmd/main.go                    CLI giriş noktası
    internal/
        config/config.go           Paylaşılan tipler, faz üreticileri
        metrics/
            latency.go             Yüzdelik hesaplama
            errors.go              Hata sınıflandırması
        discovery/
            discovery.go           Discover() genel arayüz
            katana.go              Katana tarayıcı entegrasyonu
            httpx.go               Httpx prob entegrasyonu
        runner/
            phase.go               Faz orkestrasyonu
            worker.go              Goroutine havuzu, worker başına RNG
        report/
            json.go                JSON rapor yazıcı
            html.go                Bağımsız HTML rapor
    Dockerfile                     Katana + httpx ile çok aşamalı
    .goreleaser.yml                Çapraz derleme yapılandırması
    Makefile                       Derleme, test, lint hedefleri
    config.example.yml             Faz varsayılanları referansı
```

## Gereksinimler

- Go 1.21+
- [Katana](https://github.com/projectdiscovery/katana) (ProjectDiscovery) -- isteğe bağlı, otomatik keşif için
- [Httpx](https://github.com/projectdiscovery/httpx) (ProjectDiscovery) -- isteğe bağlı, endpoint doğrulama için

```bash
bash scripts/install-tools.sh
```

## Sorumluluk Reddi

Bu araç **yalnızca yetkili güvenlik testleri ve performans değerlendirmesi** içindir. Herhangi bir yük testi çalıştırmadan önce hedef sistem sahibinden açık yazılı izin almanız gerekmektedir. Sahip olmadığınız veya test etme izniniz olmayan sistemlere karşı bu aracın yetkisiz kullanımı yerel, ulusal veya uluslararası yasaları ihlal edebilir. Yazarlar, bu aracın kötüye kullanımından veya neden olduğu zararlardan hiçbir sorumluluk kabul etmez.

## Lisans

MIT

## Yazar

**Umut ÖZEN** -- [@umutozen](https://github.com/umutozen)

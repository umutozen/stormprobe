# StormProbe

<p align="center">
  <img src="../../assets/banner.png" alt="StormProbe Banner" width="100%">
</p>

<h1 align="center">StormProbe</h1>

<p align="center">
  <strong>crawl · discover · stress · decide</strong><br>
  Özerk endpoint keşfi ve production hazırlık kararları ile modern HTTP yük testi.
</p>

<p align="center">
  <strong>Sadece benchmark değil.</strong><br>
  StormProbe sadece test etmez. Sisteminizin gerçekten daha güvenli hale gelip gelmediğini açıklar.
</p>

---

## Neden StormProbe?

Geleneksel yük test araçları sayı gösterir.<br>StormProbe karar verir.

Endpoint'leri elle seçmek ve ham gecikme tablolarını yorumlamak yerine StormProbe canlı endpoint'leri otomatik keşfeder, aşamalı stres testi yapar, doygunluk örüntülerini tespit eder ve aksiyonlu önerilerle son karar üretir.

### Geleneksel Yük Testi Araçları vs StormProbe

| Geleneksel Araçlar | StormProbe |
|---|---|
| Manuel endpoint seçimi | Canlı endpoint'leri otomatik keşfeder |
| Ham benchmark çıktısı | Son karar + darboğaz teşhisi |
| Manuel yorum gerektirir | Aksiyonlu öneriler |
| Tek stres deseni | Ramp-Up + Sustained + Spike + Recovery |
| "Bir şeyler yavaş" | "Nerede ve neden bozulduğunu söylüyor" |

---

## Ne Yapar?

StormProbe Go ile yazılmış production kalitesinde HTTP yük testi CLI'dır.

- Endpoint keşfi
- Aşamalı concurrency testi
- Gecikme yüzdelik analizi
- Throughput doygunluk tespiti
- Recovery doğrulaması
- CI/CD eşik uyarıları
- JSON + HTML raporlama
- Production hazırlık kararları

---

## Gerçek Çıktı Örneği

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

## Özellikler

| Özellik | Detay |
|---|---|
| Otomatik Keşif | Tekrarlananları eleyen Katana crawl + Httpx probe |
| 4 Aşamalı Yük Testi | Ramp-Up → Sustained → Spike → Recovery |
| Karar Motoru | Güvenli eşzamanlılık, bozulma noktası, darboğaz sınıflandırması |
| Zengin Metrikler | Avg / P50 / P95 / P99 gecikme + req/s + hata sınıflandırması |
| Throughput Analizi | Sistem tamamen çökmeden önce doygunluk tespit eder |
| Recovery Doğrulaması | Trafik patlamalarından sonra sistemin toparlanıp toparlanmadığını doğrular |
| Uyarı Eşikleri | Gecikme, hata oranı ve RPS için CI/CD uyumlu çıkış kodları |
| Çift Raporlama | JSON + self-contained HTML dashboard |
| Özel Başlıklar | Auth token'ları, tenant header'ları, çerezler |
| Docker Uyumlu | Dahili araçlarla birlikte tek imaj |
| Çapraz Platform | Linux, macOS, Windows — amd64 & arm64 |
| Saf Go Çekirdeği | Hızlı kurulum, minimum operasyonel yük |

---

## Hızlı Başlangıç

```bash
# Kurulum
go install github.com/umutozen/stormprobe/cmd/stormprobe@latest

# Otomatik keşif ile tam test
stormprobe https://example.com

# Keşfi atla, sadece kök dizini test et
stormprobe --no-discovery https://example.com

# Auth header ile API testi
stormprobe --insecure \
  -H "Authorization: Bearer TOKEN" \
  https://api.example.com

# CI/CD kapısı
stormprobe \\
  --alert-p99 500 \\
  --alert-error-rate 5 \\
  https://example.com

# JSON config dosyası ile çalıştır
stormprobe --config stormprobe.json
```

---

## Test Akışı

```text
Keşif → Ramp-Up → Sustained → Spike → Recovery → Karar
```

| Aşama | Amaç |
|---|---|
| Discovery | Canlı endpoint'leri taratır ve doğrular |
| Ramp-Up | Kademeli olarak concurrency'i artırır |
| Sustained | Kararlı durum performansını doğrular |
| Spike | Ani trafik patlaması simülasyonu |
| Recovery | Spike sonrası sistem sağlığını doğrular |
| Verdict | Güvenli concurrency + darboğaz analizi |

---

## Kurulum

### Go Install

```bash
go install github.com/umutozen/stormprobe/cmd/stormprobe@latest
```

Önerilen yol budur.

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

Bu ayar, performans regresyonlarının deployment süreçlerini iptal etmesini (fail) sağlar.

---

## Kullanım Senaryoları

**Production Öncesi** — Concurrency değişikliklerinin gecikme hedeflerini bozup bozmadığını doğrular.

**Güvenlik + Performans** — Gizli canlı endpoint'leri keşfeder ve güvenle test eder.

**Altyapı Denetimleri** — Trafik artışları olmadan önce doygunluk noktalarını ölçer.

**Regresyon Tespiti** — CI üzerinde deployment öncesi/sonrası performansı karşılaştırır.

---

## Flags

| Flag | Açıklama |
|---|---|
| `--duration` | Her fazı X süre boyunca çalıştır (ör. `30s`, `1m`) |
| `--alert-p99` | P99 gecikme eşiği (%s ms) aşılırsa 1 kodlu hata çıkışı ver |
| `--alert-error-rate` | Hata oranı eşiği (%s %%) aşılırsa 1 kodlu hata çıkışı ver |
| `--alert-rps` | Req/s eşiğin altına düşerse 1 kodlu hata çıkışı ver |
| `-H` / `--header` | Özel HTTP header (tekrarlanabilir) |
| `--insecure` | TLS doğrulamasını atla |
| `--no-discovery` | Keşfi atla, sadece kök dizini test et |
| `--endpoints` | Özel endpoint listesini dosyadan yükle |
| `--format` | Rapor formatı: `json`, `html`, `both` |
| `--concurrency-ramp` | Ramp-Up için maksimum concurrency |
| `--concurrency-sustained` | Sustained concurrency oranı |
| `--concurrency-spike` | Spike maksimum concurrency oranı |
| `--req-per-worker` | Worker başına istek sayısı |
| `--timeout` | Bir istek için en fazla beklenecek süre(timeout) |
| `--config` | JSON config dosyasından ayar yükle |
| `--auto-profile` | Sunucu yığınını otomatik algıla ve tanı profili uygula (varsayılan: açık) |
| `--profile` | Manuel profil: iis, tomcat, php |
| `--no-fingerprint` | Sunucu parmak izi tespitini atla |
| `--show-profile` | Algılanan yığın ve profili göster, ardından çık |

---

## Gereksinimler

| Araç | Not |
|---|---|
| Go 1.21+ | Go install için gereklidir |
| Katana | Opsiyonel, endpoint taraması için |
| Httpx | Opsiyonel, endpoint keşfi için |
| Docker | Opsiyonel çalışma ortamı |

---

## Uyarı

Bu araç yalnızca yetkili testler için tasarlanmıştır.

Herhangi bir hedef sistemde yük testleri çalıştırmadan önce daima açık izin alın.

İzinsiz kullanım yerel yasalara ve politikalara aykırı olabilir.

---

## Katkıda Bulunma

Katkılarınızı bekliyoruz. Geliştirme kurulumu ve katkı yönergeleri için `CONTRIBUTING.md` dosyasına başvurun.

## Güvenlik

Sorumlu güvenlik açığı bildirimi için `SECURITY.md` metnini okuyun.

## Lisans

MIT License © Umut ÖZEN

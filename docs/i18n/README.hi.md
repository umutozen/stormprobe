[English](../../README.md) · [Türkçe](README.tr.md) · [Deutsch](README.de.md) · [Francais](README.fr.md) · [Espanol](README.es.md) · [Italiano](README.it.md) · [Portugues](README.pt-BR.md) · [Russkij](README.ru.md) · [Zhongwen Jian](README.zh-CN.md) · [Zhongwen Fan](README.zh-TW.md) · [Nihongo](README.ja.md) · [Hangugeo](README.ko.md) · [Al-Arabiyya](README.ar.md) · [Hindi](README.hi.md) · [Nederlands](README.nl.md) · [Polski](README.pl.md) · [Ukrainska](README.uk.md)

# StormProbe

<p align="center">
  <img src="../../assets/banner.png" alt="StormProbe" width="100%">
</p>

**khoj karo. load karo. vishleshan karo.**

Svayam endpoint (ant bindu) khoj ke saath utpaadan-star ka HTTP load testing upakaran. Shunya nirbharata.

[![CI](https://github.com/umutozen/stormprobe/actions/workflows/ci.yml/badge.svg)](https://github.com/umutozen/stormprobe/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/umutozen/stormprobe)](https://goreportcard.com/report/github.com/umutozen/stormprobe)

---

## Visheshtaen

- **Svachalit khoj** -- Katana crawl + Httpx probe svachalit roop se sakriy endpoint khojta hai
- **Bahu-charan parikshan** -- Ramp-Up (dheere-dheere vriddhi), Sustained (nirantar), Spike (bhar shikhar), Recovery (punarprapti)
- **Vistrit maapdand** -- P50 / P95 / P99 latency (vilambita), req/s, truti vargikaran
- **Dohari report** -- JSON (machine padhaniy) + HTML (drishya dark theme dashboard)
- **Shunya nirbharata** -- Shuddh Go maanak pustakalay, koi tritiya-paksha module nahin
- **Docker taiyar** -- Katana + Httpx sahit ek aadesh
- **Cross-platform** -- GoReleaser ke madhyam se Linux, macOS, Windows binary

## Tvarit aarambh

```bash
go run ./cmd https://example.com

go build -o stormprobe ./cmd
./stormprobe --insecure https://example.com
```

### Docker

```bash
# त्वरित परीक्षण, कोई रिपोर्ट सहेजी नहीं
docker run --rm ghcr.io/umutozen/stormprobe:latest --insecure https://example.com

# ./outputs में रिपोर्ट सहेजकर पूर्ण परीक्षण (Linux / macOS)
docker run --rm -v $(pwd)/outputs:/app/outputs ghcr.io/umutozen/stormprobe:latest --insecure --format both https://example.com

# ./outputs में रिपोर्ट सहेजकर पूर्ण परीक्षण (Windows PowerShell)
docker run --rm -v "${PWD}\outputs:/app/outputs" ghcr.io/umutozen/stormprobe:latest --insecure --format both https://example.com

# उच्च-भार स्पाइक परीक्षण (500 समवर्ती उपयोगकर्ता)
docker run --rm ghcr.io/umutozen/stormprobe:latest --insecure --concurrency-spike 500 https://example.com

# डिस्कवरी छोड़ें, केवल रूट पथ का परीक्षण करें
docker run --rm ghcr.io/umutozen/stormprobe:latest --insecure --no-discovery https://example.com

# कस्टम HTTP हेडर (प्राधिकरण, कस्टम टेनेंट, आदि)
docker run --rm ghcr.io/umutozen/stormprobe:latest --insecure -H "Authorization: Bearer TOKEN" -H "X-Tenant: acme" https://api.example.com
```

## Upayog

```
stormprobe [flags] <लक्ष्य-URL>

Flags:
  --concurrency-ramp      int       Ramp-Up के लिए अधिकतम समवर्तिता (डिफ़ॉल्ट 50)
  --concurrency-sustained int       Sustained चरण के लिए समवर्तिता (डिफ़ॉल्ट 75)
  --concurrency-spike     int       Spike के लिए अधिकतम समवर्तिता (डिफ़ॉल्ट 250)
  --req-per-worker        int       प्रति worker प्रति चरण अनुरोध (डिफ़ॉल्ट 15)
  --timeout               duration  प्रति अनुरोध टाइमआउट (डिफ़ॉल्ट 10s)
  --endpoints             string    Endpoint सूची फ़ाइल (प्रति पंक्ति एक पथ)
  --no-discovery                    katana+httpx छोड़ें, केवल रूट पथ टेस्ट करें
  --output                string    रिपोर्ट आउटपुट निर्देशिका (डिफ़ॉल्ट ./outputs)
  --format                string    रिपोर्ट फॉर्मेट: json, html, both (डिफ़ॉल्ट both)
  --katana-path           string    कस्टम katana बाइनरी पथ
  --httpx-path            string    कस्टम httpx बाइनरी पथ
  --insecure                        TLS प्रमाणपत्र सत्यापन छोड़ें
  --header, -H            string    कस्टम HTTP हेडर (दोहराने योग्य): -H 'Authorization: Bearer TOKEN'
  --duration              duration  प्रति चरण अवधि (जैसे 30s, 1m)। req-per-worker को ओवरराइड करता है
  --alert-p99             float     किसी चरण में P99 लेटेंसी Xms से अधिक होने पर exit 1
  --alert-error-rate      float     किसी चरण में त्रुटि दर X%% से अधिक होने पर exit 1
  --alert-rps             float     किसी चरण में req/s X से कम होने पर exit 1
```

### Udaharan

```bash
stormprobe --insecure --concurrency-spike 500 --req-per-worker 20 https://example.com

stormprobe --insecure --endpoints endpoints.txt https://example.com

stormprobe --insecure --format json --output ./results https://example.com

stormprobe --insecure -H "Authorization: Bearer TOKEN" -H "X-Tenant: acme" https://api.example.com

stormprobe --insecure --duration 30s --no-discovery https://example.com

stormprobe --insecure --duration 1m --concurrency-spike 300 https://example.com

stormprobe --insecure --alert-p99 500 --alert-error-rate 5 https://example.com
```

## Parikshan charan

| Charan | Vivaran |
|---|---|
| **0 -- Discovery** | Katana crawl + Httpx probe, deduplicated endpoint suchi |
| **1 -- Ramp-Up** | Kramik vriddhi: 5, 15, 30, 50 virtual users (aabhaasi upayogkarta) |
| **2 -- Sustained** | Lakshya concurrency (samanvayta) par 3 lahar, girават ko maapna |
| **3 -- Spike** | Shikhar par achanak visfot, phir thandak |
| **4 -- Recovery** | Spike ke baad 10s thandak ke baad svasthya jaanch |

## Parinaam

| File | Vivaran |
|---|---|
| `stormprobe_report_<timestamp>.json` | Sabhi maapdandon sahit machine padhaniy parinaam |
| `stormprobe_report_<timestamp>.html` | Drishya dashboard, kisi bhi browser mein kholein |

## Pariyojana sanrachana

```
stormprobe/
    cmd/main.go                    CLI pravesh bindu
    internal/
        config/config.go           Sajha prakar, charan jananak, alert config
        alert/
            alert.go               Threshold checks, CI/CD exit code logic
        metrics/
            latency.go             Pratishat ganana
            errors.go              Truti vargikaran
        discovery/
            discovery.go           Discover() saarvajanik interface
            katana.go              Katana crawler ekikaran
            httpx.go               Httpx probe ekikaran
        runner/
            phase.go               Charan orchestration
            worker.go              Goroutine pool, prati worker RNG
        report/
            json.go                JSON report lekhak
            html.go                Svatantra HTML report
    Dockerfile                     Katana + Httpx sahit bahu-charan
    .goreleaser.yml                Cross-compilation vinya
    Makefile                       Build, test, lint lakshya
    config.example.yml             Charan default maan sandarbh
```

## Aavashyakataen

- Go 1.21+
- [Katana](https://github.com/projectdiscovery/katana) (ProjectDiscovery) -- vaikalpik, svachalit khoj ke liye
- [Httpx](https://github.com/projectdiscovery/httpx) (ProjectDiscovery) -- vaikalpik, endpoint satyapan ke liye

```bash
bash scripts/install-tools.sh
```

## Dayitva apvarjan

Yah upakaran **keval adhikrit suraksha parikshan aur pradarshan mulyankan** ke liye hai. Koi bhi load test chalane se pahle aapko lakshya pranali ke maalik se spasht likhit anumati prapt karni hogi. Aapke svamitv mein na hone vale ya jinke liye aapke paas parikshan ki anumati nahin hai un pranaliyon ke khilaf is upakaran ka anadhikrit upayog sthaniy, rashtriy ya antarrashtriy kanonon ka ullanghan kar sakta hai. Lekhak is upakaran ke durupayog ya iske karan hone vale nuksan ki koi jimmedari nahin lete.

## Laisens

MIT

## Lekhak

**Umut ÖZEN** — [@umutozen](https://github.com/umutozen)

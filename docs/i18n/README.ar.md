<!-- rtl -->
[English](../../README.md) · [Türkçe](README.tr.md) · [Deutsch](README.de.md) · [Francais](README.fr.md) · [Espanol](README.es.md) · [Italiano](README.it.md) · [Portugues](README.pt-BR.md) · [Russkij](README.ru.md) · [Zhongwen Jian](README.zh-CN.md) · [Zhongwen Fan](README.zh-TW.md) · [Nihongo](README.ja.md) · [Hangugeo](README.ko.md) · [Al-Arabiyya](README.ar.md) · [Hindi](README.hi.md) · [Nederlands](README.nl.md) · [Polski](README.pl.md) · [Ukrainska](README.uk.md)

# StormProbe

<p align="center">
  <img src="../../assets/banner.png" alt="StormProbe" width="100%">
</p>

**iktashif. hammil. hallil.**

Adaat ikhtibaaraat tahmiil HTTP bi mustawaa al-intaaj maa iktishaaf endpoint (nuqat al-nihaaya) dhaatiyy. Sifr tabaiyyaat.

[![CI](https://github.com/umutozen/stormprobe/actions/workflows/ci.yml/badge.svg)](https://github.com/umutozen/stormprobe/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/umutozen/stormprobe)](https://goreportcard.com/report/github.com/umutozen/stormprobe)

---

## Al-miizaat

- **Al-iktishaaf al-tiliqa'iyy** -- zahf Katana + misbaar Httpx yajid endpoint naashita tiliqa'iyyan
- **Ikhtibaaraat mutaaddidat al-maraahil** -- Ramp-Up (ziyaadat tadriijiiyya), Sustained (mustadiim), Spike (dhurwat al-himl), Recovery (istiirdaad)
- **Maqaayiis tafsiiliiyya** -- P50 / P95 / P99 latency (ta'khhur), req/s, tasniif al-akhta'
- **Taqriir muzdawaj** -- JSON (qiraa'a aaliyya) + HTML (lawhat tihakkum basariyya bi sima daakina)
- **Sifr tabaiyyaat** -- Maktabat Go al-qiyaasiyya al-naqqiyya, laa wahadaat min atraf thaalitha
- **Jaahiz lil Docker** -- Amr waahid maa Katana + Httpx mudmajaan
- **Mutaaddid al-manassaat** -- Malaffaat thunaa'iyya Linux, macOS, Windows abar GoReleaser

## Bidaaya sarii'a

```bash
go run ./cmd https://example.com

go build -o stormprobe ./cmd
./stormprobe --insecure https://example.com
```

### Docker

```bash
# اختبار سريع، لا يتم حفظ أي تقرير
docker run --rm ghcr.io/umutozen/stormprobe:latest --insecure https://example.com

# اختبار كامل مع حفظ التقارير في ./outputs (Linux / macOS)
docker run --rm -v $(pwd)/outputs:/app/outputs ghcr.io/umutozen/stormprobe:latest --insecure --format both https://example.com

# اختبار كامل مع حفظ التقارير في ./outputs (Windows PowerShell)
docker run --rm -v "${PWD}\outputs:/app/outputs" ghcr.io/umutozen/stormprobe:latest --insecure --format both https://example.com

# اختبار ذروة عالي الحمل (500 مستخدم متزامن)
docker run --rm ghcr.io/umutozen/stormprobe:latest --insecure --concurrency-spike 500 https://example.com

# تخطي الاكتشاف، اختبار المسار الجذري فقط
docker run --rm ghcr.io/umutozen/stormprobe:latest --insecure --no-discovery https://example.com

# رؤوس HTTP مخصصة (التفويض، مستأجر مخصص، إلخ)
docker run --rm ghcr.io/umutozen/stormprobe:latest --insecure -H "Authorization: Bearer TOKEN" -H "X-Tenant: acme" https://api.example.com
```

## Al-istikhdam

```
stormprobe [flags] <رابط-الهدف>

Flags:
  --concurrency-ramp      int       الحد الأقصى للتزامن في مرحلة الصعود (الافتراضي 50)
  --concurrency-sustained int       التزامن في المرحلة المستدامة (الافتراضي 75)
  --concurrency-spike     int       الحد الأقصى للتزامن في مرحلة الذروة (الافتراضي 250)
  --req-per-worker        int       الطلبات لكل عامل لكل خطوة (الافتراضي 15)
  --timeout               duration  مهلة انتظار كل طلب (الافتراضي 10s)
  --endpoints             string    ملف قائمة نقاط النهاية (مسار واحد لكل سطر)
  --no-discovery                    تجاوز katana+httpx، اختبار المسار الجذري فقط
  --output                string    دليل إخراج التقارير (الافتراضي ./outputs)
  --format                string    تنسيق التقارير: json أو html أو both (الافتراضي both)
  --katana-path           string    مسار ثنائي katana مخصص
  --httpx-path            string    مسار ثنائي httpx مخصص
  --insecure                        تخطي التحقق من شهادة TLS
  --header, -H            string    رأس HTTP مخصص (قابل للتكرار): -H 'Authorization: Bearer TOKEN'
  --duration              duration  مدة كل مرحلة (مثل 30s أو 1m). يتجاوز req-per-worker
  --alert-p99             float     exit 1 إذا تجاوز P99 قيمة Xms
  --alert-error-rate      float     exit 1 إذا تجاوزت نسبة الخطأ X%%
  --alert-rps             float     exit 1 إذا انخفض req/s عن X
```

### Amthila

```bash
stormprobe --insecure --concurrency-spike 500 --req-per-worker 20 https://example.com

stormprobe --insecure --endpoints endpoints.txt https://example.com

stormprobe --insecure --format json --output ./results https://example.com

stormprobe --insecure -H "Authorization: Bearer TOKEN" -H "X-Tenant: acme" https://api.example.com

stormprobe --insecure --duration 30s --no-discovery https://example.com

stormprobe --insecure --duration 1m --concurrency-spike 300 https://example.com

stormprobe --insecure --alert-p99 500 --alert-error-rate 5 https://example.com
```

## Maraahil al-ikhtibaaraat

| Al-marhala | Al-wasf |
|---|---|
| **0 -- Discovery** | Zahf Katana + misbaar Httpx, qaa'imat endpoint muzaalat al-takraar |
| **1 -- Ramp-Up** | Ziyaada tadriijiiyya: 5, 15, 30, 50 virtual users (mustakhdimoon iftiraadiyyoon) |
| **2 -- Sustained** | 3 mawjaat inda concurrency (tazaamun) al-hadaf, yaqiis al-tadahhur |
| **3 -- Spike** | Infijaar mufaaji' ilaa al-dhurwa thumma al-tabriid |
| **4 -- Recovery** | Fahs sihhi maabadal-spike baad 10 thawaanin min al-tabriid |

## Al-mukhrajaat

| Al-malaf | Al-wasf |
|---|---|
| `stormprobe_report_<timestamp>.json` | Nataa'ij qiraa'a aaliyya maa jamiia al-maqaayiis |
| `stormprobe_report_<timestamp>.html` | Lawhat tihakkum basariyya, iftah fii ayy mutasaffih |

## Haykal al-mashruu

```
stormprobe/
    cmd/main.go                    nuqtat dukhul CLI
    internal/
        config/config.go           anwaa mushtaraka, muwallidaat al-maraahil, alert config
        alert/
            alert.go               Threshold checks, CI/CD exit code logic
        metrics/
            latency.go             hisaab al-nisab al-mi'awiyya
            errors.go              tasniif al-akhta'
        discovery/
            discovery.go           waajhat Discover() al-aamma
            katana.go              takaamul zaahif Katana
            httpx.go               takaamul misbaar Httpx
        runner/
            phase.go               tansiiq al-maraahil
            worker.go              tajammu goroutine, RNG li-kull worker
        report/
            json.go                kaatib taqriir JSON
            html.go                taqriir HTML mustaqill
    Dockerfile                     mutaaddid al-maraahil maa Katana + Httpx
    .goreleaser.yml                iidaad al-tajmii al-mutaqaatia
    Makefile                       ahdaaf build, test, lint
    config.example.yml             marjia qiyam iftiraadiyya lil-maraahil
```

## Al-mutatalabaat

- Go 1.21+
- [Katana](https://github.com/projectdiscovery/katana) (ProjectDiscovery) -- ikhtiyaariyy, lil-iktishaaf al-tiliqa'iyy
- [Httpx](https://github.com/projectdiscovery/httpx) (ProjectDiscovery) -- ikhtiyaariyy, li-tahaqquq min endpoint

```bash
bash scripts/install-tools.sh
```

## Ikhla' mas'uliyya

Hadhihi al-adaat mukhassasa **li-ikhtibaaraat al-amn al-mu'adhdhana wa taqyiim al-ada' faqat**. Yajib alayk al-husul ala idhn khattiy sarih min malik al-nizam al-mustadaf qabl ijra' ayy ikhtibaaraat tahmiil. Qad yu'addi al-istikhdam ghayr al-musarrah li-hadhihi al-adaat didd anzhima la tamlikuha aw laysa ladayk idhn li-ikhtibariha ila mukhaalafat al-qawaniin al-mahaliyya aw al-wataniyya aw al-dawliyya. La yatahammal al-mu'allifun ayy mas'uliyya an su' al-istikhdam aw al-adrar al-naatija an hadhihi al-adaat.

## Al-tarkhiis

MIT

## Al-mu'allif

**Umut ÖZEN** -- [@umutozen](https://github.com/umutozen)

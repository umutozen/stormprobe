# StormProbe — Geliştirme Kuralları

## Yeni Özellik Eklendiğinde Zorunlu Adımlar

Her yeni özellik veya flag eklendiğinde aşağıdaki dosyalar otomatik olarak güncellenmeli:

### 1. Kod Değişikliği
- İlgili Go paketi güncellenir
- `go build ./...` ve `go test ./...` geçmeli

### 2. Dokümantasyon (ZORUNLU)
- **README.md** — Usage/Flags tablosuna ve Examples bölümüne eklenir
- **docs/i18n/*.md** — 16 dil dosyasının tamamına aynı ekleme yapılır (Go script ile)
- **RELEASE_NOTES.md** — CLI Flags tablosu ve Quick Usage güncellenir
- **CHANGELOG.md** — Unreleased bölümüne özellik eklenir

### 3. Proje Yapısı Değişirse
- **README.md** ve **docs/i18n/*.md** — Project Structure bölümü güncellenir

### 4. Git
- `git add` → `git commit` → `git push origin main`

---

## Çalıştırma (Windows PowerShell)

```powershell
go build -o stormprobe.exe ./cmd
.\stormprobe.exe --no-discovery https://example.com
```

> `stormprobe` komutu PATH'te yoksa `.\stormprobe.exe` kullan.

---

## Mevcut Özellikler

| Flag | Açıklama |
|------|----------|
| `--duration` | Her fazı X süre boyunca çalıştırır (ör. `30s`, `1m`) |
| `--alert-p99` | P99 latency eşiği aşılırsa exit 1 |
| `--alert-error-rate` | Hata oranı eşiği aşılırsa exit 1 |
| `--alert-rps` | req/s eşiğin altına düşerse exit 1 |
| `-H` / `--header` | Özel HTTP header (tekrarlanabilir) |
| `--insecure` | TLS doğrulamasını atla |
| `--no-discovery` | Katana+Httpx'i atla, sadece `/` test et |
| `--endpoints` | Özel endpoint listesi dosyası |
| `--format` | Rapor formatı: json, html, both |
| `--concurrency-ramp` | Ramp-Up peak concurrency |
| `--concurrency-sustained` | Sustained concurrency |
| `--concurrency-spike` | Spike peak concurrency |
| `--req-per-worker` | Worker başına istek sayısı |
| `--timeout` | İstek başına timeout |

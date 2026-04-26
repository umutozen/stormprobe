[English](../../README.md) · [Türkçe](README.tr.md) · [Deutsch](README.de.md) · [Francais](README.fr.md) · [Espanol](README.es.md) · [Italiano](README.it.md) · [Portugues](README.pt-BR.md) · [Russkij](README.ru.md) · [Zhongwen Jian](README.zh-CN.md) · [Zhongwen Fan](README.zh-TW.md) · [Nihongo](README.ja.md) · [Hangugeo](README.ko.md) · [Al-Arabiyya](README.ar.md) · [Hindi](README.hi.md) · [Nederlands](README.nl.md) · [Polski](README.pl.md) · [Ukrainska](README.uk.md)

# StormProbe

<p align="center">
  <img src="../../assets/banner.png" alt="StormProbe" width="100%">
</p>

**decouvrir. charger. analyser.**

Testeur de charge HTTP de qualite production avec decouverte autonome d'endpoints (points de terminaison). Zero dependance.

[![CI](https://github.com/umutozen/stormprobe/actions/workflows/ci.yml/badge.svg)](https://github.com/umutozen/stormprobe/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/umutozen/stormprobe)](https://goreportcard.com/report/github.com/umutozen/stormprobe)

---

## Fonctionnalites

- **Decouverte automatique** -- Exploration Katana + sonde Httpx trouve automatiquement les endpoints actifs
- **Tests multi-phases** -- Ramp-Up (montee en charge), Sustained (soutenu), Spike (pic de charge), Recovery (recuperation)
- **Metriques detaillees** -- P50 / P95 / P99 latency (latence), req/s, classification des erreurs
- **Double rapport** -- JSON (lisible par machine) + HTML (tableau de bord visuel theme sombre)
- **Zero dependance** -- Bibliotheque standard Go pure, aucun module tiers
- **Pret pour Docker** -- Commande unique avec Katana + Httpx integres
- **Multi-plateforme** -- Binaires Linux, macOS, Windows via GoReleaser

## Demarrage rapide

```bash
go run ./cmd https://example.com

go build -o stormprobe ./cmd
./stormprobe --insecure https://example.com
```

### Docker

```bash
# Quick test, no report saved
docker run --rm ghcr.io/umutozen/stormprobe:latest --insecure https://example.com

# Full test with reports saved to ./outputs (Linux / macOS)
docker run --rm -v $(pwd)/outputs:/app/outputs ghcr.io/umutozen/stormprobe:latest --insecure --format both https://example.com

# Full test with reports saved to ./outputs (Windows PowerShell)
docker run --rm -v "${PWD}\outputs:/app/outputs" ghcr.io/umutozen/stormprobe:latest --insecure --format both https://example.com

# High-load spike test (500 concurrent users)
docker run --rm ghcr.io/umutozen/stormprobe:latest --insecure --concurrency-spike 500 https://example.com

# Skip discovery, test root path only
docker run --rm ghcr.io/umutozen/stormprobe:latest --insecure --no-discovery https://example.com

# Custom HTTP headers (Authorization, custom tenant, etc.)
docker run --rm ghcr.io/umutozen/stormprobe:latest --insecure -H "Authorization: Bearer TOKEN" -H "X-Tenant: acme" https://api.example.com
```

## Utilisation

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
```

### Exemples

```bash
stormprobe --insecure --concurrency-spike 500 --req-per-worker 20 https://example.com

stormprobe --insecure --endpoints endpoints.txt https://example.com

stormprobe --insecure --format json --output ./results https://example.com

stormprobe --insecure -H "Authorization: Bearer TOKEN" -H "X-Tenant: acme" https://api.example.com
```

## Phases de test

| Phase | Description |
|---|---|
| **0 -- Discovery** | Exploration Katana + sonde Httpx, liste d'endpoints dedupliquee |
| **1 -- Ramp-Up** | Augmentation progressive : 5, 15, 30, 50 virtual users (utilisateurs virtuels) |
| **2 -- Sustained** | 3 vagues a la concurrency (simultaneite) cible, mesure la degradation |
| **3 -- Spike** | Rafale soudaine vers le pic, puis refroidissement |
| **4 -- Recovery** | Controle de sante post-spike apres 10s de refroidissement |

## Sortie

| Fichier | Description |
|---|---|
| `stormprobe_report_<timestamp>.json` | Resultats lisibles par machine avec toutes les metriques |
| `stormprobe_report_<timestamp>.html` | Tableau de bord visuel, ouvrir dans n'importe quel navigateur |

## Structure du projet

```
stormprobe/
    cmd/main.go                    Point d'entree CLI
    internal/
        config/config.go           Types partages, generateurs de phases
        metrics/
            latency.go             Calcul de percentiles
            errors.go              Classification des erreurs
        discovery/
            discovery.go           Interface publique Discover()
            katana.go              Integration du crawler Katana
            httpx.go               Integration de la sonde Httpx
        runner/
            phase.go               Orchestration des phases
            worker.go              Pool de goroutines, RNG par worker
        report/
            json.go                Redacteur de rapport JSON
            html.go                Rapport HTML autonome
    Dockerfile                     Multi-etapes avec Katana + Httpx
    .goreleaser.yml                Configuration de compilation croisee
    Makefile                       Cibles build, test, lint
    config.example.yml             Reference des valeurs par defaut des phases
```

## Prerequis

- Go 1.21+
- [Katana](https://github.com/projectdiscovery/katana) (ProjectDiscovery) -- optionnel, pour la decouverte automatique
- [Httpx](https://github.com/projectdiscovery/httpx) (ProjectDiscovery) -- optionnel, pour la validation d'endpoints

```bash
bash scripts/install-tools.sh
```

## Avertissement

Cet outil est destine uniquement aux **tests de securite autorises et a l'evaluation des performances**. Vous devez obtenir une autorisation ecrite explicite du proprietaire du systeme cible avant d'executer des tests de charge. L'utilisation non autorisee de cet outil contre des systemes que vous ne possedez pas ou pour lesquels vous n'avez pas d'autorisation de test peut enfreindre les lois locales, nationales ou internationales. Les auteurs n'assument aucune responsabilite pour une mauvaise utilisation ou les dommages causes par cet outil.

## Licence

MIT

## Auteur

**Umut ÖZEN** -- [@umutozen](https://github.com/umutozen)

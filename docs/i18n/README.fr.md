[English](../../README.md) · [Türkçe](README.tr.md) · [Deutsch](README.de.md) · [Francais](README.fr.md) · [Espanol](README.es.md) · [Italiano](README.it.md) · [Portugues](README.pt-BR.md) · [Russkij](README.ru.md) · [Zhongwen Jian](README.zh-CN.md) · [Zhongwen Fan](README.zh-TW.md) · [Nihongo](README.ja.md) · [Hangugeo](README.ko.md) · [Al-Arabiyya](README.ar.md) · [Hindi](README.hi.md) · [Nederlands](README.nl.md) · [Polski](README.pl.md) · [Ukrainska](README.uk.md)

<p align="center">
  <img src="../../assets/banner.png" alt="StormProbe Banner" width="100%">
</p>

<h1 align="center">StormProbe</h1>

<p align="center">
  <strong>decouvrir · charger · analyser</strong><br>
  Outil de test de charge HTTP de qualité production avec découverte autonome des points de terminaison. Zéro dépendance externe.
</p>

<p align="center">
  <a href="https://github.com/umutozen/stormprobe/actions/workflows/ci.yml">
    <img src="https://github.com/umutozen/stormprobe/actions/workflows/ci.yml/badge.svg" alt="CI">
  </a>
  <a href="https://goreportcard.com/report/github.com/umutozen/stormprobe">
    <img src="https://goreportcard.com/badge/github.com/umutozen/stormprobe" alt="Go Report Card">
  </a>
  <a href="https://github.com/umutozen/stormprobe/releases">
    <img src="https://img.shields.io/github/v/release/umutozen/stormprobe" alt="Dernière version">
  </a>
  <a href="../../LICENSE">
    <img src="https://img.shields.io/badge/license-MIT-blue.svg" alt="Licence MIT">
  </a>
</p>

---

## Qu'est-ce que StormProbe ?

StormProbe est un outil de test de charge HTTP **sans dépendance** écrit en Go pur. Il découvre automatiquement les points de terminaison de votre application avec [Katana](https://github.com/projectdiscovery/katana) et [Httpx](https://github.com/projectdiscovery/httpx), puis exécute un test de charge structuré en 4 phases — **Ramp-Up → Sustained → Spike → Recovery** — et produit des rapports JSON et HTML avec les métriques de latence P50/P95/P99.

Conçu pour :
- **Les pipelines CI/CD** — seuils d'alerte avec codes de sortie
- **Les testeurs de pénétration** — découverte automatique des points de terminaison actifs
- **Les ingénieurs DevOps** — benchmarks avant et après déploiement
- **Les équipes QA** — validation des SLA de performance

---

## Fonctionnalités

| Fonctionnalité | Détails |
|---|---|
| **Découverte automatique** | Crawl Katana + sonde Httpx, liste de points de terminaison dédupliquée |
| **Test en 4 phases** | Ramp-Up → Sustained → Spike → Recovery |
| **Métriques riches** | Latence P50/P95/P99, req/s, classification des erreurs par phase |
| **Mode durée** | Phases basées sur le temps (`--duration 30s`) au lieu du nombre de requêtes |
| **Seuils d'alerte** | `exit 1` compatible CI/CD sur violation P99, taux d'erreur ou RPS |
| **Rapports doubles** | JSON (lisible par machine) + tableau de bord HTML sombre autonome |
| **En-têtes personnalisés** | Tokens Bearer, en-têtes tenant, cookies — propagés partout |
| **Prêt pour Docker** | Image unique avec Katana + Httpx intégrés |
| **Multi-plateforme** | Linux, macOS, Windows — amd64 & arm64 |
| **Zéro dépendance** | Go 1.21+ stdlib pur, `go install` suffit |

---

## Démarrage rapide

```bash
# Installation (Go 1.21+ requis)
go install github.com/umutozen/stormprobe/cmd/stormprobe@latest

# Test de base — découverte automatique, 4 phases
stormprobe https://example.com

# Ignorer la découverte, tester uniquement le chemin racine
stormprobe --no-discovery https://example.com

# TLS + en-tête auth + rapport JSON uniquement
stormprobe --insecure -H "Authorization: Bearer TOKEN" --format json https://api.example.com

# Test basé sur la durée : 30 secondes par phase
stormprobe --insecure --duration 30s --no-discovery https://example.com

# CI/CD : exit 1 si P99 > 500ms ou erreurs > 5%
stormprobe --alert-p99 500 --alert-error-rate 5 https://example.com
```

---

## Installation

### Option 1 — `go install` *(recommandé)*

```bash
go install github.com/umutozen/stormprobe/cmd/stormprobe@latest
stormprobe --help
```

### Option 2 — Binaire précompilé *(sans Go)*

Télécharger depuis la [page Releases](https://github.com/umutozen/stormprobe/releases).

```bash
# Linux / macOS
chmod +x stormprobe
sudo mv stormprobe /usr/local/bin/

# Windows PowerShell
Move-Item stormprobe.exe C:\Windows\System32\stormprobe.exe
```

### Option 3 — Compiler depuis les sources

```bash
git clone https://github.com/umutozen/stormprobe.git
cd stormprobe
go build -o stormprobe ./cmd/stormprobe      # Linux / macOS
go build -o stormprobe.exe ./cmd/stormprobe  # Windows
```

### Option 4 — Docker

```bash
docker run --rm ghcr.io/umutozen/stormprobe:latest --insecure https://example.com

# Linux / macOS
docker run --rm -v $(pwd)/outputs:/app/outputs \
  ghcr.io/umutozen/stormprobe:latest --insecure --format both https://example.com

# Windows PowerShell
docker run --rm -v "${PWD}\outputs:/app/outputs" `
  ghcr.io/umutozen/stormprobe:latest --insecure --format both https://example.com
```

### Optionnel — Outils de découverte

```bash
bash scripts/install-tools.sh
```

> Sans Katana/Httpx, utiliser `--no-discovery` ou `--endpoints`.

---

## Référence CLI

### Concurrence & Charge

| Flag | Défaut | Description |
|---|---|---|
| `--concurrency-ramp` | `50` | Concurrence maximale pour la phase Ramp-Up |
| `--concurrency-sustained` | `75` | Concurrence pour la phase Sustained |
| `--concurrency-spike` | `250` | Concurrence maximale pour la phase Spike |
| `--req-per-worker` | `15` | Requêtes par worker (ignoré si `--duration` est défini) |
| `--duration` | `0` (désactivé) | Durée de phase basée sur le temps, ex. `30s`, `1m` |
| `--timeout` | `10s` | Délai d'attente par requête |

### Découverte

| Flag | Défaut | Description |
|---|---|---|
| `--no-discovery` | `false` | Ignorer Katana+Httpx, tester uniquement la racine |
| `--endpoints` | `""` | Charger des points de terminaison depuis un fichier |
| `--katana-path` | `""` | Chemin du binaire Katana personnalisé |
| `--httpx-path` | `""` | Chemin du binaire Httpx personnalisé |

### HTTP & Sécurité

| Flag | Défaut | Description |
|---|---|---|
| `-H`, `--header` | — | En-tête personnalisé (répétable) |
| `--insecure` | `false` | Ignorer la vérification du certificat TLS |

### Sortie

| Flag | Défaut | Description |
|---|---|---|
| `--format` | `both` | Format du rapport : `json`, `html`, `both` |
| `--output` | `./outputs` | Répertoire de sortie pour les rapports |

### Alertes CI/CD

| Flag | Défaut | Description |
|---|---|---|
| `--alert-p99` | `0` (désactivé) | exit 1 si P99 (ms) dépasse le seuil |
| `--alert-error-rate` | `0` (désactivé) | exit 1 si le taux d'erreur (%) dépasse le seuil |
| `--alert-rps` | `0` (désactivé) | exit 1 si req/s tombe en dessous du seuil |

---

## Phases de test

| Phase | Description |
|---|---|
| **0 — Discovery** | Crawl Katana + sonde Httpx, liste de points de terminaison dédupliquée |
| **1 — Ramp-Up** | Augmentation progressive de la concurrence : 5 → 15 → 30 → pic |
| **2 — Sustained** | 3 vagues à la concurrence cible, mesure la dégradation en régime permanent |
| **3 — Spike** | Augmentation soudaine jusqu'au pic de concurrence, puis refroidissement |
| **4 — Recovery** | Vérification de santé post-spike après 10s de refroidissement |

---

## Sortie & Rapports

| Fichier | Description |
|---|---|
| `stormprobe_report_<timestamp>.json` | Toutes les métriques par phase, lisible par machine |
| `stormprobe_report_<timestamp>.html` | Tableau de bord sombre autonome avec graphique SVG de latence |

---

## Avertissement

Cet outil est destiné **uniquement aux tests de sécurité autorisés et à l'évaluation des performances**. Vous devez obtenir une autorisation écrite explicite du propriétaire du système cible avant d'exécuter des tests de charge. Les auteurs n'assument **aucune responsabilité** pour tout abus ou dommage causé par cet outil.

---

## Licence

[MIT](../../LICENSE) © [Umut ÖZEN](https://github.com/umutozen)

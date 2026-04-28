# StormProbe

<p align="center">
  <img src="../../assets/banner.png" alt="StormProbe Banner" width="100%">
</p>

<h1 align="center">StormProbe</h1>

<p align="center">
  <strong>crawl · discover · stress · decide</strong><br>
  Test de charge HTTP moderne avec découverte autonome des endpoints et verdicts de production.
</p>

<p align="center">
  <strong>Pas seulement un outil de benchmarking.</strong><br>
  StormProbe ne se contente pas de tester. Il explique si votre système est réellement devenu plus sûr.
</p>

---

## Pourquoi StormProbe ?

Les testeurs de charge traditionnels donnent des chiffres.<br>StormProbe donne des décisions.

Au lieu de sélectionner manuellement des endpoints et de deviner les tableaux bruts, StormProbe détecte automatiquement les endpoints, effectue un stress par étapes et pose un diagnostic.

### Outils classiques vs StormProbe

| Outils Classiques | StormProbe |
|---|---|
| Sélection manuelle des endpoints | Découverte automatique des endpoints |
| Résultats bruts | Verdict final + diagnostic |
| Interprétation manuelle | Recommandations concrètes |
| Un seul profil de test | Montée en charge + Maintenu + Pic + Récupération |
| "Quelque chose rame" | "Voici où ça casse et pourquoi" |

---

## Qu'est-ce qu'il fait ?

StormProbe est un CLI de test de charge écrit en Go pour la production.

- Découverte d'endpoints
- Test de simultanéité
- Analyse des centiles de latence
- Détection de la saturation de la bande passante
- Validation de la récupération
- Seuils CI / CD
- Rapports JSON+HTML
- Vérifications de mise en production

---

## Exemple de sortie

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

## Fonctionnalités

| Fonctionnalité | Détails |
|---|---|
| Découverte Auto | Katana crawl + Httpx probe avec endpoints live dédupliqués |
| Test en 4 Phases | Ramp-Up → Sustained → Spike → Recovery |
| Moteur de Verdict | Concurrence sûre, point de dégradation, classification des goulots |
| Métriques Riches | Latence Avg / P50 / P95 / P99 + req/s + classification des erreurs |
| Analyse du Débit | Détecte la saturation avant les pannes sévères |
| Validation de Récupération | Confirme si le système récupère après des pics de trafic |
| Seuils d'Alerte | Codes de sortie pour CI/CD basés sur la latence, taux d'erreur et RPS |
| Rapports Doubles | JSON + tableau de bord HTML autonome |
| Headers Personnalisés | Jetons d'authentification, headers de tenant, cookies |
| Prêt pour Docker | Image unique avec outils intégrés |
| Multiplateforme | Linux, macOS, Windows — amd64 & arm64 |
| Core 100% Go | Installation rapide, surcharge opérationnelle minimale |

---

## Démarrage Rapide

```bash
# Installer
go install github.com/umutozen/stormprobe/cmd/stormprobe@latest

# Test complet avec découverte automatique
stormprobe https://example.com

# Ignorer la découverte, tester la racine
stormprobe --no-discovery https://example.com

# Test API avec authentification
stormprobe --insecure \
  -H "Authorization: Bearer TOKEN" \
  https://api.example.com

# Validation CI/CD
stormprobe \\
  --alert-p99 500 \\
  --alert-error-rate 5 \\
  https://example.com

# Utiliser un fichier de configuration JSON
stormprobe --config stormprobe.json
```

---

## Flux de Test

```text
Découverte → Montée → Maintenu → Pic → Récupération → Verdict
```

| Phase | Objectif |
|---|---|
| Discovery | Trouver des points d'accès |
| Ramp-Up | Augmentation graduelle |
| Sustained | Validation continue |
| Spike | Test de choc |
| Recovery | Contrôle après le choc |
| Verdict | Analyse finale + goulot |

---

## Installation

### Go Install

```bash
go install github.com/umutozen/stormprobe/cmd/stormprobe@latest
```

Façon recommandée.

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

Cela permet aux régressions de bloquer automatiquement les déploiements.

---

## Cas d'Utilisation

**Pré-Production** — Valider la tolérance aux pics.

**Sécurité + Perfs** — Trouver des points cachés.

**Audit Système** — Mesurer la résistance des serveurs.

**Antirégression** — Éviter la mise en prod si dégradé.

---

## Flags

| Flag | Description |
|---|---|
| `--duration` | Temps par test (ex: 30s) |
| `--alert-p99` | Erreur 1 si P99 dépasse la limite |
| `--alert-error-rate` | Erreur 1 si le taux d'erreur dépasse |
| `--alert-rps` | Erreur 1 si trop peu de requêtes |
| `-H` / `--header` | Header (ex: Auth) |
| `--insecure` | Skip TLS |
| `--no-discovery` | Tester la racine uniquement |
| `--endpoints` | Importer une liste |
| `--format` | Format du rapport |
| `--concurrency-ramp` | Concurrency de base |
| `--concurrency-sustained` | Concurrency maintenue |
| `--concurrency-spike` | Pic concurrency |
| `--req-per-worker` | Requêtes par Worker |
| `--timeout` | Timeout |
| `--config` | Charger les paramètres depuis un fichier de configuration JSON |
| `--auto-profile` | Détection automatique du stack serveur et application du profil de diagnostic (par défaut: activé) |
| `--profile` | Profil manuel: iis, tomcat, php |
| `--no-fingerprint` | Ignorer la détection du stack serveur |
| `--show-profile` | Afficher le stack détecté et le profil, puis quitter |

---

## Prérequis

| Outil | Note |
|---|---|
| Go 1.21+ | Requis pour go install |
| Katana | Optionnel pour endpoint |
| Httpx | Optionnel pour endpoint |
| Docker | Alternative docker |

---

## Clause de non-responsabilité

Outil pour des tests autorisés uniquement.

Demandez toujours l'autorisation.

Usage interdit sinon.

---

## Contribuer

Lisez le document CONTRIBUTING.md.

## Sécurité

Voir SECURITY.md.

## Licence

MIT License © Umut ÖZEN

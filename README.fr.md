# Freecode

**Assistant de codage IA unifié et multiplateforme.**

Freecode est un CLI basé sur Go qui combine le meilleur d'opencode avec des agents améliorés, des hooks, une gestion de fleet et des fonctionnalités de flux de travail—le tout en un seul produit cohérent.

## Fonctionnalités

- **11 Agents Intégrés** — Sisyphus (orchestrateur), Hephaestus (génération de code), Oracle, Librarian, Explore, Prometheus, Metis, Momus, Atlas, Multimodal-Looker, Sisyphus-Junior
- **60+ Hooks de Cycle de Vie** — Session, tool, transform, continuation, Ralph, skill hooks
- **8 Catégories de Tâches** — visual-engineering, ultrabrain, deep, artistry, quick, writing, et plus
- **Onglets de Session** — Plusieurs sessions concurrentes avec vues divisées
- **Gestion de Fleet** — Modes head/agent/client pour la coordination multi-instances
- **MCPs Intégrés** — Recherche web Exa, docs Context7, Grep.app, CLI GitHub/GitLab
- **TUI avec Support Souris** — Interface interactive complète avec éléments cliquables
- **Sécurité Avant Tout** — Tous les services liés à localhost uniquement, pas de télémétrie

## Support Plateforme

| Plateforme | État | Notes |
|-----------|------|-------|
| FreeBSD | Primaire | Go 1.25 dans les ports |
| Linux | Supporté | Empaquetage Flatpak |
| macOS | Supporté | Homebrew |
| IllumOS | Supporté | tarball |
| Windows | ❌ NON supporté | |

## Démarrage Rapide

```bash
# Construire
go build -o freecode ./cmd/freecode

# Exécuter
./freecode

# Ou installer via Homebrew (macOS)
brew install freecode
```

## Architecture

- **CLI Basé sur Go** — Binaire statique unique, sans dépendances runtime
- **Cobra CLI Framework** — Modèles CLI Go standard
- **Bubble Tea TUI** — UI de terminal composable
- **SQLite** — Stockage persistant intégré
- **chi router** — API HTTP légère

## Répertoires Clés

```
cmd/freecode/          # Point d'entrée CLI
cmd/freecode-server/  # Point d'entrée mode serveur
internal/cli/          # Commandes Cobra
internal/agent/        # 11 agents intégrés
internal/hook/         # 60+ hooks de cycle de vie
internal/session/      # Gestion de sessions, onglets
internal/ui/           # TUI Bubble Tea
internal/fleet/        # Fleet head/agent/client
internal/platform/     # Code spécifique à la plateforme
```

## Ports (Localhost Uniquement)

| Service | Port |
|---------|------|
| Serveur API | 18792 |
| Serveur MCP | 18793 |
| Web UI | 18791 |
| Fleet Head | 7842 |

## Comparaison avec opencode

Freecode est une conversion Go d'opencode avec toutes les fonctionnalités oh-my-openagent intégrées nativement:

| Fonctionnalité | opencode | freecode |
|----------------|----------|----------|
| Langage | TypeScript | Go |
| Distribution | NPM | Binaire statique |
| Agents | 7 | 11 |
| Hooks | ~20 | 60+ |
| Mode Fleet | ❌ | ✅ |
| MCPs Intégrés | ❌ | ✅ |

## Sécurité

- **Localhost uniquement** — Tous les services liés à 127.0.0.1 et ::1
- **Pas de télémétrie** — Zéro analytique ou suivi
- **Système de permissions** — Permissions d'outils configurables par agent
- **Mode YOLO** — Option pour ignorer les confirmations (désactivé par défaut)

## Documentation

Voir [AGENTS_START_HERE.md](AGENTS_START_HERE.md) pour les conseils aux agents autonomes, ou parcourir le répertoire [.plan/](.plan/) pour les documents de planification détaillés.

## Auteur

Mark LaPointe <mark@cloudbsd.org>

Tous les commits sont effectués par Mark LaPointe. Pas de co-auteurs, pas de sponsorisations.

## Licence

Licence illimitée accordée au projet.

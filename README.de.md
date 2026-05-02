# Freecode

**Einheitlicher, plattformunabhängiger KI-Programmierassistent.**

Freecode ist ein Go-basiertes CLI, das das Beste von opencode mit verbesserten Agents, Hooks, Fleet-Management und Workflow-Funktionen kombiniert—alles als ein einziges kohärentes Produkt.

## Funktionen

- **11 Integrierte Agents** — Sisyphus (Orchestrator), Hephaestus (Code-Generierung), Oracle, Librarian, Explore, Prometheus, Metis, Momus, Atlas, Multimodal-Looker, Sisyphus-Junior
- **60+ Lebenszyklus-Hooks** — Session, tool, transform, continuation, Ralph, skill hooks
- **8 Aufgabenkategorien** — visual-engineering, ultrabrain, deep, artistry, quick, writing, und mehr
- **Sitzungs-Tabbing** — Mehrere gleichzeitige Sitzungen mit geteilten Ansichten
- **Fleet-Management** — Head/agent/client Modi für Multi-Instanz-Koordination
- **Integrierte MCPs** — Exa-Websuche, Context7-Docs, Grep.app, GitHub/GitLab-CLI
- **TUI mit Maus-Unterstützung** — Vollständige interaktive Oberfläche mit klickbaren Elementen
- **Sicherheit Zuerst** — Alle Dienste nur an Localhost gebunden, keine Telemetrie

## Plattform-Unterstützung

| Plattform | Status | Hinweise |
|-----------|--------|----------|
| FreeBSD 16 | Primär | Go 1.25 in Ports |
| Linux | Unterstützt | Flatpak-Paketierung |
| macOS | Unterstützt | Homebrew |
| IllumOS | Unterstützt | tarball |
| Windows | ❌ NICHT unterstützt | |

## Schnellstart

```bash
# Bauen
go build -o freecode ./cmd/freecode

# Ausführen
./freecode

# Oder via Homebrew installieren (macOS)
brew install freecode
```

## Architektur

- **Go-basiertes CLI** — Einzelne statische Binärdatei, keine Runtime-Abhängigkeiten
- **Cobra CLI Framework** — Standard Go-CLI-Muster
- **Bubble Tea TUI** — Zusammensetzbare Terminal-UI
- **SQLite** — Integrierter persistenter Speicher
- **chi router** — Leichte HTTP-API

## Wichtige Verzeichnisse

```
cmd/freecode/          # CLI-Einstiegspunkt
cmd/freecode-server/  # Server-Modus-Einstiegspunkt
internal/cli/          # Cobra-Befehle
internal/agent/        # 11 integrierte Agents
internal/hook/         # 60+ Lebenszyklus-Hooks
internal/session/      # Sitzungsverwaltung, Tabs
internal/ui/           # Bubble Tea TUI
internal/fleet/        # Fleet head/agent/client
internal/platform/     # Plattformspezifischer Code
```

## Ports (Nur Localhost)

| Dienst | Port |
|--------|------|
| API-Server | 18792 |
| MCP-Server | 18793 |
| Web-UI | 18791 |
| Fleet Head | 7842 |

## Vergleich mit opencode

Freecode ist eine Go-Konvertierung von opencode mit allen oh-my-openagent-Funktionen nativ integriert:

| Funktion | opencode | freecode |
|----------|----------|----------|
| Sprache | TypeScript | Go |
| Verteilung | NPM | Statische Binärdatei |
| Agents | 7 | 11 |
| Hooks | ~20 | 60+ |
| Fleet-Modus | ❌ | ✅ |
| Integrierte MCPs | ❌ | ✅ |

## Sicherheit

- **Nur Localhost** — Alle Dienste an 127.0.0.1 und ::1 gebunden
- **Keine Telemetrie** — Null Analyse oder Tracking
- **Berechtigungssystem** — Konfigurierbare Tool-Berechtigungen pro Agent
- **YOLO-Modus** — Optional zum Überspringen von Bestätigungen (standardmäßig aus)

## Dokumentation

Siehe [AGENTS_START_HERE.md](AGENTS_START_HERE.md) für Anleitungen für autonome Agents, oder durchsuchen Sie das [.plan/](.plan/)-Verzeichnis für detaillierte Planungsdokumente.

## Autor

Mark LaPointe <mark@cloudbsd.org>

Alle Commits werden von Mark LaPointe durchgeführt. Keine Co-Autoren, keine Sponsorings.

## Lizenz

Unbegrenzte Lizenz für das Projekt gewährt.

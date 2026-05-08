# Freecode

**Asistente de codificación AI unificado e independiente de la plataforma.**

Freecode es un CLI basado en Go que combina lo mejor de opencode con agentes mejorados, hooks, gestión de fleet y características de flujo de trabajo—todo como un producto cohesivo.

## Características

- **11 Agentes Incorporados** — Sisyphus (orquestador), Hephaestus (generación de código), Oracle, Librarian, Explore, Prometheus, Metis, Momus, Atlas, Multimodal-Looker, Sisyphus-Junior
- **60+ Hooks de Ciclo de Vida** — Session, tool, transform, continuation, Ralph, skill hooks
- **8 Categorías de Tareas** — visual-engineering, ultrabrain, deep, artistry, quick, writing, y más
- **Pestañas de Sesión** — Múltiples sesiones concurrentes con vistas divididas
- **Gestión de Fleet** — Modos head/agent/client para coordinación multi-instancia
- **MCPs Incorporados** — Búsqueda web Exa, docs Context7, Grep.app, CLI de GitHub/GitLab
- **TUI con Soporte de Ratón** — Interfaz interactiva completa con elementos clicables
- **Seguridad Primero** — Todos los servicios vinculados a localhost únicamente, sin telemetría

## Soporte de Plataforma

| Plataforma | Estado | Notas |
|-----------|--------|-------|
| FreeBSD | Primario | Go 1.25 en ports |
| Linux | Soportado | Empaquetado Flatpak |
| macOS | Soportado | Homebrew |
| IllumOS | Soportado | tarball |
| Windows | ❌ NO soportado | |

## Inicio Rápido

```bash
# Construir
go build -o freecode ./cmd/freecode

# Ejecutar
./freecode

# O instalar vía Homebrew (macOS)
brew install freecode
```

## Arquitectura

- **CLI Basado en Go** — Binario estático único, sin dependencias de runtime
- **Cobra CLI Framework** — Patrones estándar de CLI Go
- **Bubble Tea TUI** — UI de terminal componible
- **SQLite** — Almacenamiento persistente embebido
- **chi router** — API HTTP ligera

## Directorios Clave

```
cmd/freecode/          # Punto de entrada CLI
cmd/freecode-server/  # Punto de entrada modo servidor
internal/cli/          # Comandos Cobra
internal/agent/        # 11 agentes incorporados
internal/hook/         # 60+ hooks de ciclo de vida
internal/session/      # Gestión de sesiones, pestañas
internal/ui/           # TUI Bubble Tea
internal/fleet/        # Fleet head/agent/client
internal/platform/     # Código específico de plataforma
```

## Puertos (Solo Localhost)

| Servicio | Puerto |
|----------|--------|
| Servidor API | 18792 |
| Servidor MCP | 18793 |
| Web UI | 18791 |
| Fleet Head | 7842 |

## Comparación con opencode

Freecode es una conversión Go de opencode con todas las características de oh-my-openagent integradas nativamente:

| Característica | opencode | freecode |
|----------------|----------|----------|
| Lenguaje | TypeScript | Go |
| Distribución | NPM | Binario estático |
| Agentes | 7 | 11 |
| Hooks | ~20 | 60+ |
| Modo Fleet | ❌ | ✅ |
| MCPs Incorporados | ❌ | ✅ |

## Seguridad

- **Solo localhost** — Todos los servicios se vinculan a 127.0.0.1 y ::1
- **Sin telemetría** — Cero análisis o seguimiento
- **Sistema de permisos** — Permisos de herramientas configurables por agente
- **Modo YOLO** — Opcional para omitir confirmaciones (desactivado por defecto)

## Documentación

Consulta [AGENTS_START_HERE.md](AGENTS_START_HERE.md) para guía de agentes autónomos, o explora el directorio [.plan/](.plan/) para documentos de planificación detallados.

## Autor

Mark LaPointe <mark@cloudbsd.org>

Todos los commits son realizados por Mark LaPointe. Sin co-autores, sin patrocinios.

## Licencia

Licencia ilimitada otorgada al proyecto.

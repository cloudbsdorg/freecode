# Freecode Skills

Skills are specialized knowledge areas that agents can invoke for focused task execution.

## Available Skills

| Skill | Description | Category |
|-------|-------------|----------|
| [git-master](git-master/SKILL.md) | Expert git operations, history analysis, bisect | Development Tools |
| [playwright](playwright/SKILL.md) | Browser automation, testing, web scraping | Testing |
| [frontend-ui-ux](frontend-ui-ux/SKILL.md) | UI development, accessibility, design systems | Frontend |
| [review-work](review-work/SKILL.md) | Code review, security audit, quality assessment | Quality |
| [ai-slop-remover](ai-slop-remover/SKILL.md) | Detect and fix AI-generated code patterns | Code Quality |
| [search-code](search-code/SKILL.md) | Expert code search with grep, ast-grep, LSP | Development Tools |
| [architect](architect/SKILL.md) | System design, architecture patterns, trade-offs | Architecture |

## Skill Format

Each skill is a directory containing `SKILL.md` with frontmatter:

```markdown
---
name: skill-name
description: One-line description of the skill
---

# Skill Title

Detailed skill content...
```

## Usage

Skills are loaded via the `skill` tool:

```go
task(
    category="visual-engineering",
    load_skills=["frontend-ui-ux", "playwright"],
    prompt="Build a login form with tests..."
)
```

## Discovery

Skills are auto-discovered from:
- `.skills/` directory in project root
- Configured skill paths in `~/.config/freecode/skills/`
- Remote skill repositories

## Author

Mark LaPointe <mark@cloudbsd.org>

All commits made solely by the author. No co-authors, no sponsorships.

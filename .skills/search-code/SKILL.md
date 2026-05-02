---
name: search-code
description: Expert code search using grep, ast-grep, and semantic search patterns across large codebases
---

# Search Code Skill

Find code patterns efficiently using grep, ast-grep, and structural search. This skill excels at locating code across millions of lines.

## When to Use

- Finding function/variable definitions
- Locating usage patterns across codebase
- Searching with regex for complex matches
- AST-based structural search
- Finding dead code
- Cross-reference analysis

## Core Tools

### Grep (Text Search)
```bash
# Basic search
grep -r "functionName" --include="*.go"

# Case-insensitive
grep -ri "todo" .

# With context (3 lines before/after)
grep -B3 -A3 "pattern" file.go

# Only matching parts
grep -o "pattern" file.go

# Count matches
grep -c "pattern" **/*.go
```

### AST-grep (Structural Search)
```bash
# Find all console.log calls
ast-grep search --pattern 'console.log($MSG)' -l ts

# Find React useEffect without dependencies
ast-grep search --pattern 'useEffect(() => { $$$ }, [])' -l tsx

# Replace pattern across files
ast-grep replace --pattern 'oldPattern' --rewrite 'newPattern'
```

### LSP Symbols
```bash
# Find symbol definitions
lsp_symbols --scope workspace --query "functionName"

# Find all references
lsp_find_references --file path.go --line 42
```

## Search Strategies

### Finding Definitions
1. Use LSP goto-definition for exact symbol
2. Grep for `func/class/const Name`
3. AST-grep for structural patterns

### Finding Usages
1. Grep for exact name (word boundary)
2. AST-grep for call patterns
3. LSP find references

### Dead Code Detection
```bash
# Find unt exported functions
grep -r "func [A-Z]" --only-matching

# Find unused exports
ast-grep search --lang go --pattern 'export $FUNC'
```

## Regex Quick Reference

| Pattern | Matches |
|---------|---------|
| `^func ` | Line starting with func |
| `\.go$` | Files ending in .go |
| `(?:foo\|bar)` | foo or bar (non-capture) |
| `\s` | Any whitespace |

## Integration

Use `task(category='unspecified-low', load_skills=['search-code'])` for code search tasks. Be specific about what to find and where.

## Performance Tips

- Use `--include` to limit file types
- Add `--ignore-case` only when needed
- Prefer word boundaries `\b` for symbols
- Use `--max-count` to limit results for large repos

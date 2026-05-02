---
name: git-master
description: Expert git operations including advanced workflows, history analysis, bisect, blame, rebase, and custom hooks
---

# Git Master Skill

Execute precise git operations using built-in tools. This skill handles everything from basic commits to advanced history rewriting.

## When to Use

- Running git commands (commit, push, pull, branch, merge)
- Analyzing git history (log, blame, bisect, reflog)
- Rewriting history (rebase, squash, amend)
- Managing remotes and submodules
- Finding when code was added or who wrote it
- Setting up git hooks

## Core Commands

### Status & Info
```bash
git status           # Working tree status
git log --oneline -20  # Recent commits
git branch -a        # All branches
git remote -v        # Remotes
```

### History Analysis
```bash
git blame <file>           # Line-by-line author
git log -p --follow <file>  # File history
git bisect start           # Find bad commit
git reflog                 # All ref updates
```

### Rewriting History
```bash
git rebase -i HEAD~5      # Interactive rebase
git commit --amend        # Modify last commit
git reset --soft HEAD~1   # Undo last commit
```

### Advanced
```bash
git stash push -m "message"  # Save work
git cherry-pick <commit>      # Apply single commit
git submodule update --init  # Init submodules
```

## Usage Examples

**Find who added a line:**
```
git blame path/to/file | grep "search term"
```

**Bisect for bug:**
```
git bisect start
git bisect bad HEAD
git bisect good <known-good-commit>
git bisect run <test-command>
```

**Clean up commits:**
```
git rebase -i HEAD~3  # Mark commits as squash/fixup
```

## Integration

This skill is automatically available to agents with git access. Use `task(category='quick', load_skills=['git-master'])` for git operations.

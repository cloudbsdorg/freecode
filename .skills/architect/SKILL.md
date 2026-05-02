---
name: architect
description: System design and architecture consultation - evaluate trade-offs, patterns, and long-term implications
---

# Architect Skill

Provide deep architectural analysis and guidance. Help make structural decisions that will scale and remain maintainable over time.

## When to Use

- Designing new systems or services
- Evaluating architectural patterns
- Reviewing tech stack decisions
- Planning migrations
- Assessing scalability concerns
- Making "build vs buy" decisions

## Analysis Framework

### Requirements Clarification
1. Functional requirements - what does it do?
2. Non-functional requirements - constraints (perf, scale, cost)
3. Acceptance criteria - how do we know it's done?

### Trade-off Analysis
```
Option A: [Approach]
  Pros: [Benefits]
  Cons: [Drawbacks]
  Risk: [Potential issues]

Option B: [Approach]
  ...

Recommendation: [Decision with rationale]
```

### Pattern Selection

| Pattern | Use When | Avoid When |
|---------|----------|------------|
| Microservices | Team autonomy needed | Simple CRUD apps |
| Event-driven | Loose coupling | Tight sync required |
| CQRS | Read/write separation | Simple CRUD |
| Saga | Distributed transactions | Simple workflows |

## Architecture Principles

1. **SOLID** - Single responsibility, open/closed, etc.
2. **DRY** - Don't repeat yourself
3. **KISS** - Keep it simple, stupid
4. **YAGNI** - You aren't gonna need it

## Documentation Template

```markdown
# Architecture Decision Record: [Title]

## Status
Accepted | Deprecated | Superseded

## Context
What is the issue? What's the background?

## Decision
What is the change we're making?

## Consequences
### Positive
- Benefit 1
### Negative
- Drawback 1

## Alternatives Considered
1. Option A - [brief description]
2. Option B - [brief description]
```

## Key Questions to Ask

- What problem are we solving?
- Who are the stakeholders?
- What are the success metrics?
- What's the timeline?
- What's the maintenance burden?
- What are the failure modes?

## Integration

Use `task(subagent_type='oracle', load_skills=['architect'])` for architecture consultations. The Oracle agent already has architectural expertise - this skill enhances it.

## Reference Patterns

- Clean Architecture
- Hexagonal Architecture
- Event Sourcing
- CQRS
- Repository Pattern
- Dependency Injection

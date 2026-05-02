---
name: frontend-ui-ux
description: Design and implement user interfaces with attention to UX patterns, accessibility, and visual polish
---

# Frontend UI/UX Skill

Create polished, accessible user interfaces following modern best practices. This skill bridges design intent and implementation.

## When to Use

- Building React/Solid/Vue components
- Implementing responsive layouts
- Adding animations and transitions
- Ensuring accessibility (a11y)
- Working with CSS/Tailwind styles
- Converting design mockups to code

## Core Principles

### Accessibility First
- Use semantic HTML elements
- Add proper ARIA labels
- Ensure keyboard navigation
- Maintain color contrast ratios
- Support screen readers

### Responsive Design
```css
/* Mobile-first approach */
.container { padding: 1rem; }
@media (min-width: 768px) {
  .container { padding: 2rem; }
}
```

### Component Patterns
```tsx
// Composable component structure
<Card>
  <Card.Header>Title</Card.Header>
  <Card.Body>Content</Card.Body>
  <Card.Footer>Actions</Card.Footer>
</Card>
```

## UI Patterns

### Forms
- Clear labels and error messages
- Inline validation feedback
- Disabled states during submission
- Success confirmations

### Navigation
- Consistent placement
- Active state indicators
- Breadcrumbs for deep hierarchy
- Mobile hamburger menus

### Feedback
- Loading spinners for async ops
- Toast notifications for results
- Progress bars for long operations
- Skeleton screens for content loading

## Design Tokens

```css
:root {
  --color-primary: #3b82f6;
  --color-error: #ef4444;
  --spacing-md: 1rem;
  --radius: 0.5rem;
}
```

## Integration

Use `task(category='visual-engineering', load_skills=['frontend-ui-ux'])` for frontend tasks. Provide design mockups or descriptions as input.

## Key Considerations

1. **Consistency** - Reuse design tokens and component patterns
2. **Performance** - Lazy load below-fold content
3. **Progressive enhancement** - Work without JavaScript
4. **Mobile-first** - Design for small screens first

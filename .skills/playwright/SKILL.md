---
name: playwright
description: Browser automation for testing, web scraping, and UI verification using Playwright
---

# Playwright Skill

Automate browsers for testing web applications, scraping content, and verifying UI behavior.

## When to Use

- End-to-end testing of web applications
- Scraping dynamic content from websites
- Taking screenshots of pages
- Verifying UI elements exist
- Filling forms and submitting data
- Testing user interactions (clicks, hovers, keyboard)

## Core Operations

### Navigation
```typescript
await page.goto('https://example.com')
await page.goBack()
await page.goForward()
await page.reload()
```

### Element Interaction
```typescript
await page.click('button#submit')
await page.fill('input[name="email"]', 'test@example.com')
await page.selectOption('select#country', 'US')
await page.check('input#agree')
```

### Assertions
```typescript
await expect(page.locator('h1')).toHaveText('Welcome')
await expect(page.locator('.error')).toBeHidden()
await page.waitForSelector('.loaded')
```

### Screenshot
```typescript
await page.screenshot({ path: 'screenshot.png', fullPage: true })
```

## Usage Pattern

```typescript
task(category='visual-engineering', load_skills=['playwright'], prompt="
Navigate to the login page, fill in credentials, and verify the dashboard loads.
Take a screenshot of the final state.
")
```

## Key Methods

| Method | Purpose |
|--------|---------|
| `click(selector)` | Click element |
| `fill(selector, text)` | Type text |
| `screenshot(options)` | Capture page |
| `waitForSelector(selector)` | Wait for element |
| `evaluate(fn)` | Run JS in page |

## Integration

Load this skill when task involves browser automation. The task prompt should specify exact steps and verification criteria.

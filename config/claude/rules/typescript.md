---
paths:
  - "**/*.ts"
  - "**/*.tsx"
  - "**/package.json"
---

# TypeScript

- `strict: true`. No `any` — use `unknown` and narrow it. No `@ts-ignore` without a
  comment saying what is being ignored and why.
- Parse input at the boundary into a type, then trust the type inward. Don't re-check.
- `type` for unions and shapes, `interface` only when declaration merging is the point.
- Prefer `satisfies` over a type annotation when you want inference and a check.
- ESM only. Named exports; a default export only when a framework demands one.
- No barrel `index.ts` re-exporting a directory — it defeats tree-shaking and hides cycles.
- Errors are thrown values you can discriminate, or a result union. Never a bare string.
- Don't add a runtime dependency for something the platform now does: `fetch`,
  `structuredClone`, `Intl`, `AbortController`, `node:` built-ins.

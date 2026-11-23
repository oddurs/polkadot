# Stack rules

Managed by polkadot. opencode loads this file after `AGENTS.md`, which is the
instruction file that every agent shares (`agents/AGENTS.md` in polkadot).
These rules are the language and framework detail that Claude Code gets from
its path rules.

## TypeScript

- `strict: true`. No `any` — use `unknown` and narrow. No `@ts-ignore` without
  a comment saying what is ignored and why.
- Parse input at the boundary into a type, then trust the type inward.
- `type` for unions and shapes; `interface` only when merging is the point.
- Explicit return types on anything exported.
- ESM, named exports. No barrel files re-exporting a directory.

## React and Next.js

- Server Components by default. `"use client"` only on the leaves that need it.
- Colocate state with what uses it; extract a hook when it is reused, not
  before.
- Server Components fetch directly or through server actions; client
  components use TanStack Query.
- Every fetch has a loading state and an error state. Both are visible in the
  UI, not just handled.
- Tailwind utilities, variants through `cva`/`cn`, mobile-first, theming via
  CSS variables.

## Data and security

- Migrations for every schema change. Parameterised queries, never string
  interpolation. Index what you query often. RLS on in Supabase.
- Validate input on the server, whatever the client already did.
- Secrets come from the environment. Never commit a key — and if one is
  already in the diff, say so rather than committing around it.

## Tests

- Test behaviour, not implementation. Mock at the boundary you own.
- A test with no assertion is not a test.
- No sleeps and no retries. A flaky test is a failing test.

---
name: nextjs-site
description: Scaffold a modern Next.js site for the current project — App Router, TypeScript strict, StyleX with a token-driven design system, testing, linting, CI, and a dev server on a stable non-colliding port. Structure is inferred from what the project actually is; content and visual style are left deliberately plain. Use when a project needs a website, docs site, landing page, or app shell.
argument-hint: [site kind, if you want to override what's inferred]
---

# A Next.js site for this project

Scaffold a complete, correctly configured site whose **structure** comes from what the
project is, and whose **content and appearance** stay deliberately plain — neutral tokens,
placeholder copy, no visual opinions to unpick later. You are building the frame, properly,
not designing anything.

Execute. Do not stop at a plan.

## 0. Work out what site this is

Read the project first — `README.md`, `Cargo.toml`/`package.json`/`pyproject.toml`, the
CLI's `--help`, the public API, the docs directory, `cairn roadmap` if there is one. Then
pick the shape from what you found, because a CLI tool and a SaaS app need different sites
and getting this wrong wastes the whole scaffold:

| The project is | The site is | Routes |
|---|---|---|
| A CLI or developer tool | Docs-first | `/` overview + install, `/docs/[...slug]`, `/reference` (generated from `--help`), `/changelog` |
| A library or SDK | API reference | `/` quickstart, `/docs/[...slug]`, `/api/[...symbol]`, `/examples` |
| An application | Product site + app shell | `/` landing, `/pricing`, `/docs`, `/(app)` behind auth |
| A service or API | Reference + status | `/` overview, `/docs`, `/reference` from OpenAPI, `/status` |
| A person or org | Site | `/`, `/work/[slug]`, `/writing/[slug]`, `/about` |

State which one you chose and why in one line, then build it. If the project genuinely
doesn't fit one, ask — that is the only question worth asking here.

Where the site lives: a `site/` or `www/` directory inside the repo if the project is the
product, a repo of its own if the site is the product. Default to in-repo.

## 1. Create it

```sh
pnpm create next-app@latest site --ts --app --eslint --no-tailwind --no-src-dir --import-alias "@/*"
```

**No Tailwind** — StyleX is the styling layer and two of them is a decision nobody wants to
maintain. App Router, TypeScript, ESLint.

Set `"strict": true` plus `noUncheckedIndexedAccess` and `verbatimModuleSyntax` in
`tsconfig.json`. Pin the Node version in `.node-version`.

## 2. StyleX — use the SWC compiler, not the official Next plugin

**Check this before you install.** `@stylexjs/nextjs-plugin` is the package the docs point
at and it is the wrong choice today: it is Babel-based, and it has been sitting at 0.11.x
while the StyleX core moved on to 0.19. Putting Babel into a Next 16 project also opts you
out of the SWC pipeline and makes every build slower.

Use the actively maintained Rust/SWC implementation instead:

```sh
pnpm add @stylexjs/stylex
pnpm add -D @stylexswc/nextjs-plugin @stylexswc/rs-compiler @stylexswc/postcss-plugin
```

Verify what is current before writing the config — these move:

```sh
npm view @stylexjs/stylex version
npm view @stylexswc/nextjs-plugin version peerDependencies
```

**Turbopack** (the default dev bundler in Next 16) is supported through the `/turbopack`
export, with one catch that will cost an hour if you miss it: the Turbopack loader
*compiles* StyleX but does not *extract* the CSS. PostCSS does the extraction, so both
pieces are required.

```ts
// next.config.ts
import withStylexTurbopack from '@stylexswc/nextjs-plugin/turbopack'

export default withStylexTurbopack({
  rsOptions: { dev: process.env.NODE_ENV === 'development' },
})({
  reactStrictMode: true,
  typedRoutes: true,
})
```

```js
// postcss.config.js
module.exports = {
  plugins: {
    '@stylexswc/postcss-plugin': {
      rsOptions: { dev: process.env.NODE_ENV === 'development' },
    },
    autoprefixer: {},
  },
}
```

On the webpack path instead, import the carrier stylesheet in `app/layout.tsx`:
`import '@stylexswc/webpack-plugin/stylex.css'`.

**Verify extraction works before building anything on top of it.** Style one element,
run the dev server, and confirm in DevTools that the class is atomic and the custom
properties are present. StyleX that compiles without extracting fails silently — the page
renders completely unstyled and the error message is nothing at all.

## 3. The design system comes first

This is the part that decides whether the site is still pleasant in six months, so it goes
in before any page. Follow `/design-system` for the reasoning; the shape here is:

```
site/
  app/
    layout.tsx            html, fonts, theme, skip link
    page.tsx
  design/
    tokens/
      color.stylex.ts     semantic tokens, light + dark via @media
      space.stylex.ts     one scale: 4 8 12 16 24 32 48 64
      type.stylex.ts      sizes paired with line heights, 3 weights
      radius.stylex.ts
      motion.stylex.ts    durations, easings, reduced-motion
      layer.stylex.ts     named z-index steps, never a literal
    primitives/           Box Stack Inline Grid Text — no appearance of their own
    components/           Button Input Field Link Card Dialog — behaviour + tokens only
  content/                MDX, if the site has docs
  lib/
  tests/
```

Two token layers: primitives (`gray-500`) that nothing imports directly, and semantic
tokens (`color.surface`, `color.textMuted`, `color.borderFocus`) that everything uses. Dark
mode is a `@media (prefers-color-scheme: dark)` value on each semantic token, so no
component knows a theme exists.

Keep the palette **neutral**: a grey ramp, one accent, one danger. The point is a blank
slate — a brand arrives later as a `createTheme` override and nothing else changes.

Primitive props accept token values only (`gap="4"`, never `gap="13px"`), so the scale is
enforced by the type system rather than by discipline.

For anything with real interaction behaviour — dialog, menu, tooltip, combobox — take
`react-aria-components` and style it. Do not hand-roll a focus trap.

## 4. Content and typography, plainly

Placeholder copy that says what the section is for, in plain language. No lorem ipsum
(it hides layout problems), no marketing voice, no invented claims about the project.

System font stack by default — `ui-sans-serif, system-ui, sans-serif` and
`ui-monospace, SFMono-Regular, monospace`. If the project has a wordmark or a font already,
use it via `next/font` with `display: swap` and a matched fallback. Otherwise do not pick a
typeface; that is a design decision and this skill does not make design decisions.

For docs: MDX through `@next/mdx`, with `rehype-pretty-code` (Shiki) for syntax
highlighting and `remark-gfm`. Generate the nav from the content tree rather than
maintaining a list by hand.

## 5. The dev server gets its own port

Running many projects means `3000` is always taken, and `next dev` silently walking to
3001 means you never know which one you are looking at. Give the project a port derived
from its name — stable across runs, different from every other project:

```sh
#!/bin/sh
# scripts/port — a stable, collision-resistant dev port for this project
set -eu
name=$(basename "$(cd "$(dirname "$0")/.." && pwd)")
hash=$(printf '%s' "$name" | cksum | cut -d' ' -f1)
port=$((3000 + hash % 5000))
while lsof -i ":$port" >/dev/null 2>&1; do port=$((port + 1)); done
printf '%s' "$port"
```

```json
"scripts": {
  "dev": "next dev --turbopack -p $(./scripts/port)",
  "build": "next build",
  "start": "next start -p $(./scripts/port)"
}
```

The same project gets the same port every time, and two projects effectively never
collide. Print the URL when the server starts.

## 6. Quality gates

- **Lint**: ESLint with `next/core-web-vitals` and `@typescript-eslint`. Warnings are
  errors in CI.
- **Format**: Prettier, or Biome if the repo already uses it. One formatter.
- **Types**: `tsc --noEmit` in the check step. `next build` alone does not type-check
  everything.
- **Tests**: Vitest with `@testing-library/react` for components — query by role and
  label, never by class. Playwright for the handful of end-to-end paths that matter.
  `@axe-core/playwright` on every page, failing the build on violations.
- **`scripts/task`**: wire `fmt fmt:check lint test build check` to the above so CI and
  the git hooks call the same six verbs as every other project in this repo.

## 7. The things that are always forgotten

Do these now; each is ten minutes and none of them get done later.

- `app/layout.tsx`: `lang`, a skip-to-content link, `metadataBase`.
- `generateMetadata` per route — real title and description, Open Graph, Twitter card.
- `app/opengraph-image.tsx` using `next/og`, generated from the title. No design needed.
- `app/icon.svg`, `app/apple-icon.png`, `manifest.ts`.
- `app/sitemap.ts` and `app/robots.ts`.
- `app/not-found.tsx` and `app/error.tsx` that are actually styled.
- `loading.tsx` at each route segment that fetches.
- `next/image` everywhere, with `width`/`height` or `fill` — layout shift is the easiest
  Core Web Vital to lose and the most tedious to win back.
- Server Components by default; `'use client'` only on the leaves that need it.

## 8. Prove it

1. `pnpm dev` → the assigned port, page renders **with styles**. Check the atomic classes
   and the custom properties in DevTools.
2. Toggle the OS to dark. Everything legible, nothing hardcoded.
3. `pnpm build` → clean, no type errors, no ESLint warnings.
4. Keyboard through the whole page: focus always visible, skip link works.
5. Lighthouse on the built output: accessibility 100, performance 95+ on a static page.
   Anything less is a real finding; report it rather than rounding it off.
6. `scripts/task check` green.

## 9. Report

- What kind of site you inferred and why.
- The route structure.
- The token set and what a brand would need to change (it should be one file).
- The port, and that it is stable for this project.
- The Lighthouse numbers you actually saw.
- Anything left as a placeholder that needs real content.

## Anti-patterns

- Tailwind alongside StyleX.
- `@stylexjs/nextjs-plugin` — stale, Babel-based, drops you off the SWC pipeline.
- Turbopack without the PostCSS extraction step, and an unstyled page with no error.
- Building pages before tokens, so the design system is retrofitted onto hardcoded values.
- Hex values or raw `px` in components.
- `'use client'` at the top of `layout.tsx`, making the entire tree a client component.
- Picking a typeface, a colour, or a voice. Neutral is the deliverable.
- Lorem ipsum, which hides the layout problems real copy would reveal.
- A `-p 3000` you will fight with every other project on the machine.

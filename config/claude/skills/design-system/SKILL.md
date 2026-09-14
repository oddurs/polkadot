---
name: design-system
description: Define and build the minimum viable design system for a project — token layers, layout primitives, a starter component set, accessibility and documentation — as a deliberately neutral blank slate that a brand can be applied to later. Use when starting UI work, when asked for a design system, component library, or design tokens, or when a project's UI has drifted into one-off styles.
argument-hint: [project or app to build it for]
---

# The smallest design system that is actually a system

Most design systems fail in one of two directions: a token file and three components that
nobody can build a page from, or ninety components with a Figma library and a governance
board before the product has a user.

This is the line between them. Build every layer below; stop at the end of it. Everything
here is **visually neutral on purpose** — a blank slate with correct structure, so the
brand can be applied later by changing tokens and nothing else.

## The layers

```
tokens      the values         — what a colour, a space, a duration may be
primitives  the layout         — Box, Stack, Grid, Text: no appearance of their own
components  the interactions   — Button, Input, Dialog: composed from the two above
patterns    the compositions   — a form, a page shell: made of components, not new CSS
```

A thing at one layer may only use layers above it. A component that reaches for a raw hex
value has broken the system, and one broken component is how it all unravels.

---

## 1. Tokens

Two layers, and the split is the whole idea.

**Primitives** are the raw palette. They are named for what they *are* and mean nothing on
their own: `blue-500`, `space-4`, `radius-md`. Nothing in the UI uses them directly.

**Semantic tokens** are named for what they are *for*, and reference primitives:
`color.surface`, `color.textMuted`, `color.borderFocus`, `color.dangerFg`. These are the
only tokens a component may use.

The reason: theming, dark mode, and rebranding all become one edit to the semantic layer.
When a component says `blue-500`, dark mode is a find-and-replace across the codebase and
it will be wrong somewhere.

### What to define, and what not to

| Scale | Size | Notes |
|---|---|---|
| **Colour** | 3–4 hues × 9–11 steps | One neutral ramp (most of the UI), one accent, one for danger, one for success. That is enough; a blank-slate neutral ramp plus a single accent is a legitimate starting palette. |
| **Space** | ~8 steps | One scale for margin, padding, and gap. `4 8 12 16 24 32 48 64`. No `13px`. |
| **Type size** | 6–8 steps | A ratio, not arbitrary numbers. Pair each with a line height. |
| **Font weight** | 3 | Regular, medium, bold. Not nine. |
| **Radius** | 4 | none, sm, md, full. |
| **Shadow** | 3–4 | Tied to elevation, not decoration. |
| **Motion** | 3 durations, 3 easings | And a `prefers-reduced-motion` story. |
| **Z-index** | 5 named layers | base, dropdown, sticky, overlay, toast. Never a literal `9999`. |
| **Breakpoints** | 3 | sm, md, lg. A fourth is usually someone's laptop. |

Do **not** define: per-component tokens before there are components, a type scale with
fourteen steps, six shades of success, or an opacity scale. Under-specifying is cheap to
fix; a token nobody uses is permanent.

### In StyleX

Tokens are `defineVars` in a `.stylex.ts` file, which compiles to CSS custom properties —
a flat object, since StyleX does not allow nesting. Express themes as `createTheme`
overrides of the same variables, which is how dark mode and any future brand arrive
without touching a component:

```ts
// tokens/color.stylex.ts
import * as stylex from '@stylexjs/stylex'

const DARK = '@media (prefers-color-scheme: dark)'

export const color = stylex.defineVars({
  surface:     { default: '#ffffff', [DARK]: '#0b0b0c' },
  surfaceRaised:{ default: '#f6f6f7', [DARK]: '#161618' },
  text:        { default: '#111113', [DARK]: '#ededef' },
  textMuted:   { default: '#60646c', [DARK]: '#9b9ba5' },
  border:      { default: '#e0e0e3', [DARK]: '#2a2a2e' },
  accent:      { default: '#3b5bdb', [DARK]: '#748ffc' },
  accentFg:    { default: '#ffffff', [DARK]: '#0b0b0c' },
  danger:      { default: '#c92a2a', [DARK]: '#ff8787' },
})
```

Keep one file per concern — `color`, `space`, `type`, `motion` — because a token file that
holds everything is the one nobody dares to change.

---

## 2. Primitives

Four or five components with **no appearance of their own**. They exist so that no product
code ever writes a raw `display: flex` again, and so spacing comes from the scale by
construction rather than by discipline.

| | |
|---|---|
| `Box` | the escape hatch: padding, background, border, radius — all token-typed |
| `Stack` | vertical flow with `gap` from the space scale |
| `Inline` | horizontal flow with `gap`, wrapping, alignment |
| `Grid` | columns, gap, responsive column counts |
| `Text` | every piece of text: size, weight, colour, truncation, and the right element |

Two rules make them work. Their props accept **only token values** — `gap="4"`, never
`gap="13px"` — so the type system enforces the scale. And each takes an `as` prop so the
semantic element is chosen independently of the styling: a `Text` that is an `h1` is a
heading to a screen reader and to the layout in the same breath.

---

## 3. Components

The starter set is what a real product needs before it needs anything else. Ten to twelve.

**Build these:**

`Button` (variants: primary, secondary, ghost, danger; sizes; loading and disabled states)
· `Input` · `Textarea` · `Select` · `Checkbox` · `Radio` · `Switch` · `Field` (the label,
description, error, and the wiring between them — the piece everyone forgets and then
fixes in every form) · `Dialog` · `Menu` · `Tooltip` · `Toast`.

**Defer until something needs them:** Tabs, Accordion, Table, DatePicker, Combobox,
Pagination, Avatar, Breadcrumb, Carousel. Every one of these is a week you have not yet
earned.

### Do not hand-roll behaviour

A dialog is a focus trap, an escape handler, a scroll lock, an aria-modal, inert
background content, and a focus restore. A menu is roving tabindex, type-ahead, and arrow
keys. Writing these yourself means writing them badly, and it is the single largest source
of accessibility failures in a young design system.

Take behaviour from a headless library and provide only the styling:

- **React Aria Components** (`react-aria-components`) — the default. Stable 1.x, the most
  rigorous accessibility work of the three, actively released.
- **Base UI** (`@base-ui-components/react`) — excellent API, still 1.0-rc as of now; pick
  it if you are comfortable tracking a release candidate.
- **Radix Primitives** — still solid, slower cadence.

Then every component is: headless behaviour + StyleX styles from tokens. That is the whole
implementation, and it is why twelve components is a realistic starting scope.

### Every component, without exception

- **States**: default, hover, active, focus-visible, disabled, loading, error. A component
  missing `:focus-visible` is not finished — keyboard users cannot see where they are.
- **Variants** as a typed union, never a `className` passthrough. `className` escapes the
  system, and once one component allows it the tokens are advisory.
- **`forwardRef`**, and unrecognised props spread to the underlying element.
- **Controlled and uncontrolled** both supported, the way the platform does it.
- **Keyboard operable**, with a visible focus ring that meets contrast on every surface.
- **Not responsible for its own margin.** Spacing is the parent's job — `Stack` owns the
  gap. A component with a built-in `margin-bottom` fights every layout it is ever put in.

---

## 4. Accessibility is a gate, not a phase

- Every interactive element is reachable and operable by keyboard alone, in a sensible
  order, with visible focus.
- Colour contrast: 4.5:1 for text, 3:1 for large text and for the boundaries of UI
  components. Check the semantic tokens against each surface they can appear on, in both
  themes — this is where a neutral palette quietly fails.
- Colour is never the only carrier of meaning. An error is an icon and a message.
- Respect `prefers-reduced-motion` and `prefers-color-scheme`.
- Every form control has a programmatically associated label. `Field` exists to make that
  impossible to get wrong.
- Run `axe` in the component tests, and make it fail the build. An accessibility check
  that reports instead of failing is a report nobody reads.

---

## 5. Documentation

A design system nobody can read is a component folder. Docs are part of the deliverable,
not a follow-up.

**Every component gets**, on one page:

1. **What it is and when to use it** — with the alternative it is confused with.
2. **A live example** you can interact with.
3. **Anatomy** — the parts and their names.
4. **Every variant and state**, rendered together so they can be compared.
5. **Props**, generated from the types so they cannot go stale.
6. **Do / Don't**, two or three pairs, specific to this component.
7. **Accessibility notes** — keyboard map, what is announced.

**Also document the tokens themselves**: a page that renders every colour on both themes,
the space scale to size, the type scale in situ. It is the fastest way to notice that two
tokens are the same value, or that a step is missing.

Storybook is the default for this and earns its keep through the variant matrix, the a11y
addon, and interaction tests. A `/docs` route inside the app is a fine alternative and one
fewer build to maintain — take it if the system will only ever serve one app.

---

## 6. Governance, in three rules

1. **Rule of three.** A pattern becomes a component the third time it appears. The first
   is a one-off, the second is a coincidence.
2. **New component requires a real use.** No speculative components; they are always wrong
   about the use they were built for.
3. **Changing a token is a system change.** Adding one is cheap. Changing one moves the
   whole product, and that is exactly the point — it should be deliberate.

---

## 7. Prove it

The system is real when you can build a page you did not plan for using only what exists:

1. Build a sign-in form and a settings page from primitives and components only. No raw
   CSS, no one-off values. Whatever you had to reach outside the system for is the next
   item on the list.
2. Toggle to dark. Everything remains legible and nothing is hardcoded.
3. Navigate both pages by keyboard alone. Focus is always visible and never trapped
   anywhere but a dialog.
4. `axe` clean.
5. Grep the app for hex colours, `px` values outside the token files, and `className` on
   design system components. Each hit is either a missing token or a leak.

## Anti-patterns

- Components consuming primitive tokens directly, so theming has to touch every file.
- A `className` or `style` prop on design system components. The system is now advisory.
- Component-level tokens (`buttonPrimaryHoverBackground`) before there is a second theme.
- Building Table, DatePicker, or Combobox before Button has a focus state.
- Hand-rolled dialogs, menus, and tooltips.
- Components that set their own margin.
- Documentation as a README listing prop names.
- Designing the brand into the system. Neutral tokens, brand applied as a theme — otherwise
  the rebrand is a rewrite.

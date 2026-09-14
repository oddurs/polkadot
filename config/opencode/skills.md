# opencode skills

These are not in this repository. They are third-party skill packages — most
from [vercel-labs/agent-skills](https://github.com/vercel-labs/agent-skills),
some from Supabase — and they come to about two megabytes of somebody else's
MIT-licensed writing. Vendoring that here would make this repo mostly other
people's work, and would put the attribution burden on a dotfiles repo that has
no business carrying it.

They live in `~/.config/opencode/skills/`, which `.gitignore` excludes.

| skill | from |
|---|---|
| `composition-patterns` | vercel · MIT |
| `react-best-practices` | vercel · MIT |
| `web-design-guidelines` | vercel |
| `supabase-postgres` | supabase · MIT |
| `next-best-practices` | vercel-labs/agent-skills |
| `next-cache-components` | vercel-labs/agent-skills |
| `vercel-ai-sdk` | vercel-labs/agent-skills |
| `tailwind-v4-shadcn` | vercel-labs/agent-skills |
| `tanstack-query` | vercel-labs/agent-skills |
| `forms-validation` | vercel-labs/agent-skills |
| `auth-patterns` | vercel-labs/agent-skills |
| `zustand` | vercel-labs/agent-skills |
| `testing-patterns` | MIT |
| `api-testing`, `database-design`, `debugging`, `security`, `unit-testing` | unattributed |

To restore them:

```sh
npx skills add vercel-labs/agent-skills
npx skills list                 # what is installed, and where
```

`npx skills list` reports all of the above as `Source: local`, so the ones
without a package above were copied in by hand and are not recoverable from a
registry. If they matter, they should become skills of their own in this repo
with their provenance recorded — the five unattributed ones especially.

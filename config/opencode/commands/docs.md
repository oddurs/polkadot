---
description: Generate or update documentation for code
agent: writer
---

Generate documentation for the specified code:

1. **Read the code** — Understand the purpose, API, and usage patterns
2. **Determine doc type** based on context:
   - Function/class → JSDoc/TSDoc with params, returns, examples
   - Component → Props interface, usage example, edge cases
   - API endpoint → Method, path, request/response schema, auth, errors
   - Module/package → README with install, quickstart, API reference
   - Project → README with overview, setup, development, deployment
3. **Include examples** — Real, runnable code examples that cover common use cases
4. **Skip the obvious** — Don't document self-explanatory code

Focus on: $ARGUMENTS

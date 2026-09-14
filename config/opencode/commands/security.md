---
description: Security audit for vulnerabilities and best practices
agent: reviewer
---

Perform a security audit on the specified code:

## Check for OWASP Top 10

1. **Injection** — SQL injection, XSS, command injection, SSRF
2. **Broken Auth** — Missing auth checks, weak session handling, credential exposure
3. **Sensitive Data Exposure** — Unencrypted secrets, leaked API keys, PII in logs
4. **Broken Access Control** — Missing authorization, IDOR, privilege escalation
5. **Security Misconfiguration** — Debug mode in prod, default credentials, open CORS
6. **Insecure Dependencies** — Run `npm audit` or equivalent, check for known CVEs
7. **Input Validation** — Missing server-side validation, type coercion exploits

## Additional Checks

- Environment variables for secrets (not hardcoded)
- CSRF protection on state-changing endpoints
- Rate limiting on auth endpoints
- Content Security Policy headers
- SQL parameterized queries (no string concatenation)
- File upload validation (type, size, path traversal)

Focus on: $ARGUMENTS

Report findings with severity (CRITICAL/HIGH/MEDIUM/LOW) and concrete fix suggestions.

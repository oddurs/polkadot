# Security Policy

Polkadot runs on a developer's primary machine and can install software, alter
shell configuration, and create symlinks. Treat defects in those boundaries as
security issues.

Do not open a public issue for a vulnerability involving credential exposure,
unsafe file replacement, command injection, privilege escalation, or an
untrusted installer path. Use GitHub's private security advisory flow for this
repository instead.

Include the affected command, macOS and tool versions, a minimal reproduction,
impact, and any safe mitigation you found. Redact tokens, private paths, and
machine-specific identifiers.

Only the default branch is supported. Fixes are developed privately when that
reduces exposure, then released with a short advisory once users can update.

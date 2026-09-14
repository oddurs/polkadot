---
paths:
  - "**/*.rs"
  - "**/Cargo.toml"
---

# Rust

- `unwrap` and `expect` belong in tests, in `main`, and nowhere else. If an invariant
  makes one safe, write the invariant in a comment on that line or use `expect` with it
  as the message.
- Libraries return a `thiserror` enum. Binaries use `anyhow` at the edge, and only there.
- Take `&str`, `&[T]`, `impl AsRef<Path>`. Return owned values. Don't allocate to satisfy
  a signature you control.
- Reach for a `clone()` when it makes the code clear and the profiler hasn't objected.
  Lifetimes in a public signature must pay for themselves.
- `unsafe` needs a `// SAFETY:` comment naming the invariant the caller must uphold.
- Prefer `matches!`, `let ... else`, and `if let` chains over nested matches.
- Tests live in `#[cfg(test)] mod tests` beside the code; integration tests in `tests/`.
  `insta` for anything whose expected value is more than a line.
- Pin the toolchain in `rust-toolchain.toml`. Lint with `-D warnings`.

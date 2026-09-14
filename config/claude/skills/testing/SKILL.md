---
name: testing
description: Install a complete, deterministic test setup for a project and write the tests that go with it — unit, property, snapshot, integration, and end-to-end, wired into scripts/task test and CI, with fixtures, fakes, and a hard line against flakes. Use when a project has no tests, thin tests, a flaky suite, or when asked to set up testing, add coverage, or make the suite trustworthy.
argument-hint: [path or subsystem to focus on]
---

# A suite you can believe

A test suite has one job: when it is green, you ship without thinking about it. A suite
that is sometimes red for reasons nobody investigates is worse than none, because it
teaches everyone to ignore a red build.

Execute — install the harness *and* write real tests against the code that exists. A
framework with no tests in it has not been set up.

## 0. Read what is there

The stack, the existing tests, what `scripts/task test` runs today, and — most
importantly — **what this code actually does wrong**. Read the bug tracker and
`git log --grep=fix`. The bugs a project has already had tell you where its tests are
missing, and they are the first ones to write.

## 1. The shape

Most tests are unit tests, because they are fast and precise. A few are integration
tests, because units that each work can still not work together. Very few are end to
end, because they are slow and break for reasons unrelated to your change.

| Level | Tests | Speed | How many |
|---|---|---|---|
| Unit | one function, one behaviour, no I/O | µs | most |
| Property | an invariant over generated input | ms | for anything with a law |
| Snapshot | output whose shape matters more than its text | ms | for rendering, formatting, CLI |
| Integration | the seams between real components | ms–s | one per seam |
| End to end | the real thing, the way a user runs it | s | the handful that would embarrass you |

Inverting this — mostly end-to-end, few units — is the single most common way a suite
becomes slow, flaky, and useless at saying *what* broke.

## 2. Pick the tools, don't collect them

Opinionated defaults. One choice per job; a second library doing the same job is a bug.

**Rust** — built-in `#[test]`, `cargo-nextest` as the runner (real parallelism, proper
process isolation, clear output). `insta` for snapshots, `proptest` for properties,
`rstest` for parametrised cases, `assert_cmd` + `predicates` for CLI behaviour,
`tempfile` for filesystem work, `mockall` only where a trait boundary genuinely needs a
double.

**TypeScript** — `vitest` (fast, ESM-native, one config with the app). `@testing-library`
for components — query by role and label, never by class name or test id unless there is
no alternative. `msw` to intercept HTTP at the network layer rather than stubbing your own
client. `playwright` for end to end, `fast-check` for properties.

**Python** — `pytest`, and nothing else that overlaps it. `pytest-xdist` for parallelism,
`hypothesis` for properties, `syrupy` for snapshots, `freezegun` or injected clocks for
time. Fixtures over setup methods; `tmp_path` over anything hand-rolled.

**Go** — the standard library, table-driven, `t.Parallel()` by default, `testify/require`
only if the assertions are genuinely unreadable without it. `testing/quick` or
`rapid` for properties, golden files for snapshots.

Whatever the stack: the runner is wired into `scripts/task test`, and CI runs
`scripts/task check`. Nothing in CI knows the runner's name.

## 3. Write the tests

A test has a name that states the behaviour — `rejects_a_subject_over_72_characters`,
not `test_commit_msg_2`. When it fails, the name alone should tell you what broke.

Arrange, act, assert. One behaviour per test; a test asserting five things fails for five
reasons and tells you about one.

**Start where the value is**, in this order:

1. **Every bug that has ever been fixed.** Write the test that would have caught it. If
   the fix is already merged without a test, the bug is not closed.
2. **The boundaries.** Empty, one, many, maximum, one past maximum. Zero, negative, NaN.
   Empty string, whitespace, unicode, something 10MB long. The happy path rarely breaks;
   the edges always do.
3. **The error paths.** Most untested code is error handling, and most production
   incidents are error handling. Assert the error *and its message* — a test that accepts
   any failure will pass when the code fails for a new reason.
4. **The invariants**, as property tests. Round-trips (`parse(render(x)) == x`),
   idempotence, ordering, conservation. One property test finds what fifty examples miss,
   and a failure comes with a minimal reproduction.
5. **The seams.** Each real integration point, once, with the real thing where you can.

## 4. Determinism is not negotiable

A flaky test is a failing test. Fix it or delete it — never retry it.

- **No sleeps.** Ever. Wait on the condition, not on the clock. `sleep(100)` is a race
  you have decided to lose later.
- **Inject time.** A `Clock` parameter, a frozen clock, a fake. Never `now()` inside code
  under test.
- **Inject randomness.** A seeded generator. A property test that fails prints its seed,
  and the seed reproduces it.
- **No network.** Not to localhost, not to a service you control, not "just this one".
  Intercept at the boundary; the end-to-end tests are where real I/O lives, and they are
  quarantined.
- **No shared state between tests.** Each gets its own temp dir, its own database, its own
  port from the OS (bind `:0` and ask what you got — never hardcode a port).
- **Order-independent.** Run the suite shuffled in CI. If order matters, a test is leaking.
- **No test depends on another test having run.**

The check: run the suite ten times in a row, in parallel, shuffled. Ten greens or it is
not finished.

## 5. Fixtures and doubles

Build test data with a helper that takes overrides, so a test states only what it cares
about:

```
let commit = a_commit().with_subject("x".repeat(73)).build();
```

Not thirty lines of struct literal in which the one significant field is invisible.

Fake at the boundary you own — the trait, the interface, the port. Mocking a type you do
not control couples the test to a library's internals and the test breaks on upgrade.
Prefer a real implementation where it is cheap: a temp directory beats a mock filesystem,
an in-memory SQLite beats a mock repository.

Assert on behaviour, not on calls. "The file contains the right bytes" survives a
refactor; "`write` was called twice with these arguments" does not.

## 6. Coverage

Measure it; do not target it. Coverage tells you what is *not* tested, which is useful.
A percentage gate tells people to write tests for getters, which is not.

Wire it up (`cargo-llvm-cov`, `vitest --coverage`, `pytest --cov`, `go test -cover`),
print it in CI, and look at the uncovered lines rather than the number. Uncovered error
handling is a finding. 100% coverage with no assertions about error messages is not a
green light.

## 7. Wire it in

- `scripts/task test` runs the whole suite, and is the only name CI and the hooks know.
- `pre-push` runs `scripts/task check`, so nothing red leaves the machine.
- CI runs the same line. If CI runs something local hooks don't, they will drift and one
  of them will start lying.
- The README says how to run one test, not just all of them — it is what people need at
  3pm on a Tuesday.

## 8. Prove it

Do not report a suite as installed because the files exist:

1. `scripts/task test` — green, and show the output.
2. Break something real on purpose — invert a comparison in the code under test — and
   watch a test fail with a message that names the behaviour. A suite that stays green
   when you break the code is decorative. Revert.
3. Run it ten times shuffled and parallel. Ten greens.
4. Report the runtime. A unit suite over a minute will stop being run before it is a year
   old; say so if it is already there.

## Refuse

- A test with no assertion.
- A test that retries, sleeps, or is marked flaky-tolerant.
- `#[ignore]`, `.skip`, `xit`, `t.Skip` left behind without an item saying when it comes back.
- Asserting a function was called instead of what it did.
- A snapshot committed without being read. A snapshot you didn't check is a record of a bug.
- Chasing a coverage number by testing generated code, getters, or `Debug`.

---
name: perf
description: Make a project fast and make it survive — profile before changing anything, benchmark with a regression gate, then harden the failure paths with timeouts, bounded retries, backpressure, idempotency, and graceful shutdown. Use when something is slow, when asked to optimize or profile, or when a service needs to be reliable under load, failure, or restart.
argument-hint: [what is slow, or the subsystem to harden]
---

# Fast, and still standing

Two jobs that people conflate. **Performance** is how much work it does per unit of time.
**Durability** is whether it is still correct after the network dropped, the disk filled,
the process was killed mid-write, and the same request arrived twice. A fast system that
loses data is not a good system.

Do them in that order, and never start either one by guessing.

---

# Part 1 — Performance

## 1. Measure, or don't touch it

**You do not know what is slow.** Neither do I. Every experienced engineer has spent a day
optimising something that turned out to be 2% of runtime.

Before any change:

1. **Reproduce it.** A command, a request, a workload that is slow every time.
2. **Put a number on it.** Wall clock, p50 and p99, allocations, bytes, queries. "Feels
   slow" cannot be improved and cannot be finished.
3. **Set the target.** "Under 100ms at p99 for 10k rows." Without one, optimisation has
   no end and you will keep going past the point of value.
4. **Profile.** Look at where the time goes — do not skip to the fix.

| | Profiler |
|---|---|
| Rust | `cargo flamegraph`, `samply`, `perf`; `dhat` for allocations |
| Node | `node --cpu-prof`, `--heap-prof`, Chrome DevTools; `clinic` |
| Python | `py-spy` (no instrumentation, works on a running process), `memray` |
| Go | `pprof` — CPU, heap, block, mutex |
| Browser | Performance panel, then Lighthouse for the load path |
| Database | `EXPLAIN ANALYZE` on the real query with real data volumes |

Read the flame graph before forming a theory. The answer is usually not where you assumed,
and it is usually one of: an N+1 query, a missing index, work inside a loop that belongs
outside it, an accidental O(n²) from a lookup in a list, serialising something large, or
doing sequentially what could be concurrent.

## 2. Fix in order of leverage

1. **Do it less.** Cache, memoise, batch, debounce. The fastest work is skipped work.
2. **Do it at a better time.** Move it out of the request path: precompute, background,
   lazy-load, stream instead of buffering.
3. **Use a better algorithm or a better shape.** A hash lookup instead of a scan, one
   query instead of N, an index, the right data structure. This is where the order-of-
   magnitude wins live.
4. **Do it concurrently.** Parallelism where the work is independent, and only after the
   serial version is correct.
5. **Make the operation itself cheaper.** Fewer allocations, less copying, arena or pool,
   a tighter loop. Real, but it is a constant factor — do it last.

Micro-optimising a constant factor before fixing an N+1 is the classic wasted afternoon.

## 3. Benchmark and gate

An improvement you cannot measure again next month is one you will lose next month.

- `criterion` (Rust), `vitest bench`/`tinybench` (TS), `pytest-benchmark` (Python),
  `go test -bench` — checked in, run on demand, wired behind `scripts/task`.
- Benchmark a realistic workload with realistic data volumes. A benchmark over ten rows
  measures your framework's startup.
- Record the before and after in the commit message. A perf commit without numbers is a
  refactor with a claim attached.
- Put the important ones in CI with a threshold, so a regression fails a PR instead of
  being discovered in six months. Tolerate noise: compare against a band, not a point,
  and have it fail on a sustained shift rather than a single run.

## 4. Stop at the target

When the number is met, stop. Write down what you did and what the new number is. Further
optimisation now costs clarity and buys nothing anybody asked for — and the complexity you
add is permanent while the speed is not.

---

# Part 2 — Durability

Everything here is about the paths that are not the happy one. They are the least tested
code in most projects and the cause of most incidents.

## 5. Everything that can hang, times out

Every network call, every lock, every subprocess, every query. A call with no timeout is
an unbounded wait, and one of them is enough to exhaust a pool and stall a whole service.

The timeout budget is set by the caller and passed down: if the request has 2s, a
downstream call cannot have 30s. Cancellation propagates — a client that gave up should
stop work still running on its behalf.

## 6. Retry narrowly, back off, and cap

Retry only what is **transient** (connection reset, 503, timeout, deadlock) and only what
is **idempotent**. Retrying a non-idempotent write is how a customer gets charged twice.

Exponential backoff with **jitter** — without jitter every client retries in lockstep and
you have built a thundering herd that keeps the service down. Cap the attempts and the
total elapsed time. Around a dependency that is already failing, a circuit breaker beats
retries: fail fast, recover in half-open.

## 7. Make writes idempotent

An operation that can be safely repeated removes a whole class of failure. Idempotency
keys on requests, `INSERT ... ON CONFLICT`, conditional writes with an expected version,
natural keys instead of "insert and hope". Then at-least-once delivery — which is all you
ever really get — becomes effectively-once.

## 8. Bound everything

Unbounded is a memory leak waiting for the right day: queues, buffers, caches, connection
pools, batch sizes, request bodies, result sets, recursion.

When the bound is hit, apply **backpressure** — slow the producer, reject with 429/503,
shed load. Do not buffer without limit; a queue that grows until the process dies has
turned a slow dependency into an outage. Pagination on every list endpoint, `LIMIT` on
every query that could return the whole table.

## 9. Fail as one thing

A write that is half done is worse than one that did not happen. Transactions around
multi-step state changes. Write-to-temp-then-atomic-rename for files. For work that
crosses systems, an outbox or a saga with compensation — a distributed transaction is not
available to you.

Ask of every operation: **if the process dies exactly here, what is the state?** Walk the
line after each state change. If any answer is "corrupt" or "half", restructure until the
answer is "before" or "after".

## 10. Start and stop cleanly

- **Startup**: fail fast on bad config — validate everything at boot, not on first use at
  3am. Do not accept traffic until dependencies are reachable.
- **Shutdown**: catch the signal, stop accepting new work, finish or cleanly abandon
  in-flight work within a deadline, release resources, exit. A deploy should be invisible.
- **Restart**: the process must survive being killed at any instant and come back correct.
  That is the only real test of §9, and `kill -9` in a loop is how you run it.
- **Health**: liveness ("am I running") and readiness ("should I get traffic") are
  different questions and a single endpoint answering both will take you down during a
  dependency blip.

## 11. Make failure visible

Structured logs with a correlation id that crosses service boundaries. Metrics for rate,
errors, and duration on every operation that can fail; saturation for every bounded
resource from §8. Alert on symptoms users feel — error rate, latency — not on CPU.

An error log with no context is noise. The message names what failed, what it was doing,
and what the caller should do about it.

## 12. Prove the failure paths

Assert on them like any other behaviour:

- Tests that inject failure: timeouts, connection resets, partial reads, disk full, a
  dependency returning garbage.
- Kill the process mid-write, restart, assert the state is consistent. In a loop.
- Load test to find the knee, then past it: confirm it degrades and recovers rather than
  falling over. A system that has never been pushed past its limit has an unknown limit.
- Fuzz every parser and every deserialiser (`cargo-fuzz`, `atheris`, `go-fuzz`,
  `fast-check`). Anything that reads bytes somebody else wrote gets fuzzed.
- Chaos in staging: drop the dependency, add 500ms of latency, fill the disk.

## Report

- The number before, the number after, and how you measured both.
- The profile that pointed at the cause — not the theory you started with.
- What you deliberately did not optimise, and why it wasn't worth it.
- The failure modes you hardened, and the ones still open.
- What is now gated in CI, so this does not silently come back.

## Anti-patterns

- Optimising without profiling. Almost always the wrong code.
- A benchmark that measures startup, or ten rows, or a warm cache you forgot to clear.
- Caching to hide an N+1 instead of fixing it. Now you have two problems and one of them
  is stale.
- Retries without jitter, without a cap, or on a non-idempotent operation.
- A queue, buffer, or cache with no maximum.
- `catch { }`. A swallowed error is a bug you will debug twice: once now, once in a year.
- Adding concurrency to code that is not yet correct.
- A perf PR with no numbers in it.

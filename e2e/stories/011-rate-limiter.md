---
priority: regression
estimated_runtime: 2-3min
---

# Story: 011-rate-limiter — All parts × C++ + Go

Tier 1, Strategy + Factory. Different DSA shape from 001 (Queue + HashMap vs. Sorting), so it exercises a different test-runner code path even though both are 3-part problems.

## Problem identity

| Field | Value |
|---|---|
| Problem ID | `011-rate-limiter` |
| Display title | `API Rate Limiter` |
| Tier | 1 |
| Total parts | 3 |
| Languages | cpp, go |
| Patterns | Strategy, Factory |
| DSA | Queue + HashMap |
| Why it matters | Factory pattern means Part 2 introduces a new algorithm class — exercises the cumulative-test path where Part 2 tests still call Part 1 methods. |

## Reference solutions

- C++: `problems/tier1-foundation/011-rate-limiter/solution.cpp`
- Go:  `problems/tier1-foundation/011-rate-limiter/solution.go`

Both are guaranteed to pass all 3 parts when pasted as-is.

## Execution

Execute the recipe in `e2e/COMMON.md` with the inputs above:

- `<problem-id>` = `011-rate-limiter`
- `<problem-title>` = `API Rate Limiter`
- `<total-parts>` = 3
- `<cpp-solution-path>` = `problems/tier1-foundation/011-rate-limiter/solution.cpp`
- `<go-solution-path>` = `problems/tier1-foundation/011-rate-limiter/solution.go`

Expected: 3 C++ submits + 3 Go submits, all passing. Total wall-clock 2–3 minutes.

## Story-specific assertions

- After Part 2 passes (either pass), the test output panel should show both Part 1 tests AND Part 2 tests passing (cumulative behavior — Part 2 doesn't replace Part 1).
- After Part 3 passes, the toolbar status badge reads `Solved`.

## Known caveats

- Rate-limiter tests use timestamps. The reference solutions use logical clocks, not wall-clock time, so flakes are unlikely. If a test flakes due to timing, retry once and report.
- The Go reference solution's `TierBasedFactory.Create` uses sliding-window with `windowSize = limit + 1`, not fixed-window. This is intentional — Part 3's PRO and ENTERPRISE tests issue requests with timestamps spanning more than 60 seconds, which a 60-second fixed-window can't keep in one bucket.

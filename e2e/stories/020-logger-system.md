---
priority: regression
estimated_runtime: 2-3min
---

# Story: 020-logger-system — All parts × C++ + Go

Tier 2 with **3 layered patterns** (Strategy + Observer + Singleton) — the deepest pattern stack in this suite. Tests that the dashboard handles a problem where solutions instantiate multiple cooperating types per part.

## Problem identity

| Field | Value |
|---|---|
| Problem ID | `020-logger-system` |
| Display title | `Logging Framework` (toolbar may show "Logger") |
| Tier | 2 |
| Total parts | 3 |
| Languages | cpp, go |
| Patterns | Strategy, Observer, Singleton |
| DSA | HashMap + Queue |
| Why it matters | Deepest pattern stack in the suite — Singleton enforcement in Part 1 must survive the cumulative-test recompile in Parts 2–3. |

## Reference solutions

- C++: `problems/tier2-intermediate/020-logger-system/solution.cpp`
- Go:  `problems/tier2-intermediate/020-logger-system/solution.go`

Both are guaranteed to pass all 3 parts when pasted as-is.

## Execution

Execute the recipe in `e2e/COMMON.md` with the inputs above:

- `<problem-id>` = `020-logger-system`
- `<problem-title>` = `Logging Framework`
- `<total-parts>` = 3
- `<cpp-solution-path>` = `problems/tier2-intermediate/020-logger-system/solution.cpp`
- `<go-solution-path>` = `problems/tier2-intermediate/020-logger-system/solution.go`

Expected: 3 C++ submits + 3 Go submits, all passing. Total wall-clock 2–3 minutes.

## Story-specific assertions

- After Part 3 passes, status badge reads `Solved` and parts indicator reads `3/3`.
- The cumulative test output across parts should show the Singleton instance count consistent (no test failures with messages like `expected 1 instance, got N`).

## Known caveats

- Singleton tests in Part 1 require that `ClearHistory()` resets the singleton instance — without that, state leaks between tests and `test_singleton_same_instance` produces inconsistent results. Both reference solutions implement this reset.
- C++ Singleton tests sometimes flake under `g++ -O2` due to static-initialization order. The reference solution uses Meyers' singleton (function-local static), which is order-safe. If a test fails with init-order errors, the editor likely truncated the paste — retry the paste once.

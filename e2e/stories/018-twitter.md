---
priority: regression
estimated_runtime: 1-2min
---

# Story: 018-twitter — All parts × C++ + Go

Tier 2 with **only 2 parts** — the only such problem in this suite. Specifically exercises the early-terminate path of the submit loop (no carry-forward dialog after Part 1, since Part 2 is the last).

## Problem identity

| Field | Value |
|---|---|
| Problem ID | `018-twitter` |
| Display title | `Twitter/X Feed System` (toolbar may show "Twitter") |
| Tier | 2 |
| Total parts | **2** |
| Languages | cpp, go |
| Patterns | Observer, Strategy |
| DSA | HashMap + Heap (Merge K Sorted) |
| Why it matters | News-feed problem with k-way merge. The 2-part shape is rare — testing it ensures the submit loop terminates correctly at the second-to-last part. |

## Reference solutions

- C++: `problems/tier2-intermediate/018-twitter/solution.cpp`
- Go:  `problems/tier2-intermediate/018-twitter/solution.go`

Both are guaranteed to pass both parts when pasted as-is.

## Execution

Execute the recipe in `e2e/COMMON.md` with the inputs above:

- `<problem-id>` = `018-twitter`
- `<problem-title>` = `Twitter`
- `<total-parts>` = 2
- `<cpp-solution-path>` = `problems/tier2-intermediate/018-twitter/solution.cpp`
- `<go-solution-path>` = `problems/tier2-intermediate/018-twitter/solution.go`

Expected: 2 C++ submits + 2 Go submits, all passing. Total wall-clock 1–2 minutes.

## Story-specific assertions

Because this problem has only 2 parts, the carry-forward dialog should NOT appear (in either pass). After Part 1 passes:

- A toast `Part 1 passed!` appears.
- The submit button text becomes `▶ Submit Part 2` immediately (no dialog).
- If the carry-forward dialog DOES appear, that's a regression — capture a screenshot and report.

After Part 2 passes:

- Toolbar status badge reads `Solved`.
- Parts indicator reads `2/2`.
- `All Parts Complete` banner replaces the submit button.

## Why `regression` not `critical`

If this story fails but the 3-part stories pass, the bug is specifically in the "last part" UI logic — narrow blast radius. If both fail, it's pipeline-wide.

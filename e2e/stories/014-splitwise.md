---
priority: regression
estimated_runtime: 2-3min
---

# Story: 014-splitwise — All parts × C++ + Go

Tier 2 junction problem — graph-based debt simplification (k-way greedy). Shallow design choice in Part 1 limits Part 3 algorithm efficiency. Good shakeout for the cumulative-test-compile path on a Tier 2 problem.

## Problem identity

| Field | Value |
|---|---|
| Problem ID | `014-splitwise` |
| Display title | `Splitwise Expense Sharing` (toolbar will show "Splitwise" or full name) |
| Tier | 2 |
| Total parts | 3 |
| Languages | cpp, go |
| Patterns | Strategy, Observer |
| DSA | Graph + HashMap |
| Why it matters | Junction problem — design directly affects algorithm efficiency in Part 3. Real interview at Flipkart, PhonePe, Razorpay. |

## Reference solutions

- C++: `problems/tier2-intermediate/014-splitwise/solution.cpp`
- Go:  `problems/tier2-intermediate/014-splitwise/solution.go`

Both are guaranteed to pass all 3 parts when pasted as-is.

## Execution

Execute the recipe in `e2e/COMMON.md` with the inputs above:

- `<problem-id>` = `014-splitwise`
- `<problem-title>` = `Splitwise`
- `<total-parts>` = 3
- `<cpp-solution-path>` = `problems/tier2-intermediate/014-splitwise/solution.cpp`
- `<go-solution-path>` = `problems/tier2-intermediate/014-splitwise/solution.go`

Expected: 3 C++ submits + 3 Go submits, all passing. Total wall-clock 2–3 minutes.

## Story-specific assertions

After Part 3 passes (either pass), verify on the problem page:

- Toolbar status badge reads `Solved`
- Parts indicator reads `3/3`
- The right-panel `Solution` tab is now unlocked (no padlock icon)

## Known caveats

- Splitwise's Part 3 (`simplifyDebts`) returns transactions in an order dependent on map iteration. Both reference solutions sort by amount (and by user name as a tiebreak in Go) so output is deterministic — but if a test flakes on transaction order, retry once and report.

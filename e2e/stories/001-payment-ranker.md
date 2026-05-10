---
priority: critical
estimated_runtime: 1-2min
---

# Story: 001-payment-ranker — All parts × C++ (Go pass blocked)

The smallest Tier 1 problem in the catalog. Acts as the canonical "happy path works" check. Same problem `_smoke.md` opens (no compile), so a green run here means the smoke + submit pipeline both work end-to-end.

## Problem identity

| Field | Value |
|---|---|
| Problem ID | `001-payment-ranker` |
| Display title | `Payment Method Ranker` (toolbar may show shortened form) |
| Tier | 1 |
| Total parts | 3 |
| Languages | cpp (go pass blocked — see below) |
| Patterns | Strategy, Comparator |
| DSA | Sorting |
| Why it matters | Smallest problem with both Strategy + Comparator. Fast feedback loop for "did anything break in the submit pipeline?" |

## Reference solutions

- C++: `problems/tier1-foundation/001-payment-ranker/solution.cpp`
- Go: `(blocked)` — see "Go pass blocked" section below

The C++ solution is guaranteed to pass all 3 parts when pasted as-is into the editor.

## Execution

Execute the recipe in `e2e/COMMON.md` with the inputs above:

- `<problem-id>` = `001-payment-ranker`
- `<problem-title>` = `Payment Method Ranker`
- `<total-parts>` = 3
- `<cpp-solution-path>` = `problems/tier1-foundation/001-payment-ranker/solution.cpp`
- `<go-solution-path>` = `(blocked)`

Expected: Pass 1 (C++) submits 3 parts, all passing. Pass 2 (Go) is **skipped** entirely — see reason below. Total wall-clock 1–2 minutes.

## Go pass blocked

The Go test runners for this problem use **positional struct literals** for `PaymentMethod` with different field counts across parts:

- `tests/go/part1_runner.go` uses 4 fields: `{"UPI", 0.01, 0.0, 1000}`
- `tests/go/part3_runner.go` uses 5 fields: `{"Card A", 0.10, 5.0, 300, false}`

The dashboard's Go submit pipeline is cumulative — submitting Part N compiles `solution.go` together with **all** part runners 1..N. Under that compile, no single `PaymentMethod` struct can satisfy both 4-field and 5-field positional inits, so any `solution.go` that compiles for Part 1 fails for Part 3 and vice versa.

This is an upstream issue with the Go test runners (introduced in PR #2 — "Add full Go language support across all 20 problems"). To unblock:

- **Option A (preferred):** rewrite `tests/go/part{1,2,3}_runner.go` to use **named** struct literals — e.g. `PaymentMethod{Name: "UPI", CashbackRate: 0.01, ...}`. Adding a field then doesn't break older parts. Once that's done, author `solution.go` and remove this block.
- **Option B:** make the dashboard's Go submit pipeline non-cumulative for problems with this shape change. Less general; not recommended.

Until then, this story exercises only the C++ pass — the Go pass is intentionally skipped, not a regression.

## Story-specific assertions (Pass 1 only)

After all 3 C++ parts pass, verify on the problem page:

- Toolbar status badge reads `Solved`
- Parts indicator reads `3/3`
- The right-panel `Solution` tab is now unlocked (no padlock icon)
- The `All Parts Complete` banner is visible in place of the submit button

These checks confirm the unlock-cascade fired correctly for a 3-part problem.

## Why this is `critical`

If this story's C++ pass fails, every other happy-path story will fail in the same way. Run it first when triaging a broken suite — if it's green and 014/018/020 are red, the bug is problem-specific, not pipeline-wide.

---
priority: regression
estimated_runtime: 3-4min
---

# Story: Negative paths — Failure UIs render correctly

Validates the dashboard's error and failure-rendering paths, which the happy-path stories never exercise. Uses one problem (`001-payment-ranker`) for the entire story so we don't burn time switching contexts.

Reference solution path: `problems/tier1-foundation/001-payment-ranker/solution.cpp`

## Pre-flight

- Server reachable (`/api/runner-status` returns `cpp: true`).
- Active progress file determined (per `e2e/COMMON.md` pre-flight step 2).
- Reset progress for `001-payment-ranker`: delete that key from the active progress file.

## Phase 1 — Compile error rendering

1. Open `http://localhost:3000/problem/001-payment-ranker`. Verify editor loads.
2. Click the editor, Ctrl+A, Delete. Paste:
   ```cpp
   // intentionally broken — missing class definition
   int main() { return broken; }
   ```
3. Click `▶ Submit Part 1`.
4. **Verify:**
   - Submit button briefly shows `Compiling…`.
   - Test output panel renders compile error text (containing `error:` or `was not declared` or similar g++ diagnostic).
   - A red toast appears with text matching `Tests failed for Part 1` OR `Submission failed`.
   - The parts indicator stays at `0/3` — no part is marked passed.
   - The status badge in the toolbar reads `In Progress` (because submit attempted updates `started_at`), NOT `Solved`.
5. Capture the response JSON. Verify `compilation.success === false` and `compilation.errors` is a non-empty string.

## Phase 2 — Partial-pass rendering

1. Reset progress for `001-payment-ranker` again.
2. Reload the page. Verify parts indicator back to `0/3`.
3. Paste the reference C++ solution (`problems/tier1-foundation/001-payment-ranker/solution.cpp`) but **deliberately corrupt one method's return value** to make some tests fail. For example, find the method `getRanking()` (or similar) and change `return result;` to `return {};`. Keep the rest intact so it still compiles.
4. Click `▶ Submit Part 1`.
5. **Verify:**
   - Compilation succeeds (no compile error panel).
   - Test output panel renders BOTH `PASS <name>` and `FAIL <name>` lines (mixed result).
   - The summary line reads `PART1_SUMMARY: <K>/<T>` where `K < T` (partial pass).
   - Red toast appears.
   - Part 1 is `attempted`, not `passed`.
   - Part 2 stays locked.
6. Capture the response. Verify `parts[0].all_passed === false` and `parts[0].tests` array contains both `passed: true` and `passed: false` entries.

## Phase 3 — Carry-forward "Start fresh" branch

The happy-path recipe always picks `Continue with my code`. This phase exercises the other branch.

1. Reset progress for `001-payment-ranker`.
2. Reload. Paste the correct reference solution.
3. Click `▶ Submit Part 1`. Verify it passes.
4. When the carry-forward dialog appears for Part 2, click **`Start fresh`** (not `Continue with my code`).
5. **Verify:**
   - Editor content changes — it should now contain the Part 2 starter (look for `// HINT:` or `// TODO:` markers, or a different class structure).
   - Line count differs from the reference solution.
6. **Submit attempt with starter (will fail):** click `▶ Submit Part 2` without modifying.
7. **Verify:**
   - Tests fail (the starter has TODOs, won't pass).
   - Part 2 is `attempted` (not `passed`).
   - Part 3 stays locked (gating works on failure).

## Pass criteria

- All 3 phases produce the expected failure / partial-pass UI without crashing the page.
- No JavaScript exceptions in `browser_console_messages` during any phase.
- Reset works correctly between phases.

## Why this exists

The happy-path stories validate the green flow. The whole point of a code-review tool is to render mistakes well. If error rendering or partial-pass display breaks, users get confused failures — which is much worse than no test runner at all. This story specifically exists to catch that class of bug.

---
priority: planning
estimated_runtime: n/a
---

# Story: data-testid audit — drift detector

This is **not a pass/fail story**. It documents which `data-testid` attributes the rest of the suite would prefer if/when the dashboard adds them. Run this manually when you want a punch list of what to wire up.

## Why this exists

The dashboard currently has **zero `data-testid` attributes**. Stories rely on visible text (e.g. `▶ Submit Part 1`, `All Parts Complete`). That works at this scale but:

- **Text changes break tests silently.** Renaming `▶ Submit Part 1` → `Submit (Part 1)` would break all 5 happy-path stories at once.
- **Localization would nuke the suite.** If the dashboard ever gets i18n, every text-based selector is dead.
- **Race conditions are harder.** "Wait for the submit button to read `Submitting…`" is a string-match race — testids with `data-state="submitting"` would be deterministic.

If/when you decide testids are worth adding, this file is the spec for what to add and where.

## Recommended testids

| Component | Element | Recommended testid | Extra attrs |
|---|---|---|---|
| `ProblemView.jsx` | Submit button | `submit-part-button` | `data-part="N"`, `data-state="idle\|submitting"` |
| `ProblemView.jsx` | All-parts-done banner | `all-parts-complete-banner` | – |
| `ProblemView.jsx` | Toast (success) | `toast-success` | `data-part="N"` |
| `ProblemView.jsx` | Toast (failure) | `toast-failure` | `data-part="N"` |
| `ProblemView.jsx` | Save button | `save-code-button` | `data-state="idle\|saved"` |
| `ProblemView.jsx` | Reset button | `reset-code-button` | – |
| `CarryForwardDialog.jsx` | Modal wrapper | `carry-forward-dialog` | – |
| `CarryForwardDialog.jsx` | Continue button | `carry-forward-continue` | – |
| `CarryForwardDialog.jsx` | Start fresh button | `carry-forward-start-fresh` | – |
| `PartProgressBar.jsx` | Status badge | `status-badge` | `data-status="solved\|attempted\|unsolved"` |
| `PartProgressBar.jsx` | Parts indicator | `parts-indicator` | `data-passed`, `data-total` |
| `ProblemTable.jsx` | Each problem row | `problem-row-<id>` | `data-status` |
| `ProblemTable.jsx` | Table wrapper | `problems-table-wrapper` | – |
| `DifficultyModeSelector.jsx` | Mode selector root | `difficulty-mode-selector` | `data-mode="interview\|guided\|learning"` |
| `ProblemView.jsx` | Language toggle (root) | `lang-toggle` | – |
| `ProblemView.jsx` | C++ toggle button | `lang-toggle-cpp` | `data-active="true\|false"` |
| `ProblemView.jsx` | Go toggle button | `lang-toggle-go` | `data-active="true\|false"` |
| `TestOutput.jsx` | Test output panel | `test-output-panel` | – |
| `TestOutput.jsx` | Per-part summary line | `part-summary-<N>` | `data-passed`, `data-total` |

## How to verify

After adding testids, run this story manually:

1. `pwsh scripts/e2e-up.ps1 -Detach`
2. Navigate to `http://localhost:3000/problem/001-payment-ranker`.
3. Open devtools → Elements panel.
4. For each row in the table above:
   - Search for the testid (`Ctrl+F` in Elements panel, type `data-testid="<name>"`).
   - If absent, mark it as TODO.
   - If present but on the wrong element (e.g. on a wrapper div instead of the actual button), mark it as MISALIGNED.

5. After all testids are wired up, edit `e2e/COMMON.md`:
   - Replace the visible-text selector table with the testid table from CodeJunction's COMMON.md (which already has the testid-first format).
   - Convert each story's selectors from text to testid.

## Related

CodeJunction's `e2e/COMMON.md` has already done this migration; the table there is the canonical reference for testid names. When porting back, copy that table verbatim — it's already battle-tested.

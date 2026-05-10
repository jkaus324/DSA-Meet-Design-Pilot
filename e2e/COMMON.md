# COMMON: Test-All-Parts Recipe (C++ + Go)

This file is the shared recipe every story under `e2e/stories/` references. Each story names a problem, its parts, and its reference solutions, then says "execute COMMON.md for this problem". This file owns the actual UI sequence.

## UI selector reference

The dashboard does **not** currently ship `data-testid` attributes. Stories rely on **visible-text selectors** — slightly more brittle than testids, but stable enough at this scale. If the dashboard adds testids later, see `_testid-audit.md` for the recommended names.

| Element | Visible text | File |
|---|---|---|
| Submit button | `▶ Submit Part <N>` | `ProblemView.jsx` |
| All-parts-done banner | `All Parts Complete` | `ProblemView.jsx` |
| Submitting state (transient) | `Compiling…` then `Submitting…` (same button) | `ProblemView.jsx` |
| Carry-forward dialog | modal with `Continue with my code` and `Start fresh` buttons | `CarryForwardDialog.jsx` |
| Status badge | `Solved` / `In Progress` / `Not Started` | `ProblemView.jsx` |
| Parts indicator | `<passed>/<total>` (toolbar) | `PartProgressBar.jsx` |
| Pass toast | `Part <N> passed!` | `ProblemView.jsx` |
| Fail toast | `Tests failed for Part <N>` | `ProblemView.jsx` |
| Difficulty mode | `Interview` / `Guided` / `Learning` | `DifficultyModeSelector.jsx` |
| Test output panel | `PART<N>_SUMMARY: <K>/<K>`, `PASS`/`FAIL` lines | `TestOutput.jsx` |
| Save code | `Save` → `✓ Saved` | `ProblemView.jsx` |
| Reset code | `Reset` → confirm `Reset Code` | `ProblemView.jsx` |
| Code tab / Solution tab | `Code` / `Solution` | `ProblemView.jsx` |
| Language toggle | `C++` / `Go` buttons (top-right of Code tab) | `ProblemView.jsx` |

If the agent can't find an expected element via text, take a fresh `browser_snapshot` and search for the closest match before failing — and report the discrepancy.

## Inputs the story provides

- `<problem-id>` — e.g. `014-splitwise`
- `<problem-title>` — human-facing name shown in the toolbar
- `<total-parts>` — integer N
- `<cpp-solution-path>` — repo-relative path to `solution.cpp`
- `<go-solution-path>` — repo-relative path to `solution.go`, OR the literal string `(blocked)` if the Go pass is intentionally skipped (with a reason in the story)

## Pre-flight (run once at the top of the story)

1. **Bootstrap server.** Easiest path: run `scripts/e2e-up.ps1 -Detach` (Windows) or `scripts/e2e-up.sh --detach` (POSIX). The script wipes `e2e/.tmp-progress.json` to a clean state, starts the dashboard with `PROGRESS_JSON_PATH=e2e/.tmp-progress.json`, polls `/api/runner-status`, and prints the PID + toolchain status. The user's real `progress.json` is never touched.
2. **Verify toolchain.** The script's output includes `cpp_available` and `go_available` flags from `GET /api/runner-status`. The recipe needs:
   - C++ pass requires `cpp_available: true`
   - Go pass requires `go_available: true`
   - If either compiler is missing, skip just that pass (don't abort the story).
3. **Read reference solutions** the story names from disk (or via `GET /api/problems/<id>/solution?lang=<cpp|go>`). The agent will paste these into the editor later.

## Pass 1 — C++ all parts

### 1.1 Reset progress for this problem only

Edit the active progress file (either `e2e/.tmp-progress.json` or `progress.json` per pre-flight step 2):

- If the file does not exist, skip this step (fresh state).
- If `problems["<problem-id>"]` exists, delete that single key. Preserve everything else.
- Save with the same JSON shape (2-space indent).

### 1.2 Open the problem

- Navigate to `http://localhost:3000/`.
- Take a snapshot. Find the link/card matching `<problem-title>` (problem ID `<problem-id>`).
- Click it. Verify URL becomes `/problem/<problem-id>`.
- Verify the toolbar shows the problem title and a parts indicator reading `0/<total-parts>`.

### 1.3 Set language to C++

- The default language is C++. The toggle is only visible after the editor mounts.
- After the editor is visible, look for the language toggle (`C++` and `Go` buttons near the top-right of the Code tab).
- If the active button (filled background) is `Go`, click `C++`. Otherwise no-op.

### 1.4 Paste reference C++ solution

- Click the Monaco editor area, press Ctrl+A to select all, then Delete.
- Paste the entire content of `<cpp-solution-path>` into the editor.
- Wait until the line count footer at the bottom (text `Lines: <N> | Chars: <M>`) reflects a non-trivial line count (>20 lines).

### 1.5 Submit Part 1 .. Part N (loop)

For each `part` from 1 through N:

1. Find the submit button — text is `▶ Submit Part <part>`. If you instead see the `All Parts Complete` banner, the previous part was the last; exit the loop.
2. Click it. The button text becomes `Compiling…` then `Submitting…`.
3. After the response:
   - **Success path:** a green toast appears with text `Part <part> passed!`. The test output panel shows `PART<part>_SUMMARY: <K>/<K>`.
   - **Failure path:** abort the C++ pass, capture the response details (see Failure capture below), and report. Continue to Pass 2 if Go solution is available.
4. After parts 1..N-1 succeed, a `Carry forward` dialog appears. Click **Continue with my code** (the reference solution covers all parts).
5. After Part N succeeds, no carry dialog appears; the submit area swaps to `All Parts Complete` and the toolbar status badge reads `Solved`.

### 1.6 Capture Pass 1 result

Record per-part:
- pass count / total count
- compile + run time (`time_ms` from the API response is fine)

## Pass 2 — Go all parts

If `<go-solution-path>` is `(blocked)`, **skip this entire pass** and note the reason in the final report. Otherwise continue.

### 2.1 Reset progress again

Repeat step 1.1 — delete `problems["<problem-id>"]` from the active progress file.

### 2.2 Reload the problem page

- Hard-reload the current page (Ctrl+Shift+R or `browser_navigate` to the same URL).
- Verify parts indicator is back to `0/<total-parts>`.

### 2.3 Switch to Go

- Click the `Go` toggle button.
- Wait for the editor to swap. The Monaco editor's language indicator should now read `go`.

### 2.4 Paste reference Go solution

- Select all (Ctrl+A) and delete.
- Paste the entire content of `<go-solution-path>`.

### 2.5 Submit Part 1 .. Part N (loop)

Same as 1.5, but Go. Compilation is via `go build` and slightly slower than `g++` first-run (after that, Go's build cache makes it comparable). Expect each submit to take 2–8 seconds.

### 2.6 Capture Pass 2 result

Same metrics as 1.6.

## Final report

The agent should output a markdown block at the end of the story run:

```markdown
## Run: <problem-id>  —  <YYYY-MM-DD HH:MM>

| Pass | Lang | Parts | All passed? | Total time | Notes |
|---|---|---|---|---|---|
| 1   | cpp  | N    | ✓ / ✗     | <ms>      | <if relevant> |
| 2   | go   | N    | ✓ / ✗ / SKIP | <ms>   | <reason if SKIP> |

Per-part detail (only on failure):
- Part X (cpp): <K>/<total> tests passed — <compile error or failing test names>
```

## Failure capture

On any failure, before moving on:
1. Save the JSON response from `/api/problems/<id>/submit` (status, compilation.errors, parts[].tests).
2. Take a `browser_take_screenshot` of the test output panel.
3. Read the latest `browser_console_messages` and include any errors.
4. Note exactly which step failed (e.g. `Pass 2 step 2.5, Part 3 — compile error in solution.go`).

## Edge cases the recipe must handle

- **Carry-forward dialog doesn't appear after Part N-1.** Sometimes the dialog races with the toast. Wait 500ms after the toast, take a snapshot, look for the dialog. If absent, proceed to Submit Part N+1 directly.
- **Test runner times out.** A submit may hit the 10s run timeout if the machine is busy. Retry the failing submit ONCE. If it fails twice, treat as failure.
- **Compile error on first submit after paste.** Likely indicates the editor truncated the paste. Re-select-all, re-paste, retry once.
- **Active progress file does not exist.** That's the empty-state. Skip the reset; Part 1 will be `active` by default after first interaction.
- **Server returned `cpp_available: false` or `go_available: false`.** That toolchain is missing. Skip just the pass that needs it; the other can still run.
- **Problem has only 2 parts (e.g. 018-twitter).** Loop in step 1.5/2.5 still works — after Part 2 succeeds, the `All Parts Complete` banner appears in place of `Submit Part 3`.
- **Go pass marked `(blocked)`.** Skip Pass 2 entirely. Story will explain why (typically an upstream issue with the test runner that prevents a single solution.go from satisfying all parts).

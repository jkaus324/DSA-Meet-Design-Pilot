---
priority: critical
estimated_runtime: 30s
---

# Story: Smoke — Dashboard shell loads and one problem is reachable

The 30-second sanity check. If this fails, the whole suite is doomed; don't bother running anything else.

## Goal

Verify the dashboard's UI shell renders, the home page lists problems, one problem can be opened, and the editor mounts. **No code submitted, no progress mutated.**

## Pre-flight

- **Server running** at `http://localhost:3000`. Verify with `GET /api/runner-status` returning a valid response *before launching the browser* — if the server is down, abort and report (don't spend 30s in the browser to discover it).
- No assertion on `g++` — this story does not compile anything, so `cpp` flag can be `false`.

## Steps

1. **Load home page.** Navigate to `http://localhost:3000/`.
   - Verify URL resolves (no 5xx).
   - Verify a `Navbar` is visible at the top.
   - Verify a sidebar with at least one nav item (e.g. `Problems`, `Primers`, `Stats`, `Roadmap`) is visible on non-mobile widths.
2. **Verify problems list is reachable.** Click the `Problems` nav link.
   - Verify URL becomes `/problems`.
   - Verify the page renders at least 5 problem rows.
   - Verify a filter / search bar is visible.
3. **Open a known problem.** Click the `001-payment-ranker` row (Tier 1, smallest problem in the catalog — guaranteed to be present). Note: rows may be clickable `<tr onclick>`, not anchor links — Playwright dispatches the click correctly via `.click()` on the row.
   - Verify URL becomes `/problem/001-payment-ranker`.
   - Verify the toolbar shows a problem title containing `Payment` (the full name from `problems.yml` is `Payment Method Ranker`).
   - Verify the right panel has a `Code` tab and the Monaco editor area is visible (look for any text containing `class` or `#include` in the editor — the boilerplate has those).
4. **Verify editor loaded with starter content.** The line count footer should read `Lines: <N>` where `N >= 5`.
5. **No assertions on submit, run, or progress changes.** Smoke is read-only.

## Pass criteria

All 4 steps complete without 4xx/5xx, missing elements, or hung loaders. Total time well under 60 seconds.

## When this fails

- 5xx on `/api/problems` → server is broken; check logs for stacktraces.
- Empty problem list → likely `problems.yml` is malformed or no problem folders exist on disk.
- Editor never mounts → frontend bundle is broken; rebuild with `npm run build` and retry.
- Title doesn't match → likely a problems.yml schema change. Update story or `problems.yml`.

## Why this exists

Before running the full suite (~10 minutes), you want a 30-second confirmation that the dashboard isn't fundamentally broken. If `_smoke.md` passes, the rest of the suite's failures are real signal. If `_smoke.md` fails, every other story will fail in the same uninformative way.

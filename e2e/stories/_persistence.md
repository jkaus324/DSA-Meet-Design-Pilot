---
priority: regression
estimated_runtime: 2-3min
---

# Story: Persistence — Code, notes, mode, and per-language isolation survive reload

Validates that user-visible state persists across hard reloads, and that C++ and Go editor content are stored independently per problem. None of the happy-path stories exercise this; if Ctrl+S → reload silently drops your code, users will (rightly) lose trust in the dashboard.

Uses `001-payment-ranker` for the entire story. No submit, no compile — pure state-persistence testing.

## Pre-flight

- Server reachable.
- `cpp_available: true` and `go_available: true` from `/api/runner-status` (Phase 2 below needs the Go toggle to render; if `go_available: false`, skip Phase 2 with a SKIP marker in the report).
- Reset progress for `001-payment-ranker` so we start clean.

## Phase 1 — Code persists across reload (C++)

1. Open `http://localhost:3000/problem/001-payment-ranker`.
2. Confirm language is `C++` (toggle button).
3. Replace editor content with a unique marker:
   ```cpp
   // SENTINEL_E2E_CPP_2026_05_10
   class Solution {};
   ```
4. Press Ctrl+S.
5. **Verify:** Save button briefly shows `✓ Saved` (green check), then settles to `Save`. A toast appears: `Code saved`.
6. Hard-reload the page (Ctrl+Shift+R).
7. **Verify:** editor content still contains `SENTINEL_E2E_CPP_2026_05_10`. Parts indicator unchanged.

## Phase 2 — Code persists across language switch and back

1. Continuing from Phase 1 (C++ sentinel still loaded).
2. Click the `Go` toggle.
3. **Verify:** editor swaps to a different content (Go boilerplate or empty starter beginning with `package main`). The C++ sentinel must NOT be visible — code is per-language.
4. Replace Go editor content with:
   ```go
   // SENTINEL_E2E_GO_2026_05_10
   package main
   ```
5. Press Ctrl+S. Verify save toast.
6. Click the `C++` toggle.
7. **Verify:** editor restores the C++ sentinel from Phase 1.
8. Click `Go` again.
9. **Verify:** editor restores the Go sentinel. Both languages persist independently.

## Phase 3 — Notes persist across reload

1. Stay on `001-payment-ranker`. Click the `Notes` tab in the LEFT panel.
2. Type into the textarea: `E2E_NOTES_SENTINEL_2026_05_10 — checking persistence`.
3. Click anywhere outside the textarea (blur). Verify the small `Saving…` indicator appears, then settles.
4. Hard-reload the page.
5. Click `Notes` tab again.
6. **Verify:** the textarea still contains `E2E_NOTES_SENTINEL_2026_05_10 — checking persistence`.

## Phase 4 — Difficulty mode persists

1. Stay on `001-payment-ranker`. Open the difficulty mode selector in the toolbar.
2. Switch from `Interview` to `Learning`.
3. **Verify:** if a confirm dialog appears warning about parts reset, click `OK` / continue. (Reset is fine — we have no progress to lose.)
4. Hard-reload the page.
5. **Verify:** the difficulty mode selector still reads `Learning`, NOT `Interview`. The editor contains the `learning` mode boilerplate (which has `// TODO:` markers inside method bodies).

## Cleanup

1. Reset progress for `001-payment-ranker` to clear the sentinels and notes.

## Pass criteria

- All 4 phases verify state survives reload / language toggle.
- No console errors during reloads.
- Sentinel strings are byte-identical pre- and post-reload / pre- and post-language-switch (no whitespace mangling).

## Why this exists

Code persistence is the dashboard's #1 implicit promise. If Ctrl+S looks like it works but doesn't actually persist, users lose work and stop trusting the tool. Phase 2 specifically catches a regression where the dashboard might mistakenly use a single editor buffer for both languages — a bug that's silent until a user switches languages and discovers their other-language code is gone.

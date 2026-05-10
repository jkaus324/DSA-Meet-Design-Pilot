# DSA-Meet-Design-Pilot E2E Suite

End-to-end user-story tests for the dashboard. Each story is a plain-English script that an LLM-driven browser agent (Playwright MCP) executes against a locally running server.

## Why English-driven, not code-driven

This is a solo, iteratively-built prep repo. We don't gate merges on E2E and we don't run on every commit. The cost of writing/maintaining `*.spec.ts` Playwright code per problem is higher than the value at this scale. Plain-English stories let us:

- Author and update tests in minutes, not hours
- Run them on demand when shipping a UI change
- Re-use the same stories as bug-reproduction scripts

When a story stabilizes and you want it in CI later, it can be converted to Playwright code without rewriting the intent.

## Suite shape

```
e2e/
├── README.md                       # this file
├── COMMON.md                       # selector reference + reusable "test all parts" recipe (cpp + go)
└── stories/
    ├── _smoke.md                   # critical · 30s — does the dashboard load at all?
    ├── _negative-paths.md          # regression · 3-4min — compile errors, partial pass, "Start fresh"
    ├── _persistence.md             # regression · 2-3min — code/notes/mode survive reload + per-language code isolation
    ├── _testid-audit.md            # planning  · drift-detector — lists testids the suite expects
    ├── 001-payment-ranker.md       # critical   · 1-2min — Tier 1 baseline (cpp pass; go pass blocked, see story)
    ├── 011-rate-limiter.md         # regression · 2-3min — Tier 1, Strategy + Factory
    ├── 014-splitwise.md            # regression · 2-3min — Tier 2, junction problem
    ├── 018-twitter.md              # regression · 1-2min — Tier 2, only 2-part problem in the suite
    └── 020-logger-system.md        # regression · 2-3min — Tier 2, Strategy + Observer + Singleton
```

Each story file is self-contained: it names the problem, lists its parts, points to the reference C++ and Go solutions, and tells the agent to execute the recipe in `COMMON.md`.

## Two-pass design (C++ + Go)

The repo ships both `solution.cpp` and `solution.go` reference solutions for 4 of the 5 happy-path problems. Each story runs:

- **Pass 1** — C++: paste `solution.cpp`, submit Parts 1..N, verify all pass.
- **Pass 2** — Go: reset progress, paste `solution.go`, submit Parts 1..N, verify all pass.

`001-payment-ranker`'s Go pass is currently **blocked** because that problem's Go test runners use positional struct literals with different field counts across parts (Part 1: 4 fields, Part 3: 5 fields). A single `solution.go` can't satisfy both under the cumulative-compile pipeline. The story explicitly skips Pass 2 and explains why; if the upstream Go test runners are rewritten with named struct literals, the block can be lifted by adding a `solution.go`.

## Priority tags

Each story file declares a priority in its YAML frontmatter:

| Priority | Meaning | When to run |
|---|---|---|
| `critical` | If this fails, something is fundamentally broken. Run before/after every UI change. | `_smoke.md`, `001-payment-ranker.md` |
| `regression` | Standard validation — exercises a specific code path with no extra cost. | Most stories |
| `planning` | Not run for pass/fail; documents expectations. | `_testid-audit.md` |

When asking the agent to run the suite, you can filter:

> Run all `priority: critical` stories under `e2e/stories/`.

> Run every story EXCEPT `priority: planning`.

## Prerequisites

1. **Server running with isolated progress file** — easiest path: use the bootstrap script.
   ```powershell
   # Windows PowerShell
   pwsh scripts/e2e-up.ps1 -Detach
   ```
   ```bash
   # bash / zsh
   ./scripts/e2e-up.sh --detach
   ```
   The script wipes `e2e/.tmp-progress.json` to a clean state, starts the dashboard with `PROGRESS_JSON_PATH=e2e/.tmp-progress.json`, polls `/api/runner-status`, and prints the PID + toolchain status (`cpp_available`, `go_available`).

   At startup the server logs `📝 Using progress file: <path>` when the override is active. Confirm you see that line before running stories.

   If you forget the env var, the suite falls back to editing `progress.json` at the repo root. Stories will warn before doing so.

2. **Compilers installed**
   - `g++` (C++17) for C++ pass
   - `go` (1.21+) for Go pass
   - Verify both via `GET http://localhost:3000/api/runner-status` — `cpp_available` and `go_available` should both be `true`.
   - The smoke story doesn't compile anything; runs without either.
   - If only one compiler is available, stories run only that pass.

3. **Playwright MCP installed** in your Claude Code config:
   ```bash
   claude mcp add playwright -- npx -y @playwright/mcp@latest
   ```
   First run downloads ~150 MB of Chromium. One-time.

## Running stories

In a Claude Code session inside this repo:

**Single story:**
> Run `e2e/stories/_smoke.md` against `http://localhost:3000`.

**Filtered run:**
> Run all `priority: critical` stories under `e2e/stories/`.

**Full happy-path suite:**
> Run every file under `e2e/stories/` except `_testid-audit.md`. Stop on the first hard failure but continue past flaky compile timeouts (retry once). Write each result to `e2e/runs/YYYY-MM-DD-<story-id>.md`.

### Recommended cadence

- **Before any UI change** — run `_smoke.md` (30s).
- **After any UI change** — run all `priority: critical` (~3 min).
- **Before tagging a release** — run everything except `_testid-audit.md` (~15 min).

## What each story validates

| Story | Tier | Parts | C++ | Go | What it stress-tests |
|---|---|---|---|---|---|
| `_smoke` | n/a | n/a | – | – | UI shell loads, problems list reachable, editor mounts |
| `_negative-paths` | 1 | 1 of 3 | ✓ | – | Compile-error rendering, partial-pass UI, "Start fresh" branch of carry-forward |
| `_persistence` | 1 | – | ✓ | ✓ | Code/notes/mode survive reload; per-language code isolation |
| `_testid-audit` | n/a | n/a | – | – | Drift detector — currently fails, documents what testids the suite expects |
| `001-payment-ranker` | 1 | 3 | ✓ | **blocked** | Smallest Tier 1 problem; baseline submit loop |
| `011-rate-limiter` | 1 | 3 | ✓ | ✓ | Tier 1, Strategy + Factory, queue-based DSA |
| `014-splitwise` | 2 | 3 | ✓ | ✓ | Tier 2 junction problem; graph-based debt simplification |
| `018-twitter` | 2 | 2 | ✓ | ✓ | Only 2-part problem — exercises the "fewer parts" submit terminator |
| `020-logger-system` | 2 | 3 | ✓ | ✓ | Tier 2 with 3 patterns layered (Strategy + Observer + Singleton) |

Total per full happy-path run: 14 cpp submits + 11 go submits + ~6 negative/persistence checks across the dashboard's compile + run pipeline.

## Failure reporting format

When a story fails, the agent writes:

```markdown
# Run: <story-id>  ·  <date>

## Summary
- Pass: X / Y submissions
- Failed at: <step>

## Failure detail
- Submit endpoint response: <JSON>
- Compilation errors: <if any>
- Last screenshot: <path>
- Console messages: <relevant lines>
```

Reports are written under `e2e/runs/` (gitignored).

## Adding a new story

1. Copy any existing story under `stories/`.
2. Update the YAML frontmatter (`priority`, `estimated_runtime`).
3. Change the problem id, paths, and any story-specific assertions.
4. If it's a happy-path "test all parts" story, point to `COMMON.md` and you're done.
5. If it's a different kind of check (negative, persistence, custom flow), write the steps inline like `_negative-paths.md`.

No new code, no test runner config.

## Why these 5 happy-path problems

This repo has 20 problems (Tier 1 + Tier 2 only). Picking 5 across the catalog gives broad coverage at modest authoring cost:

- **001-payment-ranker** — smallest Tier 1; same problem `_smoke` opens (cheap reuse). Cpp-only currently.
- **011-rate-limiter** — different Tier 1 patterns (Factory) and DSA shape (Queue + HashMap)
- **014-splitwise** — Tier 2 junction problem (Graph + greedy debt simplification)
- **018-twitter** — Tier 2 with only 2 parts; exercises the early-terminate path of the submit loop
- **020-logger-system** — Tier 2 with 3 layered patterns (deepest pattern stack in the suite)

Plus the 4 supplementary stories (`_smoke`, `_negative-paths`, `_persistence`, `_testid-audit`) close gaps the happy-path stories can't see.

If a regression breaks the dashboard, at least one of these 9 will fail before users see it.

## Relationship to the premium repo

This e2e suite is a slimmed-down port of the [CodeJunction](https://github.com/jkaus324/CodeJunction) premium repo's e2e suite. CodeJunction adds Java pass support, more happy-path stories (one per 20-problem bucket), and `data-testid` hooks throughout the dashboard. If you contribute to both repos, keep the suite shape aligned so stories can move freely between them.

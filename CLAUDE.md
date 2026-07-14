# CLAUDE.md — DSA-Meet-Design-Pilot

## Project Overview

**DSA Meets Design** is an interview prep platform that bridges the gap between DSA and Low-Level Design (LLD). Companies like Amazon, Flipkart, Razorpay, and Meesho now ask multi-part design questions even at SDE-1/fresher level — and no good prep resource exists for this format.

Each problem simulates a real LLD interview: a base requirement followed by 2-4 extensions that unlock progressively. The candidate's design either survives the new requirement or exposes why it doesn't.

**Owner:** Jatin Kaushal — SDE at Amazon India, active LinkedIn creator, runs free interview prep sessions on topmate.

## Product Vision & Roadmap

- **Current (Pilot):** Free/open-source repo with 20-25 problems. Acts as community builder and lead magnet via LinkedIn.
- **Full Launch:** 100 real company-tagged problems, cloud-based code execution, Next.js website with SEO/SSR. Monetized as premium content.
- **End Goal:** AlgoMaster-level platform for DSA + Design interview prep.
- **UX Benchmark:** LeetCode — 95% of target audience are LeetCode users, so the experience must feel familiar and intuitive.

## Tech Stack

- **Dashboard:** React + Vite + Tailwind CSS + shadcn/ui
- **Backend:** Express.js (server.js in dashboard/)
- **Languages supported:** C++17, Go, Java, Python, JavaScript — all five run through a single
  spec-driven harness (see Architecture below). The dashboard auto-detects which language runners
  are installed and disables submit for missing ones.
- **Code Execution:** Local runners under `harness/<lang>/`. Interpreted languages (Python, JS,
  Java) read `spec.yaml` + `tests/cases/*.yaml` at submit time; compiled languages (C++, Go)
  codegen a per-part runner that compiles against the user's solution. (Cloud judge planned for full launch.)
- **Future Migration:** Next.js for the production website

## Architecture — spec-driven, one contract per problem

Each problem is defined once by **`spec.yaml`** (types, function signatures, progressive parts)
and **`tests/cases/partN.yaml`** (language-agnostic test cases). There are NO per-problem,
per-language test files. To add a language you write exactly two things:

1. A generic runner at `harness/<lang>/` (runtime YAML for interpreted langs; a codegen script that
   emits a compilable runner for compiled langs — see `harness/cpp/codegen.py` and `harness/go/codegen.py`).
2. A boilerplate emitter in `scripts/gen_stubs.py` (`emit_<lang>`), then
   `python3 scripts/gen_stubs.py --all --lang <lang> --force`.

Then register the language in `dashboard/server.js` (`LANGS` + a submit branch) and in
`dashboard/src/pages/ProblemView.jsx` (`LANG_META`). Verify with `python3 scripts/stress_test.py`.

**Reference solutions** live at `solution.{cpp,go,java,py,js}` per problem and must pass every case.

## Project Structure

```
DSA-Meet-Design-Pilot/
├── problems/
│   ├── tier1-foundation/       # Foundation-level problems
│   │   └── XXX-problem-name/
│   │       ├── README.md       # Problem statement (LeetCode tone)
│   │       ├── DESIGN.md       # Why this pattern, what breaks without it
│   │       ├── boilerplate/cpp/ # interview / guided / learning × per part
│   │       └── tests/cpp/       # GoogleTest suites per part
│   └── tier2-intermediate/     # Intermediate-level problems
├── patterns/                   # Design pattern primers (GFG tone)
├── docs/_data/problems.yml     # Problem registry
├── dashboard/                  # React + Express dashboard app
│   ├── server.js               # Express API server
│   └── src/                    # React frontend
├── progress.json               # Local user progress (gitignored)
└── package.json                # Root — proxies to dashboard/
```

## Commands

```bash
npm install          # Install dashboard dependencies
npm run dev          # Start dev server (dashboard at :5173, API at :3000)
npm run build        # Build dashboard for production
npm start            # Start production server
```

### Prerequisites
- Node.js 18+
- g++ with C++17 support
- GoogleTest installed

## Problem Format Rules

### Structure
- Every problem has **2-4 parts** (base requirement + 1-3 extensions), depending on the problem
- Parts unlock progressively — Part N+1 only unlocks after Part N tests pass
- Problem folders are **sequentially numbered**: `001-`, `002-`, `003-`, etc. No gaps.

### Required Files Per Problem
- `README.md` — Problem statement. **Tone: LeetCode** (precise, formal, constraint-driven)
- `DESIGN.md` — Pattern explanation, what breaks without it. **Tone: GFG** (educational, thorough, beginner-friendly)
- `boilerplate/cpp/` — All 3 difficulty modes must be present for every part:
  - **Interview** — Blank slate, just problem statement and data types
  - **Guided** — Key interfaces defined, `// HINT:` comments, no pattern names
  - **Learning** — Full class structure, `// TODO:` inside method bodies only
- `tests/cpp/` — GoogleTest suites, one per part

### Problem Quality Bar
- Must be based on **real interview questions** asked at actual companies
- Must be tagged with company names for filtering
- Each problem must map to specific **design patterns** and **DSA concepts**

## Content Tone Guide

| Content Type | Tone | Reference |
|---|---|---|
| Problem statements (README.md) | LeetCode — precise, formal, constraint-driven | leetcode.com |
| Pattern primers & learning content | GFG — educational, thorough, beginner-friendly, one-stop-solution | geeksforgeeks.org |
| LinkedIn posts | Storytelling — real interview scenario, tension-building, open-ended question at the end, conversational | See style notes below |
| README.md (repo root) | Conversational, direct, punchy | Current README |

## LinkedIn Post Style

Jatin's high-performing LinkedIn format:
1. **Hook** — Start with a relatable interview story ("My friend called me after his Amazon interview...")
2. **Build tension** — Describe what the candidate did well, then the twist/failure
3. **Reveal the lesson** — Connect it to a design principle or pattern
4. **End with engagement** — Ask readers how they'd solve it / what they think
- Conversational, not academic
- Short paragraphs, line breaks for readability
- No hashtag spam — 2-3 max if any
- Target: Indian dev community preparing for product company switches

## Guardrails

- **NEVER** run `git commit` or `git push` without explicit permission from Jatin
- When adding problems, implement the changes — Jatin will review via the dashboard UI
- When suggesting architecture changes for the full launch, present as a plan first — don't implement without discussion
- Content for newbies should be comprehensive (one-stop-solution philosophy) — don't assume prior design pattern knowledge

## Primary Workflows

1. **Adding new problems** — Scrape real interview questions from the web, preprocess, transform into the multi-part format with all required files
2. **Bug fixes + UI improvements** — Fix bugs first, then improve color schemes, user flows, page redesigns
3. **Content writing** — Pattern primers, learning content for beginners, structured problem statements
4. **Architecture planning** — Post pilot-launch, plan the full paid platform (Next.js migration, cloud judge, DB, auth, payments)
5. **LinkedIn posts** — Promotional content using Jatin's storytelling style to warm up audience for the platform

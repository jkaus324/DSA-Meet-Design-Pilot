# dsa-meets-design Local Dashboard — Detailed Technical Spec
> **Author:** Jatin Kaushal | **Version:** 2.0 | **Date:** March 2026  
> **Status:** Ready for Development  
> **Build tool:** Hand this entire doc to Claude Code with: "Build this. Follow the spec exactly."  
> **Estimated build time:** 6-7 days  
> **What changed in v2.0:** Added Difficulty Mode Selector (Interview / Guided / Learning) — a tiered scaffolding system that controls how much starter code the user sees. Affects: problem file structure, data model, API, Problem View page, Stats page, and introduces a new `DifficultyModeSelector` component.

---

## Table of Contents

1. [Overview](#1-overview)
2. [User Experience Flow](#2-user-experience-flow)
3. [Difficulty Mode System](#3-difficulty-mode-system) ← NEW
4. [Tech Stack & Dependencies](#4-tech-stack--dependencies)
5. [Project File Structure](#5-project-file-structure)
6. [Data Model](#6-data-model)
7. [API Specification](#7-api-specification)
8. [Frontend — Pages & Components](#8-frontend--pages--components)
9. [Visual Design System](#9-visual-design-system)
10. [Page 1: Dashboard (Home)](#10-page-1-dashboard)
11. [Page 2: Problem View](#11-page-2-problem-view)
12. [Page 3: Pattern Primer View](#12-page-3-pattern-primer-view)
13. [Page 4: Stats View](#13-page-4-stats-view)
14. [Startup & Build Configuration](#14-startup--build-configuration)
15. [Error Handling & Edge Cases](#15-error-handling--edge-cases)
16. [Future P1 Extensions](#16-future-p1-extensions)

---

## 1. Overview

### What this is

A local web dashboard that runs on the user's machine. It provides a LeetCode-like browsing, coding, and progress tracking experience for the dsa-meets-design repo. Users browse problems, write code in an in-browser editor, choose their difficulty mode, track their progress, read rendered READMEs and DESIGN.md files — all in a clean web UI.

### What this is NOT

- Not deployed anywhere. Runs on `localhost:3000` only.
- Not a database. Progress is a local JSON file.

### The single startup experience

```bash
git clone https://github.com/USERNAME/dsa-meets-design.git
cd dsa-meets-design
npm install
npm start
# → opens http://localhost:3000 in default browser
```

That's it. Three commands from clone to dashboard.

---

## 2. User Experience Flow

```
User runs `npm start`
    → Express server starts on port 3000
    → Browser opens automatically to localhost:3000
    → Dashboard loads with all problems from problems.yml

User sees:
    → Progress ring (X/Y solved)
    → Stat cards (Tier 1: X/Y, Tier 2: X/Y, Tier 3: X/Y, Primers: X/Y)
    → Problem table (filterable by tier, pattern, company, status)
    → Pattern primers section at bottom

User clicks a problem row:
    → Problem view opens in split-pane layout
    → LEFT PANE: Problem description (README.md rendered)
    → RIGHT PANE: Code editor with difficulty mode selector at the top

User selects a difficulty mode (3 buttons above the editor):
    → 🔴 Interview — blank slate, only data models given
    → 🟡 Guided — structural hints with ??? placeholders
    → 🟢 Learning — full scaffolding with all classes pre-defined
    → Editor loads the corresponding starter code
    → User's code is preserved per-mode (switching doesn't destroy work)

User writes code in the in-browser editor:
    → Clicks "Run" to compile and test
    → Or copies terminal command to run tests locally
    → Marks problem as solved when done
    → Progress ring updates

User clicks "Stats":
    → Sees progress breakdown by tier
    → Sees patterns practiced chart
    → Sees difficulty mode distribution (how many solved per mode)
    → Sees streak and estimated time
```

---

## 3. Difficulty Mode System

### 3.1 The Problem This Solves

In a real Amazon LLD interview, the candidate is NOT handed an abstract class or interface. They must identify the right pattern, name the abstractions, and design the class hierarchy themselves. The current boilerplate gives away the design — it turns an LLD problem into a "fill in the blanks" coding exercise.

The Difficulty Mode Selector lets the same problem serve three audiences:

| Mode | Who it's for | What the user designs | What's given |
|------|-------------|----------------------|-------------|
| 🔴 Interview | Experienced devs, mock interview prep | Everything — abstractions, class hierarchy, implementation | Only the data model (structs/classes that represent domain objects) + an entry point signature |
| 🟡 Guided | Intermediate devs learning patterns | Class names, method bodies, connections between components | Structural hints with `/* ??? */` placeholders showing where abstractions go, directional comments |
| 🟢 Learning | Beginners, people who just read the pattern primer | Only method bodies (implementation logic) | Full scaffolding — all interfaces, concrete classes, and relationships pre-defined |

### 3.2 The Difficulty Mode Selector UI

The mode selector is a **row of 3 buttons** displayed at the top of the code editor panel, above the tab bar (Code / Design / AI Prompt).

```
┌──────────────────────────────────────────────────────────────┐
│  Difficulty:  [🔴 Interview]  [🟡 Guided]  [🟢 Learning]     │
├──────────────────────────────────────────────────────────────┤
│  Code    Design 🔒    AI Prompt               Save   ▶ Run  │
├──────────────────────────────────────────────────────────────┤
│  1  #include <iostream>                                      │
│  2  #include <vector>                                        │
│  ...                                                         │
└──────────────────────────────────────────────────────────────┘
```

**Button behavior:**
- Exactly one mode is active at any time
- Active button: filled background with mode color, white text
- Inactive buttons: transparent background, gray text, subtle border
- Default mode on first visit: 🟢 Learning (safe default for new users)
- The user's selected mode is persisted per-problem in `progress.json`
- Switching modes shows a confirmation dialog IF the user has edited the code: "Switching mode will load different starter code. Your current code for [current mode] will be saved. Continue?"
- The user's written code is preserved separately for each mode — so they can solve in Learning mode, then try again in Interview mode without losing work

**Button styling:**

| Mode | Active BG | Active Text | Inactive Text | Icon |
|------|-----------|-------------|---------------|------|
| Interview | `#dc2626` (red-600) | `#ffffff` | `#a3a3a3` | 🔴 |
| Guided | `#f59e0b` (amber-500) | `#ffffff` | `#a3a3a3` | 🟡 |
| Learning | `#16a34a` (green-600) | `#ffffff` | `#a3a3a3` | 🟢 |

### 3.3 Starter Code Per Mode — Concrete Example

Using Problem 001 (Payment Method Ranker) as reference:

#### 🔴 Interview Mode — Blank Slate

```cpp
#include <iostream>
#include <vector>
#include <string>
#include <algorithm>
using namespace std;

// ─── Data Model (given — do not modify) ─────────────────────

struct PaymentMethod {
    string name;
    double cashbackRate;    // e.g. 0.05 = 5%
    double transactionFee;  // in rupees
    int    usageCount;
};

// ─── Your Design Starts Here ─────────────────────────────────
//
// Design and implement a PaymentRanker that:
//   1. Ranks payment methods by the criteria described in the problem
//   2. Allows new ranking strategies to be added WITHOUT modifying
//      the ranker itself
//
// Think about:
//   - What abstraction lets you swap ranking logic at runtime?
//   - How would you add a 4th ranking criterion with zero changes
//     to existing code?
//
// Entry point (must exist for tests):
//   vector<PaymentMethod> rank(vector<PaymentMethod>& methods, ??? )
//
// ─────────────────────────────────────────────────────────────


```

**What's given:** Only the `PaymentMethod` struct (the domain data model) and guiding comments about what the entry point should look like. No interfaces, no class names, no pattern hints. The user must arrive at Strategy pattern themselves.

**Why this works for interview prep:** In an actual Amazon round, the interviewer gives you the data model and the problem statement. You design everything else. This mode replicates that exactly.

#### 🟡 Guided Mode — Structural Hints

```cpp
#include <iostream>
#include <vector>
#include <string>
#include <algorithm>
using namespace std;

// ─── Data Model (given — do not modify) ─────────────────────

struct PaymentMethod {
    string name;
    double cashbackRate;    // e.g. 0.05 = 5%
    double transactionFee;  // in rupees
    int    usageCount;
};

// ─── Strategy Interface ─────────────────────────────────────
// HINT: This interface lets you swap ranking logic at runtime.
// What method signature would let you compare two PaymentMethods?

class /* YourInterfaceName */ {
public:
    virtual bool /* yourMethodName */(const PaymentMethod& a, const PaymentMethod& b) = 0;
    virtual ~/* YourInterfaceName */() = default;
};

// ─── Concrete Strategies ────────────────────────────────────
// TODO: Implement a strategy for each ranking criterion:
//   - Rewards maximizer (highest cashback first)
//   - Low-fee seeker (lowest transaction fee first)
//   - Trust-based ranker (highest usage count first)

// ─── Ranker ─────────────────────────────────────────────────
// TODO: Implement a PaymentRanker class that:
//   - Accepts any strategy (via constructor or setter)
//   - Has a rank() method that returns sorted payment methods
//   - Does NOT know about specific ranking criteria

// class PaymentRanker {
// public:
//     PaymentRanker(/* what goes here? */);
//     vector<PaymentMethod> rank(vector<PaymentMethod> methods);
// };


```

**What's given:** The shape of the solution — the user can see that there should be an interface, concrete strategies, and a ranker class. But they name everything themselves, define the method signatures, and wire it all together. The `/* ??? */` placeholders signal "you decide what goes here."

**Why this works for intermediate learners:** They get the architectural roadmap without the implementation. They still have to make real design decisions (naming, method signatures, how the ranker uses the strategy) but they're not starting from zero.

#### 🟢 Learning Mode — Full Scaffolding

```cpp
#include <iostream>
#include <vector>
#include <string>
#include <algorithm>
using namespace std;

// ─── Data Model (given — do not modify) ─────────────────────

struct PaymentMethod {
    string name;
    double cashbackRate;    // e.g. 0.05 = 5%
    double transactionFee;  // in rupees
    int    usageCount;
};

// ─── Strategy Interface ─────────────────────────────────────

class RankingStrategy {
public:
    virtual bool compare(const PaymentMethod& a, const PaymentMethod& b) = 0;
    virtual ~RankingStrategy() = default;
};

// ─── Concrete Strategies ────────────────────────────────────
// TODO: Implement the compare() method for each strategy

class RewardsMaximizer : public RankingStrategy {
public:
    bool compare(const PaymentMethod& a, const PaymentMethod& b) override {
        // TODO: return true if 'a' should rank higher than 'b'
        // Higher cashback rate = better ranking
    }
};

class LowFeeSeeker : public RankingStrategy {
public:
    bool compare(const PaymentMethod& a, const PaymentMethod& b) override {
        // TODO: return true if 'a' should rank higher than 'b'
        // Lower transaction fee = better ranking
    }
};

class TrustBasedRanker : public RankingStrategy {
public:
    bool compare(const PaymentMethod& a, const PaymentMethod& b) override {
        // TODO: return true if 'a' should rank higher than 'b'
        // Higher usage count = better ranking
    }
};

// ─── Ranker ─────────────────────────────────────────────────

class PaymentRanker {
private:
    RankingStrategy* strategy;

public:
    PaymentRanker(RankingStrategy* strategy) : strategy(strategy) {}

    void setStrategy(RankingStrategy* newStrategy) {
        strategy = newStrategy;
    }

    vector<PaymentMethod> rank(vector<PaymentMethod> methods) {
        // TODO: Sort methods using the current strategy's compare()
        // Return the sorted vector
    }
};
```

**What's given:** Everything — the interface, all concrete classes, the ranker class with constructor and setter. The user only implements method bodies (the `// TODO` lines). This is the "fill in the blanks" mode.

**Why this works for beginners:** After reading the Strategy primer, they can see the pattern applied to a real problem and focus purely on implementation logic. The design decisions are already made for them — they just need to understand why and write the code.

### 3.4 File Structure Per Problem (Updated)

Each problem now has THREE starter files instead of one:

```
001-payment-ranker/
├── README.md
├── boilerplate/
│   └── cpp/
│       ├── interview.cpp          ← 🔴 Interview mode starter
│       ├── guided.cpp             ← 🟡 Guided mode starter
│       └── learning.cpp           ← 🟢 Learning mode starter (current boilerplate)
├── tests/
│   └── cpp/
│       ├── BasicRankingTest.cpp
│       ├── CashbackExtensionTest.cpp
│       └── EasyRefundExtensionTest.cpp
├── solution/
│   ├── cpp/
│   └── DESIGN.md
└── AI_REVIEW_PROMPT.md
```

### 3.5 Testing Across Modes — Behavior-Based Tests

**Critical architecture decision:** Tests must validate *behavior* (correct output), NOT *structure* (specific class names).

**Why:** In Interview Mode, the user might name their interface `Comparator` instead of `RankingStrategy`, or use `Sorter` instead of `PaymentRanker`. If tests check for `RankingStrategy` by name, Interview Mode becomes impossible.

**How tests work:**

All three modes share the SAME test files. Tests call a known entry point function and verify output order:

```cpp
// tests/cpp/BasicRankingTest.cpp

#include "solution.hpp"  // user's code compiled to this header
#include <cassert>

void test_rewards_ranking() {
    vector<PaymentMethod> methods = {
        {"Credit Card A", 0.05, 20.0, 100},
        {"Credit Card B", 0.10, 30.0, 50},
        {"UPI", 0.02, 0.0, 200}
    };

    // Test: rank by rewards (highest cashback first)
    auto ranked = rank_by_rewards(methods);

    assert(ranked[0].name == "Credit Card B");  // 10% cashback
    assert(ranked[1].name == "Credit Card A");  // 5% cashback
    assert(ranked[2].name == "UPI");            // 2% cashback
}
```

**The bridge:** Each mode's starter code must ultimately expose the same function signatures that tests call. This is enforced by a comment at the bottom of each starter file:

```cpp
// ─── Test Entry Points (must exist for tests to compile) ────
// Your solution must provide these functions:
//
//   vector<PaymentMethod> rank_by_rewards(vector<PaymentMethod> methods);
//   vector<PaymentMethod> rank_by_low_fee(vector<PaymentMethod> methods);
//   vector<PaymentMethod> rank_by_trust(vector<PaymentMethod> methods);
//
// How you implement them internally is up to you.
// ─────────────────────────────────────────────────────────────
```

In **Learning Mode**, these functions are pre-wired (the user just fills in `compare()` bodies). In **Interview Mode**, the user creates them from scratch. In **Guided Mode**, the comment section shows the required functions but the user builds the plumbing.

### 3.6 Writing Starter Files — Guidelines for Problem Authors

When creating a new problem, follow these rules for each mode:

**🔴 Interview Mode (`interview.cpp`):**
- Include ONLY the data model structs/classes (domain objects the problem describes)
- Include required `#include` headers
- Include a comment block describing the entry point function signatures that tests expect
- Include 2-3 "Think about" prompts as comments (questions, NOT answers)
- Do NOT mention any pattern names
- Do NOT define any interfaces or abstract classes
- Do NOT hint at class hierarchy

**🟡 Guided Mode (`guided.cpp`):**
- Include everything from Interview Mode
- Add empty class shells with `/* YourName */` placeholders
- Add `// HINT:` comments that describe the PURPOSE of each abstraction without naming the pattern
- Add `// TODO:` blocks with bullet lists of what to implement
- Show commented-out class signatures (e.g., `// class PaymentRanker { ... };`)
- OK to show the structure (interface → concrete classes → orchestrator) as long as names and signatures are left to the user

**🟢 Learning Mode (`learning.cpp`):**
- Include everything fully defined — interfaces, concrete classes, orchestrator
- Only `// TODO:` markers are inside method bodies
- The user's job is ONLY to write implementation logic
- This is the closest to a LeetCode "fill in the function" experience

### 3.7 Difficulty Mode and Design Tab Gating

The **Design 🔒** tab has different gating rules per mode:

| Mode | When Design tab unlocks |
|------|------------------------|
| 🔴 Interview | Only after marking problem as "Solved" or "Attempted" |
| 🟡 Guided | Only after marking problem as "Attempted" |
| 🟢 Learning | Available immediately (the scaffolding already reveals the design) |

This prevents Interview Mode users from accidentally spoiling the design for themselves by clicking Design before even trying.

---

## 4. Tech Stack & Dependencies

### Backend

| Dependency | Version | Purpose |
|-----------|---------|---------|
| `express` | ^4.18 | HTTP server |
| `marked` | ^12.0 | Markdown → HTML rendering |
| `marked-highlight` | ^2.1 | Syntax highlighting in code blocks |
| `highlight.js` | ^11.9 | Syntax highlighter (used by marked-highlight) |
| `js-yaml` | ^4.1 | Parse problems.yml |
| `open` | ^10.0 | Auto-open browser on startup |
| `chokidar` | ^3.6 | Watch progress.json for external changes (optional, nice-to-have) |

### Frontend

| Dependency | Purpose |
|-----------|---------|
| React 18 | UI framework |
| Tailwind CSS | Utility-first styling |
| React Router | Client-side routing (dashboard, problem view, stats) |
| Vite | Build tool + dev server for frontend |

### Why React + Vite (not served from Express directly)

The Express server is the API layer (serves data, renders markdown). The React frontend is a separate Vite app that gets built to static files and served by Express in production mode. During development, Vite dev server runs on port 5173 and proxies API calls to Express on port 3000.

**In production (what the user experiences):**
```
npm run build    ← builds React app to dashboard/dist/
npm start        ← Express serves API + static files from dashboard/dist/
```

**For the user, it's still just `npm install && npm start`.** The `npm start` script handles the build step if `dashboard/dist/` doesn't exist.

---

## 5. Project File Structure

Everything dashboard-related lives inside a `dashboard/` folder in the repo root:

```
dsa-meets-design/
├── ... (existing repo files: problems/, patterns/, lib/, etc.)
│
├── dashboard/                         ← ALL dashboard code lives here
│   ├── server.js                      ← Express API server (one file)
│   ├── package.json                   ← Dashboard dependencies
│   │
│   ├── src/                           ← React frontend source
│   │   ├── main.jsx                   ← Entry point
│   │   ├── App.jsx                    ← Router + layout
│   │   ├── index.css                  ← Tailwind imports + custom styles
│   │   │
│   │   ├── pages/
│   │   │   ├── Dashboard.jsx          ← Home page (progress + problem table)
│   │   │   ├── ProblemView.jsx        ← Single problem view (README + code editor)
│   │   │   ├── PrimerView.jsx         ← Pattern primer viewer
│   │   │   └── Stats.jsx             ← Progress stats page
│   │   │
│   │   ├── components/
│   │   │   ├── ProgressRing.jsx       ← Circular progress indicator
│   │   │   ├── StatCard.jsx           ← Tier/primer stat cards
│   │   │   ├── ProblemTable.jsx       ← Filterable problem table
│   │   │   ├── FilterBar.jsx          ← Tier/pattern/company/status filters
│   │   │   ├── StatusToggle.jsx       ← Unsolved/Attempted/Solved buttons
│   │   │   ├── DifficultyModeSelector.jsx  ← NEW: Interview/Guided/Learning buttons
│   │   │   ├── CodeEditor.jsx         ← In-browser code editor panel
│   │   │   ├── PatternBadge.jsx       ← Colored pattern name badge
│   │   │   ├── TierBadge.jsx          ← Tier indicator badge
│   │   │   ├── PrimerList.jsx         ← Pattern primers checklist
│   │   │   ├── MarkdownRenderer.jsx   ← Displays rendered markdown from API
│   │   │   ├── CopyCommand.jsx        ← Terminal command with copy button
│   │   │   ├── Navbar.jsx             ← Top navigation bar
│   │   │   └── CTABanner.jsx          ← Topmate CTA banner
│   │   │
│   │   └── lib/
│   │       ├── api.js                 ← API client (fetch wrapper)
│   │       └── constants.js           ← Pattern colors, tier labels, mode config, etc.
│   │
│   ├── vite.config.js                 ← Vite configuration (proxy to Express)
│   ├── tailwind.config.js             ← Tailwind configuration
│   ├── postcss.config.js              ← PostCSS for Tailwind
│   ├── index.html                     ← Vite HTML entry point
│   └── dist/                          ← Built frontend (generated, gitignored)
│
├── progress.json                      ← User's local progress (gitignored)
│
├── package.json                       ← Root package.json (startup scripts only)
└── .gitignore                         ← Includes progress.json and dashboard/dist/
```

### Problem Directory Structure (Updated with per-mode starters)

```
problems/
└── tier1-foundation/
    └── 001-payment-ranker/
        ├── README.md
        ├── boilerplate/
        │   └── cpp/
        │       ├── interview.cpp      ← 🔴 Blank slate starter
        │       ├── guided.cpp         ← 🟡 Structural hints starter
        │       └── learning.cpp       ← 🟢 Full scaffolding starter
        ├── tests/
        │   └── cpp/
        │       ├── BasicRankingTest.cpp
        │       ├── CashbackExtensionTest.cpp
        │       └── EasyRefundExtensionTest.cpp
        ├── solution/
        │   ├── cpp/
        │   └── DESIGN.md
        └── AI_REVIEW_PROMPT.md
```

### Why `dashboard/` is a subfolder (not root level)

The repo's primary identity is the problem content. The dashboard is a tool *for* the content. Keeping it in a subfolder means:
- The repo root stays clean (problems/, patterns/, README.md)
- Users who don't want the dashboard can ignore it entirely
- The dashboard has its own package.json and dependencies, isolated from the repo

### Root `package.json` (thin wrapper)

```json
{
  "name": "dsa-meets-design",
  "version": "1.0.0",
  "description": "The DSA + Design Pattern Integration You Were Never Taught",
  "scripts": {
    "install": "cd dashboard && npm install",
    "build": "cd dashboard && npm run build",
    "start": "cd dashboard && npm start",
    "dev": "cd dashboard && npm run dev"
  }
}
```

User runs `npm install` at root → installs dashboard dependencies.
User runs `npm start` at root → starts the dashboard.

---

## 6. Data Model

### 6.1 Source: `docs/_data/problems.yml` (read-only)

The dashboard reads this file on startup. Same file used by GitHub Pages. Adding a new problem to this file automatically makes it appear in the dashboard.

```yaml
- id: "001-payment-ranker"
  name: "Payment Method Ranker"
  tier: 1
  patterns: ["Strategy", "Comparator"]
  dsa: "Sorting"
  companies: ["Amazon", "Flipkart"]
  time_minutes: 45
  languages: ["cpp"]
  prerequisite_primer: "strategy"
  path: "problems/tier1-foundation/001-payment-ranker"

- id: "004-vending-machine"
  name: "Vending Machine"
  tier: 1
  patterns: ["State"]
  dsa: "HashMap"
  companies: ["Amazon", "Flipkart"]
  time_minutes: 45
  languages: ["cpp"]
  prerequisite_primer: "state"
  path: "problems/tier1-foundation/004-vending-machine"
```

### 6.2 Progress: `progress.json` (read-write, gitignored)

Auto-created on first status update. Lives in repo root (not inside dashboard/).

```json
{
  "version": 2,
  "created_at": "2026-04-01T10:00:00Z",
  "updated_at": "2026-04-05T14:30:00Z",
  "problems": {
    "001-payment-ranker": {
      "status": "solved",
      "difficulty_mode": "interview",
      "started_at": "2026-04-02T09:00:00Z",
      "completed_at": "2026-04-02T10:30:00Z",
      "notes": "Designed Strategy pattern from scratch. Nailed it.",
      "code": {
        "interview": "#include <iostream>\n// ... user's full Interview mode code ...",
        "guided": "",
        "learning": ""
      }
    },
    "004-vending-machine": {
      "status": "attempted",
      "difficulty_mode": "guided",
      "started_at": "2026-04-03T11:00:00Z",
      "completed_at": null,
      "notes": "",
      "code": {
        "interview": "",
        "guided": "#include <iostream>\n// ... user's Guided mode code ...",
        "learning": ""
      }
    }
  },
  "primers_read": ["strategy", "state"]
}
```

**New fields (v2):**
- `difficulty_mode`: the active mode the user last used for this problem. One of `"interview"`, `"guided"`, `"learning"`. Defaults to `"learning"` if not set.
- `code`: object with three keys — one per mode. Stores the user's written code for each mode independently. Empty string means the user hasn't edited that mode yet (load starter code from file).

**Status values:** Problems not in the file are `unsolved`. Explicit values: `"attempted"`, `"solved"`.

**Concurrency:** The file is only written by the Express server. No concurrent writers. Simple read-modify-write is sufficient.

### 6.3 Pattern Primers Metadata

Derived from the `patterns/` directory. The server scans for `*.md` files (excluding README.md) and extracts the filename as the primer name.

```json
[
  { "name": "strategy", "file": "patterns/strategy.md" },
  { "name": "state", "file": "patterns/state.md" },
  { "name": "observer", "file": "patterns/observer.md" }
]
```

---

## 7. API Specification

Base URL: `http://localhost:3000/api`

### `GET /api/problems`

Returns all problems with their progress status merged in.

**Response:**

```json
{
  "problems": [
    {
      "id": "001-payment-ranker",
      "name": "Payment Method Ranker",
      "tier": 1,
      "patterns": ["Strategy", "Comparator"],
      "dsa": "Sorting",
      "companies": ["Amazon", "Flipkart"],
      "time_minutes": 45,
      "languages": ["cpp"],
      "prerequisite_primer": "strategy",
      "path": "problems/tier1-foundation/001-payment-ranker",
      "status": "solved",
      "difficulty_mode": "interview",
      "started_at": "2026-04-02T09:00:00Z",
      "completed_at": "2026-04-02T10:30:00Z",
      "notes": "Used Strategy pattern."
    },
    {
      "id": "003-notification-system",
      "name": "Notification System",
      "tier": 1,
      "patterns": ["Observer"],
      "dsa": "Queue",
      "companies": ["Flipkart"],
      "time_minutes": 45,
      "languages": ["cpp"],
      "prerequisite_primer": "observer",
      "path": "problems/tier1-foundation/003-notification-system",
      "status": "unsolved",
      "difficulty_mode": "learning",
      "started_at": null,
      "completed_at": null,
      "notes": ""
    }
  ],
  "summary": {
    "total": 10,
    "solved": 7,
    "attempted": 1,
    "unsolved": 2
  }
}
```

### `POST /api/problems/:id/status`

Update a problem's status, notes, and/or difficulty mode.

**Request body:**

```json
{
  "status": "attempted",
  "difficulty_mode": "guided",
  "notes": "Working on Extension 1"
}
```

All fields are optional. If only `notes` is sent, status and mode stay unchanged. When status changes to `"attempted"` and `started_at` is null, server sets `started_at` to now. When status changes to `"solved"`, server sets `completed_at` to now.

**Response:** Updated problem object (same shape as in GET /api/problems).

### `GET /api/problems/:id/starter?mode=interview`

Returns the starter code for a specific difficulty mode.

**Query parameter:** `mode` — one of `interview`, `guided`, `learning`. Defaults to `learning`.

**Response:**

```json
{
  "mode": "interview",
  "code": "#include <iostream>\n#include <vector>\n..."
}
```

**Server logic:**
1. Look up the problem's `path` from problems.yml
2. Read `{path}/boilerplate/cpp/{mode}.cpp`
3. Return the file contents as a string
4. If the mode file doesn't exist, fall back to `learning.cpp` and include a `"fallback": true` flag

### `POST /api/problems/:id/code`

Save the user's code for a specific mode.

**Request body:**

```json
{
  "mode": "interview",
  "code": "#include <iostream>\n// user's code..."
}
```

**Server logic:**
1. Read progress.json
2. Set `problems[id].code[mode]` to the submitted code
3. Write progress.json

**Response:**

```json
{
  "saved": true,
  "mode": "interview"
}
```

### `GET /api/problems/:id/code?mode=interview`

Get the user's saved code for a specific mode.

**Query parameter:** `mode` — one of `interview`, `guided`, `learning`.

**Response:**

```json
{
  "mode": "interview",
  "code": "#include <iostream>\n// user's saved code...",
  "is_starter": false
}
```

**Server logic:**
1. Check progress.json for saved code at `problems[id].code[mode]`
2. If saved code exists and is non-empty → return it with `"is_starter": false`
3. If no saved code → read the starter file from disk and return it with `"is_starter": true`

### `GET /api/problems/:id/readme`

Returns the problem's README.md rendered as HTML.

**Response:**

```json
{
  "html": "<h1>Problem: Payment Method Ranker</h1><table>..."
}
```

### `GET /api/problems/:id/design`

Returns the problem's DESIGN.md rendered as HTML.

**Response:**

```json
{
  "html": "<h1>Design Walkthrough: Payment Method Ranker</h1>..."
}
```

### `GET /api/problems/:id/ai-prompt`

Returns the raw markdown content of AI_REVIEW_PROMPT.md (not rendered — user copies it as-is).

**Response:**

```json
{
  "markdown": "# Review My Solution with AI\n\nAfter attempting..."
}
```

### `GET /api/primers`

Returns all pattern primers with read status.

**Response:**

```json
{
  "primers": [
    { "name": "strategy", "read": true },
    { "name": "state", "read": true },
    { "name": "observer", "read": false },
    { "name": "singleton", "read": false },
    { "name": "composite", "read": false },
    { "name": "factory", "read": false },
    { "name": "builder", "read": false }
  ],
  "summary": {
    "total": 7,
    "read": 2
  }
}
```

### `GET /api/primers/:name`

Returns a primer's markdown rendered as HTML.

**Response:**

```json
{
  "name": "strategy",
  "html": "<h1>Strategy Pattern</h1><h2>The one-line explanation</h2>...",
  "read": true
}
```

### `POST /api/primers/:name/read`

Marks a primer as read.

**Request body:** None (empty POST).

**Response:**

```json
{
  "name": "strategy",
  "read": true
}
```

### `GET /api/stats`

Returns aggregated statistics for the stats page. Now includes difficulty mode distribution.

**Response:**

```json
{
  "overall": {
    "total": 10,
    "solved": 7,
    "attempted": 1,
    "unsolved": 2,
    "percent_complete": 70
  },
  "by_tier": {
    "1": { "total": 10, "solved": 7, "attempted": 1, "unsolved": 2 },
    "2": { "total": 0, "solved": 0, "attempted": 0, "unsolved": 0 },
    "3": { "total": 0, "solved": 0, "attempted": 0, "unsolved": 0 }
  },
  "by_pattern": {
    "Strategy": { "solved": 4, "attempted": 0, "total": 5 },
    "Observer": { "solved": 2, "attempted": 1, "total": 3 },
    "State": { "solved": 1, "attempted": 0, "total": 2 }
  },
  "by_difficulty_mode": {
    "interview": { "solved": 2, "attempted": 0 },
    "guided": { "solved": 3, "attempted": 1 },
    "learning": { "solved": 2, "attempted": 0 }
  },
  "primers": {
    "total": 7,
    "read": 2
  },
  "streak": {
    "current_days": 5,
    "last_activity": "2026-04-05"
  }
}
```

---

## 8. Frontend — Pages & Components

### Route Structure

| Route | Page | Component |
|-------|------|-----------|
| `/` | Dashboard (home) | `Dashboard.jsx` |
| `/problem/:id` | Problem view | `ProblemView.jsx` |
| `/primer/:name` | Pattern primer view | `PrimerView.jsx` |
| `/stats` | Progress stats | `Stats.jsx` |

### Component Tree

```
App.jsx
├── Navbar (always visible)
│   ├── Logo + repo name
│   ├── Nav links: Dashboard | Stats
│   └── Overall progress indicator (small)
│
├── Dashboard.jsx
│   ├── ProgressRing (circular, X/Y solved)
│   ├── StatCard × 4 (Tier 1, Tier 2, Tier 3, Primers)
│   ├── FilterBar (tier, pattern, company, status dropdowns)
│   ├── ProblemTable
│   │   └── Row × N
│   │       ├── Status icon (✓ / ○ / ~)
│   │       ├── Problem number + name (clickable → /problem/:id)
│   │       ├── PatternBadge × N
│   │       ├── TierBadge
│   │       ├── DifficultyModeBadge (🔴/🟡/🟢 — shows last used mode)
│   │       └── Company name
│   ├── PrimerList (bottom section)
│   │   └── Primer × 7 (name + read/unread, clickable → /primer/:name)
│   └── CTABanner
│
├── ProblemView.jsx (split-pane layout)
│   ├── LEFT PANE
│   │   ├── Back button (→ /)
│   │   ├── Problem header (name, tier, patterns, companies, time)
│   │   ├── Prerequisite primer link
│   │   ├── Tab bar: [Description] [Notes]
│   │   ├── MarkdownRenderer (README.md content)
│   │   ├── StatusToggle (Unsolved / Attempted / Solved)
│   │   ├── Notes textarea + save button
│   │   └── CTABanner
│   │
│   └── RIGHT PANE
│       ├── DifficultyModeSelector (3 buttons: Interview / Guided / Learning)
│       ├── Tab bar: [Code] [Design 🔒] [AI Prompt]
│       ├── CodeEditor (in-browser editor, loads code based on selected mode)
│       ├── CopyCommand (./run-tests.sh ... cpp)
│       └── Save + Run buttons
│
├── PrimerView.jsx
│   ├── Back button (→ /)
│   ├── MarkdownRenderer (primer content)
│   ├── "Mark as read" button
│   └── Link to related problems
│
└── Stats.jsx
    ├── Overall progress ring (large)
    ├── Tier breakdown bars
    ├── Pattern distribution chart (horizontal bars)
    ├── Difficulty mode distribution (NEW — horizontal bars per mode)
    ├── Primers progress
    └── Streak indicator
```

---

## 9. Visual Design System

### Style: LeetCode structure + Notion aesthetic

**Principles:**
- White/light background, generous whitespace
- Subtle borders (`1px solid #e5e5e5`), no shadows
- System font stack (no custom fonts to load)
- Colored badges for patterns (soft pastel background + dark text from same hue)
- Clean table layout for problems (LeetCode-inspired)
- Circular progress ring (LeetCode-inspired)
- Stat cards with tier breakdowns

### Tailwind Theme Extensions (`tailwind.config.js`)

```js
module.exports = {
  content: ['./src/**/*.{js,jsx}', './index.html'],
  theme: {
    extend: {
      colors: {
        // Base
        surface: '#ffffff',
        'surface-secondary': '#fafafa',
        'surface-tertiary': '#f5f5f5',
        border: '#e5e5e5',
        'border-hover': '#d4d4d4',

        // Text
        'text-primary': '#171717',
        'text-secondary': '#525252',
        'text-tertiary': '#a3a3a3',

        // Status
        'status-solved': '#16a34a',
        'status-attempted': '#f59e0b',
        'status-unsolved': '#d4d4d4',

        // Difficulty modes
        'mode-interview': '#dc2626',
        'mode-guided': '#f59e0b',
        'mode-learning': '#16a34a',

        // Accent
        accent: '#2563eb',
        'accent-light': '#eff6ff',
      }
    }
  }
}
```

### Pattern Badge Colors

Each design pattern gets a consistent color throughout the UI:

```js
// src/lib/constants.js
export const PATTERN_COLORS = {
  'Strategy':   { bg: '#eff6ff', text: '#1e40af' },  // blue
  'Observer':   { bg: '#ecfdf5', text: '#065f46' },  // green
  'State':      { bg: '#fffbeb', text: '#92400e' },  // amber
  'Singleton':  { bg: '#fef2f2', text: '#991b1b' },  // red
  'Factory':    { bg: '#f5f3ff', text: '#5b21b6' },  // purple
  'Composite':  { bg: '#ecfeff', text: '#155e75' },  // cyan
  'Builder':    { bg: '#fdf4ff', text: '#86198f' },  // pink
  'Comparator': { bg: '#f0fdf4', text: '#166534' },  // emerald
  'Decorator':  { bg: '#fff7ed', text: '#9a3412' },  // orange
  'Command':    { bg: '#faf5ff', text: '#6b21a8' },  // violet
};

export const TIER_COLORS = {
  1: { bg: '#ecfdf5', text: '#065f46', label: 'T1' },  // green — foundation
  2: { bg: '#fffbeb', text: '#92400e', label: 'T2' },  // amber — intermediate
  3: { bg: '#fef2f2', text: '#991b1b', label: 'T3' },  // red — advanced
};

export const STATUS_ICONS = {
  solved:    { icon: '✓', color: '#16a34a' },
  attempted: { icon: '~', color: '#f59e0b' },
  unsolved:  { icon: '○', color: '#d4d4d4' },
};

export const DIFFICULTY_MODES = {
  interview: { 
    label: 'Interview', 
    icon: '🔴', 
    color: '#dc2626', 
    bg: '#fef2f2',
    description: 'Blank slate — design everything from scratch'
  },
  guided: { 
    label: 'Guided', 
    icon: '🟡', 
    color: '#f59e0b', 
    bg: '#fffbeb',
    description: 'Structural hints — you name and connect the pieces'
  },
  learning: { 
    label: 'Learning', 
    icon: '🟢', 
    color: '#16a34a', 
    bg: '#ecfdf5',
    description: 'Full scaffolding — implement the method bodies'
  },
};
```

---

## 10. Page 1: Dashboard (Home)

### Layout (top to bottom)

**Section 1: Header row**
- Left: Progress ring (72px) showing X/Y solved, percentage inside
- Right: Repo name "dsa-meets-design" + subtitle "Your LLD interview prep dashboard"

**Section 2: Stat cards row** (4 cards, equal width, horizontal)
- Tier 1: `X/Y` (green number if >0, gray if 0)
- Tier 2: `X/Y`
- Tier 3: `X/Y`
- Primers: `X/Y`
- Each card: light gray background, small uppercase label, large number below

**Section 3: Filter bar**
- Row of filter controls:
  - Tier: `[All] [Tier 1] [Tier 2] [Tier 3]` — toggle buttons
  - Pattern: dropdown select → `All patterns`, `Strategy`, `Observer`, `State`, etc.
  - Company: dropdown select → `All companies`, `Amazon`, `Flipkart`, etc.
  - Status: dropdown select → `All`, `Solved`, `Attempted`, `Unsolved`
- Filters are combinable (AND logic). All default to "All".

**Section 4: Problem table**

| Column | Width | Content |
|--------|-------|---------|
| Status | 40px | Icon: ✓ (green), ~ (amber), ○ (gray) |
| # | 48px | Problem number (e.g., 001) |
| Problem | flex | Problem name (bold, clickable link to `/problem/:id`) |
| Pattern | auto | Pattern badges (colored pills) |
| Tier | 48px | Tier badge (T1/T2/T3) |
| Mode | 32px | Difficulty mode dot (🔴/🟡/🟢) — shows last used mode. Gray dot if unsolved. |
| Company | 100px | First company name (text, secondary color) |
| Time | 48px | Minutes (e.g., "45m") |

- Rows are clickable (entire row navigates to `/problem/:id`)
- Hover: light background highlight
- Sorted by: problem number (default), or clickable column headers to sort

**Section 5: Pattern primers**
- Horizontal row of primer chips: `strategy ✓ | state ✓ | observer | singleton | composite | factory | builder`
- Read primers: green check, slightly bolder
- Unread primers: gray text
- Each chip is clickable → navigates to `/primer/:name`

**Section 6: CTA banner** (bottom)
- Light gray background, subtle
- "Preparing for a product company interview? Book a mock →" with link

---

## 11. Page 2: Problem View

### Layout: Split-Pane (Left = Description, Right = Code Editor)

The problem view uses a horizontal split-pane layout similar to LeetCode:

```
┌─────────────────────────────┬──────────────────────────────────────────┐
│  ← Problems                 │  Difficulty: [🔴 Interview] [🟡] [🟢]   │
│                             ├──────────────────────────────────────────┤
│  001. Payment Method Ranker │  Code   Design 🔒   AI Prompt  Save Run │
│  T1  Strategy  Comparator   ├──────────────────────────────────────────┤
│  45 min · Amazon, Flipkart  │  1  #include <iostream>                  │
│                             │  2  #include <vector>                    │
│  ⚠ Read strategy primer →  │  3  ...                                  │
│                             │  4                                       │
│  ┌─────────────┬───────┐   │  5  struct PaymentMethod {               │
│  │ Description │ Notes │   │  6      string name;                     │
│  ├─────────────┴───────┤   │  7      ...                              │
│  │                     │   │  8  };                                    │
│  │  Problem Statement  │   │  9                                       │
│  │  ...                │   │  10 // Your design starts here            │
│  │                     │   │  ...                                      │
│  │                     │   │                                           │
│  │                     │   ├──────────────────────────────────────────┤
│  ├─────────────────────┤   │  ./run-tests.sh 001-payment-ranker  [📋] │
│  │ [○ Unsolved] [~ Attempted] [✓ Solved]                             │
│  └─────────────────────┴──────────────────────────────────────────────┘
```

### LEFT PANE — Problem Description

**Header:**
- Back arrow + "Problems" (clickable, navigates to `/`)
- Problem name: large, bold (e.g., "001. Payment Method Ranker")
- Metadata row: Tier badge + Pattern badges + "45 min · Amazon, Flipkart"
- Prerequisite primer: "⚠ Read strategy primer first →" (linked to `/primer/strategy`, shown only if primer not yet read)

**Tab bar:** `[Description]  [Notes]`

**Description tab (default):**
- Full README.md rendered as HTML via API
- Code blocks with syntax highlighting (highlight.js)
- "Before You Code" section is visually distinct (light blue/gray background box)

**Notes tab:**
- Textarea for user notes (auto-saves on blur, debounced)
- Explicit "Save" button also available

**Status toggle (always visible at bottom of left pane):**
```
[○ Unsolved]  [~ Attempted]  [✓ Solved]
```
Three buttons in a row. Active state is filled with the status color. Click to change.

When status changes from `unsolved` to `attempted` → `started_at` is recorded.
When status changes to `solved` → `completed_at` is recorded.
Status can go backward (solved → attempted → unsolved) — this clears the relevant timestamps.

### RIGHT PANE — Code Editor

**Difficulty Mode Selector (top of right pane):**

A row of 3 buttons, always visible above the tab bar:

```
Difficulty:  [🔴 Interview]  [🟡 Guided]  [🟢 Learning]
```

- Exactly one mode is active at a time
- Active button: filled with mode color, white text, slightly larger/bolder
- Inactive buttons: transparent background, gray border, gray text
- Clicking a different mode:
  - If user has edited code → show confirmation dialog: "Your code for [current mode] will be saved. Load [new mode] starter code?"
  - Auto-saves current code before switching
  - Loads saved code for new mode (or starter code if first time in that mode)
- Hovering an inactive button shows a tooltip with the mode description:
  - Interview: "Blank slate — design everything from scratch"
  - Guided: "Structural hints — you name and connect the pieces"
  - Learning: "Full scaffolding — implement the method bodies"

**Tab bar:** `[Code]  [Design 🔒]  [AI Prompt]`  + `[Save]  [▶ Run]` aligned right

**Code tab (default):**
- In-browser code editor (Monaco or similar)
- Loads code based on selected difficulty mode
- Line numbers, syntax highlighting for C++
- Auto-save on blur (debounced 2 seconds) + explicit Save button

**Design tab:**
- Gating depends on difficulty mode (see Section 3.7):
  - 🔴 Interview: locked until "Attempted" or "Solved"
  - 🟡 Guided: locked until "Attempted"
  - 🟢 Learning: always available
- When locked: shows 🔒 icon, clicking shows tooltip "Attempt the problem first"
- When unlocked: renders DESIGN.md as HTML

**AI Prompt tab:**
- Shows raw markdown of AI_REVIEW_PROMPT.md in a read-only textarea
- "Copy to clipboard" button at top right

**Terminal command (bottom of right pane):**
```
┌──────────────────────────────────────────────────┐
│ ./run-tests.sh 001-payment-ranker cpp      [📋]  │
└──────────────────────────────────────────────────┘
```
Monospace font, gray background, copy button on the right. Clicking the icon copies to clipboard and shows a brief "Copied!" tooltip.

**CTA banner (bottom of left pane):**
- "Stuck after Interview mode? Book a walkthrough →" with topmate link
- Only shown when mode is Interview and status is attempted but not solved (high-intent moment)

---

## 12. Page 3: Pattern Primer View

### Layout

**Header:**
- Back arrow + "Back to dashboard"
- Primer name: "Strategy pattern"

**Content:**
- Full primer markdown rendered as HTML
- Code blocks with syntax highlighting
- External links (GFG, Refactoring Guru) open in new tab

**Actions:**
- "Mark as read" button (green, becomes "Read ✓" after clicking)
- "Practice this pattern" section: links to all Tier 1 problems that use this pattern (derived from `prerequisite_primer` field in problems.yml)

---

## 13. Page 4: Stats View

### Layout

**Section 1: Overall progress**
- Large progress ring (120px), percentage inside, X/Y below
- "You've solved X of Y problems"

**Section 2: Tier breakdown**
Three horizontal progress bars, stacked:
```
Tier 1  ████████████░░░░░░  7/10  (70%)
Tier 2  ░░░░░░░░░░░░░░░░░░  0/20  (0%)
Tier 3  ░░░░░░░░░░░░░░░░░░  0/15  (0%)
```
Green fill for solved, amber for attempted, gray for unsolved.

**Section 3: Patterns practiced**
Horizontal bar chart showing how many problems solved per pattern:
```
Strategy  ██████  4
Observer  ███     2
State     ██      1
Composite ░       0
```

**Section 4: Difficulty mode distribution** ← NEW
Shows how many problems were solved in each mode:
```
Your solve distribution:

🔴 Interview    ████      2 solved
🟡 Guided       ██████    3 solved
🟢 Learning     ████      2 solved
```

Each bar uses the mode's color. Below the chart:
- If majority is Learning: "Try solving your next problem in Guided mode to level up!"
- If majority is Guided: "Ready for the real thing? Switch to Interview mode."
- If majority is Interview: "Interview-ready — you're designing from scratch consistently."

This nudge encourages users to progressively challenge themselves.

**Section 5: Primers**
```
Primers read: 2/7
strategy ✓  state ✓  observer  singleton  composite  factory  builder
```

**Section 6: Activity**
```
Current streak: 5 days
Last activity: April 5, 2026
```

Streak is calculated from `completed_at` and `started_at` timestamps. A "day" counts if any problem was started or completed on that date.

---

## 14. Startup & Build Configuration

### `dashboard/package.json`

```json
{
  "name": "dsa-meets-design-dashboard",
  "version": "2.0.0",
  "private": true,
  "scripts": {
    "dev": "concurrently \"node server.js\" \"vite\"",
    "build": "vite build",
    "start": "node server.js",
    "postinstall": "npm run build"
  },
  "dependencies": {
    "express": "^4.18.0",
    "marked": "^12.0.0",
    "marked-highlight": "^2.1.0",
    "highlight.js": "^11.9.0",
    "js-yaml": "^4.1.0",
    "open": "^10.0.0"
  },
  "devDependencies": {
    "react": "^18.2.0",
    "react-dom": "^18.2.0",
    "react-router-dom": "^6.22.0",
    "@vitejs/plugin-react": "^4.2.0",
    "vite": "^5.1.0",
    "tailwindcss": "^3.4.0",
    "postcss": "^8.4.0",
    "autoprefixer": "^10.4.0",
    "concurrently": "^8.2.0"
  }
}
```

### `dashboard/vite.config.js`

```js
import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';

export default defineConfig({
  plugins: [react()],
  build: {
    outDir: 'dist',
  },
  server: {
    port: 5173,
    proxy: {
      '/api': 'http://localhost:3000',
    },
  },
});
```

### `dashboard/server.js` — Startup Logic

```js
// Pseudocode for the server startup:

1. Load problems.yml from ../docs/_data/problems.yml
2. Load progress.json from ../progress.json (create empty v2 if not exists)
   - If progress.json exists with "version": 1, auto-migrate to v2 format
     (add empty "code" and "difficulty_mode" fields to each problem)
3. Scan ../patterns/*.md for primer files
4. Set up Express routes (Section 7)
5. If dist/ folder exists → serve static files from dist/
6. If dist/ folder doesn't exist → log warning: "Run 'npm run build' first"
7. Start server on port 3000 (configurable via PORT env var)
8. Auto-open browser to http://localhost:3000
9. Log: "Dashboard running at http://localhost:3000"
```

### What `npm start` does (from repo root)

```
User runs: npm start
  → root package.json: "start": "cd dashboard && npm start"
    → dashboard package.json: "start": "node server.js"
      → Express starts, serves built React app + API
      → Browser opens automatically
```

### What `npm install` does (from repo root)

```
User runs: npm install
  → root package.json: "install": "cd dashboard && npm install"
    → dashboard dependencies installed
    → postinstall: "npm run build" → Vite builds React app to dist/
```

So `git clone → npm install → npm start` just works. The build happens during install.

---

## 15. Error Handling & Edge Cases

### problems.yml doesn't exist or is malformed
- Server logs clear error: "Could not find docs/_data/problems.yml — make sure you're running from the repo root."
- Dashboard shows: "No problems found. Check that docs/_data/problems.yml exists."

### progress.json doesn't exist
- Server creates it with empty v2 default: `{ "version": 2, "problems": {}, "primers_read": [] }`
- Dashboard shows all problems as unsolved (correct behavior)

### progress.json is v1 format (migration)
- Server detects `"version": 1`
- Auto-migrates: adds `"difficulty_mode": "learning"` and `"code": { "interview": "", "guided": "", "learning": "" }` to each existing problem entry
- Updates version to 2
- Writes migrated file back

### progress.json is corrupted
- Server catches JSON parse error, backs up the corrupted file as `progress.json.bak`, creates fresh empty one
- Logs: "progress.json was corrupted. Backed up to progress.json.bak and created a fresh one."

### Starter file for a mode doesn't exist
- Example: `guided.cpp` doesn't exist for a problem (author hasn't written it yet)
- API falls back to `learning.cpp` and returns `"fallback": true`
- Frontend shows a subtle info message: "Guided mode not available for this problem. Showing Learning mode."
- Interview Mode MUST always exist — if `interview.cpp` is missing, serve an empty file with just the data model comment: `// Interview mode starter not yet available for this problem.`

### Problem README.md or DESIGN.md doesn't exist
- API returns 404 with `{ "error": "README.md not found for problem 001-payment-ranker" }`
- Frontend shows: "This problem's README hasn't been created yet."

### Port 3000 already in use
- Server catches EADDRINUSE error
- Logs: "Port 3000 is in use. Try: PORT=3001 npm start"
- Tries port 3001 automatically as fallback

### User opens dashboard before `npm install`
- The `dist/` folder won't exist. Express should log: "Frontend not built. Run 'npm install' first."

### Problems in progress.json that no longer exist in problems.yml
- Silently ignored. Progress for deleted problems stays in the file but doesn't appear in the UI.

### Code in progress.json exceeds reasonable size
- No hard limit in MVP. But if progress.json grows beyond 10MB, server logs a warning: "progress.json is getting large. Consider using 'Export Progress' to back up and reset."

---

## 16. Future P1 Extensions

These are NOT in the MVP. Documented here so the architecture doesn't block them.

| Feature | What Changes | Effort |
|---------|-------------|--------|
| **Test runner button** | New API endpoint `POST /api/problems/:id/run-tests` that executes `run-tests.sh`, streams output via Server-Sent Events. New `TestOutput` component in frontend. | 4-5 days |
| **Timer mode** | Frontend-only. Countdown timer component on ProblemView. Configurable 45/60/90 min. | 1 day |
| **Dark mode** | Tailwind dark mode classes. Toggle button in Navbar. Preference saved in localStorage. | 1 day |
| **Export progress** | `GET /api/progress/export` returns progress as downloadable JSON or CSV. | 0.5 day |
| **Test results display** | After test runner exists: show pass/fail results inline in ProblemView with green/red test names. | 1 day |
| **Search** | Frontend-only filter. Text input in FilterBar that fuzzy-matches problem names, patterns, companies. | 0.5 day |
| **Progressive hints (Approach 2 hybrid)** | For Interview Mode only: a Hints panel with 3-4 progressively revealed nudges. Tracked in progress.json as `hints_used: 0-4`. Affects Stats page scoring display. | 2-3 days |
| **Re-attempt tracking** | Allow user to solve the same problem in multiple modes (e.g., solve in Learning, then re-attempt in Interview). Track per-mode completion separately. | 1 day |
| **Mode recommendation** | Based on user's overall stats, suggest a default mode for new problems. E.g., "You've solved 10 in Learning mode — try Guided next." | 0.5 day |

### Architecture decisions that enable these later:

- **API-first:** All data flows through the Express API. The frontend never reads files directly. This means adding test execution is just a new endpoint.
- **Component-based:** Each UI element is a separate React component. Adding a timer or editor is adding a component to ProblemView, not rewriting the page.
- **progress.json is versioned:** The `"version": 2` field lets us migrate the schema later without breaking existing progress files.
- **Per-mode code storage:** The `code` object in progress.json stores code separately per mode, enabling future features like re-attempt tracking and mode-specific scoring.
- **Behavior-based tests:** Tests validate output, not class names. This means all three modes can use the same test suite without modification.

---

*End of spec. Hand this to Claude Code to build. Expected output: a working dashboard accessible via `npm install && npm start` with the difficulty mode selector functional on all problems.*
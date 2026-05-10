# DSA Meets Design
### The LLD interview gap nobody talks about — and the repo that fixes it.

---

> 🚀 **Looking for the full library?** **[CodeJunction Pro →](https://topmate.io/jatin_kaushal24/2053177)** ships **100+ machine coding + LLD problems** with progressive extensions, dual-view editorials, 9 company tracks (Amazon, Flipkart, Razorpay, PhonePe, Meesho, Swiggy, Uber, Ola, Microsoft), and a 6-week prep playbook. This pilot has **20** of them.

---

> *"He could solve any LeetCode hard. He had no idea what hit him when the interviewer asked him to make his solution extensible."*
>
> That's a real Amazon LLD round. That candidate was prepared — just not for this.

![Demo — solving Problem 001 on the dashboard](assets/Demo-DSA-Meet-Design.gif)

---

## What's the gap?

Product companies at the 1–3 year level don't just test algorithms anymore. They give you a real-world scenario, ask you to design it, then add requirements mid-interview to see if your design *survives*.

Most people fail not because they can't code — but because nobody ever taught them **how DSA lives inside good design**, and what it looks like when it breaks.

There's no LeetCode for this. GitHub repos are static code dumps. YouTube is passive. Codezym is Java-only.

This repo is different.

---

## How it works

Each problem is structured exactly like a real LLD interview round:

```
Part 1  →  Base requirement (you solve it)
Part 2  →  Extension 1 unlocks only after Part 1 tests pass
Part 3  →  Extension 2 unlocks only after Part 2 tests pass
```

No honor system. You earn the next part.

When a new part unlocks, your design either survives the new requirement — or exposes exactly why it doesn't. That's the learning.

---

## The local dashboard

Clone the repo, run two commands, and you get a LeetCode-style interface on your machine:

```bash
npm install
npm run dev
# → opens at http://localhost:5173
```

From the dashboard you can:
- Browse all problems by tier, pattern, and company
- Pick your difficulty mode per problem (Learning / Guided / Interview)
- Write C++ or Go directly in the browser editor
- Run tests locally and see pass/fail output inline
- Read Design Pattern Primers before attempting problems
- Track your progress across sessions

---

## Three difficulty modes — same problem, different starting point

| Mode | What you get | Who it's for |
|------|-------------|--------------|
| **Interview** | Blank slate — just the problem statement and data types | You know the pattern, you want real practice |
| **Guided** | Key interfaces defined, `// HINT:` comments on purpose — no pattern names | You're learning but want to think |
| **Learning** | Full class structure, `// TODO:` inside method bodies only | You want to understand before you try |

You can switch modes at any point. Most people start in Learning, attempt Guided, then test themselves in Interview mode before their actual interview.

---

## Problems (Pilot — 20 total)

### Tier 1 — Foundation

| # | Problem | Patterns | DSA | Companies |
|---|---------|----------|-----|-----------|
| 001 | Payment Method Ranker | Strategy + Comparator | Sorting | Amazon, Flipkart |
| 003 | Notification System | Observer | Queue | Flipkart, Swiggy |
| 004 | Vending Machine | State | HashMap | Amazon, Flipkart |
| 005 | Customer Issue Resolution | Strategy + Observer | HashMap + Priority Queue | PhonePe, Flipkart |
| 006 | Billing & Discount Engine | Strategy + Decorator | HashMap | Flipkart, Amazon, Meesho |
| 007 | Order Management System | State | HashMap | Meesho, PhonePe, Amazon |
| 008 | File Search System | Strategy + Composite | Tree (DFS/BFS) | Amazon, Microsoft |
| 011 | API Rate Limiter | Strategy + Factory | Queue + HashMap | Amazon, Razorpay, Uber |

### Tier 2 — Intermediate

| # | Problem | Patterns | DSA | Companies |
|---|---------|----------|-----|-----------|
| 009 | Meeting Room Scheduler | Strategy + Observer | Interval Checking + Priority Queue | Flipkart, Razorpay, Groww |
| 010 | Ride Surge Pricing Engine | Strategy + Observer | Priority Queue | Uber, Ola |
| 012 | Elevator System | State + Strategy + Command | Queue + PriorityQueue + HashMap | Adobe |
| 013 | Parking Lot System | Factory + Strategy + Singleton | HashMap + Queue | Salesforce |
| 014 | Splitwise Expense-Sharing | Strategy + Factory + Observer | HashMap + Graph + Greedy | ShareChat, Razorpay, Flipkart, Paytm |
| 015 | BookMyShow Ticket Booking | Strategy + Observer + State | HashMap + Matrix + Queue | DoorDash, BookMyShow, Swiggy, Paytm |
| 016 | Amazon Locker System | Strategy + Factory + State | HashMap + Queue | Amazon |
| 017 | LRU Cache | Singleton + Observer + Strategy | HashMap + Doubly Linked List | Kutumb |
| 018 | Simplified Twitter | Observer + Factory + Singleton | HashMap + HashSet + Heap + LinkedList | AngelOne |
| 019 | Online Auction System | Strategy + Observer + Factory + State | HashMap + PriorityQueue + Sorting | Flipkart |
| 020 | Logger System | Strategy + Observer + Factory + Singleton | HashMap + String Parsing + Queue | Amazon |
| 021 | Ride-Sharing Application | Strategy + Factory | HashMap + Graph + BFS | Flipkart |

---

## Design Pattern Primers

Before you attempt a problem, read the primer for its pattern. These aren't Wikipedia entries — they're written specifically for interview context: what the pattern is, why it fits certain problems, what breaks without it, and C++ examples you can reason about.

Primers available: **Strategy · Observer · State · Singleton**

Examples in both C++ and Go.

---

## Prerequisites

- Node.js 18+
- g++ with C++17 support (`g++ --version` to verify)
- Go 1.21+ for Go problems (`go version` to verify)

---

## Getting Started

```bash
# Clone the repo
git clone https://github.com/jkaus324/DSA-Meet-Design-Pilot.git
cd DSA-Meet-Design-Pilot

# Install dashboard dependencies
cd dashboard && npm install && cd ..

# Start the development server
npm run dev
```

Dashboard → `http://localhost:5173`  
API server → `http://localhost:3000`

---

## Project Structure

```
DSA-Meet-Design-Pilot/
├── problems/
│   ├── tier1-foundation/
│   │   ├── 001-payment-ranker/
│   │   │   ├── README.md              # Problem statement (parts format)
│   │   │   ├── DESIGN.md              # Why this pattern, what breaks without it
│   │   │   ├── boilerplate/cpp/       # interview / guided / learning × per part
│   │   │   ├── boilerplate/go/        # same three modes in Go
│   │   │   ├── tests/cpp/             # GoogleTest suites per part
│   │   │   └── tests/go/              # Go test runners per part
│   │   └── ...
│   └── tier2-intermediate/
├── patterns/                          # Design pattern primers
├── docs/_data/problems.yml            # Problem registry
├── dashboard/                         # React + Express dashboard
│   ├── server.js
│   └── src/
├── e2e/                               # Plain-English Playwright stories
├── scripts/                           # e2e-up.{ps1,sh} bootstrap helpers
└── progress.json                      # Your local progress (gitignored)
```

---

## Running the E2E suite

The `e2e/` folder contains plain-English user stories an LLM-driven browser agent (Playwright MCP) can execute against the running dashboard. Useful before tagging a release or after a UI change:

```bash
# Start the dashboard with an isolated progress file (won't touch your real progress)
pwsh scripts/e2e-up.ps1 -Detach        # Windows
./scripts/e2e-up.sh --detach           # macOS/Linux

# Then in a Claude Code session, ask the agent:
#   Run e2e/stories/_smoke.md against http://localhost:3000.
```

See `e2e/README.md` for the full suite, priorities, and authoring guide.

---

## Who this is for

**Primary:** Developers with 1–3 years experience at service companies (TCS, Infosys, Wipro, Cognizant) preparing to switch to product companies — Amazon, Flipkart, Razorpay, Meesho, PhonePe, Paytm, Groww.

You've done DSA. You haven't done LLD. This is the missing piece.

**Secondary:** Indian developers targeting Canadian product companies (Shopify, Wealthsimple, Coveo). The machine-coding format you've been practicing is the wrong format for Canadian interviews. This repo covers both.

---

## This is a pilot

20 problems. Real test suites. A working dashboard. Enough to get meaningful feedback before the full launch.

If you find a bug, open an issue. If a DESIGN.md explanation didn't click, open an issue and say what was confusing. That feedback directly shapes what gets built next.

---

## Want the full 100-problem library?

This pilot has 20 problems. The full version — **CodeJunction Pro** — ships 100+ machine coding + LLD problems pulled from real interviews at Amazon, Flipkart, Razorpay, PhonePe, Meesho, Swiggy, Uber, Ola, Microsoft, and Google.

**What you also get on top of the pilot:**
- 5x more problems, tagged by company and difficulty
- **Dual-view editorials** for every problem — LLD perspective + machine-coding perspective + the *Junction* analysis (where one wrong design choice in minute 5 kills the algorithm in minute 25)
- **9 company-specific tracks** with curated problem sequences in interview order
- **6-week prep playbook** with a structured weekly cadence
- **10 design pattern primers** (vs 4 in the pilot)
- Solutions in **C++ and Java**
- Cheat sheets and AI review prompts you can paste straight into Claude/ChatGPT

**Get it → [CodeJunction Pro on Topmate](https://topmate.io/jatin_kaushal24/2053177)** · ₹799 one-time, GitHub access within 10 minutes of payment.

---

## Want a 1:1 walkthrough or full mock interview?

The repo gives you the problems and the framework. If you want a tailored walkthrough for your specific target company — or a full mock interview with real-time feedback and a Hire/No Hire verdict:

**Book a session → [topmate.io/jatin_kaushal24](https://topmate.io/jatin_kaushal24)**

---

*Built by [Jatin Kaushal](https://www.linkedin.com/in/jatin-kaushal-0324/) — SDE at Amazon India.*  
*Shared because nobody gave me this when I needed it.*
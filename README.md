# DSA Meets Design
### The LLD interview gap nobody talks about — and the repo that fixes it.

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
└── progress.json                      # Your local progress (gitignored)
```

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

## Preparing for an interview soon?

The public repo gives you the problems and the framework. If you want a walkthrough of how to use it effectively for your specific target company — or a full mock interview with real-time feedback and a Hire/No Hire verdict:

**Book a session → [topmate.io/jatin_kaushal24](https://topmate.io/jatin_kaushal24)**

---

*Built by [Jatin Kaushal](https://www.linkedin.com/in/jatin-kaushal-0324/) — SDE at Amazon India.*  
*Shared because nobody gave me this when I needed it.*
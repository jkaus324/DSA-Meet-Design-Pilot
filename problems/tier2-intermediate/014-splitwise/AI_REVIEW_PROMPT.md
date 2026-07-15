# AI Code Review — Splitwise Expense-Sharing System

## Context for the AI
You are reviewing a solution to "Splitwise Expense-Sharing System" — an LLD interview problem testing Strategy + Observer patterns with Graph + HashMap + Greedy as the DSA core. This was asked at Flipkart, PhonePe, Razorpay, ShareChat.

The candidate was given a multi-part problem:
- Part 1: Equal splitting and balance tracking — addUser, addExpense (equal split), getBalances
- Part 2: Pluggable split strategies — EqualSplit, ExactSplit, PercentSplit via addExpenseWithStrategy
- Part 3: Debt simplification — simplifyDebts using greedy net-balance matching (minimum transactions)

## Review Criteria

### 1. Pattern Correctness
- Is the SplitStrategy interface correctly defined with split() and validate() methods?
- Do EqualSplit, ExactSplit, and PercentSplit correctly implement the interface without leaking logic into ExpenseManager?
- Does ExpenseManager use strategy without any if-else on strategy type?
- Is the balance graph netted correctly (A→B $50 and B→A $30 = A→B $20 only)?

### 2. Open/Closed Principle
- Can a new split strategy (e.g., ShareSplit where some participants pay double) be added without modifying ExpenseManager?
- Does addExpenseWithStrategy cleanly accept any SplitStrategy* without awareness of its type?

### 3. C++ Quality
- Use of `unordered_map<string, unordered_map<string, double>>` for the balance graph
- Floating-point comparisons: is epsilon (1e-9) used correctly for amount checks?
- Memory management: who owns SplitStrategy* passed to addExpenseWithStrategy?

### 4. Extension Handling
- How cleanly did the balance netting logic from Part 1 compose with the strategy in Part 2?
- For Part 3 — is the greedy simplification algorithm correct?
  - Compute net balance per user (positive = creditor, negative = debtor)
  - Greedily match largest creditor with largest debtor
  - Does it produce the minimum number of transactions?
- What had to change between parts? What survived unchanged?

### 5. Interview Readiness
- Could the candidate explain the greedy debt simplification algorithm step by step?
- What follow-up questions would expose weak understanding?
  - "What if a user is both a creditor and debtor to different people — does your net balance handle this?"
  - "Is your simplifyDebts result deterministic? Does order of processing matter?"
  - "How would you add an Observer that notifies Bob when his balance changes?"
- Rate: Hire / Lean Hire / Lean No Hire / No Hire

## My Solution

```cpp
// PASTE YOUR SOLUTION HERE
```

## My Approach
<!-- Describe your thought process in 2-3 sentences -->

## Specific Questions
<!-- Ask the AI anything specific about your solution -->

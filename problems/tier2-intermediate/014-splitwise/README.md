# Problem 014 — Splitwise Expense-Sharing System

**Tier:** 2 (Intermediate) | **Pattern:** Strategy + Observer | **DSA:** Graph + HashMap + Greedy
**Companies:** Flipkart, PhonePe, Razorpay, ShareChat | **Time:** 60 minutes

---

## Problem Statement

You are building an expense-sharing system. One user pays a group expense; the system splits the cost among participants and tracks who owes whom. The system supports multiple split strategies and can simplify a complex web of debts into the minimum number of transactions needed to settle up.

**Constraints:**
- Up to 10^3 users and 10^4 expenses
- Balances are netted: if A owes B $50 and B owes A $30, only A owes B $20 is stored
- Floating-point amounts; use double precision
- Debt simplification is greedy and minimizes transaction count, not total amount

---

## Base Requirement — Equal Splitting and Balance Tracking

Implement an `ExpenseManager` that adds users and records expenses. When an expense is added, the total is split equally among all participants (including the payer). The payer's share is subtracted from what others owe them.

**Example:**
```
addUser("alice"), addUser("bob"), addUser("charlie")
addExpense("E1", paidBy="alice", amount=300.0, participants=["alice","bob","charlie"])
// Each owes $100. Alice paid, so: bob→alice $100, charlie→alice $100

getBalances()
→  { "bob": {"alice": 100.0}, "charlie": {"alice": 100.0} }

// New expense: bob pays $60 for bob and alice
addExpense("E2", paidBy="bob", amount=60.0, participants=["bob","alice"])
// alice→bob $30, but alice was already owed $100 by bob
// Net: bob→alice $70

getBalances()["bob"]["alice"]  →  70.0
```

**Public methods:**
- `void addUser(const string& userId, const string& name)`
- `void addExpense(const string& expenseId, const string& paidBy, double amount, const vector<string>& participants)`
- `unordered_map<string, unordered_map<string, double>> getBalances()`

---

## Extension 1 — Multiple Split Strategies

Add pluggable split strategies. Adding a new strategy must require zero changes to `ExpenseManager`.

| Strategy | Rule |
|---|---|
| EqualSplit | Divide total equally; params ignored |
| ExactSplit | Each participant's share is given explicitly; amounts must sum to total |
| PercentSplit | Each participant's share as a percentage of total; percentages must sum to 100 |

**Example:**
```
// PercentSplit: alice 50%, bob 30%, charlie 20% of $200
addExpenseWithStrategy("E3", paidBy="alice", amount=200.0,
    participants=["alice","bob","charlie"],
    strategy=new PercentSplit(), params=[50.0, 30.0, 20.0])
// bob→alice $60, charlie→alice $40
```

**Public method:**
- `void addExpenseWithStrategy(const string& expenseId, const string& paidBy, double amount, const vector<string>& participants, SplitStrategy* strategy, const vector<double>& params)`

---

## Extension 2 — Debt Simplification

After many expenses, the balance graph may have many edges. Simplify all debts to the minimum number of transactions needed to settle everyone to zero.

**Algorithm:** Compute each user's net balance (total owed to them minus total they owe). Greedily match the largest creditor with the largest debtor. Transaction amount = min(|creditor net|, |debtor net|). Repeat until all balances are zero.

**Example:**
```
// Net balances: alice=+50, bob=-30, charlie=-20
simplifyDebts()
→  [("bob", "alice", 30.0), ("charlie", "alice", 20.0)]
// 2 transactions settle everyone — no matter how many original expenses existed
```

**Public method:**
- `vector<tuple<string, string, double>> simplifyDebts()`

---

## Running Tests

```bash
./run-tests.sh 014-splitwise cpp
```

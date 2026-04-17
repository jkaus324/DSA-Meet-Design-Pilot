package main

// --- Data Model (given -- do not modify) ------------------------------------

type User struct {
	ID   string
	Name string
}

type Split struct {
	UserID string
	Amount float64
}

type Expense struct {
	ID          string
	PaidBy      string
	TotalAmount float64
	Splits      []Split
}

// --- Split Strategy Interface (from Part 2) ----------------------------------

// type SplitStrategy interface {
//     SplitAmounts(total float64, participants []string, params []float64) []Split
//     Validate(total float64, participants []string, params []float64) bool
// }

// type EqualSplit struct{}
// type ExactSplit struct{}
// type PercentSplit struct{}

// --- Transaction (result of simplification) ----------------------------------
// HINT: Represent each settled debt as a struct with Debtor, Creditor, Amount.

// type Transaction struct { Debtor, Creditor string; Amount float64 }

// --- Expense Manager (extends Part 2) ----------------------------------------
// HINT: SimplifyDebts works in three steps:
//   1. Compute net[user] = sum of amounts owed TO user minus sum owed BY user.
//      Iterate m.balances: for each (debtor, creditor, amt) add amt to net[creditor]
//      and subtract from net[debtor].
//   2. Split into creditors (net > 0) and debtors (net < 0). Sort both descending.
//   3. Two-pointer greedy: settle min(creditor.amount, debtor.amount) at a time,
//      advance whichever pointer reaches zero.

// func (m *ExpenseManager) SimplifyDebts() []Transaction

// --- Global Entry Points (required by tests) --------------------------------

// func AddUser(userID, name string)
// func AddExpense(expenseID, paidBy string, amount float64, participants []string)
// func AddExpenseWithStrategy(expenseID, paidBy string, amount float64,
//     participants []string, strategy SplitStrategy, params []float64)
// func GetBalances() map[string]map[string]float64
// func SimplifyDebts() []Transaction

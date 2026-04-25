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

// --- Split Strategy Interface -----------------------------------------------
// HINT: Define a Go interface with two methods:
//         SplitAmounts(total float64, participants []string, params []float64) []Split
//         Validate(total float64, participants []string, params []float64) bool

// type SplitStrategy interface { ... }

// --- Concrete Strategies ----------------------------------------------------
// HINT: EqualSplit ignores params entirely; share = total / len(participants).
// HINT: ExactSplit uses params[i] as the amount for participants[i];
//       validate that sum(params) is within 1e-9 of total.
// HINT: PercentSplit uses params[i]/100 * total for participants[i];
//       validate that sum(params) is within 1e-9 of 100.

// type EqualSplit struct{}
// type ExactSplit struct{}
// type PercentSplit struct{}

// --- Expense Manager (extends Part 1) ----------------------------------------
// HINT: Add AddExpenseWithStrategy on top of Part 1's ExpenseManager.
//       Call strategy.Validate first; skip the expense if invalid.

// type ExpenseManager struct { ... }

// func (m *ExpenseManager) AddExpenseWithStrategy(expenseID, paidBy string, amount float64,
//     participants []string, strategy SplitStrategy, params []float64)

// --- Global Entry Points (required by tests) --------------------------------

// func AddUser(userID, name string)
// func AddExpense(expenseID, paidBy string, amount float64, participants []string)
// func AddExpenseWithStrategy(expenseID, paidBy string, amount float64,
//     participants []string, strategy SplitStrategy, params []float64)
// func GetBalances() map[string]map[string]float64

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

// --- Your Design Starts Here ------------------------------------------------
//
// Design and implement an ExpenseManager that:
//   1. Registers users
//   2. Adds an expense paid by one user, split equally among participants
//   3. Returns the current balance map: who owes whom how much
//
// Think about:
//   - How do you track net balances between pairs of users?
//   - How do you avoid double-counting when A owes B and B owes A?
//
// Entry points (must exist for tests):
//   func AddUser(userID, name string)
//   func AddExpense(expenseID, paidBy string, amount float64, participants []string)
//   func GetBalances() map[string]map[string]float64

// -------------------------------------------------------------------------

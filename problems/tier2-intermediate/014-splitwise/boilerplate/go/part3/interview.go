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

// --- Your Design Starts Here (Part 3) ---------------------------------------
//
// Extend Parts 1 & 2 to add debt simplification:
//   - SimplifyDebts() returns a minimal list of transactions that settles
//     all balances, represented as (debtor, creditor, amount) tuples.
//
// Think about:
//   - How do you compute each user's net balance (positive = owed money,
//     negative = owes money)?
//   - How do you greedily match the largest creditor with the largest debtor
//     to minimise the number of transactions?
//
// Entry points (must exist for tests — include all prior entry points too):
//   func AddUser(userID, name string)
//   func AddExpense(expenseID, paidBy string, amount float64, participants []string)
//   func AddExpenseWithStrategy(expenseID, paidBy string, amount float64,
//                               participants []string, strategy SplitStrategy, params []float64)
//   func GetBalances() map[string]map[string]float64
//   func SimplifyDebts() [][3]interface{}   // each element: {debtor, creditor, amount}

// -------------------------------------------------------------------------

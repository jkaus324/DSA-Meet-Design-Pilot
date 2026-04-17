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

// --- Your Design Starts Here (Part 2) ---------------------------------------
//
// Extend Part 1 to support pluggable split strategies:
//   - EqualSplit   : divide totalAmount evenly across all participants
//   - ExactSplit   : each participant's amount is given explicitly in params
//   - PercentSplit : each participant's share is a percentage in params (must sum to 100)
//
// Think about:
//   - What interface should a split strategy expose?
//   - How do you validate ExactSplit (params sum == total) and PercentSplit (params sum == 100)?
//
// Entry points (must exist for tests — include Part 1 entry points too):
//   func AddUser(userID, name string)
//   func AddExpense(expenseID, paidBy string, amount float64, participants []string)
//   func AddExpenseWithStrategy(expenseID, paidBy string, amount float64,
//                               participants []string, strategy SplitStrategy, params []float64)
//   func GetBalances() map[string]map[string]float64

// -------------------------------------------------------------------------

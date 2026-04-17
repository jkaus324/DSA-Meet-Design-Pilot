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

// --- Balance Tracker --------------------------------------------------------
// HINT: Use a nested map: balances[debtor][creditor] = amount.
//       When adding an expense, the payer is the creditor; each participant
//       (excluding the payer) is a debtor for their share.
// HINT: Net out opposing balances: if A owes B 50 and B already owes A 30,
//       reduce B's debt first, then record only the remaining 20 from A to B.

// --- Expense Manager --------------------------------------------------------
// HINT: Store users in map[string]User keyed by userID.
// HINT: Store expenses in []Expense.
// HINT: Provide a private updateBalance(debtor, creditor string, amount float64) helper.

// type ExpenseManager struct {
//     users    map[string]User
//     expenses []Expense
//     balances map[string]map[string]float64
// }

// func NewExpenseManager() *ExpenseManager

// func (m *ExpenseManager) AddUser(userID, name string)

// func (m *ExpenseManager) AddExpense(expenseID, paidBy string, amount float64, participants []string)
// HINT: Split equally: share = amount / len(participants). Call updateBalance for
//       each participant (debtor=participant, creditor=paidBy).

// func (m *ExpenseManager) GetBalances() map[string]map[string]float64

// --- Global Entry Points (required by tests) --------------------------------
// HINT: Use a package-level *ExpenseManager variable. Re-initialise it in
//       each test via a Reset() helper if needed.

// func AddUser(userID, name string)
// func AddExpense(expenseID, paidBy string, amount float64, participants []string)
// func GetBalances() map[string]map[string]float64

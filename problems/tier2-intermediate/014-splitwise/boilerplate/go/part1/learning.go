package main

import "math"

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

// --- Expense Manager --------------------------------------------------------

type ExpenseManager struct {
	users    map[string]User
	expenses []Expense
	balances map[string]map[string]float64
}

func NewExpenseManager() *ExpenseManager {
	// TODO: Initialise and return an ExpenseManager with empty maps/slices
	return nil
}

func (m *ExpenseManager) updateBalance(debtor, creditor string, amount float64) {
	// TODO: Skip if debtor == creditor
	// TODO: If creditor already owes debtor (balances[creditor][debtor] > 0),
	//       offset that amount first, reduce it by min(existing, amount),
	//       then remove the key if it reaches near-zero (< 1e-9).
	// TODO: If amount still > 1e-9 after offsetting, add it to balances[debtor][creditor]
	_ = math.Abs(0) // math used for 1e-9 comparisons
}

func (m *ExpenseManager) AddUser(userID, name string) {
	// TODO: Store User{ID: userID, Name: name} in m.users
}

func (m *ExpenseManager) AddExpense(expenseID, paidBy string, amount float64, participants []string) {
	// TODO: Compute equal share = amount / len(participants)
	// TODO: Build []Split, one per participant
	// TODO: Append Expense to m.expenses
	// TODO: Call updateBalance(participant, paidBy, share) for each participant
}

func (m *ExpenseManager) GetBalances() map[string]map[string]float64 {
	// TODO: Return m.balances
	return nil
}

// --- Global Entry Points (required by tests) --------------------------------

var manager *ExpenseManager

func ResetManager() {
	manager = NewExpenseManager()
}

func AddUser(userID, name string) {
	// TODO: Delegate to manager.AddUser
}

func AddExpense(expenseID, paidBy string, amount float64, participants []string) {
	// TODO: Delegate to manager.AddExpense
}

func GetBalances() map[string]map[string]float64 {
	// TODO: Delegate to manager.GetBalances
	return nil
}

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

// --- Split Strategy Interface -----------------------------------------------

type SplitStrategy interface {
	SplitAmounts(total float64, participants []string, params []float64) []Split
	Validate(total float64, participants []string, params []float64) bool
}

// --- Concrete Strategies ----------------------------------------------------

type EqualSplit struct{}

func (e EqualSplit) SplitAmounts(total float64, participants []string, params []float64) []Split {
	// TODO: Compute share = total / len(participants)
	// TODO: Return a []Split with each participant assigned share
	return nil
}

func (e EqualSplit) Validate(total float64, participants []string, params []float64) bool {
	// TODO: Return true if len(participants) > 0
	return false
}

type ExactSplit struct{}

func (e ExactSplit) SplitAmounts(total float64, participants []string, params []float64) []Split {
	// TODO: Return []Split where splits[i] = {participants[i], params[i]}
	return nil
}

func (e ExactSplit) Validate(total float64, participants []string, params []float64) bool {
	// TODO: Return false if len(params) != len(participants)
	// TODO: Sum params and return math.Abs(sum-total) < 1e-9
	_ = math.Abs(0)
	return false
}

type PercentSplit struct{}

func (p PercentSplit) SplitAmounts(total float64, participants []string, params []float64) []Split {
	// TODO: Return []Split where splits[i] = {participants[i], total * params[i] / 100}
	return nil
}

func (p PercentSplit) Validate(total float64, participants []string, params []float64) bool {
	// TODO: Return false if len(params) != len(participants)
	// TODO: Sum params and return math.Abs(sum-100) < 1e-9
	return false
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
	// TODO: Same as Part 1 — net out opposing balances, then record remainder
}

func (m *ExpenseManager) AddUser(userID, name string) {
	// TODO: Store User in m.users
}

func (m *ExpenseManager) AddExpense(expenseID, paidBy string, amount float64, participants []string) {
	// TODO: Use EqualSplit{} to split, then call updateBalance for each split
}

func (m *ExpenseManager) AddExpenseWithStrategy(expenseID, paidBy string, amount float64,
	participants []string, strategy SplitStrategy, params []float64) {
	// TODO: Call strategy.Validate; return (skip) if invalid
	// TODO: Call strategy.SplitAmounts to get []Split
	// TODO: Append Expense to m.expenses
	// TODO: Call updateBalance for each split entry
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

func AddExpenseWithStrategy(expenseID, paidBy string, amount float64,
	participants []string, strategy SplitStrategy, params []float64) {
	// TODO: Delegate to manager.AddExpenseWithStrategy
}

func GetBalances() map[string]map[string]float64 {
	// TODO: Delegate to manager.GetBalances
	return nil
}

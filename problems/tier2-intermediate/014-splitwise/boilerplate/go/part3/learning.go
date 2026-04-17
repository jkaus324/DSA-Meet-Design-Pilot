package main

import (
	"math"
	"sort"
)

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

type EqualSplit struct{}

func (e EqualSplit) SplitAmounts(total float64, participants []string, params []float64) []Split {
	// TODO: share = total / len(participants); return one Split per participant
	return nil
}
func (e EqualSplit) Validate(total float64, participants []string, params []float64) bool {
	// TODO: len(participants) > 0
	return false
}

type ExactSplit struct{}

func (e ExactSplit) SplitAmounts(total float64, participants []string, params []float64) []Split {
	// TODO: splits[i] = {participants[i], params[i]}
	return nil
}
func (e ExactSplit) Validate(total float64, participants []string, params []float64) bool {
	// TODO: len(params)==len(participants) && math.Abs(sum(params)-total) < 1e-9
	_ = math.Abs(0)
	return false
}

type PercentSplit struct{}

func (p PercentSplit) SplitAmounts(total float64, participants []string, params []float64) []Split {
	// TODO: splits[i] = {participants[i], total * params[i] / 100}
	return nil
}
func (p PercentSplit) Validate(total float64, participants []string, params []float64) bool {
	// TODO: len(params)==len(participants) && math.Abs(sum(params)-100) < 1e-9
	return false
}

// --- Transaction ------------------------------------------------------------

type Transaction struct {
	Debtor   string
	Creditor string
	Amount   float64
}

// --- Expense Manager --------------------------------------------------------

type ExpenseManager struct {
	users    map[string]User
	expenses []Expense
	balances map[string]map[string]float64
}

func NewExpenseManager() *ExpenseManager {
	// TODO: Initialise maps/slices and return
	return nil
}

func (m *ExpenseManager) updateBalance(debtor, creditor string, amount float64) {
	// TODO: Net out m.balances[creditor][debtor] first, then record remainder
}

func (m *ExpenseManager) AddUser(userID, name string) {
	// TODO: m.users[userID] = User{...}
}

func (m *ExpenseManager) AddExpense(expenseID, paidBy string, amount float64, participants []string) {
	// TODO: EqualSplit, build splits, append expense, call updateBalance
}

func (m *ExpenseManager) AddExpenseWithStrategy(expenseID, paidBy string, amount float64,
	participants []string, strategy SplitStrategy, params []float64) {
	// TODO: Validate → SplitAmounts → append expense → updateBalance for each
}

func (m *ExpenseManager) GetBalances() map[string]map[string]float64 {
	// TODO: return m.balances
	return nil
}

func (m *ExpenseManager) SimplifyDebts() []Transaction {
	// Step 1: Compute net balance per user
	// TODO: For each (debtor, creditorMap) in m.balances:
	//         net[creditor] += amount; net[debtor] -= amount

	// Step 2: Separate into creditors (net > 1e-9) and debtors (net < -1e-9)
	// TODO: Build two slices of {user, |amount|} pairs

	// Step 3: Sort both slices descending by amount
	_ = sort.Slice // use sort.Slice

	// Step 4: Two-pointer greedy matching
	// TODO: settle = min(creditor.amount, debtor.amount)
	//       append Transaction{debtor, creditor, settle}
	//       reduce both amounts; advance pointer that reached 0

	return nil
}

// --- Global Entry Points (required by tests) --------------------------------

var manager *ExpenseManager

func ResetManager() {
	manager = NewExpenseManager()
}

func AddUser(userID, name string) {
	// TODO: manager.AddUser(userID, name)
}

func AddExpense(expenseID, paidBy string, amount float64, participants []string) {
	// TODO: manager.AddExpense(...)
}

func AddExpenseWithStrategy(expenseID, paidBy string, amount float64,
	participants []string, strategy SplitStrategy, params []float64) {
	// TODO: manager.AddExpenseWithStrategy(...)
}

func GetBalances() map[string]map[string]float64 {
	// TODO: return manager.GetBalances()
	return nil
}

func SimplifyDebts() []Transaction {
	// TODO: return manager.SimplifyDebts()
	return nil
}

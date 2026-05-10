package main

import (
	"math"
	"sort"
)

// ─── Data Model ──────────────────────────────────────────────────────────────

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

type Transaction struct {
	Debtor   string
	Creditor string
	Amount   float64
}

// ─── Split Strategy Interface ────────────────────────────────────────────────

type SplitStrategy interface {
	SplitAmounts(total float64, participants []string, params []float64) []Split
	Validate(total float64, participants []string, params []float64) bool
}

type EqualSplit struct{}

func (EqualSplit) SplitAmounts(total float64, participants []string, _ []float64) []Split {
	if len(participants) == 0 {
		return nil
	}
	share := total / float64(len(participants))
	out := make([]Split, 0, len(participants))
	for _, p := range participants {
		out = append(out, Split{UserID: p, Amount: share})
	}
	return out
}

func (EqualSplit) Validate(_ float64, participants []string, _ []float64) bool {
	return len(participants) > 0
}

type ExactSplit struct{}

func (ExactSplit) SplitAmounts(_ float64, participants []string, params []float64) []Split {
	out := make([]Split, 0, len(participants))
	for i, p := range participants {
		out = append(out, Split{UserID: p, Amount: params[i]})
	}
	return out
}

func (ExactSplit) Validate(total float64, participants []string, params []float64) bool {
	if len(params) != len(participants) {
		return false
	}
	sum := 0.0
	for _, v := range params {
		sum += v
	}
	return math.Abs(sum-total) < 1e-9
}

type PercentSplit struct{}

func (PercentSplit) SplitAmounts(total float64, participants []string, params []float64) []Split {
	out := make([]Split, 0, len(participants))
	for i, p := range participants {
		out = append(out, Split{UserID: p, Amount: total * params[i] / 100.0})
	}
	return out
}

func (PercentSplit) Validate(_ float64, participants []string, params []float64) bool {
	if len(params) != len(participants) {
		return false
	}
	sum := 0.0
	for _, v := range params {
		sum += v
	}
	return math.Abs(sum-100.0) < 1e-9
}

// ─── Expense Manager ─────────────────────────────────────────────────────────

type ExpenseManager struct {
	users    map[string]User
	expenses []Expense
	balances map[string]map[string]float64
}

func NewExpenseManager() *ExpenseManager {
	return &ExpenseManager{
		users:    make(map[string]User),
		expenses: nil,
		balances: make(map[string]map[string]float64),
	}
}

// updateBalance: debtor owes creditor `amount`. Net out an existing
// reciprocal debt (creditor→debtor) before recording the remainder.
func (m *ExpenseManager) updateBalance(debtor, creditor string, amount float64) {
	if debtor == creditor || amount <= 0 {
		return
	}
	if reciprocal, ok := m.balances[creditor][debtor]; ok && reciprocal > 0 {
		offset := math.Min(reciprocal, amount)
		m.balances[creditor][debtor] -= offset
		amount -= offset
		if m.balances[creditor][debtor] < 1e-9 {
			delete(m.balances[creditor], debtor)
		}
	}
	if amount > 1e-9 {
		if m.balances[debtor] == nil {
			m.balances[debtor] = make(map[string]float64)
		}
		m.balances[debtor][creditor] += amount
	}
}

func (m *ExpenseManager) AddUser(userID, name string) {
	m.users[userID] = User{ID: userID, Name: name}
}

func (m *ExpenseManager) AddExpense(expenseID, paidBy string, amount float64, participants []string) {
	if len(participants) == 0 {
		return
	}
	splits := EqualSplit{}.SplitAmounts(amount, participants, nil)
	m.expenses = append(m.expenses, Expense{ID: expenseID, PaidBy: paidBy, TotalAmount: amount, Splits: splits})
	for _, s := range splits {
		m.updateBalance(s.UserID, paidBy, s.Amount)
	}
}

func (m *ExpenseManager) AddExpenseWithStrategy(expenseID, paidBy string, amount float64,
	participants []string, strategy SplitStrategy, params []float64) {
	if !strategy.Validate(amount, participants, params) {
		return
	}
	splits := strategy.SplitAmounts(amount, participants, params)
	m.expenses = append(m.expenses, Expense{ID: expenseID, PaidBy: paidBy, TotalAmount: amount, Splits: splits})
	for _, s := range splits {
		m.updateBalance(s.UserID, paidBy, s.Amount)
	}
}

func (m *ExpenseManager) GetBalances() map[string]map[string]float64 {
	return m.balances
}

func (m *ExpenseManager) SimplifyDebts() []Transaction {
	// Step 1: net balance per user (positive = creditor, negative = debtor)
	net := make(map[string]float64)
	for debtor, creditors := range m.balances {
		for creditor, amount := range creditors {
			net[debtor] -= amount
			net[creditor] += amount
		}
	}

	// Step 2: split into creditors and debtors (use absolute amounts for debtors)
	type userAmount struct {
		user   string
		amount float64
	}
	var creditors, debtors []userAmount
	for user, amount := range net {
		if amount > 1e-9 {
			creditors = append(creditors, userAmount{user, amount})
		} else if amount < -1e-9 {
			debtors = append(debtors, userAmount{user, -amount})
		}
	}

	// Step 3: sort descending by amount (deterministic tiebreak by name keeps output stable)
	sort.Slice(creditors, func(i, j int) bool {
		if creditors[i].amount != creditors[j].amount {
			return creditors[i].amount > creditors[j].amount
		}
		return creditors[i].user < creditors[j].user
	})
	sort.Slice(debtors, func(i, j int) bool {
		if debtors[i].amount != debtors[j].amount {
			return debtors[i].amount > debtors[j].amount
		}
		return debtors[i].user < debtors[j].user
	})

	// Step 4: greedy two-pointer matching
	var txns []Transaction
	i, j := 0, 0
	for i < len(creditors) && j < len(debtors) {
		settle := math.Min(creditors[i].amount, debtors[j].amount)
		txns = append(txns, Transaction{Debtor: debtors[j].user, Creditor: creditors[i].user, Amount: settle})
		creditors[i].amount -= settle
		debtors[j].amount -= settle
		if creditors[i].amount < 1e-9 {
			i++
		}
		if debtors[j].amount < 1e-9 {
			j++
		}
	}
	return txns
}

// ─── Global Entry Points ─────────────────────────────────────────────────────

var manager *ExpenseManager

func ResetManager() {
	manager = NewExpenseManager()
}

func AddUser(userID, name string) {
	if manager == nil {
		ResetManager()
	}
	manager.AddUser(userID, name)
}

func AddExpense(expenseID, paidBy string, amount float64, participants []string) {
	if manager == nil {
		ResetManager()
	}
	manager.AddExpense(expenseID, paidBy, amount, participants)
}

func AddExpenseWithStrategy(expenseID, paidBy string, amount float64,
	participants []string, strategy SplitStrategy, params []float64) {
	if manager == nil {
		ResetManager()
	}
	manager.AddExpenseWithStrategy(expenseID, paidBy, amount, participants, strategy, params)
}

func GetBalances() map[string]map[string]float64 {
	if manager == nil {
		return nil
	}
	return manager.GetBalances()
}

func SimplifyDebts() []Transaction {
	if manager == nil {
		return nil
	}
	return manager.SimplifyDebts()
}

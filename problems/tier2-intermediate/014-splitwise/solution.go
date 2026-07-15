// Splitwise — equal/exact/percent splits + debt simplification (Go).
package main

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
)

type SplitOp struct {
	kind string
	s1   string
	s2   string
	s3   string
	s4   string
	i1   int
}

type splitPair struct {
	user   string
	amount float64
}

type splitStrategy interface {
	split(totalAmount float64, participants []string, params []float64) []splitPair
	validate(totalAmount float64, participants []string, params []float64) bool
}

type equalSplit struct{}

func (equalSplit) split(totalAmount float64, participants []string, params []float64) []splitPair {
	share := totalAmount / float64(len(participants))
	res := make([]splitPair, 0, len(participants))
	for _, p := range participants {
		res = append(res, splitPair{p, share})
	}
	return res
}

func (equalSplit) validate(totalAmount float64, participants []string, params []float64) bool {
	return len(participants) > 0
}

type exactSplit struct{}

func (exactSplit) split(totalAmount float64, participants []string, params []float64) []splitPair {
	res := make([]splitPair, 0, len(participants))
	for i := range participants {
		res = append(res, splitPair{participants[i], params[i]})
	}
	return res
}

func (exactSplit) validate(totalAmount float64, participants []string, params []float64) bool {
	if len(params) != len(participants) {
		return false
	}
	sum := 0.0
	for _, v := range params {
		sum += v
	}
	return math.Abs(sum-totalAmount) < 1e-9
}

type percentSplit struct{}

func (percentSplit) split(totalAmount float64, participants []string, params []float64) []splitPair {
	res := make([]splitPair, 0, len(participants))
	for i := range participants {
		res = append(res, splitPair{participants[i], totalAmount * params[i] / 100.0})
	}
	return res
}

func (percentSplit) validate(totalAmount float64, participants []string, params []float64) bool {
	if len(params) != len(participants) {
		return false
	}
	sum := 0.0
	for _, v := range params {
		sum += v
	}
	return math.Abs(sum-100.0) < 1e-9
}

// orderedFloat preserves key insertion order (like a Python dict).
type orderedFloat struct {
	keys []string
	m    map[string]float64
}

func newOrderedFloat() *orderedFloat {
	return &orderedFloat{keys: []string{}, m: map[string]float64{}}
}

func (o *orderedFloat) get(k string) (float64, bool) {
	v, ok := o.m[k]
	return v, ok
}

func (o *orderedFloat) set(k string, v float64) {
	if _, ok := o.m[k]; !ok {
		o.keys = append(o.keys, k)
	}
	o.m[k] = v
}

func (o *orderedFloat) del(k string) {
	if _, ok := o.m[k]; !ok {
		return
	}
	delete(o.m, k)
	for i, kk := range o.keys {
		if kk == k {
			o.keys = append(o.keys[:i], o.keys[i+1:]...)
			break
		}
	}
}

type expenseManager struct {
	balanceKeys []string                 // debtor insertion order
	balances    map[string]*orderedFloat // debtor -> (creditor -> amount)
}

func newExpenseManager() *expenseManager {
	return &expenseManager{balanceKeys: []string{}, balances: map[string]*orderedFloat{}}
}

func (e *expenseManager) addUser(userId, name string) {}

func (e *expenseManager) updateBalance(debtor, creditor string, amount float64) {
	if debtor == creditor {
		return
	}
	if inner, ok := e.balances[creditor]; ok {
		if cur, ok2 := inner.get(debtor); ok2 && cur > 0 {
			offset := math.Min(cur, amount)
			newVal := cur - offset
			inner.set(debtor, newVal)
			amount -= offset
			if newVal < 1e-9 {
				inner.del(debtor)
			}
		}
	}
	if amount > 1e-9 {
		inner, ok := e.balances[debtor]
		if !ok {
			inner = newOrderedFloat()
			e.balances[debtor] = inner
			e.balanceKeys = append(e.balanceKeys, debtor)
		}
		prev, _ := inner.get(creditor)
		inner.set(creditor, prev+amount)
	}
}

func (e *expenseManager) addExpense(paidBy string, amount float64, participants []string) {
	splits := equalSplit{}.split(amount, participants, []float64{})
	for _, sp := range splits {
		e.updateBalance(sp.user, paidBy, sp.amount)
	}
}

func (e *expenseManager) addExpenseWithStrategy(paidBy string, amount float64, participants []string, strategy splitStrategy, params []float64) {
	if !strategy.validate(amount, participants, params) {
		return
	}
	splits := strategy.split(amount, participants, params)
	for _, sp := range splits {
		e.updateBalance(sp.user, paidBy, sp.amount)
	}
}

type simplifyTxn struct {
	from   string
	to     string
	amount float64
}

func (e *expenseManager) simplifyDebts() []simplifyTxn {
	net := newOrderedFloat()
	for _, debtor := range e.balanceKeys {
		inner := e.balances[debtor]
		for _, creditor := range inner.keys {
			amount := inner.m[creditor]
			dv, _ := net.get(debtor)
			net.set(debtor, dv-amount)
			cv, _ := net.get(creditor)
			net.set(creditor, cv+amount)
		}
	}

	type entry struct {
		user   string
		amount float64
	}
	creditors := []entry{}
	debtors := []entry{}
	for _, user := range net.keys {
		amount := net.m[user]
		if amount > 1e-9 {
			creditors = append(creditors, entry{user, amount})
		} else if amount < -1e-9 {
			debtors = append(debtors, entry{user, -amount})
		}
	}

	sort.SliceStable(creditors, func(i, j int) bool { return creditors[i].amount > creditors[j].amount })
	sort.SliceStable(debtors, func(i, j int) bool { return debtors[i].amount > debtors[j].amount })

	txns := []simplifyTxn{}
	i, j := 0, 0
	for i < len(creditors) && j < len(debtors) {
		settle := math.Min(creditors[i].amount, debtors[j].amount)
		txns = append(txns, simplifyTxn{from: debtors[j].user, to: creditors[i].user, amount: settle})
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

func splitCsv(s string) []string {
	if s == "" {
		return []string{}
	}
	return strings.Split(s, ",")
}

func splitCsvDouble(s string) []float64 {
	if s == "" {
		return []float64{}
	}
	parts := strings.Split(s, ",")
	res := make([]float64, 0, len(parts))
	for _, p := range parts {
		v, _ := strconv.ParseFloat(p, 64)
		res = append(res, v)
	}
	return res
}

func num2(v float64) string {
	return fmt.Sprintf("%.2f", v)
}

func splitwise_simulate(ops []SplitOp) []string {
	out := []string{}
	mgr := newExpenseManager()
	for _, op := range ops {
		k := op.kind
		switch k {
		case "new":
			mgr = newExpenseManager()
			out = append(out, "ok")
		case "add_user":
			mgr.addUser(op.s1, op.s2)
			out = append(out, "ok")
		case "add_expense":
			mgr.addExpense(op.s2, float64(op.i1), splitCsv(op.s3))
			out = append(out, "ok")
		case "add_eq_strat":
			mgr.addExpenseWithStrategy(op.s2, float64(op.i1), splitCsv(op.s3), equalSplit{}, []float64{})
			out = append(out, "ok")
		case "add_exact":
			mgr.addExpenseWithStrategy(op.s2, float64(op.i1), splitCsv(op.s3), exactSplit{}, splitCsvDouble(op.s4))
			out = append(out, "ok")
		case "add_pct":
			mgr.addExpenseWithStrategy(op.s2, float64(op.i1), splitCsv(op.s3), percentSplit{}, splitCsvDouble(op.s4))
			out = append(out, "ok")
		case "balance":
			v := 0.0
			if inner, ok := mgr.balances[op.s1]; ok {
				if val, ok2 := inner.get(op.s2); ok2 {
					v = val
				}
			}
			out = append(out, num2(v))
		case "validate_exact":
			ok := exactSplit{}.validate(float64(op.i1), splitCsv(op.s3), splitCsvDouble(op.s4))
			out = append(out, yesNo(ok))
		case "validate_pct":
			ok := percentSplit{}.validate(float64(op.i1), splitCsv(op.s3), splitCsvDouble(op.s4))
			out = append(out, yesNo(ok))
		case "simplify_count":
			out = append(out, fmt.Sprintf("%d", len(mgr.simplifyDebts())))
		case "simplify_total_to":
			txns := mgr.simplifyDebts()
			tot := 0.0
			for _, t := range txns {
				if t.to == op.s1 {
					tot += t.amount
				}
			}
			out = append(out, num2(tot))
		case "simplify_total_from":
			txns := mgr.simplifyDebts()
			tot := 0.0
			for _, t := range txns {
				if t.from == op.s1 {
					tot += t.amount
				}
			}
			out = append(out, num2(tot))
		case "simplify_unique_pair":
			txns := mgr.simplifyDebts()
			tot := 0.0
			for _, t := range txns {
				if t.from == op.s1 && t.to == op.s2 {
					tot += t.amount
				}
			}
			out = append(out, num2(tot))
		default:
			out = append(out, "unknown:"+k)
		}
	}
	return out
}

func yesNo(b bool) string {
	if b {
		return "yes"
	}
	return "no"
}

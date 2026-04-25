package main

import "fmt"

func part3Tests() int {
	passed := 0
	failed := 0

	test := func(name string, fn func()) {
		func() {
			defer func() {
				if r := recover(); r != nil {
					fmt.Println("FAIL", name)
					failed++
				}
			}()
			fn()
			fmt.Println("PASS", name)
			passed++
		}()
	}

	// Test 1: Simple two-person simplification
	test("test_simplify_two_person", func() {
		ResetManager()
		AddUser("alice", "Alice")
		AddUser("bob", "Bob")
		AddExpense("E1", "alice", 100.0, []string{"alice", "bob"})
		txns := SimplifyDebts()
		if len(txns) != 1 {
			panic("expected exactly 1 transaction")
		}
		if txns[0].Debtor != "bob" || txns[0].Creditor != "alice" {
			panic("wrong debtor/creditor in simplification")
		}
		if txns[0].Amount < 49.99 || txns[0].Amount > 50.01 {
			panic("wrong amount in simplification")
		}
	})

	// Test 2: Three-person: Alice paid for all, simplify to 2 transactions
	test("test_simplify_three_person_two_transactions", func() {
		ResetManager()
		AddUser("alice", "Alice")
		AddUser("bob", "Bob")
		AddUser("charlie", "Charlie")
		AddExpense("E1", "alice", 300.0, []string{"alice", "bob", "charlie"})
		txns := SimplifyDebts()
		if len(txns) != 2 {
			panic("expected 2 transactions after simplification")
		}
		totalSettled := 0.0
		for _, t := range txns {
			totalSettled += t.Amount
		}
		if totalSettled < 199.99 || totalSettled > 200.01 {
			panic("total settled amount should be 200")
		}
	})

	// Test 3: Balanced group — zero transactions
	test("test_simplify_no_debts", func() {
		ResetManager()
		AddUser("alice", "Alice")
		AddUser("bob", "Bob")
		// alice pays 100 for {alice, bob}, bob pays 100 for {alice, bob} — net zero
		AddExpense("E1", "alice", 100.0, []string{"alice", "bob"})
		AddExpense("E2", "bob", 100.0, []string{"alice", "bob"})
		txns := SimplifyDebts()
		if len(txns) != 0 {
			panic("expected 0 transactions when balances cancel out")
		}
	})

	// Test 4: Simplification reduces transaction count vs raw balances
	test("test_simplify_reduces_transactions", func() {
		ResetManager()
		AddUser("A", "A")
		AddUser("B", "B")
		AddUser("C", "C")
		// A pays for B: B owes A 50
		AddExpense("E1", "A", 100.0, []string{"A", "B"})
		// B pays for C: C owes B 50
		AddExpense("E2", "B", 100.0, []string{"B", "C"})
		// Without simplification: B→A 50, C→B 50 (2 txns)
		// With simplification: C→A 50 (1 txn) since B is balanced
		txns := SimplifyDebts()
		if len(txns) > 2 {
			panic("simplified transactions should not exceed raw transaction count")
		}
	})

	fmt.Printf("PART3_SUMMARY %d/%d\n", passed, passed+failed)
	return failed
}

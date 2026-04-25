package main

import "fmt"

func part1Tests() int {
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

	// Test 1: Equal split — payer excluded from balance
	test("test_equal_split_payer_excluded", func() {
		ResetManager()
		AddUser("alice", "Alice")
		AddUser("bob", "Bob")
		AddUser("charlie", "Charlie")
		AddExpense("E1", "alice", 300.0, []string{"alice", "bob", "charlie"})
		balances := GetBalances()
		if balances["bob"]["alice"] < 99.99 || balances["bob"]["alice"] > 100.01 {
			panic("bob should owe alice 100")
		}
		if balances["charlie"]["alice"] < 99.99 || balances["charlie"]["alice"] > 100.01 {
			panic("charlie should owe alice 100")
		}
	})

	// Test 2: Alice does not owe herself
	test("test_payer_has_no_self_debt", func() {
		ResetManager()
		AddUser("alice", "Alice")
		AddUser("bob", "Bob")
		AddExpense("E1", "alice", 200.0, []string{"alice", "bob"})
		balances := GetBalances()
		if balances["alice"]["alice"] != 0 {
			panic("alice should not owe herself")
		}
	})

	// Test 3: Bob pays; Alice owes Bob
	test("test_other_user_pays", func() {
		ResetManager()
		AddUser("alice", "Alice")
		AddUser("bob", "Bob")
		AddExpense("E1", "bob", 100.0, []string{"alice", "bob"})
		balances := GetBalances()
		if balances["alice"]["bob"] < 49.99 || balances["alice"]["bob"] > 50.01 {
			panic("alice should owe bob 50")
		}
	})

	// Test 4: Net-out opposing debts
	test("test_net_out_opposing_debts", func() {
		ResetManager()
		AddUser("alice", "Alice")
		AddUser("bob", "Bob")
		// Alice pays 100 split equally → bob owes alice 50
		AddExpense("E1", "alice", 100.0, []string{"alice", "bob"})
		// Bob pays 60 split equally → alice owes bob 30
		AddExpense("E2", "bob", 60.0, []string{"alice", "bob"})
		balances := GetBalances()
		// Net: bob owed alice 50, alice owed bob 30 → bob still owes alice 20
		if balances["bob"]["alice"] < 19.99 || balances["bob"]["alice"] > 20.01 {
			panic("net balance should be 20 from bob to alice")
		}
		if balances["alice"]["bob"] != 0 {
			panic("alice should not still owe bob after netting out")
		}
	})

	// Test 5: Multiple expenses accumulate
	test("test_multiple_expenses_accumulate", func() {
		ResetManager()
		AddUser("alice", "Alice")
		AddUser("bob", "Bob")
		AddExpense("E1", "alice", 90.0, []string{"alice", "bob"})
		AddExpense("E2", "alice", 90.0, []string{"alice", "bob"})
		balances := GetBalances()
		if balances["bob"]["alice"] < 89.99 || balances["bob"]["alice"] > 90.01 {
			panic("bob should owe alice 90 total (45+45)")
		}
	})

	// Test 6: Empty participants list is safe (no crash)
	test("test_empty_participants_no_crash", func() {
		ResetManager()
		AddUser("alice", "Alice")
		AddExpense("E1", "alice", 100.0, []string{})
		// Just verify GetBalances returns without panic
		_ = GetBalances()
	})

	fmt.Printf("PART1_SUMMARY %d/%d\n", passed, passed+failed)
	return failed
}

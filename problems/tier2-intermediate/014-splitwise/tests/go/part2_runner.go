package main

import "fmt"

func part2Tests() int {
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

	// Test 1: EqualSplit via AddExpenseWithStrategy
	test("test_equal_split_strategy", func() {
		ResetManager()
		AddUser("alice", "Alice")
		AddUser("bob", "Bob")
		AddUser("charlie", "Charlie")
		AddExpenseWithStrategy("E1", "alice", 300.0,
			[]string{"alice", "bob", "charlie"}, EqualSplit{}, []float64{})
		balances := GetBalances()
		if balances["bob"]["alice"] < 99.99 || balances["bob"]["alice"] > 100.01 {
			panic("bob should owe alice 100")
		}
	})

	// Test 2: ExactSplit with valid params
	test("test_exact_split_valid", func() {
		ResetManager()
		AddUser("alice", "Alice")
		AddUser("bob", "Bob")
		AddExpenseWithStrategy("E1", "alice", 150.0,
			[]string{"alice", "bob"}, ExactSplit{}, []float64{100.0, 50.0})
		balances := GetBalances()
		// alice's own share is 100 (she paid); bob owes 50
		if balances["bob"]["alice"] < 49.99 || balances["bob"]["alice"] > 50.01 {
			panic("bob should owe alice 50")
		}
	})

	// Test 3: ExactSplit with invalid params (sum != total) is ignored
	test("test_exact_split_invalid_ignored", func() {
		ResetManager()
		AddUser("alice", "Alice")
		AddUser("bob", "Bob")
		AddExpenseWithStrategy("E1", "alice", 150.0,
			[]string{"alice", "bob"}, ExactSplit{}, []float64{80.0, 50.0}) // sum=130 != 150
		balances := GetBalances()
		if balances["bob"]["alice"] != 0 {
			panic("invalid expense should not update balances")
		}
	})

	// Test 4: PercentSplit with valid percentages
	test("test_percent_split_valid", func() {
		ResetManager()
		AddUser("alice", "Alice")
		AddUser("bob", "Bob")
		// alice 60%, bob 40% of 200
		AddExpenseWithStrategy("E1", "alice", 200.0,
			[]string{"alice", "bob"}, PercentSplit{}, []float64{60.0, 40.0})
		balances := GetBalances()
		// alice paid 200, alice's own share=120, bob's share=80
		if balances["bob"]["alice"] < 79.99 || balances["bob"]["alice"] > 80.01 {
			panic("bob should owe alice 80")
		}
	})

	// Test 5: PercentSplit with percentages not summing to 100 is ignored
	test("test_percent_split_invalid_ignored", func() {
		ResetManager()
		AddUser("alice", "Alice")
		AddUser("bob", "Bob")
		AddExpenseWithStrategy("E1", "alice", 200.0,
			[]string{"alice", "bob"}, PercentSplit{}, []float64{60.0, 30.0}) // sum=90
		balances := GetBalances()
		if balances["bob"]["alice"] != 0 {
			panic("invalid percent split should not update balances")
		}
	})

	// Test 6: ExactSplit wrong number of params is ignored
	test("test_exact_split_param_count_mismatch_ignored", func() {
		ResetManager()
		AddUser("alice", "Alice")
		AddUser("bob", "Bob")
		AddExpenseWithStrategy("E1", "alice", 100.0,
			[]string{"alice", "bob"}, ExactSplit{}, []float64{100.0}) // only 1 param for 2 participants
		balances := GetBalances()
		if balances["bob"]["alice"] != 0 {
			panic("mismatched param count should not update balances")
		}
	})

	fmt.Printf("PART2_SUMMARY %d/%d\n", passed, passed+failed)
	return failed
}

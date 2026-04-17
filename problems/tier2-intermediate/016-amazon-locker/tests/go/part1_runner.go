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

	// Test 1: Deposit into exact-size locker
	test("test_deposit_small_into_small_locker", func() {
		InitLockerSystem()
		AddLocker("S1", SizeSmall)
		code := DepositPackage("PKG001", SizeSmall)
		if code == "" {
			panic("deposit should succeed and return a code")
		}
	})

	// Test 2: Retrieve with valid code succeeds
	test("test_retrieve_valid_code", func() {
		InitLockerSystem()
		AddLocker("S1", SizeSmall)
		code := DepositPackage("PKG001", SizeSmall)
		ok := RetrievePackage(code)
		if !ok {
			panic("retrieval with valid code should succeed")
		}
	})

	// Test 3: Retrieve with invalid code fails
	test("test_retrieve_invalid_code", func() {
		InitLockerSystem()
		ok := RetrievePackage("FAKE-CODE")
		if ok {
			panic("retrieval with invalid code should fail")
		}
	})

	// Test 4: Cannot use same code twice
	test("test_retrieve_same_code_twice", func() {
		InitLockerSystem()
		AddLocker("S1", SizeSmall)
		code := DepositPackage("PKG001", SizeSmall)
		RetrievePackage(code)
		ok := RetrievePackage(code)
		if ok {
			panic("second retrieval with same code should fail")
		}
	})

	// Test 5: Locker freed after retrieval — can accept new package
	test("test_locker_freed_after_retrieval", func() {
		InitLockerSystem()
		AddLocker("S1", SizeSmall)
		code1 := DepositPackage("PKG001", SizeSmall)
		RetrievePackage(code1)
		code2 := DepositPackage("PKG002", SizeSmall)
		if code2 == "" {
			panic("locker should be available after retrieval")
		}
	})

	// Test 6: Small package can use medium locker when small is unavailable
	test("test_small_package_upgrades_to_medium", func() {
		InitLockerSystem()
		AddLocker("M1", SizeMedium)
		code := DepositPackage("PKG001", SizeSmall)
		if code == "" {
			panic("small package should fit in medium locker")
		}
	})

	// Test 7: Medium package cannot use small locker
	test("test_medium_package_cannot_use_small_locker", func() {
		InitLockerSystem()
		AddLocker("S1", SizeSmall)
		code := DepositPackage("PKG001", SizeMedium)
		if code != "" {
			panic("medium package should not fit in small locker")
		}
	})

	// Test 8: No available locker returns empty code
	test("test_no_available_locker", func() {
		InitLockerSystem()
		AddLocker("S1", SizeSmall)
		DepositPackage("PKG001", SizeSmall)
		code := DepositPackage("PKG002", SizeSmall)
		if code != "" {
			panic("deposit should fail when no locker is available")
		}
	})

	// Test 9: Smallest-fit prefers exact size over larger
	test("test_smallest_fit_prefers_exact_size", func() {
		InitLockerSystem()
		AddLocker("S1", SizeSmall)
		AddLocker("M1", SizeMedium)
		code := DepositPackage("PKG001", SizeSmall)
		if code == "" {
			panic("deposit should succeed")
		}
		// After using S1, depositing another small should use medium (upgrade)
		code2 := DepositPackage("PKG002", SizeSmall)
		if code2 == "" {
			panic("second small deposit should use medium locker")
		}
	})

	// Test 10: Large package requires large locker
	test("test_large_package_requires_large_locker", func() {
		InitLockerSystem()
		AddLocker("S1", SizeSmall)
		AddLocker("M1", SizeMedium)
		code := DepositPackage("PKG001", SizeLarge)
		if code != "" {
			panic("large package should not fit in small or medium locker")
		}
		AddLocker("L1", SizeLarge)
		code2 := DepositPackage("PKG001", SizeLarge)
		if code2 == "" {
			panic("large package should fit in large locker")
		}
	})

	fmt.Printf("PART1_SUMMARY %d/%d\n", passed, passed+failed)
	return failed
}

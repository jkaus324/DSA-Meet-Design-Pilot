package main

import "fmt"

// MockNotificationChannel captures all notifications for test assertions.
type MockNotificationChannel struct {
	Notifications []string
}

func (m *MockNotificationChannel) Notify(packageID, message string) {
	m.Notifications = append(m.Notifications, packageID+":"+message)
}

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

	// Test 1: Notification sent on deposit
	test("test_notification_on_deposit", func() {
		InitLockerSystem()
		AddLocker("S1", SizeSmall)
		mock := &MockNotificationChannel{}
		AddNotificationChannel(mock)
		DepositPackage("PKG001", SizeSmall, 0)
		if len(mock.Notifications) == 0 {
			panic("expected notification on deposit")
		}
	})

	// Test 2: No expiry by default — CheckExpired returns empty
	test("test_no_expiry_by_default", func() {
		InitLockerSystem()
		AddLocker("S1", SizeSmall)
		DepositPackage("PKG001", SizeSmall, 0)
		expired := CheckExpired(9999999)
		if len(expired) != 0 {
			panic("no packages should expire when expiryHours is 0")
		}
	})

	// Test 3: Package expires after configured hours
	test("test_package_expires_after_configured_hours", func() {
		InitLockerSystem()
		AddLocker("S1", SizeSmall)
		SetCodeExpiry(24) // 24 hours = 86400 seconds
		DepositPackage("PKG001", SizeSmall, 0)
		// Check at 90000s (25 hours) — should be expired
		expired := CheckExpired(90000)
		if len(expired) != 1 || expired[0] != "PKG001" {
			panic("PKG001 should be expired at 90000s with 24h expiry")
		}
	})

	// Test 4: Package not expired before threshold
	test("test_package_not_expired_before_threshold", func() {
		InitLockerSystem()
		AddLocker("S1", SizeSmall)
		SetCodeExpiry(24)
		DepositPackage("PKG001", SizeSmall, 0)
		// Check at 3600s (1 hour) — not expired yet
		expired := CheckExpired(3600)
		if len(expired) != 0 {
			panic("package should not be expired at 1h with 24h expiry")
		}
	})

	// Test 5: Expired locker is freed for new deposit
	test("test_expired_locker_freed", func() {
		InitLockerSystem()
		AddLocker("S1", SizeSmall)
		SetCodeExpiry(1) // 1 hour = 3600 seconds
		DepositPackage("PKG001", SizeSmall, 0)
		CheckExpired(7200) // 2 hours later — PKG001 should expire
		code := DepositPackage("PKG002", SizeSmall, 7200)
		if code == "" {
			panic("freed locker should accept new package after expiry")
		}
	})

	// Test 6: Notification sent on expiry
	test("test_notification_on_expiry", func() {
		InitLockerSystem()
		AddLocker("S1", SizeSmall)
		mock := &MockNotificationChannel{}
		AddNotificationChannel(mock)
		SetCodeExpiry(1)
		DepositPackage("PKG001", SizeSmall, 0)
		initialCount := len(mock.Notifications)
		CheckExpired(7200)
		if len(mock.Notifications) <= initialCount {
			panic("expected expiry notification")
		}
	})

	// Test 7: Multiple packages expire in one CheckExpired call
	test("test_multiple_packages_expire", func() {
		InitLockerSystem()
		AddLocker("S1", SizeSmall)
		AddLocker("S2", SizeSmall)
		SetCodeExpiry(1)
		DepositPackage("PKG001", SizeSmall, 0)
		DepositPackage("PKG002", SizeSmall, 0)
		expired := CheckExpired(7200)
		if len(expired) != 2 {
			panic("both packages should expire")
		}
	})

	// Test 8: Multiple notification channels all receive notifications
	test("test_multiple_notification_channels", func() {
		InitLockerSystem()
		AddLocker("S1", SizeSmall)
		mock1 := &MockNotificationChannel{}
		mock2 := &MockNotificationChannel{}
		AddNotificationChannel(mock1)
		AddNotificationChannel(mock2)
		DepositPackage("PKG001", SizeSmall, 0)
		if len(mock1.Notifications) == 0 || len(mock2.Notifications) == 0 {
			panic("both channels should receive the deposit notification")
		}
	})

	fmt.Printf("PART2_SUMMARY %d/%d\n", passed, passed+failed)
	return failed
}

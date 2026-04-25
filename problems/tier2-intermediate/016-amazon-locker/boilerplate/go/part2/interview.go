package main

// --- Data Model (given -- do not modify) ------------------------------------

type LockerSize int

const (
	SizeSmall  LockerSize = iota
	SizeMedium LockerSize = iota
	SizeLarge  LockerSize = iota
)

type Package struct {
	PackageID string
	Size      LockerSize
}

type Locker struct {
	LockerID string
	Size     LockerSize
	Occupied bool
}

type DepositRecord struct {
	LockerID    string
	PackageID   string
	PickupCode  string
	DepositTime int64
}

// --- Your Design Starts Here (Part 2) ---------------------------------------
//
// Extend Part 1 to support:
//   1. Pluggable notification channels (Observer pattern):
//        - NotificationChannel interface with Notify(packageID, message string)
//        - AddNotificationChannel registers a channel; DepositPackage and
//          CheckExpired notify all registered channels.
//   2. Code expiry:
//        - SetCodeExpiry(hours int) configures a global expiry window.
//        - CheckExpired(currentTime int64) scans active deposits, frees any
//          locker whose deposit is older than expiryHours, and returns the
//          expired package IDs. Also notifies all channels per expired package.
//
// Think about:
//   - How do you store deposit timestamps for expiry checks?
//   - What interface do notification channels need?
//
// Entry points (must exist for tests — include Part 1 entry points too):
//   func InitLockerSystem()
//   func AddLocker(lockerID string, size LockerSize)
//   func DepositPackage(packageID string, size LockerSize, depositTime int64) string
//   func RetrievePackage(code string) bool
//   func SetCodeExpiry(hours int)
//   func CheckExpired(currentTime int64) []string
//   func AddNotificationChannel(ch NotificationChannel)

// -------------------------------------------------------------------------

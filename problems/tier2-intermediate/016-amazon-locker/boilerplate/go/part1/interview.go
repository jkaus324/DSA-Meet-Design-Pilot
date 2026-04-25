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

// --- Your Design Starts Here ------------------------------------------------
//
// Design and implement a LockerSystem that:
//   1. Adds lockers of various sizes
//   2. Deposits a package into the smallest available compatible locker
//      (Small fits Small; Medium fits Medium or Small; Large fits any)
//      Returns a pickup code, or "" if no compatible locker is available.
//   3. Retrieves a package by pickup code, frees the locker
//
// Think about:
//   - How do you quickly find the smallest compatible locker?
//   - What data structure lets you allocate and free lockers efficiently?
//
// Entry points (must exist for tests):
//   func InitLockerSystem()
//   func AddLocker(lockerID string, size LockerSize)
//   func DepositPackage(packageID string, size LockerSize) string
//   func RetrievePackage(code string) bool

// -------------------------------------------------------------------------

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

// --- Deposit Record ---------------------------------------------------------
// HINT: When a package is deposited, record which locker holds it and the
//       pickup code, so RetrievePackage can free the right locker.

// type DepositRecord struct {
//     LockerID  string
//     PackageID string
//     PickupCode string
// }

// --- Allocation Strategy Interface ------------------------------------------
// HINT: Define a LockerAllocator interface with one method:
//         Allocate(packageSize LockerSize, available map[LockerSize][]string) string
//       It removes and returns a lockerID from available, or "" if none fits.

// type LockerAllocator interface {
//     Allocate(packageSize LockerSize, available map[LockerSize][]string) string
// }

// --- SmallestFitAllocator ---------------------------------------------------
// HINT: Try sizes in ascending order starting from packageSize.
//       SizeSmall → try Small, then Medium, then Large.
//       SizeMedium → try Medium, then Large.
//       SizeLarge → try Large only.
// HINT: Use a slice as a queue: pop the first element with available[sz][0]
//       and advance the slice with available[sz] = available[sz][1:].

// type SmallestFitAllocator struct{}
// func (s SmallestFitAllocator) Allocate(...) string

// --- Locker System ----------------------------------------------------------
// HINT: available map[LockerSize][]string acts as per-size queues of free lockerIDs.
// HINT: codeCounter int generates sequential codes like "CODE-1", "CODE-2".

// type LockerSystem struct {
//     lockers       map[string]Locker
//     available     map[LockerSize][]string
//     activeDeposits map[string]DepositRecord  // code → record
//     allocator     LockerAllocator
//     codeCounter   int
// }

// func NewLockerSystem() *LockerSystem
// func (ls *LockerSystem) AddLocker(lockerID string, size LockerSize)
// func (ls *LockerSystem) DepositPackage(packageID string, size LockerSize) string
// func (ls *LockerSystem) RetrievePackage(code string) bool

// --- Global Entry Points (required by tests) --------------------------------

// var gSystem *LockerSystem
// func InitLockerSystem()
// func AddLocker(lockerID string, size LockerSize)
// func DepositPackage(packageID string, size LockerSize) string
// func RetrievePackage(code string) bool

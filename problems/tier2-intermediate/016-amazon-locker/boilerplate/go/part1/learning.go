package main

import "fmt"

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

type DepositRecord struct {
	LockerID   string
	PackageID  string
	PickupCode string
}

// --- Allocation Strategy ----------------------------------------------------

type LockerAllocator interface {
	Allocate(packageSize LockerSize, available map[LockerSize][]string) string
}

type SmallestFitAllocator struct{}

func (s SmallestFitAllocator) Allocate(packageSize LockerSize, available map[LockerSize][]string) string {
	// TODO: Build tryOrder based on packageSize:
	//         Small  → [Small, Medium, Large]
	//         Medium → [Medium, Large]
	//         Large  → [Large]
	// TODO: For each size in tryOrder, if available[sz] is non-empty,
	//       pop the first element and return it (update available[sz])
	// TODO: Return "" if nothing found
	return ""
}

// --- Locker System ----------------------------------------------------------

type LockerSystem struct {
	lockers        map[string]Locker
	available      map[LockerSize][]string
	activeDeposits map[string]DepositRecord
	allocator      LockerAllocator
	codeCounter    int
}

func NewLockerSystem() *LockerSystem {
	// TODO: Initialise maps; set allocator = SmallestFitAllocator{}; return
	return nil
}

func (ls *LockerSystem) generateCode() string {
	ls.codeCounter++
	return fmt.Sprintf("CODE-%d", ls.codeCounter)
}

func (ls *LockerSystem) freeLocker(lockerID string) {
	// TODO: Set ls.lockers[lockerID].Occupied = false
	// TODO: Append lockerID back to ls.available[locker.Size]
}

func (ls *LockerSystem) AddLocker(lockerID string, size LockerSize) {
	// TODO: Store Locker{LockerID, Size, Occupied: false} in ls.lockers
	// TODO: Append lockerID to ls.available[size]
}

func (ls *LockerSystem) DepositPackage(packageID string, size LockerSize) string {
	// TODO: Call ls.allocator.Allocate to get a lockerID; return "" if empty
	// TODO: Mark locker as occupied
	// TODO: Generate a pickup code
	// TODO: Store DepositRecord in ls.activeDeposits[code]
	// TODO: Return the code
	return ""
}

func (ls *LockerSystem) RetrievePackage(code string) bool {
	// TODO: Look up record in ls.activeDeposits; return false if not found
	// TODO: Call ls.freeLocker(record.LockerID)
	// TODO: Delete from ls.activeDeposits; return true
	return false
}

// --- Global Entry Points (required by tests) --------------------------------

var gSystem *LockerSystem

func InitLockerSystem() {
	gSystem = NewLockerSystem()
}

func AddLocker(lockerID string, size LockerSize) {
	// TODO: gSystem.AddLocker(lockerID, size)
}

func DepositPackage(packageID string, size LockerSize) string {
	// TODO: return gSystem.DepositPackage(packageID, size)
	return ""
}

func RetrievePackage(code string) bool {
	// TODO: return gSystem.RetrievePackage(code)
	return false
}

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

type DepositRecord struct {
	LockerID    string
	PackageID   string
	PickupCode  string
	DepositTime int64
}

// --- Notification Channel ---------------------------------------------------

type NotificationChannel interface {
	Notify(packageID, message string)
}

// --- Allocation Strategy (same as Part 1) ------------------------------------

type LockerAllocator interface {
	Allocate(packageSize LockerSize, available map[LockerSize][]string) string
}

type SmallestFitAllocator struct{}

func (s SmallestFitAllocator) Allocate(packageSize LockerSize, available map[LockerSize][]string) string {
	// TODO: Same as Part 1 — try sizes ascending from packageSize
	return ""
}

// --- Locker System ----------------------------------------------------------

type LockerSystem struct {
	lockers        map[string]Locker
	available      map[LockerSize][]string
	activeDeposits map[string]DepositRecord
	allocator      LockerAllocator
	channels       []NotificationChannel
	codeCounter    int
	expiryHours    int
}

func NewLockerSystem() *LockerSystem {
	// TODO: Initialise all fields; expiryHours=0 (disabled)
	return nil
}

func (ls *LockerSystem) generateCode() string {
	ls.codeCounter++
	return fmt.Sprintf("CODE-%d", ls.codeCounter)
}

func (ls *LockerSystem) freeLocker(lockerID string) {
	// TODO: Same as Part 1
}

func (ls *LockerSystem) notifyAll(packageID, message string) {
	// TODO: Call ch.Notify(packageID, message) for each ch in ls.channels
}

func (ls *LockerSystem) AddLocker(lockerID string, size LockerSize) {
	// TODO: Same as Part 1
}

func (ls *LockerSystem) DepositPackage(packageID string, size LockerSize, depositTime int64) string {
	// TODO: Allocate locker; return "" if none available
	// TODO: Mark occupied; generate code; store DepositRecord with depositTime
	// TODO: notifyAll(packageID, "Package "+packageID+" deposited. Code: "+code)
	// TODO: Return code
	return ""
}

func (ls *LockerSystem) RetrievePackage(code string) bool {
	// TODO: Look up code; freeLocker; delete record; return true/false
	return false
}

func (ls *LockerSystem) SetCodeExpiry(hours int) {
	// TODO: ls.expiryHours = hours
}

func (ls *LockerSystem) CheckExpired(currentTime int64) []string {
	// TODO: If expiryHours <= 0, return nil (expiry disabled)
	// TODO: Iterate ls.activeDeposits (use a keys-snapshot to avoid modifying while ranging)
	// TODO: For each record where currentTime - record.DepositTime > int64(ls.expiryHours)*3600:
	//         freeLocker, append packageID, notifyAll, delete from activeDeposits
	// TODO: Return list of expired packageIDs
	return nil
}

func (ls *LockerSystem) AddNotificationChannel(ch NotificationChannel) {
	// TODO: Append ch to ls.channels
}

// --- Global Entry Points (required by tests) --------------------------------

var gSystem *LockerSystem

func InitLockerSystem() {
	gSystem = NewLockerSystem()
}

func AddLocker(lockerID string, size LockerSize) {
	// TODO: gSystem.AddLocker(lockerID, size)
}

func DepositPackage(packageID string, size LockerSize, depositTime int64) string {
	// TODO: return gSystem.DepositPackage(packageID, size, depositTime)
	return ""
}

func RetrievePackage(code string) bool {
	// TODO: return gSystem.RetrievePackage(code)
	return false
}

func SetCodeExpiry(hours int) {
	// TODO: gSystem.SetCodeExpiry(hours)
}

func CheckExpired(currentTime int64) []string {
	// TODO: return gSystem.CheckExpired(currentTime)
	return nil
}

func AddNotificationChannel(ch NotificationChannel) {
	// TODO: gSystem.AddNotificationChannel(ch)
}

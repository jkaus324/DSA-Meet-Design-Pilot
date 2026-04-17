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

// --- Notification Channel Interface -----------------------------------------
// HINT: Define a Go interface:
//         type NotificationChannel interface {
//             Notify(packageID, message string)
//         }
// HINT: In LockerSystem, store []NotificationChannel and a notifyAll helper.

// --- Locker System (extends Part 1) ------------------------------------------
// HINT: Add expiryHours int to LockerSystem (0 means no expiry).
// HINT: DepositPackage now takes depositTime int64 to record when the deposit happened.
// HINT: CheckExpired(currentTime int64) iterates activeDeposits:
//         if expiryHours > 0 && currentTime - record.DepositTime > int64(expiryHours)*3600:
//             freeLocker, collect packageID, notifyAll, delete record

// type NotificationChannel interface { Notify(packageID, message string) }

// type LockerSystem struct {
//     lockers        map[string]Locker
//     available      map[LockerSize][]string
//     activeDeposits map[string]DepositRecord
//     allocator      LockerAllocator
//     channels       []NotificationChannel
//     codeCounter    int
//     expiryHours    int
// }

// func (ls *LockerSystem) SetCodeExpiry(hours int)
// func (ls *LockerSystem) CheckExpired(currentTime int64) []string
// func (ls *LockerSystem) AddNotificationChannel(ch NotificationChannel)

// --- Global Entry Points (required by tests) --------------------------------

// func InitLockerSystem()
// func AddLocker(lockerID string, size LockerSize)
// func DepositPackage(packageID string, size LockerSize, depositTime int64) string
// func RetrievePackage(code string) bool
// func SetCodeExpiry(hours int)
// func CheckExpired(currentTime int64) []string
// func AddNotificationChannel(ch NotificationChannel)

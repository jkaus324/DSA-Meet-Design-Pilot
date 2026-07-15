// Amazon Locker — allocation/deposit/retrieval with expiry + notifications (Go).
package main

import (
	"fmt"
)

type LockerOp struct {
	kind string
	s1   string
	s2   string
	i1   int
	i2   int
}

const (
	lockerSmall  = "SMALL"
	lockerMedium = "MEDIUM"
	lockerLarge  = "LARGE"
)

type locker struct {
	lockerId string
	size     string
	occupied bool
}

type depositRecord struct {
	lockerId    string
	packageId   string
	pickupCode  string
	depositTime int
}

func allocate(packageSize string, available map[string][]string) string {
	var tryOrder []string
	switch packageSize {
	case lockerSmall:
		tryOrder = []string{lockerSmall, lockerMedium, lockerLarge}
	case lockerMedium:
		tryOrder = []string{lockerMedium, lockerLarge}
	default:
		tryOrder = []string{lockerLarge}
	}
	for _, sz := range tryOrder {
		q := available[sz]
		if len(q) > 0 {
			id := q[0]
			available[sz] = q[1:]
			return id
		}
	}
	return ""
}

type lockerSystem struct {
	lockers        map[string]*locker
	available      map[string][]string
	depositKeys    []string // insertion order of active deposit codes
	activeDeposits map[string]*depositRecord
	chanLog        *[]string
	hasChannel     bool
	codeCounter    int
	expiryHours    int
}

func newLockerSystem() *lockerSystem {
	return &lockerSystem{
		lockers:        map[string]*locker{},
		available:      map[string][]string{lockerSmall: {}, lockerMedium: {}, lockerLarge: {}},
		depositKeys:    []string{},
		activeDeposits: map[string]*depositRecord{},
	}
}

func (ls *lockerSystem) generateCode() string {
	ls.codeCounter++
	return fmt.Sprintf("CODE-%d", ls.codeCounter)
}

func (ls *lockerSystem) notifyAll(packageId, message string) {
	if ls.hasChannel && ls.chanLog != nil {
		*ls.chanLog = append(*ls.chanLog, packageId+": "+message)
	}
}

func (ls *lockerSystem) freeLocker(lockerId string) {
	if l, ok := ls.lockers[lockerId]; ok {
		l.occupied = false
		ls.available[l.size] = append(ls.available[l.size], lockerId)
	}
}

func (ls *lockerSystem) addLocker(lockerId, size string) {
	ls.lockers[lockerId] = &locker{lockerId: lockerId, size: size}
	ls.available[size] = append(ls.available[size], lockerId)
}

func (ls *lockerSystem) addDepositKey(code string) {
	ls.depositKeys = append(ls.depositKeys, code)
}

func (ls *lockerSystem) delDepositKey(code string) {
	for i, c := range ls.depositKeys {
		if c == code {
			ls.depositKeys = append(ls.depositKeys[:i], ls.depositKeys[i+1:]...)
			break
		}
	}
}

func (ls *lockerSystem) depositPackage(packageId, size string, depositTime int) string {
	lockerId := allocate(size, ls.available)
	if lockerId == "" {
		return ""
	}
	ls.lockers[lockerId].occupied = true
	code := ls.generateCode()
	ls.activeDeposits[code] = &depositRecord{lockerId: lockerId, packageId: packageId, pickupCode: code, depositTime: depositTime}
	ls.addDepositKey(code)
	ls.notifyAll(packageId, "Package "+packageId+" deposited. Code: "+code)
	return code
}

func (ls *lockerSystem) retrievePackage(code string) bool {
	rec, ok := ls.activeDeposits[code]
	if !ok {
		return false
	}
	ls.freeLocker(rec.lockerId)
	delete(ls.activeDeposits, code)
	ls.delDepositKey(code)
	return true
}

func (ls *lockerSystem) setCodeExpiry(hours int) {
	ls.expiryHours = hours
}

func (ls *lockerSystem) checkExpired(currentTime int) []string {
	expired := []string{}
	if ls.expiryHours <= 0 {
		return expired
	}
	for _, code := range append([]string{}, ls.depositKeys...) {
		rec, ok := ls.activeDeposits[code]
		if !ok {
			continue
		}
		if currentTime-rec.depositTime > ls.expiryHours*3600 {
			ls.freeLocker(rec.lockerId)
			expired = append(expired, rec.packageId)
			ls.notifyAll(rec.packageId, "Package "+rec.packageId+" expired. Locker freed.")
			delete(ls.activeDeposits, code)
			ls.delDepositKey(code)
		}
	}
	return expired
}

func (ls *lockerSystem) addNotificationChannel() {
	ls.hasChannel = true
}

func lsizeFrom(s string) string {
	if s == "S" {
		return lockerSmall
	}
	if s == "M" {
		return lockerMedium
	}
	return lockerLarge
}

func locker_simulate(ops []LockerOp) []string {
	out := []string{}
	sys := newLockerSystem()
	codes := make([]string, 32)
	chanLog := []string{}
	var lastExpired []string
	for _, op := range ops {
		k := op.kind
		switch k {
		case "new":
			sys = newLockerSystem()
			codes = make([]string, 32)
			chanLog = []string{}
			lastExpired = nil
			out = append(out, "ok")
		case "add_locker":
			sys.addLocker(op.s1, lsizeFrom(op.s2))
			out = append(out, "ok")
		case "deposit":
			code := sys.depositPackage(op.s1, lsizeFrom(op.s2), op.i1)
			if op.i2 >= 0 && op.i2 < len(codes) {
				codes[op.i2] = code
			}
			out = append(out, code)
		case "code_at":
			out = append(out, codes[op.i2])
		case "retrieve":
			out = append(out, okFail(sys.retrievePackage(codes[op.i2])))
		case "retrieve_id":
			out = append(out, okFail(sys.retrievePackage(op.s1)))
		case "set_expiry":
			sys.setCodeExpiry(op.i1)
			out = append(out, "ok")
		case "check_expired":
			lastExpired = sys.checkExpired(op.i1)
			out = append(out, fmt.Sprintf("%d", len(lastExpired)))
		case "expired_at":
			if op.i2 >= 0 && op.i2 < len(lastExpired) {
				out = append(out, lastExpired[op.i2])
			} else {
				out = append(out, "")
			}
		case "add_chan":
			sys.addNotificationChannel()
			sys.chanLog = &chanLog
			out = append(out, "ok")
		case "chan_log_size":
			out = append(out, fmt.Sprintf("%d", len(chanLog)))
		case "chan_log_contains":
			found := false
			for _, entry := range chanLog {
				if contains(entry, op.s1) {
					found = true
					break
				}
			}
			out = append(out, yesNo(found))
		default:
			out = append(out, "unknown:"+k)
		}
	}
	return out
}

func contains(s, sub string) bool {
	if sub == "" {
		return true
	}
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}

func okFail(b bool) string {
	if b {
		return "ok"
	}
	return "fail"
}

func yesNo(b bool) string {
	if b {
		return "yes"
	}
	return "no"
}

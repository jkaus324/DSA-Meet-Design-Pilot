// BookMyShow — theater/show/booking + seat locking + contiguous search (Go).
package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type ShowOp struct {
	kind string
	s1   string
	s2   string
	s3   string
	s4   string
	i1   int
	i2   int
	i3   int
}

const (
	seatAvailable = "AVAILABLE"
	seatLocked    = "LOCKED"
	seatBooked    = "BOOKED"
)

type seat struct {
	row        int
	col        int
	status     string
	lockedBy   string
	lockExpiry int
	bookedBy   string
}

type show struct {
	id        string
	theaterId string
	movie     string
	time      string
	rows      int
	cols      int
	seats     [][]*seat
}

func newShow(id, theaterId, movie, time string, rows, cols int) *show {
	seats := make([][]*seat, rows)
	for r := 0; r < rows; r++ {
		seats[r] = make([]*seat, cols)
		for c := 0; c < cols; c++ {
			seats[r][c] = &seat{row: r, col: c, status: seatAvailable}
		}
	}
	return &show{id: id, theaterId: theaterId, movie: movie, time: time, rows: rows, cols: cols, seats: seats}
}

type theater struct {
	id   string
	name string
	city string
}

type pos struct{ r, c int }

type lockRec struct {
	id            string
	showId        string
	userId        string
	seatPositions []pos
	expiry        int
	confirmed     bool
	released      bool
}

type bookingSystem struct {
	theaters    map[string]*theater
	shows       map[string]*show
	bookings    map[string]bool
	locks       map[string]*lockRec
	cityMovies  map[string]map[string]bool
	lockCounter int
}

func newBookingSystem() *bookingSystem {
	return &bookingSystem{
		theaters:   map[string]*theater{},
		shows:      map[string]*show{},
		bookings:   map[string]bool{},
		locks:      map[string]*lockRec{},
		cityMovies: map[string]map[string]bool{},
	}
}

func (b *bookingSystem) expireSeat(s *seat, currentTime int) {
	if s.status == seatLocked && currentTime >= s.lockExpiry {
		s.status = seatAvailable
		s.lockedBy = ""
		s.lockExpiry = 0
	}
}

func (b *bookingSystem) addTheater(theaterId, name, city string) {
	b.theaters[theaterId] = &theater{id: theaterId, name: name, city: city}
}

func (b *bookingSystem) addShow(showId, theaterId, movie, time string, rows, cols int) {
	b.shows[showId] = newShow(showId, theaterId, movie, time, rows, cols)
	if t, ok := b.theaters[theaterId]; ok {
		if _, exists := b.cityMovies[t.city]; !exists {
			b.cityMovies[t.city] = map[string]bool{}
		}
		b.cityMovies[t.city][movie] = true
	}
}

func (b *bookingSystem) searchMovies(city string) []string {
	set, ok := b.cityMovies[city]
	if !ok {
		return []string{}
	}
	res := make([]string, 0, len(set))
	for m := range set {
		res = append(res, m)
	}
	sort.Strings(res)
	return res
}

func (b *bookingSystem) getAvailableSeats(showId string, currentTime int) []pos {
	s, ok := b.shows[showId]
	if !ok {
		return []pos{}
	}
	res := []pos{}
	for r := 0; r < s.rows; r++ {
		for c := 0; c < s.cols; c++ {
			b.expireSeat(s.seats[r][c], currentTime)
			if s.seats[r][c].status == seatAvailable {
				res = append(res, pos{r, c})
			}
		}
	}
	return res
}

func (b *bookingSystem) bookSeats(bookingId, showId string, seatPositions []pos, userId string, currentTime int) bool {
	s, ok := b.shows[showId]
	if !ok {
		return false
	}
	for _, p := range seatPositions {
		if p.r < 0 || p.r >= s.rows || p.c < 0 || p.c >= s.cols {
			return false
		}
		b.expireSeat(s.seats[p.r][p.c], currentTime)
		if s.seats[p.r][p.c].status != seatAvailable {
			return false
		}
	}
	for _, p := range seatPositions {
		s.seats[p.r][p.c].status = seatBooked
		s.seats[p.r][p.c].bookedBy = userId
	}
	b.bookings[bookingId] = true
	return true
}

func (b *bookingSystem) lockSeats(showId string, seatPositions []pos, userId string, ttlMinutes, currentTime int) string {
	s, ok := b.shows[showId]
	if !ok {
		return ""
	}
	for _, p := range seatPositions {
		if p.r < 0 || p.r >= s.rows || p.c < 0 || p.c >= s.cols {
			return ""
		}
		b.expireSeat(s.seats[p.r][p.c], currentTime)
		if s.seats[p.r][p.c].status != seatAvailable {
			return ""
		}
	}
	b.lockCounter++
	lockId := fmt.Sprintf("LOCK_%d", b.lockCounter)
	expiry := currentTime + ttlMinutes*60
	for _, p := range seatPositions {
		s.seats[p.r][p.c].status = seatLocked
		s.seats[p.r][p.c].lockedBy = userId
		s.seats[p.r][p.c].lockExpiry = expiry
	}
	b.locks[lockId] = &lockRec{
		id: lockId, showId: showId, userId: userId,
		seatPositions: seatPositions, expiry: expiry,
	}
	return lockId
}

func (b *bookingSystem) confirmBooking(lockId string, currentTime int) bool {
	lock, ok := b.locks[lockId]
	if !ok {
		return false
	}
	if lock.confirmed || lock.released {
		return false
	}
	if currentTime >= lock.expiry {
		if s, ok := b.shows[lock.showId]; ok {
			for _, p := range lock.seatPositions {
				b.expireSeat(s.seats[p.r][p.c], currentTime)
			}
		}
		lock.released = true
		return false
	}
	s, ok := b.shows[lock.showId]
	if !ok {
		return false
	}
	for _, p := range lock.seatPositions {
		s.seats[p.r][p.c].status = seatBooked
		s.seats[p.r][p.c].bookedBy = lock.userId
		s.seats[p.r][p.c].lockedBy = ""
		s.seats[p.r][p.c].lockExpiry = 0
	}
	bookingId := "BK_" + lockId
	b.bookings[bookingId] = true
	lock.confirmed = true
	return true
}

func (b *bookingSystem) releaseLock(lockId string, currentTime int) bool {
	lock, ok := b.locks[lockId]
	if !ok {
		return false
	}
	if lock.confirmed || lock.released {
		return false
	}
	if s, ok := b.shows[lock.showId]; ok {
		for _, p := range lock.seatPositions {
			s.seats[p.r][p.c].status = seatAvailable
			s.seats[p.r][p.c].lockedBy = ""
			s.seats[p.r][p.c].lockExpiry = 0
		}
	}
	lock.released = true
	return true
}

func (b *bookingSystem) findContiguousSeats(showId string, n, currentTime int) []pos {
	s, ok := b.shows[showId]
	if !ok {
		return []pos{}
	}
	for r := 0; r < s.rows; r++ {
		count := 0
		start := -1
		for c := 0; c < s.cols; c++ {
			b.expireSeat(s.seats[r][c], currentTime)
			if s.seats[r][c].status == seatAvailable {
				if count == 0 {
					start = c
				}
				count++
				if count == n {
					res := make([]pos, 0, n)
					for i := start; i < start+n; i++ {
						res = append(res, pos{r, i})
					}
					return res
				}
			} else {
				count = 0
				start = -1
			}
		}
	}
	return []pos{}
}

func parseSeats(s string) []pos {
	if s == "" {
		return []pos{}
	}
	out := []pos{}
	for _, token := range strings.Split(s, ";") {
		if strings.Contains(token, ",") {
			parts := strings.Split(token, ",")
			a, _ := strconv.Atoi(parts[0])
			bb, _ := strconv.Atoi(parts[1])
			out = append(out, pos{a, bb})
		}
	}
	return out
}

func show_simulate(ops []ShowOp) []string {
	out := []string{}
	sys := newBookingSystem()
	lockSlots := make([]string, 32)
	var lastContig []pos
	for _, op := range ops {
		k := op.kind
		switch k {
		case "new":
			sys = newBookingSystem()
			lockSlots = make([]string, 32)
			lastContig = nil
			out = append(out, "ok")
		case "add_theater":
			sys.addTheater(op.s1, op.s2, op.s3)
			out = append(out, "ok")
		case "add_show":
			sys.addShow(op.s1, op.s2, op.s3, op.s4, op.i1, op.i2)
			out = append(out, "ok")
		case "movies_count":
			out = append(out, fmt.Sprintf("%d", len(sys.searchMovies(op.s1))))
		case "movies_contains":
			m := sys.searchMovies(op.s1)
			found := false
			for _, x := range m {
				if x == op.s2 {
					found = true
					break
				}
			}
			out = append(out, yesNo(found))
		case "available_count":
			out = append(out, fmt.Sprintf("%d", len(sys.getAvailableSeats(op.s1, op.i1))))
		case "available_has":
			v := sys.getAvailableSeats(op.s1, op.i1)
			found := false
			for _, p := range v {
				if p.r == op.i2 && p.c == op.i3 {
					found = true
					break
				}
			}
			out = append(out, yesNo(found))
		case "book":
			ok := sys.bookSeats(op.s1, op.s2, parseSeats(op.s3), op.s4, op.i1)
			out = append(out, okFail(ok))
		case "lock":
			lid := sys.lockSeats(op.s1, parseSeats(op.s2), op.s3, op.i1, op.i2)
			if op.i3 >= 0 && op.i3 < len(lockSlots) {
				lockSlots[op.i3] = lid
			}
			out = append(out, lid)
		case "lock_at":
			out = append(out, lockSlots[op.i3])
		case "confirm":
			ok := sys.confirmBooking(lockSlots[op.i3], op.i1)
			out = append(out, okFail(ok))
		case "release":
			ok := sys.releaseLock(lockSlots[op.i3], op.i1)
			out = append(out, okFail(ok))
		case "release_id":
			ok := sys.releaseLock(op.s1, op.i1)
			out = append(out, okFail(ok))
		case "find_contig":
			lastContig = sys.findContiguousSeats(op.s1, op.i1, op.i2)
			out = append(out, fmt.Sprintf("%d", len(lastContig)))
		case "contig_at":
			if op.i1 < 0 || op.i1 >= len(lastContig) {
				out = append(out, "")
			} else {
				out = append(out, fmt.Sprintf("%d,%d", lastContig[op.i1].r, lastContig[op.i1].c))
			}
		default:
			out = append(out, "unknown:"+k)
		}
	}
	return out
}

func yesNo(b bool) string {
	if b {
		return "yes"
	}
	return "no"
}

func okFail(b bool) string {
	if b {
		return "ok"
	}
	return "fail"
}

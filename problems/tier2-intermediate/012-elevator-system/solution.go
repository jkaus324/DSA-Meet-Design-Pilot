// Elevator system — SCAN ordering + dispatch strategies (Go).
package main

import "strconv"

type ElevOp struct {
	kind string
	s1   string
	i1   int
	i2   int
}

type sortedSet struct {
	s   []int
	set map[int]bool
}

func newSortedSet() *sortedSet {
	return &sortedSet{s: []int{}, set: map[int]bool{}}
}

func (ss *sortedSet) add(x int) {
	if ss.set[x] {
		return
	}
	ss.set[x] = true
	lo, hi := 0, len(ss.s)
	for lo < hi {
		mid := (lo + hi) / 2
		if ss.s[mid] < x {
			lo = mid + 1
		} else {
			hi = mid
		}
	}
	ss.s = append(ss.s, 0)
	copy(ss.s[lo+1:], ss.s[lo:])
	ss.s[lo] = x
}

func (ss *sortedSet) remove(x int) {
	if !ss.set[x] {
		return
	}
	delete(ss.set, x)
	for i, v := range ss.s {
		if v == x {
			ss.s = append(ss.s[:i], ss.s[i+1:]...)
			break
		}
	}
}

func (ss *sortedSet) contains(x int) bool { return ss.set[x] }
func (ss *sortedSet) len() int            { return len(ss.s) }
func (ss *sortedSet) empty() bool         { return len(ss.s) == 0 }
func (ss *sortedSet) first() int          { return ss.s[0] }
func (ss *sortedSet) last() int           { return ss.s[len(ss.s)-1] }

const (
	IDLE        = "IDLE"
	MOVING_UP   = "MOVING_UP"
	MOVING_DOWN = "MOVING_DOWN"
	DOOR_OPEN   = "DOOR_OPEN"

	DIR_UP   = "UP"
	DIR_DOWN = "DOWN"
	DIR_NONE = "NONE"
)

type Elevator struct {
	id               int
	currentFloor     int
	state            string
	currentDirection string
	upRequests       *sortedSet
	downRequests     *sortedSet
}

func newElevator(id int) *Elevator {
	return &Elevator{
		id:               id,
		currentFloor:     0,
		state:            IDLE,
		currentDirection: DIR_NONE,
		upRequests:       newSortedSet(),
		downRequests:     newSortedSet(),
	}
}

func (e *Elevator) getCurrentFloor() int     { return e.currentFloor }
func (e *Elevator) getState() string         { return e.state }
func (e *Elevator) getCurrentDirection() string { return e.currentDirection }
func (e *Elevator) getPendingCount() int     { return e.upRequests.len() + e.downRequests.len() }

func iabs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func (e *Elevator) addRequest(floor int, direction string) {
	if floor == e.currentFloor && e.state == IDLE {
		e.state = DOOR_OPEN
		return
	}
	if direction == DIR_UP {
		e.upRequests.add(floor)
	} else if direction == DIR_DOWN {
		e.downRequests.add(floor)
	} else {
		if floor > e.currentFloor {
			e.upRequests.add(floor)
		} else {
			e.downRequests.add(floor)
		}
	}
	if e.state == IDLE {
		if !e.upRequests.empty() && (e.downRequests.empty() ||
			iabs(e.upRequests.first()-e.currentFloor) <= iabs(e.downRequests.last()-e.currentFloor)) {
			e.currentDirection = DIR_UP
			e.state = MOVING_UP
		} else {
			e.currentDirection = DIR_DOWN
			e.state = MOVING_DOWN
		}
	}
}

func (e *Elevator) step() {
	switch e.state {
	case IDLE:
		return
	case MOVING_UP:
		e.currentFloor++
		if e.upRequests.contains(e.currentFloor) {
			e.upRequests.remove(e.currentFloor)
			e.state = DOOR_OPEN
		}
		return
	case MOVING_DOWN:
		e.currentFloor--
		if e.downRequests.contains(e.currentFloor) {
			e.downRequests.remove(e.currentFloor)
			e.state = DOOR_OPEN
		}
		return
	case DOOR_OPEN:
		if e.currentDirection == DIR_UP {
			if !e.upRequests.empty() {
				e.state = MOVING_UP
			} else if !e.downRequests.empty() {
				e.currentDirection = DIR_DOWN
				e.state = MOVING_DOWN
			} else {
				e.currentDirection = DIR_NONE
				e.state = IDLE
			}
		} else if e.currentDirection == DIR_DOWN {
			if !e.downRequests.empty() {
				e.state = MOVING_DOWN
			} else if !e.upRequests.empty() {
				e.currentDirection = DIR_UP
				e.state = MOVING_UP
			} else {
				e.currentDirection = DIR_NONE
				e.state = IDLE
			}
		} else {
			if !e.upRequests.empty() {
				e.currentDirection = DIR_UP
				e.state = MOVING_UP
			} else if !e.downRequests.empty() {
				e.currentDirection = DIR_DOWN
				e.state = MOVING_DOWN
			} else {
				e.state = IDLE
			}
		}
	}
}

type DispatchStrategy interface {
	selectElevator(elevators []*Elevator, requestFloor int, requestDirection string) int
}

type NearestFirst struct{}

func (NearestFirst) selectElevator(elevators []*Elevator, requestFloor int, requestDirection string) int {
	bestIdx := 0
	bestScore := int(^uint(0) >> 1)
	const penalty = 10000
	for i, e := range elevators {
		dist := iabs(e.getCurrentFloor() - requestFloor)
		st := e.getState()
		d := e.getCurrentDirection()
		var score int
		if st == IDLE || d == DIR_NONE {
			score = dist
		} else if d == DIR_UP && requestDirection == DIR_UP && e.getCurrentFloor() <= requestFloor {
			score = dist
		} else if d == DIR_DOWN && requestDirection == DIR_DOWN && e.getCurrentFloor() >= requestFloor {
			score = dist
		} else {
			score = dist + penalty
		}
		if score < bestScore {
			bestScore = score
			bestIdx = i
		}
	}
	return bestIdx
}

type LeastLoaded struct{}

func (LeastLoaded) selectElevator(elevators []*Elevator, requestFloor int, requestDirection string) int {
	bestIdx := 0
	bestCount := int(^uint(0) >> 1)
	for i, e := range elevators {
		cnt := e.getPendingCount()
		if cnt < bestCount {
			bestCount = cnt
			bestIdx = i
		}
	}
	return bestIdx
}

type ElevatorSystem struct {
	elevators []*Elevator
	strategy  DispatchStrategy
}

func newElevatorSystem() *ElevatorSystem {
	return &ElevatorSystem{elevators: []*Elevator{}}
}

func (sys *ElevatorSystem) addElevator(id int) {
	sys.elevators = append(sys.elevators, newElevator(id))
}

func (sys *ElevatorSystem) setDispatchStrategy(s DispatchStrategy) { sys.strategy = s }

func (sys *ElevatorSystem) getElevator(index int) *Elevator {
	if index < 0 || index >= len(sys.elevators) {
		return nil
	}
	return sys.elevators[index]
}

func (sys *ElevatorSystem) getElevatorCount() int { return len(sys.elevators) }

func (sys *ElevatorSystem) addRequest(floor int, direction string) {
	if len(sys.elevators) == 0 {
		return
	}
	idx := 0
	if sys.strategy != nil {
		idx = sys.strategy.selectElevator(sys.elevators, floor, direction)
	}
	sys.elevators[idx].addRequest(floor, direction)
}

func (sys *ElevatorSystem) step() {
	for _, e := range sys.elevators {
		e.step()
	}
}

func dirFrom(s string) string {
	if s == "up" {
		return DIR_UP
	}
	if s == "down" {
		return DIR_DOWN
	}
	return DIR_NONE
}

func elevator_simulate(ops []ElevOp) []string {
	out := []string{}
	var single *Elevator
	var sys *ElevatorSystem
	nf := NearestFirst{}
	ll := LeastLoaded{}
	for _, op := range ops {
		switch op.kind {
		case "new_elev":
			single = newElevator(0)
			sys = nil
			out = append(out, "ok")
		case "new_sys":
			sys = newElevatorSystem()
			single = nil
			out = append(out, "ok")
		case "add_elev":
			sys.addElevator(op.i1)
			out = append(out, "ok")
		case "set_strategy":
			if op.s1 == "nearest" {
				sys.setDispatchStrategy(nf)
			} else if op.s1 == "least_loaded" {
				sys.setDispatchStrategy(ll)
			}
			out = append(out, "ok")
		case "req":
			single.addRequest(op.i1, dirFrom(op.s1))
			out = append(out, "ok")
		case "sys_req":
			sys.addRequest(op.i1, dirFrom(op.s1))
			out = append(out, "ok")
		case "elev_req":
			sys.getElevator(op.i1).addRequest(op.i2, dirFrom(op.s1))
			out = append(out, "ok")
		case "step":
			single.step()
			out = append(out, "ok")
		case "sys_step":
			sys.step()
			out = append(out, "ok")
		case "elev_step":
			sys.getElevator(op.i1).step()
			out = append(out, "ok")
		case "floor":
			out = append(out, strconv.Itoa(single.getCurrentFloor()))
		case "elev_floor":
			out = append(out, strconv.Itoa(sys.getElevator(op.i1).getCurrentFloor()))
		case "state":
			out = append(out, single.getState())
		case "elev_state":
			out = append(out, sys.getElevator(op.i1).getState())
		case "elev_pending":
			out = append(out, strconv.Itoa(sys.getElevator(op.i1).getPendingCount()))
		case "count":
			out = append(out, strconv.Itoa(sys.getElevatorCount()))
		case "elev_null":
			if sys.getElevator(op.i1) == nil {
				out = append(out, "yes")
			} else {
				out = append(out, "no")
			}
		default:
			out = append(out, "unknown:"+op.kind)
		}
	}
	return out
}

package main

import "math"

// ─── Data Model (given — do not modify) ──────────────────────────────────────

type ElevatorState int

const (
	IDLE        ElevatorState = iota
	MOVING_UP
	MOVING_DOWN
	DOOR_OPEN
)

type Direction int

const (
	UP   Direction = iota
	DOWN
	NONE
)

type Request struct {
	Floor     int
	Direction Direction
}

// ─── Elevator ─────────────────────────────────────────────────────────────────

type Elevator struct {
	currentFloor     int
	state            ElevatorState
	currentDirection Direction
	upRequests       map[int]bool // floors to visit going up
	downRequests     map[int]bool // floors to visit going down
}

func NewElevator() *Elevator {
	return &Elevator{
		currentFloor:     0,
		state:            IDLE,
		currentDirection: NONE,
		upRequests:       make(map[int]bool),
		downRequests:     make(map[int]bool),
	}
}

func (e *Elevator) GetCurrentFloor() int          { return e.currentFloor }
func (e *Elevator) GetState() ElevatorState       { return e.state }
func (e *Elevator) GetCurrentDirection() Direction { return e.currentDirection }
func (e *Elevator) GetPendingCount() int          { return len(e.upRequests) + len(e.downRequests) }

// nearestUp returns the minimum key in upRequests, or math.MaxInt32 if empty.
func (e *Elevator) nearestUp() int {
	min := math.MaxInt32
	for f := range e.upRequests {
		if f < min {
			min = f
		}
	}
	return min
}

// nearestDown returns the maximum key in downRequests, or math.MinInt32 if empty.
func (e *Elevator) nearestDown() int {
	max := math.MinInt32
	for f := range e.downRequests {
		if f > max {
			max = f
		}
	}
	return max
}

func (e *Elevator) AddRequest(floor int, direction Direction) {
	if floor == e.currentFloor && e.state == IDLE {
		// TODO: set state to DOOR_OPEN immediately
		return
	}
	if floor > e.currentFloor || (floor == e.currentFloor && direction == UP) {
		// TODO: add floor to e.upRequests
	} else {
		// TODO: add floor to e.downRequests
	}
	if e.state == IDLE {
		// TODO: decide direction: compare distance to nearest up vs nearest down
		//       start moving toward the nearer one
		upDist := math.MaxInt32
		downDist := math.MaxInt32
		if len(e.upRequests) > 0 {
			upDist = int(math.Abs(float64(e.nearestUp() - e.currentFloor)))
		}
		if len(e.downRequests) > 0 {
			downDist = int(math.Abs(float64(e.nearestDown() - e.currentFloor)))
		}
		if len(e.upRequests) > 0 && (len(e.downRequests) == 0 || upDist <= downDist) {
			// TODO: set currentDirection = UP, state = MOVING_UP
		} else {
			// TODO: set currentDirection = DOWN, state = MOVING_DOWN
		}
	}
}

func (e *Elevator) Step() {
	switch e.state {
	case IDLE:
		// nothing to do

	case MOVING_UP:
		// TODO: e.currentFloor++
		// TODO: if e.upRequests[e.currentFloor], delete it and set state = DOOR_OPEN

	case MOVING_DOWN:
		// TODO: e.currentFloor--
		// TODO: if e.downRequests[e.currentFloor], delete it and set state = DOOR_OPEN

	case DOOR_OPEN:
		// TODO: close doors and decide next state:
		//   If direction UP:
		//     if upRequests not empty → MOVING_UP
		//     else if downRequests not empty → set direction DOWN, MOVING_DOWN
		//     else → direction NONE, IDLE
		//   If direction DOWN:
		//     if downRequests not empty → MOVING_DOWN
		//     else if upRequests not empty → set direction UP, MOVING_UP
		//     else → direction NONE, IDLE
		//   If direction NONE: check both, or go IDLE
	}
}

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

// ─── Elevator (from Part 1 — assume fully implemented) ───────────────────────

type Elevator struct {
	id               int
	currentFloor     int
	state            ElevatorState
	currentDirection Direction
	upRequests       map[int]bool
	downRequests     map[int]bool
}

func NewElevator(id int) *Elevator {
	return &Elevator{
		id:               id,
		currentFloor:     0,
		state:            IDLE,
		currentDirection: NONE,
		upRequests:       make(map[int]bool),
		downRequests:     make(map[int]bool),
	}
}

func (e *Elevator) GetID() int                    { return e.id }
func (e *Elevator) GetCurrentFloor() int          { return e.currentFloor }
func (e *Elevator) GetState() ElevatorState       { return e.state }
func (e *Elevator) GetCurrentDirection() Direction { return e.currentDirection }
func (e *Elevator) GetPendingCount() int          { return len(e.upRequests) + len(e.downRequests) }

func (e *Elevator) nearestUp() int {
	min := math.MaxInt32
	for f := range e.upRequests {
		if f < min {
			min = f
		}
	}
	return min
}

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
		e.state = DOOR_OPEN
		return
	}
	if floor > e.currentFloor || (floor == e.currentFloor && direction == UP) {
		e.upRequests[floor] = true
	} else {
		e.downRequests[floor] = true
	}
	if e.state == IDLE {
		upDist := math.MaxInt32
		downDist := math.MaxInt32
		if len(e.upRequests) > 0 {
			upDist = int(math.Abs(float64(e.nearestUp() - e.currentFloor)))
		}
		if len(e.downRequests) > 0 {
			downDist = int(math.Abs(float64(e.nearestDown() - e.currentFloor)))
		}
		if len(e.upRequests) > 0 && (len(e.downRequests) == 0 || upDist <= downDist) {
			e.currentDirection = UP
			e.state = MOVING_UP
		} else {
			e.currentDirection = DOWN
			e.state = MOVING_DOWN
		}
	}
}

func (e *Elevator) Step() {
	switch e.state {
	case IDLE:
		// nothing

	case MOVING_UP:
		// TODO: e.currentFloor++
		// TODO: if e.upRequests[e.currentFloor] { delete and set DOOR_OPEN }

	case MOVING_DOWN:
		// TODO: e.currentFloor--
		// TODO: if e.downRequests[e.currentFloor] { delete and set DOOR_OPEN }

	case DOOR_OPEN:
		if e.currentDirection == UP {
			if len(e.upRequests) > 0 {
				e.state = MOVING_UP
			} else if len(e.downRequests) > 0 {
				e.currentDirection = DOWN
				e.state = MOVING_DOWN
			} else {
				e.currentDirection = NONE
				e.state = IDLE
			}
		} else if e.currentDirection == DOWN {
			if len(e.downRequests) > 0 {
				e.state = MOVING_DOWN
			} else if len(e.upRequests) > 0 {
				e.currentDirection = UP
				e.state = MOVING_UP
			} else {
				e.currentDirection = NONE
				e.state = IDLE
			}
		} else {
			if len(e.upRequests) > 0 {
				e.currentDirection = UP
				e.state = MOVING_UP
			} else if len(e.downRequests) > 0 {
				e.currentDirection = DOWN
				e.state = MOVING_DOWN
			} else {
				e.state = IDLE
			}
		}
	}
}

// ─── Dispatch Strategy Interface ──────────────────────────────────────────────

type DispatchStrategy interface {
	SelectElevator(elevators []*Elevator, requestFloor int, requestDirection Direction) int
}

// ─── NearestFirst Strategy ────────────────────────────────────────────────────

type NearestFirst struct{}

func (n *NearestFirst) SelectElevator(elevators []*Elevator, requestFloor int, requestDirection Direction) int {
	bestIdx := 0
	bestScore := math.MaxInt32
	for i, e := range elevators {
		dist := int(math.Abs(float64(e.currentFloor - requestFloor)))
		score := dist
		// TODO: penalize elevators moving in the wrong direction
		// HINT: if moving up and requestFloor < e.currentFloor, add large penalty (e.g. 1000)
		//       if moving down and requestFloor > e.currentFloor, add large penalty
		if score < bestScore {
			bestScore = score
			bestIdx = i
		}
	}
	return bestIdx
}

// ─── LeastLoaded Strategy ─────────────────────────────────────────────────────

type LeastLoaded struct{}

func (l *LeastLoaded) SelectElevator(elevators []*Elevator, requestFloor int, requestDirection Direction) int {
	bestIdx := 0
	minLoad := math.MaxInt32
	for i, e := range elevators {
		// TODO: if e.GetPendingCount() < minLoad, update bestIdx and minLoad
		_ = e
		_ = minLoad
	}
	return bestIdx
}

// ─── Elevator System ──────────────────────────────────────────────────────────

type ElevatorSystem struct {
	elevators []*Elevator
	strategy  DispatchStrategy
}

func NewElevatorSystem() *ElevatorSystem {
	return &ElevatorSystem{}
}

func (sys *ElevatorSystem) AddElevator(id int) {
	// TODO: sys.elevators = append(sys.elevators, NewElevator(id))
}

func (sys *ElevatorSystem) SetDispatchStrategy(strategy DispatchStrategy) {
	// TODO: sys.strategy = strategy
}

func (sys *ElevatorSystem) GetElevator(index int) *Elevator {
	if index < 0 || index >= len(sys.elevators) {
		return nil
	}
	return sys.elevators[index]
}

func (sys *ElevatorSystem) GetElevatorCount() int {
	return len(sys.elevators)
}

func (sys *ElevatorSystem) AddRequest(floor int, direction Direction) {
	if len(sys.elevators) == 0 {
		return
	}
	idx := 0
	if sys.strategy != nil {
		// TODO: idx = sys.strategy.SelectElevator(sys.elevators, floor, direction)
	}
	// TODO: sys.elevators[idx].AddRequest(floor, direction)
}

func (sys *ElevatorSystem) Step() {
	for _, e := range sys.elevators {
		// TODO: e.Step()
		_ = e
	}
}

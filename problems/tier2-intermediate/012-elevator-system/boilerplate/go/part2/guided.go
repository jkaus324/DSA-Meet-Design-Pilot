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

// ─── Elevator (include your Part 1 Elevator class here) ──────────────────────
// HINT: Your Elevator needs GetID(), GetCurrentFloor(), GetState(),
//       GetCurrentDirection(), and GetPendingCount() as public getters
//       so that dispatch strategies can inspect elevator state.

type Elevator struct {
	// HINT: id, currentFloor, state, currentDirection, upRequests, downRequests
}

func NewElevator(id int) *Elevator {
	return &Elevator{}
}

func (e *Elevator) GetID() int                    { return 0 }
func (e *Elevator) GetCurrentFloor() int          { return 0 }
func (e *Elevator) GetState() ElevatorState       { return IDLE }
func (e *Elevator) GetCurrentDirection() Direction { return NONE }
func (e *Elevator) GetPendingCount() int          { return 0 }

func (e *Elevator) AddRequest(floor int, direction Direction) {}
func (e *Elevator) Step()                                      {}

// ─── Dispatch Strategy Interface ──────────────────────────────────────────────
// HINT: The strategy receives a slice of all elevators plus the request
//       details, and returns the index of the best elevator.

type DispatchStrategy interface {
	SelectElevator(elevators []*Elevator, requestFloor int, requestDirection Direction) int
}

// ─── NearestFirst Strategy ────────────────────────────────────────────────────
// HINT: Calculate distance = abs(elevator.floor - requestFloor).
//       Penalize elevators moving in the wrong direction (add a large offset).
//       Prefer idle elevators or those moving toward the request.

type NearestFirst struct{}

func (n *NearestFirst) SelectElevator(elevators []*Elevator, requestFloor int, requestDirection Direction) int {
	// HINT: for each elevator compute score = distance
	//       if elevator is moving away, add 1000 as penalty
	//       return index of lowest score
	_ = math.MaxInt32
	return 0
}

// ─── LeastLoaded Strategy ─────────────────────────────────────────────────────
// HINT: Simply pick the elevator with the smallest GetPendingCount().

type LeastLoaded struct{}

func (l *LeastLoaded) SelectElevator(elevators []*Elevator, requestFloor int, requestDirection Direction) int {
	// HINT: track minLoad and bestIndex; iterate and update
	return 0
}

// ─── Elevator System ──────────────────────────────────────────────────────────
// HINT: Holds a []*Elevator and a DispatchStrategy.
//       AddRequest() uses the strategy to pick an elevator, then calls
//       that elevator's AddRequest(). Step() calls Step() on all elevators.

type ElevatorSystem struct {
	// HINT: elevators []*Elevator
	// HINT: strategy  DispatchStrategy
}

func NewElevatorSystem() *ElevatorSystem {
	return &ElevatorSystem{}
}

func (sys *ElevatorSystem) AddElevator(id int) {
	// HINT: sys.elevators = append(sys.elevators, NewElevator(id))
}

func (sys *ElevatorSystem) SetDispatchStrategy(strategy DispatchStrategy) {
	// HINT: sys.strategy = strategy
}

func (sys *ElevatorSystem) GetElevator(index int) *Elevator {
	// HINT: bounds check, return nil if invalid, else sys.elevators[index]
	return nil
}

func (sys *ElevatorSystem) GetElevatorCount() int {
	return 0
}

func (sys *ElevatorSystem) AddRequest(floor int, direction Direction) {
	// HINT: if no elevators, return
	// HINT: if strategy is nil, use index 0
	// HINT: else use strategy.SelectElevator(sys.elevators, floor, direction)
	// HINT: call sys.elevators[idx].AddRequest(floor, direction)
}

func (sys *ElevatorSystem) Step() {
	// HINT: call Step() on every elevator
}

package main

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

// ─── Your Design Starts Here ──────────────────────────────────────────────────
//
// Extend your Part 1 Elevator to support:
//   1. Multiple elevators managed by an ElevatorSystem
//   2. Pluggable dispatch strategies that decide which elevator handles
//      a new request
//   3. Two strategies: NearestFirst (nearest elevator in compatible
//      direction) and LeastLoaded (fewest pending requests)
//
// Think about:
//   - How do you define a strategy interface so new strategies can be
//     added without modifying the system?
//   - What information does a strategy need about each elevator?
//   - How does Step() work when there are multiple elevators?
//
// Entry points (must exist for tests):
//   func (sys *ElevatorSystem) AddElevator(id int)
//   func (sys *ElevatorSystem) SetDispatchStrategy(strategy DispatchStrategy)
//   func (sys *ElevatorSystem) AddRequest(floor int, direction Direction)
//   func (sys *ElevatorSystem) Step()
//   func (sys *ElevatorSystem) GetElevator(index int) *Elevator
//   func (sys *ElevatorSystem) GetElevatorCount() int
//
// ─────────────────────────────────────────────────────────────────────────────

// ─── Elevator (from Part 1) ───────────────────────────────────────────────────

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

func (e *Elevator) AddRequest(floor int, direction Direction) {}
func (e *Elevator) Step()                                      {}

// ─── Dispatch Strategy Interface ──────────────────────────────────────────────

type DispatchStrategy interface {
	SelectElevator(elevators []*Elevator, requestFloor int, requestDirection Direction) int
}

// ─── Concrete Strategies ──────────────────────────────────────────────────────

type NearestFirst struct{}

func (n *NearestFirst) SelectElevator(elevators []*Elevator, requestFloor int, requestDirection Direction) int {
	return 0
}

type LeastLoaded struct{}

func (l *LeastLoaded) SelectElevator(elevators []*Elevator, requestFloor int, requestDirection Direction) int {
	return 0
}

// ─── Elevator System ──────────────────────────────────────────────────────────

type ElevatorSystem struct {
	// TODO: add your fields here
}

func NewElevatorSystem() *ElevatorSystem {
	return &ElevatorSystem{}
}

func (sys *ElevatorSystem) AddElevator(id int)                              {}
func (sys *ElevatorSystem) SetDispatchStrategy(strategy DispatchStrategy)  {}
func (sys *ElevatorSystem) GetElevator(index int) *Elevator                 { return nil }
func (sys *ElevatorSystem) GetElevatorCount() int                           { return 0 }
func (sys *ElevatorSystem) AddRequest(floor int, direction Direction)       {}
func (sys *ElevatorSystem) Step()                                            {}

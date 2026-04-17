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
// Design and implement a single Elevator that:
//   1. Starts at floor 0 in IDLE state
//   2. Accepts external requests (floor + direction) and internal requests
//   3. Processes requests in SCAN order: serve all floors in the current
//      direction before reversing
//
// Step() behavior:
//   - IDLE + requests exist: pick direction, start moving
//   - MOVING_UP / MOVING_DOWN: move one floor; if current floor has a
//     pending request, transition to DOOR_OPEN
//   - DOOR_OPEN: close doors, resume moving or go IDLE if no requests
//
// Think about:
//   - What data structure efficiently tracks pending floors per direction?
//   - How do you decide when to reverse direction?
//   - What happens if a request arrives for the current floor while idle?
//
// Entry points (must exist for tests):
//   func (e *Elevator) AddRequest(floor int, direction Direction)
//   func (e *Elevator) Step()
//   func (e *Elevator) GetCurrentFloor() int
//   func (e *Elevator) GetState() ElevatorState
//
// ─────────────────────────────────────────────────────────────────────────────

type Elevator struct {
	// TODO: add your fields here
}

func NewElevator() *Elevator {
	return &Elevator{}
}

func (e *Elevator) GetCurrentFloor() int       { return 0 }
func (e *Elevator) GetState() ElevatorState    { return IDLE }
func (e *Elevator) GetCurrentDirection() Direction { return NONE }
func (e *Elevator) GetPendingCount() int       { return 0 }

func (e *Elevator) AddRequest(floor int, direction Direction) {}
func (e *Elevator) Step()                                      {}

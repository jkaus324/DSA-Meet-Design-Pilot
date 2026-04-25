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

// ─── Elevator ─────────────────────────────────────────────────────────────────
// HINT: Use two sorted sets (slices kept sorted) to track pending floors —
//       one for upward stops, one for downward stops.
// HINT: When moving up, check if currentFloor is in upRequests.
//       When moving down, check downRequests.
// HINT: In DOOR_OPEN, check if there are more requests in the current
//       direction. If not, check the other direction. If none, go IDLE.

type Elevator struct {
	// HINT: currentFloor      int
	// HINT: state             ElevatorState
	// HINT: currentDirection  Direction
	// HINT: upRequests   map[int]bool  // floors to visit going up
	// HINT: downRequests map[int]bool  // floors to visit going down
}

func NewElevator() *Elevator {
	// HINT: initialise all fields; currentFloor=0, state=IDLE, direction=NONE
	return &Elevator{}
}

func (e *Elevator) GetCurrentFloor() int          { return 0 }
func (e *Elevator) GetState() ElevatorState       { return IDLE }
func (e *Elevator) GetCurrentDirection() Direction { return NONE }
func (e *Elevator) GetPendingCount() int          { return 0 }

func (e *Elevator) AddRequest(floor int, direction Direction) {
	// HINT: If floor == currentFloor and state == IDLE, set state = DOOR_OPEN
	// HINT: If floor > currentFloor, add to upRequests; else add to downRequests
	// HINT: If state == IDLE and requests now exist, start moving toward nearest
}

func (e *Elevator) Step() {
	// HINT: switch on e.state:
	//   IDLE: nothing
	//   MOVING_UP: currentFloor++; if in upRequests, erase and set DOOR_OPEN
	//   MOVING_DOWN: currentFloor--; if in downRequests, erase and set DOOR_OPEN
	//   DOOR_OPEN: decide next state based on currentDirection and remaining requests
}

"""Elevator system — SCAN ordering + dispatch strategies."""


class _SortedSet:
    def __init__(self):
        self._s = []
        self._set = set()

    def add(self, x):
        if x not in self._set:
            self._set.add(x)
            # insert sorted
            lo, hi = 0, len(self._s)
            while lo < hi:
                mid = (lo + hi) // 2
                if self._s[mid] < x:
                    lo = mid + 1
                else:
                    hi = mid
            self._s.insert(lo, x)

    def remove(self, x):
        if x in self._set:
            self._set.remove(x)
            self._s.remove(x)

    def __contains__(self, x):
        return x in self._set

    def __len__(self):
        return len(self._s)

    def empty(self):
        return len(self._s) == 0

    def first(self):
        return self._s[0]

    def last(self):
        return self._s[-1]


class ElevOp:
    def __init__(self, kind, s1="", i1=0, i2=0):
        self.kind = kind
        self.s1 = s1
        self.i1 = i1
        self.i2 = i2


# State constants
IDLE = "IDLE"
MOVING_UP = "MOVING_UP"
MOVING_DOWN = "MOVING_DOWN"
DOOR_OPEN = "DOOR_OPEN"

# Direction constants
DIR_UP = "UP"
DIR_DOWN = "DOWN"
DIR_NONE = "NONE"


class Elevator:
    def __init__(self, id=0):
        self.id = id
        self.currentFloor = 0
        self.state = IDLE
        self.currentDirection = DIR_NONE
        self.upRequests = _SortedSet()
        self.downRequests = _SortedSet()

    def getId(self):
        return self.id

    def getCurrentFloor(self):
        return self.currentFloor

    def getState(self):
        return self.state

    def getCurrentDirection(self):
        return self.currentDirection

    def getPendingCount(self):
        return len(self.upRequests) + len(self.downRequests)

    def addRequest(self, floor, direction):
        if floor == self.currentFloor and self.state == IDLE:
            self.state = DOOR_OPEN
            return
        if direction == DIR_UP:
            self.upRequests.add(floor)
        elif direction == DIR_DOWN:
            self.downRequests.add(floor)
        else:
            if floor > self.currentFloor:
                self.upRequests.add(floor)
            else:
                self.downRequests.add(floor)
        if self.state == IDLE:
            if not self.upRequests.empty() and (
                self.downRequests.empty()
                or abs(self.upRequests.first() - self.currentFloor)
                <= abs(self.downRequests.last() - self.currentFloor)
            ):
                self.currentDirection = DIR_UP
                self.state = MOVING_UP
            else:
                self.currentDirection = DIR_DOWN
                self.state = MOVING_DOWN

    def step(self):
        if self.state == IDLE:
            return
        if self.state == MOVING_UP:
            self.currentFloor += 1
            if self.currentFloor in self.upRequests:
                self.upRequests.remove(self.currentFloor)
                self.state = DOOR_OPEN
            return
        if self.state == MOVING_DOWN:
            self.currentFloor -= 1
            if self.currentFloor in self.downRequests:
                self.downRequests.remove(self.currentFloor)
                self.state = DOOR_OPEN
            return
        if self.state == DOOR_OPEN:
            if self.currentDirection == DIR_UP:
                if not self.upRequests.empty():
                    self.state = MOVING_UP
                elif not self.downRequests.empty():
                    self.currentDirection = DIR_DOWN
                    self.state = MOVING_DOWN
                else:
                    self.currentDirection = DIR_NONE
                    self.state = IDLE
            elif self.currentDirection == DIR_DOWN:
                if not self.downRequests.empty():
                    self.state = MOVING_DOWN
                elif not self.upRequests.empty():
                    self.currentDirection = DIR_UP
                    self.state = MOVING_UP
                else:
                    self.currentDirection = DIR_NONE
                    self.state = IDLE
            else:
                if not self.upRequests.empty():
                    self.currentDirection = DIR_UP
                    self.state = MOVING_UP
                elif not self.downRequests.empty():
                    self.currentDirection = DIR_DOWN
                    self.state = MOVING_DOWN
                else:
                    self.state = IDLE


class NearestFirst:
    def selectElevator(self, elevators, requestFloor, requestDirection):
        bestIdx = 0
        bestScore = float("inf")
        PENALTY = 10000
        for i, e in enumerate(elevators):
            dist = abs(e.getCurrentFloor() - requestFloor)
            st = e.getState()
            d = e.getCurrentDirection()
            if st == IDLE or d == DIR_NONE:
                score = dist
            elif d == DIR_UP and requestDirection == DIR_UP and e.getCurrentFloor() <= requestFloor:
                score = dist
            elif d == DIR_DOWN and requestDirection == DIR_DOWN and e.getCurrentFloor() >= requestFloor:
                score = dist
            else:
                score = dist + PENALTY
            if score < bestScore:
                bestScore = score
                bestIdx = i
        return bestIdx


class LeastLoaded:
    def selectElevator(self, elevators, requestFloor, requestDirection):
        bestIdx = 0
        bestCount = float("inf")
        for i, e in enumerate(elevators):
            cnt = e.getPendingCount()
            if cnt < bestCount:
                bestCount = cnt
                bestIdx = i
        return bestIdx


class ElevatorSystem:
    def __init__(self):
        self.elevators = []
        self.strategy = None

    def addElevator(self, id):
        self.elevators.append(Elevator(id))

    def setDispatchStrategy(self, s):
        self.strategy = s

    def getElevator(self, index):
        if index < 0 or index >= len(self.elevators):
            return None
        return self.elevators[index]

    def getElevatorCount(self):
        return len(self.elevators)

    def addRequest(self, floor, direction):
        if not self.elevators:
            return
        idx = 0
        if self.strategy:
            idx = self.strategy.selectElevator(self.elevators, floor, direction)
        self.elevators[idx].addRequest(floor, direction)

    def step(self):
        for e in self.elevators:
            e.step()


def _dir_from(s):
    if s == "up":
        return DIR_UP
    if s == "down":
        return DIR_DOWN
    return DIR_NONE


def elevator_simulate(ops):
    out = []
    single = None
    sys = None
    nf = NearestFirst()
    ll = LeastLoaded()
    for op in ops:
        k = op.kind
        if k == "new_elev":
            single = Elevator()
            sys = None
            out.append("ok")
        elif k == "new_sys":
            sys = ElevatorSystem()
            single = None
            out.append("ok")
        elif k == "add_elev":
            sys.addElevator(op.i1)
            out.append("ok")
        elif k == "set_strategy":
            if op.s1 == "nearest":
                sys.setDispatchStrategy(nf)
            elif op.s1 == "least_loaded":
                sys.setDispatchStrategy(ll)
            out.append("ok")
        elif k == "req":
            single.addRequest(op.i1, _dir_from(op.s1))
            out.append("ok")
        elif k == "sys_req":
            sys.addRequest(op.i1, _dir_from(op.s1))
            out.append("ok")
        elif k == "elev_req":
            sys.getElevator(op.i1).addRequest(op.i2, _dir_from(op.s1))
            out.append("ok")
        elif k == "step":
            single.step()
            out.append("ok")
        elif k == "sys_step":
            sys.step()
            out.append("ok")
        elif k == "elev_step":
            sys.getElevator(op.i1).step()
            out.append("ok")
        elif k == "floor":
            out.append(str(single.getCurrentFloor()))
        elif k == "elev_floor":
            out.append(str(sys.getElevator(op.i1).getCurrentFloor()))
        elif k == "state":
            out.append(single.getState())
        elif k == "elev_state":
            out.append(sys.getElevator(op.i1).getState())
        elif k == "elev_pending":
            out.append(str(sys.getElevator(op.i1).getPendingCount()))
        elif k == "count":
            out.append(str(sys.getElevatorCount()))
        elif k == "elev_null":
            out.append("yes" if sys.getElevator(op.i1) is None else "no")
        else:
            out.append("unknown:" + k)
    return out

"""BookMyShow — theater/show/booking + seat locking + contiguous search."""


class ShowOp:
    def __init__(self, kind, s1="", s2="", s3="", s4="", i1=0, i2=0, i3=0):
        self.kind = kind
        self.s1 = s1
        self.s2 = s2
        self.s3 = s3
        self.s4 = s4
        self.i1 = i1
        self.i2 = i2
        self.i3 = i3


# Seat status constants
AVAILABLE = "AVAILABLE"
LOCKED = "LOCKED"
BOOKED = "BOOKED"


class Seat:
    def __init__(self, row, col):
        self.row = row
        self.col = col
        self.status = AVAILABLE
        self.lockedBy = ""
        self.lockExpiry = 0
        self.bookedBy = ""


class Show:
    def __init__(self, id, theaterId, movie, time, rows, cols):
        self.id = id
        self.theaterId = theaterId
        self.movie = movie
        self.time = time
        self.rows = rows
        self.cols = cols
        self.seats = [[Seat(r, c) for c in range(cols)] for r in range(rows)]


class Theater:
    def __init__(self, id, name, city):
        self.id = id
        self.name = name
        self.city = city


class BookingSystem:
    def __init__(self):
        self.theaters = {}
        self.shows = {}
        self.bookings = {}
        self.locks = {}
        self.cityMovies = {}  # city -> set of movies (use sorted set)
        self.lockCounter = 0

    def _expire_seat(self, seat, currentTime):
        if seat.status == LOCKED and currentTime >= seat.lockExpiry:
            seat.status = AVAILABLE
            seat.lockedBy = ""
            seat.lockExpiry = 0

    def addTheater(self, theaterId, name, city):
        self.theaters[theaterId] = Theater(theaterId, name, city)

    def addShow(self, showId, theaterId, movie, time, rows, cols):
        self.shows[showId] = Show(showId, theaterId, movie, time, rows, cols)
        if theaterId in self.theaters:
            city = self.theaters[theaterId].city
            self.cityMovies.setdefault(city, set()).add(movie)

    def searchMovies(self, city):
        if city not in self.cityMovies:
            return []
        # C++ uses std::set which is sorted
        return sorted(self.cityMovies[city])

    def getAvailableSeats(self, showId, currentTime=0):
        if showId not in self.shows:
            return []
        show = self.shows[showId]
        result = []
        for r in range(show.rows):
            for c in range(show.cols):
                self._expire_seat(show.seats[r][c], currentTime)
                if show.seats[r][c].status == AVAILABLE:
                    result.append((r, c))
        return result

    def bookSeats(self, bookingId, showId, seatPositions, userId, currentTime=0):
        if showId not in self.shows:
            return False
        show = self.shows[showId]
        for r, c in seatPositions:
            if r < 0 or r >= show.rows or c < 0 or c >= show.cols:
                return False
            self._expire_seat(show.seats[r][c], currentTime)
            if show.seats[r][c].status != AVAILABLE:
                return False
        for r, c in seatPositions:
            show.seats[r][c].status = BOOKED
            show.seats[r][c].bookedBy = userId
        self.bookings[bookingId] = (bookingId, showId, userId, seatPositions)
        return True

    def lockSeats(self, showId, seatPositions, userId, ttlMinutes, currentTime):
        if showId not in self.shows:
            return ""
        show = self.shows[showId]
        for r, c in seatPositions:
            if r < 0 or r >= show.rows or c < 0 or c >= show.cols:
                return ""
            self._expire_seat(show.seats[r][c], currentTime)
            if show.seats[r][c].status != AVAILABLE:
                return ""
        self.lockCounter += 1
        lockId = "LOCK_" + str(self.lockCounter)
        expiry = currentTime + ttlMinutes * 60
        for r, c in seatPositions:
            show.seats[r][c].status = LOCKED
            show.seats[r][c].lockedBy = userId
            show.seats[r][c].lockExpiry = expiry
        self.locks[lockId] = {
            "id": lockId, "showId": showId, "userId": userId,
            "seatPositions": seatPositions, "expiry": expiry,
            "confirmed": False, "released": False,
        }
        return lockId

    def confirmBooking(self, lockId, currentTime):
        if lockId not in self.locks:
            return False
        lock = self.locks[lockId]
        if lock["confirmed"] or lock["released"]:
            return False
        if currentTime >= lock["expiry"]:
            if lock["showId"] in self.shows:
                show = self.shows[lock["showId"]]
                for r, c in lock["seatPositions"]:
                    self._expire_seat(show.seats[r][c], currentTime)
            lock["released"] = True
            return False
        if lock["showId"] not in self.shows:
            return False
        show = self.shows[lock["showId"]]
        for r, c in lock["seatPositions"]:
            show.seats[r][c].status = BOOKED
            show.seats[r][c].bookedBy = lock["userId"]
            show.seats[r][c].lockedBy = ""
            show.seats[r][c].lockExpiry = 0
        bookingId = "BK_" + lockId
        self.bookings[bookingId] = (bookingId, lock["showId"], lock["userId"], lock["seatPositions"])
        lock["confirmed"] = True
        return True

    def releaseLock(self, lockId, currentTime):
        if lockId not in self.locks:
            return False
        lock = self.locks[lockId]
        if lock["confirmed"] or lock["released"]:
            return False
        if lock["showId"] in self.shows:
            show = self.shows[lock["showId"]]
            for r, c in lock["seatPositions"]:
                show.seats[r][c].status = AVAILABLE
                show.seats[r][c].lockedBy = ""
                show.seats[r][c].lockExpiry = 0
        lock["released"] = True
        return True

    def findContiguousSeats(self, showId, n, currentTime=0):
        if showId not in self.shows:
            return []
        show = self.shows[showId]
        for r in range(show.rows):
            count = 0
            start = -1
            for c in range(show.cols):
                self._expire_seat(show.seats[r][c], currentTime)
                if show.seats[r][c].status == AVAILABLE:
                    if count == 0:
                        start = c
                    count += 1
                    if count == n:
                        return [(r, i) for i in range(start, start + n)]
                else:
                    count = 0
                    start = -1
        return []


def _parse_seats(s):
    if not s:
        return []
    out = []
    for token in s.split(";"):
        if "," in token:
            parts = token.split(",")
            out.append((int(parts[0]), int(parts[1])))
    return out


def show_simulate(ops):
    out = []
    sys = BookingSystem()
    lockSlots = [""] * 32
    last_contig = []
    for op in ops:
        k = op.kind
        if k == "new":
            sys = BookingSystem()
            lockSlots = [""] * 32
            last_contig = []
            out.append("ok")
        elif k == "add_theater":
            sys.addTheater(op.s1, op.s2, op.s3)
            out.append("ok")
        elif k == "add_show":
            sys.addShow(op.s1, op.s2, op.s3, op.s4, op.i1, op.i2)
            out.append("ok")
        elif k == "movies_count":
            out.append(str(len(sys.searchMovies(op.s1))))
        elif k == "movies_contains":
            m = sys.searchMovies(op.s1)
            out.append("yes" if op.s2 in m else "no")
        elif k == "available_count":
            out.append(str(len(sys.getAvailableSeats(op.s1, op.i1))))
        elif k == "available_has":
            v = sys.getAvailableSeats(op.s1, op.i1)
            found = any(p[0] == op.i2 and p[1] == op.i3 for p in v)
            out.append("yes" if found else "no")
        elif k == "book":
            ok = sys.bookSeats(op.s1, op.s2, _parse_seats(op.s3), op.s4, op.i1)
            out.append("ok" if ok else "fail")
        elif k == "lock":
            lid = sys.lockSeats(op.s1, _parse_seats(op.s2), op.s3, op.i1, op.i2)
            if 0 <= op.i3 < len(lockSlots):
                lockSlots[op.i3] = lid
            out.append(lid)
        elif k == "lock_at":
            out.append(lockSlots[op.i3])
        elif k == "confirm":
            ok = sys.confirmBooking(lockSlots[op.i3], op.i1)
            out.append("ok" if ok else "fail")
        elif k == "release":
            ok = sys.releaseLock(lockSlots[op.i3], op.i1)
            out.append("ok" if ok else "fail")
        elif k == "release_id":
            ok = sys.releaseLock(op.s1, op.i1)
            out.append("ok" if ok else "fail")
        elif k == "find_contig":
            last_contig = sys.findContiguousSeats(op.s1, op.i1, op.i2)
            out.append(str(len(last_contig)))
        elif k == "contig_at":
            if op.i1 < 0 or op.i1 >= len(last_contig):
                out.append("")
            else:
                out.append(f"{last_contig[op.i1][0]},{last_contig[op.i1][1]}")
        else:
            out.append("unknown:" + k)
    return out

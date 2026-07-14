'use strict';

/* BookMyShow — theater/show/booking + seat locking + contiguous search. */

class ShowOp {
  constructor(kind, s1 = '', s2 = '', s3 = '', s4 = '', i1 = 0, i2 = 0, i3 = 0) {
    this.kind = kind;
    this.s1 = s1;
    this.s2 = s2;
    this.s3 = s3;
    this.s4 = s4;
    this.i1 = i1;
    this.i2 = i2;
    this.i3 = i3;
  }
}

// Seat status constants
const AVAILABLE = 'AVAILABLE';
const LOCKED = 'LOCKED';
const BOOKED = 'BOOKED';

class Seat {
  constructor(row, col) {
    this.row = row;
    this.col = col;
    this.status = AVAILABLE;
    this.lockedBy = '';
    this.lockExpiry = 0;
    this.bookedBy = '';
  }
}

class Show {
  constructor(id, theaterId, movie, time, rows, cols) {
    this.id = id;
    this.theaterId = theaterId;
    this.movie = movie;
    this.time = time;
    this.rows = rows;
    this.cols = cols;
    this.seats = [];
    for (let r = 0; r < rows; r++) {
      const row = [];
      for (let c = 0; c < cols; c++) row.push(new Seat(r, c));
      this.seats.push(row);
    }
  }
}

class Theater {
  constructor(id, name, city) {
    this.id = id;
    this.name = name;
    this.city = city;
  }
}

class BookingSystem {
  constructor() {
    this.theaters = new Map();
    this.shows = new Map();
    this.bookings = new Map();
    this.locks = new Map();
    this.cityMovies = new Map(); // city -> Set of movies
    this.lockCounter = 0;
  }

  _expire_seat(seat, currentTime) {
    if (seat.status === LOCKED && currentTime >= seat.lockExpiry) {
      seat.status = AVAILABLE;
      seat.lockedBy = '';
      seat.lockExpiry = 0;
    }
  }

  addTheater(theaterId, name, city) {
    this.theaters.set(theaterId, new Theater(theaterId, name, city));
  }

  addShow(showId, theaterId, movie, time, rows, cols) {
    this.shows.set(showId, new Show(showId, theaterId, movie, time, rows, cols));
    if (this.theaters.has(theaterId)) {
      const city = this.theaters.get(theaterId).city;
      if (!this.cityMovies.has(city)) this.cityMovies.set(city, new Set());
      this.cityMovies.get(city).add(movie);
    }
  }

  searchMovies(city) {
    if (!this.cityMovies.has(city)) return [];
    // C++ uses std::set which is sorted
    return [...this.cityMovies.get(city)].sort();
  }

  getAvailableSeats(showId, currentTime = 0) {
    if (!this.shows.has(showId)) return [];
    const show = this.shows.get(showId);
    const result = [];
    for (let r = 0; r < show.rows; r++) {
      for (let c = 0; c < show.cols; c++) {
        this._expire_seat(show.seats[r][c], currentTime);
        if (show.seats[r][c].status === AVAILABLE) {
          result.push([r, c]);
        }
      }
    }
    return result;
  }

  bookSeats(bookingId, showId, seatPositions, userId, currentTime = 0) {
    if (!this.shows.has(showId)) return false;
    const show = this.shows.get(showId);
    for (const [r, c] of seatPositions) {
      if (r < 0 || r >= show.rows || c < 0 || c >= show.cols) return false;
      this._expire_seat(show.seats[r][c], currentTime);
      if (show.seats[r][c].status !== AVAILABLE) return false;
    }
    for (const [r, c] of seatPositions) {
      show.seats[r][c].status = BOOKED;
      show.seats[r][c].bookedBy = userId;
    }
    this.bookings.set(bookingId, [bookingId, showId, userId, seatPositions]);
    return true;
  }

  lockSeats(showId, seatPositions, userId, ttlMinutes, currentTime) {
    if (!this.shows.has(showId)) return '';
    const show = this.shows.get(showId);
    for (const [r, c] of seatPositions) {
      if (r < 0 || r >= show.rows || c < 0 || c >= show.cols) return '';
      this._expire_seat(show.seats[r][c], currentTime);
      if (show.seats[r][c].status !== AVAILABLE) return '';
    }
    this.lockCounter += 1;
    const lockId = 'LOCK_' + String(this.lockCounter);
    const expiry = currentTime + ttlMinutes * 60;
    for (const [r, c] of seatPositions) {
      show.seats[r][c].status = LOCKED;
      show.seats[r][c].lockedBy = userId;
      show.seats[r][c].lockExpiry = expiry;
    }
    this.locks.set(lockId, {
      id: lockId,
      showId,
      userId,
      seatPositions,
      expiry,
      confirmed: false,
      released: false,
    });
    return lockId;
  }

  confirmBooking(lockId, currentTime) {
    if (!this.locks.has(lockId)) return false;
    const lock = this.locks.get(lockId);
    if (lock.confirmed || lock.released) return false;
    if (currentTime >= lock.expiry) {
      if (this.shows.has(lock.showId)) {
        const show = this.shows.get(lock.showId);
        for (const [r, c] of lock.seatPositions) {
          this._expire_seat(show.seats[r][c], currentTime);
        }
      }
      lock.released = true;
      return false;
    }
    if (!this.shows.has(lock.showId)) return false;
    const show = this.shows.get(lock.showId);
    for (const [r, c] of lock.seatPositions) {
      show.seats[r][c].status = BOOKED;
      show.seats[r][c].bookedBy = lock.userId;
      show.seats[r][c].lockedBy = '';
      show.seats[r][c].lockExpiry = 0;
    }
    const bookingId = 'BK_' + lockId;
    this.bookings.set(bookingId, [bookingId, lock.showId, lock.userId, lock.seatPositions]);
    lock.confirmed = true;
    return true;
  }

  releaseLock(lockId, currentTime) {
    if (!this.locks.has(lockId)) return false;
    const lock = this.locks.get(lockId);
    if (lock.confirmed || lock.released) return false;
    if (this.shows.has(lock.showId)) {
      const show = this.shows.get(lock.showId);
      for (const [r, c] of lock.seatPositions) {
        show.seats[r][c].status = AVAILABLE;
        show.seats[r][c].lockedBy = '';
        show.seats[r][c].lockExpiry = 0;
      }
    }
    lock.released = true;
    return true;
  }

  findContiguousSeats(showId, n, currentTime = 0) {
    if (!this.shows.has(showId)) return [];
    const show = this.shows.get(showId);
    for (let r = 0; r < show.rows; r++) {
      let count = 0;
      let start = -1;
      for (let c = 0; c < show.cols; c++) {
        this._expire_seat(show.seats[r][c], currentTime);
        if (show.seats[r][c].status === AVAILABLE) {
          if (count === 0) start = c;
          count += 1;
          if (count === n) {
            const res = [];
            for (let i = start; i < start + n; i++) res.push([r, i]);
            return res;
          }
        } else {
          count = 0;
          start = -1;
        }
      }
    }
    return [];
  }
}

function _parse_seats(s) {
  if (!s) return [];
  const out = [];
  for (const token of s.split(';')) {
    if (token.includes(',')) {
      const parts = token.split(',');
      out.push([parseInt(parts[0], 10), parseInt(parts[1], 10)]);
    }
  }
  return out;
}

function show_simulate(ops) {
  const out = [];
  let sys = new BookingSystem();
  let lockSlots = new Array(32).fill('');
  let lastContig = [];
  for (const op of ops) {
    const k = op.kind;
    if (k === 'new') {
      sys = new BookingSystem();
      lockSlots = new Array(32).fill('');
      lastContig = [];
      out.push('ok');
    } else if (k === 'add_theater') {
      sys.addTheater(op.s1, op.s2, op.s3);
      out.push('ok');
    } else if (k === 'add_show') {
      sys.addShow(op.s1, op.s2, op.s3, op.s4, op.i1, op.i2);
      out.push('ok');
    } else if (k === 'movies_count') {
      out.push(String(sys.searchMovies(op.s1).length));
    } else if (k === 'movies_contains') {
      const m = sys.searchMovies(op.s1);
      out.push(m.includes(op.s2) ? 'yes' : 'no');
    } else if (k === 'available_count') {
      out.push(String(sys.getAvailableSeats(op.s1, op.i1).length));
    } else if (k === 'available_has') {
      const v = sys.getAvailableSeats(op.s1, op.i1);
      const found = v.some(p => p[0] === op.i2 && p[1] === op.i3);
      out.push(found ? 'yes' : 'no');
    } else if (k === 'book') {
      const ok = sys.bookSeats(op.s1, op.s2, _parse_seats(op.s3), op.s4, op.i1);
      out.push(ok ? 'ok' : 'fail');
    } else if (k === 'lock') {
      const lid = sys.lockSeats(op.s1, _parse_seats(op.s2), op.s3, op.i1, op.i2);
      if (op.i3 >= 0 && op.i3 < lockSlots.length) lockSlots[op.i3] = lid;
      out.push(lid);
    } else if (k === 'lock_at') {
      out.push(lockSlots[op.i3]);
    } else if (k === 'confirm') {
      const ok = sys.confirmBooking(lockSlots[op.i3], op.i1);
      out.push(ok ? 'ok' : 'fail');
    } else if (k === 'release') {
      const ok = sys.releaseLock(lockSlots[op.i3], op.i1);
      out.push(ok ? 'ok' : 'fail');
    } else if (k === 'release_id') {
      const ok = sys.releaseLock(op.s1, op.i1);
      out.push(ok ? 'ok' : 'fail');
    } else if (k === 'find_contig') {
      lastContig = sys.findContiguousSeats(op.s1, op.i1, op.i2);
      out.push(String(lastContig.length));
    } else if (k === 'contig_at') {
      if (op.i1 < 0 || op.i1 >= lastContig.length) {
        out.push('');
      } else {
        out.push(`${lastContig[op.i1][0]},${lastContig[op.i1][1]}`);
      }
    } else {
      out.push('unknown:' + k);
    }
  }
  return out;
}

module.exports = { ShowOp, show_simulate };

#include <iostream>
#include <vector>
#include <string>
#include <unordered_map>
#include <set>
#include <algorithm>
using namespace std;

// ─── Data Structures ────────────────────────────────────────────────────────

enum class SeatStatus { AVAILABLE, LOCKED, BOOKED };

struct Seat {
    int row;
    int col;
    SeatStatus status;
    string lockedBy;
    long lockExpiry;
    string bookedBy;
};

struct Show {
    string id;
    string theaterId;
    string movie;
    string time;
    int rows;
    int cols;
    vector<vector<Seat>> seats;
};

struct Theater {
    string id;
    string name;
    string city;
};

struct Booking {
    string id;
    string showId;
    string userId;
    vector<pair<int,int>> seatPositions;
};

struct SeatLock {
    string id;
    string showId;
    string userId;
    vector<pair<int,int>> seatPositions;
    long expiry;
    bool confirmed;
    bool released;
};

// ─── Booking System ────────────────────────────────────────────────────────

class BookingSystem {
    unordered_map<string, Theater> theaters;
    unordered_map<string, Show> shows;
    unordered_map<string, Booking> bookings;
    unordered_map<string, SeatLock> locks;
    unordered_map<string, set<string>> cityMovies;
    int lockCounter = 0;

    void expireSeat(Seat& seat, long currentTime) {
        if (seat.status == SeatStatus::LOCKED && currentTime >= seat.lockExpiry) {
            seat.status = SeatStatus::AVAILABLE;
            seat.lockedBy = "";
            seat.lockExpiry = 0;
        }
    }

    bool isSeatAvailable(const Seat& seat, long currentTime) const {
        if (seat.status == SeatStatus::AVAILABLE) return true;
        if (seat.status == SeatStatus::LOCKED && currentTime >= seat.lockExpiry)
            return true;
        return false;
    }

public:
    void addTheater(const string& theaterId, const string& name,
                    const string& city) {
        theaters[theaterId] = {theaterId, name, city};
    }

    void addShow(const string& showId, const string& theaterId,
                 const string& movie, const string& time, int rows, int cols) {
        Show show{showId, theaterId, movie, time, rows, cols, {}};
        show.seats.resize(rows, vector<Seat>(cols));
        for (int r = 0; r < rows; r++)
            for (int c = 0; c < cols; c++)
                show.seats[r][c] = {r, c, SeatStatus::AVAILABLE, "", 0, ""};
        shows[showId] = show;
        if (theaters.count(theaterId))
            cityMovies[theaters[theaterId].city].insert(movie);
    }

    vector<string> searchMovies(const string& city) {
        if (!cityMovies.count(city)) return {};
        return vector<string>(cityMovies[city].begin(), cityMovies[city].end());
    }

    vector<pair<int,int>> getAvailableSeats(const string& showId,
                                            long currentTime = 0) {
        if (!shows.count(showId)) return {};
        auto& show = shows[showId];
        vector<pair<int,int>> result;
        for (int r = 0; r < show.rows; r++)
            for (int c = 0; c < show.cols; c++) {
                expireSeat(show.seats[r][c], currentTime);
                if (show.seats[r][c].status == SeatStatus::AVAILABLE)
                    result.push_back({r, c});
            }
        return result;
    }

    bool bookSeats(const string& bookingId, const string& showId,
                   const vector<pair<int,int>>& seatPositions,
                   const string& userId, long currentTime = 0) {
        if (!shows.count(showId)) return false;
        auto& show = shows[showId];
        // Atomic check: all seats must be available
        for (auto& [r, c] : seatPositions) {
            if (r < 0 || r >= show.rows || c < 0 || c >= show.cols)
                return false;
            expireSeat(show.seats[r][c], currentTime);
            if (show.seats[r][c].status != SeatStatus::AVAILABLE)
                return false;
        }
        // Book all seats
        for (auto& [r, c] : seatPositions) {
            show.seats[r][c].status = SeatStatus::BOOKED;
            show.seats[r][c].bookedBy = userId;
        }
        bookings[bookingId] = {bookingId, showId, userId, seatPositions};
        return true;
    }

    string lockSeats(const string& showId,
                     const vector<pair<int,int>>& seatPositions,
                     const string& userId, int ttlMinutes, long currentTime) {
        if (!shows.count(showId)) return "";
        auto& show = shows[showId];
        // Check all seats are available
        for (auto& [r, c] : seatPositions) {
            if (r < 0 || r >= show.rows || c < 0 || c >= show.cols)
                return "";
            expireSeat(show.seats[r][c], currentTime);
            if (show.seats[r][c].status != SeatStatus::AVAILABLE)
                return "";
        }
        string lockId = "LOCK_" + to_string(++lockCounter);
        long expiry = currentTime + ttlMinutes * 60;
        // Lock all seats
        for (auto& [r, c] : seatPositions) {
            show.seats[r][c].status = SeatStatus::LOCKED;
            show.seats[r][c].lockedBy = userId;
            show.seats[r][c].lockExpiry = expiry;
        }
        locks[lockId] = {lockId, showId, userId, seatPositions,
                         expiry, false, false};
        return lockId;
    }

    bool confirmBooking(const string& lockId, long currentTime) {
        if (!locks.count(lockId)) return false;
        auto& lock = locks[lockId];
        if (lock.confirmed || lock.released) return false;
        if (currentTime >= lock.expiry) {
            // Lock expired — release seats
            if (shows.count(lock.showId)) {
                auto& show = shows[lock.showId];
                for (auto& [r, c] : lock.seatPositions)
                    expireSeat(show.seats[r][c], currentTime);
            }
            lock.released = true;
            return false;
        }
        // Confirm booking
        if (!shows.count(lock.showId)) return false;
        auto& show = shows[lock.showId];
        for (auto& [r, c] : lock.seatPositions) {
            show.seats[r][c].status = SeatStatus::BOOKED;
            show.seats[r][c].bookedBy = lock.userId;
            show.seats[r][c].lockedBy = "";
            show.seats[r][c].lockExpiry = 0;
        }
        string bookingId = "BK_" + lockId;
        bookings[bookingId] = {bookingId, lock.showId, lock.userId,
                               lock.seatPositions};
        lock.confirmed = true;
        return true;
    }

    bool releaseLock(const string& lockId, long currentTime) {
        if (!locks.count(lockId)) return false;
        auto& lock = locks[lockId];
        if (lock.confirmed || lock.released) return false;
        if (shows.count(lock.showId)) {
            auto& show = shows[lock.showId];
            for (auto& [r, c] : lock.seatPositions) {
                show.seats[r][c].status = SeatStatus::AVAILABLE;
                show.seats[r][c].lockedBy = "";
                show.seats[r][c].lockExpiry = 0;
            }
        }
        lock.released = true;
        return true;
    }
};

// ─── Test Entry Points ─────────────────────────────────────────────────────

BookingSystem system_instance;

void add_theater(const string& theaterId, const string& name,
                 const string& city) {
    system_instance.addTheater(theaterId, name, city);
}

void add_show(const string& showId, const string& theaterId,
              const string& movie, const string& time, int rows, int cols) {
    system_instance.addShow(showId, theaterId, movie, time, rows, cols);
}

vector<string> search_movies(const string& city) {
    return system_instance.searchMovies(city);
}

vector<pair<int,int>> get_available_seats(const string& showId,
                                          long currentTime = 0) {
    return system_instance.getAvailableSeats(showId, currentTime);
}

bool book_seats(const string& bookingId, const string& showId,
                const vector<pair<int,int>>& seatPositions,
                const string& userId, long currentTime = 0) {
    return system_instance.bookSeats(bookingId, showId, seatPositions,
                                     userId, currentTime);
}

string lock_seats(const string& showId,
                  const vector<pair<int,int>>& seatPositions,
                  const string& userId, int ttlMinutes, long currentTime) {
    return system_instance.lockSeats(showId, seatPositions, userId,
                                     ttlMinutes, currentTime);
}

bool confirm_booking(const string& lockId, long currentTime) {
    return system_instance.confirmBooking(lockId, currentTime);
}

bool release_lock(const string& lockId, long currentTime) {
    return system_instance.releaseLock(lockId, currentTime);
}

// ─── Main ──────────────────────────────────────────────────────────────────

#ifndef RUNNING_TESTS
int main() {
    add_theater("T1", "PVR Phoenix", "Mumbai");
    add_show("S1", "T1", "Inception", "18:00", 5, 10);

    auto movies = search_movies("Mumbai");
    cout << "Movies in Mumbai: " << movies.size() << endl;

    auto seats = get_available_seats("S1");
    cout << "Available seats: " << seats.size() << endl;

    bool ok = book_seats("B1", "S1", {{0,0}, {0,1}}, "user1");
    cout << "Booked: " << (ok ? "YES" : "NO") << endl;

    return 0;
}
#endif

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

    // Part 3: Sliding window — find N contiguous available seats in earliest row
    vector<pair<int,int>> findContiguousSeats(const string& showId, int n, long currentTime = 0) {
        if (!shows.count(showId)) return {};
        auto& show = shows[showId];
        for (int r = 0; r < show.rows; r++) {
            int count = 0, start = -1;
            for (int c = 0; c < show.cols; c++) {
                expireSeat(show.seats[r][c], currentTime);
                if (show.seats[r][c].status == SeatStatus::AVAILABLE) {
                    if (count == 0) start = c;
                    count++;
                    if (count == n) {
                        vector<pair<int,int>> result;
                        for (int i = start; i < start + n; i++) result.push_back({r, i});
                        return result;
                    }
                } else {
                    count = 0;
                    start = -1;
                }
            }
        }
        return {};
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

// ─── Ops simulator (used by spec-based tests) ──────────────────────────────

#include <memory>

struct ShowOp {
    string kind;
    string s1;
    string s2;
    string s3;
    string s4;
    int    i1;
    int    i2;
    int    i3;
};

static vector<pair<int,int>> parse_seats(const string& s) {
    vector<pair<int,int>> out;
    string cur;
    auto flush = [&](){
        if (cur.empty()) return;
        size_t comma = cur.find(',');
        if (comma == string::npos) { cur.clear(); return; }
        int r = stoi(cur.substr(0, comma));
        int c = stoi(cur.substr(comma+1));
        out.push_back({r, c});
        cur.clear();
    };
    for (char ch : s) {
        if (ch == ';') flush();
        else cur.push_back(ch);
    }
    flush();
    return out;
}

vector<string> show_simulate(vector<ShowOp> ops) {
    vector<string> out;
    unique_ptr<BookingSystem> sys(new BookingSystem());
    vector<string> lockSlots(32, "");
    vector<pair<int,int>> last_contig;
    for (const auto& op : ops) {
        const string& k = op.kind;
        if (k == "new") {
            sys.reset(new BookingSystem());
            for (auto& s : lockSlots) s.clear();
            last_contig.clear();
            out.push_back("ok");
        } else if (k == "add_theater") {
            sys->addTheater(op.s1, op.s2, op.s3); out.push_back("ok");
        } else if (k == "add_show") {
            sys->addShow(op.s1, op.s2, op.s3, op.s4, op.i1, op.i2); out.push_back("ok");
        } else if (k == "movies_count") {
            out.push_back(to_string((int)sys->searchMovies(op.s1).size()));
        } else if (k == "movies_contains") {
            auto m = sys->searchMovies(op.s1);
            out.push_back(find(m.begin(), m.end(), op.s2) != m.end() ? "yes" : "no");
        } else if (k == "available_count") {
            out.push_back(to_string((int)sys->getAvailableSeats(op.s1, (long)op.i1).size()));
        } else if (k == "available_has") {
            auto v = sys->getAvailableSeats(op.s1, (long)op.i1);
            bool found = false;
            for (auto& p : v) if (p.first == op.i2 && p.second == op.i3) { found = true; break; }
            out.push_back(found ? "yes" : "no");
        } else if (k == "book") {
            out.push_back(sys->bookSeats(op.s1, op.s2, parse_seats(op.s3), op.s4, (long)op.i1) ? "ok" : "fail");
        } else if (k == "lock") {
            string lid = sys->lockSeats(op.s1, parse_seats(op.s2), op.s3, op.i1, (long)op.i2);
            if (op.i3 >= 0 && op.i3 < (int)lockSlots.size()) lockSlots[op.i3] = lid;
            out.push_back(lid);
        } else if (k == "lock_at") {
            out.push_back(lockSlots[op.i3]);
        } else if (k == "confirm") {
            out.push_back(sys->confirmBooking(lockSlots[op.i3], (long)op.i1) ? "ok" : "fail");
        } else if (k == "release") {
            out.push_back(sys->releaseLock(lockSlots[op.i3], (long)op.i1) ? "ok" : "fail");
        } else if (k == "release_id") {
            out.push_back(sys->releaseLock(op.s1, (long)op.i1) ? "ok" : "fail");
        } else if (k == "find_contig") {
            last_contig = sys->findContiguousSeats(op.s1, op.i1, (long)op.i2);
            out.push_back(to_string((int)last_contig.size()));
        } else if (k == "contig_at") {
            if (op.i1 < 0 || op.i1 >= (int)last_contig.size()) out.push_back("");
            else out.push_back(to_string(last_contig[op.i1].first) + "," + to_string(last_contig[op.i1].second));
        } else {
            out.push_back("unknown:" + k);
        }
    }
    return out;
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

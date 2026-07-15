// Part 3 Tests — Contiguous Seat Search (Sliding Window)
// Tests findContiguousSeats: locate N adjacent available seats in the same row

#include "solution.cpp"
#include <cassert>
#include <iostream>
using namespace std;

int part3_tests() {
    int passed = 0, failed = 0;

    // Test 1: Find 2 contiguous seats in a row with all seats available
    try {
        BookingSystem sys;
        sys.addTheater("T1", "PVR", "Mumbai");
        sys.addShow("S1", "T1", "Movie", "18:00", 3, 5);
        auto seats = sys.findContiguousSeats("S1", 2);
        assert(seats.size() == 2u);
        assert(seats[0].first == seats[1].first);        // same row
        assert(seats[0].second + 1 == seats[1].second); // contiguous
        cout << "PASS test_find_2_contiguous_available" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_find_2_contiguous_available" << endl;
        failed++;
    }

    // Test 2: Returns earliest (row 0) first
    try {
        BookingSystem sys;
        sys.addTheater("T1", "PVR", "Mumbai");
        sys.addShow("S1", "T1", "Movie", "18:00", 3, 5);
        auto seats = sys.findContiguousSeats("S1", 3);
        assert(seats.size() == 3u);
        assert(seats[0].first == 0); // row 0 is earliest
        cout << "PASS test_returns_earliest_row" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_returns_earliest_row" << endl;
        failed++;
    }

    // Test 3: Returns empty when N exceeds total seats in any single row
    try {
        BookingSystem sys;
        sys.addTheater("T1", "PVR", "Mumbai");
        sys.addShow("S1", "T1", "Movie", "18:00", 2, 3); // 3 cols per row
        auto seats = sys.findContiguousSeats("S1", 4);   // need 4, but only 3 per row
        assert(seats.empty());
        cout << "PASS test_returns_empty_n_exceeds_row_width" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_returns_empty_n_exceeds_row_width" << endl;
        failed++;
    }

    // Test 4: Unknown showId returns empty
    try {
        BookingSystem sys;
        auto seats = sys.findContiguousSeats("UNKNOWN", 2);
        assert(seats.empty());
        cout << "PASS test_unknown_show_returns_empty" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_unknown_show_returns_empty" << endl;
        failed++;
    }

    // Test 5: N=1 returns a single available seat
    try {
        BookingSystem sys;
        sys.addTheater("T1", "PVR", "Mumbai");
        sys.addShow("S1", "T1", "Movie", "18:00", 2, 4);
        auto seats = sys.findContiguousSeats("S1", 1);
        assert(seats.size() == 1u);
        cout << "PASS test_n_equals_1_returns_single_seat" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_n_equals_1_returns_single_seat" << endl;
        failed++;
    }

    // Test 6: Skips booked seats and finds next contiguous block
    try {
        BookingSystem sys;
        sys.addTheater("T1", "PVR", "Mumbai");
        sys.addShow("S1", "T1", "Movie", "18:00", 1, 5);
        // Book seat (0,0) directly
        sys.bookSeats("B1", "S1", {{0, 0}}, "user1");
        // Row 0: [booked, free, free, free, free] — 2-contiguous starts at col 1
        auto seats = sys.findContiguousSeats("S1", 2);
        assert(seats.size() == 2u);
        assert(seats[0].second == 1); // must start at col 1 (skip booked col 0)
        cout << "PASS test_skips_booked_seats" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_skips_booked_seats" << endl;
        failed++;
    }

    // Test 7: Expired lock makes seat available again
    try {
        BookingSystem sys;
        sys.addTheater("T1", "PVR", "Mumbai");
        sys.addShow("S1", "T1", "Movie", "18:00", 1, 2);
        // Lock both seats with ttl=1 min, at time=0 → expiry = 60
        sys.lockSeats("S1", {{0,0},{0,1}}, "user1", 1, 0);
        // At time=61, locks expired — both seats back to available
        auto seats = sys.findContiguousSeats("S1", 2, 61);
        assert(seats.size() == 2u);
        cout << "PASS test_expired_lock_seat_available" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_expired_lock_seat_available" << endl;
        failed++;
    }

    cout << "PART3_SUMMARY " << passed << "/" << (passed + failed) << endl;
    return failed;
}

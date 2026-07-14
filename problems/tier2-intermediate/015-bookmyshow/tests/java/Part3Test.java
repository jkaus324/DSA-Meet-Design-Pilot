// Part 3 Tests — BookMyShow: Search Movies + Release Lock

import java.util.*;

class Part3Test {
    public static int runTests() {
        int failures = 0;
        int total = 8;

        // Test 1: searchMovies returns movies in a city
        try {
            BookingSystem sys = new BookingSystem();
            sys.addTheater("T1", "PVR", "Mumbai");
            sys.addShow("S1", "T1", "Inception", "18:00", 5, 10);
            sys.addShow("S2", "T1", "Avengers", "21:00", 5, 10);
            List<String> movies = sys.searchMovies("Mumbai");
            if (movies == null || movies.size() != 2) throw new AssertionError("expected 2 movies, got " + (movies == null ? "null" : movies.size()));
            System.out.println("PASS test_search_movies_basic");
        } catch (Exception e) {
            System.out.println("FAIL test_search_movies_basic: " + e.getMessage());
            failures++;
        }

        // Test 2: searchMovies returns empty for unknown city
        try {
            BookingSystem sys = new BookingSystem();
            sys.addTheater("T1", "PVR", "Mumbai");
            sys.addShow("S1", "T1", "Inception", "18:00", 5, 10);
            List<String> movies = sys.searchMovies("Delhi");
            if (movies == null || !movies.isEmpty()) throw new AssertionError("expected empty for unknown city");
            System.out.println("PASS test_search_movies_unknown_city");
        } catch (Exception e) {
            System.out.println("FAIL test_search_movies_unknown_city: " + e.getMessage());
            failures++;
        }

        // Test 3: searchMovies returns sorted results
        try {
            BookingSystem sys = new BookingSystem();
            sys.addTheater("T1", "PVR", "Bangalore");
            sys.addShow("S1", "T1", "Zara", "18:00", 3, 5);
            sys.addShow("S2", "T1", "Avengers", "20:00", 3, 5);
            sys.addShow("S3", "T1", "Matrix", "22:00", 3, 5);
            List<String> movies = sys.searchMovies("Bangalore");
            if (movies == null || movies.size() != 3) throw new AssertionError("expected 3 movies");
            if (!movies.get(0).equals("Avengers")) throw new AssertionError("expected Avengers first, got " + movies.get(0));
            if (!movies.get(1).equals("Matrix")) throw new AssertionError("expected Matrix second");
            if (!movies.get(2).equals("Zara")) throw new AssertionError("expected Zara third");
            System.out.println("PASS test_search_movies_sorted");
        } catch (Exception e) {
            System.out.println("FAIL test_search_movies_sorted: " + e.getMessage());
            failures++;
        }

        // Test 4: releaseLock frees seats
        try {
            BookingSystem sys = new BookingSystem();
            sys.addTheater("T1", "PVR", "Mumbai");
            sys.addShow("S1", "T1", "Inception", "18:00", 2, 2);
            long now = 1000;
            String lockId = sys.lockSeats("S1", Arrays.asList(new int[]{0, 0}), "user1", 5, now);
            if (lockId == null || lockId.isEmpty()) throw new AssertionError("lock failed");
            List<int[]> before = sys.getAvailableSeats("S1", now);
            if (before.size() != 3) throw new AssertionError("expected 3 available before release");
            boolean released = sys.releaseLock(lockId, now);
            if (!released) throw new AssertionError("releaseLock should return true");
            List<int[]> after = sys.getAvailableSeats("S1", now);
            if (after.size() != 4) throw new AssertionError("expected 4 available after release, got " + after.size());
            System.out.println("PASS test_release_lock_frees_seats");
        } catch (Exception e) {
            System.out.println("FAIL test_release_lock_frees_seats: " + e.getMessage());
            failures++;
        }

        // Test 5: releaseLock allows other user to lock same seat
        try {
            BookingSystem sys = new BookingSystem();
            sys.addTheater("T1", "PVR", "Mumbai");
            sys.addShow("S1", "T1", "Inception", "18:00", 2, 2);
            long now = 1000;
            String lockId1 = sys.lockSeats("S1", Arrays.asList(new int[]{0, 0}), "user1", 5, now);
            sys.releaseLock(lockId1, now);
            String lockId2 = sys.lockSeats("S1", Arrays.asList(new int[]{0, 0}), "user2", 5, now);
            if (lockId2 == null || lockId2.isEmpty()) throw new AssertionError("user2 should be able to lock after release");
            System.out.println("PASS test_release_allows_relock");
        } catch (Exception e) {
            System.out.println("FAIL test_release_allows_relock: " + e.getMessage());
            failures++;
        }

        // Test 6: releaseLock on confirmed booking returns false
        try {
            BookingSystem sys = new BookingSystem();
            sys.addTheater("T1", "PVR", "Mumbai");
            sys.addShow("S1", "T1", "Inception", "18:00", 2, 2);
            long now = 1000;
            String lockId = sys.lockSeats("S1", Arrays.asList(new int[]{0, 0}), "user1", 5, now);
            sys.confirmBooking(lockId, now + 100);
            boolean released = sys.releaseLock(lockId, now + 200);
            if (released) throw new AssertionError("cannot release confirmed booking");
            System.out.println("PASS test_release_confirmed_booking_fails");
        } catch (Exception e) {
            System.out.println("FAIL test_release_confirmed_booking_fails: " + e.getMessage());
            failures++;
        }

        // Test 7: releaseLock on already-released lock returns false
        try {
            BookingSystem sys = new BookingSystem();
            sys.addTheater("T1", "PVR", "Mumbai");
            sys.addShow("S1", "T1", "Inception", "18:00", 2, 2);
            long now = 1000;
            String lockId = sys.lockSeats("S1", Arrays.asList(new int[]{0, 0}), "user1", 5, now);
            sys.releaseLock(lockId, now);
            boolean released2 = sys.releaseLock(lockId, now);
            if (released2) throw new AssertionError("double release should return false");
            System.out.println("PASS test_release_already_released");
        } catch (Exception e) {
            System.out.println("FAIL test_release_already_released: " + e.getMessage());
            failures++;
        }

        // Test 8: Multi-city search
        try {
            BookingSystem sys = new BookingSystem();
            sys.addTheater("T1", "PVR", "Mumbai");
            sys.addTheater("T2", "INOX", "Delhi");
            sys.addShow("S1", "T1", "Inception", "18:00", 3, 5);
            sys.addShow("S2", "T2", "Matrix", "20:00", 3, 5);
            sys.addShow("S3", "T2", "Avengers", "22:00", 3, 5);
            List<String> mumbai = sys.searchMovies("Mumbai");
            List<String> delhi = sys.searchMovies("Delhi");
            if (mumbai == null || mumbai.size() != 1) throw new AssertionError("Mumbai: expected 1 movie");
            if (delhi == null || delhi.size() != 2) throw new AssertionError("Delhi: expected 2 movies");
            if (!mumbai.get(0).equals("Inception")) throw new AssertionError("Mumbai movie wrong");
            System.out.println("PASS test_multi_city_search");
        } catch (Exception e) {
            System.out.println("FAIL test_multi_city_search: " + e.getMessage());
            failures++;
        }

        int passed = total - failures;
        System.out.println("PART3_SUMMARY " + passed + "/" + total);
        return failures;
    }
}

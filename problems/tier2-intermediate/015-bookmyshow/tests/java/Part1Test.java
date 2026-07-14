// Bookmyshow — Part 1 Tests
import java.util.*;
import java.util.stream.*;

class Part1Test {
    static boolean testSearchMoviesByCity() {
        try {
            BookingSystem sys = new BookingSystem();
            sys.addTheater("T1", "PVR Phoenix", "Mumbai");
            sys.addShow("S1", "T1", "Inception", "18:00", 5, 10);
            var movies = sys.searchMovies("Mumbai");
            boolean pass = movies.size() == 1
                && movies[0] == "Inception";
            System.out.println((pass ? "PASS" : "FAIL") + ": testSearchMoviesByCity");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testSearchMoviesByCity (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testMultipleMoviesSameCity() {
        try {
            BookingSystem sys = new BookingSystem();
            sys.addTheater("T1", "PVR Phoenix", "Mumbai");
            sys.addTheater("T2", "INOX", "Mumbai");
            sys.addShow("S1", "T1", "Inception", "18:00", 5, 10);
            sys.addShow("S2", "T2", "Interstellar", "20:00", 5, 10);
            var movies = sys.searchMovies("Mumbai");
            // TreeSet<String> ensures sorted order
            TreeSet<String> movieSet(movies.begin(), movies.end());
            boolean pass = movies.size() == 2
                && (movieSet.containsKey("Inception") ? 1 : 0)
                && (movieSet.containsKey("Interstellar") ? 1 : 0);
            System.out.println((pass ? "PASS" : "FAIL") + ": testMultipleMoviesSameCity");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testMultipleMoviesSameCity (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testNoMoviesInCity() {
        try {
            BookingSystem sys = new BookingSystem();
            sys.addTheater("T1", "PVR", "Mumbai");
            sys.addShow("S1", "T1", "Inception", "18:00", 5, 10);
            var movies = sys.searchMovies("Delhi");
            boolean pass = movies.isEmpty();
            System.out.println((pass ? "PASS" : "FAIL") + ": testNoMoviesInCity");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testNoMoviesInCity (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testAllSeatsAvailable() {
        try {
            BookingSystem sys = new BookingSystem();
            sys.addTheater("T1", "PVR", "Mumbai");
            sys.addShow("S1", "T1", "Inception", "18:00", 3, 4);
            var seats = sys.getAvailableSeats("S1");
            boolean pass = seats.size() == 12);  // 3 rows x 4 cols;
            System.out.println((pass ? "PASS" : "FAIL") + ": testAllSeatsAvailable");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testAllSeatsAvailable (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testBookSeatsSuccess() {
        try {
            BookingSystem sys = new BookingSystem();
            sys.addTheater("T1", "PVR", "Mumbai");
            sys.addShow("S1", "T1", "Inception", "18:00", 5, 10);
            boolean ok = sys.bookSeats("B1", "S1", {Arrays.asList(0,0), Arrays.asList(0,1), Arrays.asList(0,2)}, "user1");
            var seats = sys.getAvailableSeats("S1");
            boolean pass = ok == true
                && seats.size() == 47);  // 50 - 3 booked;
            System.out.println((pass ? "PASS" : "FAIL") + ": testBookSeatsSuccess");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testBookSeatsSuccess (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testDoubleBookingPrevention() {
        try {
            BookingSystem sys = new BookingSystem();
            sys.addTheater("T1", "PVR", "Mumbai");
            sys.addShow("S1", "T1", "Inception", "18:00", 5, 10);
            sys.bookSeats("B1", "S1", {Arrays.asList(0,0), Arrays.asList(0,1)}, "user1");
            // Try to book same seat (0,0)
            boolean ok = sys.bookSeats("B2", "S1", {Arrays.asList(0,0), Arrays.asList(0,2)}, "user2");
            boolean pass = ok == false;
            System.out.println((pass ? "PASS" : "FAIL") + ": testDoubleBookingPrevention");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testDoubleBookingPrevention (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testAtomicBooking() {
        try {
            BookingSystem sys = new BookingSystem();
            sys.addTheater("T1", "PVR", "Mumbai");
            sys.addShow("S1", "T1", "Inception", "18:00", 5, 10);
            sys.bookSeats("B1", "S1", {Arrays.asList(0,0)}, "user1");
            // Try to book Arrays.asList(0,0) (taken) and Arrays.asList(0,1) (free) — should fail entirely
            boolean ok = sys.bookSeats("B2", "S1", {Arrays.asList(0,0), Arrays.asList(0,1)}, "user2");
            // Verify Arrays.asList(0,1) is still available (atomic — nothing was booked)
            var seats = sys.getAvailableSeats("S1");
            boolean found01 = false;
            for (var _e_seats_ : seats.entrySet()) {
            var r = _e_seats_.getKey(); var c = _e_seats_.getValue();
            if (r == 0 & c == 1) found01 = true;
            }
            boolean pass = ok == false
                && found01;
            System.out.println((pass ? "PASS" : "FAIL") + ": testAtomicBooking");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testAtomicBooking (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testInvalidSeatCoordinates() {
        try {
            BookingSystem sys = new BookingSystem();
            sys.addTheater("T1", "PVR", "Mumbai");
            sys.addShow("S1", "T1", "Inception", "18:00", 3, 4);
            // Row 5 doesn't exist in a 3-row theater
            boolean ok = sys.bookSeats("B1", "S1", {Arrays.asList(5, 0)}, "user1");
            boolean pass = ok == false;
            System.out.println((pass ? "PASS" : "FAIL") + ": testInvalidSeatCoordinates");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testInvalidSeatCoordinates (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testNonexistentShow() {
        try {
            BookingSystem sys = new BookingSystem();
            boolean ok = sys.bookSeats("B1", "NONEXISTENT", {Arrays.asList(0,0)}, "user1");
            boolean pass = ok == false;
            System.out.println((pass ? "PASS" : "FAIL") + ": testNonexistentShow");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testNonexistentShow (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testDuplicateMovieDeduplicated() {
        try {
            BookingSystem sys = new BookingSystem();
            sys.addTheater("T1", "PVR Phoenix", "Mumbai");
            sys.addTheater("T2", "INOX", "Mumbai");
            sys.addShow("S1", "T1", "Inception", "15:00", 5, 10);
            sys.addShow("S2", "T2", "Inception", "18:00", 5, 10);
            var movies = sys.searchMovies("Mumbai");
            boolean pass = movies.size() == 1);  // "Inception" only once;
            System.out.println((pass ? "PASS" : "FAIL") + ": testDuplicateMovieDeduplicated");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testDuplicateMovieDeduplicated (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    public static int runTests() {
        int passed = 0, total = 0;
        total++; if (testSearchMoviesByCity()) passed++;
        total++; if (testMultipleMoviesSameCity()) passed++;
        total++; if (testNoMoviesInCity()) passed++;
        total++; if (testAllSeatsAvailable()) passed++;
        total++; if (testBookSeatsSuccess()) passed++;
        total++; if (testDoubleBookingPrevention()) passed++;
        total++; if (testAtomicBooking()) passed++;
        total++; if (testInvalidSeatCoordinates()) passed++;
        total++; if (testNonexistentShow()) passed++;
        total++; if (testDuplicateMovieDeduplicated()) passed++;
        System.out.println("PART1_SUMMARY " + passed + "/" + total);
        return passed;
    }
}

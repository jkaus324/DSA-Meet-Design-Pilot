// Bookmyshow — Part 2 Tests
import java.util.*;
import java.util.stream.*;

class Part2Test {
    static boolean testLockSeatsSuccess() {
        try {
            BookingSystem sys = new BookingSystem();
            sys.addTheater("T1", "PVR", "Mumbai");
            sys.addShow("S1", "T1", "Inception", "18:00", 5, 10);
            long now = 1000;
            String lockId = sys.lockSeats("S1", {Arrays.asList(0,0), Arrays.asList(0,1)}, "user1", 5, now);
            // Locked seats should not be available
            var seats = sys.getAvailableSeats("S1", now);
            boolean pass = !lockId.isEmpty()
                && seats.size() == 48);  // 50 - 2 locked;
            System.out.println((pass ? "PASS" : "FAIL") + ": testLockSeatsSuccess");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testLockSeatsSuccess (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testCannotLockLockedSeats() {
        try {
            BookingSystem sys = new BookingSystem();
            sys.addTheater("T1", "PVR", "Mumbai");
            sys.addShow("S1", "T1", "Inception", "18:00", 5, 10);
            long now = 1000;
            sys.lockSeats("S1", {Arrays.asList(0,0), Arrays.asList(0,1)}, "user1", 5, now);
            // Another user tries to lock same seats
            String lockId2 = sys.lockSeats("S1", {Arrays.asList(0,0)}, "user2", 5, now);
            boolean pass = lockId2.isEmpty();
            System.out.println((pass ? "PASS" : "FAIL") + ": testCannotLockLockedSeats");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testCannotLockLockedSeats (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testConfirmWithinTtl() {
        try {
            BookingSystem sys = new BookingSystem();
            sys.addTheater("T1", "PVR", "Mumbai");
            sys.addShow("S1", "T1", "Inception", "18:00", 5, 10);
            long now = 1000;
            String lockId = sys.lockSeats("S1", {Arrays.asList(0,0), Arrays.asList(0,1)}, "user1", 5, now);
            // Confirm within TTL (5 min = 300 sec)
            boolean ok = sys.confirmBooking(lockId, now + 200);
            // Seats should now be BOOKED, not available
            var seats = sys.getAvailableSeats("S1", now + 200);
            boolean pass = ok == true
                && seats.size() == 48;
            System.out.println((pass ? "PASS" : "FAIL") + ": testConfirmWithinTtl");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testConfirmWithinTtl (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testConfirmAfterTtlExpires() {
        try {
            BookingSystem sys = new BookingSystem();
            sys.addTheater("T1", "PVR", "Mumbai");
            sys.addShow("S1", "T1", "Inception", "18:00", 5, 10);
            long now = 1000;
            String lockId = sys.lockSeats("S1", {Arrays.asList(0,0), Arrays.asList(0,1)}, "user1", 5, now);
            // Try to confirm after TTL (5 min = 300 sec, try at +400)
            boolean ok = sys.confirmBooking(lockId, now + 400);
            // Seats should be available again (expired)
            var seats = sys.getAvailableSeats("S1", now + 400);
            boolean pass = ok == false
                && seats.size() == 50);  // All seats free;
            System.out.println((pass ? "PASS" : "FAIL") + ": testConfirmAfterTtlExpires");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testConfirmAfterTtlExpires (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testExpiredLockReleasesSeats() {
        try {
            BookingSystem sys = new BookingSystem();
            sys.addTheater("T1", "PVR", "Mumbai");
            sys.addShow("S1", "T1", "Inception", "18:00", 5, 10);
            long now = 1000;
            sys.lockSeats("S1", {Arrays.asList(0,0)}, "user1", 5, now);
            // After TTL, another user can lock the same seat
            long afterExpiry = now + 400;
            String lockId2 = sys.lockSeats("S1", {Arrays.asList(0,0)}, "user2", 5, afterExpiry);
            boolean pass = !lockId2.isEmpty();
            System.out.println((pass ? "PASS" : "FAIL") + ": testExpiredLockReleasesSeats");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testExpiredLockReleasesSeats (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testManualRelease() {
        try {
            BookingSystem sys = new BookingSystem();
            sys.addTheater("T1", "PVR", "Mumbai");
            sys.addShow("S1", "T1", "Inception", "18:00", 5, 10);
            long now = 1000;
            String lockId = sys.lockSeats("S1", {Arrays.asList(0,0), Arrays.asList(0,1)}, "user1", 5, now);
            boolean released = sys.releaseLock(lockId, now + 60);
            // Seats should be available again
            var seats = sys.getAvailableSeats("S1", now + 60);
            boolean pass = released == true
                && seats.size() == 50;
            System.out.println((pass ? "PASS" : "FAIL") + ": testManualRelease");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testManualRelease (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testNoConfirmAfterRelease() {
        try {
            BookingSystem sys = new BookingSystem();
            sys.addTheater("T1", "PVR", "Mumbai");
            sys.addShow("S1", "T1", "Inception", "18:00", 5, 10);
            long now = 1000;
            String lockId = sys.lockSeats("S1", {Arrays.asList(0,0)}, "user1", 5, now);
            sys.releaseLock(lockId, now + 60);
            boolean ok = sys.confirmBooking(lockId, now + 120);
            boolean pass = ok == false;
            System.out.println((pass ? "PASS" : "FAIL") + ": testNoConfirmAfterRelease");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testNoConfirmAfterRelease (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testNoDoubleConfirm() {
        try {
            BookingSystem sys = new BookingSystem();
            sys.addTheater("T1", "PVR", "Mumbai");
            sys.addShow("S1", "T1", "Inception", "18:00", 5, 10);
            long now = 1000;
            String lockId = sys.lockSeats("S1", {Arrays.asList(0,0)}, "user1", 5, now);
            boolean ok1 = sys.confirmBooking(lockId, now + 100);
            boolean ok2 = sys.confirmBooking(lockId, now + 150);
            boolean pass = ok1 == true
                && ok2 == false;
            System.out.println((pass ? "PASS" : "FAIL") + ": testNoDoubleConfirm");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testNoDoubleConfirm (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testCannotBookLockedSeat() {
        try {
            BookingSystem sys = new BookingSystem();
            sys.addTheater("T1", "PVR", "Mumbai");
            sys.addShow("S1", "T1", "Inception", "18:00", 5, 10);
            long now = 1000;
            sys.lockSeats("S1", {Arrays.asList(0,0)}, "user1", 5, now);
            // Try direct booking of locked seat
            boolean ok = sys.bookSeats("B1", "S1", {Arrays.asList(0,0)}, "user2", now + 60);
            boolean pass = ok == false;
            System.out.println((pass ? "PASS" : "FAIL") + ": testCannotBookLockedSeat");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testCannotBookLockedSeat (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testReleaseNonexistentLock() {
        try {
            BookingSystem sys = new BookingSystem();
            boolean ok = sys.releaseLock("NONEXISTENT", 1000);
            boolean pass = ok == false;
            System.out.println((pass ? "PASS" : "FAIL") + ": testReleaseNonexistentLock");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testReleaseNonexistentLock (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    public static int runTests() {
        int passed = 0, total = 0;
        total++; if (testLockSeatsSuccess()) passed++;
        total++; if (testCannotLockLockedSeats()) passed++;
        total++; if (testConfirmWithinTtl()) passed++;
        total++; if (testConfirmAfterTtlExpires()) passed++;
        total++; if (testExpiredLockReleasesSeats()) passed++;
        total++; if (testManualRelease()) passed++;
        total++; if (testNoConfirmAfterRelease()) passed++;
        total++; if (testNoDoubleConfirm()) passed++;
        total++; if (testCannotBookLockedSeat()) passed++;
        total++; if (testReleaseNonexistentLock()) passed++;
        System.out.println("PART2_SUMMARY " + passed + "/" + total);
        return passed;
    }
}

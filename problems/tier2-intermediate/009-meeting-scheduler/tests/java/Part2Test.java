// Meeting Scheduler — Part 2 Tests
import java.util.*;
import java.util.stream.*;

class Part2Test {
    static boolean testFirstAvailable() {
        try {
            MeetingScheduler s = new MeetingScheduler();
            s.addRoom({"R1", "Small", 4, false});
            s.addRoom({"R2", "Medium", 10, false});
            s.addRoom({"R3", "Large", 20, true});
            FirstAvailable fa = new FirstAvailable();
            s.setStrategy( fa);
            String roomId = s.bookWithStrategy("M1", "Standup", 540, 570, 3);
            boolean pass = roomId == "R1");  // first room that fits 3 people;
            System.out.println((pass ? "PASS" : "FAIL") + ": testFirstAvailable");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testFirstAvailable (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testFirstAvailableSkipOccupied() {
        try {
            MeetingScheduler s = new MeetingScheduler();
            s.addRoom({"R1", "Small", 4, false});
            s.addRoom({"R2", "Medium", 10, false});
            FirstAvailable fa = new FirstAvailable();
            s.setStrategy( fa);
            s.bookMeeting(Arrays.asList("M1", "Standup", 540, 570, "R1"));
            String roomId = s.bookWithStrategy("M2", "Planning", 540, 570, 3);
            boolean pass = roomId == "R2");  // R1 is occupied, falls to R2;
            System.out.println((pass ? "PASS" : "FAIL") + ": testFirstAvailableSkipOccupied");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testFirstAvailableSkipOccupied (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testBestFit() {
        try {
            MeetingScheduler s = new MeetingScheduler();
            s.addRoom({"R1", "Tiny", 2, false});
            s.addRoom({"R2", "Small", 6, false});
            s.addRoom({"R3", "Large", 20, true});
            BestFit bf = new BestFit();
            s.setStrategy( bf);
            String roomId = s.bookWithStrategy("M1", "Meeting", 540, 570, 5);
            boolean pass = roomId == "R2");  // R1 too small (2 < 5), R2 fits (6 >= 5), R3 too big;
            System.out.println((pass ? "PASS" : "FAIL") + ": testBestFit");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testBestFit (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testBestFitNoRoom() {
        try {
            MeetingScheduler s = new MeetingScheduler();
            s.addRoom({"R1", "Tiny", 2, false});
            s.addRoom({"R2", "Small", 4, false});
            BestFit bf = new BestFit();
            s.setStrategy( bf);
            String roomId = s.bookWithStrategy("M1", "Big Meeting", 540, 570, 10);
            boolean pass = roomId == "");  // no room fits 10 people;
            System.out.println((pass ? "PASS" : "FAIL") + ": testBestFitNoRoom");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testBestFitNoRoom (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testPriorityPrefersAv() {
        try {
            MeetingScheduler s = new MeetingScheduler();
            s.addRoom({"R1", "Small", 6, false});
            s.addRoom({"R2", "Medium AV", 10, true});
            s.addRoom({"R3", "Large AV", 20, true});
            PriorityBased pb = new PriorityBased();
            s.setStrategy( pb);
            String roomId = s.bookWithStrategy("M1", "Presentation", 540, 570, 5);
            boolean pass = roomId == "R2");  // smallest AV room that fits;
            System.out.println((pass ? "PASS" : "FAIL") + ": testPriorityPrefersAv");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testPriorityPrefersAv (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testPriorityFallbackNonAv() {
        try {
            MeetingScheduler s = new MeetingScheduler();
            s.addRoom({"R1", "Small", 6, false});
            s.addRoom({"R2", "Medium AV", 10, true});
            PriorityBased pb = new PriorityBased();
            s.setStrategy( pb);
            // Occupy the AV room
            s.bookMeeting(Arrays.asList("M0", "Existing", 540, 570, "R2"));
            String roomId = s.bookWithStrategy("M1", "Meeting", 540, 570, 5);
            boolean pass = roomId == "R1");  // AV room occupied, fall back to non-AV;
            System.out.println((pass ? "PASS" : "FAIL") + ": testPriorityFallbackNonAv");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testPriorityFallbackNonAv (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    public static int runTests() {
        int passed = 0, total = 0;
        total++; if (testFirstAvailable()) passed++;
        total++; if (testFirstAvailableSkipOccupied()) passed++;
        total++; if (testBestFit()) passed++;
        total++; if (testBestFitNoRoom()) passed++;
        total++; if (testPriorityPrefersAv()) passed++;
        total++; if (testPriorityFallbackNonAv()) passed++;
        System.out.println("PART2_SUMMARY " + passed + "/" + total);
        return passed;
    }
}

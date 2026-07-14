// Meeting Scheduler — Part 1 Tests
import java.util.*;
import java.util.stream.*;

class Part1Test {
    static boolean testBookMeeting() {
        try {
            MeetingScheduler s = new MeetingScheduler();
            s.addRoom({"R1", "Small Room", 4, false});
            Meeting mArrays.asList("M1", "Standup", 540, 570, "R1");
            boolean pass = s.bookMeeting(m) == true;
            System.out.println((pass ? "PASS" : "FAIL") + ": testBookMeeting");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testBookMeeting (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testConflictDetection() {
        try {
            MeetingScheduler s = new MeetingScheduler();
            s.addRoom({"R1", "Small Room", 4, false});
            s.bookMeeting(Arrays.asList("M1", "Standup", 540, 570, "R1"));
            // Overlapping: starts at 550, before M1 ends at 570
            boolean pass = s.bookMeeting(Arrays.asList("M2", "Planning", 550, 600, "R1")) == false;
            System.out.println((pass ? "PASS" : "FAIL") + ": testConflictDetection");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testConflictDetection (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testAdjacentNoConflict() {
        try {
            MeetingScheduler s = new MeetingScheduler();
            s.addRoom({"R1", "Small Room", 4, false});
            s.bookMeeting(Arrays.asList("M1", "Standup", 540, 570, "R1"));
            // Starts exactly when M1 ends — no overlap
            boolean pass = s.bookMeeting(Arrays.asList("M2", "Planning", 570, 630, "R1")) == true;
            System.out.println((pass ? "PASS" : "FAIL") + ": testAdjacentNoConflict");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testAdjacentNoConflict (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testDifferentRoomsNoConflict() {
        try {
            MeetingScheduler s = new MeetingScheduler();
            s.addRoom({"R1", "Small Room", 4, false});
            s.addRoom({"R2", "Large Room", 20, true});
            s.bookMeeting(Arrays.asList("M1", "Standup", 540, 570, "R1"));
            // Same time, different room — should succeed
            boolean pass = s.bookMeeting(Arrays.asList("M2", "Planning", 540, 570, "R2")) == true;
            System.out.println((pass ? "PASS" : "FAIL") + ": testDifferentRoomsNoConflict");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testDifferentRoomsNoConflict (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testIsAvailableFree() {
        try {
            MeetingScheduler s = new MeetingScheduler();
            s.addRoom({"R1", "Small Room", 4, false});
            boolean pass = s.isAvailable("R1", 540, 570) == true;
            System.out.println((pass ? "PASS" : "FAIL") + ": testIsAvailableFree");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testIsAvailableFree (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testIsAvailableOccupied() {
        try {
            MeetingScheduler s = new MeetingScheduler();
            s.addRoom({"R1", "Small Room", 4, false});
            s.bookMeeting(Arrays.asList("M1", "Standup", 540, 570, "R1"));
            boolean pass = s.isAvailable("R1", 550, 580) == false;
            System.out.println((pass ? "PASS" : "FAIL") + ": testIsAvailableOccupied");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testIsAvailableOccupied (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testScheduleSorted() {
        try {
            MeetingScheduler s = new MeetingScheduler();
            s.addRoom({"R1", "Small Room", 4, false});
            s.bookMeeting(Arrays.asList("M2", "Planning", 600, 660, "R1"));
            s.bookMeeting(Arrays.asList("M1", "Standup", 540, 570, "R1"));
            var sched = s.getRoomSchedule("R1");
            boolean pass = sched.size() == 2
                && sched[0].id == "M1");  // earlier meeting first
                && sched[1].id == "M2";
            System.out.println((pass ? "PASS" : "FAIL") + ": testScheduleSorted");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testScheduleSorted (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testNonexistentRoom() {
        try {
            MeetingScheduler s = new MeetingScheduler();
            boolean pass = s.bookMeeting(Arrays.asList("M1", "Standup", 540, 570, "R99")) == false;
            System.out.println((pass ? "PASS" : "FAIL") + ": testNonexistentRoom");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testNonexistentRoom (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testEmptySchedule() {
        try {
            MeetingScheduler s = new MeetingScheduler();
            s.addRoom({"R1", "Small Room", 4, false});
            var sched = s.getRoomSchedule("R1");
            boolean pass = sched.isEmpty();
            System.out.println((pass ? "PASS" : "FAIL") + ": testEmptySchedule");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testEmptySchedule (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    public static int runTests() {
        int passed = 0, total = 0;
        total++; if (testBookMeeting()) passed++;
        total++; if (testConflictDetection()) passed++;
        total++; if (testAdjacentNoConflict()) passed++;
        total++; if (testDifferentRoomsNoConflict()) passed++;
        total++; if (testIsAvailableFree()) passed++;
        total++; if (testIsAvailableOccupied()) passed++;
        total++; if (testScheduleSorted()) passed++;
        total++; if (testNonexistentRoom()) passed++;
        total++; if (testEmptySchedule()) passed++;
        System.out.println("PART1_SUMMARY " + passed + "/" + total);
        return passed;
    }
}

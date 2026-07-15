// Meeting Scheduler — Part 3 Tests
import java.util.*;
import java.util.stream.*;

class TestObserver implements MeetingObserver {
    int bookedCount = 0;
    int cancelledCount = 0;
    int rescheduledCount = 0;
    String lastBookedMeetingId;
    String lastCancelledMeetingId;
    String lastRescheduledMeetingId;
    int lastNewStart = 0;
    int lastNewEnd = 0;
    void onMeetingBooked(Meeting meeting)  {
        bookedCount++;
        lastBookedMeetingId = meeting.id;
    }
    void onMeetingCancelled(Meeting meeting)  {
        cancelledCount++;
        lastCancelledMeetingId = meeting.id;
    }
    void onMeetingRescheduled(Meeting oldMeeting,
                              Meeting newMeeting)  {
        rescheduledCount++;
        lastRescheduledMeetingId = newMeeting.id;
        lastNewStart = newMeeting.startTime;
        lastNewEnd = newMeeting.endTime;
    }
}

class Part3Test {
    static boolean testObserverOnBook() {
        try {
            MeetingScheduler s = new MeetingScheduler();
            s.addRoom({"R1", "Small", 4, false});
            TestObserver obs = new TestObserver();
            s.subscribeAttendee("M1",  obs);
            s.bookMeeting(Arrays.asList("M1", "Standup", 540, 570, "R1"));
            boolean pass = obs.bookedCount == 1
                && obs.lastBookedMeetingId == "M1";
            System.out.println((pass ? "PASS" : "FAIL") + ": testObserverOnBook");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testObserverOnBook (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testObserverOnCancel() {
        try {
            MeetingScheduler s = new MeetingScheduler();
            s.addRoom({"R1", "Small", 4, false});
            TestObserver obs = new TestObserver();
            s.subscribeAttendee("M1",  obs);
            s.bookMeeting(Arrays.asList("M1", "Standup", 540, 570, "R1"));
            s.cancelMeeting("M1");
            // Room should be free again
            boolean pass = obs.cancelledCount == 1
                && obs.lastCancelledMeetingId == "M1"
                && s.isAvailable("R1", 540, 570) == true;
            System.out.println((pass ? "PASS" : "FAIL") + ": testObserverOnCancel");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testObserverOnCancel (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testObserverOnReschedule() {
        try {
            MeetingScheduler s = new MeetingScheduler();
            s.addRoom({"R1", "Small", 4, false});
            TestObserver obs = new TestObserver();
            s.subscribeAttendee("M1",  obs);
            s.bookMeeting(Arrays.asList("M1", "Standup", 540, 570, "R1"));
            boolean ok = s.rescheduleMeeting("M1", 600, 630);
            // Old slot should be free
            // New slot should be occupied
            boolean pass = ok == true
                && obs.rescheduledCount == 1
                && obs.lastNewStart == 600
                && obs.lastNewEnd == 630
                && s.isAvailable("R1", 540, 570) == true
                && s.isAvailable("R1", 600, 630) == false;
            System.out.println((pass ? "PASS" : "FAIL") + ": testObserverOnReschedule");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testObserverOnReschedule (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testRescheduleConflict() {
        try {
            MeetingScheduler s = new MeetingScheduler();
            s.addRoom({"R1", "Small", 4, false});
            s.bookMeeting(Arrays.asList("M1", "Standup", 540, 570, "R1"));
            s.bookMeeting(Arrays.asList("M2", "Planning", 600, 660, "R1"));
            boolean ok = s.rescheduleMeeting("M1", 610, 650);
            // M1 should still be in original slot
            boolean pass = ok == false);  // conflicts with M2
                && s.isAvailable("R1", 540, 570) == false;
            System.out.println((pass ? "PASS" : "FAIL") + ": testRescheduleConflict");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testRescheduleConflict (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testMultipleObservers() {
        try {
            MeetingScheduler s = new MeetingScheduler();
            s.addRoom({"R1", "Small", 4, false});
            TestObserver obs1, obs2;
            s.subscribeAttendee("M1",  obs1);
            s.subscribeAttendee("M1",  obs2);
            s.bookMeeting(Arrays.asList("M1", "Standup", 540, 570, "R1"));
            boolean pass = obs1.bookedCount == 1
                && obs2.bookedCount == 1;
            System.out.println((pass ? "PASS" : "FAIL") + ": testMultipleObservers");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testMultipleObservers (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testCancelNonexistent() {
        try {
            MeetingScheduler s = new MeetingScheduler();
            boolean pass = s.cancelMeeting("M99") == false;
            System.out.println((pass ? "PASS" : "FAIL") + ": testCancelNonexistent");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testCancelNonexistent (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    public static int runTests() {
        int passed = 0, total = 0;
        total++; if (testObserverOnBook()) passed++;
        total++; if (testObserverOnCancel()) passed++;
        total++; if (testObserverOnReschedule()) passed++;
        total++; if (testRescheduleConflict()) passed++;
        total++; if (testMultipleObservers()) passed++;
        total++; if (testCancelNonexistent()) passed++;
        System.out.println("PART3_SUMMARY " + passed + "/" + total);
        return passed;
    }
}

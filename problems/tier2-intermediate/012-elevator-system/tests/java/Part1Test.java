// Elevator System — Part 1 Tests
import java.util.*;
import java.util.stream.*;

class Part1Test {
    static boolean testInitialState() {
        try {
            Elevator e = new Elevator();
            boolean pass = e.getCurrentFloor() == 0
                && e.getState() == ElevatorState.IDLE;
            System.out.println((pass ? "PASS" : "FAIL") + ": testInitialState");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testInitialState (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testAddUpwardRequest() {
        try {
            Elevator e = new Elevator();
            e.addRequest(3, Direction.UP);
            boolean pass = e.getState() == ElevatorState.MOVING_UP;
            System.out.println((pass ? "PASS" : "FAIL") + ": testAddUpwardRequest");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testAddUpwardRequest (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testStepMovesOneFloor() {
        try {
            Elevator e = new Elevator();
            e.addRequest(3, Direction.UP);
            e.step(); // floor 0 . 1
            e.step(); // floor 1 . 2
            boolean pass = e.getCurrentFloor() == 1
                && e.getCurrentFloor() == 2;
            System.out.println((pass ? "PASS" : "FAIL") + ": testStepMovesOneFloor");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testStepMovesOneFloor (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testDoorOpensAtTarget() {
        try {
            Elevator e = new Elevator();
            e.addRequest(2, Direction.UP);
            e.step(); // floor 0 . 1
            e.step(); // floor 1 . 2, door opens
            boolean pass = e.getCurrentFloor() == 2
                && e.getState() == ElevatorState.DOOR_OPEN;
            System.out.println((pass ? "PASS" : "FAIL") + ": testDoorOpensAtTarget");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testDoorOpensAtTarget (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testIdleAfterLastRequest() {
        try {
            Elevator e = new Elevator();
            e.addRequest(1, Direction.UP);
            e.step(); // floor 0 . 1, door opens
            e.step(); // close doors, go idle
            boolean pass = e.getState() == ElevatorState.DOOR_OPEN
                && e.getState() == ElevatorState.IDLE
                && e.getCurrentFloor() == 1;
            System.out.println((pass ? "PASS" : "FAIL") + ": testIdleAfterLastRequest");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testIdleAfterLastRequest (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testScanUpwardOrder() {
        try {
            Elevator e = new Elevator();
            e.addRequest(5, Direction.UP);
            e.addRequest(2, Direction.UP);
            // Should go up: 0.1.2(stop).3.4.5(stop)
            e.step(); // 1
            e.step(); // 2, DOOR_OPEN
            e.step(); // close doors, resume MOVING_UP
            e.step(); // 3
            e.step(); // 4
            e.step(); // 5, DOOR_OPEN
            boolean pass = e.getCurrentFloor() == 2
                && e.getState() == ElevatorState.DOOR_OPEN
                && e.getState() == ElevatorState.MOVING_UP
                && e.getCurrentFloor() == 5
                && e.getState() == ElevatorState.DOOR_OPEN;
            System.out.println((pass ? "PASS" : "FAIL") + ": testScanUpwardOrder");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testScanUpwardOrder (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testScanReversal() {
        try {
            Elevator e = new Elevator();
            e.addRequest(3, Direction.UP);
            e.addRequest(1, Direction.DOWN); // below current, goes into downRequests
            // At floor 0: upRequests={3}, downRequests={} — wait, 1 > 0 so it goes to upRequests
            // Let's use a scenario where we add a downward request after moving
            e.addRequest(3, Direction.UP);
            // Step to floor 3
            e.step(); // 1
            e.step(); // 2
            e.step(); // 3, DOOR_OPEN (first request for 3 served)
            // Now add a downward request
            e.addRequest(1, Direction.DOWN); // 1 < 3, goes to downRequests
            e.step(); // close doors, should switch to MOVING_DOWN
            e.step(); // 2
            e.step(); // 1, DOOR_OPEN
            boolean pass = e.getState() == ElevatorState.DOOR_OPEN
                && e.getState() == ElevatorState.MOVING_DOWN
                && e.getCurrentFloor() == 1
                && e.getState() == ElevatorState.DOOR_OPEN;
            System.out.println((pass ? "PASS" : "FAIL") + ": testScanReversal");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testScanReversal (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testRequestAtCurrentFloor() {
        try {
            Elevator e = new Elevator();
            e.addRequest(0, Direction.UP);
            boolean pass = e.getState() == ElevatorState.DOOR_OPEN
                && e.getCurrentFloor() == 0;
            System.out.println((pass ? "PASS" : "FAIL") + ": testRequestAtCurrentFloor");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testRequestAtCurrentFloor (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testIdleStepNoop() {
        try {
            Elevator e = new Elevator();
            e.step();
            e.step();
            boolean pass = e.getCurrentFloor() == 0
                && e.getState() == ElevatorState.IDLE
                && e.getCurrentFloor() == 0;
            System.out.println((pass ? "PASS" : "FAIL") + ": testIdleStepNoop");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testIdleStepNoop (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testMultipleStopsInOrder() {
        try {
            Elevator e = new Elevator();
            e.addRequest(5, Direction.UP);
            e.addRequest(3, Direction.UP);
            e.addRequest(1, Direction.UP);
            // Should stop at 1, 3, 5
            e.step(); // floor 1, DOOR_OPEN
            e.step(); // close, MOVING_UP
            e.step(); // floor 2
            e.step(); // floor 3, DOOR_OPEN
            e.step(); // close, MOVING_UP
            e.step(); // floor 4
            e.step(); // floor 5, DOOR_OPEN
            boolean pass = e.getCurrentFloor() == 1
                && e.getState() == ElevatorState.DOOR_OPEN
                && e.getCurrentFloor() == 3
                && e.getState() == ElevatorState.DOOR_OPEN
                && e.getCurrentFloor() == 5
                && e.getState() == ElevatorState.DOOR_OPEN;
            System.out.println((pass ? "PASS" : "FAIL") + ": testMultipleStopsInOrder");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testMultipleStopsInOrder (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    public static int runTests() {
        int passed = 0, total = 0;
        total++; if (testInitialState()) passed++;
        total++; if (testAddUpwardRequest()) passed++;
        total++; if (testStepMovesOneFloor()) passed++;
        total++; if (testDoorOpensAtTarget()) passed++;
        total++; if (testIdleAfterLastRequest()) passed++;
        total++; if (testScanUpwardOrder()) passed++;
        total++; if (testScanReversal()) passed++;
        total++; if (testRequestAtCurrentFloor()) passed++;
        total++; if (testIdleStepNoop()) passed++;
        total++; if (testMultipleStopsInOrder()) passed++;
        System.out.println("PART1_SUMMARY " + passed + "/" + total);
        return passed;
    }
}

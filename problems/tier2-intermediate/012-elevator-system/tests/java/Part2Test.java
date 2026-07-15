// Elevator System — Part 2 Tests
import java.util.*;
import java.util.stream.*;

class Part2Test {
    static boolean testAddElevators() {
        try {
            ElevatorSystem sys = new ElevatorSystem();
            sys.addElevator(1);
            sys.addElevator(2);
            boolean pass = sys.getElevatorCount() == 2
                && sys.getElevator(0) != null
                && sys.getElevator(1) != null;
            System.out.println((pass ? "PASS" : "FAIL") + ": testAddElevators");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testAddElevators (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testDefaultDispatch() {
        try {
            ElevatorSystem sys = new ElevatorSystem();
            sys.addElevator(1);
            sys.addElevator(2);
            sys.addRequest(5, Direction.UP);
            boolean pass = sys.getElevator(0).getState() != ElevatorState.IDLE
                && sys.getElevator(1).getState() == ElevatorState.IDLE;
            System.out.println((pass ? "PASS" : "FAIL") + ": testDefaultDispatch");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testDefaultDispatch (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testNearestFirstDispatch() {
        try {
            ElevatorSystem sys = new ElevatorSystem();
            NearestFirst nf = new NearestFirst();
            sys.addElevator(1);
            sys.addElevator(2);
            sys.setDispatchStrategy( nf);
            // Move elevator 0 to floor 5
            sys.getElevator(0).addRequest(5, Direction.UP);
            for (int i = 0; i < 5; i++) sys.getElevator(0).step();
            // Elevator 0 is at floor 5 (DOOR_OPEN), elevator 1 is at floor 0
            sys.getElevator(0).step(); // close doors, go idle at 5
            // Request floor 2 — elevator 1 (at 0) is closer than elevator 0 (at 5)
            sys.addRequest(2, Direction.UP);
            boolean pass = sys.getElevator(1).getState() != ElevatorState.IDLE;
            System.out.println((pass ? "PASS" : "FAIL") + ": testNearestFirstDispatch");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testNearestFirstDispatch (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testLeastLoadedDispatch() {
        try {
            ElevatorSystem sys = new ElevatorSystem();
            LeastLoaded ll = new LeastLoaded();
            sys.addElevator(1);
            sys.addElevator(2);
            sys.setDispatchStrategy( ll);
            // Give elevator 0 three requests
            sys.getElevator(0).addRequest(3, Direction.UP);
            sys.getElevator(0).addRequest(5, Direction.UP);
            sys.getElevator(0).addRequest(7, Direction.UP);
            // Elevator 0 has 3 pending, elevator 1 has 0
            sys.addRequest(4, Direction.UP);
            // Should go to elevator 1 (least loaded)
            boolean pass = sys.getElevator(1).getPendingCount() > 0;
            System.out.println((pass ? "PASS" : "FAIL") + ": testLeastLoadedDispatch");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testLeastLoadedDispatch (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testStepAllElevators() {
        try {
            ElevatorSystem sys = new ElevatorSystem();
            sys.addElevator(1);
            sys.addElevator(2);
            sys.getElevator(0).addRequest(3, Direction.UP);
            sys.getElevator(1).addRequest(2, Direction.UP);
            sys.step(); // both move one floor
            boolean pass = sys.getElevator(0).getCurrentFloor() == 1
                && sys.getElevator(1).getCurrentFloor() == 1;
            System.out.println((pass ? "PASS" : "FAIL") + ": testStepAllElevators");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testStepAllElevators (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testNearestPrefersSameDirection() {
        try {
            ElevatorSystem sys = new ElevatorSystem();
            NearestFirst nf = new NearestFirst();
            sys.addElevator(1);
            sys.addElevator(2);
            sys.setDispatchStrategy( nf);
            // Move elevator 0 to floor 3 going up (give it request for floor 8)
            sys.getElevator(0).addRequest(3, Direction.UP);
            sys.getElevator(0).addRequest(8, Direction.UP);
            for (int i = 0; i < 3; i++) sys.getElevator(0).step();
            // Elevator 0 at floor 3, DOOR_OPEN, still has request for 8 (moving up)
            sys.getElevator(0).step(); // close doors, MOVING_UP toward 8
            // Elevator 1 at floor 0
            // Request floor 5 UP — elevator 0 is moving up past it, should be preferred
            sys.addRequest(5, Direction.UP);
            // Elevator 0 should get it (moving up, will pass floor 5)
            boolean pass = sys.getElevator(0).getPendingCount() >= 2); // has floor 8 + floor 5;
            System.out.println((pass ? "PASS" : "FAIL") + ": testNearestPrefersSameDirection");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testNearestPrefersSameDirection (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testSwapStrategyRuntime() {
        try {
            ElevatorSystem sys = new ElevatorSystem();
            NearestFirst nf = new NearestFirst();
            LeastLoaded ll = new LeastLoaded();
            sys.addElevator(1);
            sys.addElevator(2);
            sys.setDispatchStrategy( nf);
            sys.addRequest(3, Direction.UP); // dispatched via NearestFirst
            sys.setDispatchStrategy( ll);
            // Now give elevator 0 more requests so elevator 1 is least loaded
            sys.getElevator(0).addRequest(5, Direction.UP);
            sys.getElevator(0).addRequest(7, Direction.UP);
            sys.addRequest(4, Direction.UP); // dispatched via LeastLoaded to elevator 1
            boolean pass = sys.getElevator(1).getPendingCount() > 0;
            System.out.println((pass ? "PASS" : "FAIL") + ": testSwapStrategyRuntime");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testSwapStrategyRuntime (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testInvalidElevatorIndex() {
        try {
            ElevatorSystem sys = new ElevatorSystem();
            sys.addElevator(1);
            boolean pass = sys.getElevator(-1) == null
                && sys.getElevator(5) == null;
            System.out.println((pass ? "PASS" : "FAIL") + ": testInvalidElevatorIndex");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testInvalidElevatorIndex (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    public static int runTests() {
        int passed = 0, total = 0;
        total++; if (testAddElevators()) passed++;
        total++; if (testDefaultDispatch()) passed++;
        total++; if (testNearestFirstDispatch()) passed++;
        total++; if (testLeastLoadedDispatch()) passed++;
        total++; if (testStepAllElevators()) passed++;
        total++; if (testNearestPrefersSameDirection()) passed++;
        total++; if (testSwapStrategyRuntime()) passed++;
        total++; if (testInvalidElevatorIndex()) passed++;
        System.out.println("PART2_SUMMARY " + passed + "/" + total);
        return passed;
    }
}

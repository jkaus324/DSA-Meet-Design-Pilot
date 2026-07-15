// Vending Machine — Part 2 Tests
import java.util.*;
import java.util.stream.*;

class Part2Test {
    static boolean testEnterMaintenanceValidPin() {
        try {
            reset();
            enterMaintenance("1234");
            boolean pass = getState() == "Maintenance";
            System.out.println((pass ? "PASS" : "FAIL") + ": testEnterMaintenanceValidPin");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testEnterMaintenanceValidPin (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testEnterMaintenanceInvalidPin() {
        try {
            reset();
            enterMaintenance("wrong");
            boolean pass = getState() == "Idle"); // should stay idle;
            System.out.println((pass ? "PASS" : "FAIL") + ": testEnterMaintenanceInvalidPin");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testEnterMaintenanceInvalidPin (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testUserBlockedInMaintenance() {
        try {
            reset();
            enterMaintenance("1234");
            selectItem("Cola"); // should print warning, not crash
            boolean pass = getState() == "Maintenance"); // stays in maintenance;
            System.out.println((pass ? "PASS" : "FAIL") + ": testUserBlockedInMaintenance");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testUserBlockedInMaintenance (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testExitMaintenance() {
        try {
            reset();
            enterMaintenance("1234");
            exitMaintenance("1234");
            boolean pass = getState() == "Idle";
            System.out.println((pass ? "PASS" : "FAIL") + ": testExitMaintenance");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testExitMaintenance (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testRestockInMaintenance() {
        try {
            reset();
            restock("Cola", 10); // should be blocked outside maintenance
            enterMaintenance("1234");
            restock("Cola", 5); // should work in maintenance
            boolean pass = true; // FIXME: no assertions found
            System.out.println((pass ? "PASS" : "FAIL") + ": testRestockInMaintenance");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testRestockInMaintenance (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    public static int runTests() {
        int passed = 0, total = 0;
        total++; if (testEnterMaintenanceValidPin()) passed++;
        total++; if (testEnterMaintenanceInvalidPin()) passed++;
        total++; if (testUserBlockedInMaintenance()) passed++;
        total++; if (testExitMaintenance()) passed++;
        total++; if (testRestockInMaintenance()) passed++;
        System.out.println("PART2_SUMMARY " + passed + "/" + total);
        return passed;
    }
}

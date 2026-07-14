// Amazon Locker — Part 1 Tests
import java.util.*;
import java.util.stream.*;

class Part1Test {
    static boolean testDepositExactSize() {
        try {
            initLockerSystem();
            addLocker("S1", LockerSize.SMALL);
            addLocker("M1", LockerSize.MEDIUM);
            addLocker("L1", LockerSize.LARGE);
            String code = depositPackage("pkg1", LockerSize.SMALL);
            boolean pass = !code.isEmpty();
            System.out.println((pass ? "PASS" : "FAIL") + ": testDepositExactSize");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testDepositExactSize (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testRetrieveValidCode() {
        try {
            initLockerSystem();
            addLocker("S1", LockerSize.SMALL);
            String code = depositPackage("pkg1", LockerSize.SMALL);
            boolean ok = retrievePackage(code);
            boolean pass = !code.isEmpty()
                && ok == true;
            System.out.println((pass ? "PASS" : "FAIL") + ": testRetrieveValidCode");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testRetrieveValidCode (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testRetrieveInvalidCode() {
        try {
            initLockerSystem();
            addLocker("S1", LockerSize.SMALL);
            depositPackage("pkg1", LockerSize.SMALL);
            boolean ok = retrievePackage("INVALID-CODE");
            boolean pass = ok == false;
            System.out.println((pass ? "PASS" : "FAIL") + ": testRetrieveInvalidCode");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testRetrieveInvalidCode (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testSmallestFitAllocation() {
        try {
            initLockerSystem();
            addLocker("M1", LockerSize.MEDIUM);
            addLocker("S1", LockerSize.SMALL);
            addLocker("L1", LockerSize.LARGE);
            // Small package should go to S1 (smallest fit)
            String code1 = depositPackage("pkg1", LockerSize.SMALL);
            // Another small package should go to M1 (next smallest available)
            String code2 = depositPackage("pkg2", LockerSize.SMALL);
            // Another small package should go to L1
            String code3 = depositPackage("pkg3", LockerSize.SMALL);
            // No more lockers — should fail
            String code4 = depositPackage("pkg4", LockerSize.SMALL);
            boolean pass = !code1.isEmpty()
                && !code2.isEmpty()
                && !code3.isEmpty()
                && code4.isEmpty();
            System.out.println((pass ? "PASS" : "FAIL") + ": testSmallestFitAllocation");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testSmallestFitAllocation (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testLargePackageNoSmallLocker() {
        try {
            initLockerSystem();
            addLocker("S1", LockerSize.SMALL);
            addLocker("S2", LockerSize.SMALL);
            String code = depositPackage("bigpkg", LockerSize.LARGE);
            boolean pass = code.isEmpty()); // no large lockers available;
            System.out.println((pass ? "PASS" : "FAIL") + ": testLargePackageNoSmallLocker");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testLargePackageNoSmallLocker (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testRetrieveFreesLocker() {
        try {
            initLockerSystem();
            addLocker("S1", LockerSize.SMALL);
            String code1 = depositPackage("pkg1", LockerSize.SMALL);
            // Locker is full — next deposit fails
            String code2 = depositPackage("pkg2", LockerSize.SMALL);
            // Retrieve pkg1 — locker freed
            // Now deposit should succeed again
            String code3 = depositPackage("pkg3", LockerSize.SMALL);
            boolean pass = !code1.isEmpty()
                && code2.isEmpty()
                && retrievePackage(code1) == true
                && !code3.isEmpty();
            System.out.println((pass ? "PASS" : "FAIL") + ": testRetrieveFreesLocker");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testRetrieveFreesLocker (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testCodeSingleUse() {
        try {
            initLockerSystem();
            addLocker("S1", LockerSize.SMALL);
            String code = depositPackage("pkg1", LockerSize.SMALL);
            boolean pass = retrievePackage(code) == true
                && retrievePackage(code) == false); // already retrieved;
            System.out.println((pass ? "PASS" : "FAIL") + ": testCodeSingleUse");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testCodeSingleUse (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    public static int runTests() {
        int passed = 0, total = 0;
        total++; if (testDepositExactSize()) passed++;
        total++; if (testRetrieveValidCode()) passed++;
        total++; if (testRetrieveInvalidCode()) passed++;
        total++; if (testSmallestFitAllocation()) passed++;
        total++; if (testLargePackageNoSmallLocker()) passed++;
        total++; if (testRetrieveFreesLocker()) passed++;
        total++; if (testCodeSingleUse()) passed++;
        System.out.println("PART1_SUMMARY " + passed + "/" + total);
        return passed;
    }
}

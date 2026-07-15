// Amazon Locker — Part 2 Tests
import java.util.*;
import java.util.stream.*;

class TestNotificationChannel extends NotificationChannel {
    void notify(String packageId, String message)  {
        notification_log.add(packageId + ": " + message);
    }
}

class Part2Test {
    static boolean testDepositNotification() {
        try {
            initLockerSystem();
            notification_log.clear();
            TestNotificationChannel channel = new TestNotificationChannel();
            addNotificationChannel( channel);
            addLocker("S1", LockerSize.SMALL);
            String code = depositPackage("pkg1", LockerSize.SMALL, 1000);
            boolean pass = !code.isEmpty()
                && notification_log.size() >= 1); // should have been notified;
            System.out.println((pass ? "PASS" : "FAIL") + ": testDepositNotification");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testDepositNotification (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testCodeExpiry() {
        try {
            initLockerSystem();
            notification_log.clear();
            addLocker("S1", LockerSize.SMALL);
            setCodeExpiry(2); // 2 hours = 7200 seconds
            String code = depositPackage("pkg1", LockerSize.SMALL, 1000);
            // Check at 1000 + 7201 = 8201 (past expiry)
            List<String> expired = checkExpired(8201);
            boolean pass = !code.isEmpty()
                && expired.size() == 1
                && expired[0] == "pkg1";
            System.out.println((pass ? "PASS" : "FAIL") + ": testCodeExpiry");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testCodeExpiry (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testNonExpiredCode() {
        try {
            initLockerSystem();
            addLocker("S1", LockerSize.SMALL);
            setCodeExpiry(2);
            String code = depositPackage("pkg1", LockerSize.SMALL, 1000);
            // Check at 1000 + 3600 = 4600 (before expiry)
            List<String> expired = checkExpired(4600);
            // Code should still work
            boolean pass = expired.isEmpty()
                && retrievePackage(code) == true;
            System.out.println((pass ? "PASS" : "FAIL") + ": testNonExpiredCode");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testNonExpiredCode (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testExpiredLockerReuse() {
        try {
            initLockerSystem();
            addLocker("S1", LockerSize.SMALL);
            setCodeExpiry(1); // 1 hour = 3600 seconds
            depositPackage("pkg1", LockerSize.SMALL, 1000);
            // Can't deposit — locker is full
            String code2 = depositPackage("pkg2", LockerSize.SMALL, 1500);
            // Expire the first package
            checkExpired(5000); // 1000 + 3600 < 5000
            // Now locker should be free
            String code3 = depositPackage("pkg3", LockerSize.SMALL, 5001);
            boolean pass = code2.isEmpty()
                && !code3.isEmpty();
            System.out.println((pass ? "PASS" : "FAIL") + ": testExpiredLockerReuse");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testExpiredLockerReuse (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testExpiryNotification() {
        try {
            initLockerSystem();
            notification_log.clear();
            TestNotificationChannel channel = new TestNotificationChannel();
            addNotificationChannel( channel);
            addLocker("S1", LockerSize.SMALL);
            setCodeExpiry(1);
            depositPackage("pkg1", LockerSize.SMALL, 1000);
            int before = notification_log.size();
            checkExpired(5000);
            boolean pass = notification_log.size() > before); // expiry notification sent;
            System.out.println((pass ? "PASS" : "FAIL") + ": testExpiryNotification");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testExpiryNotification (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testSelectiveExpiry() {
        try {
            initLockerSystem();
            addLocker("S1", LockerSize.SMALL);
            addLocker("S2", LockerSize.SMALL);
            setCodeExpiry(2);
            String code1 = depositPackage("pkg1", LockerSize.SMALL, 1000);
            String code2 = depositPackage("pkg2", LockerSize.SMALL, 5000);
            // At time 8201: pkg1 expired (1000+7200<8201), pkg2 not (5000+7200>8201)
            List<String> expired = checkExpired(8201);
            // pkg2's code should still work
            // pkg1's code should NOT work (already expired)
            boolean pass = !code1.isEmpty()
                && !code2.isEmpty()
                && expired.size() == 1
                && expired[0] == "pkg1"
                && retrievePackage(code2) == true
                && retrievePackage(code1) == false;
            System.out.println((pass ? "PASS" : "FAIL") + ": testSelectiveExpiry");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testSelectiveExpiry (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testNoExpiryDefault() {
        try {
            initLockerSystem();
            addLocker("S1", LockerSize.SMALL);
            // Don't call setCodeExpiry — default is no expiry
            String code = depositPackage("pkg1", LockerSize.SMALL, 1000);
            List<String> expired = checkExpired(999999999);
            boolean pass = expired.isEmpty()
                && retrievePackage(code) == true;
            System.out.println((pass ? "PASS" : "FAIL") + ": testNoExpiryDefault");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testNoExpiryDefault (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    public static int runTests() {
        int passed = 0, total = 0;
        total++; if (testDepositNotification()) passed++;
        total++; if (testCodeExpiry()) passed++;
        total++; if (testNonExpiredCode()) passed++;
        total++; if (testExpiredLockerReuse()) passed++;
        total++; if (testExpiryNotification()) passed++;
        total++; if (testSelectiveExpiry()) passed++;
        total++; if (testNoExpiryDefault()) passed++;
        System.out.println("PART2_SUMMARY " + passed + "/" + total);
        return passed;
    }
}

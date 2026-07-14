// Ride Surge Pricing — Part 1 Tests
import java.util.*;
import java.util.stream.*;

class Part1Test {
    static boolean testNoSurgeNormalConditions() {
        try {
            PricingContext ctx = {10.0, 20, 5, "morning", "clear"};
            double surge = calculateSurge(ctx);
            boolean pass = surge >= 1.0); // always at least 1.0x
                && surge <= 1.5); // low demand shouldn't spike;
            System.out.println((pass ? "PASS" : "FAIL") + ": testNoSurgeNormalConditions");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testNoSurgeNormalConditions (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testHighDemandCausesSurge() {
        try {
            PricingContext ctx = {10.0, 2, 10, "evening", "clear"}; // 5:1 demand ratio
            double surge = calculateSurge(ctx);
            boolean pass = surge > 1.0); // must surge;
            System.out.println((pass ? "PASS" : "FAIL") + ": testHighDemandCausesSurge");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testHighDemandCausesSurge (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testStormIncreasesSurge() {
        try {
            PricingContext ctx = {10.0, 10, 10, "morning", "storm"};
            double surge = calculateSurge(ctx);
            boolean pass = surge > 1.0;
            System.out.println((pass ? "PASS" : "FAIL") + ": testStormIncreasesSurge");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testStormIncreasesSurge (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testFareCalculation() {
        try {
            PricingContext ctx = {100.0, 5, 5, "morning", "clear"};
            RideRequest req = Arrays.asList("u1", "A", "B", "economy");
            double fare = calculateFare(req, ctx);
            double surge = calculateSurge(ctx);
            boolean pass = approxEqual(fare, 100.0 * surge);
            System.out.println((pass ? "PASS" : "FAIL") + ": testFareCalculation");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testFareCalculation (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testSurgeMinimumOne() {
        try {
            PricingContext ctx = {50.0, 100, 1, "morning", "clear"}; // plenty of drivers
            double surge = calculateSurge(ctx);
            boolean pass = surge >= 1.0;
            System.out.println((pass ? "PASS" : "FAIL") + ": testSurgeMinimumOne");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testSurgeMinimumOne (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    public static int runTests() {
        int passed = 0, total = 0;
        total++; if (testNoSurgeNormalConditions()) passed++;
        total++; if (testHighDemandCausesSurge()) passed++;
        total++; if (testStormIncreasesSurge()) passed++;
        total++; if (testFareCalculation()) passed++;
        total++; if (testSurgeMinimumOne()) passed++;
        System.out.println("PART1_SUMMARY " + passed + "/" + total);
        return passed;
    }
}

// Ride Sharing — Part 2 Tests
import java.util.*;
import java.util.stream.*;

class Part2Test {
    static boolean testMostVacantStrategy() {
        try {
            RideService service = setupTestService();
            MostVacantStrategy strategy = new MostVacantStrategy();
            String rideId = service.selectRide("Priya", "Bangalore", "Mysore", 1,  strategy);
            Ride ride = service.getRide(rideId);
            // Deepa's ride has 5 seats (most vacant)
            boolean pass = !rideId.isEmpty()
                && ride.driverId == "Deepa";
            System.out.println((pass ? "PASS" : "FAIL") + ": testMostVacantStrategy");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testMostVacantStrategy (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testPreferredVehicleStrategy() {
        try {
            RideService service = setupTestService();
            PreferredVehicleStrategy strategy = new PreferredVehicleStrategy(service.getVehicles());
            String rideId = service.selectRide("Priya", "Bangalore", "Mysore", 1,  strategy, "Swift");
            Ride ride = service.getRide(rideId);
            // Rohan has a Swift on this route
            boolean pass = !rideId.isEmpty()
                && ride.driverId == "Rohan";
            System.out.println((pass ? "PASS" : "FAIL") + ": testPreferredVehicleStrategy");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testPreferredVehicleStrategy (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testNoMatchingRoute() {
        try {
            RideService service = setupTestService();
            MostVacantStrategy strategy = new MostVacantStrategy();
            String rideId = service.selectRide("Priya", "Delhi", "Mumbai", 1,  strategy);
            boolean pass = rideId.isEmpty();
            System.out.println((pass ? "PASS" : "FAIL") + ": testNoMatchingRoute");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testNoMatchingRoute (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testSeatsDecremented() {
        try {
            RideService service = setupTestService();
            MostVacantStrategy strategy = new MostVacantStrategy();
            String rideId = service.selectRide("Priya", "Bangalore", "Mysore", 2,  strategy);
            Ride ride = service.getRide(rideId);
            // Deepa's ride: 5 total, now 3 available
            boolean pass = !rideId.isEmpty()
                && ride.availableSeats == 3;
            System.out.println((pass ? "PASS" : "FAIL") + ": testSeatsDecremented");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testSeatsDecremented (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testRidesTakenIncremented() {
        try {
            RideService service = setupTestService();
            MostVacantStrategy strategy = new MostVacantStrategy();
            service.selectRide("Priya", "Bangalore", "Mysore", 1,  strategy);
            boolean pass = service.getUser("Priya").ridesTaken == 0
                && service.getUser("Priya").ridesTaken == 1;
            System.out.println((pass ? "PASS" : "FAIL") + ": testRidesTakenIncremented");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testRidesTakenIncremented (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testCannotSelectOwnRide() {
        try {
            RideService service = setupTestService();
            MostVacantStrategy strategy = new MostVacantStrategy();
            // Deepa tries to select a Bangalore→Mysore ride, but she offered one
            // Only Rohan's ride should be a candidate for Deepa
            String rideId = service.selectRide("Deepa", "Bangalore", "Mysore", 1,  strategy);
            if (!rideId.isEmpty()) {
            Ride ride = service.getRide(rideId);
            }
            boolean pass = ride.driverId != "Deepa");  // must not select own ride;
            System.out.println((pass ? "PASS" : "FAIL") + ": testCannotSelectOwnRide");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testCannotSelectOwnRide (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testNotEnoughSeats() {
        try {
            RideService service = setupTestService();
            MostVacantStrategy strategy = new MostVacantStrategy();
            // Request 10 seats — no ride has that many
            String rideId = service.selectRide("Priya", "Bangalore", "Mysore", 10,  strategy);
            boolean pass = rideId.isEmpty();
            System.out.println((pass ? "PASS" : "FAIL") + ": testNotEnoughSeats");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testNotEnoughSeats (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testNoMatchingVehicleModel() {
        try {
            RideService service = setupTestService();
            PreferredVehicleStrategy strategy = new PreferredVehicleStrategy(service.getVehicles());
            String rideId = service.selectRide("Priya", "Bangalore", "Mysore", 1,  strategy, "BMW");
            boolean pass = rideId.isEmpty();
            System.out.println((pass ? "PASS" : "FAIL") + ": testNoMatchingVehicleModel");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testNoMatchingVehicleModel (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    public static int runTests() {
        int passed = 0, total = 0;
        total++; if (testMostVacantStrategy()) passed++;
        total++; if (testPreferredVehicleStrategy()) passed++;
        total++; if (testNoMatchingRoute()) passed++;
        total++; if (testSeatsDecremented()) passed++;
        total++; if (testRidesTakenIncremented()) passed++;
        total++; if (testCannotSelectOwnRide()) passed++;
        total++; if (testNotEnoughSeats()) passed++;
        total++; if (testNoMatchingVehicleModel()) passed++;
        System.out.println("PART2_SUMMARY " + passed + "/" + total);
        return passed;
    }
}

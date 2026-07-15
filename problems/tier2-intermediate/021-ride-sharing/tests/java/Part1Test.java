// Ride Sharing — Part 1 Tests
import java.util.*;
import java.util.stream.*;

class Part1Test {
    static boolean testAddUser() {
        try {
            RideService service = new RideService();
            service.addUser("Rohan");
            boolean pass = service.hasUser("Rohan")
                && !service.hasUser("Unknown");
            System.out.println((pass ? "PASS" : "FAIL") + ": testAddUser");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testAddUser (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testAddVehicle() {
        try {
            RideService service = new RideService();
            service.addUser("Rohan");
            service.addVehicle("Rohan", "Swift", "KA-01-1234");
            boolean pass = service.hasVehicle("KA-01-1234");
            System.out.println((pass ? "PASS" : "FAIL") + ": testAddVehicle");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testAddVehicle (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testOfferRide() {
        try {
            RideService service = new RideService();
            service.addUser("Rohan");
            service.addVehicle("Rohan", "Swift", "KA-01-1234");
            String rideId = service.offerRide("Rohan", "Bangalore", "Mysore", 3, "KA-01-1234");
            Ride ride = service.getRide(rideId);
            boolean pass = !rideId.isEmpty()
                && service.hasRide(rideId)
                && ride.origin == "Bangalore"
                && ride.destination == "Mysore"
                && ride.totalSeats == 3
                && ride.availableSeats == 3
                && ride.active == true;
            System.out.println((pass ? "PASS" : "FAIL") + ": testOfferRide");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testOfferRide (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testNoDuplicateActiveRidePerVehicle() {
        try {
            RideService service = new RideService();
            service.addUser("Rohan");
            service.addVehicle("Rohan", "Swift", "KA-01-1234");
            String ride1 = service.offerRide("Rohan", "Bangalore", "Mysore", 3, "KA-01-1234");
            String ride2 = service.offerRide("Rohan", "Bangalore", "Chennai", 2, "KA-01-1234");
            boolean pass = !ride1.isEmpty()
                && ride2.isEmpty());  // should fail — vehicle already active;
            System.out.println((pass ? "PASS" : "FAIL") + ": testNoDuplicateActiveRidePerVehicle");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testNoDuplicateActiveRidePerVehicle (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testCannotUseOthersVehicle() {
        try {
            RideService service = new RideService();
            service.addUser("Rohan");
            service.addUser("Deepa");
            service.addVehicle("Rohan", "Swift", "KA-01-1234");
            String rideId = service.offerRide("Deepa", "Bangalore", "Mysore", 2, "KA-01-1234");
            boolean pass = rideId.isEmpty());  // Deepa doesn't own KA-01-1234;
            System.out.println((pass ? "PASS" : "FAIL") + ": testCannotUseOthersVehicle");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testCannotUseOthersVehicle (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testRidesOfferedCounter() {
        try {
            RideService service = new RideService();
            service.addUser("Rohan");
            service.addVehicle("Rohan", "Swift", "KA-01-1234");
            service.addVehicle("Rohan", "XUV", "KA-01-5678");
            service.offerRide("Rohan", "Bangalore", "Mysore", 3, "KA-01-1234");
            service.offerRide("Rohan", "Bangalore", "Chennai", 2, "KA-01-5678");
            boolean pass = service.getUser("Rohan").ridesOffered == 0
                && service.getUser("Rohan").ridesOffered == 1
                && service.getUser("Rohan").ridesOffered == 2;
            System.out.println((pass ? "PASS" : "FAIL") + ": testRidesOfferedCounter");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testRidesOfferedCounter (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testOfferRideInvalidUser() {
        try {
            RideService service = new RideService();
            service.addUser("Rohan");
            service.addVehicle("Rohan", "Swift", "KA-01-1234");
            String rideId = service.offerRide("Ghost", "A", "B", 2, "KA-01-1234");
            boolean pass = rideId.isEmpty();
            System.out.println((pass ? "PASS" : "FAIL") + ": testOfferRideInvalidUser");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testOfferRideInvalidUser (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testOfferRideInvalidVehicle() {
        try {
            RideService service = new RideService();
            service.addUser("Rohan");
            String rideId = service.offerRide("Rohan", "A", "B", 2, "INVALID-REG");
            boolean pass = rideId.isEmpty();
            System.out.println((pass ? "PASS" : "FAIL") + ": testOfferRideInvalidVehicle");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testOfferRideInvalidVehicle (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    public static int runTests() {
        int passed = 0, total = 0;
        total++; if (testAddUser()) passed++;
        total++; if (testAddVehicle()) passed++;
        total++; if (testOfferRide()) passed++;
        total++; if (testNoDuplicateActiveRidePerVehicle()) passed++;
        total++; if (testCannotUseOthersVehicle()) passed++;
        total++; if (testRidesOfferedCounter()) passed++;
        total++; if (testOfferRideInvalidUser()) passed++;
        total++; if (testOfferRideInvalidVehicle()) passed++;
        System.out.println("PART1_SUMMARY " + passed + "/" + total);
        return passed;
    }
}

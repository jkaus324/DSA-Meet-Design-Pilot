// Ride Sharing — Part 3 Tests
import java.util.*;
import java.util.stream.*;

class Part3Test {
    static boolean testEndRideMarksInactive() {
        try {
            RideService service = new RideService();
            service.addUser("Rohan");
            service.addVehicle("Rohan", "Swift", "KA-01-1234");
            String rideId = service.offerRide("Rohan", "Bangalore", "Mysore", 3, "KA-01-1234");
            service.endRide(rideId);
            boolean pass = service.getRide(rideId).active == true
                && service.getRide(rideId).active == false;
            System.out.println((pass ? "PASS" : "FAIL") + ": testEndRideMarksInactive");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testEndRideMarksInactive (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testVehicleFreedAfterEnd() {
        try {
            RideService service = new RideService();
            service.addUser("Rohan");
            service.addVehicle("Rohan", "Swift", "KA-01-1234");
            String ride1 = service.offerRide("Rohan", "Bangalore", "Mysore", 3, "KA-01-1234");
            // Cannot offer again while active
            String ride2 = service.offerRide("Rohan", "Bangalore", "Chennai", 2, "KA-01-1234");
            // End first ride
            service.endRide(ride1);
            // Now can offer again
            String ride3 = service.offerRide("Rohan", "Bangalore", "Chennai", 2, "KA-01-1234");
            boolean pass = !ride1.isEmpty()
                && ride2.isEmpty()
                && !ride3.isEmpty();
            System.out.println((pass ? "PASS" : "FAIL") + ": testVehicleFreedAfterEnd");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testVehicleFreedAfterEnd (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testEndRideIdempotent() {
        try {
            RideService service = new RideService();
            service.addUser("Rohan");
            service.addVehicle("Rohan", "Swift", "KA-01-1234");
            String rideId = service.offerRide("Rohan", "Bangalore", "Mysore", 3, "KA-01-1234");
            service.endRide(rideId);
            service.endRide(rideId);  // should not crash
            boolean pass = service.getRide(rideId).active == false;
            System.out.println((pass ? "PASS" : "FAIL") + ": testEndRideIdempotent");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testEndRideIdempotent (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testEndNonexistentRide() {
        try {
            RideService service = new RideService();
            service.endRide("RIDE-999");  // should not crash
            boolean pass = true; // FIXME: no assertions found
            System.out.println((pass ? "PASS" : "FAIL") + ": testEndNonexistentRide");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testEndNonexistentRide (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testRideStats() {
        try {
            RideService service = new RideService();
            service.addUser("Rohan");
            service.addUser("Deepa");
            service.addUser("Priya");
            service.addVehicle("Rohan", "Swift", "KA-01-1234");
            service.addVehicle("Deepa", "XUV", "KA-02-5678");
            service.offerRide("Rohan", "Bangalore", "Mysore", 3, "KA-01-1234");
            service.offerRide("Deepa", "Bangalore", "Mysore", 5, "KA-02-5678");
            MostVacantStrategy strategy = new MostVacantStrategy();
            service.selectRide("Priya", "Bangalore", "Mysore", 1,  strategy);
            var stats = service.getRideStats();
            // Find each user's stats
            boolean foundRohan = false, foundDeepa = false, foundPriya = false;
            for (var _e_stats_ : stats.entrySet()) {
            var name = _e_stats_.getKey(); var counts = _e_stats_.getValue();
            if (name == "Rohan") {
            foundRohan = true;
            }
            if (name == "Deepa") {
            foundDeepa = true;
            }
            if (name == "Priya") {
            foundPriya = true;
            }
            }
            boolean pass = !stats.isEmpty()
                && counts.first == 1);   // offered 1
                && counts.second == 0);  // taken 0
                && counts.first == 1);   // offered 1
                && counts.second == 0);  // taken 0
                && counts.first == 0);   // offered 0
                && counts.second == 1);  // taken 1
                && foundRohan & foundDeepa & foundPriya;
            System.out.println((pass ? "PASS" : "FAIL") + ": testRideStats");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testRideStats (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testEndedRideNotSelectable() {
        try {
            RideService service = new RideService();
            service.addUser("Rohan");
            service.addUser("Priya");
            service.addVehicle("Rohan", "Swift", "KA-01-1234");
            String rideId = service.offerRide("Rohan", "Bangalore", "Mysore", 3, "KA-01-1234");
            service.endRide(rideId);
            MostVacantStrategy strategy = new MostVacantStrategy();
            String selected = service.selectRide("Priya", "Bangalore", "Mysore", 1,  strategy);
            boolean pass = selected.isEmpty());  // ended ride should not be selectable;
            System.out.println((pass ? "PASS" : "FAIL") + ": testEndedRideNotSelectable");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testEndedRideNotSelectable (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testFullWorkflow() {
        try {
            RideService service = new RideService();
            service.addUser("Rohan");
            service.addUser("Deepa");
            service.addVehicle("Rohan", "Swift", "KA-01-1234");
            // Rohan offers ride
            String ride1 = service.offerRide("Rohan", "Bangalore", "Mysore", 2, "KA-01-1234");
            // Deepa takes ride
            MostVacantStrategy strategy = new MostVacantStrategy();
            String selected = service.selectRide("Deepa", "Bangalore", "Mysore", 1,  strategy);
            // Rohan ends ride
            service.endRide(ride1);
            // Rohan offers new ride with same vehicle
            String ride2 = service.offerRide("Rohan", "Mysore", "Bangalore", 2, "KA-01-1234");
            // Check stats
            boolean pass = !ride1.isEmpty()
                && selected == ride1
                && service.getRide(ride1).availableSeats == 1
                && service.getRide(ride1).active == false
                && !ride2.isEmpty()
                && ride2 != ride1
                && service.getUser("Rohan").ridesOffered == 2
                && service.getUser("Deepa").ridesTaken == 1;
            System.out.println((pass ? "PASS" : "FAIL") + ": testFullWorkflow");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testFullWorkflow (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    public static int runTests() {
        int passed = 0, total = 0;
        total++; if (testEndRideMarksInactive()) passed++;
        total++; if (testVehicleFreedAfterEnd()) passed++;
        total++; if (testEndRideIdempotent()) passed++;
        total++; if (testEndNonexistentRide()) passed++;
        total++; if (testRideStats()) passed++;
        total++; if (testEndedRideNotSelectable()) passed++;
        total++; if (testFullWorkflow()) passed++;
        System.out.println("PART3_SUMMARY " + passed + "/" + total);
        return passed;
    }
}

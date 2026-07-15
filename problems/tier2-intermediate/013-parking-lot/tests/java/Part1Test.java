// Parking Lot — Part 1 Tests
import java.util.*;
import java.util.stream.*;

class Part1Test {
    static boolean testParkCar() {
        try {
            ParkingLot lot = new ParkingLot(2);
            lot.addSpot(0, SpotSize.MEDIUM);
            Vehicle car{"ABC123", VehicleType.CAR};
            Ticket t = lot.parkVehicle(car, 1000);
            boolean pass = t != null
                && t.licensePlate == "ABC123"
                && t.floor == 0;
            System.out.println((pass ? "PASS" : "FAIL") + ": testParkCar");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testParkCar (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testParkMotorcycle() {
        try {
            ParkingLot lot = new ParkingLot(1);
            lot.addSpot(0, SpotSize.SMALL);
            Vehicle moto{"MOTO1", VehicleType.MOTORCYCLE};
            Ticket t = lot.parkVehicle(moto, 1000);
            boolean pass = t != null
                && t.licensePlate == "MOTO1";
            System.out.println((pass ? "PASS" : "FAIL") + ": testParkMotorcycle");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testParkMotorcycle (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testCarNoSmallSpot() {
        try {
            ParkingLot lot = new ParkingLot(1);
            lot.addSpot(0, SpotSize.SMALL);
            Vehicle car{"CAR1", VehicleType.CAR};
            Ticket t = lot.parkVehicle(car, 1000);
            boolean pass = t == null;
            System.out.println((pass ? "PASS" : "FAIL") + ": testCarNoSmallSpot");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testCarNoSmallSpot (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testMotorcycleInMediumSpot() {
        try {
            ParkingLot lot = new ParkingLot(1);
            lot.addSpot(0, SpotSize.MEDIUM);
            Vehicle moto{"MOTO2", VehicleType.MOTORCYCLE};
            Ticket t = lot.parkVehicle(moto, 1000);
            boolean pass = t != null;
            System.out.println((pass ? "PASS" : "FAIL") + ": testMotorcycleInMediumSpot");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testMotorcycleInMediumSpot (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testTruckLargeSpot() {
        try {
            ParkingLot lot = new ParkingLot(1);
            lot.addSpot(0, SpotSize.MEDIUM);
            lot.addSpot(0, SpotSize.LARGE);
            Vehicle truck{"TRUCK1", VehicleType.TRUCK};
            Ticket t = lot.parkVehicle(truck, 1000);
            boolean pass = t != null
                && t.spotId.find("S1") != String.npos); // second spot (large;
            System.out.println((pass ? "PASS" : "FAIL") + ": testTruckLargeSpot");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testTruckLargeSpot (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testUnparkFreesSpot() {
        try {
            ParkingLot lot = new ParkingLot(1);
            lot.addSpot(0, SpotSize.MEDIUM);
            Vehicle car{"CAR2", VehicleType.CAR};
            Ticket t = lot.parkVehicle(car, 1000);
            String tid = t.ticketId;
            double fee = lot.unparkVehicle(tid, 4600); // 3600 seconds = 1 hour
            boolean pass = t != null
                && lot.getAvailableSpots(SpotSize.MEDIUM) == 0
                && fee >= 0
                && lot.getAvailableSpots(SpotSize.MEDIUM) == 1;
            System.out.println((pass ? "PASS" : "FAIL") + ": testUnparkFreesSpot");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testUnparkFreesSpot (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testInvalidTicket() {
        try {
            ParkingLot lot = new ParkingLot(1);
            double fee = lot.unparkVehicle("INVALID", 5000);
            boolean pass = fee < 0;
            System.out.println((pass ? "PASS" : "FAIL") + ": testInvalidTicket");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testInvalidTicket (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testAvailableSpotsCount() {
        try {
            ParkingLot lot = new ParkingLot(2);
            lot.addSpot(0, SpotSize.SMALL);
            lot.addSpot(0, SpotSize.MEDIUM);
            lot.addSpot(0, SpotSize.MEDIUM);
            lot.addSpot(1, SpotSize.LARGE);
            boolean pass = lot.getAvailableSpots(SpotSize.SMALL) == 1
                && lot.getAvailableSpots(SpotSize.MEDIUM) == 2
                && lot.getAvailableSpots(SpotSize.LARGE) == 1;
            System.out.println((pass ? "PASS" : "FAIL") + ": testAvailableSpotsCount");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testAvailableSpotsCount (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testAvailableSpotsByFloor() {
        try {
            ParkingLot lot = new ParkingLot(2);
            lot.addSpot(0, SpotSize.MEDIUM);
            lot.addSpot(0, SpotSize.MEDIUM);
            lot.addSpot(1, SpotSize.MEDIUM);
            boolean pass = lot.getAvailableSpotsByFloor(0, SpotSize.MEDIUM) == 2
                && lot.getAvailableSpotsByFloor(1, SpotSize.MEDIUM) == 1;
            System.out.println((pass ? "PASS" : "FAIL") + ": testAvailableSpotsByFloor");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testAvailableSpotsByFloor (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testNearestFloorAllocation() {
        try {
            ParkingLot lot = new ParkingLot(3);
            lot.addSpot(0, SpotSize.SMALL);  // floor 0 has only small
            lot.addSpot(1, SpotSize.MEDIUM); // floor 1 has medium
            lot.addSpot(2, SpotSize.MEDIUM); // floor 2 has medium
            Vehicle car{"CAR3", VehicleType.CAR};
            Ticket t = lot.parkVehicle(car, 1000);
            boolean pass = t != null
                && t.floor == 1); // floor 0 has no compatible spot, floor 1 is nearest;
            System.out.println((pass ? "PASS" : "FAIL") + ": testNearestFloorAllocation");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testNearestFloorAllocation (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testFullLot() {
        try {
            ParkingLot lot = new ParkingLot(1);
            lot.addSpot(0, SpotSize.MEDIUM);
            Vehicle car1{"C1", VehicleType.CAR};
            Vehicle car2{"C2", VehicleType.CAR};
            boolean pass = lot.parkVehicle(car1, 1000) != null
                && lot.parkVehicle(car2, 1000) == null;
            System.out.println((pass ? "PASS" : "FAIL") + ": testFullLot");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testFullLot (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    public static int runTests() {
        int passed = 0, total = 0;
        total++; if (testParkCar()) passed++;
        total++; if (testParkMotorcycle()) passed++;
        total++; if (testCarNoSmallSpot()) passed++;
        total++; if (testMotorcycleInMediumSpot()) passed++;
        total++; if (testTruckLargeSpot()) passed++;
        total++; if (testUnparkFreesSpot()) passed++;
        total++; if (testInvalidTicket()) passed++;
        total++; if (testAvailableSpotsCount()) passed++;
        total++; if (testAvailableSpotsByFloor()) passed++;
        total++; if (testNearestFloorAllocation()) passed++;
        total++; if (testFullLot()) passed++;
        System.out.println("PART1_SUMMARY " + passed + "/" + total);
        return passed;
    }
}

// Parking Lot — Part 2 Tests
import java.util.*;
import java.util.stream.*;

class Part2Test {
    static boolean testFlatRate() {
        try {
            ParkingLot lot = new ParkingLot(1);
            lot.addSpot(0, SpotSize.MEDIUM);
            FlatRate flat = new FlatRate(10.0);
            lot.setPricingStrategy( flat);
            Vehicle car{"CAR1", VehicleType.CAR};
            Ticket t = lot.parkVehicle(car, 1000, "G1");
            String tid = t.ticketId;
            double fee = lot.unparkVehicle(tid, 8600, "G2"); // 7600s = ~2.1 hours
            boolean pass = Math.abs(fee - 10.0) < 0.01;
            System.out.println((pass ? "PASS" : "FAIL") + ": testFlatRate");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testFlatRate (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testHourlyRate() {
        try {
            ParkingLot lot = new ParkingLot(1);
            lot.addSpot(0, SpotSize.MEDIUM);
            Hourly hourly = new Hourly(5.0);
            lot.setPricingStrategy( hourly);
            Vehicle car{"CAR2", VehicleType.CAR};
            Ticket t = lot.parkVehicle(car, 0, "G1");
            String tid = t.ticketId;
            // 9000 seconds = 2.5 hours, ceil = 3 hours => $15
            double fee = lot.unparkVehicle(tid, 9000, "G2");
            boolean pass = Math.abs(fee - 15.0) < 0.01;
            System.out.println((pass ? "PASS" : "FAIL") + ": testHourlyRate");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testHourlyRate (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testHourlyExactHour() {
        try {
            ParkingLot lot = new ParkingLot(1);
            lot.addSpot(0, SpotSize.MEDIUM);
            Hourly hourly = new Hourly(5.0);
            lot.setPricingStrategy( hourly);
            Vehicle car{"CAR3", VehicleType.CAR};
            Ticket t = lot.parkVehicle(car, 0, "G1");
            String tid = t.ticketId;
            double fee = lot.unparkVehicle(tid, 3600, "G2"); // exactly 1 hour => $5
            boolean pass = Math.abs(fee - 5.0) < 0.01;
            System.out.println((pass ? "PASS" : "FAIL") + ": testHourlyExactHour");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testHourlyExactHour (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testTieredBaseRate() {
        try {
            ParkingLot lot = new ParkingLot(1);
            lot.addSpot(0, SpotSize.MEDIUM);
            Tiered tiered = new Tiered(10.0, 8.0, 5.0);
            lot.setPricingStrategy( tiered);
            Vehicle car{"CAR4", VehicleType.CAR};
            Ticket t = lot.parkVehicle(car, 0, "G1");
            String tid = t.ticketId;
            double fee = lot.unparkVehicle(tid, 3000, "G2"); // 3000s < 1h, ceil=1h => $10
            boolean pass = Math.abs(fee - 10.0) < 0.01;
            System.out.println((pass ? "PASS" : "FAIL") + ": testTieredBaseRate");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testTieredBaseRate (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testTieredMidRate() {
        try {
            ParkingLot lot = new ParkingLot(1);
            lot.addSpot(0, SpotSize.MEDIUM);
            Tiered tiered = new Tiered(10.0, 8.0, 5.0);
            lot.setPricingStrategy( tiered);
            Vehicle car{"CAR5", VehicleType.CAR};
            Ticket t = lot.parkVehicle(car, 0, "G1");
            String tid = t.ticketId;
            double fee = lot.unparkVehicle(tid, 7200, "G2"); // 7200s = 2h => $10 + $8*1 = $18
            boolean pass = Math.abs(fee - 18.0) < 0.01;
            System.out.println((pass ? "PASS" : "FAIL") + ": testTieredMidRate");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testTieredMidRate (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testTieredHighRate() {
        try {
            ParkingLot lot = new ParkingLot(1);
            lot.addSpot(0, SpotSize.MEDIUM);
            Tiered tiered = new Tiered(10.0, 8.0, 5.0);
            lot.setPricingStrategy( tiered);
            Vehicle car{"CAR6", VehicleType.CAR};
            Ticket t = lot.parkVehicle(car, 0, "G1");
            String tid = t.ticketId;
            double fee = lot.unparkVehicle(tid, 18000, "G2"); // 18000s = 5h => $10 + $16 + $10 = $36
            boolean pass = Math.abs(fee - 36.0) < 0.01;
            System.out.println((pass ? "PASS" : "FAIL") + ": testTieredHighRate");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testTieredHighRate (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testSwapStrategy() {
        try {
            ParkingLot lot = new ParkingLot(1);
            lot.addSpot(0, SpotSize.MEDIUM);
            lot.addSpot(0, SpotSize.MEDIUM);
            FlatRate flat = new FlatRate(10.0);
            Hourly hourly = new Hourly(5.0);
            lot.setPricingStrategy( flat);
            Vehicle car1{"SW1", VehicleType.CAR};
            Ticket t1 = lot.parkVehicle(car1, 0, "G1");
            String tid1 = t1.ticketId;
            double fee1 = lot.unparkVehicle(tid1, 7200, "G2"); // flat = $10
            lot.setPricingStrategy( hourly);
            Vehicle car2{"SW2", VehicleType.CAR};
            Ticket t2 = lot.parkVehicle(car2, 0, "G1");
            String tid2 = t2.ticketId;
            double fee2 = lot.unparkVehicle(tid2, 7200, "G2"); // hourly: 2h * $5 = $10
            boolean pass = Math.abs(fee1 - 10.0) < 0.01
                && Math.abs(fee2 - 10.0) < 0.01;
            System.out.println((pass ? "PASS" : "FAIL") + ": testSwapStrategy");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testSwapStrategy (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testGateManagement() {
        try {
            ParkingLot lot = new ParkingLot(1);
            lot.addGate("E1", GateType.ENTRY);
            lot.addGate("E2", GateType.ENTRY);
            lot.addGate("X1", GateType.EXIT);
            var entryGates = lot.getGates(GateType.ENTRY);
            var exitGates = lot.getGates(GateType.EXIT);
            boolean pass = entryGates.size() == 2
                && exitGates.size() == 1
                && exitGates[0] == "X1";
            System.out.println((pass ? "PASS" : "FAIL") + ": testGateManagement");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testGateManagement (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testGateOnTicket() {
        try {
            ParkingLot lot = new ParkingLot(1);
            lot.addSpot(0, SpotSize.MEDIUM);
            FlatRate flat = new FlatRate(10.0);
            lot.setPricingStrategy( flat);
            lot.addGate("ENTRY1", GateType.ENTRY);
            lot.addGate("EXIT1", GateType.EXIT);
            Vehicle car{"GATE_CAR", VehicleType.CAR};
            Ticket t = lot.parkVehicle(car, 0, "ENTRY1");
            boolean pass = t != null
                && t.entryGateId == "ENTRY1";
            System.out.println((pass ? "PASS" : "FAIL") + ": testGateOnTicket");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testGateOnTicket (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testShortStayRoundsUp() {
        try {
            ParkingLot lot = new ParkingLot(1);
            lot.addSpot(0, SpotSize.SMALL);
            Hourly hourly = new Hourly(5.0);
            lot.setPricingStrategy( hourly);
            Vehicle moto{"SHORT", VehicleType.MOTORCYCLE};
            Ticket t = lot.parkVehicle(moto, 0, "G1");
            String tid = t.ticketId;
            double fee = lot.unparkVehicle(tid, 1, "G2"); // 1 second => ceil = 1 hour => $5
            boolean pass = Math.abs(fee - 5.0) < 0.01;
            System.out.println((pass ? "PASS" : "FAIL") + ": testShortStayRoundsUp");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testShortStayRoundsUp (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    public static int runTests() {
        int passed = 0, total = 0;
        total++; if (testFlatRate()) passed++;
        total++; if (testHourlyRate()) passed++;
        total++; if (testHourlyExactHour()) passed++;
        total++; if (testTieredBaseRate()) passed++;
        total++; if (testTieredMidRate()) passed++;
        total++; if (testTieredHighRate()) passed++;
        total++; if (testSwapStrategy()) passed++;
        total++; if (testGateManagement()) passed++;
        total++; if (testGateOnTicket()) passed++;
        total++; if (testShortStayRoundsUp()) passed++;
        System.out.println("PART2_SUMMARY " + passed + "/" + total);
        return passed;
    }
}

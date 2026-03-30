// Part 2 Tests -- Pricing Strategies and Gate Management
// Tests FlatRate, Hourly, Tiered pricing, strategy swapping, and gates

#include "solution.cpp"
#include <cassert>
#include <iostream>
#include <cmath>
using namespace std;

int part2_tests() {
    int passed = 0;
    int failed = 0;

    // Test 1: FlatRate pricing — same fee regardless of duration
    try {
        ParkingLot lot(1);
        lot.addSpot(0, SpotSize::MEDIUM);
        FlatRate flat(10.0);
        lot.setPricingStrategy(&flat);
        Vehicle car{"CAR1", VehicleType::CAR};
        Ticket* t = lot.parkVehicle(car, 1000, "G1");
        string tid = t->ticketId;
        double fee = lot.unparkVehicle(tid, 8600, "G2"); // 7600s = ~2.1 hours
        assert(abs(fee - 10.0) < 0.01);
        cout << "PASS test_flat_rate" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_flat_rate" << endl;
        failed++;
    }

    // Test 2: Hourly pricing — rounds up to full hours
    try {
        ParkingLot lot(1);
        lot.addSpot(0, SpotSize::MEDIUM);
        Hourly hourly(5.0); // $5/hour
        lot.setPricingStrategy(&hourly);
        Vehicle car{"CAR2", VehicleType::CAR};
        Ticket* t = lot.parkVehicle(car, 0, "G1");
        string tid = t->ticketId;
        // 9000 seconds = 2.5 hours, ceil = 3 hours => $15
        double fee = lot.unparkVehicle(tid, 9000, "G2");
        assert(abs(fee - 15.0) < 0.01);
        cout << "PASS test_hourly_rate" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_hourly_rate" << endl;
        failed++;
    }

    // Test 3: Hourly pricing — exactly 1 hour
    try {
        ParkingLot lot(1);
        lot.addSpot(0, SpotSize::MEDIUM);
        Hourly hourly(5.0);
        lot.setPricingStrategy(&hourly);
        Vehicle car{"CAR3", VehicleType::CAR};
        Ticket* t = lot.parkVehicle(car, 0, "G1");
        string tid = t->ticketId;
        double fee = lot.unparkVehicle(tid, 3600, "G2"); // exactly 1 hour => $5
        assert(abs(fee - 5.0) < 0.01);
        cout << "PASS test_hourly_exact_hour" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_hourly_exact_hour" << endl;
        failed++;
    }

    // Test 4: Tiered pricing — under 1 hour (base rate only)
    try {
        ParkingLot lot(1);
        lot.addSpot(0, SpotSize::MEDIUM);
        Tiered tiered(10.0, 8.0, 5.0); // base=$10, mid=$8/hr, high=$5/hr
        lot.setPricingStrategy(&tiered);
        Vehicle car{"CAR4", VehicleType::CAR};
        Ticket* t = lot.parkVehicle(car, 0, "G1");
        string tid = t->ticketId;
        double fee = lot.unparkVehicle(tid, 3000, "G2"); // 3000s < 1h, ceil=1h => $10
        assert(abs(fee - 10.0) < 0.01);
        cout << "PASS test_tiered_base_rate" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_tiered_base_rate" << endl;
        failed++;
    }

    // Test 5: Tiered pricing — 2 hours (base + 1 mid)
    try {
        ParkingLot lot(1);
        lot.addSpot(0, SpotSize::MEDIUM);
        Tiered tiered(10.0, 8.0, 5.0);
        lot.setPricingStrategy(&tiered);
        Vehicle car{"CAR5", VehicleType::CAR};
        Ticket* t = lot.parkVehicle(car, 0, "G1");
        string tid = t->ticketId;
        double fee = lot.unparkVehicle(tid, 7200, "G2"); // 7200s = 2h => $10 + $8*1 = $18
        assert(abs(fee - 18.0) < 0.01);
        cout << "PASS test_tiered_mid_rate" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_tiered_mid_rate" << endl;
        failed++;
    }

    // Test 6: Tiered pricing — 5 hours (base + 2*mid + 2*high)
    try {
        ParkingLot lot(1);
        lot.addSpot(0, SpotSize::MEDIUM);
        Tiered tiered(10.0, 8.0, 5.0);
        lot.setPricingStrategy(&tiered);
        Vehicle car{"CAR6", VehicleType::CAR};
        Ticket* t = lot.parkVehicle(car, 0, "G1");
        string tid = t->ticketId;
        double fee = lot.unparkVehicle(tid, 18000, "G2"); // 18000s = 5h => $10 + $16 + $10 = $36
        assert(abs(fee - 36.0) < 0.01);
        cout << "PASS test_tiered_high_rate" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_tiered_high_rate" << endl;
        failed++;
    }

    // Test 7: Swap strategy at runtime
    try {
        ParkingLot lot(1);
        lot.addSpot(0, SpotSize::MEDIUM);
        lot.addSpot(0, SpotSize::MEDIUM);
        FlatRate flat(10.0);
        Hourly hourly(5.0);
        lot.setPricingStrategy(&flat);
        Vehicle car1{"SW1", VehicleType::CAR};
        Ticket* t1 = lot.parkVehicle(car1, 0, "G1");
        string tid1 = t1->ticketId;
        double fee1 = lot.unparkVehicle(tid1, 7200, "G2"); // flat = $10
        assert(abs(fee1 - 10.0) < 0.01);
        lot.setPricingStrategy(&hourly);
        Vehicle car2{"SW2", VehicleType::CAR};
        Ticket* t2 = lot.parkVehicle(car2, 0, "G1");
        string tid2 = t2->ticketId;
        double fee2 = lot.unparkVehicle(tid2, 7200, "G2"); // hourly: 2h * $5 = $10
        assert(abs(fee2 - 10.0) < 0.01);
        cout << "PASS test_swap_strategy" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_swap_strategy" << endl;
        failed++;
    }

    // Test 8: Add and retrieve gates
    try {
        ParkingLot lot(1);
        lot.addGate("E1", GateType::ENTRY);
        lot.addGate("E2", GateType::ENTRY);
        lot.addGate("X1", GateType::EXIT);
        auto entryGates = lot.getGates(GateType::ENTRY);
        auto exitGates = lot.getGates(GateType::EXIT);
        assert(entryGates.size() == 2);
        assert(exitGates.size() == 1);
        assert(exitGates[0] == "X1");
        cout << "PASS test_gate_management" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_gate_management" << endl;
        failed++;
    }

    // Test 9: Gate IDs recorded on ticket
    try {
        ParkingLot lot(1);
        lot.addSpot(0, SpotSize::MEDIUM);
        FlatRate flat(10.0);
        lot.setPricingStrategy(&flat);
        lot.addGate("ENTRY1", GateType::ENTRY);
        lot.addGate("EXIT1", GateType::EXIT);
        Vehicle car{"GATE_CAR", VehicleType::CAR};
        Ticket* t = lot.parkVehicle(car, 0, "ENTRY1");
        assert(t != nullptr);
        assert(t->entryGateId == "ENTRY1");
        cout << "PASS test_gate_on_ticket" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_gate_on_ticket" << endl;
        failed++;
    }

    // Test 10: Hourly pricing — very short stay (1 second rounds up to 1 hour)
    try {
        ParkingLot lot(1);
        lot.addSpot(0, SpotSize::SMALL);
        Hourly hourly(5.0);
        lot.setPricingStrategy(&hourly);
        Vehicle moto{"SHORT", VehicleType::MOTORCYCLE};
        Ticket* t = lot.parkVehicle(moto, 0, "G1");
        string tid = t->ticketId;
        double fee = lot.unparkVehicle(tid, 1, "G2"); // 1 second => ceil = 1 hour => $5
        assert(abs(fee - 5.0) < 0.01);
        cout << "PASS test_short_stay_rounds_up" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_short_stay_rounds_up" << endl;
        failed++;
    }

    cout << "PART2_SUMMARY " << passed << "/" << (passed + failed) << endl;
    return failed;
}

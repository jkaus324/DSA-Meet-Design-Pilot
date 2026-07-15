// Parking Lot — Solution (Java)
import java.util.*;

class ParkOp {
    public String kind;
    public String s1;
    public String s2;
    public String s3;
    public int i1;
    public int i2;
    public int i3;

    public ParkOp(String kind, String s1, String s2, String s3, int i1, int i2, int i3) {
        this.kind = kind;
        this.s1 = s1;
        this.s2 = s2;
        this.s3 = s3;
        this.i1 = i1;
        this.i2 = i2;
        this.i3 = i3;
    }
}

enum VehicleType { MOTORCYCLE, CAR, TRUCK }
enum SpotSize { SMALL, MEDIUM, LARGE }
enum GateType { ENTRY, EXIT }

class Vehicle {
    public String licensePlate;
    public VehicleType type;

    public Vehicle(String licensePlate, VehicleType type) {
        this.licensePlate = licensePlate;
        this.type = type;
    }
}

class ParkingSpot {
    public String spotId;
    public int floor;
    public SpotSize size;
    public boolean isOccupied;
    public String vehicleLicensePlate;

    public ParkingSpot(String spotId, int floor, SpotSize size) {
        this.spotId = spotId;
        this.floor = floor;
        this.size = size;
        this.isOccupied = false;
        this.vehicleLicensePlate = "";
    }
}

class Ticket {
    public String ticketId;
    public String licensePlate;
    public String spotId;
    public int floor;
    public long entryTime;
    public String entryGateId;
    public String exitGateId;

    public Ticket() {
        this.ticketId = "";
        this.licensePlate = "";
        this.spotId = "";
        this.floor = 0;
        this.entryTime = 0L;
        this.entryGateId = "";
        this.exitGateId = "";
    }
}

class Gate {
    public String gateId;
    public GateType type;

    public Gate(String gateId, GateType type) {
        this.gateId = gateId;
        this.type = type;
    }
}

interface PricingStrategy {
    double calculateFee(long durationSeconds);
}

class FlatRate implements PricingStrategy {
    private double fee;
    public FlatRate(double fee) { this.fee = fee; }
    @Override public double calculateFee(long durationSeconds) { return fee; }
}

class Hourly implements PricingStrategy {
    private double ratePerHour;
    public Hourly(double rate) { this.ratePerHour = rate; }
    @Override public double calculateFee(long durationSeconds) {
        double hours = Math.ceil((double) durationSeconds / 3600.0);
        return ratePerHour * hours;
    }
}

class Tiered implements PricingStrategy {
    private double baseRate, midRate, highRate;
    public Tiered(double base, double mid, double high) {
        this.baseRate = base; this.midRate = mid; this.highRate = high;
    }
    @Override public double calculateFee(long durationSeconds) {
        double hours = Math.ceil((double) durationSeconds / 3600.0);
        if (hours <= 1) return baseRate;
        if (hours <= 3) return baseRate + midRate * (hours - 1);
        return baseRate + midRate * 2 + highRate * (hours - 3);
    }
}

class ParkingLot {
    List<List<ParkingSpot>> floors = new ArrayList<>();
    Map<String, Ticket> activeTickets = new LinkedHashMap<>();
    List<Gate> gates = new ArrayList<>();
    PricingStrategy strategy = null;
    int nextTicketId = 1;

    public ParkingLot(int numFloors) {
        for (int i = 0; i < numFloors; i++) floors.add(new ArrayList<>());
    }

    static SpotSize getMinSpotSize(VehicleType type) {
        switch (type) {
            case MOTORCYCLE: return SpotSize.SMALL;
            case CAR:        return SpotSize.MEDIUM;
            case TRUCK:      return SpotSize.LARGE;
        }
        return SpotSize.LARGE;
    }

    static boolean isCompatible(SpotSize spotSize, SpotSize minRequired) {
        return spotSize.ordinal() >= minRequired.ordinal();
    }

    public void addSpot(int floor, SpotSize size) {
        if (floor < 0 || floor >= floors.size()) return;
        String spotId = "F" + floor + "S" + floors.get(floor).size();
        floors.get(floor).add(new ParkingSpot(spotId, floor, size));
    }

    public void setPricingStrategy(PricingStrategy s) { this.strategy = s; }

    public void addGate(String gateId, GateType type) {
        gates.add(new Gate(gateId, type));
    }

    public List<String> getGates(GateType type) {
        List<String> result = new ArrayList<>();
        for (Gate g : gates) if (g.type == type) result.add(g.gateId);
        return result;
    }

    public Ticket parkVehicle(Vehicle vehicle, long entryTime, String gateId) {
        SpotSize minSize = getMinSpotSize(vehicle.type);
        for (int f = 0; f < floors.size(); f++) {
            List<ParkingSpot> floorSpots = floors.get(f);
            for (int s = 0; s < floorSpots.size(); s++) {
                ParkingSpot spot = floorSpots.get(s);
                if (!spot.isOccupied && isCompatible(spot.size, minSize)) {
                    spot.isOccupied = true;
                    spot.vehicleLicensePlate = vehicle.licensePlate;
                    String tid = "T" + (nextTicketId++);
                    Ticket ticket = new Ticket();
                    ticket.ticketId = tid;
                    ticket.licensePlate = vehicle.licensePlate;
                    ticket.spotId = spot.spotId;
                    ticket.floor = f;
                    ticket.entryTime = entryTime;
                    ticket.entryGateId = gateId;
                    ticket.exitGateId = "";
                    activeTickets.put(tid, ticket);
                    return ticket;
                }
            }
        }
        return null;
    }

    public double unparkVehicle(String ticketId, long exitTime, String gateId) {
        Ticket ticket = activeTickets.get(ticketId);
        if (ticket == null) return -1.0;

        ticket.exitGateId = gateId;

        for (List<ParkingSpot> floorSpots : floors) {
            for (ParkingSpot spot : floorSpots) {
                if (spot.spotId.equals(ticket.spotId) && spot.isOccupied) {
                    spot.isOccupied = false;
                    spot.vehicleLicensePlate = "";
                    break;
                }
            }
        }

        long duration = exitTime - ticket.entryTime;
        double fee;
        if (strategy != null) {
            fee = strategy.calculateFee(duration);
        } else {
            fee = (double) duration;
        }

        activeTickets.remove(ticketId);
        return fee;
    }

    public int getAvailableSpots(SpotSize size) {
        int count = 0;
        for (List<ParkingSpot> floorSpots : floors)
            for (ParkingSpot spot : floorSpots)
                if (!spot.isOccupied && spot.size == size) count++;
        return count;
    }

    public int getAvailableSpotsByFloor(int floor, SpotSize size) {
        if (floor < 0 || floor >= floors.size()) return 0;
        int count = 0;
        for (ParkingSpot spot : floors.get(floor))
            if (!spot.isOccupied && spot.size == size) count++;
        return count;
    }
}

class TicketSnap {
    public String id = "";
    public int floor = -1;
    public String spotId = "";
    public String entryGate = "";
}

public class Solution {
    private static SpotSize sizeFrom(String s) {
        if ("S".equals(s) || "small".equals(s)) return SpotSize.SMALL;
        if ("M".equals(s) || "medium".equals(s)) return SpotSize.MEDIUM;
        return SpotSize.LARGE;
    }

    private static VehicleType vtypeFrom(String s) {
        if ("M".equals(s) || "moto".equals(s) || "motorcycle".equals(s)) return VehicleType.MOTORCYCLE;
        if ("C".equals(s) || "car".equals(s)) return VehicleType.CAR;
        return VehicleType.TRUCK;
    }

    private static GateType gateFrom(String s) {
        return "entry".equals(s) ? GateType.ENTRY : GateType.EXIT;
    }

    private static String feeToStr(double f) {
        if (f < 0) return "-1";
        return String.format(java.util.Locale.US, "%.2f", f);
    }

    public static List<String> parking_simulate(List<ParkOp> ops) {
        List<String> out = new ArrayList<>();
        ParkingLot lot = null;
        String[] tickets = new String[16];
        TicketSnap[] snaps = new TicketSnap[16];
        for (int i = 0; i < 16; i++) { tickets[i] = ""; snaps[i] = new TicketSnap(); }

        for (ParkOp op : ops) {
            String k = op.kind;
            if ("new".equals(k)) {
                lot = new ParkingLot(op.i1);
                for (int i = 0; i < 16; i++) { tickets[i] = ""; snaps[i] = new TicketSnap(); }
                out.add("ok");
            } else if ("add_spot".equals(k)) {
                lot.addSpot(op.i1, sizeFrom(op.s1));
                out.add("ok");
            } else if ("add_gate".equals(k)) {
                lot.addGate(op.s1, gateFrom(op.s2));
                out.add("ok");
            } else if ("gates_count".equals(k)) {
                out.add(Integer.toString(lot.getGates(gateFrom(op.s1)).size()));
            } else if ("gate_at".equals(k)) {
                List<String> g = lot.getGates(gateFrom(op.s1));
                out.add(op.i1 >= 0 && op.i1 < g.size() ? g.get(op.i1) : "");
            } else if ("set_pricing".equals(k)) {
                PricingStrategy p = null;
                if ("flat".equals(op.s1)) p = new FlatRate((double) op.i1);
                else if ("hourly".equals(op.s1)) p = new Hourly((double) op.i1);
                else if ("tiered".equals(op.s1)) p = new Tiered((double) op.i1, (double) op.i2, (double) op.i3);
                if (p != null) lot.setPricingStrategy(p);
                out.add("ok");
            } else if ("park".equals(k)) {
                Vehicle v = new Vehicle(op.s1, vtypeFrom(op.s2));
                Ticket t = lot.parkVehicle(v, (long) op.i1, op.s3);
                if (op.i2 >= 0 && op.i2 < tickets.length) {
                    tickets[op.i2] = t != null ? t.ticketId : "";
                    if (t != null) {
                        snaps[op.i2].id = t.ticketId;
                        snaps[op.i2].floor = t.floor;
                        snaps[op.i2].spotId = t.spotId;
                        snaps[op.i2].entryGate = t.entryGateId;
                    } else {
                        snaps[op.i2] = new TicketSnap();
                    }
                }
                out.add(t != null ? t.ticketId : "");
            } else if ("ticket_at".equals(k)) {
                out.add(op.i1 >= 0 && op.i1 < tickets.length ? tickets[op.i1] : "");
            } else if ("ticket_floor".equals(k)) {
                out.add(op.i1 >= 0 && op.i1 < snaps.length ? Integer.toString(snaps[op.i1].floor) : "-1");
            } else if ("ticket_spot_id".equals(k)) {
                out.add(op.i1 >= 0 && op.i1 < snaps.length ? snaps[op.i1].spotId : "");
            } else if ("ticket_entry".equals(k)) {
                out.add(op.i1 >= 0 && op.i1 < snaps.length ? snaps[op.i1].entryGate : "");
            } else if ("unpark".equals(k)) {
                String tid = op.i1 >= 0 && op.i1 < tickets.length ? tickets[op.i1] : "";
                double fee = lot.unparkVehicle(tid, (long) op.i2, op.s1);
                out.add(feeToStr(fee));
            } else if ("unpark_id".equals(k)) {
                double fee = lot.unparkVehicle(op.s1, (long) op.i2, op.s2);
                out.add(feeToStr(fee));
            } else if ("available".equals(k)) {
                out.add(Integer.toString(lot.getAvailableSpots(sizeFrom(op.s1))));
            } else if ("available_floor".equals(k)) {
                out.add(Integer.toString(lot.getAvailableSpotsByFloor(op.i1, sizeFrom(op.s1))));
            } else {
                out.add("unknown:" + k);
            }
        }
        return out;
    }
}

// Ride Sharing — Solution (Java)
import java.util.*;

class RideOp {
    public String kind;
    public String s1, s2, s3, s4;
    public int i1, i2;
    public RideOp(String kind, String s1, String s2, String s3, String s4, int i1, int i2) {
        this.kind = kind; this.s1 = s1; this.s2 = s2; this.s3 = s3; this.s4 = s4;
        this.i1 = i1; this.i2 = i2;
    }
}

class User {
    public String id, name;
    public int ridesOffered, ridesTaken;
    public User(String id, String name) { this.id = id; this.name = name; }
}

class Vehicle {
    public String id, ownerId, model, regNumber;
    public Vehicle(String id, String ownerId, String model, String reg) {
        this.id = id; this.ownerId = ownerId; this.model = model; this.regNumber = reg;
    }
}

class Ride {
    public String id;
    public String driverId;
    public String vehicleId;
    public String origin;
    public String destination;
    public int totalSeats;
    public int availableSeats;
    public boolean active;
    public Ride(String id, String driverId, String vehicleId, String origin, String dest, int seats) {
        this.id = id; this.driverId = driverId; this.vehicleId = vehicleId;
        this.origin = origin; this.destination = dest;
        this.totalSeats = seats; this.availableSeats = seats; this.active = true;
    }
}

interface RideSelectionStrategy {
    Ride select(List<Ride> candidates, int seatsNeeded, String preference);
}

class MostVacantStrategy implements RideSelectionStrategy {
    @Override public Ride select(List<Ride> candidates, int seatsNeeded, String preference) {
        Ride best = null;
        for (Ride r : candidates) {
            if (r.availableSeats >= seatsNeeded) {
                if (best == null || r.availableSeats > best.availableSeats) best = r;
            }
        }
        return best;
    }
}

class PreferredVehicleStrategy implements RideSelectionStrategy {
    private final Map<String, Vehicle> vehicleStore;
    public PreferredVehicleStrategy(Map<String, Vehicle> vs) { this.vehicleStore = vs; }
    @Override public Ride select(List<Ride> candidates, int seatsNeeded, String preference) {
        for (Ride r : candidates) {
            if (r.availableSeats >= seatsNeeded) {
                Vehicle v = vehicleStore.get(r.vehicleId);
                if (v != null && v.model.equals(preference)) return r;
            }
        }
        return null;
    }
}

class RideService {
    Map<String, User> users = new LinkedHashMap<>();
    Map<String, Vehicle> vehicles = new LinkedHashMap<>();
    Map<String, Ride> rides = new LinkedHashMap<>();
    Map<String, String> activeVehicles = new LinkedHashMap<>();
    int rideCounter = 0;

    public void addUser(String name) {
        if (users.containsKey(name)) return;
        users.put(name, new User(name, name));
    }

    public void addVehicle(String userName, String model, String regNumber) {
        if (!users.containsKey(userName)) return;
        vehicles.put(regNumber, new Vehicle(regNumber, userName, model, regNumber));
    }

    public String offerRide(String userName, String origin, String dest, int seats, String vehicleRegNumber) {
        if (!users.containsKey(userName)) return "";
        Vehicle v = vehicles.get(vehicleRegNumber);
        if (v == null) return "";
        if (!v.ownerId.equals(userName)) return "";
        if (activeVehicles.containsKey(vehicleRegNumber)) return "";

        String rideId = "RIDE-" + (++rideCounter);
        rides.put(rideId, new Ride(rideId, userName, vehicleRegNumber, origin, dest, seats));
        activeVehicles.put(vehicleRegNumber, rideId);
        users.get(userName).ridesOffered++;
        return rideId;
    }

    public String selectRide(String passengerName, String origin, String dest, int seats,
                              RideSelectionStrategy strategy, String preference) {
        if (!users.containsKey(passengerName)) return "";

        List<Ride> candidates = new ArrayList<>();
        for (Ride r : rides.values()) {
            if (r.active && r.origin.equals(origin) && r.destination.equals(dest)
                    && r.availableSeats >= seats && !r.driverId.equals(passengerName)) {
                candidates.add(r);
            }
        }
        Ride selected = strategy.select(candidates, seats, preference);
        if (selected != null) {
            selected.availableSeats -= seats;
            users.get(passengerName).ridesTaken++;
            return selected.id;
        }
        return "";
    }

    public void endRide(String rideId) {
        Ride r = rides.get(rideId);
        if (r == null || !r.active) return;
        r.active = false;
        activeVehicles.remove(r.vehicleId);
    }

    public Map<String, Vehicle> getVehicles() { return vehicles; }

    public boolean hasUser(String name) { return users.containsKey(name); }
    public boolean hasVehicle(String reg) { return vehicles.containsKey(reg); }
    public boolean hasRide(String id) { return rides.containsKey(id); }
    public User getUser(String name) { return users.get(name); }
    public Ride getRide(String id) { return rides.get(id); }
}

public class Solution {
    public static List<String> ride_simulate(List<RideOp> ops) {
        List<String> out = new ArrayList<>();
        RideService svc = new RideService();
        String[] rideSlots = new String[32];
        for (int i = 0; i < 32; i++) rideSlots[i] = "";

        for (RideOp op : ops) {
            String k = op.kind;
            if ("new".equals(k)) {
                svc = new RideService();
                for (int i = 0; i < 32; i++) rideSlots[i] = "";
                out.add("ok");
            } else if ("add_user".equals(k)) {
                svc.addUser(op.s1); out.add("ok");
            } else if ("add_veh".equals(k)) {
                svc.addVehicle(op.s1, op.s2, op.s3); out.add("ok");
            } else if ("offer".equals(k)) {
                String rid = svc.offerRide(op.s1, op.s2, op.s3, op.i1, op.s4);
                if (op.i2 >= 0 && op.i2 < rideSlots.length) rideSlots[op.i2] = rid;
                out.add(rid);
            } else if ("ride_active".equals(k)) {
                String rid = rideSlots[op.i2];
                out.add(svc.hasRide(rid) && svc.getRide(rid).active ? "yes" : "no");
            } else if ("ride_origin".equals(k)) {
                out.add(svc.hasRide(rideSlots[op.i2]) ? svc.getRide(rideSlots[op.i2]).origin : "");
            } else if ("ride_dest".equals(k)) {
                out.add(svc.hasRide(rideSlots[op.i2]) ? svc.getRide(rideSlots[op.i2]).destination : "");
            } else if ("ride_total".equals(k)) {
                out.add(Integer.toString(svc.hasRide(rideSlots[op.i2]) ? svc.getRide(rideSlots[op.i2]).totalSeats : -1));
            } else if ("ride_avail".equals(k)) {
                out.add(Integer.toString(svc.hasRide(rideSlots[op.i2]) ? svc.getRide(rideSlots[op.i2]).availableSeats : -1));
            } else if ("ride_driver".equals(k)) {
                out.add(svc.hasRide(rideSlots[op.i2]) ? svc.getRide(rideSlots[op.i2]).driverId : "");
            } else if ("select_mv".equals(k)) {
                String rid = svc.selectRide(op.s1, op.s2, op.s3, op.i1, new MostVacantStrategy(), "");
                if (op.i2 >= 0 && op.i2 < rideSlots.length) rideSlots[op.i2] = rid;
                out.add(rid);
            } else if ("select_pv".equals(k)) {
                String rid = svc.selectRide(op.s1, op.s2, op.s3, op.i1,
                        new PreferredVehicleStrategy(svc.getVehicles()), op.s4);
                if (op.i2 >= 0 && op.i2 < rideSlots.length) rideSlots[op.i2] = rid;
                out.add(rid);
            } else if ("end".equals(k)) {
                svc.endRide(rideSlots[op.i2]); out.add("ok");
            } else if ("end_id".equals(k)) {
                svc.endRide(op.s1); out.add("ok");
            } else if ("user_offered".equals(k)) {
                out.add(svc.hasUser(op.s1) ? Integer.toString(svc.getUser(op.s1).ridesOffered) : "0");
            } else if ("user_taken".equals(k)) {
                out.add(svc.hasUser(op.s1) ? Integer.toString(svc.getUser(op.s1).ridesTaken) : "0");
            } else if ("has_user".equals(k)) {
                out.add(svc.hasUser(op.s1) ? "yes" : "no");
            } else if ("has_vehicle".equals(k)) {
                out.add(svc.hasVehicle(op.s1) ? "yes" : "no");
            } else if ("has_ride".equals(k)) {
                out.add(svc.hasRide(rideSlots[op.i2]) ? "yes" : "no");
            } else {
                out.add("unknown:" + k);
            }
        }
        return out;
    }
}

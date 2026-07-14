import java.util.*;

// Data class (given).
class PricingContext {
    public double baseFare;
    public int availableDrivers;
    public int activeRideRequests;
    public String timeOfDay;
    public String weather;

    public PricingContext(double baseFare, int availableDrivers, int activeRideRequests, String timeOfDay, String weather) {
        this.baseFare = baseFare;
        this.availableDrivers = availableDrivers;
        this.activeRideRequests = activeRideRequests;
        this.timeOfDay = timeOfDay;
        this.weather = weather;
    }
}

class RideRequest {
    public String userId;
    public String pickup;
    public String dropoff;
    public String rideType;

    public RideRequest(String userId, String pickup, String dropoff, String rideType) {
        this.userId = userId;
        this.pickup = pickup;
        this.dropoff = dropoff;
        this.rideType = rideType;
    }
}

public class Solution {
    public static double calculateSurge(PricingContext ctx) {
        // TODO: write your solution
        return 0.0;
    }

    public static double calculateFare(RideRequest req, PricingContext ctx) {
        // TODO: write your solution
        return 0.0;
    }

}

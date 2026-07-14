// Ride Surge Pricing — Solution (Java)
import java.util.*;

// ─── Data Models ─────────────────────────────────────────────────────────────

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

class PricingContext {
    public double baseFare;
    public int availableDrivers;
    public int activeRideRequests;
    public String timeOfDay;
    public String weather;

    public PricingContext(double baseFare, int availableDrivers, int activeRideRequests,
                          String timeOfDay, String weather) {
        this.baseFare = baseFare;
        this.availableDrivers = availableDrivers;
        this.activeRideRequests = activeRideRequests;
        this.timeOfDay = timeOfDay;
        this.weather = weather;
    }
}

// ─── Strategy Interface ──────────────────────────────────────────────────────

interface SurgeStrategy {
    double multiplier(PricingContext ctx);
}

class DemandSurge implements SurgeStrategy {
    @Override
    public double multiplier(PricingContext ctx) {
        if (ctx.availableDrivers == 0) return 2.5;
        double ratio = (double) ctx.activeRideRequests / ctx.availableDrivers;
        if (ratio > 3.0) return 2.0;
        if (ratio > 2.0) return 1.5;
        if (ratio > 1.5) return 1.25;
        return 1.0;
    }
}

class WeatherSurge implements SurgeStrategy {
    @Override
    public double multiplier(PricingContext ctx) {
        if ("storm".equals(ctx.weather)) return 1.8;
        if ("rain".equals(ctx.weather)) return 1.3;
        return 1.0;
    }
}

class TimeSurge implements SurgeStrategy {
    @Override
    public double multiplier(PricingContext ctx) {
        if ("evening".equals(ctx.timeOfDay)) return 1.5;
        if ("morning".equals(ctx.timeOfDay)) return 1.2;
        return 1.0;
    }
}

// ─── Observer Interface ──────────────────────────────────────────────────────

interface SurgeObserver {
    void onSurgeChange(double oldMultiplier, double newMultiplier, String rideType);
}

// ─── Pricing Engine ──────────────────────────────────────────────────────────

class PricingEngine {
    private final List<SurgeStrategy> strategies = new ArrayList<>();
    private final List<SurgeObserver> observers = new ArrayList<>();
    private double lastSurge = 1.0;
    private static final double CHANGE_THRESHOLD = 0.5;

    public void addStrategy(SurgeStrategy s) { strategies.add(s); }
    public void addObserver(SurgeObserver o) { observers.add(o); }
    public void clearObservers() { observers.clear(); }

    public double calculateSurge(PricingContext ctx, String rideType) {
        double mult = 1.0;
        for (SurgeStrategy s : strategies) {
            mult = Math.max(mult, s.multiplier(ctx));
        }
        mult = Math.min(mult, 3.0);

        if (Math.abs(mult - lastSurge) > CHANGE_THRESHOLD) {
            for (SurgeObserver o : observers) {
                o.onSurgeChange(lastSurge, mult, rideType);
            }
        }
        lastSurge = mult;
        return mult;
    }

    public double calculateFare(PricingContext ctx, String rideType) {
        return ctx.baseFare * calculateSurge(ctx, rideType);
    }
}

// ─── Solution: free-function wrappers ────────────────────────────────────────

public class Solution {
    private static PricingEngine engine = null;

    private static PricingEngine getGlobalEngine() {
        if (engine == null) {
            engine = new PricingEngine();
            engine.addStrategy(new DemandSurge());
            engine.addStrategy(new WeatherSurge());
            engine.addStrategy(new TimeSurge());
        }
        return engine;
    }

    public static double calculateSurge(PricingContext ctx) {
        return getGlobalEngine().calculateSurge(ctx, "all");
    }

    public static double calculateFare(RideRequest req, PricingContext ctx) {
        return getGlobalEngine().calculateFare(ctx, req.rideType);
    }

    public static void registerSurgeObserver(SurgeObserver obs) {
        getGlobalEngine().clearObservers();
        getGlobalEngine().addObserver(obs);
    }
}

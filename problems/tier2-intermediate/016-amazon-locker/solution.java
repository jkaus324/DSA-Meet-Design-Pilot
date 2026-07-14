// Amazon Locker — Solution (Java)
import java.util.*;

class LockerOp {
    public String kind;
    public String s1;
    public String s2;
    public int i1;
    public int i2;

    public LockerOp(String kind, String s1, String s2, int i1, int i2) {
        this.kind = kind; this.s1 = s1; this.s2 = s2; this.i1 = i1; this.i2 = i2;
    }
}

enum LockerSize { SMALL, MEDIUM, LARGE }

class Locker {
    public String lockerId;
    public LockerSize size;
    public boolean occupied;
    public Locker(String id, LockerSize sz) { this.lockerId = id; this.size = sz; this.occupied = false; }
}

class DepositRecord {
    public String lockerId;
    public String packageId;
    public String pickupCode;
    public long depositTime;
    public DepositRecord(String l, String p, String c, long t) {
        lockerId = l; packageId = p; pickupCode = c; depositTime = t;
    }
}

interface LockerAllocationStrategy {
    String allocate(LockerSize packageSize, Map<LockerSize, Deque<String>> available);
}

class SmallestFitStrategy implements LockerAllocationStrategy {
    @Override
    public String allocate(LockerSize packageSize, Map<LockerSize, Deque<String>> available) {
        List<LockerSize> tryOrder = new ArrayList<>();
        if (packageSize == LockerSize.SMALL) {
            tryOrder.add(LockerSize.SMALL); tryOrder.add(LockerSize.MEDIUM); tryOrder.add(LockerSize.LARGE);
        } else if (packageSize == LockerSize.MEDIUM) {
            tryOrder.add(LockerSize.MEDIUM); tryOrder.add(LockerSize.LARGE);
        } else {
            tryOrder.add(LockerSize.LARGE);
        }
        for (LockerSize sz : tryOrder) {
            Deque<String> q = available.get(sz);
            if (q != null && !q.isEmpty()) return q.pollFirst();
        }
        return "";
    }
}

interface NotificationChannel {
    void notify(String packageId, String message);
}

class LockerSystem {
    Map<String, Locker> lockers = new LinkedHashMap<>();
    Map<LockerSize, Deque<String>> available = new EnumMap<>(LockerSize.class);
    Map<String, DepositRecord> activeDeposits = new LinkedHashMap<>();
    LockerAllocationStrategy strategy = new SmallestFitStrategy();
    List<NotificationChannel> channels = new ArrayList<>();
    int codeCounter = 0;
    int expiryHours = 0;

    public LockerSystem() {
        for (LockerSize s : LockerSize.values()) available.put(s, new ArrayDeque<>());
    }

    private String generateCode() {
        return "CODE-" + (++codeCounter);
    }

    private void notifyAll(String packageId, String message) {
        for (NotificationChannel ch : channels) ch.notify(packageId, message);
    }

    private void freeLocker(String lockerId) {
        Locker l = lockers.get(lockerId);
        if (l != null) {
            l.occupied = false;
            available.get(l.size).addLast(lockerId);
        }
    }

    public void addLocker(String lockerId, LockerSize size) {
        lockers.put(lockerId, new Locker(lockerId, size));
        available.get(size).addLast(lockerId);
    }

    public String depositPackage(String packageId, LockerSize size, long depositTime) {
        String lockerId = strategy.allocate(size, available);
        if (lockerId.isEmpty()) return "";
        lockers.get(lockerId).occupied = true;
        String code = generateCode();
        activeDeposits.put(code, new DepositRecord(lockerId, packageId, code, depositTime));
        notifyAll(packageId, "Package " + packageId + " deposited. Code: " + code);
        return code;
    }

    public boolean retrievePackage(String code) {
        DepositRecord rec = activeDeposits.get(code);
        if (rec == null) return false;
        freeLocker(rec.lockerId);
        activeDeposits.remove(code);
        return true;
    }

    public void setCodeExpiry(int hours) { this.expiryHours = hours; }

    public List<String> checkExpired(long currentTime) {
        List<String> expired = new ArrayList<>();
        if (expiryHours <= 0) return expired;
        Iterator<Map.Entry<String, DepositRecord>> it = activeDeposits.entrySet().iterator();
        while (it.hasNext()) {
            Map.Entry<String, DepositRecord> e = it.next();
            DepositRecord rec = e.getValue();
            if (currentTime - rec.depositTime > (long) expiryHours * 3600) {
                freeLocker(rec.lockerId);
                expired.add(rec.packageId);
                notifyAll(rec.packageId, "Package " + rec.packageId + " expired. Locker freed.");
                it.remove();
            }
        }
        return expired;
    }

    public void addNotificationChannel(NotificationChannel channel) {
        channels.add(channel);
    }
}

class CapturingChannel implements NotificationChannel {
    public List<String> log;
    public CapturingChannel(List<String> log) { this.log = log; }
    @Override
    public void notify(String packageId, String message) {
        log.add(packageId + ": " + message);
    }
}

public class Solution {
    private static LockerSystem system_instance = null;
    private static List<String> chanLog = new ArrayList<>();
    private static CapturingChannel currentChan = null;

    private static LockerSize lsizeFrom(String s) {
        if ("S".equals(s)) return LockerSize.SMALL;
        if ("M".equals(s)) return LockerSize.MEDIUM;
        return LockerSize.LARGE;
    }

    public static List<String> locker_simulate(List<LockerOp> ops) {
        List<String> out = new ArrayList<>();
        String[] codes = new String[32];
        for (int i = 0; i < 32; i++) codes[i] = "";
        List<String> lastExpired = new ArrayList<>();
        chanLog = new ArrayList<>();
        currentChan = null;
        system_instance = null;

        for (LockerOp op : ops) {
            String k = op.kind;
            if ("new".equals(k)) {
                system_instance = new LockerSystem();
                for (int i = 0; i < 32; i++) codes[i] = "";
                chanLog = new ArrayList<>();
                currentChan = null;
                lastExpired = new ArrayList<>();
                out.add("ok");
            } else if ("add_locker".equals(k)) {
                if (system_instance == null) system_instance = new LockerSystem();
                system_instance.addLocker(op.s1, lsizeFrom(op.s2));
                out.add("ok");
            } else if ("deposit".equals(k)) {
                if (system_instance == null) system_instance = new LockerSystem();
                String code = system_instance.depositPackage(op.s1, lsizeFrom(op.s2), (long) op.i1);
                if (op.i2 >= 0 && op.i2 < codes.length) codes[op.i2] = code;
                out.add(code);
            } else if ("code_at".equals(k)) {
                out.add(codes[op.i2]);
            } else if ("retrieve".equals(k)) {
                out.add(system_instance != null && system_instance.retrievePackage(codes[op.i2]) ? "ok" : "fail");
            } else if ("retrieve_id".equals(k)) {
                out.add(system_instance != null && system_instance.retrievePackage(op.s1) ? "ok" : "fail");
            } else if ("set_expiry".equals(k)) {
                if (system_instance == null) system_instance = new LockerSystem();
                system_instance.setCodeExpiry(op.i1);
                out.add("ok");
            } else if ("check_expired".equals(k)) {
                if (system_instance == null) {
                    lastExpired = new ArrayList<>();
                } else {
                    lastExpired = system_instance.checkExpired((long) op.i1);
                }
                out.add(Integer.toString(lastExpired.size()));
            } else if ("expired_at".equals(k)) {
                out.add(op.i2 >= 0 && op.i2 < lastExpired.size() ? lastExpired.get(op.i2) : "");
            } else if ("add_chan".equals(k)) {
                if (system_instance == null) system_instance = new LockerSystem();
                currentChan = new CapturingChannel(chanLog);
                system_instance.addNotificationChannel(currentChan);
                out.add("ok");
            } else if ("chan_log_size".equals(k)) {
                out.add(Integer.toString(chanLog.size()));
            } else if ("chan_log_contains".equals(k)) {
                boolean found = false;
                for (String l : chanLog) if (l.contains(op.s1)) { found = true; break; }
                out.add(found ? "yes" : "no");
            } else {
                out.add("unknown:" + k);
            }
        }
        return out;
    }
}

// Notification System — Solution (Java)
import java.util.*;

class User {
    public String id;
    public String email;
    public String phone;
    public List<String> subscribedChannels;

    public User(String id, String email, String phone, List<String> subscribedChannels) {
        this.id = id;
        this.email = email;
        this.phone = phone;
        this.subscribedChannels = subscribedChannels;
    }
}

interface NotificationObserver {
    String channelName();
    void send(String userId, String message);
    void update(String event, String priority, User user);
}

class EmailNotifier implements NotificationObserver {
    public String channelName() { return "email"; }
    public void send(String userId, String message) {}
    public void update(String event, String priority, User user) {}
}

class SMSNotifier implements NotificationObserver {
    public String channelName() { return "sms"; }
    public void send(String userId, String message) {}
    public void update(String event, String priority, User user) {}
}

class PushNotifier implements NotificationObserver {
    public String channelName() { return "push"; }
    public void send(String userId, String message) {}
    public void update(String event, String priority, User user) {}
}

class PriorityFilteredObserver implements NotificationObserver {
    private final NotificationObserver inner;
    private final String minPriority;

    public PriorityFilteredObserver(NotificationObserver inner, String minPriority) {
        this.inner = inner;
        this.minPriority = minPriority;
    }

    public String channelName() { return inner.channelName(); }
    public void send(String userId, String message) { inner.send(userId, message); }

    public void update(String event, String priority, User user) {
        if (Solution.priorityLevel(priority) >= Solution.priorityLevel(minPriority)) {
            inner.update(event, priority, user);
        }
    }
}

class NotificationManager {
    private final List<NotificationObserver> observers = new ArrayList<>();

    public void subscribe(NotificationObserver obs) { observers.add(obs); }

    public void unsubscribe(String channel) {
        observers.removeIf(o -> o.channelName().equals(channel));
    }

    public void notify(String event, List<User> users) {
        for (User u : users) {
            for (NotificationObserver obs : observers) {
                if (u.subscribedChannels.contains(obs.channelName())) {
                    obs.send(u.id, event);
                }
            }
        }
    }

    public void notifyAll(String event, String priority, List<User> users) {
        for (User u : users) {
            for (NotificationObserver obs : observers) {
                if (u.subscribedChannels.contains(obs.channelName())) {
                    obs.update(event, priority, u);
                }
            }
        }
    }
}

public class Solution {
    public static final List<String> PRIORITY_ORDER =
            Arrays.asList("promotional", "info", "critical");

    public static int priorityLevel(String p) {
        int idx = PRIORITY_ORDER.indexOf(p);
        return idx < 0 ? 0 : idx;
    }

    public static void reset_service() {}

    public static void notify_event(String event, List<String> userIds, List<String> subscribedChannels) {
        List<User> users = new ArrayList<>();
        for (String uid : userIds) {
            users.add(new User(uid, uid + "@test.com", "+1-555-0000",
                    new ArrayList<>(subscribedChannels)));
        }
        NotificationManager mgr = new NotificationManager();
        mgr.subscribe(new EmailNotifier());
        mgr.subscribe(new SMSNotifier());
        mgr.subscribe(new PushNotifier());
        mgr.notify(event, users);
    }

    public static void notify_priority(String event, String priority,
                                       List<String> userIds, List<String> subscribedChannels,
                                       String minPriority) {
        List<User> users = new ArrayList<>();
        for (String uid : userIds) {
            users.add(new User(uid, uid + "@test.com", "+1-555-0000",
                    new ArrayList<>(subscribedChannels)));
        }
        String minP = (minPriority == null || minPriority.isEmpty()) ? "promotional" : minPriority;
        NotificationManager mgr = new NotificationManager();
        mgr.subscribe(new PriorityFilteredObserver(new EmailNotifier(), minP));
        mgr.subscribe(new PriorityFilteredObserver(new SMSNotifier(), minP));
        mgr.subscribe(new PriorityFilteredObserver(new PushNotifier(), minP));
        mgr.notifyAll(event, priority, users);
    }

    public static int notify_priority_level(String priority) {
        return priorityLevel(priority);
    }
}

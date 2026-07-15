// Notification System — Observer + Decorator reference solution (JavaScript).

const PRIORITY_ORDER = ["promotional", "info", "critical"];

function priorityLevel(p) {
  const idx = PRIORITY_ORDER.indexOf(p);
  return idx >= 0 ? idx : 0;
}

class User {
  constructor(id, email, phone, subscribedChannels) {
    this.id = id;
    this.email = email;
    this.phone = phone;
    this.subscribedChannels = subscribedChannels;
  }
}

class NotificationObserver {
  channelName() { throw new Error('not implemented'); }
  send(userId, message) { throw new Error('not implemented'); }
  update(event, priority, user) { throw new Error('not implemented'); }
}

class EmailNotifier extends NotificationObserver {
  channelName() { return "email"; }
  send(userId, message) {}
  update(event, priority, user) {}
}

class SMSNotifier extends NotificationObserver {
  channelName() { return "sms"; }
  send(userId, message) {}
  update(event, priority, user) {}
}

class PushNotifier extends NotificationObserver {
  channelName() { return "push"; }
  send(userId, message) {}
  update(event, priority, user) {}
}

class PriorityFilteredObserver extends NotificationObserver {
  constructor(inner, minPriority) {
    super();
    this.inner = inner;
    this.minPriority = minPriority;
  }
  channelName() { return this.inner.channelName(); }
  send(userId, message) { this.inner.send(userId, message); }
  update(event, priority, user) {
    if (priorityLevel(priority) >= priorityLevel(this.minPriority)) {
      this.inner.update(event, priority, user);
    }
  }
}

class NotificationManager {
  constructor() {
    this.observers = [];
  }
  subscribe(obs) { this.observers.push(obs); }
  unsubscribe(channel) {
    this.observers = this.observers.filter(o => o.channelName() !== channel);
  }
  notify(event, users) {
    for (const u of users) {
      for (const obs of this.observers) {
        if (u.subscribedChannels.includes(obs.channelName())) {
          obs.send(u.id, event);
        }
      }
    }
  }
  notifyAll(event, priority, users) {
    for (const u of users) {
      for (const obs of this.observers) {
        if (u.subscribedChannels.includes(obs.channelName())) {
          obs.update(event, priority, u);
        }
      }
    }
  }
}

function notify(event, users) {
  const mgr = new NotificationManager();
  mgr.subscribe(new EmailNotifier());
  mgr.subscribe(new SMSNotifier());
  mgr.subscribe(new PushNotifier());
  mgr.notify(event, users);
}

function notify_with_priority(event, priority, users, userMinPriority) {
  const minP = userMinPriority.has("*") ? userMinPriority.get("*") : "promotional";
  const mgr = new NotificationManager();
  mgr.subscribe(new PriorityFilteredObserver(new EmailNotifier(), minP));
  mgr.subscribe(new PriorityFilteredObserver(new SMSNotifier(), minP));
  mgr.subscribe(new PriorityFilteredObserver(new PushNotifier(), minP));
  mgr.notifyAll(event, priority, users);
}

function reset_service() {}

function notify_event(event, userIds, subscribedChannels) {
  const users = userIds.map(uid =>
    new User(uid, `${uid}@test.com`, "+1-555-0000", [...subscribedChannels]));
  notify(event, users);
}

function notify_priority(event, priority, userIds, subscribedChannels, minPriority) {
  const users = userIds.map(uid =>
    new User(uid, `${uid}@test.com`, "+1-555-0000", [...subscribedChannels]));
  const prefs = new Map();
  if (minPriority) {
    prefs.set("*", minPriority);
  }
  notify_with_priority(event, priority, users, prefs);
}

function notify_priority_level(p) {
  return priorityLevel(p);
}

module.exports = {
  User,
  NotificationObserver,
  EmailNotifier,
  SMSNotifier,
  PushNotifier,
  PriorityFilteredObserver,
  NotificationManager,
  reset_service,
  notify_event,
  notify_priority,
  notify_priority_level,
};

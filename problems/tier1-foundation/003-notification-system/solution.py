"""Notification System — Observer + Decorator reference solution (Python)."""

from abc import ABC, abstractmethod


PRIORITY_ORDER = ["promotional", "info", "critical"]


def priorityLevel(p):
    if p in PRIORITY_ORDER:
        return PRIORITY_ORDER.index(p)
    return 0


class User:
    def __init__(self, id, email, phone, subscribedChannels):
        self.id = id
        self.email = email
        self.phone = phone
        self.subscribedChannels = subscribedChannels


class NotificationObserver(ABC):
    @abstractmethod
    def channelName(self):
        ...

    @abstractmethod
    def send(self, userId, message):
        ...

    @abstractmethod
    def update(self, event, priority, user):
        ...


class EmailNotifier(NotificationObserver):
    def channelName(self):
        return "email"

    def send(self, userId, message):
        pass

    def update(self, event, priority, user):
        pass


class SMSNotifier(NotificationObserver):
    def channelName(self):
        return "sms"

    def send(self, userId, message):
        pass

    def update(self, event, priority, user):
        pass


class PushNotifier(NotificationObserver):
    def channelName(self):
        return "push"

    def send(self, userId, message):
        pass

    def update(self, event, priority, user):
        pass


class PriorityFilteredObserver(NotificationObserver):
    def __init__(self, inner, minPriority):
        self.inner = inner
        self.minPriority = minPriority

    def channelName(self):
        return self.inner.channelName()

    def send(self, userId, message):
        self.inner.send(userId, message)

    def update(self, event, priority, user):
        if priorityLevel(priority) >= priorityLevel(self.minPriority):
            self.inner.update(event, priority, user)


class NotificationManager:
    def __init__(self):
        self.observers = []

    def subscribe(self, obs):
        self.observers.append(obs)

    def unsubscribe(self, channel):
        self.observers = [o for o in self.observers if o.channelName() != channel]

    def notify(self, event, users):
        for u in users:
            for obs in self.observers:
                if obs.channelName() in u.subscribedChannels:
                    obs.send(u.id, event)

    def notifyAll(self, event, priority, users):
        for u in users:
            for obs in self.observers:
                if obs.channelName() in u.subscribedChannels:
                    obs.update(event, priority, u)


def notify(event, users):
    mgr = NotificationManager()
    mgr.subscribe(EmailNotifier())
    mgr.subscribe(SMSNotifier())
    mgr.subscribe(PushNotifier())
    mgr.notify(event, users)


def notify_with_priority(event, priority, users, userMinPriority):
    minP = userMinPriority.get("*", "promotional")
    mgr = NotificationManager()
    mgr.subscribe(PriorityFilteredObserver(EmailNotifier(), minP))
    mgr.subscribe(PriorityFilteredObserver(SMSNotifier(), minP))
    mgr.subscribe(PriorityFilteredObserver(PushNotifier(), minP))
    mgr.notifyAll(event, priority, users)


def reset_service():
    pass


def notify_event(event, userIds, subscribedChannels):
    users = [User(uid, f"{uid}@test.com", "+1-555-0000", list(subscribedChannels)) for uid in userIds]
    notify(event, users)


def notify_priority(event, priority, userIds, subscribedChannels, minPriority):
    users = [User(uid, f"{uid}@test.com", "+1-555-0000", list(subscribedChannels)) for uid in userIds]
    prefs = {}
    if minPriority:
        prefs["*"] = minPriority
    notify_with_priority(event, priority, users, prefs)


def notify_priority_level(p):
    return priorityLevel(p)

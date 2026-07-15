// Notification System — Observer + Decorator reference solution (Go).
package main

var priorityOrder = []string{"promotional", "info", "critical"}

func priorityLevel(p string) int {
	for i, v := range priorityOrder {
		if v == p {
			return i
		}
	}
	return 0
}

type nsUser struct {
	id                 string
	subscribedChannels []string
}

type notificationObserver interface {
	channelName() string
	update(event, priority string, user nsUser)
}

type emailNotifier struct{}

func (emailNotifier) channelName() string                       { return "email" }
func (emailNotifier) update(event, priority string, u nsUser)   {}

type smsNotifier struct{}

func (smsNotifier) channelName() string                     { return "sms" }
func (smsNotifier) update(event, priority string, u nsUser) {}

type pushNotifier struct{}

func (pushNotifier) channelName() string                     { return "push" }
func (pushNotifier) update(event, priority string, u nsUser) {}

type priorityFilteredObserver struct {
	inner       notificationObserver
	minPriority string
}

func (p priorityFilteredObserver) channelName() string { return p.inner.channelName() }
func (p priorityFilteredObserver) update(event, priority string, u nsUser) {
	if priorityLevel(priority) >= priorityLevel(p.minPriority) {
		p.inner.update(event, priority, u)
	}
}

func contains(list []string, s string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}

func reset_service() {}

func notify_event(event string, userIds []string, subscribedChannels []string) {
	observers := []notificationObserver{emailNotifier{}, smsNotifier{}, pushNotifier{}}
	for _, uid := range userIds {
		u := nsUser{id: uid, subscribedChannels: subscribedChannels}
		for _, obs := range observers {
			if contains(u.subscribedChannels, obs.channelName()) {
				// send(u.id, event) — no observable effect
			}
		}
	}
}

func notify_priority(event, priority string, userIds []string, subscribedChannels []string, minPriority string) {
	minP := minPriority
	if minP == "" {
		minP = "promotional"
	}
	observers := []notificationObserver{
		priorityFilteredObserver{emailNotifier{}, minP},
		priorityFilteredObserver{smsNotifier{}, minP},
		priorityFilteredObserver{pushNotifier{}, minP},
	}
	for _, uid := range userIds {
		u := nsUser{id: uid, subscribedChannels: subscribedChannels}
		for _, obs := range observers {
			if contains(u.subscribedChannels, obs.channelName()) {
				obs.update(event, priority, u)
			}
		}
	}
}

func notify_priority_level(p string) int {
	return priorityLevel(p)
}

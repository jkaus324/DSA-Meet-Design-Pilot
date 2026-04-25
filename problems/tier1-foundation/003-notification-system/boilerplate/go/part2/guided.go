package main

import "fmt"

// ─── Data Model (given — do not modify) ─────────────────────────────────────

type Notification struct{ UserID, Message, Channel string }
type User struct {
	ID                 string
	Email              string
	Phone              string
	SubscribedChannels []string
}

// Priority order: critical > info > promotional
var priorityOrder = []string{"promotional", "info", "critical"}

func priorityLevel(p string) int {
	for i, v := range priorityOrder {
		if v == p {
			return i
		}
	}
	return 0
}

// ─── Observer Interface ───────────────────────────────────────────────────────

type NotificationObserver interface {
	Update(event, priority string, user User)
	GetChannel() string
}

// ─── Existing Channels (copy from Part 1) ────────────────────────────────────

type EmailObserver struct{}

func (e *EmailObserver) GetChannel() string { return "email" }
func (e *EmailObserver) Update(event, priority string, user User) {
	// TODO: implement
	_ = fmt.Sprintf
}

type SMSObserver struct{}

func (s *SMSObserver) GetChannel() string { return "sms" }
func (s *SMSObserver) Update(event, priority string, user User) {
	// TODO: implement
	_ = fmt.Sprintf
}

type PushObserver struct{}

func (p *PushObserver) GetChannel() string { return "push" }
func (p *PushObserver) Update(event, priority string, user User) {
	// TODO: implement
	_ = fmt.Sprintf
}

// ─── NEW: PriorityFilteredObserver Decorator ─────────────────────────────────
// HINT: Wraps any observer and filters by priority.
// If event priority < user's minPriority for this channel, skip the notification.

type PriorityFilteredObserver struct {
	inner       NotificationObserver
	minPriority string
}

func NewPriorityFilteredObserver(obs NotificationObserver, minP string) *PriorityFilteredObserver {
	return &PriorityFilteredObserver{inner: obs, minPriority: minP}
}

func (f *PriorityFilteredObserver) GetChannel() string { return f.inner.GetChannel() }
func (f *PriorityFilteredObserver) Update(event, priority string, user User) {
	// TODO: Only call f.inner.Update() if priorityLevel(priority) >= priorityLevel(f.minPriority)
}

// ─── Notification Manager ─────────────────────────────────────────────────────

type NotificationManager struct {
	observers []NotificationObserver
}

func (m *NotificationManager) Subscribe(obs NotificationObserver) {
	m.observers = append(m.observers, obs)
}

func (m *NotificationManager) NotifyAll(event, priority string, users []User) {
	// TODO: for each user, for each observer matching user's channel, call Update()
}

// ─── Test Entry Point ────────────────────────────────────────────────────────

func Notify(event, priority string, users []User, userMinPriority map[string]string) {
	// TODO: Build manager with PriorityFilteredObserver wrappers
}

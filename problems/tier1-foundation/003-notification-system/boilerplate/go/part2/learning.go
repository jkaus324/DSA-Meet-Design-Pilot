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

// ─── Concrete Observers ───────────────────────────────────────────────────────

type EmailObserver struct{}

func (e *EmailObserver) GetChannel() string { return "email" }
func (e *EmailObserver) Update(event, priority string, user User) {
	fmt.Printf("[EMAIL] %s: %s [%s]\n", user.Email, event, priority)
}

type SMSObserver struct{}

func (s *SMSObserver) GetChannel() string { return "sms" }
func (s *SMSObserver) Update(event, priority string, user User) {
	fmt.Printf("[SMS] %s: %s [%s]\n", user.Phone, event, priority)
}

type PushObserver struct{}

func (p *PushObserver) GetChannel() string { return "push" }
func (p *PushObserver) Update(event, priority string, user User) {
	fmt.Printf("[PUSH] %s: %s [%s]\n", user.ID, event, priority)
}

// ─── PriorityFilteredObserver ─────────────────────────────────────────────────

type PriorityFilteredObserver struct {
	inner       NotificationObserver
	minPriority string
}

func NewPriorityFilteredObserver(obs NotificationObserver, minP string) *PriorityFilteredObserver {
	return &PriorityFilteredObserver{inner: obs, minPriority: minP}
}

func (f *PriorityFilteredObserver) GetChannel() string { return f.inner.GetChannel() }
func (f *PriorityFilteredObserver) Update(event, priority string, user User) {
	if priorityLevel(priority) >= priorityLevel(f.minPriority) {
		f.inner.Update(event, priority, user)
	}
}

// ─── Notification Manager ─────────────────────────────────────────────────────

type NotificationManager struct {
	observers []NotificationObserver
}

func (m *NotificationManager) Subscribe(obs NotificationObserver) {
	m.observers = append(m.observers, obs)
}

func (m *NotificationManager) NotifyAll(event, priority string, users []User) {
	for _, user := range users {
		for _, obs := range m.observers {
			for _, ch := range user.SubscribedChannels {
				if ch == obs.GetChannel() {
					obs.Update(event, priority, user)
					break
				}
			}
		}
	}
}

// ─── Test Entry Point ────────────────────────────────────────────────────────

func Notify(event, priority string, users []User, userMinPriority map[string]string) {
	// For simplicity, apply same minPriority to all users
	minP := "promotional"
	if v, ok := userMinPriority["*"]; ok {
		minP = v
	}
	email := &EmailObserver{}
	sms := &SMSObserver{}
	push := &PushObserver{}
	fe := NewPriorityFilteredObserver(email, minP)
	fs := NewPriorityFilteredObserver(sms, minP)
	fp := NewPriorityFilteredObserver(push, minP)
	mgr := &NotificationManager{}
	mgr.Subscribe(fe)
	mgr.Subscribe(fs)
	mgr.Subscribe(fp)
	mgr.NotifyAll(event, priority, users)
}

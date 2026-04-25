package main

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
	// TODO: Print formatted email notification using user.Email, event, priority
}

type SMSObserver struct{}

func (s *SMSObserver) GetChannel() string { return "sms" }
func (s *SMSObserver) Update(event, priority string, user User) {
	// TODO: Print formatted SMS notification using user.Phone, event, priority
}

type PushObserver struct{}

func (p *PushObserver) GetChannel() string { return "push" }
func (p *PushObserver) Update(event, priority string, user User) {
	// TODO: Print formatted push notification using user.ID, event, priority
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
	// TODO: Forward to f.inner only if priorityLevel(priority) >= priorityLevel(f.minPriority)
}

// ─── Notification Manager ─────────────────────────────────────────────────────

type NotificationManager struct {
	observers []NotificationObserver
}

func (m *NotificationManager) Subscribe(obs NotificationObserver) {
	// TODO: Append obs to m.observers
}

func (m *NotificationManager) NotifyAll(event, priority string, users []User) {
	// TODO: For each user, for each observer, call Update if user.SubscribedChannels contains obs.GetChannel()
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

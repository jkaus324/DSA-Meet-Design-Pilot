package main

import "fmt"

// ─── Data Model (given — do not modify) ─────────────────────────────────────

type User struct {
	ID                 string
	Email              string
	Phone              string
	SubscribedChannels []string
}

// ─── Observer Interface ───────────────────────────────────────────────────────

type NotificationObserver interface {
	Send(userID, message string)
	ChannelName() string
}

// ─── Concrete Observers ───────────────────────────────────────────────────────

type EmailNotifier struct{}

func (e *EmailNotifier) ChannelName() string { return "email" }
func (e *EmailNotifier) Send(userID, message string) {
	// TODO: Print a formatted email notification
	// e.g., "[EMAIL] To: <userID> — <message>"
	_ = fmt.Sprintf
}

type SMSNotifier struct{}

func (s *SMSNotifier) ChannelName() string { return "sms" }
func (s *SMSNotifier) Send(userID, message string) {
	// TODO: Print a formatted SMS notification
	// e.g., "[SMS] To: <userID> — <message>"
	_ = fmt.Sprintf
}

type PushNotifier struct{}

func (p *PushNotifier) ChannelName() string { return "push" }
func (p *PushNotifier) Send(userID, message string) {
	// TODO: Print a formatted push notification
	// e.g., "[PUSH] To: <userID> — <message>"
	_ = fmt.Sprintf
}

// ─── Notification Manager ────────────────────────────────────────────────────

type NotificationManager struct {
	observers []NotificationObserver
}

func (m *NotificationManager) Subscribe(obs NotificationObserver) {
	// TODO: Add observer to the list
}

func (m *NotificationManager) Unsubscribe(channel string) {
	// TODO: Remove observer matching ChannelName()
}

func (m *NotificationManager) Notify(event string, users []User) {
	// TODO: For each user, for each of their SubscribedChannels,
	//       find the matching observer and call Send()
}

// ─── Test Entry Point ────────────────────────────────────────────────────────

func Notify(event string, users []User) {
	mgr := &NotificationManager{}
	mgr.Subscribe(&EmailNotifier{})
	mgr.Subscribe(&SMSNotifier{})
	mgr.Subscribe(&PushNotifier{})
	mgr.Notify(event, users)
}

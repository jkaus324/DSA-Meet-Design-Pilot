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
// HINT: This interface represents something that "watches" for events.
// What method does it need to receive a notification?

type NotificationObserver interface {
	// HINT: method to deliver a message to a user
	Send(userID, message string)
	ChannelName() string
}

// ─── Concrete Observers ───────────────────────────────────────────────────────
// TODO: Implement one observer per notification channel:
//   - EmailNotifier
//   - SMSNotifier
//   - PushNotifier

// ─── Subject / Notification Manager ─────────────────────────────────────────
// TODO: Implement a NotificationManager that:
//   - Allows observers to subscribe/unsubscribe
//   - Notifies all subscribers when an event occurs
//   - Does NOT know which specific channels are used

type NotificationManager struct {
	// HINT: store a slice of NotificationObserver here
}

// TODO: define the parameter type — what should Subscribe accept?
func (m *NotificationManager) Subscribe(obs NotificationObserver) {
	// TODO: add observer to the list
}

func (m *NotificationManager) Unsubscribe(channel string) {
	// TODO: remove observer matching ChannelName()
}

func (m *NotificationManager) Notify(event string, users []User) {
	// TODO: For each user, for each of their SubscribedChannels,
	//       find the matching observer and call Send()
	_ = fmt.Sprintf // ensure fmt is used
}

// ─── Test Entry Point ────────────────────────────────────────────────────────
// func Notify(event string, users []User)
// ─────────────────────────────────────────────────────────────────────────────

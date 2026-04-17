package main

// ─── Data Model (given — do not modify) ─────────────────────────────────────

type Notification struct {
	UserID  string
	Message string
	Channel string
}

type User struct {
	ID                 string
	Email              string
	Phone              string
	SubscribedChannels []string
}

// ─── NEW in Extension 1 ──────────────────────────────────────────────────────
//
// The product team wants notification PRIORITIES and FILTERING:
//   - Each event now has a priority: "critical", "info", "promotional"
//   - Users can set minimum priority per channel
//     (e.g., only receive SMS for "critical" events, not "promotional")
//   - The system must respect these per-user, per-channel preferences
//
// Think about:
//   - Where does priority filtering belong in your Observer design?
//   - Is filtering a responsibility of the subject, the observer, or a decorator?
//   - How do you store per-user preferences without coupling User to channels?
//
// Entry point (must exist for tests):
//   func Notify(event, priority string, users []User, userMinPriority map[string]string)
//   // userMinPriority: userID -> min priority for their channel
//
// ─────────────────────────────────────────────────────────────────────────────

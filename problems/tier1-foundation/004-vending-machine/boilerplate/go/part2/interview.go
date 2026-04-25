package main

// ─── Data Model (given — do not modify) ─────────────────────────────────────

type Item struct {
	Name     string
	Price    float64
	Quantity int
}

// ─── NEW in Extension 1 ──────────────────────────────────────────────────────
//
// The vending machine now needs a MAINTENANCE mode:
//   - Operator can switch the machine into maintenance mode
//   - In maintenance mode: restock items, adjust prices, clear errors
//   - User-facing operations (select, pay, dispense) are blocked during maintenance
//   - Only the operator can exit maintenance mode
//
// Think about:
//   - Where does "maintenance" fit in your existing state diagram?
//   - Is it a state like Idle/PaymentPending, or a separate mode overlay?
//   - How do you prevent users from entering maintenance mode?
//
// Entry points (all from Part 1, plus):
//   func SelectItem(itemName string)
//   func InsertMoney(amount float64)
//   func Dispense()
//   func Cancel()
//   func GetState() string
//   func Reset()
//   func EnterMaintenance(operatorPin string)
//   func ExitMaintenance(operatorPin string)
//   func Restock(itemName string, quantity int)
//
// ─────────────────────────────────────────────────────────────────────────────

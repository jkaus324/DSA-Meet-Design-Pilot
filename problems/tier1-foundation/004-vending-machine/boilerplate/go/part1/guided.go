package main

// ─── Data Model (given — do not modify) ─────────────────────────────────────

type Item struct {
	Name     string
	Price    float64
	Quantity int
}

// ─── State Interface ─────────────────────────────────────────────────────────
// HINT: Each state handles its own version of each user action.
// If the action is invalid in this state, it prints an error message.

type VMState interface {
	SelectItem(vm *VendingMachine, itemName string)
	InsertMoney(vm *VendingMachine, amount float64)
	Dispense(vm *VendingMachine)
	Cancel(vm *VendingMachine)
	Name() string
}

// ─── Concrete States ─────────────────────────────────────────────────────────
// TODO: Implement each state:
//   - IdleState       — waiting for item selection
//   - SelectedState   — item chosen, waiting for payment
//   - PaidState       — payment received, ready to dispense
//   - DispensingState — currently dispensing item

// ─── Vending Machine Context ─────────────────────────────────────────────────
// TODO: Implement the VendingMachine struct that:
//   - Holds the current state
//   - Delegates all actions to the current state
//   - Has SetState() to switch between states

type VendingMachine struct {
	// HINT: store current VMState, selectedItem, insertedAmount, inventory
}

func (vm *VendingMachine) SetState(s VMState)          { /* TODO */ }
func (vm *VendingMachine) SelectItem(item string)      { /* TODO: delegate to state */ }
func (vm *VendingMachine) InsertMoney(amount float64)  { /* TODO: delegate to state */ }
func (vm *VendingMachine) Dispense()                   { /* TODO: delegate to state */ }
func (vm *VendingMachine) Cancel()                     { /* TODO: delegate to state */ }
func (vm *VendingMachine) GetState() string            { return "" /* TODO */ }

// ─────────────────────────────────────────────────────────────────────────────

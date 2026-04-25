package main

import "fmt"

// ─── Data Model (given — do not modify) ─────────────────────────────────────

type Item struct {
	Name     string
	Price    float64
	Quantity int
}

// ─── State Interface ─────────────────────────────────────────────────────────

type VMState interface {
	SelectItem(vm *VendingMachine, itemName string)
	InsertMoney(vm *VendingMachine, amount float64)
	Dispense(vm *VendingMachine)
	Cancel(vm *VendingMachine)
	Name() string
}

// ─── Vending Machine Context ─────────────────────────────────────────────────

type VendingMachine struct {
	state          VMState
	selectedItem   string
	insertedAmount float64
	inventory      map[string]Item
}

func NewVendingMachine() *VendingMachine {
	vm := &VendingMachine{
		inventory: map[string]Item{
			"Cola":  {"Cola", 25.0, 5},
			"Chips": {"Chips", 15.0, 3},
		},
	}
	vm.state = &IdleState{}
	return vm
}

func (vm *VendingMachine) SetState(s VMState)          { vm.state = s }
func (vm *VendingMachine) SelectItem(item string)      { vm.state.SelectItem(vm, item) }
func (vm *VendingMachine) InsertMoney(amount float64)  { vm.state.InsertMoney(vm, amount) }
func (vm *VendingMachine) Dispense()                   { vm.state.Dispense(vm) }
func (vm *VendingMachine) Cancel()                     { vm.state.Cancel(vm) }
func (vm *VendingMachine) GetState() string            { return vm.state.Name() }

// ─── Concrete States ─────────────────────────────────────────────────────────

type IdleState struct{}

func (s *IdleState) Name() string { return "Idle" }
func (s *IdleState) SelectItem(vm *VendingMachine, item string) {
	// TODO: Check if item exists in inventory and has quantity > 0
	//       If yes: set vm.selectedItem, transition to SelectedState
	//       If no:  print "Item not available"
	_ = fmt.Sprintf
}
func (s *IdleState) InsertMoney(vm *VendingMachine, amount float64) {
	fmt.Println("[Error] Select an item first.")
}
func (s *IdleState) Dispense(vm *VendingMachine) {
	fmt.Println("[Error] No item selected.")
}
func (s *IdleState) Cancel(vm *VendingMachine) {
	fmt.Println("[Info] Nothing to cancel.")
}

type SelectedState struct{}

func (s *SelectedState) Name() string { return "ItemSelected" }
func (s *SelectedState) SelectItem(vm *VendingMachine, item string) {
	fmt.Println("[Info] Item already selected. Cancel first.")
}
func (s *SelectedState) InsertMoney(vm *VendingMachine, amount float64) {
	// TODO: Add amount to vm.insertedAmount
	//       If insertedAmount >= item price: transition to PaidState
	//       Else: print how much more is needed
}
func (s *SelectedState) Dispense(vm *VendingMachine) {
	fmt.Println("[Error] Insert payment first.")
}
func (s *SelectedState) Cancel(vm *VendingMachine) {
	// TODO: Reset selectedItem and insertedAmount, go back to IdleState
}

type PaidState struct{}

func (s *PaidState) Name() string { return "PaymentReceived" }
func (s *PaidState) SelectItem(vm *VendingMachine, item string) {
	fmt.Println("[Error] Payment already made. Dispense or cancel.")
}
func (s *PaidState) InsertMoney(vm *VendingMachine, amount float64) {
	fmt.Println("[Error] Payment already received.")
}
func (s *PaidState) Dispense(vm *VendingMachine) {
	// TODO: Dispense the item (decrement quantity, print confirmation)
	//       Return change if overpaid
	//       Transition to IdleState
}
func (s *PaidState) Cancel(vm *VendingMachine) {
	// TODO: Refund insertedAmount, reset state, go to IdleState
}

// ─── Global machine + test entry points ─────────────────────────────────────

var globalVM = NewVendingMachine()

func SelectItem(item string)      { globalVM.SelectItem(item) }
func InsertMoney(amount float64)  { globalVM.InsertMoney(amount) }
func Dispense()                   { globalVM.Dispense() }
func Cancel()                     { globalVM.Cancel() }
func GetState() string            { return globalVM.GetState() }
func Reset()                      { globalVM = NewVendingMachine() }

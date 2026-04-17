package main

import "fmt"

// ─── Data Model (given — do not modify) ─────────────────────────────────────

type Item struct {
	Name     string
	Price    float64
	Quantity int
}

// ─── State Interface ──────────────────────────────────────────────────────────

type VMState interface {
	SelectItem(item string)
	InsertMoney(amount float64)
	Dispense()
	Cancel()
	GetName() string
}

// ─── Copy your Part 1 states here ────────────────────────────────────────────

type IdleState struct{ machine *VendingMachine }

func (s *IdleState) SelectItem(item string) { /* TODO */ }
func (s *IdleState) InsertMoney(_ float64)  { fmt.Println("Select an item first.") }
func (s *IdleState) Dispense()              { fmt.Println("Select an item first.") }
func (s *IdleState) Cancel()               { fmt.Println("Nothing to cancel.") }
func (s *IdleState) GetName() string       { return "Idle" }

// TODO: Add PaymentPendingState, DispensingState from Part 1

// ─── NEW: MaintenanceState ────────────────────────────────────────────────────
// HINT: All user-facing operations should print "Machine in maintenance" and return.
// Only operator operations (Restock, ExitMaintenance) are allowed.

type MaintenanceState struct{ machine *VendingMachine }

func (s *MaintenanceState) SelectItem(_ string) { fmt.Println("Machine in maintenance mode.") }
func (s *MaintenanceState) InsertMoney(_ float64) { fmt.Println("Machine in maintenance mode.") }
func (s *MaintenanceState) Dispense()            { fmt.Println("Machine in maintenance mode.") }
func (s *MaintenanceState) Cancel()              { fmt.Println("Machine in maintenance mode.") }
func (s *MaintenanceState) GetName() string      { return "Maintenance" }

// ─── VendingMachine ───────────────────────────────────────────────────────────

type VendingMachine struct {
	currentState   VMState
	inventory      map[string]Item
	insertedMoney  float64
	selectedItem   string
	operatorPin    string
}

func NewVendingMachine() *VendingMachine {
	vm := &VendingMachine{
		inventory:   map[string]Item{"Cola": {"Cola", 25.0, 5}, "Chips": {"Chips", 15.0, 3}},
		operatorPin: "1234",
	}
	vm.currentState = &IdleState{machine: vm}
	return vm
}

func (vm *VendingMachine) SetState(s VMState)         { vm.currentState = s }
func (vm *VendingMachine) SelectItem(item string)     { vm.currentState.SelectItem(item) }
func (vm *VendingMachine) InsertMoney(amt float64)    { vm.currentState.InsertMoney(amt) }
func (vm *VendingMachine) Dispense()                  { vm.currentState.Dispense() }
func (vm *VendingMachine) Cancel()                    { vm.currentState.Cancel() }
func (vm *VendingMachine) GetState() string           { return vm.currentState.GetName() }

func (vm *VendingMachine) EnterMaintenance(pin string) {
	// TODO: if pin == vm.operatorPin, switch to MaintenanceState
}

func (vm *VendingMachine) ExitMaintenance(pin string) {
	// TODO: if pin == vm.operatorPin and in maintenance, switch to IdleState
}

func (vm *VendingMachine) Restock(itemName string, qty int) {
	// TODO: only works in MaintenanceState
}

// ─── Global machine + test entry points ─────────────────────────────────────

var globalVM = NewVendingMachine()

func SelectItem(item string)              { globalVM.SelectItem(item) }
func InsertMoney(amount float64)          { globalVM.InsertMoney(amount) }
func Dispense()                           { globalVM.Dispense() }
func Cancel()                             { globalVM.Cancel() }
func GetState() string                    { return globalVM.GetState() }
func Reset()                              { globalVM = NewVendingMachine() }
func EnterMaintenance(pin string)         { globalVM.EnterMaintenance(pin) }
func ExitMaintenance(pin string)          { globalVM.ExitMaintenance(pin) }
func Restock(itemName string, qty int)    { globalVM.Restock(itemName, qty) }

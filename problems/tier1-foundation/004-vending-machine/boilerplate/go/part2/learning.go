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

// ─── Vending Machine Context ─────────────────────────────────────────────────

type VendingMachine struct {
	currentState  VMState
	inventory     map[string]Item
	insertedMoney float64
	selectedItem  string
	operatorPin   string
}

func NewVendingMachine() *VendingMachine {
	vm := &VendingMachine{
		inventory:   map[string]Item{"Cola": {"Cola", 25.0, 5}, "Chips": {"Chips", 15.0, 3}},
		operatorPin: "1234",
	}
	vm.currentState = &IdleState{m: vm}
	return vm
}

func (vm *VendingMachine) SetState(s VMState)       { vm.currentState = s }
func (vm *VendingMachine) SelectItem(item string)   { vm.currentState.SelectItem(item) }
func (vm *VendingMachine) InsertMoney(amt float64)  { vm.currentState.InsertMoney(amt) }
func (vm *VendingMachine) Dispense()                { vm.currentState.Dispense() }
func (vm *VendingMachine) Cancel()                  { vm.currentState.Cancel() }
func (vm *VendingMachine) GetState() string         { return vm.currentState.GetName() }

func (vm *VendingMachine) EnterMaintenance(pin string) {
	if pin == vm.operatorPin {
		vm.SetState(&MaintenanceState{m: vm})
		fmt.Println("Entered maintenance mode.")
	} else {
		fmt.Println("Invalid PIN.")
	}
}

func (vm *VendingMachine) ExitMaintenance(pin string) {
	if pin == vm.operatorPin && vm.GetState() == "Maintenance" {
		vm.SetState(&IdleState{m: vm})
		fmt.Println("Exited maintenance mode.")
	} else {
		fmt.Println("Invalid PIN or not in maintenance.")
	}
}

func (vm *VendingMachine) Restock(itemName string, qty int) {
	if vm.GetState() != "Maintenance" {
		fmt.Println("Must be in maintenance to restock.")
		return
	}
	item := vm.inventory[itemName]
	item.Quantity += qty
	vm.inventory[itemName] = item
	fmt.Printf("Restocked %s by %d\n", itemName, qty)
}

// ─── Concrete States ─────────────────────────────────────────────────────────

type IdleState struct{ m *VendingMachine }

func (s *IdleState) GetName() string { return "Idle" }
func (s *IdleState) SelectItem(item string) {
	// TODO: Check inventory, set selectedItem, transition to PaymentPendingState
}
func (s *IdleState) InsertMoney(_ float64) { fmt.Println("Select item first.") }
func (s *IdleState) Dispense()             { fmt.Println("Select item first.") }
func (s *IdleState) Cancel()               {}

type PaymentPendingState struct{ m *VendingMachine }

func (s *PaymentPendingState) GetName() string { return "PaymentPending" }
func (s *PaymentPendingState) SelectItem(_ string) {
	fmt.Println("Already selected.")
}
func (s *PaymentPendingState) InsertMoney(amt float64) {
	// TODO: Add amt, check if enough, transition to DispensingState if so
}
func (s *PaymentPendingState) Dispense() { fmt.Println("Insert money first.") }
func (s *PaymentPendingState) Cancel() {
	// TODO: Refund and go back to IdleState
}

type DispensingState struct{ m *VendingMachine }

func (s *DispensingState) GetName() string         { return "Dispensing" }
func (s *DispensingState) SelectItem(_ string)     { fmt.Println("Dispensing in progress.") }
func (s *DispensingState) InsertMoney(_ float64)   { fmt.Println("Dispensing in progress.") }
func (s *DispensingState) Cancel()                 { fmt.Println("Cannot cancel during dispensing.") }
func (s *DispensingState) Dispense() {
	// TODO: Decrement inventory, print confirmation, go to IdleState
}

type MaintenanceState struct{ m *VendingMachine }

func (s *MaintenanceState) GetName() string        { return "Maintenance" }
func (s *MaintenanceState) SelectItem(_ string)    { fmt.Println("Machine in maintenance.") }
func (s *MaintenanceState) InsertMoney(_ float64)  { fmt.Println("Machine in maintenance.") }
func (s *MaintenanceState) Dispense()              { fmt.Println("Machine in maintenance.") }
func (s *MaintenanceState) Cancel()                { fmt.Println("Machine in maintenance.") }

// ─── Global machine + test entry points ─────────────────────────────────────

var globalVM = NewVendingMachine()

func SelectItem(item string)           { globalVM.SelectItem(item) }
func InsertMoney(amount float64)       { globalVM.InsertMoney(amount) }
func Dispense()                        { globalVM.Dispense() }
func Cancel()                          { globalVM.Cancel() }
func GetState() string                 { return globalVM.GetState() }
func Reset()                           { globalVM = NewVendingMachine() }
func EnterMaintenance(pin string)      { globalVM.EnterMaintenance(pin) }
func ExitMaintenance(pin string)       { globalVM.ExitMaintenance(pin) }
func Restock(itemName string, qty int) { globalVM.Restock(itemName, qty) }

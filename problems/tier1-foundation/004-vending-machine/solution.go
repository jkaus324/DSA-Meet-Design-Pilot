// Vending Machine — State pattern reference solution (Go).
package main

const (
	vmIdle           = "Idle"
	vmPaymentPending = "PaymentPending"
	vmDispensing     = "Dispensing"
	vmMaintenance    = "Maintenance"
)

type vmItem struct {
	name     string
	price    float64
	quantity int
}

type vendingMachine struct {
	state         string
	selectedItem  string
	insertedMoney float64
	inventory     map[string]*vmItem
	operatorPin   string
}

func newVendingMachine() *vendingMachine {
	vm := &vendingMachine{
		state:       vmIdle,
		operatorPin: "1234",
		inventory:   map[string]*vmItem{},
	}
	vm.stockDefaults()
	return vm
}

func (vm *vendingMachine) stockDefaults() {
	vm.inventory["Cola"] = &vmItem{name: "Cola", price: 25.0, quantity: 5}
	vm.inventory["Chips"] = &vmItem{name: "Chips", price: 15.0, quantity: 3}
}

func (vm *vendingMachine) reset() {
	vm.state = vmIdle
	vm.selectedItem = ""
	vm.insertedMoney = 0.0
	vm.inventory = map[string]*vmItem{}
	vm.stockDefaults()
}

func (vm *vendingMachine) selectItem(item string) {
	if vm.state != vmIdle {
		return
	}
	if it, ok := vm.inventory[item]; ok && it.quantity > 0 {
		vm.selectedItem = item
		vm.state = vmPaymentPending
	}
}

func (vm *vendingMachine) insertMoney(amount float64) {
	if vm.state != vmPaymentPending {
		return
	}
	vm.insertedMoney += amount
	if vm.insertedMoney >= vm.inventory[vm.selectedItem].price {
		vm.state = vmDispensing
	}
}

func (vm *vendingMachine) dispense() {
	if vm.state != vmDispensing {
		return
	}
	vm.inventory[vm.selectedItem].quantity--
	vm.insertedMoney = 0.0
	vm.selectedItem = ""
	vm.state = vmIdle
}

func (vm *vendingMachine) cancel() {
	if vm.state != vmPaymentPending {
		return
	}
	vm.insertedMoney = 0.0
	vm.selectedItem = ""
	vm.state = vmIdle
}

func (vm *vendingMachine) enterMaintenance(pin string) {
	if pin == vm.operatorPin && vm.state == vmIdle {
		vm.state = vmMaintenance
	}
}

func (vm *vendingMachine) exitMaintenance(pin string) {
	if pin == vm.operatorPin && vm.state == vmMaintenance {
		vm.state = vmIdle
	}
}

func (vm *vendingMachine) restock(item string, qty int) {
	if vm.state != vmMaintenance {
		return
	}
	if _, ok := vm.inventory[item]; !ok {
		vm.inventory[item] = &vmItem{name: item, price: 0.0, quantity: 0}
	}
	vm.inventory[item].quantity += qty
}

var gVM = newVendingMachine()

func reset_service() {
	gVM = newVendingMachine()
}

func reset() {
	gVM.reset()
}

func getState() string {
	return gVM.state
}

func selectItem(item string) {
	gVM.selectItem(item)
}

func insertMoney(amount float64) {
	gVM.insertMoney(amount)
}

func dispense() {
	gVM.dispense()
}

func cancel() {
	gVM.cancel()
}

func enterMaintenance(pin string) {
	gVM.enterMaintenance(pin)
}

func exitMaintenance(pin string) {
	gVM.exitMaintenance(pin)
}

func restock(item string, qty int) {
	gVM.restock(item, qty)
}

func vm_get_quantity(item string) int {
	if it, ok := gVM.inventory[item]; ok {
		return it.quantity
	}
	return -1
}

"""Vending Machine — State pattern reference solution (Python)."""


class VMState:
    IDLE = "Idle"
    PAYMENT_PENDING = "PaymentPending"
    DISPENSING = "Dispensing"
    MAINTENANCE = "Maintenance"


class VendingMachine:
    def __init__(self):
        self.state = VMState.IDLE
        self.selectedItem = ""
        self.insertedMoney = 0.0
        self.inventory = {}
        self.operatorPin = "1234"
        self._stockDefaults()

    def _stockDefaults(self):
        self.inventory["Cola"] = {"name": "Cola", "price": 25.0, "quantity": 5}
        self.inventory["Chips"] = {"name": "Chips", "price": 15.0, "quantity": 3}

    def reset(self):
        self.state = VMState.IDLE
        self.selectedItem = ""
        self.insertedMoney = 0.0
        self._stockDefaults()

    def selectItem(self, item):
        if self.state != VMState.IDLE:
            return
        if item in self.inventory and self.inventory[item]["quantity"] > 0:
            self.selectedItem = item
            self.state = VMState.PAYMENT_PENDING

    def insertMoney(self, amount):
        if self.state != VMState.PAYMENT_PENDING:
            return
        self.insertedMoney += amount
        if self.insertedMoney >= self.inventory[self.selectedItem]["price"]:
            self.state = VMState.DISPENSING

    def dispense(self):
        if self.state != VMState.DISPENSING:
            return
        self.inventory[self.selectedItem]["quantity"] -= 1
        self.insertedMoney = 0.0
        self.selectedItem = ""
        self.state = VMState.IDLE

    def cancel(self):
        if self.state != VMState.PAYMENT_PENDING:
            return
        self.insertedMoney = 0.0
        self.selectedItem = ""
        self.state = VMState.IDLE

    def enterMaintenance(self, pin):
        if pin == self.operatorPin and self.state == VMState.IDLE:
            self.state = VMState.MAINTENANCE

    def exitMaintenance(self, pin):
        if pin == self.operatorPin and self.state == VMState.MAINTENANCE:
            self.state = VMState.IDLE

    def restock(self, item, qty):
        if self.state != VMState.MAINTENANCE:
            return
        if item not in self.inventory:
            self.inventory[item] = {"name": item, "price": 0.0, "quantity": 0}
        self.inventory[item]["quantity"] += qty


_g_vm = VendingMachine()


def reset_service():
    global _g_vm
    _g_vm = VendingMachine()


def reset():
    _g_vm.reset()


def getState():
    return _g_vm.state


def selectItem(item):
    _g_vm.selectItem(item)


def insertMoney(amount):
    _g_vm.insertMoney(amount)


def dispense():
    _g_vm.dispense()


def cancel():
    _g_vm.cancel()


def enterMaintenance(pin):
    _g_vm.enterMaintenance(pin)


def exitMaintenance(pin):
    _g_vm.exitMaintenance(pin)


def restock(item, qty):
    _g_vm.restock(item, qty)


def vm_get_quantity(item):
    if item not in _g_vm.inventory:
        return -1
    return _g_vm.inventory[item]["quantity"]


def vm_get_inserted_money():
    return _g_vm.insertedMoney


def vm_get_selected_item():
    return _g_vm.selectedItem

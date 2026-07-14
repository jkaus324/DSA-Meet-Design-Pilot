// Vending Machine — State pattern reference solution (JavaScript).

const VMState = {
  IDLE: "Idle",
  PAYMENT_PENDING: "PaymentPending",
  DISPENSING: "Dispensing",
  MAINTENANCE: "Maintenance",
};

class VendingMachine {
  constructor() {
    this.state = VMState.IDLE;
    this.selectedItem = "";
    this.insertedMoney = 0.0;
    this.inventory = new Map();
    this.operatorPin = "1234";
    this._stockDefaults();
  }

  _stockDefaults() {
    this.inventory.set("Cola", { name: "Cola", price: 25.0, quantity: 5 });
    this.inventory.set("Chips", { name: "Chips", price: 15.0, quantity: 3 });
  }

  reset() {
    this.state = VMState.IDLE;
    this.selectedItem = "";
    this.insertedMoney = 0.0;
    this._stockDefaults();
  }

  selectItem(item) {
    if (this.state !== VMState.IDLE) return;
    if (this.inventory.has(item) && this.inventory.get(item).quantity > 0) {
      this.selectedItem = item;
      this.state = VMState.PAYMENT_PENDING;
    }
  }

  insertMoney(amount) {
    if (this.state !== VMState.PAYMENT_PENDING) return;
    this.insertedMoney += amount;
    if (this.insertedMoney >= this.inventory.get(this.selectedItem).price) {
      this.state = VMState.DISPENSING;
    }
  }

  dispense() {
    if (this.state !== VMState.DISPENSING) return;
    this.inventory.get(this.selectedItem).quantity -= 1;
    this.insertedMoney = 0.0;
    this.selectedItem = "";
    this.state = VMState.IDLE;
  }

  cancel() {
    if (this.state !== VMState.PAYMENT_PENDING) return;
    this.insertedMoney = 0.0;
    this.selectedItem = "";
    this.state = VMState.IDLE;
  }

  enterMaintenance(pin) {
    if (pin === this.operatorPin && this.state === VMState.IDLE) {
      this.state = VMState.MAINTENANCE;
    }
  }

  exitMaintenance(pin) {
    if (pin === this.operatorPin && this.state === VMState.MAINTENANCE) {
      this.state = VMState.IDLE;
    }
  }

  restock(item, qty) {
    if (this.state !== VMState.MAINTENANCE) return;
    if (!this.inventory.has(item)) {
      this.inventory.set(item, { name: item, price: 0.0, quantity: 0 });
    }
    this.inventory.get(item).quantity += qty;
  }
}

let _g_vm = new VendingMachine();

function reset_service() {
  _g_vm = new VendingMachine();
}

function reset() {
  _g_vm.reset();
}

function getState() {
  return _g_vm.state;
}

function selectItem(item) {
  _g_vm.selectItem(item);
}

function insertMoney(amount) {
  _g_vm.insertMoney(amount);
}

function dispense() {
  _g_vm.dispense();
}

function cancel() {
  _g_vm.cancel();
}

function enterMaintenance(pin) {
  _g_vm.enterMaintenance(pin);
}

function exitMaintenance(pin) {
  _g_vm.exitMaintenance(pin);
}

function restock(item, qty) {
  _g_vm.restock(item, qty);
}

function vm_get_quantity(item) {
  if (!_g_vm.inventory.has(item)) return -1;
  return _g_vm.inventory.get(item).quantity;
}

function vm_get_inserted_money() {
  return _g_vm.insertedMoney;
}

function vm_get_selected_item() {
  return _g_vm.selectedItem;
}

module.exports = {
  VendingMachine,
  reset_service,
  reset,
  getState,
  selectItem,
  insertMoney,
  dispense,
  cancel,
  enterMaintenance,
  exitMaintenance,
  restock,
  vm_get_quantity,
  vm_get_inserted_money,
  vm_get_selected_item,
};

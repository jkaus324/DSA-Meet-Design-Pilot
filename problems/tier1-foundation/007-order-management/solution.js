// Order Management — State + Observer reference solution (JavaScript).

const OrderState = {
  Created: "Created",
  Confirmed: "Confirmed",
  Shipped: "Shipped",
  Delivered: "Delivered",
  Cancelled: "Cancelled",
};

class OrderItem {
  constructor(productId, quantity) {
    this.productId = productId;
    this.quantity = quantity;
  }
}

class Order {
  constructor(id, items, totalAmount) {
    this.id = id;
    this.items = items;
    this.totalAmount = totalAmount;
    this.state = OrderState.Created;
  }
}

class OrderObserver {
  onStateChange(orderId, fromState, toState) { throw new Error('not implemented'); }
}

class OrderManager {
  constructor() {
    this.orders = new Map();
    this.inventory = new Map();
    this.history = new Map();
    this.observers = [];
    this.nextId = 1;
  }

  _transition(orderId, expected, next_) {
    if (!this.orders.has(orderId)) return false;
    const o = this.orders.get(orderId);
    if (o.state !== expected) return false;
    const old = o.state;
    o.state = next_;
    if (!this.history.has(orderId)) this.history.set(orderId, []);
    this.history.get(orderId).push([old, next_]);
    for (const obs of this.observers) {
      obs.onStateChange(orderId, old, next_);
    }
    return true;
  }

  setInventory(productId, qty) {
    this.inventory.set(productId, qty);
  }

  getInventory(productId) {
    return this.inventory.has(productId) ? this.inventory.get(productId) : 0;
  }

  createOrder(items, totalAmount) {
    const oid = `ORD-${this.nextId}`;
    this.nextId += 1;
    for (const it of items) {
      const cur = this.inventory.has(it.productId) ? this.inventory.get(it.productId) : 0;
      this.inventory.set(it.productId, cur - it.quantity);
    }
    this.orders.set(oid, new Order(oid, [...items], totalAmount));
    this.history.set(oid, [[OrderState.Created, OrderState.Created]]);
    return oid;
  }

  confirmOrder(oid) {
    return this._transition(oid, OrderState.Created, OrderState.Confirmed);
  }

  shipOrder(oid) {
    return this._transition(oid, OrderState.Confirmed, OrderState.Shipped);
  }

  deliverOrder(oid) {
    return this._transition(oid, OrderState.Shipped, OrderState.Delivered);
  }

  cancelOrder(oid) {
    if (!this.orders.has(oid)) return false;
    const o = this.orders.get(oid);
    if (o.state !== OrderState.Created && o.state !== OrderState.Confirmed) return false;
    for (const it of o.items) {
      const cur = this.inventory.has(it.productId) ? this.inventory.get(it.productId) : 0;
      this.inventory.set(it.productId, cur + it.quantity);
    }
    const old = o.state;
    o.state = OrderState.Cancelled;
    this.history.get(oid).push([old, OrderState.Cancelled]);
    for (const obs of this.observers) {
      obs.onStateChange(oid, old, OrderState.Cancelled);
    }
    return true;
  }

  getOrderState(oid) {
    return this.orders.has(oid) ? this.orders.get(oid).state : "";
  }

  getOrderHistory(oid) {
    return this.history.has(oid) ? this.history.get(oid) : [];
  }

  addObserver(obs) {
    this.observers.push(obs);
  }
}

let _g_mgr = new OrderManager();
let _g_obs = null;

function reset_service() {
  _g_mgr = new OrderManager();
  _g_obs = null;
}

function reset_manager() {
  reset_service();
}

function set_inventory(productId, qty) {
  _g_mgr.setInventory(productId, qty);
}

function get_inventory(productId) {
  return _g_mgr.getInventory(productId);
}

function create_order_simple(productId, quantity, totalAmount) {
  return _g_mgr.createOrder([new OrderItem(productId, quantity)], totalAmount);
}

function confirm_order(oid) {
  return _g_mgr.confirmOrder(oid);
}

function ship_order(oid) {
  return _g_mgr.shipOrder(oid);
}

function deliver_order(oid) {
  return _g_mgr.deliverOrder(oid);
}

function cancel_order(oid) {
  return _g_mgr.cancelOrder(oid);
}

function get_order_state_str(oid) {
  return _g_mgr.getOrderState(oid);
}

function get_history_size(oid) {
  return _g_mgr.getOrderHistory(oid).length;
}

class _CountingObs extends OrderObserver {
  constructor() {
    super();
    this.count = 0;
    this.lastOrderId = "";
    this.lastFrom = "";
    this.lastTo = "";
  }
  onStateChange(oid, fromState, toState) {
    this.count += 1;
    this.lastOrderId = oid;
    this.lastFrom = fromState;
    this.lastTo = toState;
  }
}

function om_attach_observer() {
  _g_obs = new _CountingObs();
  _g_mgr.addObserver(_g_obs);
}

function om_observer_count() {
  return _g_obs ? _g_obs.count : 0;
}

function om_observer_last_to() {
  return _g_obs ? _g_obs.lastTo : "";
}

module.exports = {
  OrderItem,
  Order,
  OrderObserver,
  OrderManager,
  reset_service,
  reset_manager,
  set_inventory,
  get_inventory,
  create_order_simple,
  confirm_order,
  ship_order,
  deliver_order,
  cancel_order,
  get_order_state_str,
  get_history_size,
  om_attach_observer,
  om_observer_count,
  om_observer_last_to,
};

"""Order Management — State + Observer reference solution (Python)."""

from abc import ABC, abstractmethod


class OrderState:
    Created = "Created"
    Confirmed = "Confirmed"
    Shipped = "Shipped"
    Delivered = "Delivered"
    Cancelled = "Cancelled"


class OrderItem:
    def __init__(self, productId, quantity):
        self.productId = productId
        self.quantity = quantity


class Order:
    def __init__(self, id, items, totalAmount):
        self.id = id
        self.items = items
        self.totalAmount = totalAmount
        self.state = OrderState.Created


class OrderObserver(ABC):
    @abstractmethod
    def onStateChange(self, orderId, fromState, toState):
        ...


class OrderManager:
    def __init__(self):
        self.orders = {}
        self.inventory = {}
        self.history = {}
        self.observers = []
        self.nextId = 1

    def _transition(self, orderId, expected, next_):
        if orderId not in self.orders:
            return False
        o = self.orders[orderId]
        if o.state != expected:
            return False
        old = o.state
        o.state = next_
        self.history.setdefault(orderId, []).append((old, next_))
        for obs in self.observers:
            obs.onStateChange(orderId, old, next_)
        return True

    def setInventory(self, productId, qty):
        self.inventory[productId] = qty

    def getInventory(self, productId):
        return self.inventory.get(productId, 0)

    def createOrder(self, items, totalAmount):
        oid = f"ORD-{self.nextId}"
        self.nextId += 1
        for it in items:
            self.inventory[it.productId] = self.inventory.get(it.productId, 0) - it.quantity
        self.orders[oid] = Order(oid, list(items), totalAmount)
        self.history[oid] = [(OrderState.Created, OrderState.Created)]
        return oid

    def confirmOrder(self, oid):
        return self._transition(oid, OrderState.Created, OrderState.Confirmed)

    def shipOrder(self, oid):
        return self._transition(oid, OrderState.Confirmed, OrderState.Shipped)

    def deliverOrder(self, oid):
        return self._transition(oid, OrderState.Shipped, OrderState.Delivered)

    def cancelOrder(self, oid):
        if oid not in self.orders:
            return False
        o = self.orders[oid]
        if o.state not in (OrderState.Created, OrderState.Confirmed):
            return False
        for it in o.items:
            self.inventory[it.productId] = self.inventory.get(it.productId, 0) + it.quantity
        old = o.state
        o.state = OrderState.Cancelled
        self.history[oid].append((old, OrderState.Cancelled))
        for obs in self.observers:
            obs.onStateChange(oid, old, OrderState.Cancelled)
        return True

    def getOrderState(self, oid):
        return self.orders[oid].state if oid in self.orders else ""

    def getOrderHistory(self, oid):
        return self.history.get(oid, [])

    def addObserver(self, obs):
        self.observers.append(obs)


_g_mgr = OrderManager()
_g_obs = None


def reset_service():
    global _g_mgr, _g_obs
    _g_mgr = OrderManager()
    _g_obs = None


def reset_manager():
    reset_service()


def set_inventory(productId, qty):
    _g_mgr.setInventory(productId, qty)


def get_inventory(productId):
    return _g_mgr.getInventory(productId)


def create_order_simple(productId, quantity, totalAmount):
    return _g_mgr.createOrder([OrderItem(productId, quantity)], totalAmount)


def confirm_order(oid):
    return _g_mgr.confirmOrder(oid)


def ship_order(oid):
    return _g_mgr.shipOrder(oid)


def deliver_order(oid):
    return _g_mgr.deliverOrder(oid)


def cancel_order(oid):
    return _g_mgr.cancelOrder(oid)


def get_order_state_str(oid):
    return _g_mgr.getOrderState(oid)


def get_history_size(oid):
    return len(_g_mgr.getOrderHistory(oid))


class _CountingObs(OrderObserver):
    def __init__(self):
        self.count = 0
        self.lastOrderId = ""
        self.lastFrom = ""
        self.lastTo = ""

    def onStateChange(self, oid, fromState, toState):
        self.count += 1
        self.lastOrderId = oid
        self.lastFrom = fromState
        self.lastTo = toState


def om_attach_observer():
    global _g_obs
    _g_obs = _CountingObs()
    _g_mgr.addObserver(_g_obs)


def om_observer_count():
    return _g_obs.count if _g_obs else 0


def om_observer_last_to():
    return _g_obs.lastTo if _g_obs else ""

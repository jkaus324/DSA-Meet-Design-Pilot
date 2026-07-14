// Order Management — State + Observer reference solution (Go).
package main

const (
	osCreated   = "Created"
	osConfirmed = "Confirmed"
	osShipped   = "Shipped"
	osDelivered = "Delivered"
	osCancelled = "Cancelled"
)

type orderItem struct {
	productID string
	quantity  int
}

type order struct {
	id          string
	items       []orderItem
	totalAmount float64
	state       string
}

type stateChange struct {
	from string
	to   string
}

type orderObserver interface {
	onStateChange(orderID, fromState, toState string)
}

type orderManager struct {
	orders    map[string]*order
	inventory map[string]int
	history   map[string][]stateChange
	observers []orderObserver
	nextID    int
}

func newOrderManager() *orderManager {
	return &orderManager{
		orders:    map[string]*order{},
		inventory: map[string]int{},
		history:   map[string][]stateChange{},
		nextID:    1,
	}
}

func (m *orderManager) transition(orderID, expected, next string) bool {
	o, ok := m.orders[orderID]
	if !ok {
		return false
	}
	if o.state != expected {
		return false
	}
	old := o.state
	o.state = next
	m.history[orderID] = append(m.history[orderID], stateChange{old, next})
	for _, obs := range m.observers {
		obs.onStateChange(orderID, old, next)
	}
	return true
}

func (m *orderManager) setInventory(productID string, qty int) {
	m.inventory[productID] = qty
}

func (m *orderManager) getInventory(productID string) int {
	return m.inventory[productID]
}

func (m *orderManager) createOrder(items []orderItem, totalAmount float64) string {
	oid := "ORD-" + itoa(m.nextID)
	m.nextID++
	for _, it := range items {
		m.inventory[it.productID] = m.inventory[it.productID] - it.quantity
	}
	m.orders[oid] = &order{id: oid, items: items, totalAmount: totalAmount, state: osCreated}
	m.history[oid] = []stateChange{{osCreated, osCreated}}
	return oid
}

func (m *orderManager) cancelOrder(oid string) bool {
	o, ok := m.orders[oid]
	if !ok {
		return false
	}
	if o.state != osCreated && o.state != osConfirmed {
		return false
	}
	for _, it := range o.items {
		m.inventory[it.productID] = m.inventory[it.productID] + it.quantity
	}
	old := o.state
	o.state = osCancelled
	m.history[oid] = append(m.history[oid], stateChange{old, osCancelled})
	for _, obs := range m.observers {
		obs.onStateChange(oid, old, osCancelled)
	}
	return true
}

func (m *orderManager) getOrderState(oid string) string {
	if o, ok := m.orders[oid]; ok {
		return o.state
	}
	return ""
}

func (m *orderManager) getHistorySize(oid string) int {
	return len(m.history[oid])
}

func (m *orderManager) addObserver(obs orderObserver) {
	m.observers = append(m.observers, obs)
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	buf := []byte{}
	for n > 0 {
		buf = append([]byte{byte('0' + n%10)}, buf...)
		n /= 10
	}
	if neg {
		return "-" + string(buf)
	}
	return string(buf)
}

type countingObs struct {
	count       int
	lastOrderID string
	lastFrom    string
	lastTo      string
}

func (c *countingObs) onStateChange(oid, fromState, toState string) {
	c.count++
	c.lastOrderID = oid
	c.lastFrom = fromState
	c.lastTo = toState
}

var (
	gMgr = newOrderManager()
	gObs *countingObs
)

func reset_service() {
	gMgr = newOrderManager()
	gObs = nil
}

func set_inventory(productID string, qty int) {
	gMgr.setInventory(productID, qty)
}

func get_inventory(productID string) int {
	return gMgr.getInventory(productID)
}

func create_order_simple(productID string, quantity int, totalAmount float64) string {
	return gMgr.createOrder([]orderItem{{productID, quantity}}, totalAmount)
}

func confirm_order(oid string) bool {
	return gMgr.transition(oid, osCreated, osConfirmed)
}

func ship_order(oid string) bool {
	return gMgr.transition(oid, osConfirmed, osShipped)
}

func deliver_order(oid string) bool {
	return gMgr.transition(oid, osShipped, osDelivered)
}

func cancel_order(oid string) bool {
	return gMgr.cancelOrder(oid)
}

func get_order_state_str(oid string) string {
	return gMgr.getOrderState(oid)
}

func get_history_size(oid string) int {
	return gMgr.getHistorySize(oid)
}

func om_attach_observer() {
	gObs = &countingObs{}
	gMgr.addObserver(gObs)
}

func om_observer_count() int {
	if gObs != nil {
		return gObs.count
	}
	return 0
}

func om_observer_last_to() string {
	if gObs != nil {
		return gObs.lastTo
	}
	return ""
}

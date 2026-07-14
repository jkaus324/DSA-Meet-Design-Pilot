// Parking Lot — multi-floor + spot matching + pricing strategies (Go).
package main

import (
	"fmt"
	"math"
)

type ParkOp struct {
	kind string
	s1   string
	s2   string
	s3   string
	i1   int
	i2   int
	i3   int
}

// Constants
const (
	mType = "MOTORCYCLE"
	cType = "CAR"
	tType = "TRUCK"
)

const (
	sizeSmall  = 0
	sizeMedium = 1
	sizeLarge  = 2
)

type vehicle struct {
	licensePlate string
	vtype        string
}

type parkingSpot struct {
	spotId               string
	floor                int
	size                 int
	isOccupied           bool
	vehicleLicensePlate  string
}

type ticket struct {
	ticketId     string
	licensePlate string
	spotId       string
	floor        int
	entryTime    int
	entryGateId  string
	exitGateId   string
}

type gate struct {
	gateId string
	gtype  string
}

func minSpotSize(vtype string) int {
	if vtype == mType {
		return sizeSmall
	}
	if vtype == cType {
		return sizeMedium
	}
	return sizeLarge
}

func isCompatible(spotSize, minRequired int) bool {
	return spotSize >= minRequired
}

type pricingStrategy interface {
	calculateFee(durationSeconds int) float64
}

type flatRate struct{ fee float64 }

func (f flatRate) calculateFee(durationSeconds int) float64 { return f.fee }

type hourly struct{ rate float64 }

func (h hourly) calculateFee(durationSeconds int) float64 {
	hours := math.Ceil(float64(durationSeconds) / 3600.0)
	return h.rate * hours
}

type tiered struct {
	base float64
	mid  float64
	high float64
}

func (t tiered) calculateFee(durationSeconds int) float64 {
	hours := math.Ceil(float64(durationSeconds) / 3600.0)
	if hours <= 1 {
		return t.base
	}
	if hours <= 3 {
		return t.base + t.mid*(hours-1)
	}
	return t.base + t.mid*2 + t.high*(hours-3)
}

type parkingLot struct {
	floors        [][]*parkingSpot
	activeTickets map[string]*ticket
	gates         []*gate
	strategy      pricingStrategy
	nextTicketId  int
}

func newParkingLot(numFloors int) *parkingLot {
	floors := make([][]*parkingSpot, numFloors)
	for i := range floors {
		floors[i] = []*parkingSpot{}
	}
	return &parkingLot{
		floors:        floors,
		activeTickets: map[string]*ticket{},
		gates:         []*gate{},
		strategy:      nil,
		nextTicketId:  1,
	}
}

func (p *parkingLot) addSpot(floor, size int) {
	if floor < 0 || floor >= len(p.floors) {
		return
	}
	spotId := fmt.Sprintf("F%dS%d", floor, len(p.floors[floor]))
	p.floors[floor] = append(p.floors[floor], &parkingSpot{spotId: spotId, floor: floor, size: size})
}

func (p *parkingLot) setPricingStrategy(s pricingStrategy) { p.strategy = s }

func (p *parkingLot) addGate(gateId, gtype string) {
	p.gates = append(p.gates, &gate{gateId: gateId, gtype: gtype})
}

func (p *parkingLot) getGates(gtype string) []string {
	res := []string{}
	for _, g := range p.gates {
		if g.gtype == gtype {
			res = append(res, g.gateId)
		}
	}
	return res
}

func (p *parkingLot) parkVehicle(v vehicle, entryTime int, gateId string) *ticket {
	minSize := minSpotSize(v.vtype)
	for f, spots := range p.floors {
		for _, spot := range spots {
			if !spot.isOccupied && isCompatible(spot.size, minSize) {
				spot.isOccupied = true
				spot.vehicleLicensePlate = v.licensePlate
				tid := fmt.Sprintf("T%d", p.nextTicketId)
				p.nextTicketId++
				t := &ticket{
					ticketId:     tid,
					licensePlate: v.licensePlate,
					spotId:       spot.spotId,
					floor:        f,
					entryTime:    entryTime,
					entryGateId:  gateId,
					exitGateId:   "",
				}
				p.activeTickets[tid] = t
				return t
			}
		}
	}
	return nil
}

func (p *parkingLot) unparkVehicle(ticketId string, exitTime int, gateId string) float64 {
	t, ok := p.activeTickets[ticketId]
	if !ok {
		return -1.0
	}
	t.exitGateId = gateId
	for _, floorSpots := range p.floors {
		done := false
		for _, spot := range floorSpots {
			if spot.spotId == t.spotId && spot.isOccupied {
				spot.isOccupied = false
				spot.vehicleLicensePlate = ""
				done = true
				break
			}
		}
		if done {
			break
		}
	}
	duration := exitTime - t.entryTime
	var fee float64
	if p.strategy != nil {
		fee = p.strategy.calculateFee(duration)
	} else {
		fee = float64(duration)
	}
	delete(p.activeTickets, ticketId)
	return fee
}

func (p *parkingLot) getAvailableSpots(size int) int {
	count := 0
	for _, floorSpots := range p.floors {
		for _, spot := range floorSpots {
			if !spot.isOccupied && spot.size == size {
				count++
			}
		}
	}
	return count
}

func (p *parkingLot) getAvailableSpotsByFloor(floor, size int) int {
	if floor < 0 || floor >= len(p.floors) {
		return 0
	}
	count := 0
	for _, spot := range p.floors[floor] {
		if !spot.isOccupied && spot.size == size {
			count++
		}
	}
	return count
}

func sizeFrom(s string) int {
	if s == "S" || s == "small" {
		return sizeSmall
	}
	if s == "M" || s == "medium" {
		return sizeMedium
	}
	return sizeLarge
}

func vtypeFrom(s string) string {
	if s == "M" || s == "moto" || s == "motorcycle" {
		return mType
	}
	if s == "C" || s == "car" {
		return cType
	}
	return tType
}

func gateFrom(s string) string {
	if s == "entry" {
		return "ENTRY"
	}
	return "EXIT"
}

func feeToStr(f float64) string {
	if f < 0 {
		return "-1"
	}
	return fmt.Sprintf("%.2f", f)
}

type parkSnap struct {
	id        string
	floor     int
	spotId    string
	entryGate string
}

func emptySnaps() []parkSnap {
	s := make([]parkSnap, 16)
	for i := range s {
		s[i] = parkSnap{id: "", floor: -1, spotId: "", entryGate: ""}
	}
	return s
}

func parking_simulate(ops []ParkOp) []string {
	out := []string{}
	var lot *parkingLot
	tickets := make([]string, 16)
	snaps := emptySnaps()
	for _, op := range ops {
		k := op.kind
		switch k {
		case "new":
			lot = newParkingLot(op.i1)
			tickets = make([]string, 16)
			snaps = emptySnaps()
			out = append(out, "ok")
		case "add_spot":
			lot.addSpot(op.i1, sizeFrom(op.s1))
			out = append(out, "ok")
		case "add_gate":
			lot.addGate(op.s1, gateFrom(op.s2))
			out = append(out, "ok")
		case "gates_count":
			out = append(out, fmt.Sprintf("%d", len(lot.getGates(gateFrom(op.s1)))))
		case "gate_at":
			g := lot.getGates(gateFrom(op.s1))
			if op.i1 >= 0 && op.i1 < len(g) {
				out = append(out, g[op.i1])
			} else {
				out = append(out, "")
			}
		case "set_pricing":
			var p pricingStrategy
			if op.s1 == "flat" {
				p = flatRate{fee: float64(op.i1)}
			} else if op.s1 == "hourly" {
				p = hourly{rate: float64(op.i1)}
			} else if op.s1 == "tiered" {
				p = tiered{base: float64(op.i1), mid: float64(op.i2), high: float64(op.i3)}
			}
			if p != nil {
				lot.setPricingStrategy(p)
			}
			out = append(out, "ok")
		case "park":
			v := vehicle{licensePlate: op.s1, vtype: vtypeFrom(op.s2)}
			t := lot.parkVehicle(v, op.i1, op.s3)
			if op.i2 >= 0 && op.i2 < len(tickets) {
				if t != nil {
					tickets[op.i2] = t.ticketId
					snaps[op.i2] = parkSnap{id: t.ticketId, floor: t.floor, spotId: t.spotId, entryGate: t.entryGateId}
				} else {
					tickets[op.i2] = ""
					snaps[op.i2] = parkSnap{id: "", floor: -1, spotId: "", entryGate: ""}
				}
			}
			if t != nil {
				out = append(out, t.ticketId)
			} else {
				out = append(out, "")
			}
		case "ticket_at":
			if op.i1 >= 0 && op.i1 < len(tickets) {
				out = append(out, tickets[op.i1])
			} else {
				out = append(out, "")
			}
		case "ticket_floor":
			if op.i1 >= 0 && op.i1 < len(snaps) {
				out = append(out, fmt.Sprintf("%d", snaps[op.i1].floor))
			} else {
				out = append(out, "-1")
			}
		case "ticket_spot_id":
			if op.i1 >= 0 && op.i1 < len(snaps) {
				out = append(out, snaps[op.i1].spotId)
			} else {
				out = append(out, "")
			}
		case "ticket_entry":
			if op.i1 >= 0 && op.i1 < len(snaps) {
				out = append(out, snaps[op.i1].entryGate)
			} else {
				out = append(out, "")
			}
		case "unpark":
			tid := ""
			if op.i1 >= 0 && op.i1 < len(tickets) {
				tid = tickets[op.i1]
			}
			fee := lot.unparkVehicle(tid, op.i2, op.s1)
			out = append(out, feeToStr(fee))
		case "unpark_id":
			fee := lot.unparkVehicle(op.s1, op.i2, op.s2)
			out = append(out, feeToStr(fee))
		case "available":
			out = append(out, fmt.Sprintf("%d", lot.getAvailableSpots(sizeFrom(op.s1))))
		case "available_floor":
			out = append(out, fmt.Sprintf("%d", lot.getAvailableSpotsByFloor(op.i1, sizeFrom(op.s1))))
		default:
			out = append(out, "unknown:"+k)
		}
	}
	return out
}

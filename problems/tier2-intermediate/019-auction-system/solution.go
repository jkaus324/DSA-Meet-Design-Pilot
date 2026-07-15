// Auction system — register users, create auctions with strategies, place bids (Go port).
package main

import (
	"math"
	"strconv"
)

type AuctionOp struct {
	kind string
	s1   string
	s2   string
	s3   string
	i1   int
	i2   int
	i3   int
}

const (
	buyer  = "BUYER"
	seller = "SELLER"
	open   = "OPEN"
	closed = "CLOSED"
	noSale = "NO_SALE"
)

type user struct {
	userID  int
	name    string
	utype   string
}

type bid struct {
	bidderID int
	amount   float64
}

type auction struct {
	auctionID int
	sellerID  int
	item      string
	basePrice float64
	status    string
	bids      []bid
}

type auctionStrategy interface {
	acceptBid(a *auction, amount float64) bool
	getVisibleWinningBid(a *auction) float64
	shouldAutoClose(a *auction, amount float64) bool
}

type ascendingStrategy struct{}

func (ascendingStrategy) acceptBid(a *auction, amount float64) bool {
	currentHighest := a.basePrice
	for _, b := range a.bids {
		if b.amount > currentHighest {
			currentHighest = b.amount
		}
	}
	return amount > currentHighest
}

func (ascendingStrategy) getVisibleWinningBid(a *auction) float64 {
	if len(a.bids) == 0 {
		return -1
	}
	best := a.bids[0].amount
	for _, b := range a.bids {
		if b.amount > best {
			best = b.amount
		}
	}
	return best
}

func (ascendingStrategy) shouldAutoClose(a *auction, amount float64) bool { return false }

type sealedBidStrategy struct{}

func (sealedBidStrategy) acceptBid(a *auction, amount float64) bool {
	return amount > a.basePrice
}

func (sealedBidStrategy) getVisibleWinningBid(a *auction) float64 {
	if a.status == open {
		return -1
	}
	if len(a.bids) == 0 {
		return -1
	}
	best := a.bids[0].amount
	for _, b := range a.bids {
		if b.amount > best {
			best = b.amount
		}
	}
	return best
}

func (sealedBidStrategy) shouldAutoClose(a *auction, amount float64) bool { return false }

type buyNowStrategy struct{}

func (buyNowStrategy) acceptBid(a *auction, amount float64) bool {
	return amount >= a.basePrice*1.5
}

func (buyNowStrategy) getVisibleWinningBid(a *auction) float64 {
	if len(a.bids) == 0 {
		return -1
	}
	return a.bids[len(a.bids)-1].amount
}

func (buyNowStrategy) shouldAutoClose(a *auction, amount float64) bool { return true }

func createStrategy(t string) auctionStrategy {
	if t == "SEALED" {
		return sealedBidStrategy{}
	}
	if t == "BUYNOW" {
		return buyNowStrategy{}
	}
	return ascendingStrategy{}
}

type auctionSystem struct {
	nextUserID    int
	nextAuctionID int
	users         map[int]*user
	auctions      map[int]*auction
	strategies    map[int]auctionStrategy
}

func newAuctionSystem() *auctionSystem {
	return &auctionSystem{
		nextUserID:    1,
		nextAuctionID: 1,
		users:         map[int]*user{},
		auctions:      map[int]*auction{},
		strategies:    map[int]auctionStrategy{},
	}
}

func (s *auctionSystem) registerUser(name, t string) int {
	ut := buyer
	if t == "SELLER" {
		ut = seller
	}
	uid := s.nextUserID
	s.nextUserID++
	s.users[uid] = &user{uid, name, ut}
	return uid
}

func (s *auctionSystem) createAuction(sellerID int, item string, basePrice float64, strategyType string) int {
	u, ok := s.users[sellerID]
	if !ok {
		return -1
	}
	if u.utype != seller {
		return -1
	}
	aid := s.nextAuctionID
	s.nextAuctionID++
	s.auctions[aid] = &auction{auctionID: aid, sellerID: sellerID, item: item, basePrice: basePrice, status: open}
	s.strategies[aid] = createStrategy(strategyType)
	return aid
}

func (s *auctionSystem) placeBid(auctionID, buyerID int, amount float64) bool {
	a, ok := s.auctions[auctionID]
	if !ok {
		return false
	}
	u, ok := s.users[buyerID]
	if !ok {
		return false
	}
	if u.utype != buyer {
		return false
	}
	if a.status != open {
		return false
	}
	if buyerID == a.sellerID {
		return false
	}
	strat := s.strategies[auctionID]
	if !strat.acceptBid(a, amount) {
		return false
	}
	a.bids = append(a.bids, bid{buyerID, amount})
	if strat.shouldAutoClose(a, amount) {
		a.status = closed
	}
	return true
}

func (s *auctionSystem) getWinningBid(auctionID int) float64 {
	a, ok := s.auctions[auctionID]
	if !ok {
		return -1
	}
	return s.strategies[auctionID].getVisibleWinningBid(a)
}

func (s *auctionSystem) closeAuction(auctionID int) bool {
	a, ok := s.auctions[auctionID]
	if !ok {
		return false
	}
	if a.status != open {
		return false
	}
	if len(a.bids) == 0 {
		a.status = noSale
	} else {
		a.status = closed
	}
	return true
}

func (s *auctionSystem) getAuctionStatus(auctionID int) string {
	a, ok := s.auctions[auctionID]
	if !ok {
		return "UNKNOWN"
	}
	return a.status
}

func formatWinning(w float64) string {
	if w < 0 {
		return "-1"
	}
	if w == math.Trunc(w) {
		return strconv.Itoa(int(w))
	}
	return strconv.FormatFloat(w, 'f', 2, 64)
}

func auction_simulate(ops []AuctionOp) []string {
	out := []string{}
	sys := newAuctionSystem()
	userSlot := map[int]int{}
	auctionSlot := map[int]int{}
	slotGet := func(m map[int]int, k int) int {
		if v, ok := m[k]; ok {
			return v
		}
		return k
	}
	for _, op := range ops {
		switch op.kind {
		case "new":
			sys = newAuctionSystem()
			userSlot = map[int]int{}
			auctionSlot = map[int]int{}
			out = append(out, "ok")
		case "register":
			uid := sys.registerUser(op.s1, op.s2)
			userSlot[op.i1] = uid
			out = append(out, strconv.Itoa(uid))
		case "create":
			sid := slotGet(userSlot, op.i1)
			strat := op.s3
			if strat == "" {
				strat = "ASCENDING"
			}
			aid := sys.createAuction(sid, op.s2, float64(op.i3), strat)
			auctionSlot[op.i2] = aid
			out = append(out, strconv.Itoa(aid))
		case "bid":
			aid := slotGet(auctionSlot, op.i1)
			bidder := slotGet(userSlot, op.i2)
			ok := sys.placeBid(aid, bidder, float64(op.i3))
			if ok {
				out = append(out, "ok")
			} else {
				out = append(out, "fail")
			}
		case "close":
			aid := slotGet(auctionSlot, op.i1)
			if sys.closeAuction(aid) {
				out = append(out, "ok")
			} else {
				out = append(out, "fail")
			}
		case "winning":
			aid := slotGet(auctionSlot, op.i1)
			out = append(out, formatWinning(sys.getWinningBid(aid)))
		case "status":
			aid := slotGet(auctionSlot, op.i1)
			out = append(out, sys.getAuctionStatus(aid))
		case "user_id_eq":
			uid := slotGet(userSlot, op.i1)
			if uid == op.i2 {
				out = append(out, "yes")
			} else {
				out = append(out, "no")
			}
		default:
			out = append(out, "unknown:"+op.kind)
		}
	}
	return out
}

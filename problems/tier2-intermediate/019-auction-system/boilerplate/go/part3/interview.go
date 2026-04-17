package main

// UserType distinguishes buyers from sellers.
type UserType int

const (
	Buyer  UserType = iota
	Seller UserType = iota
)

// AuctionStatus represents the lifecycle state of an auction.
type AuctionStatus int

const (
	Open   AuctionStatus = iota
	Closed AuctionStatus = iota
	NoSale AuctionStatus = iota
)

// User represents a registered participant.
type User struct {
	UserId int
	Name   string
	Type   UserType
}

// Bid represents a single bid placed on an auction.
type Bid struct {
	BidderId int
	Amount   float64
}

// AuctionSystem — extend with multiple auction strategies:
//
//   ASCENDING (default):
//     - Bids must exceed current highest. Winner = highest bidder at close.
//
//   SEALED:
//     - Any bid above base price is accepted (even below previous bids).
//     - GetWinningBid returns -1 while Open (bids are hidden).
//     - Winner revealed only after close.
//
//   BUYNOW:
//     - First bid >= basePrice * 1.5 instantly wins and auto-closes.
//     - Bids below the buy-now threshold are rejected.
//
// Think about:
//   - What abstraction lets each auction type have its own bid rules?
//   - How do you create the right strategy from a string type name?
//   - How do you add a 4th strategy without modifying the system struct?
//
// Modified entry point:
//   (*AuctionSystem).CreateAuction(sellerId int, item string, basePrice float64,
//                                  strategyType string) int

type AuctionSystem struct {
}

func NewAuctionSystem() *AuctionSystem {
	return &AuctionSystem{}
}

func (a *AuctionSystem) RegisterUser(name, userType string) int {
	return -1
}

func (a *AuctionSystem) CreateAuction(sellerId int, item string, basePrice float64, strategyType string) int {
	return -1
}

func (a *AuctionSystem) PlaceBid(auctionId, buyerId int, amount float64) bool {
	return false
}

func (a *AuctionSystem) GetWinningBid(auctionId int) float64 {
	return -1
}

func (a *AuctionSystem) CloseAuction(auctionId int) bool {
	return false
}

func (a *AuctionSystem) GetAuctionStatus(auctionId int) string {
	return "UNKNOWN"
}

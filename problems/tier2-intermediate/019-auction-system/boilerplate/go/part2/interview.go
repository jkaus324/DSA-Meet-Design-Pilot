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

// AuctionSystem — extend with closing and state management:
//
//   CloseAuction(auctionId):
//     - Close an open auction. If it has bids, status -> Closed.
//     - If no bids, status -> NoSale.
//     - Return true on success, false if already closed or invalid.
//
//   GetAuctionStatus(auctionId):
//     - Return "OPEN", "CLOSED", or "NO_SALE"
//
// Think about:
//   - What state transitions are valid?
//   - What happens when someone bids on a closed auction?
//   - What does GetWinningBid return after a no-sale close?
//
// Entry points (in addition to Part 1):
//   (*AuctionSystem).CloseAuction(auctionId int) bool
//   (*AuctionSystem).GetAuctionStatus(auctionId int) string

type AuctionSystem struct {
}

func NewAuctionSystem() *AuctionSystem {
	return &AuctionSystem{}
}

func (a *AuctionSystem) RegisterUser(name, userType string) int {
	return -1
}

func (a *AuctionSystem) CreateAuction(sellerId int, item string, basePrice float64) int {
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

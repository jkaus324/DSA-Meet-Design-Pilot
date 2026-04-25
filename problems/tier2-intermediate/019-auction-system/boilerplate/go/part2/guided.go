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

// Auction holds all data for a single auction.
type Auction struct {
	AuctionId int
	SellerId  int
	Item      string
	BasePrice float64
	Status    AuctionStatus
	Bids      []Bid
}

// AuctionSystem extends Part 1 with closing and state management.
// HINT: This extends Part 1. Include all Part 1 functionality.
// HINT: State transitions: Open -> Closed (has bids), Open -> NoSale (no bids)
// HINT: Closed and NoSale are terminal states — no further transitions.
type AuctionSystem struct {
	// HINT: Same data structures as Part 1
}

func NewAuctionSystem() *AuctionSystem {
	// HINT: Initialize (same as Part 1)
	return &AuctionSystem{}
}

func (a *AuctionSystem) RegisterUser(name, userType string) int {
	// TODO: Same as Part 1
	return -1
}

func (a *AuctionSystem) CreateAuction(sellerId int, item string, basePrice float64) int {
	// TODO: Same as Part 1
	return -1
}

func (a *AuctionSystem) PlaceBid(auctionId, buyerId int, amount float64) bool {
	// TODO: Same as Part 1
	// HINT: Already checks auction is Open, so closed auctions are handled
	return false
}

func (a *AuctionSystem) GetWinningBid(auctionId int) float64 {
	// TODO: Same as Part 1
	return -1
}

func (a *AuctionSystem) CloseAuction(auctionId int) bool {
	// HINT: Only Open auctions can be closed
	// HINT: If bids exist -> Closed; if no bids -> NoSale
	// HINT: Return false if auction doesn't exist or is already closed
	return false
}

func (a *AuctionSystem) GetAuctionStatus(auctionId int) string {
	// HINT: Map status enum to string: "OPEN", "CLOSED", "NO_SALE"
	return "UNKNOWN"
}

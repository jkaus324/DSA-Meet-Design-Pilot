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

// AuctionStrategy defines the interface each auction type must implement.
// HINT: Each auction type needs different rules for:
//   1. Whether a bid should be accepted
//   2. What the visible winning bid is (sealed hides it while open)
//   3. Whether the auction should auto-close after a bid
type AuctionStrategy interface {
	AcceptBid(auction *Auction, amount float64) bool
	GetVisibleWinningBid(auction *Auction) float64
	ShouldAutoClose(auction *Auction, amount float64) bool
}

// TODO: Implement AscendingStrategy
//   HINT: AcceptBid: amount must exceed current highest (or base price)
//   HINT: GetVisibleWinningBid: return max bid amount, -1 if no bids
//   HINT: ShouldAutoClose: always false

// TODO: Implement SealedBidStrategy
//   HINT: AcceptBid: any amount above base price
//   HINT: GetVisibleWinningBid: return -1 if Open (sealed), max bid if Closed
//   HINT: ShouldAutoClose: always false

// TODO: Implement BuyNowStrategy
//   HINT: AcceptBid: amount must be >= basePrice * 1.5
//   HINT: GetVisibleWinningBid: return the bid amount (only one bid)
//   HINT: ShouldAutoClose: always true (auto-close on successful bid)

// AuctionSystem manages users and auctions with pluggable strategies.
// HINT: Store a strategy per auction (map from auctionId to AuctionStrategy)
// HINT: Delegate bid validation to the strategy
// HINT: Use a factory function to create strategies from string names
type AuctionSystem struct {
	// HINT: Same data as Part 2, plus a map from auctionId to AuctionStrategy
}

func NewAuctionSystem() *AuctionSystem {
	// TODO: Initialize
	return &AuctionSystem{}
}

func (a *AuctionSystem) RegisterUser(name, userType string) int {
	// TODO: Same as before
	return -1
}

func (a *AuctionSystem) CreateAuction(sellerId int, item string, basePrice float64, strategyType string) int {
	// TODO: Same validation as before
	// HINT: Create the appropriate strategy based on strategyType
	// HINT: Store the strategy mapped to this auctionId
	return -1
}

func (a *AuctionSystem) PlaceBid(auctionId, buyerId int, amount float64) bool {
	// TODO: Basic validations (exists, is buyer, auction is open, not seller)
	// HINT: Delegate bid acceptance to the auction's strategy
	// HINT: If strategy says auto-close, set status to Closed
	return false
}

func (a *AuctionSystem) GetWinningBid(auctionId int) float64 {
	// HINT: Delegate to the auction's strategy's GetVisibleWinningBid
	return -1
}

func (a *AuctionSystem) CloseAuction(auctionId int) bool {
	// TODO: Same as Part 2
	return false
}

func (a *AuctionSystem) GetAuctionStatus(auctionId int) string {
	// TODO: Same as Part 2
	return "UNKNOWN"
}

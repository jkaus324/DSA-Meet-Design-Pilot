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
	Open    AuctionStatus = iota
	Closed  AuctionStatus = iota
	NoSale  AuctionStatus = iota
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

// AuctionSystem — design and implement this struct so that:
//   1. Users can be registered as BUYER or SELLER (RegisterUser)
//   2. Sellers can create auctions with a base price (CreateAuction)
//   3. Buyers can place bids that must exceed the current highest (PlaceBid)
//   4. The current highest bid can be queried (GetWinningBid)
//
// Think about:
//   - How do you store users and auctions for fast lookup?
//   - How do you validate that only buyers bid and only sellers create?
//   - How do you track the current highest bid efficiently?
//   - What happens if a seller tries to bid on their own auction?
//
// Entry points (must exist for tests):
//   NewAuctionSystem() *AuctionSystem
//   (*AuctionSystem).RegisterUser(name, userType string) int
//   (*AuctionSystem).CreateAuction(sellerId int, item string, basePrice float64) int
//   (*AuctionSystem).PlaceBid(auctionId, buyerId int, amount float64) bool
//   (*AuctionSystem).GetWinningBid(auctionId int) float64

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

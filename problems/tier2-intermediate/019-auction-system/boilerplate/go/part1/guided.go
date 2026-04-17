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

// AuctionSystem manages users and auctions.
// HINT: Use maps for O(1) user and auction lookup.
// HINT: Auto-increment IDs for users and auctions.
type AuctionSystem struct {
	// HINT: You need counters for auto-assigning IDs
	// HINT: You need maps from ID to User and ID to Auction
}

func NewAuctionSystem() *AuctionSystem {
	// HINT: Initialize ID counters starting at 1
	return &AuctionSystem{}
}

func (a *AuctionSystem) RegisterUser(name, userType string) int {
	// HINT: Parse userType ("BUYER" or "SELLER") into UserType constant
	// HINT: Create User with auto-assigned ID and store in map
	return -1
}

func (a *AuctionSystem) CreateAuction(sellerId int, item string, basePrice float64) int {
	// HINT: Validate that sellerId exists and is a Seller
	// HINT: Return -1 if validation fails
	// HINT: Create auction with status Open and empty bids slice
	return -1
}

func (a *AuctionSystem) PlaceBid(auctionId, buyerId int, amount float64) bool {
	// HINT: Validate auction exists and is Open
	// HINT: Validate buyer exists and is a Buyer
	// HINT: Validate buyer is not the seller of this auction
	// HINT: Find currentHighest = max(basePrice, all existing bid amounts)
	// HINT: New bid must STRICTLY exceed currentHighest
	return false
}

func (a *AuctionSystem) GetWinningBid(auctionId int) float64 {
	// HINT: Iterate through bids to find the maximum
	// HINT: Return -1 if no bids exist
	return -1
}

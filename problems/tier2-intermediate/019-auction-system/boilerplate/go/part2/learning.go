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

// AuctionSystem manages users and auctions with lifecycle management.
type AuctionSystem struct {
	nextUserId    int
	nextAuctionId int
	users         map[int]User
	auctions      map[int]Auction
}

func NewAuctionSystem() *AuctionSystem {
	return &AuctionSystem{
		nextUserId:    1,
		nextAuctionId: 1,
		users:         make(map[int]User),
		auctions:      make(map[int]Auction),
	}
}

func (a *AuctionSystem) RegisterUser(name, userType string) int {
	// TODO: Same as Part 1 — parse type, create User, store, return ID
	return -1
}

func (a *AuctionSystem) CreateAuction(sellerId int, item string, basePrice float64) int {
	// TODO: Same as Part 1 — validate seller, create Auction, store, return ID
	return -1
}

func (a *AuctionSystem) PlaceBid(auctionId, buyerId int, amount float64) bool {
	// TODO: Same as Part 1 — all validations + bid must exceed current highest
	return false
}

func (a *AuctionSystem) GetWinningBid(auctionId int) float64 {
	// TODO: Same as Part 1 — return highest bid amount or -1
	return -1
}

func (a *AuctionSystem) CloseAuction(auctionId int) bool {
	// TODO: Check auction exists in a.auctions
	// TODO: Check auction.Status == Open (only open can close)
	// TODO: If auction.Bids is empty -> set status to NoSale
	// TODO: Else -> set status to Closed
	// TODO: Update the auction back in the map
	// TODO: Return true on success
	return false
}

func (a *AuctionSystem) GetAuctionStatus(auctionId int) string {
	// TODO: Switch on auction.Status
	// TODO: Return "OPEN", "CLOSED", or "NO_SALE"
	return "UNKNOWN"
}

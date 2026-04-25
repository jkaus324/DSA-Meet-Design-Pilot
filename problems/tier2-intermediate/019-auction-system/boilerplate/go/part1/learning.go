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
	// TODO: Parse userType ("BUYER" -> Buyer, "SELLER" -> Seller)
	// TODO: Create User{UserId: a.nextUserId, Name: name, Type: parsedType}
	// TODO: Store in a.users, increment a.nextUserId, return the ID
	return -1
}

func (a *AuctionSystem) CreateAuction(sellerId int, item string, basePrice float64) int {
	// TODO: Check a.users[sellerId] exists and its Type == Seller
	// TODO: Create Auction{AuctionId: a.nextAuctionId, SellerId: sellerId, Item: item,
	//         BasePrice: basePrice, Status: Open, Bids: []Bid{}}
	// TODO: Store in a.auctions, increment a.nextAuctionId, return the ID
	return -1
}

func (a *AuctionSystem) PlaceBid(auctionId, buyerId int, amount float64) bool {
	// TODO: Validate auctionId exists in a.auctions
	// TODO: Validate buyerId exists and is a Buyer
	// TODO: Validate auction.Status == Open
	// TODO: Validate buyerId != auction.SellerId
	// TODO: Find currentHighest = max of basePrice and all existing bid amounts
	// TODO: If amount <= currentHighest, return false
	// TODO: Append Bid{BidderId: buyerId, Amount: amount} to auction.Bids
	// TODO: Update the auction back in the map (Go maps store copies)
	// TODO: Return true
	return false
}

func (a *AuctionSystem) GetWinningBid(auctionId int) float64 {
	// TODO: If auction has no bids, return -1
	// TODO: Find and return the maximum bid amount
	return -1
}

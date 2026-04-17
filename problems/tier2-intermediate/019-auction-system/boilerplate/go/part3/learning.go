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
type AuctionStrategy interface {
	AcceptBid(auction *Auction, amount float64) bool
	GetVisibleWinningBid(auction *Auction) float64
	ShouldAutoClose(auction *Auction, amount float64) bool
}

// AscendingStrategy — bids must exceed current highest.
type AscendingStrategy struct{}

func (s *AscendingStrategy) AcceptBid(auction *Auction, amount float64) bool {
	// TODO: Find currentHighest = max(auction.BasePrice, all bid amounts)
	// TODO: Return amount > currentHighest
	return false
}

func (s *AscendingStrategy) GetVisibleWinningBid(auction *Auction) float64 {
	// TODO: Return the maximum bid amount, or -1 if no bids
	return -1
}

func (s *AscendingStrategy) ShouldAutoClose(auction *Auction, amount float64) bool {
	// TODO: Ascending never auto-closes
	return false
}

// SealedBidStrategy — any bid above base price accepted; winner hidden while open.
type SealedBidStrategy struct{}

func (s *SealedBidStrategy) AcceptBid(auction *Auction, amount float64) bool {
	// TODO: Return amount > auction.BasePrice (any bid above base is valid)
	return false
}

func (s *SealedBidStrategy) GetVisibleWinningBid(auction *Auction) float64 {
	// TODO: If auction.Status == Open, return -1 (sealed — bids are hidden)
	// TODO: If closed, return the maximum bid amount, or -1 if no bids
	return -1
}

func (s *SealedBidStrategy) ShouldAutoClose(auction *Auction, amount float64) bool {
	// TODO: Sealed bids never auto-close
	return false
}

// BuyNowStrategy — first bid >= basePrice * 1.5 instantly wins and auto-closes.
type BuyNowStrategy struct{}

func (s *BuyNowStrategy) AcceptBid(auction *Auction, amount float64) bool {
	// TODO: Return amount >= auction.BasePrice * 1.5
	return false
}

func (s *BuyNowStrategy) GetVisibleWinningBid(auction *Auction) float64 {
	// TODO: Return the bid amount if bids exist, -1 otherwise
	return -1
}

func (s *BuyNowStrategy) ShouldAutoClose(auction *Auction, amount float64) bool {
	// TODO: BuyNow always auto-closes on successful bid
	return true
}

// newStrategy is a factory that creates a strategy from a string name.
func newStrategy(strategyType string) AuctionStrategy {
	// TODO: Return &AscendingStrategy{} for "ASCENDING"
	// TODO: Return &SealedBidStrategy{} for "SEALED"
	// TODO: Return &BuyNowStrategy{} for "BUYNOW"
	// TODO: Default to &AscendingStrategy{}
	return &AscendingStrategy{}
}

// AuctionSystem manages users and auctions with pluggable strategies.
type AuctionSystem struct {
	nextUserId    int
	nextAuctionId int
	users         map[int]User
	auctions      map[int]Auction
	strategies    map[int]AuctionStrategy // auctionId -> strategy
}

func NewAuctionSystem() *AuctionSystem {
	return &AuctionSystem{
		nextUserId:    1,
		nextAuctionId: 1,
		users:         make(map[int]User),
		auctions:      make(map[int]Auction),
		strategies:    make(map[int]AuctionStrategy),
	}
}

func (a *AuctionSystem) RegisterUser(name, userType string) int {
	// TODO: Same as before — parse type, create User, store, return ID
	return -1
}

func (a *AuctionSystem) CreateAuction(sellerId int, item string, basePrice float64, strategyType string) int {
	// TODO: Validate seller exists and is a Seller
	// TODO: Create Auction with auto-assigned ID
	// TODO: Use newStrategy(strategyType) to create and store the strategy
	// TODO: Return the auction ID
	return -1
}

func (a *AuctionSystem) PlaceBid(auctionId, buyerId int, amount float64) bool {
	// TODO: Validate auctionId, buyerId, auction is Open, buyer is not seller
	// TODO: Get auction pointer (use pointer to allow mutation)
	// TODO: Delegate to a.strategies[auctionId].AcceptBid(&auction, amount)
	// TODO: If accepted, append Bid{BidderId: buyerId, Amount: amount}
	// TODO: If a.strategies[auctionId].ShouldAutoClose(...), set status to Closed
	// TODO: Update auction in map
	// TODO: Return true if accepted
	return false
}

func (a *AuctionSystem) GetWinningBid(auctionId int) float64 {
	// TODO: Get the auction, delegate to a.strategies[auctionId].GetVisibleWinningBid(&auction)
	return -1
}

func (a *AuctionSystem) CloseAuction(auctionId int) bool {
	// TODO: Same as Part 2 — validate, transition state
	return false
}

func (a *AuctionSystem) GetAuctionStatus(auctionId int) string {
	// TODO: Same as Part 2 — switch on status enum
	return "UNKNOWN"
}

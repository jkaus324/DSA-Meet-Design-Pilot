"""Auction system — register users, create auctions with strategies, place bids."""


class AuctionOp:
    def __init__(self, kind, s1="", s2="", s3="", i1=0, i2=0, i3=0):
        self.kind = kind
        self.s1 = s1
        self.s2 = s2
        self.s3 = s3
        self.i1 = i1
        self.i2 = i2
        self.i3 = i3


# User types
BUYER = "BUYER"
SELLER = "SELLER"

# Auction statuses
OPEN = "OPEN"
CLOSED = "CLOSED"
NO_SALE = "NO_SALE"


class User:
    def __init__(self, user_id, name, type_):
        self.user_id = user_id
        self.name = name
        self.type = type_


class Auction:
    def __init__(self, auction_id, seller_id, item, base_price):
        self.auction_id = auction_id
        self.seller_id = seller_id
        self.item = item
        self.base_price = base_price
        self.status = OPEN
        self.bids = []  # list of (bidder_id, amount)


class AscendingStrategy:
    def accept_bid(self, auction, amount):
        current_highest = auction.base_price
        for _, amt in auction.bids:
            if amt > current_highest:
                current_highest = amt
        return amount > current_highest

    def get_visible_winning_bid(self, auction):
        if not auction.bids:
            return -1
        return max(amt for _, amt in auction.bids)

    def should_auto_close(self, auction, amount):
        return False


class SealedBidStrategy:
    def accept_bid(self, auction, amount):
        return amount > auction.base_price

    def get_visible_winning_bid(self, auction):
        if auction.status == OPEN:
            return -1
        if not auction.bids:
            return -1
        return max(amt for _, amt in auction.bids)

    def should_auto_close(self, auction, amount):
        return False


class BuyNowStrategy:
    def accept_bid(self, auction, amount):
        return amount >= auction.base_price * 1.5

    def get_visible_winning_bid(self, auction):
        if not auction.bids:
            return -1
        return auction.bids[-1][1]

    def should_auto_close(self, auction, amount):
        return True


def create_strategy(type_):
    if type_ == "SEALED":
        return SealedBidStrategy()
    if type_ == "BUYNOW":
        return BuyNowStrategy()
    return AscendingStrategy()


class AuctionSystem:
    def __init__(self):
        self.next_user_id = 1
        self.next_auction_id = 1
        self.users = {}
        self.auctions = {}
        self.strategies = {}

    def register_user(self, name, type_):
        ut = SELLER if type_ == "SELLER" else BUYER
        uid = self.next_user_id
        self.next_user_id += 1
        self.users[uid] = User(uid, name, ut)
        return uid

    def create_auction(self, seller_id, item, base_price, strategy_type="ASCENDING"):
        if seller_id not in self.users:
            return -1
        if self.users[seller_id].type != SELLER:
            return -1
        aid = self.next_auction_id
        self.next_auction_id += 1
        self.auctions[aid] = Auction(aid, seller_id, item, base_price)
        self.strategies[aid] = create_strategy(strategy_type)
        return aid

    def place_bid(self, auction_id, buyer_id, amount):
        if auction_id not in self.auctions:
            return False
        if buyer_id not in self.users:
            return False
        if self.users[buyer_id].type != BUYER:
            return False
        auction = self.auctions[auction_id]
        if auction.status != OPEN:
            return False
        if buyer_id == auction.seller_id:
            return False
        strat = self.strategies[auction_id]
        if not strat.accept_bid(auction, amount):
            return False
        auction.bids.append((buyer_id, amount))
        if strat.should_auto_close(auction, amount):
            auction.status = CLOSED
        return True

    def get_winning_bid(self, auction_id):
        if auction_id not in self.auctions:
            return -1
        return self.strategies[auction_id].get_visible_winning_bid(self.auctions[auction_id])

    def close_auction(self, auction_id):
        if auction_id not in self.auctions:
            return False
        auction = self.auctions[auction_id]
        if auction.status != OPEN:
            return False
        if not auction.bids:
            auction.status = NO_SALE
        else:
            auction.status = CLOSED
        return True

    def get_auction_status(self, auction_id):
        if auction_id not in self.auctions:
            return "UNKNOWN"
        return self.auctions[auction_id].status


def _format_winning(w):
    if w < 0:
        return "-1"
    if float(w).is_integer():
        return str(int(w))
    return f"{w:.2f}"


def auction_simulate(ops):
    out = []
    sys = AuctionSystem()
    user_slot = {}
    auction_slot = {}
    for op in ops:
        k = op.kind
        if k == "new":
            sys = AuctionSystem()
            user_slot = {}
            auction_slot = {}
            out.append("ok")
        elif k == "register":
            uid = sys.register_user(op.s1, op.s2)
            user_slot[op.i1] = uid
            out.append(str(uid))
        elif k == "create":
            sid = user_slot.get(op.i1, op.i1)
            strat = op.s3 if op.s3 else "ASCENDING"
            aid = sys.create_auction(sid, op.s2, float(op.i3), strat)
            auction_slot[op.i2] = aid
            out.append(str(aid))
        elif k == "bid":
            aid = auction_slot.get(op.i1, op.i1)
            bid = user_slot.get(op.i2, op.i2)
            ok = sys.place_bid(aid, bid, float(op.i3))
            out.append("ok" if ok else "fail")
        elif k == "close":
            aid = auction_slot.get(op.i1, op.i1)
            out.append("ok" if sys.close_auction(aid) else "fail")
        elif k == "winning":
            aid = auction_slot.get(op.i1, op.i1)
            out.append(_format_winning(sys.get_winning_bid(aid)))
        elif k == "status":
            aid = auction_slot.get(op.i1, op.i1)
            out.append(sys.get_auction_status(aid))
        elif k == "user_id_eq":
            uid = user_slot.get(op.i1, op.i1)
            out.append("yes" if uid == op.i2 else "no")
        else:
            out.append("unknown:" + k)
    return out

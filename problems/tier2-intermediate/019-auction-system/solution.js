"use strict";
// Auction system — register users, create auctions with strategies, place bids.

class AuctionOp {
  constructor(kind, s1 = "", s2 = "", s3 = "", i1 = 0, i2 = 0, i3 = 0) {
    this.kind = kind;
    this.s1 = s1;
    this.s2 = s2;
    this.s3 = s3;
    this.i1 = i1;
    this.i2 = i2;
    this.i3 = i3;
  }
}

// User types
const BUYER = "BUYER";
const SELLER = "SELLER";

// Auction statuses
const OPEN = "OPEN";
const CLOSED = "CLOSED";
const NO_SALE = "NO_SALE";

class User {
  constructor(userId, name, type_) {
    this.user_id = userId;
    this.name = name;
    this.type = type_;
  }
}

class Auction {
  constructor(auctionId, sellerId, item, basePrice) {
    this.auction_id = auctionId;
    this.seller_id = sellerId;
    this.item = item;
    this.base_price = basePrice;
    this.status = OPEN;
    this.bids = []; // list of [bidder_id, amount]
  }
}

class AscendingStrategy {
  accept_bid(auction, amount) {
    let currentHighest = auction.base_price;
    for (const [, amt] of auction.bids) {
      if (amt > currentHighest) currentHighest = amt;
    }
    return amount > currentHighest;
  }
  get_visible_winning_bid(auction) {
    if (auction.bids.length === 0) return -1;
    let m = auction.bids[0][1];
    for (const [, amt] of auction.bids) if (amt > m) m = amt;
    return m;
  }
  should_auto_close(auction, amount) {
    return false;
  }
}

class SealedBidStrategy {
  accept_bid(auction, amount) {
    return amount > auction.base_price;
  }
  get_visible_winning_bid(auction) {
    if (auction.status === OPEN) return -1;
    if (auction.bids.length === 0) return -1;
    let m = auction.bids[0][1];
    for (const [, amt] of auction.bids) if (amt > m) m = amt;
    return m;
  }
  should_auto_close(auction, amount) {
    return false;
  }
}

class BuyNowStrategy {
  accept_bid(auction, amount) {
    return amount >= auction.base_price * 1.5;
  }
  get_visible_winning_bid(auction) {
    if (auction.bids.length === 0) return -1;
    return auction.bids[auction.bids.length - 1][1];
  }
  should_auto_close(auction, amount) {
    return true;
  }
}

function create_strategy(type_) {
  if (type_ === "SEALED") return new SealedBidStrategy();
  if (type_ === "BUYNOW") return new BuyNowStrategy();
  return new AscendingStrategy();
}

class AuctionSystem {
  constructor() {
    this.next_user_id = 1;
    this.next_auction_id = 1;
    this.users = new Map();
    this.auctions = new Map();
    this.strategies = new Map();
  }

  register_user(name, type_) {
    const ut = type_ === "SELLER" ? SELLER : BUYER;
    const uid = this.next_user_id;
    this.next_user_id += 1;
    this.users.set(uid, new User(uid, name, ut));
    return uid;
  }

  create_auction(sellerId, item, basePrice, strategyType = "ASCENDING") {
    if (!this.users.has(sellerId)) return -1;
    if (this.users.get(sellerId).type !== SELLER) return -1;
    const aid = this.next_auction_id;
    this.next_auction_id += 1;
    this.auctions.set(aid, new Auction(aid, sellerId, item, basePrice));
    this.strategies.set(aid, create_strategy(strategyType));
    return aid;
  }

  place_bid(auctionId, buyerId, amount) {
    if (!this.auctions.has(auctionId)) return false;
    if (!this.users.has(buyerId)) return false;
    if (this.users.get(buyerId).type !== BUYER) return false;
    const auction = this.auctions.get(auctionId);
    if (auction.status !== OPEN) return false;
    if (buyerId === auction.seller_id) return false;
    const strat = this.strategies.get(auctionId);
    if (!strat.accept_bid(auction, amount)) return false;
    auction.bids.push([buyerId, amount]);
    if (strat.should_auto_close(auction, amount)) {
      auction.status = CLOSED;
    }
    return true;
  }

  get_winning_bid(auctionId) {
    if (!this.auctions.has(auctionId)) return -1;
    return this.strategies
      .get(auctionId)
      .get_visible_winning_bid(this.auctions.get(auctionId));
  }

  close_auction(auctionId) {
    if (!this.auctions.has(auctionId)) return false;
    const auction = this.auctions.get(auctionId);
    if (auction.status !== OPEN) return false;
    if (auction.bids.length === 0) {
      auction.status = NO_SALE;
    } else {
      auction.status = CLOSED;
    }
    return true;
  }

  get_auction_status(auctionId) {
    if (!this.auctions.has(auctionId)) return "UNKNOWN";
    return this.auctions.get(auctionId).status;
  }
}

function _format_winning(w) {
  if (w < 0) return "-1";
  if (Number.isInteger(w)) return String(w);
  return w.toFixed(2);
}

function auction_simulate(ops) {
  const out = [];
  let sys = new AuctionSystem();
  let userSlot = new Map();
  let auctionSlot = new Map();
  const slotGet = (m, key) => (m.has(key) ? m.get(key) : key);
  for (const op of ops) {
    const k = op.kind;
    if (k === "new") {
      sys = new AuctionSystem();
      userSlot = new Map();
      auctionSlot = new Map();
      out.push("ok");
    } else if (k === "register") {
      const uid = sys.register_user(op.s1, op.s2);
      userSlot.set(op.i1, uid);
      out.push(String(uid));
    } else if (k === "create") {
      const sid = slotGet(userSlot, op.i1);
      const strat = op.s3 ? op.s3 : "ASCENDING";
      const aid = sys.create_auction(sid, op.s2, op.i3, strat);
      auctionSlot.set(op.i2, aid);
      out.push(String(aid));
    } else if (k === "bid") {
      const aid = slotGet(auctionSlot, op.i1);
      const bid = slotGet(userSlot, op.i2);
      const ok = sys.place_bid(aid, bid, op.i3);
      out.push(ok ? "ok" : "fail");
    } else if (k === "close") {
      const aid = slotGet(auctionSlot, op.i1);
      out.push(sys.close_auction(aid) ? "ok" : "fail");
    } else if (k === "winning") {
      const aid = slotGet(auctionSlot, op.i1);
      out.push(_format_winning(sys.get_winning_bid(aid)));
    } else if (k === "status") {
      const aid = slotGet(auctionSlot, op.i1);
      out.push(sys.get_auction_status(aid));
    } else if (k === "user_id_eq") {
      const uid = slotGet(userSlot, op.i1);
      out.push(uid === op.i2 ? "yes" : "no");
    } else {
      out.push("unknown:" + k);
    }
  }
  return out;
}

module.exports = {
  AuctionOp,
  User,
  Auction,
  AscendingStrategy,
  SealedBidStrategy,
  BuyNowStrategy,
  AuctionSystem,
  create_strategy,
  auction_simulate,
};

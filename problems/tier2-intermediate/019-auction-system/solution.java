// Auction System — Solution (Java)
import java.util.*;

class AuctionOp {
    public String kind;
    public String s1, s2, s3;
    public int i1, i2, i3;

    public AuctionOp(String kind, String s1, String s2, String s3, int i1, int i2, int i3) {
        this.kind = kind; this.s1 = s1; this.s2 = s2; this.s3 = s3;
        this.i1 = i1; this.i2 = i2; this.i3 = i3;
    }
}

enum UserType { BUYER, SELLER }
enum AuctionStatus { OPEN, CLOSED, NO_SALE }

class AuctionUser {
    public int userId;
    public String name;
    public UserType type;
    public AuctionUser(int id, String name, UserType t) { this.userId = id; this.name = name; this.type = t; }
}

class Bid {
    public int bidderId;
    public double amount;
    public Bid(int b, double a) { bidderId = b; amount = a; }
}

class Auction {
    public int auctionId;
    public int sellerId;
    public String item;
    public double basePrice;
    public AuctionStatus status;
    public List<Bid> bids = new ArrayList<>();

    public Auction(int aid, int sid, String item, double bp) {
        this.auctionId = aid; this.sellerId = sid;
        this.item = item; this.basePrice = bp; this.status = AuctionStatus.OPEN;
    }
}

interface AuctionStrategy {
    boolean acceptBid(Auction auction, double amount);
    double getVisibleWinningBid(Auction auction);
    boolean shouldAutoClose(Auction auction, double amount);
}

class AscendingStrategy implements AuctionStrategy {
    @Override public boolean acceptBid(Auction auction, double amount) {
        double currentHighest = auction.basePrice;
        for (Bid b : auction.bids) if (b.amount > currentHighest) currentHighest = b.amount;
        return amount > currentHighest;
    }
    @Override public double getVisibleWinningBid(Auction auction) {
        if (auction.bids.isEmpty()) return -1;
        double maxBid = auction.bids.get(0).amount;
        for (Bid b : auction.bids) if (b.amount > maxBid) maxBid = b.amount;
        return maxBid;
    }
    @Override public boolean shouldAutoClose(Auction auction, double amount) { return false; }
}

class SealedBidStrategy implements AuctionStrategy {
    @Override public boolean acceptBid(Auction auction, double amount) {
        return amount > auction.basePrice;
    }
    @Override public double getVisibleWinningBid(Auction auction) {
        if (auction.status == AuctionStatus.OPEN) return -1;
        if (auction.bids.isEmpty()) return -1;
        double maxBid = auction.bids.get(0).amount;
        for (Bid b : auction.bids) if (b.amount > maxBid) maxBid = b.amount;
        return maxBid;
    }
    @Override public boolean shouldAutoClose(Auction auction, double amount) { return false; }
}

class BuyNowStrategy implements AuctionStrategy {
    @Override public boolean acceptBid(Auction auction, double amount) {
        return amount >= auction.basePrice * 1.5;
    }
    @Override public double getVisibleWinningBid(Auction auction) {
        if (auction.bids.isEmpty()) return -1;
        return auction.bids.get(auction.bids.size() - 1).amount;
    }
    @Override public boolean shouldAutoClose(Auction auction, double amount) { return true; }
}

class AuctionSystem {
    int nextUserId = 1;
    int nextAuctionId = 1;
    Map<Integer, AuctionUser> users = new LinkedHashMap<>();
    Map<Integer, Auction> auctions = new LinkedHashMap<>();
    Map<Integer, AuctionStrategy> strategies = new LinkedHashMap<>();

    static AuctionStrategy createStrategy(String type) {
        if ("SEALED".equals(type)) return new SealedBidStrategy();
        if ("BUYNOW".equals(type)) return new BuyNowStrategy();
        return new AscendingStrategy();
    }

    public int registerUser(String name, String type) {
        UserType ut = "SELLER".equals(type) ? UserType.SELLER : UserType.BUYER;
        int id = nextUserId++;
        users.put(id, new AuctionUser(id, name, ut));
        return id;
    }

    public int createAuction(int sellerId, String item, double basePrice, String strategyType) {
        AuctionUser u = users.get(sellerId);
        if (u == null || u.type != UserType.SELLER) return -1;
        int id = nextAuctionId++;
        auctions.put(id, new Auction(id, sellerId, item, basePrice));
        strategies.put(id, createStrategy(strategyType));
        return id;
    }

    public boolean placeBid(int auctionId, int buyerId, double amount) {
        Auction auction = auctions.get(auctionId);
        if (auction == null) return false;
        AuctionUser u = users.get(buyerId);
        if (u == null || u.type != UserType.BUYER) return false;
        if (auction.status != AuctionStatus.OPEN) return false;
        if (buyerId == auction.sellerId) return false;

        AuctionStrategy strat = strategies.get(auctionId);
        if (!strat.acceptBid(auction, amount)) return false;

        auction.bids.add(new Bid(buyerId, amount));

        if (strat.shouldAutoClose(auction, amount)) auction.status = AuctionStatus.CLOSED;
        return true;
    }

    public double getWinningBid(int auctionId) {
        Auction a = auctions.get(auctionId);
        if (a == null) return -1;
        return strategies.get(auctionId).getVisibleWinningBid(a);
    }

    public boolean closeAuction(int auctionId) {
        Auction a = auctions.get(auctionId);
        if (a == null) return false;
        if (a.status != AuctionStatus.OPEN) return false;
        a.status = a.bids.isEmpty() ? AuctionStatus.NO_SALE : AuctionStatus.CLOSED;
        return true;
    }

    public String getAuctionStatus(int auctionId) {
        Auction a = auctions.get(auctionId);
        if (a == null) return "UNKNOWN";
        return a.status.name();
    }
}

public class Solution {
    public static List<String> auction_simulate(List<AuctionOp> ops) {
        List<String> out = new ArrayList<>();
        AuctionSystem sys = new AuctionSystem();
        Map<Integer, Integer> userSlot = new HashMap<>();
        Map<Integer, Integer> auctionSlot = new HashMap<>();

        for (AuctionOp op : ops) {
            String k = op.kind;
            if ("new".equals(k)) {
                sys = new AuctionSystem();
                userSlot.clear();
                auctionSlot.clear();
                out.add("ok");
            } else if ("register".equals(k)) {
                int id = sys.registerUser(op.s1, op.s2);
                userSlot.put(op.i1, id);
                out.add(Integer.toString(id));
            } else if ("create".equals(k)) {
                int sid = userSlot.containsKey(op.i1) ? userSlot.get(op.i1) : op.i1;
                String strat = (op.s3 == null || op.s3.isEmpty()) ? "ASCENDING" : op.s3;
                int aid = sys.createAuction(sid, op.s2, (double) op.i3, strat);
                auctionSlot.put(op.i2, aid);
                out.add(Integer.toString(aid));
            } else if ("bid".equals(k)) {
                int aid = auctionSlot.containsKey(op.i1) ? auctionSlot.get(op.i1) : op.i1;
                int bid = userSlot.containsKey(op.i2) ? userSlot.get(op.i2) : op.i2;
                out.add(sys.placeBid(aid, bid, (double) op.i3) ? "ok" : "fail");
            } else if ("close".equals(k)) {
                int aid = auctionSlot.containsKey(op.i1) ? auctionSlot.get(op.i1) : op.i1;
                out.add(sys.closeAuction(aid) ? "ok" : "fail");
            } else if ("winning".equals(k)) {
                int aid = auctionSlot.containsKey(op.i1) ? auctionSlot.get(op.i1) : op.i1;
                double w = sys.getWinningBid(aid);
                if (w < 0) out.add("-1");
                else if (w == (long) w) out.add(Long.toString((long) w));
                else out.add(String.format(java.util.Locale.US, "%.2f", w));
            } else if ("status".equals(k)) {
                int aid = auctionSlot.containsKey(op.i1) ? auctionSlot.get(op.i1) : op.i1;
                out.add(sys.getAuctionStatus(aid));
            } else if ("user_id_eq".equals(k)) {
                int uid = userSlot.containsKey(op.i1) ? userSlot.get(op.i1) : op.i1;
                out.add(uid == op.i2 ? "yes" : "no");
            } else {
                out.add("unknown:" + k);
            }
        }
        return out;
    }
}

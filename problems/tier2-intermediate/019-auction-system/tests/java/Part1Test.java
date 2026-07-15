// Auction System — Part 1 Tests
import java.util.*;
import java.util.stream.*;

class Part1Test {
    static boolean testRegisterAndCreate() {
        try {
            AuctionSystem sys = new AuctionSystem();
            int seller = sys.registerUser("Alice", "SELLER");
            int buyer = sys.registerUser("Bob", "BUYER");
            int auctionId = sys.createAuction(seller, "Laptop", 500.0);
            boolean pass = seller == 1
                && buyer == 2
                && auctionId == 1;
            System.out.println((pass ? "PASS" : "FAIL") + ": testRegisterAndCreate");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testRegisterAndCreate (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testPlaceValidBid() {
        try {
            AuctionSystem sys = new AuctionSystem();
            int seller = sys.registerUser("Alice", "SELLER");
            int buyer = sys.registerUser("Bob", "BUYER");
            int aId = sys.createAuction(seller, "Phone", 100.0);
            boolean result = sys.placeBid(aId, buyer, 150.0);
            boolean pass = result == true
                && sys.getWinningBid(aId) == 150.0;
            System.out.println((pass ? "PASS" : "FAIL") + ": testPlaceValidBid");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testPlaceValidBid (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testBidMustExceed() {
        try {
            AuctionSystem sys = new AuctionSystem();
            int seller = sys.registerUser("Alice", "SELLER");
            int buyer1 = sys.registerUser("Bob", "BUYER");
            int buyer2 = sys.registerUser("Charlie", "BUYER");
            int aId = sys.createAuction(seller, "Watch", 100.0);
            boolean pass = sys.placeBid(aId, buyer1, 200.0) == true
                && sys.placeBid(aId, buyer2, 150.0) == false);  // below current highest
                && sys.placeBid(aId, buyer2, 200.0) == false);  // equal, not exceeding
                && sys.placeBid(aId, buyer2, 250.0) == true);   // exceeds
                && sys.getWinningBid(aId) == 250.0;
            System.out.println((pass ? "PASS" : "FAIL") + ": testBidMustExceed");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testBidMustExceed (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testBidExceedsBasePrice() {
        try {
            AuctionSystem sys = new AuctionSystem();
            int seller = sys.registerUser("Alice", "SELLER");
            int buyer = sys.registerUser("Bob", "BUYER");
            int aId = sys.createAuction(seller, "Book", 50.0);
            boolean pass = sys.placeBid(aId, buyer, 30.0) == false);  // below base price
                && sys.placeBid(aId, buyer, 50.0) == false);  // equal to base price
                && sys.placeBid(aId, buyer, 51.0) == true);   // above base price;
            System.out.println((pass ? "PASS" : "FAIL") + ": testBidExceedsBasePrice");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testBidExceedsBasePrice (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testOnlyBuyersBid() {
        try {
            AuctionSystem sys = new AuctionSystem();
            int seller1 = sys.registerUser("Alice", "SELLER");
            int seller2 = sys.registerUser("Bob", "SELLER");
            int aId = sys.createAuction(seller1, "Tablet", 200.0);
            boolean pass = sys.placeBid(aId, seller2, 300.0) == false);  // seller cannot bid;
            System.out.println((pass ? "PASS" : "FAIL") + ": testOnlyBuyersBid");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testOnlyBuyersBid (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testOnlySellersCreate() {
        try {
            AuctionSystem sys = new AuctionSystem();
            int buyer = sys.registerUser("Bob", "BUYER");
            int aId = sys.createAuction(buyer, "Camera", 300.0);
            boolean pass = aId == -1);  // buyer cannot create auction;
            System.out.println((pass ? "PASS" : "FAIL") + ": testOnlySellersCreate");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testOnlySellersCreate (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testNoBidsReturnsNegative() {
        try {
            AuctionSystem sys = new AuctionSystem();
            int seller = sys.registerUser("Alice", "SELLER");
            int aId = sys.createAuction(seller, "Keyboard", 75.0);
            boolean pass = sys.getWinningBid(aId) == -1;
            System.out.println((pass ? "PASS" : "FAIL") + ": testNoBidsReturnsNegative");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testNoBidsReturnsNegative (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testIndependentAuctions() {
        try {
            AuctionSystem sys = new AuctionSystem();
            int seller = sys.registerUser("Alice", "SELLER");
            int buyer = sys.registerUser("Bob", "BUYER");
            int a1 = sys.createAuction(seller, "Item1", 100.0);
            int a2 = sys.createAuction(seller, "Item2", 200.0);
            sys.placeBid(a1, buyer, 150.0);
            sys.placeBid(a2, buyer, 300.0);
            boolean pass = sys.getWinningBid(a1) == 150.0
                && sys.getWinningBid(a2) == 300.0;
            System.out.println((pass ? "PASS" : "FAIL") + ": testIndependentAuctions");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testIndependentAuctions (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    public static int runTests() {
        int passed = 0, total = 0;
        total++; if (testRegisterAndCreate()) passed++;
        total++; if (testPlaceValidBid()) passed++;
        total++; if (testBidMustExceed()) passed++;
        total++; if (testBidExceedsBasePrice()) passed++;
        total++; if (testOnlyBuyersBid()) passed++;
        total++; if (testOnlySellersCreate()) passed++;
        total++; if (testNoBidsReturnsNegative()) passed++;
        total++; if (testIndependentAuctions()) passed++;
        System.out.println("PART1_SUMMARY " + passed + "/" + total);
        return passed;
    }
}

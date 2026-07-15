// Auction System — Part 3 Tests
import java.util.*;
import java.util.stream.*;

class Part3Test {
    static boolean testAscendingStrategy() {
        try {
            AuctionSystem sys = new AuctionSystem();
            int seller = sys.registerUser("Alice", "SELLER");
            int buyer1 = sys.registerUser("Bob", "BUYER");
            int buyer2 = sys.registerUser("Charlie", "BUYER");
            int aId = sys.createAuction(seller, "Laptop", 500.0, "ASCENDING");
            boolean pass = sys.placeBid(aId, buyer1, 600.0) == true
                && sys.placeBid(aId, buyer2, 550.0) == false);  // must exceed 600
                && sys.placeBid(aId, buyer2, 700.0) == true
                && sys.getWinningBid(aId) == 700.0;
            System.out.println((pass ? "PASS" : "FAIL") + ": testAscendingStrategy");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testAscendingStrategy (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testSealedBidStrategy() {
        try {
            AuctionSystem sys = new AuctionSystem();
            int seller = sys.registerUser("Alice", "SELLER");
            int buyer1 = sys.registerUser("Bob", "BUYER");
            int buyer2 = sys.registerUser("Charlie", "BUYER");
            int aId = sys.createAuction(seller, "Art", 100.0, "SEALED");
            // Any bid above base price is accepted
            // While open, winning bid is hidden
            // After close, winner is revealed
            sys.closeAuction(aId);
            boolean pass = sys.placeBid(aId, buyer1, 500.0) == true
                && sys.placeBid(aId, buyer2, 200.0) == true);  // lower than 500 but still valid
                && sys.getWinningBid(aId) == -1
                && sys.getWinningBid(aId) == 500.0);  // highest bid wins;
            System.out.println((pass ? "PASS" : "FAIL") + ": testSealedBidStrategy");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testSealedBidStrategy (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testSealedRejectsLowBids() {
        try {
            AuctionSystem sys = new AuctionSystem();
            int seller = sys.registerUser("Alice", "SELLER");
            int buyer = sys.registerUser("Bob", "BUYER");
            int aId = sys.createAuction(seller, "Vase", 200.0, "SEALED");
            boolean pass = sys.placeBid(aId, buyer, 100.0) == false);  // below base
                && sys.placeBid(aId, buyer, 200.0) == false);  // equal to base
                && sys.placeBid(aId, buyer, 201.0) == true);   // above base;
            System.out.println((pass ? "PASS" : "FAIL") + ": testSealedRejectsLowBids");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testSealedRejectsLowBids (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testBuynowStrategy() {
        try {
            AuctionSystem sys = new AuctionSystem();
            int seller = sys.registerUser("Alice", "SELLER");
            int buyer = sys.registerUser("Bob", "BUYER");
            int aId = sys.createAuction(seller, "Guitar", 100.0, "BUYNOW");
            // Buy-now price = 100 * 1.5 = 150
            // Should var-close
            boolean pass = sys.placeBid(aId, buyer, 120.0) == false);  // below buy-now price
                && sys.placeBid(aId, buyer, 150.0) == true);   // meets buy-now price
                && sys.getAuctionStatus(aId) == "CLOSED"
                && sys.getWinningBid(aId) == 150.0;
            System.out.println((pass ? "PASS" : "FAIL") + ": testBuynowStrategy");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testBuynowStrategy (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testBuynowNoBidsAfterClose() {
        try {
            AuctionSystem sys = new AuctionSystem();
            int seller = sys.registerUser("Alice", "SELLER");
            int buyer1 = sys.registerUser("Bob", "BUYER");
            int buyer2 = sys.registerUser("Charlie", "BUYER");
            int aId = sys.createAuction(seller, "Drum", 200.0, "BUYNOW");
            // Buy-now price = 200 * 1.5 = 300
            boolean pass = sys.placeBid(aId, buyer1, 300.0) == true);   // var-closes
                && sys.placeBid(aId, buyer2, 400.0) == false);  // auction already closed;
            System.out.println((pass ? "PASS" : "FAIL") + ": testBuynowNoBidsAfterClose");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testBuynowNoBidsAfterClose (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testDefaultAscending() {
        try {
            AuctionSystem sys = new AuctionSystem();
            int seller = sys.registerUser("Alice", "SELLER");
            int buyer = sys.registerUser("Bob", "BUYER");
            int aId = sys.createAuction(seller, "Mouse", 50.0);  // no strategy specified
            boolean pass = sys.placeBid(aId, buyer, 60.0) == true
                && sys.getWinningBid(aId) == 60.0);  // visible (ascending behavior;
            System.out.println((pass ? "PASS" : "FAIL") + ": testDefaultAscending");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testDefaultAscending (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testMixedStrategies() {
        try {
            AuctionSystem sys = new AuctionSystem();
            int seller = sys.registerUser("Alice", "SELLER");
            int buyer = sys.registerUser("Bob", "BUYER");
            int a1 = sys.createAuction(seller, "Item1", 100.0, "ASCENDING");
            int a2 = sys.createAuction(seller, "Item2", 100.0, "SEALED");
            int a3 = sys.createAuction(seller, "Item3", 100.0, "BUYNOW");
            sys.placeBid(a1, buyer, 200.0);
            sys.placeBid(a2, buyer, 200.0);
            sys.placeBid(a3, buyer, 150.0);
            boolean pass = sys.getWinningBid(a1) == 200.0);   // ascending: visible
                && sys.getWinningBid(a2) == -1);       // sealed: hidden while open
                && sys.getAuctionStatus(a3) == "CLOSED");  // buynow: var-closed
                && sys.getWinningBid(a3) == 150.0);   // buynow: visible after close;
            System.out.println((pass ? "PASS" : "FAIL") + ": testMixedStrategies");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testMixedStrategies (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    public static int runTests() {
        int passed = 0, total = 0;
        total++; if (testAscendingStrategy()) passed++;
        total++; if (testSealedBidStrategy()) passed++;
        total++; if (testSealedRejectsLowBids()) passed++;
        total++; if (testBuynowStrategy()) passed++;
        total++; if (testBuynowNoBidsAfterClose()) passed++;
        total++; if (testDefaultAscending()) passed++;
        total++; if (testMixedStrategies()) passed++;
        System.out.println("PART3_SUMMARY " + passed + "/" + total);
        return passed;
    }
}

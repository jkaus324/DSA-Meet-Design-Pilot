// Auction System — Part 2 Tests
import java.util.*;
import java.util.stream.*;

class Part2Test {
    static boolean testCloseWithBids() {
        try {
            AuctionSystem sys = new AuctionSystem();
            int seller = sys.registerUser("Alice", "SELLER");
            int buyer = sys.registerUser("Bob", "BUYER");
            int aId = sys.createAuction(seller, "Laptop", 500.0);
            sys.placeBid(aId, buyer, 600.0);
            boolean pass = sys.closeAuction(aId) == true
                && sys.getAuctionStatus(aId) == "CLOSED"
                && sys.getWinningBid(aId) == 600.0;
            System.out.println((pass ? "PASS" : "FAIL") + ": testCloseWithBids");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testCloseWithBids (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testCloseNoBids() {
        try {
            AuctionSystem sys = new AuctionSystem();
            int seller = sys.registerUser("Alice", "SELLER");
            int aId = sys.createAuction(seller, "Phone", 300.0);
            boolean pass = sys.closeAuction(aId) == true
                && sys.getAuctionStatus(aId) == "NO_SALE"
                && sys.getWinningBid(aId) == -1;
            System.out.println((pass ? "PASS" : "FAIL") + ": testCloseNoBids");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testCloseNoBids (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testDoubleClose() {
        try {
            AuctionSystem sys = new AuctionSystem();
            int seller = sys.registerUser("Alice", "SELLER");
            int buyer = sys.registerUser("Bob", "BUYER");
            int aId = sys.createAuction(seller, "Watch", 100.0);
            sys.placeBid(aId, buyer, 150.0);
            boolean pass = sys.closeAuction(aId) == true
                && sys.closeAuction(aId) == false);  // already closed;
            System.out.println((pass ? "PASS" : "FAIL") + ": testDoubleClose");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testDoubleClose (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testBidOnClosed() {
        try {
            AuctionSystem sys = new AuctionSystem();
            int seller = sys.registerUser("Alice", "SELLER");
            int buyer = sys.registerUser("Bob", "BUYER");
            int aId = sys.createAuction(seller, "Tablet", 200.0);
            sys.closeAuction(aId);
            boolean pass = sys.placeBid(aId, buyer, 300.0) == false;
            System.out.println((pass ? "PASS" : "FAIL") + ": testBidOnClosed");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testBidOnClosed (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testInitialStatusOpen() {
        try {
            AuctionSystem sys = new AuctionSystem();
            int seller = sys.registerUser("Alice", "SELLER");
            int aId = sys.createAuction(seller, "Camera", 400.0);
            boolean pass = sys.getAuctionStatus(aId) == "OPEN";
            System.out.println((pass ? "PASS" : "FAIL") + ": testInitialStatusOpen");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testInitialStatusOpen (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testCloseNosaleAgain() {
        try {
            AuctionSystem sys = new AuctionSystem();
            int seller = sys.registerUser("Alice", "SELLER");
            int aId = sys.createAuction(seller, "Book", 25.0);
            boolean pass = sys.closeAuction(aId) == true);   // NO_SALE
                && sys.closeAuction(aId) == false);   // already in terminal state
                && sys.getAuctionStatus(aId) == "NO_SALE";
            System.out.println((pass ? "PASS" : "FAIL") + ": testCloseNosaleAgain");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testCloseNosaleAgain (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testWinningBidAfterClose() {
        try {
            AuctionSystem sys = new AuctionSystem();
            int seller = sys.registerUser("Alice", "SELLER");
            int buyer1 = sys.registerUser("Bob", "BUYER");
            int buyer2 = sys.registerUser("Charlie", "BUYER");
            int aId = sys.createAuction(seller, "Painting", 1000.0);
            sys.placeBid(aId, buyer1, 1500.0);
            sys.placeBid(aId, buyer2, 2000.0);
            sys.closeAuction(aId);
            boolean pass = sys.getWinningBid(aId) == 2000.0
                && sys.getAuctionStatus(aId) == "CLOSED";
            System.out.println((pass ? "PASS" : "FAIL") + ": testWinningBidAfterClose");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testWinningBidAfterClose (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    public static int runTests() {
        int passed = 0, total = 0;
        total++; if (testCloseWithBids()) passed++;
        total++; if (testCloseNoBids()) passed++;
        total++; if (testDoubleClose()) passed++;
        total++; if (testBidOnClosed()) passed++;
        total++; if (testInitialStatusOpen()) passed++;
        total++; if (testCloseNosaleAgain()) passed++;
        total++; if (testWinningBidAfterClose()) passed++;
        System.out.println("PART2_SUMMARY " + passed + "/" + total);
        return passed;
    }
}

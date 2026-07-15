// Lru Cache — Part 3 Tests
import java.util.*;
import java.util.stream.*;

class TestEvictionListener implements EvictionListener {
    List<EvictionEvent> events;
    void onEviction(int key, int value, EvictionReason reason)  {
        events.add({key, value, reason});
    }
    void clear() { events.clear(); }
}

class Part3Test {
    static boolean testCapacityEvictionListener() {
        try {
            LRUCache cache = new LRUCache(2);
            TestEvictionListener listener = new TestEvictionListener();
            cache.addEvictionListener( listener);
            cache.put(1, 10, 1000);
            cache.put(2, 20, 1000);
            cache.put(3, 30, 1000);  // evicts key 1
            boolean pass = listener.events.size() == 1
                && listener.events[0].key == 1
                && listener.events[0].value == 10
                && listener.events[0].reason == EvictionReason.CAPACITY;
            System.out.println((pass ? "PASS" : "FAIL") + ": testCapacityEvictionListener");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testCapacityEvictionListener (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testTtlEvictionListener() {
        try {
            LRUCache cache = new LRUCache(5);
            TestEvictionListener listener = new TestEvictionListener();
            cache.addEvictionListener( listener);
            cache.put(1, 10, 1000, 30);  // expires at 1030
            cache.get(1, 1035);           // triggers TTL eviction
            boolean pass = listener.events.size() == 1
                && listener.events[0].key == 1
                && listener.events[0].reason == EvictionReason.TTL_EXPIRED;
            System.out.println((pass ? "PASS" : "FAIL") + ": testTtlEvictionListener");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testTtlEvictionListener (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testExplicitDeleteListener() {
        try {
            LRUCache cache = new LRUCache(5);
            TestEvictionListener listener = new TestEvictionListener();
            cache.addEvictionListener( listener);
            cache.put(1, 10, 1000);
            cache.deleteKey(1);
            boolean pass = listener.events.size() == 1
                && listener.events[0].key == 1
                && listener.events[0].value == 10
                && listener.events[0].reason == EvictionReason.EXPLICIT_DELETE;
            System.out.println((pass ? "PASS" : "FAIL") + ": testExplicitDeleteListener");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testExplicitDeleteListener (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testMultipleListeners() {
        try {
            LRUCache cache = new LRUCache(2);
            TestEvictionListener listener1 = new TestEvictionListener();
            TestEvictionListener listener2 = new TestEvictionListener();
            cache.addEvictionListener( listener1);
            cache.addEvictionListener( listener2);
            cache.put(1, 10, 1000);
            cache.put(2, 20, 1000);
            cache.put(3, 30, 1000);  // evicts key 1
            boolean pass = listener1.events.size() == 1
                && listener2.events.size() == 1
                && listener1.events[0].key == 1
                && listener2.events[0].key == 1;
            System.out.println((pass ? "PASS" : "FAIL") + ": testMultipleListeners");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testMultipleListeners (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testRemoveListener() {
        try {
            LRUCache cache = new LRUCache(2);
            TestEvictionListener listener = new TestEvictionListener();
            cache.addEvictionListener( listener);
            cache.put(1, 10, 1000);
            cache.put(2, 20, 1000);
            cache.put(3, 30, 1000);  // evicts key 1, listener notified
            cache.removeEvictionListener( listener);
            cache.put(4, 40, 1000);  // evicts key 2, listener NOT notified
            boolean pass = listener.events.size() == 1
                && listener.events.size() == 1);  // still 1, not 2;
            System.out.println((pass ? "PASS" : "FAIL") + ": testRemoveListener");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testRemoveListener (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testNoListenersNoCrash() {
        try {
            LRUCache cache = new LRUCache(1);
            cache.put(1, 10, 1000);
            cache.put(2, 20, 1000);  // evicts key 1, no listeners
            boolean pass = cache.get(2, 1000) == 20;
            System.out.println((pass ? "PASS" : "FAIL") + ": testNoListenersNoCrash");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testNoListenersNoCrash (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testSequentialEvictions() {
        try {
            LRUCache cache = new LRUCache(2);
            TestEvictionListener listener = new TestEvictionListener();
            cache.addEvictionListener( listener);
            cache.put(1, 10, 1000);
            cache.put(2, 20, 1000);
            cache.put(3, 30, 1000);  // evicts 1
            cache.put(4, 40, 1000);  // evicts 2
            boolean pass = listener.events.size() == 2
                && listener.events[0].key == 1
                && listener.events[0].value == 10
                && listener.events[1].key == 2
                && listener.events[1].value == 20;
            System.out.println((pass ? "PASS" : "FAIL") + ": testSequentialEvictions");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testSequentialEvictions (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    public static int runTests() {
        int passed = 0, total = 0;
        total++; if (testCapacityEvictionListener()) passed++;
        total++; if (testTtlEvictionListener()) passed++;
        total++; if (testExplicitDeleteListener()) passed++;
        total++; if (testMultipleListeners()) passed++;
        total++; if (testRemoveListener()) passed++;
        total++; if (testNoListenersNoCrash()) passed++;
        total++; if (testSequentialEvictions()) passed++;
        System.out.println("PART3_SUMMARY " + passed + "/" + total);
        return passed;
    }
}

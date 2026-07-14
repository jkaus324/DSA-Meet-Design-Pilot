// Logger System — Part 3 Tests
import java.util.*;
import java.util.stream.*;

class TestDestination implements LogDestination {
    LogFormatter formatter;
    String destName;
    List<String> received;
    TestDestination(String name, LogFormatter f)
        : destName(name), formatter(f) {}
    void write(LogEntry entry)  {
        received.add(formatter.format(entry));
    }
    String getName() { return destName; }
}

class FailingDestination implements LogDestination {
    void write(LogEntry entry)  {
        throw runtime_error("destination failed");
    }
    String getName() { return "failing"; }
}

class Part3Test {
    static boolean testSingleDestination() {
        try {
            Logger logger = Logger.getInstance();
            logger.clearHistory();
            logger.setLevel(LogLevel.DEBUG);
            PlainTextFormatter plainFmt = new PlainTextFormatter();
            TestDestination dest1 = new TestDestination("test1",  plainFmt);
            logger.addDestination( dest1);
            logger.log(LogLevel.INFO, "hello dest");
            logger.removeDestination( dest1);
            boolean pass = dest1.received.size() == 1
                && dest1.received[0].find("INFO") != String.npos
                && dest1.received[0].find("hello dest") != String.npos;
            System.out.println((pass ? "PASS" : "FAIL") + ": testSingleDestination");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testSingleDestination (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testMultipleDestinations() {
        try {
            Logger logger = Logger.getInstance();
            logger.clearHistory();
            logger.setLevel(LogLevel.DEBUG);
            PlainTextFormatter plainFmt = new PlainTextFormatter();
            JsonFormatter jsonFmt = new JsonFormatter();
            TestDestination dest1 = new TestDestination("plain-dest",  plainFmt);
            TestDestination dest2 = new TestDestination("json-dest",  jsonFmt);
            logger.addDestination( dest1);
            logger.addDestination( dest2);
            logger.log(LogLevel.ERROR, "multi dest");
            // dest1 uses PlainText
            // dest2 uses JSON
            logger.removeDestination( dest1);
            logger.removeDestination( dest2);
            boolean pass = dest1.received.size() == 1
                && dest2.received.size() == 1
                && dest1.received[0].find("[ERROR]") != String.npos
                && dest2.received[0].find("\"level\"") != String.npos;
            System.out.println((pass ? "PASS" : "FAIL") + ": testMultipleDestinations");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testMultipleDestinations (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testIndependentFormatters() {
        try {
            Logger logger = Logger.getInstance();
            logger.clearHistory();
            logger.setLevel(LogLevel.DEBUG);
            PlainTextFormatter plainFmt = new PlainTextFormatter();
            CsvFormatter csvFmt = new CsvFormatter();
            TestDestination destPlain = new TestDestination("plain",  plainFmt);
            TestDestination destCsv = new TestDestination("csv",  csvFmt);
            logger.addDestination( destPlain);
            logger.addDestination( destCsv);
            logger.log(LogLevel.WARN, "format test");
            // Plain has brackets, CSV does not
            logger.removeDestination( destPlain);
            logger.removeDestination( destCsv);
            boolean pass = destPlain.received[0].find("[WARN]") != String.npos
                && destCsv.received[0].find("[") == String.npos
                && destCsv.received[0].find(",") != String.npos;
            System.out.println((pass ? "PASS" : "FAIL") + ": testIndependentFormatters");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testIndependentFormatters (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testFaultIsolation() {
        try {
            Logger logger = Logger.getInstance();
            logger.clearHistory();
            logger.setLevel(LogLevel.DEBUG);
            PlainTextFormatter plainFmt = new PlainTextFormatter();
            FailingDestination failDest = new FailingDestination();
            TestDestination goodDest = new TestDestination("good",  plainFmt);
            logger.addDestination( failDest);
            logger.addDestination( goodDest);
            logger.log(LogLevel.ERROR, "fault isolation");
            // goodDest should still receive the message despite failDest throwing
            logger.removeDestination( failDest);
            logger.removeDestination( goodDest);
            boolean pass = goodDest.received.size() == 1
                && goodDest.received[0].find("fault isolation") != String.npos;
            System.out.println((pass ? "PASS" : "FAIL") + ": testFaultIsolation");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testFaultIsolation (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testRemoveDestination() {
        try {
            Logger logger = Logger.getInstance();
            logger.clearHistory();
            logger.setLevel(LogLevel.DEBUG);
            PlainTextFormatter plainFmt = new PlainTextFormatter();
            TestDestination dest = new TestDestination("removable",  plainFmt);
            logger.addDestination( dest);
            logger.log(LogLevel.INFO, "before remove");
            logger.removeDestination( dest);
            logger.log(LogLevel.INFO, "after remove");
            boolean pass = dest.received.size() == 1
                && dest.received.size() == 1);  // still 1, not 2;
            System.out.println((pass ? "PASS" : "FAIL") + ": testRemoveDestination");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testRemoveDestination (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testLevelFilteringBeforeDestinations() {
        try {
            Logger logger = Logger.getInstance();
            logger.clearHistory();
            logger.setLevel(LogLevel.ERROR);
            PlainTextFormatter plainFmt = new PlainTextFormatter();
            TestDestination dest = new TestDestination("level-test",  plainFmt);
            logger.addDestination( dest);
            logger.log(LogLevel.INFO, "filtered");
            logger.log(LogLevel.ERROR, "passes");
            logger.removeDestination( dest);
            boolean pass = dest.received.size() == 1
                && dest.received[0].find("passes") != String.npos;
            System.out.println((pass ? "PASS" : "FAIL") + ": testLevelFilteringBeforeDestinations");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testLevelFilteringBeforeDestinations (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testNoDestinationsNoCrash() {
        try {
            Logger logger = Logger.getInstance();
            logger.clearHistory();
            logger.setLevel(LogLevel.DEBUG);
            // No destinations added
            logger.log(LogLevel.INFO, "no destinations");
            // Should not crash, history may or may not record
            boolean pass = true; // FIXME: no assertions found
            System.out.println((pass ? "PASS" : "FAIL") + ": testNoDestinationsNoCrash");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testNoDestinationsNoCrash (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    public static int runTests() {
        int passed = 0, total = 0;
        total++; if (testSingleDestination()) passed++;
        total++; if (testMultipleDestinations()) passed++;
        total++; if (testIndependentFormatters()) passed++;
        total++; if (testFaultIsolation()) passed++;
        total++; if (testRemoveDestination()) passed++;
        total++; if (testLevelFilteringBeforeDestinations()) passed++;
        total++; if (testNoDestinationsNoCrash()) passed++;
        System.out.println("PART3_SUMMARY " + passed + "/" + total);
        return passed;
    }
}

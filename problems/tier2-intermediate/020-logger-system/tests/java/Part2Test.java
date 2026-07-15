// Logger System — Part 2 Tests
import java.util.*;
import java.util.stream.*;

class Part2Test {
    static boolean testPlaintextFormat() {
        try {
            PlainTextFormatter fmt = new PlainTextFormatter();
            LogEntry entry{LogLevel.ERROR, "disk full", "2024-01-15 10:30:00"};
            String result = fmt.format(entry);
            boolean pass = result.find("[2024-01-15 10:30:00]") != String.npos
                && result.find("[ERROR]") != String.npos
                && result.find("disk full") != String.npos;
            System.out.println((pass ? "PASS" : "FAIL") + ": testPlaintextFormat");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testPlaintextFormat (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testJsonFormat() {
        try {
            JsonFormatter fmt = new JsonFormatter();
            LogEntry entry{LogLevel.WARN, "low memory", "2024-01-15 10:30:00"};
            String result = fmt.format(entry);
            boolean pass = result.find("\"timestamp\"") != String.npos
                && result.find("\"level\"") != String.npos
                && result.find("\"message\"") != String.npos
                && result.find("WARN") != String.npos
                && result.find("low memory") != String.npos
                && result.find("2024-01-15 10:30:00") != String.npos;
            System.out.println((pass ? "PASS" : "FAIL") + ": testJsonFormat");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testJsonFormat (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testCsvFormat() {
        try {
            CsvFormatter fmt = new CsvFormatter();
            LogEntry entry{LogLevel.INFO, "server started", "2024-01-15 10:30:00"};
            String result = fmt.format(entry);
            // Should NOT have brackets (that's PlainText)
            boolean pass = result.find(",") != String.npos
                && result.find("INFO") != String.npos
                && result.find("server started") != String.npos
                && result.find("2024-01-15 10:30:00") != String.npos
                && result.find("[") == String.npos;
            System.out.println((pass ? "PASS" : "FAIL") + ": testCsvFormat");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testCsvFormat (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testDefaultPlaintextFormatter() {
        try {
            Logger logger = Logger.getInstance();
            logger.clearHistory();
            logger.setLevel(LogLevel.DEBUG);
            logger.setFormatter(null);  // reset to default
            logger.log(LogLevel.INFO, "default format");
            String entry = logger.getLogHistory()[0];
            boolean pass = entry.find("[") != String.npos
                && entry.find("[INFO]") != String.npos;
            System.out.println((pass ? "PASS" : "FAIL") + ": testDefaultPlaintextFormatter");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testDefaultPlaintextFormatter (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testSwitchToJsonFormatter() {
        try {
            Logger logger = Logger.getInstance();
            logger.clearHistory();
            logger.setLevel(LogLevel.DEBUG);
            JsonFormatter jsonFmt = new JsonFormatter();
            logger.setFormatter( jsonFmt);
            logger.log(LogLevel.ERROR, "json test");
            String entry = logger.getLogHistory()[0];
            // Reset formatter
            logger.setFormatter(null);
            boolean pass = entry.find("\"level\"") != String.npos
                && entry.find("\"message\"") != String.npos
                && entry.find("json test") != String.npos;
            System.out.println((pass ? "PASS" : "FAIL") + ": testSwitchToJsonFormatter");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testSwitchToJsonFormatter (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testSwitchToCsvFormatter() {
        try {
            Logger logger = Logger.getInstance();
            logger.clearHistory();
            logger.setLevel(LogLevel.DEBUG);
            CsvFormatter csvFmt = new CsvFormatter();
            logger.setFormatter( csvFmt);
            logger.log(LogLevel.WARN, "csv test");
            String entry = logger.getLogHistory()[0];
            // Reset formatter
            logger.setFormatter(null);
            boolean pass = entry.find(",") != String.npos
                && entry.find("WARN") != String.npos
                && entry.find("[") == String.npos);  // no brackets in CSV;
            System.out.println((pass ? "PASS" : "FAIL") + ": testSwitchToCsvFormatter");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testSwitchToCsvFormatter (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testLevelFilteringWithCustomFormatter() {
        try {
            Logger logger = Logger.getInstance();
            logger.clearHistory();
            logger.setLevel(LogLevel.ERROR);
            JsonFormatter jsonFmt = new JsonFormatter();
            logger.setFormatter( jsonFmt);
            logger.log(LogLevel.INFO, "should be filtered");
            logger.log(LogLevel.ERROR, "should pass");
            // Reset formatter
            logger.setFormatter(null);
            boolean pass = logger.getLogHistory().size() == 1
                && logger.getLogHistory()[0].find("should pass") != String.npos;
            System.out.println((pass ? "PASS" : "FAIL") + ": testLevelFilteringWithCustomFormatter");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testLevelFilteringWithCustomFormatter (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    public static int runTests() {
        int passed = 0, total = 0;
        total++; if (testPlaintextFormat()) passed++;
        total++; if (testJsonFormat()) passed++;
        total++; if (testCsvFormat()) passed++;
        total++; if (testDefaultPlaintextFormatter()) passed++;
        total++; if (testSwitchToJsonFormatter()) passed++;
        total++; if (testSwitchToCsvFormatter()) passed++;
        total++; if (testLevelFilteringWithCustomFormatter()) passed++;
        System.out.println("PART2_SUMMARY " + passed + "/" + total);
        return passed;
    }
}

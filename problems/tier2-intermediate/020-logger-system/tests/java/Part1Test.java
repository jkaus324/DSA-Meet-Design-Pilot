// Logger System — Part 1 Tests
import java.util.*;
import java.util.stream.*;

class Part1Test {
    static boolean testSingletonSameInstance() {
        try {
            Logger logger1 = Logger.getInstance();
            Logger logger2 = Logger.getInstance();
            logger1.clearHistory();
            boolean pass =  logger1 ==  logger2;
            System.out.println((pass ? "PASS" : "FAIL") + ": testSingletonSameInstance");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testSingletonSameInstance (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testDebugFilteredAtInfoLevel() {
        try {
            Logger logger = Logger.getInstance();
            logger.clearHistory();
            logger.setLevel(LogLevel.INFO);
            logger.log(LogLevel.DEBUG, "debug msg");
            boolean pass = logger.getLogHistory().isEmpty();
            System.out.println((pass ? "PASS" : "FAIL") + ": testDebugFilteredAtInfoLevel");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testDebugFilteredAtInfoLevel (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testInfoPassesAtInfoLevel() {
        try {
            Logger logger = Logger.getInstance();
            logger.clearHistory();
            logger.setLevel(LogLevel.INFO);
            logger.log(LogLevel.INFO, "info msg");
            // Check that the log contains the level and message
            String entry = logger.getLogHistory()[0];
            boolean pass = logger.getLogHistory().size() == 1
                && entry.find("INFO") != String.npos
                && entry.find("info msg") != String.npos;
            System.out.println((pass ? "PASS" : "FAIL") + ": testInfoPassesAtInfoLevel");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testInfoPassesAtInfoLevel (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testErrorFatalPassAtInfoLevel() {
        try {
            Logger logger = Logger.getInstance();
            logger.clearHistory();
            logger.setLevel(LogLevel.INFO);
            logger.log(LogLevel.ERROR, "error msg");
            logger.log(LogLevel.FATAL, "fatal msg");
            boolean pass = logger.getLogHistory().size() == 2
                && logger.getLogHistory()[0].find("ERROR") != String.npos
                && logger.getLogHistory()[1].find("FATAL") != String.npos;
            System.out.println((pass ? "PASS" : "FAIL") + ": testErrorFatalPassAtInfoLevel");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testErrorFatalPassAtInfoLevel (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testWarnLevelFiltersInfo() {
        try {
            Logger logger = Logger.getInstance();
            logger.clearHistory();
            logger.setLevel(LogLevel.WARN);
            logger.log(LogLevel.INFO, "should be filtered");
            logger.log(LogLevel.WARN, "should pass");
            logger.log(LogLevel.ERROR, "should also pass");
            boolean pass = logger.getLogHistory().size() == 2
                && logger.getLogHistory()[0].find("WARN") != String.npos
                && logger.getLogHistory()[1].find("ERROR") != String.npos;
            System.out.println((pass ? "PASS" : "FAIL") + ": testWarnLevelFiltersInfo");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testWarnLevelFiltersInfo (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testDebugLevelAllowsAll() {
        try {
            Logger logger = Logger.getInstance();
            logger.clearHistory();
            logger.setLevel(LogLevel.DEBUG);
            logger.log(LogLevel.DEBUG, "d");
            logger.log(LogLevel.INFO, "i");
            logger.log(LogLevel.WARN, "w");
            logger.log(LogLevel.ERROR, "e");
            logger.log(LogLevel.FATAL, "f");
            boolean pass = logger.getLogHistory().size() == 5;
            System.out.println((pass ? "PASS" : "FAIL") + ": testDebugLevelAllowsAll");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testDebugLevelAllowsAll (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testFatalLevelOnlyFatal() {
        try {
            Logger logger = Logger.getInstance();
            logger.clearHistory();
            logger.setLevel(LogLevel.FATAL);
            logger.log(LogLevel.DEBUG, "d");
            logger.log(LogLevel.INFO, "i");
            logger.log(LogLevel.WARN, "w");
            logger.log(LogLevel.ERROR, "e");
            logger.log(LogLevel.FATAL, "f");
            boolean pass = logger.getLogHistory().size() == 1
                && logger.getLogHistory()[0].find("FATAL") != String.npos;
            System.out.println((pass ? "PASS" : "FAIL") + ": testFatalLevelOnlyFatal");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testFatalLevelOnlyFatal (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testLogFormatHasBrackets() {
        try {
            Logger logger = Logger.getInstance();
            logger.clearHistory();
            logger.setLevel(LogLevel.INFO);
            logger.log(LogLevel.INFO, "check format");
            String entry = logger.getLogHistory()[0];
            // Should contain brackets around timestamp and level
            boolean pass = entry.find("[") != String.npos
                && entry.find("]") != String.npos
                && entry.find("check format") != String.npos;
            System.out.println((pass ? "PASS" : "FAIL") + ": testLogFormatHasBrackets");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testLogFormatHasBrackets (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    public static int runTests() {
        int passed = 0, total = 0;
        total++; if (testSingletonSameInstance()) passed++;
        total++; if (testDebugFilteredAtInfoLevel()) passed++;
        total++; if (testInfoPassesAtInfoLevel()) passed++;
        total++; if (testErrorFatalPassAtInfoLevel()) passed++;
        total++; if (testWarnLevelFiltersInfo()) passed++;
        total++; if (testDebugLevelAllowsAll()) passed++;
        total++; if (testFatalLevelOnlyFatal()) passed++;
        total++; if (testLogFormatHasBrackets()) passed++;
        System.out.println("PART1_SUMMARY " + passed + "/" + total);
        return passed;
    }
}

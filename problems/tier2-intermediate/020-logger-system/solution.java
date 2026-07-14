// Logger System — Solution (Java)
import java.util.*;

class LogOp {
    public String kind;
    public String s1, s2;
    public int i1;
    public LogOp(String kind, String s1, String s2, int i1) {
        this.kind = kind; this.s1 = s1; this.s2 = s2; this.i1 = i1;
    }
}

enum LogLevel {
    DEBUG(0), INFO(1), WARN(2), ERROR(3), FATAL(4);
    public final int rank;
    LogLevel(int r) { this.rank = r; }
}

class LogEntry {
    public LogLevel level;
    public String message;
    public String timestamp;
    public LogEntry(LogLevel l, String m, String t) { level = l; message = m; timestamp = t; }
}

interface LogFormatter {
    String format(LogEntry entry);
}

class PlainTextFormatter implements LogFormatter {
    @Override public String format(LogEntry e) {
        return "[" + e.timestamp + "] [" + e.level.name() + "] " + e.message;
    }
}

class JsonFormatter implements LogFormatter {
    @Override public String format(LogEntry e) {
        return "{\"timestamp\":\"" + e.timestamp + "\",\"level\":\"" + e.level.name()
                + "\",\"message\":\"" + e.message + "\"}";
    }
}

class CsvFormatter implements LogFormatter {
    @Override public String format(LogEntry e) {
        return e.timestamp + "," + e.level.name() + "," + e.message;
    }
}

interface LogDestination {
    void write(LogEntry entry);
    String getName();
}

class TestDest implements LogDestination {
    public List<String> received = new ArrayList<>();
    public LogFormatter fmt;
    public String name;
    public TestDest(String name, LogFormatter f) { this.name = name; this.fmt = f; }
    @Override public void write(LogEntry e) { received.add(fmt.format(e)); }
    @Override public String getName() { return name; }
}

class FailingDest implements LogDestination {
    @Override public void write(LogEntry e) { throw new RuntimeException("fail"); }
    @Override public String getName() { return "failing"; }
}

class Logger {
    private static final Logger INSTANCE = new Logger();

    private LogLevel minLevel = LogLevel.INFO;
    private List<String> history = new ArrayList<>();
    private List<LogDestination> destinations = new ArrayList<>();
    private final PlainTextFormatter defaultFormatter = new PlainTextFormatter();
    private LogFormatter formatter = defaultFormatter;

    private Logger() {}

    public static Logger getInstance() { return INSTANCE; }

    public void log(LogLevel level, String message) {
        if (level.rank < minLevel.rank) return;
        LogEntry entry = new LogEntry(level, message, "2024-01-15 10:30:00");
        for (LogDestination d : new ArrayList<>(destinations)) {
            try { d.write(entry); } catch (Throwable ignore) { /* fault isolation */ }
        }
        history.add(formatter.format(entry));
    }

    public void setLevel(LogLevel level) { this.minLevel = level; }

    public void setFormatter(LogFormatter f) { this.formatter = (f == null) ? defaultFormatter : f; }

    public void addDestination(LogDestination d) { destinations.add(d); }

    public void removeDestination(LogDestination d) { destinations.remove(d); }

    public LogLevel getLevel() { return minLevel; }

    public List<String> getLogHistory() { return new ArrayList<>(history); }

    public void clearHistory() {
        history.clear();
        destinations.clear();
    }
}

public class Solution {
    private static LogLevel levelFrom(String s) {
        if ("DEBUG".equals(s)) return LogLevel.DEBUG;
        if ("INFO".equals(s)) return LogLevel.INFO;
        if ("WARN".equals(s)) return LogLevel.WARN;
        if ("ERROR".equals(s)) return LogLevel.ERROR;
        return LogLevel.FATAL;
    }

    public static List<String> logger_simulate(List<LogOp> ops) {
        List<String> out = new ArrayList<>();
        Logger logger = Logger.getInstance();
        PlainTextFormatter plain = new PlainTextFormatter();
        JsonFormatter json = new JsonFormatter();
        CsvFormatter csv = new CsvFormatter();

        TestDest[] dests = new TestDest[8];
        FailingDest failing = null;

        for (LogOp op : ops) {
            String k = op.kind;
            if ("reset".equals(k)) {
                logger.clearHistory();
                logger.setLevel(LogLevel.INFO);
                logger.setFormatter(null);
                for (int i = 0; i < dests.length; i++) dests[i] = null;
                failing = null;
                out.add("ok");
            } else if ("set_level".equals(k)) {
                logger.setLevel(levelFrom(op.s1));
                out.add("ok");
            } else if ("set_formatter".equals(k)) {
                LogFormatter f = null;
                if ("json".equals(op.s1)) f = json;
                else if ("csv".equals(op.s1)) f = csv;
                else if ("plain".equals(op.s1)) f = plain;
                logger.setFormatter(f);
                out.add("ok");
            } else if ("log".equals(k)) {
                logger.log(levelFrom(op.s1), op.s2);
                out.add("ok");
            } else if ("history_size".equals(k)) {
                out.add(Integer.toString(logger.getLogHistory().size()));
            } else if ("history_contains".equals(k)) {
                List<String> h = logger.getLogHistory();
                if (op.i1 < 0 || op.i1 >= h.size()) out.add("no");
                else out.add(h.get(op.i1).contains(op.s1) ? "yes" : "no");
            } else if ("add_dest".equals(k)) {
                int a = op.s1.indexOf(':');
                int b = op.s1.indexOf(':', a + 1);
                int idx = Integer.parseInt(op.s1.substring(a + 1, b));
                String fmtName = op.s1.substring(b + 1);
                LogFormatter f = "json".equals(fmtName) ? json : "csv".equals(fmtName) ? csv : plain;
                dests[idx] = new TestDest("d" + idx, f);
                logger.addDestination(dests[idx]);
                out.add("ok");
            } else if ("rm_dest".equals(k)) {
                int idx = Integer.parseInt(op.s1.substring(op.s1.indexOf(':') + 1));
                if (dests[idx] != null) logger.removeDestination(dests[idx]);
                out.add("ok");
            } else if ("add_failing".equals(k)) {
                failing = new FailingDest();
                logger.addDestination(failing);
                out.add("ok");
            } else if ("rm_failing".equals(k)) {
                if (failing != null) logger.removeDestination(failing);
                out.add("ok");
            } else if ("dest_size".equals(k)) {
                int idx = Integer.parseInt(op.s1.substring(op.s1.indexOf(':') + 1));
                out.add(dests[idx] != null ? Integer.toString(dests[idx].received.size()) : "0");
            } else if ("dest_contains".equals(k)) {
                int idx = Integer.parseInt(op.s1.substring(op.s1.indexOf(':') + 1));
                if (dests[idx] == null) { out.add("no"); continue; }
                List<String> v = dests[idx].received;
                if (op.i1 < 0 || op.i1 >= v.size()) out.add("no");
                else out.add(v.get(op.i1).contains(op.s2) ? "yes" : "no");
            } else if ("fmt_plain".equals(k)) {
                out.add(plain.format(new LogEntry(levelFrom(op.s1), op.s2, "T")));
            } else if ("fmt_json".equals(k)) {
                out.add(json.format(new LogEntry(levelFrom(op.s1), op.s2, "T")));
            } else if ("fmt_csv".equals(k)) {
                out.add(csv.format(new LogEntry(levelFrom(op.s1), op.s2, "T")));
            } else if ("singleton_check".equals(k)) {
                out.add(Logger.getInstance() == Logger.getInstance() ? "yes" : "no");
            } else {
                out.add("unknown:" + k);
            }
        }
        return out;
    }
}

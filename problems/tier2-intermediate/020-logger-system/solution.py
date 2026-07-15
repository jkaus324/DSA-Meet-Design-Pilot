"""Logger system — singleton + level filter + pluggable formatters + destinations."""


class LogOp:
    def __init__(self, kind, s1="", s2="", i1=0):
        self.kind = kind
        self.s1 = s1
        self.s2 = s2
        self.i1 = i1


# Log levels (numeric ordering preserved)
DEBUG = 0
INFO = 1
WARN = 2
ERROR = 3
FATAL = 4

_LEVEL_NAMES = {DEBUG: "DEBUG", INFO: "INFO", WARN: "WARN", ERROR: "ERROR", FATAL: "FATAL"}


def _level_from(s):
    return {"DEBUG": DEBUG, "INFO": INFO, "WARN": WARN, "ERROR": ERROR, "FATAL": FATAL}.get(s, FATAL)


def _level_to_string(level):
    return _LEVEL_NAMES.get(level, "UNKNOWN")


class LogEntry:
    def __init__(self, level, message, timestamp):
        self.level = level
        self.message = message
        self.timestamp = timestamp


class PlainTextFormatter:
    def format(self, entry):
        return "[" + entry.timestamp + "] [" + _level_to_string(entry.level) + "] " + entry.message


class JsonFormatter:
    def format(self, entry):
        return ('{"timestamp":"' + entry.timestamp + '","level":"'
                + _level_to_string(entry.level) + '","message":"' + entry.message + '"}')


class CsvFormatter:
    def format(self, entry):
        return entry.timestamp + "," + _level_to_string(entry.level) + "," + entry.message


class TestDest:
    def __init__(self, name, formatter):
        self.name = name
        self.formatter = formatter
        self.received = []

    def write(self, entry):
        self.received.append(self.formatter.format(entry))

    def get_name(self):
        return self.name


class FailingDest:
    def write(self, entry):
        raise RuntimeError("fail")

    def get_name(self):
        return "failing"


class Logger:
    _instance = None

    def __init__(self):
        self.min_level = INFO
        self.history = []
        self.destinations = []
        self.default_formatter = PlainTextFormatter()
        self.formatter = self.default_formatter

    @classmethod
    def get_instance(cls):
        if cls._instance is None:
            cls._instance = Logger()
        return cls._instance

    def log(self, level, message):
        if level < self.min_level:
            return
        entry = LogEntry(level, message, "2024-01-15 10:30:00")
        for dest in self.destinations:
            try:
                dest.write(entry)
            except Exception:
                pass
        formatted = self.formatter.format(entry)
        self.history.append(formatted)

    def set_level(self, level):
        self.min_level = level

    def set_formatter(self, f):
        self.formatter = f if f is not None else self.default_formatter

    def add_destination(self, dest):
        self.destinations.append(dest)

    def remove_destination(self, dest):
        self.destinations = [d for d in self.destinations if d is not dest]

    def get_level(self):
        return self.min_level

    def get_log_history(self):
        return list(self.history)

    def clear_history(self):
        self.history = []
        self.destinations = []


def _fmt_for(s, plain, json, csv):
    if s == "json":
        return json
    if s == "csv":
        return csv
    if s == "plain":
        return plain
    return None


def logger_simulate(ops):
    out = []
    logger = Logger.get_instance()
    plain = PlainTextFormatter()
    json = JsonFormatter()
    csv = CsvFormatter()
    dests = [None] * 8
    failing = None

    for op in ops:
        k = op.kind
        if k == "reset":
            logger.clear_history()
            logger.set_level(INFO)
            logger.set_formatter(None)
            dests = [None] * 8
            failing = None
            out.append("ok")
        elif k == "set_level":
            logger.set_level(_level_from(op.s1))
            out.append("ok")
        elif k == "set_formatter":
            logger.set_formatter(_fmt_for(op.s1, plain, json, csv))
            out.append("ok")
        elif k == "log":
            logger.log(_level_from(op.s1), op.s2)
            out.append("ok")
        elif k == "history_size":
            out.append(str(len(logger.get_log_history())))
        elif k == "history_contains":
            h = logger.get_log_history()
            if op.i1 < 0 or op.i1 >= len(h):
                out.append("no")
            else:
                out.append("yes" if op.s1 in h[op.i1] else "no")
        elif k == "add_dest":
            # s1 = "test:<idx>:<formatter>"
            parts = op.s1.split(":")
            idx = int(parts[1])
            fmt_name = parts[2]
            dests[idx] = TestDest("d" + str(idx), _fmt_for(fmt_name, plain, json, csv))
            logger.add_destination(dests[idx])
            out.append("ok")
        elif k == "rm_dest":
            parts = op.s1.split(":")
            idx = int(parts[1])
            if dests[idx] is not None:
                logger.remove_destination(dests[idx])
            out.append("ok")
        elif k == "add_failing":
            failing = FailingDest()
            logger.add_destination(failing)
            out.append("ok")
        elif k == "rm_failing":
            if failing is not None:
                logger.remove_destination(failing)
            out.append("ok")
        elif k == "dest_size":
            parts = op.s1.split(":")
            idx = int(parts[1])
            out.append(str(len(dests[idx].received)) if dests[idx] is not None else "0")
        elif k == "dest_contains":
            parts = op.s1.split(":")
            idx = int(parts[1])
            if dests[idx] is None:
                out.append("no")
            else:
                v = dests[idx].received
                if op.i1 < 0 or op.i1 >= len(v):
                    out.append("no")
                else:
                    out.append("yes" if op.s2 in v[op.i1] else "no")
        elif k == "fmt_plain":
            e = LogEntry(_level_from(op.s1), op.s2, "T")
            out.append(plain.format(e))
        elif k == "fmt_json":
            e = LogEntry(_level_from(op.s1), op.s2, "T")
            out.append(json.format(e))
        elif k == "fmt_csv":
            e = LogEntry(_level_from(op.s1), op.s2, "T")
            out.append(csv.format(e))
        elif k == "singleton_check":
            a = Logger.get_instance()
            b = Logger.get_instance()
            out.append("yes" if a is b else "no")
        else:
            out.append("unknown:" + k)
    return out

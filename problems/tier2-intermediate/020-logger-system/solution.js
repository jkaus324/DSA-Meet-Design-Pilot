"use strict";
// Logger system — singleton + level filter + pluggable formatters + destinations.

class LogOp {
  constructor(kind, s1 = "", s2 = "", i1 = 0) {
    this.kind = kind;
    this.s1 = s1;
    this.s2 = s2;
    this.i1 = i1;
  }
}

// Log levels (numeric ordering preserved)
const DEBUG = 0;
const INFO = 1;
const WARN = 2;
const ERROR = 3;
const FATAL = 4;

const _LEVEL_NAMES = {
  [DEBUG]: "DEBUG",
  [INFO]: "INFO",
  [WARN]: "WARN",
  [ERROR]: "ERROR",
  [FATAL]: "FATAL",
};

function _level_from(s) {
  const m = { DEBUG, INFO, WARN, ERROR, FATAL };
  return Object.prototype.hasOwnProperty.call(m, s) ? m[s] : FATAL;
}

function _level_to_string(level) {
  return Object.prototype.hasOwnProperty.call(_LEVEL_NAMES, level)
    ? _LEVEL_NAMES[level]
    : "UNKNOWN";
}

class LogEntry {
  constructor(level, message, timestamp) {
    this.level = level;
    this.message = message;
    this.timestamp = timestamp;
  }
}

class PlainTextFormatter {
  format(entry) {
    return (
      "[" + entry.timestamp + "] [" + _level_to_string(entry.level) + "] " + entry.message
    );
  }
}

class JsonFormatter {
  format(entry) {
    return (
      '{"timestamp":"' +
      entry.timestamp +
      '","level":"' +
      _level_to_string(entry.level) +
      '","message":"' +
      entry.message +
      '"}'
    );
  }
}

class CsvFormatter {
  format(entry) {
    return entry.timestamp + "," + _level_to_string(entry.level) + "," + entry.message;
  }
}

class TestDest {
  constructor(name, formatter) {
    this.name = name;
    this.formatter = formatter;
    this.received = [];
  }
  write(entry) {
    this.received.push(this.formatter.format(entry));
  }
  get_name() {
    return this.name;
  }
}

class FailingDest {
  write(entry) {
    throw new Error("fail");
  }
  get_name() {
    return "failing";
  }
}

class Logger {
  constructor() {
    this.min_level = INFO;
    this.history = [];
    this.destinations = [];
    this.default_formatter = new PlainTextFormatter();
    this.formatter = this.default_formatter;
  }

  static get_instance() {
    if (Logger._instance === null) {
      Logger._instance = new Logger();
    }
    return Logger._instance;
  }

  log(level, message) {
    if (level < this.min_level) return;
    const entry = new LogEntry(level, message, "2024-01-15 10:30:00");
    for (const dest of this.destinations) {
      try {
        dest.write(entry);
      } catch (e) {
        // pass
      }
    }
    const formatted = this.formatter.format(entry);
    this.history.push(formatted);
  }

  set_level(level) {
    this.min_level = level;
  }

  set_formatter(f) {
    this.formatter = f !== null && f !== undefined ? f : this.default_formatter;
  }

  add_destination(dest) {
    this.destinations.push(dest);
  }

  remove_destination(dest) {
    this.destinations = this.destinations.filter((d) => d !== dest);
  }

  get_level() {
    return this.min_level;
  }

  get_log_history() {
    return this.history.slice();
  }

  clear_history() {
    this.history = [];
    this.destinations = [];
  }
}
Logger._instance = null;

function _fmt_for(s, plain, json, csv) {
  if (s === "json") return json;
  if (s === "csv") return csv;
  if (s === "plain") return plain;
  return null;
}

function logger_simulate(ops) {
  const out = [];
  const logger = Logger.get_instance();
  const plain = new PlainTextFormatter();
  const json = new JsonFormatter();
  const csv = new CsvFormatter();
  let dests = new Array(8).fill(null);
  let failing = null;

  for (const op of ops) {
    const k = op.kind;
    if (k === "reset") {
      logger.clear_history();
      logger.set_level(INFO);
      logger.set_formatter(null);
      dests = new Array(8).fill(null);
      failing = null;
      out.push("ok");
    } else if (k === "set_level") {
      logger.set_level(_level_from(op.s1));
      out.push("ok");
    } else if (k === "set_formatter") {
      logger.set_formatter(_fmt_for(op.s1, plain, json, csv));
      out.push("ok");
    } else if (k === "log") {
      logger.log(_level_from(op.s1), op.s2);
      out.push("ok");
    } else if (k === "history_size") {
      out.push(String(logger.get_log_history().length));
    } else if (k === "history_contains") {
      const h = logger.get_log_history();
      if (op.i1 < 0 || op.i1 >= h.length) {
        out.push("no");
      } else {
        out.push(h[op.i1].includes(op.s1) ? "yes" : "no");
      }
    } else if (k === "add_dest") {
      // s1 = "test:<idx>:<formatter>"
      const parts = op.s1.split(":");
      const idx = parseInt(parts[1], 10);
      const fmtName = parts[2];
      dests[idx] = new TestDest("d" + String(idx), _fmt_for(fmtName, plain, json, csv));
      logger.add_destination(dests[idx]);
      out.push("ok");
    } else if (k === "rm_dest") {
      const parts = op.s1.split(":");
      const idx = parseInt(parts[1], 10);
      if (dests[idx] !== null) {
        logger.remove_destination(dests[idx]);
      }
      out.push("ok");
    } else if (k === "add_failing") {
      failing = new FailingDest();
      logger.add_destination(failing);
      out.push("ok");
    } else if (k === "rm_failing") {
      if (failing !== null) {
        logger.remove_destination(failing);
      }
      out.push("ok");
    } else if (k === "dest_size") {
      const parts = op.s1.split(":");
      const idx = parseInt(parts[1], 10);
      out.push(dests[idx] !== null ? String(dests[idx].received.length) : "0");
    } else if (k === "dest_contains") {
      const parts = op.s1.split(":");
      const idx = parseInt(parts[1], 10);
      if (dests[idx] === null) {
        out.push("no");
      } else {
        const v = dests[idx].received;
        if (op.i1 < 0 || op.i1 >= v.length) {
          out.push("no");
        } else {
          out.push(v[op.i1].includes(op.s2) ? "yes" : "no");
        }
      }
    } else if (k === "fmt_plain") {
      const e = new LogEntry(_level_from(op.s1), op.s2, "T");
      out.push(plain.format(e));
    } else if (k === "fmt_json") {
      const e = new LogEntry(_level_from(op.s1), op.s2, "T");
      out.push(json.format(e));
    } else if (k === "fmt_csv") {
      const e = new LogEntry(_level_from(op.s1), op.s2, "T");
      out.push(csv.format(e));
    } else if (k === "singleton_check") {
      const a = Logger.get_instance();
      const b = Logger.get_instance();
      out.push(a === b ? "yes" : "no");
    } else {
      out.push("unknown:" + k);
    }
  }
  return out;
}

module.exports = {
  LogOp,
  LogEntry,
  PlainTextFormatter,
  JsonFormatter,
  CsvFormatter,
  Logger,
  logger_simulate,
};

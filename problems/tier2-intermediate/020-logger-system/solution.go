// Logger system — singleton + level filter + pluggable formatters + destinations (Go port).
package main

import (
	"strconv"
	"strings"
)

type LogOp struct {
	kind string
	s1   string
	s2   string
	i1   int
}

const (
	lvlDebug = 0
	lvlInfo  = 1
	lvlWarn  = 2
	lvlError = 3
	lvlFatal = 4
)

func levelFrom(s string) int {
	switch s {
	case "DEBUG":
		return lvlDebug
	case "INFO":
		return lvlInfo
	case "WARN":
		return lvlWarn
	case "ERROR":
		return lvlError
	case "FATAL":
		return lvlFatal
	default:
		return lvlFatal
	}
}

func levelToString(level int) string {
	switch level {
	case lvlDebug:
		return "DEBUG"
	case lvlInfo:
		return "INFO"
	case lvlWarn:
		return "WARN"
	case lvlError:
		return "ERROR"
	case lvlFatal:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

type logEntry struct {
	level     int
	message   string
	timestamp string
}

type formatter interface {
	format(e logEntry) string
}

type plainTextFormatter struct{}

func (plainTextFormatter) format(e logEntry) string {
	return "[" + e.timestamp + "] [" + levelToString(e.level) + "] " + e.message
}

type jsonFormatter struct{}

func (jsonFormatter) format(e logEntry) string {
	return `{"timestamp":"` + e.timestamp + `","level":"` + levelToString(e.level) + `","message":"` + e.message + `"}`
}

type csvFormatter struct{}

func (csvFormatter) format(e logEntry) string {
	return e.timestamp + "," + levelToString(e.level) + "," + e.message
}

type destination interface {
	write(e logEntry)
	getName() string
}

type testDest struct {
	name      string
	formatter formatter
	received  []string
}

func (d *testDest) write(e logEntry) {
	d.received = append(d.received, d.formatter.format(e))
}

func (d *testDest) getName() string { return d.name }

type failingDest struct{}

func (failingDest) write(e logEntry) { panic("fail") }

func (failingDest) getName() string { return "failing" }

type logger struct {
	minLevel         int
	history          []string
	destinations     []destination
	defaultFormatter formatter
	formatter        formatter
}

func newLogger() *logger {
	def := plainTextFormatter{}
	return &logger{
		minLevel:         lvlInfo,
		history:          []string{},
		destinations:     []destination{},
		defaultFormatter: def,
		formatter:        def,
	}
}

var loggerInstance *logger

func getLoggerInstance() *logger {
	if loggerInstance == nil {
		loggerInstance = newLogger()
	}
	return loggerInstance
}

func (l *logger) log(level int, message string) {
	if level < l.minLevel {
		return
	}
	entry := logEntry{level, message, "2024-01-15 10:30:00"}
	for _, dest := range l.destinations {
		func() {
			defer func() { _ = recover() }()
			dest.write(entry)
		}()
	}
	l.history = append(l.history, l.formatter.format(entry))
}

func (l *logger) setLevel(level int) { l.minLevel = level }

func (l *logger) setFormatter(f formatter) {
	if f != nil {
		l.formatter = f
	} else {
		l.formatter = l.defaultFormatter
	}
}

func (l *logger) addDestination(d destination) {
	l.destinations = append(l.destinations, d)
}

func (l *logger) removeDestination(d destination) {
	kept := []destination{}
	for _, x := range l.destinations {
		if x != d {
			kept = append(kept, x)
		}
	}
	l.destinations = kept
}

func (l *logger) getLogHistory() []string { return l.history }

func (l *logger) clearHistory() {
	l.history = []string{}
	l.destinations = []destination{}
}

func fmtFor(s string, plain, jsonF, csv formatter) formatter {
	switch s {
	case "json":
		return jsonF
	case "csv":
		return csv
	case "plain":
		return plain
	default:
		return nil
	}
}

func logger_simulate(ops []LogOp) []string {
	out := []string{}
	loggerInstance = nil
	lg := getLoggerInstance()
	plain := plainTextFormatter{}
	jsonF := jsonFormatter{}
	csv := csvFormatter{}
	dests := make([]*testDest, 8)
	var failing *failingDest

	for _, op := range ops {
		switch op.kind {
		case "reset":
			lg.clearHistory()
			lg.setLevel(lvlInfo)
			lg.setFormatter(nil)
			dests = make([]*testDest, 8)
			failing = nil
			out = append(out, "ok")
		case "set_level":
			lg.setLevel(levelFrom(op.s1))
			out = append(out, "ok")
		case "set_formatter":
			lg.setFormatter(fmtFor(op.s1, plain, jsonF, csv))
			out = append(out, "ok")
		case "log":
			lg.log(levelFrom(op.s1), op.s2)
			out = append(out, "ok")
		case "history_size":
			out = append(out, strconv.Itoa(len(lg.getLogHistory())))
		case "history_contains":
			h := lg.getLogHistory()
			if op.i1 < 0 || op.i1 >= len(h) {
				out = append(out, "no")
			} else if strings.Contains(h[op.i1], op.s1) {
				out = append(out, "yes")
			} else {
				out = append(out, "no")
			}
		case "add_dest":
			parts := strings.Split(op.s1, ":")
			idx, _ := strconv.Atoi(parts[1])
			fmtName := parts[2]
			dests[idx] = &testDest{name: "d" + strconv.Itoa(idx), formatter: fmtFor(fmtName, plain, jsonF, csv)}
			lg.addDestination(dests[idx])
			out = append(out, "ok")
		case "rm_dest":
			parts := strings.Split(op.s1, ":")
			idx, _ := strconv.Atoi(parts[1])
			if dests[idx] != nil {
				lg.removeDestination(dests[idx])
			}
			out = append(out, "ok")
		case "add_failing":
			failing = &failingDest{}
			lg.addDestination(failing)
			out = append(out, "ok")
		case "rm_failing":
			if failing != nil {
				lg.removeDestination(failing)
			}
			out = append(out, "ok")
		case "dest_size":
			parts := strings.Split(op.s1, ":")
			idx, _ := strconv.Atoi(parts[1])
			if dests[idx] != nil {
				out = append(out, strconv.Itoa(len(dests[idx].received)))
			} else {
				out = append(out, "0")
			}
		case "dest_contains":
			parts := strings.Split(op.s1, ":")
			idx, _ := strconv.Atoi(parts[1])
			if dests[idx] == nil {
				out = append(out, "no")
			} else {
				v := dests[idx].received
				if op.i1 < 0 || op.i1 >= len(v) {
					out = append(out, "no")
				} else if strings.Contains(v[op.i1], op.s2) {
					out = append(out, "yes")
				} else {
					out = append(out, "no")
				}
			}
		case "fmt_plain":
			e := logEntry{levelFrom(op.s1), op.s2, "T"}
			out = append(out, plain.format(e))
		case "fmt_json":
			e := logEntry{levelFrom(op.s1), op.s2, "T"}
			out = append(out, jsonF.format(e))
		case "fmt_csv":
			e := logEntry{levelFrom(op.s1), op.s2, "T"}
			out = append(out, csv.format(e))
		case "singleton_check":
			a := getLoggerInstance()
			b := getLoggerInstance()
			if a == b {
				out = append(out, "yes")
			} else {
				out = append(out, "no")
			}
		default:
			out = append(out, "unknown:"+op.kind)
		}
	}
	return out
}

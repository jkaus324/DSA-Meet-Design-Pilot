package main

import (
	"fmt"
	"strings"
)

// TestDestination stores formatted log lines in memory (for test assertions).
type TestDestination struct {
	formatter LogFormatter
	destName  string
	Received  []string
}

func NewTestDestination(name string, f LogFormatter) *TestDestination {
	return &TestDestination{destName: name, formatter: f}
}

func (d *TestDestination) Write(entry LogEntry) {
	d.Received = append(d.Received, d.formatter.Format(entry))
}

func (d *TestDestination) GetName() string { return d.destName }

// FailingDestination always panics in Write (tests fault isolation).
type FailingDestination struct{}

func (d *FailingDestination) Write(entry LogEntry) { panic("destination failed") }
func (d *FailingDestination) GetName() string      { return "failing" }

func part3Tests() int {
	passed := 0
	failed := 0

	test := func(name string, fn func()) {
		func() {
			defer func() {
				if r := recover(); r != nil {
					fmt.Println("FAIL", name)
					failed++
				}
			}()
			fn()
			fmt.Println("PASS", name)
			passed++
		}()
	}

	// Test 1: Register a single destination and log
	test("test_single_destination", func() {
		logger := GetInstance()
		logger.ClearHistory()
		logger = GetInstance()
		logger.SetLevel(DEBUG)

		plainFmt := &PlainTextFormatter{}
		dest1 := NewTestDestination("test1", plainFmt)
		logger.AddDestination(dest1)

		logger.Log(INFO, "hello dest")
		if len(dest1.Received) != 1 {
			panic("expected 1 received entry")
		}
		if !strings.Contains(dest1.Received[0], "INFO") {
			panic("expected INFO in received entry")
		}
		if !strings.Contains(dest1.Received[0], "hello dest") {
			panic("expected message in received entry")
		}
		logger.RemoveDestination(dest1)
	})

	// Test 2: Multiple destinations receive same log entry
	test("test_multiple_destinations", func() {
		logger := GetInstance()
		logger.ClearHistory()
		logger = GetInstance()
		logger.SetLevel(DEBUG)

		plainFmt := &PlainTextFormatter{}
		jsonFmt := &JsonFormatter{}
		dest1 := NewTestDestination("plain-dest", plainFmt)
		dest2 := NewTestDestination("json-dest", jsonFmt)
		logger.AddDestination(dest1)
		logger.AddDestination(dest2)

		logger.Log(ERROR, "multi dest")

		if len(dest1.Received) != 1 {
			panic("expected dest1 to receive 1 entry")
		}
		if len(dest2.Received) != 1 {
			panic("expected dest2 to receive 1 entry")
		}
		if !strings.Contains(dest1.Received[0], "[ERROR]") {
			panic("expected [ERROR] in plain dest")
		}
		if !strings.Contains(dest2.Received[0], `"level"`) {
			panic("expected JSON level key in json dest")
		}
		logger.RemoveDestination(dest1)
		logger.RemoveDestination(dest2)
	})

	// Test 3: Each destination has independent formatter
	test("test_independent_formatters", func() {
		logger := GetInstance()
		logger.ClearHistory()
		logger = GetInstance()
		logger.SetLevel(DEBUG)

		plainFmt := &PlainTextFormatter{}
		csvFmt := &CsvFormatter{}
		destPlain := NewTestDestination("plain", plainFmt)
		destCsv := NewTestDestination("csv", csvFmt)
		logger.AddDestination(destPlain)
		logger.AddDestination(destCsv)

		logger.Log(WARN, "format test")

		if !strings.Contains(destPlain.Received[0], "[WARN]") {
			panic("expected [WARN] in plain destination")
		}
		if strings.Contains(destCsv.Received[0], "[") {
			panic("CSV should not have brackets")
		}
		if !strings.Contains(destCsv.Received[0], ",") {
			panic("expected comma in CSV destination")
		}
		logger.RemoveDestination(destPlain)
		logger.RemoveDestination(destCsv)
	})

	// Test 4: Failing destination does not block others
	test("test_fault_isolation", func() {
		logger := GetInstance()
		logger.ClearHistory()
		logger = GetInstance()
		logger.SetLevel(DEBUG)

		plainFmt := &PlainTextFormatter{}
		failDest := &FailingDestination{}
		goodDest := NewTestDestination("good", plainFmt)
		logger.AddDestination(failDest)
		logger.AddDestination(goodDest)

		logger.Log(ERROR, "fault isolation")

		if len(goodDest.Received) != 1 {
			panic("expected goodDest to receive 1 entry despite failing destination")
		}
		if !strings.Contains(goodDest.Received[0], "fault isolation") {
			panic("expected message in goodDest")
		}
		logger.RemoveDestination(failDest)
		logger.RemoveDestination(goodDest)
	})

	// Test 5: Remove destination stops it from receiving logs
	test("test_remove_destination", func() {
		logger := GetInstance()
		logger.ClearHistory()
		logger = GetInstance()
		logger.SetLevel(DEBUG)

		plainFmt := &PlainTextFormatter{}
		dest := NewTestDestination("removable", plainFmt)
		logger.AddDestination(dest)

		logger.Log(INFO, "before remove")
		if len(dest.Received) != 1 {
			panic("expected 1 entry before remove")
		}

		logger.RemoveDestination(dest)
		logger.Log(INFO, "after remove")
		if len(dest.Received) != 1 {
			panic("expected still 1 entry after remove")
		}
	})

	// Test 6: Level filtering applies before destinations receive
	test("test_level_filtering_before_destinations", func() {
		logger := GetInstance()
		logger.ClearHistory()
		logger = GetInstance()
		logger.SetLevel(ERROR)

		plainFmt := &PlainTextFormatter{}
		dest := NewTestDestination("level-test", plainFmt)
		logger.AddDestination(dest)

		logger.Log(INFO, "filtered")
		logger.Log(ERROR, "passes")
		if len(dest.Received) != 1 {
			panic("expected 1 entry (INFO filtered)")
		}
		if !strings.Contains(dest.Received[0], "passes") {
			panic("expected passing message")
		}
		logger.RemoveDestination(dest)
	})

	// Test 7: No destinations registered — log still works (no crash)
	test("test_no_destinations_no_crash", func() {
		logger := GetInstance()
		logger.ClearHistory()
		logger = GetInstance()
		logger.SetLevel(DEBUG)
		// No destinations added
		logger.Log(INFO, "no destinations") // should not panic
	})

	fmt.Printf("PART3_SUMMARY %d/%d\n", passed, passed+failed)
	return failed
}

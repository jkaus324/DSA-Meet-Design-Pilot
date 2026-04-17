package main

import (
	"fmt"
	"strings"
)

func part2Tests() int {
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

	// Test 1: PlainTextFormatter output format
	test("test_plaintext_format", func() {
		fmt := &PlainTextFormatter{}
		entry := LogEntry{Level: ERROR, Message: "disk full", Timestamp: "2024-01-15 10:30:00"}
		result := fmt.Format(entry)
		if !strings.Contains(result, "[2024-01-15 10:30:00]") {
			panic("expected timestamp in brackets")
		}
		if !strings.Contains(result, "[ERROR]") {
			panic("expected ERROR in brackets")
		}
		if !strings.Contains(result, "disk full") {
			panic("expected message in result")
		}
	})

	// Test 2: JsonFormatter output format
	test("test_json_format", func() {
		jfmt := &JsonFormatter{}
		entry := LogEntry{Level: WARN, Message: "low memory", Timestamp: "2024-01-15 10:30:00"}
		result := jfmt.Format(entry)
		if !strings.Contains(result, `"timestamp"`) {
			panic("expected timestamp key")
		}
		if !strings.Contains(result, `"level"`) {
			panic("expected level key")
		}
		if !strings.Contains(result, `"message"`) {
			panic("expected message key")
		}
		if !strings.Contains(result, "WARN") {
			panic("expected WARN in result")
		}
		if !strings.Contains(result, "low memory") {
			panic("expected message content")
		}
		if !strings.Contains(result, "2024-01-15 10:30:00") {
			panic("expected timestamp value")
		}
	})

	// Test 3: CsvFormatter output format
	test("test_csv_format", func() {
		cfmt := &CsvFormatter{}
		entry := LogEntry{Level: INFO, Message: "server started", Timestamp: "2024-01-15 10:30:00"}
		result := cfmt.Format(entry)
		if !strings.Contains(result, ",") {
			panic("expected comma in CSV")
		}
		if !strings.Contains(result, "INFO") {
			panic("expected INFO in result")
		}
		if !strings.Contains(result, "server started") {
			panic("expected message content")
		}
		if !strings.Contains(result, "2024-01-15 10:30:00") {
			panic("expected timestamp")
		}
		if strings.Contains(result, "[") {
			panic("CSV should not have brackets")
		}
	})

	// Test 4: Logger uses PlainText by default
	test("test_default_plaintext_formatter", func() {
		logger := GetInstance()
		logger.ClearHistory()
		logger = GetInstance()
		logger.SetLevel(DEBUG)
		logger.SetFormatter(nil) // reset to default
		logger.Log(INFO, "default format")
		entry := logger.GetLogHistory()[0]
		if !strings.Contains(entry, "[") {
			panic("expected brackets (PlainText format)")
		}
		if !strings.Contains(entry, "[INFO]") {
			panic("expected [INFO] in entry")
		}
	})

	// Test 5: Switch to JSON formatter at runtime
	test("test_switch_to_json_formatter", func() {
		logger := GetInstance()
		logger.ClearHistory()
		logger = GetInstance()
		logger.SetLevel(DEBUG)
		jsonFmt := &JsonFormatter{}
		logger.SetFormatter(jsonFmt)
		logger.Log(ERROR, "json test")
		entry := logger.GetLogHistory()[0]
		if !strings.Contains(entry, `"level"`) {
			panic("expected JSON level key")
		}
		if !strings.Contains(entry, `"message"`) {
			panic("expected JSON message key")
		}
		if !strings.Contains(entry, "json test") {
			panic("expected message content")
		}
		logger.SetFormatter(nil) // reset
	})

	// Test 6: Switch to CSV formatter at runtime
	test("test_switch_to_csv_formatter", func() {
		logger := GetInstance()
		logger.ClearHistory()
		logger = GetInstance()
		logger.SetLevel(DEBUG)
		csvFmt := &CsvFormatter{}
		logger.SetFormatter(csvFmt)
		logger.Log(WARN, "csv test")
		entry := logger.GetLogHistory()[0]
		if !strings.Contains(entry, ",") {
			panic("expected comma in CSV")
		}
		if !strings.Contains(entry, "WARN") {
			panic("expected WARN in CSV")
		}
		if strings.Contains(entry, "[") {
			panic("CSV should not have brackets")
		}
		logger.SetFormatter(nil) // reset
	})

	// Test 7: Level filtering still works with custom formatter
	test("test_level_filtering_with_custom_formatter", func() {
		logger := GetInstance()
		logger.ClearHistory()
		logger = GetInstance()
		logger.SetLevel(ERROR)
		jsonFmt := &JsonFormatter{}
		logger.SetFormatter(jsonFmt)
		logger.Log(INFO, "should be filtered")
		logger.Log(ERROR, "should pass")
		history := logger.GetLogHistory()
		if len(history) != 1 {
			panic("expected 1 log entry")
		}
		if !strings.Contains(history[0], "should pass") {
			panic("expected passing message in history")
		}
		logger.SetFormatter(nil) // reset
	})

	fmt.Printf("PART2_SUMMARY %d/%d\n", passed, passed+failed)
	return failed
}

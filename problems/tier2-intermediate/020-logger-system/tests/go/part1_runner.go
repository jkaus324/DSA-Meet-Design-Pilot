package main

import (
	"fmt"
	"strings"
)

func part1Tests() int {
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

	// Test 1: Singleton — GetInstance returns same pointer
	test("test_singleton_same_instance", func() {
		logger1 := GetInstance()
		logger2 := GetInstance()
		if logger1 != logger2 {
			panic("expected same instance")
		}
		logger1.ClearHistory()
	})

	// Test 2: Default level is INFO — DEBUG messages are filtered
	test("test_debug_filtered_at_info_level", func() {
		logger := GetInstance()
		logger.ClearHistory()
		logger = GetInstance()
		logger.SetLevel(INFO)
		logger.Log(DEBUG, "debug msg")
		if len(logger.GetLogHistory()) != 0 {
			panic("expected DEBUG filtered at INFO level")
		}
	})

	// Test 3: INFO message passes at INFO level
	test("test_info_passes_at_info_level", func() {
		logger := GetInstance()
		logger.ClearHistory()
		logger = GetInstance()
		logger.SetLevel(INFO)
		logger.Log(INFO, "info msg")
		history := logger.GetLogHistory()
		if len(history) != 1 {
			panic("expected 1 log entry")
		}
		if !strings.Contains(history[0], "INFO") {
			panic("expected INFO in log entry")
		}
		if !strings.Contains(history[0], "info msg") {
			panic("expected message in log entry")
		}
	})

	// Test 4: ERROR and FATAL pass at INFO level
	test("test_error_fatal_pass_at_info_level", func() {
		logger := GetInstance()
		logger.ClearHistory()
		logger = GetInstance()
		logger.SetLevel(INFO)
		logger.Log(ERROR, "error msg")
		logger.Log(FATAL, "fatal msg")
		history := logger.GetLogHistory()
		if len(history) != 2 {
			panic("expected 2 log entries")
		}
		if !strings.Contains(history[0], "ERROR") {
			panic("expected ERROR in first entry")
		}
		if !strings.Contains(history[1], "FATAL") {
			panic("expected FATAL in second entry")
		}
	})

	// Test 5: Changing level to WARN filters out INFO
	test("test_warn_level_filters_info", func() {
		logger := GetInstance()
		logger.ClearHistory()
		logger = GetInstance()
		logger.SetLevel(WARN)
		logger.Log(INFO, "should be filtered")
		logger.Log(WARN, "should pass")
		logger.Log(ERROR, "should also pass")
		history := logger.GetLogHistory()
		if len(history) != 2 {
			panic("expected 2 log entries")
		}
		if !strings.Contains(history[0], "WARN") {
			panic("expected WARN in first entry")
		}
		if !strings.Contains(history[1], "ERROR") {
			panic("expected ERROR in second entry")
		}
	})

	// Test 6: SetLevel to DEBUG allows all messages
	test("test_debug_level_allows_all", func() {
		logger := GetInstance()
		logger.ClearHistory()
		logger = GetInstance()
		logger.SetLevel(DEBUG)
		logger.Log(DEBUG, "d")
		logger.Log(INFO, "i")
		logger.Log(WARN, "w")
		logger.Log(ERROR, "e")
		logger.Log(FATAL, "f")
		if len(logger.GetLogHistory()) != 5 {
			panic("expected 5 log entries at DEBUG level")
		}
	})

	// Test 7: SetLevel to FATAL only allows FATAL
	test("test_fatal_level_only_fatal", func() {
		logger := GetInstance()
		logger.ClearHistory()
		logger = GetInstance()
		logger.SetLevel(FATAL)
		logger.Log(DEBUG, "d")
		logger.Log(INFO, "i")
		logger.Log(WARN, "w")
		logger.Log(ERROR, "e")
		logger.Log(FATAL, "f")
		history := logger.GetLogHistory()
		if len(history) != 1 {
			panic("expected 1 log entry at FATAL level")
		}
		if !strings.Contains(history[0], "FATAL") {
			panic("expected FATAL in log entry")
		}
	})

	// Test 8: Log format contains timestamp brackets
	test("test_log_format_has_brackets", func() {
		logger := GetInstance()
		logger.ClearHistory()
		logger = GetInstance()
		logger.SetLevel(INFO)
		logger.Log(INFO, "check format")
		entry := logger.GetLogHistory()[0]
		if !strings.Contains(entry, "[") {
			panic("expected brackets in log entry")
		}
		if !strings.Contains(entry, "]") {
			panic("expected closing bracket in log entry")
		}
		if !strings.Contains(entry, "check format") {
			panic("expected message in log entry")
		}
	})

	fmt.Printf("PART1_SUMMARY %d/%d\n", passed, passed+failed)
	return failed
}

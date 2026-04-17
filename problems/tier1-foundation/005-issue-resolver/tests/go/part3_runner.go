package main

import "fmt"

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

	// Test 1: valid state transition OPEN -> IN_PROGRESS
	test("test_valid_transition_open_to_inprogress", func() {
		issues := []Issue{
			{200, "Bug", CategoryTechnical, PriorityHigh, IssueStateOpen, 0},
		}
		notifications := []string{}
		ok := TransitionIssue(&issues, 200, IssueStateInProgress, &notifications)
		if !ok {
			panic("expected transition to succeed")
		}
		if issues[0].State != IssueStateInProgress {
			panic("expected state to be IN_PROGRESS")
		}
		if len(notifications) != 1 {
			panic(fmt.Sprintf("expected 1 notification, got %d", len(notifications)))
		}
		if notifications[0] != "Issue 200: OPEN -> IN_PROGRESS" {
			panic(fmt.Sprintf("unexpected notification: %s", notifications[0]))
		}
	})

	// Test 2: invalid state transition OPEN -> CLOSED is rejected
	test("test_invalid_transition_rejected", func() {
		issues := []Issue{
			{201, "Bug", CategoryTechnical, PriorityHigh, IssueStateOpen, 0},
		}
		notifications := []string{}
		ok := TransitionIssue(&issues, 201, IssueStateClosed, &notifications)
		if ok {
			panic("expected transition to fail")
		}
		if issues[0].State != IssueStateOpen {
			panic("expected state to remain OPEN")
		}
		if len(notifications) != 0 {
			panic("expected no notifications for invalid transition")
		}
	})

	// Test 3: full lifecycle OPEN -> IN_PROGRESS -> RESOLVED -> CLOSED
	test("test_full_lifecycle", func() {
		issues := []Issue{
			{202, "Payment fail", CategoryBilling, PriorityCritical, IssueStateOpen, 1},
		}
		notifications := []string{}
		if !TransitionIssue(&issues, 202, IssueStateInProgress, &notifications) {
			panic("expected OPEN->IN_PROGRESS to succeed")
		}
		if !TransitionIssue(&issues, 202, IssueStateResolved, &notifications) {
			panic("expected IN_PROGRESS->RESOLVED to succeed")
		}
		if !TransitionIssue(&issues, 202, IssueStateClosed, &notifications) {
			panic("expected RESOLVED->CLOSED to succeed")
		}
		if issues[0].State != IssueStateClosed {
			panic("expected final state to be CLOSED")
		}
		if len(notifications) != 3 {
			panic(fmt.Sprintf("expected 3 notifications, got %d", len(notifications)))
		}
	})

	// Test 4: AssignNextPriority picks highest priority first
	test("test_priority_highest_first", func() {
		agents := []Agent{
			{0, "Alice", 0, nil},
			{1, "Bob", 0, nil},
		}
		issues := []Issue{
			{300, "Low prio", CategoryGeneral, PriorityLow, IssueStateOpen, -1},
			{301, "Critical", CategoryBilling, PriorityCritical, IssueStateOpen, -1},
			{302, "Medium", CategoryTechnical, PriorityMedium, IssueStateOpen, -1},
		}
		first := AssignNextPriority(agents, &issues)
		if first.ID != 301 {
			panic(fmt.Sprintf("expected issue 301 (CRITICAL), got %d", first.ID))
		}
		if first.AssignedAgentID == -1 {
			panic("expected issue to be assigned to an agent")
		}
	})

	// Test 5: priority tie broken by lowest issue ID
	test("test_priority_tiebreak_by_id", func() {
		agents := []Agent{
			{0, "Alice", 0, nil},
		}
		issues := []Issue{
			{403, "Third", CategoryGeneral, PriorityHigh, IssueStateOpen, -1},
			{401, "First", CategoryGeneral, PriorityHigh, IssueStateOpen, -1},
			{402, "Second", CategoryGeneral, PriorityHigh, IssueStateOpen, -1},
		}
		first := AssignNextPriority(agents, &issues)
		if first.ID != 401 {
			panic(fmt.Sprintf("expected issue 401 (lowest ID among HIGH), got %d", first.ID))
		}
	})

	// Test 6: AssignNextPriority skips already-assigned issues
	test("test_priority_skips_assigned", func() {
		agents := []Agent{
			{0, "Alice", 0, nil},
		}
		issues := []Issue{
			{500, "Assigned", CategoryGeneral, PriorityCritical, IssueStateOpen, 1},   // already assigned
			{501, "Unassigned", CategoryGeneral, PriorityLow, IssueStateOpen, -1},     // unassigned
		}
		result := AssignNextPriority(agents, &issues)
		if result.ID != 501 {
			panic(fmt.Sprintf("expected issue 501 (unassigned), got %d", result.ID))
		}
	})

	fmt.Printf("PART3_SUMMARY %d/%d\n", passed, passed+failed)
	return failed
}

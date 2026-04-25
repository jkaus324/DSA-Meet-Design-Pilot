package main

import "fmt"

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

	// Test 1: least-loaded assigns to agent with fewest issues
	test("test_least_loaded_basic", func() {
		agents := []Agent{
			{0, "Alice", 3, nil},
			{1, "Bob", 1, nil},
			{2, "Carol", 2, nil},
		}
		issues := []Issue{}
		result := AssignLeastLoaded(agents, &issues,
			Issue{100, "Help", CategoryGeneral, PriorityLow, IssueStateOpen, -1})
		if result.AssignedAgentID != 1 {
			panic(fmt.Sprintf("expected agent 1 (Bob, lowest load), got %d", result.AssignedAgentID))
		}
	})

	// Test 2: least-loaded breaks ties by lowest agent ID
	test("test_least_loaded_tiebreak", func() {
		agents := []Agent{
			{0, "Alice", 2, nil},
			{1, "Bob", 2, nil},
			{2, "Carol", 2, nil},
		}
		issues := []Issue{}
		result := AssignLeastLoaded(agents, &issues,
			Issue{101, "Tie", CategoryGeneral, PriorityLow, IssueStateOpen, -1})
		if result.AssignedAgentID != 0 {
			panic(fmt.Sprintf("expected agent 0 (Alice, tiebreak by ID), got %d", result.AssignedAgentID))
		}
	})

	// Test 3: specialist assigns to agent with matching category
	test("test_specialist_match", func() {
		agents := []Agent{
			{0, "Alice", 0, []Category{CategoryBilling}},
			{1, "Bob", 0, []Category{CategoryTechnical}},
			{2, "Carol", 0, []Category{CategoryGeneral}},
		}
		issues := []Issue{}
		result := AssignBySpecialist(agents, &issues,
			Issue{102, "Tech issue", CategoryTechnical, PriorityHigh, IssueStateOpen, -1})
		if result.AssignedAgentID != 1 {
			panic(fmt.Sprintf("expected agent 1 (Bob, TECHNICAL specialist), got %d", result.AssignedAgentID))
		}
	})

	// Test 4: specialist falls back to least-loaded when no specialist exists
	test("test_specialist_fallback", func() {
		agents := []Agent{
			{0, "Alice", 3, []Category{CategoryBilling}},
			{1, "Bob", 1, []Category{CategoryBilling}},
			{2, "Carol", 2, []Category{CategoryBilling}},
		}
		issues := []Issue{}
		result := AssignBySpecialist(agents, &issues,
			Issue{103, "Account issue", CategoryAccount, PriorityMedium, IssueStateOpen, -1})
		// No ACCOUNT specialist — falls back to least-loaded: Bob (load=1)
		if result.AssignedAgentID != 1 {
			panic(fmt.Sprintf("expected agent 1 (Bob, fallback least-loaded), got %d", result.AssignedAgentID))
		}
	})

	// Test 5: specialist picks least-loaded among multiple specialists
	test("test_specialist_least_loaded", func() {
		agents := []Agent{
			{0, "Alice", 5, []Category{CategoryBilling, CategoryTechnical}},
			{1, "Bob", 2, []Category{CategoryTechnical}},
			{2, "Carol", 8, []Category{CategoryTechnical, CategoryGeneral}},
		}
		issues := []Issue{}
		result := AssignBySpecialist(agents, &issues,
			Issue{104, "Server down", CategoryTechnical, PriorityCritical, IssueStateOpen, -1})
		if result.AssignedAgentID != 1 {
			panic(fmt.Sprintf("expected agent 1 (Bob, least-loaded TECHNICAL specialist), got %d", result.AssignedAgentID))
		}
	})

	fmt.Printf("PART2_SUMMARY %d/%d\n", passed, passed+failed)
	return failed
}

package main

import "fmt"

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

	// Test 1: round-robin assigns to agents in order
	test("test_round_robin_assignment", func() {
		agents := []Agent{
			{0, "Alice", 0, nil},
			{1, "Bob", 0, nil},
			{2, "Carol", 0, nil},
		}
		issues := []Issue{}

		i1 := AssignIssue(agents, &issues, Issue{1, "Issue A", CategoryGeneral, PriorityLow, IssueStateOpen, -1})
		i2 := AssignIssue(agents, &issues, Issue{2, "Issue B", CategoryGeneral, PriorityLow, IssueStateOpen, -1})
		i3 := AssignIssue(agents, &issues, Issue{3, "Issue C", CategoryGeneral, PriorityLow, IssueStateOpen, -1})

		if i1.AssignedAgentID != 0 {
			panic(fmt.Sprintf("expected agent 0, got %d", i1.AssignedAgentID))
		}
		if i2.AssignedAgentID != 1 {
			panic(fmt.Sprintf("expected agent 1, got %d", i2.AssignedAgentID))
		}
		if i3.AssignedAgentID != 2 {
			panic(fmt.Sprintf("expected agent 2, got %d", i3.AssignedAgentID))
		}
	})

	// Test 2: round-robin wraps around
	test("test_round_robin_wrap", func() {
		agents := []Agent{
			{0, "Alice", 0, nil},
			{1, "Bob", 0, nil},
		}
		issues := []Issue{}
		rr := &RoundRobinStrategy{}
		resolver := NewIssueResolver(rr)

		i1 := resolver.Assign(agents, &issues, Issue{10, "A", CategoryGeneral, PriorityLow, IssueStateOpen, -1})
		i2 := resolver.Assign(agents, &issues, Issue{11, "B", CategoryGeneral, PriorityLow, IssueStateOpen, -1})
		i3 := resolver.Assign(agents, &issues, Issue{12, "C", CategoryGeneral, PriorityLow, IssueStateOpen, -1})

		if i1.AssignedAgentID != 0 {
			panic(fmt.Sprintf("expected agent 0, got %d", i1.AssignedAgentID))
		}
		if i2.AssignedAgentID != 1 {
			panic(fmt.Sprintf("expected agent 1, got %d", i2.AssignedAgentID))
		}
		if i3.AssignedAgentID != 0 {
			panic(fmt.Sprintf("expected agent 0 (wrap), got %d", i3.AssignedAgentID))
		}
	})

	// Test 3: GetAgentIssues returns correct issues
	test("test_get_agent_issues", func() {
		agents := []Agent{
			{0, "Alice", 0, nil},
			{1, "Bob", 0, nil},
		}
		issues := []Issue{}
		rr := &RoundRobinStrategy{}
		resolver := NewIssueResolver(rr)

		resolver.Assign(agents, &issues, Issue{20, "X", CategoryBilling, PriorityHigh, IssueStateOpen, -1})
		resolver.Assign(agents, &issues, Issue{21, "Y", CategoryTechnical, PriorityLow, IssueStateOpen, -1})
		resolver.Assign(agents, &issues, Issue{22, "Z", CategoryGeneral, PriorityMedium, IssueStateOpen, -1})

		aliceIssues := resolver.GetAgentIssues(issues, 0)
		bobIssues := resolver.GetAgentIssues(issues, 1)

		if len(aliceIssues) != 2 {
			panic(fmt.Sprintf("expected 2 issues for Alice, got %d", len(aliceIssues)))
		}
		if len(bobIssues) != 1 {
			panic(fmt.Sprintf("expected 1 issue for Bob, got %d", len(bobIssues)))
		}
		if aliceIssues[0].ID != 20 {
			panic(fmt.Sprintf("expected issue 20 first for Alice, got %d", aliceIssues[0].ID))
		}
		if aliceIssues[1].ID != 22 {
			panic(fmt.Sprintf("expected issue 22 second for Alice, got %d", aliceIssues[1].ID))
		}
		if bobIssues[0].ID != 21 {
			panic(fmt.Sprintf("expected issue 21 for Bob, got %d", bobIssues[0].ID))
		}
	})

	// Test 4: agent CurrentLoad increments on assignment
	test("test_agent_load_increment", func() {
		agents := []Agent{
			{0, "Alice", 0, nil},
			{1, "Bob", 0, nil},
		}
		issues := []Issue{}
		rr := &RoundRobinStrategy{}
		resolver := NewIssueResolver(rr)

		resolver.Assign(agents, &issues, Issue{30, "A", CategoryGeneral, PriorityLow, IssueStateOpen, -1})
		resolver.Assign(agents, &issues, Issue{31, "B", CategoryGeneral, PriorityLow, IssueStateOpen, -1})
		resolver.Assign(agents, &issues, Issue{32, "C", CategoryGeneral, PriorityLow, IssueStateOpen, -1})

		if agents[0].CurrentLoad != 2 {
			panic(fmt.Sprintf("expected Alice load=2, got %d", agents[0].CurrentLoad))
		}
		if agents[1].CurrentLoad != 1 {
			panic(fmt.Sprintf("expected Bob load=1, got %d", agents[1].CurrentLoad))
		}
	})

	// Test 5: empty issues returns empty for GetAgentIssues
	test("test_empty_issues", func() {
		issues := []Issue{}
		rr := &RoundRobinStrategy{}
		resolver := NewIssueResolver(rr)
		result := resolver.GetAgentIssues(issues, 0)
		if len(result) != 0 {
			panic("expected empty result")
		}
	})

	fmt.Printf("PART1_SUMMARY %d/%d\n", passed, passed+failed)
	return failed
}

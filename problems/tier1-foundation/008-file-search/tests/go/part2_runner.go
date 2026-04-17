package main

import "fmt"

func buildTestTree2() *FileNode {
	mainCpp := &FileNode{Name: "main.cpp", Size: 50, Extension: "cpp", IsDirectory: false}
	utilsCpp := &FileNode{Name: "utils.cpp", Size: 120, Extension: "cpp", IsDirectory: false}
	helperH := &FileNode{Name: "helper.h", Size: 10, Extension: "h", IsDirectory: false}
	readmeMd := &FileNode{Name: "readme.md", Size: 5, Extension: "md", IsDirectory: false}
	reportPdf := &FileNode{Name: "report.pdf", Size: 200, Extension: "pdf", IsDirectory: false}
	buildSh := &FileNode{Name: "build.sh", Size: 2, Extension: "sh", IsDirectory: false}
	src := &FileNode{Name: "src", IsDirectory: true, Children: []*FileNode{mainCpp, utilsCpp, helperH}}
	docs := &FileNode{Name: "docs", IsDirectory: true, Children: []*FileNode{readmeMd, reportPdf}}
	return &FileNode{Name: "project", IsDirectory: true, Children: []*FileNode{src, docs, buildSh}}
}

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

	root := buildTestTree2()

	// Test 1: AND filter — .cpp files AND size >= 100
	test("test_and_filter_cpp_and_large", func() {
		extCriteria := &SearchByExtensionCriteria{Ext: "cpp"}
		sizeCriteria := &SearchByMinSizeCriteria{MinSize: 100}
		results := SearchComposite(root, []SearchCriteria{extCriteria, sizeCriteria}, "AND")
		if len(results) != 1 {
			panic(fmt.Sprintf("expected 1 result, got %d", len(results)))
		}
		if results[0].Name != "utils.cpp" {
			panic("expected utils.cpp")
		}
	})

	// Test 2: OR filter — .pdf files OR name contains "main"
	test("test_or_filter_pdf_or_main", func() {
		extCriteria := &SearchByExtensionCriteria{Ext: "pdf"}
		nameCriteria := &SearchByNameCriteria{Substring: "main"}
		results := SearchComposite(root, []SearchCriteria{extCriteria, nameCriteria}, "OR")
		if len(results) != 2 {
			panic(fmt.Sprintf("expected 2 results, got %d", len(results)))
		}
		foundMain, foundReport := false, false
		for _, f := range results {
			if f.Name == "main.cpp" {
				foundMain = true
			}
			if f.Name == "report.pdf" {
				foundReport = true
			}
		}
		if !foundMain || !foundReport {
			panic("should find main.cpp and report.pdf")
		}
	})

	// Test 3: AND filter with no matches — .h files AND size >= 50
	test("test_and_filter_no_match", func() {
		extCriteria := &SearchByExtensionCriteria{Ext: "h"}
		sizeCriteria := &SearchByMinSizeCriteria{MinSize: 50}
		results := SearchComposite(root, []SearchCriteria{extCriteria, sizeCriteria}, "AND")
		if len(results) != 0 {
			panic(fmt.Sprintf("expected 0 results, got %d", len(results)))
		}
	})

	// Test 4: OR filter with single criterion behaves like direct search
	test("test_single_criterion_or", func() {
		extCriteria := &SearchByExtensionCriteria{Ext: "cpp"}
		composite := SearchComposite(root, []SearchCriteria{extCriteria}, "OR")
		direct := SearchByExtension(root, "cpp")
		if len(composite) != len(direct) {
			panic(fmt.Sprintf("composite %d != direct %d", len(composite), len(direct)))
		}
	})

	// Test 5: AND filter — files >= 50KB with '.' in name
	test("test_and_filter_size_and_dot", func() {
		sizeCriteria := &SearchByMinSizeCriteria{MinSize: 50}
		nameCriteria := &SearchByNameCriteria{Substring: "."}
		results := SearchComposite(root, []SearchCriteria{sizeCriteria, nameCriteria}, "AND")
		// Files >= 50KB with '.' in name: main.cpp(50), utils.cpp(120), report.pdf(200)
		if len(results) != 3 {
			panic(fmt.Sprintf("expected 3 results, got %d", len(results)))
		}
	})

	fmt.Printf("PART2_SUMMARY %d/%d\n", passed, passed+failed)
	return failed
}

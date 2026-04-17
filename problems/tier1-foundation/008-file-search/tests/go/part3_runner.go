package main

import "fmt"

func buildTestTree3() *FileNode {
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

	root := buildTestTree3()

	// Test 1: search all files >= 1 KB, sort by name (alphabetical)
	test("test_sort_by_name", func() {
		criteria := &SearchByMinSizeCriteria{MinSize: 1}
		results := SearchAndSort(root, criteria, "name")
		if len(results) != 6 {
			panic(fmt.Sprintf("expected 6 results, got %d", len(results)))
		}
		expected := []string{"build.sh", "helper.h", "main.cpp", "readme.md", "report.pdf", "utils.cpp"}
		for i, name := range expected {
			if results[i].Name != name {
				panic(fmt.Sprintf("index %d: expected %s, got %s", i, name, results[i].Name))
			}
		}
	})

	// Test 2: search all files >= 1 KB, sort by size (largest first)
	test("test_sort_by_size", func() {
		criteria := &SearchByMinSizeCriteria{MinSize: 1}
		results := SearchAndSort(root, criteria, "size")
		if len(results) != 6 {
			panic(fmt.Sprintf("expected 6 results, got %d", len(results)))
		}
		expected := []string{"report.pdf", "utils.cpp", "main.cpp", "helper.h", "readme.md", "build.sh"}
		for i, name := range expected {
			if results[i].Name != name {
				panic(fmt.Sprintf("index %d: expected %s, got %s", i, name, results[i].Name))
			}
		}
	})

	// Test 3: search .cpp files, sort by size
	test("test_search_cpp_sort_by_size", func() {
		criteria := &SearchByExtensionCriteria{Ext: "cpp"}
		results := SearchAndSort(root, criteria, "size")
		if len(results) != 2 {
			panic(fmt.Sprintf("expected 2 results, got %d", len(results)))
		}
		if results[0].Name != "utils.cpp" {
			panic(fmt.Sprintf("expected utils.cpp first, got %s", results[0].Name))
		}
		if results[1].Name != "main.cpp" {
			panic(fmt.Sprintf("expected main.cpp second, got %s", results[1].Name))
		}
	})

	// Test 4: search all files >= 1 KB, sort by extension (alphabetical)
	test("test_sort_by_extension", func() {
		criteria := &SearchByMinSizeCriteria{MinSize: 1}
		results := SearchAndSort(root, criteria, "extension")
		if len(results) != 6 {
			panic(fmt.Sprintf("expected 6 results, got %d", len(results)))
		}
		// Extensions in order: cpp, cpp, h, md, pdf, sh
		expectedExts := []string{"cpp", "cpp", "h", "md", "pdf", "sh"}
		for i, ext := range expectedExts {
			if results[i].Extension != ext {
				panic(fmt.Sprintf("index %d: expected ext %s, got %s", i, ext, results[i].Extension))
			}
		}
	})

	// Test 5: empty results remain empty after sort
	test("test_sort_empty_results", func() {
		criteria := &SearchByExtensionCriteria{Ext: "java"}
		results := SearchAndSort(root, criteria, "name")
		if len(results) != 0 {
			panic(fmt.Sprintf("expected 0 results, got %d", len(results)))
		}
	})

	fmt.Printf("PART3_SUMMARY %d/%d\n", passed, passed+failed)
	return failed
}

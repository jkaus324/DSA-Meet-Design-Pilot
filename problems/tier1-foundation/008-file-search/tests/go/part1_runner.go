package main

import "fmt"

func buildTestTree() *FileNode {
	//  project/
	//  ├── src/
	//  │   ├── main.cpp      (50 KB)
	//  │   ├── utils.cpp     (120 KB)
	//  │   └── helper.h      (10 KB)
	//  ├── docs/
	//  │   ├── readme.md     (5 KB)
	//  │   └── report.pdf    (200 KB)
	//  └── build.sh          (2 KB)
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

	root := buildTestTree()

	// Test 1: search_by_extension — find all .cpp files
	test("test_search_by_extension_cpp", func() {
		results := SearchByExtension(root, "cpp")
		if len(results) != 2 {
			panic(fmt.Sprintf("expected 2 results, got %d", len(results)))
		}
		foundMain, foundUtils := false, false
		for _, f := range results {
			if f.Name == "main.cpp" {
				foundMain = true
			}
			if f.Name == "utils.cpp" {
				foundUtils = true
			}
		}
		if !foundMain || !foundUtils {
			panic("should find main.cpp and utils.cpp")
		}
	})

	// Test 2: search_by_extension — find .pdf files
	test("test_search_by_extension_pdf", func() {
		results := SearchByExtension(root, "pdf")
		if len(results) != 1 {
			panic(fmt.Sprintf("expected 1 result, got %d", len(results)))
		}
		if results[0].Name != "report.pdf" {
			panic("expected report.pdf")
		}
	})

	// Test 3: search_by_extension — no matches
	test("test_search_by_extension_no_match", func() {
		results := SearchByExtension(root, "java")
		if len(results) != 0 {
			panic(fmt.Sprintf("expected 0 results, got %d", len(results)))
		}
	})

	// Test 4: search_by_size — files >= 100 KB
	test("test_search_by_size_100", func() {
		results := SearchBySize(root, 100)
		if len(results) != 2 {
			panic(fmt.Sprintf("expected 2 results, got %d", len(results)))
		}
		foundUtils, foundReport := false, false
		for _, f := range results {
			if f.Name == "utils.cpp" {
				foundUtils = true
			}
			if f.Name == "report.pdf" {
				foundReport = true
			}
		}
		if !foundUtils || !foundReport {
			panic("should find utils.cpp and report.pdf")
		}
	})

	// Test 5: search_by_size — files >= 200 KB (exact match)
	test("test_search_by_size_exact", func() {
		results := SearchBySize(root, 200)
		if len(results) != 1 {
			panic(fmt.Sprintf("expected 1 result, got %d", len(results)))
		}
		if results[0].Name != "report.pdf" {
			panic("expected report.pdf")
		}
	})

	// Test 6: search_by_name — files containing "main"
	test("test_search_by_name_main", func() {
		results := SearchByName(root, "main")
		if len(results) != 1 {
			panic(fmt.Sprintf("expected 1 result, got %d", len(results)))
		}
		if results[0].Name != "main.cpp" {
			panic("expected main.cpp")
		}
	})

	// Test 7: search_by_name — files containing "report"
	test("test_search_by_name_report", func() {
		results := SearchByName(root, "report")
		if len(results) != 1 {
			panic(fmt.Sprintf("expected 1 result, got %d", len(results)))
		}
		if results[0].Name != "report.pdf" {
			panic("expected report.pdf")
		}
	})

	// Test 8: empty tree returns empty
	test("test_empty_tree", func() {
		emptyRoot := &FileNode{Name: "empty", IsDirectory: true}
		if len(SearchByExtension(emptyRoot, "cpp")) != 0 {
			panic("expected 0 for empty tree")
		}
		if len(SearchBySize(emptyRoot, 10)) != 0 {
			panic("expected 0 for empty tree")
		}
		if len(SearchByName(emptyRoot, "test")) != 0 {
			panic("expected 0 for empty tree")
		}
	})

	// Test 9: single file tree
	test("test_single_file_tree", func() {
		singleFile := &FileNode{Name: "test.txt", Size: 30, Extension: "txt", IsDirectory: false}
		singleDir := &FileNode{Name: "root", IsDirectory: true, Children: []*FileNode{singleFile}}
		results := SearchByExtension(singleDir, "txt")
		if len(results) != 1 {
			panic(fmt.Sprintf("expected 1 result, got %d", len(results)))
		}
		if results[0].Name != "test.txt" {
			panic("expected test.txt")
		}
	})

	fmt.Printf("PART1_SUMMARY %d/%d\n", passed, passed+failed)
	return failed
}

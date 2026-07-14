// File Search — Strategy + Composite reference solution (Go).
package main

import "sort"

type FileNode struct {
	name        string
	size        int
	extension   string
	isDirectory bool
	children    []*FileNode
}

type SearchCriteria interface {
	matches(f *FileNode) bool
}

type SearchByExtension struct{ ext string }

func (s SearchByExtension) matches(f *FileNode) bool {
	return !f.isDirectory && f.extension == s.ext
}

type SearchByMinSize struct{ minSize int }

func (s SearchByMinSize) matches(f *FileNode) bool {
	return !f.isDirectory && f.size >= s.minSize
}

type SearchByName struct{ substring string }

func (s SearchByName) matches(f *FileNode) bool {
	return contains(f.name, s.substring)
}

func contains(s, sub string) bool {
	if sub == "" {
		return true
	}
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}

type AndFilter struct{ criteria []SearchCriteria }

func (a AndFilter) matches(f *FileNode) bool {
	for _, c := range a.criteria {
		if !c.matches(f) {
			return false
		}
	}
	return true
}

type OrFilter struct{ criteria []SearchCriteria }

func (o OrFilter) matches(f *FileNode) bool {
	for _, c := range o.criteria {
		if c.matches(f) {
			return true
		}
	}
	return false
}

func dfs(node *FileNode, criteria SearchCriteria, results *[]*FileNode) {
	if node == nil {
		return
	}
	if !node.isDirectory && criteria.matches(node) {
		*results = append(*results, node)
	}
	for _, c := range node.children {
		dfs(c, criteria, results)
	}
}

func searchByExtension(root *FileNode, ext string) []*FileNode {
	out := []*FileNode{}
	dfs(root, SearchByExtension{ext}, &out)
	return out
}

func searchBySize(root *FileNode, minSize int) []*FileNode {
	out := []*FileNode{}
	dfs(root, SearchByMinSize{minSize}, &out)
	return out
}

func searchByName(root *FileNode, substring string) []*FileNode {
	out := []*FileNode{}
	dfs(root, SearchByName{substring}, &out)
	return out
}

func searchComposite(root *FileNode, criteria []SearchCriteria, mode string) []*FileNode {
	var f SearchCriteria
	if mode == "AND" {
		f = AndFilter{criteria}
	} else {
		f = OrFilter{criteria}
	}
	out := []*FileNode{}
	dfs(root, f, &out)
	return out
}

func searchAndSort(root *FileNode, criteria SearchCriteria, sortBy string) []*FileNode {
	out := []*FileNode{}
	dfs(root, criteria, &out)
	switch sortBy {
	case "name":
		sort.SliceStable(out, func(i, j int) bool { return out[i].name < out[j].name })
	case "size":
		sort.SliceStable(out, func(i, j int) bool { return out[i].size > out[j].size })
	case "extension":
		sort.SliceStable(out, func(i, j int) bool { return out[i].extension < out[j].extension })
	}
	return out
}

// ─── Module fixture state ───────────────────────────────────────────────────

var gRoot *FileNode

func reset_service() {
	gRoot = nil
}

func fs_build_default_tree() {
	mainCpp := &FileNode{name: "main.cpp", size: 50, extension: "cpp"}
	utilsCpp := &FileNode{name: "utils.cpp", size: 120, extension: "cpp"}
	helperH := &FileNode{name: "helper.h", size: 10, extension: "h"}
	readme := &FileNode{name: "readme.md", size: 5, extension: "md"}
	report := &FileNode{name: "report.pdf", size: 200, extension: "pdf"}
	buildSh := &FileNode{name: "build.sh", size: 2, extension: "sh"}
	src := &FileNode{name: "src", isDirectory: true, children: []*FileNode{mainCpp, utilsCpp, helperH}}
	docs := &FileNode{name: "docs", isDirectory: true, children: []*FileNode{readme, report}}
	gRoot = &FileNode{name: "project", isDirectory: true, children: []*FileNode{src, docs, buildSh}}
}

func fs_build_empty_tree() {
	gRoot = &FileNode{name: "empty", isDirectory: true}
}

func fs_build_single_file_tree() {
	f := &FileNode{name: "test.txt", size: 30, extension: "txt"}
	gRoot = &FileNode{name: "root", isDirectory: true, children: []*FileNode{f}}
}

func fs_count_by_extension(ext string) int {
	if gRoot == nil {
		return 0
	}
	return len(searchByExtension(gRoot, ext))
}

func fs_has_by_extension(ext string, name string) bool {
	if gRoot == nil {
		return false
	}
	for _, f := range searchByExtension(gRoot, ext) {
		if f.name == name {
			return true
		}
	}
	return false
}

func fs_count_by_size(minSize int) int {
	if gRoot == nil {
		return 0
	}
	return len(searchBySize(gRoot, minSize))
}

func fs_has_by_size(minSize int, name string) bool {
	if gRoot == nil {
		return false
	}
	for _, f := range searchBySize(gRoot, minSize) {
		if f.name == name {
			return true
		}
	}
	return false
}

func fs_count_by_name(sub string) int {
	if gRoot == nil {
		return 0
	}
	return len(searchByName(gRoot, sub))
}

func fs_has_by_name(sub string, name string) bool {
	if gRoot == nil {
		return false
	}
	for _, f := range searchByName(gRoot, sub) {
		if f.name == name {
			return true
		}
	}
	return false
}

func fs_count_composite_and(ext string, minSize int) int {
	if gRoot == nil {
		return 0
	}
	return len(searchComposite(gRoot, []SearchCriteria{SearchByExtension{ext}, SearchByMinSize{minSize}}, "AND"))
}

func fs_count_composite_or(ext string, minSize int) int {
	if gRoot == nil {
		return 0
	}
	return len(searchComposite(gRoot, []SearchCriteria{SearchByExtension{ext}, SearchByMinSize{minSize}}, "OR"))
}

func fs_first_sorted_by(ext string, sortBy string) string {
	if gRoot == nil {
		return ""
	}
	v := searchAndSort(gRoot, SearchByExtension{ext}, sortBy)
	if len(v) > 0 {
		return v[0].name
	}
	return ""
}

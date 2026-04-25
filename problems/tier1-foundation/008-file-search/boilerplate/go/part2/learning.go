package main

import "strings"

// ─── Data Model (given — do not modify) ──────────────────────────────────────

type FileNode struct {
	Name        string
	Size        int
	Extension   string
	IsDirectory bool
	Children    []*FileNode
}

// ─── Search Criteria Interface ────────────────────────────────────────────────

type SearchCriteria interface {
	Matches(file *FileNode) bool
}

// ─── Existing Criteria (provided — already implemented) ──────────────────────

type SearchByExtensionCriteria struct{ Ext string }

func (c *SearchByExtensionCriteria) Matches(file *FileNode) bool {
	return !file.IsDirectory && file.Extension == c.Ext
}

type SearchByMinSizeCriteria struct{ MinSize int }

func (c *SearchByMinSizeCriteria) Matches(file *FileNode) bool {
	return !file.IsDirectory && file.Size >= c.MinSize
}

type SearchByNameCriteria struct{ Substring string }

func (c *SearchByNameCriteria) Matches(file *FileNode) bool {
	return strings.Contains(file.Name, c.Substring)
}

// ─── NEW: Composite Filters ───────────────────────────────────────────────────

type AndFilter struct {
	Criteria []SearchCriteria
}

func (f *AndFilter) Matches(file *FileNode) bool {
	// TODO: return true only if ALL criteria match the file
	// HINT: iterate through criteria, return false if any doesn't match
	return false
}

type OrFilter struct {
	Criteria []SearchCriteria
}

func (f *OrFilter) Matches(file *FileNode) bool {
	// TODO: return true if ANY criterion matches the file
	// HINT: iterate through criteria, return true if any matches
	return false
}

// ─── Search Engine (provided — already implemented) ───────────────────────────

type FileSearchEngine struct{}

func (e *FileSearchEngine) dfs(node *FileNode, criteria SearchCriteria, results *[]*FileNode) {
	if node == nil {
		return
	}
	if !node.IsDirectory && criteria.Matches(node) {
		*results = append(*results, node)
	}
	for _, child := range node.Children {
		e.dfs(child, criteria, results)
	}
}

func (e *FileSearchEngine) Search(root *FileNode, criteria SearchCriteria) []*FileNode {
	var results []*FileNode
	e.dfs(root, criteria, &results)
	return results
}

// ─── Test Entry Points ────────────────────────────────────────────────────────

func SearchByExtension(root *FileNode, ext string) []*FileNode {
	return (&FileSearchEngine{}).Search(root, &SearchByExtensionCriteria{Ext: ext})
}

func SearchBySize(root *FileNode, minSize int) []*FileNode {
	return (&FileSearchEngine{}).Search(root, &SearchByMinSizeCriteria{MinSize: minSize})
}

func SearchByName(root *FileNode, substring string) []*FileNode {
	return (&FileSearchEngine{}).Search(root, &SearchByNameCriteria{Substring: substring})
}

func SearchComposite(root *FileNode, criteria []SearchCriteria, mode string) []*FileNode {
	if mode == "AND" {
		return (&FileSearchEngine{}).Search(root, &AndFilter{Criteria: criteria})
	}
	return (&FileSearchEngine{}).Search(root, &OrFilter{Criteria: criteria})
}

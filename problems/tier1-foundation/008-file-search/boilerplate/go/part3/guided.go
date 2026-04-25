package main

import (
	"sort"
	"strings"
)

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

// ─── Existing Criteria ────────────────────────────────────────────────────────

type SearchByExtensionCriteria struct{ Ext string }

func (c *SearchByExtensionCriteria) Matches(file *FileNode) bool {
	return false // TODO: implement
}

type SearchByMinSizeCriteria struct{ MinSize int }

func (c *SearchByMinSizeCriteria) Matches(file *FileNode) bool {
	return false // TODO: implement
}

type SearchByNameCriteria struct{ Substring string }

func (c *SearchByNameCriteria) Matches(file *FileNode) bool {
	_ = strings.Contains
	return false // TODO: implement
}

type AndFilter struct{ Criteria []SearchCriteria }

func (f *AndFilter) Matches(file *FileNode) bool { return false } // TODO: implement

type OrFilter struct{ Criteria []SearchCriteria }

func (f *OrFilter) Matches(file *FileNode) bool { return false } // TODO: implement

// ─── NEW: Sort Strategy ───────────────────────────────────────────────────────
// HINT: A sort strategy compares two FileNode pointers.
// It is completely independent of the search criteria.

type SortStrategy interface {
	Compare(a, b *FileNode) bool
}

// TODO: Implement SortByName (alphabetical a-z)
// TODO: Implement SortBySize (largest first)
// TODO: Implement SortByExtension (alphabetical a-z)

// ─── Search Engine ────────────────────────────────────────────────────────────

type FileSearchEngine struct{}

func (e *FileSearchEngine) dfs(node *FileNode, criteria SearchCriteria, results *[]*FileNode) {
	if node == nil {
		return
	}
	// TODO: if node is a file and matches criteria, add to results
	// TODO: recurse into all children
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

func SearchAndSort(root *FileNode, criteria SearchCriteria, sortBy string) []*FileNode {
	results := (&FileSearchEngine{}).Search(root, criteria)
	// TODO: Based on sortBy ("name", "size", "extension"), create the right
	// SortStrategy and use it to sort results.
	// HINT: Use sort.Slice with a call to strategy.Compare()
	_ = sort.Slice
	return results
}

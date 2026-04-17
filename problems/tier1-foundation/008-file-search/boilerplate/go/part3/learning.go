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

type AndFilter struct{ Criteria []SearchCriteria }

func (f *AndFilter) Matches(file *FileNode) bool {
	for _, c := range f.Criteria {
		if !c.Matches(file) {
			return false
		}
	}
	return true
}

type OrFilter struct{ Criteria []SearchCriteria }

func (f *OrFilter) Matches(file *FileNode) bool {
	for _, c := range f.Criteria {
		if c.Matches(file) {
			return true
		}
	}
	return false
}

// ─── NEW: Sort Strategy ───────────────────────────────────────────────────────

type SortStrategy interface {
	Compare(a, b *FileNode) bool
}

type SortByName struct{}

func (s *SortByName) Compare(a, b *FileNode) bool {
	// TODO: return true if a's name comes before b's name alphabetically
	return false
}

type SortBySize struct{}

func (s *SortBySize) Compare(a, b *FileNode) bool {
	// TODO: return true if a's size is larger than b's (largest first)
	return false
}

type SortByExtension struct{}

func (s *SortByExtension) Compare(a, b *FileNode) bool {
	// TODO: return true if a's extension comes before b's alphabetically
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

func SearchAndSort(root *FileNode, criteria SearchCriteria, sortBy string) []*FileNode {
	results := (&FileSearchEngine{}).Search(root, criteria)
	// TODO: Based on sortBy ("name", "size", "extension"), create the right
	// SortStrategy and use it to sort results.
	// HINT: Use sort.Slice with a call to strategy.Compare()
	var strategy SortStrategy
	switch sortBy {
	case "name":
		strategy = &SortByName{}
	case "size":
		strategy = &SortBySize{}
	case "extension":
		strategy = &SortByExtension{}
	}
	if strategy != nil {
		sort.Slice(results, func(i, j int) bool {
			return strategy.Compare(results[i], results[j])
		})
	}
	return results
}

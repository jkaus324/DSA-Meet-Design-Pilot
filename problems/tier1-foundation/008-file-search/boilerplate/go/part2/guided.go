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

// ─── Existing Criteria ────────────────────────────────────────────────────────
// TODO: Copy your Part 1 criteria here (or extend them)

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

// ─── NEW: Composite Filters ───────────────────────────────────────────────────
// HINT: An AndFilter holds a slice of other criteria.
// It returns true only if ALL criteria match.
// An OrFilter returns true if ANY criterion matches.
// Both implement SearchCriteria — this is the Composite pattern.

type AndFilter struct {
	Criteria []SearchCriteria
}

func (f *AndFilter) Matches(file *FileNode) bool {
	// TODO: return true only if ALL criteria match
	return false
}

type OrFilter struct {
	Criteria []SearchCriteria
}

func (f *OrFilter) Matches(file *FileNode) bool {
	// TODO: return true if ANY criterion matches
	return false
}

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

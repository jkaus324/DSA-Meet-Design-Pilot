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

// ─── Concrete Search Criteria ─────────────────────────────────────────────────
// TODO: Implement the Matches() method for each criterion

type SearchByExtensionCriteria struct {
	Ext string
}

func (c *SearchByExtensionCriteria) Matches(file *FileNode) bool {
	// TODO: return true if file (not directory) has the matching extension
	return false
}

type SearchByMinSizeCriteria struct {
	MinSize int
}

func (c *SearchByMinSizeCriteria) Matches(file *FileNode) bool {
	// TODO: return true if file (not directory) has size >= MinSize
	return false
}

type SearchByNameCriteria struct {
	Substring string
}

func (c *SearchByNameCriteria) Matches(file *FileNode) bool {
	// TODO: return true if file name contains Substring
	_ = strings.Contains
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
	criteria := &SearchByExtensionCriteria{Ext: ext}
	return (&FileSearchEngine{}).Search(root, criteria)
}

func SearchBySize(root *FileNode, minSize int) []*FileNode {
	criteria := &SearchByMinSizeCriteria{MinSize: minSize}
	return (&FileSearchEngine{}).Search(root, criteria)
}

func SearchByName(root *FileNode, substring string) []*FileNode {
	criteria := &SearchByNameCriteria{Substring: substring}
	return (&FileSearchEngine{}).Search(root, criteria)
}

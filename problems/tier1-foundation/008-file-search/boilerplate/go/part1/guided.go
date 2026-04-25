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
// HINT: This interface lets you swap search logic at runtime.
// What method signature would let you check if a file matches a criterion?

type SearchCriteria interface {
	// HINT: what method checks if a file matches?
	// Matches(file *FileNode) bool
}

// ─── Concrete Search Criteria ─────────────────────────────────────────────────
// TODO: Implement a criterion for each search type:
//   - Search by extension (match files with a given extension)
//   - Search by minimum size (match files >= threshold)
//   - Search by name (match files whose name contains a substring)

// ─── Search Engine ────────────────────────────────────────────────────────────
// TODO: Implement a FileSearchEngine struct that:
//   - Accepts any SearchCriteria
//   - Has a Search() method that traverses the tree using DFS
//   - Returns all files matching the criterion
//   - Does NOT know about specific search criteria

// type FileSearchEngine struct{}
// func (e *FileSearchEngine) Search(root *FileNode, criteria SearchCriteria) []*FileNode

// ─── Test Entry Points ────────────────────────────────────────────────────────

func SearchByExtension(root *FileNode, ext string) []*FileNode {
	_ = strings.Contains // available if needed
	return nil
}

func SearchBySize(root *FileNode, minSize int) []*FileNode {
	return nil
}

func SearchByName(root *FileNode, substring string) []*FileNode {
	return nil
}

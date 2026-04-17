package main

// ─── Data Model (given — do not modify) ──────────────────────────────────────

type FileNode struct {
	Name        string
	Size        int
	Extension   string
	IsDirectory bool
	Children    []*FileNode
}

// ─── Your Design Starts Here ──────────────────────────────────────────────────
//
// Design and implement a FileSearchEngine that:
//   1. Traverses a file tree using DFS
//   2. Returns files matching a given search criterion
//   3. Allows new search criteria to be added WITHOUT modifying
//      the engine itself
//
// Think about:
//   - What abstraction lets you swap search logic at runtime?
//   - How would you add a 4th search criterion with zero changes
//     to existing code?
//   - How do you traverse a tree structure recursively?
//
// Entry points (must exist for tests):
//   func SearchByExtension(root *FileNode, ext string) []*FileNode
//   func SearchBySize(root *FileNode, minSize int) []*FileNode
//   func SearchByName(root *FileNode, substring string) []*FileNode
//
// ─────────────────────────────────────────────────────────────────────────────

func SearchByExtension(root *FileNode, ext string) []*FileNode {
	return nil
}

func SearchBySize(root *FileNode, minSize int) []*FileNode {
	return nil
}

func SearchByName(root *FileNode, substring string) []*FileNode {
	return nil
}

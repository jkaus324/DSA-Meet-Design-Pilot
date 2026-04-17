package main

// ─── Data Model (given — do not modify) ──────────────────────────────────────

type FileNode struct {
	Name        string
	Size        int
	Extension   string
	IsDirectory bool
	Children    []*FileNode
}

// ─── NEW in Extension 1 ───────────────────────────────────────────────────────
//
// The product team now wants COMPOSITE filters:
// find files that are .cpp AND larger than 100KB, or files that are
// .pdf OR named "report".
//
// Think about:
//   - How do you combine criteria without modifying existing strategies?
//   - What if the product team adds a 4th criterion tomorrow?
//   - Is your Part 1 design extensible enough to handle this?
//
// Entry points (must exist for tests):
//   func SearchByExtension(root *FileNode, ext string) []*FileNode
//   func SearchBySize(root *FileNode, minSize int) []*FileNode
//   func SearchByName(root *FileNode, substring string) []*FileNode
//   func SearchComposite(root *FileNode, criteria []SearchCriteria, mode string) []*FileNode
//
// ─────────────────────────────────────────────────────────────────────────────

type SearchCriteria interface {
	Matches(file *FileNode) bool
}

func SearchByExtension(root *FileNode, ext string) []*FileNode    { return nil }
func SearchBySize(root *FileNode, minSize int) []*FileNode        { return nil }
func SearchByName(root *FileNode, substring string) []*FileNode   { return nil }
func SearchComposite(root *FileNode, criteria []SearchCriteria, mode string) []*FileNode {
	return nil
}

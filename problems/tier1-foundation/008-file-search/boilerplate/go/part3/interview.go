package main

// ─── Data Model (given — do not modify) ──────────────────────────────────────

type FileNode struct {
	Name        string
	Size        int
	Extension   string
	IsDirectory bool
	Children    []*FileNode
}

// ─── NEW in Extension 2 ───────────────────────────────────────────────────────
//
// The product team wants search results SORTED by different strategies:
// alphabetically by name, by file size (largest first), or by extension.
//
// Think about:
//   - How do you decouple search criteria from sort logic?
//   - Is the sort strategy independent of the search strategy?
//   - Can you add a new sort order without touching the search logic?
//
// Entry points (must exist for tests):
//   func SearchByExtension(root *FileNode, ext string) []*FileNode
//   func SearchBySize(root *FileNode, minSize int) []*FileNode
//   func SearchByName(root *FileNode, substring string) []*FileNode
//   func SearchComposite(root *FileNode, criteria []SearchCriteria, mode string) []*FileNode
//   func SearchAndSort(root *FileNode, criteria SearchCriteria, sortBy string) []*FileNode
//
// ─────────────────────────────────────────────────────────────────────────────

type SearchCriteria interface {
	Matches(file *FileNode) bool
}

func SearchByExtension(root *FileNode, ext string) []*FileNode  { return nil }
func SearchBySize(root *FileNode, minSize int) []*FileNode      { return nil }
func SearchByName(root *FileNode, substring string) []*FileNode { return nil }
func SearchComposite(root *FileNode, criteria []SearchCriteria, mode string) []*FileNode {
	return nil
}
func SearchAndSort(root *FileNode, criteria SearchCriteria, sortBy string) []*FileNode {
	return nil
}

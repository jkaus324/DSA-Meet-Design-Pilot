# Problem 008 — File Search System

**Tier:** 1 (Foundation) | **Pattern:** Strategy + Composite | **DSA:** Tree (DFS/BFS)
**Companies:** Amazon, Microsoft | **Time:** 45 minutes

---

## Problem Statement

You're building a file search utility similar to the Linux `find` command. Given a file system represented as a tree, you must **search for files matching specific criteria** and return the results.

Different users need different search filters:
- A **developer** wants to find all `.cpp` files
- A **sysadmin** wants to find all files larger than a certain size
- A **manager** wants to find files by name

**Your task:** Design and implement a `FileSearchEngine` that can traverse a file tree and return files matching any search criterion — and allows new criteria to be added without modifying the engine itself.

---

## Before You Code

> Read this section carefully. This is where the design thinking happens.

**Ask yourself:**
1. What varies here? The *search criterion* varies. The *traversal logic* stays the same.
2. If you used `if-else` inside `search()`, what happens when a 4th filter is added? You modify existing code — violating Open/Closed Principle.
3. How does the Strategy pattern solve this? Each search filter becomes a separate class implementing a common interface.
4. How is the file system structured? As a tree — directories contain children (files and subdirectories). You need DFS or BFS to traverse it.

**The key insight:** The search criterion is a strategy. The file tree is a composite structure. Combining them gives you a flexible, extensible search engine.

---

## Data Structures

```cpp
struct FileNode {
    std::string name;          // "main.cpp", "docs", "report.pdf"
    int size;                  // file size in KB (0 for directories)
    std::string extension;     // "cpp", "pdf", "" (empty for directories)
    bool isDirectory;          // true for folders, false for files
    std::vector<FileNode*> children;  // non-empty only for directories
};
```

---

## Base Requirement — Search by a single criterion

Implement a `FileSearchEngine` that accepts any search criterion and uses DFS to traverse the file tree, returning all files that match.

You must support these three criteria:

| Criterion | Rule |
|-----------|------|
| Search by Extension | Match files with a given extension (e.g., "cpp") |
| Search by Min Size | Match files with size >= a given threshold |
| Search by Name | Match files whose name contains a given substring |

**Design goal:** Adding a 4th criterion must require **zero changes** to `FileSearchEngine` itself.

**Entry points (tests will call these):**
```cpp
vector<FileNode*> search_by_extension(FileNode* root, const string& ext);
vector<FileNode*> search_by_size(FileNode* root, int minSize);
vector<FileNode*> search_by_name(FileNode* root, const string& substring);
```

**What to implement:**
```cpp
class SearchCriteria {
public:
    virtual bool matches(const FileNode* file) const = 0;
    virtual ~SearchCriteria() = default;
};

class SearchByExtension : public SearchCriteria { ... };
class SearchByMinSize   : public SearchCriteria { ... };
class SearchByName      : public SearchCriteria { ... };

class FileSearchEngine {
public:
    vector<FileNode*> search(FileNode* root, const SearchCriteria& criteria);
    // Uses DFS to traverse the tree, collecting files that match
};
```

---

## Extension 1 — Combine criteria with AND/OR

The product team now wants **composite filters**: find files that are `.cpp` AND larger than 100KB, or files that are `.pdf` OR named "report".

> Example: SearchByExtension("cpp") AND SearchByMinSize(100) returns only `.cpp` files that are at least 100KB.

**Design challenge:** How do you combine criteria **without modifying** `SearchByExtension`, `SearchByMinSize`, or `FileSearchEngine`?

**New entry point:**
```cpp
vector<FileNode*> search_composite(FileNode* root,
                                   const vector<SearchCriteria*>& criteria,
                                   const string& mode);  // "AND" or "OR"
```

The function accepts a list of criteria and a mode. In "AND" mode, a file must match ALL criteria. In "OR" mode, a file must match ANY criterion.

**Hint:** An `AndFilter` and an `OrFilter` each hold a list of criteria and themselves implement `SearchCriteria`. This is the Composite pattern.

---

## Extension 2 — Sort results independently

The product team wants search results to be **sorted** by different strategies: alphabetically by name, by file size (largest first), or by extension.

**New entry point:**
```cpp
vector<FileNode*> search_and_sort(FileNode* root,
                                  const SearchCriteria& criteria,
                                  const string& sortBy);  // "name", "size", "extension"
```

The function searches using any criterion, then sorts the results using the specified sort strategy. The search logic and sort logic must be **independent** — changing one doesn't affect the other.

**Design challenge:** How do you decouple the search strategy from the sort strategy? Can you add new sort orders without touching the search logic?

---

## Running Tests

```bash
./run-tests.sh 008-file-search cpp
```

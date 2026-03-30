#include <iostream>
#include <vector>
#include <string>
#include <algorithm>
using namespace std;

// ─── Data Model (given — do not modify) ─────────────────────────────────────

struct FileNode {
    string name;          // "main.cpp", "docs", "report.pdf"
    int size;             // file size in KB (0 for directories)
    string extension;     // "cpp", "pdf", "" (empty for directories)
    bool isDirectory;     // true for folders, false for files
    vector<FileNode*> children;  // non-empty only for directories
};

// ─── Search Criteria Interface ──────────────────────────────────────────────

class SearchCriteria {
public:
    virtual bool matches(const FileNode* file) const = 0;
    virtual ~SearchCriteria() = default;
};

// ─── Existing Criteria ──────────────────────────────────────────────────────
// TODO: Copy your Part 1 + Part 2 criteria here

class SearchByExtension : public SearchCriteria {
    string ext;
public:
    SearchByExtension(const string& e) : ext(e) {}
    bool matches(const FileNode* file) const override {
        return false; // TODO: implement
    }
};

class SearchByMinSize : public SearchCriteria {
    int minSize;
public:
    SearchByMinSize(int s) : minSize(s) {}
    bool matches(const FileNode* file) const override {
        return false; // TODO: implement
    }
};

class SearchByName : public SearchCriteria {
    string substring;
public:
    SearchByName(const string& s) : substring(s) {}
    bool matches(const FileNode* file) const override {
        return false; // TODO: implement
    }
};

// ─── Composite Filters ──────────────────────────────────────────────────────

class AndFilter : public SearchCriteria {
    vector<SearchCriteria*> criteria;
public:
    AndFilter(const vector<SearchCriteria*>& c) : criteria(c) {}
    bool matches(const FileNode* file) const override {
        return false; // TODO: implement
    }
};

class OrFilter : public SearchCriteria {
    vector<SearchCriteria*> criteria;
public:
    OrFilter(const vector<SearchCriteria*>& c) : criteria(c) {}
    bool matches(const FileNode* file) const override {
        return false; // TODO: implement
    }
};

// ─── NEW: Sort Strategy ─────────────────────────────────────────────────────
// HINT: A sort strategy compares two FileNode pointers.
// It is completely independent of the search criteria.

// class SortStrategy {
// public:
//     virtual bool compare(const FileNode* a, const FileNode* b) const = 0;
//     virtual ~SortStrategy() = default;
// };

// TODO: Implement SortByName (alphabetical a-z)
// TODO: Implement SortBySize (largest first)
// TODO: Implement SortByExtension (alphabetical a-z)

// ─── Search Engine ──────────────────────────────────────────────────────────

class FileSearchEngine {
    void dfs(FileNode* node, const SearchCriteria& criteria, vector<FileNode*>& results) {
        if (!node) return;
        // TODO: if node is a file and matches criteria, add to results
        // TODO: recurse into all children
    }
public:
    vector<FileNode*> search(FileNode* root, const SearchCriteria& criteria) {
        vector<FileNode*> results;
        dfs(root, criteria, results);
        return results;
    }
};

// ─── Test Entry Points ───────────────────────────────────────────────────────

vector<FileNode*> search_by_extension(FileNode* root, const string& ext) {
    SearchByExtension criteria(ext);
    return FileSearchEngine().search(root, criteria);
}

vector<FileNode*> search_by_size(FileNode* root, int minSize) {
    SearchByMinSize criteria(minSize);
    return FileSearchEngine().search(root, criteria);
}

vector<FileNode*> search_by_name(FileNode* root, const string& substring) {
    SearchByName criteria(substring);
    return FileSearchEngine().search(root, criteria);
}

vector<FileNode*> search_composite(FileNode* root, const vector<SearchCriteria*>& criteria, const string& mode) {
    if (mode == "AND") {
        AndFilter filter(criteria);
        return FileSearchEngine().search(root, filter);
    } else {
        OrFilter filter(criteria);
        return FileSearchEngine().search(root, filter);
    }
}

vector<FileNode*> search_and_sort(FileNode* root, const SearchCriteria& criteria, const string& sortBy) {
    // TODO: search using criteria, then sort results by sortBy
    // sortBy is one of: "name", "size", "extension"
    auto results = FileSearchEngine().search(root, criteria);
    // TODO: sort results based on sortBy
    return results;
}

#ifndef RUNNING_TESTS
int main() {
    cout << "Part 3: Sort strategies — implement the TODOs above." << endl;
    return 0;
}
#endif

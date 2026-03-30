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

// ─── Existing Criteria (provided — already implemented) ─────────────────────

class SearchByExtension : public SearchCriteria {
    string ext;
public:
    SearchByExtension(const string& e) : ext(e) {}
    bool matches(const FileNode* file) const override {
        return !file->isDirectory && file->extension == ext;
    }
};

class SearchByMinSize : public SearchCriteria {
    int minSize;
public:
    SearchByMinSize(int s) : minSize(s) {}
    bool matches(const FileNode* file) const override {
        return !file->isDirectory && file->size >= minSize;
    }
};

class SearchByName : public SearchCriteria {
    string substring;
public:
    SearchByName(const string& s) : substring(s) {}
    bool matches(const FileNode* file) const override {
        return file->name.find(substring) != string::npos;
    }
};

// ─── Composite Filters (provided — already implemented) ─────────────────────

class AndFilter : public SearchCriteria {
    vector<SearchCriteria*> criteria;
public:
    AndFilter(const vector<SearchCriteria*>& c) : criteria(c) {}
    bool matches(const FileNode* file) const override {
        for (auto* c : criteria) {
            if (!c->matches(file)) return false;
        }
        return true;
    }
};

class OrFilter : public SearchCriteria {
    vector<SearchCriteria*> criteria;
public:
    OrFilter(const vector<SearchCriteria*>& c) : criteria(c) {}
    bool matches(const FileNode* file) const override {
        for (auto* c : criteria) {
            if (c->matches(file)) return true;
        }
        return false;
    }
};

// ─── NEW: Sort Strategy ─────────────────────────────────────────────────────

class SortStrategy {
public:
    virtual bool compare(const FileNode* a, const FileNode* b) const = 0;
    virtual ~SortStrategy() = default;
};

class SortByName : public SortStrategy {
public:
    bool compare(const FileNode* a, const FileNode* b) const override {
        // TODO: return true if a's name comes before b's name alphabetically
        return false;
    }
};

class SortBySize : public SortStrategy {
public:
    bool compare(const FileNode* a, const FileNode* b) const override {
        // TODO: return true if a's size is larger than b's (largest first)
        return false;
    }
};

class SortByExtension : public SortStrategy {
public:
    bool compare(const FileNode* a, const FileNode* b) const override {
        // TODO: return true if a's extension comes before b's alphabetically
        return false;
    }
};

// ─── Search Engine (provided — already implemented) ─────────────────────────

class FileSearchEngine {
    void dfs(FileNode* node, const SearchCriteria& criteria, vector<FileNode*>& results) {
        if (!node) return;
        if (!node->isDirectory && criteria.matches(node)) {
            results.push_back(node);
        }
        for (auto* child : node->children) {
            dfs(child, criteria, results);
        }
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
    auto results = FileSearchEngine().search(root, criteria);
    // TODO: Based on sortBy ("name", "size", "extension"), create the right
    // SortStrategy and use it to sort results.
    // HINT: Use std::sort with a lambda that calls strategy->compare()
    return results;
}

#ifndef RUNNING_TESTS
int main() {
    cout << "Part 3: Sort strategies — implement the SortBy TODOs, then sort in search_and_sort." << endl;
    return 0;
}
#endif

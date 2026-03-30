#include <iostream>
#include <vector>
#include <string>
#include <algorithm>
using namespace std;

// ─── Data Structure ──────────────────────────────────────────────────────────

struct FileNode {
    string name;          // "main.cpp", "docs", "report.pdf"
    int size;             // file size in KB (0 for directories)
    string extension;     // "cpp", "pdf", "" (empty for directories)
    bool isDirectory;     // true for folders, false for files
    vector<FileNode*> children;  // non-empty only for directories
};

// ─── Part 1: Strategy Interface ──────────────────────────────────────────────

class SearchCriteria {
public:
    virtual bool matches(const FileNode* file) const = 0;
    virtual ~SearchCriteria() = default;
};

// ─── Concrete Search Criteria ───────────────────────────────────────────────

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

// ─── Part 2: Composite Filters ──────────────────────────────────────────────

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

// ─── Part 3: Sort Strategy ──────────────────────────────────────────────────

class SortStrategy {
public:
    virtual bool compare(const FileNode* a, const FileNode* b) const = 0;
    virtual ~SortStrategy() = default;
};

class SortByName : public SortStrategy {
public:
    bool compare(const FileNode* a, const FileNode* b) const override {
        return a->name < b->name;
    }
};

class SortBySize : public SortStrategy {
public:
    bool compare(const FileNode* a, const FileNode* b) const override {
        return a->size > b->size;
    }
};

class SortByExtension : public SortStrategy {
public:
    bool compare(const FileNode* a, const FileNode* b) const override {
        return a->extension < b->extension;
    }
};

// ─── Search Engine ──────────────────────────────────────────────────────────

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
    SortStrategy* strategy = nullptr;
    SortByName sortByName;
    SortBySize sortBySize;
    SortByExtension sortByExtension;

    if (sortBy == "name") {
        strategy = &sortByName;
    } else if (sortBy == "size") {
        strategy = &sortBySize;
    } else if (sortBy == "extension") {
        strategy = &sortByExtension;
    }

    if (strategy) {
        sort(results.begin(), results.end(), [&strategy](const FileNode* a, const FileNode* b) {
            return strategy->compare(a, b);
        });
    }
    return results;
}

// ─── Main ────────────────────────────────────────────────────────────────────

#ifndef RUNNING_TESTS
int main() {
    // Build a small file tree
    FileNode file1{"main.cpp", 50, "cpp", false, {}};
    FileNode file2{"utils.cpp", 120, "cpp", false, {}};
    FileNode file3{"report.pdf", 200, "pdf", false, {}};
    FileNode src{"src", 0, "", true, {&file1, &file2}};
    FileNode root{"project", 0, "", true, {&src, &file3}};

    cout << "Search for .cpp files:" << endl;
    auto results = search_by_extension(&root, "cpp");
    for (const auto* f : results) {
        cout << "  " << f->name << " (" << f->size << " KB)" << endl;
    }

    cout << "\nComposite AND (.cpp AND >= 100KB):" << endl;
    SearchByExtension extCriteria("cpp");
    SearchByMinSize sizeCriteria(100);
    auto composite = search_composite(&root, {&extCriteria, &sizeCriteria}, "AND");
    for (const auto* f : composite) {
        cout << "  " << f->name << " (" << f->size << " KB)" << endl;
    }

    cout << "\nAll files sorted by size (largest first):" << endl;
    SearchByMinSize allFiles(1);
    auto sorted = search_and_sort(&root, allFiles, "size");
    for (const auto* f : sorted) {
        cout << "  " << f->name << " (" << f->size << " KB)" << endl;
    }

    return 0;
}
#endif

"""File Search — Strategy + Composite + Decorator reference solution (Python)."""

from abc import ABC, abstractmethod


class FileNode:
    def __init__(self, name, size, extension, isDirectory, children=None):
        self.name = name
        self.size = size
        self.extension = extension
        self.isDirectory = isDirectory
        self.children = children if children else []


class SearchCriteria(ABC):
    @abstractmethod
    def matches(self, file):
        ...


class SearchByExtension(SearchCriteria):
    def __init__(self, ext):
        self.ext = ext

    def matches(self, file):
        return (not file.isDirectory) and file.extension == self.ext


class SearchByMinSize(SearchCriteria):
    def __init__(self, minSize):
        self.minSize = minSize

    def matches(self, file):
        return (not file.isDirectory) and file.size >= self.minSize


class SearchByName(SearchCriteria):
    def __init__(self, substring):
        self.substring = substring

    def matches(self, file):
        return self.substring in file.name


class AndFilter(SearchCriteria):
    def __init__(self, criteria):
        self.criteria = criteria

    def matches(self, file):
        return all(c.matches(file) for c in self.criteria)


class OrFilter(SearchCriteria):
    def __init__(self, criteria):
        self.criteria = criteria

    def matches(self, file):
        return any(c.matches(file) for c in self.criteria)


def _dfs(node, criteria, results):
    if node is None:
        return
    if not node.isDirectory and criteria.matches(node):
        results.append(node)
    for c in node.children:
        _dfs(c, criteria, results)


def search_by_extension(root, ext):
    out = []
    _dfs(root, SearchByExtension(ext), out)
    return out


def search_by_size(root, minSize):
    out = []
    _dfs(root, SearchByMinSize(minSize), out)
    return out


def search_by_name(root, substring):
    out = []
    _dfs(root, SearchByName(substring), out)
    return out


def search_composite(root, criteria, mode):
    f = AndFilter(criteria) if mode == "AND" else OrFilter(criteria)
    out = []
    _dfs(root, f, out)
    return out


def search_and_sort(root, criteria, sortBy):
    out = []
    _dfs(root, criteria, out)
    if sortBy == "name":
        out.sort(key=lambda f: f.name)
    elif sortBy == "size":
        out.sort(key=lambda f: -f.size)
    elif sortBy == "extension":
        out.sort(key=lambda f: f.extension)
    return out


# ─── Module fixture state ───────────────────────────────────────────────────

_g_root = None


def reset_service():
    global _g_root
    _g_root = None


def fs_build_default_tree():
    global _g_root
    main_cpp = FileNode("main.cpp", 50, "cpp", False)
    utils_cpp = FileNode("utils.cpp", 120, "cpp", False)
    helper_h = FileNode("helper.h", 10, "h", False)
    readme = FileNode("readme.md", 5, "md", False)
    report = FileNode("report.pdf", 200, "pdf", False)
    build_sh = FileNode("build.sh", 2, "sh", False)
    src = FileNode("src", 0, "", True, [main_cpp, utils_cpp, helper_h])
    docs = FileNode("docs", 0, "", True, [readme, report])
    _g_root = FileNode("project", 0, "", True, [src, docs, build_sh])


def fs_build_empty_tree():
    global _g_root
    _g_root = FileNode("empty", 0, "", True)


def fs_build_single_file_tree():
    global _g_root
    f = FileNode("test.txt", 30, "txt", False)
    _g_root = FileNode("root", 0, "", True, [f])


def fs_count_by_extension(ext):
    return len(search_by_extension(_g_root, ext)) if _g_root else 0


def fs_has_by_extension(ext, name):
    if _g_root is None:
        return False
    return any(f.name == name for f in search_by_extension(_g_root, ext))


def fs_count_by_size(minSize):
    return len(search_by_size(_g_root, minSize)) if _g_root else 0


def fs_has_by_size(minSize, name):
    if _g_root is None:
        return False
    return any(f.name == name for f in search_by_size(_g_root, minSize))


def fs_count_by_name(sub):
    return len(search_by_name(_g_root, sub)) if _g_root else 0


def fs_has_by_name(sub, name):
    if _g_root is None:
        return False
    return any(f.name == name for f in search_by_name(_g_root, sub))


def fs_count_composite_and(ext, minSize):
    if _g_root is None:
        return 0
    return len(search_composite(_g_root, [SearchByExtension(ext), SearchByMinSize(minSize)], "AND"))


def fs_count_composite_or(ext, minSize):
    if _g_root is None:
        return 0
    return len(search_composite(_g_root, [SearchByExtension(ext), SearchByMinSize(minSize)], "OR"))


def fs_first_sorted_by(ext, sortBy):
    if _g_root is None:
        return ""
    v = search_and_sort(_g_root, SearchByExtension(ext), sortBy)
    return v[0].name if v else ""

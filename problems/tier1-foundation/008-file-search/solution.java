// File Search — Solution (Java, Strategy + Composite)
import java.util.*;

class FileNode {
    public String name;
    public int size;
    public String extension;
    public boolean isDirectory;
    public List<FileNode> children;

    public FileNode(String name, int size, String extension, boolean isDirectory) {
        this(name, size, extension, isDirectory, new ArrayList<>());
    }

    public FileNode(String name, int size, String extension, boolean isDirectory, List<FileNode> children) {
        this.name = name; this.size = size; this.extension = extension;
        this.isDirectory = isDirectory; this.children = children;
    }
}

interface SearchCriteria {
    boolean matches(FileNode file);
}

class SearchByExtension implements SearchCriteria {
    private final String ext;
    public SearchByExtension(String ext) { this.ext = ext; }
    public boolean matches(FileNode file) {
        return !file.isDirectory && Objects.equals(file.extension, ext);
    }
}

class SearchByMinSize implements SearchCriteria {
    private final int minSize;
    public SearchByMinSize(int minSize) { this.minSize = minSize; }
    public boolean matches(FileNode file) {
        return !file.isDirectory && file.size >= minSize;
    }
}

class SearchByName implements SearchCriteria {
    private final String substring;
    public SearchByName(String substring) { this.substring = substring; }
    public boolean matches(FileNode file) {
        return file.name.contains(substring);
    }
}

class AndFilter implements SearchCriteria {
    private final List<SearchCriteria> criteria;
    public AndFilter(List<SearchCriteria> criteria) { this.criteria = criteria; }
    public boolean matches(FileNode file) {
        for (SearchCriteria c : criteria) if (!c.matches(file)) return false;
        return true;
    }
}

class OrFilter implements SearchCriteria {
    private final List<SearchCriteria> criteria;
    public OrFilter(List<SearchCriteria> criteria) { this.criteria = criteria; }
    public boolean matches(FileNode file) {
        for (SearchCriteria c : criteria) if (c.matches(file)) return true;
        return false;
    }
}

public class Solution {
    private static FileNode root = null;

    private static void dfs(FileNode node, SearchCriteria c, List<FileNode> results) {
        if (node == null) return;
        if (!node.isDirectory && c.matches(node)) results.add(node);
        for (FileNode ch : node.children) dfs(ch, c, results);
    }

    private static List<FileNode> searchByExtension(FileNode root, String ext) {
        List<FileNode> out = new ArrayList<>();
        dfs(root, new SearchByExtension(ext), out);
        return out;
    }

    private static List<FileNode> searchBySize(FileNode root, int minSize) {
        List<FileNode> out = new ArrayList<>();
        dfs(root, new SearchByMinSize(minSize), out);
        return out;
    }

    private static List<FileNode> searchByName(FileNode root, String sub) {
        List<FileNode> out = new ArrayList<>();
        dfs(root, new SearchByName(sub), out);
        return out;
    }

    public static void reset_service() { root = null; }

    public static void fs_build_default_tree() {
        FileNode mainCpp = new FileNode("main.cpp", 50, "cpp", false);
        FileNode utilsCpp = new FileNode("utils.cpp", 120, "cpp", false);
        FileNode helperH = new FileNode("helper.h", 10, "h", false);
        FileNode readme = new FileNode("readme.md", 5, "md", false);
        FileNode report = new FileNode("report.pdf", 200, "pdf", false);
        FileNode buildSh = new FileNode("build.sh", 2, "sh", false);
        FileNode src = new FileNode("src", 0, "", true,
                new ArrayList<>(Arrays.asList(mainCpp, utilsCpp, helperH)));
        FileNode docs = new FileNode("docs", 0, "", true,
                new ArrayList<>(Arrays.asList(readme, report)));
        root = new FileNode("project", 0, "", true,
                new ArrayList<>(Arrays.asList(src, docs, buildSh)));
    }

    public static void fs_build_empty_tree() {
        root = new FileNode("empty", 0, "", true);
    }

    public static void fs_build_single_file_tree() {
        FileNode f = new FileNode("test.txt", 30, "txt", false);
        root = new FileNode("root", 0, "", true,
                new ArrayList<>(Arrays.asList(f)));
    }

    public static int fs_count_by_extension(String ext) {
        return root == null ? 0 : searchByExtension(root, ext).size();
    }

    public static boolean fs_has_by_extension(String ext, String name) {
        if (root == null) return false;
        for (FileNode f : searchByExtension(root, ext)) if (f.name.equals(name)) return true;
        return false;
    }

    public static int fs_count_by_size(int minSize) {
        return root == null ? 0 : searchBySize(root, minSize).size();
    }

    public static boolean fs_has_by_size(int minSize, String name) {
        if (root == null) return false;
        for (FileNode f : searchBySize(root, minSize)) if (f.name.equals(name)) return true;
        return false;
    }

    public static int fs_count_by_name(String sub) {
        return root == null ? 0 : searchByName(root, sub).size();
    }

    public static boolean fs_has_by_name(String sub, String name) {
        if (root == null) return false;
        for (FileNode f : searchByName(root, sub)) if (f.name.equals(name)) return true;
        return false;
    }

    public static int fs_count_composite_and(String ext, int minSize) {
        if (root == null) return 0;
        List<SearchCriteria> criteria = new ArrayList<>();
        criteria.add(new SearchByExtension(ext));
        criteria.add(new SearchByMinSize(minSize));
        List<FileNode> out = new ArrayList<>();
        dfs(root, new AndFilter(criteria), out);
        return out.size();
    }

    public static int fs_count_composite_or(String ext, int minSize) {
        if (root == null) return 0;
        List<SearchCriteria> criteria = new ArrayList<>();
        criteria.add(new SearchByExtension(ext));
        criteria.add(new SearchByMinSize(minSize));
        List<FileNode> out = new ArrayList<>();
        dfs(root, new OrFilter(criteria), out);
        return out.size();
    }

    public static String fs_first_sorted_by(String ext, String sortBy) {
        if (root == null) return "";
        List<FileNode> v = new ArrayList<>();
        dfs(root, new SearchByExtension(ext), v);
        if (sortBy.equals("name")) v.sort(Comparator.comparing(f -> f.name));
        else if (sortBy.equals("size")) v.sort((a, b) -> Integer.compare(b.size, a.size));
        else if (sortBy.equals("extension")) v.sort(Comparator.comparing(f -> f.extension));
        return v.isEmpty() ? "" : v.get(0).name;
    }
}

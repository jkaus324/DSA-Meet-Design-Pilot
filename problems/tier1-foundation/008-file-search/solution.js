'use strict';

// File Search — Strategy + Composite + Decorator reference solution (JavaScript).

class FileNode {
  constructor(name, size, extension, isDirectory, children = null) {
    this.name = name;
    this.size = size;
    this.extension = extension;
    this.isDirectory = isDirectory;
    this.children = children ? children : [];
  }
}

class SearchByExtension {
  constructor(ext) { this.ext = ext; }
  matches(file) { return !file.isDirectory && file.extension === this.ext; }
}

class SearchByMinSize {
  constructor(minSize) { this.minSize = minSize; }
  matches(file) { return !file.isDirectory && file.size >= this.minSize; }
}

class SearchByName {
  constructor(substring) { this.substring = substring; }
  matches(file) { return file.name.includes(this.substring); }
}

class AndFilter {
  constructor(criteria) { this.criteria = criteria; }
  matches(file) { return this.criteria.every(c => c.matches(file)); }
}

class OrFilter {
  constructor(criteria) { this.criteria = criteria; }
  matches(file) { return this.criteria.some(c => c.matches(file)); }
}

function _dfs(node, criteria, results) {
  if (node === null || node === undefined) return;
  if (!node.isDirectory && criteria.matches(node)) results.push(node);
  for (const c of node.children) _dfs(c, criteria, results);
}

function search_by_extension(root, ext) {
  const out = [];
  _dfs(root, new SearchByExtension(ext), out);
  return out;
}

function search_by_size(root, minSize) {
  const out = [];
  _dfs(root, new SearchByMinSize(minSize), out);
  return out;
}

function search_by_name(root, substring) {
  const out = [];
  _dfs(root, new SearchByName(substring), out);
  return out;
}

function search_composite(root, criteria, mode) {
  const f = mode === 'AND' ? new AndFilter(criteria) : new OrFilter(criteria);
  const out = [];
  _dfs(root, f, out);
  return out;
}

function search_and_sort(root, criteria, sortBy) {
  const out = [];
  _dfs(root, criteria, out);
  if (sortBy === 'name') {
    out.sort((a, b) => (a.name < b.name ? -1 : a.name > b.name ? 1 : 0));
  } else if (sortBy === 'size') {
    out.sort((a, b) => b.size - a.size);
  } else if (sortBy === 'extension') {
    out.sort((a, b) => (a.extension < b.extension ? -1 : a.extension > b.extension ? 1 : 0));
  }
  return out;
}

// ─── Module fixture state ───────────────────────────────────────────────────

let _g_root = null;

function reset_service() {
  _g_root = null;
}

function fs_build_default_tree() {
  const main_cpp = new FileNode('main.cpp', 50, 'cpp', false);
  const utils_cpp = new FileNode('utils.cpp', 120, 'cpp', false);
  const helper_h = new FileNode('helper.h', 10, 'h', false);
  const readme = new FileNode('readme.md', 5, 'md', false);
  const report = new FileNode('report.pdf', 200, 'pdf', false);
  const build_sh = new FileNode('build.sh', 2, 'sh', false);
  const src = new FileNode('src', 0, '', true, [main_cpp, utils_cpp, helper_h]);
  const docs = new FileNode('docs', 0, '', true, [readme, report]);
  _g_root = new FileNode('project', 0, '', true, [src, docs, build_sh]);
}

function fs_build_empty_tree() {
  _g_root = new FileNode('empty', 0, '', true);
}

function fs_build_single_file_tree() {
  const f = new FileNode('test.txt', 30, 'txt', false);
  _g_root = new FileNode('root', 0, '', true, [f]);
}

function fs_count_by_extension(ext) {
  return _g_root ? search_by_extension(_g_root, ext).length : 0;
}

function fs_has_by_extension(ext, name) {
  if (_g_root === null) return false;
  return search_by_extension(_g_root, ext).some(f => f.name === name);
}

function fs_count_by_size(minSize) {
  return _g_root ? search_by_size(_g_root, minSize).length : 0;
}

function fs_has_by_size(minSize, name) {
  if (_g_root === null) return false;
  return search_by_size(_g_root, minSize).some(f => f.name === name);
}

function fs_count_by_name(sub) {
  return _g_root ? search_by_name(_g_root, sub).length : 0;
}

function fs_has_by_name(sub, name) {
  if (_g_root === null) return false;
  return search_by_name(_g_root, sub).some(f => f.name === name);
}

function fs_count_composite_and(ext, minSize) {
  if (_g_root === null) return 0;
  return search_composite(_g_root, [new SearchByExtension(ext), new SearchByMinSize(minSize)], 'AND').length;
}

function fs_count_composite_or(ext, minSize) {
  if (_g_root === null) return 0;
  return search_composite(_g_root, [new SearchByExtension(ext), new SearchByMinSize(minSize)], 'OR').length;
}

function fs_first_sorted_by(ext, sortBy) {
  if (_g_root === null) return '';
  const v = search_and_sort(_g_root, new SearchByExtension(ext), sortBy);
  return v.length ? v[0].name : '';
}

module.exports = {
  FileNode,
  SearchByExtension,
  SearchByMinSize,
  SearchByName,
  AndFilter,
  OrFilter,
  search_by_extension,
  search_by_size,
  search_by_name,
  search_composite,
  search_and_sort,
  reset_service,
  fs_build_default_tree,
  fs_build_empty_tree,
  fs_build_single_file_tree,
  fs_count_by_extension,
  fs_has_by_extension,
  fs_count_by_size,
  fs_has_by_size,
  fs_count_by_name,
  fs_has_by_name,
  fs_count_composite_and,
  fs_count_composite_or,
  fs_first_sorted_by,
};

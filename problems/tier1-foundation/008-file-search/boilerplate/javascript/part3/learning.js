// Data class (given — do not modify).

// Strategy — base strategy. Subclasses implement compare().
class Strategy {
  // Return true iff `a` ranks strictly before `b`.
  compare(a, b) { throw new Error('not implemented'); }
}

function reset_service() {
  // TODO: implement this
  return null;
}

function fs_build_default_tree() {
  // TODO: implement this
  return null;
}

function fs_build_empty_tree() {
  // TODO: implement this
  return null;
}

function fs_build_single_file_tree() {
  // TODO: implement this
  return null;
}

function fs_count_by_extension(ext) {
  // TODO: implement this
  return null;
}

function fs_has_by_extension(ext, name) {
  // TODO: implement this
  return null;
}

function fs_count_by_size(minSize) {
  // TODO: implement this
  return null;
}

function fs_has_by_size(minSize, name) {
  // TODO: implement this
  return null;
}

function fs_count_by_name(sub) {
  // TODO: implement this
  return null;
}

function fs_has_by_name(sub, name) {
  // TODO: implement this
  return null;
}

function fs_count_composite_and(ext, minSize) {
  // TODO: implement this
  return null;
}

function fs_count_composite_or(ext, minSize) {
  // TODO: implement this
  return null;
}

function fs_first_sorted_by(ext, sortBy) {
  // TODO: implement this
  return null;
}

// ── Export everything the test runner needs (do not remove) ──
module.exports = { reset_service, fs_build_default_tree, fs_build_empty_tree, fs_build_single_file_tree, fs_count_by_extension, fs_has_by_extension, fs_count_by_size, fs_has_by_size, fs_count_by_name, fs_has_by_name, fs_count_composite_and, fs_count_composite_or, fs_first_sorted_by };

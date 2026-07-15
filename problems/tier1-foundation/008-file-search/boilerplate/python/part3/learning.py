# Data class (given — do not modify).

from abc import ABC, abstractmethod

class Strategy(ABC):
    @abstractmethod
    def compare(self, a, b):
        """Return True iff `a` ranks strictly before `b`."""

def reset_service():
    # TODO: implement this
    return None

def fs_build_default_tree():
    # TODO: implement this
    return None

def fs_build_empty_tree():
    # TODO: implement this
    return None

def fs_build_single_file_tree():
    # TODO: implement this
    return None

def fs_count_by_extension(ext):
    # TODO: implement this
    return None

def fs_has_by_extension(ext, name):
    # TODO: implement this
    return None

def fs_count_by_size(minSize):
    # TODO: implement this
    return None

def fs_has_by_size(minSize, name):
    # TODO: implement this
    return None

def fs_count_by_name(sub):
    # TODO: implement this
    return None

def fs_has_by_name(sub, name):
    # TODO: implement this
    return None

def fs_count_composite_and(ext, minSize):
    # TODO: implement this
    return None

def fs_count_composite_or(ext, minSize):
    # TODO: implement this
    return None

def fs_first_sorted_by(ext, sortBy):
    # TODO: implement this
    return None

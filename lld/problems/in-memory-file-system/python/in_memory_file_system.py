"""In-Memory File System LLD — Python reference implementation.

Composite tree of Directory/File nodes, path resolution ("/", "." and ".."),
and the usual mkdir/create/write/read/ls/rm/cd operations. See ../README.md.
"""
from __future__ import annotations

import threading
from dataclasses import dataclass, field
from typing import Dict, List


class NotFoundError(Exception):
    pass


class AlreadyExistsError(Exception):
    pass


class NotDirectoryError(Exception):
    pass


class IsDirectoryError(Exception):
    pass


class DirectoryNotEmptyError(Exception):
    pass


class File:
    def __init__(self, name: str, content: str = ""):
        self.name = name
        self._content = content
        self._lock = threading.Lock()

    def is_dir(self) -> bool:
        return False

    def read(self) -> str:
        with self._lock:
            return self._content

    def write(self, content: str, append: bool) -> None:
        with self._lock:
            self._content = self._content + content if append else content


class Directory:
    def __init__(self, name: str):
        self.name = name
        self.children: Dict[str, "File | Directory"] = {}

    def is_dir(self) -> bool:
        return True


@dataclass
class Entry:
    name: str
    is_dir: bool


class FileSystem:
    def __init__(self):
        self._root = Directory("/")
        self._cwd: List[str] = []
        self._lock = threading.Lock()

    def _split_path(self, path: str, base: List[str]) -> List[str]:
        if not path:
            raise ValueError("invalid path")
        parts: List[str] = [] if path.startswith("/") else list(base)
        for seg in path.split("/"):
            if seg in ("", "."):
                continue
            if seg == "..":
                if parts:
                    parts.pop()
            else:
                parts.append(seg)
        return parts

    def _resolve_parent(self, parts: List[str]) -> Directory:
        node = self._root
        for seg in parts:
            child = node.children.get(seg)
            if child is None:
                raise NotFoundError(seg)
            if not child.is_dir():
                raise NotDirectoryError(seg)
            node = child
        return node

    def _resolve(self, parts: List[str]):
        if not parts:
            return self._root
        parent = self._resolve_parent(parts[:-1])
        node = parent.children.get(parts[-1])
        if node is None:
            raise NotFoundError(parts[-1])
        return node

    def _cwd_snapshot(self) -> List[str]:
        with self._lock:
            return list(self._cwd)

    def mkdir(self, path: str) -> None:
        parts = self._split_path(path, self._cwd_snapshot())
        if not parts:
            raise AlreadyExistsError(path)
        with self._lock:
            parent = self._resolve_parent(parts[:-1])
            name = parts[-1]
            if name in parent.children:
                raise AlreadyExistsError(path)
            parent.children[name] = Directory(name)

    def create_file(self, path: str, content: str = "") -> None:
        parts = self._split_path(path, self._cwd_snapshot())
        if not parts:
            raise IsDirectoryError(path)
        with self._lock:
            parent = self._resolve_parent(parts[:-1])
            name = parts[-1]
            if name in parent.children:
                raise AlreadyExistsError(path)
            parent.children[name] = File(name, content)

    def write_file(self, path: str, content: str, append: bool = False) -> None:
        parts = self._split_path(path, self._cwd_snapshot())
        with self._lock:
            node = self._resolve(parts)
        if not isinstance(node, File):
            raise IsDirectoryError(path)
        node.write(content, append)

    def read_file(self, path: str) -> str:
        parts = self._split_path(path, self._cwd_snapshot())
        with self._lock:
            node = self._resolve(parts)
        if not isinstance(node, File):
            raise IsDirectoryError(path)
        return node.read()

    def ls(self, path: str) -> List[Entry]:
        parts = self._split_path(path, self._cwd_snapshot())
        with self._lock:
            node = self._resolve(parts)
            if not isinstance(node, Directory):
                raise NotDirectoryError(path)
            entries = [Entry(child.name, child.is_dir()) for child in node.children.values()]
        entries.sort(key=lambda e: e.name)
        return entries

    def rm(self, path: str, recursive: bool = False) -> None:
        parts = self._split_path(path, self._cwd_snapshot())
        if not parts:
            raise ValueError("cannot remove root")
        with self._lock:
            parent = self._resolve_parent(parts[:-1])
            name = parts[-1]
            node = parent.children.get(name)
            if node is None:
                raise NotFoundError(path)
            if isinstance(node, Directory) and node.children and not recursive:
                raise DirectoryNotEmptyError(path)
            del parent.children[name]

    def cd(self, path: str) -> None:
        parts = self._split_path(path, self._cwd_snapshot())
        with self._lock:
            node = self._resolve(parts)
            if not node.is_dir():
                raise NotDirectoryError(path)
            self._cwd = parts

    def pwd(self) -> str:
        parts = self._cwd_snapshot()
        return "/" + "/".join(parts) if parts else "/"


def _demo() -> None:
    fs = FileSystem()
    fs.mkdir("/home")
    fs.mkdir("/home/docs")
    fs.create_file("/home/docs/readme.md", "hello world")

    print("pwd:", fs.pwd())
    fs.cd("/home/docs")
    print("pwd after cd:", fs.pwd())
    print("readme.md:", fs.read_file("readme.md"))

    fs.write_file("readme.md", "\nmore content", append=True)
    print("readme.md after append:", fs.read_file("readme.md"))

    entries = fs.ls("/home")
    print("ls /home:", [f"{e.name}{'/' if e.is_dir else ''}" for e in entries])


if __name__ == "__main__":
    _demo()

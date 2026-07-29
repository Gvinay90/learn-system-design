import threading

import pytest

from in_memory_file_system import (
    AlreadyExistsError,
    DirectoryNotEmptyError,
    FileSystem,
    IsDirectoryError,
    NotDirectoryError,
    NotFoundError,
)


def test_mkdir_and_ls():
    fs = FileSystem()
    fs.mkdir("/home")
    fs.mkdir("/home/docs")
    fs.create_file("/home/notes.txt", "hello")

    entries = fs.ls("/home")
    assert [e.name for e in entries] == ["docs", "notes.txt"]
    assert entries[0].is_dir is True
    assert entries[1].is_dir is False


def test_mkdir_missing_parent_fails():
    fs = FileSystem()
    with pytest.raises(NotFoundError):
        fs.mkdir("/a/b")


def test_mkdir_duplicate_fails():
    fs = FileSystem()
    fs.mkdir("/home")
    with pytest.raises(AlreadyExistsError):
        fs.mkdir("/home")


def test_write_read_and_append():
    fs = FileSystem()
    fs.create_file("/a.txt", "hello")
    assert fs.read_file("/a.txt") == "hello"

    fs.write_file("/a.txt", " world", append=True)
    assert fs.read_file("/a.txt") == "hello world"

    fs.write_file("/a.txt", "reset", append=False)
    assert fs.read_file("/a.txt") == "reset"


def test_cd_and_relative_paths():
    fs = FileSystem()
    fs.mkdir("/home")
    fs.mkdir("/home/docs")
    fs.create_file("/home/docs/readme.md", "hi")

    fs.cd("/home/docs")
    assert fs.pwd() == "/home/docs"
    assert fs.read_file("readme.md") == "hi"

    fs.cd("..")
    assert fs.pwd() == "/home"


def test_rm_non_empty_dir_requires_recursive():
    fs = FileSystem()
    fs.mkdir("/home")
    fs.create_file("/home/a.txt", "x")

    with pytest.raises(DirectoryNotEmptyError):
        fs.rm("/home", recursive=False)

    fs.rm("/home", recursive=True)
    with pytest.raises(NotFoundError):
        fs.ls("/home")


def test_read_on_directory_raises():
    fs = FileSystem()
    fs.mkdir("/home")
    with pytest.raises(IsDirectoryError):
        fs.read_file("/home")


def test_ls_on_file_raises_not_directory():
    fs = FileSystem()
    fs.create_file("/a.txt", "x")
    with pytest.raises(NotDirectoryError):
        fs.ls("/a.txt")


def test_concurrent_writes_to_same_file():
    fs = FileSystem()
    fs.create_file("/log.txt", "")
    n = 100

    def worker():
        fs.write_file("/log.txt", "x", append=True)

    threads = [threading.Thread(target=worker) for _ in range(n)]
    for t in threads:
        t.start()
    for t in threads:
        t.join()

    assert len(fs.read_file("/log.txt")) == n

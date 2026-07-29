package inmemoryfilesystem

import (
	"sync"
	"testing"
)

func TestMkdirAndLs(t *testing.T) {
	fs := NewFileSystem()
	if err := fs.Mkdir("/home"); err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if err := fs.Mkdir("/home/docs"); err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if err := fs.CreateFile("/home/notes.txt", "hello"); err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	entries, err := fs.Ls("/home")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if len(entries) != 2 || entries[0].Name != "docs" || entries[1].Name != "notes.txt" {
		t.Fatalf("unexpected entries: %+v", entries)
	}
	if !entries[0].IsDir || entries[1].IsDir {
		t.Fatalf("wrong IsDir flags: %+v", entries)
	}
}

func TestMkdirMissingParentFails(t *testing.T) {
	fs := NewFileSystem()
	if err := fs.Mkdir("/a/b"); err != ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestWriteReadAndAppend(t *testing.T) {
	fs := NewFileSystem()
	_ = fs.CreateFile("/a.txt", "hello")

	content, err := fs.ReadFile("/a.txt")
	if err != nil || content != "hello" {
		t.Fatalf("expected hello, got %q err %v", content, err)
	}

	if err := fs.WriteFile("/a.txt", " world", true); err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	content, _ = fs.ReadFile("/a.txt")
	if content != "hello world" {
		t.Fatalf("expected 'hello world', got %q", content)
	}

	if err := fs.WriteFile("/a.txt", "reset", false); err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	content, _ = fs.ReadFile("/a.txt")
	if content != "reset" {
		t.Fatalf("expected 'reset', got %q", content)
	}
}

func TestCdAndRelativePaths(t *testing.T) {
	fs := NewFileSystem()
	_ = fs.Mkdir("/home")
	_ = fs.Mkdir("/home/docs")
	_ = fs.CreateFile("/home/docs/readme.md", "hi")

	if err := fs.Cd("/home/docs"); err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if fs.Pwd() != "/home/docs" {
		t.Fatalf("expected /home/docs, got %s", fs.Pwd())
	}

	content, err := fs.ReadFile("readme.md")
	if err != nil || content != "hi" {
		t.Fatalf("expected hi, got %q err %v", content, err)
	}

	if err := fs.Cd(".."); err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if fs.Pwd() != "/home" {
		t.Fatalf("expected /home after cd .., got %s", fs.Pwd())
	}
}

func TestRmFileAndNonEmptyDir(t *testing.T) {
	fs := NewFileSystem()
	_ = fs.Mkdir("/home")
	_ = fs.CreateFile("/home/a.txt", "x")

	if err := fs.Rm("/home", false); err != ErrDirNotEmpty {
		t.Fatalf("expected ErrDirNotEmpty, got %v", err)
	}
	if err := fs.Rm("/home/a.txt", false); err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if err := fs.Rm("/home", false); err != nil {
		t.Fatalf("expected empty dir to be removable, got %v", err)
	}
	if _, err := fs.Ls("/"); err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
}

func TestRmRecursive(t *testing.T) {
	fs := NewFileSystem()
	_ = fs.Mkdir("/home")
	_ = fs.CreateFile("/home/a.txt", "x")

	if err := fs.Rm("/home", true); err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if _, err := fs.Ls("/home"); err != ErrNotFound {
		t.Fatalf("expected /home to be gone, got %v", err)
	}
}

// TestConcurrentWritesToSameFile appends from many goroutines and asserts no
// bytes are dropped/interleaved — the per-file mutex in File.write must serialize them.
func TestConcurrentWritesToSameFile(t *testing.T) {
	fs := NewFileSystem()
	_ = fs.CreateFile("/log.txt", "")

	var wg sync.WaitGroup
	const n = 100
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = fs.WriteFile("/log.txt", "x", true)
		}()
	}
	wg.Wait()

	content, _ := fs.ReadFile("/log.txt")
	if len(content) != n {
		t.Fatalf("expected %d bytes, got %d (%q)", n, len(content), content)
	}
}

// TestConcurrentMkdirSameName asserts only one of two racing mkdir calls
// for the same path succeeds.
func TestConcurrentMkdirSameName(t *testing.T) {
	fs := NewFileSystem()
	var wg sync.WaitGroup
	results := make(chan error, 2)

	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			results <- fs.Mkdir("/dup")
		}()
	}
	wg.Wait()
	close(results)

	successes := 0
	for err := range results {
		if err == nil {
			successes++
		}
	}
	if successes != 1 {
		t.Fatalf("expected exactly 1 successful mkdir, got %d", successes)
	}
}

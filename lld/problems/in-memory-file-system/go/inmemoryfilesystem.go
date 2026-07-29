// Package inmemoryfilesystem implements an in-memory hierarchical file system:
// directories and files as a Composite tree, path resolution (absolute/relative,
// "." and ".."), and the usual mkdir/create/write/read/ls/rm/cd operations.
package inmemoryfilesystem

import (
	"errors"
	"sort"
	"strings"
	"sync"
)

var (
	ErrNotFound    = errors.New("no such file or directory")
	ErrExists      = errors.New("file or directory already exists")
	ErrNotDir      = errors.New("not a directory")
	ErrIsDir       = errors.New("is a directory")
	ErrDirNotEmpty = errors.New("directory not empty")
	ErrInvalidPath = errors.New("invalid path")
)

// Node is the Composite element shared by File and Directory.
type Node interface {
	Name() string
	IsDir() bool
}

// File holds byte content behind its own mutex so concurrent writers to the
// same file (e.g. two appenders) never interleave partial writes.
type File struct {
	name    string
	mu      sync.Mutex
	content []byte
}

func (f *File) Name() string { return f.name }
func (f *File) IsDir() bool  { return false }

func (f *File) read() string {
	f.mu.Lock()
	defer f.mu.Unlock()
	return string(f.content)
}

func (f *File) write(content string, append bool) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if append {
		f.content = append2(f.content, content)
		return
	}
	f.content = []byte(content)
}

func append2(b []byte, s string) []byte {
	return append(b, s...)
}

// Directory is the Composite container: a name-indexed map of child nodes.
type Directory struct {
	name     string
	mu       sync.RWMutex
	children map[string]Node
}

func newDirectory(name string) *Directory {
	return &Directory{name: name, children: make(map[string]Node)}
}

func (d *Directory) Name() string { return d.name }
func (d *Directory) IsDir() bool  { return true }

// Entry describes one child for Ls output.
type Entry struct {
	Name  string
	IsDir bool
}

// FileSystem is the facade clients interact with; it owns the root directory
// and the current working directory (as resolved path components).
type FileSystem struct {
	root *Directory
	mu   sync.Mutex
	cwd  []string
}

func NewFileSystem() *FileSystem {
	return &FileSystem{root: newDirectory("/"), cwd: []string{}}
}

// splitPath turns a path string into clean components, resolving "." and "..".
// A leading "/" makes it absolute; otherwise it is resolved against base.
func splitPath(path string, base []string) ([]string, error) {
	if path == "" {
		return nil, ErrInvalidPath
	}
	var parts []string
	if strings.HasPrefix(path, "/") {
		parts = []string{}
	} else {
		parts = append(parts, base...)
	}
	for _, seg := range strings.Split(path, "/") {
		switch seg {
		case "", ".":
			continue
		case "..":
			if len(parts) > 0 {
				parts = parts[:len(parts)-1]
			}
		default:
			parts = append(parts, seg)
		}
	}
	return parts, nil
}

// resolveDir walks parts[:len-1] returning the parent directory node, and
// resolveAll walks the full parts returning the final node.
func (fs *FileSystem) resolveParent(parts []string) (*Directory, error) {
	dir := fs.root
	for _, seg := range parts {
		dir.mu.RLock()
		child, ok := dir.children[seg]
		dir.mu.RUnlock()
		if !ok {
			return nil, ErrNotFound
		}
		sub, isDir := child.(*Directory)
		if !isDir {
			return nil, ErrNotDir
		}
		dir = sub
	}
	return dir, nil
}

func (fs *FileSystem) resolve(parts []string) (Node, error) {
	if len(parts) == 0 {
		return fs.root, nil
	}
	parent, err := fs.resolveParent(parts[:len(parts)-1])
	if err != nil {
		return nil, err
	}
	name := parts[len(parts)-1]
	parent.mu.RLock()
	defer parent.mu.RUnlock()
	node, ok := parent.children[name]
	if !ok {
		return nil, ErrNotFound
	}
	return node, nil
}

func (fs *FileSystem) currentCwd() []string {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	out := make([]string, len(fs.cwd))
	copy(out, fs.cwd)
	return out
}

// Mkdir creates a single directory; the parent must already exist.
func (fs *FileSystem) Mkdir(path string) error {
	parts, err := splitPath(path, fs.currentCwd())
	if err != nil {
		return err
	}
	if len(parts) == 0 {
		return ErrExists
	}
	parent, err := fs.resolveParent(parts[:len(parts)-1])
	if err != nil {
		return err
	}
	name := parts[len(parts)-1]

	parent.mu.Lock()
	defer parent.mu.Unlock()
	if _, exists := parent.children[name]; exists {
		return ErrExists
	}
	parent.children[name] = newDirectory(name)
	return nil
}

// CreateFile creates a new empty (or seeded) file; the parent must exist.
func (fs *FileSystem) CreateFile(path string, content string) error {
	parts, err := splitPath(path, fs.currentCwd())
	if err != nil {
		return err
	}
	if len(parts) == 0 {
		return ErrIsDir
	}
	parent, err := fs.resolveParent(parts[:len(parts)-1])
	if err != nil {
		return err
	}
	name := parts[len(parts)-1]

	parent.mu.Lock()
	defer parent.mu.Unlock()
	if _, exists := parent.children[name]; exists {
		return ErrExists
	}
	parent.children[name] = &File{name: name, content: []byte(content)}
	return nil
}

// WriteFile overwrites (or appends to, if append is true) an existing file's content.
func (fs *FileSystem) WriteFile(path string, content string, append bool) error {
	parts, err := splitPath(path, fs.currentCwd())
	if err != nil {
		return err
	}
	node, err := fs.resolve(parts)
	if err != nil {
		return err
	}
	file, ok := node.(*File)
	if !ok {
		return ErrIsDir
	}
	file.write(content, append)
	return nil
}

// ReadFile returns a file's full content.
func (fs *FileSystem) ReadFile(path string) (string, error) {
	parts, err := splitPath(path, fs.currentCwd())
	if err != nil {
		return "", err
	}
	node, err := fs.resolve(parts)
	if err != nil {
		return "", err
	}
	file, ok := node.(*File)
	if !ok {
		return "", ErrIsDir
	}
	return file.read(), nil
}

// Ls lists the direct children of a directory, sorted by name.
func (fs *FileSystem) Ls(path string) ([]Entry, error) {
	parts, err := splitPath(path, fs.currentCwd())
	if err != nil {
		return nil, err
	}
	node, err := fs.resolve(parts)
	if err != nil {
		return nil, err
	}
	dir, ok := node.(*Directory)
	if !ok {
		return nil, ErrNotDir
	}
	dir.mu.RLock()
	defer dir.mu.RUnlock()
	entries := make([]Entry, 0, len(dir.children))
	for name, child := range dir.children {
		entries = append(entries, Entry{Name: name, IsDir: child.IsDir()})
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].Name < entries[j].Name })
	return entries, nil
}

// Rm removes a file, or a directory (only if empty, unless recursive is true).
func (fs *FileSystem) Rm(path string, recursive bool) error {
	parts, err := splitPath(path, fs.currentCwd())
	if err != nil {
		return err
	}
	if len(parts) == 0 {
		return errors.New("cannot remove root")
	}
	parent, err := fs.resolveParent(parts[:len(parts)-1])
	if err != nil {
		return err
	}
	name := parts[len(parts)-1]

	parent.mu.Lock()
	defer parent.mu.Unlock()
	node, ok := parent.children[name]
	if !ok {
		return ErrNotFound
	}
	if dir, isDir := node.(*Directory); isDir {
		dir.mu.RLock()
		empty := len(dir.children) == 0
		dir.mu.RUnlock()
		if !empty && !recursive {
			return ErrDirNotEmpty
		}
	}
	delete(parent.children, name)
	return nil
}

// Cd changes the current working directory.
func (fs *FileSystem) Cd(path string) error {
	parts, err := splitPath(path, fs.currentCwd())
	if err != nil {
		return err
	}
	node, err := fs.resolve(parts)
	if err != nil {
		return err
	}
	if !node.IsDir() {
		return ErrNotDir
	}
	fs.mu.Lock()
	fs.cwd = parts
	fs.mu.Unlock()
	return nil
}

// Pwd returns the current working directory as an absolute path string.
func (fs *FileSystem) Pwd() string {
	parts := fs.currentCwd()
	if len(parts) == 0 {
		return "/"
	}
	return "/" + strings.Join(parts, "/")
}

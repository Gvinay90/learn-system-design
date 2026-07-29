// Package searchautocomplete demonstrates the core in-memory data structure
// behind a typeahead/autocomplete service: a Trie where each node caches its
// own top-K most frequent completions, so a prefix lookup is O(prefix length)
// instead of a tree walk + sort at query time. A real system builds this
// offline from query logs and serves it read-only from many replicas — see
// the README for the full distributed design.
package searchautocomplete

import (
	"sort"
	"sync"
)

// suggestion is a single ranked completion candidate.
type suggestion struct {
	Word      string
	Frequency int
}

// trieNode is one character position in the trie. topK is a cache of this
// node's best completions, kept sorted (desc frequency, then lexicographic)
// and capped at the trie's configured K, so PrefixSearch never has to touch
// nodes below this point in the tree.
type trieNode struct {
	children map[byte]*trieNode
	topK     []suggestion
}

func newTrieNode() *trieNode {
	return &trieNode{children: make(map[byte]*trieNode)}
}

// Trie is a thread-safe prefix tree supporting Insert and PrefixSearch
// (top-K suggestions ranked by frequency).
type Trie struct {
	mu   sync.RWMutex
	root *trieNode
	k    int
}

// NewTrie creates a Trie that keeps the top k suggestions cached at every
// node. k must be >= 1.
func NewTrie(k int) *Trie {
	if k < 1 {
		k = 1
	}
	return &Trie{root: newTrieNode(), k: k}
}

// Insert adds word with the given frequency (e.g. number of times it has been
// searched). Calling Insert again for the same word increments its frequency
// rather than duplicating it.
func (t *Trie) Insert(word string, frequency int) {
	if word == "" {
		return
	}
	t.mu.Lock()
	defer t.mu.Unlock()

	node := t.root
	t.upsert(node, word, frequency)
	for i := 0; i < len(word); i++ {
		c := word[i]
		child, ok := node.children[c]
		if !ok {
			child = newTrieNode()
			node.children[c] = child
		}
		node = child
		t.upsert(node, word, frequency)
	}
}

// upsert adds/updates word's frequency in node's topK cache, re-sorting and
// trimming to the configured k.
func (t *Trie) upsert(node *trieNode, word string, frequency int) {
	found := false
	for i := range node.topK {
		if node.topK[i].Word == word {
			node.topK[i].Frequency += frequency
			found = true
			break
		}
	}
	if !found {
		node.topK = append(node.topK, suggestion{Word: word, Frequency: frequency})
	}

	sort.Slice(node.topK, func(i, j int) bool {
		if node.topK[i].Frequency != node.topK[j].Frequency {
			return node.topK[i].Frequency > node.topK[j].Frequency
		}
		return node.topK[i].Word < node.topK[j].Word
	})

	if len(node.topK) > t.k {
		node.topK = node.topK[:t.k]
	}
}

// PrefixSearch returns up to k suggestions for prefix, ranked by frequency
// (ties broken lexicographically). Returns nil if no word in the trie has
// this prefix.
func (t *Trie) PrefixSearch(prefix string) []string {
	t.mu.RLock()
	defer t.mu.RUnlock()

	node := t.root
	for i := 0; i < len(prefix); i++ {
		child, ok := node.children[prefix[i]]
		if !ok {
			return nil
		}
		node = child
	}

	out := make([]string, len(node.topK))
	for i, s := range node.topK {
		out[i] = s.Word
	}
	return out
}

package searchautocomplete

import (
	"reflect"
	"sync"
	"testing"
)

func TestPrefixMatch(t *testing.T) {
	trie := NewTrie(5)
	trie.Insert("cat", 10)
	trie.Insert("car", 5)
	trie.Insert("dog", 20)

	got := trie.PrefixSearch("ca")
	want := []string{"cat", "car"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("PrefixSearch(%q) = %v, want %v", "ca", got, want)
	}
}

func TestTopKRankingByFrequency(t *testing.T) {
	trie := NewTrie(2)
	trie.Insert("apple", 3)
	trie.Insert("app", 50)
	trie.Insert("application", 10)
	trie.Insert("apt", 1) // different branch after "ap", shares prefix "ap" not "app"

	got := trie.PrefixSearch("app")
	want := []string{"app", "application"} // top 2 by frequency: 50, 10 > apple's 3
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("PrefixSearch(%q) = %v, want %v", "app", got, want)
	}
}

func TestRepeatedInsertAccumulatesFrequency(t *testing.T) {
	trie := NewTrie(3)
	trie.Insert("go", 1)
	trie.Insert("go", 1)
	trie.Insert("go", 1)
	trie.Insert("golang", 2)

	got := trie.PrefixSearch("go")
	want := []string{"go", "golang"} // "go" frequency accumulates to 3 > "golang" 2
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("PrefixSearch(%q) = %v, want %v", "go", got, want)
	}
}

func TestNoMatch(t *testing.T) {
	trie := NewTrie(5)
	trie.Insert("hello", 1)

	if got := trie.PrefixSearch("xyz"); got != nil {
		t.Fatalf("PrefixSearch(%q) = %v, want nil", "xyz", got)
	}
	if got := trie.PrefixSearch("hellox"); got != nil {
		t.Fatalf("PrefixSearch(%q) = %v, want nil (no word extends past 'hello')", "hellox", got)
	}
}

func TestEmptyPrefixReturnsGlobalTopK(t *testing.T) {
	trie := NewTrie(2)
	trie.Insert("zeta", 1)
	trie.Insert("alpha", 100)
	trie.Insert("beta", 50)

	got := trie.PrefixSearch("")
	want := []string{"alpha", "beta"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("PrefixSearch(\"\") = %v, want %v", got, want)
	}
}

func TestConcurrentInserts(t *testing.T) {
	trie := NewTrie(5)
	words := []string{"red", "green", "blue", "redwood", "reduce"}

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			for _, w := range words {
				trie.Insert(w, 1)
			}
			// Concurrent readers shouldn't race or panic either.
			_ = trie.PrefixSearch("re")
		}(i)
	}
	wg.Wait()

	got := trie.PrefixSearch("re")
	if len(got) != 3 {
		t.Fatalf("PrefixSearch(%q) after concurrent inserts = %v, want 3 words (red, redwood, reduce)", "re", got)
	}
	// "red" should have the highest frequency (50 goroutines x 1 each = 50).
	if got[0] != "red" {
		t.Fatalf("expected 'red' to rank first with frequency 50, got order %v", got)
	}
}

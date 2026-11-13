package trie

import (
	"maps"
	"math/rand/v2"
	"slices"
	"strings"
	"testing"
)

// commonWords are the 31 common English words of Knuth's searching examples (TAOCP Vol. 3, 2nd ed.,
// §6.2.2 and the trie of §6.3). The list was not re-checked against the book in this run; every expected
// result below is derived here by hand from the list itself, so no test depends on the book's figures.
var commonWords = []string{
	"A", "AND", "ARE", "AS", "AT", "BE", "BUT", "BY", "FOR", "FROM", "HAD", "HAVE", "HE", "HER", "HIS",
	"I", "IN", "IS", "IT", "NOT", "OF", "ON", "OR", "THAT", "THE", "THIS", "TO", "WAS", "WHICH", "WITH", "YOU",
}

// wordTrie maps each word to its position in the list.
func wordTrie(words []string) *Trie[int] {
	var t Trie[int]
	for i, w := range words {
		t.Insert(w, i)
	}
	return &t
}

// collect drains a key/value iterator.
func collect[V any](seq func(func(string, V) bool)) []string {
	var keys []string
	for k := range seq {
		keys = append(keys, k)
	}
	return keys
}

// countNodes counts the trie's nodes (root included) to check pruning.
func countNodes[V any](t *Trie[V]) int {
	n := 0
	stack := []*node[V]{&t.root}
	for len(stack) > 0 {
		nd := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		n++
		stack = append(stack, nd.kids...)
	}
	return n
}

func TestTrie_knuthCommonWords(t *testing.T) {
	tr := wordTrie(commonWords)
	if tr.Len() != 31 {
		t.Fatalf("Len = %d, want 31", tr.Len())
	}
	if v, ok := tr.Search("THE"); !ok || v != 24 {
		t.Errorf(`Search("THE") = (%d, %v), want (24, true)`, v, ok)
	}
	if _, ok := tr.Search("TH"); ok {
		t.Error(`Search("TH") found a key that was never stored`)
	}
	if got := collect(tr.WithPrefix("TH")); !slices.Equal(got, []string{"THAT", "THE", "THIS"}) {
		t.Errorf(`WithPrefix("TH") = %v`, got)
	}
	if got := collect(tr.WithPrefix("H")); !slices.Equal(got, []string{"HAD", "HAVE", "HE", "HER", "HIS"}) {
		t.Errorf(`WithPrefix("H") = %v`, got)
	}
	for p, want := range map[string]bool{"WH": true, "WHICH": true, "WHICHX": false, "Q": false, "": true} {
		if got := tr.HasPrefix(p); got != want {
			t.Errorf("HasPrefix(%q) = %v, want %v", p, got, want)
		}
	}
	// The full iteration is the list in byte order (the list is already sorted).
	if got := collect(tr.WithPrefix("")); !slices.Equal(got, commonWords) {
		t.Errorf("all keys = %v", got)
	}
}

func TestTrie_insertReplace(t *testing.T) {
	var tr Trie[string]
	if tr.Insert("go", "a") {
		t.Error("first Insert reported a replacement")
	}
	if !tr.Insert("go", "b") {
		t.Error("second Insert did not report a replacement")
	}
	if v, _ := tr.Search("go"); v != "b" || tr.Len() != 1 {
		t.Errorf("after replace: %q, Len %d", v, tr.Len())
	}
}

func TestTrie_emptyKey(t *testing.T) {
	var tr Trie[int]
	if tr.HasPrefix("") {
		t.Error("empty trie HasPrefix(\"\") = true")
	}
	tr.Insert("", 7)
	if v, ok := tr.Search(""); !ok || v != 7 || !tr.HasPrefix("") || tr.Len() != 1 {
		t.Errorf("empty key: (%d, %v), Len %d", v, ok, tr.Len())
	}
	if v, ok := tr.Delete(""); !ok || v != 7 || tr.Len() != 0 || tr.HasPrefix("") {
		t.Errorf("Delete(\"\") = (%d, %v), Len %d", v, ok, tr.Len())
	}
}

func TestTrie_bytesNotRunes(t *testing.T) {
	var tr Trie[int]
	tr.Insert("é", 1) // two bytes: 0xC3 0xA9
	if !tr.HasPrefix("\xc3") {
		t.Error("the first byte of a two-byte character is a prefix")
	}
	if _, ok := tr.Search("\xc3"); ok {
		t.Error("half a character is not a key")
	}
	tr.Insert("\x00\xff", 2)
	if got := collect(tr.WithPrefix("")); !slices.Equal(got, []string{"\x00\xff", "é"}) {
		t.Errorf("byte order = %q", got)
	}
}

func TestTrie_deletePrunes(t *testing.T) {
	var tr Trie[int]
	base := countNodes(&tr) // 1: the root
	tr.Insert("car", 1)
	tr.Insert("cart", 2)
	tr.Insert("cat", 3)
	if _, ok := tr.Delete("ca"); ok {
		t.Error("Delete of a prefix that is not a key succeeded")
	}
	if _, ok := tr.Delete("cartoon"); ok {
		t.Error("Delete of an absent extension succeeded")
	}
	if v, ok := tr.Delete("cart"); !ok || v != 2 {
		t.Errorf("Delete(cart) = (%d, %v)", v, ok)
	}
	// c-a-r and c-a-t remain: root + c + a + r + t = 5 nodes.
	if n := countNodes(&tr); n != 5 {
		t.Errorf("after deleting cart: %d nodes, want 5", n)
	}
	if tr.HasPrefix("cart") {
		t.Error("pruned branch still reports a prefix")
	}
	tr.Delete("car") // a leaf whose parent has another child: only r goes
	if n := countNodes(&tr); n != 4 {
		t.Errorf("after deleting car: %d nodes, want 4", n)
	}
	tr.Delete("cat")
	if n := countNodes(&tr); n != base || tr.Len() != 0 || tr.HasPrefix("c") {
		t.Errorf("after deleting all: %d nodes, Len %d", n, tr.Len())
	}
	// A key that is an inner node: its value goes, its children stay.
	tr.Insert("ab", 1)
	tr.Insert("abc", 2)
	tr.Delete("ab")
	if _, ok := tr.Search("abc"); !ok || !tr.HasPrefix("ab") {
		t.Error("deleting an inner key removed its extension")
	}
}

func TestTrie_nilAndZero(t *testing.T) {
	var nilTrie *Trie[int]
	if nilTrie.Insert("a", 1) || nilTrie.Len() != 0 || nilTrie.HasPrefix("") {
		t.Error("nil trie stored a key")
	}
	if _, ok := nilTrie.Search("a"); ok {
		t.Error("nil Search found a key")
	}
	if _, ok := nilTrie.Delete("a"); ok {
		t.Error("nil Delete succeeded")
	}
	if got := collect(nilTrie.WithPrefix("")); len(got) != 0 {
		t.Errorf("nil WithPrefix = %v", got)
	}
}

func TestTrie_withPrefixStopsEarly(t *testing.T) {
	tr := wordTrie(commonWords)
	n := 0
	for range tr.WithPrefix("") {
		n++
		if n == 3 {
			break
		}
	}
	if n != 3 {
		t.Errorf("visited %d keys, want 3", n)
	}
}

// Oracle: random inserts and deletes over a small alphabet against a Go map; every query checked.
func TestTrie_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	randKey := func() string {
		b := make([]byte, rng.IntN(6))
		for i := range b {
			b[i] = "abc"[rng.IntN(3)]
		}
		return string(b)
	}
	for trial := range 200 {
		var tr Trie[int]
		ref := map[string]int{}
		for step := range 80 {
			k := randKey()
			if rng.IntN(3) == 0 {
				v, ok := tr.Delete(k)
				rv, rok := ref[k]
				delete(ref, k)
				if ok != rok || v != rv {
					t.Fatalf("trial %d step %d Delete(%q) = (%d, %v), want (%d, %v)", trial, step, k, v, ok, rv, rok)
				}
			} else {
				_, had := ref[k]
				if tr.Insert(k, step) != had {
					t.Fatalf("trial %d Insert(%q) replaced mismatch", trial, k)
				}
				ref[k] = step
			}
			if tr.Len() != len(ref) {
				t.Fatalf("trial %d Len %d, want %d", trial, tr.Len(), len(ref))
			}
			q := randKey()
			v, ok := tr.Search(q)
			if rv, rok := ref[q]; ok != rok || v != rv {
				t.Fatalf("trial %d Search(%q) = (%d, %v), want (%d, %v)", trial, q, v, ok, rv, rok)
			}
			var want []string
			for _, key := range slices.Sorted(maps.Keys(ref)) {
				if strings.HasPrefix(key, q) {
					want = append(want, key)
				}
			}
			if got := collect(tr.WithPrefix(q)); !slices.Equal(got, want) {
				t.Fatalf("trial %d WithPrefix(%q) = %v, want %v", trial, q, got, want)
			}
			if tr.HasPrefix(q) != (len(want) > 0) {
				t.Fatalf("trial %d HasPrefix(%q) = %v", trial, q, !(len(want) > 0))
			}
		}
		// Pruning leaves exactly the nodes of the distinct prefixes of the remaining keys.
		prefixes := map[string]bool{"": true}
		for key := range ref {
			for i := 1; i <= len(key); i++ {
				prefixes[key[:i]] = true
			}
		}
		if n := countNodes(&tr); n != len(prefixes) {
			t.Fatalf("trial %d: %d nodes, want %d", trial, n, len(prefixes))
		}
	}
}

func TestTrie_searchNoAlloc(t *testing.T) {
	tr := wordTrie(commonWords)
	allocs := testing.AllocsPerRun(100, func() {
		_, _ = tr.Search("WHICH")
		_ = tr.HasPrefix("TH")
		_ = tr.Len()
	})
	if allocs != 0 {
		t.Errorf("Search/HasPrefix/Len allocated %v times, want 0", allocs)
	}
}

// benchSizes are the collection sizes of .claude/rules/go-complexity.md.
var benchSizes = []struct {
	name string
	n    int
}{{"1e3", 1_000}, {"1e4", 10_000}, {"1e5", 100_000}}

// benchWords returns n fixed-seed lowercase words of 3–12 bytes.
func benchWords(n int) []string {
	rng := rand.New(rand.NewPCG(1, 2))
	words := make([]string, n)
	for i := range words {
		b := make([]byte, 3+rng.IntN(10))
		for j := range b {
			b[j] = byte('a' + rng.IntN(26))
		}
		words[i] = string(b)
	}
	return words
}

func BenchmarkTrieInsert(b *testing.B) {
	for _, size := range benchSizes {
		words := benchWords(size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				var tr Trie[int]
				for i, w := range words {
					tr.Insert(w, i)
				}
			}
		})
	}
}

func BenchmarkTrieSearch(b *testing.B) {
	for _, size := range benchSizes {
		words := benchWords(size.n)
		tr := wordTrie(words)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				for _, w := range words {
					_, _ = tr.Search(w)
				}
			}
		})
	}
}

func BenchmarkTrieHasPrefix(b *testing.B) {
	for _, size := range benchSizes {
		words := benchWords(size.n)
		tr := wordTrie(words)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				for _, w := range words {
					_ = tr.HasPrefix(w[:2])
				}
			}
		})
	}
}

func BenchmarkTrieDelete(b *testing.B) {
	for _, size := range benchSizes {
		words := benchWords(size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				b.StopTimer()
				tr := wordTrie(words)
				b.StartTimer()
				for _, w := range words {
					tr.Delete(w)
				}
			}
		})
	}
}

func BenchmarkTrieLen(b *testing.B) {
	tr := wordTrie(commonWords)
	b.ReportAllocs()
	for b.Loop() {
		_ = tr.Len()
	}
}

// Full iteration: O(total nodes).
func BenchmarkTrieWithPrefix(b *testing.B) {
	for _, size := range benchSizes {
		tr := wordTrie(benchWords(size.n))
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				for range tr.WithPrefix("") {
				}
			}
		})
	}
}

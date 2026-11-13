package trie

import (
	"math/rand/v2"
	"slices"
	"strings"
	"testing"
)

// naiveShortestRoot scans every root.
func naiveShortestRoot(roots []string, word string) string {
	best, found := "", false
	for _, r := range roots {
		if r != "" && strings.HasPrefix(word, r) && (!found || len(r) < len(best)) {
			best, found = r, true
		}
	}
	if !found {
		return word
	}
	return best
}

// Knuth's common words (TAOCP Vol. 3, §6.3) as roots: "THESE" stops at THE, "HERO" at HE, and "ANDES"
// at A and "ISLAND" at I, not at AND or IS: the shortest root wins. Derived by hand.
func TestShortestPrefixOf_knuthCommonWords(t *testing.T) {
	tr := wordTrie(commonWords)
	tests := []struct {
		word, key string
		ok        bool
	}{
		{"THESE", "THE", true},
		{"ANDES", "A", true}, // A is a key and a prefix of ANDES: shortest wins
		{"ISLAND", "I", true},
		{"HERO", "HE", true},
		{"ZEBRA", "", false},
		{"T", "", false},
		{"", "", false},
		{"WHICH", "WHICH", true},
	}
	for _, tt := range tests {
		key, _, ok := tr.ShortestPrefixOf(tt.word)
		if key != tt.key || ok != tt.ok {
			t.Errorf("ShortestPrefixOf(%q) = (%q, %v), want (%q, %v)", tt.word, key, ok, tt.key, tt.ok)
		}
	}
	var nilTrie *Trie[int]
	if _, _, ok := nilTrie.ShortestPrefixOf("A"); ok {
		t.Error("nil trie found a prefix")
	}
}

func TestShortestPrefixOf_noAlloc(t *testing.T) {
	tr := wordTrie(commonWords)
	if allocs := testing.AllocsPerRun(100, func() { _, _, _ = tr.ShortestPrefixOf("THESE") }); allocs != 0 {
		t.Errorf("allocated %v times, want 0", allocs)
	}
}

func TestReplaceWordsByRoot_cases(t *testing.T) {
	tests := []struct {
		name         string
		roots, words []string
		want         []string
	}{
		{"empty", nil, nil, nil},
		{"no roots", nil, []string{"keep"}, []string{"keep"}},
		{"shortest root wins", []string{"walk", "wal", "w"}, []string{"walking"}, []string{"w"}},
		{"empty root ignored", []string{"", "un"}, []string{"undo", "redo"}, []string{"un", "redo"}},
		{"word equals root", []string{"run"}, []string{"run", "ru"}, []string{"run", "ru"}},
		{"duplicates", []string{"ab", "ab"}, []string{"abc", "abc"}, []string{"ab", "ab"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ReplaceWordsByRoot(nil, tt.roots, tt.words); !slices.Equal(got, tt.want) {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

// Problem-site shape, in our own words: shorten every word of a sentence to the dictionary stem it
// begins with, leaving words with no stem alone.
func TestReplaceWordsByRoot_problemSiteShape(t *testing.T) {
	roots := []string{"cook", "bake", "bak", "oven"}
	words := strings.Fields("bakers cooked bread in ovens daily")
	got := strings.Join(ReplaceWordsByRoot(nil, roots, words), " ")
	if want := "bak cook bread in oven daily"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestReplaceWordsByRoot_appendsAndKeepsInput(t *testing.T) {
	roots := []string{"re"}
	words := []string{"redo", "do"}
	before := slices.Clone(words)
	got := ReplaceWordsByRoot([]string{"x"}, roots, words)
	if !slices.Equal(got, []string{"x", "re", "do"}) || !slices.Equal(words, before) {
		t.Errorf("got %q (words now %q)", got, words)
	}
}

func TestReplaceWordsByRoot_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	randWord := func(maxLen int) string {
		b := make([]byte, rng.IntN(maxLen+1))
		for i := range b {
			b[i] = "pq"[rng.IntN(2)]
		}
		return string(b)
	}
	for trial := range 300 {
		roots := make([]string, rng.IntN(8))
		for i := range roots {
			roots[i] = randWord(3)
		}
		words := make([]string, rng.IntN(10))
		for i := range words {
			words[i] = randWord(6)
		}
		got := ReplaceWordsByRoot(nil, roots, words)
		for i, w := range words {
			if want := naiveShortestRoot(roots, w); got[i] != want {
				t.Fatalf("trial %d roots %q word %q: %q, want %q", trial, roots, w, got[i], want)
			}
		}
	}
}

func BenchmarkShortestPrefixOf(b *testing.B) {
	for _, size := range benchSizes {
		words := benchWords(size.n)
		tr := wordTrie(words[:size.n/10])
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				for _, w := range words {
					_, _, _ = tr.ShortestPrefixOf(w)
				}
			}
		})
	}
}

func BenchmarkReplaceWordsByRoot(b *testing.B) {
	for _, size := range benchSizes {
		words := benchWords(size.n)
		roots := make([]string, size.n/10)
		for i := range roots {
			roots[i] = words[i][:3]
		}
		dst := make([]string, 0, size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				dst = ReplaceWordsByRoot(dst[:0], roots, words)
			}
		})
	}
}

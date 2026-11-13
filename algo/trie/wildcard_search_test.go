package trie

import (
	"maps"
	"math/rand/v2"
	"slices"
	"testing"
)

// naiveWildcard reports whether key matches pattern byte by byte with '.' as any byte.
func naiveWildcard(key, pattern string) bool {
	if len(key) != len(pattern) {
		return false
	}
	for i := range len(key) {
		if pattern[i] != Wildcard && pattern[i] != key[i] {
			return false
		}
	}
	return true
}

// Knuth's common words (TAOCP Vol. 3, §6.3): "H.." matches HAD and HER and HIS; ".." matches the
// twelve two-letter words; "T.E" matches only THE. Derived by hand from the list.
func TestWildcardSearch_knuthCommonWords(t *testing.T) {
	tr := wordTrie(commonWords)
	tests := map[string][]string{
		"H..":   {"HAD", "HER", "HIS"},
		"..":    {"AS", "AT", "BE", "BY", "HE", "IN", "IS", "IT", "OF", "ON", "OR", "TO"},
		"T.E":   {"THE"},
		"W...H": {"WHICH"},
		"WHICH": {"WHICH"},
		"Q..":   nil,
		".....": {"WHICH"},
		"T.":    {"TO"},
	}
	for pattern, want := range tests {
		if got := collect(tr.WildcardSearch(pattern)); !slices.Equal(got, want) {
			t.Errorf("WildcardSearch(%q) = %v, want %v", pattern, got, want)
		}
	}
}

func TestWildcardSearch_cases(t *testing.T) {
	var tr Trie[int]
	tr.Insert("", 0)
	tr.Insert("a.c", 1)
	tr.Insert("abc", 2)
	tr.Insert("é", 3) // 2 bytes
	tr.Insert("ab", 4)
	tests := []struct {
		name, pattern string
		want          []string
	}{
		{"empty pattern matches the empty key", "", []string{""}},
		{"literal dot is matched by a dot", "a.c", []string{"a.c", "abc"}},
		{"exact length only", "a.", []string{"ab"}},
		{"all dots", "...", []string{"a.c", "abc"}},
		{"one dot does not match a two-byte character", ".", nil},
		{"two dots match a two-byte character", "..", []string{"ab", "é"}},
		{"half a character", "\xc3.", []string{"é"}},
		{"too long", "....", nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := collect(tr.WildcardSearch(tt.pattern)); !slices.Equal(got, tt.want) {
				t.Errorf("WildcardSearch(%q) = %q, want %q", tt.pattern, got, tt.want)
			}
		})
	}
}

// Problem-site shape, in our own words: a word dictionary is asked whether any stored word fits a
// pattern with unknown letters; stopping at the first hit answers it.
func TestWildcardSearch_problemSiteShape(t *testing.T) {
	var tr Trie[struct{}]
	for _, w := range []string{"red", "rod", "ran"} {
		tr.Insert(w, struct{}{})
	}
	exists := func(p string) bool {
		for range tr.WildcardSearch(p) {
			return true
		}
		return false
	}
	for p, want := range map[string]bool{"r.d": true, "..n": true, "re.": true, "r..d": false, "x..": false, "...": true} {
		if got := exists(p); got != want {
			t.Errorf("exists(%q) = %v, want %v", p, got, want)
		}
	}
	var nilTrie *Trie[int]
	if got := collect(nilTrie.WildcardSearch("...")); len(got) != 0 {
		t.Errorf("nil trie matched %v", got)
	}
}

func TestWildcardSearch_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 300 {
		var tr Trie[int]
		ref := map[string]int{}
		for i := range rng.IntN(50) {
			b := make([]byte, rng.IntN(5))
			for j := range b {
				b[j] = "ab."[rng.IntN(3)]
			}
			tr.Insert(string(b), i)
			ref[string(b)] = i
		}
		p := make([]byte, rng.IntN(5))
		for j := range p {
			p[j] = "ab."[rng.IntN(3)]
		}
		var want []string
		for _, key := range slices.Sorted(maps.Keys(ref)) {
			if naiveWildcard(key, string(p)) {
				want = append(want, key)
			}
		}
		if got := collect(tr.WildcardSearch(string(p))); !slices.Equal(got, want) {
			t.Fatalf("trial %d pattern %q: %v, want %v", trial, p, got, want)
		}
	}
}

func BenchmarkWildcardSearch(b *testing.B) {
	for _, size := range benchSizes {
		tr := wordTrie(benchWords(size.n))
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				for range tr.WildcardSearch("a.c.") {
				}
			}
		})
	}
}

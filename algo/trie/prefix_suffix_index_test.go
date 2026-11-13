package trie

import (
	"math/rand/v2"
	"slices"
	"strings"
	"testing"
)

func naivePrefixSuffix(words []string, prefix, suffix string) []int {
	var out []int
	for i, w := range words {
		if strings.HasPrefix(w, prefix) && strings.HasSuffix(w, suffix) {
			out = append(out, i)
		}
	}
	return out
}

// The permuterm example of Manning, Raghavan & Schütze, Introduction to Information Retrieval, §3.2.1:
// the query m*n finds man and men, and also moron (via the rotation n$moro). Here: prefix "m", suffix "n".
func TestPrefixSuffixIndex_permutermExample(t *testing.T) {
	words := []string{"man", "moron", "men", "hello", "mean"}
	x := NewPrefixSuffixIndex(words)
	if got := x.Lookup(nil, "m", "n"); !slices.Equal(got, []int{0, 1, 2, 4}) {
		t.Errorf(`Lookup("m", "n") = %v, want [0 1 2 4]`, got)
	}
	if got := x.Lookup(nil, "he", "llo"); !slices.Equal(got, []int{3}) {
		t.Errorf(`Lookup("he", "llo") = %v, want [3]`, got)
	}
	if i, ok := x.LookupLast("m", "n"); !ok || i != 4 {
		t.Errorf(`LookupLast("m", "n") = (%d, %v), want (4, true)`, i, ok)
	}
}

func TestPrefixSuffixIndex_cases(t *testing.T) {
	words := []string{"abc", "", "abc", "bc", "a", "abcabc"}
	x := NewPrefixSuffixIndex(words)
	tests := []struct {
		name, prefix, suffix string
		want                 []int
	}{
		{"empty matches all", "", "", []int{0, 1, 2, 3, 4, 5}},
		{"overlap inside a word", "ab", "bc", []int{0, 2, 5}},
		{"prefix and suffix are the whole word", "abc", "abc", []int{0, 2, 5}},
		{"prefix longer than some words", "abca", "", []int{5}},
		{"suffix only", "", "c", []int{0, 2, 3, 5}},
		{"prefix only", "a", "", []int{0, 2, 4, 5}},
		{"one-letter word as both", "a", "a", []int{4}},
		{"no match", "z", "", nil},
		{"too long together is fine when overlapping", "abcab", "cabc", []int{5}},
		{"suffix longer than every word", "", "xabcabc", nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := x.Lookup(nil, tt.prefix, tt.suffix); !slices.Equal(got, tt.want) {
				t.Errorf("Lookup(%q, %q) = %v, want %v", tt.prefix, tt.suffix, got, tt.want)
			}
			i, ok := x.LookupLast(tt.prefix, tt.suffix)
			if wantOK := len(tt.want) > 0; ok != wantOK || (ok && i != tt.want[len(tt.want)-1]) || (!ok && i != -1) {
				t.Errorf("LookupLast(%q, %q) = (%d, %v)", tt.prefix, tt.suffix, i, ok)
			}
		})
	}
}

// Problem-site shape, in our own words: a word list answers "the latest-added word that starts with p
// and ends with s", or nothing.
func TestPrefixSuffixIndex_problemSiteShape(t *testing.T) {
	x := NewPrefixSuffixIndex([]string{"lemon", "melon", "lemonade"})
	for _, q := range []struct {
		p, s string
		want int
		ok   bool
	}{{"le", "on", 0, true}, {"le", "e", 2, true}, {"m", "n", 1, true}, {"x", "", -1, false}} {
		if i, ok := x.LookupLast(q.p, q.s); i != q.want || ok != q.ok {
			t.Errorf("LookupLast(%q, %q) = (%d, %v), want (%d, %v)", q.p, q.s, i, ok, q.want, q.ok)
		}
	}
}

func TestPrefixSuffixIndex_emptyAndNil(t *testing.T) {
	for _, x := range []*PrefixSuffixIndex{nil, {}, NewPrefixSuffixIndex(nil)} {
		if got := x.Lookup(nil, "", ""); len(got) != 0 {
			t.Errorf("empty index Lookup = %v", got)
		}
		if i, ok := x.LookupLast("", ""); ok || i != -1 {
			t.Errorf("empty index LookupLast = (%d, %v)", i, ok)
		}
	}
	x := NewPrefixSuffixIndex([]string{"ab"})
	if got := x.Lookup([]int{9}, "a", "b"); !slices.Equal(got, []int{9, 0}) {
		t.Errorf("append = %v, want [9 0]", got)
	}
}

func TestPrefixSuffixIndex_bytes(t *testing.T) {
	// "é" is C3 A9: a suffix of one byte (A9) matches; the separator symbol is not a byte, so words
	// holding every byte value (including '$' and 0xFF) cannot confuse it.
	words := []string{"café", "a$b", "\xff\x00"}
	x := NewPrefixSuffixIndex(words)
	if got := x.Lookup(nil, "", "\xa9"); !slices.Equal(got, []int{0}) {
		t.Errorf("byte suffix: %v", got)
	}
	if got := x.Lookup(nil, "a", "b"); !slices.Equal(got, []int{1}) {
		t.Errorf("dollar word: %v", got)
	}
	if got := x.Lookup(nil, "\xff", "\x00"); !slices.Equal(got, []int{2}) {
		t.Errorf("extreme bytes: %v", got)
	}
}

func TestPrefixSuffixIndex_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	randWord := func(maxLen int) string {
		b := make([]byte, rng.IntN(maxLen+1))
		for i := range b {
			b[i] = "uv"[rng.IntN(2)]
		}
		return string(b)
	}
	for trial := range 300 {
		words := make([]string, rng.IntN(25))
		for i := range words {
			words[i] = randWord(6)
		}
		x := NewPrefixSuffixIndex(words)
		for range 10 {
			p, s := randWord(3), randWord(3)
			want := naivePrefixSuffix(words, p, s)
			if got := x.Lookup(nil, p, s); !slices.Equal(got, want) {
				t.Fatalf("trial %d Lookup(%q, %q) over %q = %v, want %v", trial, p, s, words, got, want)
			}
			i, ok := x.LookupLast(p, s)
			if len(want) == 0 && (ok || i != -1) || len(want) > 0 && (!ok || i != want[len(want)-1]) {
				t.Fatalf("trial %d LookupLast(%q, %q) = (%d, %v), want %v", trial, p, s, i, ok, want)
			}
		}
	}
}

func TestPrefixSuffixIndex_lookupLastNoAlloc(t *testing.T) {
	x := NewPrefixSuffixIndex([]string{"apple", "ample"})
	if allocs := testing.AllocsPerRun(100, func() { _, _ = x.LookupLast("a", "le") }); allocs != 0 {
		t.Errorf("allocated %v times, want 0", allocs)
	}
}

// Build cost is Σ L²; benchWords are 3–12 bytes.
func BenchmarkNewPrefixSuffixIndex(b *testing.B) {
	for _, size := range benchSizes {
		words := benchWords(size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_ = NewPrefixSuffixIndex(words)
			}
		})
	}
}

func BenchmarkPrefixSuffixIndexLookup(b *testing.B) {
	for _, size := range benchSizes {
		x := NewPrefixSuffixIndex(benchWords(size.n))
		dst := make([]int, 0, size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				dst = x.Lookup(dst[:0], "a", "e")
			}
		})
	}
}

func BenchmarkPrefixSuffixIndexLookupLast(b *testing.B) {
	for _, size := range benchSizes {
		x := NewPrefixSuffixIndex(benchWords(size.n))
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = x.LookupLast("ab", "yz")
			}
		})
	}
}

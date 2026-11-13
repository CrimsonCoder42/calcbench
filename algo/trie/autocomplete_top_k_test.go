package trie

import (
	"cmp"
	"errors"
	"math"
	"math/rand/v2"
	"slices"
	"strings"
	"testing"
)

// byWeightDesc ranks higher weights first.
func byWeightDesc(a, b int) int { return cmp.Compare(b, a) }

func keysOf[V any](entries []Entry[V]) []string {
	keys := make([]string, len(entries))
	for i, e := range entries {
		keys[i] = e.Key
	}
	return keys
}

// On Knuth's common words (TAOCP Vol. 3, §6.3), weighted by list position, the two heaviest "H" words
// are HIS (14) and HER (13); in byte order the first two are HAD and HAVE. Derived by hand.
func TestAutocompleteTopK_knuthCommonWords(t *testing.T) {
	tr := wordTrie(commonWords)
	got, err := tr.AutocompleteTopK(nil, "H", 2, byWeightDesc)
	if err != nil {
		t.Fatal(err)
	}
	if want := []Entry[int]{{"HIS", 14}, {"HER", 13}}; !slices.Equal(got, want) {
		t.Errorf("by weight = %v, want %v", got, want)
	}
	got, _ = tr.AutocompleteTopK(nil, "H", 2, nil)
	if want := []Entry[int]{{"HAD", 10}, {"HAVE", 11}}; !slices.Equal(got, want) {
		t.Errorf("by byte order = %v, want %v", got, want)
	}
}

func TestAutocompleteTopK_cases(t *testing.T) {
	var tr Trie[int]
	for k, v := range map[string]int{"apple": 5, "app": 9, "apply": 5, "apt": 1, "bat": 100, "ap": 5} {
		tr.Insert(k, v)
	}
	tests := []struct {
		name   string
		prefix string
		k      int
		rank   func(a, b int) int
		want   []string
	}{
		{"ties broken by byte order", "ap", 4, byWeightDesc, []string{"app", "ap", "apple", "apply"}},
		{"k larger than matches", "ap", 50, byWeightDesc, []string{"app", "ap", "apple", "apply", "apt"}},
		{"k zero", "ap", 0, byWeightDesc, nil},
		{"one", "ap", 1, byWeightDesc, []string{"app"}},
		{"no match", "zz", 3, byWeightDesc, nil},
		{"prefix is a key", "apple", 3, byWeightDesc, []string{"apple"}},
		{"empty prefix", "", 1, byWeightDesc, []string{"bat"}},
		{"ascending weight", "ap", 2, cmp.Compare[int], []string{"apt", "ap"}},
		{"nil rank byte order", "ap", 3, nil, []string{"ap", "app", "apple"}},
		{"huge k", "ap", math.MaxInt, nil, []string{"ap", "app", "apple", "apply", "apt"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tr.AutocompleteTopK(nil, tt.prefix, tt.k, tt.rank)
			if err != nil {
				t.Fatal(err)
			}
			if keys := keysOf(got); !slices.Equal(keys, tt.want) && !(len(keys) == 0 && len(tt.want) == 0) {
				t.Errorf("got %v, want %v", keys, tt.want)
			}
		})
	}
}

// Problem-site shape, in our own words: a search box suggests at most three stored words that begin
// with what has been typed so far, alphabetically, after each keystroke.
func TestAutocompleteTopK_problemSiteShape(t *testing.T) {
	var tr Trie[struct{}]
	for _, w := range []string{"monitor", "mobile", "mouse", "moneypot", "mousepad", "money"} {
		tr.Insert(w, struct{}{})
	}
	want := [][]string{
		{"mobile", "money", "moneypot"},
		{"mobile", "money", "moneypot"},
		{"mouse", "mousepad"},
		{"mouse", "mousepad"},
		{"mouse", "mousepad"},
	}
	for i := range want {
		got, _ := tr.AutocompleteTopK(nil, "mouse"[:i+1], 3, nil)
		if keys := keysOf(got); !slices.Equal(keys, want[i]) {
			t.Errorf("typed %q: %v, want %v", "mouse"[:i+1], keys, want[i])
		}
	}
}

func TestAutocompleteTopK_appendsAndErrors(t *testing.T) {
	tr := wordTrie(commonWords)
	dst := []Entry[int]{{"keep", -1}}
	got, err := tr.AutocompleteTopK(dst, "B", 2, byWeightDesc)
	if err != nil || !slices.Equal(got, []Entry[int]{{"keep", -1}, {"BY", 7}, {"BUT", 6}}) {
		t.Errorf("append = %v, %v", got, err)
	}
	got, err = tr.AutocompleteTopK(dst, "B", -1, byWeightDesc)
	if !errors.Is(err, ErrNegativeK) || len(got) != 1 {
		t.Errorf("k=-1: (%v, %v), want ErrNegativeK and dst unchanged", got, err)
	}
	var nilTrie *Trie[int]
	if got, err := nilTrie.AutocompleteTopK(nil, "", 3, nil); err != nil || len(got) != 0 {
		t.Errorf("nil trie: (%v, %v)", got, err)
	}
}

// Oracle: sort the matching keys by (rank, key) and take k.
func TestAutocompleteTopK_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 300 {
		var tr Trie[int]
		ref := map[string]int{}
		for range rng.IntN(60) {
			b := make([]byte, rng.IntN(5))
			for i := range b {
				b[i] = "xyz"[rng.IntN(3)]
			}
			w := rng.IntN(5) // heavy ties
			tr.Insert(string(b), w)
			ref[string(b)] = w
		}
		prefix := "xyz"[:rng.IntN(2)]
		k := rng.IntN(8)
		var want []Entry[int]
		for key, w := range ref {
			if strings.HasPrefix(key, prefix) {
				want = append(want, Entry[int]{key, w})
			}
		}
		slices.SortFunc(want, func(a, b Entry[int]) int {
			return cmp.Or(byWeightDesc(a.Value, b.Value), strings.Compare(a.Key, b.Key))
		})
		want = want[:min(k, len(want))]
		got, _ := tr.AutocompleteTopK(nil, prefix, k, byWeightDesc)
		if !slices.Equal(got, want) && !(len(got) == 0 && len(want) == 0) {
			t.Fatalf("trial %d prefix %q k %d: %v, want %v", trial, prefix, k, got, want)
		}
	}
}

func BenchmarkAutocompleteTopK(b *testing.B) {
	for _, size := range benchSizes {
		words := benchWords(size.n)
		tr := wordTrie(words)
		dst := make([]Entry[int], 0, 10)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				dst, _ = tr.AutocompleteTopK(dst[:0], "a", 10, byWeightDesc)
			}
		})
	}
}

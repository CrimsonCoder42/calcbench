package graph

import (
	"math/rand/v2"
	"slices"
	"strings"
	"testing"
)

// oneLetterApart reports whether a and b have the same length and differ in exactly one byte.
func oneLetterApart(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	diff := 0
	for i := range len(a) {
		if a[i] != b[i] {
			diff++
		}
	}
	return diff == 1
}

// Hand example (catalog/graph.md, ALG-GRAPH-021), in the sense of Knuth's LADDERS (The Stanford GraphBase,
// 1993): cold and warm differ in all four letters, so a ladder has at least 5 words; BFS with this
// dictionary order finds cold cord word ward warm.
func TestWordLadder_handExample(t *testing.T) {
	dict := []string{"cord", "card", "ward", "warm", "word", "worm", "wore", "core"}
	got := WordLadder("cold", "warm", dict)
	if want := []string{"cold", "cord", "word", "ward", "warm"}; !slices.Equal(got, want) {
		t.Errorf("WordLadder = %v, want %v", got, want)
	}
}

func TestWordLadder_cases(t *testing.T) {
	tests := []struct {
		name          string
		start, target string
		dict          []string
		want          []string
	}{
		{"start equals target", "same", "same", nil, []string{"same"}},
		{"empty words equal", "", "", nil, []string{""}},
		{"one step", "hot", "hat", []string{"hat"}, []string{"hot", "hat"}},
		{"target missing from dictionary", "hot", "hat", []string{"hit"}, nil},
		{"length mismatch", "hot", "heat", []string{"heat"}, nil},
		{"no connection", "aaa", "bbb", []string{"bbb", "abc"}, nil},
		{"empty dictionary", "ab", "cd", nil, nil},
		{"start in dictionary, duplicates, other lengths", "lead", "gold",
			[]string{"lead", "load", "load", "goad", "gold", "go", "golden"}, []string{"lead", "load", "goad", "gold"}},
		{"bytes: é and è share their first byte, so they are one byte apart", "é", "è", []string{"è"}, []string{"é", "è"}},
		{"bytes: é and ж differ in both bytes, so not one step", "é", "ж", []string{"ж"}, nil},
	}
	for _, tt := range tests {
		got := WordLadder(tt.start, tt.target, tt.dict)
		if !slices.Equal(got, tt.want) || (got == nil) != (tt.want == nil) {
			t.Errorf("%s: %q, want %q", tt.name, got, tt.want)
		}
	}
}

func TestWordLadder_inputNotModified(t *testing.T) {
	dict := []string{"cat", "cot", "cog", "dog"}
	before := slices.Clone(dict)
	_ = WordLadder("bat", "dog", dict)
	if !slices.Equal(dict, before) {
		t.Error("dictionary modified")
	}
}

// Oracle: a plain BFS comparing every pair of words; ladders must have the same length and be valid.
func TestWordLadder_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	word := func() string {
		var b strings.Builder
		for range 3 {
			b.WriteByte("abc"[rng.IntN(3)])
		}
		return b.String()
	}
	for run := 0; run < 300; run++ {
		start, target := word(), word()
		dict := make([]string, rng.IntN(12))
		for i := range dict {
			dict[i] = word()
		}
		if rng.IntN(2) == 0 {
			dict = append(dict, target)
		}
		// Slow oracle: words = start + dictionary; BFS with an O(N) scan per vertex.
		words := append([]string{start}, dict...)
		dist := make([]int, len(words))
		for i := range dist {
			dist[i] = -1
		}
		dist[0] = 0
		for queue := []int{0}; len(queue) > 0; queue = queue[1:] {
			u := queue[0]
			for v := 1; v < len(words); v++ {
				if dist[v] == -1 && oneLetterApart(words[u], words[v]) {
					dist[v] = dist[u] + 1
					queue = append(queue, v)
				}
			}
		}
		want := -1
		if start == target {
			want = 0
		}
		for v := 1; v < len(words); v++ {
			if words[v] == target && dist[v] >= 0 && (want == -1 || dist[v] < want) {
				want = dist[v]
			}
		}
		got := WordLadder(start, target, dict)
		if (want == -1) != (got == nil) || (got != nil && len(got) != want+1) {
			t.Fatalf("run %d: %q→%q in %q: got %q, want %d steps", run, start, target, dict, got, want)
		}
		for i := 1; i < len(got); i++ {
			if !oneLetterApart(got[i-1], got[i]) || !slices.Contains(dict, got[i]) {
				t.Fatalf("run %d: invalid step %q→%q", run, got[i-1], got[i])
			}
		}
	}
}

func BenchmarkWordLadder(b *testing.B) {
	// N random 5-letter words over a 6-letter alphabet (dense enough to connect), N = 1e3, 1e4, 1e5.
	for _, size := range benchSizes {
		rng := rand.New(rand.NewPCG(1, 2))
		dict := make([]string, size.n)
		for i := range dict {
			buf := make([]byte, 5)
			for k := range buf {
				buf[k] = "abcdef"[rng.IntN(6)]
			}
			dict[i] = string(buf)
		}
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_ = WordLadder("aaaaa", dict[len(dict)-1], dict)
			}
		})
	}
}

package text

import (
	"math/rand/v2"
	"testing"
)

// naiveMinimumWindow is the oracle: try every window, shortest first, leftmost first.
func naiveMinimumWindow(text, target string) (int, int, bool) {
	if target == "" {
		return 0, 0, true
	}
	var need [256]int
	for i := range len(target) {
		need[target[i]]++
	}
	for l := len(target); l <= len(text); l++ {
		for s := 0; s+l <= len(text); s++ {
			var have [256]int
			for i := s; i < s+l; i++ {
				have[text[i]]++
			}
			ok := true
			for c := range need {
				if have[c] < need[c] {
					ok = false
					break
				}
			}
			if ok {
				return s, l, true
			}
		}
	}
	return 0, 0, false
}

// Hand-derived: in xaybzcab the shortest window holding a, b and c is "cab" at 5.
func TestMinimumWindowSubstring_example(t *testing.T) {
	if s, l, ok := MinimumWindowSubstring("xaybzcab", "abc"); s != 5 || l != 3 || !ok {
		t.Errorf("MinimumWindowSubstring(xaybzcab, abc) = (%d, %d, %v), want (5, 3, true)", s, l, ok)
	}
}

// Our own wording of the usual exercise cases.
func TestMinimumWindowSubstring_cases(t *testing.T) {
	tests := []struct {
		name, text, target string
		start, length      int
		found              bool
	}{
		{"empty target", "abc", "", 0, 0, true},
		{"empty text", "", "a", 0, 0, false},
		{"absent byte", "abc", "d", 0, 0, false},
		{"target longer than text", "a", "aa", 0, 0, false},
		{"multiplicity counts", "abab", "aa", 0, 3, true},
		{"whole text", "abc", "cba", 0, 3, true},
		{"single byte", "xyz", "y", 1, 1, true},
		{"tie goes left", "ab_ba", "ab", 0, 2, true},
		{"case matters", "aA", "A", 1, 1, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, l, ok := MinimumWindowSubstring(tt.text, tt.target)
			if s != tt.start || l != tt.length || ok != tt.found {
				t.Errorf("MinimumWindowSubstring(%q, %q) = (%d, %d, %v), want (%d, %d, %v)", tt.text, tt.target, s, l, ok, tt.start, tt.length, tt.found)
			}
		})
	}
}

func TestMinimumWindowSubstring_noAllocation(t *testing.T) {
	if allocs := testing.AllocsPerRun(100, func() { _, benchSink, _ = MinimumWindowSubstring("xaybzcab", "abc") }); allocs != 0 {
		t.Errorf("MinimumWindowSubstring allocates %v times, want 0", allocs)
	}
}

func TestMinimumWindowSubstring_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 500 {
		text := randomText(rng, rng.IntN(30), "abc")
		target := randomText(rng, rng.IntN(5), "abc")
		gs, gl, gok := MinimumWindowSubstring(text, target)
		ws, wl, wok := naiveMinimumWindow(text, target)
		if gs != ws || gl != wl || gok != wok {
			t.Fatalf("trial %d: MinimumWindowSubstring(%q, %q) = (%d, %d, %v), want (%d, %d, %v)", trial, text, target, gs, gl, gok, ws, wl, wok)
		}
	}
}

func BenchmarkMinimumWindowSubstring(b *testing.B) {
	rng := rand.New(rand.NewPCG(1, 2))
	for _, size := range benchSizes {
		text := randomText(rng, size.n, "abcdefghij")
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, benchSink, _ = MinimumWindowSubstring(text, "jihgfedcba")
			}
		})
	}
}

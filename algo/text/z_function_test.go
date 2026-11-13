package text

import (
	"math/rand/v2"
	"slices"
	"testing"
)

// naiveZ is the oracle: compare the string with each suffix byte by byte.
func naiveZ(s string) []int {
	out := make([]int, len(s))
	for i := range len(s) {
		k := 0
		for i+k < len(s) && s[k] == s[i+k] {
			k++
		}
		out[i] = k
	}
	return out
}

// Gusfield, Algorithms on Strings, Trees, and Sequences, 1997, §1.3: for S = aabcaabxaaz,
// Z₅ = 3, Z₆ = 1, Z₇ = Z₈ = 0, Z₉ = 2 (1-based); the remaining values are hand-derived.
func TestZFunction_gusfieldExample(t *testing.T) {
	got := ZFunction(nil, "aabcaabxaaz")
	if want := []int{11, 1, 0, 0, 3, 1, 0, 0, 2, 1, 0}; !slices.Equal(got, want) {
		t.Errorf("ZFunction(aabcaabxaaz) = %v, want %v", got, want)
	}
}

func TestZFunction_cases(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want []int
	}{
		{"empty", "", []int{}},
		{"one", "q", []int{1}},
		{"two equal", "qq", []int{2, 1}},
		{"two different", "qr", []int{2, 0}},
		{"all equal", "aaaaa", []int{5, 4, 3, 2, 1}},
		{"periodic", "abab", []int{4, 0, 2, 0}},
		{"pattern$text search", "ab$xaby", []int{7, 0, 0, 0, 2, 0, 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ZFunction(nil, tt.s); !slices.Equal(got, tt.want) {
				t.Errorf("ZFunction(%q) = %v, want %v", tt.s, got, tt.want)
			}
		})
	}
}

func TestZFunction_bytesAndAppend(t *testing.T) {
	if got, want := ZFunction([]int{5}, []byte("aa")), []int{5, 2, 1}; !slices.Equal(got, want) {
		t.Errorf("ZFunction([5], aa) = %v, want %v", got, want)
	}
}

func TestZFunction_noAllocation(t *testing.T) {
	buf := make([]int, 0, 16)
	if allocs := testing.AllocsPerRun(100, func() { buf = ZFunction(buf[:0], "aabcaabxaaz") }); allocs != 0 {
		t.Errorf("ZFunction with a large enough dst allocates %v times, want 0", allocs)
	}
}

func TestZFunction_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 500 {
		s := randomText(rng, rng.IntN(50), "abc"[:1+trial%3])
		if got, want := ZFunction(nil, s), naiveZ(s); !slices.Equal(got, want) {
			t.Fatalf("trial %d: ZFunction(%q) = %v, want %v", trial, s, got, want)
		}
	}
}

func BenchmarkZFunction(b *testing.B) {
	rng := rand.New(rand.NewPCG(1, 2))
	for _, size := range benchSizes {
		s := randomText(rng, size.n, "ab")
		buf := make([]int, 0, size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				buf = ZFunction(buf[:0], s)
			}
			benchSink = len(buf)
		})
	}
}

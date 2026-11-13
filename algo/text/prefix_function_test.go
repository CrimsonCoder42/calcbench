package text

import (
	"math/rand/v2"
	"slices"
	"strings"
	"testing"
)

// benchSizes are the collection sizes of .claude/rules/go-complexity.md. Shared by the benchmarks.
var benchSizes = []struct {
	name string
	n    int
}{{"1e3", 1_000}, {"1e4", 10_000}, {"1e5", 100_000}}

// benchSink keeps benchmark results observable. Shared by the benchmarks.
var benchSink int

// randomText returns n bytes drawn uniformly from alphabet. Shared by the tests and benchmarks.
func randomText(rng *rand.Rand, n int, alphabet string) string {
	var b strings.Builder
	b.Grow(n)
	for range n {
		b.WriteByte(alphabet[rng.IntN(len(alphabet))])
	}
	return b.String()
}

// naiveSearch is the slow obvious oracle: every start where pattern occurs, by direct comparison.
func naiveSearch(text, pattern string) []int {
	var out []int
	for i := 0; i+len(pattern) <= len(text); i++ {
		if text[i:i+len(pattern)] == pattern {
			out = append(out, i)
		}
	}
	return out
}

// naivePrefixFunction is the oracle: for each prefix, the longest proper border by trying every length.
func naivePrefixFunction(s string) []int {
	out := make([]int, len(s))
	for q := range len(s) {
		for k := q; k > 0; k-- {
			if s[:k] == s[q+1-k:q+1] {
				out[q] = k
				break
			}
		}
	}
	return out
}

// Cormen et al., Introduction to Algorithms, 3rd ed., §32.4, Figure 32.10: for P = ababaca the prefix
// function is 0 0 1 2 3 0 1.
func TestPrefixFunction_clrsFigure32_10(t *testing.T) {
	if got, want := PrefixFunction(nil, "ababaca"), []int{0, 0, 1, 2, 3, 0, 1}; !slices.Equal(got, want) {
		t.Errorf("PrefixFunction(ababaca) = %v, want %v", got, want)
	}
}

func TestPrefixFunction_cases(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want []int
	}{
		{"empty", "", []int{}},
		{"one byte", "a", []int{0}},
		{"two equal", "aa", []int{0, 1}},
		{"two different", "ab", []int{0, 0}},
		{"all equal", "aaaa", []int{0, 1, 2, 3}},
		{"no border", "abcd", []int{0, 0, 0, 0}},
		{"KMP 1977 pattern", "abcabcacab", []int{0, 0, 0, 1, 2, 3, 4, 0, 1, 2}},
		{"period 3", "xyzxyzxy", []int{0, 0, 0, 1, 2, 3, 4, 5}},
		{"high bytes", "\xff\x00\xff", []int{0, 0, 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PrefixFunction(nil, tt.s); !slices.Equal(got, tt.want) {
				t.Errorf("PrefixFunction(%q) = %v, want %v", tt.s, got, tt.want)
			}
		})
	}
}

func TestPrefixFunction_bytesAndAppend(t *testing.T) {
	dst := []int{-7}
	got := PrefixFunction(dst, []byte("abab"))
	if want := []int{-7, 0, 0, 1, 2}; !slices.Equal(got, want) {
		t.Errorf("PrefixFunction(dst=[-7], abab) = %v, want %v", got, want)
	}
}

func TestPrefixFunction_noAllocation(t *testing.T) {
	buf := make([]int, 0, 16)
	if allocs := testing.AllocsPerRun(100, func() { buf = PrefixFunction(buf[:0], "abcabcacab") }); allocs != 0 {
		t.Errorf("PrefixFunction with a large enough dst allocates %v times, want 0", allocs)
	}
}

func TestPrefixFunction_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 500 {
		s := randomText(rng, rng.IntN(40), "ab"[:1+trial%2])
		if got, want := PrefixFunction(nil, s), naivePrefixFunction(s); !slices.Equal(got, want) {
			t.Fatalf("trial %d: PrefixFunction(%q) = %v, want %v", trial, s, got, want)
		}
	}
}

func BenchmarkPrefixFunction(b *testing.B) {
	rng := rand.New(rand.NewPCG(1, 2))
	for _, size := range benchSizes {
		s := randomText(rng, size.n, "ab")
		buf := make([]int, 0, size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				buf = PrefixFunction(buf[:0], s)
			}
			benchSink = len(buf)
		})
	}
}

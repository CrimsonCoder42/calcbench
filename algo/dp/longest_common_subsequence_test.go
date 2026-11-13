package dp

import (
	"math"
	"math/rand/v2"
	"testing"
)

// isSubsequence reports whether sub occurs in seq as a subsequence (greedy match).
func isSubsequence[T comparable](sub, seq []T) bool {
	k := 0
	for _, x := range seq {
		if k < len(sub) && sub[k] == x {
			k++
		}
	}
	return k == len(sub)
}

// lcsLengthSlow fills the full (m+1)×(n+1) table of CLRS §15.4.
func lcsLengthSlow(a, b []byte) int {
	c := make([][]int, len(a)+1)
	for i := range c {
		c[i] = make([]int, len(b)+1)
	}
	for i := 1; i <= len(a); i++ {
		for j := 1; j <= len(b); j++ {
			if a[i-1] == b[j-1] {
				c[i][j] = c[i-1][j-1] + 1
			} else {
				c[i][j] = max(c[i-1][j], c[i][j-1])
			}
		}
	}
	return c[len(a)][len(b)]
}

func randomBytes(rng *rand.Rand, n int, alphabet string) []byte {
	out := make([]byte, n)
	for i := range out {
		out[i] = alphabet[rng.IntN(len(alphabet))]
	}
	return out
}

// CLRS 3rd ed., §15.4, Figure 15.8: X = ABCBDAB, Y = BDCABA have an LCS of length 4 (BCBA is one).
func TestLCSLength_clrsExample(t *testing.T) {
	if got := LCSLength([]byte("ABCBDAB"), []byte("BDCABA")); got != 4 {
		t.Errorf("LCSLength(ABCBDAB, BDCABA) = %d, want 4", got)
	}
}

func TestLongestCommonSubsequence_clrsExample(t *testing.T) {
	x, y := []byte("ABCBDAB"), []byte("BDCABA")
	got := LongestCommonSubsequence(x, y)
	if len(got) != 4 || !isSubsequence(got, x) || !isSubsequence(got, y) {
		t.Errorf("LongestCommonSubsequence(ABCBDAB, BDCABA) = %q, want a common subsequence of length 4", got)
	}
}

func TestLCSLength_cases(t *testing.T) {
	tests := []struct {
		name string
		a, b string
		want int
	}{
		{"both empty", "", "", 0},
		{"one empty", "abc", "", 0},
		{"one element equal", "a", "a", 1},
		{"one element different", "a", "b", 0},
		{"identical", "abcde", "abcde", 5},
		{"reversed", "abcd", "dcba", 1},
		{"duplicates", "aaaa", "aa", 2},
		{"disjoint alphabets", "abc", "xyz", 0},
		{"interleaved", "axbycz", "abc", 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, b := []byte(tt.a), []byte(tt.b)
			if got := LCSLength(a, b); got != tt.want {
				t.Errorf("LCSLength(%q, %q) = %d, want %d", tt.a, tt.b, got, tt.want)
			}
			if got := LCSLength(b, a); got != tt.want {
				t.Errorf("LCSLength(%q, %q) = %d, want %d (symmetric)", tt.b, tt.a, got, tt.want)
			}
			got := LongestCommonSubsequence(a, b)
			if got == nil || len(got) != tt.want || !isSubsequence(got, a) || !isSubsequence(got, b) {
				t.Errorf("LongestCommonSubsequence(%q, %q) = %#v, want a common subsequence of length %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestLCSLength_nanMatchesNothing(t *testing.T) {
	a := []float64{math.NaN(), 1, 2}
	if got := LCSLength(a, a); got != 2 {
		t.Errorf("LCSLength with NaN = %d, want 2 (NaN != NaN)", got)
	}
	if got := LongestCommonSubsequence(a, a); len(got) != 2 || got[0] != 1 || got[1] != 2 {
		t.Errorf("LongestCommonSubsequence with NaN = %v, want [1 2]", got)
	}
}

func TestLCSLength_allocations(t *testing.T) {
	a, b := []byte("ABCBDAB"), []byte("BDCABA")
	if n := testing.AllocsPerRun(100, func() { LCSLength(a, b) }); n != 1 {
		t.Errorf("LCSLength allocates %v times, want 1 (one row)", n)
	}
	if n := testing.AllocsPerRun(100, func() { LCSLength(a, nil) }); n != 0 {
		t.Errorf("LCSLength(a, nil) allocates %v times, want 0", n)
	}
}

func TestLongestCommonSubsequence_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 400 {
		a := randomBytes(rng, rng.IntN(40), "abc")
		b := randomBytes(rng, rng.IntN(40), "abcd")
		want := lcsLengthSlow(a, b)
		if got := LCSLength(a, b); got != want {
			t.Fatalf("trial %d: LCSLength(%q, %q) = %d, want %d", trial, a, b, got, want)
		}
		got := LongestCommonSubsequence(a, b)
		if len(got) != want || !isSubsequence(got, a) || !isSubsequence(got, b) {
			t.Fatalf("trial %d: LongestCommonSubsequence(%q, %q) = %q, want length %d", trial, a, b, got, want)
		}
	}
}

func BenchmarkLCSLength(b *testing.B) {
	for _, size := range quadraticSizes {
		rng := rand.New(rand.NewPCG(1, 2))
		x, y := randomBytes(rng, size.n, "ACGT"), randomBytes(rng, size.n, "ACGT")
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_ = LCSLength(x, y)
			}
		})
	}
}

func BenchmarkLongestCommonSubsequence(b *testing.B) {
	for _, size := range quadraticSizes {
		rng := rand.New(rand.NewPCG(1, 2))
		x, y := randomBytes(rng, size.n, "ACGT"), randomBytes(rng, size.n, "ACGT")
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_ = LongestCommonSubsequence(x, y)
			}
		})
	}
}

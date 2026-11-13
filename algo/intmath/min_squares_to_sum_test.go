package intmath

import (
	"math"
	"math/rand/v2"
	"testing"
)

// Hardy & Wright §20.5 and §20.10: 7 = 4 + 1 + 1 + 1 needs four squares; 4^a(8b + 7) always does.
func TestMinSquaresToSum_sourceExample(t *testing.T) {
	for _, n := range []uint64{7, 15, 23, 28, 60, 112} {
		if got := MinSquaresToSum(n); got != 4 {
			t.Errorf("MinSquaresToSum(%d) = %d, want 4", n, got)
		}
	}
}

func TestMinSquaresToSum_cases(t *testing.T) {
	tests := []struct {
		n    uint64
		want int
	}{
		{0, 0}, {1, 1}, {2, 2}, {3, 3}, {4, 1}, {5, 2}, {6, 3}, {12, 3}, {13, 2}, {25, 1}, {50, 2},
		{(1<<32 - 1) * (1<<32 - 1), 1},
		{math.MaxUint64, 4}, // 2^64 − 1 ≡ 7 (mod 8)
		{1 << 63, 2},        // 2^63 = 2·4^31 = 4^31·(1 + 1)
		{1 << 62, 1},
		{3 << 62, 3}, // 3·4^31: m = 3 is neither a square, 8b + 7, nor a sum of two squares
	}
	for _, tt := range tests {
		if got := MinSquaresToSum(tt.n); got != tt.want {
			t.Errorf("MinSquaresToSum(%d) = %d, want %d", tt.n, got, tt.want)
		}
	}
}

// Oracle: the dynamic programme best[n] = 1 + min over squares s ≤ n of best[n − s].
func TestMinSquaresToSum_oracle(t *testing.T) {
	const limit = 20_000
	best := make([]int, limit+1)
	for n := 1; n <= limit; n++ {
		best[n] = math.MaxInt
		for a := 1; a*a <= n; a++ {
			best[n] = min(best[n], best[n-a*a]+1)
		}
		if got := MinSquaresToSum(uint64(n)); got != best[n] {
			t.Fatalf("MinSquaresToSum(%d) = %d, want %d", n, got, best[n])
		}
	}
	// Large n: a "2" answer must come with a witness; build sums of two squares and check they are found.
	rng := rand.New(rand.NewPCG(1, 2))
	for range 200 {
		a, b := rng.Uint64N(1<<20), rng.Uint64N(1<<20) // n < 2⁴¹: each call ≤ ~1.5e6 steps (H2: keep the test fast)
		n := a*a + b*b
		want := 2
		if a == 0 || b == 0 {
			want = 1
		}
		if s := ISqrt(n); s*s == n {
			want = 1
		}
		if got := MinSquaresToSum(n); got != want {
			t.Fatalf("MinSquaresToSum(%d² + %d²) = %d, want %d", a, b, got, want)
		}
	}
}

func BenchmarkMinSquaresToSum(b *testing.B) {
	// Worst case: n ≡ 1 (mod 4) and not a sum of two squares (n = 3·p, p ≡ 3 (mod 4), 3 ∤ p: the prime 3
	// appears to an odd power), so the two-pointer scan runs to the end: about √n steps.
	for _, c := range []struct {
		name string
		n    uint64
	}{{"1e6", 1_000_005}, {"1e10", 10_000_000_005}, {"1e14", 100_000_000_000_005}} {
		b.Run(c.name, func(b *testing.B) {
			b.ReportAllocs()
			s := 0
			for b.Loop() {
				s += MinSquaresToSum(c.n)
			}
			benchSink = uint64(s)
		})
	}
}

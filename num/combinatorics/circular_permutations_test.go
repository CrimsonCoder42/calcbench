package combinatorics

import (
	"errors"
	"math"
	"math/big"
	"math/rand/v2"
	"testing"
)

// OEIS A111492, row n = 5, k = 1…5: 5, 10, 20, 30, 24 (a(n,k) = (k−1)!·C(n,k)); OEIS A144151, row n = 5,
// k = 1…5: 5, 10, 10, 15, 12 (C(n,k) for k ≤ 2, else C(n,k)·(k−1)!/2).
func TestCircularPermutations_sourceExample(t *testing.T) {
	directed := []uint64{5, 10, 20, 30, 24}
	undirected := []uint64{5, 10, 10, 15, 12}
	for k := 1; k <= 5; k++ {
		if got, err := CircularPermutations(5, k); err != nil || got != directed[k-1] {
			t.Errorf("CircularPermutations(5, %d) = (%d, %v), want (%d, nil)", k, got, err, directed[k-1])
		}
		if got, err := CircularPermutationsUndirected(5, k); err != nil || got != undirected[k-1] {
			t.Errorf("CircularPermutationsUndirected(5, %d) = (%d, %v), want (%d, nil)", k, got, err, undirected[k-1])
		}
	}
}

func TestCircularPermutations_cases(t *testing.T) {
	tests := []struct {
		name               string
		n, k               int
		directed, undirect uint64
	}{
		{"four of six", 6, 4, 90, 45},
		{"all six: 5! and 5!/2", 6, 6, 120, 60},
		{"one of one", 1, 1, 1, 1},
		{"two of two", 2, 2, 1, 1},
		{"all 21: 20!", 21, 21, 2432902008176640000, 1216451004088320000},
		{"max int, k = 1", math.MaxInt, 1, math.MaxInt, math.MaxInt},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := CircularPermutations(tt.n, tt.k); err != nil || got != tt.directed {
				t.Errorf("CircularPermutations(%d, %d) = (%d, %v), want (%d, nil)", tt.n, tt.k, got, err, tt.directed)
			}
			if got, err := CircularPermutationsUndirected(tt.n, tt.k); err != nil || got != tt.undirect {
				t.Errorf("CircularPermutationsUndirected(%d, %d) = (%d, %v), want (%d, nil)", tt.n, tt.k, got, err, tt.undirect)
			}
		})
	}
	// C(22,20)·19! exceeds 2⁶⁴ − 1, its half does not: the undirected count must not overflow falsely.
	if _, err := CircularPermutations(22, 20); !errors.Is(err, ErrOverflow) {
		t.Errorf("CircularPermutations(22, 20) err = %v, want ErrOverflow", err)
	}
	if got, err := CircularPermutationsUndirected(22, 20); err != nil || got != 14050009097220096000 {
		t.Errorf("CircularPermutationsUndirected(22, 20) = (%d, %v), want (14050009097220096000, nil)", got, err)
	}
	for _, c := range [][2]int{{22, 22}, {math.MaxInt, 2}, {100, 50}} {
		if _, err := CircularPermutations(c[0], c[1]); !errors.Is(err, ErrOverflow) {
			t.Errorf("CircularPermutations(%d, %d) err = %v, want ErrOverflow", c[0], c[1], err)
		}
		if _, err := CircularPermutationsUndirected(c[0], c[1]); !errors.Is(err, ErrOverflow) {
			t.Errorf("CircularPermutationsUndirected(%d, %d) err = %v, want ErrOverflow", c[0], c[1], err)
		}
	}
}

func TestCircularPermutations_errors(t *testing.T) {
	for _, fn := range []func(int, int) (uint64, error){CircularPermutations, CircularPermutationsUndirected} {
		for _, c := range [][2]int{{-1, 1}, {3, -1}, {math.MinInt, 0}} {
			if got, err := fn(c[0], c[1]); !errors.Is(err, ErrNegative) || got != 0 {
				t.Errorf("(%d, %d) = (%d, %v), want (0, ErrNegative)", c[0], c[1], got, err)
			}
		}
		for _, c := range [][2]int{{5, 0}, {0, 0}, {5, 6}, {0, 1}} {
			if got, err := fn(c[0], c[1]); !errors.Is(err, ErrOutOfRange) || got != 0 {
				t.Errorf("(%d, %d) = (%d, %v), want (0, ErrOutOfRange)", c[0], c[1], got, err)
			}
		}
	}
}

// Against P(n,k)/k and, for k ≥ 3, P(n,k)/(2k), in math/big.
func TestCircularPermutations_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for range 500 {
		n := 1 + rng.IntN(40)
		k := 1 + rng.IntN(n)
		p := new(big.Int).MulRange(int64(n-k+1), int64(n))
		directed := new(big.Int).Quo(p, big.NewInt(int64(k)))
		got, err := CircularPermutations(n, k)
		wantU64(t, "CircularPermutations", got, err, directed)
		undirected := directed
		if k >= 3 {
			undirected = new(big.Int).Quo(directed, big.NewInt(2))
		}
		got, err = CircularPermutationsUndirected(n, k)
		wantU64(t, "CircularPermutationsUndirected", got, err, undirected)
	}
}

func TestCircularPermutations_noAlloc(t *testing.T) {
	if a := testing.AllocsPerRun(100, func() { sinkU, _ = CircularPermutations(20, 10) }); a != 0 {
		t.Errorf("CircularPermutations allocates %v times, want 0", a)
	}
	if a := testing.AllocsPerRun(100, func() { sinkU, _ = CircularPermutationsUndirected(20, 10) }); a != 0 {
		t.Errorf("CircularPermutationsUndirected allocates %v times, want 0", a)
	}
}

func BenchmarkCircularPermutations(b *testing.B) {
	b.Run("n=20,k=10", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			sinkU, _ = CircularPermutations(20, 10)
		}
	})
}

func BenchmarkCircularPermutationsUndirected(b *testing.B) {
	b.Run("n=20,k=10", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			sinkU, _ = CircularPermutationsUndirected(20, 10)
		}
	})
}

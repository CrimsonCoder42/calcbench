package combinatorics

import (
	"errors"
	"math"
	"math/big"
	"math/rand/v2"
	"testing"
)

// DLMF Table 26.3.1: C(6,3) = 20, C(8,2) = 28, C(10,5) = 252.
func TestBinomial_sourceExample(t *testing.T) {
	for _, c := range []struct {
		n, k int
		want uint64
	}{{6, 3, 20}, {8, 2, 28}, {10, 5, 252}} {
		if got, err := Binomial(c.n, c.k); err != nil || got != c.want {
			t.Errorf("Binomial(%d, %d) = (%d, %v), want (%d, nil)", c.n, c.k, got, err, c.want)
		}
	}
}

func TestBinomial_cases(t *testing.T) {
	tests := []struct {
		name string
		n, k int
		want uint64
	}{
		{"zero of zero", 0, 0, 1},
		{"k zero", 5, 0, 1},
		{"k equals n", 5, 5, 1},
		{"k above n is zero (DLMF 26.3.2)", 5, 6, 0},
		{"k above n, n zero", 0, 3, 0},
		{"poker hands", 52, 5, 2598960},
		{"largest central that fits", 67, 33, 14226520737620288370},
		{"max int choose 1", math.MaxInt, 1, math.MaxInt},
		{"max int choose max int", math.MaxInt, math.MaxInt, 1},
		{"max int choose max int minus 1", math.MaxInt, math.MaxInt - 1, math.MaxInt},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := Binomial(tt.n, tt.k); err != nil || got != tt.want {
				t.Errorf("Binomial(%d, %d) = (%d, %v), want (%d, nil)", tt.n, tt.k, got, err, tt.want)
			}
		})
	}
}

// Per k, the largest n with C(n,k) ≤ 2⁶⁴ − 1 (doc/catalog/combinatorics.md, exact integers): n fits,
// n+1 overflows. This proves no false overflow at the edge.
func TestBinomial_overflowEdges(t *testing.T) {
	for _, e := range []struct{ k, n int }{{2, 6074001000}, {3, 4801280}, {4, 145056}, {5, 18580}, {10, 386}, {20, 86}, {30, 68}, {33, 67}} {
		for _, n := range []int{e.n, e.n + 1} {
			got, err := Binomial(n, e.k)
			wantU64(t, "Binomial", got, err, new(big.Int).Binomial(int64(n), int64(e.k)))
		}
	}
	if _, err := Binomial(68, 34); !errors.Is(err, ErrOverflow) {
		t.Errorf("Binomial(68, 34) err = %v, want ErrOverflow", err)
	}
	if _, err := Binomial(math.MaxInt, math.MaxInt/2); !errors.Is(err, ErrOverflow) {
		t.Errorf("Binomial(MaxInt, MaxInt/2) err = %v, want ErrOverflow", err)
	}
}

func TestBinomial_errors(t *testing.T) {
	for _, c := range [][2]int{{-1, 0}, {0, -1}, {math.MinInt, 3}, {3, math.MinInt}} {
		if got, err := Binomial(c[0], c[1]); !errors.Is(err, ErrNegative) || got != 0 {
			t.Errorf("Binomial(%d, %d) = (%d, %v), want (0, ErrNegative)", c[0], c[1], got, err)
		}
	}
}

func TestBinomial_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for range 600 {
		n := rng.IntN(200)
		k := rng.IntN(n + 6)
		got, err := Binomial(n, k)
		wantU64(t, "Binomial", got, err, new(big.Int).Binomial(int64(n), int64(k)))
	}
}

func TestBinomial_noAlloc(t *testing.T) {
	if a := testing.AllocsPerRun(100, func() { sinkU, _ = Binomial(67, 33) }); a != 0 {
		t.Errorf("Binomial allocates %v times, want 0", a)
	}
}

func BenchmarkBinomial(b *testing.B) {
	for _, c := range []struct {
		name string
		n, k int
	}{{"n=52,k=5", 52, 5}, {"n=67,k=33", 67, 33}} {
		b.Run(c.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				sinkU, _ = Binomial(c.n, c.k)
			}
		})
	}
}

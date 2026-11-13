package dp

import (
	"errors"
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// clrsPrices is the price table of CLRS 3rd ed., §15.1, Figure 15.1 (lengths 1 … 10).
var clrsPrices = []int64{1, 5, 8, 9, 10, 17, 17, 20, 24, 30}

// rodSlow is CUT-ROD, the plain recursion of CLRS §15.1 (exponential; n ≤ 14), allowing waste.
func rodSlow(prices []int64, n int) int64 {
	var best int64
	for i := 1; i <= n; i++ {
		v := rodSlow(prices, n-i)
		if i <= len(prices) {
			v += prices[i-1]
		}
		best = max(best, v)
	}
	return best
}

// CLRS §15.1 lists the optimal revenues r1 … r10 = 1 5 8 10 13 17 18 22 25 30, and (EXTENDED-BOTTOM-UP-CUT-ROD
// with PRINT-CUT-ROD-SOLUTION) the cuts 4 → 2+2 and 7 → 1+6.
func TestRodCutting_clrsFigure15_1(t *testing.T) {
	want := []int64{0, 1, 5, 8, 10, 13, 17, 18, 22, 25, 30}
	for n, w := range want {
		got, pieces, err := RodCutting(clrsPrices, n)
		if err != nil || got != w {
			t.Errorf("RodCutting(clrs, %d) = (%d, %v), want %d", n, got, err, w)
		}
		sum, rev := 0, int64(0)
		for _, p := range pieces {
			sum += p
			rev += clrsPrices[p-1]
		}
		if sum > n || rev != got {
			t.Errorf("RodCutting(clrs, %d) pieces %v: length %d, revenue %d", n, pieces, sum, rev)
		}
	}
	for _, tt := range []struct {
		n    int
		want []int
	}{{4, []int{2, 2}}, {7, []int{1, 6}}, {10, []int{10}}} {
		if _, pieces, _ := RodCutting(clrsPrices, tt.n); !slices.Equal(pieces, tt.want) {
			t.Errorf("RodCutting(clrs, %d) pieces = %v, want %v", tt.n, pieces, tt.want)
		}
	}
}

func TestRodCutting_cases(t *testing.T) {
	tests := []struct {
		name   string
		prices []int64
		n      int
		want   int64
		pieces []int
	}{
		{"length zero", clrsPrices, 0, 0, []int{}},
		{"no prices: all waste", nil, 5, 0, []int{}},
		{"rod longer than the table; ties to the shortest first piece", []int64{2, 5}, 5, 12, []int{1, 2, 2}},
		{"zero prices", []int64{0, 0}, 3, 0, []int{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, pieces, err := RodCutting(tt.prices, tt.n)
			if err != nil || got != tt.want || pieces == nil || !slices.Equal(pieces, tt.pieces) {
				t.Errorf("RodCutting(%v, %d) = (%d, %#v, %v), want (%d, %v, nil)", tt.prices, tt.n, got, pieces, err, tt.want, tt.pieces)
			}
		})
	}
}

func TestRodCutting_errors(t *testing.T) {
	tests := []struct {
		name   string
		prices []int64
		n      int
		want   error
	}{
		{"negative length", clrsPrices, -1, ErrNegative},
		{"negative price", []int64{1, -1}, 3, ErrNegative},
		{"too large", clrsPrices, 1 << 40, ErrTooLarge},
		{"max int", clrsPrices, math.MaxInt, ErrTooLarge},
		{"overflow", []int64{math.MaxInt64 / 2}, 3, ErrOverflow},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, _, err := RodCutting(tt.prices, tt.n); !errors.Is(err, tt.want) {
				t.Errorf("RodCutting error = %v, want %v", err, tt.want)
			}
		})
	}
}

func TestRodCutting_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 300 {
		prices := benchInt64sFrom(rng, rng.IntN(8), 30)
		n := rng.IntN(13)
		got, pieces, err := RodCutting(prices, n)
		if want := rodSlow(prices, n); err != nil || got != want {
			t.Fatalf("trial %d: RodCutting(%v, %d) = (%d, %v), want %d", trial, prices, n, got, err, want)
		}
		sum, rev := 0, int64(0)
		for _, p := range pieces {
			sum += p
			rev += prices[p-1]
		}
		if sum > n || rev != got {
			t.Fatalf("trial %d: pieces %v do not give %d within %d", trial, pieces, got, n)
		}
	}
}

func benchInt64sFrom(rng *rand.Rand, n int, limit int64) []int64 {
	out := make([]int64, n)
	for i := range out {
		out[i] = rng.Int64N(limit)
	}
	return out
}

// Prices for every length up to the rod: the full O(n²) recurrence.
func BenchmarkRodCutting(b *testing.B) {
	for _, size := range quadraticSizes {
		prices := benchInt64s(size.n, 1000)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _, _ = RodCutting(prices, size.n)
			}
		})
	}
}

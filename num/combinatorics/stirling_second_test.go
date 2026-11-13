package combinatorics

import (
	"errors"
	"math"
	"math/big"
	"testing"
)

// DLMF Table 26.8.2: S(6,3) = 90, S(10,5) = 42 525.
func TestStirlingSecond_sourceExample(t *testing.T) {
	for _, c := range []struct {
		n, k int
		want uint64
	}{{6, 3, 90}, {10, 5, 42525}} {
		if got, err := StirlingSecond(c.n, c.k); err != nil || got != c.want {
			t.Errorf("StirlingSecond(%d, %d) = (%d, %v), want (%d, nil)", c.n, c.k, got, err, c.want)
		}
		if got, err := StirlingSecondBig(c.n, c.k); err != nil || got.Uint64() != c.want {
			t.Errorf("StirlingSecondBig(%d, %d) = (%v, %v), want (%d, nil)", c.n, c.k, got, err, c.want)
		}
	}
}

func TestStirlingSecond_cases(t *testing.T) {
	tests := []struct {
		name string
		n, k int
		want uint64
	}{
		{"zero zero", 0, 0, 1},
		{"k zero, n positive", 5, 0, 0},
		{"k above n", 3, 5, 0},
		{"k equals n", math.MaxInt, math.MaxInt, 1},
		{"k = 1", math.MaxInt, 1, 1},
		{"k = n−1 is C(n,2)", 1_000_000_000, 999_999_999, 499999999500000000},
		{"{64 2} = 2⁶³ − 1", 64, 2, 1<<63 - 1},
		{"{65 2} = 2⁶⁴ − 1, the band's widest fitting case", 65, 2, math.MaxUint64},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := StirlingSecond(tt.n, tt.k); err != nil || got != tt.want {
				t.Errorf("StirlingSecond(%d, %d) = (%d, %v), want (%d, nil)", tt.n, tt.k, got, err, tt.want)
			}
		})
	}
	for _, c := range [][2]int{{66, 2}, {100, 50}, {math.MaxInt, 2}, {math.MaxInt, math.MaxInt - 2}} {
		if got, err := StirlingSecond(c[0], c[1]); !errors.Is(err, ErrOverflow) || got != 0 {
			t.Errorf("StirlingSecond(%d, %d) = (%d, %v), want (0, ErrOverflow)", c[0], c[1], got, err)
		}
	}
}

func TestStirlingSecond_errors(t *testing.T) {
	for _, c := range [][2]int{{-1, 0}, {0, -1}, {math.MinInt, 1}} {
		if got, err := StirlingSecond(c[0], c[1]); !errors.Is(err, ErrNegative) || got != 0 {
			t.Errorf("StirlingSecond(%d, %d) = (%d, %v), want (0, ErrNegative)", c[0], c[1], got, err)
		}
		if got, err := StirlingSecondBig(c[0], c[1]); !errors.Is(err, ErrNegative) || got != nil {
			t.Errorf("StirlingSecondBig(%d, %d) = (%v, %v), want (nil, ErrNegative)", c[0], c[1], got, err)
		}
	}
}

// Every (n, k) with n ≤ 90, k ≤ n+1, against the naive big triangle.
func TestStirlingSecond_oracle(t *testing.T) {
	const top = 90
	_, second := stirlingTriangle(top)
	for n := 0; n <= top; n++ {
		for k := 0; k <= n+1; k++ {
			exact := new(big.Int)
			if k <= n {
				exact = second[n][k]
			}
			got, err := StirlingSecond(n, k)
			wantU64(t, "StirlingSecond", got, err, exact)
			gotBig, err := StirlingSecondBig(n, k)
			if err != nil || gotBig.Cmp(exact) != 0 {
				t.Fatalf("StirlingSecondBig(%d, %d) = (%v, %v), want %v", n, k, gotBig, err, exact)
			}
		}
	}
}

// DLMF 26.8.6, the explicit sum {n k} = (1/k!)Σⱼ (−1)^(k−j) C(k,j) jⁿ, in math/big: an oracle independent
// of the recurrence, for wider n than the triangle.
func TestStirlingSecond_explicitSum(t *testing.T) {
	for _, c := range [][2]int{{200, 3}, {200, 197}, {150, 75}} {
		n, k := c[0], c[1]
		sum := new(big.Int)
		for j := 0; j <= k; j++ {
			term := new(big.Int).Exp(big.NewInt(int64(j)), big.NewInt(int64(n)), nil)
			term.Mul(term, new(big.Int).Binomial(int64(k), int64(j)))
			if (k-j)%2 == 1 {
				sum.Sub(sum, term)
			} else {
				sum.Add(sum, term)
			}
		}
		sum.Quo(sum, bigFactorial(k))
		got, err := StirlingSecondBig(n, k)
		if err != nil || got.Cmp(sum) != 0 {
			t.Errorf("StirlingSecondBig(%d, %d) = (%v, %v), want %v", n, k, got, err, sum)
		}
	}
}

func TestStirlingSecond_noAlloc(t *testing.T) {
	if a := testing.AllocsPerRun(100, func() { sinkU, _ = StirlingSecond(40, 35) }); a != 0 {
		t.Errorf("StirlingSecond allocates %v times, want 0", a)
	}
}

func BenchmarkStirlingSecond(b *testing.B) {
	for _, c := range []struct {
		name string
		n, k int
	}{{"n=26,k=13", 26, 13}, {"n=65,k=2", 65, 2}} {
		b.Run(c.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				sinkU, _ = StirlingSecond(c.n, c.k)
			}
		})
	}
}

func BenchmarkStirlingSecondBig(b *testing.B) {
	for _, c := range []struct {
		name string
		n, k int
	}{{"n=100,k=50", 100, 50}, {"n=1000,k=500", 1000, 500}} {
		b.Run(c.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				sinkBig, _ = StirlingSecondBig(c.n, c.k)
			}
		})
	}
}

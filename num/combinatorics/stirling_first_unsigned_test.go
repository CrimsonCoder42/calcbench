package combinatorics

import (
	"errors"
	"math"
	"math/big"
	"testing"
)

// stirlingTriangle returns the full triangles [i j] and {i j}, i, j ≤ n, in math/big, built naively
// from DLMF 26.8.18 (unsigned) and 26.8.22. The oracle for both kinds.
func stirlingTriangle(n int) (first, second [][]*big.Int) {
	first = make([][]*big.Int, n+1)
	second = make([][]*big.Int, n+1)
	for i := range n + 1 {
		first[i] = make([]*big.Int, n+1)
		second[i] = make([]*big.Int, n+1)
		for j := range n + 1 {
			first[i][j], second[i][j] = new(big.Int), new(big.Int)
		}
	}
	first[0][0].SetInt64(1)
	second[0][0].SetInt64(1)
	for i := 1; i <= n; i++ {
		for j := 1; j <= i; j++ {
			t := new(big.Int).Mul(big.NewInt(int64(i-1)), first[i-1][j])
			first[i][j].Add(t, first[i-1][j-1])
			t = new(big.Int).Mul(big.NewInt(int64(j)), second[i-1][j])
			second[i][j].Add(t, second[i-1][j-1])
		}
	}
	return first, second
}

// DLMF Table 26.8.1: s(6,3) = −225 and s(10,5) = −269 325, so [6 3] = 225 and [10 5] = 269 325.
func TestStirlingFirstUnsigned_sourceExample(t *testing.T) {
	for _, c := range []struct {
		n, k int
		want uint64
	}{{6, 3, 225}, {10, 5, 269325}} {
		if got, err := StirlingFirstUnsigned(c.n, c.k); err != nil || got != c.want {
			t.Errorf("StirlingFirstUnsigned(%d, %d) = (%d, %v), want (%d, nil)", c.n, c.k, got, err, c.want)
		}
		if got, err := StirlingFirstUnsignedBig(c.n, c.k); err != nil || got.Uint64() != c.want {
			t.Errorf("StirlingFirstUnsignedBig(%d, %d) = (%v, %v), want (%d, nil)", c.n, c.k, got, err, c.want)
		}
	}
}

func TestStirlingFirstUnsigned_cases(t *testing.T) {
	tests := []struct {
		name string
		n, k int
		want uint64
	}{
		{"zero zero", 0, 0, 1},
		{"k zero, n positive", 5, 0, 0},
		{"k above n", 3, 5, 0},
		{"n zero, k positive", 0, 3, 0},
		{"k equals n", math.MaxInt, math.MaxInt, 1},
		{"k = n−1 is C(n,2)", 1_000_000_000, 999_999_999, 499999999500000000},
		{"k = 1 is (n−1)!", 21, 1, 2432902008176640000},
		{"k = 2", 5, 2, 50},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := StirlingFirstUnsigned(tt.n, tt.k); err != nil || got != tt.want {
				t.Errorf("StirlingFirstUnsigned(%d, %d) = (%d, %v), want (%d, nil)", tt.n, tt.k, got, err, tt.want)
			}
		})
	}
	for _, c := range [][2]int{{22, 1}, {100, 2}, {100, 36}, {math.MaxInt, 2}, {math.MaxInt, math.MaxInt - 2}, {math.MaxInt, math.MaxInt - 1}} {
		if got, err := StirlingFirstUnsigned(c[0], c[1]); !errors.Is(err, ErrOverflow) || got != 0 {
			t.Errorf("StirlingFirstUnsigned(%d, %d) = (%d, %v), want (0, ErrOverflow)", c[0], c[1], got, err)
		}
	}
}

func TestStirlingFirstUnsigned_errors(t *testing.T) {
	for _, c := range [][2]int{{-1, 0}, {0, -1}, {math.MinInt, 1}} {
		if got, err := StirlingFirstUnsigned(c[0], c[1]); !errors.Is(err, ErrNegative) || got != 0 {
			t.Errorf("StirlingFirstUnsigned(%d, %d) = (%d, %v), want (0, ErrNegative)", c[0], c[1], got, err)
		}
		if got, err := StirlingFirstUnsignedBig(c[0], c[1]); !errors.Is(err, ErrNegative) || got != nil {
			t.Errorf("StirlingFirstUnsignedBig(%d, %d) = (%v, %v), want (nil, ErrNegative)", c[0], c[1], got, err)
		}
	}
}

// Every (n, k) with n ≤ 90, k ≤ n+1, against the naive big triangle: the uint64 value or ErrOverflow, and
// the big twin exactly (both of its band layouts: k ≤ n−k and k > n−k).
func TestStirlingFirstUnsigned_oracle(t *testing.T) {
	const top = 90
	first, _ := stirlingTriangle(top)
	for n := 0; n <= top; n++ {
		for k := 0; k <= n+1; k++ {
			exact := new(big.Int)
			if k <= n {
				exact = first[n][k]
			}
			got, err := StirlingFirstUnsigned(n, k)
			wantU64(t, "StirlingFirstUnsigned", got, err, exact)
			gotBig, err := StirlingFirstUnsignedBig(n, k)
			if err != nil || gotBig.Cmp(exact) != 0 {
				t.Fatalf("StirlingFirstUnsignedBig(%d, %d) = (%v, %v), want %v", n, k, gotBig, err, exact)
			}
		}
	}
}

func TestStirlingFirstUnsigned_noAlloc(t *testing.T) {
	if a := testing.AllocsPerRun(100, func() { sinkU, _ = StirlingFirstUnsigned(40, 35) }); a != 0 {
		t.Errorf("StirlingFirstUnsigned allocates %v times, want 0", a)
	}
}

func BenchmarkStirlingFirstUnsigned(b *testing.B) {
	for _, c := range []struct {
		name string
		n, k int
	}{{"n=20,k=10", 20, 10}, {"n=80,k=70", 80, 70}} {
		b.Run(c.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				sinkU, _ = StirlingFirstUnsigned(c.n, c.k)
			}
		})
	}
}

func BenchmarkStirlingFirstUnsignedBig(b *testing.B) {
	for _, c := range []struct {
		name string
		n, k int
	}{{"n=100,k=50", 100, 50}, {"n=1000,k=500", 1000, 500}} {
		b.Run(c.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				sinkBig, _ = StirlingFirstUnsignedBig(c.n, c.k)
			}
		})
	}
}

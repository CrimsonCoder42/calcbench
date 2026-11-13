package dp

import (
	"errors"
	"math"
	"math/rand/v2"
	"strconv"
	"testing"
)

// chainSlow is RECURSIVE-MATRIX-CHAIN of CLRS §15.3 (exponential; n ≤ 8).
func chainSlow(dims []int, i, j int) int64 {
	if i == j {
		return 0
	}
	best := int64(math.MaxInt64)
	for k := i; k < j; k++ {
		best = min(best, chainSlow(dims, i, k)+chainSlow(dims, k+1, j)+int64(dims[i]*dims[k+1]*dims[j+1]))
	}
	return best
}

// parensCost parses a parenthesisation such as "((A1A2)A3)" and returns its multiplication count and
// the dimensions (rows, cols) of its product, checking that it names A1 … An once each, in order.
func parensCost(t *testing.T, s string, dims []int) int64 {
	t.Helper()
	pos, next := 0, 1
	var parse func() (int64, int, int)
	parse = func() (int64, int, int) {
		if pos >= len(s) {
			t.Fatalf("parens %q: unexpected end", s)
		}
		if s[pos] == 'A' {
			pos++
			start := pos
			for pos < len(s) && s[pos] >= '0' && s[pos] <= '9' {
				pos++
			}
			k, err := strconv.Atoi(s[start:pos])
			if err != nil || k != next {
				t.Fatalf("parens %q: matrix A%s out of order (want A%d)", s, s[start:pos], next)
			}
			next++
			return 0, dims[k-1], dims[k]
		}
		if s[pos] != '(' {
			t.Fatalf("parens %q: unexpected %q at %d", s, s[pos], pos)
		}
		pos++
		c1, r1, k1 := parse()
		c2, r2, k2 := parse()
		if pos >= len(s) || s[pos] != ')' || k1 != r2 {
			t.Fatalf("parens %q: bad product at %d", s, pos)
		}
		pos++
		return c1 + c2 + int64(r1*k1*k2), r1, k2
	}
	cost, _, _ := parse()
	if pos != len(s) || next != len(dims) {
		t.Fatalf("parens %q: trailing text or missing matrices", s)
	}
	return cost
}

// CLRS 3rd ed., §15.2, Figure 15.5: dimensions 30 35 15 5 10 20 25 need m[1,6] = 15125 scalar
// multiplications with the parenthesisation ((A1(A2A3))((A4A5)A6)).
func TestMatrixChainOrder_clrsFigure15_5(t *testing.T) {
	cost, parens, err := MatrixChainOrder([]int{30, 35, 15, 5, 10, 20, 25})
	if err != nil || cost != 15125 || parens != "((A1(A2A3))((A4A5)A6))" {
		t.Errorf("MatrixChainOrder(CLRS) = (%d, %q, %v), want (15125, ((A1(A2A3))((A4A5)A6)), nil)", cost, parens, err)
	}
}

func TestMatrixChainOrder_cases(t *testing.T) {
	tests := []struct {
		name   string
		dims   []int
		cost   int64
		parens string
	}{
		{"one matrix", []int{10, 20}, 0, "A1"},
		{"two matrices", []int{10, 20, 30}, 6000, "(A1A2)"},
		{"three, left first", []int{10, 100, 5, 50}, 7500, "((A1A2)A3)"}, // CLRS §15.2 opening example
		{"all ones: ties to the smallest split", []int{1, 1, 1, 1}, 2, "(A1(A2A3))"},
		{"extreme dimension fits", []int{math.MaxInt32, 1, 1}, math.MaxInt32, "(A1A2)"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cost, parens, err := MatrixChainOrder(tt.dims)
			if err != nil || cost != tt.cost || parens != tt.parens {
				t.Errorf("MatrixChainOrder(%v) = (%d, %q, %v), want (%d, %q, nil)", tt.dims, cost, parens, err, tt.cost, tt.parens)
			}
		})
	}
}

func TestMatrixChainOrder_errors(t *testing.T) {
	long := make([]int, 5000)
	for i := range long {
		long[i] = 1
	}
	tests := []struct {
		name string
		dims []int
		want error
	}{
		{"no dims", nil, ErrEmpty},
		{"one dim", []int{5}, ErrEmpty},
		{"zero dim", []int{3, 0, 4}, ErrNotPositive},
		{"negative dim", []int{-3, 4}, ErrNotPositive},
		{"tables too large", long, ErrTooLarge},
		{"overflow", []int{1 << 30, 1 << 30, 1 << 30, 1 << 30}, ErrOverflow},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, _, err := MatrixChainOrder(tt.dims); !errors.Is(err, tt.want) {
				t.Errorf("MatrixChainOrder error = %v, want %v", err, tt.want)
			}
		})
	}
}

func TestMatrixChainOrder_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 300 {
		dims := make([]int, 2+rng.IntN(7))
		for i := range dims {
			dims[i] = 1 + rng.IntN(30)
		}
		cost, parens, err := MatrixChainOrder(dims)
		if want := chainSlow(dims, 0, len(dims)-2); err != nil || cost != want {
			t.Fatalf("trial %d: MatrixChainOrder(%v) = (%d, %v), want %d", trial, dims, cost, err, want)
		}
		if got := parensCost(t, parens, dims); got != cost {
			t.Fatalf("trial %d: parens %q cost %d, reported %d", trial, parens, got, cost)
		}
	}
}

// n matrices: O(n³) time, so n = 50, 100, 200 (8× the work per doubling).
func BenchmarkMatrixChainOrder(b *testing.B) {
	for _, n := range []int{50, 100, 200} {
		rng := rand.New(rand.NewPCG(1, 2))
		dims := make([]int, n+1)
		for i := range dims {
			dims[i] = 1 + rng.IntN(100)
		}
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _, _ = MatrixChainOrder(dims)
			}
		})
	}
}

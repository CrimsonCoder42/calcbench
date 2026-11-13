package integer

import (
	"errors"
	"math"
	"math/big"
	"math/rand/v2"
	"testing"
)

// bigCeilDiv is the oracle: ⌈a/b⌉ = −⌊−a/b⌋ (Concrete Mathematics §3.1 eq. (3.4); Exercise 3.12), with ⌊n/d⌋ for d > 0 taken from
// big.Int.Div, which is Euclidean and so equals floor division when the divisor is positive.
func bigCeilDiv(a, b int64) *big.Int {
	n, d := big.NewInt(a), big.NewInt(b)
	if d.Sign() < 0 {
		n.Neg(n)
		d.Neg(d)
	}
	n.Neg(n)
	q := new(big.Int).Div(n, d)
	return q.Neg(q)
}

// Worked examples from the BRIEF, each checked by hand against ⌈x⌉ = −⌊−x⌋ (Concrete Mathematics §3.1 eq. (3.4); Exercise 3.12):
// ⌈10/3⌉ = 4, ⌈9/3⌉ = 3, ⌈0/5⌉ = 0, ⌈−7/2⌉ = ⌈−3.5⌉ = −3, ⌈MaxInt64/1⌉ = MaxInt64, ⌈(2⁵³+1)/1⌉ = 2⁵³+1.
func TestCeilDiv_workedExamples(t *testing.T) {
	tests := []struct{ a, b, want int64 }{
		{10, 3, 4},
		{9, 3, 3},
		{0, 5, 0},
		{-7, 2, -3},
		{math.MaxInt64, 1, math.MaxInt64},
		{1<<53 + 1, 1, 1<<53 + 1},
	}
	for _, tc := range tests {
		if got, err := CeilDiv(tc.a, tc.b); err != nil || got != tc.want {
			t.Errorf("CeilDiv(%d, %d) = (%d, %v), want %d", tc.a, tc.b, got, err, tc.want)
		}
	}
}

func TestCeilDiv_cases(t *testing.T) {
	tests := []struct {
		name       string
		a, b, want int64
	}{
		{"+/+ inexact", 7, 2, 4},
		{"+/- inexact", 7, -2, -3},
		{"-/+ inexact", -7, 2, -3},
		{"-/- inexact", -7, -2, 4},
		{"+/+ exact", 8, 2, 4},
		{"+/- exact", 8, -2, -4},
		{"-/+ exact", -8, 2, -4},
		{"-/- exact", -8, -2, 4},
		{"small positive", 1, 5, 1},
		{"small negative", -1, 5, 0},
		{"small negative over negative", -1, -5, 1},
		{"small positive over negative", 1, -5, 0},
		{"zero over negative", 0, -5, 0},
		{"MinInt64 / 1", math.MinInt64, 1, math.MinInt64},
		{"MinInt64 / 2", math.MinInt64, 2, math.MinInt64 / 2},
		{"MaxInt64 / -1", math.MaxInt64, -1, -math.MaxInt64},
		{"MaxInt64 / 2", math.MaxInt64, 2, 1 << 62},
		{"MinInt64 / MinInt64", math.MinInt64, math.MinInt64, 1},
		{"MaxInt64 / MaxInt64", math.MaxInt64, math.MaxInt64, 1},
		{"1 / MinInt64", 1, math.MinInt64, 0},
		{"-1 / MinInt64", -1, math.MinInt64, 1},
		{"MaxInt64 / MinInt64", math.MaxInt64, math.MinInt64, 0},
		{"MinInt64 / MaxInt64", math.MinInt64, math.MaxInt64, -1},
		{"MinInt64 / -2", math.MinInt64, -2, 1 << 62},
	}
	for _, tc := range tests {
		if got, err := CeilDiv(tc.a, tc.b); err != nil || got != tc.want {
			t.Errorf("%s: CeilDiv(%d, %d) = (%d, %v), want %d", tc.name, tc.a, tc.b, got, err, tc.want)
		}
	}
}

func TestCeilDiv_errors(t *testing.T) {
	for _, a := range []int64{0, 1, -1, math.MaxInt64, math.MinInt64} {
		if _, err := CeilDiv(a, 0); !errors.Is(err, ErrZeroDenominator) {
			t.Errorf("CeilDiv(%d, 0): err = %v, want ErrZeroDenominator", a, err)
		}
	}
	if _, err := CeilDiv(math.MinInt64, -1); !errors.Is(err, ErrOverflow) {
		t.Errorf("CeilDiv(MinInt64, -1): err = %v, want ErrOverflow", err)
	}
}

func TestCeilDiv_oracle(t *testing.T) {
	r := rand.New(rand.NewPCG(1, 2))
	extremes := []int64{0, 1, -1, 2, -2, 3, -3, math.MaxInt64, math.MaxInt64 - 1, math.MinInt64, math.MinInt64 + 1}
	gen := func() int64 {
		switch r.IntN(3) {
		case 0:
			return extremes[r.IntN(len(extremes))]
		case 1:
			return r.Int64N(2001) - 1000
		default:
			return int64(r.Uint64())
		}
	}
	for range 600 {
		a, b := gen(), gen()
		if b == 0 || (a == math.MinInt64 && b == -1) {
			continue
		}
		got, err := CeilDiv(a, b)
		if want := bigCeilDiv(a, b); err != nil || !want.IsInt64() || want.Int64() != got {
			t.Fatalf("CeilDiv(%d, %d) = (%d, %v), want %s", a, b, got, err, want)
		}
	}
}

func FuzzCeilDiv(f *testing.F) {
	f.Add(int64(10), int64(3))
	f.Add(int64(-7), int64(2))
	f.Add(int64(math.MinInt64), int64(-1))
	f.Add(int64(math.MinInt64), int64(math.MinInt64))
	f.Add(int64(1), int64(0))
	f.Fuzz(func(t *testing.T, a, b int64) {
		got, err := CeilDiv(a, b)
		switch {
		case b == 0:
			if !errors.Is(err, ErrZeroDenominator) {
				t.Fatalf("CeilDiv(%d, 0): err = %v, want ErrZeroDenominator", a, err)
			}
			return
		case a == math.MinInt64 && b == -1:
			if !errors.Is(err, ErrOverflow) {
				t.Fatalf("CeilDiv(MinInt64, -1): err = %v, want ErrOverflow", err)
			}
			return
		}
		if want := bigCeilDiv(a, b); err != nil || !want.IsInt64() || want.Int64() != got {
			t.Fatalf("CeilDiv(%d, %d) = (%d, %v), want %s", a, b, got, err, want)
		}
	})
}

func TestCeilDiv_noAlloc(t *testing.T) {
	if a := testing.AllocsPerRun(100, func() { _, _ = CeilDiv(-7, 2) }); a != 0 {
		t.Errorf("CeilDiv allocates %v", a)
	}
}

func BenchmarkCeilDiv(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_, _ = CeilDiv(-1234567890123, 97)
	}
}

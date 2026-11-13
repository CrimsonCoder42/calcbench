package business

import (
	"errors"
	"math"
	"testing"
)

func TestReturnOnAssets_sourceExample(t *testing.T) {
	cases := []struct {
		num, den  int64
		want, tol float64
	}{
		// Principles of Finance §6.6, Clear Lake Sporting Goods: $35,000 ÷ $225,000, printed 0.16 (tolerance: half of the printed last digit)
		{3500000, 22500000, 0.16, 0.005},
	}
	for _, c := range cases {
		got, err := ReturnOnAssets(c.num, c.den)
		if err != nil {
			t.Fatalf("%d/%d: %v", c.num, c.den, err)
		}
		if math.Abs(got-c.want) > c.tol {
			t.Errorf("%d/%d = %v, want %v ± %v", c.num, c.den, got, c.want, c.tol)
		}
		if exact := float64(c.num) / float64(c.den); got != exact {
			t.Errorf("%d/%d = %v, want the correctly rounded quotient %v", c.num, c.den, got, exact)
		}
	}
}

func TestReturnOnAssets_table(t *testing.T) {
	cases := []struct {
		name     string
		num, den int64
		want     float64
	}{
		{"zero numerator", 0, 5, 0},
		{"one", 7, 7, 1},
		{"max numerator", maxExactCents, 1, maxExactCents},
		{"max denominator", 1, maxExactCents, 1.0 / maxExactCents},
		{"negative numerator", -3, 4, -0.75},
		{"min numerator", -maxExactCents, 1, -maxExactCents},
	}
	for _, c := range cases {
		got, err := ReturnOnAssets(c.num, c.den)
		if err != nil || got != c.want {
			t.Errorf("%s: got %v, %v; want %v", c.name, got, err, c.want)
		}
	}
}

func TestReturnOnAssets_errors(t *testing.T) {
	cases := []struct {
		name     string
		num, den int64
		want     error
	}{
		{"zero denominator", 1, 0, ErrZeroDenominator},
		{"numerator above 2^53", maxExactCents + 1, 1, ErrOutOfRange},
		{"numerator below -2^53", -maxExactCents - 1, 1, ErrOutOfRange},
		{"denominator above 2^53", 1, maxExactCents + 1, ErrOutOfRange},
		{"denominator min int64", 1, math.MinInt64, ErrOutOfRange},
		{"negative denominator", 1, -1, ErrNegative},
	}
	for _, c := range cases {
		if _, err := ReturnOnAssets(c.num, c.den); !errors.Is(err, c.want) {
			t.Errorf("%s: err = %v, want %v", c.name, err, c.want)
		}
	}
}

func TestReturnOnAssets_oracle(t *testing.T) {
	checkRatioOracle(t, ReturnOnAssets, -maxExactCents, 1)
}

func TestReturnOnAssets_allocs(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkFloat, _ = ReturnOnAssets(3500000, 22500000) }); n != 0 {
		t.Fatalf("allocs = %v, want 0", n)
	}
}

func BenchmarkReturnOnAssets(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = ReturnOnAssets(3500000, 22500000)
	}
}

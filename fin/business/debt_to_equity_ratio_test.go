package business

import (
	"errors"
	"math"
	"testing"
)

func TestDebtToEquityRatio_sourceExample(t *testing.T) {
	cases := []struct {
		num, den  int64
		want, tol float64
	}{
		// Principles of Finance §6.4, Clear Lake Sporting Goods: $150,000 ÷ $100,000 = 1.5
		{15000000, 10000000, 1.5, 0},
	}
	for _, c := range cases {
		got, err := DebtToEquityRatio(c.num, c.den)
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

func TestDebtToEquityRatio_table(t *testing.T) {
	cases := []struct {
		name     string
		num, den int64
		want     float64
	}{
		{"zero numerator", 0, 5, 0},
		{"one", 7, 7, 1},
		{"max numerator", maxExactCents, 1, maxExactCents},
		{"max denominator", 1, maxExactCents, 1.0 / maxExactCents},
		{"negative denominator", 3, -4, -0.75},
	}
	for _, c := range cases {
		got, err := DebtToEquityRatio(c.num, c.den)
		if err != nil || got != c.want {
			t.Errorf("%s: got %v, %v; want %v", c.name, got, err, c.want)
		}
	}
}

func TestDebtToEquityRatio_errors(t *testing.T) {
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
		{"negative numerator", -1, 1, ErrNegative},
	}
	for _, c := range cases {
		if _, err := DebtToEquityRatio(c.num, c.den); !errors.Is(err, c.want) {
			t.Errorf("%s: err = %v, want %v", c.name, err, c.want)
		}
	}
}

func TestDebtToEquityRatio_oracle(t *testing.T) {
	checkRatioOracle(t, DebtToEquityRatio, 0, -maxExactCents)
}

func TestDebtToEquityRatio_allocs(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkFloat, _ = DebtToEquityRatio(15000000, 10000000) }); n != 0 {
		t.Fatalf("allocs = %v, want 0", n)
	}
}

func BenchmarkDebtToEquityRatio(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = DebtToEquityRatio(15000000, 10000000)
	}
}

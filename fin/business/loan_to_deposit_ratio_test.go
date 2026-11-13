package business

import (
	"errors"
	"math"
	"testing"
)

func TestLoanToDepositRatio_sourceExample(t *testing.T) {
	cases := []struct {
		num, den  int64
		want, tol float64
	}{
		// hand-checked: $800,000 ÷ $1,000,000 = 0.8 (the UBPR guide states the formula without numbers)
		{80000000, 100000000, 0.8, 0},
	}
	for _, c := range cases {
		got, err := LoanToDepositRatio(c.num, c.den)
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

func TestLoanToDepositRatio_table(t *testing.T) {
	cases := []struct {
		name     string
		num, den int64
		want     float64
	}{
		{"zero numerator", 0, 5, 0},
		{"one", 7, 7, 1},
		{"max numerator", maxExactCents, 1, maxExactCents},
		{"max denominator", 1, maxExactCents, 1.0 / maxExactCents},
	}
	for _, c := range cases {
		got, err := LoanToDepositRatio(c.num, c.den)
		if err != nil || got != c.want {
			t.Errorf("%s: got %v, %v; want %v", c.name, got, err, c.want)
		}
	}
}

func TestLoanToDepositRatio_errors(t *testing.T) {
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
		{"negative denominator", 1, -1, ErrNegative},
	}
	for _, c := range cases {
		if _, err := LoanToDepositRatio(c.num, c.den); !errors.Is(err, c.want) {
			t.Errorf("%s: err = %v, want %v", c.name, err, c.want)
		}
	}
}

func TestLoanToDepositRatio_oracle(t *testing.T) {
	checkRatioOracle(t, LoanToDepositRatio, 0, 1)
}

func TestLoanToDepositRatio_allocs(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkFloat, _ = LoanToDepositRatio(80000000, 100000000) }); n != 0 {
		t.Fatalf("allocs = %v, want 0", n)
	}
}

func BenchmarkLoanToDepositRatio(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = LoanToDepositRatio(80000000, 100000000)
	}
}

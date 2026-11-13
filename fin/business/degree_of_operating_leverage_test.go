package business

import (
	"errors"
	"math"
	"testing"
)

func TestDegreeOfOperatingLeverage_sourceExample(t *testing.T) {
	cases := []struct {
		num, den  int64
		want, tol float64
	}{
		// Principles of Accounting Vol. 2 §3.5, Company A: $148,000 ÷ $85,000, printed 1.74 (tolerance: half of the printed last digit)
		{14800000, 8500000, 1.74, 0.005},
		// Principles of Accounting Vol. 2 §3.5, Company B: $210,000 ÷ $85,000, printed 2.47
		{21000000, 8500000, 2.47, 0.005},
	}
	for _, c := range cases {
		got, err := DegreeOfOperatingLeverage(c.num, c.den)
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

func TestDegreeOfOperatingLeverage_table(t *testing.T) {
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
		{"negative denominator", 3, -4, -0.75},
	}
	for _, c := range cases {
		got, err := DegreeOfOperatingLeverage(c.num, c.den)
		if err != nil || got != c.want {
			t.Errorf("%s: got %v, %v; want %v", c.name, got, err, c.want)
		}
	}
}

func TestDegreeOfOperatingLeverage_errors(t *testing.T) {
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
	}
	for _, c := range cases {
		if _, err := DegreeOfOperatingLeverage(c.num, c.den); !errors.Is(err, c.want) {
			t.Errorf("%s: err = %v, want %v", c.name, err, c.want)
		}
	}
}

func TestDegreeOfOperatingLeverage_oracle(t *testing.T) {
	checkRatioOracle(t, DegreeOfOperatingLeverage, -maxExactCents, -maxExactCents)
}

func TestDegreeOfOperatingLeverage_allocs(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkFloat, _ = DegreeOfOperatingLeverage(14800000, 8500000) }); n != 0 {
		t.Fatalf("allocs = %v, want 0", n)
	}
}

func BenchmarkDegreeOfOperatingLeverage(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = DegreeOfOperatingLeverage(14800000, 8500000)
	}
}

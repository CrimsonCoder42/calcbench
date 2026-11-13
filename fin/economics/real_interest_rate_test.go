package economics

import (
	"errors"
	"math"
	"math/big"
	"testing"
)

// Hand-checked from OpenStax, Principles of Finance, §3.4 "Interest Rates": $1,000 saved at 6 % becomes
// $1,060 while the TV rises 2 % to $1,020; the exact real gain is 1060/1020 − 1 = 0.04/1.02. The book
// quotes the approximation 4 %. Tolerance 1e-16: two roundings on values near 0.04.
func TestRealInterestRate_example(t *testing.T) {
	got, err := RealInterestRate(0.06, 0.02)
	if err != nil {
		t.Fatal(err)
	}
	if want := 0.04 / 1.02; math.Abs(got-want) > 1e-16 {
		t.Fatalf("got %v, want %v", got, want)
	}
	if math.Abs(got-(1060.0/1020.0-1)) > 1e-15 {
		t.Fatalf("got %v, want 1060/1020 − 1", got)
	}
}

// Fisher 1896, ch. II §3, eq. (2): 1 + i = (1 + r)(1 + π). The round trip must return the nominal
// rate to 1e-15 (a few ulps of values below 1).
func TestRealInterestRate_fisherRoundTrip(t *testing.T) {
	for _, c := range [][2]float64{{0.06, 0.02}, {0.15, 0.12}, {0.01, 0.05}, {-0.5, 0.3}, {0.2, -0.1}, {0, 0}} {
		r, err := RealInterestRate(c[0], c[1])
		if err != nil {
			t.Fatal(err)
		}
		if back := (1+r)*(1+c[1]) - 1; math.Abs(back-c[0]) > 1e-15 {
			t.Errorf("i=%v π=%v: (1+r)(1+π)−1 = %v", c[0], c[1], back)
		}
	}
}

// Oracle: (1+i)/(1+π) − 1 in 256-bit big.Float, compared to 2 ulps of the float64 result.
func TestRealInterestRate_bigOracle(t *testing.T) {
	for _, c := range [][2]float64{{0.0500001, 0.05}, {0.03, 0.0299999999}, {0.9, 0.1}, {1e-12, 2e-12}} {
		got, _ := RealInterestRate(c[0], c[1])
		one := new(big.Float).SetPrec(256).SetInt64(1)
		num := new(big.Float).SetPrec(256).Add(one, big.NewFloat(c[0]))
		den := new(big.Float).SetPrec(256).Add(one, big.NewFloat(c[1]))
		w, _ := new(big.Float).SetPrec(256).Sub(new(big.Float).SetPrec(256).Quo(num, den), one).Float64()
		if math.Abs(got-w) > 2*math.Abs(math.Nextafter(w, math.Inf(1))-w) {
			t.Errorf("i=%v π=%v: got %v, want %v", c[0], c[1], got, w)
		}
	}
}

func TestRealInterestRate_cases(t *testing.T) {
	cases := []struct {
		name               string
		nominal, inflation float64
		want               float64
	}{
		{"zeroBoth", 0, 0, 0},
		{"equalRates", 0.07, 0.07, 0},
		{"zeroInflation", 0.05, 0, 0.05},
		{"deflation", 0, -0.5, 1},
	}
	for _, c := range cases {
		got, err := RealInterestRate(c.nominal, c.inflation)
		if err != nil || got != c.want {
			t.Errorf("%s: got %v, %v; want %v", c.name, got, err, c.want)
		}
	}
}

func TestRealInterestRate_errors(t *testing.T) {
	cases := []struct {
		name               string
		nominal, inflation float64
		want               error
	}{
		{"nominalNaN", math.NaN(), 0, ErrNotFinite},
		{"inflationInf", 0, math.Inf(1), ErrNotFinite},
		{"nominalNegInfIsRate", math.Inf(-1), 0, ErrRate},
		{"nominalMinusOne", -1, 0, ErrRate},
		{"inflationMinusOne", 0.05, -1, ErrRate},
		{"inflationBelowMinusOne", 0.05, -2, ErrRate},
		{"resultOverflow", 1e308, -0.9999999999999999, ErrOverflow},
	}
	for _, c := range cases {
		if _, err := RealInterestRate(c.nominal, c.inflation); !errors.Is(err, c.want) {
			t.Errorf("%s: err = %v, want %v", c.name, err, c.want)
		}
	}
}

func TestRealInterestRate_noAlloc(t *testing.T) {
	if a := testing.AllocsPerRun(100, func() { _, _ = RealInterestRate(0.06, 0.02) }); a != 0 {
		t.Fatalf("allocs = %v, want 0", a)
	}
}

func BenchmarkRealInterestRate(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = RealInterestRate(0.06, 0.02)
	}
}

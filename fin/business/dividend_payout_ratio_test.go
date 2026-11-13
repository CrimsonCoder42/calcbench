package business

import (
	"errors"
	"math"
	"testing"
)

// Hand-checked (the CFA ratio list gives no numbers): dividends $250,000 on net income $1,000,000 →
// 250,000 ÷ 1,000,000 = 0.25. Exact: both inputs and the quotient 1/4 are exact in float64.
func TestDividendPayoutRatio_sourceExample(t *testing.T) {
	got, err := DividendPayoutRatio(25000000, 100000000)
	if err != nil || got != 0.25 {
		t.Fatalf("got %v, %v; want 0.25", got, err)
	}
}

func TestDividendPayoutRatio_table(t *testing.T) {
	cases := []struct {
		name    string
		div, ni int64
		want    float64
	}{
		{"no dividend", 0, 5, 0},
		{"all paid out", 7, 7, 1},
		{"more than earnings", 3, 2, 1.5},
		{"max dividend", maxExactCents, 1, maxExactCents},
		{"max income", 1, maxExactCents, 1.0 / maxExactCents},
	}
	for _, c := range cases {
		got, err := DividendPayoutRatio(c.div, c.ni)
		if err != nil || got != c.want {
			t.Errorf("%s: got %v, %v; want %v", c.name, got, err, c.want)
		}
	}
}

func TestDividendPayoutRatio_errors(t *testing.T) {
	cases := []struct {
		name    string
		div, ni int64
		want    error
	}{
		{"zero income", 1, 0, ErrZeroDenominator},
		{"loss", 1, -1, ErrNegative},
		{"negative dividends", -1, 1, ErrNegative},
		{"dividends above 2^53", maxExactCents + 1, 1, ErrOutOfRange},
		{"income above 2^53", 1, maxExactCents + 1, ErrOutOfRange},
		{"income min int64", 1, math.MinInt64, ErrOutOfRange},
	}
	for _, c := range cases {
		if _, err := DividendPayoutRatio(c.div, c.ni); !errors.Is(err, c.want) {
			t.Errorf("%s: err = %v, want %v", c.name, err, c.want)
		}
	}
}

func TestDividendPayoutRatio_oracle(t *testing.T) {
	checkRatioOracle(t, DividendPayoutRatio, 0, 1)
}

func TestDividendPayoutRatio_allocs(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkFloat, _ = DividendPayoutRatio(25000000, 100000000) }); n != 0 {
		t.Fatalf("allocs = %v, want 0", n)
	}
}

func BenchmarkDividendPayoutRatio(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = DividendPayoutRatio(25000000, 100000000)
	}
}

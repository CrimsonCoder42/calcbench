package loan

import (
	"errors"
	"testing"
)

// Hand-checked from the OCC definition (NOI ÷ annual debt service): $125,000 / $100,000 = 1.25 (exact
// in binary).
func TestDebtServiceCoverageRatio_example(t *testing.T) {
	got, err := DebtServiceCoverageRatio(12_500_000, 10_000_000)
	if err != nil || got != 1.25 {
		t.Fatalf("got %v, %v; want 1.25", got, err)
	}
}

func TestDebtServiceCoverageRatio_cases(t *testing.T) {
	cases := []struct {
		name    string
		noi, ds int64
		want    float64
	}{
		{"negativeNOI", -5_000_000, 10_000_000, -0.5},
		{"zeroNOI", 0, 100, 0},
		{"breakEven", 100, 100, 1},
		{"maxExact", maxExactCents, 1, maxExactCents},
		{"minExact", -maxExactCents, maxExactCents, -1},
	}
	for _, c := range cases {
		got, err := DebtServiceCoverageRatio(c.noi, c.ds)
		if err != nil || got != c.want {
			t.Errorf("%s: got %v, %v; want %v", c.name, got, err, c.want)
		}
	}
}

func TestDebtServiceCoverageRatio_errors(t *testing.T) {
	cases := []struct {
		name    string
		noi, ds int64
		want    error
	}{
		{"noiAboveMax", maxExactCents + 1, 100, ErrOutOfRange},
		{"noiBelowMin", -maxExactCents - 1, 100, ErrOutOfRange},
		{"dsAboveMax", 100, maxExactCents + 1, ErrOutOfRange},
		{"dsNegative", 100, -1, ErrNegative},
		{"dsZero", 100, 0, ErrZeroDenominator},
	}
	for _, c := range cases {
		if _, err := DebtServiceCoverageRatio(c.noi, c.ds); !errors.Is(err, c.want) {
			t.Errorf("%s: err = %v, want %v", c.name, err, c.want)
		}
	}
}

func TestDebtServiceCoverageRatio_noAlloc(t *testing.T) {
	if a := testing.AllocsPerRun(100, func() { sinkFloat, _ = DebtServiceCoverageRatio(12_500_000, 10_000_000) }); a != 0 {
		t.Fatalf("allocs = %v, want 0", a)
	}
}

func BenchmarkDebtServiceCoverageRatio(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = DebtServiceCoverageRatio(12_500_000, 10_000_000)
	}
}

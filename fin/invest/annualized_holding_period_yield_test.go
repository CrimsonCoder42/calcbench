package invest

import (
	"math"
	"testing"
)

// OpenStax §15.1 annualizes a three-month return of 14.62 % with m = 4 to 72.60 %; 365/4 days is not
// whole, so the source's formula is checked at m = 5 (73 days), hand-checked as 1.1462⁵ − 1 =
// 0.9783449586159608 (double precision; compared to 1e-14), and at m = 1 (365 days).
func TestAnnualizedHoldingPeriodYield_handChecked(t *testing.T) {
	got, err := AnnualizedHoldingPeriodYield(0.1462, 73)
	if err != nil || !near(got, 0.9783449586159608, 1e-14) {
		t.Errorf("73 days: %v, %v; want 0.978345", got, err)
	}
	got, err = AnnualizedHoldingPeriodYield(0.1462, 365)
	if err != nil || !near(got, 0.1462, 1e-15) {
		t.Errorf("365 days: %v, %v; want 0.1462", got, err)
	}
	// The source's own m = 4 example through the same expression: (1.1462)⁴ − 1 = 0.7260 (±5e-5).
	if v := math.Expm1(4 * math.Log1p(0.1462)); !near(v, 0.7260, 5e-5) {
		t.Errorf("m = 4: %v, want 0.7260", v)
	}
}

func TestAnnualizedHoldingPeriodYield_cases(t *testing.T) {
	got, err := AnnualizedHoldingPeriodYield(0, 30)
	if err != nil || got != 0 {
		t.Errorf("zero: %v, %v", got, err)
	}
	// Small returns keep their digits: (1 + 1e-12)^365 − 1 ≈ 3.65e-10 (next term ≈ 6.7e-20).
	got, err = AnnualizedHoldingPeriodYield(1e-12, 1)
	if err != nil || !near(got, 3.65e-10, 1e-18) {
		t.Errorf("tiny: %v, %v", got, err)
	}
	got, err = AnnualizedHoldingPeriodYield(-0.5, 365)
	if err != nil || !near(got, -0.5, 1e-15) {
		t.Errorf("loss: %v, %v", got, err)
	}
	// Longer than a year: 730 days, 21 % → √1.21 − 1 = 0.1.
	got, err = AnnualizedHoldingPeriodYield(0.21, 730)
	if err != nil || !near(got, 0.1, 1e-15) {
		t.Errorf("two years: %v, %v", got, err)
	}
}

func TestAnnualizedHoldingPeriodYield_errors(t *testing.T) {
	for _, bad := range []float64{math.NaN(), math.Inf(1), math.Inf(-1), -1, -2} {
		_, err := AnnualizedHoldingPeriodYield(bad, 30)
		wantErr(t, "bad HPY", err, ErrRate)
	}
	_, err := AnnualizedHoldingPeriodYield(0.1, 0)
	wantErr(t, "days 0", err, ErrPeriods)
	_, err = AnnualizedHoldingPeriodYield(1e300, 1)
	wantErr(t, "overflow", err, ErrOverflow)
}

func TestAnnualizedHoldingPeriodYield_noAlloc(t *testing.T) {
	if a := testing.AllocsPerRun(100, func() { sinkFloat, sinkErr = AnnualizedHoldingPeriodYield(0.1462, 73) }); a != 0 {
		t.Errorf("allocs = %v, want 0", a)
	}
}

func BenchmarkAnnualizedHoldingPeriodYield(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, sinkErr = AnnualizedHoldingPeriodYield(0.1462, 73)
	}
}

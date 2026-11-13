package invest

import (
	"math"
	"testing"
)

// OpenStax Principles of Finance §17.1–17.3, Bluebonnet Industries: equity 1,000,000 × $15, debt
// 0.97 × $5,000,000 = $4,850,000, rd 6.312 %, T 21 % (after-tax 4.986 %), re 13.4 % (CAPM) → 11.33 %.
// The source rounds the weights (24 %/75.6 %) and the after-tax cost (4.99 %), so its figure is
// matched to ±5e-4; the exact value is 0.1134430368 (±1e-10).
func TestWeightedAverageCostOfCapital_openStax(t *testing.T) {
	got, err := WeightedAverageCostOfCapital(1_500_000_000, 485_000_000, 0.134, 0.06312, 0.21)
	if err != nil || !near(got, 0.1133, 5e-4) || !near(got, 0.1134430368, 1e-10) {
		t.Errorf("WACC = %v, %v; want 0.1134430368", got, err)
	}
}

func TestWeightedAverageCostOfCapital_cases(t *testing.T) {
	for _, c := range []struct {
		name      string
		e, d      int64
		re, rd, t float64
		want      float64
	}{
		{"all equity", 1000, 0, 0.12, 0.05, 0.3, 0.12},
		{"all debt", 0, 1000, 0.12, 0.05, 0.2, 0.04},
		{"half and half, no tax", 500, 500, 0.10, 0.06, 0, 0.08},
		{"tax 100 %", 500, 500, 0.10, 0.06, 1, 0.05},
		{"extremes", maxCents, maxCents, 0.1, 0.1, 0, 0.1},
	} {
		got, err := WeightedAverageCostOfCapital(c.e, c.d, c.re, c.rd, c.t)
		if err != nil || !near(got, c.want, 1e-15) {
			t.Errorf("%s: %v, %v; want %v", c.name, got, err, c.want)
		}
	}
}

func TestWeightedAverageCostOfCapital_errors(t *testing.T) {
	_, err := WeightedAverageCostOfCapital(0, 0, 0.1, 0.05, 0.2)
	wantErr(t, "no capital", err, ErrZeroDenominator)
	_, err = WeightedAverageCostOfCapital(-1, 100, 0.1, 0.05, 0.2)
	wantErr(t, "equity < 0", err, ErrNegative)
	_, err = WeightedAverageCostOfCapital(100, -1, 0.1, 0.05, 0.2)
	wantErr(t, "debt < 0", err, ErrNegative)
	_, err = WeightedAverageCostOfCapital(maxCents+1, 100, 0.1, 0.05, 0.2)
	wantErr(t, "equity > 2^53", err, ErrOutOfRange)
	for _, bad := range []float64{math.NaN(), math.Inf(1), -1} {
		_, err = WeightedAverageCostOfCapital(100, 100, bad, 0.05, 0.2)
		wantErr(t, "bad re", err, ErrRate)
		_, err = WeightedAverageCostOfCapital(100, 100, 0.1, bad, 0.2)
		wantErr(t, "bad rd", err, ErrRate)
	}
	_, err = WeightedAverageCostOfCapital(1, 1000000, 1e6, 1e308, 0) // 1e6·1e308 is beyond float64
	wantErr(t, "overflow", err, ErrOverflow)
	_, err = WeightedAverageCostOfCapital(100, 100, 0.1, 0.05, math.NaN())
	wantErr(t, "tax NaN", err, ErrNotFinite)
	_, err = WeightedAverageCostOfCapital(100, 100, 0.1, 0.05, -0.1)
	wantErr(t, "tax < 0", err, ErrDomain)
	_, err = WeightedAverageCostOfCapital(100, 100, 0.1, 0.05, 1.1)
	wantErr(t, "tax > 1", err, ErrDomain)
}

func TestWeightedAverageCostOfCapital_noAlloc(t *testing.T) {
	if a := testing.AllocsPerRun(100, func() {
		sinkFloat, sinkErr = WeightedAverageCostOfCapital(1_500_000_000, 485_000_000, 0.134, 0.06312, 0.21)
	}); a != 0 {
		t.Errorf("allocs = %v, want 0", a)
	}
}

func BenchmarkWeightedAverageCostOfCapital(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, sinkErr = WeightedAverageCostOfCapital(1_500_000_000, 485_000_000, 0.134, 0.06312, 0.21)
	}
}

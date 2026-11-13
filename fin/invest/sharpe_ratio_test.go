package invest

import (
	"math"
	"testing"
)

// OpenStax Principles of Finance §15.4, McKinley Investment Management: 20 % return, 3 % risk-free,
// σ 26 % → 17 %/26 % = 0.65 (two decimals: ±5e-3); exact 0.6538461538.
func TestSharpeRatio_openStax(t *testing.T) {
	got, err := SharpeRatio(0.20, 0.03, 0.26)
	if err != nil || !near(got, 0.65, 5e-3) || !near(got, 0.17/0.26, 1e-15) {
		t.Errorf("SharpeRatio = %v, %v; want 0.6538", got, err)
	}
}

func TestSharpeRatio_cases(t *testing.T) {
	for _, c := range []struct{ rp, rf, sd, want float64 }{
		{0.03, 0.03, 0.1, 0},
		{0.01, 0.03, 0.1, -0.2},
		{0.13, 0.03, 0.1, 1},
	} {
		got, err := SharpeRatio(c.rp, c.rf, c.sd)
		if err != nil || !near(got, c.want, 1e-15) {
			t.Errorf("SharpeRatio(%v, %v, %v) = %v, %v; want %v", c.rp, c.rf, c.sd, got, err, c.want)
		}
	}
}

func TestSharpeRatio_errors(t *testing.T) {
	_, err := SharpeRatio(0.1, 0.03, 0)
	wantErr(t, "sd 0", err, ErrZeroDenominator)
	_, err = SharpeRatio(0.1, 0.03, -0.1)
	wantErr(t, "sd < 0", err, ErrNegative)
	for _, bad := range []float64{math.NaN(), math.Inf(1)} {
		_, err = SharpeRatio(0.1, 0.03, bad)
		wantErr(t, "bad sd", err, ErrNotFinite)
	}
	for _, bad := range []float64{math.NaN(), math.Inf(-1), -1} {
		_, err = SharpeRatio(bad, 0.03, 0.1)
		wantErr(t, "bad rp", err, ErrRate)
		_, err = SharpeRatio(0.1, bad, 0.1)
		wantErr(t, "bad rf", err, ErrRate)
	}
	_, err = SharpeRatio(1e300, 0, 1e-300)
	wantErr(t, "overflow", err, ErrOverflow)
}

func TestSharpeRatio_noAlloc(t *testing.T) {
	if a := testing.AllocsPerRun(100, func() { sinkFloat, sinkErr = SharpeRatio(0.20, 0.03, 0.26) }); a != 0 {
		t.Errorf("allocs = %v, want 0", a)
	}
}

func BenchmarkSharpeRatio(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, sinkErr = SharpeRatio(0.20, 0.03, 0.26)
	}
}

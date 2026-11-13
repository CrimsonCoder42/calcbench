package invest

import (
	"math"
	"testing"
)

// OpenStax Principles of Finance §15.3, Delta Air Lines: Rf 3.36 %, Rm 11.64 %, β 1.39 → 14.87 %
// (two decimals of a percent: ±5e-5); exact arithmetic gives 0.148692.
func TestCAPMExpectedReturn_openStax(t *testing.T) {
	got, err := CAPMExpectedReturn(0.0336, 1.39, 0.1164)
	if err != nil || !near(got, 0.1487, 5e-5) || !near(got, 0.148692, 1e-15) {
		t.Errorf("CAPMExpectedReturn = %v, %v; want 0.148692", got, err)
	}
}

// OpenStax §17.3 IBM: Rf 3 %, β 0.97, market risk premium 5 % → 7.85 %.
func TestCAPMExpectedReturn_openStaxIBM(t *testing.T) {
	got, err := CAPMExpectedReturn(0.03, 0.97, 0.08)
	if err != nil || !near(got, 0.0785, 1e-15) {
		t.Errorf("IBM: %v, %v; want 0.0785", got, err)
	}
}

func TestCAPMExpectedReturn_cases(t *testing.T) {
	for _, c := range []struct {
		name         string
		rf, beta, rm float64
		want         float64
	}{
		{"beta 0", 0.03, 0, 0.10, 0.03},
		{"beta 1", 0.03, 1, 0.10, 0.10},
		{"negative beta", 0.03, -1, 0.10, -0.04},
		{"all zero", 0, 0, 0, 0},
	} {
		got, err := CAPMExpectedReturn(c.rf, c.beta, c.rm)
		if err != nil || !near(got, c.want, 1e-15) {
			t.Errorf("%s: %v, %v; want %v", c.name, got, err, c.want)
		}
	}
}

func TestCAPMExpectedReturn_errors(t *testing.T) {
	for _, bad := range []float64{math.NaN(), math.Inf(1), math.Inf(-1), -1} {
		_, err := CAPMExpectedReturn(bad, 1, 0.1)
		wantErr(t, "bad rf", err, ErrRate)
		_, err = CAPMExpectedReturn(0.03, 1, bad)
		wantErr(t, "bad rm", err, ErrRate)
	}
	for _, bad := range []float64{math.NaN(), math.Inf(1)} {
		_, err := CAPMExpectedReturn(0.03, bad, 0.1)
		wantErr(t, "bad beta", err, ErrNotFinite)
	}
	_, err := CAPMExpectedReturn(0, math.MaxFloat64, 1e300)
	wantErr(t, "overflow", err, ErrOverflow)
}

func TestCAPMExpectedReturn_noAlloc(t *testing.T) {
	if a := testing.AllocsPerRun(100, func() { sinkFloat, sinkErr = CAPMExpectedReturn(0.0336, 1.39, 0.1164) }); a != 0 {
		t.Errorf("allocs = %v, want 0", a)
	}
}

func BenchmarkCAPMExpectedReturn(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, sinkErr = CAPMExpectedReturn(0.0336, 1.39, 0.1164)
	}
}

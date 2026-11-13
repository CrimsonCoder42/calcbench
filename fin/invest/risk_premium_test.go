package invest

import (
	"math"
	"testing"
)

// OpenStax Principles of Finance §15.3: 11.64 % − 3.36 % = 8.28 %.
func TestRiskPremium_openStax(t *testing.T) {
	got, err := RiskPremium(0.1164, 0.0336)
	if err != nil || !near(got, 0.0828, 1e-15) {
		t.Errorf("RiskPremium = %v, %v; want 0.0828", got, err)
	}
}

func TestRiskPremium_cases(t *testing.T) {
	for _, c := range []struct{ e, rf, want float64 }{{0.05, 0.05, 0}, {0.02, 0.05, -0.03}, {0, 0, 0}} {
		got, err := RiskPremium(c.e, c.rf)
		if err != nil || !near(got, c.want, 1e-15) {
			t.Errorf("RiskPremium(%v, %v) = %v, %v; want %v", c.e, c.rf, got, err, c.want)
		}
	}
}

func TestRiskPremium_errors(t *testing.T) {
	for _, bad := range []float64{math.NaN(), math.Inf(1), math.Inf(-1), -1} {
		_, err := RiskPremium(bad, 0.03)
		wantErr(t, "bad expected", err, ErrRate)
		_, err = RiskPremium(0.1, bad)
		wantErr(t, "bad rf", err, ErrRate)
	}
}

func TestRiskPremium_noAlloc(t *testing.T) {
	if a := testing.AllocsPerRun(100, func() { sinkFloat, sinkErr = RiskPremium(0.1164, 0.0336) }); a != 0 {
		t.Errorf("allocs = %v, want 0", a)
	}
}

func BenchmarkRiskPremium(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, sinkErr = RiskPremium(0.1164, 0.0336)
	}
}

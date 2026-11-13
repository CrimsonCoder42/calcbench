package invest

import (
	"math"
	"testing"
)

// Confidence 0.95 and 0.99 are stated parameter values; z = Φ⁻¹(c) from NIST e-Handbook §1.3.6.7.1
// (1.645 and 2.326 to three decimals). Tolerance 1e-12: Φ⁻¹ is accurate to ~1e-16, 1 − c rounds.
func TestParametricVaR_handWorked(t *testing.T) {
	got, err := ParametricVaR(0, 0.02, 0.95)
	if err != nil || !near(got, 0.02*1.6448536269514722, 1e-12) || !near(got, 0.0329, 1e-4) {
		t.Errorf("95 %%: %v, %v; want 0.0328970725", got, err)
	}
	got, err = ParametricVaR(0.001, 0.02, 0.99)
	if err != nil || !near(got, 0.02*2.3263478740408408-0.001, 1e-12) {
		t.Errorf("99 %% with drift: %v, %v; want %v", got, err, 0.02*2.3263478740408408-0.001)
	}
}

func TestParametricVaR_cases(t *testing.T) {
	got, err := ParametricVaR(0.01, 0, 0.95)
	if err != nil || got != -0.01 {
		t.Errorf("σ = 0: %v, %v; want -0.01", got, err)
	}
	got, err = ParametricVaR(0.01, 0.02, 0.5)
	if err != nil || !near(got, -0.01, 1e-15) {
		t.Errorf("c = 0.5: %v, %v; want -μ", got, err)
	}
}

func TestParametricVaR_errors(t *testing.T) {
	_, err := ParametricVaR(math.NaN(), 0.02, 0.95)
	wantErr(t, "NaN mean", err, ErrNotFinite)
	_, err = ParametricVaR(0, math.Inf(1), 0.95)
	wantErr(t, "Inf sd", err, ErrNotFinite)
	_, err = ParametricVaR(0, -0.02, 0.95)
	wantErr(t, "negative sd", err, ErrNegative)
	for _, bad := range []float64{0, 1, 2, math.NaN()} {
		_, err = ParametricVaR(0, 0.02, bad)
		wantErr(t, "bad confidence", err, ErrDomain)
	}
	_, err = ParametricVaR(0, 1e308, 0.999)
	wantErr(t, "overflow", err, ErrOverflow)
}

func TestParametricVaR_noAlloc(t *testing.T) {
	if a := testing.AllocsPerRun(100, func() { sinkFloat, sinkErr = ParametricVaR(0, 0.02, 0.95) }); a != 0 {
		t.Errorf("allocs = %v, want 0", a)
	}
}

func BenchmarkParametricVaR(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, sinkErr = ParametricVaR(0, 0.02, 0.95)
	}
}

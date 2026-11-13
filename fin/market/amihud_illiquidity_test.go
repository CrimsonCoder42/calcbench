package market

import (
	"math"
	"testing"
)

// Amihud 2002 eq. (1); hand-checked: |0.02|/1e6 = 2e-8, |−0.01|/5e5 = 2e-8, mean 2e-8. Tolerance 1e-22
// (a few ulps of 2e-8).
func TestAmihudIlliquidity_handChecked(t *testing.T) {
	got, err := AmihudIlliquidity([]float64{0.02, -0.01}, []float64{1e6, 5e5})
	if err != nil || !near(got, 2e-8, 1e-22) {
		t.Errorf("AmihudIlliquidity = %v, %v; want 2e-8", got, err)
	}
}

func TestAmihudIlliquidity_cases(t *testing.T) {
	got, err := AmihudIlliquidity([]float64{0}, []float64{1})
	if err != nil || got != 0 {
		t.Errorf("zero return: %v, %v", got, err)
	}
	got, err = AmihudIlliquidity([]float64{0.5, -0.5, 0.25}, []float64{1, 1, 1})
	if err != nil || !near(got, 1.25/3, 1e-16) {
		t.Errorf("three days: %v, %v; want 1.25/3", got, err)
	}
}

func TestAmihudIlliquidity_errors(t *testing.T) {
	_, err := AmihudIlliquidity(nil, nil)
	wantErr(t, "empty", err, ErrEmpty)
	_, err = AmihudIlliquidity([]float64{1}, []float64{1, 2})
	wantErr(t, "lengths", err, ErrLengthMismatch)
	_, err = AmihudIlliquidity([]float64{math.NaN()}, []float64{1})
	wantErr(t, "NaN return", err, ErrNotFinite)
	_, err = AmihudIlliquidity([]float64{0.1}, []float64{-1})
	wantErr(t, "negative volume", err, ErrNegative)
	_, err = AmihudIlliquidity([]float64{0.1}, []float64{0})
	wantErr(t, "zero volume", err, ErrDomain)
	_, err = AmihudIlliquidity([]float64{1}, []float64{1e-320})
	wantErr(t, "overflow", err, ErrOverflow)
}

func BenchmarkAmihudIlliquidity(b *testing.B) {
	r, v := make([]float64, 250), make([]float64, 250)
	for i := range r {
		r[i], v[i] = float64(i%11-5)*0.001, 1e6+float64(i)*1e3
	}
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, sinkErr = AmihudIlliquidity(r, v)
	}
}

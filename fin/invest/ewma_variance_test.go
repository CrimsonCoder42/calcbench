package invest

import (
	"fmt"
	"math"
	"testing"
)

// λ = 0.94 is a stated parameter value: the RiskMetrics Technical Document (1996) daily decay factor.
func TestEWMAVariance_handWorked(t *testing.T) {
	got, err := EWMAVariance([]float64{0.01, -0.02}, 0.94, 0.0001)
	if err != nil || !near(got, 0.000118, 1e-17) {
		t.Errorf("EWMAVariance = %v, %v; want 0.000118", got, err)
	}
}

// Oracle: the unrolled weighted sum λⁿσ²₀ + (1 − λ)Σ λⁿ⁻¹⁻ᵗ rₜ².
func TestEWMAVariance_oracle(t *testing.T) {
	r := genReturns(200)
	const lambda, v0 = 0.97, 0.0004
	for n := 0; n <= len(r); n += 17 {
		want := math.Pow(lambda, float64(n)) * v0
		for i, x := range r[:n] {
			want += (1 - lambda) * math.Pow(lambda, float64(n-1-i)) * x * x
		}
		got, err := EWMAVariance(r[:n], lambda, v0)
		if err != nil || !near(got, want, 1e-15) {
			t.Errorf("n=%d: %v, %v; want %v", n, got, err, want)
		}
	}
}

func TestEWMAVariance_empty(t *testing.T) {
	got, err := EWMAVariance(nil, 0.94, 0.0003)
	if err != nil || got != 0.0003 {
		t.Errorf("empty: %v, %v; want the seed 0.0003", got, err)
	}
}

func TestEWMAVariance_errors(t *testing.T) {
	_, err := EWMAVariance([]float64{math.NaN()}, 0.94, 0)
	wantErr(t, "NaN return", err, ErrNotFinite)
	for _, bad := range []float64{0, 1, -0.5, 1.5, math.NaN(), math.Inf(1)} {
		_, err = EWMAVariance([]float64{0.01}, bad, 0)
		wantErr(t, "bad lambda", err, ErrDomain)
	}
	_, err = EWMAVariance([]float64{0.01}, 0.94, math.Inf(1))
	wantErr(t, "Inf seed", err, ErrNotFinite)
	_, err = EWMAVariance([]float64{0.01}, 0.94, -1e-6)
	wantErr(t, "negative seed", err, ErrNegative)
	_, err = EWMAVariance([]float64{1e200}, 0.94, 0)
	wantErr(t, "overflow", err, ErrOverflow)
}

func TestEWMAVariance_noAlloc(t *testing.T) {
	r := []float64{0.01, -0.02}
	if a := testing.AllocsPerRun(100, func() { sinkFloat, sinkErr = EWMAVariance(r, 0.94, 0.0001) }); a != 0 {
		t.Errorf("allocs = %v, want 0", a)
	}
}

func BenchmarkEWMAVariance(b *testing.B) {
	for _, n := range benchSizes {
		r := genReturns(n)
		b.Run(fmt.Sprint(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				sinkFloat, sinkErr = EWMAVariance(r, 0.94, 0.0001)
			}
		})
	}
}

package prob

import (
	"math"
	"testing"
)

// NIST §1.3.6.6.9 (mean e^(σ²/2), SD √(e^(σ²)(e^(σ²) − 1)) for m = 1): μ = 0, σ = 1 → (√e, (e − 1)e) =
// (1.6487212707001282, 4.670774270471604). Hand-checked; tolerance 4 ulps.
func TestLognormalMeanVariance_sourceExample(t *testing.T) {
	m, v, err := LognormalMeanVariance(0, 1)
	if err != nil || !within(m, 1.6487212707001282, 1e-15) || !within(v, 4.670774270471604, 1e-15) {
		t.Errorf("LognormalMeanVariance(0, 1) = (%v, %v, %v), want (1.6487212707001282, 4.670774270471604, nil)", m, v, err)
	}
}

func TestLognormalMeanVariance_cases(t *testing.T) {
	tests := []struct {
		name           string
		mu, sigma      float64
		mean, variance float64
	}{
		{"small sigma keeps the variance", 0, 1e-10, math.Exp(0.5e-20), 1e-20},
		{"scaled by e^μ", 2, 1, math.Exp(2.5), math.Expm1(1) * math.Exp(5)},
		{"large sigma", -200, 20, 1, math.Expm1(400)},
		{"variance from logs (expm1(σ²) overflows alone)", -800, 30, math.Exp(-350), math.Exp(200)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, v, err := LognormalMeanVariance(tt.mu, tt.sigma)
			if err != nil || !within(m, tt.mean, 1e-14) || !within(v, tt.variance, 1e-12) {
				t.Errorf("LognormalMeanVariance(%v, %v) = (%v, %v, %v), want (%v, %v)", tt.mu, tt.sigma, m, v, err, tt.mean, tt.variance)
			}
		})
	}
}

func TestLognormalMeanVariance_errors(t *testing.T) {
	_, _, err := LognormalMeanVariance(0, 0)
	requireErr(t, "sigma = 0", err, ErrOutOfRange)
	_, _, err = LognormalMeanVariance(nan, 1)
	requireErr(t, "NaN", err, ErrNotFinite)
	_, _, err = LognormalMeanVariance(710, 1)
	requireErr(t, "mean overflows", err, ErrOverflow)
	_, _, err = LognormalMeanVariance(0, 30)
	requireErr(t, "variance overflows", err, ErrOverflow)
}

func TestLognormalMeanVariance_noAlloc(t *testing.T) {
	requireNoAlloc(t, "LognormalMeanVariance", func() { sinkFloat, _, _ = LognormalMeanVariance(0, 1) })
}

func BenchmarkLognormalMeanVariance(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _, _ = LognormalMeanVariance(0, 1)
	}
}

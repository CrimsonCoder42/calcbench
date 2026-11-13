package prob

import (
	"math"
	"testing"
)

// NIST §1.3.6.6.1 (f(x) = exp(−(x−μ)²/(2σ²))/(σ√(2π))): x = 0, μ = 0, σ = 1 → 1/√(2π) =
// 0.3989422804014327. Further values from 150-digit decimal arithmetic for this test: φ(1) =
// 0.24197072451914334, φ(7.3) = 1.0693837871541648e-12 (at the float64 nearest 7.3), φ(30) = 1.4736461348785476e-196.
// Tolerance 2e-15 (the tail path splits z so that z² is exact).
func TestNormalPDF_sourceExample(t *testing.T) {
	for _, c := range []struct{ x, want float64 }{
		{0, 0.3989422804014327}, {1, 0.24197072451914334}, {-7.3, 1.0693837871541648e-12}, {30, 1.4736461348785476e-196},
	} {
		if got, err := NormalPDF(c.x, 0, 1); err != nil || !within(got, c.want, 2e-15) {
			t.Errorf("NormalPDF(%v, 0, 1) = (%v, %v), want (%v, nil)", c.x, got, err, c.want)
		}
	}
}

func TestNormalPDF_cases(t *testing.T) {
	tests := []struct {
		name              string
		x, mean, sd, want float64
	}{
		{"scaled", 10, 10, 2, 0.3989422804014327 / 2},
		{"far tail underflows", 100, 0, 1, 0},
		{"difference overflows to 0", math.MaxFloat64, -math.MaxFloat64, 1, 0},
		{"huge sd", 0, 0, 1e300, 0.3989422804014327e-300},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := NormalPDF(tt.x, tt.mean, tt.sd); err != nil || !within(got, tt.want, 2e-15) {
				t.Errorf("NormalPDF(%v, %v, %v) = (%v, %v), want %v", tt.x, tt.mean, tt.sd, got, err, tt.want)
			}
		})
	}
}

func TestNormalPDF_errors(t *testing.T) {
	_, err := NormalPDF(0, 0, 0)
	requireErr(t, "sd = 0", err, ErrOutOfRange)
	_, err = NormalPDF(0, 0, -1)
	requireErr(t, "sd < 0", err, ErrOutOfRange)
	_, err = NormalPDF(0, 0, 5e-324)
	requireErr(t, "density overflows", err, ErrOverflow)
	_, err = NormalPDF(nan, 0, 1)
	requireErr(t, "x NaN", err, ErrNotFinite)
	_, err = NormalPDF(0, inf, 1)
	requireErr(t, "mean Inf", err, ErrNotFinite)
}

// Oracle: exp(−z²/2)/√(2π) for |z| < 5, where rounding z² costs at most 13 ulps. Tolerance 4e-15.
func TestNormalPDF_oracle(t *testing.T) {
	rng := newTestRand()
	for trial := range 500 {
		z := (rng.Float64() - 0.5) * 9.9
		want := math.Exp(-0.5*z*z) / math.Sqrt(2*math.Pi)
		if got, err := NormalPDF(z, 0, 1); err != nil || !within(got, want, 4e-15) {
			t.Fatalf("trial %d: NormalPDF(%v, 0, 1) = (%v, %v), want %v", trial, z, got, err, want)
		}
	}
}

func TestNormalPDF_noAlloc(t *testing.T) {
	requireNoAlloc(t, "NormalPDF", func() { sinkFloat, _ = NormalPDF(1, 0, 1) })
}

func BenchmarkNormalPDF(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = NormalPDF(1, 0, 1)
	}
}

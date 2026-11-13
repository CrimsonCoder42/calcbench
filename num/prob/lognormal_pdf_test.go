package prob

import (
	"math"
	"testing"
)

// NIST §1.3.6.6.9 (f(x) = exp(−(ln x − μ)²/(2σ²))/(xσ√(2π)), μ = ln m): μ = 0, σ = 1, x = 1 → 1/√(2π) =
// 0.3989422804014327 (hand-checked); μ = 1, σ = 0.5, x = 2 → 0.33046456598348395 (150-digit decimal
// arithmetic for this test). Tolerance 4e-15.
func TestLognormalPDF_sourceExample(t *testing.T) {
	for _, c := range []struct{ x, mu, sigma, want float64 }{{1, 0, 1, 0.3989422804014327}, {2, 1, 0.5, 0.33046456598348395}} {
		if got, err := LognormalPDF(c.x, c.mu, c.sigma); err != nil || !within(got, c.want, 4e-15) {
			t.Errorf("LognormalPDF(%v, %v, %v) = (%v, %v), want (%v, nil)", c.x, c.mu, c.sigma, got, err, c.want)
		}
	}
}

func TestLognormalPDF_cases(t *testing.T) {
	tests := []struct {
		name               string
		x, mu, sigma, want float64
	}{
		{"x = 0", 0, 0, 1, 0},
		{"x < 0", -1, 0, 1, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LognormalPDF(tt.x, tt.mu, tt.sigma)
			if err != nil || !within(got, tt.want, 1e-12) {
				t.Errorf("LognormalPDF(%v, %v, %v) = (%v, %v), want %v", tt.x, tt.mu, tt.sigma, got, err, tt.want)
			}
		})
	}
}

// x·σ below the normal range takes the logarithmic path; the expected value is the same density written
// in logs. Tolerance 1e-12 (exp of an argument near −260 carries about 260 ulps).
func TestLognormalPDF_logPath(t *testing.T) {
	x, sigma := 1e-300, 1e-10
	mu := math.Log(x) - 30*sigma
	z := (math.Log(x) - mu) / sigma
	want := math.Exp(-0.5*z*z-math.Log(x)-math.Log(sigma)) / math.Sqrt(2*math.Pi)
	if got, err := LognormalPDF(x, mu, sigma); err != nil || !within(got, want, 1e-12) {
		t.Errorf("LognormalPDF(1e-300, …, 1e-10) = (%v, %v), want %v", got, err, want)
	}
}

func TestLognormalPDF_errors(t *testing.T) {
	_, err := LognormalPDF(1, 0, 0)
	requireErr(t, "sigma = 0", err, ErrOutOfRange)
	_, err = LognormalPDF(nan, 0, 1)
	requireErr(t, "x NaN", err, ErrNotFinite)
	_, err = LognormalPDF(1, inf, 1)
	requireErr(t, "mu Inf", err, ErrNotFinite)
	_, err = LognormalPDF(1e-300, math.Log(1e-300), 1e-10)
	requireErr(t, "density overflows", err, ErrOverflow)
}

// Oracle: the change of variables f_X(x) = φ((ln x − μ)/σ)/(xσ), with φ from NormalPDF. Tolerance 4e-15.
func TestLognormalPDF_oracle(t *testing.T) {
	rng := newTestRand()
	for trial := range 400 {
		x, mu, sigma := 0.01+10*rng.Float64(), rng.Float64()-0.5, 0.1+2*rng.Float64()
		phi, _ := NormalPDF((math.Log(x)-mu)/sigma, 0, 1)
		want := phi / (x * sigma)
		if got, err := LognormalPDF(x, mu, sigma); err != nil || !within(got, want, 4e-15) {
			t.Fatalf("trial %d: LognormalPDF(%v, %v, %v) = (%v, %v), want %v", trial, x, mu, sigma, got, err, want)
		}
	}
}

func TestLognormalPDF_noAlloc(t *testing.T) {
	requireNoAlloc(t, "LognormalPDF", func() { sinkFloat, _ = LognormalPDF(2, 1, 0.5) })
}

func BenchmarkLognormalPDF(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = LognormalPDF(2, 1, 0.5)
	}
}

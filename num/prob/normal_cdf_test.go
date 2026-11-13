package prob

import (
	"math"
	"testing"
)

// NIST §1.3.6.7.1: the area from 0 to 1.96 is 0.47500, so Φ(1.96) = 0.97500 to 5 decimals. Full
// precision from 150-digit decimal arithmetic for this test: Φ(1.96) = 0.97500210485177952,
// Φ(−1.96) = 0.024997895148220435, Φ(−5) = 2.8665157187919391e-07. Tolerance 1e-14 relative (erfc's
// argument z/√2 is rounded, which costs about 2z² ulps).
func TestNormalCDF_sourceExample(t *testing.T) {
	got, err := NormalCDF(1.96, 0, 1)
	if err != nil || math.Abs(got-0.5-0.47500) > 5e-6 {
		t.Errorf("NormalCDF(1.96) − 0.5 = %v (err %v), want 0.47500 (NIST table)", got-0.5, err)
	}
	for _, c := range []struct{ x, want float64 }{{1.96, 0.97500210485177952}, {-1.96, 0.024997895148220435}, {-5, 2.8665157187919391e-07}} {
		if got, err := NormalCDF(c.x, 0, 1); err != nil || !within(got, c.want, 1e-14) {
			t.Errorf("NormalCDF(%v, 0, 1) = (%v, %v), want (%v, nil)", c.x, got, err, c.want)
		}
	}
}

func TestNormalCDF_cases(t *testing.T) {
	tests := []struct {
		name              string
		x, mean, sd, want float64
	}{
		{"at the mean", 3, 3, 2, 0.5},
		{"far lower", -40, 0, 1, 0},
		{"far upper", 40, 0, 1, 1},
		{"difference overflows", -math.MaxFloat64, math.MaxFloat64, 1, 0},
		{"one sd, scaled", 12, 10, 2, 0.8413447460685429},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := NormalCDF(tt.x, tt.mean, tt.sd); err != nil || !within(got, tt.want, 1e-15) {
				t.Errorf("NormalCDF(%v, %v, %v) = (%v, %v), want %v", tt.x, tt.mean, tt.sd, got, err, tt.want)
			}
		})
	}
}

func TestNormalCDF_errors(t *testing.T) {
	_, err := NormalCDF(0, 0, 0)
	requireErr(t, "sd = 0", err, ErrOutOfRange)
	_, err = NormalCDF(nan, 0, 1)
	requireErr(t, "NaN", err, ErrNotFinite)
}

// Oracle: Φ(z) + Q(z) = 1 and Φ(−z) = Q(z) for random z. Tolerance 1e-15 for the sum; 2 ulps for the symmetry.
func TestNormalCDF_oracle(t *testing.T) {
	rng := newTestRand()
	for trial := range 500 {
		z := (rng.Float64() - 0.5) * 20
		p, _ := NormalCDF(z, 0, 1)
		q, _ := NormalSurvival(z, 0, 1)
		m, _ := NormalCDF(-z, 0, 1)
		if absDiff(p+q, 1) > 1e-15 || !within(m, q, 5e-16) {
			t.Fatalf("trial %d: z = %v: Φ = %v, Q = %v, Φ(−z) = %v", trial, z, p, q, m)
		}
	}
}

func TestNormalCDF_noAlloc(t *testing.T) {
	requireNoAlloc(t, "NormalCDF", func() { sinkFloat, _ = NormalCDF(1.96, 0, 1) })
}

func BenchmarkNormalCDF(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = NormalCDF(1.96, 0, 1)
	}
}

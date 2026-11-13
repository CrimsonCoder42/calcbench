package prob

import (
	"math"
	"testing"
)

// Ross §4.8 / Casella–Berger §3.2 (geometric, trials): p = 0.2, k = 3 → 1 − 0.8³ = 0.488. Hand-checked;
// tolerance 4 ulps.
func TestGeometricTrialsCDF_sourceExample(t *testing.T) {
	if got, err := GeometricTrialsCDF(3, 0.2); err != nil || !within(got, 0.488, 1e-15) {
		t.Errorf("GeometricTrialsCDF(3, 0.2) = (%v, %v), want (0.488, nil)", got, err)
	}
}

func TestGeometricTrialsCDF_cases(t *testing.T) {
	tests := []struct {
		name string
		k    int
		p    float64
		want float64
	}{
		{"k = 0", 0, 0.2, 0},
		{"k = 1", 1, 0.2, 0.2},
		{"p = 1", 1, 1, 1},
		{"tiny p keeps precision", 1, 1e-20, 1e-20},
		{"huge k", math.MaxInt, 0.5, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := GeometricTrialsCDF(tt.k, tt.p); err != nil || !within(got, tt.want, 1e-15) {
				t.Errorf("GeometricTrialsCDF(%d, %v) = (%v, %v), want %v", tt.k, tt.p, got, err, tt.want)
			}
		})
	}
}

func TestGeometricTrialsCDF_errors(t *testing.T) {
	_, err := GeometricTrialsCDF(1, 0)
	requireErr(t, "p = 0", err, ErrOutOfRange)
	_, err = GeometricTrialsCDF(1, inf)
	requireErr(t, "Inf", err, ErrNotFinite)
}

// Oracle: the CDF equals the compensated sum of GeometricTrialsPMF over 1…k. Tolerance 1e-13.
func TestGeometricTrialsCDF_oracle(t *testing.T) {
	rng := newTestRand()
	for trial := range 300 {
		p := 0.01 + 0.99*rng.Float64()
		k := 1 + rng.IntN(200)
		var s compensatedSum
		for i := 1; i <= k; i++ {
			v, _ := GeometricTrialsPMF(i, p)
			s.add(v)
		}
		if got, err := GeometricTrialsCDF(k, p); err != nil || !within(got, s.total(), 1e-13) {
			t.Fatalf("trial %d: GeometricTrialsCDF(%d, %v) = (%v, %v), want %v", trial, k, p, got, err, s.total())
		}
	}
}

func TestGeometricTrialsCDF_noAlloc(t *testing.T) {
	requireNoAlloc(t, "GeometricTrialsCDF", func() { sinkFloat, _ = GeometricTrialsCDF(3, 0.2) })
}

func BenchmarkGeometricTrialsCDF(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = GeometricTrialsCDF(3, 0.2)
	}
}

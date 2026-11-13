package prob

import (
	"math"
	"testing"
)

// Inverting the trials-convention CDF of Ross §4.8 / Casella–Berger §3.2: p = 0.2, u = 0.5 → 4, since
// F(3) = 0.488 < 0.5 ≤ F(4) = 0.5904. Hand-checked.
func TestGeometricTrialsQuantile_sourceExample(t *testing.T) {
	if got, err := GeometricTrialsQuantile(0.5, 0.2); err != nil || got != 4 {
		t.Errorf("GeometricTrialsQuantile(0.5, 0.2) = (%v, %v), want (4, nil)", got, err)
	}
}

func TestGeometricTrialsQuantile_cases(t *testing.T) {
	tests := []struct {
		name string
		u, p float64
		want int
	}{
		{"u = 0", 0, 0.2, 1},
		{"u = p", 0.2, 0.2, 1},
		{"just above p", 0.2000000001, 0.2, 2},
		{"p = 1", 0.99, 1, 1},
		{"p = 1, u = 1", 1, 1, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := GeometricTrialsQuantile(tt.u, tt.p); err != nil || got != tt.want {
				t.Errorf("GeometricTrialsQuantile(%v, %v) = (%v, %v), want %v", tt.u, tt.p, got, err, tt.want)
			}
		})
	}
}

func TestGeometricTrialsQuantile_errors(t *testing.T) {
	_, err := GeometricTrialsQuantile(1, 0.5)
	requireErr(t, "u = 1, p < 1", err, ErrOutOfRange)
	_, err = GeometricTrialsQuantile(0.5, 0)
	requireErr(t, "p = 0", err, ErrOutOfRange)
	_, err = GeometricTrialsQuantile(0.5, 1e-300)
	requireErr(t, "k beyond int", err, ErrOverflow)
	_, err = GeometricTrialsQuantile(nan, 0.5)
	requireErr(t, "u NaN", err, ErrNotFinite)
	_, err = GeometricTrialsQuantile(1.5, 0.5)
	requireErr(t, "u > 1", err, ErrProbability)
}

// Oracle: the defining property F(k−1) < u ≤ F(k) with this package's CDF, for random u and p.
func TestGeometricTrialsQuantile_oracle(t *testing.T) {
	rng := newTestRand()
	for trial := range 1000 {
		u, p := rng.Float64(), math.Pow(10, -6*rng.Float64())
		k, err := GeometricTrialsQuantile(u, p)
		if err != nil {
			t.Fatal(err)
		}
		hi, _ := GeometricTrialsCDF(k, p)
		lo, _ := GeometricTrialsCDF(k-1, p)
		if !(lo < u && u <= hi) {
			t.Fatalf("trial %d: GeometricTrialsQuantile(%v, %v) = %d, F(k−1) = %v, F(k) = %v", trial, u, p, k, lo, hi)
		}
	}
}

func TestGeometricTrialsQuantile_noAlloc(t *testing.T) {
	requireNoAlloc(t, "GeometricTrialsQuantile", func() { sinkInt, _ = GeometricTrialsQuantile(0.5, 0.2) })
}

func BenchmarkGeometricTrialsQuantile(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkInt, _ = GeometricTrialsQuantile(0.5, 0.2)
	}
}

package prob

import (
	"math"
	"testing"
)

// Ross §4.8 / Casella–Berger §3.2 (geometric, trials convention p(1−p)^(x−1)): p = 0.2, k = 3 →
// 0.8²·0.2 = 0.128. Hand-checked; tolerance 4 ulps (exp of a log1p product).
func TestGeometricTrialsPMF_sourceExample(t *testing.T) {
	if got, err := GeometricTrialsPMF(3, 0.2); err != nil || !within(got, 0.128, 1e-15) {
		t.Errorf("GeometricTrialsPMF(3, 0.2) = (%v, %v), want (0.128, nil)", got, err)
	}
}

func TestGeometricTrialsPMF_cases(t *testing.T) {
	tests := []struct {
		name string
		k    int
		p    float64
		want float64
	}{
		{"first trial", 1, 0.2, 0.2},
		{"k = 0 outside support", 0, 0.2, 0},
		{"negative k", -5, 0.2, 0},
		{"p = 1, k = 1", 1, 1, 1},
		{"p = 1, k = 2", 2, 1, 0},
		{"tiny p", 1_000_000, 1e-6, 1e-6 * math.Exp(999_999*math.Log1p(-1e-6))},
		{"huge k underflows", math.MaxInt, 0.5, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := GeometricTrialsPMF(tt.k, tt.p); err != nil || !within(got, tt.want, 1e-15) {
				t.Errorf("GeometricTrialsPMF(%d, %v) = (%v, %v), want %v", tt.k, tt.p, got, err, tt.want)
			}
		})
	}
}

func TestGeometricTrialsPMF_errors(t *testing.T) {
	_, err := GeometricTrialsPMF(1, 0)
	requireErr(t, "p = 0", err, ErrOutOfRange)
	_, err = GeometricTrialsPMF(1, nan)
	requireErr(t, "NaN", err, ErrNotFinite)
	_, err = GeometricTrialsPMF(1, 1.5)
	requireErr(t, "p > 1", err, ErrProbability)
}

// Oracle: math.Pow(1−p, k−1)·p for small k, where Pow is accurate. Tolerance 1e-13 relative (1 − p is
// rounded in the oracle, an error of k·2⁻⁵³).
func TestGeometricTrialsPMF_oracle(t *testing.T) {
	rng := newTestRand()
	for trial := range 400 {
		p := 0.01 + 0.99*rng.Float64()
		k := 1 + rng.IntN(100)
		want := math.Pow(1-p, float64(k-1)) * p
		if want < 1e-290 {
			continue
		}
		if got, err := GeometricTrialsPMF(k, p); err != nil || !within(got, want, 1e-13) {
			t.Fatalf("trial %d: GeometricTrialsPMF(%d, %v) = (%v, %v), want %v", trial, k, p, got, err, want)
		}
	}
}

func TestGeometricTrialsPMF_noAlloc(t *testing.T) {
	requireNoAlloc(t, "GeometricTrialsPMF", func() { sinkFloat, _ = GeometricTrialsPMF(3, 0.2) })
}

func BenchmarkGeometricTrialsPMF(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = GeometricTrialsPMF(3, 0.2)
	}
}

package prob

import "testing"

// JKK ch. 4: λ = 1: F(1) = 2/e ≈ 0.7358 < 0.9 ≤ F(2) = 2.5/e ≈ 0.9197 → 2. Hand-checked.
func TestPoissonQuantile_sourceExample(t *testing.T) {
	if got, err := PoissonQuantile(0.9, 1); err != nil || got != 2 {
		t.Errorf("PoissonQuantile(0.9, 1) = (%v, %v), want 2", got, err)
	}
}

func TestPoissonQuantile_leastK(t *testing.T) {
	for _, lambda := range []float64{0.01, 1, 37.5, 1e5} {
		for _, u := range []float64{1e-12, 0.05, 0.5, 0.95, 1 - 1e-12} {
			k, err := PoissonQuantile(u, lambda)
			hi, _ := PoissonCDF(k, lambda)
			lo := 0.0
			if k > 0 {
				lo, _ = PoissonCDF(k-1, lambda)
			}
			if err != nil || !(lo < u && u <= hi) {
				t.Errorf("λ %v u %v: k %v, F(k−1) %v, F(k) %v, err %v", lambda, u, k, lo, hi, err)
			}
		}
	}
}

func TestPoissonQuantile_errors(t *testing.T) {
	if got, err := PoissonQuantile(0, 3); err != nil || got != 0 {
		t.Errorf("u = 0: (%v, %v)", got, err)
	}
	if got, err := PoissonQuantile(1, 0); err != nil || got != 0 {
		t.Errorf("λ = 0: (%v, %v)", got, err)
	}
	_, err := PoissonQuantile(1, 3)
	requireErr(t, "u = 1", err, ErrOutOfRange)
	_, err = PoissonQuantile(0.5, -1)
	requireErr(t, "λ < 0", err, ErrOutOfRange)
	_, err = PoissonQuantile(2, 1)
	requireErr(t, "u > 1", err, ErrProbability)
	_, err = PoissonQuantile(0.5, inf)
	requireErr(t, "λ Inf", err, ErrNotFinite)
}

func BenchmarkPoissonQuantile(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkInt, _ = PoissonQuantile(0.9, 50)
	}
}

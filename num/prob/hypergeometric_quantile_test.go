package prob

import "testing"

// JKK ch. 6: N = 52, K = 4, n = 5: F(0) = C(48,5)/C(52,5) ≈ 0.6588, F(1) ≈ 0.9583, F(4) = 1.
func TestHypergeometricQuantile_sourceExample(t *testing.T) {
	cases := []struct {
		u    float64
		want int
	}{{0, 0}, {0.5, 0}, {0.65, 0}, {0.9, 1}, {0.96, 2}, {1, 4}}
	for _, c := range cases {
		if got, err := HypergeometricQuantile(c.u, 52, 4, 5); err != nil || got != c.want {
			t.Errorf("u %v: (%v, %v), want %v", c.u, got, err, c.want)
		}
	}
	// Definition: F(q) ≥ u > F(q − 1), and the support bottom when n + K > N.
	for _, u := range []float64{1e-9, 0.1, 0.37, 0.5, 0.83, 0.999999} {
		q, err := HypergeometricQuantile(u, 100, 60, 70)
		if err != nil || q < 30 {
			t.Fatalf("u %v: (%v, %v)", u, q, err)
		}
		c, _ := HypergeometricCDF(q, 100, 60, 70)
		cm, _ := HypergeometricCDF(q-1, 100, 60, 70)
		if c < u || cm >= u {
			t.Errorf("u %v: q %d, F(q) %v, F(q−1) %v", u, q, c, cm)
		}
	}
	if q, _ := HypergeometricQuantile(0, 100, 60, 70); q != 30 {
		t.Errorf("u = 0: %v, want 30", q)
	}
}

func TestHypergeometricQuantile_errors(t *testing.T) {
	_, err := HypergeometricQuantile(nan, 10, 3, 4)
	requireErr(t, "u NaN", err, ErrNotFinite)
	_, err = HypergeometricQuantile(1.5, 10, 3, 4)
	requireErr(t, "u > 1", err, ErrProbability)
	_, err = HypergeometricQuantile(0.5, 10, 11, 4)
	requireErr(t, "K > N", err, ErrOutOfRange)
}

func TestHypergeometricQuantile_noAlloc(t *testing.T) {
	requireNoAlloc(t, "HypergeometricQuantile", func() { sinkInt, _ = HypergeometricQuantile(0.7, 1000, 300, 200) })
}

func BenchmarkHypergeometricQuantile(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkInt, _ = HypergeometricQuantile(0.7, 1000, 300, 200)
	}
}

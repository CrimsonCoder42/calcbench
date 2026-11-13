package prob

import (
	"math"
	"testing"
)

// JKK 3rd ed. §11.2.20: F(1) = 1/ζ(2) = 6/π² ≈ 0.6079, F(2) = 1.25·6/π² ≈ 0.7599.
func TestZipfQuantile(t *testing.T) {
	if got, err := ZipfQuantile(0.7, 2); err != nil || got != 2 {
		t.Errorf("ZipfQuantile(0.7, 2) = (%v, %v), want 2", got, err)
	}
	if got, _ := ZipfQuantile(0, 2); got != 1 {
		t.Errorf("u = 0: %v, want 1", got)
	}
	if got, _ := ZipfQuantile(6/(math.Pi*math.Pi)*0.999, 2); got != 1 {
		t.Errorf("u just below F(1): %v, want 1", got)
	}
	// The defining property against ZipfCDF, both below and above the k = 1000 switch (a = 1.1, u = 0.9:
	// k ≈ 10¹⁰, the bisection branch).
	for _, a := range []float64{1.1, 1.5, 2, 3.5, 50} {
		for _, u := range []float64{0.01, 0.3, 0.5, 0.9, 0.99, 0.999, 0.999999} {
			if a == 1.1 && u > 0.9 {
				continue // the quantile exceeds 2^62 (k^(−0.1) ≈ 0.1·ζ(1.1)·(1 − u))
			}
			k, err := ZipfQuantile(u, a)
			if err != nil {
				t.Errorf("a = %v, u = %v: %v", a, u, err)
				continue
			}
			hi, _ := ZipfCDF(k, a)
			lo, _ := ZipfCDF(k-1, a)
			if !(hi >= u && lo < u) {
				t.Errorf("a = %v, u = %v: k = %d, F(k−1) = %v, F(k) = %v", a, u, k, lo, hi)
			}
		}
	}
	_, err := ZipfQuantile(1, 2)
	requireErr(t, "u = 1", err, ErrOutOfRange)
	_, err = ZipfQuantile(1.5, 2)
	requireErr(t, "u > 1", err, ErrProbability)
	_, err = ZipfQuantile(0.5, 1)
	requireErr(t, "a = 1", err, ErrOutOfRange)
	_, err = ZipfQuantile(math.NaN(), 2)
	requireErr(t, "u NaN", err, ErrNotFinite)
	_, err = ZipfQuantile(1-1e-15, 1.0001)
	requireErr(t, "beyond 2^62", err, ErrOutOfRange)
	requireNoAlloc(t, "ZipfQuantile", func() { sinkInt, _ = ZipfQuantile(0.9, 2.5) })
}

func BenchmarkZipfQuantile(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkInt, _ = ZipfQuantile(0.9, 2.5)
	}
}

package prob

import (
	"math"
	"testing"
)

// JKK 3rd ed. §7.1: P(X = k) = −p^k/(k ln(1 − p)); p = 0.5 gives 1/(2 ln 2) at k = 1, 1/(8 ln 2) at k = 2.
func TestLogSeries(t *testing.T) {
	if got, err := LogSeriesPMF(1, 0.5); err != nil || !within(got, 0.5/math.Ln2, 1e-15) {
		t.Errorf("PMF(1) = (%v, %v)", got, err)
	}
	if got, _ := LogSeriesPMF(2, 0.5); !within(got, 0.125/math.Ln2, 1e-15) {
		t.Errorf("PMF(2) = %v", got)
	}
	if got, _ := LogSeriesCDF(2, 0.5); !within(got, 0.625/math.Ln2, 1e-15) {
		t.Errorf("CDF(2) = %v", got)
	}
	if got, _ := LogSeriesPMF(0, 0.5); got != 0 {
		t.Errorf("PMF(0) = %v", got)
	}
	if got, _ := LogSeriesCDF(0, 0.5); got != 0 {
		t.Errorf("CDF(0) = %v", got)
	}
	// CDF agrees with the summed PMF, reaches 1, and the quantile inverts it.
	for _, p := range []float64{1e-6, 0.3, 0.9, 0.999} {
		sum := 0.0
		for k := 1; k <= 60; k++ {
			f, _ := LogSeriesPMF(k, p)
			sum += f
			c, _ := LogSeriesCDF(k, p)
			if !within(c, sum, 1e-13) {
				t.Errorf("p = %v, k = %d: CDF %v, Σ PMF %v", p, k, c, sum)
			}
		}
		for _, u := range []float64{0.01, 0.5, 0.9, 0.99} {
			k, err := LogSeriesQuantile(u, p)
			if err != nil {
				t.Fatalf("Q(%v, %v): %v", u, p, err)
			}
			hi, _ := LogSeriesCDF(k, p)
			lo, _ := LogSeriesCDF(k-1, p)
			if hi < u || lo >= u {
				t.Errorf("p = %v, u = %v: k = %d, F(k−1) = %v, F(k) = %v", p, u, k, lo, hi)
			}
		}
		// Moments against the series Σ k f(k) and Σ (k − μ)² f(k).
		m, v, err := LogSeriesMeanVariance(p)
		if err != nil {
			t.Fatal(err)
		}
		var m1, m2 float64
		for k := 1; k <= 200000; k++ {
			f, _ := LogSeriesPMF(k, p)
			if f == 0 {
				break
			}
			m1 += float64(k) * f
			m2 += (float64(k) - m) * (float64(k) - m) * f
		}
		if !within(m1, m, 1e-10) || !within(m2, v, 1e-9) {
			t.Errorf("p = %v: series (%v, %v), closed form (%v, %v)", p, m1, m2, m, v)
		}
	}
	if c, _ := LogSeriesCDF(math.MaxInt, 0.5); c != 1 {
		t.Errorf("CDF(MaxInt) = %v", c)
	}
	if k, _ := LogSeriesQuantile(0, 0.5); k != 1 {
		t.Errorf("Q(0) = %v", k)
	}
	if m, v, _ := LogSeriesMeanVariance(0.5); !within(m, 1/math.Ln2, 1e-15) || !within(v, 0.804021100772319, 1e-14) {
		t.Errorf("MV = (%v, %v)", m, v)
	}
	_, err := LogSeriesPMF(1, 0)
	requireErr(t, "p = 0", err, ErrOutOfRange)
	_, err = LogSeriesCDF(1, 1)
	requireErr(t, "p = 1", err, ErrOutOfRange)
	_, err = LogSeriesPMF(1, nan)
	requireErr(t, "p NaN", err, ErrNotFinite)
	_, err = LogSeriesQuantile(1, 0.5)
	requireErr(t, "u = 1", err, ErrOutOfRange)
	_, err = LogSeriesQuantile(2, 0.5)
	requireErr(t, "u > 1", err, ErrProbability)
	_, err = LogSeriesQuantile(nan, 0.5)
	requireErr(t, "u NaN", err, ErrNotFinite)
	_, err = LogSeriesQuantile(0.5, 2)
	requireErr(t, "quantile p", err, ErrOutOfRange)
	_, _, err = LogSeriesMeanVariance(-1)
	requireErr(t, "MV p", err, ErrOutOfRange)
	_, _, err = LogSeriesMeanVariance(1 - 0x1p-53)
	if err != nil {
		requireErr(t, "MV near 1", err, ErrOverflow)
	}
	requireNoAlloc(t, "LogSeriesQuantile", func() { sinkInt, _ = LogSeriesQuantile(0.9, 0.5) })
}

func BenchmarkLogSeriesPMF(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = LogSeriesPMF(5, 0.7)
	}
}

func BenchmarkLogSeriesCDF(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = LogSeriesCDF(20, 0.7)
	}
}

func BenchmarkLogSeriesQuantile(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkInt, _ = LogSeriesQuantile(0.9, 0.7)
	}
}

func BenchmarkLogSeriesMeanVariance(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _, _ = LogSeriesMeanVariance(0.7)
	}
}

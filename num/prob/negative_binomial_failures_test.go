package prob

import (
	"math"
	"testing"
)

// JKK ch. 5 eq. 5.1: P(Y = k) = C(k+r−1, k)pʳ(1−p)ᵏ; r = 2, p = ½: 0.25, 0.25, 0.1875, 0.125.
func TestNegativeBinomialFailuresPMF_sourceExample(t *testing.T) {
	want := []float64{0.25, 0.25, 0.1875, 0.125}
	for k, w := range want {
		if got, err := NegativeBinomialFailuresPMF(k, 2, 0.5); err != nil || !within(got, w, 1e-14) {
			t.Errorf("k %d: (%v, %v), want %v", k, got, err, w)
		}
	}
	if got, _ := NegativeBinomialFailuresPMF(-1, 2, 0.5); got != 0 {
		t.Errorf("k < 0: %v", got)
	}
	// Equals the trials form shifted by r, and the exact rational oracle.
	for k := range 30 {
		got, _ := NegativeBinomialFailuresPMF(k, 3, 0.3)
		tr, _ := NegativeBinomialTrialsPMF(k+3, 3, 0.3)
		if got != tr || !within(got, exactNegativeBinomialPMF(k+3, 3, 0.3), 1e-13) {
			t.Errorf("k %d: %v vs trials %v", k, got, tr)
		}
	}
	_, err := NegativeBinomialFailuresPMF(math.MaxInt, 2, 0.5)
	requireErr(t, "k + r overflow", err, ErrOutOfRange)
	_, err = NegativeBinomialFailuresPMF(1, 0, 0.5)
	requireErr(t, "r = 0", err, ErrOutOfRange)
	_, err = NegativeBinomialFailuresPMF(1, 2, 0)
	requireErr(t, "p = 0", err, ErrOutOfRange)
	_, err = NegativeBinomialFailuresPMF(1, 2, nan)
	requireErr(t, "p NaN", err, ErrNotFinite)
	_, err = NegativeBinomialFailuresPMF(1, 2, 1.5)
	requireErr(t, "p > 1", err, ErrProbability)
}

func TestNegativeBinomialFailuresCDF_sourceExample(t *testing.T) {
	if got, err := NegativeBinomialFailuresCDF(2, 2, 0.5); err != nil || !within(got, 0.6875, 1e-14) {
		t.Errorf("(2, 2, ½) = (%v, %v)", got, err)
	}
	if got, _ := NegativeBinomialFailuresCDF(-1, 2, 0.5); got != 0 {
		t.Errorf("k < 0: %v", got)
	}
	// CDF differences are the PMF.
	prev := 0.0
	for k := range 40 {
		c, _ := NegativeBinomialFailuresCDF(k, 4, 0.35)
		p, _ := NegativeBinomialFailuresPMF(k, 4, 0.35)
		if math.Abs(c-prev-p) > 1e-14 {
			t.Errorf("k %d: F(k) − F(k−1) = %v, pmf %v", k, c-prev, p)
		}
		prev = c
	}
	_, err := NegativeBinomialFailuresCDF(math.MaxInt, 2, 0.5)
	requireErr(t, "k + r overflow", err, ErrOutOfRange)
	_, err = NegativeBinomialFailuresCDF(1, 2, nan)
	requireErr(t, "p NaN", err, ErrNotFinite)
}

func TestNegativeBinomialFailuresQuantile(t *testing.T) {
	if got, err := NegativeBinomialFailuresQuantile(0.6, 2, 0.5); err != nil || got != 2 {
		t.Errorf("(0.6, 2, ½) = (%v, %v), want 2", got, err)
	}
	for _, u := range []float64{1e-6, 0.1, 0.5, 0.9, 0.999999} {
		q, err := NegativeBinomialFailuresQuantile(u, 5, 0.2)
		if err != nil {
			t.Fatalf("u %v: %v", u, err)
		}
		c, _ := NegativeBinomialFailuresCDF(q, 5, 0.2)
		cm, _ := NegativeBinomialFailuresCDF(q-1, 5, 0.2)
		if c < u || cm >= u {
			t.Errorf("u %v: q %d, F(q) %v, F(q−1) %v", u, q, c, cm)
		}
	}
	if q, _ := NegativeBinomialFailuresQuantile(0, 2, 0.5); q != 0 {
		t.Errorf("u = 0: %v", q)
	}
	if q, _ := NegativeBinomialFailuresQuantile(1, 2, 1); q != 0 {
		t.Errorf("p = 1: %v", q)
	}
	_, err := NegativeBinomialFailuresQuantile(1, 2, 0.5)
	requireErr(t, "u = 1", err, ErrOutOfRange)
	_, err = NegativeBinomialFailuresQuantile(-0.1, 2, 0.5)
	requireErr(t, "u < 0", err, ErrProbability)
	_, err = NegativeBinomialFailuresQuantile(nan, 2, 0.5)
	requireErr(t, "u NaN", err, ErrNotFinite)
	_, err = NegativeBinomialFailuresQuantile(0.5, 0, 0.5)
	requireErr(t, "r = 0", err, ErrOutOfRange)
}

func TestNegativeBinomialFailuresMeanVariance(t *testing.T) {
	if m, v, err := NegativeBinomialFailuresMeanVariance(2, 0.5); err != nil || m != 2 || v != 4 {
		t.Errorf("(2, ½) = (%v, %v, %v)", m, v, err)
	}
	// Against the PMF: Σk·P and Σk²·P.
	var s1, s2 float64
	for k := range 400 {
		p, _ := NegativeBinomialFailuresPMF(k, 3, 0.4)
		s1 += float64(k) * p
		s2 += float64(k*k) * p
	}
	m, v, _ := NegativeBinomialFailuresMeanVariance(3, 0.4)
	if !within(s1, m, 1e-12) || !within(s2-s1*s1, v, 1e-11) {
		t.Errorf("sums (%v, %v), closed (%v, %v)", s1, s2-s1*s1, m, v)
	}
	_, _, err := NegativeBinomialFailuresMeanVariance(1<<62, 1e-300)
	requireErr(t, "overflow", err, ErrOverflow)
	_, _, err = NegativeBinomialFailuresMeanVariance(0, 0.5)
	requireErr(t, "r = 0", err, ErrOutOfRange)
}

func TestNegativeBinomialFailures_noAlloc(t *testing.T) {
	requireNoAlloc(t, "NegativeBinomialFailuresPMF", func() { sinkFloat, _ = NegativeBinomialFailuresPMF(7, 3, 0.3) })
	requireNoAlloc(t, "NegativeBinomialFailuresCDF", func() { sinkFloat, _ = NegativeBinomialFailuresCDF(7, 3, 0.3) })
	requireNoAlloc(t, "NegativeBinomialFailuresQuantile", func() { sinkInt, _ = NegativeBinomialFailuresQuantile(0.7, 3, 0.3) })
}

func BenchmarkNegativeBinomialFailuresPMF(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = NegativeBinomialFailuresPMF(7, 3, 0.3)
	}
}

func BenchmarkNegativeBinomialFailuresCDF(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = NegativeBinomialFailuresCDF(7, 3, 0.3)
	}
}

func BenchmarkNegativeBinomialFailuresQuantile(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkInt, _ = NegativeBinomialFailuresQuantile(0.7, 3, 0.3)
	}
}

func BenchmarkNegativeBinomialFailuresMeanVariance(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _, _ = NegativeBinomialFailuresMeanVariance(3, 0.3)
	}
}

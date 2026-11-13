package prob

import (
	"math"
	"testing"
)

// JKK 3rd ed. §11.2.20: P(X = k) = k^(−a)/ζ(a); a = 2 gives 6/π² at k = 1 and (5/4)·6/π² for F(2).
func TestZipf(t *testing.T) {
	z2 := math.Pi * math.Pi / 6
	if got, err := ZipfPMF(1, 2); err != nil || !within(got, 1/z2, 1e-14) {
		t.Errorf("PMF(1) = (%v, %v)", got, err)
	}
	if got, _ := ZipfCDF(2, 2); !within(got, 1.25/z2, 1e-14) {
		t.Errorf("CDF(2) = %v", got)
	}
	if got, _ := ZipfPMF(0, 2); got != 0 {
		t.Errorf("PMF(0) = %v", got)
	}
	if got, _ := ZipfCDF(0, 2); got != 0 {
		t.Errorf("CDF(0) = %v", got)
	}
	// The Euler–Maclaurin tail matches the direct sum across the k = 1000 switch.
	for _, a := range []float64{1.1, 2, 3.5, 50} {
		var s compensatedSum
		for i := 1001; i <= 200000; i++ {
			s.add(math.Pow(float64(i), -a))
		}
		direct := s.total() + zipfTail(200001, a)
		if em := zipfTail(1001, a); !within(em, direct, 1e-12) {
			t.Errorf("a = %v: tail %v, direct %v", a, em, direct)
		}
		lo, _ := ZipfCDF(1000, a)
		hi, _ := ZipfCDF(1001, a)
		p, _ := ZipfPMF(1001, a)
		if math.Abs(hi-lo-p) > 1e-13 {
			t.Errorf("a = %v: F(1001) − F(1000) = %v, PMF = %v", a, hi-lo, p)
		}
	}
	if got, _ := ZipfCDF(math.MaxInt, 2); !within(got, 1, 1e-15) {
		t.Errorf("CDF(MaxInt) = %v", got)
	}
	if got, _ := ZipfCDF(5000, 150); got != 1 {
		t.Errorf("CDF large a = %v", got)
	}
	// ζ(3) = 1.2020569031595942 (A&S Table 23.3), ζ(4) = π⁴/90.
	z3, z4 := 1.2020569031595942, math.Pow(math.Pi, 4)/90
	if m, v, err := ZipfMeanVariance(4); err != nil || !within(m, z3/z4, 1e-14) || !within(v, z2/z4-(z3/z4)*(z3/z4), 1e-13) {
		t.Errorf("MV(4) = (%v, %v, %v)", m, v, err)
	}
	_, err := ZipfPMF(1, 1)
	requireErr(t, "a = 1", err, ErrOutOfRange)
	_, err = ZipfCDF(1, nan)
	requireErr(t, "a NaN", err, ErrNotFinite)
	_, _, err = ZipfMeanVariance(3)
	requireErr(t, "MV a = 3", err, ErrOutOfRange)
	_, _, err = ZipfMeanVariance(inf)
	requireErr(t, "MV a Inf", err, ErrNotFinite)
	requireNoAlloc(t, "ZipfCDF", func() { sinkFloat, _ = ZipfCDF(50, 2) })
}

func BenchmarkZipfPMF(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = ZipfPMF(7, 2.5)
	}
}

func BenchmarkZipfCDF(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = ZipfCDF(50, 2.5)
	}
}

func BenchmarkZipfMeanVariance(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _, _ = ZipfMeanVariance(4.5)
	}
}

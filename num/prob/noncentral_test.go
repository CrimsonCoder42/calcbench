package prob

import (
	"math"
	"testing"
)

// noncentralRange returns the Poisson indices that carry all but far below 2⁻⁵³ of the mass of mean μ.
func noncentralRange(mu float64) (int, int) {
	w := 40*math.Sqrt(mu+1) + 50
	return int(math.Max(0, mu-w)), int(mu + w)
}

// Brute-force mixture oracles (JKB vol. 2 ch. 29 and 30): every term from the central function directly.
func chiMixture(x, df, lambda float64, central func(x, df float64) (float64, error)) float64 {
	mu := lambda / 2
	lo, hi := noncentralRange(mu)
	var s compensatedSum
	for j := lo; j <= hi; j++ {
		w, _ := PoissonPMF(j, mu)
		c, _ := central(x, df+2*float64(j))
		s.add(w * c)
	}
	return s.total()
}

func fMixture(x, df1, df2, lambda float64) float64 {
	mu := lambda / 2
	lo, hi := noncentralRange(mu)
	var s compensatedSum
	for j := lo; j <= hi; j++ {
		w, _ := PoissonPMF(j, mu)
		d := df1 + 2*float64(j)
		c, _ := FDistributionCDF(df1*x/d, d, df2)
		s.add(w * c)
	}
	return s.total()
}

// JKB vol. 2 ch. 29: with df = 1, X = (Z + √λ)², so F(x) = Φ(√x − √λ) − Φ(−√x − √λ) and
// f(x) = (φ(√x − √λ) + φ(√x + √λ))/(2√x).
func TestNoncentralChiSquareDF1(t *testing.T) {
	for _, lambda := range []float64{0.01, 0.5, 1, 4, 25, 100} {
		for _, x := range []float64{0.001, 0.1, 1, 3, 10, 50, 200} {
			sx, sl := math.Sqrt(x), math.Sqrt(lambda)
			a, _ := NormalCDF(sx-sl, 0, 1)
			b, _ := NormalCDF(-sx-sl, 0, 1)
			got, err := NoncentralChiSquareCDF(x, 1, lambda)
			if err != nil || math.Abs(got-(a-b)) > 1e-13 {
				t.Errorf("CDF(%v, 1, %v) = (%v, %v), want %v", x, lambda, got, err, a-b)
			}
			p, _ := NormalPDF(sx-sl, 0, 1)
			q, _ := NormalPDF(sx+sl, 0, 1)
			want := (p + q) / (2 * sx)
			gotPDF, err := NoncentralChiSquarePDF(x, 1, lambda)
			if err != nil || (want > 1e-280 && !within(gotPDF, want, 1e-12)) {
				t.Errorf("PDF(%v, 1, %v) = (%v, %v), want %v", x, lambda, gotPDF, err, want)
			}
		}
	}
	if got, _ := NoncentralChiSquarePDF(1, 1, 1); !within(got, 0.2264666234573104, 1e-13) {
		t.Errorf("PDF example %v", got)
	}
	if got, _ := NoncentralChiSquareCDF(1, 1, 1); math.Abs(got-0.4772498680518208) > 1e-14 {
		t.Errorf("CDF example %v", got)
	}
}

func TestNoncentralChiSquareMixture(t *testing.T) {
	for _, df := range []float64{0.5, 2, 3, 10, 40, 300} {
		for _, lambda := range []float64{0.3, 2, 15, 80, 1e4} {
			mean := df + lambda
			sd := math.Sqrt(2 * (df + 2*lambda))
			for _, s := range []float64{-3, -1, 0, 1, 3, 8} {
				x := mean + s*sd
				if x <= 0 {
					continue
				}
				want := chiMixture(x, df, lambda, ChiSquareCDF)
				got, err := NoncentralChiSquareCDF(x, df, lambda)
				if err != nil || math.Abs(got-want) > 1e-12 {
					t.Errorf("CDF(%v, %v, %v) = (%v, %v), mixture %v", x, df, lambda, got, err, want)
				}
				wantPDF := chiMixture(x, df, lambda, ChiSquarePDF)
				gotPDF, err := NoncentralChiSquarePDF(x, df, lambda)
				if err != nil || !within(gotPDF, wantPDF, 1e-10) {
					t.Errorf("PDF(%v, %v, %v) = (%v, %v), mixture %v", x, df, lambda, gotPDF, err, wantPDF)
				}
			}
		}
	}
	// λ = 0 is the central distribution.
	c, _ := ChiSquareCDF(4, 3)
	if got, _ := NoncentralChiSquareCDF(4, 3, 0); got != c {
		t.Errorf("λ = 0: %v, central %v", got, c)
	}
	if got, _ := NoncentralChiSquarePDF(0, 2, 2); !within(got, math.Exp(-1)/2, 1e-15) {
		t.Errorf("PDF(0, 2, 2) = %v", got)
	}
	if got, _ := NoncentralChiSquarePDF(0, 3, 2); got != 0 {
		t.Errorf("PDF(0, 3, 2) = %v", got)
	}
	_, err := NoncentralChiSquarePDF(0, 1, 2)
	requireErr(t, "PDF x = 0, df < 2", err, ErrOutOfRange)
	_, err = NoncentralChiSquareCDF(1, 3, -1)
	requireErr(t, "λ < 0", err, ErrOutOfRange)
	_, err = NoncentralChiSquareCDF(1, 3, 2e5)
	requireErr(t, "λ > 1e5", err, ErrOutOfRange)
	_, err = NoncentralChiSquareCDF(math.NaN(), 3, 1)
	requireErr(t, "x NaN", err, ErrNotFinite)
	_, err = NoncentralChiSquarePDF(1, 0, 1)
	requireErr(t, "df = 0", err, ErrOutOfRange)
	requireNoAlloc(t, "NoncentralChiSquareCDF", func() { sinkFloat, _ = NoncentralChiSquareCDF(12, 5, 7) })
	requireNoAlloc(t, "NoncentralChiSquarePDF", func() { sinkFloat, _ = NoncentralChiSquarePDF(12, 5, 7) })
}

// JKB vol. 2 ch. 29: mean ν + λ, variance 2(ν + 2λ); ch. 30: the noncentral F moments.
func TestNoncentralMeanVariance(t *testing.T) {
	if m, v, err := NoncentralChiSquareMeanVariance(3, 2); err != nil || m != 5 || v != 14 {
		t.Errorf("chi-square (3, 2) = (%v, %v, %v)", m, v, err)
	}
	m, v, err := NoncentralFMeanVariance(4, 10, 2)
	if err != nil || !within(m, 1.875, 1e-15) || !within(v, 1250.0/384, 1e-15) {
		t.Errorf("F (4, 10, 2) = (%v, %v, %v)", m, v, err)
	}
	// λ = 0 agrees with the central F.
	cm, cv, _ := FDistributionMeanVariance(4, 10)
	if m, v, _ := NoncentralFMeanVariance(4, 10, 0); !within(m, cm, 1e-15) || !within(v, cv, 1e-14) {
		t.Errorf("F λ = 0: (%v, %v), central (%v, %v)", m, v, cm, cv)
	}
	_, _, err = NoncentralFMeanVariance(4, 4, 1)
	requireErr(t, "F df2 = 4", err, ErrOutOfRange)
	_, _, err = NoncentralChiSquareMeanVariance(3, math.Inf(1))
	requireErr(t, "λ Inf", err, ErrNotFinite)
}

// Lenth (1987) AS 226 sums the same mixture; the oracle forms every term with FDistributionCDF.
func TestNoncentralFCDF(t *testing.T) {
	for _, df := range [][2]float64{{1, 1}, {2, 5}, {4, 10}, {10, 3}, {30, 60}} {
		for _, lambda := range []float64{0.5, 3, 20, 150} {
			for _, x := range []float64{0.01, 0.3, 1, 2, 5, 20, 200} {
				want := fMixture(x, df[0], df[1], lambda)
				got, err := NoncentralFCDF(x, df[0], df[1], lambda)
				if err != nil || math.Abs(got-want) > 1e-12 {
					t.Errorf("NoncentralFCDF(%v, %v, %v, %v) = (%v, %v), mixture %v", x, df[0], df[1], lambda, got, err, want)
				}
			}
		}
	}
	c, _ := FDistributionCDF(1, 2, 2)
	if got, _ := NoncentralFCDF(1, 2, 2, 0); got != c || !within(got, 0.5, 1e-15) {
		t.Errorf("λ = 0: %v, central %v", got, c)
	}
	if got, _ := NoncentralFCDF(-1, 2, 2, 1); got != 0 {
		t.Errorf("x < 0: %v", got)
	}
	_, err := NoncentralFCDF(1, 0, 2, 1)
	requireErr(t, "df1 = 0", err, ErrOutOfRange)
	_, err = NoncentralFCDF(1, 2, 2, -1)
	requireErr(t, "λ < 0", err, ErrOutOfRange)
	requireNoAlloc(t, "NoncentralFCDF", func() { sinkFloat, _ = NoncentralFCDF(2, 4, 10, 3) })
}

func BenchmarkNoncentralChiSquarePDF(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = NoncentralChiSquarePDF(12, 5, 7)
	}
}

func BenchmarkNoncentralChiSquareCDF(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = NoncentralChiSquareCDF(12, 5, 7)
	}
}

func BenchmarkNoncentralChiSquareMeanVariance(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _, _ = NoncentralChiSquareMeanVariance(5, 7)
	}
}

func BenchmarkNoncentralFCDF(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = NoncentralFCDF(2, 4, 10, 3)
	}
}

func BenchmarkNoncentralFMeanVariance(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _, _ = NoncentralFMeanVariance(4, 10, 3)
	}
}

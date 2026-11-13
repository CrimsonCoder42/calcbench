package stats

import (
	"fmt"
	"math"
)

// BinRule selects a rule for the number or the width of equal-width histogram classes.
type BinRule int

// The four bin rules of HistogramBinEdges.
const (
	// BinSturges: k = ⌈1 + log₂ n⌉ classes (Sturges 1926).
	BinSturges BinRule = 1 + iota
	// BinScott: width h = (24√π)^(1/3)·s·n^(−1/3) ≈ 3.49·s·n^(−1/3), s the sample standard deviation (Scott 1979).
	BinScott
	// BinFreedmanDiaconis: width h = 2·IQR·n^(−1/3), IQR by Hyndman–Fan type 7 (Freedman & Diaconis 1981).
	BinFreedmanDiaconis
	// BinDoane: k = ⌈1 + log₂ n + log₂(1 + |g₁|/σ_g₁)⌉, g₁ the moment skewness and
	// σ_g₁ = √(6(n − 2)/((n + 1)(n + 3))) (Doane 1976).
	BinDoane
)

// maxBins bounds the number of classes HistogramBinEdges returns (a width rule on data with a few far
// outliers can ask for an absurd count).
const maxBins = 1 << 20

// HistogramBinEdges returns k + 1 equal-width class edges spanning [min, max] of values, with k chosen
// by rule: edges[0] = min and edges[k] = max exactly, edges[i] = min + (i/k)·(max − min). For a width
// rule (Scott, Freedman–Diaconis) k = ⌈(max − min)/h⌉, at least 1. The edges feed FrequencyDistribution.
// values is read, never modified (Freedman–Diaconis selects on a copy).
// Errors: ErrEmpty when values is empty; ErrTooFew when there are fewer than 2 values (3 for BinDoane);
// ErrNotFinite when a value is NaN or ±Inf; ErrZeroVariance when every value is equal, or the
// interquartile range is 0 for BinFreedmanDiaconis; ErrOverflow when (BinScott) the standard deviation
// is beyond float64; ErrOutOfRange when rule is not one of the four, or
// the rule asks for more than 2^20 classes, or the range is too narrow for k distinct edges.
//
// ID:      NUM-STATS-043
// Time:    O(n)   n = len(values); a few passes (Freedman–Diaconis adds two linear-time selections)
// Space:   O(k) for the result (one allocation); O(n) extra for BinFreedmanDiaconis only
// Source:  Sturges, "The choice of a class interval", Journal of the American Statistical Association, 1926, 21(153):65–66, doi:10.1080/01621459.1926.10502161 (k = 1 + log₂ n)
// Source:  Scott, "On optimal and data-based histograms", Biometrika, 1979, 66(3):605–610, doi:10.1093/biomet/66.3.605 (eq. for h*: (24√π)^(1/3) σ n^(−1/3) ≈ 3.49 s n^(−1/3))
// Source:  Freedman, Diaconis, "On the histogram as a density estimator: L₂ theory", Zeitschrift für Wahrscheinlichkeitstheorie und verwandte Gebiete, 1981, 57(4):453–476, doi:10.1007/BF01025868 (h = 2·IQR·n^(−1/3))
// Source:  Doane, "Aesthetic frequency classifications", The American Statistician, 1976, 30(4):181–183, doi:10.1080/00031305.1976.10479172 (Sturges corrected by the skewness term)
// Tags:    histogram, bins, bin-width, sturges, scott, freedman-diaconis, doane, frequency-distribution, descriptive
// Units:   values in any one unit; edges in that unit
// Example: values 1 … 8, BinSturges → k = ⌈1 + 3⌉ = 4 classes, edges 1, 2.75, 4.5, 6.25, 8 (hand-checked)
func HistogramBinEdges(values []float64, rule BinRule) ([]float64, error) {
	if rule < BinSturges || rule > BinDoane {
		return nil, fmt.Errorf("%w: rule = %d", ErrOutOfRange, int(rule))
	}
	minLen := 2
	if rule == BinDoane {
		minLen = 3
	}
	if err := checkSample(values, minLen); err != nil {
		return nil, err
	}
	lo, hi := values[0], values[0]
	for _, v := range values[1:] {
		lo, hi = min(lo, v), max(hi, v)
	}
	if lo == hi {
		return nil, fmt.Errorf("%w: every value equals %v", ErrZeroVariance, lo)
	}
	n := float64(len(values))
	halfRange := hi/2 - lo/2 // cannot overflow
	var k float64
	switch rule {
	case BinSturges:
		k = math.Ceil(1 + math.Log2(n))
	case BinDoane:
		g1 := skewnessScaled(values, halfRange)
		sigma := math.Sqrt(6 * (n - 2) / ((n + 1) * (n + 3)))
		k = math.Ceil(1 + math.Log2(n) + math.Log2(1+math.Abs(g1)/sigma))
	case BinScott, BinFreedmanDiaconis:
		var spread float64 // h/2 divided into halfRange keeps both in range
		if rule == BinScott {
			s, err := SampleStandardDeviation(values)
			if err != nil {
				return nil, err
			}
			spread = math.Cbrt(24*math.Sqrt(math.Pi)) * s / 2
		} else {
			q, _ := Percentiles(values, []float64{0.25, 0.75}, HyndmanFan7) // cannot fail: checked
			iqr := q[1]/2 - q[0]/2
			if iqr == 0 {
				return nil, fmt.Errorf("%w: interquartile range is 0", ErrZeroVariance)
			}
			spread = 2 * iqr
		}
		k = math.Ceil(halfRange / (spread / math.Cbrt(n)))
	}
	if !(k <= maxBins) {
		return nil, fmt.Errorf("%w: rule asks for %v classes, want ≤ %d", ErrOutOfRange, k, maxBins)
	}
	kb := max(int(k), 1)
	edges := make([]float64, kb+1)
	for i := range edges {
		edges[i] = lerp(lo, hi, float64(i)/float64(kb))
	}
	edges[kb] = hi
	for i := 1; i <= kb; i++ {
		if !(edges[i] > edges[i-1]) {
			return nil, fmt.Errorf("%w: range %v..%v too narrow for %d distinct classes", ErrOutOfRange, lo, hi, kb)
		}
	}
	return edges, nil
}

// skewnessScaled returns the moment skewness g₁ = m₃/m₂^(3/2) of finite values that are not all equal,
// computed on deviations divided by scale > 0 (g₁ is scale-free, and the division keeps the cubes finite).
func skewnessScaled(values []float64, scale float64) float64 {
	m := meanOf(values)
	var s2, s3 compensatedSum
	for _, v := range values {
		d := (v/2 - m/2) / scale // |d| ≤ 1
		s2.add(d * d)
		s3.add(d * d * d)
	}
	n := float64(len(values))
	m2, m3 := s2.total()/n, s3.total()/n
	return m3 / (m2 * math.Sqrt(m2))
}

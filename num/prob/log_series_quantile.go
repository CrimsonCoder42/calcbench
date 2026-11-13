package prob

import "fmt"

// logSeriesMaxSteps caps LogSeriesQuantile's walk; past it the quantile is reported as ErrOutOfRange.
const logSeriesMaxSteps = 1 << 28

// LogSeriesQuantile returns the smallest k ≥ 1 with P(X ≤ k) ≥ u for the log-series distribution. It
// walks the PMF upward with the ratio p·i/(i + 1), summing compensated, so the result agrees with
// LogSeriesCDF (NUM-PROB-135) up to rounding. u = 0 returns 1.
// Errors: ErrNotFinite when u or p is NaN or ±Inf; ErrProbability when u is outside [0, 1]; ErrOutOfRange
// unless 0 < p < 1, when u = 1 (no finite quantile), or past 2²⁸ steps (p extremely close to 1).
//
// ID:      NUM-PROB-136
// Time:    O(k)   k = the result
// Space:   O(1), no allocation
// Source:  Johnson, Kemp, Kotz, Univariate Discrete Distributions, 3rd ed., Wiley, 2005, ch. 7, §7.1 (the PMF; inversion by sequential search)
// Tags:    probability, distribution, log-series, logarithmic, quantile, inverse-cdf, discrete
// Units:   u, p: probabilities; result: a count (≥ 1)
// Example: p = 0.5: u = 0.8 → 2 (F(1) ≈ 0.7213 < 0.8 ≤ F(2) ≈ 0.9017)
func LogSeriesQuantile(u, p float64) (int, error) {
	if err := checkProbability("u", u); err != nil {
		return 0, err
	}
	if err := checkLogSeries(p); err != nil {
		return 0, err
	}
	if u == 1 {
		return 0, fmt.Errorf("%w: u = 1, quantile infinite", ErrOutOfRange)
	}
	t, _ := LogSeriesPMF(1, p)
	var s compensatedSum
	for k := 1; k <= logSeriesMaxSteps; k++ {
		s.add(t)
		if s.total() >= u {
			return k, nil
		}
		t *= p * float64(k) / float64(k+1)
		if t == 0 {
			return k + 1, nil // the remaining mass is below the smallest subnormal
		}
	}
	return 0, fmt.Errorf("%w: quantile for u = %v, p = %v needs more than 2^28 steps", ErrOutOfRange, u, p)
}

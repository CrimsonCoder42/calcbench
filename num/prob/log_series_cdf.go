package prob

import "math"

// LogSeriesCDF returns P(X ≤ k) = Σᵢ₌₁ᵏ −pⁱ/(i·ln(1 − p)) for the log-series distribution. The terms fall
// monotonically (ratio p·i/(i + 1) < p), so the compensated sum stops once the rest of the tail, at most
// term·p/(1 − p), is below 2⁻⁶⁰ of the sum; the result is then min(1, sum). k < 1 gives 0.
// Errors: ErrNotFinite when p is NaN or ±Inf; ErrOutOfRange unless 0 < p < 1.
//
// ID:      NUM-PROB-135
// Time:    O(min(k, T))   T ≈ ln(2⁻⁶⁰(1 − p))/ln p, the terms until the tail is negligible
// Space:   O(1), no allocation
// Source:  Johnson, Kemp, Kotz, Univariate Discrete Distributions, 3rd ed., Wiley, 2005, ch. 7, §7.1 (the PMF), summed
// Tags:    probability, distribution, log-series, logarithmic, cdf, discrete
// Units:   k: a count; p: a probability; result: a probability in [0, 1]
// Example: p = 0.5, k = 2 → (0.5 + 0.125)/ln 2 ≈ 0.9016844005556022
func LogSeriesCDF(k int, p float64) (float64, error) {
	if err := checkLogSeries(p); err != nil {
		return 0, err
	}
	if k < 1 {
		return 0, nil
	}
	c := -1 / math.Log1p(-p)
	var s compensatedSum
	pi := p // pⁱ
	for i := 1; i <= k; i++ {
		t := c * pi / float64(i)
		s.add(t)
		if t*p <= 0x1p-60*(1-p)*s.total() {
			break
		}
		pi *= p
	}
	return min(1, s.total()), nil
}

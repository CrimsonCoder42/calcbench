package prob

import (
	"fmt"
	"math"
)

// checkLogSeries returns ErrNotFinite for NaN or ±Inf and ErrOutOfRange unless 0 < p < 1.
func checkLogSeries(p float64) error {
	if err := checkFinite("p", p); err != nil {
		return err
	}
	if p <= 0 || p >= 1 {
		return fmt.Errorf("%w: p = %v, want 0 < p < 1", ErrOutOfRange, p)
	}
	return nil
}

// LogSeriesPMF returns P(X = k) = −p^k / (k·ln(1 − p)), k = 1, 2, …, the logarithmic (log-series)
// distribution. ln(1 − p) is math.Log1p(−p) and p^k is exp(k·ln p), so small p and large k keep their
// accuracy; k < 1 gives 0.
// Errors: ErrNotFinite when p is NaN or ±Inf; ErrOutOfRange unless 0 < p < 1.
//
// ID:      NUM-PROB-134
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Johnson, Kemp, Kotz, Univariate Discrete Distributions, 3rd ed., Wiley, 2005, ch. 7 "Logarithmic and Lagrangian Distributions", §7.1 (the PMF)
// Tags:    probability, distribution, log-series, logarithmic, pmf, discrete, species-abundance
// Units:   k: a count (≥ 1); p: a probability; result: a probability in [0, 1]
// Example: p = 0.5, k = 1 → 0.5/ln 2 ≈ 0.7213475204444817
func LogSeriesPMF(k int, p float64) (float64, error) {
	if err := checkLogSeries(p); err != nil {
		return 0, err
	}
	if k < 1 {
		return 0, nil
	}
	return -math.Exp(float64(k)*math.Log(p)) / (float64(k) * math.Log1p(-p)), nil
}

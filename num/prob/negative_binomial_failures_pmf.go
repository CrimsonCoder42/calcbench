package prob

import (
	"fmt"
	"math"
)

// negativeBinomialTrialsFor converts a failures count k to the trials count k + r, or ErrOutOfRange when
// that exceeds the int range.
func negativeBinomialTrialsFor(k, r int) (int, error) {
	if k > math.MaxInt-r {
		return 0, fmt.Errorf("%w: k + r = %d + %d exceeds the int range", ErrOutOfRange, k, r)
	}
	return k + r, nil
}

// NegativeBinomialFailuresPMF returns P(Y = k) = C(k + r − 1, k)·pʳ·(1 − p)ᵏ, where Y counts the FAILURES
// before the r-th success (support 0, 1, 2, …; the numpy and scipy convention); 0 for k < 0. Y = X − r for
// X the trials count of NegativeBinomialTrialsPMF (NUM-PROB-032), which evaluates it (saddle-point form).
// Errors: ErrOutOfRange when r < 1, p = 0, or k + r exceeds the int range; ErrNotFinite when p is NaN or
// ±Inf; ErrProbability when p is outside [0, 1].
//
// ID:      NUM-PROB-099
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Johnson, Kemp, Kotz, Univariate Discrete Distributions, 3rd ed., Wiley, 2005, ch. 5 "Negative Binomial Distribution", §5.1 (definition, eq. 5.1)
// Tags:    probability, distribution, negative-binomial, pmf, discrete, failures, overdispersion
// Units:   k: a count of failures; r: a count of successes; p and result: probabilities in [0, 1]
// Example: r = 2, p = 0.5, k = 2 → C(3,2)·0.25·0.25 = 0.1875
func NegativeBinomialFailuresPMF(k, r int, p float64) (float64, error) {
	if err := checkNegativeBinomial(r, p); err != nil {
		return 0, err
	}
	if k < 0 {
		return 0, nil
	}
	x, err := negativeBinomialTrialsFor(k, r)
	if err != nil {
		return 0, err
	}
	return NegativeBinomialTrialsPMF(x, r, p)
}

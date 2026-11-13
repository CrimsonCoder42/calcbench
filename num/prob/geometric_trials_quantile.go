package prob

import (
	"fmt"
	"math"
)

// GeometricTrialsQuantile returns the smallest k ≥ 1 with P(X ≤ k) ≥ u, X the number of TRIALS up to and
// including the first success (D-013). The estimate ⌈log1p(−u)/log1p(−p)⌉ is corrected by single steps
// until GeometricTrialsCDF(k−1) < u ≤ GeometricTrialsCDF(k) holds for the CDF as this package computes it.
// The failures-count quantile is k − 1.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrProbability when u or p is outside [0, 1];
// ErrOutOfRange when p = 0, or u = 1 with p < 1 (no finite quantile); ErrOverflow when k exceeds the int
// range.
//
// ID:      NUM-PROB-057
// Time:    O(1 + k·2⁻⁵⁰)   k = the result (the float64 estimate is off by at most a few units in the last place of k, each a correction step)
// Space:   O(1), no allocation
// Source:  Ross, A First Course in Probability, 10th ed., Pearson, 2019, §4.8 (other discrete distributions: the geometric random variable), inverted
// Source:  Casella, Berger, Statistical Inference, 2nd ed., Duxbury, 2002, §3.2 "Discrete Distributions" (geometric distribution), inverted
// Tags:    probability, distribution, geometric, quantile, inverse-cdf, trials
// Units:   u, p: probabilities; result: a count of trials
// Example: p = 0.2, u = 0.5 → 4 (F(3) = 0.488 < 0.5 ≤ F(4) = 0.5904)
func GeometricTrialsQuantile(u, p float64) (int, error) {
	if err := checkProbability("u", u); err != nil {
		return 0, err
	}
	if err := checkGeometric(p); err != nil {
		return 0, err
	}
	if p == 1 || u <= p {
		return 1, nil // F(1) = p
	}
	if u == 1 {
		return 0, fmt.Errorf("%w: u = 1 with p < 1, quantile infinite", ErrOutOfRange)
	}
	est := math.Ceil(math.Log1p(-u) / math.Log1p(-p))
	if !(est < math.MaxInt-2) { // also catches a NaN or +Inf estimate
		return 0, fmt.Errorf("%w: quantile for u = %v, p = %v", ErrOverflow, u, p)
	}
	k := max(1, int(est))
	for k > 1 && geometricTrialsCDF(k-1, p) >= u {
		k--
	}
	for geometricTrialsCDF(k, p) < u {
		if k == math.MaxInt {
			return 0, fmt.Errorf("%w: quantile for u = %v, p = %v", ErrOverflow, u, p)
		}
		k++
	}
	return k, nil
}

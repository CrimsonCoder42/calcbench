package prob

import "fmt"

// ChebyshevBoundEpsilon returns min(1, σ²/ε²), Chebyshev's upper bound on P(|X − μ| ≥ ε) for X with
// variance σ². It is a bound, not a probability.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrOutOfRange when variance < 0 or epsilon ≤ 0.
//
// ID:      NUM-PROB-021
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Ross, A First Course in Probability, 10th ed., Pearson, 2019, §8.2 "Chebyshev's Inequality and the Weak Law of Large Numbers"
// Source:  Casella, Berger, Statistical Inference, 2nd ed., Duxbury, 2002, §3.6 "Inequalities and Identities" (P(g(X) ≥ r) ≤ E g(X)/r for g ≥ 0)
// Tags:    probability, chebyshev, inequality, tail-bound
// Units:   variance: the caller's unit squared; epsilon: that unit; result: a probability bound in [0, 1]
// Example: variance 25, ε = 10 → 0.25, so a mean-50 variable has P(40 < X < 60) ≥ 0.75
func ChebyshevBoundEpsilon(variance, epsilon float64) (float64, error) {
	if err := checkFinite("variance", variance); err != nil {
		return 0, err
	}
	if variance < 0 {
		return 0, fmt.Errorf("%w: variance = %v, want >= 0", ErrOutOfRange, variance)
	}
	if err := checkPositive("epsilon", epsilon); err != nil {
		return 0, err
	}
	// Two divisions instead of ε²: ε² can overflow or underflow where the ratio is representable.
	return min(1, variance/epsilon/epsilon), nil
}

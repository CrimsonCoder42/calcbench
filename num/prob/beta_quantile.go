package prob

import "fmt"

// BetaQuantile returns the u-quantile of the beta distribution with shapes α and β: the x in [0, 1] with
// I_x(α, β) = u. It is found by bisection on [0, 1] (on the survival for u > ½) until the bracket is within
// 4 ulps. u = 0 returns 0 and u = 1 returns 1.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrProbability when u is outside [0, 1];
// ErrOutOfRange when α ≤ 0 or β ≤ 0, or the incomplete beta does not converge.
//
// ID:      NUM-PROB-075
// Time:    O(B·(√max(α, β) + log(1/ε)))   B ≤ 1,300 bisection steps (about 60 for a result above 1e-3)
// Space:   O(1)
// Source:  Johnson, Kotz, Balakrishnan, Continuous Univariate Distributions, vol. 2, 2nd ed., Wiley, 1995, ch. 25 "Beta Distributions", §1 (cdf), inverted
// Source:  Burden, Faires, Numerical Analysis, 9th ed., Brooks/Cole, 2011, §2.1 (bisection)
// Tags:    probability, distribution, beta, quantile, inverse-cdf, ppf
// Units:   dimensionless
// Example: α = 2, β = 1, u = 0.25 → √0.25 = 0.5 (I_x(2, 1) = x²)
func BetaQuantile(u, alpha, beta float64) (float64, error) {
	if err := checkProbability("u", u); err != nil {
		return 0, err
	}
	if err := checkBetaShapes(alpha, beta); err != nil {
		return 0, err
	}
	if u == 0 || u == 1 {
		return u, nil
	}
	x, err := quantileByBisection(u, 0, 1, func(x float64) (float64, float64, error) {
		if x <= 0 {
			return 0, 1, nil
		}
		if x >= 1 {
			return 1, 0, nil
		}
		return betaCDFSurvival(x, alpha, beta)
	})
	if err != nil {
		return 0, err
	}
	if x > 1 {
		return 0, fmt.Errorf("%w: beta quantile above 1", ErrOutOfRange) // unreachable: cdf(1) = 1 ≥ u
	}
	return x, nil
}

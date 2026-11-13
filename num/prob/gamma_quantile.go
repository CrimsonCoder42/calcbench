package prob

import "fmt"

// GammaQuantile returns the u-quantile of the gamma distribution with shape k and scale θ: the x with
// P(X ≤ x) = u. It has no closed form; it is found by bisection on GammaCDF (u ≤ ½) or on the survival
// Q(k, x/θ) (u > ½, so the upper tail keeps its accuracy), on a bracket that starts at [0, kθ] and doubles,
// until the bracket is within 4 ulps. u = 0 returns 0.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrProbability when u is outside [0, 1];
// ErrOutOfRange when u = 1 (an infinite quantile), shape ≤ 0 or scale ≤ 0, or the incomplete gamma does not
// converge; ErrOverflow when the quantile exceeds the float64 range.
//
// ID:      NUM-PROB-071
// Time:    O(B·(√k + log(1/ε)))   B ≤ 1,300 bisection steps (about 70 for a normal-range result), k = shape
// Space:   O(1); the bisection closure may allocate once
// Source:  Johnson, Kotz, Balakrishnan, Continuous Univariate Distributions, vol. 1, 2nd ed., Wiley, 1994, ch. 17 "Gamma Distributions", §1 (cdf), inverted
// Source:  Burden, Faires, Numerical Analysis, 9th ed., Brooks/Cole, 2011, §2.1 (bisection)
// Tags:    probability, distribution, gamma, quantile, inverse-cdf, ppf
// Units:   u a probability in [0, 1); scale and result: one unit; shape dimensionless
// Example: k = 1, θ = 1, u = 0.5 → ln 2 ≈ 0.6931471805599453 (the exponential case)
func GammaQuantile(u, shape, scale float64) (float64, error) {
	if err := checkProbability("u", u); err != nil {
		return 0, err
	}
	if err := checkGamma(shape, scale); err != nil {
		return 0, err
	}
	if u == 1 {
		return 0, fmt.Errorf("%w: u = 1, quantile infinite", ErrOutOfRange)
	}
	if u == 0 {
		return 0, nil
	}
	return quantileByBisection(u, 0, shape*scale, func(x float64) (float64, float64, error) {
		if x <= 0 {
			return 0, 1, nil
		}
		return gammaCDFSurvival(x, shape, scale)
	})
}

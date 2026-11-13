package prob

import "fmt"

// StudentTQuantile returns the u-quantile of Student's t distribution with ν = df degrees of freedom: the t
// with P(T ≤ t) = u. By symmetry it solves P(T ≤ −s) = min(u, 1 − u) for s ≥ 0 by bisection on the
// two-sided tail 2·P(T ≤ −s) (held exactly, so a tiny u keeps its accuracy) until the bracket is within 4
// ulps, and returns −s for u < ½ and s for u > ½; u = ½ returns 0.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrProbability when u is outside [0, 1];
// ErrOutOfRange when u is 0 or 1 (an infinite quantile), df ≤ 0, or the incomplete beta does not converge;
// ErrOverflow when the quantile exceeds the float64 range.
//
// ID:      NUM-PROB-079
// Time:    O(B·(√ν + log(1/ε)))   B ≤ 1,300 bisection steps (about 70 in practice), ν = df
// Space:   O(1)
// Source:  Johnson, Kotz, Balakrishnan, Continuous Univariate Distributions, vol. 2, 2nd ed., Wiley, 1995, ch. 28 "t-Distributions", §2 (cdf), inverted
// Source:  Burden, Faires, Numerical Analysis, 9th ed., Brooks/Cole, 2011, §2.1 (bisection)
// Tags:    probability, distribution, student-t, t-distribution, quantile, inverse-cdf, critical-value
// Units:   dimensionless
// Example: ν = 1, u = 0.975 → tan(0.475π) ≈ 12.706204736174698; ν = 10, u = 0.975 → 2.228 (NIST e-Handbook §1.3.6.7.2 table)
func StudentTQuantile(u, df float64) (float64, error) {
	if err := checkProbability("u", u); err != nil {
		return 0, err
	}
	if err := checkPositive("df", df); err != nil {
		return 0, err
	}
	if u == 0 || u == 1 {
		return 0, fmt.Errorf("%w: u = %v, quantile infinite", ErrOutOfRange, u)
	}
	if u == 0.5 {
		return 0, nil
	}
	p := min(u, 1-u) // exact for u in (½, 1) (Sterbenz)
	// |T| has P(|T| > s) = 2·tail(s); solve for its (1 − 2p)-quantile with the upper tail 2p held exactly.
	s, err := quantileByBisectionTail(1-2*p, 2*p, 0, 1, func(s float64) (float64, float64, error) {
		tail, _, err := studentTTail(s, df)
		return 1 - 2*tail, 2 * tail, err
	})
	if err != nil {
		return 0, err
	}
	if u < 0.5 {
		return -s, nil
	}
	return s, nil
}

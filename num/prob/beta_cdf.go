package prob

import "github.com/CrimsonCoder42/calcbench/num/special"

// betaCDFSurvival returns I_x(α, β) and 1 − I_x(α, β) = I_(1−x)(β, α) for validated shapes, 0 ≤ x ≤ 1;
// the smaller of the two is computed directly.
func betaCDFSurvival(x, alpha, beta float64) (cdf, survival float64, err error) {
	return betaSplit(x, 1-x, alpha, beta)
}

// betaSplit is betaCDFSurvival with the complement xc = 1 − x supplied by a caller that holds it more
// accurately than 1 − x can be formed.
func betaSplit(x, xc, alpha, beta float64) (cdf, survival float64, err error) {
	if x < alpha/(alpha+beta) {
		cdf, err = special.RegularizedBeta(x, alpha, beta)
		return cdf, 1 - cdf, wrapIf("incomplete beta", err)
	}
	survival, err = special.RegularizedBeta(xc, beta, alpha)
	return 1 - survival, survival, wrapIf("incomplete beta", err)
}

// BetaCDF returns P(X ≤ x) for the beta distribution with shapes α and β: the regularized incomplete beta
// function I_x(α, β) (num/special NUM-SPEC-005), 0 for x ≤ 0 and 1 for x ≥ 1.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrOutOfRange when α ≤ 0 or β ≤ 0, or when the
// continued fraction does not converge (wrapping special.ErrNoConvergence).
//
// ID:      NUM-PROB-074
// Time:    O(√max(α, β) + log(1/ε))   ε = 2⁻⁵²
// Space:   O(1), no allocation
// Source:  Johnson, Kotz, Balakrishnan, Continuous Univariate Distributions, vol. 2, 2nd ed., Wiley, 1995, ch. 25 "Beta Distributions", §1
// Source:  NIST DLMF, Release 1.2.1, §8.17.2 (I_x(a, b)), §8.17.4 (symmetry)
// Tags:    probability, distribution, beta, cdf, continuous, incomplete-beta
// Units:   dimensionless; result a probability in [0, 1]
// Example: α = 2, β = 3, x = 0.5 → 11/16 = 0.6875
func BetaCDF(x, alpha, beta float64) (float64, error) {
	if err := checkBetaShapes(alpha, beta); err != nil {
		return 0, err
	}
	if err := checkFinite("x", x); err != nil {
		return 0, err
	}
	if x <= 0 {
		return 0, nil
	}
	if x >= 1 {
		return 1, nil
	}
	c, _, err := betaCDFSurvival(x, alpha, beta)
	if err != nil {
		return 0, err
	}
	return c, nil
}

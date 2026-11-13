package prob

import "github.com/CrimsonCoder42/calcbench/num/special"

// gammaCDFSurvival returns P(X ≤ x) = P(k, x/θ) and P(X > x) = Q(k, x/θ) for validated arguments, x ≥ 0.
func gammaCDFSurvival(x, shape, scale float64) (cdf, survival float64, err error) {
	z := x / scale
	if z < shape+1 {
		cdf, err = special.RegularizedGammaP(shape, z)
		return cdf, 1 - cdf, wrapIf("gamma cdf", err)
	}
	survival, err = special.RegularizedGammaQ(shape, z)
	return 1 - survival, survival, wrapIf("gamma survival", err)
}

// wrapIf returns specialErr(what, err) when err is not nil.
func wrapIf(what string, err error) error {
	if err != nil {
		return specialErr(what, err)
	}
	return nil
}

// GammaCDF returns P(X ≤ x) for the gamma distribution with shape k and scale θ: the regularized lower
// incomplete gamma function P(k, x/θ) (num/special NUM-SPEC-001), and 0 for x ≤ 0.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrOutOfRange when shape ≤ 0 or scale ≤ 0, or when
// the incomplete gamma does not converge (wrapping special.ErrNoConvergence).
//
// ID:      NUM-PROB-070
// Time:    O(√k + log(1/ε))   k = shape, ε = 2⁻⁵² (the incomplete gamma's terms)
// Space:   O(1), no allocation
// Source:  Johnson, Kotz, Balakrishnan, Continuous Univariate Distributions, vol. 1, 2nd ed., Wiley, 1994, ch. 17 "Gamma Distributions", §1
// Source:  NIST DLMF, Release 1.2.1, §8.2.4 (P = γ/Γ), §8.7.1 (series), §8.9.2 (continued fraction)
// Tags:    probability, distribution, gamma, cdf, continuous, incomplete-gamma, erlang
// Units:   x and scale: one unit; shape dimensionless; result: a probability in [0, 1]
// Example: k = 1, θ = 1, x = 1 → 1 − e⁻¹ ≈ 0.6321205588285577 (DLMF 8.4.4)
func GammaCDF(x, shape, scale float64) (float64, error) {
	if err := checkGamma(shape, scale); err != nil {
		return 0, err
	}
	if err := checkFinite("x", x); err != nil {
		return 0, err
	}
	if x <= 0 {
		return 0, nil
	}
	c, _, err := gammaCDFSurvival(x, shape, scale)
	if err != nil {
		return 0, err
	}
	return c, nil
}

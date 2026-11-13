package prob

// ChiSquarePDF returns the chi-square density with k = df degrees of freedom,
// f(x) = x^(k/2 − 1)·e^(−x/2) / (2^(k/2)·Γ(k/2)) for x > 0 and 0 for x < 0: the gamma density with shape
// k/2 and scale 2, evaluated by GammaPDF (NUM-PROB-069) in log space. At x = 0 it is 0 for k > 2, ½ for
// k = 2, and unbounded (ErrOutOfRange) for k < 2.
// Errors: ErrNotFinite when x or df is NaN or ±Inf; ErrOutOfRange when df ≤ 0, or at x = 0 with df < 2;
// ErrOverflow when the density exceeds the float64 range.
//
// ID:      NUM-PROB-096
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Johnson, Kotz, Balakrishnan, Continuous Univariate Distributions, vol. 1, 2nd ed., Wiley, 1994, ch. 18 "Chi-Square Distributions Including Chi and Rayleigh", §1 (density)
// Tags:    probability, distribution, chi-square, pdf, continuous, goodness-of-fit
// Units:   x and df dimensionless; result: probability per unit of x
// Example: df = 2, x = 2 → ½e⁻¹ ≈ 0.18393972058572117
func ChiSquarePDF(x, df float64) (float64, error) {
	if err := checkPositive("df", df); err != nil {
		return 0, err
	}
	return GammaPDF(x, df/2, 2)
}

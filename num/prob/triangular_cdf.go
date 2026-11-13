package prob

// TriangularCDF returns P(X ≤ x) for the triangular distribution on [a, b] with mode c:
// (x − a)²/((b − a)(c − a)) for a ≤ x ≤ c and 1 − (b − x)²/((b − a)(b − c)) for c < x ≤ b; 0 below a and 1
// above b. Each product is formed as a product of ratios, so no intermediate overflows.
// Errors: as TriangularPDF.
//
// ID:      NUM-PROB-124
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Johnson, Kotz, Balakrishnan, Continuous Univariate Distributions, vol. 2, 2nd ed., Wiley, 1995, ch. 26 "Uniform (Rectangular) Distributions", §9 (triangular distributions)
// Tags:    probability, distribution, triangular, cdf, continuous, pert
// Units:   x, low, mode and high: one unit; result: a probability in [0, 1]
// Example: a = 0, c = 1, b = 2, x = 0.5 → 0.125
func TriangularCDF(x, low, mode, high float64) (float64, error) {
	w, err := checkTriangular(low, mode, high)
	if err != nil {
		return 0, err
	}
	if err := checkFinite("x", x); err != nil {
		return 0, err
	}
	switch {
	case x <= low:
		return 0, nil
	case x >= high:
		return 1, nil
	case x <= mode:
		return min(1, (x-low)/w*((x-low)/(mode-low))), nil
	default:
		return max(0, 1-(high-x)/w*((high-x)/(high-mode))), nil
	}
}

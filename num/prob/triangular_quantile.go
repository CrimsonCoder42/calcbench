package prob

import "math"

// TriangularQuantile returns the inverse of TriangularCDF: with F_c = (c − a)/(b − a),
// a + (b − a)·√(u·F_c) for u < F_c and b − (b − a)·√((1 − u)·(b − c)/(b − a)) otherwise. u = 0 returns a,
// u = 1 returns b.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrProbability when u is outside [0, 1]; as
// TriangularPDF otherwise.
//
// ID:      NUM-PROB-125
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Johnson, Kotz, Balakrishnan, Continuous Univariate Distributions, vol. 2, 2nd ed., Wiley, 1995, ch. 26 "Uniform (Rectangular) Distributions", §9 (triangular distributions, inverted)
// Tags:    probability, distribution, triangular, quantile, inverse-cdf, pert, monte-carlo
// Units:   u: a probability; low, mode, high and result: one unit
// Example: a = 0, c = 1, b = 2, u = 0.125 → 0.5
func TriangularQuantile(u, low, mode, high float64) (float64, error) {
	if err := checkProbability("u", u); err != nil {
		return 0, err
	}
	w, err := checkTriangular(low, mode, high)
	if err != nil {
		return 0, err
	}
	fc := (mode - low) / w
	if u < fc {
		return min(mode, low+w*math.Sqrt(u*fc)), nil
	}
	return max(mode, high-w*math.Sqrt((1-u)*((high-mode)/w))), nil
}

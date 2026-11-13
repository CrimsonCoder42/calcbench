package prob

import "github.com/CrimsonCoder42/calcbench/algo/reduce"

// UniformCDF returns P(X ≤ x) = (x − a)/(b − a) for X uniform on [a, b]: 0 below a and 1 above b.
// When b − a exceeds the float64 range, the quotient is taken of the halved terms (x/2 − a/2)/(b/2 − a/2).
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrOutOfRange when a ≥ b.
//
// ID:      NUM-PROB-045
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  NIST/SEMATECH, e-Handbook of Statistical Methods, §1.3.6.6.2 "Uniform Distribution", https://www.itl.nist.gov/div898/handbook/eda/section3/eda3662.htm (read 2026-11-12)
// Source:  Ross, A First Course in Probability, 10th ed., Pearson, 2019, §5.3 (uniform random variable)
// Tags:    probability, distribution, uniform, cdf, continuous
// Units:   x, a, b: one unit; result: a probability in [0, 1]
// Example: a = 0, b = 10, x = 3 → 0.3
func UniformCDF(x, a, b float64) (float64, error) {
	w, err := checkUniform(a, b)
	if err != nil {
		return 0, err
	}
	if err := checkFinite("x", x); err != nil {
		return 0, err
	}
	switch {
	case x <= a:
		return 0, nil
	case x >= b:
		return 1, nil
	}
	if !reduce.AllFinite(w) {
		return min(1, (x/2-a/2)/(b/2-a/2)), nil // halving is exact for |a|, |b| ≥ 2¹⁰²³
	}
	return min(1, (x-a)/w), nil // a < x < b, so x − a ≤ b − a is finite
}

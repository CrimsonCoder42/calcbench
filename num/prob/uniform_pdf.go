package prob

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// checkUniform validates a and b finite with a < b and returns the width b − a. The width is +Inf when
// b − a exceeds the float64 range (then a < 0 < b and both |a|, |b| ≥ 2¹⁰²³); callers switch to the
// halved endpoints a/2, b/2 (exact there), whose difference is always finite.
func checkUniform(a, b float64) (float64, error) {
	if err := checkFinite("a", a); err != nil {
		return 0, err
	}
	if err := checkFinite("b", b); err != nil {
		return 0, err
	}
	if a >= b {
		return 0, fmt.Errorf("%w: a = %v >= b = %v", ErrOutOfRange, a, b)
	}
	return b - a, nil
}

// UniformPDF returns the density 1/(b − a) of the continuous uniform distribution on [a, b] at x, and 0
// outside [a, b]. When b − a exceeds the float64 range the density is computed as ½/(b/2 − a/2); it is
// then at least 2⁻¹⁰²⁵ (about 2.8e-309, subnormal) and never underflows to 0, so every finite a < b gives
// a nonzero density.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrOutOfRange when a ≥ b; ErrOverflow when the
// density itself exceeds the float64 range (b − a below about 5.6e-309, subnormal).
//
// ID:      NUM-PROB-044
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  NIST/SEMATECH, e-Handbook of Statistical Methods, §1.3.6.6.2 "Uniform Distribution", https://www.itl.nist.gov/div898/handbook/eda/section3/eda3662.htm (read 2026-11-12)
// Source:  Ross, A First Course in Probability, 10th ed., Pearson, 2019, §5.3 (uniform random variable)
// Tags:    probability, distribution, uniform, pdf, continuous
// Units:   x, a, b: one unit; result: probability per that unit
// Example: a = 0, b = 10, x = 3 → 0.1
func UniformPDF(x, a, b float64) (float64, error) {
	w, err := checkUniform(a, b)
	if err != nil {
		return 0, err
	}
	if err := checkFinite("x", x); err != nil {
		return 0, err
	}
	if x < a || x > b {
		return 0, nil
	}
	if !reduce.AllFinite(w) {
		return 0.5 / (b/2 - a/2), nil
	}
	d := 1 / w
	if !reduce.AllFinite(d) {
		return 0, fmt.Errorf("%w: 1/(b − a) for b − a = %v", ErrOverflow, w)
	}
	return d, nil
}

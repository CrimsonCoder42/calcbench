package prob

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// checkLocationScale validates a location (finite) and a scale (finite, > 0).
func checkLocationScale(location, scale float64) error {
	if err := checkFinite("location", location); err != nil {
		return err
	}
	return checkPositive("scale", scale)
}

// checkQuantileOpen validates u in [0, 1] (ErrNotFinite, ErrProbability) and rejects u = 0 and u = 1 with
// ErrOutOfRange, for a distribution whose support is unbounded on both sides.
func checkQuantileOpen(u float64) error {
	if err := checkProbability("u", u); err != nil {
		return err
	}
	if u == 0 || u == 1 {
		return fmt.Errorf("%w: u = %v, quantile infinite", ErrOutOfRange, u)
	}
	return nil
}

// GumbelPDF returns the Gumbel (type I extreme value, maximum) density
// f(x) = (1/β)·exp(−z − e^(−z)), z = (x − μ)/β, with μ the location and β the scale.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrOutOfRange when scale ≤ 0; ErrOverflow when the
// density exceeds the float64 range (scale below about 10⁻³⁰⁸).
//
// ID:      NUM-PROB-103
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Johnson, Kotz, Balakrishnan, Continuous Univariate Distributions, vol. 2, 2nd ed., Wiley, 1995, ch. 22 "Extreme Value Distributions", §1 (type 1, eq. 22.1)
// Tags:    probability, distribution, gumbel, extreme-value, pdf, continuous, flood, maxima
// Units:   x, location and scale: one unit; result: probability per that unit
// Example: μ = 0, β = 1, x = 0 → e⁻¹ ≈ 0.36787944117144233
func GumbelPDF(x, location, scale float64) (float64, error) {
	if err := checkLocationScale(location, scale); err != nil {
		return 0, err
	}
	if err := checkFinite("x", x); err != nil {
		return 0, err
	}
	z := (x - location) / scale
	if !reduce.AllFinite(z) {
		return 0, nil
	}
	d := math.Exp(-z-math.Exp(-z)) / scale
	if !reduce.AllFinite(d) {
		return 0, fmt.Errorf("%w: density at x = %v, scale = %v", ErrOverflow, x, scale)
	}
	return d, nil
}

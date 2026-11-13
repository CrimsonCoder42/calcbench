package prob

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// LogisticPDF returns the logistic density f(x) = e^(−z)/(s·(1 + e^(−z))²), z = (x − μ)/s, μ the location
// and s the scale, evaluated with e^(−|z|) (the density is symmetric) so it never overflows inside.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrOutOfRange when scale ≤ 0; ErrOverflow when the
// density exceeds the float64 range.
//
// ID:      NUM-PROB-111
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Johnson, Kotz, Balakrishnan, Continuous Univariate Distributions, vol. 2, 2nd ed., Wiley, 1995, ch. 23 "Logistic Distribution", §1 (definition)
// Tags:    probability, distribution, logistic, sigmoid, pdf, continuous, elo
// Units:   x, location and scale: one unit; result: probability per that unit
// Example: μ = 0, s = 1, x = 0 → ¼
func LogisticPDF(x, location, scale float64) (float64, error) {
	if err := checkLocationScale(location, scale); err != nil {
		return 0, err
	}
	if err := checkFinite("x", x); err != nil {
		return 0, err
	}
	e := math.Exp(-math.Abs((x - location) / scale))
	d := e / (1 + e) / (1 + e) / scale
	if !reduce.AllFinite(d) {
		return 0, fmt.Errorf("%w: density at x = %v, scale = %v", ErrOverflow, x, scale)
	}
	return d, nil
}

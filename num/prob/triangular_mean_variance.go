package prob

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// TriangularMeanVariance returns the mean (a + b + c)/3 and the variance (a² + b² + c² − ab − ac − bc)/18 of
// the triangular distribution on [a, b] with mode c; the variance is formed as (w² + v² − w·v)/18 with
// w = b − a, v = c − a (shift-free, no cancellation of large squares).
// Errors: as TriangularPDF; ErrOverflow when the mean or the variance exceeds the float64 range.
//
// ID:      NUM-PROB-126
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Johnson, Kotz, Balakrishnan, Continuous Univariate Distributions, vol. 2, 2nd ed., Wiley, 1995, ch. 26 "Uniform (Rectangular) Distributions", §9 (triangular distributions, moments)
// Tags:    probability, distribution, triangular, mean, variance, pert
// Units:   low, mode, high and mean: one unit; variance: that unit squared
// Example: a = 0, c = 1, b = 2 → (1, 1/6)
func TriangularMeanVariance(low, mode, high float64) (mean, variance float64, err error) {
	w, err := checkTriangular(low, mode, high)
	if err != nil {
		return 0, 0, err
	}
	mean = low/3 + mode/3 + high/3
	v := mode - low
	variance = w*(w/18) + v*((v-w)/18)
	if !reduce.AllFinite(mean) || !reduce.AllFinite(variance) {
		return 0, 0, fmt.Errorf("%w: low = %v, mode = %v, high = %v", ErrOverflow, low, mode, high)
	}
	return mean, variance, nil
}

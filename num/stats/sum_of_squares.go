package stats

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// SumOfSquares returns the sum of squared deviations of values about their mean, Σ(x − x̄)², the
// numerator of the variance. It uses the corrected two-pass formula Σd² − (Σd)²/n (d = x − x̄) with
// compensated sums, never the one-pass Σx² − (Σx)²/n, which cancels catastrophically.
// values is read, never modified.
// A result below the smallest float64 rounds to a subnormal or 0 like any float64 result (the spread
// itself is computed on rescaled deviations and is never lost); for such data use the standard deviation.
// Errors: ErrEmpty when values is empty; ErrNotFinite when a value is NaN or ±Inf; ErrOverflow when the
// sum exceeds the float64 range.
//
// ID:      NUM-STATS-025
// Time:    O(n)   n = len(values); two passes (four only when a plain sum overflows or every deviation is below 2^-480)
// Space:   O(1) extra, no allocation
// Source:  NIST, Dataplot Reference Manual, "SUM OF SQUARES FROM MEAN", https://www.itl.nist.gov/div898/software/dataplot/refman2/auxillar/ssqmean.htm (read 2026-11-12)
// Source:  Chan, Golub, LeVeque, "Algorithms for computing the sample variance: analysis and recommendations", The American Statistician, 1983, 37(3):242–247, doi:10.1080/00031305.1983.10483115 (the corrected two-pass algorithm)
// Tags:    sum-of-squares, deviation, dispersion, variance, descriptive, two-pass
// Units:   values in any one unit; result in that unit squared
// Example: NIST StRD Univariate NumAcc1 → 2, that is (n − 1)·s² with the certified s = 1
func SumOfSquares(values []float64) (float64, error) {
	if err := checkSample(values, 1); err != nil {
		return 0, err
	}
	ss, e := centredSquares(values, meanOf(values))
	r := math.Ldexp(ss, 2*e)
	if !reduce.AllFinite(r) {
		return 0, fmt.Errorf("%w: sum of squares of values", ErrOverflow)
	}
	return r, nil
}

// checkSample returns ErrEmpty for no values, ErrTooFew for fewer than minLen, and ErrNotFinite for a
// NaN or ±Inf value.
func checkSample(values []float64, minLen int) error {
	n := len(values)
	if n == 0 {
		return fmt.Errorf("%w: values", ErrEmpty)
	}
	if n < minLen {
		return fmt.Errorf("%w: len(values) = %d, want ≥ %d", ErrTooFew, n, minLen)
	}
	return checkFinite("values", values)
}

// varianceOf returns Σ(x − x̄)²/divisor, or ErrOverflow. values are checked, finite, non-empty.
func varianceOf(values []float64, divisor float64) (float64, error) {
	ss, e := centredSquares(values, meanOf(values))
	r := math.Ldexp(ss/divisor, 2*e)
	if !reduce.AllFinite(r) {
		return 0, fmt.Errorf("%w: variance of values", ErrOverflow)
	}
	return r, nil
}

// standardDeviationOf returns √(Σ(x − x̄)²/divisor), or ErrOverflow. It stays finite when only the
// variance would overflow. values are checked, finite, non-empty.
func standardDeviationOf(values []float64, divisor float64) (float64, error) {
	sd, e := scaledStandardDeviation(values, divisor)
	r := math.Ldexp(sd, e)
	if !reduce.AllFinite(r) {
		return 0, fmt.Errorf("%w: standard deviation of values", ErrOverflow)
	}
	return r, nil
}

// scaledStandardDeviation returns sd and e with √(Σ(x − x̄)²/divisor) = sd·2^e.
func scaledStandardDeviation(values []float64, divisor float64) (sd float64, e int) {
	ss, e := centredSquares(values, meanOf(values))
	return math.Sqrt(ss / divisor), e
}

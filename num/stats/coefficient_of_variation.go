package stats

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// CoefficientOfVariation returns s/x̄, the sample standard deviation (NUM-STATS-019) divided by the mean
// (NUM-STATS-001). It keeps the mean's sign; it is meaningful for data on a ratio scale (a unit with a
// true zero). For 100·s/|x̄| use PercentRelativeStandardDeviation. values is read, never modified.
// Errors: ErrEmpty when values is empty; ErrTooFew when it has one value; ErrNotFinite when a value is
// NaN or ±Inf; ErrZeroMean when the mean is 0; ErrOverflow when the ratio exceeds the float64 range.
//
// ID:      NUM-STATS-003
// Time:    O(n)   n = len(values); three passes (mean, then the centred sums)
// Space:   O(1) extra, no allocation
// Source:  NIST, Dataplot Reference Manual, "COEFFICIENT OF VARIATION", https://www.itl.nist.gov/div898/software/dataplot/refman2/auxillar/coefvari.htm (read 2026-11-12)
// Tags:    coefficient-of-variation, cv, relative-dispersion, dispersion, descriptive
// Units:   values in any one ratio-scale unit; result dimensionless
// Example: NIST StRD Univariate Michelso → 0.0790105478190518/299.8524 = 2.6349813e-04 (from the certified s and mean)
func CoefficientOfVariation(values []float64) (float64, error) {
	return sdOverMean(values)
}

// sdOverMean returns s/x̄ for CoefficientOfVariation and PercentRelativeStandardDeviation.
func sdOverMean(values []float64) (float64, error) {
	if err := checkSample(values, 2); err != nil {
		return 0, err
	}
	mean := meanOf(values)
	if mean == 0 {
		return 0, fmt.Errorf("%w: mean of values", ErrZeroMean)
	}
	ss, e := centredSquares(values, mean)
	sd := math.Sqrt(ss / float64(len(values)-1)) // s = sd·2^e
	var r float64
	if e > 0 {
		r = math.Ldexp(sd/mean, e) // overflowing spread: sd/mean ≤ the result, so it cannot overflow first
	} else {
		// Small spread (e ≤ 0): scale the mean instead. |mean| ≤ 2^54·max|d| when a deviation is
		// nonzero, so mean·2^-e is finite; sd·2^e alone could be subnormal and lose bits.
		r = sd / math.Ldexp(mean, -e)
	}
	if !reduce.AllFinite(r) {
		return 0, fmt.Errorf("%w: standard deviation / mean of values", ErrOverflow)
	}
	return r, nil
}

package stats

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// PercentRelativeStandardDeviation returns the percent relative standard deviation (%RSD) of values,
// 100·s/|x̄|, with s the sample standard deviation (NUM-STATS-019) and x̄ the mean (NUM-STATS-001).
// It is never negative; for the signed ratio s/x̄ use CoefficientOfVariation. values is read, never modified.
// Errors: ErrEmpty when values is empty; ErrTooFew when it has one value; ErrNotFinite when a value is
// NaN or ±Inf; ErrZeroMean when the mean is 0; ErrOverflow when the result exceeds the float64 range.
//
// ID:      NUM-STATS-030
// Time:    O(n)   n = len(values); three passes (mean, then the centred sums)
// Space:   O(1) extra, no allocation
// Source:  NIST, Dataplot Reference Manual, "RELATIVE STANDARD DEVIATION", https://www.itl.nist.gov/div898/software/dataplot/refman2/ch2/relat_sd.pdf (read 2026-11-12)
// Source:  NIST, Dataplot Reference Manual, "COEFFICIENT OF VARIATION", https://www.itl.nist.gov/div898/software/dataplot/refman2/auxillar/coefvari.htm (read 2026-11-12) (relsd, the percent relative standard deviation)
// Tags:    relative-standard-deviation, rsd, percent-rsd, coefficient-of-variation, dispersion, descriptive
// Units:   values in any one ratio-scale unit; result in percent
// Example: NIST StRD Univariate Michelso → 100·0.0790105478190518/299.8524 = 0.026349813 % (from the certified s and mean)
func PercentRelativeStandardDeviation(values []float64) (float64, error) {
	r, err := sdOverMean(values)
	if err != nil {
		return 0, err
	}
	p := 100 * math.Abs(r)
	if !reduce.AllFinite(p) {
		return 0, fmt.Errorf("%w: percent relative standard deviation of values", ErrOverflow)
	}
	return p, nil
}

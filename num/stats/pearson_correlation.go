package stats

import (
	"fmt"
	"math"
)

// PearsonCorrelation returns the Pearson product-moment correlation coefficient of the paired samples
// x and y, r = S_xy/(√S_xx·√S_yy) with S_xx = Σ(x − x̄)², S_yy = Σ(y − ȳ)², S_xy = Σ(x − x̄)(y − ȳ).
// The result is clamped to [−1, 1] against rounding. No significance test. x and y are read, never modified.
// Errors: ErrEmpty when x is empty; ErrLengthMismatch when len(y) ≠ len(x); ErrTooFew when there is one
// pair; ErrNotFinite when a value is NaN or ±Inf; ErrZeroVariance when every x, or every y, is equal.
//
// ID:      NUM-STATS-014
// Time:    O(n)   n = len(x); three passes (means, then the centred sums)
// Space:   O(1) extra, no allocation
// Source:  NIST, Dataplot Reference Manual, "CORRELATION", https://www.itl.nist.gov/div898/software/dataplot/refman2/auxillar/correlat.htm (read 2026-11-12)
// Tags:    correlation, pearson, product-moment, association, descriptive
// Units:   x and y in any units; result dimensionless in [−1, 1]
// Example: NIST StRD Norris → r = +√R² = 0.999996872936967, from the certified R² = 0.999993745883712 (hand-checked)
func PearsonCorrelation(x, y []float64) (float64, error) {
	if err := checkPaired(x, y); err != nil {
		return 0, err
	}
	sxx, syy, sxy, _, _ := centredProducts(x, y, meanOf(x), meanOf(y))
	if sxx == 0 {
		return 0, fmt.Errorf("%w: every x equals %v", ErrZeroVariance, x[0])
	}
	if syy == 0 {
		return 0, fmt.Errorf("%w: every y equals %v", ErrZeroVariance, y[0])
	}
	r := sxy / (math.Sqrt(sxx) * math.Sqrt(syy)) // the scales 2^(ex+ey)/(2^ex·2^ey) cancel
	return min(max(r, -1), 1), nil
}

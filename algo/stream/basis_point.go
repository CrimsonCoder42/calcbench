package stream

import (
	"fmt"
	"math"
)

// basisPointsPerUnit is the definition of the basis point: 1 bp = 1/10 000 = 0.01%.
const basisPointsPerUnit = 10_000

// BasisPointsToFraction returns bp/10 000, the fraction (0.01 = 1%) that bp basis points express. The
// division is one correctly rounded operation, so 25 bp gives the float64 nearest 0.0025.
// Errors: ErrNotFinite when bp is NaN or ±Inf.
//
// ID:      ALG-STREAM-015
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Hull, Options, Futures, and Other Derivatives, 10th ed., Pearson, 2018, Glossary of Terms, "Basis Point" (one hundredth of one percent, 0.01%)
// Tags:    basis-point, bp, bps, fraction, percent, rate, spread, unit-conversion, trading
// Units:   bp in basis points; result a fraction (dimensionless)
// Example: BasisPointsToFraction(25) = 0.0025
func BasisPointsToFraction(bp float64) (float64, error) {
	if math.IsNaN(bp) || math.IsInf(bp, 0) {
		return 0, fmt.Errorf("%w: bp = %v", ErrNotFinite, bp)
	}
	return bp / basisPointsPerUnit, nil
}

// FractionToBasisPoints returns fraction × 10 000, the basis points that a fraction (0.01 = 1%)
// expresses, as one correctly rounded product.
// Errors: ErrNotFinite when fraction is NaN or ±Inf; ErrOverflow when the result exceeds the float64
// range (|fraction| above about 1.8e304).
//
// ID:      ALG-STREAM-016
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Hull, Options, Futures, and Other Derivatives, 10th ed., Pearson, 2018, Glossary of Terms, "Basis Point" (one hundredth of one percent, 0.01%)
// Tags:    basis-point, bp, bps, fraction, percent, rate, spread, unit-conversion, trading
// Units:   fraction dimensionless; result in basis points
// Example: FractionToBasisPoints(0.0025) = 25
func FractionToBasisPoints(fraction float64) (float64, error) {
	if math.IsNaN(fraction) || math.IsInf(fraction, 0) {
		return 0, fmt.Errorf("%w: fraction = %v", ErrNotFinite, fraction)
	}
	bp := fraction * basisPointsPerUnit
	if math.IsInf(bp, 0) {
		return 0, fmt.Errorf("%w: fraction = %v", ErrOverflow, fraction)
	}
	return bp, nil
}

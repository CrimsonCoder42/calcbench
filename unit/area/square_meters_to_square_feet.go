package area

import (
	"fmt"
	"math"
)

// SquareMetersToSquareFeet returns the area in square feet of an area in square meters:
// ft² = m² / 0.09290304.
// Errors: ErrNotFinite for NaN or ±Inf; ErrNegative below 0; ErrOverflow when the result is beyond float64.
//
// ID:      UNT-AREA-002
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Thompson A, Taylor BN. Guide for the Use of the International System of Units (SI). NIST Special Publication 811, 2008 ed., Appendix B.8, "foot squared (ft²) to meter squared (m²): 9.290 304 E−02" (exact), inverted
// Tags:    unit, area, square-meters, square-feet, conversion
// Units:   squareMeters m²; result ft²
// Range:   squareMeters ≥ 0
// Example: 100 m² → 1076.391041670972 ft² (12.5 m × 8 m floor)
func SquareMetersToSquareFeet(squareMeters float64) (float64, error) {
	if err := checkArea("squareMeters", squareMeters); err != nil {
		return 0, err
	}
	f := squareMeters / squareMetersPerSquareFoot
	if math.IsInf(f, 0) {
		return 0, fmt.Errorf("%w: squareMeters = %v", ErrOverflow, squareMeters)
	}
	return f, nil
}

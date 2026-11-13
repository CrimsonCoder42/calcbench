package area

import (
	"errors"
	"fmt"
	"math"
)

var (
	// ErrNotFinite reports a NaN or infinite input.
	ErrNotFinite = errors.New("area: value is NaN or infinite")
	// ErrNegative reports a negative area.
	ErrNegative = errors.New("area: value is negative")
	// ErrOverflow reports a result beyond float64.
	ErrOverflow = errors.New("area: result overflows float64")
)

// squareMetersPerSquareFoot is exact: (0.3048 m)² (NIST SP 811, B.8, "foot squared").
const squareMetersPerSquareFoot = 0.09290304

func checkArea(name string, v float64) error {
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return fmt.Errorf("%w: %s = %v", ErrNotFinite, name, v)
	}
	if v < 0 {
		return fmt.Errorf("%w: %s = %v", ErrNegative, name, v)
	}
	return nil
}

// SquareFeetToSquareMeters returns the area in square meters of an area in square feet: m² = ft² × 0.09290304.
// Errors: ErrNotFinite for NaN or ±Inf; ErrNegative below 0.
//
// ID:      UNT-AREA-001
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Thompson A, Taylor BN. Guide for the Use of the International System of Units (SI). NIST Special Publication 811, 2008 ed., Appendix B.8, "foot squared (ft²) to meter squared (m²): 9.290 304 E−02" (exact)
// Tags:    unit, area, square-feet, square-meters, conversion
// Units:   squareFeet ft²; result m²
// Range:   squareFeet ≥ 0
// Example: 1076.391041670972 ft² → 100 m²; 1 ft² → 0.09290304 m²
func SquareFeetToSquareMeters(squareFeet float64) (float64, error) {
	if err := checkArea("squareFeet", squareFeet); err != nil {
		return 0, err
	}
	return squareFeet * squareMetersPerSquareFoot, nil
}

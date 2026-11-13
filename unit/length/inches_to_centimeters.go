package length

import (
	"fmt"
	"math"
)

// InchesToCentimeters returns the length in centimeters of a length in inches: l/cm = 2.54 · l/in.
// Errors: ErrNotFinite for NaN or ±Inf; ErrNegative below 0 in; ErrOverflow above MaxFloat64/2.54 in.
//
// ID:      UNT-LEN-002
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Thompson A, Taylor BN. Guide for the Use of the International System of Units (SI). NIST Special Publication 811, 2008 ed., Appendix B.8 "Factors for units listed alphabetically", row "inch (in) to centimeter (cm) 2.54 E+00" (boldface: exact)
// Tags:    unit, length, height, inch, centimeter, conversion
// Units:   inches in; result cm
// Range:   inches ≥ 0
// Example: 1 in → 2.54 cm; 70 in → 177.8 cm; 12 in → 30.48 cm
func InchesToCentimeters(inches float64) (float64, error) {
	if err := checkLength("inches", inches); err != nil {
		return 0, err
	}
	cm := inches * centimetersPerInch
	if math.IsInf(cm, 0) {
		return 0, fmt.Errorf("%w: inches = %v", ErrOverflow, inches)
	}
	return cm, nil
}

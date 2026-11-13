package length

import (
	"errors"
	"fmt"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

var (
	// ErrNotFinite is returned when an input is NaN or ±Inf.
	ErrNotFinite = errors.New("length: value is NaN or infinite")
	// ErrNegative is returned when a length is below zero.
	ErrNegative = errors.New("length: value is negative")
	// ErrOverflow is returned when a finite input gives a result beyond the float64 range.
	ErrOverflow = errors.New("length: result overflows float64")
)

// centimetersPerInch is the international inch in centimeters, 2.54 exactly (NIST SP 811, 2008, App. B.8,
// "inch (in) to centimeter (cm) 2.54 E+00", boldface = exact).
const centimetersPerInch = 2.54

// checkLength returns ErrNotFinite for NaN or ±Inf and ErrNegative below zero, naming the parameter.
func checkLength(name string, v float64) error {
	if !reduce.AllFinite(v) {
		return fmt.Errorf("%w: %s = %v", ErrNotFinite, name, v)
	}
	if v < 0 {
		return fmt.Errorf("%w: %s = %v", ErrNegative, name, v)
	}
	return nil
}

// CentimetersToInches returns the length in inches of a length in centimeters: l/in = (l/cm)/2.54.
// Errors: ErrNotFinite for NaN or ±Inf; ErrNegative below 0 cm.
//
// ID:      UNT-LEN-001
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Thompson A, Taylor BN. Guide for the Use of the International System of Units (SI). NIST Special Publication 811, 2008 ed., Appendix B.8 "Factors for units listed alphabetically", row "inch (in) to centimeter (cm) 2.54 E+00" (boldface: exact)
// Tags:    unit, length, height, centimeter, inch, conversion
// Units:   centimeters cm; result in
// Range:   centimeters ≥ 0
// Example: 2.54 cm → 1 in; 180 cm → 70.86614173228347 in; 30.48 cm → 12 in
func CentimetersToInches(centimeters float64) (float64, error) {
	if err := checkLength("centimeters", centimeters); err != nil {
		return 0, err
	}
	// The factor is above 1, so a finite input gives a finite result.
	return centimeters / centimetersPerInch, nil
}

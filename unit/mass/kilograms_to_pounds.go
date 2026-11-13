package mass

import (
	"errors"
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

var (
	// ErrNotFinite is returned when an input is NaN or ±Inf.
	ErrNotFinite = errors.New("mass: value is NaN or infinite")
	// ErrNegative is returned when a mass is below zero.
	ErrNegative = errors.New("mass: value is negative")
	// ErrOverflow is returned when a finite input gives a result beyond the float64 range.
	ErrOverflow = errors.New("mass: result overflows float64")
)

// kilogramsPerPound is the international avoirdupois pound in kilograms, 0.453 592 37 exactly by
// definition (NIST SP 811, 2008, App. B.8, row "pound (avoirdupois) (lb)", factor 4.535 924 E-01, and
// footnote 22: "The exact conversion factor is 4.535 923 7 E-01").
const kilogramsPerPound = 0.45359237

// checkMass returns ErrNotFinite for NaN or ±Inf and ErrNegative below zero, naming the parameter.
func checkMass(name string, v float64) error {
	if !reduce.AllFinite(v) {
		return fmt.Errorf("%w: %s = %v", ErrNotFinite, name, v)
	}
	if v < 0 {
		return fmt.Errorf("%w: %s = %v", ErrNegative, name, v)
	}
	return nil
}

// KilogramsToPounds returns the mass in pounds (avoirdupois) of a mass in kilograms:
// m/lb = (m/kg)/0.453 592 37.
// Errors: ErrNotFinite for NaN or ±Inf; ErrNegative below 0 kg; ErrOverflow above 0.453 592 37·MaxFloat64 kg.
//
// ID:      UNT-MASS-001
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Thompson A, Taylor BN. Guide for the Use of the International System of Units (SI). NIST Special Publication 811, 2008 ed., Appendix B.8 "Factors for units listed alphabetically", row "pound (avoirdupois) (lb) to kilogram (kg) 4.535 924 E-01" and footnote 22: "The exact conversion factor is 4.535 923 7 E-01"
// Tags:    unit, mass, weight, kilogram, pound, avoirdupois, conversion
// Units:   kilograms kg; result lb
// Range:   kilograms ≥ 0
// Example: 0.45359237 kg → 1 lb; 1 kg → 2.2046226218487757 lb; 100 kg → 220.46226218487757 lb
func KilogramsToPounds(kilograms float64) (float64, error) {
	if err := checkMass("kilograms", kilograms); err != nil {
		return 0, err
	}
	// The result exceeds MaxFloat64 only for kilograms above 0.453 592 37·MaxFloat64.
	lb := kilograms / kilogramsPerPound
	if math.IsInf(lb, 0) {
		return 0, fmt.Errorf("%w: kilograms = %v", ErrOverflow, kilograms)
	}
	return lb, nil
}

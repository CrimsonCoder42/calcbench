package cardio

import (
	"errors"
	"fmt"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// The package's sentinel errors. They live in the file of the first function that uses them
// (doc/ASSUMPTIONS.md A-008). Every error wraps one of these with the offending parameter; test with errors.Is.
var (
	// ErrNotFinite: an input is NaN or ±Inf.
	ErrNotFinite = errors.New("cardio: value is NaN or infinite")
	// ErrOutOfRange: an input is outside the range the source validated.
	ErrOutOfRange = errors.New("cardio: input outside the range the source validated")
	// ErrDomain: the inputs are valid one by one but the result is beyond float64 for them together.
	ErrDomain = errors.New("cardio: inputs valid one by one but the formula is undefined for them together")
)

// checkFinite returns ErrNotFinite for NaN or ±Inf.
func checkFinite(name string, v float64) error {
	if !reduce.AllFinite(v) {
		return fmt.Errorf("%w: %s = %v", ErrNotFinite, name, v)
	}
	return nil
}

// checkRange returns ErrOutOfRange unless lo ≤ v ≤ hi.
func checkRange(name string, v, lo, hi float64) error {
	if v < lo || v > hi {
		return fmt.Errorf("%w: %s = %v, want %v ≤ %s ≤ %v", ErrOutOfRange, name, v, lo, name, hi)
	}
	return nil
}

// checkPositive returns ErrOutOfRange unless v > 0.
func checkPositive(name string, v float64) error {
	if !(v > 0) {
		return fmt.Errorf("%w: %s = %v, want > 0", ErrOutOfRange, name, v)
	}
	return nil
}

// MaxHeartRateTanaka returns the age-predicted maximal heart rate, in beats/min, by Tanaka's equation
// HRmax = 208 − 0.7 × age.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range.
//
// ID:      MED-CARDIO-001
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Tanaka H, Monahan KD, Seals DR. Age-predicted maximal heart rate revisited. J Am Coll Cardiol. 2001;37(1):153-156. doi:10.1016/S0735-1097(00)01054-8. PMID 11153730. Abstract: "HRmax was strongly related to age (r = −0.90), using the equation of 208 − 0.7 × age"
// Tags:    cardio, heart-rate, maximal-heart-rate, exercise, tanaka
// Units:   ageYears years (decimal allowed); result beats/min
// Range:   ageYears 18–81 (the laboratory study's 514 subjects, ages 18 to 81; the meta-analysis pooled 18,712 adults in 351 studies)
// Example: Tanaka 2001 (full text): at 70 years, 208 − 0.7 × 70 = 159 beats/min
func MaxHeartRateTanaka(ageYears float64) (float64, error) {
	if err := checkFinite("ageYears", ageYears); err != nil {
		return 0, err
	}
	if err := checkRange("ageYears", ageYears, 18, 81); err != nil {
		return 0, err
	}
	return 208 - 0.7*ageYears, nil
}

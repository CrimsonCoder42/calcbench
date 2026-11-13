package anthropometry

import (
	"errors"
	"fmt"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

var (
	// ErrNotFinite is returned when an input is NaN or ±Inf; it is checked before any range.
	ErrNotFinite = errors.New("anthropometry: value is NaN or infinite")
	// ErrOutOfRange is returned when an input is outside the range the source validated.
	ErrOutOfRange = errors.New("anthropometry: input outside the range the source validated")
	// ErrDomain is returned when inputs valid one by one give a result beyond float64 together.
	ErrDomain = errors.New("anthropometry: inputs valid one by one but the formula is undefined for them together")
)

// checkFinite returns ErrNotFinite, wrapped with name, when v is NaN or ±Inf.
func checkFinite(name string, v float64) error {
	if !reduce.AllFinite(v) {
		return fmt.Errorf("%w: %s = %v", ErrNotFinite, name, v)
	}
	return nil
}

// checkPositive returns ErrOutOfRange, wrapped with name, unless v > 0.
func checkPositive(name string, v float64) error {
	if !(v > 0) {
		return fmt.Errorf("%w: %s = %v, want > 0", ErrOutOfRange, name, v)
	}
	return nil
}

// checkResult returns ErrDomain when a result computed from finite inputs is not finite.
func checkResult(v float64) error {
	if !reduce.AllFinite(v) {
		return fmt.Errorf("%w: result is beyond float64", ErrDomain)
	}
	return nil
}

// BodyMassIndex returns the body mass index, in kg/m², from weight in kilograms and height in metres:
// weightKg / heightM².
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range; ErrDomain when the result is beyond
// float64 (a height so small that its square underflows, or a weight so large that the quotient overflows).
//
// ID:      MED-ANTHRO-001
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Centers for Disease Control and Prevention. BMI Frequently Asked Questions. https://www.cdc.gov/bmi/faq/index.html (last reviewed 2024-06-28; read 2026-11-12): "weight (kg)/[height (m)]2"
// Tags:    anthropometry, bmi, body-mass-index, metric
// Units:   weightKg kg; heightM m; result kg/m²
// Range:   weightKg > 0; heightM > 0 (the source states no limits)
// Example: 70 kg, 1.75 m → 70 / 3.0625 = 22.857142857142858 (HAND-DERIVED)
func BodyMassIndex(weightKg, heightM float64) (float64, error) {
	if err := checkFinite("weightKg", weightKg); err != nil {
		return 0, err
	}
	if err := checkFinite("heightM", heightM); err != nil {
		return 0, err
	}
	if err := checkPositive("weightKg", weightKg); err != nil {
		return 0, err
	}
	if err := checkPositive("heightM", heightM); err != nil {
		return 0, err
	}
	r := weightKg / (heightM * heightM)
	if err := checkResult(r); err != nil {
		return 0, err
	}
	return r, nil
}

package algebra

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// Sign returns the signum of x: −1 for x < 0, 0 for x = 0 (either signed zero), +1 for x > 0.
// Errors: ErrNotFinite for NaN (±Inf give ∓1/±1).
//
// ID:      NUM-ALGEB-025
// Time:    O(1)
// Space:   O(1)
// Source:  NIST Digital Library of Mathematical Functions, §1.16(iv) (the signum function sign x)
// Tags:    sign, signum, scalar, algebra
// Units:   x in any unit; the result dimensionless
// Example: −2.5 → −1; 0 → 0; 3 → 1 (hand-checked)
func Sign(x float64) (float64, error) {
	switch {
	case x > 0:
		return 1, nil
	case x < 0:
		return -1, nil
	case x == 0:
		return 0, nil
	}
	return 0, fmt.Errorf("%w: x = NaN", ErrNotFinite)
}

// Heaviside returns the Heaviside step H(x): 0 for x < 0, 1 for x > 0, and valueAtZero at x = 0 (DLMF
// leaves H(0) to the user; ½ is common).
// Errors: ErrNotFinite for NaN x or a NaN or ±Inf valueAtZero (±Inf x give 0 or 1).
//
// ID:      NUM-ALGEB-026
// Time:    O(1)
// Space:   O(1)
// Source:  NIST Digital Library of Mathematical Functions, §1.16(iv), eq. 1.16.13 (the Heaviside function)
// Tags:    heaviside, step-function, unit-step, scalar, algebra
// Units:   x in any unit; the result dimensionless
// Example: −1 → 0; 2 → 1; 0 with valueAtZero ½ → ½ (hand-checked)
func Heaviside(x, valueAtZero float64) (float64, error) {
	if !reduce.AllFinite(valueAtZero) {
		return 0, fmt.Errorf("%w: valueAtZero = %v", ErrNotFinite, valueAtZero)
	}
	switch {
	case x > 0:
		return 1, nil
	case x < 0:
		return 0, nil
	case x == 0:
		return valueAtZero, nil
	}
	return 0, fmt.Errorf("%w: x = NaN", ErrNotFinite)
}

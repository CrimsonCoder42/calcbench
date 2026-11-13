package hepatic

import (
	"errors"
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

var (
	// ErrNotFinite is returned when an input is NaN or ±Inf; it is checked before any range.
	ErrNotFinite = errors.New("hepatic: value is NaN or infinite")
	// ErrOutOfRange is returned when an input is outside the range the source validated.
	ErrOutOfRange = errors.New("hepatic: input outside the range the source validated")
	// ErrDomain is returned when inputs valid one by one give a result beyond float64 together.
	ErrDomain = errors.New("hepatic: inputs valid one by one but the formula is undefined for them together")
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

// FIB4 returns the FIB-4 index, dimensionless: (age × AST) / (platelets × √ALT). No category is returned;
// the source's abstract states its two cutoffs for its own population (832 HIV/HCV-coinfected persons,
// APRICOT), which are not applied here.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range; ErrDomain when the result is beyond float64.
//
// ID:      MED-HEPAT-001
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Sterling RK, Lissen E, Clumeck N, et al. Development of a simple noninvasive index to predict significant fibrosis in patients with HIV/HCV coinfection. Hepatology. 2006;43(6):1317-1325. doi:10.1002/hep.21178. PMID 16729309. Abstract formula "age ([yr] x AST [U/L]) / ((PLT [10^9/L]) x (ALT [U/L])^(1/2))"
// Tags:    hepatic, liver, fibrosis, fib-4, ast, alt, platelets
// Units:   ageYears years (decimal); astUL U/L; altUL U/L; plateletsPer10e9L 10^9/L; result dimensionless
// Range:   every input > 0 (population: 832 HIV/HCV-coinfected persons (APRICOT); input limits not stated in the source)
// Example: age 50, AST 40, ALT 25, platelets 200 → 50 × 40 / (200 × 5) = 2.0 (hand-derived)
func FIB4(ageYears, astUL, altUL, plateletsPer10e9L float64) (float64, error) {
	if err := checkFinite("ageYears", ageYears); err != nil {
		return 0, err
	}
	if err := checkFinite("astUL", astUL); err != nil {
		return 0, err
	}
	if err := checkFinite("altUL", altUL); err != nil {
		return 0, err
	}
	if err := checkFinite("plateletsPer10e9L", plateletsPer10e9L); err != nil {
		return 0, err
	}
	if err := checkPositive("ageYears", ageYears); err != nil {
		return 0, err
	}
	if err := checkPositive("astUL", astUL); err != nil {
		return 0, err
	}
	if err := checkPositive("altUL", altUL); err != nil {
		return 0, err
	}
	if err := checkPositive("plateletsPer10e9L", plateletsPer10e9L); err != nil {
		return 0, err
	}
	r := (ageYears * astUL) / (plateletsPer10e9L * math.Sqrt(altUL))
	if !reduce.AllFinite(r) || r == 0 {
		// A product overflowed or underflowed on its own; the ratio form keeps the in-range digits.
		r = (ageYears / plateletsPer10e9L) * (astUL / math.Sqrt(altUL))
	}
	if !reduce.AllFinite(r) {
		return 0, fmt.Errorf("%w: result is beyond float64", ErrDomain)
	}
	return r, nil
}

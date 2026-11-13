package chemistry

import (
	"errors"
	"fmt"
	"math"
)

var (
	// ErrNotFinite reports a NaN or infinite input.
	ErrNotFinite = errors.New("chemistry: value is NaN or infinite")
	// ErrNegative reports a negative quantity.
	ErrNegative = errors.New("chemistry: value is negative")
	// ErrZero reports a zero divisor.
	ErrZero = errors.New("chemistry: divisor is zero")
	// ErrOverflow reports a result beyond float64.
	ErrOverflow = errors.New("chemistry: result overflows float64")
)

// checkNonNegative rejects NaN, ±Inf and negative values, naming the parameter.
func checkNonNegative(name string, v float64) error {
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return fmt.Errorf("%w: %s = %v", ErrNotFinite, name, v)
	}
	if v < 0 {
		return fmt.Errorf("%w: %s = %v", ErrNegative, name, v)
	}
	return nil
}

// AmountFromMass returns the amount of substance of a sample, n = m / M, in moles.
// Errors: ErrNotFinite; ErrNegative for a negative mass or molar mass; ErrZero for molar mass 0; ErrOverflow.
//
// ID:      SCI-CHEM-001
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Flowers P, Theopold K, Langley R, Robinson WR. Chemistry 2e. OpenStax, 2019, https://openstax.org/details/books/chemistry-2e, §3.1 "Formula Mass and the Mole Concept" (moles = mass / molar mass); IUPAC, Compendium of Chemical Terminology (Gold Book), "molar mass", M = m/n
// Tags:    mole, amount-of-substance, molar-mass, stoichiometry, chemistry
// Units:   massGrams g; molarMassGramsPerMole g/mol; result mol
// Example: 100 g of NaCl (58.44 g/mol) → 1.711156742…  mol
func AmountFromMass(massGrams, molarMassGramsPerMole float64) (float64, error) {
	if err := checkNonNegative("massGrams", massGrams); err != nil {
		return 0, err
	}
	if err := checkNonNegative("molarMassGramsPerMole", molarMassGramsPerMole); err != nil {
		return 0, err
	}
	if molarMassGramsPerMole == 0 {
		return 0, fmt.Errorf("%w: molarMassGramsPerMole = 0", ErrZero)
	}
	n := massGrams / molarMassGramsPerMole
	if math.IsInf(n, 0) {
		return 0, fmt.Errorf("%w: massGrams / molarMassGramsPerMole", ErrOverflow)
	}
	return n, nil
}

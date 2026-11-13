package chemistry

import (
	"fmt"
	"math"
)

// Molarity returns the amount concentration of a solution, c = n / V, in mol/L (M).
// Errors: ErrNotFinite; ErrNegative for a negative amount or volume; ErrZero for volume 0; ErrOverflow.
//
// ID:      SCI-CHEM-002
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Flowers P, Theopold K, Langley R, Robinson WR. Chemistry 2e. OpenStax, 2019, https://openstax.org/details/books/chemistry-2e, §3.3 "Molarity" (M = mol solute / L solution); IUPAC, Compendium of Chemical Terminology (Gold Book), "amount concentration", c = n/V
// Tags:    molarity, concentration, amount-concentration, solution, chemistry
// Units:   amountMoles mol; volumeLiters L; result mol/L
// Example: 0.133 mol in 0.355 L → 0.37465 mol/L
func Molarity(amountMoles, volumeLiters float64) (float64, error) {
	if err := checkNonNegative("amountMoles", amountMoles); err != nil {
		return 0, err
	}
	if err := checkNonNegative("volumeLiters", volumeLiters); err != nil {
		return 0, err
	}
	if volumeLiters == 0 {
		return 0, fmt.Errorf("%w: volumeLiters = 0", ErrZero)
	}
	c := amountMoles / volumeLiters
	if math.IsInf(c, 0) {
		return 0, fmt.Errorf("%w: amountMoles / volumeLiters", ErrOverflow)
	}
	return c, nil
}

// MolarityFromMass returns the molarity of a solution made from a mass of solute, c = m / (M · V), in mol/L:
// AmountFromMass then Molarity.
// Errors: as AmountFromMass and Molarity.
//
// ID:      SCI-CHEM-003
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Flowers P, Theopold K, Langley R, Robinson WR. Chemistry 2e. OpenStax, 2019, https://openstax.org/details/books/chemistry-2e, §3.1 (moles = mass / molar mass) and §3.3 "Molarity" (M = mol solute / L solution)
// Tags:    molarity, concentration, molar-mass, solution, chemistry
// Units:   massGrams g; molarMassGramsPerMole g/mol; volumeLiters L; result mol/L
// Example: 100 g NaCl (58.44 g/mol) in 2 L → 0.855578371 mol/L
func MolarityFromMass(massGrams, molarMassGramsPerMole, volumeLiters float64) (float64, error) {
	n, err := AmountFromMass(massGrams, molarMassGramsPerMole)
	if err != nil {
		return 0, err
	}
	return Molarity(n, volumeLiters)
}

// MassForMolarity returns the mass of solute that makes a volume of solution at a molarity, m = c · V · M,
// in grams.
// Errors: ErrNotFinite; ErrNegative for any negative input; ErrOverflow.
//
// ID:      SCI-CHEM-004
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Flowers P, Theopold K, Langley R, Robinson WR. Chemistry 2e. OpenStax, 2019, https://openstax.org/details/books/chemistry-2e, §3.3 "Molarity" (mol solute = M × L solution) and §3.1 (mass = moles × molar mass)
// Tags:    molarity, concentration, molar-mass, solution, preparation, chemistry
// Units:   molarityMolPerLiter mol/L; volumeLiters L; molarMassGramsPerMole g/mol; result g
// Example: 0.5 mol/L × 2 L × 58.44 g/mol → 58.44 g
func MassForMolarity(molarityMolPerLiter, volumeLiters, molarMassGramsPerMole float64) (float64, error) {
	for _, p := range []struct {
		name string
		v    float64
	}{{"molarityMolPerLiter", molarityMolPerLiter}, {"volumeLiters", volumeLiters}, {"molarMassGramsPerMole", molarMassGramsPerMole}} {
		if err := checkNonNegative(p.name, p.v); err != nil {
			return 0, err
		}
	}
	m := molarityMolPerLiter * volumeLiters * molarMassGramsPerMole
	if math.IsInf(m, 0) {
		return 0, fmt.Errorf("%w: molarity · volume · molar mass", ErrOverflow)
	}
	return m, nil
}

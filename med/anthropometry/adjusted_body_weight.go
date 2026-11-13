package anthropometry

import "fmt"

// AdjustedBodyWeight returns the adjusted body weight, in kg: idealWeightKg + 0.4 × (actualWeightKg −
// idealWeightKg), the source's "normalized body mass plus 40% of the adipose mass". The 0.4 was derived for
// the distribution volume of gentamicin and tobramycin (the source's context); the function returns a
// weight, nothing else. The result lies between the two inputs, so it never overflows.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range.
//
// ID:      MED-ANTHRO-005
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Schwartz SN, Pazin GJ, Lyon JA, Ho M, Pasculle AW. A controlled investigation of the pharmacokinetics of gentamicin and tobramycin in obese subjects. J Infect Dis. 1978;138(4):499-505. doi:10.1093/infdis/138.4.499. PMID 712111. Abstract
// Tags:    anthropometry, body-weight, adjusted-body-weight, dosing-weight
// Units:   idealWeightKg kg; actualWeightKg kg; result kg
// Range:   idealWeightKg > 0; actualWeightKg > idealWeightKg (the adipose mass is actual − ideal; the source's population: "13 obese and 13 non-obese subjects")
// Example: ideal 70 kg, actual 120 kg → 70 + 0.4 × 50 = 90 kg (HAND-DERIVED)
func AdjustedBodyWeight(idealWeightKg, actualWeightKg float64) (float64, error) {
	if err := checkFinite("idealWeightKg", idealWeightKg); err != nil {
		return 0, err
	}
	if err := checkFinite("actualWeightKg", actualWeightKg); err != nil {
		return 0, err
	}
	if err := checkPositive("idealWeightKg", idealWeightKg); err != nil {
		return 0, err
	}
	if !(actualWeightKg > idealWeightKg) {
		return 0, fmt.Errorf("%w: actualWeightKg = %v, want > idealWeightKg = %v", ErrOutOfRange, actualWeightKg, idealWeightKg)
	}
	return idealWeightKg + 0.4*(actualWeightKg-idealWeightKg), nil
}

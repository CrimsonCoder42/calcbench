package anthropometry

import "fmt"

// EstimatedWeightInfantBestGuess returns the estimated weight of an infant under 12 months, in kg, by the
// Best Guess formula for that band: (ageMonths + 9) / 2. The paper's other bands (1–4 years, 5–14 years)
// are not built here.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range.
//
// ID:      MED-ANTHRO-007
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Tinning K, Acworth J. Make your Best Guess: an updated method for paediatric weight estimation in emergencies. Emerg Med Australas. 2007;19(6):528-534. doi:10.1111/j.1742-6723.2007.01026.x. PMID 18021105. Abstract, the "<12 months" formula
// Tags:    anthropometry, weight-estimation, pediatric, infant, best-guess
// Units:   ageMonths months (decimal); result kg
// Range:   0 ≤ ageMonths < 12 (the "<12 months" band; the study had 70,181 emergency department presentations aged 0–14 years)
// Example: 6 months → (6 + 9) / 2 = 7.5 kg; 0 months → 4.5 kg (HAND-DERIVED)
func EstimatedWeightInfantBestGuess(ageMonths float64) (float64, error) {
	if err := checkFinite("ageMonths", ageMonths); err != nil {
		return 0, err
	}
	if ageMonths < 0 || ageMonths >= 12 {
		return 0, fmt.Errorf("%w: ageMonths = %v, want 0 ≤ ageMonths < 12", ErrOutOfRange, ageMonths)
	}
	return (ageMonths + 9) / 2, nil
}

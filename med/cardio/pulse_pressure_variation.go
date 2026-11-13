package cardio

import (
	"fmt"
	"math"
)

// PulsePressureVariationPercent returns the respiratory variation of the arterial pulse pressure, ΔPP, in
// percent, by Michard 2000: the difference between the maximal and the minimal pulse pressure over one
// mechanical breath divided by the mean of the two, expressed as a percentage:
// ΔPP = 100 × (PPmax − PPmin) / ((PPmax + PPmin) / 2). The abstract reports that a ΔPP of 13% discriminated
// the study's responders from nonresponders to volume expansion; that cut-off is documented here and not
// applied (the abstract does not say on which side the value 13 itself falls).
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange when pulsePressureMinMmHg ≤ 0 or
// pulsePressureMaxMmHg < pulsePressureMinMmHg; ErrDomain when their sum overflows float64.
//
// ID:      MED-CARDIO-046
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Michard F, Boussat S, Chemla D, et al. Relation between respiratory changes in arterial pulse pressure and fluid responsiveness in septic patients with acute circulatory failure. Am J Respir Crit Care Med. 2000;162(1):134-138. doi:10.1164/ajrccm.162.1.9903035. PMID 10903232. Abstract: "calculated as the difference between Pp(max) and Pp(min) divided by the mean of the two values and were expressed as a percentage"
// Tags:    cardio, hemodynamics, pulse-pressure, ppv, fluid-responsiveness, ventilation
// Units:   pulsePressureMaxMmHg, pulsePressureMinMmHg mmHg (systolic minus diastolic, over one breath); result percent
// Range:   pulsePressureMaxMmHg ≥ pulsePressureMinMmHg > 0; the study's population: mechanically ventilated adults with sepsis-related circulatory failure
// Example: PPmax 50, PPmin 43 → 100 × 7 / 46.5 = 15.0537… % (HAND-DERIVED from the abstract's definition)
func PulsePressureVariationPercent(pulsePressureMaxMmHg, pulsePressureMinMmHg float64) (float64, error) {
	if err := checkFinite("pulsePressureMaxMmHg", pulsePressureMaxMmHg); err != nil {
		return 0, err
	}
	if err := checkFinite("pulsePressureMinMmHg", pulsePressureMinMmHg); err != nil {
		return 0, err
	}
	if err := checkPositive("pulsePressureMinMmHg", pulsePressureMinMmHg); err != nil {
		return 0, err
	}
	if pulsePressureMaxMmHg < pulsePressureMinMmHg {
		return 0, fmt.Errorf("%w: pulsePressureMaxMmHg = %v below pulsePressureMinMmHg = %v",
			ErrOutOfRange, pulsePressureMaxMmHg, pulsePressureMinMmHg)
	}
	// (max − min) / ((max + min)/2) × 100 = 200 × ((max − min) / (max + min)). The ratio is in [0, 1), so only
	// the sum can overflow, for values beyond any pressure; that case is reported.
	sum := pulsePressureMaxMmHg + pulsePressureMinMmHg
	if math.IsInf(sum, 1) {
		return 0, fmt.Errorf("%w: pulsePressureMaxMmHg + pulsePressureMinMmHg overflows", ErrDomain)
	}
	return 200 * ((pulsePressureMaxMmHg - pulsePressureMinMmHg) / sum), nil
}

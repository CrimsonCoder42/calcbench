package respiratory

import "math"

// CURB65Score returns the CURB-65 score (0–5) for community-acquired pneumonia: one point each for confusion,
// urea > 7 mmol/L, respiratory rate ≥ 30/min, low blood pressure (systolic < 90 mmHg or diastolic ≤ 60 mmHg)
// and age ≥ 65 years.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range.
//
// ID:      MED-RESP-011
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Lim WS, van der Eerden MM, Laing R, et al. Defining community acquired pneumonia severity on presentation to hospital: an international derivation and validation study. Thorax. 2003;58(5):377-382. doi:10.1136/thorax.58.5.377. PMID 12728155, PMC1746657, "Confusion", "Urea >7 mmol/l", "Respiratory rate ⩾30/min", "low systolic (<90 mm Hg) or diastolic (⩽60 mm Hg) Blood pressure", "age ⩾65 years", one point each
// Tags:    respiratory, pneumonia, cap, curb-65, severity, score
// Units:   ureaMmolL serum urea mmol/L (a BUN in mg/dL is not the same quantity; convert before calling); respiratoryRatePerMin breaths/min; systolicMmHg, diastolicMmHg mmHg; ageYears years; result points 0–5
// Range:   ureaMmolL ≥ 0; respiratoryRatePerMin ≥ 0; systolicMmHg > 0; diastolicMmHg > 0; ageYears ≥ 0 (the text read states no limits; adults with pneumonia were studied)
// Example: confused, urea 8, RR 24, BP 110/70, age 70 → confusion 1 + urea 1 + age 1 = 3 (hand-derived from the abstract's criteria)
func CURB65Score(confusion bool, ureaMmolL, respiratoryRatePerMin, systolicMmHg, diastolicMmHg, ageYears float64) (int, error) {
	if err := checkClosedRange("ureaMmolL", ureaMmolL, 0, math.Inf(1)); err != nil {
		return 0, err
	}
	if err := checkClosedRange("respiratoryRatePerMin", respiratoryRatePerMin, 0, math.Inf(1)); err != nil {
		return 0, err
	}
	if err := checkPositiveFinite("systolicMmHg", systolicMmHg); err != nil {
		return 0, err
	}
	if err := checkPositiveFinite("diastolicMmHg", diastolicMmHg); err != nil {
		return 0, err
	}
	if err := checkClosedRange("ageYears", ageYears, 0, math.Inf(1)); err != nil {
		return 0, err
	}
	n := 0
	if confusion {
		n++
	}
	if ureaMmolL > 7 {
		n++
	}
	if respiratoryRatePerMin >= 30 {
		n++
	}
	if systolicMmHg < 90 || diastolicMmHg <= 60 {
		n++
	}
	if ageYears >= 65 {
		n++
	}
	return n, nil
}

package hepatic

import "math"

// MELDNa returns the MELD-Na score, dimensionless and not rounded: MELD + 1.32 × (137 − Na) − 0.033 × MELD ×
// (137 − Na), with MELD from MED-HEPAT-002 (its bounds and dialysis rule) and sodium bounded to 125–137 mmol/L
// (values outside are clamped, the source's rule). The sodium term is applied at every MELD value: the
// formula as quoted in the source states no MELD threshold, so none is applied here (the OPTN allocation
// policy's condition on MELD is not part of this function; ASSUMPTION A-M9).
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range.
//
// ID:      MED-HEPAT-003
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Kim WR, Mannalithara A, Heimbach JK, et al. MELD 3.0: the Model for End-Stage Liver Disease updated for the modern era. Gastroenterology. 2021;161(6):1887-1895.e4. doi:10.1053/j.gastro.2021.08.050. PMID 34481845. Methods (full text PMC8608337), the MELDNa formula as quoted there
// Tags:    hepatic, liver, meld, meld-na, sodium, bilirubin, inr, creatinine
// Units:   bilirubinMgDL mg/dL; inrRatio dimensionless; creatinineMgDL mg/dL; sodiumMmolL mmol/L; dialysisTwiceInPastWeek true when dialysed at least twice in the past week; result dimensionless score
// Range:   every numeric input > 0 (limits not stated by the source; the bounds above are its rules, applied as clamps)
// Example: bilirubin 1, INR 1, creatinine 1, Na 125 → 6.43 + 1.32 × 12 − 0.033 × 6.43 × 12 = 19.72372 (hand-derived)
func MELDNa(bilirubinMgDL, inrRatio, creatinineMgDL, sodiumMmolL float64, dialysisTwiceInPastWeek bool) (float64, error) {
	if err := checkFinite("sodiumMmolL", sodiumMmolL); err != nil {
		return 0, err
	}
	if err := checkMELDInputs(bilirubinMgDL, inrRatio, creatinineMgDL); err != nil {
		return 0, err
	}
	if err := checkPositive("sodiumMmolL", sodiumMmolL); err != nil {
		return 0, err
	}
	m := meld(bilirubinMgDL, inrRatio, creatinineMgDL, dialysisTwiceInPastWeek)
	gap := 137 - math.Min(math.Max(sodiumMmolL, 125), 137)
	return m + 1.32*gap - 0.033*m*gap, nil
}

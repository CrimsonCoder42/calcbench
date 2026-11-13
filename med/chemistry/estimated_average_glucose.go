package chemistry

import "fmt"

// EstimatedAverageGlucose returns the estimated average glucose (eAG), in mg/dL, from HbA1c in percent by
// the ADAG regression: 28.7 × A1C − 46.7. It is evaluated as (287 × A1C − 467) / 10, the same equation with
// the decimal coefficients scaled to integers, so a whole-percent A1C gives the correctly rounded value
// (6 % → 125.5 exactly, where 28.7 × 6 − 46.7 in float64 gives 125.49999999999999).
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range.
//
// ID:      MED-CHEM-005
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Nathan DM, Kuenen J, Borg R, Zheng H, Schoenfeld D, Heine RJ; A1c-Derived Average Glucose Study Group. Translating the A1C assay into estimated average glucose values. Diabetes Care. 2008;31(8):1473-1478. doi:10.2337/dc08-0545. PMID 18540046. Results: "AG(mg/dl) = 28.7 × A1C − 46.7" (full text PMC2742903)
// Tags:    chemistry, glucose, a1c, hba1c, eag, diabetes
// Units:   a1cPercent percent (NGSP); result mg/dL
// Range:   a1cPercent 4–12 (the study's lowest reported A1C band starts at 4 %; its table tabulates 5–12; ASSUMPTION A-M4)
// Example: A1C 7 → 28.7 × 7 − 46.7 = 154.2 mg/dL; the paper's table prints 154
func EstimatedAverageGlucose(a1cPercent float64) (float64, error) {
	if err := checkA1C(a1cPercent); err != nil {
		return 0, err
	}
	return (287*a1cPercent - 467) / 10, nil
}

// EstimatedAverageGlucoseMmolL returns the estimated average glucose (eAG), in mmol/L, from HbA1c in percent
// by the ADAG regression in SI units: 1.5944 × A1C − 2.594. It is evaluated as (15944 × A1C − 25940) / 10000,
// the same equation with the decimal coefficients scaled to integers.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range.
//
// ID:      MED-CHEM-006
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Nathan DM, Kuenen J, Borg R, Zheng H, Schoenfeld D, Heine RJ; A1c-Derived Average Glucose Study Group. Translating the A1C assay into estimated average glucose values. Diabetes Care. 2008;31(8):1473-1478. doi:10.2337/dc08-0545. PMID 18540046. Results, the mmol/L form 1.5944 × A1C − 2.594 (full text PMC2742903)
// Tags:    chemistry, glucose, a1c, hba1c, eag, diabetes, si
// Units:   a1cPercent percent (NGSP); result mmol/L
// Range:   a1cPercent 4–12 (ASSUMPTION A-M4)
// Example: A1C 7 → 1.5944 × 7 − 2.594 = 8.5668 mmol/L; the paper's table prints 8.6
func EstimatedAverageGlucoseMmolL(a1cPercent float64) (float64, error) {
	if err := checkA1C(a1cPercent); err != nil {
		return 0, err
	}
	return (15944*a1cPercent - 25940) / 10000, nil
}

// checkA1C validates an A1C in percent against the 4–12 % Range of MED-CHEM-005 and -006.
func checkA1C(a1cPercent float64) error {
	if err := checkFinite("a1cPercent", a1cPercent); err != nil {
		return err
	}
	if a1cPercent < 4 || a1cPercent > 12 {
		return fmt.Errorf("%w: a1cPercent = %v, want 4 ≤ A1C ≤ 12", ErrOutOfRange, a1cPercent)
	}
	return nil
}

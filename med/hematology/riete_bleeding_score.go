package hematology

// RIETEBleedingScore returns the RIETE score for major bleeding during anticoagulation for venous
// thromboembolism, 0–8: "2 points to recent bleeding, 1.5 to abnormal creatinine levels or anemia, 1 point to
// the remaining variables" — the six predictors being age > 75 years, recent bleeding, cancer, creatinine
// > 1.2 mg/dL, anemia and pulmonary embolism at baseline. Anemia is the caller's finding (its definition is the
// paper's).
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range.
//
// ID:      MED-HEMAT-009
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Ruíz-Giménez N, Suárez C, González R, et al. Predictive variables for major bleeding events in patients presenting with documented acute venous thromboembolism. Findings from the RIETE Registry. Thromb Haemost. 2008;100(1):26-31. doi:10.1160/TH08-03-0193, PMID 18612534, abstract
// Tags:    hematology, bleeding, anticoagulation, vte, riete, score
// Units:   ageYears years; creatinineMgDL mg/dL; findings present/absent; result points 0–8 in steps of 0.5
// Range:   0 ≤ ageYears ≤ 130; 0 ≤ creatinineMgDL ≤ 30 (physical limits)
// Example: recent bleeding, creatinine 1.3, anemia → 2 + 1.5 + 1.5 = 5; age 76, cancer, PE → 3; all six → 8
func RIETEBleedingScore(ageYears float64, recentBleeding, cancer bool, creatinineMgDL float64, anemia, pulmonaryEmbolism bool) (float64, error) {
	if err := checkFloat("ageYears", ageYears, 0, 130); err != nil {
		return 0, err
	}
	if err := checkFloat("creatinineMgDL", creatinineMgDL, 0, 30); err != nil {
		return 0, err
	}
	// Half-points counted in units of 0.5 so the sum is exact.
	halves := pts(recentBleeding, 4) + pts(creatinineMgDL > 1.2, 3) + pts(anemia, 3) +
		pts(ageYears > 75, 2) + pts(cancer, 2) + pts(pulmonaryEmbolism, 2)
	return float64(halves) / 2, nil
}

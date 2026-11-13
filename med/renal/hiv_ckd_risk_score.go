package renal

// HIVCKDRiskScore returns the Scherzer 2014 risk score for incident chronic kidney disease over 5 years in HIV
// infection, 0–15 (Table 2 "Points Assigned to CKD Risk Factor Categories"): age 19–39 → 0, 40–49 → 2, 50–59 → 4,
// 60–90 → 6; glucose > 140 mg/dL → 2; systolic blood pressure > 140 mmHg → 1; hypertension → 2; triglycerides
// > 200 mg/dL → 1; proteinuria → 2; CD4+ count < 200 cells/µL → 1. The paper reports a 5-year CKD risk from
// less than 1% at 0 points to 16% at 9 or more points in persons not using tenofovir; it prints no category
// labels, so none is returned.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange for an age outside 19–90 or a negative laboratory value.
//
// ID:      MED-RENAL-021
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Scherzer R, Gandhi M, Estrella MM, et al. A chronic kidney disease risk score to determine tenofovir safety in a prospective cohort of HIV-positive male veterans. AIDS. 2014;28(9):1289-1295. doi:10.1097/QAD.0000000000000258. PMID 24922479. PMC4188545, Table 2 "Points Assigned to CKD Risk Factor Categories" and abstract ("less than 1% (zero points) to 16% (≥9 points)")
// Tags:    renal, ckd, hiv, scherzer, risk score
// Units:   ageYears years; glucoseMgDL mg/dL; systolicMmHg mmHg; triglyceridesMgDL mg/dL; cd4CellsPerUL cells/µL; result points
// Range:   19 ≤ ageYears ≤ 90 (the table's age rows); glucoseMgDL, systolicMmHg, triglyceridesMgDL, cd4CellsPerUL ≥ 0 (HIV-positive male veterans, the source's population)
// Example: age 55, glucose 150, SBP 145, hypertension, triglycerides 180, no proteinuria, CD4 350 → 4+2+1+2+0+0+0 = 9 (HAND-DERIVED from Table 2)
func HIVCKDRiskScore(ageYears, glucoseMgDL, systolicMmHg float64, hypertension bool, triglyceridesMgDL float64, proteinuria bool, cd4CellsPerUL float64) (int, error) {
	for _, e := range [...]error{
		checkFinite("ageYears", ageYears), checkFinite("glucoseMgDL", glucoseMgDL),
		checkFinite("systolicMmHg", systolicMmHg), checkFinite("triglyceridesMgDL", triglyceridesMgDL),
		checkFinite("cd4CellsPerUL", cd4CellsPerUL),
		checkBetween("ageYears", ageYears, 19, 90), checkNonNegative("glucoseMgDL", glucoseMgDL),
		checkNonNegative("systolicMmHg", systolicMmHg), checkNonNegative("triglyceridesMgDL", triglyceridesMgDL),
		checkNonNegative("cd4CellsPerUL", cd4CellsPerUL),
	} {
		if e != nil {
			return 0, e
		}
	}
	s := 0
	switch {
	case ageYears >= 60:
		s += 6
	case ageYears >= 50:
		s += 4
	case ageYears >= 40:
		s += 2
	}
	if glucoseMgDL > 140 {
		s += 2
	}
	if systolicMmHg > 140 {
		s++
	}
	if hypertension {
		s += 2
	}
	if triglyceridesMgDL > 200 {
		s++
	}
	if proteinuria {
		s += 2
	}
	if cd4CellsPerUL < 200 {
		s++
	}
	return s, nil
}

package renal

import "fmt"

// McMahonScore returns the McMahon risk prediction score for renal replacement therapy or in-hospital death in
// rhabdomyolysis, 0–19 in steps of 0.5 (McMahon 2013 Table 3): age > 50 to ≤ 70 → 1.5, > 70 to ≤ 80 → 2.5,
// > 80 → 3; female sex → 1; initial creatinine 1.4–2.2 mg/dL → 1.5, > 2.2 → 3; initial calcium < 7.5 mg/dL → 2;
// initial CPK > 40 000 U/L → 2; origin not seizures, syncope, exercise, statins or myositis → 3; initial
// phosphate 4.0–5.4 mg/dL → 1.5, > 5.4 → 3; initial bicarbonate < 19 mEq/L → 2.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange for an age outside 18–120 or a negative laboratory value.
//
// ID:      MED-RENAL-020
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  McMahon GM, Zeng X, Waikar SS. A risk prediction score for kidney failure or mortality in rhabdomyolysis. JAMA Intern Med. 2013;173(19):1821-1828. doi:10.1001/jamainternmed.2013.9774. PMID 24000014. PMC5152583, Table 3 "Risk Score" (Score column) and Results ("a score of less than 5 … 2.3% … a score of more than 10 … 61.2%")
// Tags:    renal, rhabdomyolysis, aki, mcmahon, score
// Units:   ageYears years; creatinineMgDL mg/dL; calciumMgDL mg/dL; cpkUL U/L; phosphateMgDL mg/dL; bicarbonateMEqL mEq/L; result points
// Range:   18 ≤ ageYears ≤ 120 (adult hospital cohorts, the source's population; 120 is a physical cap); laboratory values ≥ 0
// Example: age 75, male, creatinine 2.5, calcium 7.0, CPK 50 000, origin trauma, phosphate 5.0, bicarbonate 18 → 2.5+0+3+2+2+3+1.5+2 = 16 (HAND-DERIVED from Table 3)
func McMahonScore(ageYears float64, female bool, creatinineMgDL, calciumMgDL, cpkUL float64, originNotListed bool, phosphateMgDL, bicarbonateMEqL float64) (float64, error) {
	for _, e := range [...]error{
		checkFinite("ageYears", ageYears), checkFinite("creatinineMgDL", creatinineMgDL),
		checkFinite("calciumMgDL", calciumMgDL), checkFinite("cpkUL", cpkUL),
		checkFinite("phosphateMgDL", phosphateMgDL), checkFinite("bicarbonateMEqL", bicarbonateMEqL),
		checkBetween("ageYears", ageYears, 18, 120), checkNonNegative("creatinineMgDL", creatinineMgDL),
		checkNonNegative("calciumMgDL", calciumMgDL), checkNonNegative("cpkUL", cpkUL),
		checkNonNegative("phosphateMgDL", phosphateMgDL), checkNonNegative("bicarbonateMEqL", bicarbonateMEqL),
	} {
		if e != nil {
			return 0, e
		}
	}
	s := 0.0
	switch {
	case ageYears > 80:
		s += 3
	case ageYears > 70:
		s += 2.5
	case ageYears > 50:
		s += 1.5
	}
	if female {
		s++
	}
	switch {
	case creatinineMgDL > 2.2:
		s += 3
	case creatinineMgDL >= 1.4:
		s += 1.5
	}
	if calciumMgDL < 7.5 {
		s += 2
	}
	if cpkUL > 40000 {
		s += 2
	}
	if originNotListed {
		s += 3
	}
	switch {
	case phosphateMgDL > 5.4:
		s += 3
	case phosphateMgDL >= 4.0:
		s += 1.5
	}
	if bicarbonateMEqL < 19 {
		s += 2
	}
	return s, nil
}

// checkBetween returns ErrOutOfRange, wrapped with name, unless lo ≤ v ≤ hi. NaN passes; check finite first.
func checkBetween(name string, v, lo, hi float64) error {
	if v < lo || v > hi {
		return fmt.Errorf("%w: %s = %v, want %v–%v", ErrOutOfRange, name, v, lo, hi)
	}
	return nil
}

// checkNonNegative returns ErrOutOfRange, wrapped with name, when v < 0. NaN passes; check finite first.
func checkNonNegative(name string, v float64) error {
	if v < 0 {
		return fmt.Errorf("%w: %s = %v, want ≥ 0", ErrOutOfRange, name, v)
	}
	return nil
}

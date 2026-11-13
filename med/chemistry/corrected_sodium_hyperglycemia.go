package chemistry

// CorrectedSodiumHyperglycemia returns the serum sodium corrected for hyperglycaemia, in mEq/L: measured Na⁺
// plus 1.6 mEq/L for each 100 mg/dL of glucose above 100 mg/dL, i.e. Na + 1.6 × (glucose − 100)/100 (the
// factor 1.6 is Katz's). At glucose ≤ 100 mg/dL nothing is added and the measured sodium is returned.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range; ErrDomain when the result is beyond float64.
//
// ID:      MED-CHEM-013
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Kitabchi AE, Umpierrez GE, Miles JM, Fisher JN. Hyperglycemic crises in adult patients with diabetes. Diabetes Care. 2009;32(7):1335-1343. doi:10.2337/dc09-9032. PMID 19564476. PMC2699725, Diagnosis: "For each 100 mg/dl of glucose above 100 mg/dl, add 1.6 mEq to sodium value for corrected serum value"
// Tags:    chemistry, sodium, corrected-sodium, hyperglycemia, glucose, katz, dka, hhs
// Units:   sodiumMEqL mEq/L; glucoseMgDL mg/dL; result mEq/L
// Range:   every input > 0 (the source states no limits)
// Example: Na 130, glucose 600 → 130 + 1.6 × 5 = 138 mEq/L (hand-derived from the quoted sentence)
func CorrectedSodiumHyperglycemia(sodiumMEqL, glucoseMgDL float64) (float64, error) {
	if err := checkFinite("sodiumMEqL", sodiumMEqL); err != nil {
		return 0, err
	}
	if err := checkFinite("glucoseMgDL", glucoseMgDL); err != nil {
		return 0, err
	}
	if err := checkPositive("sodiumMEqL", sodiumMEqL); err != nil {
		return 0, err
	}
	if err := checkPositive("glucoseMgDL", glucoseMgDL); err != nil {
		return 0, err
	}
	excess := glucoseMgDL - 100
	if excess < 0 {
		excess = 0
	}
	r := sodiumMEqL + 1.6*excess/100
	if err := checkResult(r); err != nil {
		return 0, err
	}
	return r, nil
}

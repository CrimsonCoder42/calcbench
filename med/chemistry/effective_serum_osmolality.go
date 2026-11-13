package chemistry

// EffectiveSerumOsmolality returns the effective serum osmolality, in mOsm/kg, by the ADA consensus form:
// 2 × measured Na⁺ + glucose/18 (urea is left out, as an ineffective osmole). For the total calculated
// osmolality with urea see MED-CHEM-003 (Dorwart–Chalmers).
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range; ErrDomain when the result is beyond float64.
//
// ID:      MED-CHEM-012
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Kitabchi AE, Umpierrez GE, Miles JM, Fisher JN. Hyperglycemic crises in adult patients with diabetes. Diabetes Care. 2009;32(7):1335-1343. doi:10.2337/dc09-9032. PMID 19564476. PMC2699725, Table 1 footnote: "Effective serum osmolality: 2[measured Na+ (mEq/l)] + glucose (mg/dl)/18"
// Tags:    chemistry, osmolality, effective-osmolality, tonicity, sodium, glucose, hhs, dka
// Units:   sodiumMEqL mEq/L; glucoseMgDL mg/dL; result mOsm/kg
// Range:   every input > 0 (the source states no limits)
// Example: Na 140, glucose 180 → 280 + 10 = 290 mOsm/kg (hand-derived from the Table 1 footnote)
func EffectiveSerumOsmolality(sodiumMEqL, glucoseMgDL float64) (float64, error) {
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
	r := 2*sodiumMEqL + glucoseMgDL/18
	if err := checkResult(r); err != nil {
		return 0, err
	}
	return r, nil
}

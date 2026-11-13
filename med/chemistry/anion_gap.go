package chemistry

// AnionGap returns the serum anion gap, in mEq/L: Na⁺ − (Cl⁻ + HCO₃⁻), with no albumin adjustment (for the
// albumin-adjusted gap pass this result to MED-CHEM-007 or -008). A result ≤ 0 is returned as computed.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range; ErrDomain when the result is beyond float64.
//
// ID:      MED-CHEM-011
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Kitabchi AE, Umpierrez GE, Miles JM, Fisher JN. Hyperglycemic crises in adult patients with diabetes. Diabetes Care. 2009;32(7):1335-1343. doi:10.2337/dc09-9032. PMID 19564476. PMC2699725, Table 1 footnote: "Anion gap: (Na+) − [(Cl− + HCO3− (mEq/l)]"
// Tags:    chemistry, anion-gap, sodium, chloride, bicarbonate, electrolytes
// Units:   sodiumMEqL, chlorideMEqL, bicarbonateMEqL mEq/L (= mmol/L for these monovalent ions); result mEq/L
// Range:   every input > 0 (the source states no limits)
// Example: Na 140, Cl 104, HCO₃ 24 → 140 − 128 = 12 mEq/L (hand-derived from the Table 1 footnote)
func AnionGap(sodiumMEqL, chlorideMEqL, bicarbonateMEqL float64) (float64, error) {
	if err := checkFinite("sodiumMEqL", sodiumMEqL); err != nil {
		return 0, err
	}
	if err := checkFinite("chlorideMEqL", chlorideMEqL); err != nil {
		return 0, err
	}
	if err := checkFinite("bicarbonateMEqL", bicarbonateMEqL); err != nil {
		return 0, err
	}
	if err := checkPositive("sodiumMEqL", sodiumMEqL); err != nil {
		return 0, err
	}
	if err := checkPositive("chlorideMEqL", chlorideMEqL); err != nil {
		return 0, err
	}
	if err := checkPositive("bicarbonateMEqL", bicarbonateMEqL); err != nil {
		return 0, err
	}
	r := sodiumMEqL - (chlorideMEqL + bicarbonateMEqL)
	if err := checkResult(r); err != nil {
		return 0, err
	}
	return r, nil
}

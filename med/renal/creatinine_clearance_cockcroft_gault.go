package renal

import "fmt"

// CreatinineClearanceCockcroftGault returns the predicted creatinine clearance, in mL/min, by the
// Cockcroft–Gault formula: (140 − age) × weightKg / (72 × Scr mg/dL), × 0.85 if female (the paper's
// "15% less in females"; the derivation population was adult males).
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range; ErrDomain when the result is beyond float64.
//
// ID:      MED-RENAL-005
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Cockcroft DW, Gault MH. Prediction of creatinine clearance from serum creatinine. Nephron. 1976;16(1):31-41. doi:10.1159/000180580. PMID 1244564. Formula as printed on the publisher abstract page: "Ccr = (140 – age) (wt kg)/72 × Scr(mg/100ml)" "(15% less in females)"
// Tags:    renal, creatinine-clearance, cockcroft-gault, creatinine
// Units:   ageYears years (decimal); weightKg kg; creatinineMgDL mg/dL (= mg/100 mL); female the source's sex factor; result mL/min
// Range:   ageYears 18–92 (the 249 persons aged 18–92 the formula was derived from); weightKg > 0; creatinineMgDL > 0
// Example: male, 40 y, 72 kg, Scr 1.0 → (140 − 40) × 72 / (72 × 1.0) = 100 mL/min; female, same → 85 (hand-derived)
func CreatinineClearanceCockcroftGault(ageYears, weightKg, creatinineMgDL float64, female bool) (float64, error) {
	if err := checkCockcroftGault(ageYears, weightKg, "creatinineMgDL", creatinineMgDL); err != nil {
		return 0, err
	}
	r := creatinineClearanceCockcroftGault(ageYears, weightKg, creatinineMgDL, female)
	if err := checkResult(r); err != nil {
		return 0, err
	}
	return r, nil
}

// CreatinineClearanceCockcroftGaultSI returns the predicted creatinine clearance, in mL/min, from serum
// creatinine in µmol/L: creatinine ÷ 88.4 gives mg/dL, then the Cockcroft–Gault formula (MED-RENAL-005).
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range; ErrDomain when the result is beyond float64.
//
// ID:      MED-RENAL-006
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Cockcroft DW, Gault MH. Prediction of creatinine clearance from serum creatinine. Nephron. 1976;16(1):31-41. doi:10.1159/000180580. PMID 1244564. Formula as printed on the publisher abstract page: "Ccr = (140 – age) (wt kg)/72 × Scr(mg/100ml)" "(15% less in females)"
// Source:  Kidney Disease: Improving Global Outcomes (KDIGO) CKD Work Group. KDIGO 2024 Clinical Practice Guideline for the Evaluation and Management of Chronic Kidney Disease. Kidney Int. 2024;105(4S):S117–S314. doi:10.1016/j.kint.2023.10.018. p. S126, "Conversion factors of conventional units to SI units" (creatinine 88.4)
// Tags:    renal, creatinine-clearance, cockcroft-gault, creatinine, si
// Units:   ageYears years (decimal); weightKg kg; creatinineUmolL µmol/L; female the source's sex factor; result mL/min
// Range:   ageYears 18–92; weightKg > 0; creatinineUmolL > 0
// Example: male, 40 y, 72 kg, Scr 88.4 µmol/L (= 1.0 mg/dL) → 100 mL/min (hand-derived)
func CreatinineClearanceCockcroftGaultSI(ageYears, weightKg, creatinineUmolL float64, female bool) (float64, error) {
	if err := checkCockcroftGault(ageYears, weightKg, "creatinineUmolL", creatinineUmolL); err != nil {
		return 0, err
	}
	r := creatinineClearanceCockcroftGault(ageYears, weightKg, creatinineUmolL/creatinineUmolLPerMgDL, female)
	if err := checkResult(r); err != nil {
		return 0, err
	}
	return r, nil
}

// checkCockcroftGault validates the inputs of MED-RENAL-005 and -006: every finiteness check first, then ranges.
func checkCockcroftGault(ageYears, weightKg float64, creatinineName string, creatinine float64) error {
	if err := checkFinite("ageYears", ageYears); err != nil {
		return err
	}
	if err := checkFinite("weightKg", weightKg); err != nil {
		return err
	}
	if err := checkFinite(creatinineName, creatinine); err != nil {
		return err
	}
	if ageYears < 18 || ageYears > 92 {
		return fmt.Errorf("%w: ageYears = %v, want 18 ≤ age ≤ 92", ErrOutOfRange, ageYears)
	}
	if err := checkPositive("weightKg", weightKg); err != nil {
		return err
	}
	return checkPositive(creatinineName, creatinine)
}

// creatinineClearanceCockcroftGault is the 1976 formula on validated inputs.
func creatinineClearanceCockcroftGault(ageYears, weightKg, creatinineMgDL float64, female bool) float64 {
	r := (140 - ageYears) * weightKg / (72 * creatinineMgDL)
	if female {
		r *= 0.85
	}
	return r
}

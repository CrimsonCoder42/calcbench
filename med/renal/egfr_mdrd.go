package renal

import "math"

// EGFRMDRD returns the estimated glomerular filtration rate, in mL/min/1.73 m², by the 4-variable MDRD Study
// equation re-expressed for IDMS-standardized serum creatinine:
// 175 × Scr^−1.154 × age^−0.203 × 1.212 [if black] × 0.742 [if female]. black carries the source's own race
// coefficient and label.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range; ErrDomain when the result is beyond float64.
//
// ID:      MED-RENAL-022
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Levey AS, Coresh J, Greene T, et al. Using standardized serum creatinine values in the modification of diet in renal disease study equation for estimating glomerular filtration rate. Ann Intern Med. 2006;145(4):247-254. doi:10.7326/0003-4819-145-4-200608150-00004. PMID 16908915. Abstract: "GFR = 175 x standardized S(cr)(-1.154) x age(-0.203) x 1.212 (if black) x 0.742 (if female)"
// Source:  Levey AS, Coresh J, Greene T, et al. Expressing the Modification of Diet in Renal Disease Study equation for estimating glomerular filtration rate with standardized serum creatinine values. Clin Chem. 2007;53(4):766-772. doi:10.1373/clinchem.2006.077180. PMID 17332152. Abstract (same equation)
// Tags:    renal, egfr, mdrd, creatinine
// Units:   creatinineMgDL mg/dL (IDMS-standardized); ageYears years (decimal); result mL/min/1.73 m²
// Range:   ageYears ≥ 18 (1628 adults with chronic kidney disease, MDRD Study); creatinineMgDL > 0
// Example: Scr 1.0, age 60, female, not black → 175 × 1 × 60^−0.203 × 0.742 = 175 × 0.435547 × 0.742 ≈ 56.556 (hand-derived)
func EGFRMDRD(creatinineMgDL, ageYears float64, female, black bool) (float64, error) {
	if err := checkMDRD("creatinineMgDL", creatinineMgDL, ageYears); err != nil {
		return 0, err
	}
	return egfrMDRD(creatinineMgDL, ageYears, female, black)
}

// EGFRMDRDSI returns the MDRD Study eGFR (MED-RENAL-022), in mL/min/1.73 m², from IDMS-standardized serum
// creatinine in µmol/L: creatinine ÷ 88.4 gives mg/dL, then the same equation.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range; ErrDomain when the result is beyond float64.
//
// ID:      MED-RENAL-023
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Levey AS, Coresh J, Greene T, et al. Using standardized serum creatinine values in the modification of diet in renal disease study equation for estimating glomerular filtration rate. Ann Intern Med. 2006;145(4):247-254. doi:10.7326/0003-4819-145-4-200608150-00004. PMID 16908915. Abstract equation
// Source:  Kidney Disease: Improving Global Outcomes (KDIGO) CKD Work Group. KDIGO 2024 Clinical Practice Guideline for the Evaluation and Management of Chronic Kidney Disease. Kidney Int. 2024;105(4S):S117–S314. doi:10.1016/j.kint.2023.10.018. p. S126, "Conversion factors of conventional units to SI units" (creatinine 88.4)
// Tags:    renal, egfr, mdrd, creatinine, si
// Units:   creatinineUmolL µmol/L (IDMS-standardized); ageYears years (decimal); result mL/min/1.73 m²
// Range:   ageYears ≥ 18; creatinineUmolL > 0
// Example: Scr 88.4 µmol/L (= 1.0 mg/dL), age 60, female, not black → ≈ 56.556 (as MED-RENAL-022)
func EGFRMDRDSI(creatinineUmolL, ageYears float64, female, black bool) (float64, error) {
	if err := checkMDRD("creatinineUmolL", creatinineUmolL, ageYears); err != nil {
		return 0, err
	}
	return egfrMDRD(creatinineUmolL/creatinineUmolLPerMgDL, ageYears, female, black)
}

// checkMDRD validates the inputs of MED-RENAL-022 and -023: finiteness first, then ranges.
func checkMDRD(creatinineName string, creatinine, ageYears float64) error {
	for _, e := range [...]error{
		checkFinite(creatinineName, creatinine), checkFinite("ageYears", ageYears),
		checkPositive(creatinineName, creatinine), checkAdultAge(ageYears),
	} {
		if e != nil {
			return e
		}
	}
	return nil
}

// egfrMDRD is the shared equation; creatinine in mg/dL.
func egfrMDRD(creatinineMgDL, ageYears float64, female, black bool) (float64, error) {
	r := 175 * math.Pow(creatinineMgDL, -1.154) * math.Pow(ageYears, -0.203)
	if black {
		r *= 1.212
	}
	if female {
		r *= 0.742
	}
	if err := checkResult(r); err != nil {
		return 0, err
	}
	return r, nil
}

package renal

import "fmt"

// EGFRSchwartzBedside returns the estimated glomerular filtration rate, in mL/min/1.73 m², of a child by
// the bedside CKiD Schwartz equation: 0.413 × heightCm / creatinineMgDL. ageYears is taken only to
// enforce the Range.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range; ErrDomain when the result is beyond float64.
//
// ID:      MED-RENAL-007
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Schwartz GJ, Muñoz A, Schneider MF, et al. New equations to estimate GFR in children with CKD. J Am Soc Nephrol. 2009;20(3):629-637. doi:10.1681/ASN.2008030287. PMID 19158356. Abstract: "bedside calculation of 0.413*(height/serum creatinine)"
// Tags:    renal, egfr, schwartz, pediatric, creatinine
// Units:   heightCm cm; creatinineMgDL mg/dL; ageYears years (decimal); result mL/min/1.73 m²
// Range:   ageYears 1–16 (349 children aged 1 to 16, CKiD); heightCm > 0; creatinineMgDL > 0
// Example: height 100 cm, Scr 0.413 → 0.413 × 100 / 0.413 = 100 (hand-derived)
func EGFRSchwartzBedside(heightCm, creatinineMgDL, ageYears float64) (float64, error) {
	if err := checkSchwartz(heightCm, "creatinineMgDL", creatinineMgDL, ageYears); err != nil {
		return 0, err
	}
	r := 0.413 * heightCm / creatinineMgDL
	if err := checkResult(r); err != nil {
		return 0, err
	}
	return r, nil
}

// EGFRSchwartzBedsideSI returns the estimated glomerular filtration rate, in mL/min/1.73 m², of a child from
// serum creatinine in µmol/L: creatinine ÷ 88.4 gives mg/dL, then the bedside Schwartz equation (MED-RENAL-007).
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range; ErrDomain when the result is beyond float64.
//
// ID:      MED-RENAL-008
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Schwartz GJ, Muñoz A, Schneider MF, et al. New equations to estimate GFR in children with CKD. J Am Soc Nephrol. 2009;20(3):629-637. doi:10.1681/ASN.2008030287. PMID 19158356. Abstract: "bedside calculation of 0.413*(height/serum creatinine)"
// Source:  Kidney Disease: Improving Global Outcomes (KDIGO) CKD Work Group. KDIGO 2024 Clinical Practice Guideline for the Evaluation and Management of Chronic Kidney Disease. Kidney Int. 2024;105(4S):S117–S314. doi:10.1016/j.kint.2023.10.018. p. S126, "Conversion factors of conventional units to SI units" (creatinine 88.4)
// Tags:    renal, egfr, schwartz, pediatric, creatinine, si
// Units:   heightCm cm; creatinineUmolL µmol/L; ageYears years (decimal); result mL/min/1.73 m²
// Range:   ageYears 1–16; heightCm > 0; creatinineUmolL > 0
// Example: height 100 cm, Scr 88.4 µmol/L (= 1.0 mg/dL) → 0.413 × 100 / 1.0 = 41.3 (hand-derived)
func EGFRSchwartzBedsideSI(heightCm, creatinineUmolL, ageYears float64) (float64, error) {
	if err := checkSchwartz(heightCm, "creatinineUmolL", creatinineUmolL, ageYears); err != nil {
		return 0, err
	}
	r := 0.413 * heightCm / (creatinineUmolL / creatinineUmolLPerMgDL)
	if err := checkResult(r); err != nil {
		return 0, err
	}
	return r, nil
}

// checkSchwartz validates the inputs of MED-RENAL-007 and -008: every finiteness check first, then ranges.
func checkSchwartz(heightCm float64, creatinineName string, creatinine, ageYears float64) error {
	if err := checkFinite("heightCm", heightCm); err != nil {
		return err
	}
	if err := checkFinite(creatinineName, creatinine); err != nil {
		return err
	}
	if err := checkFinite("ageYears", ageYears); err != nil {
		return err
	}
	if err := checkPositive("heightCm", heightCm); err != nil {
		return err
	}
	if err := checkPositive(creatinineName, creatinine); err != nil {
		return err
	}
	if ageYears < 1 || ageYears > 16 {
		return fmt.Errorf("%w: ageYears = %v, want 1 ≤ age ≤ 16", ErrOutOfRange, ageYears)
	}
	return nil
}

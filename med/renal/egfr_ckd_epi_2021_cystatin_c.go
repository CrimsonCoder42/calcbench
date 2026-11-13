package renal

import "math"

// EGFRCKDEPI2021CystatinC returns the estimated glomerular filtration rate, in mL/min/1.73 m², from serum
// creatinine in mg/dL and serum cystatin C in mg/L by the 2021 CKD-EPI creatinine–cystatin C equation
// (race-free): 135 × min(Scr/κ, 1)^a1 × max(Scr/κ, 1)^−0.544 × min(Scys/0.8, 1)^−0.323 ×
// max(Scys/0.8, 1)^−0.778 × 0.9961^age × 0.963 [if female], with κ = 0.7 (female) or 0.9 (male) and
// a1 = −0.219 (female) or −0.144 (male).
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range (every in-range input gives a finite result).
//
// ID:      MED-RENAL-003
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Inker LA, Eneanya ND, Coresh J, et al. New creatinine- and cystatin C-based equations to estimate GFR without race. N Engl J Med. 2021;385(19):1737-1749. doi:10.1056/NEJMoa2102953. PMID 34554658. Table 2, eGFRcr-cys(AS) column (full text PMC8822996)
// Tags:    renal, egfr, ckd-epi, creatinine, cystatin-c
// Units:   creatinineMgDL mg/dL; cystatinCMgL mg/L; ageYears years (decimal); female the source's sex coefficient; result mL/min/1.73 m²
// Range:   ageYears ≥ 18 ("All the participants were 18 years of age or older"); creatinineMgDL > 0; cystatinCMgL > 0; no upper limits stated by the source
// Example: female, 50 y, Scr 0.7, Scys 0.8 → 135 × 0.9961^50 × 0.963 = 106.9318450861178 (hand-derived from Table 2)
func EGFRCKDEPI2021CystatinC(creatinineMgDL, cystatinCMgL, ageYears float64, female bool) (float64, error) {
	if err := checkFinite("creatinineMgDL", creatinineMgDL); err != nil {
		return 0, err
	}
	if err := checkFinite("cystatinCMgL", cystatinCMgL); err != nil {
		return 0, err
	}
	if err := checkFinite("ageYears", ageYears); err != nil {
		return 0, err
	}
	if err := checkPositive("creatinineMgDL", creatinineMgDL); err != nil {
		return 0, err
	}
	if err := checkPositive("cystatinCMgL", cystatinCMgL); err != nil {
		return 0, err
	}
	if err := checkAdultAge(ageYears); err != nil {
		return 0, err
	}
	return egfrCKDEPI2021CystatinC(creatinineMgDL, cystatinCMgL, ageYears, female), nil
}

// EGFRCKDEPI2021CystatinCSI returns the estimated glomerular filtration rate, in mL/min/1.73 m², from serum
// creatinine in µmol/L and serum cystatin C in mg/L: creatinine ÷ 88.4 gives mg/dL, then the 2021 CKD-EPI
// creatinine–cystatin C equation (MED-RENAL-003).
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range; ErrDomain when the result is beyond float64.
//
// ID:      MED-RENAL-004
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Inker LA, Eneanya ND, Coresh J, et al. New creatinine- and cystatin C-based equations to estimate GFR without race. N Engl J Med. 2021;385(19):1737-1749. doi:10.1056/NEJMoa2102953. PMID 34554658. Table 2, eGFRcr-cys(AS) column (full text PMC8822996)
// Source:  Kidney Disease: Improving Global Outcomes (KDIGO) CKD Work Group. KDIGO 2024 Clinical Practice Guideline for the Evaluation and Management of Chronic Kidney Disease. Kidney Int. 2024;105(4S):S117–S314. doi:10.1016/j.kint.2023.10.018. p. S126, "Conversion factors of conventional units to SI units" (creatinine 88.4)
// Tags:    renal, egfr, ckd-epi, creatinine, cystatin-c, si
// Units:   creatinineUmolL µmol/L; cystatinCMgL mg/L; ageYears years (decimal); female the source's sex coefficient; result mL/min/1.73 m²
// Range:   ageYears ≥ 18; creatinineUmolL > 0; cystatinCMgL > 0; no upper limits stated by the source
// Example: female, 50 y, Scr 61.88 µmol/L (= 0.7 mg/dL), Scys 0.8 → 106.9318450861178 (MED-RENAL-003), to 1e-12 relative
func EGFRCKDEPI2021CystatinCSI(creatinineUmolL, cystatinCMgL, ageYears float64, female bool) (float64, error) {
	if err := checkFinite("creatinineUmolL", creatinineUmolL); err != nil {
		return 0, err
	}
	if err := checkFinite("cystatinCMgL", cystatinCMgL); err != nil {
		return 0, err
	}
	if err := checkFinite("ageYears", ageYears); err != nil {
		return 0, err
	}
	if err := checkPositive("creatinineUmolL", creatinineUmolL); err != nil {
		return 0, err
	}
	if err := checkPositive("cystatinCMgL", cystatinCMgL); err != nil {
		return 0, err
	}
	if err := checkAdultAge(ageYears); err != nil {
		return 0, err
	}
	r := egfrCKDEPI2021CystatinC(creatinineUmolL/creatinineUmolLPerMgDL, cystatinCMgL, ageYears, female)
	if err := checkResult(r); err != nil {
		return 0, err
	}
	return r, nil
}

// egfrCKDEPI2021CystatinC is the Table 2 eGFRcr-cys(AS) equation on validated inputs.
func egfrCKDEPI2021CystatinC(creatinineMgDL, cystatinCMgL, ageYears float64, female bool) float64 {
	kappa, a1, sexFactor := 0.9, -0.144, 1.0
	if female {
		kappa, a1, sexFactor = 0.7, -0.219, 0.963
	}
	cr := creatinineMgDL / kappa
	cys := cystatinCMgL / 0.8
	return 135 * math.Pow(math.Min(cr, 1), a1) * math.Pow(math.Max(cr, 1), -0.544) *
		math.Pow(math.Min(cys, 1), -0.323) * math.Pow(math.Max(cys, 1), -0.778) *
		math.Pow(0.9961, ageYears) * sexFactor
}

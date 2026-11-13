package renal

import (
	"errors"
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

var (
	// ErrNotFinite is returned when an input is NaN or ±Inf; it is checked before any range.
	ErrNotFinite = errors.New("renal: value is NaN or infinite")
	// ErrOutOfRange is returned when an input is outside the range the source validated.
	ErrOutOfRange = errors.New("renal: input outside the range the source validated")
	// ErrDomain is returned when inputs valid one by one give a result beyond float64 together.
	ErrDomain = errors.New("renal: inputs valid one by one but the formula is undefined for them together")
)

// KDIGO 2024 CKD guideline, Kidney Int. 2024;105(4S):S117–S314, doi:10.1016/j.kint.2023.10.018, p. S126
// "Conversion factors of conventional units to SI units" (read 2026-11-12): creatinine 1 mg/dL = 88.4 µmol/L.
const creatinineUmolLPerMgDL = 88.4

// checkFinite returns ErrNotFinite, wrapped with name, when v is NaN or ±Inf.
func checkFinite(name string, v float64) error {
	if !reduce.AllFinite(v) {
		return fmt.Errorf("%w: %s = %v", ErrNotFinite, name, v)
	}
	return nil
}

// checkPositive returns ErrOutOfRange, wrapped with name, unless v > 0.
func checkPositive(name string, v float64) error {
	if !(v > 0) {
		return fmt.Errorf("%w: %s = %v, want > 0", ErrOutOfRange, name, v)
	}
	return nil
}

// checkAdultAge returns ErrOutOfRange unless ageYears ≥ 18 (the CKD-EPI 2021 population).
func checkAdultAge(ageYears float64) error {
	if ageYears < 18 {
		return fmt.Errorf("%w: ageYears = %v, want ≥ 18", ErrOutOfRange, ageYears)
	}
	return nil
}

// checkResult returns ErrDomain when a result computed from finite inputs is not finite.
func checkResult(v float64) error {
	if !reduce.AllFinite(v) {
		return fmt.Errorf("%w: result is beyond float64", ErrDomain)
	}
	return nil
}

// EGFRCKDEPI2021 returns the estimated glomerular filtration rate, in mL/min/1.73 m², from serum creatinine
// in mg/dL by the 2021 CKD-EPI creatinine equation (race-free):
// 142 × min(Scr/κ, 1)^α × max(Scr/κ, 1)^−1.200 × 0.9938^age × 1.012 [if female], with κ = 0.7 (female) or
// 0.9 (male) and α = −0.241 (female) or −0.302 (male).
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range (every in-range input gives a finite result).
//
// ID:      MED-RENAL-001
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Inker LA, Eneanya ND, Coresh J, et al. New creatinine- and cystatin C-based equations to estimate GFR without race. N Engl J Med. 2021;385(19):1737-1749. doi:10.1056/NEJMoa2102953. PMID 34554658. Table 2, eGFRcr(AS) column (full text PMC8822996)
// Tags:    renal, egfr, ckd-epi, creatinine
// Units:   creatinineMgDL mg/dL; ageYears years (decimal, completed age as the caller has it); female the source's sex coefficient; result mL/min/1.73 m²
// Range:   ageYears ≥ 18 ("All the participants were 18 years of age or older"); creatinineMgDL > 0; no upper limits stated by the source
// Example: female, 50 y, Scr 0.7 → 142 × 0.9938^50 × 1.012 = 105.29760114914401 (hand-derived from Table 2)
func EGFRCKDEPI2021(creatinineMgDL, ageYears float64, female bool) (float64, error) {
	if err := checkFinite("creatinineMgDL", creatinineMgDL); err != nil {
		return 0, err
	}
	if err := checkFinite("ageYears", ageYears); err != nil {
		return 0, err
	}
	if err := checkPositive("creatinineMgDL", creatinineMgDL); err != nil {
		return 0, err
	}
	if err := checkAdultAge(ageYears); err != nil {
		return 0, err
	}
	return egfrCKDEPI2021(creatinineMgDL, ageYears, female), nil
}

// EGFRCKDEPI2021SI returns the estimated glomerular filtration rate, in mL/min/1.73 m², from serum creatinine
// in µmol/L: creatinine ÷ 88.4 gives mg/dL, then the 2021 CKD-EPI creatinine equation (MED-RENAL-001).
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range; ErrDomain when the result is beyond float64.
//
// ID:      MED-RENAL-002
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Inker LA, Eneanya ND, Coresh J, et al. New creatinine- and cystatin C-based equations to estimate GFR without race. N Engl J Med. 2021;385(19):1737-1749. doi:10.1056/NEJMoa2102953. PMID 34554658. Table 2, eGFRcr(AS) column (full text PMC8822996)
// Source:  Kidney Disease: Improving Global Outcomes (KDIGO) CKD Work Group. KDIGO 2024 Clinical Practice Guideline for the Evaluation and Management of Chronic Kidney Disease. Kidney Int. 2024;105(4S):S117–S314. doi:10.1016/j.kint.2023.10.018. p. S126, "Conversion factors of conventional units to SI units" (creatinine 88.4)
// Tags:    renal, egfr, ckd-epi, creatinine, si
// Units:   creatinineUmolL µmol/L; ageYears years (decimal); female the source's sex coefficient; result mL/min/1.73 m²
// Range:   ageYears ≥ 18; creatinineUmolL > 0; no upper limits stated by the source
// Example: female, 50 y, Scr 61.88 µmol/L (= 0.7 mg/dL) → 105.29760114914401 (MED-RENAL-001), to 1e-12 relative
func EGFRCKDEPI2021SI(creatinineUmolL, ageYears float64, female bool) (float64, error) {
	if err := checkFinite("creatinineUmolL", creatinineUmolL); err != nil {
		return 0, err
	}
	if err := checkFinite("ageYears", ageYears); err != nil {
		return 0, err
	}
	if err := checkPositive("creatinineUmolL", creatinineUmolL); err != nil {
		return 0, err
	}
	if err := checkAdultAge(ageYears); err != nil {
		return 0, err
	}
	r := egfrCKDEPI2021(creatinineUmolL/creatinineUmolLPerMgDL, ageYears, female)
	if err := checkResult(r); err != nil {
		return 0, err
	}
	return r, nil
}

// egfrCKDEPI2021 is the Table 2 eGFRcr(AS) equation on validated inputs.
func egfrCKDEPI2021(creatinineMgDL, ageYears float64, female bool) float64 {
	kappa, alpha, sexFactor := 0.9, -0.302, 1.0
	if female {
		kappa, alpha, sexFactor = 0.7, -0.241, 1.012
	}
	ratio := creatinineMgDL / kappa
	return 142 * math.Pow(math.Min(ratio, 1), alpha) * math.Pow(math.Max(ratio, 1), -1.200) *
		math.Pow(0.9938, ageYears) * sexFactor
}

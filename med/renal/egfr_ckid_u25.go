package renal

import (
	"fmt"
	"math"
)

// ckidU25CreatinineK returns K of the CKiD U25 height/serum-creatinine equation (Pierce 2021 Table 2,
// age-dependent model): males 39.0 × 1.008^(age−12) for ages 1–12, 39.0 × 1.045^(age−12) for 12–18, 50.8 from 18;
// females 36.1 × 1.008^(age−12), 36.1 × 1.023^(age−12), 41.4.
func ckidU25CreatinineK(ageYears float64, female bool) float64 {
	switch {
	case ageYears >= 18 && female:
		return 41.4
	case ageYears >= 18:
		return 50.8
	case ageYears < 12 && female:
		return 36.1 * math.Pow(1.008, ageYears-12)
	case ageYears < 12:
		return 39.0 * math.Pow(1.008, ageYears-12)
	case female:
		return 36.1 * math.Pow(1.023, ageYears-12)
	}
	return 39.0 * math.Pow(1.045, ageYears-12)
}

// ckidU25CystatinCK returns K of the CKiD U25 cystatin C equation (Pierce 2021 Table 2, age-dependent model):
// males 87.2 × 1.011^(age−15) for ages 1–15, 87.2 × 0.960^(age−15) for 15–18, 77.1 from 18; females
// 79.9 × 1.004^(age−12) for 1–12, 79.9 × 0.974^(age−12) for 12–18, 68.3 from 18.
func ckidU25CystatinCK(ageYears float64, female bool) float64 {
	switch {
	case ageYears >= 18 && female:
		return 68.3
	case ageYears >= 18:
		return 77.1
	case female && ageYears < 12:
		return 79.9 * math.Pow(1.004, ageYears-12)
	case female:
		return 79.9 * math.Pow(0.974, ageYears-12)
	case ageYears < 15:
		return 87.2 * math.Pow(1.011, ageYears-15)
	}
	return 87.2 * math.Pow(0.960, ageYears-15)
}

// checkU25Age returns ErrOutOfRange unless 1 ≤ ageYears ≤ 25 (the CKiD U25 population, "aged 1 to 25 years").
func checkU25Age(ageYears float64) error {
	if ageYears < 1 || ageYears > 25 {
		return fmt.Errorf("%w: ageYears = %v, want 1–25", ErrOutOfRange, ageYears)
	}
	return nil
}

// checkU25Creatinine validates the inputs of MED-RENAL-017 and -018: finiteness first, then ranges.
func checkU25Creatinine(heightM float64, creatinineName string, creatinine, ageYears float64) error {
	for _, e := range [...]error{
		checkFinite("heightM", heightM), checkFinite(creatinineName, creatinine), checkFinite("ageYears", ageYears),
		checkPositive("heightM", heightM), checkPositive(creatinineName, creatinine), checkU25Age(ageYears),
	} {
		if e != nil {
			return e
		}
	}
	return nil
}

// EGFRCKiDU25Creatinine returns the estimated glomerular filtration rate, in mL/min/1.73 m², of a child or young
// adult by the CKiD U25 height/serum-creatinine equation: K × heightM / creatinineMgDL, with the age- and
// sex-dependent K of Pierce 2021 Table 2.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range; ErrDomain when the result is beyond float64.
//
// ID:      MED-RENAL-017
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Pierce CB, Muñoz A, Ng DK, Warady BA, Furth SL, Schwartz GJ. Age- and sex-dependent clinical equations to estimate glomerular filtration rates in children and young adults with chronic kidney disease. Kidney Int. 2021;99(4):948-956. doi:10.1016/j.kint.2020.10.047. PMID 33301749. PMC9083470, "GFR = K × [height (m)/sCr (mg/dL)]" and Table 2 (age-dependent K); abstract (K 35.7 at age 1 to 50.8 from 18, males; 33.1 to 41.4, females)
// Tags:    renal, egfr, ckid, u25, pediatric, creatinine
// Units:   heightM m; creatinineMgDL mg/dL; ageYears years (decimal); result mL/min/1.73 m²
// Range:   1 ≤ ageYears ≤ 25 (CKiD participants aged 1 to 25 years); heightM > 0; creatinineMgDL > 0
// Example: male aged 1, height 0.75 m, Scr 0.3 → 39.0 × 1.008^−11 × 0.75 / 0.3 = 35.73 × 2.5 ≈ 89.3 (K matches the abstract's 35.7)
func EGFRCKiDU25Creatinine(heightM, creatinineMgDL, ageYears float64, female bool) (float64, error) {
	if err := checkU25Creatinine(heightM, "creatinineMgDL", creatinineMgDL, ageYears); err != nil {
		return 0, err
	}
	r := ckidU25CreatinineK(ageYears, female) * heightM / creatinineMgDL
	if err := checkResult(r); err != nil {
		return 0, err
	}
	return r, nil
}

// EGFRCKiDU25CreatinineSI returns the CKiD U25 height/serum-creatinine eGFR, in mL/min/1.73 m², from serum
// creatinine in µmol/L: creatinine ÷ 88.4 gives mg/dL, then the equation of MED-RENAL-017.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range; ErrDomain when the result is beyond float64.
//
// ID:      MED-RENAL-018
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Pierce CB, Muñoz A, Ng DK, Warady BA, Furth SL, Schwartz GJ. Age- and sex-dependent clinical equations to estimate glomerular filtration rates in children and young adults with chronic kidney disease. Kidney Int. 2021;99(4):948-956. doi:10.1016/j.kint.2020.10.047. PMID 33301749. PMC9083470, Table 2
// Source:  Kidney Disease: Improving Global Outcomes (KDIGO) CKD Work Group. KDIGO 2024 Clinical Practice Guideline for the Evaluation and Management of Chronic Kidney Disease. Kidney Int. 2024;105(4S):S117–S314. doi:10.1016/j.kint.2023.10.018. p. S126, "Conversion factors of conventional units to SI units" (creatinine 88.4)
// Tags:    renal, egfr, ckid, u25, pediatric, creatinine, si
// Units:   heightM m; creatinineUmolL µmol/L; ageYears years (decimal); result mL/min/1.73 m²
// Range:   1 ≤ ageYears ≤ 25; heightM > 0; creatinineUmolL > 0
// Example: female aged 20, height 1.60 m, Scr 70.72 µmol/L (= 0.8 mg/dL) → 41.4 × 1.60 / 0.8 = 82.8 (hand-derived)
func EGFRCKiDU25CreatinineSI(heightM, creatinineUmolL, ageYears float64, female bool) (float64, error) {
	if err := checkU25Creatinine(heightM, "creatinineUmolL", creatinineUmolL, ageYears); err != nil {
		return 0, err
	}
	r := ckidU25CreatinineK(ageYears, female) * heightM / (creatinineUmolL / creatinineUmolLPerMgDL)
	if err := checkResult(r); err != nil {
		return 0, err
	}
	return r, nil
}

// EGFRCKiDU25CystatinC returns the estimated glomerular filtration rate, in mL/min/1.73 m², of a child or young
// adult by the CKiD U25 cystatin C equation: K × 1 / cystatinCMgL, with the age- and sex-dependent K of Pierce
// 2021 Table 2.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range; ErrDomain when the result is beyond float64.
//
// ID:      MED-RENAL-019
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Pierce CB, Muñoz A, Ng DK, Warady BA, Furth SL, Schwartz GJ. Age- and sex-dependent clinical equations to estimate glomerular filtration rates in children and young adults with chronic kidney disease. Kidney Int. 2021;99(4):948-956. doi:10.1016/j.kint.2020.10.047. PMID 33301749. PMC9083470, "GFR = K × [1/cysC (mg/L)]" and Table 2 (age-dependent K)
// Tags:    renal, egfr, ckid, u25, pediatric, cystatin-c
// Units:   cystatinCMgL mg/L; ageYears years (decimal); result mL/min/1.73 m²
// Range:   1 ≤ ageYears ≤ 25; cystatinCMgL > 0
// Example: male aged 20, cystatin C 1.0 mg/L → 77.1 (Table 2, age ≥ 18)
func EGFRCKiDU25CystatinC(cystatinCMgL, ageYears float64, female bool) (float64, error) {
	for _, e := range [...]error{
		checkFinite("cystatinCMgL", cystatinCMgL), checkFinite("ageYears", ageYears),
		checkPositive("cystatinCMgL", cystatinCMgL), checkU25Age(ageYears),
	} {
		if e != nil {
			return 0, e
		}
	}
	r := ckidU25CystatinCK(ageYears, female) / cystatinCMgL
	if err := checkResult(r); err != nil {
		return 0, err
	}
	return r, nil
}

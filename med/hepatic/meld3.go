package hepatic

import "math"

// MELD3 returns the MELD 3.0 score, an integer: 1.33 (if female) + 4.56 × ln(bilirubin) + 0.82 × (137 − Na)
// − 0.24 × (137 − Na) × ln(bilirubin) + 9.09 × ln(INR) + 11.14 × ln(creatinine) + 1.85 × (3.5 − albumin)
// − 1.83 × (3.5 − albumin) × ln(creatinine) + 6, rounded to the nearest integer (the source's rule; it states
// no tie rule, so math.Round, half away from zero, is used: ASSUMPTION A-M6). The source's bounds are applied
// as clamps: bilirubin, INR and creatinine below 1.0 set to 1.0; creatinine at most 3.0; sodium 125–137;
// albumin 1.5–3.5. There is no dialysis input: the equation and bounds as quoted in the source state no
// dialysis rule (ASSUMPTION A-M10).
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range.
//
// ID:      MED-HEPAT-004
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Kim WR, Mannalithara A, Heimbach JK, et al. MELD 3.0: the Model for End-Stage Liver Disease updated for the modern era. Gastroenterology. 2021;161(6):1887-1895.e4. doi:10.1053/j.gastro.2021.08.050. PMID 34481845. Methods/Results, the MELD 3.0 equation (full text PMC8608337)
// Tags:    hepatic, liver, meld, meld-3, sodium, albumin, bilirubin, inr, creatinine
// Units:   female the source's sex term; bilirubinMgDL mg/dL; sodiumMmolL mmol/L; inrRatio dimensionless; creatinineMgDL mg/dL; albuminGDL g/dL; result dimensionless integer score
// Range:   every numeric input > 0 (limits not stated by the source; the bounds above are its rules, applied as clamps)
// Example: male, bilirubin 1, Na 137, INR 1, creatinine 1, albumin 3.5 → 6 (hand-derived: every other term is 0)
func MELD3(female bool, bilirubinMgDL, sodiumMmolL, inrRatio, creatinineMgDL, albuminGDL float64) (int, error) {
	if err := checkFinite("bilirubinMgDL", bilirubinMgDL); err != nil {
		return 0, err
	}
	if err := checkFinite("sodiumMmolL", sodiumMmolL); err != nil {
		return 0, err
	}
	if err := checkFinite("inrRatio", inrRatio); err != nil {
		return 0, err
	}
	if err := checkFinite("creatinineMgDL", creatinineMgDL); err != nil {
		return 0, err
	}
	if err := checkFinite("albuminGDL", albuminGDL); err != nil {
		return 0, err
	}
	if err := checkPositive("bilirubinMgDL", bilirubinMgDL); err != nil {
		return 0, err
	}
	if err := checkPositive("sodiumMmolL", sodiumMmolL); err != nil {
		return 0, err
	}
	if err := checkPositive("inrRatio", inrRatio); err != nil {
		return 0, err
	}
	if err := checkPositive("creatinineMgDL", creatinineMgDL); err != nil {
		return 0, err
	}
	if err := checkPositive("albuminGDL", albuminGDL); err != nil {
		return 0, err
	}
	lnBili := math.Log(math.Max(bilirubinMgDL, 1))
	lnINR := math.Log(math.Max(inrRatio, 1))
	lnCr := math.Log(math.Min(math.Max(creatinineMgDL, 1), 3.0))
	naGap := 137 - math.Min(math.Max(sodiumMmolL, 125), 137)
	albGap := 3.5 - math.Min(math.Max(albuminGDL, 1.5), 3.5)
	score := 4.56*lnBili + 0.82*naGap - 0.24*naGap*lnBili + 9.09*lnINR + 11.14*lnCr +
		1.85*albGap - 1.83*albGap*lnCr + 6
	if female {
		score += 1.33
	}
	return int(math.Round(score)), nil
}

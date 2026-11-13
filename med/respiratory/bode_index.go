package respiratory

import "fmt"

// BODEIndex returns the BODE index (0–10) for COPD from body-mass index, airflow obstruction, dyspnea and
// exercise capacity: FEV₁ ≥ 65 % predicted → 0, 50–64 → 1, 36–49 → 2, ≤ 35 → 3; six-minute walk distance
// ≥ 350 m → 0, 250–349 → 1, 150–249 → 2, ≤ 149 → 3; MMRC dyspnea scale 0–1 → 0, 2 → 1, 3 → 2, 4 → 3; body-mass
// index > 21 → 0, ≤ 21 → 1 kg/m². The table prints whole numbers; a value between two bands is compared with
// the lower edge of the upper band (FEV₁ 64.5 % → 1, 35.5 % → 3; 349.5 m → 1).
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range.
//
// ID:      MED-RESP-046
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Celli BR, Cote CG, Marin JM, et al. The body-mass index, airflow obstruction, dyspnea, and exercise capacity index in chronic obstructive pulmonary disease. N Engl J Med. 2004;350(10):1005-1012. doi:10.1056/NEJMoa021322. PMID 14999112 (primary; full text not opened, abstract: "a multidimensional 10-point scale")
// Source:  Funk GC, Kirchheiner K, Burghuber OC, Hartl S. BODE index versus GOLD classification for explaining anxious and depressive symptoms in patients with COPD — a cross-sectional study. Respir Res. 2009;10:1. doi:10.1186/1465-9921-10-1. PMID 19134192, PMC2626579, Table 1
// Source:  Kamath S, Kumar A, Panda SK, Samanta RP. J Family Med Prim Care. 2020, PMC7842444, BODE table (identical cut-offs; its reference to Celli 2004 NEJM 350 prints wrong pages)
// Verified: secondary (Funk 2009 Respir Res PMC2626579 Table 1, ref 24 Celli 2004; Kamath 2020 PMC7842444, ref Celli 2004)
// Tags:    respiratory, copd, bode, prognosis, score
// Units:   fev1PercentPredicted % predicted (post-bronchodilator); walkDistanceM metres in 6 minutes; mmrcGrade 0–4; bmiKgM2 kg/m²; result points 0–10
// Range:   fev1PercentPredicted > 0; walkDistanceM ≥ 0; mmrcGrade 0–4; bmiKgM2 > 0; adults with COPD
// Example: FEV₁ 45 %, 300 m, MMRC 2, BMI 20 → 2 + 1 + 1 + 1 = 5 (hand-derived from the table)
func BODEIndex(fev1PercentPredicted, walkDistanceM float64, mmrcGrade int, bmiKgM2 float64) (int, error) {
	if err := checkPositiveFinite("fev1PercentPredicted", fev1PercentPredicted); err != nil {
		return 0, err
	}
	if err := checkFiniteNonNegative("walkDistanceM", walkDistanceM); err != nil {
		return 0, err
	}
	if mmrcGrade < 0 || mmrcGrade > 4 {
		return 0, fmt.Errorf("%w: mmrcGrade = %d, want 0 ≤ mmrcGrade ≤ 4", ErrOutOfRange, mmrcGrade)
	}
	if err := checkPositiveFinite("bmiKgM2", bmiKgM2); err != nil {
		return 0, err
	}
	s := 0
	switch {
	case fev1PercentPredicted >= 65:
	case fev1PercentPredicted >= 50:
		s++
	case fev1PercentPredicted >= 36:
		s += 2
	default:
		s += 3
	}
	switch {
	case walkDistanceM >= 350:
	case walkDistanceM >= 250:
		s++
	case walkDistanceM >= 150:
		s += 2
	default:
		s += 3
	}
	if mmrcGrade >= 2 {
		s += mmrcGrade - 1
	}
	if bmiKgM2 <= 21 {
		s++
	}
	return s, nil
}

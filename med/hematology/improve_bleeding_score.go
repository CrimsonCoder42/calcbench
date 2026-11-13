package hematology

import "fmt"

// IMPROVERenal is the renal-function band of the IMPROVE bleeding score.
type IMPROVERenal int

// Renal bands of the IMPROVE bleeding score (GFR in mL/min/m² as the tables print it).
const (
	IMPROVERenalGFR60OrMore IMPROVERenal = iota // GFR ≥ 60, 0 points (the reference band)
	IMPROVERenalGFR30To59                       // "Moderate renal failure (GFR 30–59 vs. ≥60)", 1 point
	IMPROVERenalGFRUnder30                      // "Severe renal failure (GFR <30 vs. ≥60)", 2.5 points
)

// IMPROVEBleedingScore returns the IMPROVE bleeding risk score for acutely ill medical inpatients, 0–30.5:
// moderate renal failure (GFR 30–59) 1; male 1; age 40–84 years 1.5; current cancer 2; rheumatic disease 2;
// central venous catheter 2; ICU/CCU 2.5; severe renal failure (GFR < 30) 2.5; hepatic failure (INR > 1.5) 2.5;
// age ≥ 85 years 3.5; platelet count < 50 × 10⁹/L 4; bleeding in the 3 months before admission 4; active
// gastroduodenal ulcer 4.5. Both sources mark ≥ 7 as the high-risk cut-off; no label is returned.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range.
//
// ID:      MED-HEMAT-016
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Decousus H, Tapson VF, Bergmann JF, et al. Factors at admission associated with bleeding risk in medical patients: findings from the IMPROVE investigators. Chest. 2011;139(1):69-79. doi:10.1378/chest.09-3081, PMID 20453069
// Source:  Arakaki D, Iwata M, Terasawa T. External validation and update of the International Medical Prevention Registry on Venous Thromboembolism bleeding risk score for predicting bleeding in acutely ill hospitalized medical patients: a retrospective single-center cohort study in Japan. Thromb J. 2024;22(1):31. doi:10.1186/s12959-024-00603-w, PMID 38549086, PMC10976666, Table 1
// Source:  Chindamo MC, Marques MA. Bleeding risk assessment for venous thromboembolism prophylaxis. J Vasc Bras. 2021;20:e20200109. doi:10.1590/1677-5449.200109, PMID 34093680, PMC8147884, Table 2
// Verified: secondary (Arakaki et al. 2024, PMC10976666; Chindamo & Marques 2021, PMC8147884)
// Tags:    hematology, bleeding, venous thromboembolism, prophylaxis, improve, score
// Units:   ageYears years; renal band as the constants; findings present/absent; result points 0–30.5 (multiples of 0.5, exact in binary)
// Range:   0 ≤ ageYears ≤ 130 (physical limit); renal one of the IMPROVERenal constants (acutely ill medical inpatients, the derivation population)
// Example: man 86, GFR 25, ICU, active ulcer → 1 + 3.5 + 2.5 + 2.5 + 4.5 = 14; woman 39, GFR 70, nothing else → 0 (Table 1 of Arakaki et al.)
func IMPROVEBleedingScore(ageYears float64, male bool, renal IMPROVERenal, currentCancer, rheumaticDisease, centralVenousCatheter, icuOrCCU, hepaticFailureINROver1_5, plateletsUnder50E9PerL, bleedingWithin3Months, activeGastroduodenalUlcer bool) (float64, error) {
	if err := checkFloat("ageYears", ageYears, 0, 130); err != nil {
		return 0, err
	}
	s := 0.0
	switch renal {
	case IMPROVERenalGFR60OrMore:
	case IMPROVERenalGFR30To59:
		s++
	case IMPROVERenalGFRUnder30:
		s += 2.5
	default:
		return 0, fmt.Errorf("%w: renal = %d, want 0–2", ErrOutOfRange, renal)
	}
	switch {
	case ageYears >= 85:
		s += 3.5
	case ageYears >= 40:
		s += 1.5
	}
	s += ptsF(male, 1) + ptsF(currentCancer, 2) + ptsF(rheumaticDisease, 2) + ptsF(centralVenousCatheter, 2) +
		ptsF(icuOrCCU, 2.5) + ptsF(hepaticFailureINROver1_5, 2.5) + ptsF(plateletsUnder50E9PerL, 4) +
		ptsF(bleedingWithin3Months, 4) + ptsF(activeGastroduodenalUlcer, 4.5)
	return s, nil
}

// ptsF returns w when present is true and 0 otherwise (half-point weights).
func ptsF(present bool, w float64) float64 {
	if present {
		return w
	}
	return 0
}

package cardio

import "fmt"

// CHADS2Factors are the yes/no items of CHADS₂ besides age (Gage 2001 abstract).
type CHADS2Factors struct {
	CongestiveHeartFailure bool // 1 point
	Hypertension           bool // 1 point
	Diabetes               bool // 1 point
	PriorStrokeOrTIA       bool // "history of stroke or transient ischemic attack": 2 points
}

// CHADS2Score returns the CHADS₂ score, 0–6: one point each for congestive heart failure, hypertension, age
// 75 years or older and diabetes, two for prior stroke or TIA.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange when ageYears is outside 65–95.
//
// ID:      MED-CARDIO-019
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Gage BF, Waterman AD, Shannon W, Boechler M, Rich MW, Radford MJ. Validation of clinical classification schemes for predicting stroke: results from the National Registry of Atrial Fibrillation. JAMA. 2001;285(22):2864-2870. doi:10.1001/jama.285.22.2864. PMID 11401607. Abstract: "assigning 1 point each for the presence of congestive heart failure, hypertension, age 75 years or older, and diabetes mellitus and by assigning 2 points for history of stroke or transient ischemic attack"
// Tags:    cardio, atrial-fibrillation, stroke, chads2, gage, risk-score
// Units:   ageYears years; result points 0–6
// Range:   ageYears 65–95 ("1733 Medicare beneficiaries aged 65 to 95 years" with nonrheumatic AF)
// Example: age 80, hypertension → 2 (abstract: 4.0 strokes per 100 person-years at a score of 2)
func CHADS2Score(ageYears float64, f CHADS2Factors) (int, error) {
	if err := checkFinite("ageYears", ageYears); err != nil {
		return 0, err
	}
	if err := checkRange("ageYears", ageYears, 65, 95); err != nil {
		return 0, err
	}
	s := 0
	for _, present := range [...]bool{f.CongestiveHeartFailure, f.Hypertension, ageYears >= 75, f.Diabetes} {
		if present {
			s++
		}
	}
	if f.PriorStrokeOrTIA {
		s += 2
	}
	return s, nil
}

// CHADS2StrokeRatePer100PersonYears returns the stroke rate per 100 person-years without antithrombotic
// therapy that Gage 2001 reports for a CHADS₂ score: 1.9, 2.8, 4.0, 5.9, 8.5, 12.5, 18.2 for 0 to 6.
// Errors: ErrOutOfRange when score is outside 0–6.
//
// ID:      MED-CARDIO-020
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Gage BF, et al. JAMA. 2001;285(22):2864-2870. doi:10.1001/jama.285.22.2864. PMID 11401607. Abstract: "1.9 (95% CI, 1.2-3.0) for a score of 0; 2.8 (95% CI, 2.0-3.8) for 1; 4.0 (95% CI, 3.1-5.1) for 2; 5.9 (95% CI, 4.6-7.3) for 3; 8.5 (95% CI, 6.3-11.1) for 4; 12.5 (95% CI, 8.2-17.5) for 5; and 18.2 (95% CI, 10.5-27.4) for 6"
// Tags:    cardio, atrial-fibrillation, stroke, chads2, event-rate
// Units:   score points 0–6; result strokes per 100 person-years
// Range:   0 ≤ score ≤ 6
// Example: 6 → 18.2 (abstract)
func CHADS2StrokeRatePer100PersonYears(score int) (float64, error) {
	if score < 0 || score > 6 {
		return 0, fmt.Errorf("%w: score = %d, want 0–6", ErrOutOfRange, score)
	}
	return chads2StrokeRate[score], nil
}

var chads2StrokeRate = [7]float64{1.9, 2.8, 4.0, 5.9, 8.5, 12.5, 18.2}

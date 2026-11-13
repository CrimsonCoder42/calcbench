package cardio

import "fmt"

// OESILFindings are the three yes/no OESIL items besides age (Colivicchi 2003 abstract).
type OESILFindings struct {
	CardiovascularDiseaseHistory bool // "cardiovascular disease in clinical history"
	SyncopeWithoutProdromes      bool // "syncope without prodromes"
	AbnormalECG                  bool // "abnormal electrocardiogram"
}

// OESILScore returns the OESIL syncope score, 0–4: one point for age > 65 years and one for each finding.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange when ageYears is outside 18–120.
//
// ID:      MED-CARDIO-010
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Colivicchi F, Ammirati F, Melina D, Guido V, Imperoli G, Santini M. Development and prospective validation of a risk stratification system for patients with syncope in the emergency department: the OESIL risk score. Eur Heart J. 2003;24(9):811-819. doi:10.1016/S0195-668X(02)00827-8. PMID 12727148. Abstract: "(1) age >65 years; (2) cardiovascular disease in clinical history; (3) syncope without prodromes; and (4) abnormal electrocardiogram"
// Tags:    cardio, syncope, emergency, oesil, risk-score
// Units:   ageYears years; result points 0–4
// Range:   ageYears 18–120 (derivation 270 persons, mean age 59.5, Italy; the abstract states no age limits)
// Example: age 70, abnormal ECG → 2 (abstract: 12-month mortality 19.6% for 2 points)
func OESILScore(ageYears float64, f OESILFindings) (int, error) {
	if err := checkFinite("ageYears", ageYears); err != nil {
		return 0, err
	}
	if err := checkRange("ageYears", ageYears, 18, 120); err != nil {
		return 0, err
	}
	s := 0
	for _, present := range [...]bool{ageYears > 65, f.CardiovascularDiseaseHistory, f.SyncopeWithoutProdromes, f.AbnormalECG} {
		if present {
			s++
		}
	}
	return s, nil
}

// OESILMortalityPercent returns the 12-month total mortality, in percent, that Colivicchi 2003 reports for an
// OESIL score: 0 → 0; 1 → 0.8; 2 → 19.6; 3 → 34.7; 4 → 57.1.
// Errors: ErrOutOfRange when score is outside 0–4.
//
// ID:      MED-CARDIO-011
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Colivicchi F, et al. Eur Heart J. 2003;24(9):811-819. doi:10.1016/S0195-668X(02)00827-8. PMID 12727148. Abstract: "0% for a score of 0, 0.8% for 1 point; 19.6% for 2 points; 34.7% for 3 points; 57.1% for 4 points"
// Tags:    cardio, syncope, oesil, mortality
// Units:   score points 0–4; result percent
// Range:   0 ≤ score ≤ 4
// Example: 3 → 34.7 (abstract)
func OESILMortalityPercent(score int) (float64, error) {
	if score < 0 || score > 4 {
		return 0, fmt.Errorf("%w: score = %d, want 0–4", ErrOutOfRange, score)
	}
	return oesilMortality[score], nil
}

var oesilMortality = [5]float64{0, 0.8, 19.6, 34.7, 57.1}

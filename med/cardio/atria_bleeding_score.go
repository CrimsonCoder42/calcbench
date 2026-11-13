package cardio

// ATRIABleedingFactors are the ATRIA bleeding items besides age (Fang 2011 abstract and Table 3; the caller
// decides each item by the paper's definitions).
type ATRIABleedingFactors struct {
	Anemia             bool // 3 points ("hemoglobin < 13g/dL in men and < 12 g/dL in women")
	SevereRenalDisease bool // 3 points ("e.g., glomerular filtration rate <30 ml/min or dialysis-dependent")
	PriorBleeding      bool // 1 point
	Hypertension       bool // 1 point
}

// ATRIABleedingScore returns the ATRIA hemorrhage score for persons with atrial fibrillation on warfarin,
// 0–10: anemia 3, severe renal disease 3, age ≥ 75 years 2, prior bleeding 1, hypertension 1.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange when ageYears is outside 18–120.
//
// ID:      MED-CARDIO-021
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Fang MC, Go AS, Chang Y, et al. A new risk scheme to predict warfarin-associated hemorrhage: the ATRIA (Anticoagulation and Risk Factors in Atrial Fibrillation) Study. J Am Coll Cardiol. 2011;58(4):395-401. doi:10.1016/j.jacc.2011.03.031. PMID 21757117. Abstract: "anemia (3 points), severe renal disease (e.g., glomerular filtration rate <30 ml/min or dialysis-dependent, 3 points), age ≥75 years (2 points), prior bleeding (1 point), and hypertension (1 point)"
// Tags:    cardio, atrial-fibrillation, bleeding, warfarin, atria, risk-score
// Units:   ageYears years; result points 0–10
// Range:   ageYears 18–120 (9,186 adults with AF on warfarin, Kaiser Permanente Northern California; no age limits stated)
// Example: age 80, anemia → 5 (abstract: "high risk (5 to 10 points)")
func ATRIABleedingScore(ageYears float64, f ATRIABleedingFactors) (int, error) {
	if err := checkFinite("ageYears", ageYears); err != nil {
		return 0, err
	}
	if err := checkRange("ageYears", ageYears, 18, 120); err != nil {
		return 0, err
	}
	s := 0
	add := func(present bool, pts int) {
		if present {
			s += pts
		}
	}
	add(f.Anemia, 3)
	add(f.SevereRenalDisease, 3)
	add(ageYears >= 75, 2)
	add(f.PriorBleeding, 1)
	add(f.Hypertension, 1)
	return s, nil
}

// ATRIABleedingCategory is a risk group of Fang 2011; its value is the abstract's wording.
type ATRIABleedingCategory string

// Fang 2011 groups: "low risk (0 to 3 points)", "intermediate risk (4 points)", "high risk (5 to 10 points)".
const (
	ATRIABleedingLow          ATRIABleedingCategory = "low risk"          // 0–3, major hemorrhage 0.8%/yr
	ATRIABleedingIntermediate ATRIABleedingCategory = "intermediate risk" // 4, 2.6%/yr
	ATRIABleedingHigh         ATRIABleedingCategory = "high risk"         // 5–10, 5.8%/yr
)

// ATRIABleedingGroup returns the Fang 2011 group for an ATRIA bleeding score: 0–3 "low risk", 4
// "intermediate risk", 5–10 "high risk".
// Errors: ErrOutOfRange when score is outside 0–10.
//
// ID:      MED-CARDIO-022
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Fang MC, et al. J Am Coll Cardiol. 2011;58(4):395-401. doi:10.1016/j.jacc.2011.03.031. PMID 21757117. Abstract: "major hemorrhage rates were 0.8% for low risk (0 to 3 points), 2.6% for intermediate risk (4 points), and 5.8% for high risk (5 to 10 points)"
// Tags:    cardio, atrial-fibrillation, bleeding, atria, category
// Units:   score points 0–10; result the abstract's group words
// Range:   0 ≤ score ≤ 10
// Example: 3 → "low risk"; 4 → "intermediate risk"; 5 → "high risk"
func ATRIABleedingGroup(score int) (ATRIABleedingCategory, error) {
	switch {
	case score < 0 || score > 10:
		return "", errScore(score, 0, 10)
	case score <= 3:
		return ATRIABleedingLow, nil
	case score == 4:
		return ATRIABleedingIntermediate, nil
	default:
		return ATRIABleedingHigh, nil
	}
}

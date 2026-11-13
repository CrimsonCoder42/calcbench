package cardio

// ORBITBleedingFactors are the ORBIT items besides age (O'Brien 2015, Table 2; the caller decides each item by
// the paper's definitions).
type ORBITBleedingFactors struct {
	ReducedHemoglobinOrAnemia  bool // "Reduced haemoglobin/Hct/anaemia": 2
	BleedingHistory            bool // "Bleeding history": 2
	InsufficientKidneyFunction bool // "Insufficient kidney function" (eGFR < 60 mL/min/1.73 m²): 1
	AntiplateletTreatment      bool // "Treatment with antiplatelets": 1
}

// ORBITBleedingScore returns the ORBIT bleeding score for atrial fibrillation on oral anticoagulation, 0–7:
// older age (75+ years) 1; reduced haemoglobin/Hct/anaemia 2; bleeding history 2; insufficient kidney
// function 1; treatment with antiplatelets 1.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange when ageYears is outside 18–120.
//
// ID:      MED-CARDIO-025
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  O'Brien EC, Simon DN, Thomas LE, et al. The ORBIT bleeding score: a simple bedside score to assess bleeding risk in atrial fibrillation. Eur Heart J. 2015;36(46):3258-3264. doi:10.1093/eurheartj/ehv476. PMID 26424865. PMC4670965, Table 2
// Tags:    cardio, atrial-fibrillation, bleeding, anticoagulation, orbit, risk-score
// Units:   ageYears years; result points 0–7
// Range:   ageYears 18–120 (7,411 ORBIT-AF participants on OAC, US, median age 75; no age limits stated)
// Example: age 78, antiplatelet treatment → 2 (Table 3: "low" 0–2)
func ORBITBleedingScore(ageYears float64, f ORBITBleedingFactors) (int, error) {
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
	add(ageYears >= 75, 1)
	add(f.ReducedHemoglobinOrAnemia, 2)
	add(f.BleedingHistory, 2)
	add(f.InsufficientKidneyFunction, 1)
	add(f.AntiplateletTreatment, 1)
	return s, nil
}

// ORBITBleedingCategory is a risk group of O'Brien 2015; its value is the paper's wording.
type ORBITBleedingCategory string

// O'Brien 2015, Table 3: low (0–2) 2.4, medium (3) 4.7, high (≥4) 8.1 bleeds per 100 person-years.
const (
	ORBITBleedingLow    ORBITBleedingCategory = "low"
	ORBITBleedingMedium ORBITBleedingCategory = "medium"
	ORBITBleedingHigh   ORBITBleedingCategory = "high"
)

// ORBITBleedingGroup returns the O'Brien 2015 group for an ORBIT score: 0–2 "low", 3 "medium", ≥ 4 "high".
// Errors: ErrOutOfRange when score is outside 0–7.
//
// ID:      MED-CARDIO-026
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  O'Brien EC, et al. Eur Heart J. 2015;36(46):3258-3264. doi:10.1093/eurheartj/ehv476. PMID 26424865. PMC4670965, Table 3: low (0–2) 2.4, medium (3) 4.7, high (≥4) 8.1 "bleeds per 100 patient-years"
// Tags:    cardio, atrial-fibrillation, bleeding, orbit, category
// Units:   score points 0–7; result the paper's group word
// Range:   0 ≤ score ≤ 7
// Example: 2 → "low"; 3 → "medium"; 4 → "high"
func ORBITBleedingGroup(score int) (ORBITBleedingCategory, error) {
	switch {
	case score < 0 || score > 7:
		return "", errScore(score, 0, 7)
	case score <= 2:
		return ORBITBleedingLow, nil
	case score == 3:
		return ORBITBleedingMedium, nil
	default:
		return ORBITBleedingHigh, nil
	}
}

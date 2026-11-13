package cardio

import "math"

// IEMortalityFactors are the yes/no items of the ICE 6-month mortality score besides age (Park 2016, Table 5,
// points in the comments).
type IEMortalityFactors struct {
	DialysisHistory          bool // 3
	Nosocomial               bool // 2
	ProstheticValve          bool // "Prosthetic IE": 1
	SymptomsOver1Month       bool // "Symptoms >1 month before admission": −1
	StaphylococcusAureus     bool // 1
	ViridansStreptococci     bool // −2
	AorticVegetation         bool // 1
	MitralVegetation         bool // 1
	NYHA3or4HeartFailure     bool // 3
	Stroke                   bool // 2
	ParavalvularComplication bool // 2
	PersistentBacteremia     bool // 2
	SurgicalTreatment        bool // −2
}

// IEMortalityScore returns the simplified 6-month mortality risk score for infective endocarditis: a constant
// 4, age points (≤ 45 y 0; 46–60 2; 61–70 3; > 70 4) and the item points of Table 5 (cells checked against the
// PMC XML of Table 5 on 2026-11-12); the table allows totals from −1 to 26.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange when ageYears is below 18.
//
// ID:      MED-CARDIO-035
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Park LP, Chu VH, Peterson G, et al. Validated risk score for predicting 6-month mortality in infective endocarditis. J Am Heart Assoc. 2016;5(4):e003016. doi:10.1161/JAHA.115.003016. PMID 27091179. PMC4859286, Table 5 "Simplified Risk Score Calculation for 6-Month Mortality in IE"
// Tags:    cardio, infective-endocarditis, mortality, ice, risk-score
// Units:   ageYears years (completed: 45.5 scores as ≤ 45 only if the caller rounds down); result points
// Range:   ageYears ≥ 18 (4,049 adults with definite IE, ICE-PCS, 64 centres in 28 countries; no upper age stated by the source, so none is enforced: A-M1)
// Example: age 65, S. aureus, mitral vegetation, stroke → 4 + 3 + 1 + 1 + 2 = 11 (HAND-DERIVED from Table 5)
func IEMortalityScore(ageYears float64, f IEMortalityFactors) (int, error) {
	if err := checkFinite("ageYears", ageYears); err != nil {
		return 0, err
	}
	if err := checkRange("ageYears", ageYears, 18, math.MaxFloat64); err != nil {
		return 0, err
	}
	s := 4
	switch {
	case ageYears > 70:
		s += 4
	case ageYears > 60:
		s += 3
	case ageYears > 45:
		s += 2
	}
	add := func(present bool, pts int) {
		if present {
			s += pts
		}
	}
	add(f.DialysisHistory, 3)
	add(f.Nosocomial, 2)
	add(f.ProstheticValve, 1)
	add(f.SymptomsOver1Month, -1)
	add(f.StaphylococcusAureus, 1)
	add(f.ViridansStreptococci, -2)
	add(f.AorticVegetation, 1)
	add(f.MitralVegetation, 1)
	add(f.NYHA3or4HeartFailure, 3)
	add(f.Stroke, 2)
	add(f.ParavalvularComplication, 2)
	add(f.PersistentBacteremia, 2)
	add(f.SurgicalTreatment, -2)
	return s, nil
}

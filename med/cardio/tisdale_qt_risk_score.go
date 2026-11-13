package cardio

import "fmt"

// TisdaleQTFactors are the yes/no items of the Tisdale QT-prolongation risk score besides age, sex and drug
// count (Tisdale 2013, risk score table).
type TisdaleQTFactors struct {
	LoopDiuretic      bool // loop diuretic: 1
	PotassiumLE35     bool // serum K+ ≤ 3.5 mEq/L: 2
	AdmissionQTcGE450 bool // admission QTc ≥ 450 ms: 2
	AcuteMI           bool // acute MI: 2
	Sepsis            bool // sepsis: 3
	HeartFailure      bool // heart failure: 3
}

// TisdaleQTRiskScore returns the Tisdale risk score for QTc prolongation in cardiac critical care, 0–21:
// age ≥ 68 y 1; female 1; loop diuretic 1; K+ ≤ 3.5 mEq/L 2; admission QTc ≥ 450 ms 2; acute MI 2; one
// QTc-prolonging drug 3; ≥ 2 QTc-prolonging drugs 3 more (6 in all); sepsis 3; heart failure 3.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange when ageYears is outside 18–120 or qtcProlongingDrugs < 0.
//
// ID:      MED-CARDIO-016
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Tisdale JE, Jaynes HA, Kingery JR, et al. Development and validation of a risk score to predict QT interval prolongation in hospitalized patients. Circ Cardiovasc Qual Outcomes. 2013;6(4):479-487. doi:10.1161/CIRCOUTCOMES.113.000152. PMID 23716032. PMC3788679, risk score table (maximum score 21: the one-drug and ≥2-drug items both score 3)
// Tags:    cardio, ecg, qt-prolongation, drug-induced, critical-care, tisdale, risk-score
// Units:   ageYears years; qtcProlongingDrugs count; result points 0–21
// Range:   ageYears 18–120 (900 cardiac critical-care admissions, Indianapolis; the paper's age limits were not read); qtcProlongingDrugs ≥ 0
// Example: woman aged 70, two QTc-prolonging drugs, heart failure → 1 + 1 + 6 + 3 = 11 (HAND-DERIVED)
func TisdaleQTRiskScore(ageYears float64, female bool, qtcProlongingDrugs int, f TisdaleQTFactors) (int, error) {
	if err := checkFinite("ageYears", ageYears); err != nil {
		return 0, err
	}
	if err := checkRange("ageYears", ageYears, 18, 120); err != nil {
		return 0, err
	}
	if qtcProlongingDrugs < 0 {
		return 0, fmt.Errorf("%w: qtcProlongingDrugs = %d, want ≥ 0", ErrOutOfRange, qtcProlongingDrugs)
	}
	s := 0
	add := func(present bool, pts int) {
		if present {
			s += pts
		}
	}
	add(ageYears >= 68, 1)
	add(female, 1)
	add(f.LoopDiuretic, 1)
	add(f.PotassiumLE35, 2)
	add(f.AdmissionQTcGE450, 2)
	add(f.AcuteMI, 2)
	add(qtcProlongingDrugs >= 1, 3)
	add(qtcProlongingDrugs >= 2, 3)
	add(f.Sepsis, 3)
	add(f.HeartFailure, 3)
	return s, nil
}

// TisdaleQTRiskCategory is a Tisdale 2013 risk group; its value is the paper's group name.
type TisdaleQTRiskCategory string

// Tisdale 2013 groups: low < 7, moderate 7–10, high ≥ 11.
const (
	TisdaleQTLow      TisdaleQTRiskCategory = "Low"      // score < 7
	TisdaleQTModerate TisdaleQTRiskCategory = "Moderate" // score 7–10
	TisdaleQTHigh     TisdaleQTRiskCategory = "High"     // score ≥ 11
)

// TisdaleQTRiskGroup returns the Tisdale 2013 risk group for a score: < 7 "Low", 7–10 "Moderate", ≥ 11 "High".
// Errors: ErrOutOfRange when score is outside 0–21.
//
// ID:      MED-CARDIO-017
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Tisdale JE, et al. Circ Cardiovasc Qual Outcomes. 2013;6(4):479-487. doi:10.1161/CIRCOUTCOMES.113.000152. PMID 23716032. PMC3788679: low < 7 (15% with QTc prolongation, derivation), moderate 7–10 (37%), high ≥ 11 (73%)
// Tags:    cardio, ecg, qt-prolongation, tisdale, category
// Units:   score points 0–21; result the paper's group name
// Range:   0 ≤ score ≤ 21
// Example: 6 → "Low"; 7 → "Moderate"; 11 → "High"
func TisdaleQTRiskGroup(score int) (TisdaleQTRiskCategory, error) {
	switch {
	case score < 0 || score > 21:
		return "", fmt.Errorf("%w: score = %d, want 0–21", ErrOutOfRange, score)
	case score < 7:
		return TisdaleQTLow, nil
	case score <= 10:
		return TisdaleQTModerate, nil
	default:
		return TisdaleQTHigh, nil
	}
}

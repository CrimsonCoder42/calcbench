package cardio

// MarburgHeartFindings are the four yes/no findings of the Marburg Heart Score besides age/sex (Bösner 2010,
// score table, verbatim labels in the comments).
type MarburgHeartFindings struct {
	KnownVascularDisease    bool // "Known clinical vascular disease"
	PainWorseDuringExercise bool // "Pain worse during exercise"
	PainNotReproducible     bool // "Pain not reproducible by palpation"
	PersonAssumesCardiac    bool // the person assumes the pain is cardiac (Bösner 2010, fifth item)
}

// MarburgHeartScore returns the Marburg Heart Score, 0–5: one point for "Age/sex (female ≥ 65, male ≥ 55)" and
// one for each finding present. The paper's cut-off: "Negative result ≤ 2 points".
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange when ageYears is outside 35–93.
//
// ID:      MED-CARDIO-007
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Bösner S, Haasenritter J, Becker A, et al. Ruling out coronary artery disease in primary care: development and validation of a simple prediction rule. CMAJ. 2010;182(12):1295-1300. doi:10.1503/cmaj.100212. PMID 20603345. PMC2934794, score table: one point each for "Age/sex (female ≥ 65, male ≥ 55)", "Known clinical vascular disease", "Pain worse during exercise", "Pain not reproducible by palpation", "Patient assumes pain is of cardiac origin"; "Negative result ≤ 2 points"
// Tags:    cardio, chest-pain, coronary-artery-disease, primary-care, marburg, risk-score
// Units:   ageYears years; female the source's sex; result points 0–5
// Range:   ageYears 35–93 (the derivation cohort's ages, 1,199 analysed, Hesse, Germany)
// Example: man aged 60, pain worse on exercise → 2 (HAND-DERIVED)
func MarburgHeartScore(ageYears float64, female bool, f MarburgHeartFindings) (int, error) {
	if err := checkFinite("ageYears", ageYears); err != nil {
		return 0, err
	}
	if err := checkRange("ageYears", ageYears, 35, 93); err != nil {
		return 0, err
	}
	ageSex := (female && ageYears >= 65) || (!female && ageYears >= 55)
	s := 0
	for _, present := range [...]bool{ageSex, f.KnownVascularDisease, f.PainWorseDuringExercise,
		f.PainNotReproducible, f.PersonAssumesCardiac} {
		if present {
			s++
		}
	}
	return s, nil
}

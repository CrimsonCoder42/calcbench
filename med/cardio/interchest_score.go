package cardio

// INTERCHESTFindings are the five yes/no findings of the INTERCHEST rule besides age/sex (Aerts 2017 abstract).
type INTERCHESTFindings struct {
	PhysicianSuspectsSerious bool // "attending physician suspected a serious diagnosis" (+1)
	HistoryOfCAD             bool // "history of CAD" (+1)
	PainOnExertion           bool // "pain brought on by exertion" (+1)
	PainFeelsLikePressure    bool // "pain feels like 'pressure'" (+1)
	PainReproducible         bool // "pain reproducible by palpation" (−1)
}

// INTERCHESTScore returns the INTERCHEST chest-pain score, −1 to 5. The paper's rule: "CAD was considered
// absent if the prediction score is <2".
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange when ageYears is outside 18–120.
//
// ID:      MED-CARDIO-008
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Aerts M, Minalu G, Bösner S, et al. Pooled individual patient data from five countries were used to derive a clinical prediction rule for coronary artery disease in primary care. J Clin Epidemiol. 2017;81:120-128. doi:10.1016/j.jclinepi.2016.09.011. PMID 27773828. Abstract: "age ≥55 (males) or ≥65 (females) (+1); attending physician suspected a serious diagnosis (+1); history of CAD (+1); pain brought on by exertion (+1); pain feels like 'pressure' (+1); pain reproducible by palpation (-1)"
// Tags:    cardio, chest-pain, coronary-artery-disease, primary-care, interchest, risk-score
// Units:   ageYears years; female the source's sex; result points −1 to 5
// Range:   ageYears 18–120 (3,099 primary-care persons; the abstract states no age limits, only the thresholds are applied)
// Example: woman aged 70, pain on exertion, reproducible by palpation → 1 + 1 − 1 = 1 (HAND-DERIVED)
func INTERCHESTScore(ageYears float64, female bool, f INTERCHESTFindings) (int, error) {
	if err := checkFinite("ageYears", ageYears); err != nil {
		return 0, err
	}
	if err := checkRange("ageYears", ageYears, 18, 120); err != nil {
		return 0, err
	}
	s := 0
	if (female && ageYears >= 65) || (!female && ageYears >= 55) {
		s++
	}
	for _, present := range [...]bool{f.PhysicianSuspectsSerious, f.HistoryOfCAD, f.PainOnExertion, f.PainFeelsLikePressure} {
		if present {
			s++
		}
	}
	if f.PainReproducible {
		s--
	}
	return s, nil
}

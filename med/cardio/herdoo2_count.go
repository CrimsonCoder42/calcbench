package cardio

// HERDOO2Count returns the number of HERDOO2 characteristics present in a woman after 5–7 months of
// anticoagulation for a first unprovoked VTE, 0–4: hyperpigmentation, edema or redness of either leg; D-dimer
// ≥ 250 µg/L while taking warfarin (VIDAS assay); body mass index ≥ 30 kg/m²; age ≥ 65 years. Rodger 2008's
// low-risk rule is "0 or 1" of these, in women only ("no combination of clinical predictors" for men).
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange when ageYears is outside 18–120, or bmiKgM2 or
// dDimerUgL is negative.
//
// ID:      MED-CARDIO-043
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Rodger MA, Kahn SR, Wells PS, et al. Identifying unprovoked thromboembolism patients at low risk for recurrence who can discontinue anticoagulant therapy. CMAJ. 2008;179(5):417-426. doi:10.1503/cmaj.080493. PMID 18725614. PMC2518177: "Women with 0 or 1 of the following characteristics: hyperpigmentation, edema or redness of either leg; D-dimer ≥ 250 μg/L while taking warfarin; body mass index ≥ 30 kg/m²; or age ≥ 65 years"
// Tags:    cardio, vte, recurrence, anticoagulation, herdoo2, women
// Units:   ageYears years; bmiKgM2 kg/m²; dDimerUgL µg/L (VIDAS, measured on warfarin); result count 0–4
// Range:   ageYears 18–120 (600 adults with a first unprovoked proximal DVT or PE, 12 centres in 4 countries); bmiKgM2, dDimerUgL ≥ 0
// Example: age 70, BMI 31, D-dimer 200, no leg signs → 2 (HAND-DERIVED)
func HERDOO2Count(ageYears, bmiKgM2, dDimerUgL float64, legSigns bool) (int, error) {
	for _, p := range [...]struct {
		n      string
		v      float64
		lo, hi float64
	}{{"ageYears", ageYears, 18, 120}, {"bmiKgM2", bmiKgM2, 0, maxFinite}, {"dDimerUgL", dDimerUgL, 0, maxFinite}} {
		if err := checkFinite(p.n, p.v); err != nil {
			return 0, err
		}
		if err := checkRange(p.n, p.v, p.lo, p.hi); err != nil {
			return 0, err
		}
	}
	n := 0
	for _, present := range [...]bool{legSigns, dDimerUgL >= 250, bmiKgM2 >= 30, ageYears >= 65} {
		if present {
			n++
		}
	}
	return n, nil
}

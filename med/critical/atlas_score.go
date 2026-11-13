package critical

// ATLASScore returns the five-component ATLAS score for Clostridium difficile infection, 0–10: age < 60 years
// 0, 60–79 1, ≥ 80 2; treatment with systemic antibiotics during CDI therapy (≥ 1 day) no 0, yes 2; leukocyte
// count < 16,000 0, 16,000–25,000 1, > 25,000 2; serum albumin > 35 g/L 0, 26–35 1, ≤ 25 2; serum creatinine
// ≤ 120 µmol/L 0, 121–179 1, ≥ 180 2. Between two printed integers (albumin 25–26 g/L, creatinine 120–121 and
// 179–180 µmol/L) a value is scored by the band whose printed edge it has passed: albumin > 25 scores 1,
// creatinine > 120 scores 1 and ≥ 180 scores 2.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange for an age outside 0–130 or a negative count or
// concentration.
//
// ID:      MED-CRIT-004
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Miller MA, Louie T, Mullane K, et al. Derivation and validation of a simple clinical bedside score (ATLAS) for Clostridium difficile infection which predicts response to therapy. BMC Infect Dis. 2013;13:148. doi:10.1186/1471-2334-13-148. PMID 23530807, PMC3618004, Table "Clinical and laboratory variables, along with their respective values and points" and the abstract's "five-component ATLAS score (age, treatment with systemic antibiotics, leukocyte count, albumin and serum creatinine as a measure of renal function) ... scores of 0 to 10 inclusive"
// Tags:    critical, infectious, clostridium difficile, atlas, score
// Units:   ageYears years; leukocytesPerUL cells/µL (the table's "16,000"); albuminGL g/L; creatinineUmolL µmol/L; result points 0–10
// Range:   0 ≤ ageYears ≤ 130; leukocytesPerUL, albuminGL, creatinineUmolL ≥ 0 (physical limits; adults with CDI, the study's population)
// Example: age 82, antibiotics, WBC 26,000, albumin 25 g/L, creatinine 180 µmol/L → 10; age 59, none, 15,999, 36, 120 → 0 (HAND-DERIVED from the table)
func ATLASScore(ageYears float64, systemicAntibiotics bool, leukocytesPerUL, albuminGL, creatinineUmolL float64) (int, error) {
	for _, e := range [...]error{
		finite("ageYears", ageYears), finite("leukocytesPerUL", leukocytesPerUL),
		finite("albuminGL", albuminGL), finite("creatinineUmolL", creatinineUmolL),
		within("ageYears", ageYears, 0, 130), atLeastZero("leukocytesPerUL", leukocytesPerUL),
		atLeastZero("albuminGL", albuminGL), atLeastZero("creatinineUmolL", creatinineUmolL),
	} {
		if e != nil {
			return 0, e
		}
	}
	s := 0
	switch {
	case ageYears >= 80:
		s += 2
	case ageYears >= 60:
		s++
	}
	if systemicAntibiotics {
		s += 2
	}
	switch {
	case leukocytesPerUL > 25000:
		s += 2
	case leukocytesPerUL >= 16000:
		s++
	}
	switch {
	case albuminGL <= 25:
		s += 2
	case albuminGL <= 35:
		s++
	}
	switch {
	case creatinineUmolL >= 180:
		s += 2
	case creatinineUmolL > 120:
		s++
	}
	return s, nil
}

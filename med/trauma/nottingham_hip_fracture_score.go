package trauma

import "fmt"

// NottinghamHipFractureScore returns the Nottingham Hip Fracture Score (NHFS), 0–10, with haemoglobin in g/dL:
// age 66–85 years 3 points, ≥ 86 years 4 points (below 66: 0); male 1; admission haemoglobin ≤ 10 g/dL 1;
// abbreviated mental test score ≤ 6 (of 10) 1; living in an institution 1; two or more comorbidities 1;
// malignancy 1. The 30-day mortality equation of the 2012 recalibration is not here (see the ledger).
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange for ageYears outside 0–130, haemoglobinGDL ≤ 0,
// amtsScore outside 0–10 or comorbidityCount < 0.
//
// ID:      MED-TRAUMA-021
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Maxwell MJ, Moran CG, Moppett IK. Development and validation of a preoperative scoring system to predict 30 day mortality in patients undergoing hip fracture surgery. Br J Anaesth. 2008;101(4):511-517. doi:10.1093/bja/aen236. PMID 18723517, abstract (the seven variables and their thresholds: "age (66-85 yr, > or =86 yr), sex (male), number of co-morbidities (> or =2), mini-mental test score (< or =6 out of 10), admission haemoglobin concentration (< or =10 g dl(-1)), living in an institution, and presence of malignant disease"); point weights not in the abstract
// Source:  Lisk et al. Calcif Tissue Int. 2020;107(4):319-326. PMC7497295, Methods: "age (66–85 years = 3, ≥ 86 years = 4), sex (male = 1), admission haemoglobin (≤ 100 g/l = 1), AMTS (≤ 6 = 1), residence (living in an institution = 1), co-morbidities (≥ 2 = 1), active malignancy within 20 years (yes = 1)"; cites Maxwell 2008 (ref 9)
// Source:  Karres et al. Arch Orthop Trauma Surg. 2023;143(7):4125-4132. PMC10293368, Table 2 (age < 66 = 0, 66–85 = 3, ≥ 86 = 4; male, haemoglobin ≤ 10 g dl⁻¹, cognitive impairment, institution, comorbidities ≥ 2, malignancy 1 each; maximum 10); cites Maxwell 2008
// Verified: secondary (Lisk 2020 PMC7497295; Karres 2023 PMC10293368)
// Tags:    trauma, hip fracture, nhfs, nottingham, orthopaedic, score
// Units:   ageYears years; haemoglobinGDL g/dL; amtsScore points 0–10; comorbidityCount count; result points 0–10
// Range:   0 ≤ ageYears ≤ 130 (the abstract states no age limits; under 66 scores 0); haemoglobinGDL > 0; 0 ≤ amtsScore ≤ 10; comorbidityCount ≥ 0
// Example: age 90, female, Hb 9.5 g/dL, AMTS 8, from home, 2 comorbidities, no malignancy → 4 + 1 + 1 = 6 (HAND-DERIVED from the points)
func NottinghamHipFractureScore(ageYears float64, male bool, haemoglobinGDL float64, amtsScore int, livesInInstitution bool, comorbidityCount int, malignancy bool) (int, error) {
	if err := finite("haemoglobinGDL", haemoglobinGDL); err != nil {
		return 0, err
	}
	if !(haemoglobinGDL > 0) {
		return 0, fmt.Errorf("%w: haemoglobinGDL = %v, want > 0", ErrOutOfRange, haemoglobinGDL)
	}
	return nhfs(ageYears, male, haemoglobinGDL <= 10, amtsScore, livesInInstitution, comorbidityCount, malignancy)
}

// NottinghamHipFractureScoreSI returns the Nottingham Hip Fracture Score (NHFS), 0–10, with haemoglobin in g/L:
// the points of NottinghamHipFractureScore with the haemoglobin point at ≤ 100 g/L, the threshold Lisk 2020
// prints in g/L (no conversion is applied).
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange for ageYears outside 0–130, haemoglobinGL ≤ 0,
// amtsScore outside 0–10 or comorbidityCount < 0.
//
// ID:      MED-TRAUMA-022
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Maxwell MJ, Moran CG, Moppett IK. Br J Anaesth. 2008;101(4):511-517. doi:10.1093/bja/aen236. PMID 18723517, abstract (variables and thresholds)
// Source:  Lisk et al. Calcif Tissue Int. 2020;107(4):319-326. PMC7497295, Methods: "admission haemoglobin (≤ 100 g/l = 1)" and the other points; cites Maxwell 2008
// Source:  Karres et al. Arch Orthop Trauma Surg. 2023;143(7):4125-4132. PMC10293368, Table 2; cites Maxwell 2008
// Verified: secondary (Lisk 2020 PMC7497295; Karres 2023 PMC10293368)
// Tags:    trauma, hip fracture, nhfs, nottingham, orthopaedic, score, si units
// Units:   ageYears years; haemoglobinGL g/L; amtsScore points 0–10; comorbidityCount count; result points 0–10
// Range:   0 ≤ ageYears ≤ 130; haemoglobinGL > 0; 0 ≤ amtsScore ≤ 10; comorbidityCount ≥ 0
// Example: age 70, male, Hb 100 g/L, AMTS 4, institution, 1 comorbidity, malignancy → 3 + 1 + 1 + 1 + 1 + 1 = 8 (HAND-DERIVED)
func NottinghamHipFractureScoreSI(ageYears float64, male bool, haemoglobinGL float64, amtsScore int, livesInInstitution bool, comorbidityCount int, malignancy bool) (int, error) {
	if err := finite("haemoglobinGL", haemoglobinGL); err != nil {
		return 0, err
	}
	if !(haemoglobinGL > 0) {
		return 0, fmt.Errorf("%w: haemoglobinGL = %v, want > 0", ErrOutOfRange, haemoglobinGL)
	}
	return nhfs(ageYears, male, haemoglobinGL <= 100, amtsScore, livesInInstitution, comorbidityCount, malignancy)
}

// nhfs checks the remaining inputs and sums the NHFS points.
func nhfs(ageYears float64, male, lowHaemoglobin bool, amtsScore int, livesInInstitution bool, comorbidityCount int, malignancy bool) (int, error) {
	if err := checkAge(ageYears); err != nil {
		return 0, err
	}
	if amtsScore < 0 || amtsScore > 10 {
		return 0, fmt.Errorf("%w: amtsScore = %d, want 0–10", ErrOutOfRange, amtsScore)
	}
	if comorbidityCount < 0 {
		return 0, fmt.Errorf("%w: comorbidityCount = %d, want ≥ 0", ErrOutOfRange, comorbidityCount)
	}
	s := 0
	switch {
	case ageYears >= 86:
		s = 4
	case ageYears >= 66:
		s = 3
	}
	s += count(male, lowHaemoglobin, amtsScore <= 6, livesInInstitution, comorbidityCount >= 2, malignancy)
	return s, nil
}

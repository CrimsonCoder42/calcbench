package trauma

// OttawaKneeFindings are the Ottawa knee rule variables besides age (Stiell 1995 abstract).
type OttawaKneeFindings struct {
	FibularHeadTenderness       bool // "(2) tenderness at the head of the fibula"
	IsolatedPatellarTenderness  bool // "(3) isolated tenderness of the patella"
	UnableToFlex90              bool // "(4) inability to flex to 90 degrees"
	UnableToBearWeightFourSteps bool // "(5) inability to bear weight both immediately and in the ED (four steps)"
}

// OttawaKneeRule returns the number of Ottawa knee rule variables present (0–5): age 55 years or older,
// tenderness at the head of the fibula, isolated tenderness of the patella, inability to flex to 90°, and
// inability to bear weight both immediately and in the ED (four steps). The rule is positive when the count is
// at least 1.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange when ageYears is outside 0–130.
//
// ID:      MED-TRAUMA-008
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Stiell IG, Greenberg GH, Wells GA, et al. Derivation of a decision rule for the use of radiography in acute knee injuries. Ann Emerg Med. 1995;26(4):405-413. doi:10.1016/s0196-0644(95)70106-0. PMID 7574120. Abstract: "The derived decision rule included the following variables: (1) age 55 years or older, (2) tenderness at the head of the fibula, (3) isolated tenderness of the patella, (4) inability to flex to 90 degrees, and (5) inability to bear weight both immediately and in the ED (four steps)."
// Tags:    trauma, knee, ottawa knee rule, decision rule, radiography
// Units:   ageYears years; result count 0–5
// Range:   0 ≤ ageYears ≤ 130 (physical limits; acute knee injuries in adults, the study's population, which the caller establishes)
// Example: age 55, nothing else → 1; age 54.9, none → 0 (HAND-DERIVED from "age 55 years or older")
func OttawaKneeRule(ageYears float64, f OttawaKneeFindings) (int, error) {
	if err := checkAge(ageYears); err != nil {
		return 0, err
	}
	return count(ageYears >= 55, f.FibularHeadTenderness, f.IsolatedPatellarTenderness, f.UnableToFlex90,
		f.UnableToBearWeightFourSteps), nil
}

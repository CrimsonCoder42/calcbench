package respiratory

import "fmt"

// checkIntRange returns ErrOutOfRange unless lo ≤ v ≤ hi.
func checkIntRange(name string, v, lo, hi int) error {
	if v < lo || v > hi {
		return fmt.Errorf("%w: %s = %d, want %d ≤ %s ≤ %d", ErrOutOfRange, name, v, lo, name, hi)
	}
	return nil
}

// CHOSENScore returns the CHOSEN (COVID-19 Hospitalization Of Short-stay Expected to be Negative) risk score
// (0–55) of Levine 2020, Table 2, from age, oxygen saturation and serum albumin: age 18–45 → 5, 46–59 → 2,
// 60–73 → 1, > 73 → 0; SpO₂ < 94 → 0, 94–96 → 9, 97–98 → 14, > 98 → 21; albumin < 2.8 → 0, 2.8–3.3 → 5,
// 3.4–3.7 → 15, > 3.7 → 29 g/dL. The table prints albumin to one decimal; a value between 3.3 and 3.4 is
// compared as ≤ 3.3 or not, so 3.35 scores 15. The paper's cut point is 30 ("At a cut point of 30, the score had
// a sensitivity of 83% and specificity of 82%"); no category label is attached, so none is returned.
// Errors: ErrNotFinite for a NaN or ±Inf albumin; ErrOutOfRange outside Range.
//
// ID:      MED-RESP-027
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Levine DM, Lipsitz SR, Co Z, Song W, Dykes PC, Samal L. Derivation of a clinical risk score to predict 14-day occurrence of hypoxia, ICU admission, and death among patients with coronavirus disease 2019. J Gen Intern Med. 2021;36(3):730-737. doi:10.1007/s11606-020-06353-5. PMID 33274414, PMC7713904, Table 2 (points) and Results ("a risk score between 0 and 55")
// Tags:    respiratory, covid-19, chosen, discharge, score
// Units:   ageYears completed years; spo2Percent % (integer, as the table bands); albuminGDL g/dL; result points 0–55
// Range:   ageYears 18–120 ("All adult patients"); spo2Percent 0–100; albuminGDL > 0
// Example: age 40, SpO₂ 99, albumin 4.0 → 5 + 21 + 29 = 55; age 80, SpO₂ 92, albumin 2.5 → 0 (hand-derived from Table 2)
func CHOSENScore(ageYears, spo2Percent int, albuminGDL float64) (int, error) {
	if err := checkIntRange("ageYears", ageYears, 18, 120); err != nil {
		return 0, err
	}
	if err := checkSpO2Percent(spo2Percent); err != nil {
		return 0, err
	}
	if err := checkPositiveFinite("albuminGDL", albuminGDL); err != nil {
		return 0, err
	}
	s := 0
	switch {
	case ageYears <= 45:
		s += 5
	case ageYears <= 59:
		s += 2
	case ageYears <= 73:
		s++
	}
	switch {
	case spo2Percent > 98:
		s += 21
	case spo2Percent >= 97:
		s += 14
	case spo2Percent >= 94:
		s += 9
	}
	switch {
	case albuminGDL > 3.7:
		s += 29
	case albuminGDL > 3.3:
		s += 15
	case albuminGDL >= 2.8:
		s += 5
	}
	return s, nil
}

package cardio

// SanFranciscoSyncopeRuleCount returns how many of the San Francisco Syncope Rule's five items are present,
// 0–5, and whether any is present: an abnormal ECG, a complaint of shortness of breath, hematocrit less than
// 30%, systolic blood pressure less than 90 mm Hg, or a history of congestive heart failure.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range.
//
// ID:      MED-CARDIO-063
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Quinn JV, Stiell IG, McDermott DA, Sellers KL, Kohn MA, Wells GA. Derivation of the San Francisco Rule to predict patients with short-term serious outcomes. Ann Emerg Med. 2004;43(2):224-232. doi:10.1016/S0196-0644(03)00823-0. PMID 14747812. Abstract: the rule considers "an abnormal ECG, a complaint of shortness of breath, hematocrit less than 30%, systolic blood pressure less than 90 mm Hg, or a history of congestive heart failure"
// Tags:    cardio, syncope, emergency, san-francisco, decision-rule
// Units:   hematocritPercent %; systolicMmHg mmHg; result count 0–5 and any-present flag
// Range:   hematocritPercent 5–75; systolicMmHg 30–300 (ED visits for syncope or near syncope)
// Example: hematocrit 29%, SBP 120, no other item → (1, true)
func SanFranciscoSyncopeRuleCount(abnormalECG, shortnessOfBreath, historyOfCHF bool, hematocritPercent, systolicMmHg float64) (count int, anyPresent bool, err error) {
	if err := checkBounded([]boundedInput{
		{"hematocritPercent", hematocritPercent, 5, 75},
		{"systolicMmHg", systolicMmHg, 30, 300},
	}); err != nil {
		return 0, false, err
	}
	for _, p := range [...]bool{abnormalECG, shortnessOfBreath, historyOfCHF, hematocritPercent < 30, systolicMmHg < 90} {
		if p {
			count++
		}
	}
	return count, count > 0, nil
}

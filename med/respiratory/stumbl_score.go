package respiratory

// STUMBLScore returns the blunt chest wall trauma risk score of Battle 2014 (STUMBL), Table 4: age "1 per
// additional 10-year increase starting at 10 years of age" (⌊age/10⌋), rib fractures "3 per additional rib
// fracture", chronic lung disease 5, pre-injury anticoagulant use 4, and oxygen saturation "2 per 5% decrease
// in oxygen saturations starting at 94%": 95 % or more → 0, 90–94 → 2, 85–89 → 4, 80–84 → 6, and so on. This
// reading of "starting at 94%" (94 % is the first scored value) reproduces the paper's worked example (87 % →
// 4); the other reading (94 % scores 0) also gives 4 there and differs only at 94 and at 89/84/… (ledger
// stumbl_battle_score).
// Errors: ErrNotFinite for a NaN or ±Inf age; ErrOutOfRange outside Range.
//
// ID:      MED-RESP-040
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Battle C, Hutchings H, Lovett S, et al. Predicting outcomes after blunt chest wall trauma: development and external validation of a new prognostic model. Crit Care. 2014;18(3):R98. doi:10.1186/cc13873. PMID 24887537, PMC4095687, Table 4 (risk score column) and the worked example
// Tags:    respiratory, trauma, chest-wall, rib-fracture, stumbl, score
// Units:   ageYears years; ribFractures count; spo2Percent % (integer); the other items bool; result points
// Range:   ageYears 18–120 ("Patients were excluded if they were under 18 years of age"); ribFractures 0–24; spo2Percent 0–100
// Example: "a 67-year-old (score 6) patient with three rib fractures (score 9), a history of chronic lung disease (score 5) and oxygen saturations of 87% (score 4) would have a final risk score of 24 points" (Battle 2014)
func STUMBLScore(ageYears float64, ribFractures int, chronicLungDisease, preInjuryAnticoagulant bool, spo2Percent int) (int, error) {
	if err := checkClosedRange("ageYears", ageYears, 18, 120); err != nil {
		return 0, err
	}
	if err := checkIntRange("ribFractures", ribFractures, 0, 24); err != nil {
		return 0, err
	}
	if err := checkSpO2Percent(spo2Percent); err != nil {
		return 0, err
	}
	s := int(ageYears/10) + 3*ribFractures
	if chronicLungDisease {
		s += 5
	}
	if preInjuryAnticoagulant {
		s += 4
	}
	if spo2Percent < 95 {
		s += 2 * ((95 - spo2Percent + 4) / 5)
	}
	return s, nil
}

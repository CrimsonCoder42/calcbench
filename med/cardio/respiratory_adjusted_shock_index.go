package cardio

// RespiratoryAdjustedShockIndex returns RASI, heart rate / systolic blood pressure × (respiratory rate / 10).
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange when an input is ≤ 0.
//
// ID:      MED-CARDIO-033
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Caputo N, Reilly J, Kanter M, West J. A retrospective analysis of the respiratory adjusted shock index to determine the presence of occult shock in trauma patients. J Trauma Acute Care Surg. 2018;84(4):674-678. doi:10.1097/TA.0000000000001761. PMID 29251715. Abstract: "heart rate/systolic blood pressure*(RR/10)"
// Tags:    cardio, hemodynamics, trauma, shock-index, rasi
// Units:   heartRateBpm beats/min; systolicMmHg mmHg; respiratoryRatePerMin breaths/min; result unitless
// Range:   inputs > 0 (3,093 trauma cases, 2012–2016; no limits stated)
// Example: HR 100, SBP 100, RR 20 → 1 × 2 = 2 (HAND-DERIVED)
func RespiratoryAdjustedShockIndex(heartRateBpm, systolicMmHg, respiratoryRatePerMin float64) (float64, error) {
	if err := checkPositivePair("heartRateBpm", heartRateBpm, "systolicMmHg", systolicMmHg); err != nil {
		return 0, err
	}
	if err := checkFinite("respiratoryRatePerMin", respiratoryRatePerMin); err != nil {
		return 0, err
	}
	if err := checkPositive("respiratoryRatePerMin", respiratoryRatePerMin); err != nil {
		return 0, err
	}
	return heartRateBpm / systolicMmHg * (respiratoryRatePerMin / 10), nil
}

package cardio

// ModifiedShockIndex returns the modified shock index, heart rate divided by mean arterial pressure. The paper
// reports MSI > 1.3 or < 0.7 as the thresholds associated with mortality; it does not state how the mean
// pressure was obtained, so the caller passes the mean arterial pressure it has.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange when an input is ≤ 0.
//
// ID:      MED-CARDIO-056
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Liu YC, Liu JH, Fang ZA, et al. Modified shock index and mortality rate of emergency patients. World J Emerg Med. 2012;3(2):114-117. doi:10.5847/wjem.j.issn.1920-8642.2012.02.006. PMID 25215048. PMC4129788, Methods: "MSI, which is a ratio of heart rate to mean blood pressure (MAP)"; Results: thresholds 0.7 and 1.3
// Tags:    cardio, shock-index, hemodynamics, emergency
// Units:   heartRateBpm beats/min; meanArterialMmHg mmHg; result beats/min per mmHg
// Range:   inputs > 0 (adults presenting to the ED; no limits stated)
// Example: HR 110, MAP 70 → 1.5714 (HAND-DERIVED)
func ModifiedShockIndex(heartRateBpm, meanArterialMmHg float64) (float64, error) {
	if err := checkPositivePair("heartRateBpm", heartRateBpm, "meanArterialMmHg", meanArterialMmHg); err != nil {
		return 0, err
	}
	return heartRateBpm / meanArterialMmHg, nil
}

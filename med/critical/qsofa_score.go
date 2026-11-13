package critical

// QSOFAScore returns the quick SOFA (qSOFA) score, 0–3: one point each for systolic hypotension (≤ 100 mmHg),
// tachypnea (≥ 22/min) and altered mentation. Sepsis-3 names a score of 2 or more as its qSOFA criterion
// (Singer 2016).
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange for a negative rate or pressure.
//
// ID:      MED-CRIT-008
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Seymour CW, Liu VX, Iwashyna TJ, et al. Assessment of clinical criteria for sepsis: for the Third International Consensus Definitions for Sepsis and Septic Shock (Sepsis-3). JAMA. 2016;315(8):762-774. doi:10.1001/jama.2016.0288. PMID 26903335. Abstract: "qSOFA score (range, 0-3 points, with 1 point each for systolic hypotension [≤100 mm Hg], tachypnea [≥22/min], or altered mentation)"
// Source:  Singer M, Deutschman CS, Seymour CW, et al. The Third International Consensus Definitions for Sepsis and Septic Shock (Sepsis-3). JAMA. 2016;315(8):801-810. doi:10.1001/jama.2016.0287. PMID 26903338. Abstract: "respiratory rate of 22/min or greater, altered mentation, or systolic blood pressure of 100 mm Hg or less"
// Tags:    critical, sepsis, qsofa, sepsis-3, score
// Units:   respirationRatePerMin breaths/min; systolicMmHg mmHg; result points 0–3
// Range:   respirationRatePerMin ≥ 0; systolicMmHg ≥ 0 (physical limits; adults with suspected infection, the source's population)
// Example: RR 22, SBP 100, altered → 3; RR 21, SBP 101, alert → 0 (from the abstract's thresholds)
func QSOFAScore(respirationRatePerMin, systolicMmHg float64, alteredMentation bool) (int, error) {
	for _, e := range [...]error{
		finite("respirationRatePerMin", respirationRatePerMin), finite("systolicMmHg", systolicMmHg),
		atLeastZero("respirationRatePerMin", respirationRatePerMin), atLeastZero("systolicMmHg", systolicMmHg),
	} {
		if e != nil {
			return 0, e
		}
	}
	s := 0
	if systolicMmHg <= 100 {
		s++
	}
	if respirationRatePerMin >= 22 {
		s++
	}
	if alteredMentation {
		s++
	}
	return s, nil
}

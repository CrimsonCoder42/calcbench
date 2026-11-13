package respiratory

// SpO2FiO2Ratio returns the SpO₂/FiO₂ (S/F) ratio: oxygen saturation by pulse oximetry, in percent, divided by
// the fraction of inspired oxygen. The source compared it with PaO₂/FiO₂ only at saturations ≤ 97 %, and the
// function enforces that limit.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range.
//
// ID:      MED-RESP-015
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Rice TW, Wheeler AP, Bernard GR, et al. Comparison of the SpO2/FIO2 ratio and the PaO2/FIO2 ratio in patients with acute lung injury or ARDS. Chest. 2007;132(2):410-417. doi:10.1378/chest.07-0617. PMID 17573487, abstract: "Corresponding measurements of Spo(2) (values </= 97%)"
// Tags:    respiratory, sf-ratio, spo2-fio2, oxygenation, ards, pulse-oximetry
// Units:   spo2Percent %; fio2Fraction fraction 0.21–1.0; result % per unit FiO₂, as the source reports it
// Range:   0 < spo2Percent ≤ 97; 0.21 ≤ fio2Fraction ≤ 1.0
// Example: SpO₂ 94 % on FiO₂ 0.4 → 235 (hand-derived; the abstract pairs S/F 235 with P/F 200)
func SpO2FiO2Ratio(spo2Percent, fio2Fraction float64) (float64, error) {
	if err := checkPositiveAtMost("spo2Percent", spo2Percent, 97); err != nil {
		return 0, err
	}
	if err := checkFiO2Fraction(fio2Fraction); err != nil {
		return 0, err
	}
	return spo2Percent / fio2Fraction, nil
}

// SpO2FiO2FromPaO2FiO2Rice returns the S/F ratio that the Rice 2007 regression S/F = 64 + 0.84 × (P/F) predicts
// for a PaO₂/FiO₂ ratio.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange for a ratio ≤ 0.
//
// ID:      MED-RESP-016
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Rice TW, Wheeler AP, Bernard GR, et al. Comparison of the SpO2/FIO2 ratio and the PaO2/FIO2 ratio in patients with acute lung injury or ARDS. Chest. 2007;132(2):410-417. doi:10.1378/chest.07-0617. PMID 17573487, abstract: "S/F = 64 + 0.84 x (P/F)"
// Tags:    respiratory, sf-ratio, pf-ratio, regression, ards
// Units:   pfRatioMmHg mmHg; result S/F (% per unit FiO₂)
// Range:   pfRatioMmHg > 0
// Example: P/F 200 → 64 + 168 = 232; P/F 300 → 316 (hand-derived from the equation; the abstract reports the fitted pairs 235 ↔ 200 and 315 ↔ 300)
func SpO2FiO2FromPaO2FiO2Rice(pfRatioMmHg float64) (float64, error) {
	if err := checkPositiveFinite("pfRatioMmHg", pfRatioMmHg); err != nil {
		return 0, err
	}
	// 0.84 × a finite value is finite and far from MaxFloat64 after adding 64: no overflow check needed.
	return 64 + 0.84*pfRatioMmHg, nil
}

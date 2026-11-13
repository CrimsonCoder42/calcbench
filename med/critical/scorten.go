package critical

// SCORTEN returns SCORTEN, 0–7: the number of its seven risk factors present, each scoring one point (the
// abstract: "for each ... score point the odds ratio was 3.45"): age above 40 y; malignancy; tachycardia above
// 120 per min; initial epidermal detachment above 10% of body surface area; serum urea above 10 mmol/L; serum
// glucose above 14 mmol/L; bicarbonate below 20 mmol/L. Every threshold is strict, as the source words it.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range.
//
// ID:      MED-CRIT-003
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Bastuji-Garin S, Fouchard N, Bertocchi M, Roujeau JC, Revuz J, Wolkenstein P. SCORTEN: a severity-of-illness score for toxic epidermal necrolysis. J Invest Dermatol. 2000;115(2):149-153. doi:10.1046/j.1523-1747.2000.00061.x, PMID 10951229, abstract
// Tags:    critical, scorten, toxic epidermal necrolysis, score
// Units:   ageYears years; heartRatePerMin beats/min; detachmentPercentBSA % of body surface area; ureaMmolL, glucoseMmolL, bicarbonateMmolL mmol/L; result points 0–7
// Range:   0 ≤ ageYears ≤ 130; heartRatePerMin ≥ 0; 0 ≤ detachmentPercentBSA ≤ 100; ureaMmolL, glucoseMmolL, bicarbonateMmolL ≥ 0 (physical limits; the abstract states no validated range)
// Example: age 41, malignancy, HR 121, detachment 11%, urea 10.1, glucose 14.1, bicarbonate 19.9 → 7; age 40, none, HR 120, 10%, 10, 14, 20 → 0 (HAND-DERIVED from the abstract's thresholds)
func SCORTEN(ageYears float64, malignancy bool, heartRatePerMin, detachmentPercentBSA, ureaMmolL, glucoseMmolL, bicarbonateMmolL float64) (int, error) {
	for _, e := range [...]error{
		finite("ageYears", ageYears), finite("heartRatePerMin", heartRatePerMin),
		finite("detachmentPercentBSA", detachmentPercentBSA), finite("ureaMmolL", ureaMmolL),
		finite("glucoseMmolL", glucoseMmolL), finite("bicarbonateMmolL", bicarbonateMmolL),
		within("ageYears", ageYears, 0, 130), atLeastZero("heartRatePerMin", heartRatePerMin),
		within("detachmentPercentBSA", detachmentPercentBSA, 0, 100), atLeastZero("ureaMmolL", ureaMmolL),
		atLeastZero("glucoseMmolL", glucoseMmolL), atLeastZero("bicarbonateMmolL", bicarbonateMmolL),
	} {
		if e != nil {
			return 0, e
		}
	}
	n := 0
	if ageYears > 40 {
		n++
	}
	if malignancy {
		n++
	}
	if heartRatePerMin > 120 {
		n++
	}
	if detachmentPercentBSA > 10 {
		n++
	}
	if ureaMmolL > 10 {
		n++
	}
	if glucoseMmolL > 14 {
		n++
	}
	if bicarbonateMmolL < 20 {
		n++
	}
	return n, nil
}

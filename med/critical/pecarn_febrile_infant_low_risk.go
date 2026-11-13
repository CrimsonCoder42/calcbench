package critical

// PECARNFebrileInfantLowRisk reports whether a febrile infant 60 days and younger meets the PECARN rule's
// low-risk criteria for serious bacterial infection: a negative urinalysis result, an absolute neutrophil
// count of 4090/µL or less, and serum procalcitonin of 1.71 ng/mL or less (all three).
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange for ageDays outside 0–60 or a negative count or
// concentration.
//
// ID:      MED-CRIT-007
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Kuppermann N, Dayan PS, Levine DA, et al. A clinical prediction rule to identify febrile infants 60 days and younger at low risk for serious bacterial infections. JAMA Pediatr. 2019;173(4):342-351. doi:10.1001/jamapediatrics.2018.5501. PMID 30776077. Abstract: "negative urinalysis result, an ANC of 4090/µL or less, and serum procalcitonin of 1.71 ng/mL or less"
// Tags:    critical, pediatrics, febrile infant, pecarn, serious bacterial infection, decision rule
// Units:   ageDays days; ancPerUL cells/µL; procalcitoninNgML ng/mL; result boolean
// Range:   0 ≤ ageDays ≤ 60 ("febrile infants 60 days and younger"); ancPerUL ≥ 0; procalcitoninNgML ≥ 0
// Example: negative urinalysis, ANC 4090, PCT 1.71 → true; ANC 4091 → false (from the abstract's thresholds)
func PECARNFebrileInfantLowRisk(ageDays float64, urinalysisPositive bool, ancPerUL, procalcitoninNgML float64) (bool, error) {
	for _, e := range [...]error{
		finite("ageDays", ageDays), finite("ancPerUL", ancPerUL), finite("procalcitoninNgML", procalcitoninNgML),
		within("ageDays", ageDays, 0, 60), atLeastZero("ancPerUL", ancPerUL),
		atLeastZero("procalcitoninNgML", procalcitoninNgML),
	} {
		if e != nil {
			return false, e
		}
	}
	return !urinalysisPositive && ancPerUL <= 4090 && procalcitoninNgML <= 1.71, nil
}

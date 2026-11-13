package chemistry

// SerumOsmolalityDorwart returns the calculated serum osmolality, in mOsm/kg, by the Dorwart–Chalmers
// equation: 1.86 × Na + glucose/18 + BUN/2.8 + 9.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range; ErrDomain when the result is beyond float64.
//
// ID:      MED-CHEM-003
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Dorwart WV, Chalmers L. Comparison of methods for calculating serum osmolality from chemical concentrations, and the prognostic value of such calculations. Clin Chem. 1975;21(2):190-194. doi:10.1093/clinchem/21.2.190. PMID 1112025. Abstract: "Osmolality = 1.86 Na + (Glu/18) + (BUN/2.8) + 9."
// Tags:    chemistry, osmolality, sodium, glucose, bun
// Units:   sodiumMmolL mmol/L; glucoseMgDL mg/dL; ureaNitrogenMgDL mg/dL (blood urea nitrogen); result mOsm/kg
// Range:   every input > 0 (715 hospital sera; limits not stated by the source)
// Example: Na 140, glucose 90, BUN 14 → 260.4 + 5 + 5 + 9 = 279.4 mOsm/kg (hand-derived)
func SerumOsmolalityDorwart(sodiumMmolL, glucoseMgDL, ureaNitrogenMgDL float64) (float64, error) {
	if err := checkOsmolality(sodiumMmolL, "glucoseMgDL", glucoseMgDL, "ureaNitrogenMgDL", ureaNitrogenMgDL); err != nil {
		return 0, err
	}
	r := 1.86*sodiumMmolL + glucoseMgDL/18 + ureaNitrogenMgDL/2.8 + 9
	if err := checkResult(r); err != nil {
		return 0, err
	}
	return r, nil
}

// SerumOsmolalityDorwartSI returns the calculated serum osmolality, in mOsm/kg, from glucose and urea in
// mmol/L: 1.86 × Na + glucose + urea + 9. This is the Dorwart–Chalmers equation (MED-CHEM-003) with the
// paper's own divisors 18 and 2.8 read as its mg/dL → mmol/L conversions for glucose and urea nitrogen;
// it is exact algebra on the paper's equation, with no other conversion factor.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range; ErrDomain when the result is beyond float64.
//
// ID:      MED-CHEM-004
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Dorwart WV, Chalmers L. Comparison of methods for calculating serum osmolality from chemical concentrations, and the prognostic value of such calculations. Clin Chem. 1975;21(2):190-194. doi:10.1093/clinchem/21.2.190. PMID 1112025. Abstract: "Osmolality = 1.86 Na + (Glu/18) + (BUN/2.8) + 9."
// Tags:    chemistry, osmolality, sodium, glucose, urea, si
// Units:   sodiumMmolL mmol/L; glucoseMmolL mmol/L; ureaMmolL mmol/L; result mOsm/kg
// Range:   every input > 0 (limits not stated by the source)
// Example: Na 140, glucose 5 mmol/L (= 90 mg/dL), urea 5 mmol/L (= BUN 14 mg/dL) → 260.4 + 5 + 5 + 9 = 279.4 mOsm/kg (hand-derived)
func SerumOsmolalityDorwartSI(sodiumMmolL, glucoseMmolL, ureaMmolL float64) (float64, error) {
	if err := checkOsmolality(sodiumMmolL, "glucoseMmolL", glucoseMmolL, "ureaMmolL", ureaMmolL); err != nil {
		return 0, err
	}
	r := 1.86*sodiumMmolL + glucoseMmolL + ureaMmolL + 9
	if err := checkResult(r); err != nil {
		return 0, err
	}
	return r, nil
}

// checkOsmolality validates the inputs of MED-CHEM-003 and -004: every finiteness check first, then ranges.
func checkOsmolality(sodiumMmolL float64, glucoseName string, glucose float64, ureaName string, urea float64) error {
	if err := checkFinite("sodiumMmolL", sodiumMmolL); err != nil {
		return err
	}
	if err := checkFinite(glucoseName, glucose); err != nil {
		return err
	}
	if err := checkFinite(ureaName, urea); err != nil {
		return err
	}
	if err := checkPositive("sodiumMmolL", sodiumMmolL); err != nil {
		return err
	}
	if err := checkPositive(glucoseName, glucose); err != nil {
		return err
	}
	return checkPositive(ureaName, urea)
}

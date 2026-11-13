package hepatic

// GlasgowAlcoholicHepatitisScore returns the Glasgow alcoholic hepatitis score (GAHS), 5–12, from Forrest
// 2005 Table 4: age < 50 years 1, ≥ 50 2; white cell count (WCC) < 15 ×10⁹/L 1, ≥ 15 2; urea < 5 mmol/L 1,
// ≥ 5 2; PT ratio < 1.5 1, 1.5–2.0 2, > 2.0 3; bilirubin < 125 µmol/L 1, 125–250 2, > 250 3. The inputs take
// the source's SI units; no band is returned.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range.
//
// ID:      MED-HEPAT-006
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Forrest EH, Evans CD, Stewart S, et al. Analysis of factors predictive of mortality in alcoholic hepatitis and derivation and validation of the Glasgow alcoholic hepatitis score. Gut. 2005;54(8):1174-1179. doi:10.1136/gut.2004.050781. PMID 16009691. PMC1774903, Table 4 (score 1 / 2 / 3 per variable)
// Tags:    hepatic, liver, alcoholic-hepatitis, gahs, glasgow, bilirubin, urea, prothrombin, score
// Units:   ageYears years; wccE9PerL ×10⁹/L; ureaMmolL mmol/L; ptRatio prothrombin time ratio (dimensionless); bilirubinUmolL µmol/L; result points 5–12
// Range:   every input > 0 (people with alcoholic hepatitis; input limits not stated in the source)
// Example: age 55, WCC 16, urea 7, PT ratio 1.8, bilirubin 300 → 2 + 2 + 2 + 2 + 3 = 11 (Table 4)
func GlasgowAlcoholicHepatitisScore(ageYears, wccE9PerL, ureaMmolL, ptRatio, bilirubinUmolL float64) (int, error) {
	vals := [...]struct {
		name string
		v    float64
	}{
		{"ageYears", ageYears}, {"wccE9PerL", wccE9PerL}, {"ureaMmolL", ureaMmolL},
		{"ptRatio", ptRatio}, {"bilirubinUmolL", bilirubinUmolL},
	}
	for _, x := range vals {
		if err := checkFinite(x.name, x.v); err != nil {
			return 0, err
		}
	}
	for _, x := range vals {
		if err := checkPositive(x.name, x.v); err != nil {
			return 0, err
		}
	}
	points := 5 // each variable scores at least 1
	if ageYears >= 50 {
		points++
	}
	if wccE9PerL >= 15 {
		points++
	}
	if ureaMmolL >= 5 {
		points++
	}
	switch {
	case ptRatio > 2.0:
		points += 2
	case ptRatio >= 1.5:
		points++
	}
	switch {
	case bilirubinUmolL > 250:
		points += 2
	case bilirubinUmolL >= 125:
		points++
	}
	return points, nil
}

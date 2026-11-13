package respiratory

// LightsCriteria returns how many of Light's three criteria a pleural effusion meets (pleural fluid/serum
// protein ratio > 0.5; pleural fluid/serum LDH ratio > 0.6; pleural fluid LDH > two-thirds of the upper limit of
// normal serum LDH) and whether it is an exudate by the rule "Exudative pleural effusions meet a minimum of one of
// the following criteria, but transudative pleural effusions meet none". The ratio edges are compared without
// division: pleural > 0.5 × serum, pleural > 0.6 × serum, 3 × pleural > 2 × upper limit.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range.
//
// ID:      MED-RESP-047
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Light RW, Macgregor MI, Luchsinger PC, Ball WC Jr. Pleural effusions: the diagnostic separation of transudates and exudates. Ann Intern Med. 1972;77(4):507-513. doi:10.7326/0003-4819-77-4-507. PMID 4642731 (primary; no abstract, full text not opened)
// Source:  Porcel JM, Porcel L, Palma R, Bielsa S. Revisiting Light's criteria: a validated blood-free triple combination matches diagnostic accuracy in over 7000 patients. ERJ Open Res. 2026;12(2):00821-2025. doi:10.1183/23120541.00821-2025. PMID 41846683, PMC12991015, Introduction (the three criteria, "at least one")
// Source:  Subedi AR, Shrestha TM, Shakya YL, Budha B. Comparative evaluation of pleural fluid cholesterol and total protein versus Light's criteria in differentiation of exudates and transudates. Clin Respir J. 2025;19(9):e70121. doi:10.1111/crj.70121. PMID 40867084, PMC12390823, Methods
// Verified: secondary (Porcel 2026 ERJ Open Res, PMC12991015; Subedi 2025 Clin Respir J, PMC12390823)
// Tags:    respiratory, pleural-effusion, lights-criteria, exudate, transudate
// Units:   pleuralProtein and serumProtein one protein unit (g/dL or g/L, the same for both); pleuralLDH, serumLDH and serumLDHUpperLimit one enzyme unit (U/L); result count 0–3 and the rule's outcome
// Range:   every input > 0
// Example: pleural protein 4.0 of serum 6.0 (0.67), LDH 150 of serum 200 (0.75), upper limit 250 → 2 criteria, exudate (hand-derived)
func LightsCriteria(pleuralProtein, serumProtein, pleuralLDH, serumLDH, serumLDHUpperLimit float64) (met int, exudate bool, err error) {
	for _, v := range [5]struct {
		name string
		x    float64
	}{{"pleuralProtein", pleuralProtein}, {"serumProtein", serumProtein}, {"pleuralLDH", pleuralLDH}, {"serumLDH", serumLDH}, {"serumLDHUpperLimit", serumLDHUpperLimit}} {
		if err := checkPositiveFinite(v.name, v.x); err != nil {
			return 0, false, err
		}
	}
	if pleuralProtein > 0.5*serumProtein {
		met++
	}
	if pleuralLDH > 0.6*serumLDH {
		met++
	}
	if 3*pleuralLDH > 2*serumLDHUpperLimit {
		met++
	}
	return met, met >= 1, nil
}

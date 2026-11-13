package endocrine

// CPeptideGlucoseRatio returns the C-peptide to glucose ratio (CPR), dimensionless as the source states it:
// C-peptide (ng/mL) / glucose (mg/dL) × 100. The source defines it for paired postprandial samples; the
// caller supplies the pair. Not rounded.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range; ErrDomain when the result is beyond float64.
//
// ID:      MED-ENDO-004
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Saisho Y. Postprandial C-peptide to glucose ratio as a marker of β cell function: implication for the management of type 2 diabetes. Int J Mol Sci. 2016;17(5):744. doi:10.3390/ijms17050744. PMID 27196896. PMC4881566, section 4: "C-peptide (ng/mL)/glucose (mg/dL) (×100)"
// Tags:    endocrine, diabetes, c-peptide, glucose, beta-cell
// Units:   cPeptideNgML ng/mL; glucoseMgDL mg/dL; result C-peptide (ng/mL) per glucose (mg/dL) × 100
// Range:   both > 0 (input limits not stated in the source)
// Example: C-peptide 3.0 ng/mL, glucose 150 mg/dL → 3.0 / 150 × 100 = 2.0 (HAND-DERIVED from section 4)
func CPeptideGlucoseRatio(cPeptideNgML, glucoseMgDL float64) (float64, error) {
	if err := checkFinite("cPeptideNgML", cPeptideNgML); err != nil {
		return 0, err
	}
	if err := checkFinite("glucoseMgDL", glucoseMgDL); err != nil {
		return 0, err
	}
	if err := checkPositive("cPeptideNgML", cPeptideNgML); err != nil {
		return 0, err
	}
	if err := checkPositive("glucoseMgDL", glucoseMgDL); err != nil {
		return 0, err
	}
	r := cPeptideNgML / glucoseMgDL * 100
	if err := checkResult(r); err != nil {
		return 0, err
	}
	return r, nil
}

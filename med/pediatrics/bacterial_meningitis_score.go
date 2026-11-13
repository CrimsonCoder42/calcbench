package pediatrics

// BacterialMeningitisScore returns the Bacterial Meningitis Score, 0–6: 2 points for a positive CSF Gram
// stain and 1 point for each of CSF absolute neutrophil count ≥ 1000 cells/mm³, CSF protein ≥ 80 mg/dL,
// peripheral blood absolute neutrophil count ≥ 10,000 cells/mm³ and seizure before or at the time of
// presentation. The 2007 validation calls a child with none of the criteria (score 0) "very low risk"; that
// reading is the caller's and is not returned.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range.
//
// ID:      MED-PEDS-007
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Nigrovic LE, Kuppermann N, Malley R. Development and validation of a multivariable predictive model to distinguish bacterial from aseptic meningitis in children in the post-Haemophilus influenzae era. Pediatrics. 2002;110(4):712-719. doi:10.1542/peds.110.4.712. PMID 12359784, abstract ("attributing 2 points for a positive Gram stain and 1 point for each of the other variables"; the five predictors and thresholds)
// Source:  Nigrovic LE, Kuppermann N, Macias CG, et al. Clinical prediction rule for identifying children with cerebrospinal fluid pleocytosis at very low risk of bacterial meningitis. JAMA. 2007;297(1):52-60. doi:10.1001/jama.297.1.52. PMID 17200475, abstract (validation; population; risk "very low (0.1%)" with none of the criteria)
// Tags:    pediatrics, meningitis, cerebrospinal fluid, bacterial meningitis score, nigrovic, score
// Units:   csfNeutrophilsPerMm3, bloodNeutrophilsPerMm3 cells/mm³ (= cells/µL); csfProteinMgDL mg/dL; findings present/absent; result points 0–6
// Range:   children aged 29 days to 19 years with CSF pleocytosis (CSF white cells ≥ 10 cells/µL), the 2007 population (not checked: no age input); 0 ≤ csfNeutrophilsPerMm3 ≤ 1e6, 0 ≤ csfProteinMgDL ≤ 5000, 0 ≤ bloodNeutrophilsPerMm3 ≤ 1e6 (physical limits)
// Example: Gram stain positive, CSF ANC 1000, protein 79.9, blood ANC 9999, no seizure → 3 (2 + 1); nothing present → 0
func BacterialMeningitisScore(csfGramStainPositive bool, csfNeutrophilsPerMm3, csfProteinMgDL, bloodNeutrophilsPerMm3 float64, seizureAtOrBeforePresentation bool) (int, error) {
	if err := checkFloat("csfNeutrophilsPerMm3", csfNeutrophilsPerMm3, 0, 1e6); err != nil {
		return 0, err
	}
	if err := checkFloat("csfProteinMgDL", csfProteinMgDL, 0, 5000); err != nil {
		return 0, err
	}
	if err := checkFloat("bloodNeutrophilsPerMm3", bloodNeutrophilsPerMm3, 0, 1e6); err != nil {
		return 0, err
	}
	score := 0
	if csfGramStainPositive {
		score += 2
	}
	for _, present := range [...]bool{
		csfNeutrophilsPerMm3 >= 1000, csfProteinMgDL >= 80, bloodNeutrophilsPerMm3 >= 10000, seizureAtOrBeforePresentation,
	} {
		if present {
			score++
		}
	}
	return score, nil
}

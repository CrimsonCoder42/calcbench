package hematology

// mysecPointsPerYear is MYSEC-PM's age weight: "0.15 points to any year of age" (Passamonti 2017, abstract).
const mysecPointsPerYear = 0.15

// MYSECPMScore returns the MYSEC-PM score for myelofibrosis secondary to polycythemia vera or essential
// thrombocythemia: 2 points each for hemoglobin < 11 g/dL, circulating blasts ≥ 3% and a CALR-unmutated
// genotype; 1 point each for platelets < 150 × 10⁹/L and constitutional symptoms; plus 0.15 points per year
// of age. The paper's risk-group cut-offs on this score are not built here (not in the abstract).
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range.
//
// ID:      MED-HEMAT-008
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Passamonti F, Giorgino T, Mora B, et al. A clinical-molecular prognostic model to predict survival in patients with post polycythemia vera and post essential thrombocythemia myelofibrosis. Leukemia. 2017;31(12):2726-2731. doi:10.1038/leu.2017.169, PMID 28561069, abstract
// Tags:    hematology, myelofibrosis, mysec-pm, prognostic score
// Units:   ageYears years; hemoglobinGDL g/dL; blastsPercent %; plateletsE9PerL × 10⁹/L; result points
// Range:   0 ≤ ageYears ≤ 130; 0 ≤ hemoglobinGDL ≤ 30; 0 ≤ blastsPercent ≤ 100; 0 ≤ plateletsE9PerL ≤ 5000 (physical limits)
// Example: age 60, Hb 10.9, blasts 3, CALR-unmutated, platelets 149, symptoms → 2+2+2+1+1+9 = 17; age 60, none → 9 (HAND-DERIVED from the abstract)
func MYSECPMScore(ageYears, hemoglobinGDL, blastsPercent float64, calrUnmutated bool, plateletsE9PerL float64, constitutionalSymptoms bool) (float64, error) {
	for _, e := range [...]error{
		checkFloat("ageYears", ageYears, 0, 130),
		checkFloat("hemoglobinGDL", hemoglobinGDL, 0, 30),
		checkFloat("blastsPercent", blastsPercent, 0, 100),
		checkFloat("plateletsE9PerL", plateletsE9PerL, 0, 5000),
	} {
		if e != nil {
			return 0, e
		}
	}
	p := pts(hemoglobinGDL < 11, 2) + pts(blastsPercent >= 3, 2) + pts(calrUnmutated, 2) +
		pts(plateletsE9PerL < 150, 1) + pts(constitutionalSymptoms, 1)
	return float64(p) + mysecPointsPerYear*ageYears, nil
}

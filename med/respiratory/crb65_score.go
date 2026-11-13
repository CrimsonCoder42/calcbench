package respiratory

import "math"

// CRB65Score returns the CRB-65 score (0–4) for community-acquired pneumonia, CURB-65 without the urea item:
// one point each for confusion, respiratory rate ≥ 30/min, systolic blood pressure < 90 mmHg or diastolic
// ≤ 60 mmHg, and age ≥ 65 years.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range.
//
// ID:      MED-RESP-045
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Lim WS, van der Eerden MM, Laing R, et al. Defining community acquired pneumonia severity on presentation to hospital: an international derivation and validation study. Thorax. 2003;58(5):377-382. doi:10.1136/thorax.58.5.377. PMID 12728155, PMC1746657 (primary named by both secondaries; the text reached defines CURB-65, not CRB-65)
// Source:  Dwyer R, Hedlund J, Henriques-Normark B, Kalin M. Improvement of CRB-65 as a prognostic tool in adult patients with community-acquired pneumonia. BMJ Open Respir Res. 2014;1(1):e000038. doi:10.1136/bmjresp-2014-000038. PMID 25478185, PMC4212804, Methods (CRB-65 criteria)
// Source:  Chen YX, Wang JY, Guo SB. Use of CRB-65 and quick Sepsis-related Organ Failure Assessment to predict site of care and mortality in pneumonia patients in the emergency department: a retrospective study. Crit Care. 2016;20(1):167. doi:10.1186/s13054-016-1351-0. PMID 27250351, PMC4888495, Methods ("CRB-65 (confusion, respiratory rate ≥30/minute, systolic blood pressure <90 mmHg or diastolic blood pressure ≤60 mmHg, age ≥65 years)")
// Verified: secondary (Dwyer 2014 BMJ Open Respir Res, PMC4212804; Chen 2016 Crit Care, PMC4888495)
// Tags:    respiratory, pneumonia, cap, crb-65, severity, score
// Units:   respiratoryRatePerMin breaths/min; systolicMmHg, diastolicMmHg mmHg; ageYears years; result points 0–4
// Range:   respiratoryRatePerMin ≥ 0; systolicMmHg > 0; diastolicMmHg > 0; ageYears ≥ 0 (adults with pneumonia were studied; no limits stated)
// Example: confused, RR 32, BP 110/70, age 70 → confusion 1 + RR 1 + age 1 = 3 (hand-derived from the criteria)
func CRB65Score(confusion bool, respiratoryRatePerMin, systolicMmHg, diastolicMmHg, ageYears float64) (int, error) {
	if err := checkClosedRange("respiratoryRatePerMin", respiratoryRatePerMin, 0, math.Inf(1)); err != nil {
		return 0, err
	}
	if err := checkPositiveFinite("systolicMmHg", systolicMmHg); err != nil {
		return 0, err
	}
	if err := checkPositiveFinite("diastolicMmHg", diastolicMmHg); err != nil {
		return 0, err
	}
	if err := checkClosedRange("ageYears", ageYears, 0, math.Inf(1)); err != nil {
		return 0, err
	}
	n := 0
	if confusion {
		n++
	}
	if respiratoryRatePerMin >= 30 {
		n++
	}
	if systolicMmHg < 90 || diastolicMmHg <= 60 {
		n++
	}
	if ageYears >= 65 {
		n++
	}
	return n, nil
}

package cardio

// HEARScore returns the HEAR score, 0–8: the HEART score without its troponin element (History, ECG, Age,
// Risk factors, scored as in HEARTScore). Smith 2020 reports HEAR scores of 0 and 1 as its very-low-risk group.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range.
//
// ID:      MED-CARDIO-059
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Smith LM, Ashburn NP, Snavely AC, et al. Identification of very low-risk acute chest pain patients without troponin testing. Emerg Med J. 2020;37(11):690-695. doi:10.1136/emermed-2020-209698. PMID 32753395. PMC7952041, Methods (HEAR = HEART without troponin) and Results (HEAR 0–1 very low risk)
// Source:  Mahler SA, et al. Circ Cardiovasc Qual Outcomes. 2015;8(2):195-203. doi:10.1161/CIRCOUTCOMES.114.001384. PMID 25737484. PMC4413911, Table 2 (element points)
// Tags:    cardio, chest-pain, acs, hear-score, heart-score, risk-score
// Units:   history, ecg, riskFactors points 0–2; ageYears years; result points 0–8
// Range:   each rated item 0–2; ageYears 18–110
// Example: history 0, ECG 0, age 40, risk factors 1 → 1 (HAND-DERIVED)
func HEARScore(history, ecg int, ageYears float64, riskFactors int) (int, error) {
	return hearPoints(history, ecg, ageYears, riskFactors)
}

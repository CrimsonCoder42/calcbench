package respiratory

import "fmt"

// COPDAssessmentTestScore returns the COPD Assessment Test (CAT) total (0–40), the sum of its eight item scores:
// "Item scores range from 0 to 5 points resulting in a CAT total score ranging from 0 to 40 points". The item
// wording is not embedded (licence: GlaxoSmithKline, the CAT's copyright holder); items[i] is the caller's
// 0–5 answer to item i+1.
// Errors: ErrOutOfRange for an item outside 0–5.
//
// ID:      MED-RESP-048
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Jones PW, Harding G, Berry P, Wiklund I, Chen WH, Kline Leidy N. Development and first validation of the COPD Assessment Test. Eur Respir J. 2009;34(3):648-654. doi:10.1183/09031936.00102509. PMID 19720809 (primary; abstract: "eight items", "40-point scale"; full text not opened)
// Source:  Nishimura K, Kusunose M, Shibayama A, Nakayasu K. Patient Relat Outcome Meas. 2026, PMC13596913 ("eight items rated from 0 to 5 … Total scores range from 0 to 40"; ref 21 Jones 2009)
// Source:  Houben-Wilke S, Janssen DJA, Franssen FME, Vanfleteren LEGW, et al. Contribution of individual COPD assessment test (CAT) items to CAT total score and effects of pulmonary rehabilitation on CAT scores. Health Qual Life Outcomes. 2018;16(1):205. doi:10.1186/s12955-018-1034-4. PMID 30376861, PMC6208036, Methods
// Verified: secondary (Houben-Wilke 2018 PMC6208036, ref 2 Jones 2009; Nishimura 2026 PMC13596913, ref 21 Jones 2009)
// Tags:    respiratory, copd, cat, health-status, questionnaire, score
// Units:   items eight integers 0–5; result points 0–40
// Range:   each item 0–5
// Example: items 2, 3, 1, 4, 2, 1, 3, 2 → 18 (hand-derived: the sum)
func COPDAssessmentTestScore(items [8]int) (int, error) {
	s := 0
	for i, v := range items {
		if v < 0 || v > 5 {
			return 0, fmt.Errorf("%w: items[%d] = %d, want 0 ≤ item ≤ 5", ErrOutOfRange, i, v)
		}
		s += v
	}
	return s, nil
}

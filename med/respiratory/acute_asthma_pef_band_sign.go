package respiratory

import "math"

// AcuteAsthmaPEFBand is a peak-flow cell of SIGN 158 Table 15 "Levels of severity of acute asthma attacks in
// adults"; its value is the table's own cell text, verbatim.
type AcuteAsthmaPEFBand string

// The three PEF cells of SIGN 158 (November 2024, corrections May 2025) Table 15, verbatim. Each is one
// criterion of its row, never the whole row: "Moderate acute asthma" also needs "Increasing symptoms" and
// "No features of acute severe asthma"; "Acute severe asthma" is "Any one of" four features, PEF being one;
// "Life-threatening asthma" is "SpO2 <92% plus any one of the following", PEF being one of the
// measurements.
const (
	AcuteAsthmaPEFModerate        AcuteAsthmaPEFBand = "PEF >50–75% best or predicted"
	AcuteAsthmaPEFAcuteSevere     AcuteAsthmaPEFBand = "PEF 33–50% best or predicted"
	AcuteAsthmaPEFLifeThreatening AcuteAsthmaPEFBand = "PEF <33% best or predicted"
)

// Level returns the Table 15 row that holds the cell, verbatim ("Moderate acute asthma", "Acute severe
// asthma", "Life-threatening asthma"), or "" for a value that is not one of the three constants. The row
// name is where the cell sits, not a level by itself (see the constants).
func (b AcuteAsthmaPEFBand) Level() string {
	switch b {
	case AcuteAsthmaPEFModerate:
		return "Moderate acute asthma"
	case AcuteAsthmaPEFAcuteSevere:
		return "Acute severe asthma"
	case AcuteAsthmaPEFLifeThreatening:
		return "Life-threatening asthma"
	}
	return ""
}

// AcuteAsthmaPEFBandSIGN returns the PEF cell of SIGN 158 Table 15 (adults) that a peak expiratory flow,
// in percent of best or predicted, falls in: < 33 "PEF <33% best or predicted"; 33 to 50 inclusive "PEF
// 33–50% best or predicted"; above 50 to 75 inclusive "PEF >50–75% best or predicted"; above 75 "" (the
// table has no PEF cell above 75 %). SIGN 158 §9.2.1 takes predicted PEF only when the recent best PEF
// (within two years) is unknown.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange for pefPercentBestOrPredicted < 0 or ageYears < 18.
//
// ID:      MED-RESP-054
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Scottish Intercollegiate Guidelines Network, British Thoracic Society. SIGN 158: British guideline on the management of asthma. Revised edition November 2024, corrections May 2025. https://www.sign.ac.uk/media/2342/sign-158-2024-british-guideline-on-the-management-of-asthma-may-2025.pdf (read 2026-11-12), Table 15 "Levels of severity of acute asthma attacks in adults" and §9.2.1
// Tags:    respiratory, asthma, peak-flow, pef, acute-asthma, severity, sign, bts
// Units:   pefPercentBestOrPredicted percent of the best (or, if unknown, predicted) PEF; ageYears years (decimal); result the table's PEF cell text or ""
// Range:   pefPercentBestOrPredicted ≥ 0; ageYears ≥ 18 (Table 15 is for adults; children have Table 17, not ported here)
// Example: 32.9 → "PEF <33% best or predicted"; 33 and 50 → "PEF 33–50% best or predicted"; 50.1 and 75 → "PEF >50–75% best or predicted"; 75.1 → ""
func AcuteAsthmaPEFBandSIGN(pefPercentBestOrPredicted, ageYears float64) (AcuteAsthmaPEFBand, error) {
	if err := checkClosedRange("pefPercentBestOrPredicted", pefPercentBestOrPredicted, 0, math.Inf(1)); err != nil {
		return "", err
	}
	if err := checkClosedRange("ageYears", ageYears, 18, math.Inf(1)); err != nil {
		return "", err
	}
	switch p := pefPercentBestOrPredicted; {
	case p < 33:
		return AcuteAsthmaPEFLifeThreatening, nil
	case p <= 50:
		return AcuteAsthmaPEFAcuteSevere, nil
	case p <= 75:
		return AcuteAsthmaPEFModerate, nil
	}
	return "", nil
}

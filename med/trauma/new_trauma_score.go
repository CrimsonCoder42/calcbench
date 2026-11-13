package trauma

import (
	"fmt"
	"math"
)

// NewTraumaScoreResult holds the coded values and the three quantities Jeong 2017 defines from them.
type NewTraumaScoreResult struct {
	SBPCode             int     // 0–4, Table 1
	SpO2Code            int     // 0–4, Table 1
	NTS                 float64 // 0.4006 × GCS + 0.2983 × SBP code + 0.8709 × SpO2 code
	TriageNTS           float64 // GCS + SBP code + SpO2 code, 3–23
	SurvivalProbability float64 // 1 / (1 + e^−b), b = −6.5406 + NTS
}

// NewTraumaScore returns the New Trauma Score (NTS) of Jeong 2017: the actual Glasgow Coma Scale total (3–15)
// in place of a GCS code, systolic blood pressure coded 4 (110–149), 3 (≥150), 2 (90–109), 1 (70–89), 0 (<70),
// and SpO2 coded 4 (≥94), 3 (80–93), 2 (60–79), 1 (40–59), 0 (<40); the weighted NTS, the unweighted triage
// sum T-NTS, and the probability of survival.
// Errors: ErrOutOfRange for gcs outside 3–15, systolicMmHg < 0 or spo2Percent outside 0–100.
//
// ID:      MED-TRAUMA-006
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Jeong JH, Park YJ, Kim DH, et al. The new trauma score (NTS): a modification of the revised trauma score for better trauma mortality prediction. BMC Surg. 2017;17(1):77. doi:10.1186/s12893-017-0272-4. PMID 28673278, PMC5496419, Table 1 "Modification of the Revised Trauma Score" (codes) and the equations "NTS = (0.4006 x GCS) + (0.2983 x SBP_NTS) + (0.8709 x SpO2_NTS)", "Ps = 1 / (1 + e^-b), b = −6.5406 + (0.4006 x GCS) + (0.2983 x SBP_NTS) + (0.8709 x SpO2_NTS)", "T-NTS = GCS + SBP_NTS + SpO2_NTS"
// Tags:    trauma, new trauma score, nts, revised trauma score, survival probability
// Units:   gcs points 3–15; systolicMmHg mmHg (integer, as Table 1 prints); spo2Percent % (integer); NTS and TriageNTS points; SurvivalProbability fraction 0–1
// Range:   3 ≤ gcs ≤ 15; systolicMmHg ≥ 0; 0 ≤ spo2Percent ≤ 100
// Example: GCS 15, SBP 120, SpO2 98 → codes 4 and 4, NTS = 6.009 + 1.1932 + 3.4836 = 10.6858, T-NTS 23, Ps = 1/(1+e^−4.1452) = 0.98440 (HAND-DERIVED from the equations)
func NewTraumaScore(gcs, systolicMmHg, spo2Percent int) (NewTraumaScoreResult, error) {
	if err := checkGCS(gcs); err != nil {
		return NewTraumaScoreResult{}, err
	}
	if systolicMmHg < 0 {
		return NewTraumaScoreResult{}, fmt.Errorf("%w: systolicMmHg = %d, want ≥ 0", ErrOutOfRange, systolicMmHg)
	}
	if spo2Percent < 0 || spo2Percent > 100 {
		return NewTraumaScoreResult{}, fmt.Errorf("%w: spo2Percent = %d, want 0–100", ErrOutOfRange, spo2Percent)
	}
	var r NewTraumaScoreResult
	switch s := systolicMmHg; {
	case s < 70:
		r.SBPCode = 0
	case s <= 89:
		r.SBPCode = 1
	case s <= 109:
		r.SBPCode = 2
	case s <= 149:
		r.SBPCode = 4
	default: // ≥ 150
		r.SBPCode = 3
	}
	switch p := spo2Percent; {
	case p < 40:
		r.SpO2Code = 0
	case p <= 59:
		r.SpO2Code = 1
	case p <= 79:
		r.SpO2Code = 2
	case p <= 93:
		r.SpO2Code = 3
	default: // ≥ 94
		r.SpO2Code = 4
	}
	r.NTS = 0.4006*float64(gcs) + 0.2983*float64(r.SBPCode) + 0.8709*float64(r.SpO2Code)
	r.TriageNTS = float64(gcs + r.SBPCode + r.SpO2Code)
	r.SurvivalProbability = 1 / (1 + math.Exp(-(-6.5406 + r.NTS)))
	return r, nil
}

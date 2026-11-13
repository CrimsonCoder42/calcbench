package trauma

import "fmt"

// RevisedTraumaScoreResult holds the coded values of the Revised Trauma Score and the two totals made from them.
type RevisedTraumaScoreResult struct {
	GCSCode   int     // 0–4
	SBPCode   int     // 0–4
	RRCode    int     // 0–4
	TriageRTS int     // GCS code + SBP code + RR code, 0–12
	RTS       float64 // 0.9368 × GCS code + 0.7326 × SBP code + 0.2908 × RR code, 0–7.8408
}

// RevisedTraumaScore returns the Revised Trauma Score of Champion 1989: the Glasgow Coma Scale total coded
// 4 (13–15), 3 (9–12), 2 (6–8), 1 (4–5), 0 (3); systolic blood pressure coded 4 (> 89), 3 (76–89), 2 (50–75),
// 1 (1–49), 0 (0); respiratory rate coded 4 (10–29), 3 (> 29), 2 (6–9), 1 (1–5), 0 (0); the unweighted sum
// of the codes (T-RTS, 0–12) and the weighted RTS = 0.9368 × GCS code + 0.7326 × SBP code + 0.2908 × RR code.
// Errors: ErrOutOfRange for gcs outside 3–15, systolicMmHg < 0 or respiratoryRatePerMin < 0.
//
// ID:      MED-TRAUMA-020
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Champion HR, Sacco WJ, Copes WS, Gann DS, Gennarelli TA, Flanagan ME. A revision of the Trauma Score. J Trauma. 1989;29(5):623-629. doi:10.1097/00005373-198905000-00017. PMID 2657085 (primary; abstract names GCS, SBP and RR; full text not open)
// Source:  Yousefi et al. BMC Emerg Med. 2024;24:163. PMC11382384, Table 1 (coded values) and the equation "RTS = (0.9368 × GCS value) + (0.7326 × SBP value) + (0.2908 × RR value)"; cites Champion 1989
// Source:  Yu et al. BMJ Open. 2021;11(3):e041882. PMC7959230, Table 1 (coded values), "The range of scores is 0–12"; cites Champion 1989 (ref 13)
// Source:  Heiranizadeh et al. 2026. PMC13045056, "RTS = (GCS code × 0.9368) + (SBP code × 0.7326) + (RR code × 0.2908)"; cites Champion 1989
// Verified: secondary (Yousefi 2024 PMC11382384 Table 1; Yu 2021 PMC7959230 Table 1; coefficients also Heiranizadeh 2026 PMC13045056)
// Tags:    trauma, revised trauma score, rts, t-rts, triage, gcs
// Units:   gcs points 3–15; systolicMmHg mmHg (integer, as Table 1 prints); respiratoryRatePerMin breaths/min (integer); codes and TriageRTS points; RTS points 0–7.8408
// Range:   3 ≤ gcs ≤ 15; systolicMmHg ≥ 0; respiratoryRatePerMin ≥ 0
// Example: GCS 15, SBP 120, RR 16 → codes 4, 4, 4, T-RTS 12, RTS 3.7472 + 2.9304 + 1.1632 = 7.8408 (HAND-DERIVED from the equation)
func RevisedTraumaScore(gcs, systolicMmHg, respiratoryRatePerMin int) (RevisedTraumaScoreResult, error) {
	if err := checkGCS(gcs); err != nil {
		return RevisedTraumaScoreResult{}, err
	}
	if systolicMmHg < 0 {
		return RevisedTraumaScoreResult{}, fmt.Errorf("%w: systolicMmHg = %d, want ≥ 0", ErrOutOfRange, systolicMmHg)
	}
	if respiratoryRatePerMin < 0 {
		return RevisedTraumaScoreResult{}, fmt.Errorf("%w: respiratoryRatePerMin = %d, want ≥ 0", ErrOutOfRange, respiratoryRatePerMin)
	}
	var r RevisedTraumaScoreResult
	switch {
	case gcs >= 13:
		r.GCSCode = 4
	case gcs >= 9:
		r.GCSCode = 3
	case gcs >= 6:
		r.GCSCode = 2
	case gcs >= 4:
		r.GCSCode = 1
	default: // 3
		r.GCSCode = 0
	}
	switch s := systolicMmHg; {
	case s > 89:
		r.SBPCode = 4
	case s >= 76:
		r.SBPCode = 3
	case s >= 50:
		r.SBPCode = 2
	case s >= 1:
		r.SBPCode = 1
	default: // 0
		r.SBPCode = 0
	}
	switch rr := respiratoryRatePerMin; {
	case rr > 29:
		r.RRCode = 3
	case rr >= 10:
		r.RRCode = 4
	case rr >= 6:
		r.RRCode = 2
	case rr >= 1:
		r.RRCode = 1
	default: // 0
		r.RRCode = 0
	}
	r.TriageRTS = r.GCSCode + r.SBPCode + r.RRCode
	r.RTS = 0.9368*float64(r.GCSCode) + 0.7326*float64(r.SBPCode) + 0.2908*float64(r.RRCode)
	return r, nil
}

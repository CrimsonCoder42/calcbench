package critical

import "fmt"

// SICResult holds the points of each SIC item, their total and whether the SIC criterion is met.
type SICResult struct {
	PTINRPoints    int  // 0–2
	PlateletPoints int  // 0–2
	SOFAPoints     int  // 0–2
	Total          int  // 0–6
	SIC            bool // "total score is 4 or more with total score of prothrombin time and coagulation exceeding 2"
}

// SICScore returns the sepsis-induced coagulopathy (SIC) score of Iba 2017 Table 3: PT-INR ≤ 1.2 0, > 1.2 1,
// > 1.4 2; platelet count ≥ 150 × 10⁹/L 0, < 150 1, < 100 2; total SOFA (the sum of the respiratory,
// cardiovascular, hepatic and renal SOFA items) 0 → 0, 1 → 1, ≥ 2 → 2; and the SIC criterion (total ≥ 4 with
// PT-INR plus platelet points exceeding 2).
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange for a PT-INR not above zero, a negative platelet count, or
// a four-item SOFA sum outside 0–16.
//
// ID:      MED-CRIT-011
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Iba T, Nisio MD, Levy JH, Kitamura N, Thachil J. New criteria for sepsis-induced coagulopathy (SIC) following the revised sepsis definition: a retrospective analysis of a nationwide survey. BMJ Open. 2017;7(9):e017046. doi:10.1136/bmjopen-2017-017046. PMID 28963294, PMC5623518, Table 3 "Scoring for the diagnosis of sepsis-induced coagulopathy" and its footnotes ("Diagnosed as sepsis-induced coagulopathy when the total score is 4 or more with total score of prothrombin time and coagulation exceeding 2"; "Total SOFA is the sum the four items (respiratory SOFA, cardiovascular SOFA, hepatic SOFA, renal SOFA). The score of total SOFA is defined as 2 if the total score exceeded 2.")
// Tags:    critical, sepsis, coagulopathy, sic, score
// Units:   ptINR ratio; plateletsE9PerL 10⁹/L; sofaFourItems points 0–16; result points 0–6
// Range:   ptINR > 0; plateletsE9PerL ≥ 0; 0 ≤ sofaFourItems ≤ 16 (four SOFA items of 0–4 each)
// Example: INR 1.5, platelets 90, SOFA 3 → 2+2+2 = 6, SIC; INR 1.3, platelets 120, SOFA 2 → 1+1+2 = 4, PT+platelet 2, not SIC (HAND-DERIVED from Table 3)
func SICScore(ptINR, plateletsE9PerL float64, sofaFourItems int) (SICResult, error) {
	for _, e := range [...]error{
		finite("ptINR", ptINR), finite("plateletsE9PerL", plateletsE9PerL), atLeastZero("plateletsE9PerL", plateletsE9PerL),
	} {
		if e != nil {
			return SICResult{}, e
		}
	}
	if !(ptINR > 0) {
		return SICResult{}, fmt.Errorf("%w: ptINR = %v, want > 0", ErrOutOfRange, ptINR)
	}
	if sofaFourItems < 0 || sofaFourItems > 16 {
		return SICResult{}, fmt.Errorf("%w: sofaFourItems = %d, want 0–16", ErrOutOfRange, sofaFourItems)
	}
	var r SICResult
	switch {
	case ptINR > 1.4:
		r.PTINRPoints = 2
	case ptINR > 1.2:
		r.PTINRPoints = 1
	}
	switch {
	case plateletsE9PerL < 100:
		r.PlateletPoints = 2
	case plateletsE9PerL < 150:
		r.PlateletPoints = 1
	}
	r.SOFAPoints = min(sofaFourItems, 2)
	r.Total = r.PTINRPoints + r.PlateletPoints + r.SOFAPoints
	r.SIC = r.Total >= 4 && r.PTINRPoints+r.PlateletPoints > 2
	return r, nil
}

package cardio

import "fmt"

// GillmoreStage is an ATTR-CM stage of Gillmore 2018 (1, 2 or 3 for stages I, II, III).
type GillmoreStage int

// GillmoreATTRStage returns the Gillmore (National Amyloidosis Centre) stage of transthyretin cardiac
// amyloidosis: stage I NT-proBNP ≤ 3000 ng/L and eGFR ≥ 45 mL/min; stage III NT-proBNP > 3000 ng/L and
// eGFR < 45 mL/min; stage II otherwise.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange when an input is negative.
//
// ID:      MED-CARDIO-036
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Gillmore JD, Damy T, Fontana M, et al. A new staging system for cardiac transthyretin amyloidosis. Eur Heart J. 2018;39(30):2799-2806. doi:10.1093/eurheartj/ehx589. PMID 29048471. Abstract: "Stage I was defined as NT-proBNP ≤3000 ng/L and eGFR ≥45 ml/min, Stage III was defined as NT-proBNP >3000 ng/L and eGFR <45 ml/min", the remainder stage II
// Tags:    cardio, amyloidosis, attr-cm, staging, gillmore
// Units:   ntProBNPNgL ng/L; eGFRMlMin mL/min (as the abstract writes it); result stage 1–3
// Range:   inputs ≥ 0 (869 persons with ATTR-CM, UK National Amyloidosis Centre)
// Example: NT-proBNP 3000, eGFR 45 → 1; NT-proBNP 3001, eGFR 44.9 → 3; NT-proBNP 3001, eGFR 60 → 2
func GillmoreATTRStage(ntProBNPNgL, eGFRMlMin float64) (GillmoreStage, error) {
	for _, p := range [...]struct {
		n string
		v float64
	}{{"ntProBNPNgL", ntProBNPNgL}, {"eGFRMlMin", eGFRMlMin}} {
		if err := checkFinite(p.n, p.v); err != nil {
			return 0, err
		}
		if p.v < 0 {
			return 0, fmt.Errorf("%w: %s = %v, want ≥ 0", ErrOutOfRange, p.n, p.v)
		}
	}
	switch {
	case ntProBNPNgL <= 3000 && eGFRMlMin >= 45:
		return 1, nil
	case ntProBNPNgL > 3000 && eGFRMlMin < 45:
		return 3, nil
	default:
		return 2, nil
	}
}

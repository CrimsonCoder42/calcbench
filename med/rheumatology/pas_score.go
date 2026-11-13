package rheumatology

import (
	"fmt"
	"math"
)

// PASScore returns the PAS (Wolfe 2005 activity scale), 0–9.997: "(HAQ × 3.33 + pain VAS + PtGA VAS)/3".
// With the HAQ score it is PAS; with the HAQ-II score in the same place it is PAS-II (the formula is the
// same). The HAQ or HAQ-II score itself is computed by the caller from its own instrument.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange for HAQ outside 0–3 or a VAS outside 0–10.
//
// ID:      MED-RHEUM-009
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Wolfe F, Michaud K, Pincus T. A composite disease activity scale for clinical practice, observational studies, and clinical trials: the patient activity scale (PAS/PAS-II). J Rheumatol. 2005;32(12):2410-2415. PMID 16331773 (primary; not opened)
// Source:  Anderson J, Caplan L, Yazdany J, et al. Rheumatoid arthritis disease activity measures: American College of Rheumatology recommendations for use in clinical practice. Arthritis Care Res (Hoboken). 2012;64(5):640-647. doi:10.1002/acr.21649, PMID 22473918, PMC4028066, Table 1 and text
// Source:  England BR, Tiong BK, Bergman MJ, et al. 2019 Update of the American College of Rheumatology Recommended Rheumatoid Arthritis Disease Activity Measures. Arthritis Care Res (Hoboken). 2019;71(12):1540-1555. doi:10.1002/acr.24042, PMID 31709779, PMC6884664, Table 1
// Verified: secondary (Anderson 2012 PMC4028066; England 2019 PMC6884664; also Salaffi 2021 PMC8477094)
// Tags:    rheumatology, rheumatoid arthritis, pas, pas-ii, haq, disease activity
// Units:   haqScore HAQ or HAQ-II points 0–3; painVAS and globalVAS 0–10; result PAS points
// Range:   haqScore 0–3; painVAS 0–10; globalVAS 0–10
// Example: HAQ 1.0, pain 5, global 4 → (3.33 + 5 + 4)/3 = 4.11 (hand-derived)
func PASScore(haqScore, painVAS, globalVAS float64) (float64, error) {
	for _, in := range [...]struct {
		name string
		v    float64
		hi   float64
	}{{"haqScore", haqScore, 3}, {"painVAS", painVAS, 10}, {"globalVAS", globalVAS, 10}} {
		if math.IsNaN(in.v) || math.IsInf(in.v, 0) {
			return 0, fmt.Errorf("%w: %s = %v", ErrNotFinite, in.name, in.v)
		}
		if in.v < 0 || in.v > in.hi {
			return 0, fmt.Errorf("%w: %s = %v, want 0–%v", ErrOutOfRange, in.name, in.v, in.hi)
		}
	}
	return (haqScore*3.33 + painVAS + globalVAS) / 3, nil
}

// PASActivity is a PAS / PAS-II disease-activity category; values are Anderson 2012 Table 3's labels.
type PASActivity string

const (
	PASRemission PASActivity = "Remission"   // 0.00–0.25
	PASLow       PASActivity = "Low/minimal" // 0.26–3.70
	PASModerate  PASActivity = "Moderate"    // 3.71 to < 8.0
	PASHigh      PASActivity = "High/severe" // 8.00–10.00
)

// PASActivityCategory returns the PAS / PAS-II category: Remission 0.00–0.25, Low/minimal 0.26–3.70,
// Moderate 3.71 to < 8.0, High/severe 8.00–10.00. The printed bands leave gaps (0.25–0.26, 3.70–3.71); a
// value in a gap goes to the band below (≤ 0.25 Remission, ≤ 3.70 Low/minimal), as England 2019 writes
// "≤0.25" and "3.71 to <8.0".
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside 0–10.
//
// ID:      MED-RHEUM-010
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Wolfe F, Michaud K, Pincus T. J Rheumatol. 2005;32(12):2410-2415. PMID 16331773 (primary; not opened)
// Source:  Anderson J, Caplan L, Yazdany J, et al. Arthritis Care Res (Hoboken). 2012;64(5):640-647. doi:10.1002/acr.21649, PMID 22473918, PMC4028066, Table 3
// Source:  England BR, Tiong BK, Bergman MJ, et al. Arthritis Care Res (Hoboken). 2019;71(12):1540-1555. doi:10.1002/acr.24042, PMID 31709779, PMC6884664, Table 1
// Verified: secondary (Anderson 2012 PMC4028066 Table 3; England 2019 PMC6884664 Table 1)
// Tags:    rheumatology, rheumatoid arthritis, pas, pas-ii, category
// Units:   pas PAS points; result the source's label
// Range:   pas 0–10
// Example: 0.25 → Remission; 0.26 → Low/minimal; 3.71 → Moderate; 8.0 → High/severe
func PASActivityCategory(pas float64) (PASActivity, error) {
	if math.IsNaN(pas) || math.IsInf(pas, 0) {
		return "", fmt.Errorf("%w: pas = %v", ErrNotFinite, pas)
	}
	switch {
	case pas < 0 || pas > 10:
		return "", fmt.Errorf("%w: pas = %v, want 0–10", ErrOutOfRange, pas)
	case pas <= 0.25:
		return PASRemission, nil
	case pas <= 3.70:
		return PASLow, nil
	case pas < 8.0:
		return PASModerate, nil
	}
	return PASHigh, nil
}

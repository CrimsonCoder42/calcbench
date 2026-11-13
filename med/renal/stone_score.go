package renal

import "fmt"

// STONENauseaVomiting is the nausea-and-vomiting item of the STONE score (Moore 2014 Table 3).
type STONENauseaVomiting int

const (
	// STONENoNauseaVomiting is "None": 0 points.
	STONENoNauseaVomiting STONENauseaVomiting = iota
	// STONENauseaAlone is "Nausea alone": 1 point.
	STONENauseaAlone
	// STONEVomitingAlone is "Vomiting alone": 2 points.
	STONEVomitingAlone
)

// STONEScore returns the STONE score, 0–13, from Moore 2014 Table 3: sex (female 0, male 2); timing, the
// duration of pain to presentation (> 24 h 0, 6–24 h 1, < 6 h 3); race (Black 0, non-black 3); nausea and
// vomiting (none 0, nausea alone 1, vomiting alone 2); hematuria on urine dipstick (absent 0, present 3).
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange for a negative duration or an unknown nausea-vomiting value.
//
// ID:      MED-RENAL-013
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Moore CL, Bomann S, Daniels B, et al. Derivation and validation of a clinical prediction rule for uncomplicated ureteral stone—the STONE score: retrospective and prospective observational cohort studies. BMJ. 2014;348:g2191. doi:10.1136/bmj.g2191. PMID 24671981, PMC3966515, Table 3 "STONE score by factors and categories" (Points column)
// Tags:    renal, urology, ureteral stone, stone score, score
// Units:   painDurationHours hours; result points 0–13
// Range:   painDurationHours ≥ 0 (adults with suspected renal colic, the study's population)
// Example: male, pain 3 h, non-black, vomiting alone, hematuria → 2+3+3+2+3 = 13; female, 30 h, Black, none, absent → 0 (from Table 3)
func STONEScore(male bool, painDurationHours float64, nonBlack bool, nv STONENauseaVomiting, hematuria bool) (int, error) {
	if err := checkFinite("painDurationHours", painDurationHours); err != nil {
		return 0, err
	}
	if painDurationHours < 0 {
		return 0, fmt.Errorf("%w: painDurationHours = %v, want ≥ 0", ErrOutOfRange, painDurationHours)
	}
	s := 0
	if male {
		s += 2
	}
	switch {
	case painDurationHours < 6:
		s += 3
	case painDurationHours <= 24:
		s++
	}
	if nonBlack {
		s += 3
	}
	switch nv {
	case STONENoNauseaVomiting:
	case STONENauseaAlone:
		s++
	case STONEVomitingAlone:
		s += 2
	default:
		return 0, fmt.Errorf("%w: nv = %d, want 0–2", ErrOutOfRange, nv)
	}
	if hematuria {
		s += 3
	}
	return s, nil
}

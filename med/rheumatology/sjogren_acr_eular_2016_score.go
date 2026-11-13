package rheumatology

import "errors"

var (
	// ErrNotFinite is returned when a float input is NaN or ±Inf.
	ErrNotFinite = errors.New("rheumatology: value is NaN or infinite")
	// ErrOutOfRange is returned when an input is outside the range its source's table allows.
	ErrOutOfRange = errors.New("rheumatology: input outside the range the source validated")
)

// pts returns w when present is true and 0 otherwise.
func pts(present bool, w int) int {
	if present {
		return w
	}
	return 0
}

// SjogrenACREULAR2016Score returns the 2016 ACR/EULAR classification score for primary Sjögren's syndrome,
// 0–9: the sum of the weights of the items present (Table 3): labial salivary gland with focal lymphocytic
// sialadenitis and focus score ≥ 1 foci/4 mm², 3; anti-SSA/Ro positive, 3; ocular staining score ≥ 5 (or van
// Bijsterveld ≥ 4) on at least one eye, 1; Schirmer ≤ 5 mm/5 min on at least one eye, 1; unstimulated whole
// saliva flow ≤ 0.1 mL/min, 1. The source applies it to a person meeting its inclusion criterion and none of
// its exclusions, and classifies at "Score ≥ 4"; both are the caller's to apply.
//
// ID:      MED-RHEUM-001
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Shiboski CH, Shiboski SC, Seror R, et al. 2016 American College of Rheumatology/European League Against Rheumatism classification criteria for primary Sjögren's syndrome. Arthritis Rheumatol. 2017;69(1):35-45. doi:10.1002/art.39859, PMID 27785888, Table 3
// Tags:    rheumatology, sjogren, classification criteria, acr, eular, score
// Units:   each item present or absent; result points 0–9
// Range:   any combination of the five items
// Example: biopsy focus score ≥ 1 and Schirmer ≤ 5 → 4; anti-SSA alone → 3; all five → 9 (Table 3 weights)
func SjogrenACREULAR2016Score(lsgFocusScoreAtLeast1, antiSSAPositive, ocularStainingAtLeast5, schirmerAtMost5mm, unstimulatedSalivaAtMost01mLMin bool) int {
	return pts(lsgFocusScoreAtLeast1, 3) + pts(antiSSAPositive, 3) + pts(ocularStainingAtLeast5, 1) +
		pts(schirmerAtMost5mm, 1) + pts(unstimulatedSalivaAtMost01mLMin, 1)
}

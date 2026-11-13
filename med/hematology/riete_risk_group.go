package hematology

import (
	"fmt"
	"math"
)

// RIETEGroup is a RIETE bleeding risk group; its value is the source's words, verbatim.
type RIETEGroup string

// Ruíz-Giménez 2008 abstract: "scored 0 points (low risk); ... 1-4 points (intermediate); ... >4 points (high risk)".
const (
	RIETELow          RIETEGroup = "low risk"     // 0
	RIETEIntermediate RIETEGroup = "intermediate" // 1–4
	RIETEHigh         RIETEGroup = "high risk"    // > 4
)

// RIETERiskGroup returns the RIETE risk group for a score: 0 "low risk"; 1–4 "intermediate"; > 4 "high risk".
// The score moves in steps of 0.5, so the only score between 0 and 1 is 0.5; the abstract does not place it,
// and it cannot occur (the smallest non-zero item is 1 point).
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside 0–8 or for a score that is not a multiple of 0.5
// or equals 0.5.
//
// ID:      MED-HEMAT-010
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Ruíz-Giménez N, Suárez C, González R, et al. Predictive variables for major bleeding events in patients presenting with documented acute venous thromboembolism. Findings from the RIETE Registry. Thromb Haemost. 2008;100(1):26-31. doi:10.1160/TH08-03-0193, PMID 18612534, abstract
// Tags:    hematology, bleeding, anticoagulation, vte, riete, category
// Units:   score points; result the source's words
// Range:   score in {0, 1, 1.5, …, 8}
// Example: 0 → "low risk"; 1 → "intermediate"; 4 → "intermediate"; 4.5 → "high risk" (abstract)
func RIETERiskGroup(score float64) (RIETEGroup, error) {
	if math.IsNaN(score) || math.IsInf(score, 0) {
		return "", fmt.Errorf("%w: score = %v", ErrNotFinite, score)
	}
	if score < 0 || score > 8 || score == 0.5 || math.Mod(score*2, 1) != 0 {
		return "", fmt.Errorf("%w: score = %v, want 0 or 1–8 in steps of 0.5", ErrOutOfRange, score)
	}
	switch {
	case score == 0:
		return RIETELow, nil
	case score <= 4:
		return RIETEIntermediate, nil
	}
	return RIETEHigh, nil
}

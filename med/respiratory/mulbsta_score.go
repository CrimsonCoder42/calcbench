package respiratory

import (
	"fmt"
	"math"
)

// MuLBSTASmoking is the smoking item of the MuLBSTA score.
type MuLBSTASmoking int

// The smoking states the source scores: "acute-smoker (3 points), quit-smoker (2 points)"; any other history 0.
const (
	MuLBSTANonSmoker   MuLBSTASmoking = iota // 0 points
	MuLBSTAAcuteSmoker                       // "acute-smoker", 3 points
	MuLBSTAQuitSmoker                        // "quit-smoker", 2 points
)

// MuLBSTAFindings are the MuLBSTA items, worded as in the source's abstract.
type MuLBSTAFindings struct {
	MultilobularInfiltrates bool           // "multilobular infiltrates (5 points)"
	LymphocytesE9PerL       float64        // lymphocyte count ×10⁹/L; "lymphocyte ≤ 0.8∗109/L (4 points)"
	BacterialCoinfection    bool           // "bacterial coinfection (4 points)"
	Smoking                 MuLBSTASmoking // "acute-smoker (3 points), quit-smoker (2 points)"
	Hypertension            bool           // "hypertension (2 points)"
	AgeYears                float64        // "age ≥60 years (2 points)"
}

// MuLBSTAScore returns the MuLBSTA score (0–20) for viral pneumonia: multilobular infiltrates 5, lymphocytes
// ≤ 0.8 ×10⁹/L 4, bacterial coinfection 4, acute smoker 3 or quit smoker 2, hypertension 2, age ≥ 60 years 2.
// The source's abstract reports "12 points was used as a cut-off value for mortality risk stratification".
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range or for an unknown smoking state.
//
// ID:      MED-RESP-020
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Guo L, Wei D, Zhang X, et al. Clinical Features Predicting Mortality Risk in Patients With Viral Pneumonia: The MuLBSTA Score. Front Microbiol. 2019;10:2752. doi:10.3389/fmicb.2019.02752. PMID 31849894, abstract (items, points and the 12-point cut-off)
// Tags:    respiratory, pneumonia, viral, mulbsta, mortality, score
// Units:   LymphocytesE9PerL ×10⁹ cells/L; AgeYears years; result points 0–20
// Range:   LymphocytesE9PerL ≥ 0; AgeYears ≥ 0 (no limits stated in the abstract)
// Example: multilobular, lymphocytes 0.6, no coinfection, quit smoker, hypertension, age 65 → 5 + 4 + 2 + 2 + 2 = 15 (hand-derived from the abstract's points)
func MuLBSTAScore(f MuLBSTAFindings) (int, error) {
	if err := checkClosedRange("LymphocytesE9PerL", f.LymphocytesE9PerL, 0, math.Inf(1)); err != nil {
		return 0, err
	}
	if err := checkClosedRange("AgeYears", f.AgeYears, 0, math.Inf(1)); err != nil {
		return 0, err
	}
	n := 0
	switch f.Smoking {
	case MuLBSTANonSmoker:
	case MuLBSTAAcuteSmoker:
		n += 3
	case MuLBSTAQuitSmoker:
		n += 2
	default:
		return 0, fmt.Errorf("%w: Smoking = %d, want MuLBSTANonSmoker, MuLBSTAAcuteSmoker or MuLBSTAQuitSmoker", ErrOutOfRange, f.Smoking)
	}
	if f.MultilobularInfiltrates {
		n += 5
	}
	if f.LymphocytesE9PerL <= 0.8 {
		n += 4
	}
	if f.BacterialCoinfection {
		n += 4
	}
	if f.Hypertension {
		n += 2
	}
	if f.AgeYears >= 60 {
		n += 2
	}
	return n, nil
}

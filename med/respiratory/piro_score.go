package respiratory

import (
	"fmt"
	"math"
)

// PIROFindings are the PIRO items for community-acquired pneumonia, worded as in the source's abstract. Each
// present item scores one point; AgeYears scores when "age >70 years".
type PIROFindings struct {
	Comorbidity              bool    // "Comorbidities (chronic obstructive pulmonary disease, immunocompromise)"
	AgeYears                 float64 // "age >70 years"
	MultilobarOpacities      bool    // "multilobar opacities in chest radiograph"
	Shock                    bool    // "shock"
	SevereHypoxemia          bool    // "severe hypoxemia"
	AcuteRenalFailure        bool    // "acute renal failure"
	Bacteremia               bool    // "bacteremia"
	AcuteRespiratoryDistress bool    // "acute respiratory distress syndrome"
}

// PIROScore returns the PIRO score for community-acquired pneumonia (0–8): "one point was given for each present
// feature (range, 0-8 points)".
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange for a negative age.
//
// ID:      MED-RESP-023
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Rello J, Rodriguez A, Lisboa T, Gallego M, Lujan M, Wunderink R. PIRO score for community-acquired pneumonia: a new prediction rule for assessment of severity in intensive care unit patients with community-acquired pneumonia. Crit Care Med. 2009;37(2):456-462. doi:10.1097/CCM.0b013e318194b021. PMID 19114916, abstract (the eight variables, one point each)
// Tags:    respiratory, pneumonia, cap, piro, icu, severity, score
// Units:   AgeYears years; the other items bool; result points 0–8
// Range:   AgeYears ≥ 0 (no limits stated in the abstract)
// Example: age 75, shock, bacteremia → 3 (hand-derived: one point per feature)
func PIROScore(f PIROFindings) (int, error) {
	if err := checkClosedRange("AgeYears", f.AgeYears, 0, math.Inf(1)); err != nil {
		return 0, err
	}
	n := 0
	for _, present := range [8]bool{f.Comorbidity, f.AgeYears > 70, f.MultilobarOpacities, f.Shock,
		f.SevereHypoxemia, f.AcuteRenalFailure, f.Bacteremia, f.AcuteRespiratoryDistress} {
		if present {
			n++
		}
	}
	return n, nil
}

// PIRORisk is the PIRO risk level for a score.
type PIRORisk string

// The levels, verbatim from the abstract: "a) Low, 0-2 points; b) Mild, 3 points; c) high, 4 points; and
// d) Very high, 5-8 points".
const (
	PIROLow      PIRORisk = "Low"
	PIROMild     PIRORisk = "Mild"
	PIROHigh     PIRORisk = "high"
	PIROVeryHigh PIRORisk = "Very high"
)

// PIRORiskLevel returns the PIRO risk level for a score: 0–2 Low, 3 Mild, 4 high, 5–8 Very high.
// Errors: ErrOutOfRange for a score outside 0–8.
//
// ID:      MED-RESP-024
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Rello J, Rodriguez A, Lisboa T, Gallego M, Lujan M, Wunderink R. PIRO score for community-acquired pneumonia: a new prediction rule for assessment of severity in intensive care unit patients with community-acquired pneumonia. Crit Care Med. 2009;37(2):456-462. doi:10.1097/CCM.0b013e318194b021. PMID 19114916, abstract (the four risk levels)
// Tags:    respiratory, pneumonia, cap, piro, category
// Units:   score points 0–8; result the source's level
// Range:   score 0–8
// Example: 2 → "Low"; 3 → "Mild"; 4 → "high"; 5 → "Very high" (the abstract's levels)
func PIRORiskLevel(score int) (PIRORisk, error) {
	switch {
	case score < 0 || score > 8:
		return "", fmt.Errorf("%w: score = %d, want 0 ≤ score ≤ 8", ErrOutOfRange, score)
	case score <= 2:
		return PIROLow, nil
	case score == 3:
		return PIROMild, nil
	case score == 4:
		return PIROHigh, nil
	default:
		return PIROVeryHigh, nil
	}
}

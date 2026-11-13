package neuro

import (
	"fmt"
	"math"
)

// DRAGONScore returns the DRAGON score (0–10) for outcome after IV thrombolysis: (hyper)Dense cerebral artery
// sign / early infarct signs on admission CT (both 2, either 1, none 0); prestroke modified Rankin Scale > 1
// (1); Age ≥ 80 years 2, 65–79 years 1; Glucose > 8 mmol/L (> 144 mg/dL) 1; Onset-to-treatment > 90 min 1;
// baseline NIHSS > 15 3, 10–15 2, 5–9 1, 0–4 0.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range.
//
// ID:      MED-NEURO-005
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Strbian D, Meretoja A, Ahlhelm FJ, et al. Predicting outcome of IV thrombolysis-treated ischemic stroke patients: the DRAGON score. Neurology. 2012;78(6):427-432. doi:10.1212/WNL.0b013e318245d2a9, PMID 22311929, abstract (components and points, quoted in full)
// Tags:    neuro, stroke, thrombolysis, dragon, outcome, score
// Units:   ageYears years; glucoseMmolL mmol/L (the source prints > 8 mmol/L and > 144 mg/dL; mg/dL ÷ 18 before calling); onsetToTreatmentMinutes minutes; nihss NIHSS points; result points 0–10
// Range:   ageYears ≥ 0; glucoseMmolL ≥ 0; onsetToTreatmentMinutes ≥ 0; nihss ≥ 0
// Example: dense sign and early infarct signs, mRS > 1, age 82, glucose 9, OTT 100 min, NIHSS 18 → 2+1+2+1+1+3 = 10 (hand-derived from the abstract)
func DRAGONScore(denseArterySign, earlyInfarctSigns, prestrokeMRSAbove1 bool, ageYears, glucoseMmolL, onsetToTreatmentMinutes float64, nihss int) (int, error) {
	if err := checkRange("ageYears", ageYears, 0, math.MaxFloat64); err != nil {
		return 0, err
	}
	if err := checkRange("glucoseMmolL", glucoseMmolL, 0, math.MaxFloat64); err != nil {
		return 0, err
	}
	if err := checkRange("onsetToTreatmentMinutes", onsetToTreatmentMinutes, 0, math.MaxFloat64); err != nil {
		return 0, err
	}
	if nihss < 0 {
		return 0, fmt.Errorf("%w: nihss = %d, want ≥ 0", ErrOutOfRange, nihss)
	}
	n := 0
	if denseArterySign {
		n++
	}
	if earlyInfarctSigns {
		n++
	}
	if prestrokeMRSAbove1 {
		n++
	}
	switch {
	case ageYears >= 80:
		n += 2
	case ageYears >= 65:
		n++
	}
	if glucoseMmolL > 8 {
		n++
	}
	if onsetToTreatmentMinutes > 90 {
		n++
	}
	switch {
	case nihss > 15:
		n += 3
	case nihss >= 10:
		n += 2
	case nihss >= 5:
		n++
	}
	return n, nil
}

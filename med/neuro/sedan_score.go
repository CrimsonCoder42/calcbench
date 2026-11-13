package neuro

import (
	"fmt"
	"math"
)

// SEDANScore returns the SEDAN score (0–6) for symptomatic intracranial hemorrhage after thrombolysis:
// baseline blood Sugar 8.1–12.0 mmol/L 1, > 12.0 mmol/L 2; Early infarct signs 1; (hyper)Dense cerebral
// artery sign 1; Age > 75 years 1; NIH Stroke Scale on admission ≥ 10 1. The source prints the glucose bands
// to 0.1 mmol/L; a value above 8.0 and below 8.1 is counted in the 8.1–12.0 band (no gap is left).
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range.
//
// ID:      MED-NEURO-007
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Strbian D, Engelter S, Michel P, et al. Symptomatic intracranial hemorrhage after stroke thrombolysis: the SEDAN score. Ann Neurol. 2012;71(5):634-641. doi:10.1002/ana.23546, PMID 22522478, abstract ("baseline blood Sugar (glucose; 8.1-12.0 mmol/l [145-216 mg/dl] = 1; >12.0 mmol/l [>216 mg/dl] = 2), Early infarct signs (yes = 1) and (hyper)Dense cerebral artery sign (yes = 1) on admission computed tomography scan, Age (>75 years = 1), and NIH Stroke Scale on admission (≥10 = 1)")
// Tags:    neuro, stroke, thrombolysis, intracranial hemorrhage, sedan, score
// Units:   glucoseMmolL mmol/L (mg/dL ÷ 18 before calling); ageYears years; nihss NIHSS points; result points 0–6
// Range:   glucoseMmolL ≥ 0; ageYears ≥ 0; nihss ≥ 0
// Example: glucose 13, early infarct signs, no dense sign, age 80, NIHSS 12 → 2 + 1 + 0 + 1 + 1 = 5 (hand-derived from the abstract)
func SEDANScore(glucoseMmolL float64, earlyInfarctSigns, denseArterySign bool, ageYears float64, nihss int) (int, error) {
	if err := checkRange("glucoseMmolL", glucoseMmolL, 0, math.MaxFloat64); err != nil {
		return 0, err
	}
	if err := checkRange("ageYears", ageYears, 0, math.MaxFloat64); err != nil {
		return 0, err
	}
	if nihss < 0 {
		return 0, fmt.Errorf("%w: nihss = %d, want ≥ 0", ErrOutOfRange, nihss)
	}
	n := 0
	switch {
	case glucoseMmolL > 12.0:
		n += 2
	case glucoseMmolL > 8.0:
		n++
	}
	if earlyInfarctSigns {
		n++
	}
	if denseArterySign {
		n++
	}
	if ageYears > 75 {
		n++
	}
	if nihss >= 10 {
		n++
	}
	return n, nil
}

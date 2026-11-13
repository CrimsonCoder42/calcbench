package neuro

import "math"

// ICHScore returns the ICH Score (0–6): Glasgow Coma Scale 3–4 2 points, 5–12 1, 13–15 0; age ≥ 80 years 1;
// infratentorial origin 1; ICH volume ≥ 30 cm³ 1; intraventricular hemorrhage 1.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range.
//
// ID:      MED-NEURO-006
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Hemphill JC 3rd, Bonovich DC, Besmertis L, Manley GT, Johnston SC. The ICH score: a simple, reliable grading scale for intracerebral hemorrhage. Stroke. 2001;32(4):891-897. doi:10.1161/01.str.32.4.891, PMID 11283388, abstract ("GCS score 3 to 4 (=2 points), 5 to 12 (=1), 13 to 15 (=0); age >/=80 years yes (=1) ...; infratentorial origin yes (=1) ...; ICH volume >/=30 cm(3) (=1) ...; and intraventricular hemorrhage yes (=1)")
// Tags:    neuro, stroke, intracerebral hemorrhage, ich, gcs, score
// Units:   gcsTotal GCS points 3–15; ageYears years; volumeML mL (= cm³); result points 0–6
// Range:   gcsTotal 3–15; ageYears ≥ 0; volumeML ≥ 0
// Example: GCS 10, age 82, supratentorial, 35 mL, IVH → 1 + 1 + 0 + 1 + 1 = 4 (hand-derived from the abstract)
func ICHScore(gcsTotal int, ageYears float64, infratentorial bool, volumeML float64, intraventricular bool) (int, error) {
	if err := checkInt("gcsTotal", gcsTotal, 3, 15); err != nil {
		return 0, err
	}
	if err := checkRange("ageYears", ageYears, 0, math.MaxFloat64); err != nil {
		return 0, err
	}
	if err := checkRange("volumeML", volumeML, 0, math.MaxFloat64); err != nil {
		return 0, err
	}
	n := 0
	switch {
	case gcsTotal <= 4:
		n += 2
	case gcsTotal <= 12:
		n++
	}
	if ageYears >= 80 {
		n++
	}
	if infratentorial {
		n++
	}
	if volumeML >= 30 {
		n++
	}
	if intraventricular {
		n++
	}
	return n, nil
}

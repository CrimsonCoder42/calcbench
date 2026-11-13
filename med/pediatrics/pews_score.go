package pediatrics

import "fmt"

// PEWSFindings holds the Brighton Paediatric Early Warning Score items (Monaghan 2005, as reproduced by Shafi
// 2020 and Malde 2024). Each row is scored 0–3 by the caller from the chart's cell descriptions.
type PEWSFindings struct {
	Behaviour               int  // 0 "Playing/appropriate"; 1 "Sleeping"; 2 "Irritable"; 3 "Lethargic/confused or reduced response to pain"
	Cardiovascular          int  // 0 "Pink or capillary refill 1-2 s" … 3 "Gray and mottled or capillary refill ≥5 s or tachycardia ≥30 bpm above normal rate or bradycardia"
	Respiratory             int  // 0 "Within normal parameters, no retractions or tracheal tug" … 3 (see Shafi 2020)
	QuarterHourlyNebulisers bool // "¼ hourly nebulizers": 2
	PostoperativeVomiting   bool // "persistent vomiting following surgery": 2
}

// PEWSScore returns the Brighton Paediatric Early Warning Score, 0–13: behaviour, cardiovascular and
// respiratory rows 0–3 each, and "Score 2 extra for persistent vomiting following surgery or ¼ hourly
// nebulizers", each scored separately (the 0–13 total of Miranda 2017 requires both). Action bands differ
// between the sources and are not built.
// Errors: ErrOutOfRange for a row outside 0–3.
//
// ID:      MED-PEDS-006
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Monaghan A. Detecting and managing deterioration in children. Paediatr Nurs. 2005;17(1):32-35. doi:10.7748/paed2005.02.17.1.32.c964, PMID 15751446 (primary; not opened)
// Source:  Shafi OM, Diego Rondon JD, Gulati G. Can the Pediatric Early Warning Score (PEWS) Predict Hospital Length of Stay? Cureus. 2020;12(11):e11339. doi:10.7759/cureus.11339, PMID 33304675, PMC7719480, PEWS table
// Source:  Malde S, Jain P, Revathi N, et al. Evaluation of Pediatric Early Warning Score (PEWS) in Unintentional Childhood Injuries Admitted to the Critical Care Unit. Cureus. 2024;16(7):e65312. doi:10.7759/cureus.65312, PMID 39184619, PMC11344085, Methods (components and the extra 2 points)
// Verified: secondary (Shafi 2020 PMC7719480; Malde 2024 PMC11344085; total 0–13: Miranda 2017 Rev Lat Am Enfermagem 25:e2912, PMC5511004)
// Tags:    pediatrics, early warning, pews, brighton, score
// Units:   rows points 0–3; result points 0–13
// Range:   hospitalised children (the source population); each row 0–3
// Example: irritable (2), capillary refill 3 s (1), respiratory row 0, quarter-hourly nebulisers (2) → 5 (hand-derived)
func PEWSScore(f PEWSFindings) (int, error) {
	for _, in := range [...]struct {
		name string
		v    int
	}{{"Behaviour", f.Behaviour}, {"Cardiovascular", f.Cardiovascular}, {"Respiratory", f.Respiratory}} {
		if in.v < 0 || in.v > 3 {
			return 0, fmt.Errorf("%w: %s = %d, want 0–3", ErrOutOfRange, in.name, in.v)
		}
	}
	s := f.Behaviour + f.Cardiovascular + f.Respiratory
	if f.QuarterHourlyNebulisers {
		s += 2
	}
	if f.PostoperativeVomiting {
		s += 2
	}
	return s, nil
}

package neuro

// FASTEDScore returns the Field Assessment Stroke Triage for Emergency Destination score (0–9), the sum of its
// five item points: facial palsy 0–1; arm weakness 0–2; speech changes 0–2; eye deviation 0–2; denial /
// neglect 0–2. The source reports "A FAST-ED ≥4"; applying it is the caller's.
// Errors: ErrOutOfRange when an item is outside its range, naming the item.
//
// ID:      MED-NEURO-009
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Lima FO, Silva GS, Furie KL, et al. Field Assessment Stroke Triage for Emergency Destination: a simple and accurate prehospital scale to detect large vessel occlusion strokes. Stroke. 2016;47(8):1997-2002. doi:10.1161/STROKEAHA.116.013301, PMID 27364531, PMC4961538, Table 1 (item levels and points)
// Tags:    neuro, stroke, prehospital, large vessel occlusion, fast-ed, score
// Units:   each item points as in Table 1; result points 0–9
// Range:   facialPalsy 0–1; armWeakness, speechChanges, eyeDeviation, denialNeglect 0–2
// Example: facial 1, arm 2, speech 1, eye 1, neglect 0 → 5 (hand-derived from Table 1)
func FASTEDScore(facialPalsy, armWeakness, speechChanges, eyeDeviation, denialNeglect int) (int, error) {
	if err := checkInt("facialPalsy", facialPalsy, 0, 1); err != nil {
		return 0, err
	}
	if err := checkInt("armWeakness", armWeakness, 0, 2); err != nil {
		return 0, err
	}
	if err := checkInt("speechChanges", speechChanges, 0, 2); err != nil {
		return 0, err
	}
	if err := checkInt("eyeDeviation", eyeDeviation, 0, 2); err != nil {
		return 0, err
	}
	if err := checkInt("denialNeglect", denialNeglect, 0, 2); err != nil {
		return 0, err
	}
	return facialPalsy + armWeakness + speechChanges + eyeDeviation + denialNeglect, nil
}

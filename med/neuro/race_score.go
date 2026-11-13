package neuro

// RACEScore returns the Rapid Arterial oCclusion Evaluation scale total (0–9), the sum of its five item
// scores: facial palsy 0–2; arm motor function 0–2; leg motor function 0–2; gaze 0–1; aphasia or agnosia 0–2.
// The source reports "RACE scale≥5"; applying it is the caller's.
// Errors: ErrOutOfRange when an item is outside its range, naming the item.
//
// ID:      MED-NEURO-010
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Pérez de la Ossa N, Carrera D, Gorchs M, et al. Design and validation of a prehospital stroke scale to predict large arterial occlusion: the Rapid Arterial Occlusion Evaluation scale. Stroke. 2014;45(1):87-91. doi:10.1161/STROKEAHA.113.003071, PMID 24281224, abstract ("facial palsy (scored 0-2), arm motor function (0-2), leg motor function (0-2), gaze (0-1), and aphasia or agnosia (0-2)")
// Tags:    neuro, stroke, prehospital, large vessel occlusion, race, score
// Units:   each item points; result points 0–9
// Range:   facialPalsy, armMotor, legMotor, aphasiaOrAgnosia 0–2; gaze 0–1
// Example: facial 1, arm 2, leg 1, gaze 1, aphasia 0 → 5 (hand-derived from the abstract)
func RACEScore(facialPalsy, armMotor, legMotor, gaze, aphasiaOrAgnosia int) (int, error) {
	if err := checkInt("facialPalsy", facialPalsy, 0, 2); err != nil {
		return 0, err
	}
	if err := checkInt("armMotor", armMotor, 0, 2); err != nil {
		return 0, err
	}
	if err := checkInt("legMotor", legMotor, 0, 2); err != nil {
		return 0, err
	}
	if err := checkInt("gaze", gaze, 0, 1); err != nil {
		return 0, err
	}
	if err := checkInt("aphasiaOrAgnosia", aphasiaOrAgnosia, 0, 2); err != nil {
		return 0, err
	}
	return facialPalsy + armMotor + legMotor + gaze + aphasiaOrAgnosia, nil
}

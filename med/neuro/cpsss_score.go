package neuro

import "fmt"

// CPSSSScore returns the Cincinnati Prehospital Stroke Severity Scale (0–4) from four NIH Stroke Scale item
// scores: 2 points for conjugate gaze (NIHSS gaze ≥ 1); 1 point for arm weakness (NIHSS arm ≥ 2); 1 point for
// abnormal level of consciousness questions and commands (NIHSS LOC questions ≥ 1 and LOC commands ≥ 1).
// Errors: ErrOutOfRange when an item score is negative.
//
// ID:      MED-NEURO-004
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Katz BS, McMullan JT, Sucharew H, Adeoye O, Broderick JP. Design and validation of a prehospital scale to predict stroke severity: Cincinnati Prehospital Stroke Severity Scale. Stroke. 2015;46(6):1508-1512. doi:10.1161/STROKEAHA.115.008804, PMID 25899242, abstract ("2 points for presence of conjugate gaze (NIHSS≥1); 1 point for presence of arm weakness (NIHSS≥2); and 1 point for presence abnormal level of consciousness commands and questions (NIHSS level of consciousness≥1 each)")
// Tags:    neuro, stroke, prehospital, nihss, cp-sss, score
// Units:   each input the NIHSS item score (points), worst arm for nihssArm; result points 0–4
// Range:   every item ≥ 0 (the NIHSS form's upper limits are not checked here)
// Example: gaze 1, arm 3, LOC questions 1, LOC commands 0 → 2 + 1 + 0 = 3 (hand-derived from the abstract)
func CPSSSScore(nihssGaze, nihssArm, nihssLOCQuestions, nihssLOCCommands int) (int, error) {
	for _, p := range [4]struct {
		name string
		v    int
	}{{"nihssGaze", nihssGaze}, {"nihssArm", nihssArm}, {"nihssLOCQuestions", nihssLOCQuestions}, {"nihssLOCCommands", nihssLOCCommands}} {
		if p.v < 0 {
			return 0, fmt.Errorf("%w: %s = %d, want ≥ 0", ErrOutOfRange, p.name, p.v)
		}
	}
	n := 0
	if nihssGaze >= 1 {
		n += 2
	}
	if nihssArm >= 2 {
		n++
	}
	if nihssLOCQuestions >= 1 && nihssLOCCommands >= 1 {
		n++
	}
	return n, nil
}

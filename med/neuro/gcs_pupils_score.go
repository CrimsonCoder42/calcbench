package neuro

// GCSPupilsScore returns the GCS-Pupils score (1–15): the Glasgow Coma Scale total (3–15) minus the number of
// non-reacting pupils (0, 1 or 2). The GCS total itself is MED-SCALE-003 GlasgowComaScore.
// Errors: ErrOutOfRange outside Range.
//
// ID:      MED-NEURO-013
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Brennan PM, Murray GD, Teasdale GM. Simplifying the use of prognostic information in traumatic brain injury. Part 1: The GCS-Pupils score: an extended index of clinical severity. J Neurosurg. 2018;128(6):1612-1620. doi:10.3171/2017.12.JNS172780, PMID 29631516, abstract ("GCS score [range 3-15] minus the number of nonreacting pupils [0, 1, or 2]), which we call the GCS-Pupils score (GCS-P; range 1-15)")
// Tags:    neuro, head injury, tbi, gcs, pupils, gcs-p, score
// Units:   gcsTotal points 3–15; nonReactingPupils count 0–2; result points 1–15
// Range:   gcsTotal 3–15; nonReactingPupils 0–2
// Example: GCS 3, both pupils non-reacting → 1; GCS 15, both reacting → 15 (abstract's range ends)
func GCSPupilsScore(gcsTotal, nonReactingPupils int) (int, error) {
	if err := checkInt("gcsTotal", gcsTotal, 3, 15); err != nil {
		return 0, err
	}
	if err := checkInt("nonReactingPupils", nonReactingPupils, 0, 2); err != nil {
		return 0, err
	}
	return gcsTotal - nonReactingPupils, nil
}

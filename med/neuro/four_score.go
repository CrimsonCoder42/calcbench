package neuro

// FOURScore returns the Full Outline of UnResponsiveness (FOUR) score, 0–16: the sum of its four components,
// eye response, motor response, brainstem reflexes and respiration, each rated 0–4.
// Errors: ErrOutOfRange when a component is outside 0–4, naming it.
//
// ID:      MED-NEURO-026
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Wijdicks EF, Bamlet WR, Maramattom BV, Manno EM, McClelland RL. Validation of a new coma scale: the FOUR score. Ann Neurol. 2005;58(4):585-593. doi:10.1002/ana.20611, PMID 16178024 (abstract: four components, "maximal score of 4" each)
// Source:  Fischer M, Rüegg S, Czaplinski A, et al. Inter-rater reliability of the Full Outline of UnResponsiveness score and the Glasgow Coma Scale in critically ill patients: a prospective observational study. Crit Care. 2010;14(2):R64. doi:10.1186/cc8963, PMC2887186, Table 1 (eye, motor, brainstem, respiration rated 4 to 0; "FOUR score 16")
// Source:  Peng J, Deng Y, Chen F, et al. Validation of the Chinese version of the FOUR score in the assessment of neurosurgical patients with different level of consciousness. BMC Neurol. 2015;15:254. doi:10.1186/s12883-015-0508-9, PMC4676150, Background ("All the four categories are scored from 0 to 4 points")
// Verified: secondary (Fischer 2010 PMC2887186; Peng 2015 PMC4676150)
// Tags:    neuro, coma, consciousness, four score, score
// Units:   component points 0–4; result points 0–16
// Range:   each component 0–4
// Example: E4 M4 B4 R4 → 16; E0 M0 B0 R0 → 0
func FOURScore(eye, motor, brainstem, respiration int) (int, error) {
	if err := checkInt("eye", eye, 0, 4); err != nil {
		return 0, err
	}
	if err := checkInt("motor", motor, 0, 4); err != nil {
		return 0, err
	}
	if err := checkInt("brainstem", brainstem, 0, 4); err != nil {
		return 0, err
	}
	if err := checkInt("respiration", respiration, 0, 4); err != nil {
		return 0, err
	}
	return eye + motor + brainstem + respiration, nil
}

package neuro

// LAMSScore returns the Los Angeles Motor Scale, 0–5: facial droop (0 absent, 1 present) plus arm drift
// (0 absent, 1 drifts down, 2 falls rapidly) plus grip strength (0 normal, 1 weak grip, 2 no grip).
// Errors: ErrOutOfRange when an item is outside its range, naming it.
//
// ID:      MED-NEURO-025
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Llanes JN, Kidwell CS, Starkman S, Leary MC, Eckstein M, Saver JL. The Los Angeles Motor Scale (LAMS): a new measure to characterize stroke severity in the field. Prehosp Emerg Care. 2004;8(1):46-50. doi:10.1080/312703002806, PMID 14691787 (abstract: facial weakness, arm strength and grip, total 0–5; item points not in the abstract)
// Source:  Arthur KC, et al. J Clin Med. 2023;12(10):3374. doi:10.3390/jcm12103374, PMC10219214 ("0 or 1 for facial strength (normal or droop); 0, 1, or 2 for arm strength (normal, drifts down, falls rapidly), and 0, 1, or 2 for grip strength (normal, weak, no grip)")
// Source:  Thu PW, Yu YP, Wang JH, Liew CL. Tzu Chi Med J. 2025;37(3):311-320. doi:10.4103/tcmj.tcmj_191_24, PMC12306891, Table 1 ("Facial droop (1)"; "Drift down (1)"; "Rapid fall (2)"; "Weak grip (1)"; "No grip (2)")
// Verified: secondary (Arthur 2023 PMC10219214; Thu 2025 PMC12306891)
// Tags:    neuro, stroke, prehospital, lams, score
// Units:   item points; result points 0–5
// Range:   facialDroop 0–1; armDrift 0–2; gripStrength 0–2
// Example: droop present, arm falls rapidly, no grip → 5; every item 0 → 0
func LAMSScore(facialDroop, armDrift, gripStrength int) (int, error) {
	if err := checkInt("facialDroop", facialDroop, 0, 1); err != nil {
		return 0, err
	}
	if err := checkInt("armDrift", armDrift, 0, 2); err != nil {
		return 0, err
	}
	if err := checkInt("gripStrength", gripStrength, 0, 2); err != nil {
		return 0, err
	}
	return facialDroop + armDrift + gripStrength, nil
}

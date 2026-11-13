package respiratory

// PARSFindings are the six factors of the Pediatric Asthma Risk Score (Biagini Myers 2019, Table III), worded
// as in the table.
type PARSFindings struct {
	ParentalAsthma          bool // "Parental asthma" 2
	EczemaBeforeAge3        bool // "Eczema before age 3 years" 2
	WheezingApartFromColds  bool // "Wheezing apart from colds" 3
	EarlyWheezingBeforeAge3 bool // "Early wheezing (before age 3 years)" 3
	SPTPositive2OrMore      bool // "SPT+ to ≥2 aero and/or food allergens" 2
	AfricanAmericanRace     bool // "African American Race" 2
}

// PARSScore returns the Pediatric Asthma Risk Score (0–14): parental asthma 2, eczema before age 3 years 2,
// wheezing apart from colds 3, early wheezing (before age 3 years) 3, skin prick test positive to ≥ 2 aero
// and/or food allergens 2, African American race 2. "The scores range from 0–14 with scores of 1 and 13
// unattainable given the weighting of the ORs."
//
// ID:      MED-RESP-042
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Biagini Myers JM, Schauberger E, He H, et al. A Pediatric Asthma Risk Score to better predict asthma development in young children. J Allergy Clin Immunol. 2019;143(5):1803-1810.e2. doi:10.1016/j.jaci.2018.09.037. PMID 30554722, PMC6504569, Table III and Methods (weights, range 0–14)
// Tags:    respiratory, asthma, pediatric, pars, risk, score
// Units:   every item bool; result points 0–14
// Range:   any combination of items (children followed from birth to age 7 in the derivation cohort)
// Example: parental asthma and early wheezing → 2 + 3 = 5 (hand-derived from Table III)
func PARSScore(f PARSFindings) int {
	s := 0
	for _, t := range [6]struct {
		on  bool
		pts int
	}{{f.ParentalAsthma, 2}, {f.EczemaBeforeAge3, 2}, {f.WheezingApartFromColds, 3}, {f.EarlyWheezingBeforeAge3, 3},
		{f.SPTPositive2OrMore, 2}, {f.AfricanAmericanRace, 2}} {
		if t.on {
			s += t.pts
		}
	}
	return s
}

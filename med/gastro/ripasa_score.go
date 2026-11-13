package gastro

// RIPASAFindings holds the 15 RIPASA parameters of Chong 2010. The three two-level parameters (age, sex,
// duration) are booleans for the level that scores 1 point; the other level scores 0.5.
type RIPASAFindings struct {
	AgeUnder40                bool // "less than 40 years is 1 point; greater than 40 years is 0.5 point"
	Male                      bool // "male is 1 point; female is 0.5 point"
	RightIliacFossaPain       bool // 0.5 point
	MigrationOfPainToRIF      bool // 0.5 point
	NauseaAndVomiting         bool // 1 point
	Anorexia                  bool // 1 point
	SymptomsUnder48Hours      bool // "less than 48 hours is 1 point; more than 48 hours is 0.5 point"
	RightIliacFossaTenderness bool // 1 point
	Guarding                  bool // 2 points
	ReboundTenderness         bool // 1 point
	RovsingSign               bool // 2 points
	Fever                     bool // 1 point
	RaisedWhiteCellCount      bool // 1 point
	NegativeUrinalysis        bool // 1 point
	ForeignNationalIDCard     bool // "foreign national registration identity card (1 point)"
}

// RIPASAScore returns the RIPASA score for acute appendicitis, 1.5–16 in steps of 0.5, from the 15
// parameters of Chong 2010. Age exactly 40 and duration exactly 48 h are not assigned by the source (its
// levels are "less than" and "greater than"); the caller chooses the level. No band is returned; the source's
// "optimal cut-off threshold score from the ROC was 7.5".
// Errors: none; every combination of findings is valid.
//
// ID:      MED-GASTRO-005
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Chong CF, Adi MI, Thien A, et al. Development of the RIPASA score: a new appendicitis scoring system for the diagnosis of acute appendicitis. Singapore Med J. 2010;51(3):220-225. PMID 20428744. Abstract: the 15 parameters and their scores
// Tags:    gastro, appendicitis, ripasa, score
// Units:   result points 1.5–16 (half-point steps, exact in float64)
// Range:   no numeric input (population as described in the source; not restated here)
// Example: male, age under 40, under 48 h, RIF pain, RIF tenderness, guarding, raised WBC → 1 + 1 + 1 + 0.5 + 1 + 2 + 1 = 7.5 (HAND-DERIVED from the abstract)
func RIPASAScore(f RIPASAFindings) float64 {
	score := 0.0
	two := func(on bool) {
		if on {
			score += 1
		} else {
			score += 0.5
		}
	}
	add := func(on bool, p float64) {
		if on {
			score += p
		}
	}
	two(f.AgeUnder40)
	two(f.Male)
	two(f.SymptomsUnder48Hours)
	add(f.RightIliacFossaPain, 0.5)
	add(f.MigrationOfPainToRIF, 0.5)
	add(f.NauseaAndVomiting, 1)
	add(f.Anorexia, 1)
	add(f.RightIliacFossaTenderness, 1)
	add(f.Guarding, 2)
	add(f.ReboundTenderness, 1)
	add(f.RovsingSign, 2)
	add(f.Fever, 1)
	add(f.RaisedWhiteCellCount, 1)
	add(f.NegativeUrinalysis, 1)
	add(f.ForeignNationalIDCard, 1)
	return score
}

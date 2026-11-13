package neuro

// MoCADomains are the seven domain sub-totals of the Montreal Cognitive Assessment with their maxima.
type MoCADomains struct {
	VisuospatialExecutive int // 0–5
	Naming                int // 0–3
	Attention             int // 0–6
	Language              int // 0–3
	Abstraction           int // 0–2
	DelayedRecall         int // 0–5
	Orientation           int // 0–6
}

// MoCATotal returns the Montreal Cognitive Assessment total, 0–30: the sum of the seven domain sub-totals plus
// one point when the person has 12 years of education or fewer, capped at the instrument's maximum of 30. The
// source reports a cut-off of 26; applying it is the caller's. Items not embedded; licence: the MoCA is
// © Z. Nasreddine (MoCA Test Inc.), and its administration and form are the licence holder's.
// Errors: ErrOutOfRange when a domain is outside its range, naming it.
//
// ID:      MED-NEURO-035
// Time:    O(1)   7 domains, fixed
// Space:   O(1), no allocation
// Source:  Nasreddine ZS, Phillips NA, Bédirian V, et al. The Montreal Cognitive Assessment, MoCA: a brief screening tool for mild cognitive impairment. J Am Geriatr Soc. 2005;53(4):695-699. doi:10.1111/j.1532-5415.2005.53221.x, PMID 15817019 (abstract: cut-off 26; domain points not in the abstract)
// Source:  Kwon CY. BMJ Open. 2026;16(6):e114046. doi:10.1136/bmjopen-2025-114046, PMC13331112 ("visuospatial/executive function (5 points), naming (3 points), attention (6 points), language (3 points), abstraction (2 points), delayed recall (5 points) and orientation (6 points)"; "One additional point is added for individuals with ≤12 years of education")
// Source:  Ruiz-Marrara J, Antunes LG, Bataglion C, et al. J Oral Rehabil. 2026;53(1):76-88. doi:10.1111/joor.70064, PMC12705290 ("total score, which can range from 0 to 30"; "an additional point was given to individuals with 12 or fewer years of education")
// Verified: secondary (Kwon 2026 PMC13331112; Ruiz-Marrara 2026 PMC12705290)
// Tags:    neuro, cognition, moca, screening, score
// Units:   domain points; twelveYearsOrLessEducation yes/no; result points 0–30
// Range:   each domain within the range on MoCADomains
// Example: domains summing to 25 with ≤12 years of education → 26; all maxima with ≤12 years → 30 (capped)
func MoCATotal(d MoCADomains, twelveYearsOrLessEducation bool) (int, error) {
	items := [...]struct {
		name string
		v    int
		max  int
	}{
		{"VisuospatialExecutive", d.VisuospatialExecutive, 5}, {"Naming", d.Naming, 3},
		{"Attention", d.Attention, 6}, {"Language", d.Language, 3}, {"Abstraction", d.Abstraction, 2},
		{"DelayedRecall", d.DelayedRecall, 5}, {"Orientation", d.Orientation, 6},
	}
	total := 0
	for _, x := range items {
		if err := checkInt(x.name, x.v, 0, x.max); err != nil {
			return 0, err
		}
		total += x.v
	}
	if twelveYearsOrLessEducation && total < 30 {
		total++
	}
	return total, nil
}

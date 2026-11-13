package neuro

// FourATCategory is a 4AT score band; its value is the source's words, verbatim.
type FourATCategory string

// Bellelli 2014: "A score of 0 is intended to suggest that delirium and/or moderate to severe cognitive
// impairment is unlikely"; "Scores between 1 and 3 are intended to suggest possible moderate to severe general
// cognitive impairment"; "A score of 4 or above suggests possible delirium".
const (
	FourATUnlikely                    FourATCategory = "delirium and/or moderate to severe cognitive impairment is unlikely" // 0
	FourATPossibleCognitiveImpairment FourATCategory = "possible moderate to severe general cognitive impairment"            // 1–3
	FourATPossibleDelirium            FourATCategory = "possible delirium"                                                   // ≥ 4
)

// FourATCategoryOf returns the 4AT band for a total (0–12): 0, 1–3, ≥ 4. The four items' points were not read in
// the primary (they are in its supplement), so the total is the caller's; only the published bands are here.
// Errors: ErrOutOfRange when total is outside 0–12.
//
// ID:      MED-NEURO-015
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Bellelli G, Morandi A, Davis DH, et al. Validation of the 4AT, a new instrument for rapid delirium screening: a study in 234 hospitalised older people. Age Ageing. 2014;43(4):496-502. doi:10.1093/ageing/afu021, PMID 24590568, PMC4066613, Methods (score "0 to 12"; the three bands, quoted)
// Tags:    neuro, delirium, cognition, 4at, category
// Units:   total points 0–12; result the source's words
// Range:   total 0–12
// Example: 0 → unlikely; 3 → possible cognitive impairment; 4 → "possible delirium" (source)
func FourATCategoryOf(total int) (FourATCategory, error) {
	if err := checkInt("total", total, 0, 12); err != nil {
		return "", err
	}
	switch {
	case total == 0:
		return FourATUnlikely, nil
	case total <= 3:
		return FourATPossibleCognitiveImpairment, nil
	}
	return FourATPossibleDelirium, nil
}

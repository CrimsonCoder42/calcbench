package neuro

// TokuhashiGroup is a revised Tokuhashi total-score group; its value is the source's words, verbatim.
type TokuhashiGroup string

// Tokuhashi 2005 abstract: "total score of 8 or less (predicted survival period, less than 6 months)";
// "9 to 11 (predicted survival period, 6 months or more)"; "12 or more (predicted survival period, 1 year or more)".
const (
	TokuhashiUnder6Months  TokuhashiGroup = "predicted survival period, less than 6 months" // 0–8
	Tokuhashi6MonthsOrMore TokuhashiGroup = "predicted survival period, 6 months or more"   // 9–11
	Tokuhashi1YearOrMore   TokuhashiGroup = "predicted survival period, 1 year or more"     // 12–15
)

// TokuhashiRevisedGroup returns the revised Tokuhashi group for a total score (0–15): ≤ 8, 9–11, ≥ 12. The
// six parameters' points that make the total were not read in the primary, so the total is the caller's;
// only the published groups are here.
// Errors: ErrOutOfRange when total is outside 0–15.
//
// ID:      MED-NEURO-014
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Tokuhashi Y, Matsuzaki H, Oda H, Oshima M, Ryu J. A revised scoring system for preoperative evaluation of metastatic spine tumor prognosis. Spine. 2005;30(19):2186-2191. doi:10.1097/01.brs.0000180401.06919.a5, PMID 16205345, abstract ("the total score was 15 points"; the three groups, quoted)
// Tags:    neuro, oncology, spine metastasis, tokuhashi, category
// Units:   total points 0–15; result the source's words
// Range:   total 0–15
// Example: 8 → "predicted survival period, less than 6 months"; 9 → "... 6 months or more"; 12 → "... 1 year or more" (abstract)
func TokuhashiRevisedGroup(total int) (TokuhashiGroup, error) {
	if err := checkInt("total", total, 0, 15); err != nil {
		return "", err
	}
	switch {
	case total <= 8:
		return TokuhashiUnder6Months, nil
	case total <= 11:
		return Tokuhashi6MonthsOrMore, nil
	}
	return Tokuhashi1YearOrMore, nil
}

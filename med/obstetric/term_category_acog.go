package obstetric

import (
	"errors"
	"fmt"
)

var (
	// ErrOutOfRange is returned when an input is outside the range the source defines or validated.
	ErrOutOfRange = errors.New("obstetric: input outside the range the source validated")
)

// TermCategory is a gestational-age category of ACOG Committee Opinion No. 579, with the source's label.
type TermCategory string

// The four categories of ACOG Committee Opinion No. 579, labels verbatim.
const (
	TermEarly    TermCategory = "early term" // 37 0/7 through 38 6/7 weeks
	TermFull     TermCategory = "full term"  // 39 0/7 through 40 6/7 weeks
	TermLate     TermCategory = "late term"  // 41 0/7 through 41 6/7 weeks
	TermPostterm TermCategory = "postterm"   // 42 0/7 weeks and beyond
)

// TermCategoryACOG returns the ACOG Committee Opinion No. 579 category for a gestational age given as
// completed weeks and days (weeks w/7): early term 37 0/7–38 6/7, full term 39 0/7–40 6/7, late term
// 41 0/7–41 6/7, postterm 42 0/7 and beyond. A gestational age before 37 0/7 weeks is not one of these
// categories and returns ErrOutOfRange. How the gestational age was estimated is the caller's.
// Errors: ErrOutOfRange outside Range.
//
// ID:      MED-OB-001
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  ACOG Committee Opinion No 579: Definition of term pregnancy. Obstet Gynecol. 2013;122(5):1139-1140. doi:10.1097/01.AOG.0000437385.88715.4a. PMID 24150030, abstract: "early term (37 0/7 weeks of gestation through 38 6/7 weeks of gestation), full term (39 0/7 weeks of gestation through 40 6/7 weeks of gestation), late term (41 0/7 weeks of gestation through 41 6/7 weeks of gestation), and postterm (42 0/7 weeks of gestation and beyond)"
// Tags:    obstetric, pregnancy, gestational-age, term, acog
// Units:   gestationalWeeks completed weeks; gestationalDays days past the completed weeks (0–6); result category
// Range:   gestationalWeeks ≥ 37 (37 0/7 weeks and beyond); gestationalDays 0–6
// Example: 38 6/7 → "early term"; 39 0/7 → "full term"; 41 6/7 → "late term"; 42 0/7 → "postterm" (PMID 24150030)
func TermCategoryACOG(gestationalWeeks, gestationalDays int) (TermCategory, error) {
	if gestationalDays < 0 || gestationalDays > 6 {
		return "", fmt.Errorf("%w: gestationalDays = %d, want 0–6", ErrOutOfRange, gestationalDays)
	}
	switch {
	case gestationalWeeks < 37:
		return "", fmt.Errorf("%w: gestationalWeeks = %d, want ≥ 37 (37 0/7 weeks and beyond)", ErrOutOfRange, gestationalWeeks)
	case gestationalWeeks <= 38:
		return TermEarly, nil
	case gestationalWeeks <= 40:
		return TermFull, nil
	case gestationalWeeks == 41:
		return TermLate, nil
	default:
		return TermPostterm, nil
	}
}

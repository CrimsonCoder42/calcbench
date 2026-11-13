package scales

import "fmt"

// WASTScore returns the Woman Abuse Screening Tool total, 8–24: each of the 8 items is recoded "to reflect a
// higher score for higher reported frequency of experiences" (3 = "A lot of tension" / "Great difficulty" /
// "Often"; 2 = the middle option; 1 = "No tension" / "No difficulty" / "Never") and summed "for individuals
// who answered all 8 items". The source document gives no cut-off for the 8-item form; none is built. Items
// are not embedded.
// Errors: ErrOutOfRange for an item outside 1–3.
//
// ID:      MED-SCALE-013
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Brown JB, Lent B, Brett PJ, Sas G, Pederson LL. Development of the Woman Abuse Screening Tool for use in family practice. Fam Med. 1996;28(6):422-428. PMID 8791071 (primary; not opened)
// Source:  Basile KC, Hertz MF, Back SE. Intimate Partner Violence and Sexual Violence Victimization Assessment Instruments for Use in Healthcare Settings, Version 1. Atlanta (GA): Centers for Disease Control and Prevention, National Center for Injury Prevention and Control; 2007. https://stacks.cdc.gov/view/cdc/44660, WAST pp. 72–74 (scoring) (read 2026-11-12)
// Source:  Aydemir A, Aksoy İ, Şahin E, et al. Psychometric validation of the Turkish version of the Woman Abuse Screening Tool (WAST) in emergency department settings. BMC Womens Health. 2026;26(1):485. doi:10.1186/s12905-026-04663-0, PMID 42415006, PMC13629004, Methods ("1 (often) to 3 (never)", reverse-scored)
// Verified: secondary (CDC Basile 2007 pp. 72–74; Aydemir 2026 PMC13629004; range 8–24 also Al Ubaidi 2021 PMC8483099)
// Tags:    scales, intimate partner violence, screening, wast, score
// Units:   items points 1–3 each (higher = more frequent / more tension); result points 8–24
// Range:   each item 1–3
// Example: items 2,2,2,1,1,1,2,1 → 12 (hand-derived)
func WASTScore(items [8]int) (int, error) {
	s := 0
	for i, v := range items {
		if v < 1 || v > 3 {
			return 0, fmt.Errorf("%w: items[%d] = %d, want 1–3", ErrOutOfRange, i, v)
		}
		s += v
	}
	return s, nil
}

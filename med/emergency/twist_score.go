package emergency

// TWISTScore returns the TWIST (Testicular Workup for Ischemia and Suspected Torsion) score, 0–7: "testicular
// swelling (2 points), hard testicle (2), absent cremasteric reflex (1), nausea/vomiting (1) and high riding
// testis (1)". The abstract reports "Cutoffs for low and high risk were 2 and 5 points"; the categories are
// not built here.
//
// ID:      MED-EMERG-001
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Barbosa JA, Tiseo BC, Barayan GA, et al. Development and initial validation of a scoring system to diagnose testicular torsion in children. J Urol. 2013;189(5):1859-1864. doi:10.1016/j.juro.2012.10.056, PMID 23103800, abstract
// Tags:    emergency, urology, pediatric, testicular torsion, twist, score
// Units:   each finding present or absent; result points 0–7
// Range:   any combination of the five findings
// Example: swelling and hard testicle → 4; all five → 7; nausea alone → 1 (abstract)
func TWISTScore(testicularSwelling, hardTesticle, absentCremastericReflex, nauseaOrVomiting, highRidingTestis bool) int {
	s := 0
	if testicularSwelling {
		s += 2
	}
	if hardTesticle {
		s += 2
	}
	if absentCremastericReflex {
		s++
	}
	if nauseaOrVomiting {
		s++
	}
	if highRidingTestis {
		s++
	}
	return s
}

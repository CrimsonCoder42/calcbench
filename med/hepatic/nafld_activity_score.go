package hepatic

import "fmt"

// NAFLDActivityScore returns the NAFLD activity score (NAS), 0–8: the unweighted sum of the histological
// grades for steatosis (0–3: < 5 %, 5–33 %, 34–66 %, > 66 % of hepatocytes), lobular inflammation (0–3: no
// foci, < 2, 2–4, > 4 foci per 200× field) and hepatocyte ballooning (0–2: none, few, many cells). Fibrosis is
// staged separately and is not part of the NAS. No band is returned.
// Errors: ErrOutOfRange when a grade is outside its range.
//
// ID:      MED-HEPAT-008
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Kleiner DE, Brunt EM, Van Natta M, et al. Design and validation of a histological scoring system for nonalcoholic fatty liver disease. Hepatology. 2005;41(6):1313-1321. doi:10.1002/hep.20701. PMID 15915461 (full text not opened; the PubMed abstract prints "lobular inflammation (0-2)", which disagrees with both secondaries)
// Source:  Dowman JK, Tomlinson JW, Newsome PN. Systematic review: the diagnosis and staging of non-alcoholic fatty liver disease and non-alcoholic steatohepatitis. Aliment Pharmacol Ther. 2011. PMID 21198708. PMC3080668: "steatosis (0–3), lobular inflammation (0–3) and hepatocyte ballooning (0–2)" [Kleiner NAS]
// Source:  Benedict M, Zhang X. Non-alcoholic fatty liver disease: an expanded review. World J Hepatol. 2017. PMID 28652891. PMC5468341, Table 1 (grade definitions for each component)
// Verified: secondary (Dowman 2011 PMC3080668; Benedict 2017 PMC5468341)
// Tags:    hepatic, liver, nafld, nash, histology, nas, steatosis, ballooning, score
// Units:   steatosisGrade 0–3; lobularInflammationGrade 0–3; ballooningGrade 0–2; result points 0–8
// Range:   steatosisGrade 0–3; lobularInflammationGrade 0–3; ballooningGrade 0–2 (liver biopsies graded by a pathologist)
// Example: steatosis 2, lobular inflammation 2, ballooning 1 → 5 (hand-derived)
func NAFLDActivityScore(steatosisGrade, lobularInflammationGrade, ballooningGrade int) (int, error) {
	if steatosisGrade < 0 || steatosisGrade > 3 {
		return 0, fmt.Errorf("%w: steatosisGrade = %d, want 0–3", ErrOutOfRange, steatosisGrade)
	}
	if lobularInflammationGrade < 0 || lobularInflammationGrade > 3 {
		return 0, fmt.Errorf("%w: lobularInflammationGrade = %d, want 0–3", ErrOutOfRange, lobularInflammationGrade)
	}
	if ballooningGrade < 0 || ballooningGrade > 2 {
		return 0, fmt.Errorf("%w: ballooningGrade = %d, want 0–2", ErrOutOfRange, ballooningGrade)
	}
	return steatosisGrade + lobularInflammationGrade + ballooningGrade, nil
}

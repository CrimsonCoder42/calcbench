package hematology

// ICCSystemicMastocytosis reports whether the ICC 2022 criteria for systemic mastocytosis are met: the major
// criterion ("Multifocal dense infiltrates of tryptase- and/or CD117 positive mast cells (≥15 mast cells in
// aggregates)" in bone marrow and/or other extracutaneous organs), or, "In the absence of the major criterion,
// at least 3 of the following 4 minor criteria": > 25% of mast cells spindle shaped or atypical immature;
// mast cells express CD25, CD2 and/or CD30; KIT D816V or another activating KIT mutation; serum tryptase
// "persistently >20 ng/mL". "In cases of SM-AMN an elevated tryptase does not count as a SM minor criterion":
// with smAMN true the tryptase input is not counted.
// Errors: none; every combination of the criteria is defined.
//
// ID:      MED-HEMAT-013
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Arber DA, Orazi A, Hasserjian RP, et al. International Consensus Classification of Myeloid Neoplasms and Acute Leukemias: integrating morphologic, clinical, and genomic data. Blood. 2022;140(11):1200-1228. doi:10.1182/blood.2022015850, PMID 35767897, PMC9479031, Table 11 (systemic mastocytosis)
// Tags:    hematology, mastocytosis, systemic mastocytosis, icc, diagnostic criteria
// Units:   criteria present/absent; result true when the rule is met
// Range:   any combination of the criteria
// Example: major only → true; three minor criteria without the major → true; morphology, markers and tryptase with SM-AMN → false (Table 11)
func ICCSystemicMastocytosis(major, minorMorphology, minorMarkers, minorKIT, minorTryptaseOver20NgML, smAMN bool) bool {
	if major {
		return true
	}
	n := pts(minorMorphology, 1) + pts(minorMarkers, 1) + pts(minorKIT, 1)
	if !smAMN {
		n += pts(minorTryptaseOver20NgML, 1)
	}
	return n >= 3
}

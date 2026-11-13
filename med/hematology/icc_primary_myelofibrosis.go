package hematology

import "fmt"

// ICCPMFMarrow is the bone-marrow major criterion of the ICC primary myelofibrosis criteria.
type ICCPMFMarrow int

// Bone-marrow findings of Arber 2022, Table 5 (major criterion 1); comments abridge the source's text.
const (
	ICCPMFMarrowNeither     ICCPMFMarrow = iota // neither biopsy finding below
	ICCPMFMarrowPrefibrotic                     // "megakaryocytic proliferation and atypia, bone marrow fibrosis grade < 2, increased age-adjusted BM cellularity, granulocytic proliferation, and (often) decreased erythropoiesis"
	ICCPMFMarrowOvert                           // "megakaryocytic proliferation and atypia, accompanied by reticulin and/or collagen fibrosis grades 2 or 3"
)

// ICCPMFCategory is the ICC primary myelofibrosis stage whose criteria are met; values are the source's
// abbreviations, verbatim.
type ICCPMFCategory string

// Stages of Arber 2022, Table 5.
const (
	ICCPMFNotMet      ICCPMFCategory = ""          // the rule is not met
	ICCPMFPrefibrotic ICCPMFCategory = "pre-PMF"   // early/prefibrotic primary myelofibrosis
	ICCPMFOvert       ICCPMFCategory = "overt PMF" // overt primary myelofibrosis
)

// ICCPrimaryMyelofibrosis returns the ICC 2022 primary myelofibrosis stage whose criteria are met: "The
// diagnosis of pre-PMF or overt PMF requires all 3 major criteria and at least 1 minor criterion confirmed in 2
// consecutive determinations". Major criteria: the marrow finding of the stage; "JAK2, CALR, or MPL mutation or
// another clonal marker or absence of reactive" fibrosis (clonalMarker); the criteria for the other myeloid
// neoplasms are not met (otherNeoplasmsExcluded). Minor criteria: anemia not attributed to a comorbid
// condition; leukocytosis ≥ 11 × 10⁹/L; palpable splenomegaly; LDH above the reference range; and, for overt
// PMF only, leukoerythroblastosis. The two consecutive determinations are the caller's to confirm.
// Errors: ErrOutOfRange for a marrow value other than the constants.
//
// ID:      MED-HEMAT-012
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Arber DA, Orazi A, Hasserjian RP, et al. International Consensus Classification of Myeloid Neoplasms and Acute Leukemias: integrating morphologic, clinical, and genomic data. Blood. 2022;140(11):1200-1228. doi:10.1182/blood.2022015850, PMID 35767897, PMC9479031, Table 5 (primary myelofibrosis)
// Tags:    hematology, myelofibrosis, pmf, icc, diagnostic criteria
// Units:   marrow finding as the constants; criteria present/absent; result the source's stage abbreviation, "" when not met
// Range:   marrow one of the ICCPMFMarrow constants
// Example: prefibrotic marrow, JAK2, exclusion met, anemia → "pre-PMF"; overt marrow, CALR, exclusion met, leukoerythroblastosis only → "overt PMF"; prefibrotic marrow with leukoerythroblastosis as the only minor → "" (Table 5)
func ICCPrimaryMyelofibrosis(marrow ICCPMFMarrow, clonalMarker, otherNeoplasmsExcluded, anemia, leukocytosisAtLeast11E9PerL, palpableSplenomegaly, ldhAboveReference, leukoerythroblastosis bool) (ICCPMFCategory, error) {
	if marrow < ICCPMFMarrowNeither || marrow > ICCPMFMarrowOvert {
		return ICCPMFNotMet, fmt.Errorf("%w: marrow = %d, want 0–2", ErrOutOfRange, marrow)
	}
	if marrow == ICCPMFMarrowNeither || !clonalMarker || !otherNeoplasmsExcluded {
		return ICCPMFNotMet, nil
	}
	minor := anemia || leukocytosisAtLeast11E9PerL || palpableSplenomegaly || ldhAboveReference
	if marrow == ICCPMFMarrowOvert {
		if minor || leukoerythroblastosis {
			return ICCPMFOvert, nil
		}
		return ICCPMFNotMet, nil
	}
	if minor {
		return ICCPMFPrefibrotic, nil
	}
	return ICCPMFNotMet, nil
}

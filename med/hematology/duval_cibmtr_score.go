package hematology

import "fmt"

// DuvalDiseaseGroup is the disease group row of the Duval 2010 AML score.
type DuvalDiseaseGroup int

// Disease group rows of Duval 2010, Table 5 (AML); each comment is the row text, verbatim.
const (
	DuvalPIFOrCR1Over6Months DuvalDiseaseGroup = iota // "PIF or duration of first CR > 6 months", 0 points
	DuvalCR1Under6Months                              // "Duration of first CR < 6 months", 1 point
)

// DuvalDonorGroup is the HLA match group row of the Duval 2010 AML score.
type DuvalDonorGroup int

// HLA match group rows of Duval 2010, Table 5 (AML); each comment is the row text, verbatim.
const (
	DuvalDonorSiblingOrMatchedUnrelated  DuvalDonorGroup = iota // "HLA identical sibling or well matched or partially matched unrelated", 0 points
	DuvalDonorMismatchedUnrelated                               // "Mismatched unrelated", 1 point
	DuvalDonorRelatedNotIdenticalSibling                        // "Related other than HLA identical sibling", 2 points
)

// DuvalCIBMTRScore returns the Duval (CIBMTR) pre-transplant score for acute myeloid leukaemia in relapse or
// primary induction failure, 0–6: duration of first CR < 6 months 1; poor-risk cytogenetics before HSCT 1;
// mismatched unrelated donor 1, related donor other than an HLA-identical sibling 2; circulating blasts present
// 1; Karnofsky or Lansky score < 90 1.
// Errors: ErrOutOfRange for a disease or donor group other than the constants.
//
// ID:      MED-HEMAT-011
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Duval M, Klein JP, He W, et al. Hematopoietic stem-cell transplantation for acute leukemia in relapse or primary induction failure. J Clin Oncol. 2010;28(23):3730-3738. doi:10.1200/JCO.2010.28.8852, PMID 20625136, PMC2917308, Table 5 (AML), and the abstract's "predictive score of 0 had 42% OS at 3 years, whereas OS was 6% for a score ≥ 3"
// Tags:    hematology, aml, leukemia, hsct, transplant, cibmtr, duval, prognostic score
// Units:   disease and donor group as the table's rows; findings present/absent; result points 0–6
// Range:   disease one of the DuvalDiseaseGroup constants; donor one of the DuvalDonorGroup constants (AML, myeloablative HSCT in relapse or primary induction failure, the study's population)
// Example: PIF, intermediate cytogenetics, HLA-identical sibling, blasts present, Karnofsky 80 → 2; CR1 < 6 months, poor, related non-identical, blasts, Karnofsky 70 → 6 (HAND-DERIVED from Table 5)
func DuvalCIBMTRScore(disease DuvalDiseaseGroup, poorCytogenetics bool, donor DuvalDonorGroup, circulatingBlasts, performanceBelow90 bool) (int, error) {
	if disease < DuvalPIFOrCR1Over6Months || disease > DuvalCR1Under6Months {
		return 0, fmt.Errorf("%w: disease = %d, want 0–1", ErrOutOfRange, disease)
	}
	if donor < DuvalDonorSiblingOrMatchedUnrelated || donor > DuvalDonorRelatedNotIdenticalSibling {
		return 0, fmt.Errorf("%w: donor = %d, want 0–2", ErrOutOfRange, donor)
	}
	return int(disease) + pts(poorCytogenetics, 1) + int(donor) + pts(circulatingBlasts, 1) + pts(performanceBelow90, 1), nil
}

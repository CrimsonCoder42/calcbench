package rheumatology

import "fmt"

// SLICC2012Classified reports whether the SLICC 2012 classification rule for systemic lupus erythematosus is met:
// "four of the criteria listed in Table 3, including at least one clinical criterion and one immunologic
// criterion", or "biopsy-proven nephritis compatible with SLE and with ANA or anti-dsDNA antibodies".
// clinicalCount is how many of the 11 clinical criteria are met; immunologicCount how many of the 6
// immunologic criteria (the criteria's definitions are the source's; the caller counts them).
// Errors: ErrOutOfRange when clinicalCount is outside 0–11 or immunologicCount outside 0–6.
//
// ID:      MED-RHEUM-005
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Petri M, Orbai AM, Alarcón GS, et al. Derivation and validation of the Systemic Lupus International Collaborating Clinics classification criteria for systemic lupus erythematosus. Arthritis Rheum. 2012;64(8):2677-2686. doi:10.1002/art.34473, PMID 22553077, Table 3 and its rule
// Tags:    rheumatology, lupus, sle, slicc, classification criteria
// Units:   clinicalCount 0–11 and immunologicCount 0–6 criteria met; result true when the rule is met
// Range:   0 ≤ clinicalCount ≤ 11; 0 ≤ immunologicCount ≤ 6
// Example: 3 clinical + 1 immunologic → true; 4 clinical + 0 immunologic → false; 0 + 0 with biopsy nephritis and ANA → true (the rule)
func SLICC2012Classified(clinicalCount, immunologicCount int, biopsyNephritisWithANAOrDsDNA bool) (bool, error) {
	if clinicalCount < 0 || clinicalCount > 11 {
		return false, fmt.Errorf("%w: clinicalCount = %d, want 0–11", ErrOutOfRange, clinicalCount)
	}
	if immunologicCount < 0 || immunologicCount > 6 {
		return false, fmt.Errorf("%w: immunologicCount = %d, want 0–6", ErrOutOfRange, immunologicCount)
	}
	if biopsyNephritisWithANAOrDsDNA {
		return true, nil
	}
	return clinicalCount >= 1 && immunologicCount >= 1 && clinicalCount+immunologicCount >= 4, nil
}

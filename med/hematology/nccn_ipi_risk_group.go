package hematology

import "fmt"

// NCCNIPIGroup is an NCCN-IPI risk group; its value is the source's label, verbatim.
type NCCNIPIGroup string

// Zhou 2014 Table 3 risk groups.
const (
	NCCNIPILow              NCCNIPIGroup = "Low"               // 0–1
	NCCNIPILowIntermediate  NCCNIPIGroup = "Low-intermediate"  // 2–3
	NCCNIPIHighIntermediate NCCNIPIGroup = "High-intermediate" // 4–5
	NCCNIPIHigh             NCCNIPIGroup = "High"              // ≥ 6
)

// NCCNIPIRiskGroup returns the NCCN-IPI risk group for a score: 0–1 "Low"; 2–3 "Low-intermediate"; 4–5
// "High-intermediate"; ≥ 6 "High".
// Errors: ErrOutOfRange when score is outside 0–8.
//
// ID:      MED-HEMAT-002
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Zhou Z, Sehn LH, Rademaker AW, et al. An enhanced International Prognostic Index (NCCN-IPI) for patients with diffuse large B-cell lymphoma treated in the rituximab era. Blood. 2014;123(6):837-842. doi:10.1182/blood-2013-09-524108, PMID 24264230, Table 3
// Tags:    hematology, lymphoma, dlbcl, nccn-ipi, category
// Units:   score points 0–8; result the source's label
// Range:   0 ≤ score ≤ 8
// Example: 1 → "Low"; 2 → "Low-intermediate"; 5 → "High-intermediate"; 6 → "High" (Table 3)
func NCCNIPIRiskGroup(score int) (NCCNIPIGroup, error) {
	if score < 0 || score > 8 {
		return "", fmt.Errorf("%w: score = %d, want 0–8", ErrOutOfRange, score)
	}
	switch {
	case score <= 1:
		return NCCNIPILow, nil
	case score <= 3:
		return NCCNIPILowIntermediate, nil
	case score <= 5:
		return NCCNIPIHighIntermediate, nil
	}
	return NCCNIPIHigh, nil
}

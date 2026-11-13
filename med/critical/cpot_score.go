package critical

import "fmt"

// CPOTScore returns the Critical-Care Pain Observation Tool score, 0–8: the sum of four indicators each rated
// 0–2 — facial expression (0 relaxed or neutral, 1 tense, 2 grimacing); body movements (0 absence of movement
// or normal position, 1 protection, 2 restlessness or agitation); muscle tension (0 relaxed, 1 tense or rigid,
// 2 very tense or rigid); and compliance with the ventilator for intubated adults (0 fully compliant, 1 coughing but
// tolerating, 2 fighting ventilator) or vocalization for extubated adults. Item wording is the source's; only the sum is
// computed here.
// Errors: ErrOutOfRange when an indicator is outside 0–2.
//
// ID:      MED-CRIT-012
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Gélinas C, Fillion L, Puntillo KA, Viens C, Fortier M. Validation of the critical-care pain observation tool in adult patients. Am J Crit Care. 2006;15(4):420-427. PMID 16823021
// Source:  Afenigus AD. Evaluating pain in non-verbal critical care patients: a narrative review of the critical care pain observation tool and its clinical applications. Front Pain Res (Lausanne). 2024;5:1481085. doi:10.3389/fpain.2024.1481085, PMID 39474515, PMC11518793, Table 1
// Source:  Asriyanto LF, Chayati N. Cross-cultural adaptation and validation of the Indonesian version of the Critical-care Pain Observation Tool. Int J Nurs Sci. 2024;11(1):113-119. PMID 38352280, PMC10859591, the original CPOT's four indicators 0–2, total 0–8
// Verified: secondary (Afenigus 2024, PMC11518793; Asriyanto & Chayati 2024, PMC10859591)
// Tags:    critical, pain, intensive care, cpot, score
// Units:   each indicator points 0–2; result points 0–8
// Range:   each indicator 0–2 (adults in intensive care unable to self-report, the validation population)
// Example: grimacing 2, protection 1, tense 1, coughing but tolerating 1 → 5; all relaxed → 0 (Table 1 of Afenigus 2024)
func CPOTScore(facialExpression, bodyMovements, muscleTension, ventilatorComplianceOrVocalization int) (int, error) {
	for _, e := range [...]error{
		cpotItem("facialExpression", facialExpression), cpotItem("bodyMovements", bodyMovements),
		cpotItem("muscleTension", muscleTension), cpotItem("ventilatorComplianceOrVocalization", ventilatorComplianceOrVocalization),
	} {
		if e != nil {
			return 0, e
		}
	}
	return facialExpression + bodyMovements + muscleTension + ventilatorComplianceOrVocalization, nil
}

// cpotItem returns ErrOutOfRange, wrapped with name, unless 0 ≤ v ≤ 2.
func cpotItem(name string, v int) error {
	if v < 0 || v > 2 {
		return fmt.Errorf("%w: %s = %d, want 0–2", ErrOutOfRange, name, v)
	}
	return nil
}

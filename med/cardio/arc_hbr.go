package cardio

import "fmt"

// ARCHBRHighBleedingRisk reports whether the ARC-HBR consensus definition of high bleeding risk at the time of
// PCI is met: "at least 1 major or 2 minor criteria". The caller counts the criteria of Urban 2019 Table 1
// that are present (major, e.g. anticipated long-term oral anticoagulation, eGFR < 30 mL/min, hemoglobin
// < 11 g/dL, platelet count < 100 × 10⁹/L; minor, e.g. age ≥ 75 years, eGFR 30–59 mL/min, hemoglobin
// 11–12.9 g/dL in men or 11–11.9 g/dL in women).
// Errors: ErrOutOfRange when a count is negative.
//
// ID:      MED-CARDIO-053
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Urban P, Mehran R, Colleran R, et al. Defining high bleeding risk in patients undergoing percutaneous coronary intervention: a consensus document from the Academic Research Consortium for High Bleeding Risk. Circulation. 2019;140(3):240-261. doi:10.1161/CIRCULATIONAHA.119.040167. PMID 31116032. PMC6636810, Table 1 and text: HBR "if at least 1 major or 2 minor criteria are met"
// Tags:    cardio, pci, bleeding, arc-hbr, decision-rule
// Units:   majorCriteria, minorCriteria counts; result true when the HBR definition is met
// Range:   counts ≥ 0 (persons undergoing PCI)
// Example: 0 major, 2 minor → true; 0 major, 1 minor → false
func ARCHBRHighBleedingRisk(majorCriteria, minorCriteria int) (bool, error) {
	if majorCriteria < 0 || minorCriteria < 0 {
		return false, fmt.Errorf("%w: majorCriteria = %d, minorCriteria = %d, want ≥ 0", ErrOutOfRange, majorCriteria, minorCriteria)
	}
	return majorCriteria >= 1 || minorCriteria >= 2, nil
}

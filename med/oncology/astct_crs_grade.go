package oncology

import (
	"errors"
	"fmt"
)

var (
	// ErrOutOfRange reports an input outside the values the source defines.
	ErrOutOfRange = errors.New("oncology: input outside the range the source validated")
)

// CRSHypotension is the hypotension parameter of the ASTCT CRS grading table; each constant's comment is the
// table's cell, verbatim.
type CRSHypotension int

// Hypotension rows of Lee 2019, Table 2.
const (
	CRSHypotensionNone                 CRSHypotension = iota // "None" (grade 1 column)
	CRSHypotensionNoVasopressors                             // "Not requiring vasopressors" (grade 2)
	CRSHypotensionOneVasopressor                             // "Requiring one vasopressor with or without vasopressin" (grade 3)
	CRSHypotensionMultipleVasopressors                       // "Requiring multiple vasopressors (excluding vasopressin)" (grade 4)
)

// CRSHypoxia is the hypoxia parameter of the ASTCT CRS grading table; each constant's comment is the table's
// cell, verbatim.
type CRSHypoxia int

// Hypoxia rows of Lee 2019, Table 2.
const (
	CRSHypoxiaNone             CRSHypoxia = iota // "None" (grade 1 column)
	CRSHypoxiaLowFlow                            // "Requiring low-flow nasal cannula or blow-by" (≤ 6 L/min) (grade 2)
	CRSHypoxiaHighFlowOrMask                     // "Requiring high-flow nasal cannula, facemask, non-rebreather mask, or Venturi mask" (> 6 L/min) (grade 3)
	CRSHypoxiaPositivePressure                   // "Requiring positive pressure (eg: CPAP, BiPAP, intubation and mechanical ventilation)" (grade 4)
)

// ASTCTCRSGrade returns the ASTCT consensus grade of cytokine release syndrome, 1–4: with the fever condition
// met, grade 1 is fever alone and grades 2–4 follow the more severe of the hypotension and hypoxia rows ("CRS
// grade is determined by the more severe event: hypotension or hypoxia"). The fever condition is
// "temperature ≥38°C not attributable to any other cause" (feverAtLeast38C), or, after antipyretics or
// anti-cytokine therapy such as tocilizumab or steroids, fever "is no longer required to grade subsequent CRS
// severity" (afterAntipyreticOrAnticytokine). When neither holds the table assigns no grade and the result is
// 0. Organ toxicities "do not influence CRS grading" and are not an input. Grade 5 (death due to CRS) is not
// computed.
// Errors: ErrOutOfRange for a hypotension or hypoxia value other than the constants.
//
// ID:      MED-ONCO-001
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Lee DW, Santomasso BD, Locke FL, et al. ASTCT Consensus Grading for Cytokine Release Syndrome and Neurologic Toxicity Associated with Immune Effector Cells. Biol Blood Marrow Transplant. 2019;25(4):625-638. doi:10.1016/j.bbmt.2018.12.758, PMID 30592986, PMC12180426, Table 2 (ASTCT CRS Consensus Grading) and its footnotes
// Tags:    oncology, immune effector cell, car-t, cytokine release syndrome, crs, astct, grade
// Units:   findings present/absent; hypotension and hypoxia as the table's rows; result grade 0–4 (0 = fever condition unmet)
// Range:   hypotension one of the CRSHypotension constants; hypoxia one of the CRSHypoxia constants
// Example: fever, no hypotension, low-flow cannula → 2; fever, one vasopressor, low-flow → 3; no fever and no prior antipyretic or anti-cytokine therapy → 0 (Table 2)
func ASTCTCRSGrade(feverAtLeast38C, afterAntipyreticOrAnticytokine bool, hypotension CRSHypotension, hypoxia CRSHypoxia) (int, error) {
	if hypotension < CRSHypotensionNone || hypotension > CRSHypotensionMultipleVasopressors {
		return 0, fmt.Errorf("%w: hypotension = %d, want 0–3", ErrOutOfRange, hypotension)
	}
	if hypoxia < CRSHypoxiaNone || hypoxia > CRSHypoxiaPositivePressure {
		return 0, fmt.Errorf("%w: hypoxia = %d, want 0–3", ErrOutOfRange, hypoxia)
	}
	if !feverAtLeast38C && !afterAntipyreticOrAnticytokine {
		return 0, nil
	}
	return 1 + max(int(hypotension), int(hypoxia)), nil
}

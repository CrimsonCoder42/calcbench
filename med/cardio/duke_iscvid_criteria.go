package cardio

import "fmt"

// DukeISCVIDCategory is a 2023 Duke-ISCVID infective endocarditis category, the source's own label.
type DukeISCVIDCategory string

// The three 2023 Duke-ISCVID categories (Fowler 2023, Table 1 definitions).
const (
	DukeISCVIDDefinite DukeISCVIDCategory = "Definite IE"
	DukeISCVIDPossible DukeISCVIDCategory = "Possible IE"
	DukeISCVIDRejected DukeISCVIDCategory = "Rejected IE"
)

// DukeISCVIDFindings are the counts and flags the 2023 Duke-ISCVID classification is built from.
type DukeISCVIDFindings struct {
	PathologicCriterion bool // a pathologic criterion met (microorganisms or active endocarditis in tissue)
	MajorCriteria       int  // major criteria met, 0–3 (microbiologic, imaging, surgical)
	MinorCriteria       int  // minor criteria met, 0–7 (predisposition, fever, vascular, immunologic, microbiologic, imaging, physical examination)
	// RejectionFinding: "firm alternate diagnosis", "lack of recurrence despite antibiotic therapy for less
	// than 4 d", or no pathologic or macroscopic evidence at surgery or autopsy after < 4 d of antibiotics.
	RejectionFinding bool
}

// DukeISCVIDCategoryOf returns the 2023 Duke-ISCVID category: "Definite IE" for a pathologic criterion, or 2
// major, or 1 major and 3 minor, or 5 minor criteria; otherwise "Rejected IE" when a rejection finding is
// present; otherwise "Possible IE" for 1 major and 1 minor, or 3 minor criteria; otherwise "Rejected IE"
// ("does not meet criteria for possible IE"). Definite is tested first, as the table lists it first.
// Errors: ErrOutOfRange when a count is outside its range.
//
// ID:      MED-CARDIO-054
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Fowler VG Jr, Durack DT, Selton-Suty C, et al. The 2023 Duke-International Society for Cardiovascular Infectious Diseases criteria for infective endocarditis: updating the modified Duke criteria. Clin Infect Dis. 2023;77(4):518-526. doi:10.1093/cid/ciad271. PMID 37138445. PMC10681650, Table 1 (Definite, Possible, Rejected)
// Tags:    cardio, infective-endocarditis, duke, iscvid, diagnostic-criteria
// Units:   counts of criteria; result category label
// Range:   MajorCriteria 0–3; MinorCriteria 0–7
// Example: 1 major, 3 minor → "Definite IE"; 1 major, 1 minor → "Possible IE"
func DukeISCVIDCategoryOf(f DukeISCVIDFindings) (DukeISCVIDCategory, error) {
	if f.MajorCriteria < 0 || f.MajorCriteria > 3 {
		return "", fmt.Errorf("%w: MajorCriteria = %d, want 0–3", ErrOutOfRange, f.MajorCriteria)
	}
	if f.MinorCriteria < 0 || f.MinorCriteria > 7 {
		return "", fmt.Errorf("%w: MinorCriteria = %d, want 0–7", ErrOutOfRange, f.MinorCriteria)
	}
	maj, mnr := f.MajorCriteria, f.MinorCriteria
	switch {
	case f.PathologicCriterion || maj >= 2 || (maj >= 1 && mnr >= 3) || mnr >= 5:
		return DukeISCVIDDefinite, nil
	case f.RejectionFinding:
		return DukeISCVIDRejected, nil
	case (maj >= 1 && mnr >= 1) || mnr >= 3:
		return DukeISCVIDPossible, nil
	}
	return DukeISCVIDRejected, nil
}

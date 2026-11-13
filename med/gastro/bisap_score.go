package gastro

import (
	"errors"
	"fmt"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

var (
	// ErrNotFinite is returned when an input is NaN or ±Inf; it is checked before any range.
	ErrNotFinite = errors.New("gastro: value is NaN or infinite")
	// ErrOutOfRange is returned when an input is outside the range the source validated.
	ErrOutOfRange = errors.New("gastro: input outside the range the source validated")
)

// checkFinite returns ErrNotFinite, wrapped with name, when v is NaN or ±Inf.
func checkFinite(name string, v float64) error {
	if !reduce.AllFinite(v) {
		return fmt.Errorf("%w: %s = %v", ErrNotFinite, name, v)
	}
	return nil
}

// BISAPFindings holds the five BISAP items, each assessed during the first 24 h (Wu 2008).
type BISAPFindings struct {
	BUNMgDL              float64 // blood urea nitrogen, mg/dL; > 25 mg/dL: 1 point
	ImpairedMentalStatus bool    // impaired mental status: 1 point
	SIRS                 bool    // systemic inflammatory response syndrome present: 1 point
	AgeYears             float64 // completed age, years; > 60 years: 1 point
	PleuralEffusion      bool    // pleural effusion present: 1 point
}

// BISAPScore returns the Bedside Index for Severity in Acute Pancreatitis (BISAP), 0–5: one point each for
// BUN > 25 mg/dL, impaired mental status, SIRS, age > 60 years and a pleural effusion, during the first 24 h.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range.
//
// ID:      MED-GASTRO-001
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Wu BU, Johannes RS, Sun X, Tabak Y, Conwell DL, Banks PA. The early prediction of mortality in acute pancreatitis: a large population-based study. Gut. 2008;57(12):1698-1703. doi:10.1136/gut.2008.152702. PMID 18519429. Abstract: "One point is assigned for the presence of each of the following during the first 24 h: blood urea nitrogen (BUN) >25 mg/dl; impaired mental status; systemic inflammatory response syndrome (SIRS); age >60 years; or the presence of a pleural effusion"
// Tags:    gastro, pancreatitis, bisap, severity, score
// Units:   BUNMgDL mg/dL; AgeYears years; result points 0–5
// Range:   BUNMgDL ≥ 0; AgeYears ≥ 0 (derivation: 17,992 cases of acute pancreatitis, 212 hospitals, 2000–2001; input limits not stated in the source)
// Example: BUN 30, age 65, no other item → 2; BUN 25, age 60 → 0 (HAND-DERIVED from the abstract)
func BISAPScore(f BISAPFindings) (int, error) {
	if err := checkFinite("BUNMgDL", f.BUNMgDL); err != nil {
		return 0, err
	}
	if err := checkFinite("AgeYears", f.AgeYears); err != nil {
		return 0, err
	}
	if !(f.BUNMgDL >= 0) {
		return 0, fmt.Errorf("%w: BUNMgDL = %v, want ≥ 0", ErrOutOfRange, f.BUNMgDL)
	}
	if !(f.AgeYears >= 0) {
		return 0, fmt.Errorf("%w: AgeYears = %v, want ≥ 0", ErrOutOfRange, f.AgeYears)
	}
	points := 0
	if f.BUNMgDL > 25 {
		points++
	}
	if f.ImpairedMentalStatus {
		points++
	}
	if f.SIRS {
		points++
	}
	if f.AgeYears > 60 {
		points++
	}
	if f.PleuralEffusion {
		points++
	}
	return points, nil
}

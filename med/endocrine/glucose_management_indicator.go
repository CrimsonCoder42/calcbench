package endocrine

import (
	"errors"
	"fmt"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

var (
	// ErrNotFinite is returned when an input is NaN or ±Inf; it is checked before any range.
	ErrNotFinite = errors.New("endocrine: value is NaN or infinite")
	// ErrOutOfRange is returned when an input is outside the range the source validated.
	ErrOutOfRange = errors.New("endocrine: input outside the range the source validated")
	// ErrDomain is returned when inputs valid one by one give a result beyond float64 together.
	ErrDomain = errors.New("endocrine: inputs valid one by one but the formula is undefined for them together")
)

// checkFinite returns ErrNotFinite, wrapped with name, when v is NaN or ±Inf.
func checkFinite(name string, v float64) error {
	if !reduce.AllFinite(v) {
		return fmt.Errorf("%w: %s = %v", ErrNotFinite, name, v)
	}
	return nil
}

// checkPositive returns ErrOutOfRange, wrapped with name, unless v > 0.
func checkPositive(name string, v float64) error {
	if !(v > 0) {
		return fmt.Errorf("%w: %s = %v, want > 0", ErrOutOfRange, name, v)
	}
	return nil
}

// checkFinitePositive applies checkFinite then checkPositive.
func checkFinitePositive(name string, v float64) error {
	if err := checkFinite(name, v); err != nil {
		return err
	}
	return checkPositive(name, v)
}

// checkResult returns ErrDomain when a result computed from finite inputs is not finite.
func checkResult(v float64) error {
	if !reduce.AllFinite(v) {
		return fmt.Errorf("%w: result is beyond float64", ErrDomain)
	}
	return nil
}

// GlucoseManagementIndicatorPercent returns the glucose management indicator (GMI), in percent (the A1C
// scale), from the mean sensor glucose in mg/dL: 3.31 + 0.02392 × mean glucose. Not rounded.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range.
//
// ID:      MED-ENDO-001
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Bergenstal RM, Beck RW, Close KL, et al. Glucose Management Indicator (GMI): a new term for estimating A1C from continuous glucose monitoring. Diabetes Care. 2018;41(11):2275-2280. doi:10.2337/dc18-1581. PMID 30224348. PMC6196826, "How Is GMI Calculated?": "GMI (%) = 3.31 + 0.02392 × [mean glucose in mg/dL]"; Table 1
// Tags:    endocrine, diabetes, gmi, cgm, glucose, a1c
// Units:   meanGlucoseMgDL mg/dL (mean of continuous glucose monitoring readings); result percent
// Range:   meanGlucoseMgDL > 0 (derivation: 528 people with CGM data, up to 91 days; input limits not stated in the source)
// Example: 150 mg/dL → 3.31 + 3.588 = 6.898 (Table 1 prints 6.9%)
func GlucoseManagementIndicatorPercent(meanGlucoseMgDL float64) (float64, error) {
	if err := checkFinitePositive("meanGlucoseMgDL", meanGlucoseMgDL); err != nil {
		return 0, err
	}
	// 0.02392 × MaxFloat64 is finite, so no finite input overflows.
	return 3.31 + 0.02392*meanGlucoseMgDL, nil
}

// GlucoseManagementIndicatorMmolMol returns the glucose management indicator (GMI), in mmol/mol (the IFCC A1C
// scale), from the mean sensor glucose in mmol/L: 12.71 + 4.70587 × mean glucose. Not rounded.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range; ErrDomain when the result is beyond float64.
//
// ID:      MED-ENDO-002
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Bergenstal RM, Beck RW, Close KL, et al. Glucose Management Indicator (GMI): a new term for estimating A1C from continuous glucose monitoring. Diabetes Care. 2018;41(11):2275-2280. doi:10.2337/dc18-1581. PMID 30224348. PMC6196826, "How Is GMI Calculated?": "GMI (mmol/mol) = 12.71 + 4.70587 × [mean glucose in mmol/L]"
// Tags:    endocrine, diabetes, gmi, cgm, glucose, a1c, si
// Units:   meanGlucoseMmolL mmol/L (mean of continuous glucose monitoring readings); result mmol/mol
// Range:   meanGlucoseMmolL > 0 (input limits not stated in the source)
// Example: 10 mmol/L → 12.71 + 47.0587 = 59.7687 (HAND-DERIVED from the source equation)
func GlucoseManagementIndicatorMmolMol(meanGlucoseMmolL float64) (float64, error) {
	if err := checkFinitePositive("meanGlucoseMmolL", meanGlucoseMmolL); err != nil {
		return 0, err
	}
	r := 12.71 + 4.70587*meanGlucoseMmolL
	if err := checkResult(r); err != nil {
		return 0, err
	}
	return r, nil
}

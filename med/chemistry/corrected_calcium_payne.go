package chemistry

import (
	"errors"
	"fmt"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

var (
	// ErrNotFinite is returned when an input is NaN or ±Inf; it is checked before any range.
	ErrNotFinite = errors.New("chemistry: value is NaN or infinite")
	// ErrOutOfRange is returned when an input is outside the range the source validated.
	ErrOutOfRange = errors.New("chemistry: input outside the range the source validated")
	// ErrDomain is returned when inputs valid one by one give a result beyond float64 together.
	ErrDomain = errors.New("chemistry: inputs valid one by one but the formula is undefined for them together")
)

// KDIGO 2024 CKD guideline, Kidney Int. 2024;105(4S):S117–S314, doi:10.1016/j.kint.2023.10.018, p. S126
// "Conversion factors of conventional units to SI units" (read 2026-11-12): calcium 1 mg/dL = 0.2495 mmol/L.
const calciumMmolLPerMgDL = 0.2495

// gLPerGDL is grams per litre in one gram per decilitre, by definition of the decilitre (1 L = 10 dL).
const gLPerGDL = 10

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

// checkResult returns ErrDomain when a result computed from finite inputs is not finite.
func checkResult(v float64) error {
	if !reduce.AllFinite(v) {
		return fmt.Errorf("%w: result is beyond float64", ErrDomain)
	}
	return nil
}

// CorrectedCalciumPayne returns the albumin-adjusted serum calcium, in mg/dL, by Payne's formula:
// calcium − albumin + 4.0 (calcium in mg/100 mL, albumin in g/100 mL). The source states no limits; inputs
// far from those it studied can give a result ≤ 0, which is returned as computed.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range.
//
// ID:      MED-CHEM-001
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Payne RB, Little AJ, Williams RB, Milner JR. Interpretation of serum calcium in patients with abnormal serum proteins. Br Med J. 1973;4(5893):643-646. doi:10.1136/bmj.4.5893.643. PMID 4758544. Abstract: "Adjusted calcium = calcium − albumin + 4.0, where calcium is in mg/100 ml and albumin in g/100 ml."
// Tags:    chemistry, calcium, albumin, corrected-calcium, payne
// Units:   calciumMgDL mg/dL (= mg/100 mL); albuminGDL g/dL (= g/100 mL); result mg/dL
// Range:   calciumMgDL > 0; albuminGDL > 0 (population: 200 consecutive liver-function-test specimens; limits not stated by the source)
// Example: calcium 8.0, albumin 2.0 → 8.0 − 2.0 + 4.0 = 10.0 mg/dL (hand-derived)
func CorrectedCalciumPayne(calciumMgDL, albuminGDL float64) (float64, error) {
	if err := checkPayne("calciumMgDL", calciumMgDL, "albuminGDL", albuminGDL); err != nil {
		return 0, err
	}
	return correctedCalciumPayne(calciumMgDL, albuminGDL), nil
}

// CorrectedCalciumPayneSI returns the albumin-adjusted serum calcium, in mmol/L, by Payne's formula
// (MED-CHEM-001): calcium mmol/L ÷ 0.2495 gives mg/dL and albumin g/L ÷ 10 gives g/dL; the mg/dL result is
// multiplied by 0.2495 back to mmol/L.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range; ErrDomain when the result is beyond float64.
//
// ID:      MED-CHEM-002
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Payne RB, Little AJ, Williams RB, Milner JR. Interpretation of serum calcium in patients with abnormal serum proteins. Br Med J. 1973;4(5893):643-646. doi:10.1136/bmj.4.5893.643. PMID 4758544. Abstract: "Adjusted calcium = calcium − albumin + 4.0, where calcium is in mg/100 ml and albumin in g/100 ml."
// Source:  Kidney Disease: Improving Global Outcomes (KDIGO) CKD Work Group. KDIGO 2024 Clinical Practice Guideline for the Evaluation and Management of Chronic Kidney Disease. Kidney Int. 2024;105(4S):S117–S314. doi:10.1016/j.kint.2023.10.018. p. S126, "Conversion factors of conventional units to SI units" (calcium 0.2495)
// Tags:    chemistry, calcium, albumin, corrected-calcium, payne, si
// Units:   calciumMmolL mmol/L; albuminGL g/L; result mmol/L
// Range:   calciumMmolL > 0; albuminGL > 0 (limits not stated by the source)
// Example: calcium 1.996 mmol/L (= 8.0 mg/dL), albumin 20 g/L (= 2.0 g/dL) → 10.0 × 0.2495 = 2.495 mmol/L (hand-derived)
func CorrectedCalciumPayneSI(calciumMmolL, albuminGL float64) (float64, error) {
	if err := checkPayne("calciumMmolL", calciumMmolL, "albuminGL", albuminGL); err != nil {
		return 0, err
	}
	r := correctedCalciumPayne(calciumMmolL/calciumMmolLPerMgDL, albuminGL/gLPerGDL) * calciumMmolLPerMgDL
	if err := checkResult(r); err != nil {
		return 0, err
	}
	return r, nil
}

// checkPayne validates the inputs of MED-CHEM-001 and -002: both finiteness checks first, then ranges.
func checkPayne(calciumName string, calcium float64, albuminName string, albumin float64) error {
	if err := checkFinite(calciumName, calcium); err != nil {
		return err
	}
	if err := checkFinite(albuminName, albumin); err != nil {
		return err
	}
	if err := checkPositive(calciumName, calcium); err != nil {
		return err
	}
	return checkPositive(albuminName, albumin)
}

// correctedCalciumPayne is Payne's 1973 formula on validated inputs in mg/dL and g/dL.
func correctedCalciumPayne(calciumMgDL, albuminGDL float64) float64 {
	return calciumMgDL - albuminGDL + 4.0
}

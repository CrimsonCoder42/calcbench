package respiratory

import (
	"errors"
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// The package's sentinel errors. They live in the file of the first function that uses them
// (doc/ASSUMPTIONS.md A-008). Every error wraps one of these with the offending parameter; test with errors.Is.
var (
	// ErrNotFinite: an input is NaN or ±Inf.
	ErrNotFinite = errors.New("respiratory: value is NaN or infinite")
	// ErrOutOfRange: an input is outside the range the source validated.
	ErrOutOfRange = errors.New("respiratory: input outside the range the source validated")
	// ErrDomain: the inputs are valid one by one but the result is undefined or beyond float64 for them together.
	ErrDomain = errors.New("respiratory: inputs valid one by one but the formula is undefined for them together")
)

// checkPositiveFinite returns ErrNotFinite for NaN or ±Inf and ErrOutOfRange unless v > 0.
func checkPositiveFinite(name string, v float64) error {
	if !reduce.AllFinite(v) {
		return fmt.Errorf("%w: %s = %v", ErrNotFinite, name, v)
	}
	if !(v > 0) {
		return fmt.Errorf("%w: %s = %v, want > 0", ErrOutOfRange, name, v)
	}
	return nil
}

// DiffusingCapacityRoughtonForster returns the lung diffusing capacity, in mL/min/mmHg, from its two
// components by the Roughton–Forster equation 1/DL = 1/DM + 1/(θ·VC): DM the membrane diffusing capacity,
// θ the gas uptake of 1 mL of blood per minute per mmHg, VC the pulmonary capillary blood volume.
// It is computed as DM·θVC / (DM + θVC), the form without three reciprocals.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange for a non-positive input; ErrDomain when θ·VC or
// DM + θ·VC overflows float64.
//
// ID:      MED-RESP-001
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Zavorsky GS, Hsia CCW, Hughes JMB, et al. Standardisation and application of the single-breath determination of nitric oxide uptake in the lung. Eur Respir J. 2017;49(2):1600962. doi:10.1183/13993003.00962-2016. p. 2 "Origins of DLNO" and Figure 1 legend: "Roughton–Forster equation, 1/DL=1/DM+1/(θ·VC)"
// Source:  Roughton FJW, Forster RE. Relative importance of diffusion and chemical reaction rates in determining rate of exchange of gases in the human lung, with special reference to true diffusing capacity of pulmonary membrane and volume of blood in the lung capillaries. J Appl Physiol. 1957;11(2):290-302. doi:10.1152/jappl.1957.11.2.290. PMID 13475180
// Tags:    respiratory, diffusing-capacity, dlco, roughton-forster, gas-exchange
// Units:   membraneDLMLMinMmHg mL/min/mmHg; thetaMLMinMmHgPerMLBlood mL/min/mmHg per mL of blood; capillaryBloodVolumeML mL; result mL/min/mmHg
// Range:   every input > 0 (the relation is physiological, with no population limits stated in the source)
// Example: DM 50, θ 1.0, VC 75 → 50·75/125 = 30 mL/min/mmHg (hand-derived)
func DiffusingCapacityRoughtonForster(membraneDLMLMinMmHg, thetaMLMinMmHgPerMLBlood, capillaryBloodVolumeML float64) (float64, error) {
	if err := checkPositiveFinite("membraneDLMLMinMmHg", membraneDLMLMinMmHg); err != nil {
		return 0, err
	}
	if err := checkPositiveFinite("thetaMLMinMmHgPerMLBlood", thetaMLMinMmHgPerMLBlood); err != nil {
		return 0, err
	}
	if err := checkPositiveFinite("capillaryBloodVolumeML", capillaryBloodVolumeML); err != nil {
		return 0, err
	}
	blood := thetaMLMinMmHgPerMLBlood * capillaryBloodVolumeML
	if math.IsInf(blood, 0) || math.IsInf(membraneDLMLMinMmHg+blood, 0) {
		return 0, fmt.Errorf("%w: θ·VC or DM + θ·VC overflows", ErrDomain)
	}
	// Harmonic form: DM·B/(DM+B) = DM / (1 + DM/B) avoids overflow of the product DM·B.
	return membraneDLMLMinMmHg / (1 + membraneDLMLMinMmHg/blood), nil
}

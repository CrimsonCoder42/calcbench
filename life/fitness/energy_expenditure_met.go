package fitness

import "fmt"

// EnergyExpenditureMET returns the energy expended in an activity, in kcal, from its MET value, the body mass
// and the duration: MET × bodyMassKg × durationHours, using the Compendium's definition 1 MET ≈ 1.0 kcal/kg/h.
// The MET value of an activity is the caller's, read from the Compendium's activity table (not held here).
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange when met ≤ 0, bodyMassKg ≤ 0 or durationHours < 0;
// ErrDomain when the result is beyond float64.
//
// ID:      LIF-FIT-015
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Herrmann SD, Willis EA, Ainsworth BE, et al. 2024 Adult Compendium of Physical Activities: a third update of the energy costs of human activities. J Sport Health Sci. 2024;13(1):6-12. doi:10.1016/j.jshs.2023.10.010. PMC10818145. §2.5: 1 MET = "3.5 mL/kg/min or approximately 1.0 kcal/kg/h"
// Tags:    fitness, met, energy-expenditure, calories, compendium, physical-activity
// Units:   met METs (multiples of resting metabolic rate); bodyMassKg kg; durationHours h; result kcal
// Range:   met > 0; bodyMassKg > 0; durationHours ≥ 0 (0 h gives 0 kcal)
// Example: 6.0 METs, 70 kg, 0.5 h → 6.0 × 70 × 0.5 = 210 kcal (HAND-DERIVED from the §2.5 definition)
func EnergyExpenditureMET(met, bodyMassKg, durationHours float64) (float64, error) {
	if err := checkFinite("met", met); err != nil {
		return 0, err
	}
	if err := checkFinite("bodyMassKg", bodyMassKg); err != nil {
		return 0, err
	}
	if err := checkFinite("durationHours", durationHours); err != nil {
		return 0, err
	}
	if !(met > 0) {
		return 0, fmt.Errorf("%w: met = %v, want > 0", ErrOutOfRange, met)
	}
	if !(bodyMassKg > 0) {
		return 0, fmt.Errorf("%w: bodyMassKg = %v, want > 0", ErrOutOfRange, bodyMassKg)
	}
	if durationHours < 0 {
		return 0, fmt.Errorf("%w: durationHours = %v, want ≥ 0", ErrOutOfRange, durationHours)
	}
	e := met * bodyMassKg * durationHours
	if err := checkResult(e); err != nil {
		return 0, err
	}
	return e, nil
}

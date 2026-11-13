package temperature

import (
	"fmt"
	"math"
)

// CelsiusToKelvin returns the thermodynamic temperature, in K, of a Celsius temperature: T/K = t/°C + 273.15.
// Errors: ErrNotFinite for NaN or ±Inf; ErrBelowAbsoluteZero below −273.15 °C.
//
// ID:      UNT-TEMP-003
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Thompson A, Taylor BN. Guide for the Use of the International System of Units (SI). NIST Special Publication 811, 2008 ed., Appendix B.9 "Temperature": "T/K = t/°C + 273.15"
// Tags:    unit, temperature, celsius, kelvin, conversion
// Units:   degreesC °C; result K
// Range:   degreesC ≥ −273.15 (absolute zero)
// Example: 0 °C → 273.15 K; 100 °C → 373.15 K; −273.15 °C → 0 K
func CelsiusToKelvin(degreesC float64) (float64, error) {
	if math.IsNaN(degreesC) || math.IsInf(degreesC, 0) {
		return 0, fmt.Errorf("%w: degreesC = %v", ErrNotFinite, degreesC)
	}
	if degreesC < absoluteZeroC {
		return 0, fmt.Errorf("%w: degreesC = %v, want ≥ %v", ErrBelowAbsoluteZero, degreesC, absoluteZeroC)
	}
	return degreesC + kelvinOffset, nil
}

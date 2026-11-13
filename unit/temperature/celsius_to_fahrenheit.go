package temperature

import (
	"fmt"
	"math"
)

// CelsiusToFahrenheit returns the Fahrenheit temperature, in °F, of a Celsius temperature:
// t/°F = 1.8 · t/°C + 32, the inverse of NIST's t/°C = (t/°F − 32)/1.8.
// Errors: ErrNotFinite for NaN or ±Inf; ErrBelowAbsoluteZero below −273.15 °C; ErrOverflow when the result is
// beyond float64 (degreesC above about 9.99e307).
//
// ID:      UNT-TEMP-002
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Thompson A, Taylor BN. Guide for the Use of the International System of Units (SI). NIST Special Publication 811, 2008 ed., Appendix B.9 "Temperature": "t/°C = (t/°F - 32)/1.8", solved for t/°F
// Tags:    unit, temperature, celsius, fahrenheit, conversion
// Units:   degreesC °C; result °F
// Range:   degreesC ≥ −273.15 (absolute zero)
// Example: 100 °C → 212 °F; 0 °C → 32 °F; 37 °C → 98.6 °F; −40 °C → −40 °F
func CelsiusToFahrenheit(degreesC float64) (float64, error) {
	if math.IsNaN(degreesC) || math.IsInf(degreesC, 0) {
		return 0, fmt.Errorf("%w: degreesC = %v", ErrNotFinite, degreesC)
	}
	if degreesC < absoluteZeroC {
		return 0, fmt.Errorf("%w: degreesC = %v, want ≥ %v", ErrBelowAbsoluteZero, degreesC, absoluteZeroC)
	}
	f := degreesC*fahrenheitPerCelsius + fahrenheitOffset
	if math.IsInf(f, 0) {
		return 0, fmt.Errorf("%w: degreesC = %v", ErrOverflow, degreesC)
	}
	return f, nil
}

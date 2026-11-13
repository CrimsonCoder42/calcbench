package temperature

import (
	"fmt"
	"math"
)

// KelvinToCelsius returns the Celsius temperature, in °C, of a thermodynamic temperature: t/°C = T/K − 273.15.
// Errors: ErrNotFinite for NaN or ±Inf; ErrBelowAbsoluteZero below 0 K.
//
// ID:      UNT-TEMP-004
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Thompson A, Taylor BN. Guide for the Use of the International System of Units (SI). NIST Special Publication 811, 2008 ed., Appendix B.9 "Temperature": "t/°C = T/K - 273.15"
// Tags:    unit, temperature, kelvin, celsius, conversion
// Units:   kelvin K; result °C
// Range:   kelvin ≥ 0
// Example: 273.15 K → 0 °C; 373.15 K → 100 °C; 0 K → −273.15 °C
func KelvinToCelsius(kelvin float64) (float64, error) {
	if math.IsNaN(kelvin) || math.IsInf(kelvin, 0) {
		return 0, fmt.Errorf("%w: kelvin = %v", ErrNotFinite, kelvin)
	}
	if kelvin < 0 {
		return 0, fmt.Errorf("%w: kelvin = %v, want ≥ 0", ErrBelowAbsoluteZero, kelvin)
	}
	return kelvin - kelvinOffset, nil
}

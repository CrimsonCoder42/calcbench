package temperature

import (
	"errors"
	"fmt"
	"math"
)

var (
	// ErrNotFinite is returned when an input is NaN or ±Inf.
	ErrNotFinite = errors.New("temperature: value is NaN or infinite")
	// ErrBelowAbsoluteZero is returned when an input is colder than 0 K.
	ErrBelowAbsoluteZero = errors.New("temperature: value is below absolute zero")
	// ErrOverflow: a finite input whose result is beyond float64.
	ErrOverflow = errors.New("temperature: result overflows float64")
)

// Absolute zero on each scale. NIST SP 811 (2008) App. B.9: T/K = t/°C + 273.15 and
// T/K = (t/°F + 459.67)/1.8, so T = 0 K is t = −273.15 °C and t = −459.67 °F.
const (
	absoluteZeroC = -273.15
	absoluteZeroF = -459.67
	// fahrenheitPerCelsius is the size ratio of the two degrees, 1.8 exactly (NIST SP 811 App. B.9).
	fahrenheitPerCelsius = 1.8
	// fahrenheitOffset is the Fahrenheit reading at 0 °C, 32 exactly (NIST SP 811 App. B.9).
	fahrenheitOffset = 32.0
	// kelvinOffset is the kelvin reading at 0 °C, 273.15 exactly (NIST SP 811 App. B.9).
	kelvinOffset = 273.15
)

// FahrenheitToCelsius returns the Celsius temperature, in °C, of a Fahrenheit temperature:
// t/°C = (t/°F − 32)/1.8.
// Errors: ErrNotFinite for NaN or ±Inf; ErrBelowAbsoluteZero below −459.67 °F.
//
// ID:      UNT-TEMP-001
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Thompson A, Taylor BN. Guide for the Use of the International System of Units (SI). NIST Special Publication 811, 2008 ed., Appendix B.9 "Temperature": "t/°C = (t/°F - 32)/1.8"
// Tags:    unit, temperature, fahrenheit, celsius, conversion
// Units:   degreesF °F; result °C
// Range:   degreesF ≥ −459.67 (absolute zero)
// Example: 212 °F → 100 °C; 32 °F → 0 °C; 98.6 °F → 37 °C; −40 °F → −40 °C
func FahrenheitToCelsius(degreesF float64) (float64, error) {
	if math.IsNaN(degreesF) || math.IsInf(degreesF, 0) {
		return 0, fmt.Errorf("%w: degreesF = %v", ErrNotFinite, degreesF)
	}
	if degreesF < absoluteZeroF {
		return 0, fmt.Errorf("%w: degreesF = %v, want ≥ %v", ErrBelowAbsoluteZero, degreesF, absoluteZeroF)
	}
	return (degreesF - fahrenheitOffset) / fahrenheitPerCelsius, nil
}

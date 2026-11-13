// Package temperature converts temperatures between the degree Fahrenheit, the degree Celsius and the kelvin,
// by the exact relations of the NIST Guide to the SI (NIST SP 811, 2008, Appendix B.9, "Temperature").
//
// ID short form: TEMP (IDs are UNT-TEMP-NNN).
//
// Units: every input and result names its scale (degreesF, degreesC, kelvin). A temperature below absolute
// zero (0 K = −273.15 °C = −459.67 °F) returns ErrBelowAbsoluteZero; NaN and ±Inf return ErrNotFinite.
// Results are not rounded.
//
// Functions:
//
//	UNT-TEMP-001  FahrenheitToCelsius  t/°C = (t/°F − 32)/1.8
//	UNT-TEMP-002  CelsiusToFahrenheit  t/°F = 1.8 · t/°C + 32
//	UNT-TEMP-003  CelsiusToKelvin      T/K = t/°C + 273.15
//	UNT-TEMP-004  KelvinToCelsius      t/°C = T/K − 273.15
package temperature

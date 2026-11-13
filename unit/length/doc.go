// Package length converts lengths between the centimeter and the inch, by the exact relation of the NIST
// Guide to the SI (NIST SP 811, 2008, Appendix B.8): 1 in = 2.54 cm exactly.
//
// ID short form: LEN (IDs are UNT-LEN-NNN).
//
// Units: every input and result names its unit (centimeters, inches). A negative length returns
// ErrNegative; NaN and ±Inf return ErrNotFinite; a finite input whose result is beyond float64 returns
// ErrOverflow. Results are not rounded.
//
// Functions:
//
//	UNT-LEN-001  CentimetersToInches  l/in = (l/cm)/2.54
//	UNT-LEN-002  InchesToCentimeters  l/cm = 2.54 · l/in
package length

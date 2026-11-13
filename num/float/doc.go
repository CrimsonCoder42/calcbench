// Package float holds IEEE 754 double-precision helpers that Go's math package lacks: tolerance
// comparison, the unit in the last place, non-finite replacement and two machine constants.
// math.Nextafter, math.Copysign, math.MaxFloat64 and math.SmallestNonzeroFloat64 are Go's own and are not
// wrapped (D-012).
//
// ID short form: FLT (IDs are NUM-FLT-NNN).
//
// Functions and constants:
//
//	NUM-FLT-001  MachineEpsilon    2⁻⁵², the gap between 1 and the next double
//	NUM-FLT-002  SmallestNormal    2⁻¹⁰²², the smallest positive normal double
//	NUM-FLT-003  IsClose           symmetric closeness test |a−b| ≤ max(rel·max(|a|,|b|), abs)
//	NUM-FLT-004  AllClose          IsClose for every pair of two slices
//	NUM-FLT-005  Ulp               unit in the last place of x
//	NUM-FLT-006  ReplaceNonFinite  NaN and ±Inf replaced by given values
package float

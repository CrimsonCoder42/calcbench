// Package reduce folds a slice to one value, or to its running values: compensated sum, product,
// cumulative sum and product, finite differences, the index of the minimum or maximum, and a finite check.
// Every package in calcbench may import it; it imports nothing from calcbench (doc/DECISIONS.md D-081).
//
// ID short form: RED (IDs are ALG-RED-NNN).
//
// Floats: a NaN or ±Inf input returns ErrNotFinite (wrapped with its index); a finite input whose result,
// or a partial result, is beyond ±MaxFloat64 returns ErrOverflow. No function returns a silent NaN or Inf.
// Sums are Neumaier-compensated (Neumaier 1974; Higham, Accuracy and Stability of Numerical Algorithms,
// 2nd ed., SIAM, 2002, §4.3), so the error does not grow with n for well-conditioned data. Products keep
// the exponent apart from the mantissa (math.Frexp), so an intermediate product never overflows or
// underflows while the true result is representable.
//
// Output: functions that produce a slice write it into dst, reused when cap(dst) is large enough, else one
// new slice is allocated. dst may be the input itself (the same first element) to work in place; any other
// overlap is unsupported. The input is never modified unless it is dst.
//
// Functions:
//
//	ALG-RED-001  Sum          Neumaier-compensated sum of a slice
//	ALG-RED-002  NeumaierAdd  one step of a Neumaier-compensated running sum
//	ALG-RED-003  TwoSum       Knuth's error-free sum: s + e = a + b exactly
//	ALG-RED-004  Prod         product of a slice, with no spurious overflow or underflow
//	ALG-RED-005  CumSum       compensated running sums
//	ALG-RED-006  CumProd      running products
//	ALG-RED-007  Diff         the n-th forward difference
//	ALG-RED-008  ArgMin       index of the first minimum
//	ALG-RED-009  ArgMax       index of the first maximum
//	ALG-RED-010  AllFinite    whether no value is NaN or ±Inf
package reduce

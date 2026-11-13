package float

// MachineEpsilon is the float64 machine epsilon 2⁻⁵², the gap between 1 and the next larger double
// (numpy finfo(float64).eps). Go's math package has no such constant.
//
// ID:      NUM-FLT-001
// Time:    O(1)   a constant
// Space:   O(1), no allocation
// Source:  IEEE Std 754-2019, IEEE Standard for Floating-Point Arithmetic, §3.3 and Table 3.5 (binary64: p = 53, emax = 1023)
// Source:  Higham, Accuracy and Stability of Numerical Algorithms, 2nd ed., SIAM, 2002, §2.1 (ε = β^(1−t))
// Tags:    constant, floating-point, machine-epsilon, ieee-754, finfo
// Units:   dimensionless
// Example: 2.220446049250313e-16
const MachineEpsilon = 0x1p-52

// SmallestNormal is the smallest positive normal float64, 2⁻¹⁰²² (numpy finfo(float64).tiny /
// smallest_normal). Go's math.SmallestNonzeroFloat64 is the smallest subnormal, 2⁻¹⁰⁷⁴, a different value.
//
// ID:      NUM-FLT-002
// Time:    O(1)   a constant
// Space:   O(1), no allocation
// Source:  IEEE Std 754-2019, IEEE Standard for Floating-Point Arithmetic, §3.3 and Table 3.5 (binary64: emin = −1022)
// Tags:    constant, floating-point, smallest-normal, ieee-754, finfo, underflow
// Units:   dimensionless
// Example: 2.2250738585072014e-308
const SmallestNormal = 0x1p-1022

package logs

// GoldenRatio is the golden ratio φ = (1+√5)/2, the positive root of x² = x + 1, to 63 significant
// digits (an untyped constant; as a float64 it rounds to the nearest double, which equals math.Phi).
//
// ID:      NUM-LOGS-001
// Time:    O(1)   a constant
// Space:   O(1), no allocation
// Source:  Knuth, The Art of Computer Programming, Vol. 1: Fundamental Algorithms, 3rd ed., Addison-Wesley, 1997, §1.2.8 (Fibonacci numbers; φ = (1+√5)/2)
// Source:  OEIS A001622, "Decimal expansion of golden ratio phi (or tau) = (1 + sqrt(5))/2", https://oeis.org/A001622 (read 2026-11-12)
// Tags:    constant, golden-ratio, phi, irrational, fibonacci
// Units:   dimensionless
// Example: 1.6180339887498948482045868343656381177203091798057628621354486 (OEIS A001622, first 62 digits)
const GoldenRatio = 1.61803398874989484820458683436563811772030917980576286213544862

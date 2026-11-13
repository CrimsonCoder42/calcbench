package linalg

import "fmt"

// MatPowerInteger returns Aᵏ for an n×n matrix and any integer k: MatPower (NUM-LINALG-122) for k ≥ 0, and
// (A⁻¹)⁻ᵏ for k < 0 with A⁻¹ from MatInverse, in dst (which must not overlap a).
// Errors: ErrOutOfRange when |k| > 2⁴⁰; ErrSingular for k < 0 and a singular A; otherwise as MatPower and
// MatInverse.
//
// ID:      NUM-LINALG-183
// Time:    O(n³·log|k|)   at most 2⌊log₂|k|⌋ products, plus one O(n³) inverse when k < 0
// Space:   O(n²) extra; O(n²) for the result, one allocation only if dst lacks capacity
// Source:  Knuth, The Art of Computer Programming, Vol. 2: Seminumerical Algorithms, 3rd ed., Addison-Wesley, 1997, §4.6.3, Algorithm A; Strang, Introduction to Linear Algebra, 5th ed., 2016, §2.5 (A⁻ᵏ = (A⁻¹)ᵏ)
// Tags:    matrix, n-dimensional, power, negative-power, inverse, repeated-squaring, linear-algebra
// Units:   a in unit A; the result in Aᵏ
// Example: [2 0; 0 4]⁻² → [0.25 0; 0 0.0625] (hand-checked)
func MatPowerInteger(dst, a []float64, n, k int) ([]float64, error) {
	if k >= 0 {
		return MatPower(dst, a, n, k)
	}
	if k < -maxMatPower {
		return nil, fmt.Errorf("%w: k = %d, want |k| ≤ 2^40", ErrOutOfRange, k)
	}
	inv, err := MatInverse(nil, a, n)
	if err != nil {
		return nil, err
	}
	return MatPower(dst, inv, n, -k)
}

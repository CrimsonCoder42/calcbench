package linalg

import "github.com/CrimsonCoder42/calcbench/algo/reduce"

// conditionExact returns ‖A‖·‖A⁻¹‖ for the given norm, A⁻¹ from MatInverse.
func conditionExact(a []float64, n int, norm func([]float64, int, int) (float64, error)) (float64, error) {
	inv, err := MatInverse(nil, a, n)
	if err != nil {
		return 0, err
	}
	na, err := norm(a, n, n)
	if err != nil {
		return 0, err
	}
	ni, err := norm(inv, n, n)
	if err != nil {
		return 0, err
	}
	c := na * ni
	if !reduce.AllFinite(c) {
		return 0, errOverflow("condition number")
	}
	return c, nil
}

// MatCondition1 returns the exact 1-norm condition number κ₁(A) = ‖A‖₁·‖A⁻¹‖₁ of an n×n matrix, with A⁻¹
// formed by MatInverse (MatConditionEstimate1, NUM-LINALG-100, estimates it in O(n²) after the LU).
// Errors: ErrSingular when A is singular; otherwise as MatInverse and MatNorm1; ErrOverflow when κ leaves
// the float64 range.
//
// ID:      NUM-LINALG-180
// Time:    O(n³)   n = order of the matrix
// Space:   O(n²) for the factor and the inverse
// Source:  Higham, Accuracy and Stability of Numerical Algorithms, 2nd ed., SIAM, 2002, §6.2 (κ(A) = ‖A‖‖A⁻¹‖ for a subordinate norm) and §7.1
// Tags:    matrix, n-dimensional, condition-number, one-norm, conditioning, linear-algebra
// Units:   a in unit A; the result dimensionless
// Example: [1 2; 3 4] → ‖A‖₁ = 6, ‖A⁻¹‖₁ = 3.5 → 21 (A⁻¹ = [−2 1; 1.5 −0.5]; hand-checked)
func MatCondition1(a []float64, n int) (float64, error) { return conditionExact(a, n, MatNorm1) }

// MatConditionInf returns the exact ∞-norm condition number κ∞(A) = ‖A‖∞·‖A⁻¹‖∞ of an n×n matrix.
// Errors: as MatCondition1.
//
// ID:      NUM-LINALG-181
// Time:    O(n³)   n = order of the matrix
// Space:   O(n²) for the factor and the inverse
// Source:  Higham, Accuracy and Stability of Numerical Algorithms, 2nd ed., SIAM, 2002, §6.2 and §7.1
// Tags:    matrix, n-dimensional, condition-number, infinity-norm, conditioning, linear-algebra
// Units:   a in unit A; the result dimensionless
// Example: [1 2; 3 4] → ‖A‖∞ = 7, ‖A⁻¹‖∞ = 3 → 21 (hand-checked)
func MatConditionInf(a []float64, n int) (float64, error) { return conditionExact(a, n, MatNormInf) }

// MatConditionFrobenius returns κ_F(A) = ‖A‖_F·‖A⁻¹‖_F of an n×n matrix.
// Errors: as MatCondition1.
//
// ID:      NUM-LINALG-182
// Time:    O(n³)   n = order of the matrix
// Space:   O(n²) for the factor and the inverse
// Source:  Higham, Accuracy and Stability of Numerical Algorithms, 2nd ed., SIAM, 2002, §6.2 (condition numbers in the Frobenius norm)
// Tags:    matrix, n-dimensional, condition-number, frobenius, conditioning, linear-algebra
// Units:   a in unit A; the result dimensionless
// Example: [1 2; 3 4] → √30·√7.5 = 15 (hand-checked)
func MatConditionFrobenius(a []float64, n int) (float64, error) {
	return conditionExact(a, n, MatNormFrobenius)
}

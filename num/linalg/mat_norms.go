package linalg

import "github.com/CrimsonCoder42/calcbench/algo/reduce"

// MatNormFrobenius returns the Frobenius norm ‖A‖_F = √(Σᵢⱼ aᵢⱼ²) of a rows×cols matrix, scaled by a power
// of two as VecNorm is, so no intermediate square overflows or underflows.
// Errors: ErrDimension when len(a) disagrees with rows×cols; ErrNotFinite for a NaN or ±Inf element;
// ErrOverflow when the norm is outside the float64 range.
//
// ID:      NUM-LINALG-163
// Time:    O(r·c)   r = rows, c = cols
// Space:   O(1), no allocation
// Source:  Golub, Van Loan, Matrix Computations, 4th ed., Johns Hopkins University Press, 2013, §2.3 "Matrix Norms" (the Frobenius norm ‖A‖_F)
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §11.2 "Norms and Condition Numbers" (‖A‖_F² is the sum of all aᵢⱼ²)
// Tags:    matrix, n-dimensional, norm, frobenius, linear-algebra
// Units:   a in unit A; the result in A
// Example: [1 2; 3 4] → √30 ≈ 5.477225575 (hand-checked)
func MatNormFrobenius(a []float64, rows, cols int) (float64, error) {
	if err := checkMatrix("a", a, rows, cols); err != nil {
		return 0, err
	}
	return VecNorm(a)
}

// MatNorm1 returns the 1-norm ‖A‖₁ = maxⱼ Σᵢ |aᵢⱼ|, the largest absolute column sum, of a rows×cols matrix.
// Each column sum is compensated (Neumaier) so the result does not depend on summation order.
// Errors: ErrDimension when len(a) disagrees with rows×cols; ErrNotFinite for a NaN or ±Inf element;
// ErrOverflow when a column sum is outside the float64 range.
//
// ID:      NUM-LINALG-164
// Time:    O(r·c)   r = rows, c = cols
// Space:   O(1), no allocation
// Source:  Golub, Van Loan, Matrix Computations, 4th ed., Johns Hopkins University Press, 2013, §2.3 "Matrix Norms" (‖A‖₁ = max over j of Σᵢ |aᵢⱼ|)
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §11.2 "Norms and Condition Numbers" (the l¹ norm of A: largest column sum)
// Tags:    matrix, n-dimensional, norm, one-norm, column-sum, linear-algebra
// Units:   a in unit A; the result in A
// Example: [1 −2; 3 4] → max(1 + 3, 2 + 4) = 6 (hand-checked)
func MatNorm1(a []float64, rows, cols int) (float64, error) {
	if err := checkMatrix("a", a, rows, cols); err != nil {
		return 0, err
	}
	if i := firstNotFinite(a); i >= 0 {
		return 0, notFiniteAt("a", i)
	}
	best := 0.0
	for j := range cols {
		s, c := 0.0, 0.0
		for i := range rows {
			s, c = reduce.NeumaierAdd(s, c, abs(a[i*cols+j]))
		}
		best = max(best, s+c)
	}
	if !reduce.AllFinite(best) {
		return 0, errOverflow("column sum")
	}
	return best, nil
}

// MatNormInf returns the infinity norm ‖A‖∞ = maxᵢ Σⱼ |aᵢⱼ|, the largest absolute row sum, of a rows×cols
// matrix, each row sum compensated (Neumaier).
// Errors: ErrDimension when len(a) disagrees with rows×cols; ErrNotFinite for a NaN or ±Inf element;
// ErrOverflow when a row sum is outside the float64 range.
//
// ID:      NUM-LINALG-165
// Time:    O(r·c)   r = rows, c = cols
// Space:   O(1), no allocation
// Source:  Golub, Van Loan, Matrix Computations, 4th ed., Johns Hopkins University Press, 2013, §2.3 "Matrix Norms" (‖A‖∞ = max over i of Σⱼ |aᵢⱼ|)
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §11.2 "Norms and Condition Numbers" (the l∞ norm of A: largest row sum)
// Tags:    matrix, n-dimensional, norm, infinity-norm, row-sum, linear-algebra
// Units:   a in unit A; the result in A
// Example: [1 −2; 3 4] → max(1 + 2, 3 + 4) = 7 (hand-checked)
func MatNormInf(a []float64, rows, cols int) (float64, error) {
	if err := checkMatrix("a", a, rows, cols); err != nil {
		return 0, err
	}
	if i := firstNotFinite(a); i >= 0 {
		return 0, notFiniteAt("a", i)
	}
	best := 0.0
	for i := range rows {
		s, c := 0.0, 0.0
		for _, x := range a[i*cols : i*cols+cols] {
			s, c = reduce.NeumaierAdd(s, c, abs(x))
		}
		best = max(best, s+c)
	}
	if !reduce.AllFinite(best) {
		return 0, errOverflow("row sum")
	}
	return best, nil
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

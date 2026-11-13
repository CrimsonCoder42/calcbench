package linalg

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// General matrices are []float64 in row-major order with their dimensions passed alongside: element
// (i, j) of an r×c matrix is a[i*c+j]. A result is written into dst, reused when its capacity suffices
// and allocated once otherwise; dst must not overlap an input unless the function says it may.

// checkMatrix reports ErrDimension unless rows, cols ≥ 0, rows·cols does not overflow, and
// len(a) = rows·cols.
func checkMatrix(name string, a []float64, rows, cols int) error {
	if rows < 0 || cols < 0 {
		return fmt.Errorf("%w: %s is %d×%d", ErrDimension, name, rows, cols)
	}
	if rows > 0 && cols > math.MaxInt/rows {
		return fmt.Errorf("%w: %s is %d×%d (too large)", ErrDimension, name, rows, cols)
	}
	if len(a) != rows*cols {
		return fmt.Errorf("%w: len(%s) = %d, want %d×%d = %d", ErrDimension, name, len(a), rows, cols, rows*cols)
	}
	return nil
}

// firstNotFinite returns the index of the first NaN or ±Inf in a, or −1.
func firstNotFinite(a []float64) int {
	for i, x := range a {
		if !reduce.AllFinite(x) {
			return i
		}
	}
	return -1
}

// MatIdentity returns the n×n identity matrix in dst.
// Errors: ErrDimension when n < 0 or n·n overflows.
//
// ID:      NUM-LINALG-089
// Time:    O(n²)   n = order of the matrix
// Space:   O(1) extra; O(n²) for the result, one allocation only if dst lacks capacity
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §2.5 "Inverse Matrices" (the identity I)
// Tags:    matrix, n-dimensional, identity, linear-algebra
// Units:   dimensionless
// Example: n = 2 → [1 0; 0 1]
func MatIdentity(dst []float64, n int) ([]float64, error) {
	if n < 0 || (n > 0 && n > math.MaxInt/n) {
		return nil, fmt.Errorf("%w: n = %d", ErrDimension, n)
	}
	dst = grow(dst, n*n)
	clear(dst)
	for i := range n {
		dst[i*n+i] = 1
	}
	return dst, nil
}

// MatMul returns the product of the rows×inner matrix a and the inner×cols matrix b, a rows×cols matrix,
// in dst (which must not overlap a or b). Each element is accumulated with fused multiply-adds in i-k-j
// order (stride-1 access to b and dst).
// Errors: ErrDimension when a slice's length disagrees with its dimensions; ErrNotFinite for a NaN or ±Inf
// element; ErrOverflow when an element of the product is outside the float64 range.
//
// ID:      NUM-LINALG-090
// Time:    O(rows·inner·cols)
// Space:   O(1) extra; O(rows·cols) for the result, one allocation only if dst lacks capacity
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §2.4 "Rules for Matrix Operations" (row times column)
// Tags:    matrix, n-dimensional, multiply, product, gemm, linear-algebra
// Units:   a in unit A, b in unit B; the result in A·B
// Example: [1 2 3; 4 5 6] (2×3) · [7 8; 9 10; 11 12] (3×2) → [58 64; 139 154] (hand-checked)
func MatMul(dst, a []float64, rows, inner int, b []float64, cols int) ([]float64, error) {
	if err := checkMatrix("a", a, rows, inner); err != nil {
		return nil, err
	}
	if err := checkMatrix("b", b, inner, cols); err != nil {
		return nil, err
	}
	if rows > 0 && cols > math.MaxInt/rows {
		return nil, fmt.Errorf("%w: result is %d×%d (too large)", ErrDimension, rows, cols)
	}
	dst = grow(dst, rows*cols)
	clear(dst)
	for i := range rows {
		ci := dst[i*cols : i*cols+cols]
		for k := range inner {
			aik := a[i*inner+k]
			bk := b[k*cols : k*cols+cols]
			for j := range ci {
				ci[j] = math.FMA(aik, bk[j], ci[j])
			}
		}
	}
	if i := firstNotFinite(dst); i >= 0 {
		if j := firstNotFinite(a); j >= 0 {
			return nil, fmt.Errorf("%w: a[%d]", ErrNotFinite, j)
		}
		if j := firstNotFinite(b); j >= 0 {
			return nil, fmt.Errorf("%w: b[%d]", ErrNotFinite, j)
		}
		return nil, fmt.Errorf("%w: element %d", ErrOverflow, i)
	}
	return dst, nil
}

// MatVecMul returns the product of the rows×cols matrix a and the column vector x (length cols), a
// vector of length rows, in dst (which must not overlap a or x).
// Errors: ErrDimension when a's length disagrees with its dimensions or len(x) ≠ cols; ErrNotFinite for
// a NaN or ±Inf element; ErrOverflow when an element of the result is outside the float64 range.
//
// ID:      NUM-LINALG-091
// Time:    O(rows·cols)
// Space:   O(1) extra; O(rows) for the result, one allocation only if dst lacks capacity
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §2.4 "Rules for Matrix Operations" (matrix times vector)
// Tags:    matrix, n-dimensional, matrix-vector, multiply, gemv, linear-algebra
// Units:   a in unit A, x in unit B; the result in A·B
// Example: [1 2 3; 4 5 6]·(1, 0, −1) → (−2, −2) (hand-checked)
func MatVecMul(dst, a []float64, rows, cols int, x []float64) ([]float64, error) {
	if err := checkMatrix("a", a, rows, cols); err != nil {
		return nil, err
	}
	if len(x) != cols {
		return nil, fmt.Errorf("%w: len(x) = %d, want %d", ErrDimension, len(x), cols)
	}
	dst = grow(dst, rows)
	for i := range rows {
		s := 0.0
		r := a[i*cols : i*cols+cols]
		for j, v := range r {
			s = math.FMA(v, x[j], s)
		}
		dst[i] = s
	}
	if i := firstNotFinite(dst); i >= 0 {
		if j := firstNotFinite(a); j >= 0 {
			return nil, fmt.Errorf("%w: a[%d]", ErrNotFinite, j)
		}
		if j := firstNotFinite(x); j >= 0 {
			return nil, fmt.Errorf("%w: x[%d]", ErrNotFinite, j)
		}
		return nil, fmt.Errorf("%w: element %d", ErrOverflow, i)
	}
	return dst, nil
}

// MatTranspose returns the cols×rows transpose of the rows×cols matrix a in dst (which must not overlap
// a).
// Errors: ErrDimension when a's length disagrees with its dimensions; ErrNotFinite for a NaN or ±Inf
// element.
//
// ID:      NUM-LINALG-092
// Time:    O(rows·cols)
// Space:   O(1) extra; O(rows·cols) for the result, one allocation only if dst lacks capacity
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §2.7 "Transposes and Permutations"
// Tags:    matrix, n-dimensional, transpose, linear-algebra
// Units:   the unit of a
// Example: [1 2 3; 4 5 6] → [1 4; 2 5; 3 6]
func MatTranspose(dst, a []float64, rows, cols int) ([]float64, error) {
	if err := checkMatrix("a", a, rows, cols); err != nil {
		return nil, err
	}
	if i := firstNotFinite(a); i >= 0 {
		return nil, fmt.Errorf("%w: a[%d]", ErrNotFinite, i)
	}
	dst = grow(dst, rows*cols)
	for i := range rows {
		for j := range cols {
			dst[j*rows+i] = a[i*cols+j]
		}
	}
	return dst, nil
}

// MatTrace returns the sum of the diagonal of the n×n matrix a, summed by Sum2 (as accurate as in twice
// the working precision).
// Errors: ErrDimension when len(a) ≠ n·n; ErrNotFinite for a NaN or ±Inf diagonal element; ErrOverflow
// when the sum is outside the float64 range.
//
// ID:      NUM-LINALG-093
// Time:    O(n)   n = order of the matrix
// Space:   O(1), no allocation
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §6.1 "Introduction to Eigenvalues" (the trace equals the sum of the eigenvalues)
// Source:  Ogita, Rump, Oishi, "Accurate Sum and Dot Product", SIAM Journal on Scientific Computing, 2005, 26(6):1955–1988, doi:10.1137/030601818, Algorithm 4.4 (Sum2)
// Tags:    matrix, n-dimensional, trace, linear-algebra
// Units:   the unit of a
// Example: [1 2 3; 4 5 6; 7 8 9] → 15
func MatTrace(a []float64, n int) (float64, error) {
	if err := checkMatrix("a", a, n, n); err != nil {
		return 0, err
	}
	var p, s float64
	for i := range n {
		x := a[i*n+i]
		if !reduce.AllFinite(x) {
			return 0, fmt.Errorf("%w: a[%d]", ErrNotFinite, i*n+i)
		}
		var q float64
		p, q = reduce.TwoSum(p, x)
		s += q
	}
	r := p + s
	if !reduce.AllFinite(r) {
		return 0, errOverflow("trace")
	}
	return r, nil
}

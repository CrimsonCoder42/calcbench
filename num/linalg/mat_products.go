package linalg

import (
	"fmt"
	"math"
)

// VecOuter returns the outer product u·vᵀ of an m-vector u and an n-vector v, the m×n matrix with element
// (i, j) = uᵢ·vⱼ, in dst (which must not overlap u or v).
// Errors: ErrNotFinite for a NaN or ±Inf element; ErrOverflow when a product is outside the float64 range;
// ErrDimension when m·n overflows.
//
// ID:      NUM-LINALG-120
// Time:    O(m·n)   m = len(u), n = len(v)
// Space:   O(1) extra; O(m·n) for the result, one allocation only if dst lacks capacity
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §2.4 "Rules for Matrix Operations" (a column times a row is a rank-one matrix)
// Tags:    matrix, vector, n-dimensional, outer-product, rank-one, linear-algebra
// Units:   u in unit U, v in unit V; the result in U·V
// Example: u = (1, 2), v = (3, 4, 5) → [3 4 5; 6 8 10] (hand-checked)
func VecOuter(dst, u, v []float64) ([]float64, error) {
	m, n := len(u), len(v)
	if m > 0 && n > math.MaxInt/m {
		return nil, fmt.Errorf("%w: result is %d×%d (too large)", ErrDimension, m, n)
	}
	if i := firstNotFinite(u); i >= 0 {
		return nil, notFiniteAt("u", i)
	}
	if i := firstNotFinite(v); i >= 0 {
		return nil, notFiniteAt("v", i)
	}
	dst = grow(dst, m*n)
	for i, ui := range u {
		row := dst[i*n : i*n+n]
		for j, vj := range v {
			row[j] = ui * vj
		}
	}
	return finishElementwise(dst, "outer product")
}

// MatKronecker returns the Kronecker product A ⊗ B of the m×n matrix a and the p×q matrix b: the mp×nq
// block matrix whose (i, j) block is aᵢⱼ·B, in dst (which must not overlap a or b).
// Errors: ErrDimension when a length disagrees with its dimensions or mp·nq overflows; ErrNotFinite for a
// NaN or ±Inf element; ErrOverflow when a product is outside the float64 range.
//
// ID:      NUM-LINALG-121
// Time:    O(m·n·p·q)
// Space:   O(1) extra; O(m·n·p·q) for the result, one allocation only if dst lacks capacity
// Source:  Van Loan, "The ubiquitous Kronecker product", Journal of Computational and Applied Mathematics, 2000, 123(1–2):85–100, doi:10.1016/S0377-0427(00)00393-9, §1 (B ⊗ C as the block matrix of bᵢⱼC)
// Tags:    matrix, n-dimensional, kronecker, tensor-product, block, linear-algebra
// Units:   a in unit A, b in unit B; the result in A·B
// Example: [1 2; 3 4] ⊗ [0 5; 6 7] → [0 5 0 10; 6 7 12 14; 0 15 0 20; 18 21 24 28] (hand-checked)
func MatKronecker(dst, a []float64, m, n int, b []float64, p, q int) ([]float64, error) {
	if err := checkMatrix("a", a, m, n); err != nil {
		return nil, err
	}
	if err := checkMatrix("b", b, p, q); err != nil {
		return nil, err
	}
	rows, ok1 := mulInt(m, p)
	cols, ok2 := mulInt(n, q)
	size, ok3 := mulInt(rows, cols)
	if !ok1 || !ok2 || !ok3 {
		return nil, fmt.Errorf("%w: result is (%d·%d)×(%d·%d) (too large)", ErrDimension, m, p, n, q)
	}
	if i := firstNotFinite(a); i >= 0 {
		return nil, notFiniteAt("a", i)
	}
	if i := firstNotFinite(b); i >= 0 {
		return nil, notFiniteAt("b", i)
	}
	dst = grow(dst, size)
	for i := range m {
		for j := range n {
			aij := a[i*n+j]
			for r := range p {
				out := dst[(i*p+r)*cols+j*q : (i*p+r)*cols+j*q+q]
				for c, brc := range b[r*q : r*q+q] {
					out[c] = aij * brc
				}
			}
		}
	}
	return finishElementwise(dst, "Kronecker product")
}

// mulInt returns a·b for non-negative a, b and whether it fits in an int.
func mulInt(a, b int) (int, bool) {
	if a == 0 || b == 0 {
		return 0, true
	}
	if a > math.MaxInt/b {
		return 0, false
	}
	return a * b, true
}

// maxMatPower caps the exponent of MatPower at 2⁴⁰ (at most 80 products), a policy: any k beyond it is
// treated as an input error even when Aᵏ is finite (for example a rotation or a projector).
const maxMatPower = 1 << 40

// MatPower returns Aᵏ for an n×n matrix and an integer k ≥ 0 (A⁰ = I) by binary powering (repeated
// squaring), in dst (which must not overlap a).
// Errors: ErrDimension when n < 1 or len(a) ≠ n²; ErrOutOfRange when k < 0 or k > 2⁴⁰ (a policy cap, see
// maxMatPower); ErrNotFinite for a NaN or ±Inf element; ErrOverflow when an element leaves the float64 range.
//
// ID:      NUM-LINALG-122
// Time:    O(n³·log k)   at most 2⌊log₂ k⌋ matrix products
// Space:   O(n²) extra (two work buffers); O(n²) for the result, one allocation only if dst lacks capacity
// Source:  Knuth, The Art of Computer Programming, Vol. 2: Seminumerical Algorithms, 3rd ed., Addison-Wesley, 1997, §4.6.3 "Evaluation of Powers", Algorithm A (right-to-left binary method)
// Source:  Golub, Van Loan, Matrix Computations, 4th ed., Johns Hopkins University Press, 2013, §2.3 "Matrix Norms" (‖AB‖ ≤ ‖A‖‖B‖ bounds the growth of each product)
// Tags:    matrix, n-dimensional, power, repeated-squaring, exponentiation, linear-algebra
// Units:   a in unit A; the result in Aᵏ
// Example: [1 1; 1 0]⁵ → [8 5; 5 3] (Fibonacci numbers; hand-checked)
func MatPower(dst, a []float64, n, k int) ([]float64, error) {
	if n < 1 {
		return nil, fmt.Errorf("%w: n = %d", ErrDimension, n)
	}
	if err := checkMatrix("a", a, n, n); err != nil {
		return nil, err
	}
	if k < 0 || k > maxMatPower {
		return nil, fmt.Errorf("%w: k = %d, want 0 ≤ k ≤ 2^40", ErrOutOfRange, k)
	}
	if i := firstNotFinite(a); i >= 0 {
		return nil, notFiniteAt("a", i)
	}
	dst, _ = MatIdentity(dst, n)
	if k == 0 {
		return dst, nil
	}
	base := make([]float64, 2*n*n)
	z, tmp := base[:n*n], base[n*n:]
	copy(z, a)
	for {
		if k&1 == 1 {
			if _, err := MatMul(tmp, dst, n, n, z, n); err != nil {
				return nil, err
			}
			copy(dst, tmp)
		}
		k >>= 1
		if k == 0 {
			return dst, nil
		}
		if _, err := MatMul(tmp, z, n, n, z, n); err != nil {
			return nil, err
		}
		copy(z, tmp)
	}
}

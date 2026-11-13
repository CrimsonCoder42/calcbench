package linalg

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/algo/dp"
)

// MatMulChain returns the product M₁·M₂·…·Mₖ (numpy's multi_dot) of k row-major matrices, Mᵢ being
// dims[i−1]×dims[i] (len(dims) = k+1), in the parenthesisation with the fewest scalar multiplications,
// found by dp.MatrixChainOrder (ALG-DP-019); each product is MatMul. The result is dims[0]×dims[k], in a
// new slice (a single matrix is copied).
// Errors: ErrDimension when len(mats) < 1, len(dims) ≠ len(mats)+1, a dimension is < 1, or len(Mᵢ) ≠
// dims[i−1]·dims[i]; ErrNotFinite for a NaN or ±Inf entry; ErrOverflow from a product; an error from
// dp.MatrixChainOrder (too long a chain).
//
// ID:      NUM-LINALG-184
// Time:    O(k³ + P)   k = number of matrices, P = the optimal scalar-multiplication count
// Space:   O(k²) for the order tables; O(size of the intermediates) for the products; recursion depth ≤ k
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §15.2 "Matrix-chain multiplication" (MATRIX-CHAIN-ORDER, PRINT-OPTIMAL-PARENS)
// Tags:    matrix, n-dimensional, matrix-chain, multi-dot, product, dynamic-programming, linear-algebra
// Units:   the product of the matrices' units
// Example: 10×100, 100×5, 5×50 → computed as (M₁M₂)M₃, 7500 multiplications instead of 75000 (CLRS §15.2)
func MatMulChain(mats [][]float64, dims []int) ([]float64, error) {
	k := len(mats)
	if k < 1 || len(dims) != k+1 {
		return nil, fmt.Errorf("%w: %d matrices, len(dims) = %d", ErrDimension, k, len(dims))
	}
	for i, d := range dims {
		if d < 1 {
			return nil, fmt.Errorf("%w: dims[%d] = %d", ErrDimension, i, d)
		}
	}
	for i, m := range mats {
		if err := checkMatrix(fmt.Sprintf("mats[%d]", i), m, dims[i], dims[i+1]); err != nil {
			return nil, err
		}
		if j := firstNotFinite(m); j >= 0 {
			return nil, fmt.Errorf("%w: mats[%d][%d]", ErrNotFinite, i, j)
		}
	}
	if k == 1 {
		return append([]float64(nil), mats[0]...), nil
	}
	_, parens, err := dp.MatrixChainOrder(dims)
	if err != nil {
		return nil, err
	}
	pos := 0
	prod, _, _, err := chainEval(parens, &pos, mats, dims)
	return prod, err
}

// chainEval evaluates the parenthesisation string at *pos ("A3" or "(XY)") and returns the product with
// its rows and columns. Leaves return the input matrix itself (read only).
func chainEval(s string, pos *int, mats [][]float64, dims []int) ([]float64, int, int, error) {
	if s[*pos] == 'A' {
		*pos++
		idx := 0
		for *pos < len(s) && s[*pos] >= '0' && s[*pos] <= '9' {
			idx = idx*10 + int(s[*pos]-'0')
			*pos++
		}
		return mats[idx-1], dims[idx-1], dims[idx], nil
	}
	*pos++ // '('
	left, r, inner, err := chainEval(s, pos, mats, dims)
	if err != nil {
		return nil, 0, 0, err
	}
	right, _, c, err := chainEval(s, pos, mats, dims)
	if err != nil {
		return nil, 0, 0, err
	}
	*pos++ // ')'
	p, err := MatMul(nil, left, r, inner, right, c)
	return p, r, c, err
}

package linalg

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// MatConditionEstimate1 returns an estimate of the 1-norm condition number κ₁(A) = ‖A‖₁·‖A⁻¹‖₁ of the
// n×n matrix a, without forming A⁻¹. ‖A⁻¹‖₁ is estimated by Hager's method with Higham's refinements: at
// most five iterations, each one solve with A and one with Aᵀ from a single LU factorisation, stopping
// when the estimate stops growing or the sign vector repeats, then Higham's alternating test vector
// bᵢ = (−1)ⁱ(1 + i/(n−1)). The estimate never exceeds the true κ₁ and is usually within a factor of 3 of
// it. Roughly log₁₀ κ₁ decimal digits are lost when solving Ax = b or inverting A; κ₁ near 1/ε (about
// 4.5e15) means the matrix is numerically singular.
// Errors: ErrSingular when a is exactly singular (κ₁ = ∞); ErrOverflow when the estimate exceeds the
// float64 range; otherwise as LUDecompose.
//
// ID:      NUM-LINALG-100
// Time:    O(n³)   n = order of the matrix (the LU factorisation; each iteration is O(n²))
// Space:   O(n²) for the factor and O(n) work
// Source:  Hager, "Condition Estimates", SIAM Journal on Scientific and Statistical Computing, 1984, 5(2):311–316, doi:10.1137/0905023
// Source:  Higham, "FORTRAN codes for estimating the one-norm of a real or complex matrix, with applications to condition estimation", ACM Transactions on Mathematical Software, 1988, 14(4):381–396, doi:10.1145/50063.214386
// Tags:    matrix, n-dimensional, condition-number, estimate, lu, stability, linear-algebra
// Units:   dimensionless
// Example: diag(1, 1e−8) → 1e8 (exact for a diagonal matrix; hand-checked)
func MatConditionEstimate1(a []float64, n int) (float64, error) {
	f, err := LUDecompose(a, n)
	if err != nil {
		return 0, err
	}
	if f.singular {
		return 0, fmt.Errorf("%w: zero pivot", ErrSingular)
	}
	normA := 0.0
	for j := range n {
		var p, s float64
		for i := range n {
			var q float64
			p, q = reduce.TwoSum(p, math.Abs(a[i*n+j]))
			s += q
		}
		normA = max(normA, p+s)
	}
	buf := make([]float64, 5*n)
	x, y, xi, z, tmp := buf[:n], buf[n:2*n], buf[2*n:3*n], buf[3*n:4*n], buf[4*n:]
	for i := range x {
		x[i] = 1 / float64(n)
	}
	est := 0.0
	for iter := range 5 {
		f.solvePermuted(y, x) // y = A⁻¹x
		ny := norm1(y)
		if iter > 0 && ny <= est {
			break
		}
		est = ny
		repeated := iter > 0
		for i, v := range y {
			s := 1.0
			if v < 0 {
				s = -1
			}
			if xi[i] != s {
				repeated = false
			}
			xi[i] = s
		}
		if repeated {
			break
		}
		copy(z, xi)
		f.solveTransposeInPlace(z, tmp) // z = A⁻ᵀξ
		j, zmax, ztx := 0, -1.0, 0.0
		for i, v := range z {
			if av := math.Abs(v); av > zmax {
				j, zmax = i, av
			}
			ztx = math.FMA(v, x[i], ztx)
		}
		if iter > 0 && zmax <= ztx {
			break
		}
		clear(x)
		x[j] = 1
	}
	for i := range n {
		v := 1.0
		if n > 1 {
			v += float64(i) / float64(n-1)
		}
		if i%2 == 1 {
			v = -v
		}
		x[i] = v
	}
	f.solvePermuted(y, x)
	est = max(est, 2*norm1(y)/(3*float64(n)))
	k := normA * est
	if !reduce.AllFinite(k) {
		return 0, errOverflow("κ₁")
	}
	return k, nil
}

// solvePermuted writes A⁻¹b into y (y must not overlap b).
func (f *LU) solvePermuted(y, b []float64) {
	for i := range f.n {
		y[i] = b[f.piv[i]]
	}
	f.solveInPlace(y)
}

// norm1 returns Σ|vᵢ| by Sum2.
func norm1(v []float64) float64 {
	var p, s float64
	for _, x := range v {
		var q float64
		p, q = reduce.TwoSum(p, math.Abs(x))
		s += q
	}
	return p + s
}

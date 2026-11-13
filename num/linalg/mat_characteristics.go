package linalg

import (
	"errors"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
	"github.com/CrimsonCoder42/calcbench/num/algebra"
)

// Mat2Info bundles the characteristics of a 2×2 matrix: determinant, trace, numerical rank, the inverse
// when it exists (Invertible), and the eigenvalues (real or a complex pair, sorted by real part).
type Mat2Info struct {
	Det, Trace  float64
	Rank        int
	Invertible  bool
	Inverse     Mat2
	Eigenvalues [2]complex128
}

// Mat3Info bundles the characteristics of a 3×3 matrix: determinant, trace, numerical rank, the inverse
// when it exists (Invertible), the characteristic polynomial (ascending), and the eigenvalues (sorted by
// real part, then imaginary part).
type Mat3Info struct {
	Det, Trace  float64
	Rank        int
	Invertible  bool
	Inverse     Mat3
	CharPoly    [4]float64
	Eigenvalues [3]complex128
}

// Mat2Characteristics returns the determinant (Mat2Det), trace (Mat2Trace), rank (as MatRank: complete pivoting, tol < 0
// selects max(n)·ε·max|aᵢⱼ|), inverse when it exists (Mat2Inverse) and eigenvalues (Mat2Eigen) of a 2×2 matrix in
// one result, built from the existing functions.
// Errors: ErrNotFinite for a NaN or ±Inf element or tol; ErrOverflow when a quantity is outside the float64
// range.
//
// ID:      NUM-LINALG-134
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §5.3 "Cramer's Rule, Inverses, and Volumes" (determinant and inverse) and §6.1 "Introduction to Eigenvalues" (trace, determinant and eigenvalues of a 2×2)
// Tags:    matrix, 2x2, characteristics, determinant, trace, rank, inverse, eigenvalues, linear-algebra
// Units:   m in unit A; Det in A², Trace and eigenvalues in A, Inverse in 1/A
// Example: [4 1; 2 3] → det 10, trace 7, rank 2, inverse [0.3 −0.1; −0.2 0.4], eigenvalues 2, 5 (hand-checked)
func Mat2Characteristics(m Mat2, tol float64) (Mat2Info, error) {
	det, err := Mat2Det(m)
	if err != nil {
		return Mat2Info{}, err
	}
	tr, err := Mat2Trace(m)
	if err != nil {
		return Mat2Info{}, err
	}
	rank, err := rankFixed(m[:], 2, tol)
	if err != nil {
		return Mat2Info{}, err
	}
	eig, err := Mat2Eigen(m)
	if err != nil {
		return Mat2Info{}, err
	}
	info := Mat2Info{Det: det, Trace: tr, Rank: rank, Eigenvalues: eig.Values}
	inv, err := Mat2Inverse(m)
	switch {
	case err == nil:
		info.Invertible, info.Inverse = true, inv
	case !errors.Is(err, ErrSingular):
		return Mat2Info{}, err
	}
	return info, nil
}

// Mat3Characteristics returns the determinant (Mat3Det), trace (Mat3Trace), rank (as MatRank: complete pivoting, tol < 0
// selects max(n)·ε·max|aᵢⱼ|), inverse when it exists (Mat3Inverse), characteristic polynomial (Mat3CharPoly) and
// eigenvalues (the roots of that polynomial by CubicRoots, NUM-ALGEB-009, each accurate to its own
// conditioning) of a 3×3 matrix in one result. For a symmetric matrix Mat3SymmetricEigenvalues is the more
// accurate eigenvalue source.
// Errors: ErrNotFinite for a NaN or ±Inf element or tol; ErrOverflow when a quantity is outside the float64
// range.
//
// ID:      NUM-LINALG-135
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §5.3 "Cramer's Rule, Inverses, and Volumes" and §6.1 "Introduction to Eigenvalues" (eigenvalues are the roots of det(A − λI) = 0; Mat3CharPoly returns the monic det(λI − A) = −det(A − λI), same roots)
// Tags:    matrix, 3x3, characteristics, determinant, trace, rank, inverse, eigenvalues, linear-algebra
// Units:   m in unit A; Det in A³, Trace and eigenvalues in A, Inverse in 1/A
// Example: diag(1, 2, 3) → det 6, trace 6, rank 3, inverse diag(1, 1/2, 1/3), eigenvalues 1, 2, 3 (hand-checked)
func Mat3Characteristics(m Mat3, tol float64) (Mat3Info, error) {
	det, err := Mat3Det(m)
	if err != nil {
		return Mat3Info{}, err
	}
	tr, err := Mat3Trace(m)
	if err != nil {
		return Mat3Info{}, err
	}
	rank, err := rankFixed(m[:], 3, tol)
	if err != nil {
		return Mat3Info{}, err
	}
	p, err := Mat3CharPoly(m)
	if err != nil {
		return Mat3Info{}, err
	}
	roots, err := algebra.CubicRoots(p[3], p[2], p[1], p[0])
	if err != nil {
		return Mat3Info{}, errOverflow("eigenvalue")
	}
	info := Mat3Info{Det: det, Trace: tr, Rank: rank, CharPoly: p, Eigenvalues: roots}
	inv, err := Mat3Inverse(m)
	switch {
	case err == nil:
		info.Invertible, info.Inverse = true, inv
	case !errors.Is(err, ErrSingular):
		return Mat3Info{}, err
	}
	return info, nil
}

// rankFixed is MatRank (complete pivoting on a power-of-two-scaled copy, the same default tolerance) for an
// n×n matrix with n ≤ 3, on a stack array so the fixed-size functions allocate nothing.
func rankFixed(a []float64, n int, tol float64) (int, error) {
	if !reduce.AllFinite(tol) {
		return 0, errNotFinite("tol")
	}
	mx, _ := maxAbs(a)
	if mx == 0 {
		return 0, nil
	}
	if tol < 0 {
		tol = float64(n) * 0x1p-52 * mx
	}
	_, e := pow2Scale(mx)
	var w [9]float64
	for i, x := range a {
		w[i] = math.Ldexp(x, -e)
	}
	tol = math.Ldexp(tol, -e)
	rank := 0
	for k := range n {
		p, q, best := k, k, -1.0
		for i := k; i < n; i++ {
			for j := k; j < n; j++ {
				if v := math.Abs(w[i*n+j]); v > best {
					p, q, best = i, j, v
				}
			}
		}
		if best <= tol {
			break
		}
		for j := range n {
			w[k*n+j], w[p*n+j] = w[p*n+j], w[k*n+j]
		}
		for i := range n {
			w[i*n+k], w[i*n+q] = w[i*n+q], w[i*n+k]
		}
		rank++
		for i := k + 1; i < n; i++ {
			l := w[i*n+k] / w[k*n+k]
			for j := k + 1; j < n; j++ {
				w[i*n+j] = math.FMA(-l, w[k*n+j], w[i*n+j])
			}
			w[i*n+k] = 0
		}
	}
	return rank, nil
}

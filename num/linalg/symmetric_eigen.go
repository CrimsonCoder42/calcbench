package linalg

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// Mat2SymmetricEigenvalues returns the two eigenvalues lo ≤ hi of the symmetric matrix [a b; b c]. They
// are the roots of λ² − (a + c)λ + (ac − b²) = 0; the one of larger magnitude is
// (a + c)/2 ± hypot((a − c)/2, b) with the sign of (a + c), and the other is det/that, so neither
// suffers cancellation (the determinant itself is computed within 2 ulps by Kahan's algorithm).
// Errors: ErrNotSymmetric when m₀₁ ≠ m₁₀ (compared exactly); ErrNotFinite for a NaN or ±Inf element;
// ErrOverflow when an eigenvalue is outside the float64 range.
//
// ID:      NUM-LINALG-064
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §6.1 "Introduction to Eigenvalues" (λ₁ + λ₂ = trace, λ₁λ₂ = det) and §6.4 "Symmetric Matrices" (real eigenvalues)
// Source:  Golub, Van Loan, Matrix Computations, 4th ed., Johns Hopkins University Press, 2013, §8.5 "Jacobi Methods" (the 2-by-2 symmetric Schur decomposition)
// Tags:    matrix, 2x2, eigenvalues, symmetric, closed-form, linear-algebra
// Units:   m in any unit; the eigenvalues in that unit
// Example: [2 1; 1 2] → (1, 3) (hand-checked)
func Mat2SymmetricEigenvalues(m Mat2) (lo, hi float64, err error) {
	if !reduce.AllFinite(m[:]...) {
		return 0, 0, errNotFinite("m")
	}
	if m[1] != m[2] {
		return 0, 0, fmt.Errorf("%w: m01 = %v, m10 = %v", ErrNotSymmetric, m[1], m[2])
	}
	// Scale by an exact power of two (eigenvalues scale with the matrix) so that neither a/2 underflows
	// nor the squares overflow; the eigenvalues are scaled back at the end.
	mx := max(math.Abs(m[0]), math.Abs(m[1]), math.Abs(m[3]))
	if mx == 0 {
		return 0, 0, nil
	}
	_, e := pow2Scale(mx)
	a, b, c := math.Ldexp(m[0], -e), math.Ldexp(m[1], -e), math.Ldexp(m[3], -e)
	t := a/2 + c/2
	r := math.Hypot(a/2-c/2, b)
	big := t + math.Copysign(r, t)
	var other float64
	if big != 0 {
		other = det2(a, b, b, c) / big
	}
	lo, hi = math.Ldexp(min(big, other), e), math.Ldexp(max(big, other), e)
	if !reduce.AllFinite(lo) || !reduce.AllFinite(hi) {
		return 0, 0, errOverflow("eigenvalue")
	}
	return lo, hi, nil
}

// Mat3SymmetricEigenvalues returns the three eigenvalues of a symmetric 3×3 matrix in ascending order,
// by Smith's closed form: with m = trace/3, K = A − mI, q = det(K)/2 and p = (Σ Kᵢⱼ²)/6, the angle
// φ = ⅓·atan2(√(p³ − q²), q) gives the eigenvalues m + 2√p cos φ and m − √p(cos φ ± √3 sin φ). The
// matrix is first scaled by an exact power of two (eigenvalues scale with it) so no square overflows.
// For distinct eigenvalues the absolute error is a few ulps of the largest |eigenvalue|. For two equal
// or nearly equal eigenvalues p³ − q² nearly vanishes and its rounding error passes through a square
// root, so those two carry an absolute error up to about √ε·max|λ| (≈ 1.5e−8 relative); an eigenvalue
// much smaller than the largest also loses relative accuracy. An iterative symmetric eigensolver (LATER)
// does not have these limits.
// Errors: ErrNotSymmetric unless m₀₁ = m₁₀, m₀₂ = m₂₀ and m₁₂ = m₂₁ exactly; ErrNotFinite for a NaN or
// ±Inf element; ErrOverflow when an eigenvalue is outside the float64 range.
//
// ID:      NUM-LINALG-076
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Smith, "Eigenvalues of a symmetric 3 × 3 matrix", Communications of the ACM, 1961, 4(4):168, doi:10.1145/355578.366316
// Tags:    matrix, 3x3, eigenvalues, symmetric, closed-form, principal-moments, linear-algebra
// Units:   m in any unit; the eigenvalues in that unit
// Example: [2 0 0; 0 3 4; 0 4 9] → (1, 2, 11) (hand-checked: the lower block has trace 12 and det 11)
func Mat3SymmetricEigenvalues(m Mat3) ([3]float64, error) {
	if !reduce.AllFinite(m[:]...) {
		return [3]float64{}, errNotFinite("m")
	}
	if m[1] != m[3] || m[2] != m[6] || m[5] != m[7] {
		return [3]float64{}, fmt.Errorf("%w: m", ErrNotSymmetric)
	}
	mx := 0.0
	for _, x := range m {
		mx = max(mx, math.Abs(x))
	}
	if mx == 0 {
		return [3]float64{}, nil
	}
	sc, e := pow2Scale(mx)
	var a Mat3
	for i, x := range m {
		a[i] = x * sc
	}
	var ev [3]float64
	off := a[1]*a[1] + a[2]*a[2] + a[5]*a[5]
	if off == 0 {
		ev = [3]float64{a[0], a[4], a[8]}
	} else {
		mean := (a[0] + a[4] + a[8]) / 3
		k := a
		k[0] -= mean
		k[4] -= mean
		k[8] -= mean
		p := (k[0]*k[0] + k[4]*k[4] + k[8]*k[8] + 2*off) / 6
		c := cross3(row3(&k, 1), row3(&k, 2))
		q := dot3(row3(&k, 0), c) / 2
		phi := math.Atan2(math.Sqrt(max(p*p*p-q*q, 0)), q) / 3
		sp := math.Sqrt(p)
		sphi, cphi := math.Sincos(phi)
		ev = [3]float64{
			mean + 2*sp*cphi,
			mean - sp*(cphi+math.Sqrt(3)*sphi),
			mean - sp*(cphi-math.Sqrt(3)*sphi),
		}
	}
	// Sort three values ascending.
	if ev[0] > ev[1] {
		ev[0], ev[1] = ev[1], ev[0]
	}
	if ev[1] > ev[2] {
		ev[1], ev[2] = ev[2], ev[1]
	}
	if ev[0] > ev[1] {
		ev[0], ev[1] = ev[1], ev[0]
	}
	for i := range ev {
		ev[i] = math.Ldexp(ev[i], e)
		if !reduce.AllFinite(ev[i]) {
			return [3]float64{}, errOverflow("eigenvalue")
		}
	}
	return ev, nil
}

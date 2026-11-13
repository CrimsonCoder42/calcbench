package linalg

import (
	"fmt"
	"math"
)

// maxFrancisIterations is the number of double-shift QR steps allowed per eigenvalue before
// MatEigenvalues reports ErrNoConvergence (LAPACK's dhseqr allows 30; doubled here, with exceptional
// shifts every 10 steps).
const maxFrancisIterations = 60

// MatEigenvalues returns the n eigenvalues of a general real n×n matrix (symmetric or not) as complex
// numbers, in valuesDst. The matrix is scaled by an exact power of two, balanced (Parlett–Reinsch, radix
// 2, exact), reduced to upper Hessenberg form by Householder reflections, and the Hessenberg matrix is
// driven to real Schur form by the Francis implicit double-shift QR iteration, deflating 1×1 and 2×2
// blocks. A complex-conjugate pair is returned in adjacent positions, the positive imaginary part
// first; the order otherwise follows deflation (no sort). Eigenvectors are not formed (use
// MatSymmetricEigenJacobi for a symmetric matrix). a is not modified.
// Errors: ErrDimension when n < 1 or len(a) ≠ n²; ErrNotFinite for a NaN or ±Inf element;
// ErrNoConvergence when one eigenvalue needs more than 60 QR steps.
//
// ID:      NUM-LINALG-169
// Time:    O(n³)   n = order; about 10n³ flops for the reduction and 2–3 QR steps per eigenvalue of O(n²) each
// Space:   O(n²) for the working copy; O(n) for the result, one allocation only if valuesDst lacks capacity
// Source:  Golub, Van Loan, Matrix Computations, 4th ed., Johns Hopkins University Press, 2013, §7.4 "The Hessenberg and Real Schur Forms" (Algorithm 7.4.2, Householder reduction) and §7.5 "The Practical QR Algorithm" (Algorithm 7.5.1, Francis QR step; §7.5.7 balancing)
// Source:  Martin, Peters, Wilkinson, "The QR Algorithm for Real Hessenberg Matrices", Numerische Mathematik, 1970, 14:219–231, doi:10.1007/BF02163331 (deflation tests and exceptional shifts)
// Tags:    matrix, n-dimensional, eigenvalues, nonsymmetric, hessenberg, francis, qr-algorithm, schur, linear-algebra
// Units:   a in unit A; eigenvalues in A
// Example: [0 −1; 1 0] (rotation by 90°) → ±i; [2 0 0; 0 3 4; 0 4 9] → 1, 2, 11 (hand-checked)
func MatEigenvalues(valuesDst []complex128, a []float64, n int) ([]complex128, error) {
	if n < 1 {
		return nil, fmt.Errorf("%w: n = %d", ErrDimension, n)
	}
	if err := checkMatrix("a", a, n, n); err != nil {
		return nil, err
	}
	if i := firstNotFinite(a); i >= 0 {
		return nil, notFiniteAt("a", i)
	}
	big, _ := maxAbs(a)
	_, e := pow2Scale(max(big, 0x1p-1000))
	h := make([]float64, n*n+2*n)
	wr, wi := h[n*n:n*n+n], h[n*n+n:]
	h = h[:n*n]
	for i, x := range a {
		h[i] = math.Ldexp(x, -e)
	}
	balance(h, n)
	hessenberg(h, nil, n)
	if err := francisQR(h, n, wr, wi); err != nil {
		return nil, err
	}
	if cap(valuesDst) >= n {
		valuesDst = valuesDst[:n]
	} else {
		valuesDst = make([]complex128, n)
	}
	for i := range n {
		valuesDst[i] = complex(math.Ldexp(wr[i], e), math.Ldexp(wi[i], e))
	}
	return valuesDst, nil
}

// balance applies the Parlett–Reinsch diagonal similarity D⁻¹AD with D a power of two (exact), so that
// row and column norms are close; it leaves the eigenvalues unchanged (GvL §7.5.7).
func balance(h []float64, n int) {
	for pass := 0; pass < 100; pass++ {
		done := true
		for i := range n {
			c, r := 0.0, 0.0
			for j := range n {
				if j != i {
					c += math.Abs(h[j*n+i])
					r += math.Abs(h[i*n+j])
				}
			}
			if c == 0 || r == 0 {
				continue
			}
			s := c + r
			f := 1.0
			for g := r / 2; c < g; c *= 4 {
				f *= 2
			}
			for g := r * 2; c > g; c /= 4 {
				f /= 2
			}
			if (c+r)/f < 0.95*s {
				done = false
				g := 1 / f
				for j := range n {
					h[i*n+j] *= g
				}
				for j := range n {
					h[j*n+i] *= f
				}
			}
		}
		if done {
			return
		}
	}
}

// hessenberg reduces h in place to upper Hessenberg form by Householder similarity transforms
// (GvL Algorithm 7.4.2); entries below the first subdiagonal are set to zero. When q is non-nil (n×n,
// holding the identity on entry) it is overwritten by the orthogonal Q with A = Q·H·Qᵀ (GvL §7.4.3).
func hessenberg(h, q []float64, n int) {
	v := make([]float64, n)
	for k := 0; k+2 < n; k++ {
		scale := 0.0
		for i := k + 1; i < n; i++ {
			scale = max(scale, math.Abs(h[i*n+k]))
		}
		if scale == 0 {
			continue
		}
		sigma := 0.0
		for i := k + 1; i < n; i++ {
			v[i] = h[i*n+k] / scale
			sigma += v[i] * v[i]
		}
		alpha := math.Copysign(math.Sqrt(sigma), v[k+1])
		v[k+1] += alpha
		vtv := 2 * alpha * v[k+1] // ‖v‖² = σ + 2α·x₁ + α² = 2α(x₁ + α)
		if vtv == 0 {
			continue
		}
		beta := 2 / vtv
		// H·A: rows k+1…n−1.
		for j := k; j < n; j++ {
			s := 0.0
			for i := k + 1; i < n; i++ {
				s += v[i] * h[i*n+j]
			}
			s *= beta
			for i := k + 1; i < n; i++ {
				h[i*n+j] -= s * v[i]
			}
		}
		// A·H: columns k+1…n−1.
		for i := range n {
			s := 0.0
			for j := k + 1; j < n; j++ {
				s += h[i*n+j] * v[j]
			}
			s *= beta
			for j := k + 1; j < n; j++ {
				h[i*n+j] -= s * v[j]
			}
		}
		if q != nil { // Q ← Q·Hₖ: columns k+1…n−1.
			for i := range n {
				s := 0.0
				for j := k + 1; j < n; j++ {
					s += q[i*n+j] * v[j]
				}
				s *= beta
				for j := k + 1; j < n; j++ {
					q[i*n+j] -= s * v[j]
				}
			}
		}
		for i := k + 2; i < n; i++ {
			h[i*n+k] = 0
		}
	}
}

// francisQR computes the eigenvalues of the upper Hessenberg matrix h (destroyed) by the implicit
// double-shift QR iteration with deflation (GvL Algorithm 7.5.1; Martin, Peters, Wilkinson 1970).
func francisQR(h []float64, n int, wr, wi []float64) error {
	const eps = 0x1p-52
	at := func(i, j int) *float64 { return &h[i*n+j] }
	anorm := 0.0
	for i := range n {
		for j := max(i-1, 0); j < n; j++ {
			anorm += math.Abs(h[i*n+j])
		}
	}
	t := 0.0
	nn := n - 1
	for nn >= 0 {
		its := 0
		for {
			l := nn
			for ; l >= 1; l-- {
				s := math.Abs(*at(l-1, l-1)) + math.Abs(*at(l, l))
				if s == 0 {
					s = anorm
				}
				if math.Abs(*at(l, l-1)) <= eps*s {
					*at(l, l-1) = 0
					break
				}
			}
			x := *at(nn, nn)
			if l == nn { // one root found
				wr[nn], wi[nn] = x+t, 0
				nn--
				break
			}
			y := *at(nn-1, nn-1)
			w := *at(nn, nn-1) * *at(nn-1, nn)
			if l == nn-1 { // a 2×2 block deflates
				p := 0.5 * (y - x)
				q := p*p + w
				z := math.Sqrt(math.Abs(q))
				x += t
				if q >= 0 {
					z = p + math.Copysign(z, p)
					wr[nn-1], wr[nn] = x+z, x+z
					if z != 0 {
						wr[nn] = x - w/z
					}
					wi[nn-1], wi[nn] = 0, 0
				} else {
					wr[nn-1], wr[nn] = x+p, x+p
					wi[nn-1], wi[nn] = z, -z
				}
				nn -= 2
				break
			}
			if its == maxFrancisIterations {
				return fmt.Errorf("%w: eigenvalue %d after %d QR steps", ErrNoConvergence, nn, its)
			}
			if its > 0 && its%10 == 0 { // exceptional shift
				t += x
				for i := 0; i <= nn; i++ {
					*at(i, i) -= x
				}
				s := math.Abs(*at(nn, nn-1)) + math.Abs(*at(nn-1, nn-2))
				x = 0.75 * s
				y = x
				w = -0.4375 * s * s
			}
			its++
			// Find two consecutive small subdiagonal elements and form the first column of the shift.
			var m int
			var p, q, r float64
			for m = nn - 2; m >= l; m-- {
				z := *at(m, m)
				r = x - z
				s := y - z
				p = (r*s-w) / *at(m+1, m) + *at(m, m+1)
				q = *at(m+1, m+1) - z - r - s
				r = *at(m+2, m+1)
				s = math.Abs(p) + math.Abs(q) + math.Abs(r)
				p, q, r = p/s, q/s, r/s
				if m == l {
					break
				}
				u := math.Abs(*at(m, m-1)) * (math.Abs(q) + math.Abs(r))
				v := math.Abs(p) * (math.Abs(*at(m-1, m-1)) + math.Abs(z) + math.Abs(*at(m+1, m+1)))
				if u <= eps*v {
					break
				}
			}
			for i := m + 2; i <= nn; i++ {
				*at(i, i-2) = 0
				if i != m+2 {
					*at(i, i-3) = 0
				}
			}
			// Double QR step on rows l…nn and columns m…nn (bulge chase by 3×3 reflectors).
			for k := m; k <= nn-1; k++ {
				if k != m {
					p = *at(k, k-1)
					q = *at(k+1, k-1)
					r = 0
					if k != nn-1 {
						r = *at(k+2, k-1)
					}
					x = math.Abs(p) + math.Abs(q) + math.Abs(r)
					if x != 0 {
						p, q, r = p/x, q/x, r/x
					}
				}
				s := math.Copysign(math.Sqrt(p*p+q*q+r*r), p)
				if s == 0 {
					continue
				}
				if k == m {
					if l != m {
						*at(k, k-1) = -*at(k, k-1)
					}
				} else {
					*at(k, k-1) = -s * x
				}
				p += s
				x = p / s
				y = q / s
				z := r / s
				q /= p
				r /= p
				for j := k; j <= nn; j++ {
					pp := *at(k, j) + q**at(k+1, j)
					if k != nn-1 {
						pp += r * *at(k+2, j)
						*at(k+2, j) -= pp * z
					}
					*at(k+1, j) -= pp * y
					*at(k, j) -= pp * x
				}
				top := min(nn, k+3)
				for i := l; i <= top; i++ {
					pp := x**at(i, k) + y**at(i, k+1)
					if k != nn-1 {
						pp += z * *at(i, k+2)
						*at(i, k+2) -= pp * r
					}
					*at(i, k+1) -= pp * q
					*at(i, k) -= pp
				}
			}
		}
	}
	return nil
}

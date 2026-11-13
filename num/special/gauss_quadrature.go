package special

import (
	"fmt"
	"math"
	"slices"
)

// GaussLegendre returns the n nodes (the roots of P_n, ascending) and weights of the n-point
// Gauss–Legendre rule ∫₋₁¹ f(x) dx ≈ Σ wᵢ f(xᵢ), exact for polynomials of degree ≤ 2n−1.
// Nodes are the eigenvalues of the symmetric tridiagonal Jacobi matrix (Golub & Welsch), each polished
// by Newton steps on the orthonormal recurrence; weights are the Christoffel numbers 1/Σₖ qₖ(xᵢ)².
//
// ID:      NUM-SPEC-015
// Time:    O(n²)   n = number of points (implicit QL on the tridiagonal matrix, then O(n) per node)
// Space:   O(n) for the two results and one scratch slice
// Source:  Golub & Welsch, "Calculation of Gauss quadrature rules", Math. Comp. 23 (1969) 221–230, doi:10.1090/S0025-5718-69-99647-1
// Source:  NIST DLMF, Release 1.2.1, §3.5(v) (Gauss quadrature, eq. 3.5.15 and Christoffel numbers), §18.9 (Table 18.9.1 recurrence)
// Tags:    special-function, quadrature, gauss-legendre, orthogonal-polynomial, roots, integration
// Units:   dimensionless
// Range:   n ≥ 1. Errors: ErrOutOfRange (n < 1), ErrNoConvergence (QL iteration exceeded 60 sweeps on one eigenvalue).
// Example: n = 3: nodes −√(3/5), 0, √(3/5); weights 5/9, 8/9, 5/9 (A&S Table 25.4)
func GaussLegendre(n int) (nodes, weights []float64, err error) {
	return gaussRule(n, 2, func(k int) float64 { return 0 }, func(k int) float64 {
		fk := float64(k)
		return fk / math.Sqrt(4*fk*fk-1)
	})
}

// GaussHermite returns the nodes (roots of H_n, ascending) and weights of the n-point Gauss–Hermite
// rule ∫ e^(−x²) f(x) dx ≈ Σ wᵢ f(xᵢ) over the real line, by the Golub–Welsch method.
//
// ID:      NUM-SPEC-016
// Time:    O(n²)   n = number of points
// Space:   O(n)
// Source:  Golub & Welsch, Math. Comp. 23 (1969) 221–230, doi:10.1090/S0025-5718-69-99647-1
// Source:  NIST DLMF, Release 1.2.1, §3.5(v), §18.9 (Table 18.9.1); A&S Table 25.10 (test values)
// Tags:    special-function, quadrature, gauss-hermite, orthogonal-polynomial, roots, integration
// Units:   dimensionless
// Range:   n ≥ 1. Weights below about 1e-308 underflow to 0. Errors: ErrOutOfRange, ErrNoConvergence.
// Example: n = 2: nodes ∓1/√2, weights √π/2 each
func GaussHermite(n int) (nodes, weights []float64, err error) {
	return gaussRule(n, math.Sqrt(math.Pi), func(k int) float64 { return 0 }, func(k int) float64 {
		return math.Sqrt(float64(k) / 2)
	})
}

// GaussHermiteE returns the nodes (roots of He_n, ascending) and weights of the n-point probabilists'
// Gauss–Hermite rule ∫ e^(−x²/2) f(x) dx ≈ Σ wᵢ f(xᵢ), by the Golub–Welsch method.
//
// ID:      NUM-SPEC-017
// Time:    O(n²)   n = number of points
// Space:   O(n)
// Source:  Golub & Welsch, Math. Comp. 23 (1969) 221–230, doi:10.1090/S0025-5718-69-99647-1
// Source:  NIST DLMF, Release 1.2.1, §3.5(v), §18.9 (Table 18.9.1)
// Tags:    special-function, quadrature, gauss-hermite, orthogonal-polynomial, roots, probability
// Units:   dimensionless
// Range:   n ≥ 1. Errors: ErrOutOfRange, ErrNoConvergence.
// Example: n = 2: nodes ∓1, weights √(2π)/2 each
func GaussHermiteE(n int) (nodes, weights []float64, err error) {
	return gaussRule(n, math.Sqrt(2*math.Pi), func(k int) float64 { return 0 }, func(k int) float64 {
		return math.Sqrt(float64(k))
	})
}

// GaussLaguerre returns the nodes (roots of L_n, ascending) and weights of the n-point Gauss–Laguerre
// rule ∫₀^∞ e^(−x) f(x) dx ≈ Σ wᵢ f(xᵢ), by the Golub–Welsch method.
//
// ID:      NUM-SPEC-018
// Time:    O(n²)   n = number of points
// Space:   O(n)
// Source:  Golub & Welsch, Math. Comp. 23 (1969) 221–230, doi:10.1090/S0025-5718-69-99647-1
// Source:  NIST DLMF, Release 1.2.1, §3.5(v), §18.9 (Table 18.9.1, α = 0); A&S Table 25.9 (test values)
// Tags:    special-function, quadrature, gauss-laguerre, orthogonal-polynomial, roots, integration
// Units:   dimensionless
// Range:   n ≥ 1. Weights below about 1e-308 underflow to 0. Errors: ErrOutOfRange, ErrNoConvergence.
// Example: n = 2: nodes 2 ∓ √2, weights (2 ± √2)/4
func GaussLaguerre(n int) (nodes, weights []float64, err error) {
	return gaussRule(n, 1, func(k int) float64 { return float64(2*k + 1) }, func(k int) float64 {
		return float64(k)
	})
}

// GaussChebyshevT returns the nodes (roots of T_n, ascending) and weights of the n-point Gauss–Chebyshev
// rule of the first kind ∫₋₁¹ f(x)/√(1−x²) dx ≈ Σ wᵢ f(xᵢ), in closed form:
// xᵢ = cos((2i−1)π/(2n)), wᵢ = π/n.
//
// ID:      NUM-SPEC-019
// Time:    O(n)   n = number of points
// Space:   O(n) for the results
// Source:  Abramowitz & Stegun, Handbook of Mathematical Functions, 1964, eq. 25.4.38
// Source:  NIST DLMF, Release 1.2.1, §18.5 (eq. 18.5.1: T_n(cos θ) = cos nθ, whose zeros these are)
// Tags:    special-function, quadrature, gauss-chebyshev, orthogonal-polynomial, roots
// Units:   dimensionless
// Range:   n ≥ 1. Errors: ErrOutOfRange.
// Example: n = 2: nodes ∓1/√2, weights π/2 each
func GaussChebyshevT(n int) (nodes, weights []float64, err error) {
	if n < 1 {
		return nil, nil, fmt.Errorf("%w: n = %d, want ≥ 1", ErrOutOfRange, n)
	}
	nodes = make([]float64, n)
	weights = make([]float64, n)
	for i := 0; i < n; i++ {
		// ascending: i-th node uses k = n − i
		k := n - i
		nodes[i] = math.Cos(float64(2*k-1) * math.Pi / float64(2*n))
		weights[i] = math.Pi / float64(n)
	}
	return nodes, weights, nil
}

// GaussChebyshevU returns the nodes (roots of U_n, ascending) and weights of the n-point Gauss–Chebyshev
// rule of the second kind ∫₋₁¹ f(x)√(1−x²) dx ≈ Σ wᵢ f(xᵢ), in closed form:
// xᵢ = cos(iπ/(n+1)), wᵢ = π/(n+1)·sin²(iπ/(n+1)).
//
// ID:      NUM-SPEC-020
// Time:    O(n)   n = number of points
// Space:   O(n) for the results
// Source:  Abramowitz & Stegun, Handbook of Mathematical Functions, 1964, eq. 25.4.40
// Source:  NIST DLMF, Release 1.2.1, §18.5 (eq. 18.5.2: U_n(cos θ) = sin((n+1)θ)/sin θ)
// Tags:    special-function, quadrature, gauss-chebyshev, orthogonal-polynomial, roots
// Units:   dimensionless
// Range:   n ≥ 1. Errors: ErrOutOfRange.
// Example: n = 1: node 0, weight π/2
func GaussChebyshevU(n int) (nodes, weights []float64, err error) {
	if n < 1 {
		return nil, nil, fmt.Errorf("%w: n = %d, want ≥ 1", ErrOutOfRange, n)
	}
	nodes = make([]float64, n)
	weights = make([]float64, n)
	h := math.Pi / float64(n+1)
	for i := 0; i < n; i++ {
		k := n - i
		s := math.Sin(float64(k) * h)
		nodes[i] = math.Cos(float64(k) * h)
		weights[i] = h * s * s
	}
	return nodes, weights, nil
}

// gaussRule builds the n-point rule for the weight with total mass mu0 and orthonormal recurrence
// x qₖ = βₖ₊₁ qₖ₊₁ + αₖ qₖ + βₖ qₖ₋₁ (alpha(k), k ≥ 0; beta(k), k ≥ 1).
func gaussRule(n int, mu0 float64, alpha, beta func(k int) float64) (nodes, weights []float64, err error) {
	if n < 1 {
		return nil, nil, fmt.Errorf("%w: n = %d, want ≥ 1", ErrOutOfRange, n)
	}
	nodes = make([]float64, n)
	off := make([]float64, n)
	for k := 0; k < n; k++ {
		nodes[k] = alpha(k)
		if k+1 < n {
			off[k] = beta(k + 1)
		}
	}
	if err := tridiagonalEigenvalues(nodes, off); err != nil {
		return nil, nil, err
	}
	slices.Sort(nodes)
	weights = off[:n] // reuse the scratch slice for the weights
	q0 := 1 / math.Sqrt(mu0)
	betaN := beta(n)
	for i, x := range nodes {
		// Newton polish on q_n(x) = 0 (at most 3 steps; the eigenvalue is already within a few ulps of ‖J‖).
		for step := 0; step < 3; step++ {
			qn, dqn, _ := orthonormalAt(n, x, q0, betaN, alpha, beta)
			if dqn == 0 || math.IsNaN(qn) || math.IsInf(dqn, 0) {
				break
			}
			dx := qn / dqn
			x -= dx
			if math.Abs(dx) <= epsilon*math.Abs(x) {
				break
			}
		}
		nodes[i] = x
		_, _, sumSq := orthonormalAt(n, x, q0, betaN, alpha, beta)
		weights[i] = 1 / sumSq
	}
	return nodes, weights, nil
}

// orthonormalAt runs the orthonormal recurrence to degree n at x and returns q_n(x), q_n'(x) and
// Σ_(k<n) q_k(x)² (the reciprocal Christoffel number).
func orthonormalAt(n int, x, q0, betaN float64, alpha, beta func(k int) float64) (qn, dqn, sumSq float64) {
	prev, cur := 0.0, q0
	dprev, dcur := 0.0, 0.0
	bk := 0.0 // β_k, with β_0 = 0
	for k := 0; k < n; k++ {
		sumSq += cur * cur
		bNext := betaN
		if k+1 < n {
			bNext = beta(k + 1)
		}
		ak := alpha(k)
		next := ((x-ak)*cur - bk*prev) / bNext
		dnext := (cur + (x-ak)*dcur - bk*dprev) / bNext
		prev, cur = cur, next
		dprev, dcur = dcur, dnext
		bk = bNext
	}
	return cur, dcur, sumSq
}

// tridiagonalEigenvalues overwrites d with the eigenvalues of the symmetric tridiagonal matrix with
// diagonal d and off-diagonal e (e[i] couples i and i+1; e[n−1] is scratch), by the implicit QL method
// with Wilkinson-type shifts (Golub & Van Loan, Matrix Computations, 4th ed., §8.3; Numerical Recipes
// 3rd ed. §11.4, tqli without eigenvectors). e is destroyed.
func tridiagonalEigenvalues(d, e []float64) error {
	n := len(d)
	if n > 0 {
		e[n-1] = 0
	}
	for l := 0; l < n; l++ {
		for iter := 0; ; iter++ {
			m := l
			for ; m < n-1; m++ {
				dd := math.Abs(d[m]) + math.Abs(d[m+1])
				if math.Abs(e[m]) <= epsilon*dd {
					break
				}
			}
			if m == l {
				break
			}
			if iter >= 60 {
				return fmt.Errorf("%w: tridiagonal QL, eigenvalue %d", ErrNoConvergence, l)
			}
			g := (d[l+1] - d[l]) / (2 * e[l])
			r := math.Hypot(g, 1)
			g = d[m] - d[l] + e[l]/(g+math.Copysign(r, g))
			s, c, p := 1.0, 1.0, 0.0
			underflow := false
			for i := m - 1; i >= l; i-- {
				f := s * e[i]
				b := c * e[i]
				r = math.Hypot(f, g)
				e[i+1] = r
				if r == 0 {
					d[i+1] -= p
					e[m] = 0
					underflow = true
					break
				}
				s = f / r
				c = g / r
				g = d[i+1] - p
				r = (d[i]-g)*s + 2*c*b
				p = s * r
				d[i+1] = g + p
				g = c*r - b
			}
			if underflow {
				continue
			}
			d[l] -= p
			e[l] = g
			e[m] = 0
		}
	}
	return nil
}

// Package special evaluates the special functions of mathematical physics and statistics on real
// arguments: incomplete gamma and beta, beta, Bessel I₀, sinc, Riemann zeta, and the classical
// orthogonal polynomials with their Gauss quadrature rules. Gamma, log-gamma, erf and their kin are
// Go's math package and are not wrapped here (D-012).
//
// ID short form: SPEC (IDs are NUM-SPEC-NNN).
//
// special imports nothing from num/ (build brief R1): num/prob and num/signal import it.
// Every function is pure, allocation-free unless it returns a slice, and returns a wrapped sentinel
// (ErrNotFinite, ErrOutOfRange, ErrOverflow, ErrNoConvergence) instead of a NaN.
//
// Functions:
//
//	NUM-SPEC-001  RegularizedGammaP  regularized lower incomplete gamma P(a, x)
//	NUM-SPEC-002  RegularizedGammaQ  regularized upper incomplete gamma Q(a, x)
//	NUM-SPEC-003  Beta               beta function B(a, b)
//	NUM-SPEC-004  LogBeta            ln B(a, b)
//	NUM-SPEC-005  RegularizedBeta    regularized incomplete beta I_x(a, b)
//	NUM-SPEC-006  BesselI0           modified Bessel function I₀(x)
//	NUM-SPEC-007  Sinc               normalized sinc sin(πx)/(πx)
//	NUM-SPEC-008  RiemannZeta        Riemann zeta ζ(s) for s > 1
//	NUM-SPEC-009  LegendreP          Legendre polynomial P_n(x)
//	NUM-SPEC-010  ChebyshevT         Chebyshev polynomial of the first kind T_n(x)
//	NUM-SPEC-011  ChebyshevU         Chebyshev polynomial of the second kind U_n(x)
//	NUM-SPEC-012  HermiteH           physicists' Hermite polynomial H_n(x)
//	NUM-SPEC-013  HermiteHe          probabilists' Hermite polynomial He_n(x)
//	NUM-SPEC-014  LaguerreL          Laguerre polynomial L_n(x)
//	NUM-SPEC-015  GaussLegendre      Gauss–Legendre nodes (roots of P_n) and weights
//	NUM-SPEC-016  GaussHermite       Gauss–Hermite nodes (roots of H_n) and weights
//	NUM-SPEC-017  GaussHermiteE      probabilists' Gauss–Hermite nodes (roots of He_n) and weights
//	NUM-SPEC-018  GaussLaguerre      Gauss–Laguerre nodes (roots of L_n) and weights
//	NUM-SPEC-019  GaussChebyshevT    Gauss–Chebyshev (first kind) nodes (roots of T_n) and weights
//	NUM-SPEC-020  GaussChebyshevU    Gauss–Chebyshev (second kind) nodes (roots of U_n) and weights
package special

// Package algebra evaluates, multiplies and solves polynomials: Horner evaluation with the derivative,
// the product of two polynomials, lines in slope-intercept form, an accurate quadratic discriminant,
// quadratic and cubic roots in closed form, and all roots of any degree by Aberth–Ehrlich iteration.
//
// ID short form: ALGEB (IDs are NUM-ALGEB-NNN).
//
// Conventions: a polynomial is coeffs []float64 in ASCENDING order, coeffs[i] multiplies xⁱ, so the
// degree is len(coeffs)−1. Every value is dimensionless. Inputs are only read. NaN or ±Inf in returns
// ErrNotFinite; finite inputs with a result outside the float64 range return ErrOverflow. "No real
// roots" is a result (n = 0), not an error. A zero leading coefficient is ErrDegenerate: no function
// silently drops to a lower degree. Errors are wrapped with the parameter name; test with errors.Is.
//
// Functions:
//
//	NUM-ALGEB-001  PolynomialEvaluate            value of a polynomial at x by Horner's rule
//	NUM-ALGEB-002  PolynomialValueAndDerivative  value and first derivative of a polynomial at x in one pass
//	NUM-ALGEB-003  PolynomialMultiply            coefficients of the product of two polynomials
//	NUM-ALGEB-004  LinearEvaluate                y = m·x + b with one rounding
//	NUM-ALGEB-005  LinearSolve                   x from y = m·x + b
//	NUM-ALGEB-006  QuadraticDiscriminant         b² − 4ac to within 2 ulps
//	NUM-ALGEB-007  QuadraticRealRoots            real roots of a·x² + b·x + c without cancellation
//	NUM-ALGEB-008  QuadraticRoots                both complex roots of a·x² + b·x + c
//	NUM-ALGEB-009  CubicRoots                    the three complex roots of a cubic, Newton-polished
//	NUM-ALGEB-010  PolynomialRoots               all complex roots of a real polynomial (Aberth–Ehrlich)
//	NUM-ALGEB-011  PolynomialAdd                 coefficients of p + q
//	NUM-ALGEB-012  PolynomialSubtract            coefficients of p − q
//	NUM-ALGEB-013  PolynomialDivide              quotient and remainder of p ÷ q
//	NUM-ALGEB-014  PolynomialPower               coefficients of pᵏ by binary powering
//	NUM-ALGEB-015  PolynomialDerivative          coefficients of the m-th derivative
//	NUM-ALGEB-016  PolynomialIntegral            coefficients of the antiderivative with P(0) given
//	NUM-ALGEB-017  PolynomialFromRoots           real monic polynomial from its (conjugate-paired) roots
//	NUM-ALGEB-018  PolynomialValueFromRoots      ∏(x − rᵢ) without expanding
//	NUM-ALGEB-019  PolynomialCompanion           companion matrix of a polynomial
//	NUM-ALGEB-020  PolynomialTrim                drop small trailing coefficients
//	NUM-ALGEB-021  PolynomialEvaluate2D          p(x, y) by nested Horner
//	NUM-ALGEB-022  RoundDecimals                 round to d decimal places, ties to even
//	NUM-ALGEB-023  FloorDivMod                   floored quotient and remainder
//	NUM-ALGEB-024  FloorMod                      remainder with the sign of the divisor
//	NUM-ALGEB-025  Sign                          signum −1, 0, +1
//	NUM-ALGEB-026  Heaviside                     Heaviside step with a chosen H(0)
//	NUM-ALGEB-027  PolynomialEvaluateGrid2D      p(x, y) over the Cartesian product of xs and ys
//	NUM-ALGEB-028  PolynomialIntegralOrder       m-th repeated integral with lower bound and constants
//	NUM-ALGEB-029  PercentOf                     the part that is p % of a total
//	NUM-ALGEB-030  PercentFromPart               what percent a part is of a total
//	NUM-ALGEB-031  TotalFromPercent              the total of which a part is p %
//	NUM-ALGEB-032  PercentChange                 percent increase or decrease from an original value
//	NUM-ALGEB-033  ApplyPercentChange            a value increased or decreased by p %
//	NUM-ALGEB-034  PercentDifference             symmetric percent difference against the mean
//	NUM-ALGEB-035  PercentError                  relative error as a percent of the true value
//	NUM-ALGEB-036  ProportionFourthTerm          d in a : b = c : d by cross products
//	NUM-ALGEB-037  RoundDecimalsMode             round to d decimal places under an IEEE 754 rounding mode
package algebra

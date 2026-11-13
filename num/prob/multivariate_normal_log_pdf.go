package prob

// MultivariateNormalLogPDF returns ln f(x) = −½·(n·ln 2π + ln det Σ + (x − μ)ᵀΣ⁻¹(x − μ)) for the n-variate
// normal with Σ given by its lower-triangular Cholesky factor L (row-major n×n, as NUM-LINALG-152 Cholesky
// returns it; the strict upper triangle is not read). ln det Σ = 2·Σᵢ ln Lᵢᵢ, summed in logs, so the
// result stays finite where MultivariateNormalPDF underflows to 0 (far tails, high dimension).
// Errors: as MultivariateNormalPDF, except that no density overflow is possible: ErrDimension,
// ErrNotFinite, ErrOutOfRange (a diagonal entry of cholL ≤ 0), ErrOverflow (the quadratic form).
//
// ID:      NUM-PROB-095
// Time:    O(n²)   n = len(x) (forward substitution)
// Space:   O(1) extra for n ≤ 16 (scratch on the stack, no allocation); O(n), one allocation, above
// Source:  Kotz, Balakrishnan, Johnson, Continuous Multivariate Distributions, vol. 1: Models and Applications, 2nd ed., Wiley, 2000, ch. 45 "Multivariate Normal Distributions", §1 (density)
// Source:  Golub, Van Loan, Matrix Computations, 4th ed., Johns Hopkins, 2013, §3.1.1 (forward substitution), §4.2 (Cholesky)
// Tags:    probability, distribution, multivariate-normal, gaussian, log-pdf, log-likelihood, multivariate, cholesky
// Units:   x and mean: one unit per coordinate; result: natural log of probability per unit volume
// Example: n = 2, μ = 0, Σ = [[4,2],[2,3]] (L = [[2,0],[1,√2]]), x = (1,1) → −3/16 − ln(2π) − ½·ln 8 ≈ −3.0638
func MultivariateNormalLogPDF(x, mean, cholL []float64) (float64, error) {
	return multivariateNormalLog(x, mean, cholL)
}

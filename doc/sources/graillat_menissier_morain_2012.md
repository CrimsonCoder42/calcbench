# Graillat and Ménissier-Morain 2012

**Citation**: Graillat, Ménissier-Morain, "Accurate summation, dot product and polynomial evaluation in complex floating point arithmetic", *Information and Computation*, 2012, 216:57–71, doi:10.1016/j.ic.2011.09.003.

| Taken | Locator | Used by |
|---|---|---|
| TwoSum, TwoProduct (FMA form), and their complex versions | Algorithms 2.1, 3.1, 3.4 | NUM-ALGEB-010 (`compensatedHorner`, polynomial_roots.go:245, :279, :296) |
| Compensated Horner scheme in complex arithmetic | Algorithms 5.2, 5.4 | NUM-ALGEB-010 |
| Error bound of the compensated Horner scheme, used as the noise level of p(z) (γ̃₂ₙ ≈ 4√2·n·u, taken as ≤ 6n·u) | Theorem 5.5, eq. (5.17) | NUM-ALGEB-010 (`atNoise`, polynomial_roots.go:231) |

Our own, not from the paper: the error terms are summed in plain floating point (not AccSum), and p′ is compensated by the same recurrence; the convergence test |Δ| ≤ 4u·|z| is ours and documented as such.

How checked: `TestPolynomialRoots_wilkinson20` (polynomial_roots_test.go:134; every root within 1.74e-4 of k, against 8192-bit Newton references) and the DLMF §1.11 quartic Example (`TestPolynomialRoots_workedExample`).

Status: bibliographic record confirmed (DOI). Full paper PDF read first-hand, pp. 57–68 (convexoptimization.com/TOOLS/Graillat.pdf), 2026-11-12.

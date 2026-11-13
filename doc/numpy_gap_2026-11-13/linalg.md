# numpy gap: numpy.linalg and numpy.polynomial (agent A3, 2026-11-13)

READ-ONLY analysis. Inventory: `calcbench_functions.md` (grep by meaning), signatures confirmed by grep on
`num/linalg/*.go` and `num/algebra/*.go` (`go doc` could not be bounded: no `timeout` on this macOS, so not run).

**Pages read (WebFetch, 2026-11-13):** routines.linalg.html · routines.polynomials.html ·
routines.polynomials.polynomial.html · routines.polynomials.chebyshev.html · routines.polynomials.legendre.html ·
routines.polynomials.laguerre.html · routines.polynomials.hermite.html · routines.polynomials.hermite_e.html ·
routines.polynomials.polyutils.html · routines.polynomials.poly1d.html.
**Pages that failed:** routines.polynomials-poly1d.html (the URL in the brief) returned 404; the correct page is
routines.polynomials.poly1d.html, which was read instead.

Conventions noted: calcbench matrices are flat row-major `[]float64` with explicit dims; polynomial coefficients
are ASCENDING (`coeffs[i]` multiplies xⁱ), matching numpy.polynomial, the reverse of legacy `np.polyval`/`np.roots`
(descending). Real-valued only: complex/Hermitian variants are out of scope by the scalar/[]float64 rule.

**Counts (table rows; a family row covers several names):** HAVE 32 · PARTIAL 16 · MISSING, IN SCOPE 28 · OUT OF SCOPE 8

## numpy.linalg (routines.linalg)

| numpy name | what it computes | calcbench match | verdict | proposed home + primary source |
|---|---|---|---|---|
| matmul / `@` | matrix product | NUM-LINALG-090 MatMul (+057/066/078 fixed sizes) | HAVE | stacked/broadcast form is ndarray machinery |
| dot | dot product / matrix product | NUM-LINALG-039 VecDot, -040 VecDotCompensated, -090 MatMul, -091 MatVecMul | HAVE | |
| linalg.multi_dot | chained product in the cheapest order | ALG-DP-019 MatrixChainOrder (order only) + NUM-LINALG-090 MatMul | PARTIAL | no one routine: num/linalg `MatMulChain` using ALG-DP-019; Cormen et al., CLRS 3rd ed., §15.2 |
| vdot | dot of flattened vectors (conjugating complex) | NUM-LINALG-039 VecDot (real case; a flat matrix is already a vector) | HAVE | complex conjugation out of scope |
| vecdot | vector dot along last axis | NUM-LINALG-039 VecDot | HAVE | batching is broadcasting |
| linalg.vecdot | same (array-API name) | NUM-LINALG-039 VecDot | HAVE | |
| inner | inner product (1-D: dot; 2-D: A·Bᵀ) | NUM-LINALG-039 VecDot; 2-D = MatMul(A, MatTranspose(B)) | HAVE | 2-D form only by composition (an extra transpose copy) |
| outer | u·vᵀ | NUM-LINALG-120 VecOuter | HAVE | |
| linalg.outer | same | NUM-LINALG-120 VecOuter | HAVE | |
| linalg.matmul | same as matmul | NUM-LINALG-090 MatMul | HAVE | |
| matvec | A·v | NUM-LINALG-091 MatVecMul | HAVE | |
| vecmat | vᵀ·A | match?: MatVecMul(MatTranspose(A), v) | PARTIAL | no direct routine (costs a transpose); num/linalg `VecMatMul`; Golub & Van Loan 4th ed. §1.1 (gaxpy/row-oriented products) |
| tensordot | contraction over chosen axes of n-D tensors | fixed cases: axes=0 → VecOuter (-120), axes=1 → MatMul (-090) | OUT OF SCOPE | general n-D contraction is tensor machinery |
| linalg.tensordot | same | as above | OUT OF SCOPE | |
| linalg.matrix_power | Aⁿ for integer n (negative via inverse) | NUM-LINALG-122 MatPower (k ≥ 0 only, ErrOutOfRange for k < 0) | PARTIAL | negative n missing: MatInverse then MatPower; Golub & Van Loan 4th ed. §3.2 (inverse via LU), Knuth TAOCP Vol. 2, 3rd ed., §4.6.3 (binary powering) |
| kron | Kronecker product | NUM-LINALG-121 MatKronecker | HAVE | |
| linalg.cross | 3-vector cross product | NUM-LINALG-016 Vec3Cross (2-D: -006 Vec2Cross) | HAVE | |
| linalg.cholesky | A = L·Lᵀ | NUM-LINALG-152 Cholesky (+153 CholeskySolve) | HAVE | upper=True form is a transpose |
| linalg.qr | QR factorisation | NUM-LINALG-154 QRHouseholder (thin), -155 GramSchmidtModified | HAVE | mode='complete' (full m×m Q) and mode='r' not offered; reduced is numpy's default |
| linalg.svd | A = U·Σ·Vᵀ | NUM-LINALG-160 SVDJacobi (thin; wantVectors) | PARTIAL | full_matrices=True (numpy's default, square m×m U) missing: complete U by Householder completion of the thin basis; Golub & Van Loan 4th ed. §2.4, §8.6.3 |
| linalg.svdvals | singular values only | NUM-LINALG-160 SVDJacobi with wantVectors=false | HAVE | |
| linalg.eig | eigenvalues + right eigenvectors of a general square matrix | NUM-LINALG-149 Mat2Eigen (2×2 only), -150 MatSymmetricEigenJacobi (symmetric only), -151 MatPowerIteration (dominant pair) | PARTIAL | general nonsymmetric n×n missing: num/linalg `MatEigen` by Hessenberg reduction + Francis double-shift QR + inverse iteration; Golub & Van Loan 4th ed. §7.4–7.5, §7.6.1 |
| linalg.eigh | symmetric eigenvalues + vectors | NUM-LINALG-150 MatSymmetricEigenJacobi (+064, -076 closed forms) | HAVE | Hermitian complex out of scope |
| linalg.eigvals | eigenvalues of a general matrix | NUM-LINALG-149 (2×2), -148 CharPoly + NUM-ALGEB-010 PolynomialRoots (ill-conditioned route) | PARTIAL | same as eig: Francis QR on the Hessenberg form; Golub & Van Loan 4th ed. §7.5 |
| linalg.eigvalsh | symmetric eigenvalues | NUM-LINALG-150 (wantVectors=false), -064, -076 | HAVE | |
| linalg.norm | vector or matrix norm, ord selectable | NUM-LINALG-041/042/043/044 VecNorm/1/Inf/P; -163/164/165/166 MatNormFrobenius/1/Inf/2 | PARTIAL | missing: nuclear norm (Σσᵢ), matrix ord −1/−2/−∞ (min col sum, σ_min, min row sum), vector ord 0 (count non-zero), −∞ (min abs), p < 1; num/linalg; Golub & Van Loan 4th ed. §2.3, §2.4 |
| linalg.matrix_norm | matrix norm | NUM-LINALG-163..166 | PARTIAL | nuclear and negative-ord norms missing (as above); Horn & Johnson, Matrix Analysis 2nd ed., §5.6 |
| linalg.vector_norm | vector norm | NUM-LINALG-041..044 | PARTIAL | ord 0, −∞ and p < 1 missing; Golub & Van Loan 4th ed. §2.2 |
| linalg.cond | condition number in a chosen norm | NUM-LINALG-162 MatCondition2 (2-norm), -100 MatConditionEstimate1 (estimate) | PARTIAL | exact κ in 1, ∞, Frobenius norms (‖A‖·‖A⁻¹‖) missing; num/linalg; Higham, Accuracy and Stability of Numerical Algorithms 2nd ed., §6.2, §7.1 |
| linalg.det | determinant | NUM-LINALG-097 MatDet (+060/069/081 fixed sizes) | HAVE | |
| linalg.matrix_rank | numerical rank | NUM-LINALG-099 MatRank (complete pivoting, tol) | HAVE | numpy uses SVD; the SVD-based rank is composable from -160 |
| linalg.slogdet | sign and log|det| | none (MatDet keeps mantissa + exponent internally but returns det, ErrOverflow past float64) | MISSING, IN SCOPE | num/linalg `MatLogDet` (sign, log|det|) from the same LU; Golub & Van Loan 4th ed. §3.2; Higham 2nd ed. §14.6 (determinant computation) |
| trace | sum of a diagonal | NUM-LINALG-093 MatTrace (+062/071/083) | HAVE | offset diagonals not offered (trivial) |
| linalg.trace | same | NUM-LINALG-093 MatTrace | HAVE | |
| linalg.solve | solve A·x = b | NUM-LINALG-096 MatSolve, -094/095 LU, -153, -157/158, -159, -144/145 | HAVE | multiple right-hand sides (B matrix) by repeated LUSolve |
| linalg.tensorsolve | tensor equation | none | OUT OF SCOPE | n-D tensor machinery |
| linalg.lstsq | least-squares solution (min-norm, rank-deficient OK; returns residuals, rank, s) | NUM-LINALG-156 LeastSquaresQR (full rank, m ≥ n only; X + ResidualNorm) | PARTIAL | rank-deficient and underdetermined minimum-norm solution missing: num/linalg `LeastSquaresSVD` via -160/-161 or column-pivoted QR; Golub & Van Loan 4th ed. §5.4–5.5 |
| linalg.inv | inverse | NUM-LINALG-098 MatInverse (+061/070/082) | HAVE | |
| linalg.pinv | Moore–Penrose pseudo-inverse | NUM-LINALG-161 MatPseudoInverse | HAVE | |
| linalg.tensorinv | tensor 'inverse' | none | OUT OF SCOPE | n-D tensor machinery |
| diagonal | extract a (offset) diagonal | none | MISSING, IN SCOPE | low: num/linalg `MatDiagonal(a, rows, cols, offset)`; definitional (Golub & Van Loan 4th ed. §1.2.1 band/diagonal notation) |
| linalg.diagonal | same | none | MISSING, IN SCOPE | same row as above |
| linalg.matrix_transpose | transpose | NUM-LINALG-092 MatTranspose | HAVE | |

Also on the page, out of scope: `einsum` (general engine; its fixed cases are VecDot, VecOuter, MatMul, MatTrace,
MatHadamard -119, all HAVE), `einsum_path`, `LinAlgError` (calcbench uses sentinel errors).

## numpy.polynomial.polynomial (power series) and legacy poly1d

| numpy name | what it computes | calcbench match | verdict | proposed home + primary source |
|---|---|---|---|---|
| polyadd | p + q | none | MISSING, IN SCOPE | num/algebra `PolynomialAdd`; Knuth TAOCP Vol. 2, 3rd ed., §4.6 |
| polysub | p − q | none | MISSING, IN SCOPE | num/algebra `PolynomialSubtract`; Knuth TAOCP Vol. 2, 3rd ed., §4.6 |
| polymulx | p·x | NUM-ALGEB-003 PolynomialMultiply with [0, 1] | HAVE | by composition (a shift) |
| polymul | p·q | NUM-ALGEB-003 PolynomialMultiply | HAVE | O(nm) direct; FFT form is A4's area |
| polydiv | quotient and remainder | none | MISSING, IN SCOPE | num/algebra `PolynomialDivide`; Knuth TAOCP Vol. 2, 3rd ed., §4.6.1 (Algorithm D) |
| polypow | pⁿ | match?: repeated NUM-ALGEB-003 | MISSING, IN SCOPE | low: num/algebra `PolynomialPower` by binary powering; Knuth TAOCP Vol. 2, 3rd ed., §4.6.3 |
| polyval | p(x) at points | NUM-ALGEB-001 PolynomialEvaluate (scalar x, Horner + FMA) | HAVE | slice-of-x form is the vec-layer question |
| polyval2d / polygrid2d | bivariate p(x, y) from a coefficient matrix; on a grid | none | MISSING, IN SCOPE | low: num/algebra `PolynomialEvaluate2D` (nested Horner over [][]float64); Knuth TAOCP Vol. 2, 3rd ed., §4.6.4 |
| polyval3d / polyvalnd / polygrid3d | 3-D / n-D polynomials | none | OUT OF SCOPE | need 3-D/n-D coefficient tensors |
| polyder | coefficients of the m-th derivative (scaled) | NUM-ALGEB-002 PolynomialValueAndDerivative (value of p′ at x only) | PARTIAL | derivative coefficient vector and m-th order missing: num/algebra `PolynomialDerivative`; Knuth TAOCP Vol. 2, 3rd ed., §4.6 (formal derivative) |
| polyint | antiderivative coefficients (order m, constants k, lower bound) | none | MISSING, IN SCOPE | num/algebra `PolynomialIntegral`; DLMF §1.4 (term-by-term integration of xⁿ); confirm subsection |
| polyfromroots | monic polynomial with given roots | none | MISSING, IN SCOPE | num/algebra `PolynomialFromRoots` (real roots, conjugate pairs paired so coefficients stay real); Knuth TAOCP Vol. 2, 3rd ed., §4.6 |
| polyroots | roots | NUM-ALGEB-010 PolynomialRoots (Aberth–Ehrlich), -007/-008/-009 closed forms | HAVE | numpy uses companion eigenvalues; Aberth is a better-conditioned method |
| polyvalfromroots | Π(x − rᵢ) at x | none | MISSING, IN SCOPE | low: num/algebra `PolynomialEvaluateFromRoots`; definitional (Knuth TAOCP Vol. 2 §4.6) |
| polyvander | Vandermonde matrix | none | MISSING, IN SCOPE | num/linalg `Vandermonde`; Golub & Van Loan 4th ed. §4.6 (Vandermonde systems); confirm section |
| polyvander2d / polyvander3d | pseudo-Vandermonde for 2-D/3-D points | none | OUT OF SCOPE | multivariate design matrices; revisit only with a regression need |
| polycompanion | companion matrix | none (MatCharPoly goes the other way, -148) | MISSING, IN SCOPE | num/algebra or num/linalg `PolynomialCompanion`; Edelman & Murakami, "Polynomial roots from companion matrix eigenvalues", Math. Comp. 64 (1995) 763–776 |
| polyfit | weighted least-squares polynomial fit of degree n | NUM-STATS-009 LinearRegression (degree 1); NUM-LINALG-156 LeastSquaresQR on a hand-built Vandermonde | PARTIAL | general degree, weights, domain scaling missing: num/stats or num/linalg `PolynomialFit` = Vandermonde + Householder QR (never normal equations); Golub & Van Loan 4th ed. §5.3; Björck, Numerical Methods for Least Squares Problems, SIAM 1996, §4.? (to confirm) |
| polytrim | drop small trailing coefficients | none | MISSING, IN SCOPE | low: num/algebra `PolynomialTrim`; definitional, no formula source needed beyond the tolerance policy (doc/rules numeric.md) |
| polyline | [off, scl] as a polynomial | none | OUT OF SCOPE | a constructor, no maths |
| np.poly (legacy) | coefficients from roots, or char. polynomial of a matrix | NUM-LINALG-146/147/148 CharPoly (matrix input) | PARTIAL | roots input = polyfromroots (MISSING above) |
| np.polyval (legacy) | p(x), descending coefficients | NUM-ALGEB-001 (ascending order) | HAVE | order reversal only |
| np.roots (legacy) | roots, descending coefficients | NUM-ALGEB-010 | HAVE | |
| np.polyfit (legacy) | least-squares fit, optional cov | NUM-STATS-009 (degree 1 only) | PARTIAL | as polyfit; coefficient covariance (σ²(VᵀV)⁻¹) also missing; Golub & Van Loan 4th ed. §5.3 |
| np.polyder (legacy) | derivative coefficients | NUM-ALGEB-002 (value only) | PARTIAL | as polyder |
| np.polyint (legacy) | antiderivative coefficients | none | MISSING, IN SCOPE | as polyint |
| np.polyadd / np.polysub (legacy) | sum / difference | none | MISSING, IN SCOPE | as polyadd / polysub |
| np.polymul (legacy) | product | NUM-ALGEB-003 | HAVE | |
| np.polydiv (legacy) | quotient and remainder | none | MISSING, IN SCOPE | as polydiv |

## Orthogonal-polynomial families (one row per family per group; every name listed)

None of the five families exists in calcbench (grep for Chebyshev/Legendre/Laguerre/Hermite finds only NUM-PROB-020/021
Chebyshev's inequality, unrelated). Proposed home for all: **NEW `num/orthopoly` (SHORT `ORTHO`)**, or inside
num/algebra if Devin prefers no new package; the Gauss rules could instead go to the planned num/calculus
(integration). Primary sources: DLMF Ch. 18 (§18.3 definitions, §18.9 recurrences and derivatives, §18.17 integrals),
Clenshaw, "A note on the summation of Chebyshev series", MTAC 9 (1955) 118–120 (evaluation), Golub & Welsch,
"Calculation of Gauss quadrature rules", Math. Comp. 23 (1969) 221–230 + DLMF §3.5 (quadrature), Trefethen,
Approximation Theory and Approximation Practice, SIAM 2013, ch. 2–3 (Chebyshev points, interpolation).

| numpy name | what it computes | calcbench match | verdict | proposed home + primary source |
|---|---|---|---|---|
| Chebyshev core: chebadd, chebsub, chebmulx, chebmul, chebdiv, chebpow, chebval, chebder, chebint, chebfromroots, chebroots, chebvander, chebcompanion, chebfit, chebtrim, chebline, cheb2poly, poly2cheb | Chebyshev-series arithmetic, Clenshaw evaluation, calculus, roots, fit, basis conversion | none | MISSING, IN SCOPE | num/orthopoly; chebval by Clenshaw 1955; der/int DLMF §18.9; fit Golub & Van Loan §5.3 on chebvander |
| Chebyshev nodes and quadrature: chebpts1, chebpts2, chebinterpolate, chebgauss, chebweight | Chebyshev points (1st/2nd kind), interpolant at them, Gauss–Chebyshev rule, weight 1/√(1−x²) | none | MISSING, IN SCOPE | num/orthopoly; Trefethen ATAP ch. 2–3; DLMF §3.5(v) |
| Legendre core: legadd, legsub, legmulx, legmul, legdiv, legpow, legval, legder, legint, legfromroots, legroots, legvander, legcompanion, legfit, legtrim, legline, leg2poly, poly2leg | Legendre-series counterparts | none | MISSING, IN SCOPE | num/orthopoly; DLMF §18.3, §18.9 |
| Legendre quadrature: leggauss, legweight | Gauss–Legendre nodes and weights; weight 1 | none | MISSING, IN SCOPE | num/orthopoly or num/calculus; Golub & Welsch 1969; DLMF §3.5(v) |
| Laguerre core: lagadd, lagsub, lagmulx, lagmul, lagdiv, lagpow, lagval, lagder, lagint, lagfromroots, lagroots, lagvander, lagcompanion, lagfit, lagtrim, lagline, lag2poly, poly2lag | Laguerre-series counterparts | none | MISSING, IN SCOPE | num/orthopoly; DLMF §18.3, §18.9 |
| Laguerre quadrature: laggauss, lagweight | Gauss–Laguerre rule on [0, ∞); weight e⁻ˣ | none | MISSING, IN SCOPE | Golub & Welsch 1969; DLMF §3.5(v) |
| Hermite (physicists') core: hermadd, hermsub, hermmulx, hermmul, hermdiv, hermpow, hermval, hermder, hermint, hermfromroots, hermroots, hermvander, hermcompanion, hermfit, hermtrim, hermline, herm2poly, poly2herm | Hermite Hₙ series counterparts | none | MISSING, IN SCOPE | num/orthopoly; DLMF §18.3, §18.9 |
| Hermite quadrature: hermgauss, hermweight | Gauss–Hermite rule; weight e^(−x²) | none | MISSING, IN SCOPE | Golub & Welsch 1969; DLMF §3.5(v) |
| HermiteE (probabilists') core: hermeadd, hermesub, hermemulx, hermemul, hermediv, hermepow, hermeval, hermeder, hermeint, hermefromroots, hermeroots, hermevander, hermecompanion, hermefit, hermetrim, hermeline, herme2poly, poly2herme | Heₙ series counterparts | none | MISSING, IN SCOPE | num/orthopoly; DLMF §18.3 (Heₙ), §18.9 |
| HermiteE quadrature: hermegauss, hermeweight | Gauss–HermiteE rule; weight e^(−x²/2) | none | MISSING, IN SCOPE | Golub & Welsch 1969; DLMF §3.5(v) |
| all families, multivariate: {cheb,leg,lag,herm,herme}val2d/val3d/valnd/grid2d/grid3d/vander2d/vander3d | n-D series evaluation and design matrices | none | OUT OF SCOPE | n-D coefficient tensors (2-D evaluation could follow polyval2d if ever needed) |
| polyutils: mapdomain, mapparms, getdomain | affine map between intervals, its (off, scl), the data's [min, max] | none | MISSING, IN SCOPE | low: helper for every *fit; num/orthopoly; definitional affine map (Trefethen ATAP ch. 2, change of interval) |

## Out of scope, listed once

The convenience classes `Polynomial`, `Chebyshev`, `Legendre`, `Laguerre`, `Hermite`, `HermiteE`, `poly1d` (types
with methods); the constants `polydomain/polyzero/polyone/polyx` and each family's `*domain/*zero/*one/*x`;
polyutils `as_series`, `trimseq` (array plumbing; `trimcoef` = polytrim above); `einsum`, `einsum_path`;
`tensordot`, `tensorsolve`, `tensorinv` (n-D); `LinAlgError`; stacked/batched forms of every linalg routine and
complex/Hermitian variants (real `[]float64` only).

## Hand loops in calcbench

- fin/interest/net_present_value.go, fin/interest/internal_rate_of_return.go, fin/invest/modified_internal_rate_of_return.go,
  fin/invest/profitability_index.go, fin/invest/macaulay_duration.go: each evaluates a polynomial in 1/(1+r) by its own
  Horner loop (NPV is p(v), v = 1/(1+r)); NUM-ALGEB-001/002 PolynomialEvaluate/ValueAndDerivative do the same maths
  (IRR's Newton step needs exactly ValueAndDerivative). Not a defect (fin may import num), a reuse candidate.
- NUM-STATS-009 LinearRegression is a degree-1 polyfit by its own sums; a future PolynomialFit should agree with it
  at degree 1 (a cross-check test, not a replacement; the public API is additive).
- NUM-LINALG-135 Mat3Characteristics / -148 CharPoly + NUM-ALGEB-010 is today the only n×n route to nonsymmetric
  eigenvalues; a Francis QR `MatEigen` would replace that route for anything beyond 3×3.

## Most important MISSING / PARTIAL in scope (ranked)

1. **polyfit** (general degree, weights): Vandermonde + QR (`Vandermonde`, `PolynomialFit`). Highest consumer value.
2. **eig / eigvals** for general nonsymmetric n×n (Hessenberg + Francis QR), GvL §7.5.
3. **lstsq** rank-deficient / underdetermined minimum-norm (SVD or pivoted QR), GvL §5.4–5.5.
4. **Polynomial arithmetic set**: polyadd, polysub, polydiv, polyder (coefficients), polyint, polyfromroots, polycompanion.
5. **Gauss quadrature nodes/weights** (leggauss first) and **Chebyshev** chebval/chebfit/chebpts/chebinterpolate,
   a NEW num/orthopoly (ORTHO) package for Devin to approve; also slogdet (small, already half-built inside MatDet).

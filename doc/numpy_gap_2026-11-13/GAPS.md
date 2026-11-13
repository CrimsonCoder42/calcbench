# GAPS: what numpy has that calcbench is missing (2026-11-13)

This file merges every MISSING, IN SCOPE row (M) and every PARTIAL row (P) from math.md (A1), stats.md (A2),
linalg.md (A3) and fft_special.md (A4). The agent files hold 156 such rows (39 P, 117 M).

When the same gap appeared more than once, it is listed once here. The merges were:
- i0 and sinc, which appeared in both A1 and A4;
- around and round;
- ediff1d and diff;
- cumulative_prod and cumprod;
- linalg.diagonal and diagonal;
- the legacy np.poly* routines and their numpy.polynomial twins;
- matrix_norm and vector_norm, folded into norm;
- remainder and mod.

That leaves **142 rows in 16 groups**. Each family row in num/orthopoly covers about 18 routine names.

Where agents proposed different homes for the same thing, the lead picked one; the agent files keep each agent's
original proposal. "NEW" marks a package that is not in doc/taxonomy.md, which makes it a decision for Devin.
Sources are given in short form; the full citation is in the agent file named in brackets.

| # | Home | Rows | M | P |
|---|---|---|---|---|
| 1 | num/prob (distributions) | 29 | 23 | 6 |
| 2 | NaN-skipping family (one policy decision, then num/stats + the reduction home) | 17 | 17 | 0 |
| 3 | NEW num/signal (SIGNAL): FFT, windows, convolution | 17 | 16 | 1 |
| 4 | num/linalg (matrix routines) | 13 | 3 | 10 |
| 5 | num/algebra (polynomial arithmetic) | 12 | 10 | 2 |
| 6 | NEW num/orthopoly (ORTHO): Chebyshev, Legendre, Laguerre, Hermite, quadrature | 11 | 11 | 0 |
| 7 | Slice reductions and elementwise (num/linalg Vec*, algo/prefix, algo/selection) | 9 | 6 | 3 |
| 8 | num/stats | 7 | 3 | 4 |
| 9 | algo/sorting, algo/hashing (unique, sets, counting) | 7 | 3 | 4 |
| 10 | num/algebra (scalar elementwise: rounding, floored division, sign) | 6 | 5 | 1 |
| 11 | NEW num/float (FLOAT): closeness, ulp, machine constants | 4 | 3 | 1 |
| 12 | num/calculus (planned in the taxonomy): gradient, trapezoid, interp | 3 | 2 | 1 |
| 13 | NEW num/special (SPECL): i0, sinc (plus the incomplete gamma/beta and zeta that group 1 needs) | 2 | 2 | 0 |
| 14 | num/integer: gcd, lcm | 2 | 2 | 0 |
| 15 | num/logs: logaddexp, logaddexp2 | 2 | 2 | 0 |
| 16 | num/trig: unwrap | 1 | 0 | 1 |
| | **Total** | **142** | **108** | **34** |

## 1. num/prob: distributions, maths only (29) [fft_special.md]

| numpy | V | calcbench today | missing / proposed + source |
|---|---|---|---|
| binomial | P | NUM-PROB-026/027/028 | quantile (search on the CDF); Johnson, Kemp & Kotz 3rd ed. ch. 3 |
| poisson | P | NUM-PROB-041/042/043 | quantile; JKK ch. 4 |
| hypergeometric | P | NUM-PROB-035/036/037 | quantile; JKK ch. 6 |
| negative_binomial | P | NUM-PROB-032/033/034 (counts trials) | failures convention, quantile; JKK ch. 5 |
| chisquare | P | NUM-PROB-063/064/065 | PDF, mean and variance; DLMF §8.2 |
| multinomial | P | NUM-COMBI-008 (coefficient) | PMF in log space; Johnson, Kotz & Balakrishnan 1997 ch. 35 |
| gamma, standard_gamma | M | unexported `regularizedGamma`, num/prob/chi_square_cdf.go:45 | PDF/CDF/quantile; DLMF §8.2.4, 8.7.1, 8.9.2 |
| beta | M | none | needs the regularized incomplete beta (group 13); DLMF §8.17.22; Numerical Recipes (NR) 3rd ed. §6.4 |
| standard_t | M | none | via the incomplete beta; A&S 26.7.1; Hill 1970 CACM Alg. 396 |
| f | M | none | via the incomplete beta; A&S 26.6.2 |
| dirichlet | M | none | PDF in log space; Kotz, Balakrishnan & Johnson vol. 1 ch. 49 |
| gumbel, laplace, logistic, pareto, power, rayleigh, standard_cauchy, triangular, weibull | M (9) | none | closed forms; Johnson, Kotz & Balakrishnan, Continuous Univariate Distributions 2nd ed., chapter per family |
| logseries | M | none | JKK ch. 7 |
| multivariate_hypergeometric | M | none | product of binomial coefficients; JKB 1997 ch. 39 |
| multivariate_normal | M | NUM-LINALG-152 Cholesky exists | PDF via Cholesky and log-det; KBJ vol. 1 ch. 45 |
| noncentral_chisquare | M | none | Ding 1992, AS 275 (low priority) |
| noncentral_f | M | none | Lenth 1987, AS 226 (low priority) |
| vonmises | M | none | needs i0; Hill 1977, ACM TOMS Alg. 518 |
| wald | M | NUM-PROB-048 NormalCDF | Chhikara & Folks 1989 ch. 2 |
| zipf | M | none | needs the Riemann zeta (DLMF §25.2) |

## 2. NaN-skipping family (17): one policy decision [math.md, stats.md]

numeric.md treats a NaN input as an error. Every row here conflicts with that rule, so a single D-id unlocks or
closes all 17 rows. Two options:
- **(A)** One exported filter, for example `FiniteValues(dst, values)`, that callers compose with the existing
  functions. A2 recommends this.
- **(B)** Explicitly named `...SkipNaN` variants. A1 recommends this.

The routines: nansum · nanprod · nancumsum · nancumprod · nanmax · nanmin · nanargmax · nanargmin · nanmean ·
nanvar · nanstd · nanmedian · nanpercentile · nanquantile · fmax · fmin · nan_to_num.
Source for the NaN semantics: IEEE 754-2019 §6.2 and §9.6 (maxNum/minNum in the 2008 edition, §5.3.1).

## 3. NEW num/signal (SIGNAL) (17) [fft_special.md, stats.md, math.md]

No FFT, DFT or window code exists anywhere in the tree.

| numpy | V | proposed + source |
|---|---|---|
| fft, ifft | M | radix-2 Cooley & Tukey 1965, Math. Comp. 19:297–301; any n by Bluestein 1970, IEEE AU-18(4):451–455 |
| rfft, irfft | M | NR 3rd ed. §12.3.1 |
| hfft, ihfft | M | Oppenheim & Schafer 3rd ed. §8.6 (low priority) |
| fftfreq, rfftfreq | M | Oppenheim & Schafer §8.1 |
| fftshift, ifftshift | M | cyclic shift by ⌊n/2⌋ and ⌈n/2⌉ |
| bartlett, blackman, hamming, hanning | M | Harris 1978, Proc. IEEE 66(1):51–83, §V eqs. (20), (32), (30a), (27a) |
| kaiser | M | Kaiser 1974; Harris 1978 eq. (46); needs i0 (group 13) |
| correlate | M | Oppenheim & Schafer §2.3; A2's other proposed home was num/series |
| convolve | P | NUM-ALGEB-003 PolynomialMultiply is mode 'full'; missing modes 'same'/'valid' and an FFT path for long inputs |

## 4. num/linalg: matrix routines (13) [linalg.md]

| numpy | V | calcbench today | missing / proposed + source |
|---|---|---|---|
| eig, eigvals | P (2) | NUM-LINALG-149 (2×2), -150 (symmetric) | general nonsymmetric n×n: Hessenberg + Francis QR, GvL 4th ed. §7.5 (already LATER in S17) |
| lstsq | P | NUM-LINALG-156 (full rank, m ≥ n) | minimum-norm rank-deficient or underdetermined solve via -160/-161; GvL §5.4–5.5 |
| svd | P | NUM-LINALG-160 (thin) | full U (full_matrices); GvL §2.4 |
| norm (incl. matrix_norm, vector_norm) | P | NUM-LINALG-041…044, 163…166 | nuclear norm; negative orders; vector ord 0; Horn & Johnson 2nd ed. §5.6 |
| cond | P | NUM-LINALG-162, -100 (estimate) | exact κ in the 1, ∞ and Frobenius norms; Higham 2nd ed. §6.2 |
| matrix_power | P | NUM-LINALG-122 (k ≥ 0) | negative k via MatInverse |
| multi_dot | P | ALG-DP-019 + MatMul | one `MatMulChain` routine; CLRS 3rd ed. §15.2 |
| vecmat | P | MatVecMul on a transpose | direct `VecMatMul`; GvL §1.1 |
| polyfit (np.polyfit) | P | NUM-STATS-009 (degree 1 only) | any degree, weights, coefficient covariance: Vandermonde + Householder QR (-156); GvL §5.3 |
| slogdet | M | none (MatDet keeps the mantissa and exponent internally) | `MatLogDet`; Higham §14.6 |
| diagonal | M | none | `MatDiagonal(a, rows, cols, offset)`; GvL §1.2.1 |
| polyvander | M | none | `Vandermonde`; GvL §4.6 (A3 asks for this section to be confirmed) |

## 5. num/algebra: polynomial arithmetic (12) [linalg.md]

Coefficients are ascending, as in NUM-ALGEB-001…010.

| numpy | V | proposed + source |
|---|---|---|
| polyadd, polysub | M (2) | Knuth TAOCP vol. 2 3rd ed. §4.6 |
| polydiv | M | Knuth §4.6.1 Algorithm D |
| polypow | M | binary powering of NUM-ALGEB-003; Knuth §4.6.3 |
| polyder | P | NUM-ALGEB-002 returns only the value; missing the derivative coefficients and m-th order |
| polyint | M | term by term; DLMF §1.4 (A3 asks for the subsection to be confirmed) |
| polyfromroots | M | conjugate pairs kept together so the coefficients stay real; Knuth §4.6 |
| polyvalfromroots | M | definitional |
| polycompanion | M | Edelman & Murakami 1995, Math. Comp. 64:763–776 |
| polytrim | M | tolerance policy per numeric.md |
| polyval2d / polygrid2d | M | nested Horner; Knuth §4.6.4 |
| np.poly | P | NUM-LINALG-146/147/148 take a matrix; the roots input is polyfromroots |

## 6. NEW num/orthopoly (ORTHO) (11) [linalg.md]

calcbench has none of these routines. One row per family per group; each "core" row is roughly 18 names:
add, sub, mulx, mul, div, pow, val, der, int, fromroots, roots, vander, companion, fit, trim, line, x2poly, poly2x.
- Chebyshev core: Clenshaw 1955; DLMF §18.9
- Chebyshev points and quadrature (chebpts1/2, chebinterpolate, chebgauss, chebweight): Trefethen, ATAP ch. 2–3
- Legendre core, and leggauss/legweight
- Laguerre core, and laggauss/lagweight
- Hermite core, and hermgauss/hermweight
- HermiteE core, and hermegauss/hermeweight
- polyutils (mapdomain, mapparms, getdomain)

Sources: DLMF §18.3 and §18.9; Golub & Welsch 1969, Math. Comp. 23:221–230 for the Gauss nodes; DLMF §3.5(v).

## 7. Slice reductions and elementwise (9) [math.md, stats.md]

| numpy | V | calcbench today | proposed + source |
|---|---|---|---|
| sum | P | ALG-PREFIX-001 PrefixSums (last element); no exported Sum | compensated `VecSum`; Neumaier 1974, ZAMM 54:39–51; Higham §4.3 |
| prod | M | none (ALG-PREFIX-010 ProductExceptSelf is different) | `VecProduct` with scaling; Higham §3.1 |
| cumprod (+ cumulative_prod) | M | none | algo/prefix `PrefixProducts` |
| diff (+ ediff1d) | M | none exported | algo/prefix `Differences(values, order)`; A&S §25.1 |
| argmax | M | none | algo/selection `ArgMax`; Knuth TAOCP vol. 1 §1.2.10 Algorithm M |
| argmin | M | none | algo/selection `ArgMin`; same source |
| argpartition | M | ALG-SELECT-004 returns the value, not the index | `KthSmallestIndex`; CLRS §9.2–9.3 |
| multiply | P | NUM-LINALG-037 VecScale (scalar only) | `VecHadamard` for two slices; Horn & Johnson, Topics §5.0 (S17 built MatHadamard, NUM-LINALG-119, for matrices) |
| divide | P | NUM-LINALG-038 VecDiv (by a scalar only) | elementwise slice ÷ slice with zero checks |

## 8. num/stats (7) [stats.md]

| numpy | V | calcbench today | missing / proposed + source |
|---|---|---|---|
| cov | M | none at all | `SampleCovariance(x, y)` and `CovarianceMatrix`; NIST e-Handbook §6.5.4.1; Welford 1962 |
| percentile, quantile | P (2) | NUM-STATS-015 (all nine Hyndman–Fan types) | numpy's lower/higher/nearest/midpoint (not in H&F: need a source or a decision); weights; several p in one sort |
| corrcoef | P | NUM-STATS-014 (two variables) | k×k matrix; NR 3rd ed. §14.5 |
| histogram | P | NUM-STATS-008 FrequencyDistribution | density, weights, automatic bins; out-of-range values are an error by policy |
| histogram_bin_edges | M | none | Sturges 1926; Scott 1979; Freedman & Diaconis 1981; Doane 1976 |
| histogram2d | M | none | `FrequencyDistribution2D`; Scott 2015 §3.4 |

## 9. algo/sorting and algo/hashing (7) [stats.md]

| numpy | V | calcbench today | missing / proposed + source |
|---|---|---|---|
| bincount | P | ALG-HASH-001 FrequencyCount (a map) | dense `[]int` (and weighted) result indexed by value |
| lexsort | P | ALG-SORT-009 SortIndices with a compare func | `LexSortIndices(keys [][]float64)`; Knuth TAOCP vol. 3 §5.2 |
| unique_all | P | none combined | first-occurrence indices and inverse; Knuth vol. 3 §5.2 |
| unique_counts | P | ALG-HASH-001 (unordered) | sorted values and counts as two slices |
| unique_inverse | M | none | label encoding via ALG-SORT-009 |
| isin | M | none | `IsIn`, O(n + m); Knuth vol. 3 §6.4 |
| setxor1d | M | can be composed from two set calls | `SetSymmetricDifference`; Knuth vol. 3 §6.4 |

## 10. num/algebra: scalar elementwise (6) [math.md]

| numpy | V | proposed + source |
|---|---|---|
| round / around | P | the Go stdlib has math.RoundToEven at 0 decimals only; `RoundDecimals`; IEEE 754-2019 §5.3.1; Muller et al. 2018 |
| floor_divide, divmod | M (2) | `FloorDivMod`; Knuth TAOCP vol. 1 §1.2.4 |
| mod / remainder | M | `FloorMod` (math.Mod truncates); Knuth §1.2.4 |
| sign | M | DLMF §1.16(iv) |
| heaviside | M | DLMF §1.16.13 |

## 11. NEW num/float (FLOAT) (4) [math.md]

| numpy | V | proposed + source |
|---|---|---|
| isclose | M | `IsClose(a, b, relTol, absTol)`; Knuth TAOCP vol. 2 §4.2.2 eqs. (21)–(24). Decision: Knuth's symmetric form or numpy's asymmetric one |
| allclose | M | `AllClose` over slices; same source |
| spacing | M | `Ulp`; Muller et al. 2018 §2.6.1 |
| finfo | P | named `MachineEpsilon` and `SmallestNormal`; IEEE 754-2019 Table 3.5 |

## 12. num/calculus (planned) (3) [math.md]

| numpy | V | proposed + source |
|---|---|---|
| gradient | M | Fornberg 1988, Math. Comp. 51(184):699–706 |
| trapezoid | M | A&S §25.4.2; NR 3rd ed. §4.1 |
| interp | P | NUM-LINALG-049 VecLerp interpolates within one interval only; a tabulated interpolation (search, then lerp) is missing; NR §3.1–3.3 (A1 also offered a NEW num/interp) |

## 13. NEW num/special (SPECL) (2) [fft_special.md]

| numpy | V | proposed + source |
|---|---|---|
| i0 | M | DLMF §10.25.2 and §10.40.1; A&S 9.8.1–9.8.2 |
| sinc | M | Oppenheim & Schafer §2.9; a Taylor branch near 0 |

This package would also hold the non-numpy primitives that groups 1 and 3 need: the regularized incomplete gamma
(today unexported in num/prob), the regularized incomplete beta, and the Riemann zeta.

## 14. num/integer (2) [math.md]

gcd and lcm are both MISSING. The taxonomy already reserves them for num/integer (D-005), but nothing is exported.
Private `gcd` helpers sit at num/integer/modular_arithmetic.go:38, num/combinatorics/binomial.go:60 and
algo/geometry2d/max_points_on_line.go:62. Source: Knuth TAOCP vol. 2 §4.5.2, Algorithms A and B.

## 15. num/logs (2) [math.md]

logaddexp and logaddexp2 are both MISSING, along with a slice LogSumExp. Source: Blanchard, Higham & Higham 2021,
IMA J. Numer. Anal. 41(4):2311–2330.

## 16. num/trig (1) [math.md]

unwrap is PARTIAL: NUM-TRIG-005 AngleDifference is the per-step kernel, but there is no sequence form and no
custom period. Source: Itoh 1982, Applied Optics 21(14):2470.

## Verdicts after the build (S18, doc/numpy_build_2026-11-13/, 2026-11-13)

This section records the build's verdict for every row above.

**Tally of the 142 rows:**
- **BUILT 112.**
- **PARTLY BUILT 10.** Part of the row is built; the rest was DROPPED or needs a source.
- **DROPPED 20.** 17 under D2, 2 under D1 (hfft and ihfft) and 1 under D1 (polyutils).
- **NEEDS SOURCE 0 whole rows.** Two parts of rows need a source: weights for percentile/quantile and for histogram.

S18 built 208 functions. Its tests run in the reviewer step (PROGRESS.md).

| § | Row(s) | Verdict |
|---|---|---|
| 1 | gamma, standard_gamma | BUILT NUM-PROB-069…072 (Gamma PDF/CDF/Quantile/MeanVariance) |
| 1 | beta · standard_t · f | BUILT 073…076 · 077…080 · 081…084 (FDistribution*) |
| 1 | binomial · poisson · hypergeometric quantile | BUILT 085 · 086 · 098 |
| 1 | weibull · standard_cauchy | BUILT 087…090 · 091…093 (Cauchy has no mean or variance) |
| 1 | multivariate_normal | BUILT 094/095 PDF/LogPDF; it takes the caller's Cholesky factor (D-083) |
| 1 | chisquare | BUILT 096/097 (PDF, MeanVariance); the CDF now uses NUM-SPEC-001/002 |
| 1 | negative_binomial | BUILT 099…102 (failures convention) |
| 1 | gumbel, laplace, logistic, pareto, rayleigh, triangular, power | BUILT 103…130. Pareto is the first kind; numpy's pareto is Lomax = Pareto I − 1, noted in the doc comment |
| 1 | multinomial · dirichlet · multivariate_hypergeometric | BUILT 131 · 132 · 133; mean and variance for the first two are 149/150 |
| 1 | logseries · wald · zipf | BUILT 134…137 · 138…141 · 144…146 and 147 (Zipf quantile) |
| 1 | vonmises | PARTLY BUILT: PDF 142, CDF 143, quantile 148. Mean and variance need BesselI1, which is outside D1's list |
| 1 | noncentral_chisquare · noncentral_f | BUILT 151…153 · 154/155 (noncentral F has CDF and MeanVariance only; λ ≤ 1e5) |
| 2 | all 17 NaN-skipping rows | DROPPED (D2) |
| 3 | fft, ifft, rfft, irfft, fftfreq, rfftfreq, fftshift, ifftshift | BUILT NUM-SIG-001…008 |
| 3 | bartlett, blackman, hamming, hanning, kaiser | BUILT NUM-SIG-012, 011, 010, 009, 013 |
| 3 | convolve · correlate | BUILT NUM-SIG-014 · 015 (full/same/valid) |
| 3 | hfft, ihfft | DROPPED (not in D1's list) |
| 4 | eig, eigvals | BUILT NUM-LINALG-169 MatEigenvalues, 191 MatEigenvectors |
| 4 | lstsq · svd · slogdet · diagonal · polyvander | BUILT 173 LeastSquaresMinNorm · 189 SVDFull · 174 MatLogDet · 186 MatDiagonal · 170 Vandermonde |
| 4 | norm · cond | BUILT 175…179 (nuclear, −1, −∞, −2, vector ord 0) · 180…182 (exact 1, ∞, Frobenius) |
| 4 | matrix_power · multi_dot · vecmat | BUILT 183 MatPowerInteger · 184 MatMulChain · 185 VecMatMul |
| 4 | polyfit | BUILT 171 PolynomialFit, 190 PolynomialFitCovariance (num/linalg, D-082) |
| 4 | (D8 extras) matrix_rank, pinv, kron | REUSED NUM-LINALG-099, 161, 121 |
| 5 | polyadd, polysub, polydiv, polypow, polyder, polyint | BUILT NUM-ALGEB-011…016, 028 PolynomialIntegralOrder |
| 5 | polyfromroots, polyvalfromroots, polycompanion, polytrim, polyval2d/polygrid2d | BUILT 017…021, 027 |
| 5 | np.poly | BUILT through 017 (roots input); REUSED NUM-LINALG-146…148 for a matrix |
| 5 | roots by companion eigenvalues | BUILT NUM-LINALG-172 PolynomialRootsCompanion (D-082) |
| 6 | Chebyshev, Legendre, Laguerre, Hermite, HermiteE core (5 rows) | PARTLY BUILT: evaluation NUM-SPEC-009…014. Coefficient arithmetic (add, mul, der, int, fit, vander, companion, trim, line, conversions) DROPPED (D1: evaluation and roots only) |
| 6 | Legendre, Laguerre, Hermite, HermiteE quadrature | BUILT NUM-SPEC-015, 018, 016, 017 (nodes = roots, Golub–Welsch) |
| 6 | Chebyshev points and quadrature | PARTLY BUILT: gauss 019/020. chebpts and chebinterpolate DROPPED (D1) |
| 6 | polyutils | DROPPED (D1) |
| 7 | sum, prod, cumprod, diff, argmin, argmax | BUILT ALG-RED-001, 004, 006, 007, 008, 009 (also CumSum 005, NeumaierAdd 002, TwoSum 003, AllFinite 010) |
| 7 | argpartition · multiply · divide | BUILT ALG-SELECT-005 · NUM-LINALG-187 VecHadamard · 188 VecDivElementwise |
| 8 | cov · corrcoef | BUILT NUM-STATS-038…040 · 041 |
| 8 | percentile, quantile | PARTLY BUILT: 042 Percentiles (several p, one sort). lower/higher/nearest/midpoint DROPPED (D4). Weights NEEDS SOURCE |
| 8 | histogram | PARTLY BUILT: 044 FrequencyDensity. Weights NEEDS SOURCE. Automatic bins via 043 |
| 8 | histogram_bin_edges · histogram2d | BUILT 043 · 045; also digitize, NUM-STATS-046 |
| 9 | bincount · lexsort · unique_all/counts/inverse · isin · setxor1d | BUILT ALG-HASH-014/015 · ALG-SORT-011 · 012/013/014 · ALG-HASH-016 · 017 |
| 10 | round, floor_divide/divmod, mod, sign, heaviside | BUILT NUM-ALGEB-022…026 |
| 11 | isclose · allclose · spacing · finfo | BUILT NUM-FLT-003 (symmetric, D3) · 004 · 005 Ulp · 001/002 constants; nan_to_num is 006 ReplaceNonFinite (D1 keeps it) |
| 12 | trapezoid · gradient · interp | BUILT NUM-CALC-001/002 · 003 · 004 (ErrOutOfRange outside the table, D-085) |
| 13 | i0 · sinc | BUILT NUM-SPEC-006 · 007 (also 001…005 and 008: incomplete gamma and beta, beta, zeta) |
| 14 | gcd · lcm | BUILT NUM-INTEG-004 · 005 (also ExtendedGCD 006) |
| 15 | logaddexp · logaddexp2 | BUILT NUM-LOGS-002 · 003 (also LogSumExp 004) |
| 16 | unwrap | BUILT NUM-TRIG-032 UnwrapAngles |

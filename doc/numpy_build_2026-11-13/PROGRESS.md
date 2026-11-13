# Progress (re-runnable: skip what is marked finished)

Builders append rows: `ID | Function | file | GAPS row | BUILT` · `GAPS row | REUSED <ID>` ·
`GAPS row | NEEDS SOURCE (why)` · `GAPS row | DROPPED (D2/D4)`. The lead marks stages.

## Stages
- [x] Lead: brief, baseline, taxonomy rows reserved (planned)
- [x] W1 B1 algo/reduce + num/integer + algo sets
- [x] W1 B2 num/prob
- [x] W1 B3 num/signal
- [x] W1 B4 num/linalg + num/algebra
- [x] W1 B5 num/special + num/float
- [x] W1 B6 num/stats + num/calculus + num/logs + num/trig
- [x] W2 C1 algo/** · C2 num/** · C3 fin/life/unit · C4 med/**
- [x] Lead: INDEX, taxonomy, doc.go headers, GAPS verdicts, DECISIONS, CHANGELOG, SCOPE
- [x] Reviewer: finished 2026-11-13 — 71 ok + 5 FAIL (15 findings) → owners fixed → rerun of the 5: all ok, plain and -race; index 1,505|1,505|1,505

## Rows
B5: R4 API READY
ALG-RED-001 | Sum | algo/reduce/sum.go | D5; §7 sum | BUILT
ALG-RED-002 | NeumaierAdd | algo/reduce/neumaier_add.go | D5 (R2) | BUILT
ALG-RED-003 | TwoSum | algo/reduce/two_sum.go | D5 (R2) | BUILT
ALG-RED-004 | Prod | algo/reduce/prod.go | D5; §7 prod | BUILT
ALG-RED-005 | CumSum | algo/reduce/cum_sum.go | D5 | BUILT
ALG-RED-006 | CumProd | algo/reduce/cum_prod.go | D5; §7 cumprod (+ cumulative_prod) | BUILT
ALG-RED-007 | Diff | algo/reduce/diff.go | D5; §7 diff (+ ediff1d) | BUILT
ALG-RED-008 | ArgMin | algo/reduce/arg_min.go | D5; §7 argmin | BUILT
ALG-RED-009 | ArgMax | algo/reduce/arg_max.go | D5; §7 argmax | BUILT
ALG-RED-010 | AllFinite | algo/reduce/all_finite.go | D5 (R2) | BUILT
NUM-INTEG-004 | GCD | num/integer/gcd.go | D6; §14 gcd | BUILT
NUM-INTEG-005 | LCM | num/integer/lcm.go | D6; §14 lcm | BUILT
NUM-INTEG-006 | ExtendedGCD | num/integer/extended_gcd.go | D6 (R3) | BUILT
ALG-SELECT-005 | KthSmallestIndex | algo/selection/kth_smallest_index.go | §7 argpartition | BUILT
ALG-SORT-011 | LexSortIndices | algo/sorting/lex_sort_indices.go | §9 lexsort | BUILT
ALG-SORT-012 | UniqueAll | algo/sorting/unique_all.go | §9 unique_all | BUILT
ALG-SORT-013 | UniqueCounts | algo/sorting/unique_counts.go | §9 unique_counts | BUILT
ALG-SORT-014 | UniqueInverse | algo/sorting/unique_inverse.go | §9 unique_inverse | BUILT
ALG-HASH-014 | BinCount | algo/hashing/bin_count.go | §9 bincount | BUILT
ALG-HASH-015 | BinCountWeighted | algo/hashing/bin_count_weighted.go | §9 bincount (weights) | BUILT
ALG-HASH-016 | IsIn | algo/hashing/is_in.go | §9 isin | BUILT
ALG-HASH-017 | SetSymmetricDifference | algo/hashing/set_symmetric_difference.go | §9 setxor1d | BUILT
NUM-SPEC-001 | RegularizedGammaP | num/special/regularized_gamma.go | §13 (+§1 needs) incomplete gamma | BUILT
NUM-SPEC-002 | RegularizedGammaQ | num/special/regularized_gamma.go | §13 (+§1 needs) incomplete gamma | BUILT
NUM-SPEC-003 | Beta | num/special/beta.go | D1 special beta | BUILT
NUM-SPEC-004 | LogBeta | num/special/beta.go | D1 special beta | BUILT
NUM-SPEC-005 | RegularizedBeta | num/special/regularized_beta.go | §13 (+§1 needs) incomplete beta | BUILT
NUM-SPEC-006 | BesselI0 | num/special/bessel_i0.go | §13 i0 | BUILT
NUM-SPEC-007 | Sinc | num/special/sinc.go | §13 sinc | BUILT
NUM-SPEC-008 | RiemannZeta | num/special/riemann_zeta.go | §13 zeta (for §1) | BUILT
NUM-SPEC-009 | LegendreP | num/special/legendre_p.go | §6 Legendre core (val) | BUILT
NUM-SPEC-010 | ChebyshevT | num/special/chebyshev_t.go | §6 Chebyshev core (val) | BUILT
NUM-SPEC-011 | ChebyshevU | num/special/chebyshev_u.go | §6 Chebyshev (second kind val) | BUILT
NUM-SPEC-012 | HermiteH | num/special/hermite_h.go | §6 Hermite core (val) | BUILT
NUM-SPEC-013 | HermiteHe | num/special/hermite_he.go | §6 HermiteE core (val) | BUILT
NUM-SPEC-014 | LaguerreL | num/special/laguerre_l.go | §6 Laguerre core (val) | BUILT
NUM-SPEC-015 | GaussLegendre | num/special/gauss_quadrature.go | §6 leggauss/legweight + legroots | BUILT
NUM-SPEC-016 | GaussHermite | num/special/gauss_quadrature.go | §6 hermgauss + hermroots | BUILT
NUM-SPEC-017 | GaussHermiteE | num/special/gauss_quadrature.go | §6 hermegauss + hermeroots | BUILT
NUM-SPEC-018 | GaussLaguerre | num/special/gauss_quadrature.go | §6 laggauss + lagroots | BUILT
NUM-SPEC-019 | GaussChebyshevT | num/special/gauss_quadrature.go | §6 chebgauss + chebroots | BUILT
NUM-SPEC-020 | GaussChebyshevU | num/special/gauss_quadrature.go | §6 Chebyshev U roots/rule | BUILT
NUM-FLT-001 | MachineEpsilon | num/float/machine_constants.go | §11 finfo | BUILT
NUM-FLT-002 | SmallestNormal | num/float/machine_constants.go | §11 finfo | BUILT
NUM-FLT-003 | IsClose | num/float/is_close.go | §11 isclose (D3 symmetric) | BUILT
NUM-FLT-004 | AllClose | num/float/is_close.go | §11 allclose | BUILT
NUM-FLT-005 | Ulp | num/float/ulp.go | §11 spacing | BUILT
NUM-FLT-006 | ReplaceNonFinite | num/float/replace_non_finite.go | D1 nan_to_num | BUILT
D1 float nextafter | REUSED Go math.Nextafter (D-012)
D1 float copysign | REUSED Go math.Copysign (D-012)
D1 special gamma, lgamma | REUSED Go math.Gamma, math.Lgamma (D-012)
§6 coefficient-series arithmetic per family (add/sub/mul/div/der/int/fit/vander/companion/trim, polyutils, chebpts/chebinterpolate) | NOT BUILT: D1 reduces §6 to evaluation and roots; lead to mark OUT OF SCOPE or route to num/algebra
NUM-LOGS-002 | LogAddExp | num/logs/log_add_exp.go | §15 logaddexp | BUILT
NUM-LOGS-003 | LogAddExp2 | num/logs/log_add_exp.go | §15 logaddexp2 | BUILT
NUM-LOGS-004 | LogSumExp | num/logs/log_sum_exp.go | §15 LogSumExp | BUILT
NUM-TRIG-032 | UnwrapAngles | num/trig/unwrap_angles.go | §16 unwrap | BUILT
NUM-CALC-001 | Trapezoid | num/calculus/trapezoid.go | §12 trapezoid (x given) | BUILT
NUM-CALC-002 | TrapezoidUniform | num/calculus/trapezoid.go | §12 trapezoid (dx) | BUILT
NUM-CALC-003 | Gradient | num/calculus/gradient.go | §12 gradient | BUILT
NUM-CALC-004 | Interp | num/calculus/interp.go | §12 interp (outside the table = ErrOutOfRange, not numpy's clamp) | BUILT
NUM-STATS-038 | SampleCovariance | num/stats/covariance.go | §8 cov | BUILT
NUM-STATS-039 | PopulationCovariance | num/stats/covariance.go | §8 cov (bias=True) | BUILT
NUM-STATS-040 | CovarianceMatrix | num/stats/covariance.go | §8 cov (k variables, row-major []float64 as num/linalg) | BUILT
NUM-STATS-041 | CorrelationMatrix | num/stats/correlation_matrix.go | §8 corrcoef | BUILT
NUM-STATS-042 | Percentiles | num/stats/percentiles.go | §8 percentile/quantile, several p in one sort | BUILT
NUM-STATS-043 | HistogramBinEdges | num/stats/histogram_bin_edges.go | §8 histogram_bin_edges (Sturges, Scott, FD, Doane) + histogram automatic bins | BUILT
NUM-STATS-044 | FrequencyDensity | num/stats/frequency_density.go | §8 histogram density | BUILT
NUM-STATS-045 | FrequencyDistribution2D | num/stats/frequency_distribution_2d.go | §8 histogram2d | BUILT
NUM-STATS-046 | Digitize | num/stats/digitize.go | §8 digitize | BUILT
§8 percentile/quantile lower, higher, nearest, midpoint | DROPPED (D4)
§8 percentile/quantile weights | NEEDS SOURCE (Hyndman & Fan 1996 define unweighted sample quantiles only; no primary source for a weighted H&F family found)
§8 histogram weights | NEEDS SOURCE (weighted bin sums: no primary statistical source cited in GAPS; trivial once one is chosen)
§2 nanmean, nanvar, nanstd, nanmedian, nanpercentile, nanquantile (num/stats) | DROPPED (D2)
NUM-LINALG-169 | MatEigenvalues | num/linalg/eigen_general.go | §4 eig/eigvals (general) | BUILT
NUM-LINALG-170 | Vandermonde | num/linalg/vandermonde.go | §4 polyvander | BUILT
NUM-LINALG-171 | PolynomialFit | num/linalg/polynomial_fit.go | §4 polyfit (any degree, weights; covariance not built) | BUILT
NUM-LINALG-172 | PolynomialRootsCompanion | num/linalg/polynomial_roots_companion.go | §5 roots via companion eigenvalues (R1) | BUILT
NUM-LINALG-173 | LeastSquaresMinNorm | num/linalg/least_squares_min_norm.go | §4 lstsq rank-deficient | BUILT
NUM-LINALG-174 | MatLogDet | num/linalg/mat_log_det.go | §4 slogdet | BUILT
NUM-LINALG-175 | MatNormNuclear | num/linalg/mat_norms_extra.go | §4 norm nuclear | BUILT
NUM-LINALG-176 | MatNormNeg1 | num/linalg/mat_norms_extra.go | §4 norm ord −1 | BUILT
NUM-LINALG-177 | MatNormNegInf | num/linalg/mat_norms_extra.go | §4 norm ord −inf | BUILT
NUM-LINALG-178 | MatNormNeg2 | num/linalg/mat_norms_extra.go | §4 norm ord −2 | BUILT
NUM-LINALG-179 | VecNorm0 | num/linalg/mat_norms_extra.go | §4 vector norm ord 0 | BUILT
NUM-LINALG-180 | MatCondition1 | num/linalg/mat_condition_exact.go | §4 cond 1 | BUILT
NUM-LINALG-181 | MatConditionInf | num/linalg/mat_condition_exact.go | §4 cond ∞ | BUILT
NUM-LINALG-182 | MatConditionFrobenius | num/linalg/mat_condition_exact.go | §4 cond fro | BUILT
NUM-LINALG-183 | MatPowerInteger | num/linalg/mat_power_integer.go | §4 matrix_power negative k | BUILT
NUM-LINALG-184 | MatMulChain | num/linalg/mat_mul_chain.go | §4 multi_dot (via ALG-DP-019) | BUILT
NUM-LINALG-185 | VecMatMul | num/linalg/vec_mat_mul.go | §4 vecmat | BUILT
NUM-LINALG-186 | MatDiagonal | num/linalg/mat_diagonal.go | §4 diagonal | BUILT
NUM-LINALG-187 | VecHadamard | num/linalg/vec_elementwise.go | §7 multiply | BUILT
NUM-LINALG-188 | VecDivElementwise | num/linalg/vec_elementwise.go | §7 divide | BUILT
NUM-LINALG-189 | SVDFull | num/linalg/svd_full.go | §4 svd full_matrices | BUILT
NUM-ALGEB-011 | PolynomialAdd | num/algebra/polynomial_add.go | §5 polyadd | BUILT
NUM-ALGEB-012 | PolynomialSubtract | num/algebra/polynomial_add.go | §5 polysub | BUILT
NUM-ALGEB-013 | PolynomialDivide | num/algebra/polynomial_divide.go | §5 polydiv | BUILT
NUM-ALGEB-014 | PolynomialPower | num/algebra/polynomial_power.go | §5 polypow | BUILT
NUM-ALGEB-015 | PolynomialDerivative | num/algebra/polynomial_derivative.go | §5 polyder | BUILT
NUM-ALGEB-016 | PolynomialIntegral | num/algebra/polynomial_integral.go | §5 polyint (one order per call) | BUILT
NUM-ALGEB-017 | PolynomialFromRoots | num/algebra/polynomial_from_roots.go | §5 polyfromroots | BUILT
NUM-ALGEB-018 | PolynomialValueFromRoots | num/algebra/polynomial_from_roots.go | §5 polyvalfromroots | BUILT
NUM-ALGEB-019 | PolynomialCompanion | num/algebra/polynomial_companion.go | §5 polycompanion | BUILT
NUM-ALGEB-020 | PolynomialTrim | num/algebra/polynomial_trim.go | §5 polytrim | BUILT
NUM-ALGEB-021 | PolynomialEvaluate2D | num/algebra/polynomial_evaluate_2d.go | §5 polyval2d (polygrid2d not built) | BUILT
NUM-ALGEB-022 | RoundDecimals | num/algebra/round_decimals.go | §10 round/around | BUILT
NUM-ALGEB-023 | FloorDivMod | num/algebra/floor_div_mod.go | §10 floor_divide, divmod | BUILT
NUM-ALGEB-024 | FloorMod | num/algebra/floor_div_mod.go | §10 mod/remainder | BUILT
NUM-ALGEB-025 | Sign | num/algebra/sign.go | §10 sign | BUILT
NUM-ALGEB-026 | Heaviside | num/algebra/sign.go | §10 heaviside | BUILT
§4 matrix_rank | REUSED NUM-LINALG-099
§4 pinv | REUSED NUM-LINALG-161
§4 kron | REUSED NUM-LINALG-121
§5 polymul | REUSED NUM-ALGEB-003
§5 np.poly (from roots input) | REUSED NUM-ALGEB-017 (+ NUM-LINALG-146/147/148 for a matrix)
NUM-LINALG-190 | PolynomialFitCovariance | num/linalg/polynomial_fit_covariance.go | §4 polyfit cov=True (B4 pass 2) | BUILT
NUM-LINALG-191 | MatEigenvectors | num/linalg/eigenvectors_general.go | §4 eig eigenvectors (general, inverse iteration on Hessenberg; B4 pass 2) | BUILT
NUM-ALGEB-027 | PolynomialEvaluateGrid2D | num/algebra/polynomial_evaluate_grid_2d.go | §5 polygrid2d (B4 pass 2) | BUILT
NUM-ALGEB-028 | PolynomialIntegralOrder | num/algebra/polynomial_integral_order.go | §5 polyint m, lbnd, k (B4 pass 2) | BUILT
NUM-SIG-001 | FFT | num/signal/fft.go | §3 fft | BUILT
NUM-SIG-002 | IFFT | num/signal/fft.go | §3 ifft | BUILT
NUM-SIG-003 | RFFT | num/signal/rfft.go | §3 rfft | BUILT
NUM-SIG-004 | IRFFT | num/signal/rfft.go | §3 irfft | BUILT
NUM-SIG-005 | FFTFreq | num/signal/fft_freq.go | §3 fftfreq | BUILT
NUM-SIG-006 | RFFTFreq | num/signal/fft_freq.go | §3 rfftfreq | BUILT
NUM-SIG-007 | FFTShift | num/signal/fft_shift.go | §3 fftshift | BUILT
NUM-SIG-008 | IFFTShift | num/signal/fft_shift.go | §3 ifftshift | BUILT
NUM-SIG-009 | Hanning | num/signal/cosine_window.go | §3 hanning | BUILT
NUM-SIG-010 | Hamming | num/signal/cosine_window.go | §3 hamming | BUILT
NUM-SIG-011 | Blackman | num/signal/cosine_window.go | §3 blackman | BUILT
NUM-SIG-012 | Bartlett | num/signal/bartlett.go | §3 bartlett | BUILT
NUM-SIG-013 | Kaiser | num/signal/kaiser.go | §3 kaiser | BUILT
NUM-SIG-014 | Convolve | num/signal/convolve.go | §3 convolve (modes full/same/valid; FFT path at min(n,m) ≥ 64) | BUILT
NUM-SIG-015 | Correlate | num/signal/convolve.go | §3 correlate | BUILT
§3 hfft, ihfft | NOT BUILT (low priority, time; next IDs NUM-SIG-016/017)
NUM-PROB-069 | GammaPDF | num/prob/gamma_pdf.go | §1 gamma | BUILT
NUM-PROB-070 | GammaCDF | num/prob/gamma_cdf.go | §1 gamma | BUILT
NUM-PROB-071 | GammaQuantile | num/prob/gamma_quantile.go | §1 gamma | BUILT
NUM-PROB-072 | GammaMeanVariance | num/prob/gamma_mean_variance.go | §1 gamma | BUILT
NUM-PROB-073 | BetaPDF | num/prob/beta_pdf.go | §1 beta | BUILT
NUM-PROB-074 | BetaCDF | num/prob/beta_cdf.go | §1 beta | BUILT
NUM-PROB-075 | BetaQuantile | num/prob/beta_quantile.go | §1 beta | BUILT
NUM-PROB-076 | BetaMeanVariance | num/prob/beta_mean_variance.go | §1 beta | BUILT
NUM-PROB-077 | StudentTPDF | num/prob/student_t_pdf.go | §1 standard_t | BUILT
NUM-PROB-078 | StudentTCDF | num/prob/student_t_cdf.go | §1 standard_t | BUILT
NUM-PROB-079 | StudentTQuantile | num/prob/student_t_quantile.go | §1 standard_t | BUILT
NUM-PROB-080 | StudentTMeanVariance | num/prob/student_t_mean_variance.go | §1 standard_t | BUILT
NUM-PROB-081 | FDistributionPDF | num/prob/f_distribution_pdf.go | §1 f | BUILT
NUM-PROB-082 | FDistributionCDF | num/prob/f_distribution_cdf.go | §1 f | BUILT
NUM-PROB-083 | FDistributionQuantile | num/prob/f_distribution_quantile.go | §1 f | BUILT
NUM-PROB-084 | FDistributionMeanVariance | num/prob/f_distribution_mean_variance.go | §1 f | BUILT
NUM-PROB-085 | BinomialQuantile | num/prob/binomial_quantile.go | §1 binomial | BUILT
NUM-PROB-086 | PoissonQuantile | num/prob/poisson_quantile.go | §1 poisson | BUILT
NUM-PROB-087 | WeibullPDF | num/prob/weibull_pdf.go | §1 weibull | BUILT
NUM-PROB-088 | WeibullCDF | num/prob/weibull_cdf.go | §1 weibull | BUILT
NUM-PROB-089 | WeibullQuantile | num/prob/weibull_quantile.go | §1 weibull | BUILT
NUM-PROB-090 | WeibullMeanVariance | num/prob/weibull_mean_variance.go | §1 weibull | BUILT
NUM-PROB-091 | CauchyPDF | num/prob/cauchy_pdf.go | §1 standard_cauchy (location-scale) | BUILT
NUM-PROB-092 | CauchyCDF | num/prob/cauchy_cdf.go | §1 standard_cauchy | BUILT
NUM-PROB-093 | CauchyQuantile | num/prob/cauchy_quantile.go | §1 standard_cauchy (no MeanVariance: undefined) | BUILT
§1 multivariate_normal | NEEDS DECISION (D-009/R1: num/prob may not import num/linalg; options: (a) take a caller-supplied Cholesky factor L, no import, recommended; (b) move MVN PDF to num/linalg or num/stats; lost: a private Cholesky copy in prob, two homes for one algorithm)
§1 chi-square private regularizedGamma (chi_square_cdf.go) | LEFT for wave 2: not identical by construction to special.RegularizedGammaP/Q (different evaluation), not replaced
§1 hypergeometric quantile, negative binomial failures + quantile, chi-square PDF/MeanVariance, multinomial, dirichlet, gumbel, laplace, logistic, pareto, power, rayleigh, triangular, logseries, multivariate_hypergeometric, noncentral chi-square/F, vonmises, wald, zipf | NOT STARTED (B2 out of time)
NUM-PROB-094 | MultivariateNormalPDF | num/prob/multivariate_normal_pdf.go | §1 multivariate_normal (D-083: takes Cholesky L) | BUILT
NUM-PROB-095 | MultivariateNormalLogPDF | num/prob/multivariate_normal_log_pdf.go | §1 multivariate_normal logpdf | BUILT
NUM-PROB-096 | ChiSquarePDF | num/prob/chi_square_pdf.go | §1 chisquare | BUILT
NUM-PROB-097 | ChiSquareMeanVariance | num/prob/chi_square_mean_variance.go | §1 chisquare | BUILT
NUM-PROB-098 | HypergeometricQuantile | num/prob/hypergeometric_quantile.go | §1 hypergeometric | BUILT
NUM-PROB-099 | NegativeBinomialFailuresPMF | num/prob/negative_binomial_failures_pmf.go | §1 negative_binomial | BUILT
NUM-PROB-100 | NegativeBinomialFailuresCDF | num/prob/negative_binomial_failures_cdf.go | §1 negative_binomial | BUILT
NUM-PROB-101 | NegativeBinomialFailuresQuantile | num/prob/negative_binomial_failures_quantile.go | §1 negative_binomial | BUILT
NUM-PROB-102 | NegativeBinomialFailuresMeanVariance | num/prob/negative_binomial_failures_mean_variance.go | §1 negative_binomial | BUILT
NUM-PROB-103 | GumbelPDF | num/prob/gumbel_pdf.go | §1 gumbel | BUILT
NUM-PROB-104 | GumbelCDF | num/prob/gumbel_cdf.go | §1 gumbel | BUILT
NUM-PROB-105 | GumbelQuantile | num/prob/gumbel_quantile.go | §1 gumbel | BUILT
NUM-PROB-106 | GumbelMeanVariance | num/prob/gumbel_mean_variance.go | §1 gumbel | BUILT
NUM-PROB-107 | LaplacePDF | num/prob/laplace_pdf.go | §1 laplace | BUILT
NUM-PROB-108 | LaplaceCDF | num/prob/laplace_cdf.go | §1 laplace | BUILT
NUM-PROB-109 | LaplaceQuantile | num/prob/laplace_quantile.go | §1 laplace | BUILT
NUM-PROB-110 | LaplaceMeanVariance | num/prob/laplace_mean_variance.go | §1 laplace | BUILT
NUM-PROB-111 | LogisticPDF | num/prob/logistic_pdf.go | §1 logistic | BUILT
NUM-PROB-112 | LogisticCDF | num/prob/logistic_cdf.go | §1 logistic | BUILT
NUM-PROB-113 | LogisticQuantile | num/prob/logistic_quantile.go | §1 logistic | BUILT
NUM-PROB-114 | LogisticMeanVariance | num/prob/logistic_mean_variance.go | §1 logistic | BUILT
NUM-PROB-115 | ParetoPDF | num/prob/pareto_pdf.go | §1 pareto (first kind; numpy's pareto is Lomax = Pareto − 1, note) | BUILT
NUM-PROB-116 | ParetoCDF | num/prob/pareto_cdf.go | §1 pareto | BUILT
NUM-PROB-117 | ParetoQuantile | num/prob/pareto_quantile.go | §1 pareto | BUILT
NUM-PROB-118 | ParetoMeanVariance | num/prob/pareto_mean_variance.go | §1 pareto | BUILT
NUM-PROB-119 | RayleighPDF | num/prob/rayleigh_pdf.go | §1 rayleigh | BUILT
NUM-PROB-120 | RayleighCDF | num/prob/rayleigh_cdf.go | §1 rayleigh | BUILT
NUM-PROB-121 | RayleighQuantile | num/prob/rayleigh_quantile.go | §1 rayleigh | BUILT
NUM-PROB-122 | RayleighMeanVariance | num/prob/rayleigh_mean_variance.go | §1 rayleigh | BUILT
NUM-PROB-123 | TriangularPDF | num/prob/triangular_pdf.go | §1 triangular | BUILT
NUM-PROB-124 | TriangularCDF | num/prob/triangular_cdf.go | §1 triangular | BUILT
NUM-PROB-125 | TriangularQuantile | num/prob/triangular_quantile.go | §1 triangular | BUILT
NUM-PROB-126 | TriangularMeanVariance | num/prob/triangular_mean_variance.go | §1 triangular | BUILT
NUM-PROB-127 | PowerFunctionPDF | num/prob/power_function_pdf.go | §1 power | BUILT
NUM-PROB-128 | PowerFunctionCDF | num/prob/power_function_cdf.go | §1 power | BUILT
NUM-PROB-129 | PowerFunctionQuantile | num/prob/power_function_quantile.go | §1 power | BUILT
NUM-PROB-130 | PowerFunctionMeanVariance | num/prob/power_function_mean_variance.go | §1 power | BUILT
§1 multinomial, dirichlet, multivariate_hypergeometric, logseries, wald, vonmises, zipf, noncentral_chisquare, noncentral_f | NOT STARTED (B2 pass 2 out of time; next ID NUM-PROB-131)
B5b: LOADER READY
NUM-PROB-131 | MultinomialPMF | num/prob/multinomial_pmf.go | §1 multinomial | BUILT
NUM-PROB-132 | DirichletPDF | num/prob/dirichlet_pdf.go | §1 dirichlet | BUILT
NUM-PROB-133 | MultivariateHypergeometricPMF | num/prob/multivariate_hypergeometric_pmf.go | §1 multivariate_hypergeometric | BUILT
NUM-PROB-134 | LogSeriesPMF | num/prob/log_series_pmf.go | §1 logseries | BUILT
NUM-PROB-135 | LogSeriesCDF | num/prob/log_series_cdf.go | §1 logseries | BUILT
NUM-PROB-136 | LogSeriesQuantile | num/prob/log_series_quantile.go | §1 logseries | BUILT
NUM-PROB-137 | LogSeriesMeanVariance | num/prob/log_series_mean_variance.go | §1 logseries | BUILT
NUM-PROB-138 | WaldPDF | num/prob/wald_pdf.go | §1 wald | BUILT
NUM-PROB-139 | WaldCDF | num/prob/wald_cdf.go | §1 wald | BUILT
NUM-PROB-140 | WaldQuantile | num/prob/wald_quantile.go | §1 wald | BUILT
NUM-PROB-141 | WaldMeanVariance | num/prob/wald_mean_variance.go | §1 wald | BUILT
NUM-PROB-142 | VonMisesPDF | num/prob/von_mises_pdf.go | §1 vonmises | BUILT
NUM-PROB-143 | VonMisesCDF | num/prob/von_mises_cdf.go | §1 vonmises (Hill 1977 series branch for all κ ≤ 1e6; no quantile/MeanVariance: circular variance needs BesselI1, absent in num/special) | BUILT
NUM-PROB-144 | ZipfPMF | num/prob/zipf_pmf.go | §1 zipf | BUILT
NUM-PROB-145 | ZipfCDF | num/prob/zipf_cdf.go | §1 zipf (no quantile built) | BUILT
NUM-PROB-146 | ZipfMeanVariance | num/prob/zipf_mean_variance.go | §1 zipf | BUILT
§1 noncentral_chisquare, noncentral_f; multinomial/dirichlet MeanVariance (vector results), von Mises quantile/moments, Zipf quantile | NOT STARTED (B2 pass 3 out of time; next ID NUM-PROB-147)
num/prob wave-2 consolidation | isFinite → reduce.AllFinite (52 files), compensatedSum.add → reduce.NeumaierAdd | DONE (B2 pass 3)
num/prob/chi_square_cdf.go Part 3 | NOT DONE: "B5b: LOADER READY" absent from PROGRESS.md
num/prob/chi_square_cdf.go Part 3 | DONE (B2 pass 4): ChiSquareCDF/Survival/Quantile use special.RegularizedGammaP/Q; old private gamma kept verbatim as test oracle (chi_square_cdf_private_gamma_test.go, TestChiSquareCDF_matchesPrivateGamma; 1e-14 in the bulk, documented wider allowance: a ≤ 15 far tails 4·E·2⁻⁵³, large df √a·2⁻⁵²); quantile density now ChiSquarePDF
NUM-PROB-147 | ZipfQuantile | num/prob/zipf_quantile.go | §1 zipf | BUILT
NUM-PROB-148 | VonMisesQuantile | num/prob/von_mises_quantile.go | §1 vonmises | BUILT
VonMises mean/variance | NOT BUILT (needs BesselI1, outside D1's list)
NUM-PROB-149 | MultinomialMeanVariance | num/prob/multinomial_mean_variance.go | §1 multinomial | BUILT
NUM-PROB-150 | DirichletMeanVariance | num/prob/dirichlet_mean_variance.go | §1 dirichlet | BUILT
NUM-PROB-151 | NoncentralChiSquarePDF | num/prob/noncentral_chi_square_pdf.go | §1 noncentral_chisquare | BUILT
NUM-PROB-152 | NoncentralChiSquareCDF | num/prob/noncentral_chi_square_cdf.go | §1 noncentral_chisquare | BUILT
NUM-PROB-153 | NoncentralChiSquareMeanVariance | num/prob/noncentral_chi_square_mean_variance.go | §1 noncentral_chisquare | BUILT
NUM-PROB-154 | NoncentralFCDF | num/prob/noncentral_f_cdf.go | §1 noncentral_f | BUILT
NUM-PROB-155 | NoncentralFMeanVariance | num/prob/noncentral_f_mean_variance.go | §1 noncentral_f | BUILT

// Package prob holds the rules of probability, odds conversions, expectation bounds, and the PMF or PDF,
// CDF, mean and variance, and quantile of the common discrete and continuous distributions.
//
// ID short form: PROB (IDs are NUM-PROB-NNN).
//
// Conventions (doc/DECISIONS.md D-013): plain functions per distribution, no Distribution interface.
// Probabilities are float64 in [0, 1]; counts are int. A continuous variable carries the caller's unit,
// its location and scale share that unit, and a rate is per that unit. Mean and variance come back as one
// pair from <Dist>MeanVariance. The geometric and negative binomial distributions count TRIALS (support
// 1, 2, … and r, r+1, …), and their names say so; the failures count is X − 1 or X − r, and the
// NegativeBinomialFailures* functions give it directly (support 0, 1, 2, …, the numpy/scipy convention).
// The multivariate normal takes the caller's lower-triangular Cholesky factor L of Σ (row-major, as
// NUM-LINALG-152 Cholesky returns it): num/prob does not import num/linalg (D-009, D-083).
//
// Values, not errors: a PMF at an integer outside its support is 0; a CDF below the support is 0 and at
// or above its top is 1.
//
// Distributions: where slices of probabilities must form a distribution, their sum must be 1 within
// n·2⁻⁵² (n = the number of probabilities), the first-order error bound of a recursive sum, doubled.
//
// Accuracy: binomial, negative binomial, hypergeometric and Poisson terms use Loader's saddle-point form
// (no factorial, power or ln Γ difference). Their CDFs sum the short tail with compensation and stop below
// 2⁻⁵³ of the sum, capped at 2²⁸ terms (ErrOutOfRange beyond). The incomplete beta and gamma functions now
// live in num/special (NUM-SPEC-001/002 RegularizedGammaP/Q, RegularizedBeta) and the gamma, beta, Student
// t and F CDFs use them, as do the chi-square CDF, survival and quantile (num/special's prefactor now has
// Loader's saddle-point form, so the private incomplete gamma here was retired; see
// doc/numpy_build_2026-11-13/consolidation.md). Normal tails use erfc on the correct side; the normal
// quantile is Wichura's PPND16.
//
// Every function is pure. Scalar functions allocate nothing when they succeed (an error allocates its
// wrapped message); the two that return slices append to a caller's dst. No function modifies its input.
// Errors are wrapped with the offending parameter; test with errors.Is.
//
// Functions:
//
//	NUM-PROB-001  ProbabilityFromCounts               favorable/total for equally likely outcomes
//	NUM-PROB-002  ComplementProbability               1 − P(A)
//	NUM-PROB-003  IndependentJointProbability         P(A)·P(B) for independent events
//	NUM-PROB-004  JointProbability                    P(A given B)·P(B), the multiplication rule
//	NUM-PROB-005  UnionProbability                    P(A) + P(B) − P(A ∩ B)
//	NUM-PROB-006  ConditionalProbability              P(A ∩ B)/P(B)
//	NUM-PROB-007  TotalProbability                    Σ P(A given Bᵢ)·P(Bᵢ) over a partition
//	NUM-PROB-008  BayesTheorem                        P(B given A)·P(A)/P(B)
//	NUM-PROB-009  BayesPosterior                      posteriors over a partition, appended to dst
//	NUM-PROB-010  PositivePredictiveValue             P(condition given a positive test) from Se, Sp, prevalence
//	NUM-PROB-011  NegativePredictiveValue             P(no condition given a negative test) from Se, Sp, prevalence
//	NUM-PROB-012  OddsFromProbability                 odds in favour p/(1 − p)
//	NUM-PROB-013  ProbabilityFromOdds                 probability o/(1 + o) from odds in favour
//	NUM-PROB-014  DecimalOddsFromProbability          fair decimal odds 1/p
//	NUM-PROB-015  ImpliedProbability                  probability 1/d implied by decimal odds
//	NUM-PROB-016  Overround                           booksum Σ 1/dᵢ minus 1 of one market
//	NUM-PROB-017  NormalizedImpliedProbabilities      implied probabilities with the overround removed proportionally
//	NUM-PROB-018  ExpectedValue                       E[X] of a discrete distribution
//	NUM-PROB-019  DiscreteVariance                    Var X of a discrete distribution, two-pass
//	NUM-PROB-020  ChebyshevBound                      Chebyshev's bound 1/k² on a deviation of k or more standard deviations
//	NUM-PROB-021  ChebyshevBoundEpsilon               Chebyshev's bound σ²/ε² on a deviation of ε or more
//	NUM-PROB-022  MarkovBound                         Markov's bound E[X]/a on P(X ≥ a)
//	NUM-PROB-023  BernoulliPMF                        Bernoulli P(X = k)
//	NUM-PROB-024  BernoulliCDF                        Bernoulli P(X ≤ k)
//	NUM-PROB-025  BernoulliMeanVariance               Bernoulli mean and variance
//	NUM-PROB-026  BinomialPMF                         binomial P(X = k), saddle-point form
//	NUM-PROB-027  BinomialCDF                         binomial P(X ≤ k), compensated tail sum
//	NUM-PROB-028  BinomialMeanVariance                binomial mean and variance
//	NUM-PROB-029  GeometricTrialsPMF                  geometric P(X = k), X = trials to the first success
//	NUM-PROB-030  GeometricTrialsCDF                  geometric P(X ≤ k), X = trials to the first success
//	NUM-PROB-031  GeometricTrialsMeanVariance         geometric mean and variance, trials convention
//	NUM-PROB-032  NegativeBinomialTrialsPMF           negative binomial P(X = x), X = trials to the r-th success
//	NUM-PROB-033  NegativeBinomialTrialsCDF           negative binomial P(X ≤ x), X = trials to the r-th success
//	NUM-PROB-034  NegativeBinomialTrialsMeanVariance  negative binomial mean and variance, trials convention
//	NUM-PROB-035  HypergeometricPMF                   hypergeometric P(X = k), saddle-point form
//	NUM-PROB-036  HypergeometricCDF                   hypergeometric P(X ≤ k), compensated tail sum
//	NUM-PROB-037  HypergeometricMeanVariance          hypergeometric mean and variance
//	NUM-PROB-038  DiscreteUniformPMF                  discrete uniform P(X = k) on a…b
//	NUM-PROB-039  DiscreteUniformCDF                  discrete uniform P(X ≤ k) on a…b
//	NUM-PROB-040  DiscreteUniformMeanVariance         discrete uniform mean and variance on a…b
//	NUM-PROB-041  PoissonPMF                          Poisson P(X = k), saddle-point form
//	NUM-PROB-042  PoissonCDF                          Poisson P(X ≤ k), compensated tail sum
//	NUM-PROB-043  PoissonMeanVariance                 Poisson mean and variance
//	NUM-PROB-044  UniformPDF                          continuous uniform density on [a, b]
//	NUM-PROB-045  UniformCDF                          continuous uniform P(X ≤ x) on [a, b]
//	NUM-PROB-046  UniformMeanVariance                 continuous uniform mean and variance on [a, b]
//	NUM-PROB-047  NormalPDF                           normal density
//	NUM-PROB-048  NormalCDF                           normal P(X ≤ x), accurate in the lower tail
//	NUM-PROB-049  NormalSurvival                      normal P(X > x), accurate in the upper tail
//	NUM-PROB-050  NormalIntervalProbability           normal P(lo ≤ X ≤ hi), tail-aware
//	NUM-PROB-051  LognormalPDF                        lognormal density
//	NUM-PROB-052  LognormalCDF                        lognormal P(X ≤ x)
//	NUM-PROB-053  LognormalMeanVariance               lognormal mean and variance
//	NUM-PROB-054  ExponentialPDF                      exponential density, rate form
//	NUM-PROB-055  ExponentialCDF                      exponential P(X ≤ x), rate form
//	NUM-PROB-056  ExponentialMeanVariance             exponential mean and variance, rate form
//	NUM-PROB-057  GeometricTrialsQuantile             least k with P(X ≤ k) ≥ u, trials convention
//	NUM-PROB-058  DiscreteUniformQuantile             least k in a…b with P(X ≤ k) ≥ u
//	NUM-PROB-059  UniformQuantile                     continuous uniform quantile a + u(b − a)
//	NUM-PROB-060  ExponentialQuantile                 exponential quantile −ln(1 − u)/λ
//	NUM-PROB-061  NormalQuantile                      normal quantile by Wichura's PPND16
//	NUM-PROB-062  LognormalQuantile                   lognormal quantile exp(μ + σΦ⁻¹(u))
//	NUM-PROB-063  ChiSquareCDF             chi-square P(X ≤ x), regularized incomplete gamma P(ν/2, x/2)
//	NUM-PROB-064  ChiSquareSurvival        chi-square P(X > x), the p-value, accurate in the upper tail
//	NUM-PROB-065  ChiSquareQuantile        chi-square quantile (critical value), safeguarded Newton
//	NUM-PROB-066  LikelihoodRatioPositive  LR+ = sensitivity/(1 − specificity)
//	NUM-PROB-067  LikelihoodRatioNegative  LR− = (1 − sensitivity)/specificity
//	NUM-PROB-068  PostTestProbability      probability after a result with likelihood ratio LR (odds × LR)
//	NUM-PROB-069  GammaPDF                 gamma density, shape and scale, log space
//	NUM-PROB-070  GammaCDF                 gamma P(X ≤ x), regularized incomplete gamma P(k, x/θ)
//	NUM-PROB-071  GammaQuantile            gamma quantile, bisection on the CDF or survival
//	NUM-PROB-072  GammaMeanVariance        gamma mean kθ and variance kθ²
//	NUM-PROB-073  BetaPDF                  beta density on [0, 1], log space
//	NUM-PROB-074  BetaCDF                  beta P(X ≤ x), regularized incomplete beta I_x(α, β)
//	NUM-PROB-075  BetaQuantile             beta quantile, bisection on [0, 1]
//	NUM-PROB-076  BetaMeanVariance         beta mean and variance
//	NUM-PROB-077  StudentTPDF              Student t density
//	NUM-PROB-078  StudentTCDF              Student t P(T ≤ t), incomplete beta, tail-accurate
//	NUM-PROB-079  StudentTQuantile         Student t quantile (critical value), symmetric bisection
//	NUM-PROB-080  StudentTMeanVariance     Student t mean 0 and variance ν/(ν − 2)
//	NUM-PROB-081  FDistributionPDF         Snedecor F density
//	NUM-PROB-082  FDistributionCDF         Snedecor F P(X ≤ x), incomplete beta
//	NUM-PROB-083  FDistributionQuantile    Snedecor F quantile (critical value), bisection
//	NUM-PROB-084  FDistributionMeanVariance  Snedecor F mean and variance
//	NUM-PROB-085  BinomialQuantile         least k with binomial P(X ≤ k) ≥ u
//	NUM-PROB-086  PoissonQuantile          least k with Poisson P(X ≤ k) ≥ u
//	NUM-PROB-087  WeibullPDF               Weibull density, shape and scale
//	NUM-PROB-088  WeibullCDF               Weibull P(X ≤ x) = −expm1(−(x/λ)^k)
//	NUM-PROB-089  WeibullQuantile          Weibull quantile λ(−ln(1 − u))^(1/k)
//	NUM-PROB-090  WeibullMeanVariance      Weibull mean and variance via Γ
//	NUM-PROB-091  CauchyPDF                Cauchy density, location and scale
//	NUM-PROB-092  CauchyCDF                Cauchy P(X ≤ x) = atan2(1, −z)/π
//	NUM-PROB-093  CauchyQuantile           Cauchy quantile x₀ + γ tan(π(u − ½)), tail-accurate
//	NUM-PROB-094  MultivariateNormalPDF    n-variate normal density from the Cholesky factor of Σ
//	NUM-PROB-095  MultivariateNormalLogPDF  n-variate normal log density from the Cholesky factor of Σ
//	NUM-PROB-096  ChiSquarePDF             chi-square density (gamma with shape k/2, scale 2)
//	NUM-PROB-097  ChiSquareMeanVariance    chi-square mean k and variance 2k
//	NUM-PROB-098  HypergeometricQuantile   least k with hypergeometric P(X ≤ k) ≥ u
//	NUM-PROB-099  NegativeBinomialFailuresPMF  negative binomial, FAILURES before the r-th success
//	NUM-PROB-100  NegativeBinomialFailuresCDF  negative binomial failures P(Y ≤ k)
//	NUM-PROB-101  NegativeBinomialFailuresQuantile  least k with failures P(Y ≤ k) ≥ u
//	NUM-PROB-102  NegativeBinomialFailuresMeanVariance  failures mean r(1 − p)/p and variance r(1 − p)/p²
//	NUM-PROB-103  GumbelPDF                Gumbel (type I extreme value, maximum) density
//	NUM-PROB-104  GumbelCDF                Gumbel P(X ≤ x) = exp(−e^(−z))
//	NUM-PROB-105  GumbelQuantile           Gumbel quantile μ − β ln(−ln u)
//	NUM-PROB-106  GumbelMeanVariance       Gumbel mean μ + γβ and variance π²β²/6
//	NUM-PROB-107  LaplacePDF               Laplace (double exponential) density
//	NUM-PROB-108  LaplaceCDF               Laplace P(X ≤ x)
//	NUM-PROB-109  LaplaceQuantile          Laplace quantile
//	NUM-PROB-110  LaplaceMeanVariance      Laplace mean μ and variance 2b²
//	NUM-PROB-111  LogisticPDF              logistic density
//	NUM-PROB-112  LogisticCDF              logistic P(X ≤ x), the sigmoid, tail-accurate
//	NUM-PROB-113  LogisticQuantile         logistic quantile, the logit
//	NUM-PROB-114  LogisticMeanVariance     logistic mean μ and variance π²s²/3
//	NUM-PROB-115  ParetoPDF                Pareto (first kind) density
//	NUM-PROB-116  ParetoCDF                Pareto P(X ≤ x) = 1 − (x_m/x)^α
//	NUM-PROB-117  ParetoQuantile           Pareto quantile x_m(1 − u)^(−1/α)
//	NUM-PROB-118  ParetoMeanVariance       Pareto mean and variance (α > 2)
//	NUM-PROB-119  RayleighPDF              Rayleigh density
//	NUM-PROB-120  RayleighCDF              Rayleigh P(X ≤ x) = 1 − e^(−x²/2σ²)
//	NUM-PROB-121  RayleighQuantile         Rayleigh quantile σ√(−2 ln(1 − u))
//	NUM-PROB-122  RayleighMeanVariance     Rayleigh mean σ√(π/2) and variance (4 − π)σ²/2
//	NUM-PROB-123  TriangularPDF            triangular density on [a, b], mode c
//	NUM-PROB-124  TriangularCDF            triangular P(X ≤ x)
//	NUM-PROB-125  TriangularQuantile       triangular quantile
//	NUM-PROB-126  TriangularMeanVariance   triangular mean (a + b + c)/3 and variance
//	NUM-PROB-127  PowerFunctionPDF         power-function density a·x^(a−1) on [0, 1] (beta(a, 1))
//	NUM-PROB-128  PowerFunctionCDF         power-function P(X ≤ x) = x^a
//	NUM-PROB-129  PowerFunctionQuantile    power-function quantile u^(1/a)
//	NUM-PROB-130  PowerFunctionMeanVariance  power-function mean and variance
//	NUM-PROB-131  MultinomialPMF           multinomial P(X = k), log space
//	NUM-PROB-132  DirichletPDF             Dirichlet density on the simplex, log space
//	NUM-PROB-133  MultivariateHypergeometricPMF  multivariate hypergeometric P(X = k), log space
//	NUM-PROB-134  LogSeriesPMF             log-series P(X = k) = −p^k/(k ln(1 − p))
//	NUM-PROB-135  LogSeriesCDF             log-series P(X ≤ k)
//	NUM-PROB-136  LogSeriesQuantile        log-series quantile by sequential search
//	NUM-PROB-137  LogSeriesMeanVariance    log-series mean and variance
//	NUM-PROB-138  WaldPDF                  inverse Gaussian (Wald) density
//	NUM-PROB-139  WaldCDF                  inverse Gaussian P(X ≤ x) via Φ, overflow-free
//	NUM-PROB-140  WaldQuantile             inverse Gaussian quantile by bisection
//	NUM-PROB-141  WaldMeanVariance         inverse Gaussian mean μ and variance μ³/λ
//	NUM-PROB-142  VonMisesPDF              von Mises (circular normal) density
//	NUM-PROB-143  VonMisesCDF              von Mises P(X ≤ x), Hill's series
//	NUM-PROB-144  ZipfPMF                  Zipf (zeta) P(X = k) = k^(−a)/ζ(a)
//	NUM-PROB-145  ZipfCDF                  Zipf P(X ≤ k)
//	NUM-PROB-146  ZipfMeanVariance         Zipf mean and variance (a > 3)
//	NUM-PROB-147  ZipfQuantile             Zipf quantile, search on the CDF
//	NUM-PROB-148  VonMisesQuantile         von Mises quantile by bisection on the CDF
//	NUM-PROB-149  MultinomialMeanVariance  multinomial mean vector and covariance matrix
//	NUM-PROB-150  DirichletMeanVariance    Dirichlet mean vector and covariance matrix
//	NUM-PROB-151  NoncentralChiSquarePDF   noncentral chi-square density, Poisson mixture from its peak
//	NUM-PROB-152  NoncentralChiSquareCDF   noncentral chi-square P(X ≤ x), Ding AS 275
//	NUM-PROB-153  NoncentralChiSquareMeanVariance  noncentral chi-square mean df + λ, variance 2(df + 2λ)
//	NUM-PROB-154  NoncentralFCDF           noncentral F P(X ≤ x), Lenth AS 226
//	NUM-PROB-155  NoncentralFMeanVariance  noncentral F mean and variance (df2 > 4)
package prob

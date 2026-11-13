# S9: Probability (`num/prob`, SHORT `PROB`)

Status: SHIPPED 2026-11-12
Work folder: `doc/math_build_2026-11-12/` (hand-off: `notes_prob.md`)
Decisions: D-013 (trials convention for the geometric, in the names; plain functions; odds here; Cauchy only if
sourced; log binomial from num/combinatorics where the lgamma route was planned), D-015 (an unconfirmed citation
holds back only its own function) · Assumptions: proposed in `notes_prob.md` (A-PROB-1…4).
Catalog: [`doc/catalog/prob.md`](../catalog/prob.md) Step 3 holds the full entry per function (formula, inputs,
errors, worked example); this spec records only what the build fixes or changes.

## Functions

| ID | Name | Signature | Units |
|---|---|---|---|
| NUM-PROB-001 | `ProbabilityFromCounts` | `func ProbabilityFromCounts(favorable, total int) (float64, error)` | favorable, total: counts of outcomes; result: a probability in [0, 1] |
| NUM-PROB-002 | `ComplementProbability` | `func ComplementProbability(p float64) (float64, error)` | p and result: probabilities in [0, 1] |
| NUM-PROB-003 | `IndependentJointProbability` | `func IndependentJointProbability(pA, pB float64) (float64, error)` | pA, pB and result: probabilities in [0, 1] |
| NUM-PROB-004 | `JointProbability` | `func JointProbability(pAGivenB, pB float64) (float64, error)` | pAGivenB, pB and result: probabilities in [0, 1] |
| NUM-PROB-005 | `UnionProbability` | `func UnionProbability(pA, pB, pAAndB float64) (float64, error)` | pA, pB, pAAndB and result: probabilities in [0, 1] |
| NUM-PROB-006 | `ConditionalProbability` | `func ConditionalProbability(pAAndB, pB float64) (float64, error)` | pAAndB, pB and result: probabilities in [0, 1] |
| NUM-PROB-007 | `TotalProbability` | `func TotalProbability(pAGivenB, pB []float64) (float64, error)` | every input and the result: probabilities in [0, 1] |
| NUM-PROB-008 | `BayesTheorem` | `func BayesTheorem(pBGivenA, pA, pB float64) (float64, error)` | pBGivenA, pA, pB and result: probabilities in [0, 1] |
| NUM-PROB-009 | `BayesPosterior` | `func BayesPosterior(dst, likelihoods, priors []float64) ([]float64, error)` | every input and output: probabilities in [0, 1] |
| NUM-PROB-010 | `PositivePredictiveValue` | `func PositivePredictiveValue(sensitivity, specificity, prevalence float64) (float64, error)` | sensitivity, specificity, prevalence and result: probabilities in [0, 1] |
| NUM-PROB-011 | `NegativePredictiveValue` | `func NegativePredictiveValue(sensitivity, specificity, prevalence float64) (float64, error)` | sensitivity, specificity, prevalence and result: probabilities in [0, 1] |
| NUM-PROB-012 | `OddsFromProbability` | `func OddsFromProbability(p float64) (float64, error)` | p: a probability in [0, 1); result: odds in favour, a ratio ≥ 0 |
| NUM-PROB-013 | `ProbabilityFromOdds` | `func ProbabilityFromOdds(odds float64) (float64, error)` | odds: a ratio ≥ 0; result: a probability in [0, 1] |
| NUM-PROB-014 | `DecimalOddsFromProbability` | `func DecimalOddsFromProbability(p float64) (float64, error)` | p: a probability in (0, 1]; result: decimal odds ≥ 1 (units returned per unit staked) |
| NUM-PROB-015 | `ImpliedProbability` | `func ImpliedProbability(decimalOdds float64) (float64, error)` | decimalOdds: units returned per unit staked, ≥ 1; result: a probability in (0, 1] |
| NUM-PROB-016 | `Overround` | `func Overround(decimalOdds []float64) (float64, error)` | decimalOdds: units returned per unit staked, each ≥ 1; result: a dimensionless fraction |
| NUM-PROB-017 | `NormalizedImpliedProbabilities` | `func NormalizedImpliedProbabilities(dst, decimalOdds []float64) ([]float64, error)` | decimalOdds: units returned per unit staked, each ≥ 1; result: probabilities in (0, 1] |
| NUM-PROB-018 | `ExpectedValue` | `func ExpectedValue(outcomes, probabilities []float64) (float64, error)` | outcomes and result: the caller's unit; probabilities: in [0, 1] |
| NUM-PROB-019 | `DiscreteVariance` | `func DiscreteVariance(outcomes, probabilities []float64) (float64, error)` | outcomes: the caller's unit; probabilities: in [0, 1]; result: that unit squared |
| NUM-PROB-020 | `ChebyshevBound` | `func ChebyshevBound(k float64) (float64, error)` | k: a number of standard deviations, dimensionless; result: a probability bound in [0, 1] |
| NUM-PROB-021 | `ChebyshevBoundEpsilon` | `func ChebyshevBoundEpsilon(variance, epsilon float64) (float64, error)` | variance: the caller's unit squared; epsilon: that unit; result: a probability bound in [0, 1] |
| NUM-PROB-022 | `MarkovBound` | `func MarkovBound(mean, threshold float64) (float64, error)` | mean and threshold: the same unit; result: a probability bound in [0, 1] |
| NUM-PROB-023 | `BernoulliPMF` | `func BernoulliPMF(k int, p float64) (float64, error)` | k: an outcome, 0 or 1; p and result: probabilities in [0, 1] |
| NUM-PROB-024 | `BernoulliCDF` | `func BernoulliCDF(k int, p float64) (float64, error)` | k: an outcome; p and result: probabilities in [0, 1] |
| NUM-PROB-025 | `BernoulliMeanVariance` | `func BernoulliMeanVariance(p float64) (mean, variance float64, err error)` | p: a probability; mean: a probability; variance: dimensionless |
| NUM-PROB-026 | `BinomialPMF` | `func BinomialPMF(k, n int, p float64) (float64, error)` | k, n: counts of successes and trials; p and result: probabilities in [0, 1] |
| NUM-PROB-027 | `BinomialCDF` | `func BinomialCDF(k, n int, p float64) (float64, error)` | k, n: counts of successes and trials; p and result: probabilities in [0, 1] |
| NUM-PROB-028 | `BinomialMeanVariance` | `func BinomialMeanVariance(n int, p float64) (mean, variance float64, err error)` | n: a count of trials; p: a probability; mean: a count of successes; variance: that count squared |
| NUM-PROB-029 | `GeometricTrialsPMF` | `func GeometricTrialsPMF(k int, p float64) (float64, error)` | k: a count of trials; p and result: probabilities in [0, 1] |
| NUM-PROB-030 | `GeometricTrialsCDF` | `func GeometricTrialsCDF(k int, p float64) (float64, error)` | k: a count of trials; p and result: probabilities in [0, 1] |
| NUM-PROB-031 | `GeometricTrialsMeanVariance` | `func GeometricTrialsMeanVariance(p float64) (mean, variance float64, err error)` | p: a probability; mean: a count of trials; variance: that count squared |
| NUM-PROB-032 | `NegativeBinomialTrialsPMF` | `func NegativeBinomialTrialsPMF(x, r int, p float64) (float64, error)` | x: a count of trials; r: a count of successes; p and result: probabilities in [0, 1] |
| NUM-PROB-033 | `NegativeBinomialTrialsCDF` | `func NegativeBinomialTrialsCDF(x, r int, p float64) (float64, error)` | x: a count of trials; r: a count of successes; p and result: probabilities in [0, 1] |
| NUM-PROB-034 | `NegativeBinomialTrialsMeanVariance` | `func NegativeBinomialTrialsMeanVariance(r int, p float64) (mean, variance float64, err error)` | r: a count of successes; p: a probability; mean: a count of trials; variance: that count squared |
| NUM-PROB-035 | `HypergeometricPMF` | `func HypergeometricPMF(k, population, successes, draws int) (float64, error)` | k, population, successes, draws: counts of items; result: a probability in [0, 1] |
| NUM-PROB-036 | `HypergeometricCDF` | `func HypergeometricCDF(k, population, successes, draws int) (float64, error)` | k, population, successes, draws: counts of items; result: a probability in [0, 1] |
| NUM-PROB-037 | `HypergeometricMeanVariance` | `func HypergeometricMeanVariance(population, successes, draws int) (mean, variance float64, err error)` | population, successes, draws: counts of items; mean: a count of successes; variance: that count squared |
| NUM-PROB-038 | `DiscreteUniformPMF` | `func DiscreteUniformPMF(k, a, b int) (float64, error)` | k, a, b: integer outcomes; result: a probability in [0, 1] |
| NUM-PROB-039 | `DiscreteUniformCDF` | `func DiscreteUniformCDF(k, a, b int) (float64, error)` | k, a, b: integer outcomes; result: a probability in [0, 1] |
| NUM-PROB-040 | `DiscreteUniformMeanVariance` | `func DiscreteUniformMeanVariance(a, b int) (mean, variance float64, err error)` | a, b: integer outcomes; mean: that unit; variance: that unit squared |
| NUM-PROB-041 | `PoissonPMF` | `func PoissonPMF(k int, lambda float64) (float64, error)` | k: a count of events; lambda: the mean count per interval; result: a probability in [0, 1] |
| NUM-PROB-042 | `PoissonCDF` | `func PoissonCDF(k int, lambda float64) (float64, error)` | k: a count of events; lambda: the mean count per interval; result: a probability in [0, 1] |
| NUM-PROB-043 | `PoissonMeanVariance` | `func PoissonMeanVariance(lambda float64) (mean, variance float64, err error)` | lambda: the mean count per interval; mean: a count; variance: that count squared |
| NUM-PROB-044 | `UniformPDF` | `func UniformPDF(x, a, b float64) (float64, error)` | x, a, b: one unit; result: probability per that unit |
| NUM-PROB-045 | `UniformCDF` | `func UniformCDF(x, a, b float64) (float64, error)` | x, a, b: one unit; result: a probability in [0, 1] |
| NUM-PROB-046 | `UniformMeanVariance` | `func UniformMeanVariance(a, b float64) (mean, variance float64, err error)` | a, b: one unit; mean: that unit; variance: that unit squared |
| NUM-PROB-047 | `NormalPDF` | `func NormalPDF(x, mean, sd float64) (float64, error)` | x, mean, sd: one unit; result: probability per that unit |
| NUM-PROB-048 | `NormalCDF` | `func NormalCDF(x, mean, sd float64) (float64, error)` | x, mean, sd: one unit; result: a probability in [0, 1] |
| NUM-PROB-049 | `NormalSurvival` | `func NormalSurvival(x, mean, sd float64) (float64, error)` | x, mean, sd: one unit; result: a probability in [0, 1] |
| NUM-PROB-050 | `NormalIntervalProbability` | `func NormalIntervalProbability(lo, hi, mean, sd float64) (float64, error)` | lo, hi, mean, sd: one unit; result: a probability in [0, 1] |
| NUM-PROB-051 | `LognormalPDF` | `func LognormalPDF(x, mu, sigma float64) (float64, error)` | x: the caller's unit, > 0; mu: natural log of that unit; sigma: dimensionless; result: probability per unit of x |
| NUM-PROB-052 | `LognormalCDF` | `func LognormalCDF(x, mu, sigma float64) (float64, error)` | x: the caller's unit; mu: natural log of that unit; sigma: dimensionless; result: a probability in [0, 1] |
| NUM-PROB-053 | `LognormalMeanVariance` | `func LognormalMeanVariance(mu, sigma float64) (mean, variance float64, err error)` | mu: natural log of the caller's unit; sigma: dimensionless; mean: that unit; variance: that unit squared |
| NUM-PROB-054 | `ExponentialPDF` | `func ExponentialPDF(x, rate float64) (float64, error)` | x: a time or length unit; rate: per that unit; result: probability per that unit |
| NUM-PROB-055 | `ExponentialCDF` | `func ExponentialCDF(x, rate float64) (float64, error)` | x: a time or length unit; rate: per that unit; result: a probability in [0, 1] |
| NUM-PROB-056 | `ExponentialMeanVariance` | `func ExponentialMeanVariance(rate float64) (mean, variance float64, err error)` | rate: per a time or length unit; mean: that unit; variance: that unit squared |
| NUM-PROB-057 | `GeometricTrialsQuantile` | `func GeometricTrialsQuantile(u, p float64) (int, error)` | u, p: probabilities; result: a count of trials |
| NUM-PROB-058 | `DiscreteUniformQuantile` | `func DiscreteUniformQuantile(u float64, a, b int) (int, error)` | u: a probability; a, b and result: integer outcomes |
| NUM-PROB-059 | `UniformQuantile` | `func UniformQuantile(u, a, b float64) (float64, error)` | u: a probability in [0, 1]; a, b and result: one unit |
| NUM-PROB-060 | `ExponentialQuantile` | `func ExponentialQuantile(u, rate float64) (float64, error)` | u: a probability in [0, 1); rate: per a time or length unit; result: that unit |
| NUM-PROB-061 | `NormalQuantile` | `func NormalQuantile(u, mean, sd float64) (float64, error)` | u: a probability in (0, 1); mean, sd and result: one unit |
| NUM-PROB-062 | `LognormalQuantile` | `func LognormalQuantile(u, mu, sigma float64) (float64, error)` | u: a probability in (0, 1); mu: natural log of the caller's unit; sigma: dimensionless; result: that unit |

Sentinels: `ErrNotFinite`, `ErrProbability`, `ErrOutOfRange`, `ErrInconsistent`, `ErrNotDistribution`,
`ErrOverflow`, each wrapped with the parameter name. No function modifies its input; every scalar function
allocates nothing on success; 009 and 017 append to `dst`.

## Package-wide choices (changes from the catalog)

- **Names (D-013)**: 029/030/031/057 are `GeometricTrials*`. 032/033/034 are `NegativeBinomialTrials*`: the
  Pascal form counts trials too, and the common library default (R, SciPy) counts failures, so the name says
  which (assumption A-PROB-1). Hypergeometric parameters are named `population`, `successes`, `draws`
  (N, K, n in the formula) so that `N` and `n` do not differ only by case.
- **Distribution tolerance (catalog Q5)**: a probability slice must sum to 1 within n·2⁻⁵² (n·u, u = 2⁻⁵³,
  doubled for inputs that were themselves rounded once). Assumption A-PROB-2.
- **Log binomial (D-013)**: the Loader saddle-point route (δ = stirlerr, D₀ = bd0) is kept for 026, 032, 035,
  041. It is more accurate than exp(LogBinomial + k·ln p + (n−k)·ln q): that sum's absolute error is about
  2⁻⁵³·n·|ln q| in the exponent (a relative error of 1e-13 in the PMF at n = 1000), while Loader's is about
  2⁻⁵³·D₀, which is near 0 at the mode. `num/combinatorics` is therefore not imported. **Lost**: the
  LogBinomial route (simpler, less accurate in the body of the distribution).
- **δ table**: δ(m) for m = 1…15 was computed for this build from its definition with 60-digit decimal
  arithmetic (Python `decimal`) and agrees with Loader's Appendix A table to float64 precision. Above 15, the
  Stirling series (DLMF 5.11.1 coefficients 1/12, 1/360, 1/1260, 1/1680, 1/1188) cut at Loader's thresholds.
- **Catalog sign**: the catalog writes the saddle-point exponent as −δ(n) + δ(k) + …; the correct form (Loader
  eq. 5 and Appendix A) is +δ(n) − δ(k) − δ(n−k) − D₀(k, np) − D₀(n−k, nq). Built correctly.
- **Tail sums (027, 033, 036, 042)**: sum the side of k away from the mean from one saddle-point term outward
  with the ratio recurrence, compensated, and stop when the rest (bounded by t·r/(1−r) for a log-concave PMF)
  is below 2⁻⁵³ of the sum. Capped at 2²⁸ terms → ErrOutOfRange (σ above roughly 10⁷), so every call returns
  in bounded time (assumption A-PROB-3). **Lost**: the regularized incomplete beta/gamma (O(1)), which has no
  home yet (catalog Q6); the naive sum from i = 0 (underflows, O(n)).
- **Normal PDF tail**: for |z| ≥ 5, z is split on a 2⁻¹⁶ grid so that z² is exact (full relative accuracy to
  the underflow point). **Lost**: exp(−z²/2) directly (relative error about z²·2⁻⁵³).
- **Normal quantile (061)**: Wichura PPND16, coefficients typed from the published AS 241 listing
  (StatLib apstat/241, two fetches agreed digit for digit). **Lost**: Acklam's or Moro's approximations (1e-9,
  3e-9 accuracy), and Newton refinement on erfc (an iteration for no gain at 1e-16).
- **Quantiles 057, 058**: the closed-form estimate is corrected by single steps until F(k−1) < u ≤ F(k) holds
  for this package's own CDF; GeometricTrialsCDF(1, p) returns p exactly so that the u ≤ p shortcut agrees.
- **Odds (012, 013)**: Ross §3.3's definition of odds could not be confirmed online; the source is Bland,
  Altman, "The odds ratio", BMJ 2000 (read on PMC).
- **Decimal odds (014–017), fix 2026-11-12**: Štrumbelj 2014 (IJF) could not be read (ScienceDirect 403, no
  preprint found), so it is no longer a Source. The formulas are carried by three peer-reviewed papers read
  first-hand (D-015 second-source rule): Wheatcroft, J. Sports Analytics 2021 (preprint arXiv:2001.09097v1 §2.1:
  rᵢ = 1/Oᵢ and eq. (1) π = Σ 1/Oᵢ − 1); Egidi, Pauli, Torelli, Statistical Modelling 2018 (preprint
  arXiv:1802.08848 §2, eq. (2.1): booksum β = Σ inverse odds, πᵢ = oᵢ/β); Wunderlich, Memmert, PLOS ONE 2018
  (Method, "Transferring betting odds to probabilities": no margin → inverse odds are probabilities; basic
  normalisation). None held.

## Per function

Every entry is as cataloged (formula, inputs, errors, worked example, best approach) except as follows.

### Rules 001–022
- **Sources**: Ross citations are at section level (the proposition, equation and example numbers could not be
  confirmed online; the Pearson 10th-edition table of contents confirms the sections). Casella–Berger is cited
  at section level too (fix 2026-11-12): §1.2 "Basics of Probability Theory" (002, 005, 007) and §1.3
  "Conditional Probability and Independence" (003, 006, 008, 009), titles confirmed from the TOC; the result
  numbers 1.2.8, 1.2.9, 1.2.11, 1.3.2, 1.3.5, 1.3.7 came from a search summary and are no longer cited. CB
  Theorem 3.6.1 was not confirmed either, so 020–022 cite CB §3.6 "Inequalities and Identities".
- **010, 011 (fix 2026-11-12)**: Altman & Bland 1994 (BMJ 309:102) is a page scan that could not be read, so it
  is no longer a Source; Ross §3.3 carries the formula. The 2×2 example (TP 231, FP 32, FN 27, TN 54) is kept as
  a hand-checked table, not attributed.
- **Worked examples**: Ross's Example 3a, 5a and 2a numbers are not cited; the values (die 7/2, 35/12; the
  mean-50 production bounds 2/3 and 3/4) are hand-checked.
- **009, 017**: `dst` first, append semantics, as algo/window. **016**: the −1 is added inside the compensated
  sum. **018/019**: products added with their FMA-recovered rounding error; 018 retries with outcomes scaled by
  2⁻² when a partial sum overflows. **019 (fix 2026-11-12)**: the corrected two-pass algorithm, Σ pᵢdᵢ² −
  (Σ pᵢdᵢ)², Chan, Golub, LeVeque, The American Statistician 1983, 37(3):242–247, eq. (1.7) (read first-hand),
  with the weights 1/N replaced by pᵢ.
- **022 MarkovBound**: parameter `threshold` (the catalog's `a`).

### Discrete 023–043
- **026**: tails relative error about 2⁻⁵³·D₀ (documented; the catalog's "near machine precision for every n"
  holds near the mode).
- **027, 033, 036, 042**: as above; ErrOutOfRange past 2²⁸ terms.
- **035**: three saddle-point terms combined in their exponents (one exp), q = n/N.
- **038–040, 058 (fix 2026-11-12)**: every int range is accepted; the span b − a is exact in uint64 and
  m = b − a + 1 (up to 2⁶⁴) is held in float64, so ErrOverflow is no longer returned (the catalog's
  "b − a + 1 > MaxInt is ErrOverflow" reported representable answers as errors). The mean is a + ⌊(b − a)/2⌋
  in integers plus ½ when b − a is odd (a/2 + b/2 in float64 gave 0, not −½, on MinInt…MaxInt). Above 2⁵³
  points the quantile resolves to the float64 CDF's own steps.
- **041 worked example**: the catalog's λ = k = 1000 value 0.012614611348719664 is off by 1.5e-13 relative; the
  correct value, 0.0126146113487215 (50-digit decimal), is used.

### Continuous 044–056, quantiles 057–062
- **044–046, 059**: NIST §1.3.6.6.2 gives the general PDF, mean and SD; the CDF and PPF only for the standard
  form, so Ross §5.3 is the second source.
- **044, 045, 059 (fix 2026-11-12)**: when b − a overflows float64 (a < 0 < b, |a| or |b| ≥ 2¹⁰²³) they work
  on the halved endpoints, which are exact there: PDF ½/(b/2 − a/2) (≥ 2⁻¹⁰²⁵, subnormal, never 0), CDF
  (x/2 − a/2)/(b/2 − a/2), quantile 2·(a/2 + u·(b/2 − a/2)). ErrOverflow remains only for a density above the
  float64 range (b − a subnormal). **046** keeps ErrOverflow when the variance (b − a)²/12 overflows (from
  b − a ≈ 1.2e154), which covers every overflowing width; the mean alone would be representable, but the pair
  is not returned partially. **Lost**: ErrOverflow whenever b − a overflows (rejected representable results).
- **048–050**: DLMF §7.1 gives Φ(z) = ½erfc(−z/√2) and Q(z) = ½erfc(z/√2) (notation section, no equation number),
  §7.2 eqs. 7.2.1–7.2.2 define erf and erfc.
- **051**: x·σ below the normal range (or beyond it) switches to a log-domain formula.
- **053**: when expm1(σ²) or exp(2μ + σ²) over- or underflows alone, the variance comes from its logarithm.

## Out of scope

Cauchy (NIST §1.3.6.6.3 exists with PDF, CDF and PPF; left LATER per D-013), Student t, χ², F (inference
slice), binomial and Poisson quantiles, Shin's method, the held odds formats (decimal/fractional/American: no
primary source found), random sampling (a sampler would take a `rand.Source`), log-PMFs.

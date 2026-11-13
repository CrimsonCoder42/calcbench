# Catalog: probability (`num/prob`, SHORT `PROB`)

Catalog pass, 2026-11-12. Docs only: nothing is built until Devin approves this file.
Input: the vCalc "Probability" list, NAMES ONLY (no vCalc page was fetched; nothing else was taken from it).
Work folder: [../prob_2026-11-12/](../prob_2026-11-12/README.md). Shared-file rows: `PENDING_SHARED.md` there.
Read, not edited: [stats.md](stats.md) (NUM-STATS-001 ArithmeticMean, 003 CoefficientOfVariation (RSD),
028 ZScore, 029 RawScoreFromZScore, 020 SummaryStatistics and 002 RelativeFrequency are stats' and are not
repeated here).

Citation status, used in every entry: **CHECKED** = the locator was confirmed on the web on 2026-11-12.
**CONFIRM IN PRINT** = the book and chapter are right, but the theorem, equation or example number has to be
checked against the printed copy before the build cites it. The two textbooks are cited as:
- **[CB]** Casella, Berger, *Statistical Inference*, 2nd ed., Duxbury, 2002.
- **[Ross]** Ross, *A First Course in Probability*, 10th ed., Pearson, 2019.
- **[NIST]** NIST/SEMATECH e-Handbook of Statistical Methods, §1.3.6.6 "Gallery of Distributions",
  https://www.itl.nist.gov/div898/handbook/eda/section3/eda366.htm (read 2026-11-12). Sections CHECKED:
  1.3.6.6.7 Exponential, 1.3.6.6.9 Lognormal, 1.3.6.6.18 Binomial, 1.3.6.6.19 Poisson, 1.3.6.7.1 normal
  table. Inferred from the page order, not opened: 1.3.6.6.1 Normal, 1.3.6.6.2 Uniform, 1.3.6.6.3 Cauchy. The
  gallery has no geometric, negative binomial or hypergeometric page.
- **[Loader]** Loader, "Fast and Accurate Computation of Binomial Probabilities", technical report, 2000
  (revised 2002), https://r-project.org/doc/reports/CLoader-dbinom-2002.pdf [CHECKED that it exists; read its
  sections before the build].

## Step 1: clean

The input's compressed lines are expanded so that every row the vCalc list holds has its own line: "(x2)"
gives two rows, the poker line gives 20 (10 hands × fraction and percentage), and "Roll a N-sided" gives one
row per N. That makes **139 rows**.

| # | Name as given | Decision | Home | Reason |
|---|---|---|---|---|
| 1 | Collections: D&D 1.0 - Dice Rolls | MOVE | `life/games` | a D&D collection; its dice maths is the DiceRoll row (#20) |
| 2 | Collections: LM 33.2 Calculating randomness | MOVE | `sci/physics` | an LM textbook collection, routed as the request says |
| 3 | Collections: LM 33.3 Probability distributions | MOVE | `sci/physics` | an LM textbook collection, routed as the request says |
| 4 | Collections: LM 33.5 Applications of calculus | MOVE | `sci/physics` | an LM textbook collection, routed as the request says |
| 5 | Collections: LM 33.6 Summary | MOVE | `sci/physics` | an LM textbook collection, routed as the request says |
| 6 | Collections: Poker Calculator (five card) | MOVE | `life/games` | poker |
| 7 | Arithmetic Mean | MOVE | `num/stats` | already NUM-STATS-001 in `doc/catalog/stats.md` |
| 8 | Bayes Theorem | KEEP | `num/prob` → NUM-PROB-008 | Bayes' theorem, scalar form |
| 9 | Bayes' Theorem for Disease Testing | KEEP | `num/prob` → NUM-PROB-010, 011 | PPV and NPV from sensitivity, specificity, prevalence |
| 10 | Bayes' Theorem for Event Probability | KEEP | `num/prob` → NUM-PROB-009 | posterior over a partition (total probability in the denominator) |
| 11 | Binomial Coefficient | MOVE | `num/combinatorics` | combinatorics |
| 12 | Cauchy-Lorentz Distribution | LATER | `num/prob`, later | no consumer yet and no mean or variance, so it does not fit the set (Q4) |
| 13 | Chebyshev's Inequality(1) | KEEP | `num/prob` → NUM-PROB-020 | the k-standard-deviations form |
| 14 | Chebyshev's Inequality(1) | MERGE | → #13 | duplicate |
| 15 | Chebyshev's Inequality(2) | KEEP | `num/prob` → NUM-PROB-021 | the variance-and-ε form |
| 16 | Chebyshev's Inequality(2) | MERGE | → #15 | duplicate |
| 17 | Chi-Square (Stats) | MOVE | `num/stats` | the chi-square test; stats lists it LATER with the inference slice |
| 18 | Circular r-Permutation of n Elements | MOVE | `num/combinatorics` | combinatorics |
| 19 | Coin Flip | DROP | — | random, not a pure function |
| 20 | D20 | MOVE | `life/games` | the one DiceRoll row: the pure form is the distribution of a dice total; one die is NUM-PROB-038 |
| 21 | D4 | MERGE | → #20 | dice |
| 22 | DnD Dice Roller | MERGE | → #20 | dice |
| 23 | Displaced Gaussian | MERGE | `num/prob` → NUM-PROB-047 | normal PDF with a mean and a standard deviation |
| 24 | Drake Equation (original) | MOVE | `sci/astronomy` | astronomy |
| 25 | Expected Value(Bernoulli) | KEEP | `num/prob` → NUM-PROB-025 | Bernoulli mean |
| 26 | Expected Value(Binomial) | KEEP | `num/prob` → NUM-PROB-028 | binomial mean |
| 27 | Expected Value(Discrete Uniform) | KEEP | `num/prob` → NUM-PROB-040 | discrete uniform mean |
| 28 | Expected Value(Geometric) | KEEP | `num/prob` → NUM-PROB-031 | geometric mean (convention Q1) |
| 29 | Expected Value(Geometric) | MERGE | → #28 | duplicate |
| 30 | Expected Value(Hypergeometric) | KEEP | `num/prob` → NUM-PROB-037 | hypergeometric mean |
| 31 | Expected Value(Pascal Distribution) | KEEP | `num/prob` → NUM-PROB-034 | negative binomial (Pascal) mean |
| 32 | Expected Value(Pascal) | MERGE | → #31 | duplicate |
| 33 | Friedman Equation | ASK DEVIN | not planned | cosmology, not probability; the request marks it ask-first |
| 34 | Gamma Function | MOVE | `num/combinatorics` | special function; stats.md routes it there too |
| 35 | Gaussian | KEEP | `num/prob` → NUM-PROB-047 | normal PDF |
| 36 | Gaussian Distribution | MERGE | `num/prob` → NUM-PROB-047 | same PDF |
| 37 | Gaussian Distribution - Displaced | MERGE | `num/prob` → NUM-PROB-047 | same PDF with a mean and a standard deviation |
| 38 | Happiness Probability | ASK DEVIN | not planned | no recognised formula; the request marks it ask-first |
| 39 | Lognormal | KEEP | `num/prob` → NUM-PROB-051 (and 052, 053, 062) | lognormal |
| 40 | Lognormal Distribution | MERGE | → #39 | duplicate |
| 41 | Magic the Gathering - Opening Hand Mana Curve | MOVE | `life/games` | a card game; it calls the hypergeometric (NUM-PROB-035, 036) |
| 42 | Magic the Gathering - Opening Hand Mana Curve | MERGE | → #41 | duplicate |
| 43 | Multiple Values from Gaussian Distribution | DROP | — | random sampling, not a pure function |
| 44 | Multiple Values from Gaussian Distribution (as string) | DROP | — | random sampling, not a pure function |
| 45 | Normal Distribution (PDF) | MERGE | `num/prob` → NUM-PROB-047 | same PDF |
| 46 | P(event) = F / T | KEEP | `num/prob` → NUM-PROB-001 | classical probability from counts |
| 47 | P(n,S) - Binomial Probability | KEEP | `num/prob` → NUM-PROB-026 | binomial PMF |
| 48 | Percent Relative Standard Deviation | MOVE | `num/stats` | NUM-STATS-003 and its percent variant |
| 49 | Permutations | MOVE | `num/combinatorics` | combinatorics |
| 50 | PowerBall | MOVE | `life/games` | a lottery |
| 51 | Powerball Chances | MOVE | `life/games` | a lottery |
| 52 | Probability of Fatality | ASK DEVIN | not planned | no recognised formula; the request marks it ask-first |
| 53 | Probability of Flush losing (fraction) | MOVE | `life/games` | poker; one function per hand |
| 54 | Probability of Flush losing (percentage) | MERGE | → #53 | the same number ×100; one function per hand |
| 55 | Probability of Four of a Kind losing (fraction) | MOVE | `life/games` | poker; one function per hand |
| 56 | Probability of Four of a Kind losing (percentage) | MERGE | → #55 | the same number ×100; one function per hand |
| 57 | Probability of Full House losing (fraction) | MOVE | `life/games` | poker; one function per hand |
| 58 | Probability of Full House losing (percentage) | MERGE | → #57 | the same number ×100; one function per hand |
| 59 | Probability of High Card losing (fraction) | MOVE | `life/games` | poker; one function per hand |
| 60 | Probability of High Card losing (percentage) | MERGE | → #59 | the same number ×100; one function per hand |
| 61 | Probability of One Pair losing (fraction) | MOVE | `life/games` | poker; one function per hand |
| 62 | Probability of One Pair losing (percentage) | MERGE | → #61 | the same number ×100; one function per hand |
| 63 | Probability of Royal Flush losing (fraction) | MOVE | `life/games` | poker; one function per hand |
| 64 | Probability of Royal Flush losing (percentage) | MERGE | → #63 | the same number ×100; one function per hand |
| 65 | Probability of Straight Flush losing (fraction) | MOVE | `life/games` | poker; one function per hand |
| 66 | Probability of Straight Flush losing (percentage) | MERGE | → #65 | the same number ×100; one function per hand |
| 67 | Probability of Straight losing (fraction) | MOVE | `life/games` | poker; one function per hand |
| 68 | Probability of Straight losing (percentage) | MERGE | → #67 | the same number ×100; one function per hand |
| 69 | Probability of Three of a kind losing (fraction) | MOVE | `life/games` | poker; one function per hand |
| 70 | Probability of Three of a kind losing (percentage) | MERGE | → #69 | the same number ×100; one function per hand |
| 71 | Probability of Two Pair losing (fraction) | MOVE | `life/games` | poker; one function per hand |
| 72 | Probability of Two Pair losing (percentage) | MERGE | → #71 | the same number ×100; one function per hand |
| 73 | Probability of Picking NCAA Basketball Tournament | MOVE | `life/sports` | a bracket; sports |
| 74 | Random 0 OR 1 | DROP | — | random, not a pure function |
| 75 | Random Integer | DROP | — | random, not a pure function |
| 76 | Random Integer [0 , Upper] | DROP | — | random, not a pure function |
| 77 | Random Real Number in Range | DROP | — | random, not a pure function |
| 78 | Random Real [0.0, 1.0] | DROP | — | random, not a pure function |
| 79 | Random Real [0.0, Upper] | DROP | — | random, not a pure function |
| 80 | Relative Standard Deviation | MOVE | `num/stats` | NUM-STATS-003 |
| 81 | Roll a 2-sided Virtual Dice / Roll a d2 / Roll a D-2 | DROP | — | random, not a pure function |
| 82 | Roll a 3-sided Virtual Dice / Roll a d3 / Roll a D-3 | DROP | — | random, not a pure function |
| 83 | Roll a 4-sided Virtual Dice / Roll a d4 / Roll a D-4 | DROP | — | random, not a pure function |
| 84 | Roll a 5-sided Virtual Dice / Roll a d5 / Roll a D-5 | DROP | — | random, not a pure function |
| 85 | Roll a 6-sided Virtual Dice / Roll a d6 / Roll a D-6 | DROP | — | random, not a pure function |
| 86 | Roll a 7-sided Virtual Dice / Roll a d7 / Roll a D-7 | DROP | — | random, not a pure function |
| 87 | Roll a 8-sided Virtual Dice / Roll a d8 / Roll a D-8 | DROP | — | random, not a pure function |
| 88 | Roll a 9-sided Virtual Dice / Roll a d9 / Roll a D-9 | DROP | — | random, not a pure function |
| 89 | Roll a 10-sided Virtual Dice / Roll a d10 / Roll a D-10 | DROP | — | random, not a pure function |
| 90 | Roll a 11-sided Virtual Dice / Roll a d11 / Roll a D-11 | DROP | — | random, not a pure function |
| 91 | Roll a 12-sided Virtual Dice / Roll a d12 / Roll a D-12 | DROP | — | random, not a pure function |
| 92 | Roll a 20-sided Virtual Dice / Roll a d20 / Roll a D-20 | DROP | — | random, not a pure function |
| 93 | Roll a 50-sided Virtual Dice / Roll a d50 / Roll a D-50 | DROP | — | random, not a pure function |
| 94 | Roll a 100-sided Virtual Dice / Roll a d100 / Roll a D-100 | DROP | — | random, not a pure function |
| 95 | Roll Two Virtual Cube Dice | DROP | — | random, not a pure function |
| 96 | Roll Three Virtual Cube Dice | DROP | — | random, not a pure function |
| 97 | Roll Four Virtual Cube Dice | DROP | — | random, not a pure function |
| 98 | Roll Five Virtual Cube Dice | DROP | — | random, not a pure function |
| 99 | Seager Equation | MOVE | `sci/astronomy` | astronomy |
| 100 | Sod's Law | ASK DEVIN | not planned | not a formula with a source; the request marks it ask-first |
| 101 | Standard Normal Distribution | MERGE | `num/prob` → NUM-PROB-047 | the PDF at mean 0, sd 1 |
| 102 | Standard Normal Distribution (CDF) | KEEP | `num/prob` → NUM-PROB-048 | normal CDF |
| 103 | Stirling's Formula | MOVE | `num/combinatorics` | an approximation to n! |
| 104 | Variance(Bernoulli) | KEEP | `num/prob` → NUM-PROB-025 | bernoulli variance, returned with the mean |
| 105 | Variance(Bernoulli) | MERGE | → #104 | duplicate |
| 106 | Variance(Binomial) | KEEP | `num/prob` → NUM-PROB-028 | binomial variance, returned with the mean |
| 107 | Variance(Binomial) | MERGE | → #106 | duplicate |
| 108 | Variance(Geometric) | KEEP | `num/prob` → NUM-PROB-031 | geometric variance, returned with the mean |
| 109 | Variance(Geometric) | MERGE | → #108 | duplicate |
| 110 | Variance(Hypergeometric) | KEEP | `num/prob` → NUM-PROB-037 | hypergeometric variance, returned with the mean |
| 111 | Variance(Hypergeometric) | MERGE | → #110 | duplicate |
| 112 | Variance(Pascal) | KEEP | `num/prob` → NUM-PROB-034 | pascal variance, returned with the mean |
| 113 | Variance(Pascal) | MERGE | → #112 | duplicate |
| 114 | Variance(Uniform) | KEEP | `num/prob` → NUM-PROB-040 (continuous: 046) | uniform variance, returned with the mean |
| 115 | Variance(Uniform) | MERGE | → #114 | duplicate |
| 116 | Virtual Dice | DROP | — | random, not a pure function |
| 117 | Data items: five card hands | DROP | — | a table of poker results |
| 118 | Data items: Powerball Possibilities | DROP | — | a table of lottery results |
| 119 | Data items: Probability of Flush (fraction) | DROP | — | a table of poker results |
| 120 | Data items: Probability of Flush (percentage) | DROP | — | a table of poker results |
| 121 | Data items: Probability of Four of a Kind (fraction) | DROP | — | a table of poker results |
| 122 | Data items: Probability of Four of a Kind (percentage) | DROP | — | a table of poker results |
| 123 | Data items: Probability of Full House (fraction) | DROP | — | a table of poker results |
| 124 | Data items: Probability of Full House (percentage) | DROP | — | a table of poker results |
| 125 | Data items: Probability of High Card (fraction) | DROP | — | a table of poker results |
| 126 | Data items: Probability of High Card (percentage) | DROP | — | a table of poker results |
| 127 | Data items: Probability of One Pair (fraction) | DROP | — | a table of poker results |
| 128 | Data items: Probability of One Pair (percentage) | DROP | — | a table of poker results |
| 129 | Data items: Probability of Royal Flush (fraction) | DROP | — | a table of poker results |
| 130 | Data items: Probability of Royal Flush (percentage) | DROP | — | a table of poker results |
| 131 | Data items: Probability of Straight Flush (fraction) | DROP | — | a table of poker results |
| 132 | Data items: Probability of Straight Flush (percentage) | DROP | — | a table of poker results |
| 133 | Data items: Probability of Straight (fraction) | DROP | — | a table of poker results |
| 134 | Data items: Probability of Straight (percentage) | DROP | — | a table of poker results |
| 135 | Data items: Probability of Three of a kind (fraction) | DROP | — | a table of poker results |
| 136 | Data items: Probability of Three of a kind (percentage) | DROP | — | a table of poker results |
| 137 | Data items: Probability of Two Pair (fraction) | DROP | — | a table of poker results |
| 138 | Data items: Probability of Two Pair (percentage) | DROP | — | a table of poker results |
| 139 | Dataset: Unit Normal Distribution | DROP | — | a table; NUM-PROB-048 computes it |

Counts: KEEP 22 · MERGE 29 · MOVE 32 · DROP 51 · LATER 1 (Cauchy) · ASK DEVIN 4 (Friedman equation,
Happiness probability, Probability of fatality, Sod's law: "not planned", as the request says).

### Rows that stats.md routed to `num/prob` (covered here, so they are not built twice)

| stats.md row | Covered by |
|---|---|
| Conditional Probability of B given A (#32) | NUM-PROB-006 |
| Law of Total Probability (#72) | NUM-PROB-007 |
| Probability of A and B (#110) | NUM-PROB-003 |
| Probability of A or B (#111) | NUM-PROB-005 |
| Probability of Non-independent Events (#112) | NUM-PROB-004 |
| P(y) from z, P(y) Left and Right of z, P(y) right of z (#97–#99) | NUM-PROB-048, 049 |
| Probability between raw scores, between z scores (#108, #109) | NUM-PROB-050 |
| t Distribution Lookup, t-test critical value lookup (#157, #159) | LATER (Student t goes with the inference slice) |

## Step 2: the standard set

What a probability package must have, independently of the list. "Rows" are the Step 1 rows that each entry
covers. Several odds conversions and Shin's method are held back: see the end of the table.

| ID | Function | Rows |
|---|---|---|
| **Rules** | | |
| NUM-PROB-001 | ProbabilityFromCounts | #46 |
| NUM-PROB-002 | ComplementProbability | — |
| NUM-PROB-003 | IndependentJointProbability | stats #110 |
| NUM-PROB-004 | JointProbability | stats #112 |
| NUM-PROB-005 | UnionProbability | stats #111 |
| NUM-PROB-006 | ConditionalProbability | stats #32 |
| NUM-PROB-007 | TotalProbability | stats #72 |
| NUM-PROB-008 | BayesTheorem | #8 |
| NUM-PROB-009 | BayesPosterior | #10 |
| NUM-PROB-010 | PositivePredictiveValue | #9 |
| NUM-PROB-011 | NegativePredictiveValue | #9 |
| NUM-PROB-012 | OddsFromProbability | — |
| NUM-PROB-013 | ProbabilityFromOdds | — |
| NUM-PROB-014 | DecimalOddsFromProbability | — |
| NUM-PROB-015 | ImpliedProbability | — |
| NUM-PROB-016 | Overround | — |
| NUM-PROB-017 | NormalizedImpliedProbabilities | — |
| NUM-PROB-018 | ExpectedValue | — |
| NUM-PROB-019 | DiscreteVariance | — |
| NUM-PROB-020 | ChebyshevBound | #13, #14 |
| NUM-PROB-021 | ChebyshevBoundEpsilon | #15, #16 |
| NUM-PROB-022 | MarkovBound | — |
| **Discrete** | | |
| NUM-PROB-023 / 024 / 025 | BernoulliPMF / BernoulliCDF / BernoulliMeanVariance | #25, #104, #105 |
| NUM-PROB-026 / 027 / 028 | BinomialPMF / BinomialCDF / BinomialMeanVariance | #47, #26, #106, #107 |
| NUM-PROB-029 / 030 / 031 | GeometricPMF / GeometricCDF / GeometricMeanVariance | #28, #29, #108, #109 |
| NUM-PROB-032 / 033 / 034 | NegativeBinomialPMF / NegativeBinomialCDF / NegativeBinomialMeanVariance | #31, #32, #112, #113 |
| NUM-PROB-035 / 036 / 037 | HypergeometricPMF / HypergeometricCDF / HypergeometricMeanVariance | #30, #110, #111; used by #41 (life/games) |
| NUM-PROB-038 / 039 / 040 | DiscreteUniformPMF / DiscreteUniformCDF / DiscreteUniformMeanVariance | #27, #114, #115; one die (#20) |
| NUM-PROB-041 / 042 / 043 | PoissonPMF / PoissonCDF / PoissonMeanVariance | — |
| **Continuous** | | |
| NUM-PROB-044 / 045 / 046 | UniformPDF / UniformCDF / UniformMeanVariance | #114 if vCalc's "Uniform" is the continuous one |
| NUM-PROB-047 | NormalPDF | #23, #35, #36, #37, #45, #101 |
| NUM-PROB-048 | NormalCDF | #102, #139; stats #97–#99 |
| NUM-PROB-049 | NormalSurvival | stats #98, #99 |
| NUM-PROB-050 | NormalIntervalProbability | stats #108, #109 |
| NUM-PROB-051 / 052 / 053 | LognormalPDF / LognormalCDF / LognormalMeanVariance | #39, #40 |
| NUM-PROB-054 / 055 / 056 | ExponentialPDF / ExponentialCDF / ExponentialMeanVariance | — |
| **Quantiles** | | |
| NUM-PROB-057 | GeometricQuantile | — |
| NUM-PROB-058 | DiscreteUniformQuantile | — |
| NUM-PROB-059 | UniformQuantile | — |
| NUM-PROB-060 | ExponentialQuantile | — |
| NUM-PROB-061 | NormalQuantile | — |
| NUM-PROB-062 | LognormalQuantile | #39 |
| **Held, no ID** | | |
| — | DecimalOddsFromFractional, DecimalOddsFromAmerican, AmericanOddsFromDecimal | SOURCE NOT FOUND (Q7) |
| — | FractionalOddsFromDecimal | LATER: it needs a bounded-denominator rational approximation, and its source is not found either |
| — | Shin's method for removing the overround | LATER: Shin 1993 needs a 1-D root find; basic normalisation (017) comes first |
| — | Student t, chi-square, F | LATER, with the inferential tests |
| — | binomial and Poisson quantiles | LATER: no cited closed or standard numeric form; a search on the CDF is possible later |

Mean and variance come back as one pair, `(mean, variance float64, err error)`, from `<Dist>MeanVariance`. They
are always derived together, and vCalc's separate "Expected Value(X)" and "Variance(X)" rows both point at the
pair (Q2). Random sampling is out of scope: everything here is pure. A sampler, if one is ever wanted, takes a
`rand.Source` argument (PENDING_SHARED.md §5).

## Step 3: research (62 functions)

Conventions shared by every entry (Q2 gives the API reasons):
- Plain functions, no Distribution interface.
- Probabilities are `float64` in [0, 1], counts are `int`, and every input is dimensionless unless stated.
- A continuous variable `x` carries the caller's unit. Location and scale parameters (μ, σ, a, b) use the same
  unit, and a rate uses its inverse.
- Proposed sentinel errors, each wrapped with the parameter name:
  - `ErrNotFinite`: a NaN or ±Inf input.
  - `ErrProbability`: a probability outside [0, 1].
  - `ErrOutOfRange`: a count or parameter outside its domain, e.g. k > n, σ ≤ 0, n < 0.
  - `ErrInconsistent`: inputs that cannot come from one probability space, e.g. P(A∩B) > P(A).
  - `ErrNotDistribution`: weights that are not a distribution (Q5).
  - `ErrOverflow`: finite inputs give a non-finite result.
- A PMF at an integer outside its support returns 0 with no error; that is the value, not a failure.

### ProbabilityFromCounts      (ID reserved: NUM-PROB-001)
Formula: P(E) = f / t, where f = the number of equally likely outcomes in E and t = the total number of outcomes.
Inputs and units: `favorable, total int`, 0 ≤ favorable ≤ total, total ≥ 1.
Output: `float64` in [0, 1].
Errors: total < 1, favorable < 0 or favorable > total → ErrOutOfRange.
Primary source: [Ross] §2.5 "Sample Spaces Having Equally Likely Outcomes" [CONFIRM IN PRINT].
Worked example: HAND-CHECKED: a fair die, E = {even}: f = 3, t = 6 → 0.5.
Variants: none.
Best approach: one division; float64(f)/float64(t) is exact for any int below 2⁵³. Time O(1), Space O(1).
Notes: it must NOT accept non-equally-likely outcomes; that assumption belongs to the caller and the doc comment says so.

### ComplementProbability      (ID reserved: NUM-PROB-002)
Formula: P(Aᶜ) = 1 − P(A).
Inputs and units: `p float64` in [0, 1].   Output: `float64` in [0, 1].
Errors: ErrNotFinite; ErrProbability.
Primary source: [Ross] §2.4, Proposition 4.1; [CB] Theorem 1.2.8(c) [both CONFIRM IN PRINT].
Worked example: HAND-CHECKED: 0.25 → 0.75.
Variants: none.   Best approach: one subtraction (exact for p ∈ [0.5, 1] by Sterbenz). Time O(1), Space O(1).
Notes: a tiny complement of a p near 1 cannot be recovered here; callers needing upper tails use the survival functions (049).

### IndependentJointProbability      (ID reserved: NUM-PROB-003)
Formula: P(A ∩ B) = P(A)·P(B), for independent A and B.
Inputs and units: `pA, pB float64` in [0, 1].   Output: `float64` in [0, 1].
Errors: ErrNotFinite; ErrProbability.
Primary source: [Ross] §3.4, eq. (4.1); [CB] Definition 1.3.7 [both CONFIRM IN PRINT].
Worked example: HAND-CHECKED: two fair coins both heads: 0.5 · 0.5 = 0.25.
Variants: n events: a product over a slice, LATER if asked.   Best approach: one multiply. O(1), O(1).
Notes: it must NOT check independence; it cannot be checked from two numbers.

### JointProbability      (ID reserved: NUM-PROB-004)
Formula: P(A ∩ B) = P(A | B)·P(B), the multiplication rule for dependent events.
Inputs and units: `pAGivenB, pB float64` in [0, 1].   Output: `float64` in [0, 1].
Errors: ErrNotFinite; ErrProbability.
Primary source: [Ross] §3.2, eq. (2.2) (multiplication rule) [CONFIRM IN PRINT].
Worked example: HAND-CHECKED: draw two cards without replacement, both aces: P(B) = 4/52 and P(A|B) = 3/51,
  product 12/2652 = 1/221 ≈ 0.004524886877828055.
Variants: none.   Best approach: one multiply. O(1), O(1).
Notes: none.

### UnionProbability      (ID reserved: NUM-PROB-005)
Formula: P(A ∪ B) = P(A) + P(B) − P(A ∩ B). Mutually exclusive events pass P(A ∩ B) = 0.
Inputs and units: `pA, pB, pAAndB float64` in [0, 1].   Output: `float64` in [0, 1].
Errors: ErrNotFinite; ErrProbability; pAAndB > min(pA, pB) or pA + pB − pAAndB > 1 → ErrInconsistent.
Primary source: [Ross] §2.4, Proposition 4.3; [CB] Theorem 1.2.9(b) [both CONFIRM IN PRINT].
Worked example: HAND-CHECKED: one card, A = heart (13/52), B = king (4/52), A∩B = 1/52 → 16/52 = 4/13.
Variants: inclusion–exclusion for n events LATER.
Best approach: (pA − pAAndB) + pB, which subtracts the smaller term first. O(1), O(1).
Notes: the consistency checks use exact comparisons; no tolerance is invented.

### ConditionalProbability      (ID reserved: NUM-PROB-006)
Formula: P(A | B) = P(A ∩ B) / P(B), P(B) > 0.
Inputs and units: `pAAndB, pB float64` in [0, 1].   Output: `float64` in [0, 1].
Errors: pB = 0 → ErrOutOfRange (undefined); pAAndB > pB → ErrInconsistent; ErrNotFinite; ErrProbability.
Primary source: [Ross] §3.2, eq. (2.1); [CB] Definition 1.3.2 [both CONFIRM IN PRINT].
Worked example: HAND-CHECKED: one die, A = {6}, B = {even}: (1/6)/(1/2) = 1/3.
Variants: none.   Best approach: one division. O(1), O(1).
Notes: it must NOT return NaN for 0/0.

### TotalProbability      (ID reserved: NUM-PROB-007)
Formula: P(A) = Σᵢ P(A | Bᵢ)·P(Bᵢ), over a partition B₁…Bₙ (Σ P(Bᵢ) = 1).
Inputs and units: `pAGivenB, pB []float64`, equal length n ≥ 1, each in [0, 1]; pB must be a distribution.
Output: `float64` in [0, 1].
Errors: length mismatch or n = 0 → ErrOutOfRange; ErrProbability; ErrNotFinite; pB not summing to 1 → ErrNotDistribution (Q5).
Primary source: [Ross] §3.3 (law of total probability) [CONFIRM IN PRINT]; [CB] Theorem 1.2.11(a) [CONFIRM IN PRINT].
Worked example: HAND-CHECKED: two urns chosen with P = 0.5 each, red with P = 0.3 and 0.6 → 0.45.
Variants: none.
Best approach: products added with Neumaier compensated summation (Higham, *Accuracy and Stability of
  Numerical Algorithms*, 2nd ed., SIAM, 2002, §4.3; already a calcbench source). Time O(n), Space O(1).
Notes: the result is clamped to [0, 1] only if it overshoots by rounding within the Q5 bound; a larger overshoot is an error.

### BayesTheorem      (ID reserved: NUM-PROB-008)
Formula: P(A | B) = P(B | A)·P(A) / P(B), P(B) > 0.
Inputs and units: `pBGivenA, pA, pB float64` in [0, 1].   Output: `float64` in [0, 1].
Errors: pB = 0 → ErrOutOfRange; a result > 1 (that is, P(B|A)·P(A) > P(B)) → ErrInconsistent; ErrProbability; ErrNotFinite.
Primary source: [Ross] §3.3, eq. (3.1) and Bayes's formula eq. (3.4) [CONFIRM IN PRINT]; [CB] Theorem 1.3.5 (Bayes' Rule) [CONFIRM IN PRINT].
Worked example: HAND-CHECKED: P(B|A) = 0.9, P(A) = 0.01, P(B) = 0.0585 → 0.009/0.0585 = 0.15384615384615383 in float64 (= 2/13).
Variants: the partition form is 009.   Best approach: multiply, then divide. O(1), O(1).
Notes: none.

### BayesPosterior      (ID reserved: NUM-PROB-009)
Formula: P(Aⱼ | B) = P(B | Aⱼ)·P(Aⱼ) / Σᵢ P(B | Aᵢ)·P(Aᵢ), over a partition A₁…Aₙ.
Inputs and units: `likelihoods, priors []float64`, equal length n ≥ 1, each in [0, 1]; priors must be a distribution (Q5).
Output: `[]float64`, the n posteriors (sum 1).
Errors: as 007; a zero denominator (B impossible under every Aᵢ) → ErrOutOfRange.
Primary source: as 008, the general form (Bayes's formula over a partition).
Worked example: HAND-CHECKED: priors [0.5, 0.5], likelihoods [0.3, 0.6] → [1/3, 2/3].
Variants: a `dst []float64` to avoid allocation, like the S1 window functions.
Best approach: one compensated sum for the denominator, then n divisions. Time O(n); Space O(n) for the result, one allocation.
Notes: none.

### PositivePredictiveValue      (ID reserved: NUM-PROB-010)
Formula: PPV = Se·π / (Se·π + (1 − Sp)·(1 − π)), where Se = sensitivity, Sp = specificity, π = prevalence
  (pre-test probability). This is Bayes' theorem with total probability in the denominator.
Inputs and units: `sensitivity, specificity, prevalence float64`, each in [0, 1].
Output: `float64` in [0, 1], P(disease | positive test).
Errors: ErrProbability; ErrNotFinite; a zero denominator (Se·π = 0 and (1−Sp)(1−π) = 0) → ErrOutOfRange.
Primary source: Altman, Bland, "Statistics Notes: Diagnostic tests 2: predictive values", BMJ, 1994,
  309(6947):102, doi:10.1136/bmj.309.6947.102 [DOI CHECKED; article text not readable online, formula CONFIRM IN PRINT].
  Second: as 008.
Worked example: the liver-scan data of Altman and Bland (BMJ 1994, 308:1552 and 309:102). The PPV of
  231/263 = 0.88 was confirmed from a search excerpt of the article. The 2×2 counts (TP 231, FP 32, FN 27,
  TN 54) are CONFIRM IN PRINT. Then Se = 231/258, Sp = 54/86, π = 258/344, and PPV = 231/263 =
  0.8783269961977186 exactly. HAND-CHECKED with exact fractions.
Variants: likelihood ratios LR+ = Se/(1−Sp) and LR− = (1−Se)/Sp are LATER, if the medical subject asks.
Best approach: one formula. O(1), O(1).
Notes: it states a number; it must NOT label a result (no "good test" wording). The medical subject calls this function.

### NegativePredictiveValue      (ID reserved: NUM-PROB-011)
Formula: NPV = Sp·(1 − π) / (Sp·(1 − π) + (1 − Se)·π).
Inputs, output, errors, source: as 010.
Worked example: the same liver-scan data: NPV = 54/81 = 2/3 = 0.6666666666666666. HAND-CHECKED with exact fractions.
Variants: as 010.   Best approach: O(1), O(1).   Notes: as 010.

### OddsFromProbability      (ID reserved: NUM-PROB-012)
Formula: odds in favour o = p / (1 − p).
Inputs and units: `p float64` in [0, 1).   Output: `float64` ≥ 0 (the "o to 1" odds in favour).
Errors: p = 1 → ErrOutOfRange (infinite odds); ErrProbability; ErrNotFinite.
Primary source: [Ross] §3.3, the definition of the odds of an event [CONFIRM IN PRINT].
Worked example: HAND-CHECKED: p = 0.75 → 3 (3 to 1).
Variants: odds against = (1 − p)/p, as the caller's 1/o.   Best approach: one division. O(1), O(1).
Notes: none.

### ProbabilityFromOdds      (ID reserved: NUM-PROB-013)
Formula: p = o / (1 + o).
Inputs and units: `odds float64` ≥ 0.   Output: `float64` in [0, 1).
Errors: odds < 0 → ErrOutOfRange; ErrNotFinite.
Primary source: as 012.
Worked example: HAND-CHECKED: 3 → 0.75.
Variants: none.
Best approach: o/(1+o); for o > 2⁵³ the result is 1 − 1/(1+o), which rounds to 1, and that is correct. O(1), O(1).
Notes: none.

### DecimalOddsFromProbability      (ID reserved: NUM-PROB-014)
Formula: d = 1 / p, the fair decimal (European) odds: the total returned per unit staked.
Inputs and units: `p float64` in (0, 1].   Output: `float64` ≥ 1.
Errors: p = 0 → ErrOutOfRange; ErrProbability; ErrNotFinite.
Primary source: Štrumbelj, "On determining probability forecasts from betting odds", International Journal of
  Forecasting, 2014, 30(4):934–943 [journal, volume and pages CHECKED; DOI 10.1016/j.ijforecast.2014.02.008 CONFIRM],
  the definition of decimal odds and inverse odds [section CONFIRM IN PRINT].
Worked example: HAND-CHECKED: p = 0.4 → 2.5.
Variants: none.   Best approach: one division. O(1), O(1).
Notes: these are fair odds; a bookmaker's quoted odds include an overround (016).

### ImpliedProbability      (ID reserved: NUM-PROB-015)
Formula: π = 1 / d, the probability implied by decimal odds d, still including the overround.
Inputs and units: `decimalOdds float64` ≥ 1.   Output: `float64` in (0, 1].
Errors: d < 1 → ErrOutOfRange; ErrNotFinite.
Primary source: as 014 (inverse odds).
Worked example: HAND-CHECKED: 2.5 → 0.4; 1.9 → 0.5263157894736842.
Variants: overround removed: 017.   Best approach: O(1), O(1).
Notes: the sports subject calls this.

### Overround      (ID reserved: NUM-PROB-016)
Formula: booksum B = Σᵢ 1/dᵢ over all outcomes of one market; overround = B − 1.
Inputs and units: `decimalOdds []float64`, n ≥ 2, each ≥ 1.   Output: `float64`, the overround as a fraction
  (0.05 = 5 %). It can be negative for an arbitrage market.
Errors: n < 2 → ErrOutOfRange; any d < 1 → ErrOutOfRange; ErrNotFinite.
Primary source: as 014 (booksum, the sum of inverse odds) [section CONFIRM IN PRINT].
Worked example: HAND-CHECKED: [1.9, 1.9] → 2/1.9 − 1 = 0.05263157894736836.
Variants: none.   Best approach: compensated sum. Time O(n), Space O(1).
Notes: it must NOT assume the market is complete; the caller passes every outcome.

### NormalizedImpliedProbabilities      (ID reserved: NUM-PROB-017)
Formula: πᵢ = (1/dᵢ) / Σⱼ (1/dⱼ), the "basic normalisation" that removes the overround proportionally.
Inputs and units: as 016.   Output: `[]float64`, n probabilities summing to 1.
Errors: as 016.
Primary source: as 014; Štrumbelj compares basic normalisation with Shin's method and finds Shin's more
  accurate. Shin's method is LATER (Shin, "Measuring the incidence of insider trading in a market for
  state-contingent claims", Economic Journal, 1993, 103(420):1141–1153 [DOI CONFIRM]).
Worked example: HAND-CHECKED: [1.9, 1.9] → [0.5, 0.5]; [1.5, 2.75] → (0.6667, 0.3636)/1.0303 = [0.6470588235294118, 0.35294117647058826].
Variants: Shin's method LATER; a `dst` slice.
Best approach: one compensated sum, then n divisions. Time O(n); Space O(n) for the result.
Notes: the doc comment says this is the proportional method and that it is known to be biased toward longshots (per the source).

### ExpectedValue      (ID reserved: NUM-PROB-018)
Formula: E[X] = Σᵢ xᵢ·pᵢ, a discrete distribution with outcomes xᵢ and probabilities pᵢ.
Inputs and units: `outcomes []float64` (the caller's unit), `probabilities []float64` (each in [0, 1], a
  distribution, Q5), equal length n ≥ 1.
Output: `float64`, in the unit of the outcomes.
Errors: length mismatch or n = 0 → ErrOutOfRange; ErrProbability; ErrNotDistribution; ErrNotFinite; ErrOverflow.
Primary source: [Ross] §4.3, the definition of expectation, and Example 3a (a fair die) [CONFIRM IN PRINT].
Worked example: [Ross] §4.3 Example 3a: a fair die → 3.5 [example number CONFIRM IN PRINT].
Variants: none.
Best approach: products with `math.FMA` accumulated by Neumaier compensated summation (Higham §4.3). Time O(n), Space O(1).
Notes: the life/games and sports subjects call this for bet expectations; money stays in the caller's integer cents.

### DiscreteVariance      (ID reserved: NUM-PROB-019)
Formula: Var X = Σᵢ pᵢ·(xᵢ − μ)², μ = E[X] from 018.
Inputs, errors: as 018.   Output: `float64` ≥ 0, in the unit squared.
Primary source: [Ross] §4.5, the definition of variance, and Example 5a (a die: 35/12) [CONFIRM IN PRINT].
Worked example: [Ross] §4.5 Example 5a: a fair die → 35/12 = 2.9166666666666665 [CONFIRM IN PRINT].
Variants: the standard deviation is the caller's `math.Sqrt`.
Best approach: two passes (the mean, then the centred sum), both compensated. It must NOT use E[X²] − μ²,
  which cancels catastrophically (Higham §1.9). Time O(n), Space O(1).
Notes: none.

### ChebyshevBound      (ID reserved: NUM-PROB-020)
Formula: P(|X − μ| ≥ kσ) ≤ 1/k², for any X with finite variance and k > 0. Returns min(1, 1/k²).
Inputs and units: `k float64` > 0 (the number of standard deviations).   Output: `float64` in (0, 1], the upper bound.
Errors: k ≤ 0 → ErrOutOfRange; ErrNotFinite.
Primary source: [Ross] §8.2, Proposition 2.2 (Chebyshev's inequality) [CONFIRM IN PRINT]; [CB] Theorem 3.6.1 [CONFIRM IN PRINT].
Worked example: HAND-CHECKED: k = 2 → 0.25; k = 0.5 → 1 (the bound is trivial below k = 1).
Variants: the ε form is 021.   Best approach: O(1), O(1).
Notes: it is a bound, not a probability; the doc comment says so.

### ChebyshevBoundEpsilon      (ID reserved: NUM-PROB-021)
Formula: P(|X − μ| ≥ ε) ≤ σ²/ε², returned as min(1, σ²/ε²).
Inputs and units: `variance float64` ≥ 0 (the unit squared), `epsilon float64` > 0 (the unit).   Output: `float64` in [0, 1].
Errors: variance < 0 or epsilon ≤ 0 → ErrOutOfRange; ErrNotFinite.
Primary source: as 020.
Worked example: [Ross] §8.2 Example 2a(b): mean 50, variance 25, ε = 10 → 25/100 = 0.25, so
  P(40 < X < 60) ≥ 0.75 [example CONFIRM IN PRINT].
Variants: none.   Best approach: (σ/ε)² if σ²/ε² would overflow. O(1), O(1).
Notes: none.

### MarkovBound      (ID reserved: NUM-PROB-022)
Formula: P(X ≥ a) ≤ E[X]/a, for X ≥ 0 and a > 0. Returns min(1, E[X]/a).
Inputs and units: `mean float64` ≥ 0, `a float64` > 0 (same unit).   Output: `float64` in [0, 1].
Errors: mean < 0 or a ≤ 0 → ErrOutOfRange; ErrNotFinite.
Primary source: [Ross] §8.2, Proposition 2.1 (Markov's inequality) [CONFIRM IN PRINT].
Worked example: [Ross] §8.2 Example 2a(a): mean 50, a = 75 → 2/3 = 0.6666666666666666 [CONFIRM IN PRINT].
Variants: none.   Best approach: O(1), O(1).   Notes: valid only for a non-negative X; the caller asserts it.

### BernoulliPMF      (ID reserved: NUM-PROB-023)
Formula: P(X = 1) = p, P(X = 0) = 1 − p; 0 elsewhere.
Inputs and units: `k int`, `p float64` in [0, 1].   Output: `float64`.
Errors: ErrProbability; ErrNotFinite.
Primary source: [Ross] §4.6; [CB] §3.2 (Bernoulli) [both CONFIRM IN PRINT].
Worked example: HAND-CHECKED: p = 0.3: k = 1 → 0.3, k = 0 → 0.7, k = 2 → 0.
Variants: none.   Best approach: O(1), O(1).   Notes: none.

### BernoulliCDF      (ID reserved: NUM-PROB-024)
Formula: F(k) = 0 for k < 0; 1 − p for 0 ≤ k < 1; 1 for k ≥ 1.
Inputs and units: `k int`, `p float64` in [0, 1].   Output: `float64`.
Errors, source: as 023.
Worked example: HAND-CHECKED: p = 0.3, k = 0 → 0.7.
Variants: none.   Best approach: O(1), O(1).   Notes: none.

### BernoulliMeanVariance      (ID reserved: NUM-PROB-025)
Formula: mean = p, variance = p(1 − p).
Inputs and units: `p float64` in [0, 1].   Output: `(mean, variance float64, err error)`.
Errors, source: as 023.
Worked example: HAND-CHECKED: p = 0.3 → (0.3, 0.21).
Variants: none.   Best approach: O(1), O(1).   Notes: none.

### BinomialPMF      (ID reserved: NUM-PROB-026)
Formula: P(X = k) = C(n, k)·p^k·(1 − p)^(n−k), k = 0…n (C(n, k) = n!/(k!(n−k)!)). Evaluated as Loader's
  saddle-point form: P = exp(−stirlerr(n) + stirlerr(k) + stirlerr(n−k) + ... − bd0(k, np) − bd0(n−k, n(1−p)))
  · √(n / (2π·k·(n−k))). Here stirlerr(m) = ln m! − ln(√(2πm)(m/e)^m) is the Stirling-series remainder, and
  bd0(x, M) = x ln(x/M) + M − x is the deviance term, computed without cancellation. The edge cases k = 0 and
  k = n are evaluated directly with `math.Log1p`.
Inputs and units: `k, n int`, n ≥ 0; `p float64` in [0, 1].   Output: `float64` in [0, 1].
Errors: n < 0 → ErrOutOfRange; ErrProbability; ErrNotFinite. k outside 0…n → 0, nil.
Primary source: [Loader] §2–3 (saddle-point algorithm, stirlerr and bd0) [sections CONFIRM on reading];
  the formula itself: [NIST] §1.3.6.6.18 [CHECKED]; [Ross] §4.6.
Worked example: HAND-CHECKED: n = 10, p = 0.5, k = 5 → 252/1024 = 0.24609375. Large n, where naive
  C(n,k)·p^k overflows or underflows: n = 1000, p = 0.5, k = 500 → 0.0252250181783608 (checked in the test
  against an exact `math/big` rational).
Variants: a log-PMF for very small values, LATER.
Best approach: Loader's method gives near machine precision for every n. The rejected alternative is
  exp(lgamma differences): `math.Lgamma` terms of size n·ln n cancel and lose about log10(n ln n) digits.
  Time O(1) (stirlerr is a table for m ≤ 15, a short series above); Space O(1).
Notes: it must NOT compute n! or C(n, k) directly; it must NOT return a value > 1 from rounding (clamp only
  at the last step, documented).

### BinomialCDF      (ID reserved: NUM-PROB-027)
Formula: F(k) = Σ_{i=0..k} C(n, i)·p^i·(1 − p)^(n−i).
Inputs and units: as 026.   Output: `float64` in [0, 1].
Errors: as 026; k < 0 → 0; k ≥ n → 1.
Primary source: [NIST] §1.3.6.6.18 (CDF) [CHECKED]; terms by [Loader].
Worked example: HAND-CHECKED: n = 10, p = 0.5, k = 5 → 638/1024 = 0.623046875.
Variants: an O(1)-in-n version through the regularized incomplete beta, F(k) = I_{1−p}(n − k, k + 1)
  (DLMF §8.17, https://dlmf.nist.gov/8.17 [eq. number CONFIRM]), waits on Q6.
Best approach: start at the PMF of the term nearest the mode on the requested side, then use the ratio
  recurrence P(i−1) = P(i)·i(1−p)/((n−i+1)p) and add the terms smallest-last with compensated summation.
  Stop when a term falls below u·sum (u = 2⁻⁵³). If k is above the mean, sum the upper tail instead and
  return 1 − tail, which has no cancellation there. Time O(min(k+1, n−k, c·√(np(1−p)))) terms, c ≈ 40 for
  double precision; Space O(1).
Notes: it must NOT sum from i = 0 in naive order for large n (underflowing leading terms).

### BinomialMeanVariance      (ID reserved: NUM-PROB-028)
Formula: mean = np, variance = np(1 − p).
Inputs, errors: as 026.   Output: `(mean, variance float64, err error)`.
Primary source: [NIST] §1.3.6.6.18 (mean np, standard deviation √(np(1−p))) [CHECKED]; [Ross] §4.6.1.
Worked example: HAND-CHECKED: n = 10, p = 0.5 → (5, 2.5).
Variants: none.   Best approach: O(1), O(1).   Notes: none.

### GeometricPMF      (ID reserved: NUM-PROB-029)
Formula (TRIALS convention, Q1): X = the number of the trial on which the first success occurs,
  P(X = k) = (1 − p)^(k−1)·p, k = 1, 2, …
Inputs and units: `k int`, `p float64` in (0, 1].   Output: `float64`.
Errors: p = 0 → ErrOutOfRange; ErrProbability; ErrNotFinite. k < 1 → 0.
Primary source: [Ross] §4.8.1; [CB] §3.2 (geometric, number of trials) [both CONFIRM IN PRINT].
Worked example: HAND-CHECKED: p = 0.2, k = 3 → 0.8²·0.2 = 0.128.
Variants: the failures convention (support 0, 1, …) is X − 1; a separate `GeometricFailuresPMF` only if Q1 says so.
Best approach: exp((k−1)·log1p(−p))·p, which is accurate for small p. O(1), O(1).
Notes: none.

### GeometricCDF      (ID reserved: NUM-PROB-030)
Formula: F(k) = 1 − (1 − p)^k, k ≥ 1.
Inputs, errors, source: as 029.   Output: `float64`.
Worked example: HAND-CHECKED: p = 0.2, k = 3 → 1 − 0.512 = 0.488.
Variants: as 029.   Best approach: −expm1(k·log1p(−p)), with no cancellation for small p. O(1), O(1).
Notes: none.

### GeometricMeanVariance      (ID reserved: NUM-PROB-031)
Formula: mean = 1/p, variance = (1 − p)/p².
Inputs, errors, source: as 029.   Output: `(mean, variance float64, err error)`.
Worked example: HAND-CHECKED: p = 0.2 → (5, 20).
Variants: failures convention: mean (1 − p)/p, same variance.   Best approach: O(1), O(1).   Notes: none.

### NegativeBinomialPMF      (ID reserved: NUM-PROB-032)
Formula (Pascal, TRIALS convention, matching the geometric): X = the trial on which the r-th success occurs,
  P(X = x) = C(x−1, r−1)·p^r·(1 − p)^(x−r), x = r, r+1, … Evaluated as (r/x)·BinomialPMF(r; x, p), using the
  identity C(x−1, r−1) = (r/x)·C(x, r) (HAND-DERIVED, checked symbolically).
Inputs and units: `x, r int`, r ≥ 1; `p float64` in (0, 1].   Output: `float64`.
Errors: r < 1 or p = 0 → ErrOutOfRange; ErrProbability; ErrNotFinite. x < r → 0.
Primary source: [Ross] §4.8.2 (negative binomial); [CB] §3.2 [both CONFIRM IN PRINT]; [Loader] for the binomial term.
Worked example: HAND-CHECKED: r = 2, p = 0.5, x = 4 → C(3,1)·0.25·0.25 = 0.1875.
Variants: the failures form Y = X − r; a real-valued r (Pólya) LATER.
Best approach: one Loader binomial term. O(1), O(1).
Notes: none.

### NegativeBinomialCDF      (ID reserved: NUM-PROB-033)
Formula: P(X ≤ x) = P(at least r successes in x trials) = Σ_{i=r..x} C(x, i)·p^i·(1 − p)^(x−i) (HAND-DERIVED from the definition).
Inputs, errors, source: as 032.   Output: `float64`.
Worked example: HAND-CHECKED: r = 2, p = 0.5, x = 4 → 0.25 + 0.25 + 0.1875 = 0.6875.
Variants: the incomplete-beta form waits on Q6.
Best approach: the upper-tail summation of 027 (binomial terms by the ratio recurrence, compensated). Time and Space as 027.
Notes: none.

### NegativeBinomialMeanVariance      (ID reserved: NUM-PROB-034)
Formula: mean = r/p, variance = r(1 − p)/p².
Inputs, errors, source: as 032.   Output: `(mean, variance float64, err error)`.
Worked example: HAND-CHECKED: r = 2, p = 0.5 → (4, 4).
Variants: failures convention: mean r(1−p)/p, same variance.   Best approach: O(1), O(1).   Notes: none.

### HypergeometricPMF      (ID reserved: NUM-PROB-035)
Formula: P(X = k) = C(K, k)·C(N−K, n−k) / C(N, n): k successes in n draws without replacement from N items,
  K of which are successes; k from max(0, n+K−N) to min(n, K). Evaluated as
  b(k; K, q)·b(n−k; N−K, q) / b(n; N, q) with q = n/N, where b is the Loader binomial PMF. The identity holds
  for any q in (0, 1), because the powers of q cancel (HAND-DERIVED, checked symbolically).
Inputs and units: `k, N, K, n int`, 0 ≤ K ≤ N, 0 ≤ n ≤ N.   Output: `float64`.
Errors: N < 0, K or n outside [0, N] → ErrOutOfRange. k outside its support → 0.
Primary source: [Ross] §4.8.3 (hypergeometric); [CB] §3.2 [both CONFIRM IN PRINT]; terms by [Loader].
Worked example: HAND-CHECKED: one ace in a 5-card hand: N = 52, K = 4, n = 5, k = 1 →
  4·194580/2598960 = 3243/10829 = 0.2994736356080894.
Variants: none.
Best approach: three Loader terms, O(1), with no cancellation. The alternative, a difference of
  log-binomial coefficients (lgamma), loses digits as in 026. Space O(1).
Notes: the life/games subject calls this for card-hand probabilities (row #41, poker rows).

### HypergeometricCDF      (ID reserved: NUM-PROB-036)
Formula: F(k) = Σ over i from the lower support bound to k of P(X = i).
Inputs, errors, source: as 035.   Output: `float64`.
Worked example: HAND-CHECKED: N = 52, K = 4, n = 5, k = 1 → (C(48,5) + 4·C(48,4))/C(52,5) = 0.958315633945886.
Variants: none.
Best approach: one PMF at the nearest end, then the ratio recurrence
  P(i+1)/P(i) = (K−i)(n−i)/((i+1)(N−K−n+i+1)), compensated; sum the shorter side. Time O(min(k − lo + 1, hi − k)), Space O(1).
Notes: none.

### HypergeometricMeanVariance      (ID reserved: NUM-PROB-037)
Formula: mean = nK/N; variance = n·(K/N)·((N−K)/N)·((N−n)/(N−1)), with variance 0 when N = 1.
Inputs, errors, source: as 035.   Output: `(mean, variance float64, err error)`.
Worked example: HAND-CHECKED: N = 52, K = 4, n = 5 → (0.38461538461538464, 940/2873 = 0.32718412808910546).
Variants: none.   Best approach: O(1), O(1); every factor is a ratio in [0, 1], so nothing overflows.   Notes: N = 0 → ErrOutOfRange.

### DiscreteUniformPMF      (ID reserved: NUM-PROB-038)
Formula: P(X = k) = 1/(b − a + 1) for integers a ≤ k ≤ b; 0 elsewhere.
Inputs and units: `k, a, b int`, a ≤ b.   Output: `float64`.
Errors: a > b → ErrOutOfRange; b − a + 1 overflowing int → ErrOverflow.
Primary source: [CB] §3.2 (discrete uniform on 1…N) [CONFIRM IN PRINT], shifted to a…b (HAND-DERIVED).
Worked example: HAND-CHECKED: a fair die (1, 6), k = 4 → 1/6.
Variants: none.   Best approach: O(1), O(1).   Notes: one die of the DiceRoll row (#20).

### DiscreteUniformCDF      (ID reserved: NUM-PROB-039)
Formula: F(k) = (k − a + 1)/(b − a + 1) for a ≤ k ≤ b; 0 below a; 1 above b.
Inputs, errors, source: as 038.   Output: `float64`.
Worked example: HAND-CHECKED: die, k = 4 → 4/6 = 0.6666666666666666.
Variants: none.   Best approach: O(1), O(1).   Notes: none.

### DiscreteUniformMeanVariance      (ID reserved: NUM-PROB-040)
Formula: mean = (a + b)/2; variance = ((b − a + 1)² − 1)/12.
Inputs, errors, source: as 038.   Output: `(mean, variance float64, err error)`.
Worked example: HAND-CHECKED: die → (3.5, 35/12 = 2.9166666666666665), the same values as Ross's die (019).
Variants: none.
Best approach: compute in float64, with mean = a/2 + b/2 + corrections, so that no int overflow occurs. O(1), O(1).
Notes: none.

### PoissonPMF      (ID reserved: NUM-PROB-041)
Formula: P(X = k) = e^(−λ)·λ^k / k!, k = 0, 1, … Evaluated as Loader's form
  exp(−stirlerr(k) − bd0(k, λ)) / √(2πk) for k ≥ 1, and e^(−λ) for k = 0.
Inputs and units: `k int`, `lambda float64` ≥ 0 (the mean count per interval).   Output: `float64`.
Errors: λ < 0 → ErrOutOfRange; ErrNotFinite. k < 0 → 0. λ = 0: P(0) = 1.
Primary source: [NIST] §1.3.6.6.19 [CHECKED]; [Loader] (Poisson by the same saddle point) [section CONFIRM on reading]; [Ross] §4.7.
Worked example: HAND-CHECKED: λ = 2, k = 0 → e⁻² = 0.1353352832366127. λ = 1000, k = 1000 →
  0.012614611348719664 (the naive λ^k/k! overflows).
Variants: none.   Best approach: Loader, O(1), O(1).   Notes: it must NOT form λ^k or k!.

### PoissonCDF      (ID reserved: NUM-PROB-042)
Formula: F(k) = Σ_{i=0..k} e^(−λ) λ^i / i!.
Inputs, errors, source: as 041 ([NIST] §1.3.6.6.19 gives the sum).   Output: `float64`.
Worked example: HAND-CHECKED: λ = 2, k = 2 → 5e⁻² = 0.6766764161830635.
Variants: the regularized upper incomplete gamma Q(k+1, λ) (DLMF §8.2, https://dlmf.nist.gov/8.2 [eq. CONFIRM]) waits on Q6.
Best approach: as 027, using the ratio P(i−1) = P(i)·i/λ, compensated, and summing the short side. Time
  O(min(k+1, c·√λ + |k − λ|)), Space O(1).
Notes: none.

### PoissonMeanVariance      (ID reserved: NUM-PROB-043)
Formula: mean = λ, variance = λ.
Inputs, errors, source: as 041.   Output: `(mean, variance float64, err error)`.
Worked example: HAND-CHECKED: λ = 2 → (2, 2).
Variants: none.   Best approach: O(1), O(1).   Notes: none.

### UniformPDF      (ID reserved: NUM-PROB-044)
Formula: f(x) = 1/(b − a) for a ≤ x ≤ b; 0 elsewhere.
Inputs and units: `x, a, b float64` (one unit), a < b.   Output: `float64`, per unit of x.
Errors: a ≥ b → ErrOutOfRange; ErrNotFinite; b − a overflowing → ErrOverflow.
Primary source: [NIST] §1.3.6.6.2 [section inferred, CONFIRM]; [Ross] §5.3.
Worked example: HAND-CHECKED: a = 0, b = 10, x = 3 → 0.1.
Variants: none.   Best approach: O(1), O(1).   Notes: none.

### UniformCDF      (ID reserved: NUM-PROB-045)
Formula: F(x) = (x − a)/(b − a), clamped to [0, 1] outside [a, b].
Inputs, errors, source: as 044.   Output: `float64`.
Worked example: HAND-CHECKED: a = 0, b = 10, x = 3 → 0.3.
Variants: none.   Best approach: O(1), O(1).   Notes: none.

### UniformMeanVariance      (ID reserved: NUM-PROB-046)
Formula: mean = (a + b)/2; variance = (b − a)²/12.
Inputs, errors, source: as 044.   Output: `(mean, variance float64, err error)`.
Worked example: HAND-CHECKED: a = 0, b = 10 → (5, 8.333333333333334).
Variants: none.   Best approach: a + (b − a)/2 for the mean, so that a + b cannot overflow. O(1), O(1).   Notes: none.

### NormalPDF      (ID reserved: NUM-PROB-047)
Formula: f(x) = exp(−z²/2) / (σ√(2π)), z = (x − μ)/σ. The standard normal is μ = 0, σ = 1.
Inputs and units: `x, mean float64` (one unit), `sd float64` > 0 (the same unit).   Output: `float64`, per unit of x.
Errors: sd ≤ 0 → ErrOutOfRange; ErrNotFinite.
Primary source: [NIST] §1.3.6.6.1 [section inferred, CONFIRM]; [Ross] §5.4.
Worked example: HAND-CHECKED: x = 0, μ = 0, σ = 1 → 1/√(2π) = 0.3989422804014327.
Variants: none (every Gaussian row is this function).
Best approach: O(1), O(1); for |z| > ~38.6 the result underflows to 0, which is correct.
Notes: it must NOT compute (x−μ)² before dividing by σ (overflow); compute z first.

### NormalCDF      (ID reserved: NUM-PROB-048)
Formula: Φ(z) = ½·erfc(−z/√2), z = (x − μ)/σ. This is the definition of Φ rewritten with the complementary
  error function (DLMF eqs. 7.2.1–7.2.2, https://dlmf.nist.gov/7.2, by the substitution t = u/√2; HAND-DERIVED).
Inputs, errors: as 047.   Output: `float64` in [0, 1].
Primary source: [NIST] §1.3.6.6.1; DLMF §7.2 [CONFIRM eq. numbers]; [Ross] §5.4.
Worked example: [NIST] §1.3.6.7.1 table: area 0 to 1.96 = 0.47500, so Φ(1.96) = 0.97500 [CHECKED]. Float64:
  0.9750021048517795, and Φ(−1.96) = 0.02499789514822043.
Variants: upper tail is 049; an interval is 050.
Best approach: `math.Erfc` for every z, which is accurate in the lower tail where ½(1 + erf) would cancel. O(1), O(1).
Notes: it must NOT use 1 − Φ for upper tails (049).

### NormalSurvival      (ID reserved: NUM-PROB-049)
Formula: Q(z) = 1 − Φ(z) = ½·erfc(z/√2).
Inputs, errors, source: as 048.   Output: `float64` in [0, 1].
Worked example: HAND-CHECKED: z = 1.96 → 0.02499789514822043; z = 10 → 7.619853024160593e-24 (1 − Φ(10) gives 0).
Variants: none.   Best approach: `math.Erfc`. O(1), O(1).   Notes: the stats rows "P(y) right of z" call this.

### NormalIntervalProbability      (ID reserved: NUM-PROB-050)
Formula: P(lo ≤ X ≤ hi) = Φ(z_hi) − Φ(z_lo), computed as Q(z_lo) − Q(z_hi) when both z are above 0.
Inputs and units: `lo, hi, mean float64` (one unit), `sd float64` > 0.   Output: `float64` in [0, 1].
Errors: lo > hi → ErrOutOfRange; as 047.
Primary source: as 048.
Worked example: HAND-CHECKED: ±1.96 → 0.950004209703559; z from 8 to 9 → 6.219831985865867e-16 (Φ − Φ gives 0).
Variants: in z units: the caller passes mean 0, sd 1.
Best approach: pick the tail by sign, so that the subtraction is between small numbers. O(1), O(1).
Notes: none.

### LognormalPDF      (ID reserved: NUM-PROB-051)
Formula: f(x) = exp(−(ln x − μ)²/(2σ²)) / (xσ√(2π)), x > 0. μ and σ are the mean and standard deviation of
  ln X (NIST's alternative form, μ = ln m, where m is the scale and median).
Inputs and units: `x float64` > 0 (the caller's unit); `mu float64` (log of that unit); `sigma float64` > 0
  (dimensionless).   Output: `float64`, per unit of x.
Errors: sigma ≤ 0 → ErrOutOfRange; ErrNotFinite. x ≤ 0 → 0.
Primary source: [NIST] §1.3.6.6.9 (both parameterisations given; NIST prefers m) [CHECKED].
Worked example: HAND-CHECKED: μ = 0, σ = 1, x = 1 → 0.3989422804014327.
Variants: (σ, m) with m = e^μ: the caller passes math.Log(m) (Q2).   Best approach: O(1), O(1).   Notes: none.

### LognormalCDF      (ID reserved: NUM-PROB-052)
Formula: F(x) = Φ((ln x − μ)/σ), x > 0; 0 for x ≤ 0.
Inputs, errors, source: as 051 ([NIST] §1.3.6.6.9 gives Φ(ln x / σ) for the standard form).   Output: `float64`.
Worked example: HAND-CHECKED: μ = 0, σ = 1, x = 1 → 0.5.
Variants: none.   Best approach: calls 048. O(1), O(1).   Notes: none.

### LognormalMeanVariance      (ID reserved: NUM-PROB-053)
Formula: mean = exp(μ + σ²/2); variance = expm1(σ²)·exp(2μ + σ²).
Inputs, errors, source: as 051. NIST gives the mean e^(σ²/2) and SD √(e^(σ²)(e^(σ²) − 1)) for m = 1; the
  general μ follows by scaling (HAND-DERIVED).   Output: `(mean, variance float64, err error)`.
Worked example: HAND-CHECKED: μ = 0, σ = 1 → (1.6487212707001282, 4.670774270471604).
Variants: none.   Best approach: `math.Expm1` for small σ. Overflow → ErrOverflow. O(1), O(1).   Notes: none.

### ExponentialPDF      (ID reserved: NUM-PROB-054)
Formula: f(x) = λ·e^(−λx), x ≥ 0, where λ is the RATE (NIST uses the scale β = 1/λ; Q2).
Inputs and units: `x float64` ≥ 0 (a time or length unit); `rate float64` > 0 (per that unit).   Output: `float64`.
Errors: rate ≤ 0 → ErrOutOfRange; ErrNotFinite. x < 0 → 0.
Primary source: [NIST] §1.3.6.6.7 [CHECKED]; [Ross] §5.5 (rate λ) [CONFIRM IN PRINT].
Worked example: HAND-CHECKED: λ = 2, x = 1 → 2e⁻² = 0.2706705664732254.
Variants: the scale form, by the caller's 1/β.   Best approach: O(1), O(1).   Notes: none.

### ExponentialCDF      (ID reserved: NUM-PROB-055)
Formula: F(x) = 1 − e^(−λx) = −expm1(−λx).
Inputs, errors, source: as 054 ([NIST] gives 1 − e^(−x/β)).   Output: `float64`.
Worked example: HAND-CHECKED: λ = 2, x = 1 → 0.8646647167633873.
Variants: none.   Best approach: `math.Expm1`, which is accurate for small λx. O(1), O(1).   Notes: none.

### ExponentialMeanVariance      (ID reserved: NUM-PROB-056)
Formula: mean = 1/λ; variance = 1/λ².
Inputs, errors, source: as 054 ([NIST]: mean β, SD β).   Output: `(mean, variance float64, err error)`.
Worked example: HAND-CHECKED: λ = 2 → (0.5, 0.25).
Variants: none.   Best approach: O(1), O(1).   Notes: none.

### GeometricQuantile      (ID reserved: NUM-PROB-057)
Formula: the smallest k ≥ 1 with F(k) ≥ u: k = ⌈ln(1 − u)/ln(1 − p)⌉ = ⌈log1p(−u)/log1p(−p)⌉, then corrected
  by ±1 so that F(k−1) < u ≤ F(k) holds exactly (HAND-DERIVED from 030).
Inputs and units: `u float64` in [0, 1); `p float64` in (0, 1].   Output: `(int, error)`.
Errors: u = 1 with p < 1 → ErrOutOfRange (no finite quantile); a k above MaxInt → ErrOverflow; as 029.
Primary source: as 029 (by inverting 030).
Worked example: HAND-CHECKED: p = 0.2, u = 0.5 → 4 (F(3) = 0.488 < 0.5 ≤ F(4) = 0.5904).
Variants: failures convention: k − 1.   Best approach: O(1), O(1).   Notes: none.

### DiscreteUniformQuantile      (ID reserved: NUM-PROB-058)
Formula: k = a + ⌈u·(b − a + 1)⌉ − 1 for u in (0, 1]; a for u = 0.
Inputs and units: `u float64` in [0, 1], `a, b int`.   Output: `(int, error)`.
Errors: as 038; ErrProbability.
Primary source: as 038 (by inverting 039).
Worked example: HAND-CHECKED: die, u = 0.5 → 3.
Variants: none.   Best approach: O(1), O(1), with the same exact ±1 check as 057.   Notes: none.

### UniformQuantile      (ID reserved: NUM-PROB-059)
Formula: x = a + u·(b − a).
Inputs and units: `u float64` in [0, 1]; `a, b float64`, a < b.   Output: `float64`.
Errors: as 044; ErrProbability.
Primary source: as 044 (by inverting 045).
Worked example: HAND-CHECKED: a = 0, b = 10, u = 0.3 → 3.
Variants: none.   Best approach: `math.FMA(u, b−a, a)`. O(1), O(1).   Notes: none.

### ExponentialQuantile      (ID reserved: NUM-PROB-060)
Formula: x = −ln(1 − u)/λ = −log1p(−u)/λ.
Inputs and units: `u float64` in [0, 1); `rate float64` > 0.   Output: `float64`.
Errors: u = 1 → ErrOutOfRange; as 054.
Primary source: [NIST] §1.3.6.6.7, PPF G(p) = −β ln(1 − p) [CHECKED].
Worked example: HAND-CHECKED: λ = 2, u = 0.5 → ln 2/2 = 0.34657359027997264.
Variants: none.   Best approach: `math.Log1p`. O(1), O(1).   Notes: none.

### NormalQuantile      (ID reserved: NUM-PROB-061)
Formula: x = μ + σ·Φ⁻¹(u), where Φ⁻¹ is computed by Wichura's PPND16: rational minimax approximations of
  degree 7 in three regions (|u − 0.5| ≤ 0.425; then r = √(−ln min(u, 1−u)), with r ≤ 5 or r > 5).
Inputs and units: `u float64` in (0, 1); `mean float64`; `sd float64` > 0.   Output: `float64`.
Errors: u ≤ 0 or u ≥ 1 → ErrOutOfRange (±Inf is not returned); as 047.
Primary source: Wichura, "Algorithm AS 241: The Percentage Points of the Normal Distribution", Journal of
  the Royal Statistical Society, Series C (Applied Statistics), 1988, 37(3):477–484, doi:10.2307/2347330
  [CHECKED]. The coefficients are typed from the paper's printed tables, not from any library.
Worked example: [NIST] §1.3.6.7.1 critical values [CHECKED]: u = 0.975 → 1.960; 0.995 → 2.576;
  0.999 → 3.090; 0.025 → −1.960 (3 decimals). Full precision is HAND-CHECKED as the round trip
  Φ(Φ⁻¹(u)) = u within the accuracy Wichura states for PPND16 (about 1 part in 10¹⁶).
Variants: none.   Best approach: PPND16 (no iteration). O(1), O(1).
Notes: it must NOT copy coefficients from a code listing; they come from the paper's tables and are checked
  by the round-trip test.

### LognormalQuantile      (ID reserved: NUM-PROB-062)
Formula: x = exp(μ + σ·Φ⁻¹(u)).
Inputs and units: `u float64` in (0, 1); `mu float64`; `sigma float64` > 0.   Output: `float64` > 0.
Errors: as 061; an overflowing exp → ErrOverflow.
Primary source: [NIST] §1.3.6.6.9 with 061.
Worked example: HAND-CHECKED: μ = 0, σ = 1, u = 0.5 → 1 (the median, m = e^μ).
Variants: none.   Best approach: calls 061. O(1), O(1).   Notes: none.

### Held: odds formats (SOURCE NOT FOUND, no ID, no formula)
DecimalOddsFromFractional, DecimalOddsFromAmerican, AmericanOddsFromDecimal (FractionalOddsFromDecimal is LATER).
The sports project needs them, but every source found for the fractional and American (moneyline)
conventions was a betting site or a converter page; none is a primary source under
`.claude/rules/sources.md`. No formula is written here until Devin names an acceptable source, for example a
gaming regulator's rule book or a peer-reviewed paper that defines moneyline odds (Q7).

## Step 4: build order and dependencies

Build order, lowest first: 001–022 (the rules) → 023–043 (discrete) → 044–056 (continuous) → 057–062 (quantiles).
Inside the discrete block, the unexported Loader helpers (stirlerr, bd0) come first; 026, 032, 035 and 041 call them.

| What | From where | Status |
|---|---|---|
| `math.Erfc`, `math.Log1p`, `math.Expm1`, `math.FMA`, `math.Lgamma` | standard library | available, not BLOCKED |
| Neumaier compensated summation | unexported here, same algorithm and source as S1's RollingSum (Higham §4.3; Neumaier 1974) | S1 keeps its own unexported; nothing exported to call, so nothing BLOCKED |
| Log binomial coefficient, gamma, Stirling | `num/combinatorics` | NOT needed: the Loader route avoids them. If Devin prefers the lgamma route for 026/035, those rows become BLOCKED ON num/combinatorics (request in PENDING_SHARED.md §4) |
| Regularized incomplete beta and gamma (O(1) binomial, Poisson and negative binomial CDFs; later t, χ², F) | no home yet (Q6) | the summation versions here do not need them; the faster variants wait on Q6 |
| Any S1 algorithm | — | none used |

Offers to other subjects (they call `num/prob`, never the reverse):
- **life/games**:
  - poker and Magic hands → HypergeometricPMF/CDF (035, 036)
  - dice → DiscreteUniformPMF (038); a dice-total distribution stays theirs
  - Powerball → ProbabilityFromCounts (001), with the counts from their combinatorics calls
  - bet expectation → ExpectedValue (018)
- **life/sports**:
  - odds → 012–017 (and the held odds formats, once Q7 is answered)
  - expectation → 018
  - normal-model win probabilities → 048–050
- **med/*** (diagnostic tests): PositivePredictiveValue (010), NegativePredictiveValue (011), BayesTheorem (008).
  These state numbers only.
- **num/stats**: NormalCDF/Survival/Quantile (048, 049, 061) for its z-score and confidence work.

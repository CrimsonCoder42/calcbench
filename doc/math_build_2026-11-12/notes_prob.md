# Hand-off: num/prob (NUM-PROB-001…062), 2026-11-12

Builder: prob. Owned paths only: `num/prob/**`, `doc/specs/s9_prob.md`, this file.
Checks run while building: `gofmt -l num/prob` (clean), `go build ./...` (ok), `go vet ./num/prob/` (ok).
Tests and benchmarks NOT run (lead verifies). A scratch module outside the repo printed 14 key values
(binomial, Poisson, hypergeometric, normal quantile and tails), and each matched its reference.
62 functions built, 0 held back. 63 source files (62 concepts + doc.go), 62 test files.

## (a) INDEX.md rows (final form, generated from the doc comments and doc.go)

> Superseded by fix_prob.md (2026-11-12): the Source cells of 002, 003, 005–011, 014–017 and 019 changed (Štrumbelj and Altman–Bland dropped; Casella–Berger at section level); use the lines there.

| ID | Function | Package | Does (one line) | Time | Space | Source | Tags |
|---|---|---|---|---|---|---|---|
| NUM-PROB-001 | ProbabilityFromCounts | num/prob | favorable/total for equally likely outcomes | O(1) | O(1), no allocation | Ross, A First Course in Probability, 10th ed., Pearson, 2019, §2.5 "Sample Spaces Having Equally Likely Outcomes" | probability, classical, counting, equally-likely |
| NUM-PROB-002 | ComplementProbability | num/prob | 1 − P(A) | O(1) | O(1), no allocation | Ross, A First Course in Probability, 10th ed., Pearson, 2019, §2.4 "Some Simple Propositions"; Casella, Berger, Statistical Inference, 2nd ed., Duxbury, 2002, Theorem 1.2.8(c) | probability, complement, rules |
| NUM-PROB-003 | IndependentJointProbability | num/prob | P(A)·P(B) for independent events | O(1) | O(1), no allocation | Ross, A First Course in Probability, 10th ed., Pearson, 2019, §3.4 "Independent Events"; Casella, Berger, Statistical Inference, 2nd ed., Duxbury, 2002, Definition 1.3.7 | probability, joint, independence, intersection, rules |
| NUM-PROB-004 | JointProbability | num/prob | P(A given B)·P(B), the multiplication rule | O(1) | O(1), no allocation | Ross, A First Course in Probability, 10th ed., Pearson, 2019, §3.2 "Conditional Probabilities" (the multiplication rule) | probability, joint, multiplication-rule, conditional, rules |
| NUM-PROB-005 | UnionProbability | num/prob | P(A) + P(B) − P(A ∩ B) | O(1) | O(1), no allocation | Ross, A First Course in Probability, 10th ed., Pearson, 2019, §2.4 "Some Simple Propositions"; Casella, Berger, Statistical Inference, 2nd ed., Duxbury, 2002, Theorem 1.2.9(b) | probability, union, inclusion-exclusion, addition-rule, rules |
| NUM-PROB-006 | ConditionalProbability | num/prob | P(A ∩ B)/P(B) | O(1) | O(1), no allocation | Ross, A First Course in Probability, 10th ed., Pearson, 2019, §3.2 "Conditional Probabilities"; Casella, Berger, Statistical Inference, 2nd ed., Duxbury, 2002, Definition 1.3.2 | probability, conditional, rules |
| NUM-PROB-007 | TotalProbability | num/prob | Σ P(A given Bᵢ)·P(Bᵢ) over a partition | O(n)   n = len(pB) | O(1), no allocation | Ross, A First Course in Probability, 10th ed., Pearson, 2019, §3.3 "Bayes's Formula" (total probability over a partition); Casella, Berger, Statistical Inference, 2nd ed., Duxbury, 2002, Theorem 1.2.11(a); Neumaier, "Rundungsfehleranalyse einiger Verfahren zur Summation endlicher Summen", ZAMM, 1974, 54(1):39–51, doi:10.1002/zamm.19740540106 | probability, total-probability, partition, conditional, compensated-summation, rules |
| NUM-PROB-008 | BayesTheorem | num/prob | P(B given A)·P(A)/P(B) | O(1) | O(1), no allocation | Ross, A First Course in Probability, 10th ed., Pearson, 2019, §3.3 "Bayes's Formula"; Casella, Berger, Statistical Inference, 2nd ed., Duxbury, 2002, Theorem 1.3.5 (Bayes' Rule) | probability, bayes, conditional, posterior, rules |
| NUM-PROB-009 | BayesPosterior | num/prob | posteriors over a partition, appended to dst | O(n)   n = len(priors) | O(1) extra; O(n) for the result, one allocation only if dst lacks capacity | Ross, A First Course in Probability, 10th ed., Pearson, 2019, §3.3 "Bayes's Formula" (over a partition); Casella, Berger, Statistical Inference, 2nd ed., Duxbury, 2002, Theorem 1.3.5 (Bayes' Rule) | probability, bayes, posterior, partition, compensated-summation, rules |
| NUM-PROB-010 | PositivePredictiveValue | num/prob | P(condition given a positive test) from Se, Sp, prevalence | O(1) | O(1), no allocation | Altman, Bland, "Statistics Notes: Diagnostic tests 2: predictive values", BMJ, 1994, 309(6947):102, doi:10.1136/bmj.309.6947.102; Ross, A First Course in Probability, 10th ed., Pearson, 2019, §3.3 "Bayes's Formula" | probability, bayes, diagnostic-test, predictive-value, ppv, medical |
| NUM-PROB-011 | NegativePredictiveValue | num/prob | P(no condition given a negative test) from Se, Sp, prevalence | O(1) | O(1), no allocation | Altman, Bland, "Statistics Notes: Diagnostic tests 2: predictive values", BMJ, 1994, 309(6947):102, doi:10.1136/bmj.309.6947.102; Ross, A First Course in Probability, 10th ed., Pearson, 2019, §3.3 "Bayes's Formula" | probability, bayes, diagnostic-test, predictive-value, npv, medical |
| NUM-PROB-012 | OddsFromProbability | num/prob | odds in favour p/(1 − p) | O(1) | O(1), no allocation | Bland, Altman, "Statistics Notes: The odds ratio", BMJ, 2000, 320(7247):1468, doi:10.1136/bmj.320.7247.1468 (odds = P(event)/P(no event)) | probability, odds, conversion, betting |
| NUM-PROB-013 | ProbabilityFromOdds | num/prob | probability o/(1 + o) from odds in favour | O(1) | O(1), no allocation | Bland, Altman, "Statistics Notes: The odds ratio", BMJ, 2000, 320(7247):1468, doi:10.1136/bmj.320.7247.1468 (odds = P(event)/P(no event)) | probability, odds, conversion, betting |
| NUM-PROB-014 | DecimalOddsFromProbability | num/prob | fair decimal odds 1/p | O(1) | O(1), no allocation | Štrumbelj, "On determining probability forecasts from betting odds", International Journal of Forecasting, 2014, 30(4):934–943, doi:10.1016/j.ijforecast.2014.02.008 | probability, odds, decimal-odds, betting, conversion |
| NUM-PROB-015 | ImpliedProbability | num/prob | probability 1/d implied by decimal odds | O(1) | O(1), no allocation | Štrumbelj, "On determining probability forecasts from betting odds", International Journal of Forecasting, 2014, 30(4):934–943, doi:10.1016/j.ijforecast.2014.02.008 | probability, odds, decimal-odds, implied-probability, betting, conversion |
| NUM-PROB-016 | Overround | num/prob | booksum Σ 1/dᵢ minus 1 of one market | O(n)   n = len(decimalOdds) | O(1), no allocation | Štrumbelj, "On determining probability forecasts from betting odds", International Journal of Forecasting, 2014, 30(4):934–943, doi:10.1016/j.ijforecast.2014.02.008 | probability, odds, overround, booksum, margin, vig, betting |
| NUM-PROB-017 | NormalizedImpliedProbabilities | num/prob | implied probabilities with the overround removed proportionally | O(n)   n = len(decimalOdds) | O(1) extra; O(n) for the result, one allocation only if dst lacks capacity | Štrumbelj, "On determining probability forecasts from betting odds", International Journal of Forecasting, 2014, 30(4):934–943, doi:10.1016/j.ijforecast.2014.02.008 | probability, odds, implied-probability, overround, normalization, betting |
| NUM-PROB-018 | ExpectedValue | num/prob | E[X] of a discrete distribution | O(n)   n = len(outcomes) | O(1), no allocation | Ross, A First Course in Probability, 10th ed., Pearson, 2019, §4.3 (expected value); Neumaier, "Rundungsfehleranalyse einiger Verfahren zur Summation endlicher Summen", ZAMM, 1974, 54(1):39–51, doi:10.1002/zamm.19740540106 | probability, expectation, expected-value, mean, discrete, compensated-summation |
| NUM-PROB-019 | DiscreteVariance | num/prob | Var X of a discrete distribution, two-pass | O(n)   n = len(outcomes) | O(1), no allocation | Ross, A First Course in Probability, 10th ed., Pearson, 2019, §4.5 (variance); Higham, Accuracy and Stability of Numerical Algorithms, 2nd ed., SIAM, 2002, §1.9 "Computing the Sample Variance" | probability, variance, discrete, dispersion, compensated-summation |
| NUM-PROB-020 | ChebyshevBound | num/prob | Chebyshev's bound 1/k² on a deviation of k or more standard deviations | O(1) | O(1), no allocation | Ross, A First Course in Probability, 10th ed., Pearson, 2019, §8.2 "Chebyshev's Inequality and the Weak Law of Large Numbers"; Casella, Berger, Statistical Inference, 2nd ed., Duxbury, 2002, §3.6 "Inequalities and Identities" (P(g(X) ≥ r) ≤ E g(X)/r for g ≥ 0) | probability, chebyshev, inequality, tail-bound |
| NUM-PROB-021 | ChebyshevBoundEpsilon | num/prob | Chebyshev's bound σ²/ε² on a deviation of ε or more | O(1) | O(1), no allocation | Ross, A First Course in Probability, 10th ed., Pearson, 2019, §8.2 "Chebyshev's Inequality and the Weak Law of Large Numbers"; Casella, Berger, Statistical Inference, 2nd ed., Duxbury, 2002, §3.6 "Inequalities and Identities" (P(g(X) ≥ r) ≤ E g(X)/r for g ≥ 0) | probability, chebyshev, inequality, tail-bound |
| NUM-PROB-022 | MarkovBound | num/prob | Markov's bound E[X]/a on P(X ≥ a) | O(1) | O(1), no allocation | Ross, A First Course in Probability, 10th ed., Pearson, 2019, §8.2 "Chebyshev's Inequality and the Weak Law of Large Numbers" (Markov's inequality); Casella, Berger, Statistical Inference, 2nd ed., Duxbury, 2002, §3.6 "Inequalities and Identities" (P(g(X) ≥ r) ≤ E g(X)/r for g ≥ 0) | probability, markov, inequality, tail-bound |
| NUM-PROB-023 | BernoulliPMF | num/prob | Bernoulli P(X = k) | O(1) | O(1), no allocation | Ross, A First Course in Probability, 10th ed., Pearson, 2019, §4.6 (Bernoulli and binomial random variables); Casella, Berger, Statistical Inference, 2nd ed., Duxbury, 2002, §3.2 "Discrete Distributions" (Bernoulli distribution) | probability, distribution, bernoulli, pmf, discrete |
| NUM-PROB-024 | BernoulliCDF | num/prob | Bernoulli P(X ≤ k) | O(1) | O(1), no allocation | Ross, A First Course in Probability, 10th ed., Pearson, 2019, §4.6 (Bernoulli and binomial random variables); Casella, Berger, Statistical Inference, 2nd ed., Duxbury, 2002, §3.2 "Discrete Distributions" (Bernoulli distribution) | probability, distribution, bernoulli, cdf, discrete |
| NUM-PROB-025 | BernoulliMeanVariance | num/prob | Bernoulli mean and variance | O(1) | O(1), no allocation | Ross, A First Course in Probability, 10th ed., Pearson, 2019, §4.6 (Bernoulli and binomial random variables); Casella, Berger, Statistical Inference, 2nd ed., Duxbury, 2002, §3.2 "Discrete Distributions" (Bernoulli distribution) | probability, distribution, bernoulli, mean, variance, expected-value |
| NUM-PROB-026 | BinomialPMF | num/prob | binomial P(X = k), saddle-point form | O(1) | O(1), no allocation | Loader, "Fast and Accurate Computation of Binomial Probabilities", technical report, 2000 (rev. 2002), https://www.r-project.org/doc/reports/CLoader-dbinom-2002.pdf (read 2026-11-12), §2 eqs. (3)–(7) and Appendix A (stirlerr, bd0); NIST/SEMATECH, e-Handbook of Statistical Methods, §1.3.6.6.18 "Binomial Distribution", https://www.itl.nist.gov/div898/handbook/eda/section3/eda366i.htm (read 2026-11-12); Ross, A First Course in Probability, 10th ed., Pearson, 2019, §4.6 (Bernoulli and binomial random variables) | probability, distribution, binomial, pmf, discrete, saddle-point |
| NUM-PROB-027 | BinomialCDF | num/prob | binomial P(X ≤ k), compensated tail sum | O(min(k + 1, n − k, c·(σ + 1)))   σ = √(np(1 − p)); c ≈ 10 in double precision (the terms beyond c·σ from the mean are below 2⁻⁵³ by Bernstein's inequality) | O(1), no allocation | NIST/SEMATECH, e-Handbook of Statistical Methods, §1.3.6.6.18 "Binomial Distribution", https://www.itl.nist.gov/div898/handbook/eda/section3/eda366i.htm (read 2026-11-12); Loader, "Fast and Accurate Computation of Binomial Probabilities", technical report, 2000 (rev. 2002), https://www.r-project.org/doc/reports/CLoader-dbinom-2002.pdf (read 2026-11-12), §2 eqs. (3)–(7) and Appendix A (stirlerr, bd0) | probability, distribution, binomial, cdf, discrete, compensated-summation |
| NUM-PROB-028 | BinomialMeanVariance | num/prob | binomial mean and variance | O(1) | O(1), no allocation | NIST/SEMATECH, e-Handbook of Statistical Methods, §1.3.6.6.18 "Binomial Distribution", https://www.itl.nist.gov/div898/handbook/eda/section3/eda366i.htm (read 2026-11-12); Ross, A First Course in Probability, 10th ed., Pearson, 2019, §4.6 (Bernoulli and binomial random variables) | probability, distribution, binomial, mean, variance, expected-value |
| NUM-PROB-029 | GeometricTrialsPMF | num/prob | geometric P(X = k), X = trials to the first success | O(1) | O(1), no allocation | Ross, A First Course in Probability, 10th ed., Pearson, 2019, §4.8 (other discrete distributions: the geometric random variable); Casella, Berger, Statistical Inference, 2nd ed., Duxbury, 2002, §3.2 "Discrete Distributions" (geometric distribution) | probability, distribution, geometric, pmf, discrete, trials, waiting-time |
| NUM-PROB-030 | GeometricTrialsCDF | num/prob | geometric P(X ≤ k), X = trials to the first success | O(1) | O(1), no allocation | Ross, A First Course in Probability, 10th ed., Pearson, 2019, §4.8 (other discrete distributions: the geometric random variable); Casella, Berger, Statistical Inference, 2nd ed., Duxbury, 2002, §3.2 "Discrete Distributions" (geometric distribution) | probability, distribution, geometric, cdf, discrete, trials, waiting-time |
| NUM-PROB-031 | GeometricTrialsMeanVariance | num/prob | geometric mean and variance, trials convention | O(1) | O(1), no allocation | Ross, A First Course in Probability, 10th ed., Pearson, 2019, §4.8 (other discrete distributions: the geometric random variable); Casella, Berger, Statistical Inference, 2nd ed., Duxbury, 2002, §3.2 "Discrete Distributions" (geometric distribution) | probability, distribution, geometric, mean, variance, expected-value, trials |
| NUM-PROB-032 | NegativeBinomialTrialsPMF | num/prob | negative binomial P(X = x), X = trials to the r-th success | O(1) | O(1), no allocation | Ross, A First Course in Probability, 10th ed., Pearson, 2019, §4.8 (other discrete distributions: the negative binomial random variable); Casella, Berger, Statistical Inference, 2nd ed., Duxbury, 2002, §3.2 "Discrete Distributions" (negative binomial distribution); Loader, "Fast and Accurate Computation of Binomial Probabilities", technical report, 2000 (rev. 2002), https://www.r-project.org/doc/reports/CLoader-dbinom-2002.pdf (read 2026-11-12), §2 eqs. (3)–(7) and Appendix A (stirlerr, bd0) | probability, distribution, negative-binomial, pascal, pmf, discrete, trials, waiting-time |
| NUM-PROB-033 | NegativeBinomialTrialsCDF | num/prob | negative binomial P(X ≤ x), X = trials to the r-th success | O(min(r, x − r + 1, c·(σ + 1)))   σ = √(xp(1 − p)); c ≈ 10 (as BinomialCDF on x trials) | O(1), no allocation | Ross, A First Course in Probability, 10th ed., Pearson, 2019, §4.8 (other discrete distributions: the negative binomial random variable); Casella, Berger, Statistical Inference, 2nd ed., Duxbury, 2002, §3.2 "Discrete Distributions" (negative binomial distribution); Loader, "Fast and Accurate Computation of Binomial Probabilities", technical report, 2000 (rev. 2002), https://www.r-project.org/doc/reports/CLoader-dbinom-2002.pdf (read 2026-11-12), §2 eqs. (3)–(7) and Appendix A (stirlerr, bd0) | probability, distribution, negative-binomial, pascal, cdf, discrete, trials |
| NUM-PROB-034 | NegativeBinomialTrialsMeanVariance | num/prob | negative binomial mean and variance, trials convention | O(1) | O(1), no allocation | Ross, A First Course in Probability, 10th ed., Pearson, 2019, §4.8 (other discrete distributions: the negative binomial random variable); Casella, Berger, Statistical Inference, 2nd ed., Duxbury, 2002, §3.2 "Discrete Distributions" (negative binomial distribution) | probability, distribution, negative-binomial, pascal, mean, variance, expected-value, trials |
| NUM-PROB-035 | HypergeometricPMF | num/prob | hypergeometric P(X = k), saddle-point form | O(1) | O(1), no allocation | Ross, A First Course in Probability, 10th ed., Pearson, 2019, §4.8 (other discrete distributions: the hypergeometric random variable); Casella, Berger, Statistical Inference, 2nd ed., Duxbury, 2002, §3.2 "Discrete Distributions" (hypergeometric distribution); Loader, "Fast and Accurate Computation of Binomial Probabilities", technical report, 2000 (rev. 2002), https://www.r-project.org/doc/reports/CLoader-dbinom-2002.pdf (read 2026-11-12), §2 eqs. (3)–(7) and Appendix A (stirlerr, bd0) | probability, distribution, hypergeometric, pmf, discrete, sampling-without-replacement, cards |
| NUM-PROB-036 | HypergeometricCDF | num/prob | hypergeometric P(X ≤ k), compensated tail sum | O(min(k − lo + 1, hi − k, c·(σ + 1)))   lo, hi = the support bounds, σ = the standard deviation (HypergeometricMeanVariance), c ≈ 10 | O(1), no allocation | Ross, A First Course in Probability, 10th ed., Pearson, 2019, §4.8 (other discrete distributions: the hypergeometric random variable); Casella, Berger, Statistical Inference, 2nd ed., Duxbury, 2002, §3.2 "Discrete Distributions" (hypergeometric distribution) | probability, distribution, hypergeometric, cdf, discrete, sampling-without-replacement, cards |
| NUM-PROB-037 | HypergeometricMeanVariance | num/prob | hypergeometric mean and variance | O(1) | O(1), no allocation | Ross, A First Course in Probability, 10th ed., Pearson, 2019, §4.8 (other discrete distributions: the hypergeometric random variable); Casella, Berger, Statistical Inference, 2nd ed., Duxbury, 2002, §3.2 "Discrete Distributions" (hypergeometric distribution) | probability, distribution, hypergeometric, mean, variance, expected-value |
| NUM-PROB-038 | DiscreteUniformPMF | num/prob | discrete uniform P(X = k) on a…b | O(1) | O(1), no allocation | Casella, Berger, Statistical Inference, 2nd ed., Duxbury, 2002, §3.2 "Discrete Distributions" (discrete uniform distribution on 1…N, shifted here to a…b) | probability, distribution, discrete-uniform, pmf, discrete, dice |
| NUM-PROB-039 | DiscreteUniformCDF | num/prob | discrete uniform P(X ≤ k) on a…b | O(1) | O(1), no allocation | Casella, Berger, Statistical Inference, 2nd ed., Duxbury, 2002, §3.2 "Discrete Distributions" (discrete uniform distribution on 1…N, shifted here to a…b) | probability, distribution, discrete-uniform, cdf, discrete, dice |
| NUM-PROB-040 | DiscreteUniformMeanVariance | num/prob | discrete uniform mean and variance on a…b | O(1) | O(1), no allocation | Casella, Berger, Statistical Inference, 2nd ed., Duxbury, 2002, §3.2 "Discrete Distributions" (discrete uniform distribution on 1…N, shifted here to a…b) | probability, distribution, discrete-uniform, mean, variance, expected-value, dice |
| NUM-PROB-041 | PoissonPMF | num/prob | Poisson P(X = k), saddle-point form | O(1) | O(1), no allocation | NIST/SEMATECH, e-Handbook of Statistical Methods, §1.3.6.6.19 "Poisson Distribution", https://www.itl.nist.gov/div898/handbook/eda/section3/eda366j.htm (read 2026-11-12); Loader, "Fast and Accurate Computation of Binomial Probabilities", technical report, 2000 (rev. 2002), https://www.r-project.org/doc/reports/CLoader-dbinom-2002.pdf (read 2026-11-12), §2 eqs. (3)–(7) and Appendix A (stirlerr, bd0); Ross, A First Course in Probability, 10th ed., Pearson, 2019, §4.7 (Poisson random variable) | probability, distribution, poisson, pmf, discrete, counts, saddle-point |
| NUM-PROB-042 | PoissonCDF | num/prob | Poisson P(X ≤ k), compensated tail sum | O(min(k + 1, c·(√λ + 1)))   c ≈ 10 in double precision (terms beyond c·√λ from λ are below 2⁻⁵³) | O(1), no allocation | NIST/SEMATECH, e-Handbook of Statistical Methods, §1.3.6.6.19 "Poisson Distribution", https://www.itl.nist.gov/div898/handbook/eda/section3/eda366j.htm (read 2026-11-12); Loader, "Fast and Accurate Computation of Binomial Probabilities", technical report, 2000 (rev. 2002), https://www.r-project.org/doc/reports/CLoader-dbinom-2002.pdf (read 2026-11-12), §2 eqs. (3)–(7) and Appendix A (stirlerr, bd0) | probability, distribution, poisson, cdf, discrete, counts, compensated-summation |
| NUM-PROB-043 | PoissonMeanVariance | num/prob | Poisson mean and variance | O(1) | O(1), no allocation | NIST/SEMATECH, e-Handbook of Statistical Methods, §1.3.6.6.19 "Poisson Distribution", https://www.itl.nist.gov/div898/handbook/eda/section3/eda366j.htm (read 2026-11-12); Ross, A First Course in Probability, 10th ed., Pearson, 2019, §4.7 (Poisson random variable) | probability, distribution, poisson, mean, variance, expected-value |
| NUM-PROB-044 | UniformPDF | num/prob | continuous uniform density on [a, b] | O(1) | O(1), no allocation | NIST/SEMATECH, e-Handbook of Statistical Methods, §1.3.6.6.2 "Uniform Distribution", https://www.itl.nist.gov/div898/handbook/eda/section3/eda3662.htm (read 2026-11-12); Ross, A First Course in Probability, 10th ed., Pearson, 2019, §5.3 (uniform random variable) | probability, distribution, uniform, pdf, continuous |
| NUM-PROB-045 | UniformCDF | num/prob | continuous uniform P(X ≤ x) on [a, b] | O(1) | O(1), no allocation | NIST/SEMATECH, e-Handbook of Statistical Methods, §1.3.6.6.2 "Uniform Distribution", https://www.itl.nist.gov/div898/handbook/eda/section3/eda3662.htm (read 2026-11-12); Ross, A First Course in Probability, 10th ed., Pearson, 2019, §5.3 (uniform random variable) | probability, distribution, uniform, cdf, continuous |
| NUM-PROB-046 | UniformMeanVariance | num/prob | continuous uniform mean and variance on [a, b] | O(1) | O(1), no allocation | NIST/SEMATECH, e-Handbook of Statistical Methods, §1.3.6.6.2 "Uniform Distribution", https://www.itl.nist.gov/div898/handbook/eda/section3/eda3662.htm (read 2026-11-12); Ross, A First Course in Probability, 10th ed., Pearson, 2019, §5.3 (uniform random variable) | probability, distribution, uniform, mean, variance, expected-value, continuous |
| NUM-PROB-047 | NormalPDF | num/prob | normal density | O(1) | O(1), no allocation | NIST/SEMATECH, e-Handbook of Statistical Methods, §1.3.6.6.1 "Normal Distribution", https://www.itl.nist.gov/div898/handbook/eda/section3/eda3661.htm (read 2026-11-12); Ross, A First Course in Probability, 10th ed., Pearson, 2019, §5.4 (normal random variables) | probability, distribution, normal, gaussian, pdf, continuous, bell-curve |
| NUM-PROB-048 | NormalCDF | num/prob | normal P(X ≤ x), accurate in the lower tail | O(1) | O(1), no allocation | NIST/SEMATECH, e-Handbook of Statistical Methods, §1.3.6.6.1 "Normal Distribution", https://www.itl.nist.gov/div898/handbook/eda/section3/eda3661.htm (read 2026-11-12); NIST DLMF, https://dlmf.nist.gov (read 2026-11-12), §7.2 eqs. 7.2.1–7.2.2 (erf, erfc) and §7.1 (Φ(z) = ½erfc(−z/√2), Q(z) = ½erfc(z/√2)); Ross, A First Course in Probability, 10th ed., Pearson, 2019, §5.4 (normal random variables) | probability, distribution, normal, gaussian, cdf, continuous, z-score, phi |
| NUM-PROB-049 | NormalSurvival | num/prob | normal P(X > x), accurate in the upper tail | O(1) | O(1), no allocation | NIST/SEMATECH, e-Handbook of Statistical Methods, §1.3.6.6.1 "Normal Distribution", https://www.itl.nist.gov/div898/handbook/eda/section3/eda3661.htm (read 2026-11-12); NIST DLMF, https://dlmf.nist.gov (read 2026-11-12), §7.2 eqs. 7.2.1–7.2.2 (erf, erfc) and §7.1 (Φ(z) = ½erfc(−z/√2), Q(z) = ½erfc(z/√2)) | probability, distribution, normal, gaussian, survival, upper-tail, continuous, p-value |
| NUM-PROB-050 | NormalIntervalProbability | num/prob | normal P(lo ≤ X ≤ hi), tail-aware | O(1) | O(1), no allocation | NIST/SEMATECH, e-Handbook of Statistical Methods, §1.3.6.6.1 "Normal Distribution", https://www.itl.nist.gov/div898/handbook/eda/section3/eda3661.htm (read 2026-11-12); NIST DLMF, https://dlmf.nist.gov (read 2026-11-12), §7.2 eqs. 7.2.1–7.2.2 (erf, erfc) and §7.1 (Φ(z) = ½erfc(−z/√2), Q(z) = ½erfc(z/√2)) | probability, distribution, normal, gaussian, interval, continuous, z-score |
| NUM-PROB-051 | LognormalPDF | num/prob | lognormal density | O(1) | O(1), no allocation | NIST/SEMATECH, e-Handbook of Statistical Methods, §1.3.6.6.9 "Lognormal Distribution", https://www.itl.nist.gov/div898/handbook/eda/section3/eda3669.htm (read 2026-11-12) | probability, distribution, lognormal, pdf, continuous |
| NUM-PROB-052 | LognormalCDF | num/prob | lognormal P(X ≤ x) | O(1) | O(1), no allocation | NIST/SEMATECH, e-Handbook of Statistical Methods, §1.3.6.6.9 "Lognormal Distribution", https://www.itl.nist.gov/div898/handbook/eda/section3/eda3669.htm (read 2026-11-12) | probability, distribution, lognormal, cdf, continuous |
| NUM-PROB-053 | LognormalMeanVariance | num/prob | lognormal mean and variance | O(1) | O(1), no allocation | NIST/SEMATECH, e-Handbook of Statistical Methods, §1.3.6.6.9 "Lognormal Distribution", https://www.itl.nist.gov/div898/handbook/eda/section3/eda3669.htm (read 2026-11-12) (the m = 1 form, scaled by e^μ) | probability, distribution, lognormal, mean, variance, expected-value |
| NUM-PROB-054 | ExponentialPDF | num/prob | exponential density, rate form | O(1) | O(1), no allocation | NIST/SEMATECH, e-Handbook of Statistical Methods, §1.3.6.6.7 "Exponential Distribution", https://www.itl.nist.gov/div898/handbook/eda/section3/eda3667.htm (read 2026-11-12); Ross, A First Course in Probability, 10th ed., Pearson, 2019, §5.5 (exponential random variables) | probability, distribution, exponential, pdf, continuous, waiting-time, decay |
| NUM-PROB-055 | ExponentialCDF | num/prob | exponential P(X ≤ x), rate form | O(1) | O(1), no allocation | NIST/SEMATECH, e-Handbook of Statistical Methods, §1.3.6.6.7 "Exponential Distribution", https://www.itl.nist.gov/div898/handbook/eda/section3/eda3667.htm (read 2026-11-12); Ross, A First Course in Probability, 10th ed., Pearson, 2019, §5.5 (exponential random variables) | probability, distribution, exponential, cdf, continuous, waiting-time, decay |
| NUM-PROB-056 | ExponentialMeanVariance | num/prob | exponential mean and variance, rate form | O(1) | O(1), no allocation | NIST/SEMATECH, e-Handbook of Statistical Methods, §1.3.6.6.7 "Exponential Distribution", https://www.itl.nist.gov/div898/handbook/eda/section3/eda3667.htm (read 2026-11-12); Ross, A First Course in Probability, 10th ed., Pearson, 2019, §5.5 (exponential random variables) | probability, distribution, exponential, mean, variance, expected-value |
| NUM-PROB-057 | GeometricTrialsQuantile | num/prob | least k with P(X ≤ k) ≥ u, trials convention | O(1 + k·2⁻⁵⁰)   k = the result (the float64 estimate is off by at most a few units in the last place of k, each a correction step) | O(1), no allocation | Ross, A First Course in Probability, 10th ed., Pearson, 2019, §4.8 (other discrete distributions: the geometric random variable), inverted; Casella, Berger, Statistical Inference, 2nd ed., Duxbury, 2002, §3.2 "Discrete Distributions" (geometric distribution), inverted | probability, distribution, geometric, quantile, inverse-cdf, trials |
| NUM-PROB-058 | DiscreteUniformQuantile | num/prob | least k in a…b with P(X ≤ k) ≥ u | O(1 + m·2⁻⁵²)   m = b − a + 1 (the float64 estimate is off by at most a few units in the last place of m, each a correction step) | O(1), no allocation | Casella, Berger, Statistical Inference, 2nd ed., Duxbury, 2002, §3.2 "Discrete Distributions" (discrete uniform distribution on 1…N, shifted here to a…b), inverted | probability, distribution, discrete-uniform, quantile, inverse-cdf, dice |
| NUM-PROB-059 | UniformQuantile | num/prob | continuous uniform quantile a + u(b − a) | O(1) | O(1), no allocation | NIST/SEMATECH, e-Handbook of Statistical Methods, §1.3.6.6.2 "Uniform Distribution", https://www.itl.nist.gov/div898/handbook/eda/section3/eda3662.htm (read 2026-11-12); Ross, A First Course in Probability, 10th ed., Pearson, 2019, §5.3 (uniform random variable) | probability, distribution, uniform, quantile, inverse-cdf, ppf |
| NUM-PROB-060 | ExponentialQuantile | num/prob | exponential quantile −ln(1 − u)/λ | O(1) | O(1), no allocation | NIST/SEMATECH, e-Handbook of Statistical Methods, §1.3.6.6.7 "Exponential Distribution", https://www.itl.nist.gov/div898/handbook/eda/section3/eda3667.htm (read 2026-11-12) (PPF G(p) = −β ln(1 − p), β = 1/λ) | probability, distribution, exponential, quantile, inverse-cdf, ppf |
| NUM-PROB-061 | NormalQuantile | num/prob | normal quantile by Wichura's PPND16 | O(1) | O(1), no allocation | Wichura, "Algorithm AS 241: The Percentage Points of the Normal Distribution", Journal of the Royal Statistical Society, Series C (Applied Statistics), 1988, 37(3):477–484, doi:10.2307/2347330 (PPND16); NIST/SEMATECH, e-Handbook of Statistical Methods, §1.3.6.7.1 "Cumulative Distribution Function of the Standard Normal Distribution", https://www.itl.nist.gov/div898/handbook/eda/section3/eda3671.htm (read 2026-11-12) | probability, distribution, normal, gaussian, quantile, inverse-cdf, ppf, critical-value, z-score |
| NUM-PROB-062 | LognormalQuantile | num/prob | lognormal quantile exp(μ + σΦ⁻¹(u)) | O(1) | O(1), no allocation | NIST/SEMATECH, e-Handbook of Statistical Methods, §1.3.6.6.9 "Lognormal Distribution", https://www.itl.nist.gov/div898/handbook/eda/section3/eda3669.htm (read 2026-11-12) (PPF exp(σΦ⁻¹(p)), scaled by e^μ); Wichura, "Algorithm AS 241: The Percentage Points of the Normal Distribution", Journal of the Royal Statistical Society, Series C (Applied Statistics), 1988, 37(3):477–484, doi:10.2307/2347330 (PPND16) | probability, distribution, lognormal, quantile, inverse-cdf, ppf |

## (b) doc/sources/ records (for the lead to fold in)

### `ross_first_course_probability.md` (new)
- Citation: Ross, *A First Course in Probability*, 10th ed., Pearson, 2019.
- Taken (section level only): §2.4 "Some Simple Propositions" → 002, 005; §2.5 "Sample Spaces Having Equally
  Likely Outcomes" → 001; §3.2 "Conditional Probabilities" → 004, 006; §3.3 "Bayes's Formula" → 007, 008, 009,
  010, 011; §3.4 "Independent Events" → 003; §4.3 (expected value) → 018; §4.5 (variance) → 019; §4.6 (Bernoulli
  and binomial) → 023–026, 028; §4.7 (Poisson) → 041, 043; §4.8 (geometric, negative binomial, hypergeometric)
  → 029–037, 057; §5.3 (uniform) → 044–046, 059; §5.4 (normal) → 047, 048; §5.5 (exponential) → 054–056;
  §8.2 "Chebyshev's Inequality and the Weak Law of Large Numbers" → 020, 021, 022.
- Checked: section titles against Pearson's 10th-edition product page TOC
  (pearson.com/en-us/subject-catalog/p/first-course-in-probability-a/P200000006334), read 2026-11-12. §4.8.1–4.8.3
  subsections are confirmed only by a search summary, so 4.8 is cited. Proposition, equation and example
  numbers (Prop 4.1, 4.3; eqs. 2.1, 2.2, 3.1, 3.4, 4.1; Examples 3a, 5a, 2a; Props 8.2.1, 8.2.2) were NOT
  confirmed and are not cited. Worked values are hand-checked in the tests.

### `casella_berger_statistical_inference.md` (new)

> Superseded by fix_prob.md (2026-11-12): theorem and definition numbers dropped; NUM-PROB-002…009 now cite §1.2/§1.3 (record is `doc/sources/casella_berger_2002.md`).
- Citation: Casella, Berger, *Statistical Inference*, 2nd ed., Duxbury, 2002.
- Taken: Theorem 1.2.8(c) → 002; Theorem 1.2.9(b) → 005; Theorem 1.2.11(a) → 007; Definition 1.3.2 → 006;
  Theorem 1.3.5 (Bayes' Rule) → 008, 009; Definition 1.3.7 → 003; §3.2 "Discrete Distributions" (discrete
  uniform 1…N, Bernoulli, geometric in the trials form p(1−p)^(x−1), negative binomial in the trials form
  C(x−1,r−1)pʳ(1−p)^(x−r), hypergeometric) → 023–025, 029–040, 057, 058; §3.6 "Inequalities and Identities"
  (P(g(X) ≥ r) ≤ E g(X)/r) → 020, 021, 022.
- Checked: 1.2.8, 1.2.9, 1.2.11 quoted on hlt.inesc-id.pt (Statistical_Inference/Probability_Theory);
  1.3.2, 1.3.5, 1.3.7 via USU STAT 6710 notes (usu.edu/math/symanzik) and a second set of slides, a search
  summary only; §3.2 and §3.6 titles and pages (85, 121) from a library TOC (typeset.io, ProofWiki book page,
  HTW Berlin TOC PDF); §3.2 contents via Purdue STAT 610 lecture 3, which follows CB closely. Theorem 3.6.1's
  number was NOT confirmed and is not cited.

### `nist_ehandbook.md` (add a "§1.3.6 distributions" section)
- Taken: §1.3.6.6.1 Normal (eda3661.htm: PDF; CDF and PPF "computed numerically") → 047–050, 061;
  §1.3.6.6.2 Uniform (eda3662.htm: general PDF 1/(B−A), mean (A+B)/2, SD √((B−A)²/12); CDF and PPF in the
  standard form only) → 044–046, 059; §1.3.6.6.7 Exponential (eda3667.htm: CDF 1 − e^(−x/β), PPF −β ln(1−p),
  mean and SD β, λ = 1/β) → 054–056, 060; §1.3.6.6.9 Lognormal (eda3669.htm: CDF Φ(ln x/σ), PPF exp(σΦ⁻¹(p)),
  mean e^(σ²/2), SD √(e^(σ²)(e^(σ²)−1)) for m = 1) → 051–053, 062; §1.3.6.6.18 Binomial (eda366i.htm: PMF, CDF as
  the sum, mean np, SD √(np(1−p))) → 026–028; §1.3.6.6.19 Poisson (eda366j.htm: PMF, mean λ, SD √λ) → 041–043;
  §1.3.6.7.1 normal table (eda3671.htm: row 1.9, col 0.06 = 0.47500; critical values 0.975 → 1.960,
  0.995 → 2.576, 0.999 → 3.090) → 048, 049, 050, 061.
- Checked: every page opened 2026-11-12. Tests reproduce 0.47500 (to 5 decimals) and the three critical values
  (to 3 decimals). Also found: §1.3.6.6.3 Cauchy (eda3663.htm) gives the PDF, standard CDF 0.5 + arctan(x)/π and
  PPF −cot(πp); Cauchy is sourced but left LATER (D-013).

### `loader_2000.md` (new)
- Citation: Loader, "Fast and Accurate Computation of Binomial Probabilities", technical report, 2000 (revised
  February 5, 2002), https://www.r-project.org/doc/reports/CLoader-dbinom-2002.pdf (read 2026-11-12).
- Taken: §2, eq. (3) the saddle-point identity; eq. (4) Stirling–De Moivre with remainder δ(n); eq. (5)
  p(x; n, x/n); eq. (6) the deviance D₀ and its series (x−np)²/(x+np) + 2x Σ v^(2j+1)/(2j+1), used when
  |x − np| < 0.1(x + np); eq. (7) Poisson; Appendix A: the δ table for n ≤ 15, the series terms 1/1680 and 1/1188
  and the thresholds 35, 80, 500, and the assembled dbinom/dpois. → 026, 027, 032, 033, 035, 041, 042 (and the
  terms inside 036).
- Checked: report read in full. The δ table was recomputed independently with 60-digit decimal arithmetic. Tests:
  exact big.Rat binomial, negative binomial and hypergeometric oracles, plus the n = 1000 and λ = 1000 values.

### `strumbelj_2014.md` (new)

> Superseded by fix_prob.md (2026-11-12): Štrumbelj 2014 is no longer cited; NUM-PROB-014…017 cite Wheatcroft 2021, Egidi–Pauli–Torelli 2018 and Wunderlich–Memmert 2018. Record retired to `_delete/doc/sources/strumbelj_2014.md`.
- Citation: Štrumbelj, "On determining probability forecasts from betting odds", International Journal of
  Forecasting, 2014, 30(4):934–943, doi:10.1016/j.ijforecast.2014.02.008.
- Taken: decimal odds as the inverse of probability, the booksum Σ 1/dᵢ (overround = booksum − 1), and basic
  normalisation πᵢ = (1/dᵢ)/booksum. Also the finding that Shin's method beats basic normalisation (stated in
  017's doc comment). → 014, 015, 016, 017.
- Checked: DOI, title, volume and pages via the Crossref API. The article text was not readable (ScienceDirect
  returned 403), so no section locator is cited. The booksum and basic-normalisation definitions are attributed
  to it by two independent secondary sources: Springer article s10479-022-04722-3 and arXiv:1802.08848.

### `altman_bland_1994.md` (new)

> Superseded by fix_prob.md (2026-11-12): Altman–Bland 1994 is no longer cited; NUM-PROB-010/011 cite Ross §3.3 only. Record retired to `_delete/doc/sources/altman_bland_1994.md`.
- Citation: Altman, Bland, "Statistics Notes: Diagnostic tests 2: predictive values", BMJ, 1994,
  309(6947):102, doi:10.1136/bmj.309.6947.102 (PMC2540558); companion "Diagnostic tests 1", BMJ 1994;308:1552.
- Taken: PPV and NPV from sensitivity, specificity and prevalence, and the liver-scan example (TP 231, FP 32,
  FN 27, TN 54; PPV 231/263 ≈ 0.88). → 010, 011.
- Checked: citation via PMC metadata. The counts come from a secondary reproduction of "Diagnostic tests 1"
  (docs.switzerlandomics.ch). PPV 0.88 comes from a search snippet. The full text is a scanned image and was
  NOT read; the Bayes form is also covered by Ross §3.3. Tests: exact-fraction PPV 231/263 and NPV 2/3.

### `bland_altman_2000.md` (new)
- Citation: Bland, Altman, "Statistics Notes: The odds ratio", BMJ, 2000, 320(7247):1468,
  doi:10.1136/bmj.320.7247.1468 (PMC1127651).
- Taken: "The odds is the ratio of the probability that the event of interest occurs to the probability that
  it does not" → 012, 013.
- Checked: PMC full text read 2026-11-12 (quote above).

### `wichura_1988.md` (new)
- Citation: Wichura, "Algorithm AS 241: The Percentage Points of the Normal Distribution", Journal of the
  Royal Statistical Society, Series C (Applied Statistics), 1988, 37(3):477–484, doi:10.2307/2347330.
- Taken: PPND16 (split 0.425 and 5, constants 0.180625 and 1.6, the coefficients A0–A7, B1–B7, C0–C7, D1–D7,
  E0–E7, F1–F7) → 061, 062.
- Checked: Crossref (vol. 37, no. 3, first page 477; the end page 484 is not in Crossref). The coefficients
  came from the journal's published listing at StatLib (lib.stat.cmu.edu/apstat/241, header "ALGORITHM AS241
  APPL. STATIST. (1988) VOL. 37, NO. 3"). Two fetches agreed digit for digit, and the copy in
  `normal_quantile.go` was compared digit by digit. Tests: NIST critical values and a 2000-point round trip
  through NormalCDF and NormalSurvival in all three regions.

### `nist_dlmf.md` (add to the existing record)
- Taken: §7.1 (notation: Φ(z) = ½erfc(−z/√2), Q(z) = ½erfc(z/√2), no equation number) and §7.2 eqs.
  7.2.1–7.2.2 (erf, erfc) → 048, 049, 050. The Stirling-series coefficients of eq. 5.11.1 (already recorded) are
  used by the δ series inside 026/032/035/041.
- Checked: dlmf.nist.gov/7.2 and /7.1 read 2026-11-12.

### `neumaier_1974.md`, `higham_2002.md` (add IDs)
- Neumaier → 007, 018 (and the compensated sums inside 009, 016, 017, 027, 033, 036, 042).
- Higham §1.9 "Computing the Sample Variance" (confirmed from a library TOC) → 019. Higham §4.3 is cited in the
  unexported compensatedSum comment.

## (c) Citations: confirmed, not confirmed, held back

> Superseded by fix_prob.md (2026-11-12): the three "NOT confirmed, still cited" items (Štrumbelj, Altman–Bland, CB 1.3.x numbers) are no longer cited.

- Confirmed, with how: see (b). Every function has at least one confirmed source, so none is held back.
- Narrowed to what could be confirmed:
  - Ross: every proposition, equation and example number → section level.
  - CB: Theorem 3.6.1 → §3.6.
  - Ross §3.3 "odds" → replaced by Bland & Altman 2000.
  - DLMF "7.2.1–7.2.2 for Φ" → §7.1 + §7.2.
  - Loader: "§§2–3" → §2 eqs. (3)–(7) and Appendix A.
- NOT confirmed, still cited:
  - Štrumbelj: section locators (none cited).
  - Altman–Bland: formula text not read first-hand (Ross §3.3 covers the formula).
  - CB 1.3.2/1.3.5/1.3.7: from a search summary of course notes.
- Held odds formats (decimal↔fractional↔American): no primary source found. Every definition with conversion
  formulas was on a betting or industry site; no regulator rule book or peer-reviewed paper turned up. They
  stay held, with no ID.
- Cauchy: now sourced (NIST §1.3.6.6.3), but left LATER per D-013.

## Proposed for doc/ASSUMPTIONS.md (decisions made without Devin, told to proceed)

- A-PROB-1: The negative binomial is named `NegativeBinomialTrials*`, matching D-013's geometric. R and SciPy count
  failures, so a bare name would mislead. Lost: the catalog's `NegativeBinomial*`.
- A-PROB-2: Package distribution tolerance |Σpᵢ − 1| ≤ n·2⁻⁵² (catalog Q5). Lost: an exact sum (rejects 0.1 +
  0.2 + 0.7-style inputs after a caller's rounding) and a fixed 1e-9 (an invented constant).
- A-PROB-3: CDF tail sums are capped at 2²⁸ terms → ErrOutOfRange (σ above roughly 10⁷), so that no call runs
  unbounded until Q6 gives the incomplete beta/gamma a home. Lost: no cap (a call could run for minutes).
- A-PROB-4: Parameter names `population, successes, draws` (hypergeometric), `threshold` (Markov),
  `a, b` (uniform, as the sources). Lost: N, K, n (differ only by case).

## (d) FOUND, NOT FIXED

- `doc/catalog/prob.md` 041 worked example: 0.012614611348719664 is off by 1.5e-13 relative. The correct value
  is 0.0126146113487215 (50-digit decimal), and the test uses it.
- `doc/catalog/prob.md` 026 formula: the saddle-point exponent is written −stirlerr(n) + stirlerr(k) + …; the
  correct sign is +δ(n) − δ(k) − δ(n−k) (Loader eq. 5). The code is correct.
- `doc/catalog/prob.md` 026: "near machine precision for every n" holds near the mode only. In the tails the
  relative error is about 2⁻⁵³·D₀; the doc comment says so.
- `doc/catalog/prob.md` header points to the work folder `doc/prob_2026-11-12/` and its PENDING_SHARED.md,
  neither of which exists.
- `doc/catalog/prob.md` 044: the NIST uniform page gives the CDF and PPF only in the standard form (noted in the
  spec; Ross §5.3 is the second source).
- PROGRESS.md line for the lead: `| num/prob | BUILT: 62 functions NUM-PROB-001…062, none held; Ross cited at
  section level; odds (012/013) cite Bland & Altman 2000; negative binomial named …Trials (A-PROB-1) |`

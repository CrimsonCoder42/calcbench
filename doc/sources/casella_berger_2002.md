# Casella and Berger, Statistical Inference

**Citation**: Casella, Berger, *Statistical Inference*, 2nd ed., Duxbury, 2002.

| Taken | Locator | Used by |
|---|---|---|
| Variance and standard deviation | §2.3 "Moments and Moment Generating Functions" | NUM-STATS-017, -022 |
| Complement, union, law of total probability | §1.2 "Basics of Probability Theory" | NUM-PROB-002, -005, -007 |
| Conditional probability, independence | §1.3 "Conditional Probability and Independence" | NUM-PROB-003, -006 |
| Bayes' Rule | §1.3 "Conditional Probability and Independence" (Bayes' Rule) | NUM-PROB-008, -009 |
| Discrete distributions: Bernoulli; geometric in the trials form p(1−p)^(x−1); negative binomial in the trials form C(x−1, r−1)pʳ(1−p)^(x−r); hypergeometric; discrete uniform on 1…N (shifted to a…b) | §3.2 "Discrete Distributions" | NUM-PROB-023…025, 029…040, 057, 058 |
| P(g(X) ≥ r) ≤ E g(X)/r for g ≥ 0 (Markov, Chebyshev) | §3.6 "Inequalities and Identities" | NUM-PROB-020, -021, -022 |

How checked: worked values hand-checked in the tests; exact big.Rat oracles for the discrete distributions.

Status: NUM-PROB-002…009 are cited at section level since the fix pass of 2026-11-12 (fix_prob.md §2): §1.2 and §1.3 titles confirmed from the table of contents. The result numbers formerly cited (Theorems 1.2.8(c), 1.2.9(b), 1.2.11(a), 1.3.5; Definitions 1.3.2, 1.3.7) were never confirmed first-hand and are no longer cited. §3.2 and §3.6 titles and pages (85, 121) from library TOCs (typeset.io, ProofWiki book page, HTW Berlin TOC PDF); §3.2 contents via Purdue STAT 610 lecture 3. §2.3 title via the table of contents only. Theorem 3.6.1's number and Definition 2.3.2 NOT confirmed and not cited.

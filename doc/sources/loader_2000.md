# Loader 2000

**Citation**: Loader, "Fast and Accurate Computation of Binomial Probabilities", technical report, 2000 (revised February 5, 2002), https://www.r-project.org/doc/reports/CLoader-dbinom-2002.pdf (read 2026-11-12).

Taken: §2 eq. (3), the saddle-point identity; eq. (4), Stirling–De Moivre with remainder δ(n); eq. (5), p(x; n, x/n); eq. (6), the deviance D₀ and its series (x−np)²/(x+np) + 2x Σ v^(2j+1)/(2j+1), used when |x − np| < 0.1(x + np); eq. (7), Poisson; Appendix A: the δ table for n ≤ 15, the series terms 1/1680 and 1/1188, the thresholds 35, 80, 500, and the assembled dbinom/dpois (stirlerr, bd0).

Used by NUM-PROB-026, -027, -032, -033, -035, -041, -042 (and the terms inside -036).

Checked: report read in full; the δ table recomputed independently with 60-digit decimal arithmetic. Tests: exact big.Rat binomial, negative binomial and hypergeometric oracles, plus the n = 1000 and λ = 1000 values (λ = 1000 test value 0.0126146113487215 from 50-digit decimal; the catalog's 0.012614611348719664 is off by 1.5e-13 relative). Status: confirmed. The saddle-point exponent sign is +δ(n) − δ(k) − δ(n−k) (eq. 5); the catalog's sign is wrong, the code is right. Relative error near machine precision holds near the mode only; in the tails it is about 2⁻⁵³·D₀ (stated in the doc comment).

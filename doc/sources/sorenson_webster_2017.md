# Sorenson and Webster 2017

**Citation**: Sorenson, Webster, "Strong pseudoprimes to twelve prime bases", *Mathematics of Computation*, 2017 (online 2016-06-02), 86(304):985–1003, doi:10.1090/mcom/3134.

Taken: ψ₁₂ = 318665857834031151167461 > 2⁶⁴, so strong tests to the 12 primes ≤ 37 decide primality exactly on uint64. Used by NUM-INTEG-001 IsPrime. Checked: ψ₁₂ re-verified by script to be a strong pseudoprime to all 12 bases and > 2⁶⁴; `TestIsPrime_strongPseudoprimes`, `TestIsPrime_oracle` (vs math/big ProbablyPrime(0)). Status: confirmed on Crossref 2026-11-12.

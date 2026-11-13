# C2 (pending clustering.go:37, :121) → NUM-STATS-033 ChiSquareGoodnessOfFit, NUM-STATS-034 ChiSquareGoodnessOfFitUniform (PORTED)
- Citation: NIST/SEMATECH e-Handbook of Statistical Methods, §1.3.5.15 "Chi-Square Goodness-of-Fit Test", https://www.itl.nist.gov/div898/handbook/eda/section3/eda35f.htm.
- Transcribed: "χ² = Σᵢ₌₁ᵏ (Oᵢ − Eᵢ)²/Eᵢ", Oᵢ "the number of observations that fall in bin i", Eᵢ "expected to fall in bin i", "Eᵢ = N(F(yᵢ₊₁) − F(yᵢ))"; "(k − c) degrees of freedom where k is the number of non-empty cells and c = the number of estimated parameters"; "rejected if χ² > χ²₁₋α,k₋c".
- Uniform: F discrete uniform over k bins → Eᵢ = N/k, df = k − 1. Day-part buckets and the 7.815 / 12.592 flags (app policy) not ported.
- Checks: hand example 10, 20, 30, 40 → 20; source-shaped 12, 3, 3, 2 → 13.2 and 1×6, 8 → 21.
- Route: WebFetch, read 2026-11-12.

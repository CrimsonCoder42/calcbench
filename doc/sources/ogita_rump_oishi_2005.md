# Ogita

**Citation**: Ogita, Rump, Oishi, "Accurate Sum and Dot Product", SIAM Journal on Scientific Computing, 2005, 26(6):1955–1988, doi:10.1137/030601818.

**Taken**: Algorithm 3.1 TwoSum, Algorithm 3.5 TwoProductFMA, Algorithm 4.4 Sum2, Algorithm 5.3 Dot2, Proposition 5.5 (|res − xᵀy| ≤ eps|xᵀy| + γ²ₙ|xᵀ||y|).

**Used by**: NUM-LINALG-040, NUM-LINALG-041, NUM-LINALG-042, NUM-LINALG-044, NUM-LINALG-093.

**Checked**: PDF (https://www.tuhh.de/ti3/paper/rump/OgRuOi05.pdf) text extracted and read 2026-11-12; DOI on Crossref. Test: Proposition 5.5's bound against math/big on 300 ill-conditioned dot products (TestVecDotCompensated_oracle). Build note: Go may fuse a·b into a following addition; twoProduct and twoSum use explicit float64 conversions to keep the transformations error-free.

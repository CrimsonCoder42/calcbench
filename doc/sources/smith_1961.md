# Smith

**Citation**: Smith, "Eigenvalues of a symmetric 3 × 3 matrix", Communications of the ACM, 1961, 4(4):168, doi:10.1145/355578.366316.

**Taken**: The trigonometric closed form for the eigenvalues of a symmetric 3×3 matrix (m = tr/3, q = det(A − mI)/2, p = Σ(A − mI)ᵢⱼ²/6, φ = ⅓ atan(√(p³ − q²)/q)).

**Used by**: NUM-LINALG-076.

**Checked**: DOI and pages confirmed on Crossref 2026-11-12. The paper's text was not readable online; the formula was checked by the oracle test R·D·Rᵀ → D on 500 random rotations (TestMatSymmetricEigenvalues_oracle), which found the √ε accuracy limit for repeated eigenvalues now documented.

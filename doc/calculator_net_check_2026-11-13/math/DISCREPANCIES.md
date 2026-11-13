# math: discrepancies (A-MATH)

| slug | inputs | ours | theirs | verdict | citation |
|---|---|---|---|---|---|
| percent-error-calculator | observed 10, true 11 | 9.0909 % | −9.0909 % | THEIRS DIFFERS (they keep the sign of observed − true) | Burden, Faires & Burden, Numerical Analysis, 10th ed., 2016, §1.2, Def. 1.15 (absolute value) |
| percent-error-calculator | observed 5, true 0 | ErrDegenerate | "Infiniti" | THEIRS DIFFERS (undefined for true = 0) | same, Def. 1.15 requires p ≠ 0 |
| half-life-calculator | T½ 5730 → τ, λ | τ 8266.64258429376, λ 1.2096809433856e-4 | τ 8266.6425838161, λ 1.2096809434555e-4 | THEIRS DIFFERS (they use ln 2 ≈ 0.6931471806 in λ and τ) | Abramson, College Algebra 2e, OpenStax 2021, §6.7 (t½ = −ln 2/k, exact ln 2); Ling et al., University Physics Vol. 3, §10.3 eqs. 10.15–10.16 |
| quadratic-formula-calculator | a 1, b −1e8, c 1 | small root 1e−8 | 7.4505805969238e−9 | THEIRS DIFFERS (cancellation in −b + √(b² − 4ac); roots must multiply to c/a = 1) | Press et al., Numerical Recipes in C, 2nd ed., 1992, §5.6, eqs. 5.6.4–5.6.5; Higham, Accuracy and Stability, 2nd ed., 2002, §1.8 |
| sample-size-calculator | 99 %, d 2 %, p 30 %, N ∞ / 90 %, d 5 %, p 50 %, N 10000 | 3484 / 264 | 3495 / 266 | THEIRS DIFFERS (z rounded to 2.58 and 1.65) | Cochran, Sampling Techniques, 3rd ed., 1977, §4.4; Wichura AS 241, Appl. Stat. 1988, 37(3):477–484 (z₀.₉₉₅ = 2.5758293, z₀.₉₅ = 1.6448536) |
| statistics-calculator, mean-median-mode-range-calculator | 1.5, −2.25, 7, 7, 3.125, −2.25, 0 | GeometricMean ErrNegative | 0 | THEIRS DIFFERS (domain: ours requires non-negative values; product here is 0 because of the 0) | NIST Dataplot Reference Manual, "GEOMETRIC MEAN", (∏Xᵢ)^(1/n) |
| pythagorean-theorem-calculator | b 100000, c 100000.001 | a 14.142135686251251 (exact for the stored inputs) | 14.142135691166 | THEIRS DIFFERS (c² − b² cancellation) | Goldberg, ACM Computing Surveys, 1991, 23(1):5–48, doi:10.1145/103162.103163, §1.4 (x² − y² = (x − y)(x + y)) |
| surface-area-calculator | ellipsoid semi-axes 1, 2, 3.5 | 56.14575157027773 (exact) | ≈ 56.193275605353 | THEIRS DIFFERS (p = 1.6 power-mean approximation, ≈ 1 % class) | NIST DLMF eq. 19.33.1 (S = 3V R_G(a⁻², b⁻², c⁻²)); Carlson, Numer. Algorithms 1995, 10(1):13–26 |

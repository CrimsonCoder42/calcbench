# Discrepancies: ours vs calculator.net

Merged by the lead from each folder's DISCREPANCIES.md (2026-11-13). Verdicts: OURS WRONG (fixed), THEIRS DIFFERS (ours stays, source cited), UNRESOLVED (for Devin). The primary source decides; no constant was changed to match their number.

## financial


| slug | inputs | ours | theirs | verdict | citation |
|---|---|---|---|---|---|
| mortgage-calculator | 400k, 20 % down, 30 y, 7.263 % | total of payments 786,879.68 | 786,883.11 | THEIRS DIFFERS (they multiply the unrounded payment by n; ours sums the rounded schedule, last payment adjusted) | 12 CFR 1026.18(h) total of payments = sum of scheduled payments; D-022 |
| amortization-calculator | 200k, 15 y, 6 %, +200/month | interest saved 18,230.54; total interest 103,788.82 | 18,230.42; 103,788.46 | THEIRS DIFFERS (unrounded running cents vs per-period half-even rounding) | D-022; 12 CFR 1026.18(h) |
| refinance-calculator | balance 250,000, payment 1,800, 7 % | 285.69 → 286 payments (balance after 285 = 1,242.95) | 285 months, total 513,000 | THEIRS DIFFERS (they truncate the payment count and drop the final partial payment) | NCEES FE Reference Handbook 10.0.1 p. 230 (P/A,i,n); balance recursion FIN-LOAN-003 |
| apr-calculator | 25,750 note (750 financed fee), 300 upfront, 5 y, 8.5 % | APR 10.28726 % | 10.288 % | THEIRS DIFFERS (0.0007 pp; independent bisection agrees with ours) | 12 CFR 1026 Appendix J; 1026.22(a)(2) tolerance |
| bond-calculator | price 950, face 1,000, 4 % semiannual coupon, 10 y | YTM 4.63032 % | 4.6302 % | THEIRS DIFFERS (0.0001 pp; independent bisection 4.630325 %; their yield prices the bond at 950.0096) | QIA 3rd ed. ch. 1 §6.1 Eq. 11, §7.1; OpenStax Principles of Finance §10.2 |
| payback-period-calculator | −100,000; 30,000 growing 5 %/yr, 5 y; 10 % | payback 3.1562; discounted 3.9209 | 3.144; 3.943 | THEIRS DIFFERS (their growing-flow form uses an unidentified within-year rule; their listed-flow form matches ours exactly) | Dahlquist & Knight, Principles of Finance, OpenStax 2022, §16.1 and §16.4 |
| business-loan-calculator | 10,000, 10 %, 60 mo, fees 500 + 750 | APR 15.9326 % | 15.931 % | THEIRS DIFFERS (0.0016 pp; independent bisection 15.93265 %) | 12 CFR 1026 Appendix J; 1026.22(a)(2) |
| personal-loan-calculator | 8,000, 13.5 %, 36 mo, 300 fee deducted | APR 16.2026 % | 16.205 % | THEIRS DIFFERS (0.0024 pp; independent bisection 16.2026–16.2032 %) | 12 CFR 1026 Appendix J; 1026.22(a)(2) |

## fitness-and-health


| slug | inputs | ours | theirs | verdict | citation / note |
|---|---|---|---|---|---|
| body-fat-calculator | m waist 34, neck 15, h 70 in | 17.5131 | 17.4 | THEIRS DIFFERS (code withdrawn, NEEDS SOURCE) | ours: DoD direct %fat eq., Foulis et al. Front Physiol 2023;14:1183836 Methods. Theirs reproduces the Hodgdon–Beckett 1984 body-density form + Siri (495/D − 450). |
| body-fat-calculator | f waist 30, hip 38, neck 13, h 65 in | 28.5563 | 28.3 | THEIRS DIFFERS (code withdrawn, NEEDS SOURCE) | same as above |
| bmi-calculator | 100 lb, 72 in | "Underweight" | "Severe thinness" | THEIRS DIFFERS (labels) | ours: CDC adult categories (MED-ANTHRO-003); theirs: WHO-style labels, same cut points |
| ideal-weight-calculator | man 70 in / woman 64 in, Hamwi | 166.0 / 120.0 lb | 165.3 / 119.7 lb | THEIRS DIFFERS (code withdrawn, NEEDS SOURCE) | ours: Hamwi lb form, Peterson et al. AJCN 2016;103:1197 Table 3. Theirs: metric restatement 48 + 2.7 kg/in, 45.5 + 2.2 kg/in |
| macro-calculator | m 25 y, 1740 kcal | AMDR protein 43.5–152.25, carb 195.75–282.75, fat 38.67–67.67 g | 62–149, 186–303, 40–69 g | THEIRS DIFFERS | ours: NASEM 2024 Letter Report Table 3-1 (IOM 2002/2005 AMDR), LIF-NUTR-003. Theirs: unstated scheme, not the AMDR % |
| carbohydrate-calculator | 2060.4 kcal, 65 % | 334.81 g | 357 g | THEIRS DIFFERS | ours: 4 kcal/g, 21 CFR 101.9(c)(1)(i)(B); theirs: 3.75 kcal/g (monosaccharide-equivalent factor); their 40–75 % band is not the AMDR 45–65 % |
| protein-calculator | m 25 y 160 lb 5ft10 ×1.375 | 59.02–206.58 g (10–35 % of 2360.9 kcal) | 58–202 g | THEIRS DIFFERS | same % rule (NASEM AMDR, LIF-NUTR-003); their implied energy ≈ 2310 kcal, basis not shown |
| fat-intake-calculator | 2060.4 kcal, 20–35 % | 45.79–80.13 g | 47–82 g | THEIRS DIFFERS | ours: 9 kcal/g, 21 CFR 101.9(c)(1)(i)(B); theirs reproduced by ≈ 8.8 kcal/g |
| gfr-calculator | m 50 y SCr 0.9 | 104.049 (CKD-EPI 2021) | 99.2 | THEIRS DIFFERS | ours: Inker NEJM 2021 (MED-RENAL-001); theirs: CKD-EPI 2009 with race term (Levey Ann Intern Med 2009) |
| due-date-calculator | LMP 2026-07-04, cycle 32 | EDD 2027-04-10, GA 18w6d | Apr 14, 2027, 18w2d | THEIRS DIFFERS | ours: ACOG CO 700 LMP + 280 d, no cycle adjustment stated; theirs adds (cycle − 28) days |

## math


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

## other


| slug | inputs | ours | theirs | verdict | citation |
|---|---|---|---|---|---|
| ip-subnet-calculator | 10.0.0.1/31 | 2 usable hosts (10.0.0.0, 10.0.0.1) | 0 usable, range NA | THEIRS DIFFERS | RFC 3021 §2 (31-bit prefixes on point-to-point links: both addresses usable); theirs applies the RFC 950 total − 2 rule |

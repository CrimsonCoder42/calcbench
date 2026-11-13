# pythagorean-theorem-calculator (A-MATH, 2026-11-13)

CLASSIFY: in scope.
FIELDS (1 load): calform av (2), asv, bv, bsv, cv (2), csv (5) — each side as v·√sv (default c = 2√5); any two sides give the third.
MATCH: HAVE NUM-GEOM-020 RightTriangleHypotenuse, NUM-GEOM-021 RightTriangleLeg.
COMPARE (scratch/pythagorean-theorem-calculator.go; 3 submissions)

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| a 3, b 4 | c 5 | 5 | 0 | match |
| b 5, c 13 | a 12 | 12 | 0 | match |
| b 100000, c 100000.001 | a 14.142135686251251 | 14.142135691166 (√200.00000190735) | 4.9e−9 (rel 3.5e−10) | THEIRS DIFFERS |

DISCREPANCY: THEIRS DIFFERS. c² − b² cancels catastrophically when c ≈ b: their radicand 200.00000190735 carries the rounding error of c² ≈ 10¹⁰ (ulp 1.9e−6). Ours forms (c − b)(c + b), whose factors are exact here; the exact leg of the stored inputs is 14.14213568625125063… (50-digit decimal check), equal to ours. (The decimal inputs themselves give 14.1421356590863; both sides see the binary value of 100000.001, so that difference is input representation, not the formula.) Ours stays. Citation for the cancellation-free form: Goldberg, "What Every Computer Scientist Should Know About Floating-Point Arithmetic", ACM Computing Surveys, 1991, 23(1):5–48, doi:10.1145/103162.103163, §1.4 "Cancellation" (x² − y² computed as (x − y)(x + y)).

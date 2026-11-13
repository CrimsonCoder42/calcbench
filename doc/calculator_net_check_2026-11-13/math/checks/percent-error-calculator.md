# percent-error-calculator (A-MATH, 2026-11-13)

CLASSIFY: in scope.
FIELDS (1 load): calform observedvalue (default 10), truevalue (default 11).
MATCH: NONE (no relative/percent error in INDEX.md).
RESEARCH: Burden, Faires & Burden, Numerical Analysis, 10th ed., Cengage 2016, §1.2, Definition 1.15: relative error |p − p*|/|p|, p ≠ 0; Example 1(a) p = 0.3000×10¹, p* = 0.3100×10¹ → 0.3333×10⁻¹.
BUILD: NUM-ALGEB-035 PercentError, num/algebra/percent_error.go; tests in num/algebra/percent_test.go.

COMPARE (scratch/percent-error-calculator.go; 4 submissions)

| inputs (observed, true) | ours | theirs | diff | verdict |
|---|---|---|---|---|
| 3.1, 3 (Burden Ex. 1a) | 3.333333333333336 % | 3.3333 % | 0 at 4 dp | match |
| 10, 11 (their default) | 9.090909090909092 % | −9.0909 % | sign | THEIRS DIFFERS |
| −11, −10 | 10 % | 10 % | 0 | match |
| 5, 0 | ErrDegenerate | "Infiniti" | n/a | THEIRS DIFFERS |

DISCREPANCY: THEIRS DIFFERS. They report the signed (observed − true)/|true| (their −11/−10 case gives +10 %, so the denominator is |true|). Burden Def. 1.15 defines relative error with the absolute value of the difference, so ours stays. A true value of 0 has no relative error by the definition (p ≠ 0); ours returns ErrDegenerate, theirs prints an infinity.

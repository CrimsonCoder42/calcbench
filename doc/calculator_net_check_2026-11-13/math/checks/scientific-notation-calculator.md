# scientific-notation-calculator (A-MATH, 2026-11-13)

FIELDS (1 load): calform cvtnum (1568938, convert to a × 10ⁿ); calform2 x (1.23), xpow (7), y (3.45), ypow (2), cp (20, significant figures shown) for +, −, ×, ÷ of two numbers in scientific notation.
CLASSIFY: SKIPPED. Normalised a × 10ⁿ (1 ≤ |a| < 10) to any number of significant figures is strconv.FormatFloat(x, 'e', digits−1, 64) in the standard library, correctly rounded; the arithmetic is ordinary float64 arithmetic. No formula to build, no comparison run.

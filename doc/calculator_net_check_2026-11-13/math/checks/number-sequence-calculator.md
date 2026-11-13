# number-sequence-calculator (A-MATH, 2026-11-13)

CLASSIFY: in scope.
FIELDS (1 load): calform afirstnumber (2), afactor (5, common difference), athenumber (20, n); calform2 gfirstnumber (2), gfactor (5, ratio), gthenumber (12); calform3 fthenumber (10, Fibonacci n). Outputs: first terms, n-th term, sum through n; geometric also the infinite sum when |r| < 1.
MATCH: NONE for arithmetic/geometric; HAVE ALG-DP-010 Fibonacci (n ≤ 93).
RESEARCH: Abramson, College Algebra 2e, OpenStax 2021, §9.2 (aₙ = a₁ + d(n − 1)), §9.3 (aₙ = a₁rⁿ⁻¹), §9.4 (Sₙ = n(a₁ + aₙ)/2; Sₙ = a₁(1 − rⁿ)/(1 − r)).
BUILD: NUM-SERIES-020 ArithmeticSequenceTerm, -021 ArithmeticSeriesSum (num/series/arithmetic_sequence.go); -022 GeometricSequenceTerm, -023 GeometricSeriesSum (num/series/geometric_sequence.go; expm1 form for r near 1); tests num/series/sequence_test.go. Infinite geometric sum a₁/(1 − r) not built (one division; noted in LEAD.md).

COMPARE (scratch/number-sequence-calculator.go; 7 submissions, 2 re-runs to read the sum line)

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| arithmetic a₁ 2, d 5, n 20 | a₂₀ 97 (S₂₀ 990) | 97 | 0 | match |
| arithmetic a₁ −1.5, d 0.25, n 7 | a₇ 0, S₇ −5.25 | 0, −5.25 | 0 | match |
| geometric a₁ 2, r 5, n 12 | a₁₂ 97656250 (S 122070312) | 97656250 | 0 | match |
| geometric a₁ 3, r 0.5, n 10 | a₁₀ 0.005859375, S₁₀ 5.994140625 | 0.005859375, 5.994140625 | 0 | match |
| Fibonacci n 50 | 12586269025 | 12586269025 | 0 | match |

# gcf-calculator (A-MATH, 2026-11-13)

CLASSIFY: in scope.
FIELDS (1 load): calform numberinputs (textarea, comma list; default "330, 75, 450, 225").
MATCH: HAVE NUM-INTEG-004 GCD (pairwise; a list is a left fold, scratch/lcm-calculator.go computes both).
COMPARE (scratch/lcm-calculator.go plus the cases below by hand; 4 submissions)

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| 40902, 24140 (Knuth §4.5.2 example) | 34 | 34 | 0 | match |
| 33, 99, 3 | 3 (hand: gcd(33, 99) = 33, gcd(33, 3) = 3) | 3 | 0 | match |
| 2⁴⁰, 3·2³⁰ | 1073741824 | 1073741824 | 0 | match |
| 17, 19 | 1 | 1 | 0 | match |

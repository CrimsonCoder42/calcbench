# quadratic-formula-calculator (A-MATH, 2026-11-13)

CLASSIFY: in scope.
FIELDS (1 load): calform av (1), bv (1), cv (1/4); fractions accepted as text.
MATCH: HAVE NUM-ALGEB-007 QuadraticRealRoots, NUM-ALGEB-008 QuadraticRoots (NUM-ALGEB-006 discriminant). Fraction text input is a UI concern (callers convert; NUM-INTEG-007..011 for exact fractions).
COMPARE (scratch/quadratic-formula-calculator.go; 5 submissions)

| inputs a, b, c | ours | theirs | diff | verdict |
|---|---|---|---|---|
| 1, 1, 1/4 | −0.5 (double, n = 1) | −0.5 | 0 | match |
| 2, −7, 3 | 0.5, 3 | 0.5, 3 | 0 | match |
| 1, 2, 5 | −1 ± 2i | −1 ± 2i | 0 | match |
| 1, 0, −2 | ±1.4142135623730951 | ±1.4142135623731 | 0 at 13 dp | match |
| 1, −1e8, 1 | 1e−8, 1e8 | 7.4505805969238e−9, 1e8 | small root 25 % off | THEIRS DIFFERS |

DISCREPANCY: THEIRS DIFFERS. For b² ≫ 4ac the textbook form (−b + √(b² − 4ac))/2a loses the small root to cancellation; the roots must multiply to c/a = 1, so the small root is 1e−8 (ours, Vieta-consistent). Press et al., Numerical Recipes in C, 2nd ed., 1992, §5.6, eqs. (5.6.4)–(5.6.5) and Higham, Accuracy and Stability of Numerical Algorithms, 2nd ed., 2002, §1.8 give the cancellation-free form ours uses. Ours stays.

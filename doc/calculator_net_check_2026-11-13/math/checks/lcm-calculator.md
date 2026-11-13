# lcm-calculator (A-MATH, 2026-11-13)

CLASSIFY: in scope.
FIELDS (1 load): calform numberinputs (textarea, comma list; default "330, 75, 450, 225").
MATCH: HAVE NUM-INTEG-005 LCM (pairwise, overflow-checked). A list is a left fold LCM(LCM(a, b), c)… (scratch does this); no list function built (a fold, no new maths).
COMPARE (scratch/lcm-calculator.go; 3 submissions)

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| 330, 75, 450, 225 | 4950 | 4950 | 0 | match |
| 17, 19, 23, 29, 31, 37, 41, 43 | 435656388001 | 435656388001 | 0 | match |
| 2⁴⁰, 3·2³⁰ | 3298534883328 | 3298534883328 | 0 | match |

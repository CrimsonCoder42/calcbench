# ratio-calculator (A-MATH, 2026-11-13)

CLASSIFY: in scope.
FIELDS (1 load): calform t1 (3) : b1 (4) = t2 (600) : b2 (empty) — any one blank is solved; calform2 c2t1 (250) : c2b1 (280), c2type {Enlarge, Shrink}, c2time (2.5) — scale a ratio by a factor (both terms ×/÷ the factor, no new maths).
MATCH: NONE.
RESEARCH: Marecek, Anthony-Smith & Honeycutt Mathis, Prealgebra 2e, OpenStax 2020, §6.5 "Solve Proportions and their Applications" (cross products a·d = b·c).
BUILD: NUM-ALGEB-036 ProportionFourthTerm, num/algebra/proportion.go; tests num/algebra/proportion_test.go. Other blank positions: reorder the proportion (doc comment). Scaling (calform2) not built: two multiplications.

COMPARE (scratch/ratio-calculator.go; 3 submissions)

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| 3 : 4 = 600 : ? | 800 | 800 | 0 | match |
| 7 : 3 = 10 : ? | 4.285714285714286 | 4.2857142857143 | 0 at 13 dp | match |
| −2 : 5 = 4 : ? | −10 | −10 | 0 | match |

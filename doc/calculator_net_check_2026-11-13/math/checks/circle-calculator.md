# circle-calculator (A-MATH, 2026-11-13)

CLASSIFY: in scope.
FIELDS (1 load): calform radius, diameter, circumference, area (default area 100) — any one gives the rest.
MATCH: HAVE NUM-GEOM-034 CircleCircumference, NUM-GEOM-035 CircleArea (from the radius). The inverse forms (radius from circumference C/2π, from area √(A/π)) are one stdlib operation each; not built.
COMPARE (scratch/circle-calculator.go; 5 submissions, 2 re-runs to read the lines the first filter cut off)

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| r 5 | C 31.4159265358979, A 78.5398163397448 | 31.415926535898, 78.539816339745 | 0 at 12 dp | match |
| r 123456.789 | C 775701.88271637, A 47882831830.7088 | 775701.88271637, 47882831830.709 | 0 at displayed figs | match |
| r 0.001 | d 0.002 (C, A not read) | d 0.002 | 0 | match (diameter only) |

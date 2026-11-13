# triangle-calculator (A-MATH, 2026-11-13)

CLASSIFY: in scope.
FIELDS (1 load): calform vx, vy, vz (sides; from the results vy = a, vx = b, vz = c), va, vb, vc (angles A, B, C; default vc 60, vx 1, vy 1 = SAS), angleunits {d, r}.
MATCH: HAVE NUM-TRIG-022 SolveTriangleSSS, -023 SAS, -024 ASA, -025 AAS, -026 SSA, NUM-GEOM-016 TriangleAreaHeron.
COMPARE (scratch/triangle-calculator.go; 3 submissions)

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| SSS a 4, b 3, c 6 (their labels) | angles 36.3360575146°, 26.3843297494°, 117.279612736° | 36.336°, 26.384°, 117.28° | 0 at displayed figs | match |
| ASA A 40°, B 75°, c 10 | a 7.09237655111, b 10.6578122808, C 65° | 7.09238, 10.65781, 65° | 0 at 5 dp | match |
| SSS 1e6, 1e6, 0.001 (needle) | A = B = 89.9999999714°, C 5.72957795131e−8° | 90°, 90°, 0°; labelled "Right Isosceles" | 0 at displayed figs | match (numbers); their "Right" label is a display rounding artefact (A + B + C = 180 needs C > 0), not a value we compute |

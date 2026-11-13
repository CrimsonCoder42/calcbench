# Stair Calculator (row 13)

CLASSIFY: in scope (right-triangle geometry). The building-code limits (IRC R311.7.5: riser ≤ 7¾ in,
tread ≥ 10 in) are a model code adopted by jurisdiction; not built this session.
FIELDS: calform (fixed run): run 10 inch, totalrun 15 foot, totalheight 10 foot (unit selects foot/inch/meter/
centimeter). calform2 (fixed rise/step count): run2 10 in, totalrun2 15 ft, totalheight2 10 ft, risesetting r|s,
rise2 7 in, step2 15, treadsetting n|y, treadthickness2 1 in, headroomsetting n|y, flooropen 20 ft,
floorthickness 6 in, headroom 8 ft, mountsetting s|f.
MATCH: HAVE NUM-GEOM-020 RightTriangleHypotenuse (stringer), NUM-TRIG-027 SolveRightTriangleLegs (angle).
Their stringer convention (read from their output): legs = treads × run and total height − one rise.
COMPARE (ours: checks/scratch/stair-calculator.go; total height 120 in, run 10 in):

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| 12 treads (13 risers) | rise 9.23 in, stringer 13 ft 7.31 in, height 110.77 in, angle 42.71° | 9.23 in, 13 ft 7.31 in, 9 ft 2.77 in, 42.71° | 0 | MATCH |
| 15 treads (16 risers) | rise 7.50 in, stringer 15 ft 7.50 in, angle 36.87° | 7.5 in, 15 ft 7.5 in, 36.87° | 0 | MATCH |

Submissions: 1 (the one page lists every step count).
CLOSE: MATCH (geometry; IRC limits not built, LEAD.md).

# area-calculator (A-MATH, 2026-11-13)

CLASSIFY: in scope (their unit selects are unit conversion, not compared).
FIELDS (1 load): calform1 reclength, recheight; calform2 triedge1..3; calform3 tratop, trabottom, traheight; calform4 circleradius; calform5 sectorradius, sectorangle (+unit); calform6 majoraxes, minoraxes (semi-axes, from their πab); calform7 parbase, parheight; each with a "<field>unit" select.
MATCH: HAVE NUM-GEOM-022, -025, -035, -043 and RectangleArea, TriangleAreaHeron (NUM-GEOM-016), TrapezoidArea, CircleArea, SectorArea, EllipseArea, ParallelogramAreaSAS (num/geometry).
COMPARE (scratch/area-calculator.go; 4 submissions)

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| trapezoid 5, 9.5, h 4 | 29 | 29 | 0 | match |
| sector r 6, 75° | 23.56194490192345 | 23.561944901923 | 0 at 12 dp | match |
| ellipse semi-axes 7, 2.5 | 54.97787143782138 | 54.977871437821 | 0 at 12 dp | match |
| triangle 3, 4, 6 | 5.332682251925386 | not read (first lines show s = 6.5 only) | n/a | — |

# slope-calculator (A-MATH, 2026-11-13)

CLASSIFY: in scope.
FIELDS (1 load): calform x11 (1), y11 (1), x12 (2), y12 (2) — two points; calform2 x21 (1), y21 (1), d2 (5, distance), m2 (0.75, slope), th2 (angle) — second point from one point, distance and slope or angle.
MATCH: HAVE NUM-GEOM-003 Slope, NUM-GEOM-004 SlopeAngleRad, PointDistance (num/geometry/point_distance.go).
COMPARE (scratch/slope-calculator.go; 2 submissions)

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| (−3, 4), (5, −2) | m −0.75; angle −36.8698976458° | m −0.75; θ 323.13010235416° | 360° | match after convention: they report θ in [0°, 360°), ours in (−180°, 180°] (atan2); −36.8698976458 + 360 = 323.1301023542 |
| (2.5, 7), (2.5, 9) vertical | Slope → error "slope of a vertical line is undefined"; angle 90° | m "infinity"; θ 90° | n/a | convention: ours returns the documented error, theirs prints infinity; angle matches |

Outcome: MATCH (conventions written above).

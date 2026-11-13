# Roofing Calculator (row 14)

CLASSIFY: in scope (geometry); price is one multiplication; shingle bundles per square are product data.
FIELDS: calform: acarea 2000, acareaunit feet (feet/meters/yards/inches/centimeters), roofpitch select 1..24
(rise per 12, default 6), angle 25 (alternative input), eaves 1 foot, price (empty). calform2: mcarea 2000 feet.
MATCH: HAVE NUM-GEOM-020 RightTriangleHypotenuse: roof area = plan area × hypot(12, rise)/12.
The eaves-overhang handling needs the plan perimeter, which the form does not take; their convention not
inferred, not compared (eaves set to 0).
COMPARE (ours: checks/scratch/roofing-calculator.go):

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| 2000 ft², 6/12, eaves 0 | 2236.07 ft² | 2,236 ft² | 0 at 0 dp | MATCH |
| 2000 ft², 12/12, eaves 0 | 2828.43 ft² | 2,828 ft² | 0 at 0 dp | MATCH |

Submissions: 2.
CLOSE: MATCH (geometry).

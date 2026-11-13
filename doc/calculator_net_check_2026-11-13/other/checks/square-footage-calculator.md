# Square Footage Calculator (row 12)

CLASSIFY: in scope (area); the price-per-unit field is one multiplication.
FIELDS: calform1 rectangle: reclength 30, reclengthunit feet (feet/yards/inches/miles/meters/kilometers/
centimeters), recheight 20, recheightunit feet, recquantity 1, recprice (empty), recpriceunit feet.
calform12 rectangle border: recblength 30, recbheight 20, recborder 2 (feet). Further forms for other shapes
(circle, triangle, trapezoid, sector...) not listed here.
MATCH: HAVE NUM-GEOM-022 RectangleArea, NUM-GEOM-035 CircleArea. PARTIAL: no m² ↔ ft² conversion.
BUILD: unit/area (planned AREA): UNT-AREA-001 SquareFeetToSquareMeters, UNT-AREA-002 SquareMetersToSquareFeet,
NIST SP 811 (2008) B.8, 1 ft² = 0.09290304 m² exactly.
COMPARE (ours: checks/scratch/square-footage-calculator.go):

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| 12.5 m × 8 m, in ft² | 1076.391041671 | 1076.391041671 | 0 at 9 dp | MATCH |

Submissions: 1.
CLOSE: BUILT UNT-AREA-001, UNT-AREA-002 (and MATCH).

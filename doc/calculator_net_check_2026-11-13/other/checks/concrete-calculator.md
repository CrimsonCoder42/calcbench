# Concrete Calculator (row 10)

CLASSIFY: in scope (volume); the weight and bag counts rest on an assumed density (their 133 lb/ft³) and bag
yields, which are product data, not a primary source.
FIELDS: calcf1 slab: slablength 5 foot, slabwidth 2.5 foot, slabthick 5 inch, slabquantity 1 (unit selects foot/
inch/yard/meter/centimeter). calcf2 hole: holediameter 2.5 ft, holedepth 6 ft, holequantity 1. calcf3 tube:
tubediameter 5 ft, tubeindiameter 4 ft, tubedepth 6 in. calcf4 curb and gutter: curbdepth 4 in, curbgwidth 10 in,
curbheight 4 in, curbfthick 5 in, curblength 10 ft. calcf5 stairs: stairrun 12 in, stairrise 6 in, stairwidth
50 in, stairpdepth 5 in, stairquantity 5.
MATCH: HAVE NUM-GEOM-022 RectangleArea, NUM-GEOM-048 PrismVolume, NUM-GEOM-049 CylinderVolume (volume).
PARTIAL: no foot/yard length or ft³/yd³ volume conversions in unit/ yet (unit/length has cm↔in only; unit/volume
is planned); weight/bag counts not built (density: no primary source in hand; ACI 211.1 tabulates fresh
concrete unit weight by aggregate size, not checked this session: NEEDS SOURCE for that part).
COMPARE (ours: checks/scratch/concrete-calculator.go; 1 in = 1/12 ft, 1 yd³ = 27 ft³):

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| slab 5 ft × 2.5 ft × 5 in | 5.2083 ft³, 0.1929 yd³ | 5.21 ft³, 0.19 yd³ | 0 at 2 dp | MATCH |
| hole Ø 2.5 ft × 6 ft | 29.4524 ft³, 1.0908 yd³ | 29.45 ft³, 1.09 yd³ | 0 at 2 dp | MATCH |

Submissions: 2.
CLOSE: MATCH (volume; weight/bags NEEDS SOURCE).

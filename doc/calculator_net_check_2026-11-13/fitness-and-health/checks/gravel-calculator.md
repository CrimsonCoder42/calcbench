# gravel-calculator (HANDOFF from A-OTHER; notes ../other/checks/gravel-calculator.md)

## FIELDS (form calform, 2026-11-13)
areasetting a (area) | r (rectangle) | c (circle) · totalarea 200 · totallength 20 / totalwidth 15 · totaldiameter 10 (units f|m|i|c|y) ·
gdepth 3 + unit · gdensityselect g | gs | ds | ws | own · gdensity 1600 + unit kgm|gcm|lbf · price, bag size.
Results: volume ft³, weight lb, area ft² / yd² / m².

## MATCH
HAVE NUM-GEOM-035 CircleArea (num/geometry/circle.go:35) and NUM-GEOM-048 PrismVolume; mass = volume × caller density.
Their preset densities are material data (no single primary source, A-OTHER note) — caller input. Conversions in scratch only.

## COMPARE (scratch: checks/scratch/gravel-calculator.go)
| inputs | ours | theirs | verdict |
|---|---|---|---|
| circle Ø 10 ft, 3 in deep, 1600 kg/m³ | 78.540 ft², 19.6350 ft³, 1961.23 lb | 79 ft², 19.6 ft³, 1,961 lb | MATCH |
Submissions: 1 fields + 1 run.

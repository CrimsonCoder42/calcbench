# mulch-calculator (HANDOFF from A-OTHER; notes ../other/checks/mulch-calculator.md)

## FIELDS (form calform, 2026-11-13)
areasetting a (area) | d (up to 12 length × width rows l1..l12, w1..w12, u1..u12) · totalarea 200 + unit f|m|y|a · mdepth 2 + unit i|c|f|m|y ·
price + priceunit (per ft³, yd³, m³, L, bags). Result: volume (ft³, also yd³ / m³ / bags).

## MATCH
HAVE NUM-GEOM-048 PrismVolume (area × depth). Unit factors (12 in/ft, 27 ft³/yd³) are exact (NIST SP 811 B.8, per A-OTHER); unit/volume does
not exist yet, so the conversions stay in the scratch only. No life/garden package created.

## COMPARE (scratch: checks/scratch/mulch-calculator.go)
| inputs | ours | theirs | verdict |
|---|---|---|---|
| 200 ft², 2 in | 33.3333 ft³ (1.2346 yd³) | 33.3 ft³ | MATCH |
Submissions: 1 fields + 1 run.

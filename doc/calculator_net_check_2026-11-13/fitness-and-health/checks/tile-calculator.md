# tile-calculator (HANDOFF from A-OTHER; their notes: ../other/checks/tile-calculator.md)

## FIELDS (form calform, 2026-11-13)
tilelength 12 / tilewidth 12 (inch|foot|meter|cm|yard) · areasetting d (dimensions) | a (area) · totalarea + unit · totallength 20 / totalwidth 15
(foot...) · gapsize 0.25 (inch...) · boxsize 12 · price + priceunit. Results: tiles, boxes, area, a 5–10 % extra range.

## MATCH
HAVE NUM-GEOM-022 RectangleArea; tiles = ceil(floor area / ((tile length + gap) × (tile width + gap))) is caller arithmetic on it.
No new function: the gap model (one gap per tile per side) and the 5–10 % waste are trade conventions with no primary source (A-OTHER's
note: TCNA Handbook covers installation, not waste), so they stay caller inputs. No life/home package created.

## COMPARE (scratch: checks/scratch/tile-calculator.go)
| inputs | ours | theirs | verdict |
|---|---|---|---|
| 20 × 15 ft floor, 12 × 12 in tile, 0.25 in gap | 43200 / 150.0625 = 287.88 → 288 | 288 | MATCH |
| +5 % / +10 % | 302.4 → 303 / 316.8 → 317 | 303 – 317 | MATCH (convention) |
Submissions: 1 fields + 1 run.

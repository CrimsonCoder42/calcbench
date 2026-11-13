# protein-calculator

## FIELDS (form calform, 2026-11-13)
cage 25 · csex m|f · cheightfeet 5 / cheightinch 10 · cpound 160 · metric cheightmeter 180, ckg 60 · cactivity 1.2 / 1.375 (default) / ... / 1.9 ·
cformula m | k · cfatpct 20. Results: three lines — a "1.0–1.8 g/kg"-type range, 10–35 % of energy, a "safe lower limit" g/day.

## MATCH
HAVE LIF-NUTR-003 (10–35 % of energy, 4 kcal/g), LIF-NUTR-005 (RDA), LIF-NUTR-007 (ISSN 1.4–2.0 g/kg; not shown by them).
Not built: the g/kg ranges they attribute to named organisations (no primary read this session); the WHO/FAO/UNU 0.83 g/kg safe level
(WHO TRS 935, 2007) would be a candidate (not read, not built).

## COMPARE (scratch: checks/scratch/protein-calculator.go)
| inputs | ours | theirs | verdict |
|---|---|---|---|
| m 25 y, 5'10", 160 lb, ×1.375; protein 10–35 % of energy | energy 2360.9 → 59.02–206.58 g | 58–202 g | THEIRS DIFFERS |
| same, 0.83 g/kg (hand, not a function) | 60.2 g | 60 g | (reference only) |
Their 10–35 % bounds imply an energy of about 2,310 kcal (58 × 4 / 0.10 = 2320; 202 × 4 / 0.35 = 2309), not the 2,361 their own calorie
page gives for these inputs (Mifflin simplified × 1.375). The % → grams rule is the same as ours; the hidden energy basis differs. Ours stays.
Submissions: 1 fields + 1 run.

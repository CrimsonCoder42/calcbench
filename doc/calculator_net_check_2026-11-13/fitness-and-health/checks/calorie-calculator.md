# calorie-calculator

## FIELDS (form calform, 2026-11-13)
cage 25 (ages 15-80) · csex m|f (m) · cheightfeet 5 · cheightinch 10 · cpound 165 · metric: cheightmeter 180, ckg 65 ·
cactivity select: 1 (BMR), 1.2, 1.375, 1.465 (default), 1.55, 1.725, 1.9 · coutunit c|j · cformula m (Mifflin St Jeor, default) | h (Revised Harris-Benedict) | k (Katch-McArdle) · cfatpct 20.
Result: kcal/day (integer display) at the chosen activity factor, plus weight-change rows.

## MATCH
HAVE LIF-NUTR-002 RestingEnergyExpenditureMifflinStJeorSimplified (10W + 6.25H − 5A + 5 / −161), which is what their default uses.
LIF-NUTR-001 (9.99W, −4.92A, 166·sex − 161 form) differs from theirs by 1–4 kcal on the cases below — a second published form of the same
paper, not an error.
NOT BUILT (NEEDS SOURCE this session):
- activity multipliers 1.2/1.375/1.465/1.55/1.725/1.9: no primary source located; the FAO/WHO/UNU 2004 PAL bands (Human energy
  requirements, FNTRS 1, 2004, Table 5.1: 1.40-1.69 / 1.70-1.99 / 2.00-2.40) are a different scheme.
- Revised Harris-Benedict (Roza AM, Shizgal HM, Am J Clin Nutr 1984;40:168-182): paper not seen; coefficients not shipped from memory.
- Katch-McArdle (370 + 21.6 × LBM): primary (McArdle/Katch textbook) not seen.
Unit conversions in scratch: lb × 0.45359237, in × 2.54.

## COMPARE (scratch: checks/scratch/calorie-calculator.go)
| inputs | ours LIF-NUTR-002 | ours LIF-NUTR-001 | theirs | diff (002) | verdict |
|---|---|---|---|---|---|
| m 25 y, 5'10", 165 lb, BMR | 1739.677 | 1740.929 | 1,740 | 0 at integer | MATCH |
| f 40 y, 5'5", 140 lb, BMR | 1305.904 | 1308.469 | 1,306 | 0 | MATCH |
| m 60 y, 6'2", 220 lb, ×1.2 | 2253.184 | 2257.746 | 2,253 | 0 | MATCH |
Submissions: 1 fields + 3 runs.

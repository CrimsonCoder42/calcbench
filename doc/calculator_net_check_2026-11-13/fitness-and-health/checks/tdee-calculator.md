# tdee-calculator

## FIELDS (form calform, 2026-11-13)
cage 25 · csex m|f · cheightfeet 5 / cheightinch 10 · cpound 165 · metric cheightmeter 180, ckg 65 · cactivity 1 / 1.2 / 1.375 / 1.465 (default)
/ 1.55 / 1.725 / 1.9 · coutunit c|j · cformula m | k · cfatpct 20. Results: TDEE kcal/day, BMI, weight-change rows, macro tabs.

## MATCH
PARTIAL LIF-NUTR-002 (energy) and MED-ANTHRO-002 (BMI). The activity multiplier table stays NEEDS SOURCE (calorie-calculator.md, LEAD.md);
the factor is supplied by hand in scratch only.

## COMPARE (scratch: checks/scratch/tdee-calculator.go)
| inputs | ours | theirs | verdict |
|---|---|---|---|
| f 40 y, 5'5", 140 lb, ×1.465 | 1305.904 × 1.465 = 1913.150 kcal; BMI 23.295 | 1,913; BMI 23.3 | MATCH |
Submissions: 1 fields + 1 run.

# carbohydrate-calculator

## FIELDS (form calform, 2026-11-13)
cage 25 · csex m|f · cheightfeet 5 / cheightinch 10 · cpound 160 · metric cheightmeter 180, ckg 60 · cactivity 1.2 / 1.375 (default) /
1.465 / 1.55 / 1.725 / 1.9 · cformula m | k · cfatpct 20.
Results: energy per goal (maintain, −1 lb/wk, −2 lb/wk, ...) and carbohydrate g/day at 40 / 55 / 65 / 75 % of energy.

## MATCH
PARTIAL LIF-NUTR-003 AMDRGrams (carbohydrate 45–65 %, 4 kcal/g) + LIF-NUTR-002 energy. Activity factor supplied in scratch only (NEEDS SOURCE).

## COMPARE (scratch: checks/scratch/carbohydrate-calculator.go)
| inputs | ours | theirs | verdict |
|---|---|---|---|
| m 25 y, 5'10", 160 lb, ×1.2 → energy | 2060.397 kcal | 2,060 | MATCH |
| carbohydrate at 65 % of 2060.4 kcal | 334.81 g (AMDR high, 4 kcal/g) | 357 g | THEIRS DIFFERS |
| carbohydrate at 55 % | 283.30 g (4 kcal/g, hand) | 302 g | THEIRS DIFFERS |
Their grams are reproduced exactly by 3.75 kcal/g (2060 × 0.40 / 3.75 = 219.7 → 220; × 0.55 → 302.2; × 0.65 → 357.1; × 0.75 → 412.1).
3.75 kcal/g is the energy factor for carbohydrate expressed as monosaccharide equivalents; ours is the 4 kcal/g general factor of
21 CFR 101.9(c)(1)(i)(B) (LIF-NUTR-003 second Source). Their 40–75 % band is also not the AMDR 45–65 % (NASEM Table 3-1). Ours stays.
Submissions: 1 fields + 1 run.

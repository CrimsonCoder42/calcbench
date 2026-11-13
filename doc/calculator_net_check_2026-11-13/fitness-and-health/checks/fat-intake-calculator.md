# fat-intake-calculator

## FIELDS (form calform, 2026-11-13)
cage 25 · csex m|f · cheightfeet 5 / cheightinch 10 · cpound 160 · metric cheightmeter 180, ckg 60 · cactivity 1.2 / 1.375 (default) / ... / 1.9 ·
cformula m | k · cfatpct 20. Results per goal: energy, fat 20–35 % (g), saturated fat < 10 % and < 7 % (g).

## MATCH
PARTIAL LIF-NUTR-003 AMDRGrams (adult fat 20–35 %, 9 kcal/g). Saturated-fat caps not built (no primary read: DGA 2020-2025 10 % would be
the candidate).

## COMPARE (scratch: checks/scratch/fat-intake-calculator.go)
| inputs | ours | theirs | verdict |
|---|---|---|---|
| m 25 y, 5'10", 160 lb, ×1.2 → energy | 2060.397 kcal | 2,060 | MATCH |
| fat 20–35 % | 45.79–80.13 g (9 kcal/g) | 47–82 g | THEIRS DIFFERS |
Their grams are reproduced by about 8.8 kcal/g (2060 × 0.20 / 8.8 = 46.8 → 47; × 0.35 / 8.8 = 81.9 → 82); 9 kcal/g gives 45.8–80.1.
Ours: the 9 kcal/g general factor of 21 CFR 101.9(c)(1)(i)(B) (LIF-NUTR-003). Ours stays.
Submissions: 1 fields + 1 run.

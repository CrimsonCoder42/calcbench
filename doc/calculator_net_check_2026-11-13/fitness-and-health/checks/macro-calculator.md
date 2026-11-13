# macro-calculator

## FIELDS (form calform, 2026-11-13)
cage 25 (18-80) · csex m|f · cheightfeet 5 / cheightinch 10 · cpound 165 · metric cheightmeter 180, ckg 65 · cactivity 1 / 1.2 / 1.375 /
1.465 (default) / 1.55 / 1.725 / 1.9 · cgoal m (maintain) | l | l1 | l2 | g | g1 | g2 · cformula m (Mifflin) | k (Katch-McArdle) · cfatpct 20.
Results: kcal/day, then tabs Balanced / Low Fat / Low Carb / High Protein / Create Your Own with protein, carbs, fat g/day and ranges,
sugar and saturated-fat caps.

## MATCH
PARTIAL LIF-NUTR-003 AMDRGrams (+ LIF-NUTR-002 for the energy). Their "diet" tabs, goal offsets, sugar and saturated-fat caps: not ours;
activity factor / Katch-McArdle NEEDS SOURCE (calorie-calculator.md).

## COMPARE (scratch: checks/scratch/macro-calculator.go)
| inputs | ours | theirs (Balanced tab) | verdict |
|---|---|---|---|
| m 25 y, 5'10", 165 lb, BMR, maintain → energy | 1739.677 kcal (LIF-NUTR-002) | 1,740 | MATCH |
| protein range at 1740 kcal | 43.5–152.25 g (AMDR 10–35 %) | 62–149 g (point 106) | THEIRS DIFFERS |
| carbohydrate range | 195.75–282.75 g (45–65 %) | 186–303 g (point 232) | THEIRS DIFFERS |
| fat range | 38.67–67.67 g (20–35 %) | 40–69 g (point 49) | THEIRS DIFFERS |
Their ranges are not the adult AMDR percentages at the displayed energy (their bounds work out to about 14–34 % protein, 43–70 % carbohydrate,
21–36 % fat). Their scheme is not stated; ours is NASEM Table 3-1 (LIF-NUTR-003 Source). Ours stays.
Submissions: 1 fields + 1 run.
